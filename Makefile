.PHONY: help infra infra-down backend builder deployer frontend dev clean deploy logs test

help:
	@echo "Dejavu - Deployment Platform Commands"
	@echo ""
	@echo "Infrastructure:"
	@echo "  make infra          - Start infrastructure (Postgres, Redis, NATS, MinIO)"
	@echo "  make infra-down     - Stop infrastructure"
	@echo ""
	@echo "Development:"
	@echo "  make backend        - Run backend API"
	@echo "  make builder        - Run builder worker"
	@echo "  make deployer       - Run deployer worker"
	@echo "  make frontend       - Run frontend dev server"
	@echo "  make dev            - Run all services"
	@echo ""
	@echo "Kubernetes:"
	@echo "  make deploy         - Deploy to Kubernetes"
	@echo "  make k8s-infra      - Setup Kubernetes infrastructure"
	@echo "  make logs           - Show logs"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback migrations"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make test           - Run tests"

# Infrastructure
infra:
	docker compose -f infra-compose.yml up -d

infra-down:
	docker compose -f infra-compose.yml down

# Development
backend:
	cd backend && go run cmd/api/main.go

builder:
	cd builder && go run cmd/worker/main.go

deployer:
	cd deployer && go run cmd/worker/main.go

frontend:
	cd frontend && npm run dev

dev:
	@echo "Starting all services..."
	$(MAKE) infra
	@echo "Waiting for infrastructure to be ready..."
	@sleep 5
	@echo "Start backend, builder, deployer, and frontend in separate terminals"

# Kubernetes
k8s-infra:
	kubectl apply -f infra/kubernetes/namespace.yaml
	kubectl apply -f infra/kubernetes/traefik/
	kubectl apply -f infra/kubernetes/cert-manager/

deploy:
	kubectl apply -f infra/kubernetes/

logs:
	kubectl logs -f -l app=dejavu-api -n dejavu-system

# Database
migrate-up:
	cd backend && go run cmd/migrate/main.go up

migrate-down:
	cd backend && go run cmd/migrate/main.go down

# Utilities
clean:
	rm -rf backend/bin
	rm -rf builder/bin
	rm -rf deployer/bin
	rm -rf frontend/.next
	rm -rf frontend/out

test:
	cd backend && go test ./...
	cd builder && go test ./...
	cd deployer && go test ./...
	cd frontend && npm test

