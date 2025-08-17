# Makefile for LoL Match Exporter

.PHONY: help build test clean run docker-build docker-run deploy deploy-staging rollback health-check lint security-scan

# Variables
BINARY_NAME=lol-match-exporter
DOCKER_IMAGE=lol-match-exporter-fullstack
VERSION?=latest
BUILD_NUMBER?=$(shell date +%Y%m%d-%H%M%S)
GO_FILES=$(shell find . -name "*.go" -type f)
DEPLOY_HOST=51.178.17.78
DEPLOY_USER=debian

# Couleurs pour l'affichage
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
CYAN=\033[0;36m
NC=\033[0m

help: ## Show this help message
	@echo "$(GREEN)LoL Match Exporter - Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-25s$(NC) %s\n", $$1, $$2}'

# =============================================================================
# BUILD & TEST
# =============================================================================

build: ## Build the Go binary
	@echo "$(YELLOW)Building $(BINARY_NAME)...$(NC)"
	go mod tidy
	go build -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o $(BINARY_NAME) .
	@echo "$(GREEN)Build completed$(NC)"

test: ## Run tests with coverage
	@echo "$(YELLOW)Running tests...$(NC)"
	go test -v -race -coverprofile=coverage.out ./...
	@if command -v go >/dev/null 2>&1; then \
		go tool cover -html=coverage.out -o coverage.html; \
		echo "$(GREEN)Tests completed - coverage report: coverage.html$(NC)"; \
	fi

test-integration: ## Run integration tests
	@echo "$(YELLOW)Running integration tests...$(NC)"
	go test -v -tags=integration ./tests/...
	@echo "$(GREEN)Integration tests completed$(NC)"

lint: ## Run linting tools
	@echo "$(YELLOW)Running linters...$(NC)"
	go vet ./...
	go fmt ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(RED)golangci-lint not installed, skipping$(NC)"; \
	fi
	@echo "$(GREEN)Linting completed$(NC)"

security-scan: ## Run security scans
	@echo "$(YELLOW)Running security scans...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "$(RED)gosec not installed, skipping Go security scan$(NC)"; \
	fi
	@if command -v trivy >/dev/null 2>&1 && docker images $(DOCKER_IMAGE):latest >/dev/null 2>&1; then \
		trivy image $(DOCKER_IMAGE):latest; \
	else \
		echo "$(RED)trivy not installed or image not found, skipping Docker security scan$(NC)"; \
	fi
	@echo "$(GREEN)Security scans completed$(NC)"

clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning...$(NC)"
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	go clean
	@if command -v docker >/dev/null 2>&1; then \
		docker system prune -f 2>/dev/null || true; \
	fi
	@echo "$(GREEN)Cleanup completed$(NC)"

# =============================================================================
# LOCAL DEVELOPMENT
# =============================================================================

run: build ## Run the application locally
	@echo "$(YELLOW)Running $(BINARY_NAME)...$(NC)"
	./$(BINARY_NAME)

dev: ## Run in development mode with hot reload
	@echo "$(YELLOW)Starting development server...$(NC)"
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "$(RED)air not installed. Install with: go install github.com/cosmtrek/air@latest$(NC)"; \
		$(MAKE) run; \
	fi

frontend-dev: ## Start frontend development server
	@echo "$(YELLOW)Starting frontend development server...$(NC)"
	@if [ -d "web" ] && [ -f "web/package.json" ]; then \
		cd web && npm install && npm run dev; \
	else \
		echo "$(RED)Frontend not found or not configured$(NC)"; \
	fi

# =============================================================================
# DOCKER
# =============================================================================

docker-build: ## Build Docker image
	@echo "$(YELLOW)Building Docker image $(DOCKER_IMAGE):$(VERSION)...$(NC)"
	docker build -f Dockerfile.simple-fullstack -t $(DOCKER_IMAGE):$(VERSION) .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest
	@echo "$(GREEN)Docker image built$(NC)"

docker-build-dev: ## Build development Docker image
	@echo "$(YELLOW)Building development Docker image...$(NC)"
	docker build -f Dockerfile.debug -t $(DOCKER_IMAGE):dev .
	@echo "$(GREEN)Development Docker image built$(NC)"

docker-run: docker-build ## Run Docker container locally
	@echo "$(YELLOW)Running Docker container locally...$(NC)"
	docker run --rm -p 80:80 -p 443:443 --name lol-exporter-local $(DOCKER_IMAGE):$(VERSION)

docker-compose-up: ## Start with docker-compose
	@echo "$(YELLOW)Starting with docker-compose...$(NC)"
	docker-compose -f docker-compose.complete.yml up -d
	@echo "$(GREEN)Application started at https://localhost$(NC)"

docker-compose-down: ## Stop docker-compose services
	@echo "$(YELLOW)Stopping docker-compose services...$(NC)"
	docker-compose -f docker-compose.complete.yml down
	@echo "$(GREEN)Services stopped$(NC)"

docker-logs: ## Show Docker container logs
	@if command -v docker >/dev/null 2>&1; then \
		docker logs -f lol-fullstack-app 2>/dev/null || docker-compose -f docker-compose.complete.yml logs -f; \
	else \
		echo "$(RED)Docker not available$(NC)"; \
	fi

# =============================================================================
# DEPLOYMENT
# =============================================================================

deploy-staging: docker-build ## Deploy to staging
	@echo "$(YELLOW)Deploying to staging environment...$(NC)"
	@if [ -f "./scripts/deploy-automated.sh" ]; then \
		chmod +x ./scripts/deploy-automated.sh; \
		./scripts/deploy-automated.sh staging $(BUILD_NUMBER); \
	else \
		echo "$(RED)Deployment script not found$(NC)"; \
	fi
	@echo "$(GREEN)Staging deployment completed$(NC)"

deploy: docker-build ## Deploy to production
	@echo "$(YELLOW)Deploying to production...$(NC)"
	@echo "$(RED)WARNING: This will deploy to production. Continue? [y/N]$(NC)" && read -r ans && [ "$${ans:-N}" = "y" ]
	@if [ -f "./scripts/deploy-automated.sh" ]; then \
		chmod +x ./scripts/deploy-automated.sh; \
		./scripts/deploy-automated.sh production $(BUILD_NUMBER); \
	else \
		echo "$(RED)Deployment script not found$(NC)"; \
	fi
	@echo "$(GREEN)Production deployment completed$(NC)"

deploy-quick: ## Quick deploy (skips tests and scans)
	@echo "$(YELLOW)Quick deployment to production...$(NC)"
	@echo "$(RED)WARNING: Skipping tests and security scans. Continue? [y/N]$(NC)" && read -r ans && [ "$${ans:-N}" = "y" ]
	$(MAKE) docker-build
	@if [ -f "./scripts/deploy-automated.sh" ]; then \
		chmod +x ./scripts/deploy-automated.sh; \
		./scripts/deploy-automated.sh production $(BUILD_NUMBER); \
	fi
	@echo "$(GREEN)Quick deployment completed$(NC)"

rollback: ## Rollback to previous version
	@echo "$(YELLOW)Rolling back to previous version...$(NC)"
	@if [ -f "./scripts/deploy-automated.sh" ]; then \
		chmod +x ./scripts/deploy-automated.sh; \
		./scripts/deploy-automated.sh rollback; \
	else \
		echo "$(RED)Deployment script not found$(NC)"; \
	fi
	@echo "$(GREEN)Rollback completed$(NC)"

# =============================================================================
# MONITORING & HEALTH
# =============================================================================

health-check: ## Check application health
	@echo "$(YELLOW)Checking application health...$(NC)"
	@if [ -f "./scripts/deploy-automated.sh" ]; then \
		chmod +x ./scripts/deploy-automated.sh; \
		./scripts/deploy-automated.sh health; \
	else \
		echo "Checking directly..."; \
		curl -s https://herald.lol/api/health || echo "$(RED)Health check failed$(NC)"; \
	fi
	@echo "$(GREEN)Health check completed$(NC)"

logs: ## Show remote application logs
	@echo "$(YELLOW)Fetching remote logs...$(NC)"
	ssh $(DEPLOY_USER)@$(DEPLOY_HOST) "docker logs lol-fullstack-app --tail 100 -f"

status: ## Show deployment status
	@echo "$(YELLOW)Checking deployment status...$(NC)"
	@ssh $(DEPLOY_USER)@$(DEPLOY_HOST) "docker ps | grep lol-fullstack || echo 'Container not running'"
	@echo ""
	@if curl -s https://herald.lol/api/health >/dev/null 2>&1; then \
		echo "$(GREEN)API Health: OK$(NC)"; \
	else \
		echo "$(RED)API Health: FAILED$(NC)"; \
	fi

monitor: ## Start monitoring dashboard
	@echo "$(YELLOW)Starting monitoring...$(NC)"
	@echo "Opening monitoring URLs:"
	@echo "- Application: https://herald.lol"
	@echo "- Health API: https://herald.lol/api/health"
	@echo ""
	@echo "Press Ctrl+C to stop monitoring"
	@while true; do \
		if curl -s https://herald.lol/api/health >/dev/null 2>&1; then \
			echo "$$(date): $(GREEN)OK$(NC)"; \
		else \
			echo "$$(date): $(RED)FAIL$(NC)"; \
		fi; \
		sleep 30; \
	done

# =============================================================================
# CI/CD PIPELINE
# =============================================================================

ci-test: ## Run CI tests
	@echo "$(YELLOW)Running CI test suite...$(NC)"
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) security-scan
	@echo "$(GREEN)CI tests completed$(NC)"

