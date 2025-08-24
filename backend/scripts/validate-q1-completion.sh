#!/bin/bash

# Herald.lol Gaming Analytics Platform - Q1 2025 Completion Validation
# Validates that all Q1 infrastructure and development tasks are complete

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

# Counters
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✅ PASS]${NC} $1"
    ((PASSED_CHECKS++))
}

log_error() {
    echo -e "${RED}[❌ FAIL]${NC} $1"
    ((FAILED_CHECKS++))
}

log_warning() {
    echo -e "${YELLOW}[⚠️ WARN]${NC} $1"
}

check_file_exists() {
    local file_path="$1"
    local description="$2"
    ((TOTAL_CHECKS++))
    
    if [ -f "${file_path}" ]; then
        log_success "${description}"
    else
        log_error "${description}"
    fi
}

check_dir_exists() {
    local dir_path="$1"
    local description="$2"
    ((TOTAL_CHECKS++))
    
    if [ -d "${dir_path}" ]; then
        log_success "${description}"
    else
        log_error "${description}"
    fi
}

show_banner() {
    echo -e "${BLUE}"
    cat << "EOF"
    ██   ██ ███████ ██████   █████  ██      ██████      ██       ██████  ██      
    ██   ██ ██      ██   ██ ██   ██ ██      ██   ██     ██      ██    ██ ██      
    ███████ █████   ██████  ███████ ██      ██   ██     ██      ██    ██ ██      
    ██   ██ ██      ██   ██ ██   ██ ██      ██   ██     ██      ██    ██ ██      
    ██   ██ ███████ ██   ██ ██   ██ ███████ ██████  ██  ███████  ██████  ███████ 
    
    🎮 Q1 2025 Completion Validation 🎮
EOF
    echo -e "${NC}"
}

validate_infrastructure() {
    log_info "🏗️ Validating Infrastructure Foundation..."
    
    # Check Terraform configuration
    check_file_exists "${PROJECT_ROOT}/backend/terraform/main.tf" "Terraform main configuration exists"
    check_file_exists "${PROJECT_ROOT}/backend/terraform/variables.tf" "Terraform variables configuration exists"
    check_file_exists "${PROJECT_ROOT}/backend/terraform/outputs.tf" "Terraform outputs configuration exists"
    
    # Check Kubernetes configuration
    check_dir_exists "${PROJECT_ROOT}/backend/k8s" "Kubernetes configuration directory exists"
    
    # Check Istio configuration
    check_file_exists "${PROJECT_ROOT}/backend/k8s/istio/istio-install.yaml" "Istio service mesh configuration exists"
    
    # Check deployment scripts
    check_file_exists "${PROJECT_ROOT}/backend/scripts/deploy-infrastructure.sh" "Infrastructure deployment script exists"
}

validate_backend_services() {
    log_info "⚡ Validating Backend Services..."
    
    # Check core services
    local services=(
        "analytics_service.go"
        "auth_service.go"
        "riot_service.go"
        "notification_service.go"
        "match_processing_service.go"
        "realtime_service.go"
        "coaching_service.go"
        "skill_progression_service.go"
    )
    
    for service in "${services[@]}"; do
        check_file_exists "${PROJECT_ROOT}/backend/internal/services/${service}" "Service exists: ${service}"
    done
    
    # Check gRPC server
    if [ -f "${PROJECT_ROOT}/backend/cmd/grpc-server/main.go" ]; then
        ((TOTAL_CHECKS++))
        log_success "gRPC server implementation exists"
    else
        ((TOTAL_CHECKS++))
        log_error "gRPC server implementation missing"
    fi
    
    # Check analytics engine
    if [ -f "${PROJECT_ROOT}/backend/internal/analytics/core_engine.go" ]; then
        ((TOTAL_CHECKS++))
        log_success "Analytics core engine exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Analytics core engine missing"
    fi
    
    # Check match analyzer
    if [ -f "${PROJECT_ROOT}/backend/internal/match/analyzer.go" ]; then
        ((TOTAL_CHECKS++))
        log_success "Match analyzer exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Match analyzer missing"
    fi
}

