#!/bin/bash

# Herald.lol Gaming Analytics Platform - VPS Production Deployment Script
# Automated deployment on VPS with Docker Compose

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Default values
ENVIRONMENT="production"
DOMAIN="herald.lol"
SKIP_SSL=false
SKIP_BUILD=false
QUICK_DEPLOY=false

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

show_banner() {
    echo -e "${BLUE}"
    cat << "EOF"
    ██   ██ ███████ ██████   █████  ██      ██████      ██       ██████  ██      
    ██   ██ ██      ██   ██ ██   ██ ██      ██   ██     ██      ██    ██ ██      
    ███████ █████   ██████  ███████ ██      ██   ██     ██      ██    ██ ██      
    ██   ██ ██      ██   ██ ██   ██ ██      ██   ██     ██      ██    ██ ██      
    ██   ██ ███████ ██   ██ ██   ██ ███████ ██████  ██  ███████  ██████  ███████ 
    
    🎮 Gaming Analytics Platform - VPS Production Deployment 🎮
EOF
    echo -e "${NC}"
}

show_help() {
    cat << EOF
Herald.lol VPS Production Deployment Script

Usage: $0 [OPTIONS]

OPTIONS:
    -d, --domain DOMAIN        Domain name (default: herald.lol)
    -s, --skip-ssl            Skip SSL certificate generation
    -b, --skip-build          Skip Docker image builds
    -q, --quick               Quick deploy (skip builds and SSL)
    -h, --help                Show this help message

EXAMPLES:
    # Full production deployment
    $0 -d herald.lol

    # Quick deploy for updates
    $0 --quick

    # Deploy without SSL (for testing)
    $0 --skip-ssl
EOF
}

check_prerequisites() {
    log_info "Checking prerequisites for VPS deployment..."
    
    # Check if Docker is installed and running
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        log_info "Install with: curl -fsSL https://get.docker.com | sh"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        log_error "Docker daemon is not running"
        log_info "Start with: sudo systemctl start docker"
        exit 1
    fi
    
    # Check if Docker Compose is available
    if ! docker compose version &> /dev/null; then
        log_error "Docker Compose is not available"
        log_info "Install Docker Compose v2"
        exit 1
    fi
    
    # Check if running as root or in docker group
    if [ "$EUID" -ne 0 ] && ! groups | grep -q docker; then
        log_error "User must be root or in docker group"
        log_info "Add user to docker group: sudo usermod -aG docker $USER"
        exit 1
    fi
    
    # Check available disk space (at least 10GB)
    local available_space=$(df / | awk 'NR==2 {print $4}')
    if [ "${available_space}" -lt 10485760 ]; then # 10GB in KB
        log_warning "Less than 10GB disk space available"
    fi
    
    # Check available memory (at least 4GB)
    local available_memory=$(free -m | awk 'NR==2{print $7}')
    if [ "${available_memory}" -lt 4000 ]; then
        log_warning "Less than 4GB memory available"
    fi
    
    log_success "All prerequisites satisfied"
}

setup_environment() {
    log_info "Setting up production environment..."
    
    cd "${PROJECT_ROOT}"
    
    # Create .env.production if it doesn't exist
    if [ ! -f ".env.production" ]; then
        log_info "Creating production environment file..."
        cat << EOF > .env.production
# Herald.lol Production Environment
POSTGRES_PASSWORD=$(openssl rand -base64 32)
RIOT_API_KEY=RGAPI-your-riot-api-key-here
JWT_SECRET=$(openssl rand -base64 64)
GRAFANA_PASSWORD=$(openssl rand -base64 16)
DOMAIN=${DOMAIN}
EOF
        log_warning "Please update .env.production with your actual Riot API key"
    fi
    
    # Create necessary directories
    mkdir -p logs nginx/ssl monitoring/prometheus monitoring/grafana
    
    log_success "Environment setup complete"
}

