# Maushold Microservices Makefile

.PHONY: help build-all build-player build-monster build-battle build-ranking \
        docker-up docker-down docker-logs \
        k8s-deploy k8s-status k8s-logs k8s-cleanup \
		kong-setup kong-status kong-routes kong-plugins \
        test lint tidy clean

# Default target
help:
	@echo "Maushold Microservices - Available Commands:"
	@echo ""
	@echo "Local Development (Docker Compose):"
	@echo "  make docker-up          - Start all services with Docker Compose"
	@echo "  make docker-down        - Stop all services"
	@echo "  make docker-logs        - View logs (specify SERVICE=name)"
	@echo "  make docker-restart     - Restart all services"
	@echo ""
	@echo "Build Commands:"
	@echo "  make build-all          - Build all service images"
	@echo "  make build-player       - Build player service"
	@echo "  make build-monster      - Build monster service"
	@echo "  make build-battle       - Build battle service"
	@echo "  make build-ranking      - Build ranking service"
	@echo ""
	@echo "Kubernetes:"
	@echo "  make k8s-deploy         - Deploy to Kubernetes"
	@echo "  make k8s-status         - Check deployment status"
	@echo "  make k8s-logs           - View logs (specify SERVICE=name)"
	@echo "  make k8s-port-forward   - Setup port forwarding"
	@echo "  make k8s-cleanup        - Delete all resources"
	@echo ""
	@echo "Development:"
	@echo "  make tidy               - Run go mod tidy on all services"
	@echo "  make test               - Run tests"
	@echo "  make lint               - Run linter"
	@echo "  make clean              - Clean build artifacts"
	@echo ""
	@echo "Consul & Monitoring:"
	@echo "  make consul-ui          - Port-forward Consul UI (localhost:8500)"
	@echo "  make rabbitmq-ui        - Port-forward RabbitMQ UI (localhost:15672)"

# Variables
SERVICES = player-service monster-service battle-service ranking-service
NAMESPACE = maushold
IMAGE_TAG ?= latest

# Docker Compose Commands
docker-up:
	@echo "🚀 Starting all services with Docker Compose..."
	docker-compose up --build -d
	@echo "✅ Services started! Access at:"
	@echo "   Player:  http://localhost:8001"
	@echo "   Monster: http://localhost:8002"
	@echo "   Battle:  http://localhost:8003"
	@echo "   Ranking: http://localhost:8004"
	@echo "   Frontend: http://localhost:3000"

docker-down:
	@echo "🛑 Stopping all services..."
	docker-compose down

docker-restart:
	@echo "🔄 Restarting services..."
	docker-compose restart

docker-logs:
ifdef SERVICE
	docker-compose logs -f $(SERVICE)
else
	docker-compose logs -f
endif

# Build Commands
build-all: build-player build-monster build-battle build-ranking
	@echo "✅ All services built successfully!"

build-player:
	@echo "🔨 Building player-service..."
	cd services/player-service && docker build -t maushold/player-service:$(IMAGE_TAG) .

build-monster:
	@echo "🔨 Building monster-service..."
	cd services/monster-service && docker build -t maushold/monster-service:$(IMAGE_TAG) .

build-battle:
	@echo "🔨 Building battle-service..."
	cd services/battle-service && docker build -t maushold/battle-service:$(IMAGE_TAG) .

build-ranking:
	@echo "🔨 Building ranking-service..."
	cd services/ranking-service && docker build -t maushold/ranking-service:$(IMAGE_TAG) .

# Kubernetes Commands
k8s-deploy: build-all
	@echo "🚀 Deploying to Kubernetes..."
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/configmap.yaml
	kubectl apply -f k8s/secrets.yaml
	@echo "⏳ Deploying infrastructure..."
	kubectl apply -f k8s/consul.yaml
	kubectl apply -f k8s/redis.yaml
	kubectl apply -f k8s/rabbitmq.yaml
	kubectl apply -f k8s/databases.yaml
	@echo "⏳ Waiting for databases to be ready..."
	sleep 30
	@echo "⏳ Deploying services..."
	kubectl apply -f k8s/player-service.yaml
	kubectl apply -f k8s/monster-service.yaml
	kubectl apply -f k8s/battle-service.yaml
	kubectl apply -f k8s/ranking-service.yaml
	@echo "✅ Deployment complete!"
	@echo ""
	@echo "Check status with: make k8s-status"
	@echo "Access services at:"
	@echo "   Player:  http://localhost:30001"
	@echo "   Monster: http://localhost:30002"
	@echo "   Battle:  http://localhost:30003"
	@echo "   Ranking: http://localhost:30004"

k8s-status:
	@echo "📊 Kubernetes Status:"
	@echo ""
	@echo "Pods:"
	kubectl get pods -n $(NAMESPACE)
	@echo ""
	@echo "Services:"
	kubectl get svc -n $(NAMESPACE)
	@echo ""
	@echo "Deployments:"
	kubectl get deployments -n $(NAMESPACE)

