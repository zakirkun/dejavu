# 🚀 Panduan Deployment Dejavu

Panduan lengkap untuk deploy Dejavu platform dari awal.

## Prerequisites

- Docker & Docker Compose
- Kubernetes cluster (minikube, k3s, k3d, atau cloud provider)
- kubectl configured
- Go 1.21+
- Node.js 18+
- Git

## Setup Development Environment

### 1. Clone Repository

```bash
git clone <your-repo-url>
cd dejavu
```

### 2. Start Infrastructure Dependencies

```bash
# Start PostgreSQL, Redis, NATS, MinIO
make infra

# Wait for services to be ready (30 seconds)
sleep 30
```

### 3. Setup Database

```bash
# Create .env file untuk backend
cp backend/.env.example backend/.env

# Edit backend/.env dan sesuaikan dengan environment Anda

# Run migrations
cd backend
go mod download
go run cmd/migrate/main.go up
cd ..
```

### 4. Start Backend Services

**Terminal 1 - Backend API:**
```bash
cd backend
go run cmd/api/main.go
```

**Terminal 2 - Builder Worker:**
```bash
cd builder
cp .env.example .env
go mod download
go run cmd/worker/main.go
```

**Terminal 3 - Deployer Worker:**
```bash
cd deployer
cp .env.example .env
go mod download
go run cmd/worker/main.go
```

### 5. Start Frontend

**Terminal 4 - Next.js:**
```bash
cd frontend
npm install
npm run dev
```

### 6. Setup Local DNS

Untuk wildcard subdomain `*.dejavu.local`:

**macOS:**
```bash
# Install dnsmasq
brew install dnsmasq

# Configure
echo "address=/.dejavu.local/127.0.0.1" | sudo tee -a /usr/local/etc/dnsmasq.conf

# Start service
sudo brew services start dnsmasq

# Add resolver
sudo mkdir -p /etc/resolver
sudo tee /etc/resolver/dejavu.local <<EOF
nameserver 127.0.0.1
EOF

# Test
ping test.dejavu.local
```

**Linux:**
```bash
# Install dnsmasq
sudo apt install dnsmasq

# Configure
echo "address=/.dejavu.local/127.0.0.1" | sudo tee -a /etc/dnsmasq.conf

# Restart
sudo systemctl restart dnsmasq

# Test
dig test.dejavu.local
```

Lihat [infra/dns-local/README.md](infra/dns-local/README.md) untuk panduan lengkap.

## Kubernetes Deployment

### 1. Setup Kubernetes Cluster

Jika belum punya cluster, gunakan minikube:

```bash
# Install minikube
brew install minikube  # macOS
# atau
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# Start cluster
minikube start --cpus=4 --memory=8192

# Enable ingress
minikube addons enable ingress
```

### 2. Deploy Core Infrastructure

```bash
# Create namespaces
kubectl apply -f infra/kubernetes/namespace.yaml

# Deploy Traefik
kubectl apply -f infra/kubernetes/traefik/

# Deploy Docker Registry
kubectl apply -f infra/kubernetes/docker-registry/

# Wait for Traefik to be ready
kubectl wait --for=condition=ready pod -l app=traefik -n dejavu-system --timeout=300s
```

### 3. Deploy Cert-Manager (Optional for production)

```bash
# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Wait for cert-manager
kubectl wait --for=condition=ready pod -l app=cert-manager -n cert-manager --timeout=300s

# Apply ClusterIssuers
kubectl apply -f infra/kubernetes/cert-manager/
```

### 4. Deploy Monitoring Stack (Optional)

```bash
# Deploy Prometheus
kubectl apply -f infra/kubernetes/monitoring/prometheus.yaml

# Deploy Loki
kubectl apply -f infra/kubernetes/monitoring/loki.yaml

# Deploy Promtail (log collector)
kubectl apply -f infra/kubernetes/monitoring/promtail.yaml

# Deploy Grafana
kubectl apply -f infra/kubernetes/monitoring/grafana.yaml

# Access Grafana
kubectl port-forward svc/grafana -n dejavu-system 3000:3000
# Open http://localhost:3000
# Username: admin, Password: admin
```

### 5. Configure Docker Registry Access

```bash
# Get registry IP
kubectl get svc docker-registry -n dejavu-system

# Add to /etc/hosts (untuk local development)
echo "$(minikube ip) registry.dejavu.local" | sudo tee -a /etc/hosts

# Configure Docker to trust insecure registry
# macOS: Docker Desktop > Preferences > Docker Engine
# Linux: /etc/docker/daemon.json
{
  "insecure-registries": ["registry.dejavu.local:5000"]
}

# Restart Docker
```

### 6. Deploy Application Services

**Option A: Manual Deploy**