generate_ssl_certificates() {
    if [ "${SKIP_SSL}" = true ]; then
        log_info "Skipping SSL certificate generation"
        return
    fi
    
    log_info "Setting up SSL certificates with Let's Encrypt..."
    
    # Install certbot if not present
    if ! command -v certbot &> /dev/null; then
        log_info "Installing Certbot..."
        if command -v apt-get &> /dev/null; then
            sudo apt-get update
            sudo apt-get install -y certbot
        elif command -v yum &> /dev/null; then
            sudo yum install -y certbot
        else
            log_error "Cannot install certbot automatically"
            exit 1
        fi
    fi
    
    # Generate certificates for all domains
    local domains=(
        "${DOMAIN}"
        "www.${DOMAIN}"
        "api.${DOMAIN}"
        "ws.${DOMAIN}"
        "grpc.${DOMAIN}"
        "monitoring.${DOMAIN}"
    )
    
    local domain_args=""
    for domain in "${domains[@]}"; do
        domain_args="${domain_args} -d ${domain}"
    done
    
    # Request certificate
    sudo certbot certonly --standalone ${domain_args} \
        --agree-tos \
        --register-unsafely-without-email \
        --non-interactive
    
    log_success "SSL certificates generated"
}

build_images() {
    if [ "${SKIP_BUILD}" = true ]; then
        log_info "Skipping Docker image builds"
        return
    fi
    
    log_info "Building Herald gaming Docker images..."
    
    cd "${PROJECT_ROOT}"
    
    # Build backend image
    log_info "Building backend image..."
    docker build -f backend/Dockerfile.dev -t herald-backend:production ./backend
    
    # Build frontend image
    log_info "Building frontend image..."
    docker build -f frontend/Dockerfile.dev -t herald-frontend:production ./frontend
    
    log_success "Docker images built successfully"
}

deploy_services() {
    log_info "Deploying Herald gaming services..."
    
    cd "${PROJECT_ROOT}"
    
    # Stop existing services
    log_info "Stopping existing services..."
    docker compose -f docker-compose.production.yml down || true
    
    # Clean up old containers and images
    docker system prune -f
    
    # Deploy production stack
    log_info "Starting Herald gaming production stack..."
    docker compose -f docker-compose.production.yml --env-file .env.production up -d
    
    # Wait for services to be healthy
    log_info "Waiting for services to be ready..."
    sleep 30
    
    # Check service health
    local services=("herald-postgres-prod" "herald-redis-prod" "herald-api-prod" "herald-frontend-prod" "herald-nginx-prod")
    for service in "${services[@]}"; do
        if docker ps | grep -q "${service}"; then
            log_success "✅ ${service} is running"
        else
            log_error "❌ ${service} failed to start"
            docker logs "${service}" --tail 20
        fi
    done
    
    log_success "Herald gaming services deployed"
}

setup_monitoring() {
    log_info "Setting up monitoring stack..."
    
    # Create Prometheus configuration
    cat << EOF > monitoring/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "gaming_alerts.yml"

scrape_configs:
  - job_name: 'herald-api'
    static_configs:
      - targets: ['herald-api:8080']
    metrics_path: '/metrics'
    
  - job_name: 'herald-postgres'
    static_configs:
      - targets: ['herald-postgres-exporter:9187']
    
  - job_name: 'herald-redis'
    static_configs:
      - targets: ['herald-redis-exporter:9121']
    
  - job_name: 'node-exporter'
    static_configs:
      - targets: ['herald-node-exporter:9100']

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093
EOF

    log_success "Monitoring configuration complete"
}