ci-build: ## Build for CI
	@echo "$(YELLOW)Building for CI...$(NC)"
	$(MAKE) build
	$(MAKE) docker-build
	@echo "$(GREEN)CI build completed$(NC)"

ci-deploy: ## Full CI/CD pipeline
	@echo "$(YELLOW)Running full CI/CD pipeline...$(NC)"
	$(MAKE) ci-test
	$(MAKE) ci-build
	$(MAKE) deploy
	$(MAKE) health-check
	@echo "$(GREEN)CI/CD pipeline completed$(NC)"

# =============================================================================
# UTILITIES
# =============================================================================

ssh: ## SSH into production server
	ssh $(DEPLOY_USER)@$(DEPLOY_HOST)

backup: ## Backup production data
	@echo "$(YELLOW)Creating backup...$(NC)"
	ssh $(DEPLOY_USER)@$(DEPLOY_HOST) "docker exec lol-fullstack-app tar czf /tmp/backup-$$(date +%Y%m%d-%H%M%S).tar.gz /app/data 2>/dev/null || echo 'No data to backup'"
	@echo "$(GREEN)Backup completed$(NC)"

install-tools: ## Install development tools
	@echo "$(YELLOW)Installing development tools...$(NC)"
	go install github.com/cosmtrek/air@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2; \
	fi
	@echo "$(GREEN)Development tools installed$(NC)"