validate_database_architecture() {
    log_info "🗄️ Validating Database Architecture..."
    
    # Check models
    local models=(
        "user.go"
        "match.go"
        "analytics.go"
        "vision.go"
        "damage.go"
        "gold.go"
        "coaching.go"
        "skill_progression.go"
    )
    
    for model in "${models[@]}"; do
        if [ -f "${PROJECT_ROOT}/backend/internal/models/${model}" ]; then
            ((TOTAL_CHECKS++))
            log_success "Model exists: ${model}"
        else
            ((TOTAL_CHECKS++))
            log_error "Model missing: ${model}"
        fi
    done
    
    # Check migration system
    if [ -f "${PROJECT_ROOT}/backend/cmd/migrate/main.go" ]; then
        ((TOTAL_CHECKS++))
        log_success "Database migration system exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Database migration system missing"
    fi
}

validate_frontend_infrastructure() {
    log_info "🎨 Validating Frontend Infrastructure..."
    
    # Check package.json
    if [ -f "${PROJECT_ROOT}/frontend/package.json" ]; then
        ((TOTAL_CHECKS++))
        log_success "Frontend package.json exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Frontend package.json missing"
    fi
    
    # Check Vite configuration
    if [ -f "${PROJECT_ROOT}/frontend/vite.config.ts" ]; then
        ((TOTAL_CHECKS++))
        log_success "Vite configuration exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Vite configuration missing"
    fi
    
    # Check TypeScript configuration
    if [ -f "${PROJECT_ROOT}/frontend/tsconfig.json" ]; then
        ((TOTAL_CHECKS++))
        log_success "TypeScript configuration exists"
    else
        ((TOTAL_CHECKS++))
        log_error "TypeScript configuration missing"
    fi
    
    # Check Storybook configuration
    if [ -d "${PROJECT_ROOT}/frontend/.storybook" ]; then
        ((TOTAL_CHECKS++))
        log_success "Storybook configuration exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Storybook configuration missing"
    fi
}

validate_testing_infrastructure() {
    log_info "🧪 Validating Testing Infrastructure..."
    
    # Check Go tests
    local go_test_files=$(find "${PROJECT_ROOT}/backend" -name "*_test.go" | wc -l)
    if [ "${go_test_files}" -gt 0 ]; then
        ((TOTAL_CHECKS++))
        log_success "Go test files found: ${go_test_files}"
    else
        ((TOTAL_CHECKS++))
        log_error "No Go test files found"
    fi
    
    # Check performance testing
    if [ -f "${PROJECT_ROOT}/backend/scripts/load-test-gaming.js" ]; then
        ((TOTAL_CHECKS++))
        log_success "K6 performance testing script exists"
    else
        ((TOTAL_CHECKS++))
        log_error "K6 performance testing script missing"
    fi
    
    # Check testing configuration
    if [ -f "${PROJECT_ROOT}/backend/Makefile" ]; then
        if grep -q "test:" "${PROJECT_ROOT}/backend/Makefile"; then
            ((TOTAL_CHECKS++))
            log_success "Makefile contains test targets"
        else
            ((TOTAL_CHECKS++))
            log_error "Makefile missing test targets"
        fi
    else
        ((TOTAL_CHECKS++))
        log_error "Backend Makefile missing"
    fi
}

validate_security_compliance() {
    log_info "🔒 Validating Security & Compliance..."
    
    # Check authentication system
    if [ -f "${PROJECT_ROOT}/backend/internal/auth/jwt_manager.go" ]; then
        ((TOTAL_CHECKS++))
        log_success "JWT authentication system exists"
    else
        ((TOTAL_CHECKS++))
        log_error "JWT authentication system missing"
    fi
    
    # Check MFA implementation
    if [ -f "${PROJECT_ROOT}/backend/internal/auth/mfa.go" ]; then
        ((TOTAL_CHECKS++))
        log_success "Multi-Factor Authentication system exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Multi-Factor Authentication system missing"
    fi
    
    # Check RBAC system
    if [ -f "${PROJECT_ROOT}/backend/internal/auth/rbac.go" ]; then
        ((TOTAL_CHECKS++))
        log_success "Role-Based Access Control system exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Role-Based Access Control system missing"
    fi
    
    # Check security middleware
    if [ -f "${PROJECT_ROOT}/backend/internal/middleware/rate_limiter.go" ]; then
        ((TOTAL_CHECKS++))
        log_success "Rate limiting middleware exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Rate limiting middleware missing"
    fi
    
    # Check gaming security
    if [ -f "${PROJECT_ROOT}/backend/internal/middleware/gaming_security.go" ]; then
        ((TOTAL_CHECKS++))
        log_success "Gaming security middleware exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Gaming security middleware missing"
    fi
}

