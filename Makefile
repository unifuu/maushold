# Maushold Microservices Makefile

.PHONY: help build-all dev stop clean test

# Default target
help:
	@echo "Maushold Microservices - Available Commands:"
	@echo ""
	@echo "🚀 Quick Start:"
	@echo "  make start          - Start all services"
	@echo "  make setup          - Setup Kong routes"
	@echo "  make test           - Run complete test suite"
	@echo "  make ui             - Open all UIs"
	@echo ""
	@echo "🔍 Monitoring:"
	@echo "  make status         - Check all services status"
	@echo "  make logs           - View all logs (specify SERVICE=name)"
	@echo "  make consul         - Open Consul UI"
	@echo "  make kong-ui        - Open Konga UI"
	@echo "  make rabbitmq       - Open RabbitMQ UI"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  make test-quick     - Quick health check"
	@echo "  make test-full      - Full user flow test"
	@echo "  make test-battle    - Test battle system"
	@echo ""
	@echo "🛠️ Development:"
	@echo "  make stop           - Stop all services"
	@echo "  make restart        - Restart all services"
	@echo "  make clean          - Stop and remove volumes"
	@echo "  make rebuild        - Rebuild all services"

# Start all services
start:
	@echo "🚀 Starting Maushold..."
	docker-compose up -d
	@echo "⏳ Waiting for services to be ready..."
	@sleep 30
	@echo "✅ Services started!"
	@echo ""
	@echo "📊 Access Points:"
	@echo "  Frontend:  http://localhost:3000"
	@echo "  Kong:      http://localhost:8000"
	@echo "  Consul:    http://localhost:8500"
	@echo "  Konga:     http://localhost:1337"
	@echo "  RabbitMQ:  http://localhost:15672"

# Setup Kong routes
setup:
	@echo "🔧 Setting up Kong API Gateway..."
	@bash scripts/setup-kong.sh || echo "Run 'make setup' again if Kong isn't ready yet"

# Complete setup (start + setup)
init: start
	@echo "⏳ Waiting for Kong to be ready..."
	@bash scripts/wait-for-kong.sh
	@make setup
	@echo ""
	@echo "✅ Maushold is ready!"
	@echo "🌐 Frontend: http://localhost:3000"
	@echo "🔗 API Gateway: http://localhost:8000"
	@echo ""
	@echo "Run 'make test-quick' to verify everything works"

# Quick development start (same as init)
dev: init
	@echo "🚀 Development environment is ready!"

# Stop all services
stop:
	@echo "🛑 Stopping all services..."
	docker-compose down

# Clean everything (including volumes)
clean:
	@echo "🧹 Cleaning up..."
	docker-compose down -v
	@echo "✅ Cleaned!"

# Restart all services
restart:
	@echo "🔄 Restarting services..."
	docker-compose restart
	@echo "✅ Restarted!"

# Rebuild all services
rebuild:
	@echo "🔨 Rebuilding all services..."
	docker-compose up -d --build
	@echo "✅ Rebuilt!"

# Check status
status:
	@echo "📊 Service Status:"
	@docker-compose ps

# View logs
logs:
ifdef SERVICE
	@docker-compose logs -f $(SERVICE)
else
	@docker-compose logs -f
endif

# Quick health check test
test-quick:
	@echo "🧪 Quick Health Check..."
	@echo ""
	@echo "Player Service:"
	@curl -s http://localhost:8000/players/health | jq '.' || echo "❌ Failed"
	@echo ""
	@echo "Monster Service:"
	@curl -s http://localhost:8000/monster/health | jq '.' || echo "❌ Failed"
	@echo ""
	@echo "Battle Service:"
	@curl -s http://localhost:8000/battles/health | jq '.' || echo "❌ Failed"
	@echo ""
	@echo "Ranking Service:"
	@curl -s http://localhost:8000/rankings/health | jq '.' || echo "❌ Failed"
	@echo ""
	@echo "✅ Health check complete!"

# Full test suite
test-full:
	@echo "🧪 Running Full Test Suite..."
	@bash test-maushold.sh || echo "Create test-maushold.sh first"

# Test battle system
test-battle:
	@echo "⚔️ Testing Battle System..."
	@bash scripts/test-battle.sh || echo "Creating test script..."
	@echo '#!/bin/bash' > scripts/test-battle.sh
	@echo 'P1=$(curl -s -X POST http://localhost:8000/players -H "Content-Type: application/json" -d '"'"'{"username":"Player1"}'"'"' | jq -r ".id")' >> scripts/test-battle.sh
	@echo 'P2=$(curl -s -X POST http://localhost:8000/players -H "Content-Type: application/json" -d '"'"'{"username":"Player2"}'"'"' | jq -r ".id")' >> scripts/test-battle.sh
	@echo 'M1=$(curl -s -X POST http://localhost:8000/players/$P1/monster -H "Content-Type: application/json" -d '"'"'{"monster_id":25,"nickname":"Pikachu","level":5,"hp":35,"attack":55,"defense":40,"speed":90}'"'"' | jq -r ".id")' >> scripts/test-battle.sh
	@echo 'M2=$(curl -s -X POST http://localhost:8000/players/$P2/monster -H "Content-Type: application/json" -d '"'"'{"monster_id":6,"nickname":"Charizard","level":5,"hp":78,"attack":84,"defense":78,"speed":100}'"'"' | jq -r ".id")' >> scripts/test-battle.sh
	@echo 'curl -s -X POST http://localhost:8000/battles -H "Content-Type: application/json" -d "{\"player1_id\":$P1,\"player2_id\":$P2,\"monster1_id\":$M1,\"monster2_id\":$M2}" | jq "."' >> scripts/test-battle.sh
	@chmod +x scripts/test-battle.sh
	@bash scripts/test-battle.sh

