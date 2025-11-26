package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dejavu/deployer/internal/k8s"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

type Worker struct {
	nats       *nats.Conn
	js         nats.JetStreamContext
	k8sClient  *k8s.Client
	db         *sql.DB
	namespace  string
	baseDomain string
}

type BuildCompleteEvent struct {
	DeploymentID string `json:"deployment_id"`
	ImageURL     string `json:"image_url"`
	Success      bool   `json:"success"`
	Logs         string `json:"logs"`
}

func New() (*Worker, error) {
	// Connect to NATS
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	// Connect to Kubernetes
	k8sClient, err := k8s.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %w", err)
	}

	// Connect to database
	db, err := connectDB()
	if err != nil {
		return nil, err
	}

	namespace := os.Getenv("K8S_NAMESPACE")
	if namespace == "" {
		namespace = "dejavu-apps"
	}

	baseDomain := os.Getenv("BASE_DOMAIN")
	if baseDomain == "" {
		baseDomain = "dejavu.local"
	}

	return &Worker{
		nats:       nc,
		js:         js,
		k8sClient:  k8sClient,
		db:         db,
		namespace:  namespace,
		baseDomain: baseDomain,
	}, nil
}

func (w *Worker) Start() error {
	_, err := w.js.Subscribe("BUILDS.complete", func(msg *nats.Msg) {
		var event BuildCompleteEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Error parsing event: %v", err)
			msg.Ack()
			return
		}

		log.Printf("🚀 Deploying: %s", event.DeploymentID)
		w.processDeploy(event)
		msg.Ack()
	})

	return err
}

func (w *Worker) processDeploy(event BuildCompleteEvent) {
	ctx := context.Background()

	// Update status to deploying
	w.updateDeploymentStatus(event.DeploymentID, "deploying")
	w.updateDeploymentLogs(event.DeploymentID, event.Logs)

	if !event.Success {
		w.updateDeploymentStatus(event.DeploymentID, "error")
		return
	}

	// Update image URL
	w.updateDeploymentImage(event.DeploymentID, event.ImageURL)

	// Get deployment info
	deployment, err := w.getDeployment(event.DeploymentID)
	if err != nil {
		log.Printf("Error getting deployment: %v", err)
		w.updateDeploymentStatus(event.DeploymentID, "error")
		return
	}

	// 1. Create namespace if not exists
	if err := w.k8sClient.EnsureNamespace(ctx, w.namespace); err != nil {
		log.Printf("Error creating namespace: %v", err)
		w.updateDeploymentStatus(event.DeploymentID, "error")
		return
	}

	// 2. Create deployment
	deploymentName := fmt.Sprintf("app-%s", event.DeploymentID[:8])
	if err := w.k8sClient.CreateDeployment(ctx, w.namespace, deploymentName, event.ImageURL, nil); err != nil {
		log.Printf("Error creating deployment: %v", err)
		w.updateDeploymentStatus(event.DeploymentID, "error")
		return
	}

	// 3. Create service
	if err := w.k8sClient.CreateService(ctx, w.namespace, deploymentName, 80); err != nil {
		log.Printf("Error creating service: %v", err)
		w.updateDeploymentStatus(event.DeploymentID, "error")
		return
	}

	// 4. Create ingress
	host := fmt.Sprintf("%s.%s", deployment.Subdomain, w.baseDomain)
	if err := w.k8sClient.CreateIngress(ctx, w.namespace, deploymentName, host, deploymentName, 80); err != nil {
		log.Printf("Error creating ingress: %v", err)
		w.updateDeploymentStatus(event.DeploymentID, "error")
		return
	}

	// 5. Create HPA
	if err := w.k8sClient.CreateHPA(ctx, w.namespace, deploymentName, 2, 10); err != nil {
		log.Printf("Error creating HPA: %v", err)
		// HPA is optional, continue anyway
	}

	// Update status to ready
	w.updateDeploymentStatus(event.DeploymentID, "ready")
	log.Printf("✅ Deployment %s is ready at %s", event.DeploymentID, host)
}

type Deployment struct {
	ID        string
	Subdomain string
}

func (w *Worker) getDeployment(id string) (*Deployment, error) {
	var d Deployment
	err := w.db.QueryRow(
		"SELECT id, subdomain FROM deployments WHERE id = $1",
		id,
	).Scan(&d.ID, &d.Subdomain)
	return &d, err
}

func (w *Worker) updateDeploymentStatus(id, status string) {
	_, err := w.db.Exec(
		"UPDATE deployments SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		status, id,
	)
	if err != nil {
		log.Printf("Error updating deployment status: %v", err)
	}
}

func (w *Worker) updateDeploymentImage(id, imageURL string) {
	_, err := w.db.Exec(
		"UPDATE deployments SET image_url = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		imageURL, id,
	)
	if err != nil {
		log.Printf("Error updating deployment image: %v", err)
	}
}

func (w *Worker) updateDeploymentLogs(id, logs string) {
	_, err := w.db.Exec(
		"UPDATE deployments SET build_logs = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		logs, id,
	)
	if err != nil {
		log.Printf("Error updating deployment logs: %v", err)
	}
}

func (w *Worker) Close() {
	if w.nats != nil {
		w.nats.Close()
	}
	if w.db != nil {
		w.db.Close()
	}
}

func connectDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