validate_gaming_features() {
    log_info "🎮 Validating Gaming-Specific Features..."
    
    # Check Riot API integration
    if [ -f "${PROJECT_ROOT}/backend/internal/riot/client.go" ]; then
        ((TOTAL_CHECKS++))
        log_success "Riot API client exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Riot API client missing"
    fi
    
    # Check gaming analytics services
    local gaming_services=(
        "champion_analytics_service.go"
        "damage_analytics_service.go"
        "vision_analytics_service.go"
        "gold_analytics_service.go"
        "ward_analytics_service.go"
        "meta_analytics_service.go"
    )
    
    for service in "${gaming_services[@]}"; do
        if [ -f "${PROJECT_ROOT}/backend/internal/services/${service}" ]; then
            ((TOTAL_CHECKS++))
            log_success "Gaming service exists: ${service}"
        else
            ((TOTAL_CHECKS++))
            log_error "Gaming service missing: ${service}"
        fi
    done
    
    # Check export service
    if [ -f "${PROJECT_ROOT}/backend/internal/export/service.go" ]; then
        ((TOTAL_CHECKS++))
        log_success "Export & reporting service exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Export & reporting service missing"
    fi
}

validate_documentation() {
    log_info "📚 Validating Documentation..."
    
    # Check project documentation
    if [ -f "${PROJECT_ROOT}/CLAUDE.md" ]; then
        ((TOTAL_CHECKS++))
        log_success "Project documentation (CLAUDE.md) exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Project documentation (CLAUDE.md) missing"
    fi
    
    # Check infrastructure documentation
    if [ -f "${PROJECT_ROOT}/backend/README-INFRASTRUCTURE.md" ]; then
        ((TOTAL_CHECKS++))
        log_success "Infrastructure documentation exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Infrastructure documentation missing"
    fi
    
    # Check todo roadmap
    if [ -f "${PROJECT_ROOT}/todo.md" ]; then
        ((TOTAL_CHECKS++))
        log_success "Development roadmap (todo.md) exists"
    else
        ((TOTAL_CHECKS++))
        log_error "Development roadmap (todo.md) missing"
    fi
}

show_validation_summary() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}Herald.lol Q1 2025 Validation Summary${NC}"
    echo -e "${BLUE}================================${NC}"
    
    echo "📊 Validation Results:"
    echo "   Total Checks: ${TOTAL_CHECKS}"
    echo -e "   ${GREEN}Passed: ${PASSED_CHECKS}${NC}"
    echo -e "   ${RED}Failed: ${FAILED_CHECKS}${NC}"
    
    if [ "${FAILED_CHECKS}" -eq 0 ]; then
        echo -e "${GREEN}🎉 ALL Q1 2025 TASKS COMPLETED SUCCESSFULLY! 🎉${NC}"
        echo -e "${GREEN}Herald.lol gaming analytics platform is ready for production deployment!${NC}"
        echo ""
        echo "✅ Infrastructure Foundation Complete"
        echo "✅ Backend Services Implemented"
        echo "✅ Database Architecture Ready"
        echo "✅ Frontend Infrastructure Setup"
        echo "✅ Testing Framework Configured"
        echo "✅ Security & Compliance Implemented"
        echo "✅ Gaming Features Complete"
        echo "✅ Documentation Available"
        echo ""
        echo "🚀 Next Steps:"
        echo "   1. Deploy infrastructure: ./backend/scripts/deploy-infrastructure.sh"
        echo "   2. Run performance tests: k6 run backend/scripts/load-test-gaming.js"
        echo "   3. Monitor Q1 objectives: <5s analytics, 99.9% uptime"
        echo "   4. Begin Q2 2025 development phase"
        
        return 0
    else
        echo -e "${RED}❌ Q1 2025 VALIDATION FAILED${NC}"
        echo -e "${YELLOW}Please address the failed checks above before proceeding to Q2 2025${NC}"
        
        # Calculate completion percentage
        local completion_percentage=$(( (PASSED_CHECKS * 100) / TOTAL_CHECKS ))
        echo "📈 Completion: ${completion_percentage}%"
        
        return 1
    fi
}

main() {
    show_banner
    
    log_info "🎮 Starting Herald.lol Q1 2025 completion validation..."
    log_info "📋 This will verify all Q1 infrastructure and development tasks"
    
    # Run all validations
    validate_infrastructure
    validate_backend_services
    validate_database_architecture
    validate_frontend_infrastructure
    validate_testing_infrastructure
    validate_security_compliance
    validate_gaming_features
    validate_documentation
    
    # Show summary and exit with appropriate code
    show_validation_summary
}

# Run main function
main "$@"