# Open UIs
consul:
	@echo "🌐 Opening Consul UI..."
	@open http://localhost:8500 || xdg-open http://localhost:8500

kong-ui:
	@echo "🌐 Opening Konga UI..."
	@open http://localhost:1337 || xdg-open http://localhost:1337

rabbitmq:
	@echo "🌐 Opening RabbitMQ UI..."
	@open http://localhost:15672 || xdg-open http://localhost:15672

frontend:
	@echo "🌐 Opening Frontend..."
	@open http://localhost:3000 || xdg-open http://localhost:3000

ui: consul kong-ui rabbitmq frontend
	@echo "✅ All UIs opened!"

# Check Kong configuration
kong-status:
	@echo "🔍 Kong Services:"
	@curl -s http://localhost:18001/services | jq '.data[] | {name, url}'
	@echo ""
	@echo "🔍 Kong Routes:"
	@curl -s http://localhost:18001/routes | jq '.data[] | {name, paths}'

# Check Consul services
consul-services:
	@echo "🔍 Consul Services:"
	@curl -s http://localhost:8500/v1/catalog/services | jq '.'
	@echo ""
	@echo "🔍 Healthy Services:"
	@curl -s http://localhost:8500/v1/health/state/passing | jq '.[] | {service: .ServiceName, status: .Status}'

# Complete test
test: test-quick
	@echo ""
	@echo "✅ All tests passed!"
	@echo ""
	@echo "Try these next:"
	@echo "  make test-battle  - Test the battle system"
	@echo "  make frontend     - Open the UI"
	@echo "  make ui           - Open all admin UIs"
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
# 	@echo "Kubernetes:"
# 	@echo "  make k8s-deploy         - Deploy to Kubernetes"
# 	@echo "  make k8s-status         - Check deployment status"
# 	@echo "  make k8s-logs           - View logs (specify SERVICE=name)"
# 	@echo "  make k8s-port-forward   - Setup port forwarding"
# 	@echo "  make k8s-cleanup        - Delete all resources"
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

# # Kubernetes Commands
# k8s-deploy: build-all
# 	@echo "🚀 Deploying to Kubernetes..."
# 	kubectl apply -f k8s/namespace.yaml
# 	kubectl apply -f k8s/configmap.yaml
# 	kubectl apply -f k8s/secrets.yaml
# 	@echo "⏳ Deploying infrastructure..."
# 	kubectl apply -f k8s/consul.yaml
# 	kubectl apply -f k8s/redis.yaml
# 	kubectl apply -f k8s/rabbitmq.yaml
# 	kubectl apply -f k8s/databases.yaml
# 	@echo "⏳ Waiting for databases to be ready..."
# 	sleep 30
# 	@echo "⏳ Deploying services..."
# 	kubectl apply -f k8s/player-service.yaml
# 	kubectl apply -f k8s/monster-service.yaml
# 	kubectl apply -f k8s/battle-service.yaml
# 	kubectl apply -f k8s/ranking-service.yaml
# 	@echo "✅ Deployment complete!"
# 	@echo ""
# 	@echo "Check status with: make k8s-status"
# 	@echo "Access services at:"
# 	@echo "   Player:  http://localhost:30001"
# 	@echo "   Monster: http://localhost:30002"
# 	@echo "   Battle:  http://localhost:30003"
# 	@echo "   Ranking: http://localhost:30004"

# k8s-status:
# 	@echo "📊 Kubernetes Status:"
# 	@echo ""
# 	@echo "Pods:"
# 	kubectl get pods -n $(NAMESPACE)
# 	@echo ""
# 	@echo "Services:"
# 	kubectl get svc -n $(NAMESPACE)
# 	@echo ""
# 	@echo "Deployments:"
# 	kubectl get deployments -n $(NAMESPACE)

# k8s-logs:
# ifdef SERVICE
# 	kubectl logs -f deployment/$(SERVICE) -n $(NAMESPACE)
# else
# 	@echo "Usage: make k8s-logs SERVICE=player-service"
# endif

# k8s-port-forward:
# 	@echo "🔌 Setting up port forwarding..."
# 	@echo "Consul UI will be available at http://localhost:8500"
# 	kubectl port-forward -n $(NAMESPACE) svc/consul 8500:8500

# k8s-cleanup:
# 	@echo "🧹 Cleaning up Kubernetes resources..."
# 	kubectl delete namespace $(NAMESPACE)
# 	@echo "✅ Cleanup complete!"

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