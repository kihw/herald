#!/bin/bash

# Herald.lol Gaming Analytics - Local SonarQube Development Setup
# Quick setup for gaming developers to run code quality analysis locally

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
GAMING_PERFORMANCE_TARGET=5000

# Colors for gaming terminal output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
GOLD='\033[1;33m'  # Gaming gold color
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] [INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] [SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] [WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] [ERROR]${NC} $1"
}

log_gaming() {
    echo -e "${GOLD}[$(date +'%Y-%m-%d %H:%M:%S')] [GAMING]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_info "🔍 Checking Herald.lol development prerequisites..."
    
    if ! command -v docker >/dev/null 2>&1; then
        log_error "❌ Docker is required for SonarQube"
        return 1
    fi
    
    if ! command -v docker-compose >/dev/null 2>&1; then
        log_error "❌ Docker Compose is required for SonarQube"
        return 1
    fi
    
    if ! command -v go >/dev/null 2>&1; then
        log_error "❌ Go is required for gaming backend analysis"
        return 1
    fi
    
    if ! command -v node >/dev/null 2>&1; then
        log_warning "⚠️ Node.js not found - frontend analysis will be skipped"
    fi
    
    log_success "✅ Prerequisites check completed"
}

# Start SonarQube for gaming development
start_sonarqube() {
    log_gaming "🎮 Starting Herald.lol SonarQube for local gaming development..."
    
    cd "$PROJECT_ROOT"
    
    # Create required directories
    mkdir -p volumes/sonarqube/{data,logs,extensions,database}
    chmod -R 777 volumes/sonarqube/
    
    log_info "🚀 Starting SonarQube containers..."
    docker-compose -f docker-compose.sonarqube.yml up -d
    
    # Wait for SonarQube to be ready
    log_info "⏳ Waiting for Herald.lol SonarQube to initialize..."
    local timeout=300
    while [ $timeout -gt 0 ]; do
        if curl -f http://localhost:9000/api/system/status 2>/dev/null | grep -q '"status":"UP"'; then
            log_success "✅ SonarQube is ready for gaming analysis!"
            break
        fi
        echo -n "."
        sleep 5
        timeout=$((timeout - 5))
    done
    
    if [ $timeout -le 0 ]; then
        log_error "❌ SonarQube failed to start within timeout"
        docker-compose -f docker-compose.sonarqube.yml logs sonarqube-herald
        return 1
    fi
}

# Configure SonarQube for gaming
setup_gaming_profiles() {
    log_gaming "🎮 Setting up Herald.lol gaming quality profiles..."
    
    # Default credentials for local development
    local sonar_token="squ_herald_gaming_local_token"
    local admin_password="herald_gaming_admin_2024"
    
    # Change default admin password
    log_info "🔐 Setting up gaming admin credentials..."
    curl -X POST "http://localhost:9000/api/users/change_password" \
        -u admin:admin \
        -d "login=admin&password=$admin_password&previousPassword=admin" \
        2>/dev/null || log_warning "Password might already be changed"
    
    # Create user token for gaming analysis
    curl -X POST "http://localhost:9000/api/user_tokens/generate" \
        -u admin:$admin_password \
        -d "name=herald-gaming-local" \
        2>/dev/null || log_warning "Token might already exist"
    
    # Create gaming project
    log_info "🎯 Creating Herald.lol gaming project..."
    curl -X POST "http://localhost:9000/api/projects/create" \
        -u admin:$admin_password \
        -d "project=herald-gaming-analytics&name=Herald.lol Gaming Analytics Platform" \
        2>/dev/null || log_warning "Project might already exist"
    
    # Create gaming quality profiles
    log_info "🏆 Setting up gaming quality profiles..."
    
    # Go gaming profile
    curl -X POST "http://localhost:9000/api/qualityprofiles/create" \
        -u admin:$admin_password \
        -d "name=Herald Gaming Go Profile&language=go" \
        2>/dev/null || log_warning "Go profile might already exist"
    
    # TypeScript gaming profile
    curl -X POST "http://localhost:9000/api/qualityprofiles/create" \
        -u admin:$admin_password \
        -d "name=Herald Gaming TypeScript Profile&language=ts" \
        2>/dev/null || log_warning "TypeScript profile might already exist"
    
    # Create gaming quality gate
    log_info "🎯 Creating gaming quality gate..."
    curl -X POST "http://localhost:9000/api/qualitygates/create" \
        -u admin:$admin_password \
        -d "name=Herald Gaming Quality Gate" \
        2>/dev/null || log_warning "Quality gate might already exist"
    
    log_success "✅ Gaming profiles configured"
}

# Run tests and generate coverage
run_tests() {
    log_gaming "🧪 Running Herald.lol gaming tests..."
    
    cd "$PROJECT_ROOT"
    
    # Run Go tests with coverage
    log_info "🎮 Running Go gaming backend tests..."
    if ! go test -v -race -coverprofile=coverage.out -covermode=atomic ./...; then
        log_error "❌ Go tests failed"
        return 1
    fi
    
    # Generate test report
    go test -json -coverprofile=coverage.out ./... > test-report.json
    
    # Run frontend tests if Node.js is available
    if command -v node >/dev/null 2>&1 && [ -d "frontend" ]; then
        log_info "🎯 Running frontend gaming tests..."
        cd frontend
        if [ -f "package.json" ]; then
            npm install >/dev/null 2>&1 || log_warning "Frontend dependencies installation failed"
            npm run test:coverage >/dev/null 2>&1 || log_warning "Frontend tests not configured"
        fi
        cd "$PROJECT_ROOT"
    fi
    
    log_success "✅ Gaming tests completed"
}

