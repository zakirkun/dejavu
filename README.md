# 🚀 Dejavu - Platform Deployment SaaS

Platform hosting & deployment otomatis mirip Vercel, dibangun dengan Go, Kubernetes, dan Next.js.

## 📋 Daftar Isi

- [Fitur](#fitur)
- [Arsitektur](#arsitektur)
- [Tech Stack](#tech-stack)
- [Quick Start](#quick-start)
- [Development](#development)
- [Deployment](#deployment)

## ✨ Fitur

- 🔄 **Auto Deployment** - Deploy otomatis dari Git repository
- 🏗️ **Framework Detection** - Deteksi otomatis Next.js, Node, Bun, Go, PHP
- 🌐 **Wildcard Subdomain** - Setiap deployment dapat subdomain unik (*.dejavu.id)
- 📊 **Real-time Logs** - Lihat build & deployment logs secara real-time
- ⚡ **Zero Downtime** - Rolling update tanpa downtime
- 📈 **Auto Scaling** - Horizontal pod autoscaling otomatis
- 💳 **Billing System** - Credit-based usage tracking
- 🔐 **TLS Otomatis** - HTTPS dengan cert-manager
- 📦 **Build Caching** - Build lebih cepat dengan caching
- 🎯 **Multi-Project** - Kelola banyak project dalam satu dashboard

## 🏛️ Arsitektur

```
User → Frontend → Backend API → NATS → Builder Worker → Docker Registry
                                   ↓
                              Deployer Worker → Kubernetes → Traefik → *.dejavu.id
```

### Komponen

- **Backend** - REST API (Go Fiber)
- **Builder** - Build worker yang clone, detect, dan build project
- **Deployer** - Deploy worker yang manage Kubernetes resources
- **Frontend** - Dashboard & landing page (Next.js)
- **Infrastructure** - Kubernetes dengan Traefik & Cert-Manager

## 🛠️ Tech Stack

### Backend
- Go 1.21+
- Fiber (Web Framework)
- PostgreSQL (Database)
- Redis (Cache)
- NATS JetStream (Message Queue)
- MinIO (Object Storage)

### Frontend
- Next.js 14+
- Tailwind CSS
- Shadcn/ui
- WebSocket

### Infrastructure
- Docker
- Kubernetes
- Traefik (Ingress)
- Cert-Manager (TLS)
- Docker Registry

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Kubernetes (minikube, k3s, atau kind)
- Go 1.21+
- Node.js 18+
- kubectl

### 1. Clone Repository

```bash
git clone https://github.com/zakirkun/dejavu.git
cd dejavu
```

### 2. Setup Environment

```bash
# Copy environment files
cp backend/.env.example backend/.env
cp builder/.env.example builder/.env
cp deployer/.env.example deployer/.env
cp frontend/.env.example frontend/.env.local
```

### 3. Start Infrastructure

```bash
# Start PostgreSQL, Redis, NATS, MinIO
make infra
```

### 4. Run Migrations

```bash
make migrate-up
```

### 5. Start Services

Di terminal terpisah, jalankan:

```bash
# Terminal 1 - Backend API
make backend

# Terminal 2 - Builder Worker
make builder

# Terminal 3 - Deployer Worker
make deployer

# Terminal 4 - Frontend
make frontend
```

### 6. Akses Dashboard

Buka browser: http://localhost:3000

## 💻 Development

### Struktur Folder

```
dejavu/
├── backend/          # Go API
├── builder/          # Build worker
├── deployer/         # Deploy worker
├── infra/            # Kubernetes manifests
├── frontend/         # Next.js app
├── infra-compose.yml # Docker Compose
└── Makefile          # Commands
```

### Backend API

```bash
cd backend
go run cmd/api/main.go
```

API akan berjalan di: http://localhost:8080

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend akan berjalan di: http://localhost:3000

### Builder Worker

```bash
cd builder
go run cmd/worker/main.go
```

### Deployer Worker

```bash
cd deployer
go run cmd/worker/main.go
```

## 🎯 API Endpoints

### Authentication
- `POST /api/auth/register` - Register user baru
- `POST /api/auth/login` - Login user

### Projects
- `GET /api/projects` - List semua projects
- `POST /api/projects` - Buat project baru
- `GET /api/projects/:id` - Detail project
- `PUT /api/projects/:id` - Update project
- `DELETE /api/projects/:id` - Hapus project

### Deployments
- `POST /api/deploy` - Trigger deployment
- `GET /api/deploy/:id` - Status deployment
- `GET /api/deploy/:id/logs` - Stream logs (WebSocket)

## 🐳 Docker Compose Services

- **PostgreSQL** - Port 5432
- **Redis** - Port 6379
- **NATS** - Port 4222
- **MinIO** - Port 9000 (API), 9001 (Console)
- **Adminer** - Port 8081

## ☸️ Kubernetes Deployment

### 1. Setup Infrastructure

```bash
make k8s-infra
```

### 2. Deploy Services

```bash
make deploy
```

### 3. Check Logs

```bash
make logs
```

## 🌐 Local DNS Setup

Untuk menggunakan wildcard subdomain `*.dejavu.local`:

### macOS/Linux

```bash
# Install dnsmasq
brew install dnsmasq  # macOS
# atau
sudo apt install dnsmasq  # Ubuntu

# Configure
echo "address=/.dejavu.local/127.0.0.1" | sudo tee -a /etc/dnsmasq.conf

# Restart
sudo brew services restart dnsmasq  # macOS
# atau
sudo systemctl restart dnsmasq  # Linux
```

### Windows

Edit `C:\Windows\System32\drivers\etc\hosts`:

```
127.0.0.1 dejavu.local
127.0.0.1 app1.dejavu.local
127.0.0.1 app2.dejavu.local
```

## 🧪 Testing

```bash
make test
```

## 📝 License

MIT License

## 🤝 Contributing

Contributions are welcome! Please open an issue or submit a PR.

## 📧 Support

Untuk pertanyaan atau dukungan, hubungi: support@dejavu.id