run_health_checks() {
    log_info "Running comprehensive health checks..."
    
    local checks_passed=0
    local total_checks=5
    
    # Check database connection
    if docker exec herald-postgres-prod pg_isready -U herald_admin > /dev/null 2>&1; then
        log_success "✅ PostgreSQL connection"
        ((checks_passed++))
    else
        log_error "❌ PostgreSQL connection failed"
    fi
    
    # Check Redis connection
    if docker exec herald-redis-prod redis-cli ping | grep -q PONG; then
        log_success "✅ Redis connection"
        ((checks_passed++))
    else
        log_error "❌ Redis connection failed"
    fi
    
    # Check API health
    sleep 10  # Wait for API to be ready
    if curl -f -s http://localhost:8080/health > /dev/null 2>&1; then
        log_success "✅ API health check"
        ((checks_passed++))
    else
        log_error "❌ API health check failed"
    fi
    
    # Check frontend
    if curl -f -s http://localhost:3000 > /dev/null 2>&1; then
        log_success "✅ Frontend accessibility"
        ((checks_passed++))
    else
        log_error "❌ Frontend accessibility failed"
    fi
    
    # Check SSL certificates
    if [ "${SKIP_SSL}" = false ] && [ -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ]; then
        log_success "✅ SSL certificates present"
        ((checks_passed++))
    elif [ "${SKIP_SSL}" = true ]; then
        log_info "⏭️ SSL checks skipped"
        ((checks_passed++))
    else
        log_error "❌ SSL certificates missing"
    fi
    
    # Summary
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}Health Check Summary${NC}"
    echo -e "${BLUE}================================${NC}"
    echo "Passed: ${checks_passed}/${total_checks}"
    
    if [ "${checks_passed}" -eq "${total_checks}" ]; then
        log_success "🎉 All health checks passed!"
        return 0
    else
        log_error "❌ Some health checks failed"
        return 1
    fi
}

show_deployment_info() {
    echo -e "${GREEN}================================${NC}"
    echo -e "${GREEN}Herald.lol Production Deployment Complete!${NC}"
    echo -e "${GREEN}================================${NC}"
    echo
    echo "🎮 Gaming Analytics Platform URLs:"
    
    if [ "${SKIP_SSL}" = false ]; then
        echo "  🌐 Main Site: https://${DOMAIN}"
        echo "  📊 API: https://api.${DOMAIN}"
        echo "  🔗 WebSocket: wss://ws.${DOMAIN}"
        echo "  📈 Monitoring: https://monitoring.${DOMAIN}"
    else
        echo "  🌐 Main Site: http://${DOMAIN}"
        echo "  📊 API: http://api.${DOMAIN}:8080"
        echo "  📈 Monitoring: http://monitoring.${DOMAIN}:3001"
    fi
    
    echo
    echo "🔧 Management Commands:"
    echo "  📋 View logs: docker compose -f docker-compose.production.yml logs -f"
    echo "  🔄 Restart: docker compose -f docker-compose.production.yml restart"
    echo "  ⛔ Stop: docker compose -f docker-compose.production.yml down"
    echo "  📊 Status: docker compose -f docker-compose.production.yml ps"
    echo
    echo "🎯 Performance Targets:"
    echo "  ⚡ API Response: <500ms"
    echo "  🎮 Gaming Analytics: <5s"
    echo "  📈 Uptime: 99.9%"
    echo
    log_success "Herald.lol gaming analytics platform is now live! 🎮🚀"
}

main() {
    show_banner
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--domain)
                DOMAIN="$2"
                shift 2
                ;;
            -s|--skip-ssl)
                SKIP_SSL=true
                shift
                ;;
            -b|--skip-build)
                SKIP_BUILD=true
                shift
                ;;
            -q|--quick)
                QUICK_DEPLOY=true
                SKIP_BUILD=true
                SKIP_SSL=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    log_info "Starting Herald.lol VPS production deployment"
    log_info "Domain: ${DOMAIN}"
    log_info "Skip SSL: ${SKIP_SSL}"
    log_info "Skip Build: ${SKIP_BUILD}"
    
    # Execute deployment steps
    check_prerequisites
    setup_environment
    
    if [ "${SKIP_SSL}" = false ]; then
        generate_ssl_certificates
    fi
    
    if [ "${SKIP_BUILD}" = false ]; then
        build_images
    fi
    
    setup_monitoring
    deploy_services
    
    # Run health checks
    if run_health_checks; then
        show_deployment_info
        log_success "🎮 Herald.lol production deployment completed successfully!"
        exit 0
    else
        log_error "❌ Deployment completed with issues"
        log_info "Check logs with: docker compose -f docker-compose.production.yml logs"
        exit 1
    fi
}

# Trap to cleanup on exit
trap 'log_error "Deployment interrupted"' INT TERM

# Run main function
main "$@"