version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Build Number: $(BUILD_NUMBER)"
	@echo "Go Version: $$(go version)"
	@if command -v docker >/dev/null 2>&1; then \
		echo "Docker Version: $$(docker --version)"; \
	else \
		echo "Docker Version: Not installed"; \
	fi

# =============================================================================
# WINDOWS COMPATIBILITY
# =============================================================================

deploy-win: docker-build ## Deploy from Windows
	@echo "$(YELLOW)Deploying from Windows...$(NC)"
	@powershell -Command "Write-Host 'WARNING: This will deploy to production. Continue? [y/N]' -ForegroundColor Red -NoNewline; $$ans = Read-Host; if ($$ans -eq 'y') { exit 0 } else { exit 1 }"
	@if exist "scripts\\deploy-automated.sh" ( \
		bash scripts/deploy-automated.sh production $(BUILD_NUMBER) \
	) else ( \
		echo "$(RED)Deployment script not found$(NC)" \
	)

health-check-win: ## Windows health check
	@powershell -Command "try { Invoke-RestMethod 'https://herald.lol/api/health' | Out-Null; Write-Host 'API Health: OK' -ForegroundColor Green } catch { Write-Host 'API Health: FAILED' -ForegroundColor Red }"

# =============================================================================
# DEFAULT TARGET
# =============================================================================

.DEFAULT_GOAL := help
