package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dejavu/builder/internal/detector"
	"github.com/dejavu/builder/internal/runner"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type Worker struct {
	nats          *nats.Conn
	js            nats.JetStreamContext
	workspaceDir  string
	cacheDir      string
	registryURL   string
	registryUser  string
	registryPass  string
}

type DeploymentEvent struct {
	DeploymentID string `json:"deployment_id"`
	ProjectID    string `json:"project_id"`
	RepoURL      string `json:"repo_url"`
	BuildCommand string `json:"build_command"`
	OutputDir    string `json:"output_dir"`
	CommitHash   string `json:"commit_hash"`
}

type BuildCompleteEvent struct {
	DeploymentID string `json:"deployment_id"`
	ImageURL     string `json:"image_url"`
	Success      bool   `json:"success"`
	Logs         string `json:"logs"`
}

func New() (*Worker, error) {
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

	workspaceDir := os.Getenv("WORKSPACE_DIR")
	if workspaceDir == "" {
		workspaceDir = "/tmp/dejavu-builds"
	}

	cacheDir := os.Getenv("CACHE_DIR")
	if cacheDir == "" {
		cacheDir = "/tmp/dejavu-cache"
	}

	// Create directories
	os.MkdirAll(workspaceDir, 0755)
	os.MkdirAll(cacheDir, 0755)

	return &Worker{
		nats:         nc,
		js:           js,
		workspaceDir: workspaceDir,
		cacheDir:     cacheDir,
		registryURL:  os.Getenv("REGISTRY_URL"),
		registryUser: os.Getenv("REGISTRY_USERNAME"),
		registryPass: os.Getenv("REGISTRY_PASSWORD"),
	}, nil
}

func (w *Worker) Start() error {
	_, err := w.js.Subscribe("DEPLOYMENTS.request", func(msg *nats.Msg) {
		var event DeploymentEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Error parsing event: %v", err)
			msg.Ack()
			return
		}

		log.Printf("📦 Processing deployment: %s", event.DeploymentID)
		w.processBuild(event)
		msg.Ack()
	})

	return err
}

func (w *Worker) processBuild(event DeploymentEvent) {
	logs := ""
	success := false
	imageURL := ""

	defer func() {
		// Publish build complete event
		completeEvent := BuildCompleteEvent{
			DeploymentID: event.DeploymentID,
			ImageURL:     imageURL,
			Success:      success,
			Logs:         logs,
		}

		data, _ := json.Marshal(completeEvent)
		w.js.Publish("BUILDS.complete", data)
	}()

	// 1. Clone repository
	buildID := uuid.New().String()[:8]
	buildPath := filepath.Join(w.workspaceDir, buildID)
	defer os.RemoveAll(buildPath)

	logs += fmt.Sprintf("Cloning repository: %s\n", event.RepoURL)
	if err := w.cloneRepo(event.RepoURL, buildPath, event.CommitHash); err != nil {
		logs += fmt.Sprintf("Error cloning: %v\n", err)
		return
	}

	// 2. Detect framework
	logs += "Detecting framework...\n"
	framework := detector.Detect(buildPath)
	logs += fmt.Sprintf("Detected: %s\n", framework)

	// 3. Build project
	logs += "Building project...\n"
	buildRunner := runner.GetRunner(framework)
	if err := buildRunner.Build(buildPath, event.BuildCommand); err != nil {
		logs += fmt.Sprintf("Build failed: %v\n", err)
		return
	}
	logs += "Build completed successfully\n"

	// 4. Build Docker image
	logs += "Building Docker image...\n"
	imageName := fmt.Sprintf("%s/dejavu/%s", w.registryURL, event.ProjectID)
	imageTag := fmt.Sprintf("%s:%s", imageName, buildID)

	if err := w.buildDockerImage(buildPath, imageTag, framework, event.OutputDir); err != nil {
		logs += fmt.Sprintf("Docker build failed: %v\n", err)
		return
	}

	// 5. Push to registry
	logs += "Pushing to registry...\n"
	if err := w.pushImage(imageTag); err != nil {
		logs += fmt.Sprintf("Push failed: %v\n", err)
		return
	}

	logs += "✅ Deployment build complete\n"
	success = true
	imageURL = imageTag
}

func (w *Worker) cloneRepo(repoURL, destination, commitHash string) error {
	args := []string{"clone", "--depth", "1"}
	if commitHash != "" {
		args = append(args, "--branch", commitHash)
	}
	args = append(args, repoURL, destination)

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, output)
	}
	return nil
}

func (w *Worker) buildDockerImage(buildPath, imageTag, framework, outputDir string) error {
	// Create Dockerfile
	dockerfile := w.generateDockerfile(framework, outputDir)
	dockerfilePath := filepath.Join(buildPath, "Dockerfile.dejavu")
	if err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644); err != nil {
		return err
	}

	// Build image
	cmd := exec.Command("docker", "build", "-f", dockerfilePath, "-t", imageTag, ".")
	cmd.Dir = buildPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, output)
	}

	return nil
}

func (w *Worker) pushImage(imageTag string) error {
	// Login to registry
	if w.registryUser != "" && w.registryPass != "" {
		loginCmd := exec.Command("docker", "login", w.registryURL, "-u", w.registryUser, "-p", w.registryPass)
		if output, err := loginCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("registry login failed: %v: %s", err, output)
		}
	}

	// Push image
	cmd := exec.Command("docker", "push", imageTag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, output)
	}

	return nil
}

func (w *Worker) generateDockerfile(framework, outputDir string) string {
	switch framework {
	case "nextjs":
		return `FROM node:18-alpine
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]`

	case "nodejs":
		return fmt.Sprintf(`FROM node:18-alpine
WORKDIR /app
COPY . .
RUN npm install
EXPOSE 3000
CMD ["node", "index.js"]`)

	case "static":
		return fmt.Sprintf(`FROM nginx:alpine
COPY %s /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]`, outputDir)

	case "go":
		return `FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]`

	case "php":
		return `FROM php:8.2-apache
COPY . /var/www/html/
EXPOSE 80`

	default:
		return fmt.Sprintf(`FROM nginx:alpine
COPY %s /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]`, outputDir)
	}
}

func (w *Worker) Close() {
	if w.nats != nil {
		w.nats.Close()
	}
}