k8s-logs:
ifdef SERVICE
	kubectl logs -f deployment/$(SERVICE) -n $(NAMESPACE)
else
	@echo "Usage: make k8s-logs SERVICE=player-service"
endif

k8s-port-forward:
	@echo "🔌 Setting up port forwarding..."
	@echo "Consul UI will be available at http://localhost:8500"
	kubectl port-forward -n $(NAMESPACE) svc/consul 8500:8500

k8s-cleanup:
	@echo "🧹 Cleaning up Kubernetes resources..."
	kubectl delete namespace $(NAMESPACE)
	@echo "✅ Cleanup complete!"

# Consul & Monitoring
consul-ui:
	@echo "🔌 Port-forwarding Consul UI..."
	@echo "Access at: http://localhost:8500"
	kubectl port-forward -n $(NAMESPACE) svc/consul 8500:8500

rabbitmq-ui:
	@echo "🔌 Port-forwarding RabbitMQ Management..."
	@echo "Access at: http://localhost:15672"
	@echo "Default credentials: maushold / changeme"
	kubectl port-forward -n $(NAMESPACE) svc/rabbitmq 15672:15672

# Development Commands
tidy:
	@echo "📦 Running go mod tidy on all services..."
	@for service in $(SERVICES); do \
		echo "  ↳ $$service"; \
		cd services/$$service && go mod tidy && cd ../..; \
	done
	@echo "✅ Done!"

test:
	@echo "🧪 Running tests..."
	@for service in $(SERVICES); do \
		echo "  ↳ Testing $$service"; \
		cd services/$$service && go test ./... && cd ../..; \
	done

lint:
	@echo "🔍 Running linter..."
	@for service in $(SERVICES); do \
		echo "  ↳ Linting $$service"; \
		cd services/$$service && golangci-lint run && cd ../..; \
	done

clean:
	@echo "🧹 Cleaning build artifacts..."
	@for service in $(SERVICES); do \
		cd services/$$service && rm -f $$service && cd ../..; \
	done
	docker system prune -f
	@echo "✅ Cleaned!"

# Kong Commands
kong-setup:
	@echo "🚀 Setting up Kong API Gateway..."
	chmod +x scripts/setup-kong.sh
	./scripts/setup-kong.sh

kong-status:
	@echo "📊 Kong Status:"
	@curl -s http://localhost:8001/ | jq '.'

kong-routes:
	@echo "🔍 Kong Routes:"
	@curl -s http://localhost:8001/routes | jq '.data[] | {name: .name, paths: .paths, service: .service.id}'

kong-services:
	@echo "🔍 Kong Services:"
	@curl -s http://localhost:8001/services | jq '.data[] | {name: .name, url: .url}'

kong-plugins:
	@echo "🔌 Kong Plugins:"
	@curl -s http://localhost:8001/plugins | jq '.data[] | {name: .name, enabled: .enabled}'

kong-test:
	@echo "🧪 Testing Kong Gateway..."
	@echo "Player Service:"
	@curl -s http://localhost:8000/api/players | jq '.'
	@echo ""
	@echo "Monster Service:"
	@curl -s http://localhost:8000/api/monster | jq '.'
	@echo ""
	@echo "Health Check:"
	@curl -s http://localhost:8000/api/players/health | jq '.'

# Quick Commands
dev: docker-up
	@echo "💻 Development environment started!"

stop: docker-down
	@echo "🛑 Development environment stopped!"

# Database Commands
db-migrate:
	@echo "🗄️ Running database migrations..."
	@echo "Migrations run automatically on service startup"

# Check Prerequisites
check:
	@echo "🔍 Checking prerequisites..."
	@command -v docker >/dev/null 2>&1 || { echo "❌ Docker not found"; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "❌ Docker Compose not found"; exit 1; }
	@command -v kubectl >/dev/null 2>&1 || { echo "❌ kubectl not found"; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo "❌ Go not found"; exit 1; }
	@echo "✅ All prerequisites installed!"

# Setup Commands
setup:
	@echo "🔧 Setting up project..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "⚠️  Created .env file - please update passwords!"; \
	fi
	@echo "📦 Installing Go dependencies..."
	make tidy
	@echo "✅ Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Edit .env file with secure passwords"
	@echo "  2. Run 'make dev' to start development environment"
	@echo "  3. Or run 'make k8s-deploy' to deploy to Kubernetes"

# Complete Workflow
all: check setup build-all
	@echo "🎉 Project ready!"

# Minikube specific commands
minikube-start:
	@echo "🚀 Starting Minikube..."
	minikube start --cpus=4 --memory=8192
	@echo "🐳 Configuring Docker environment..."
	eval $$(minikube docker-env)

minikube-stop:
	@echo "🛑 Stopping Minikube..."
	minikube stop

minikube-dashboard:
	@echo "📊 Opening Kubernetes Dashboard..."
	minikube dashboard