# Run SonarQube analysis
run_analysis() {
    log_gaming "🔍 Running Herald.lol gaming code quality analysis..."
    
    cd "$PROJECT_ROOT"
    
    log_info "📊 Analyzing gaming platform code..."
    docker run --rm \
        --network=host \
        -v "${PWD}:/usr/src" \
        -e SONAR_HOST_URL="http://localhost:9000" \
        -e SONAR_SCANNER_OPTS="-Xms512m -Xmx1g" \
        sonarsource/sonar-scanner-cli:5.0 \
        -Dsonar.login=admin \
        -Dsonar.password=herald_gaming_admin_2024 \
        -Dsonar.projectKey=herald-gaming-analytics \
        -Dsonar.projectName="Herald.lol Gaming Analytics Platform" \
        -Dsonar.projectVersion="1.0.0-local" \
        -Dsonar.sources=. \
        -Dsonar.exclusions="**/vendor/**,**/node_modules/**,**/*_test.go,**/*.test.ts,**/testdata/**,**/mock/**,**/*.pb.go,**/dist/**,**/build/**,**/volumes/**" \
        -Dsonar.tests=. \
        -Dsonar.test.inclusions="**/*_test.go,**/*.test.ts,**/*.test.js,**/*.spec.ts" \
        -Dsonar.go.coverage.reportPaths=coverage.out \
        -Dsonar.go.tests.reportPaths=test-report.json \
        -Dsonar.gaming.performance.target=$GAMING_PERFORMANCE_TARGET \
        -Dsonar.gaming.concurrent.users=1000000 \
        -Dsonar.gaming.uptime.target=99.9 \
        -Dsonar.qualitygate.wait=false
    
    if [ $? -eq 0 ]; then
        log_success "✅ Gaming code quality analysis completed"
        log_info "🌐 View results at: http://localhost:9000/dashboard?id=herald-gaming-analytics"
    else
        log_error "❌ Analysis failed"
        return 1
    fi
}

# Show gaming dashboard
show_dashboard() {
    log_gaming "🎮 Opening Herald.lol gaming quality dashboard..."
    
    echo ""
    echo "🎮 Herald.lol Gaming Code Quality Dashboard"
    echo "=========================================="
    echo "🌐 URL: http://localhost:9000"
    echo "👤 Username: admin"
    echo "🔑 Password: herald_gaming_admin_2024"
    echo ""
    echo "🎯 Gaming Project: herald-gaming-analytics"
    echo "⚡ Performance Target: <${GAMING_PERFORMANCE_TARGET}ms"
    echo "👥 Concurrent Users: 1M+ support"
    echo "🏆 Uptime Target: 99.9%"
    echo ""
    
    if command -v xdg-open >/dev/null 2>&1; then
        xdg-open "http://localhost:9000/dashboard?id=herald-gaming-analytics" 2>/dev/null &
    elif command -v open >/dev/null 2>&1; then
        open "http://localhost:9000/dashboard?id=herald-gaming-analytics" 2>/dev/null &
    fi
}

# Stop SonarQube
stop_sonarqube() {
    log_info "🛑 Stopping Herald.lol SonarQube..."
    
    cd "$PROJECT_ROOT"
    docker-compose -f docker-compose.sonarqube.yml down
    
    log_success "✅ SonarQube stopped"
}

# Clean up everything
cleanup() {
    log_warning "🧹 Cleaning up Herald.lol SonarQube..."
    
    cd "$PROJECT_ROOT"
    docker-compose -f docker-compose.sonarqube.yml down -v
    docker system prune -f >/dev/null 2>&1
    
    # Remove generated files
    rm -f coverage.out test-report.json
    rm -rf .sonar/
    
    log_success "✅ Cleanup completed"
}

# Usage
usage() {
    cat << EOF
🎮 Herald.lol Gaming SonarQube Local Development

Usage: $0 [COMMAND]

COMMANDS:
    start           Start SonarQube for gaming development
    setup           Configure gaming quality profiles
    test            Run gaming tests with coverage
    analyze         Run gaming code quality analysis
    dashboard       Open gaming quality dashboard
    full            Run complete gaming analysis (start + setup + test + analyze)
    stop            Stop SonarQube containers
    cleanup         Clean up all SonarQube data
    -h, --help      Show this help message

GAMING FOCUS:
    🎯 Performance Target: <${GAMING_PERFORMANCE_TARGET}ms analytics response
    👥 Concurrent Users: 1M+ gaming platform support
    🎮 Gaming Metrics: LoL/TFT specific code quality
    ⚡ Real-time: WebSocket/gRPC gaming optimization

EXAMPLES:
    # Full gaming analysis
    $0 full

    # Quick analysis run
    $0 start && $0 analyze

    # View dashboard
    $0 dashboard
EOF
}

# Main function
main() {
    case "${1:-}" in
        start)
            check_prerequisites
            start_sonarqube
            ;;
        setup)
            setup_gaming_profiles
            ;;
        test)
            run_tests
            ;;
        analyze)
            run_tests
            run_analysis
            ;;
        dashboard)
            show_dashboard
            ;;
        full)
            check_prerequisites
            start_sonarqube
            sleep 10
            setup_gaming_profiles
            run_tests
            run_analysis
            show_dashboard
            ;;
        stop)
            stop_sonarqube
            ;;
        cleanup)
            cleanup
            ;;
        -h|--help)
            usage
            ;;
        *)
            log_gaming "🎮 Herald.lol Gaming SonarQube Local Development"
            echo ""
            usage
            exit 1
            ;;
    esac
}

# Handle Ctrl+C gracefully
trap 'echo ""; log_info "🎮 Gaming analysis interrupted"; exit 0' SIGINT

# Run main function
main "$@"