Build dan push images:
```bash
# Backend
cd backend
docker build -t registry.dejavu.local:5000/dejavu/backend:latest .
docker push registry.dejavu.local:5000/dejavu/backend:latest

# Builder
cd ../builder
docker build -t registry.dejavu.local:5000/dejavu/builder:latest .
docker push registry.dejavu.local:5000/dejavu/builder:latest

# Deployer
cd ../deployer
docker build -t registry.dejavu.local:5000/dejavu/deployer:latest .
docker push registry.dejavu.local:5000/dejavu/deployer:latest

# Frontend
cd ../frontend
docker build -t registry.dejavu.local:5000/dejavu/frontend:latest .
docker push registry.dejavu.local:5000/dejavu/frontend:latest
```

**Option B: Using Skaffold (Recommended)**

```bash
# Install Skaffold
brew install skaffold  # macOS
# atau
curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
sudo install skaffold /usr/local/bin/

# Deploy dengan Skaffold
skaffold dev
```

## Production Deployment

### Cloud Providers

**AWS (EKS):**
```bash
# Create cluster
eksctl create cluster --name dejavu --region us-east-1 --nodes 3

# Deploy
kubectl apply -f infra/kubernetes/
```

**GCP (GKE):**
```bash
# Create cluster
gcloud container clusters create dejavu --num-nodes=3

# Deploy
kubectl apply -f infra/kubernetes/
```

**DigitalOcean (DOKS):**
```bash
# Create via Dashboard or doctl
doctl kubernetes cluster create dejavu --count 3

# Deploy
kubectl apply -f infra/kubernetes/
```

### DNS Configuration

Untuk production dengan wildcard subdomain:

1. Point `*.dejavu.id` ke LoadBalancer IP:
```bash
# Get LoadBalancer IP
kubectl get svc traefik -n dejavu-system

# Add DNS A record
*.dejavu.id → <LOADBALANCER_IP>
```

2. Update BASE_DOMAIN di environment variables:
```bash
# deployer/.env
BASE_DOMAIN=dejavu.id

# frontend/.env
NEXT_PUBLIC_API_URL=https://api.dejavu.id
```

### SSL/TLS Setup

Cert-Manager akan otomatis request certificates dari Let's Encrypt:

```bash
# Check certificate status
kubectl get certificate -n dejavu-apps

# Check certificate details
kubectl describe certificate <cert-name> -n dejavu-apps
```

## Troubleshooting

### Builder tidak bisa clone private repos

```bash
# Add SSH key atau GitHub token
kubectl create secret generic git-credentials \
  --from-literal=username=<github-username> \
  --from-literal=password=<github-token> \
  -n dejavu-system
```

### Docker build fails dengan "no space left"

```bash
# Clean up Docker
docker system prune -a

# Increase disk di minikube
minikube delete
minikube start --disk-size=50g
```

### Deployment stuck di "pending"

```bash
# Check pod logs
kubectl logs -f <pod-name> -n dejavu-apps

# Check events
kubectl get events -n dejavu-apps --sort-by='.lastTimestamp'

# Check resource availability
kubectl top nodes
```

### Logs tidak muncul di Grafana

```bash
# Check Promtail logs
kubectl logs -f daemonset/promtail -n dejavu-system

# Check Loki logs
kubectl logs -f deployment/loki -n dejavu-system

# Restart Promtail
kubectl rollout restart daemonset/promtail -n dejavu-system
```

## Backup & Recovery

### Database Backup

```bash
# Backup PostgreSQL
kubectl exec -it postgres-pod -n dejavu-system -- \
  pg_dump -U dejavu dejavu > backup-$(date +%Y%m%d).sql

# Restore
kubectl exec -i postgres-pod -n dejavu-system -- \
  psql -U dejavu dejavu < backup-20240101.sql
```

### Persistent Volume Backup

```bash
# List PVCs
kubectl get pvc -n dejavu-system

# Backup using Velero (recommended)
velero install
velero backup create dejavu-backup
```

## Monitoring & Alerts

### Setup Prometheus Alerts

```yaml
# prometheus-rules.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-rules
  namespace: dejavu-system
data:
  alert.rules: |
    groups:
    - name: dejavu
      rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status="500"}[5m]) > 0.05
        annotations:
          summary: "High error rate detected"
```

### Grafana Dashboards

Akses Grafana:
```bash
kubectl port-forward svc/grafana -n dejavu-system 3000:3000
```

Import dashboard untuk:
- Kubernetes cluster metrics
- Application metrics
- Deployment logs
- Build metrics

## Performance Tuning

### Horizontal Pod Autoscaler

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: dejavu-api
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: dejavu-api
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### Database Connection Pooling

```go
// backend/pkg/database/database.go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

## Security Best Practices

1. **Gunakan secrets untuk sensitive data:**
```bash
kubectl create secret generic dejavu-secrets \
  --from-literal=db-password=<password> \
  --from-literal=jwt-secret=<secret>
```

2. **Enable RBAC:**
```bash
kubectl apply -f infra/kubernetes/rbac/
```

3. **Network Policies:**
```bash
kubectl apply -f infra/kubernetes/network-policies/
```

4. **Image scanning:**
```bash
trivy image registry.dejavu.id/dejavu/backend:latest
```

## Support

Untuk pertanyaan atau issue:
- GitHub Issues: <your-repo-url>/issues
- Email: support@dejavu.id
- Documentation: https://docs.dejavu.id

