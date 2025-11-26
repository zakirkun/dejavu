# 🪟 Panduan Start Dejavu di Windows

## Prerequisites

✅ Docker Desktop dengan Kubernetes enabled
✅ PostgreSQL, Redis, NATS, MinIO running (via Docker Compose)
✅ Go 1.21+ installed
✅ Node.js 18+ installed

## Quick Start

### 1. Start Infrastructure

```powershell
# Terminal 1 - Start dependencies
cd C:\Users\"MyBook Hype AMD"\workarea\dejavu
docker compose -f infra-compose.yml up -d

# Wait 30 seconds for services to start
Start-Sleep -Seconds 30
```

### 2. Run Database Migrations

```powershell
# Terminal 1 (lanjutan)
cd backend
go run cmd/migrate/main.go up
```

### 3. Start Backend API

```powershell
# Terminal 2 - Backend API
cd C:\Users\"MyBook Hype AMD"\workarea\dejavu\backend
go run cmd/api/main.go

# Backend akan running di: http://localhost:8080
```

### 4. Start Builder Worker

```powershell
# Terminal 3 - Builder Worker
cd C:\Users\"MyBook Hype AMD"\workarea\dejavu\builder
go run cmd/worker/main.go

# Builder akan listen untuk build jobs dari NATS
```

### 5. Start Deployer Worker

```powershell
# Terminal 4 - Deployer Worker
cd C:\Users\"MyBook Hype AMD"\workarea\dejavu\deployer

# IMPORTANT: Set KUBECONFIG environment variable
$env:KUBECONFIG = "C:\Users\MyBook Hype AMD\.kube\config"

go run cmd/worker/main.go

# Deployer akan listen untuk deployment jobs dari NATS
```

### 6. Start Frontend

```powershell
# Terminal 5 - Frontend Next.js
cd C:\Users\"MyBook Hype AMD"\workarea\dejavu\frontend
npm install  # first time only
npm run dev

# Frontend akan running di: http://localhost:3000
```

## Environment Variables for Deployer

Jika tidak mau set `$env:KUBECONFIG` setiap kali, tambahkan ke PowerShell profile:

```powershell
# Edit profile
notepad $PROFILE

# Tambahkan line ini:
$env:KUBECONFIG = "C:\Users\MyBook Hype AMD\.kube\config"

# Save dan reload:
. $PROFILE
```

## Verification

### Check Services Running:

```powershell
# Check infrastructure
docker ps

# Check Kubernetes
kubectl cluster-info

# Test backend API
curl http://localhost:8080/health

# Test frontend
Start-Process "http://localhost:3000"
```

### Check Logs:

```powershell
# Infrastructure logs
docker compose -f infra-compose.yml logs -f

# Kubernetes logs (if deploying)
kubectl get pods -A
kubectl logs -f <pod-name> -n dejavu-apps
```

## Troubleshooting

### Issue: KUBECONFIG not found

**Solution:**
```powershell
$env:KUBECONFIG = "C:\Users\MyBook Hype AMD\.kube\config"
```

### Issue: Port already in use

**Solution:**
```powershell
# Find process using port 8080
netstat -ano | findstr :8080

# Kill the process
taskkill /PID <PID> /F
```

### Issue: Docker containers not starting

**Solution:**
```powershell
# Restart Docker Desktop
# Or clean up:
docker system prune -a
docker compose -f infra-compose.yml down
docker compose -f infra-compose.yml up -d
```

### Issue: Go modules error

**Solution:**
```powershell
cd backend  # or builder, deployer
go mod tidy
go mod download
```

### Issue: npm install fails

**Solution:**
```powershell
cd frontend
Remove-Item -Recurse -Force node_modules
Remove-Item package-lock.json
npm install
```

## Quick Commands

```powershell
# Stop all services
docker compose -f infra-compose.yml down

# Restart infrastructure
docker compose -f infra-compose.yml restart

# View all logs
docker compose -f infra-compose.yml logs -f

# Clean everything
docker compose -f infra-compose.yml down -v
Remove-Item -Recurse -Force backend/bin, builder/bin, deployer/bin
```

## Production Deployment

Untuk production di Windows Server, gunakan:
- Windows Server with Docker
- Kubernetes (AKS, EKS, atau on-premise)
- Lihat [DEPLOYMENT.md](DEPLOYMENT.md) untuk panduan lengkap

## Support

Jika ada issue, check:
1. Docker Desktop running
2. Kubernetes enabled di Docker Desktop
3. All .env files exist
4. Port 8080, 3000 tidak dipakai
5. Antivirus tidak block Go/Node

---

**Happy Deploying! 🚀**

