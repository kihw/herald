# Herald.lol Development Makefile

.PHONY: help dev build test clean logs restart

# Default target
help:
	@echo "Herald.lol Development Commands:"
	@echo ""
	@echo "  dev         Start development environment"
	@echo "  build       Build all services"
	@echo "  test        Run all tests"
	@echo "  clean       Clean up containers and volumes"
	@echo "  logs        Show service logs"
	@echo "  restart     Restart all services"
	@echo "  db-reset    Reset database with fresh data"
	@echo "  db-migrate  Run database migrations"
	@echo ""

# Start development environment
dev:
	@echo "🚀 Starting Herald.lol development environment..."
	docker-compose -f docker-compose.dev.yml up -d
	@echo "✅ Services started:"
	@echo "   - API: http://localhost:8080"
	@echo "   - Frontend: http://localhost:3000" 
	@echo "   - App: http://localhost (NGINX)"
	@echo "   - Grafana: http://localhost:3001"
	@echo "   - Prometheus: http://localhost:9090"

# Build all services
build:
	@echo "🔨 Building Herald.lol services..."
	docker-compose -f docker-compose.dev.yml build

# Run tests
test:
	@echo "🧪 Running tests..."
	docker-compose -f docker-compose.dev.yml exec herald-api go test ./... -v
	docker-compose -f docker-compose.dev.yml exec herald-frontend npm test -- --coverage --watchAll=false

# Clean up
clean:
	@echo "🧹 Cleaning up Herald.lol environment..."
	docker-compose -f docker-compose.dev.yml down -v
	docker system prune -f
	docker volume prune -f

# Show logs
logs:
	docker-compose -f docker-compose.dev.yml logs -f

# Show logs for specific service
logs-%:
	docker-compose -f docker-compose.dev.yml logs -f $*

# Restart services
restart:
	@echo "🔄 Restarting Herald.lol services..."
	docker-compose -f docker-compose.dev.yml restart

# Reset database
db-reset:
	@echo "🗄️ Resetting database..."
	docker-compose -f docker-compose.dev.yml stop postgres
	docker-compose -f docker-compose.dev.yml rm -f postgres
	docker volume rm herald_postgres_data || true
	docker-compose -f docker-compose.dev.yml up -d postgres
	@echo "⏳ Waiting for database..."
	sleep 10
	make db-migrate

# Run database migrations
db-migrate:
	@echo "📊 Running database migrations..."
	docker-compose -f docker-compose.dev.yml exec herald-api go run cmd/migrate/main.go

# Setup development environment
setup:
	@echo "⚙️ Setting up Herald.lol development environment..."
	cp .env.development .env
	mkdir -p backend frontend database monitoring/grafana/dashboards
	@echo "✅ Development environment ready!"
	@echo "Next steps:"
	@echo "1. Add your Riot API key to .env file"
	@echo "2. Run 'make dev' to start services"

# Quick development workflow
dev-quick: build dev
	@echo "🎮 Herald.lol development environment ready!"

# Production build (for testing)
build-prod:
	@echo "🚀 Building production images..."
	docker-compose -f docker-compose.prod.yml build

# Health check
health:
	@echo "🏥 Checking Herald.lol services health..."
	@curl -s http://localhost:8080/health || echo "API: ❌ DOWN"
	@curl -s http://localhost:3000 > /dev/null && echo "Frontend: ✅ UP" || echo "Frontend: ❌ DOWN"
	@curl -s http://localhost:9090/-/healthy > /dev/null && echo "Prometheus: ✅ UP" || echo "Prometheus: ❌ DOWN"
	@curl -s http://localhost:3001/api/health > /dev/null && echo "Grafana: ✅ UP" || echo "Grafana: ❌ DOWN"