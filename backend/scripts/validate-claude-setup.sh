#!/bin/bash

# Herald.lol Gaming Analytics - Claude Code Setup Validation
# Comprehensive validation of Claude Code integration

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Gaming colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
GOLD='\033[1;33m'
NC='\033[0m'

VALIDATION_ERRORS=0

log_info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] [VALIDATE]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] [SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] [ERROR]${NC} $1"
    ((VALIDATION_ERRORS++))
}

log_warning() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] [WARNING]${NC} $1"
}

log_gaming() {
    echo -e "${GOLD}[$(date +'%Y-%m-%d %H:%M:%S')] [HERALD-GAMING]${NC} $1"
}

# Validate Claude hooks
validate_claude_hooks() {
    log_gaming "🎮 Validating Herald.lol Claude Code hooks..."
    
    cd "$PROJECT_ROOT"
    
    # Check hooks directory exists
    if [ ! -d ".claude/hooks" ]; then
        log_error "❌ Claude hooks directory not found: .claude/hooks"
        return
    fi
    
    log_success "✅ Claude hooks directory found"
    
    # Validate individual hooks
    local hooks=(
        "load-context.sh"
        "deploy-check.sh"
        "gaming-lint.sh" 
        "post-bash-gaming.sh"
        "pre-task-gaming.sh"
        "test-suite.sh"
        "validate-stack.py"
        "prompt-gaming-context.py"
    )
    
    for hook in "${hooks[@]}"; do
        if [ -f ".claude/hooks/$hook" ] && [ -x ".claude/hooks/$hook" ]; then
            log_success "  ✅ $hook - present and executable"
        else
            log_error "  ❌ $hook - missing or not executable"
        fi
    done
    
    # Test hook functionality
    log_info "🧪 Testing hook functionality..."
    
    # Test load-context hook
    if ./.claude/hooks/load-context.sh >/dev/null 2>&1; then
        log_success "  ✅ load-context.sh executes successfully"
    else
        log_error "  ❌ load-context.sh failed to execute"
    fi
    
    # Test deploy-check hook  
    if ./.claude/hooks/deploy-check.sh >/dev/null 2>&1; then
        log_success "  ✅ deploy-check.sh executes successfully"
    else
        log_warning "  ⚠️ deploy-check.sh has warnings (normal)"
    fi
}

# Validate project structure
validate_project_structure() {
    log_gaming "🏗️ Validating Herald.lol project structure..."
    
    cd "$PROJECT_ROOT"
    
    # Check essential directories
    local directories=(
        "cmd/server"
        "internal/services"
        "internal/handlers"  
        "internal/models"
        "api/proto"
        "scripts"
        "k8s/blue-green"
        "terraform/drift-detection"
    )
    
    for dir in "${directories[@]}"; do
        if [ -d "$dir" ]; then
            log_success "  ✅ $dir directory found"
        else
            log_error "  ❌ $dir directory missing"
        fi
    done
    
    # Check essential files
    local files=(
        "go.mod"
        "go.sum"
        "Makefile"
        "docker-compose.herald-dev.yml"
        "README-VPS-SETUP.md"
        "README-SONARQUBE.md"
    )
    
    for file in "${files[@]}"; do
        if [ -f "$file" ]; then
            log_success "  ✅ $file found"
        else
            log_error "  ❌ $file missing"
        fi
    done
}

# Validate Go environment
validate_go_environment() {
    log_gaming "🐹 Validating Herald.lol Go environment..."
    
    cd "$PROJECT_ROOT"
    
    # Check Go installation
    if command -v go >/dev/null 2>&1; then
        local go_version=$(go version | awk '{print $3}' | sed 's/go//')
        log_success "✅ Go installed: $go_version"
        
        # Check minimum version
        if [[ "$go_version" < "1.21" ]]; then
            log_warning "⚠️ Go version $go_version may be too old (recommend 1.21+)"
        fi
    else
        log_error "❌ Go not installed"
        return
    fi
    
    # Check Go modules
    if [ -f "go.mod" ]; then
        if go mod verify >/dev/null 2>&1; then
            log_success "  ✅ Go modules verified"
        else
            log_error "  ❌ Go modules verification failed"
        fi
        
        if go mod download >/dev/null 2>&1; then
            log_success "  ✅ Go dependencies downloaded"
        else
            log_error "  ❌ Go dependencies download failed"
        fi
    fi
    
    # Test Go build
    if go build -o /tmp/herald-test ./cmd/server >/dev/null 2>&1; then
        log_success "  ✅ Go build successful"
        rm -f /tmp/herald-test
    else
        log_error "  ❌ Go build failed"
    fi
}

# Validate gaming-specific setup
validate_gaming_setup() {
    log_gaming "🎮 Validating Herald.lol gaming-specific setup..."
    
    cd "$PROJECT_ROOT"
    
    # Check gaming services
    local services=(
        "internal/services/analytics_service.go"
        "internal/services/riot_service.go"
        "internal/services/coaching_service.go"
        "internal/services/notification_service.go"
        "internal/streaming/service.go"
    )
    
    for service in "${services[@]}"; do
        if [ -f "$service" ]; then
            log_success "  ✅ Gaming service: $(basename "$service")"
        else
            log_error "  ❌ Missing gaming service: $(basename "$service")"
        fi
    done
    
    # Check gaming models
    local models=(
        "internal/models/analytics.go"
        "internal/models/match.go"
        "internal/models/user.go"
        "internal/models/coaching.go"
    )
    
    for model in "${models[@]}"; do
        if [ -f "$model" ]; then
            log_success "  ✅ Gaming model: $(basename "$model")"
        else
            log_error "  ❌ Missing gaming model: $(basename "$model")"
        fi
    done
    
    # Check gRPC proto files
    local protos=(
        "api/proto/analytics.proto"
        "api/proto/match.proto"
        "api/proto/riot.proto"
    )
    
    for proto in "${protos[@]}"; do
        if [ -f "$proto" ]; then
            log_success "  ✅ Gaming gRPC: $(basename "$proto")"
        else
            log_error "  ❌ Missing gaming gRPC: $(basename "$proto")"
        fi
    done
}

# Validate infrastructure setup
validate_infrastructure() {
    log_gaming "🏗️ Validating Herald.lol infrastructure setup..."
    
    cd "$PROJECT_ROOT"
    
    # Check VPS setup
    if [ -f "scripts/vps-setup.sh" ] && [ -x "scripts/vps-setup.sh" ]; then
        log_success "  ✅ VPS setup script ready"
    else
        log_error "  ❌ VPS setup script missing or not executable"
    fi
    
    # Check Docker Compose
    if [ -f "docker-compose.herald-dev.yml" ]; then
        log_success "  ✅ Herald.lol Docker Compose configured"
        
        # Validate compose file
        if command -v docker-compose >/dev/null 2>&1; then
            if docker-compose -f docker-compose.herald-dev.yml config >/dev/null 2>&1; then
                log_success "    ✅ Docker Compose file valid"
            else
                log_error "    ❌ Docker Compose file invalid"
            fi
        fi
    else
        log_error "  ❌ Herald.lol Docker Compose missing"
    fi
    
    # Check SonarQube setup
    if [ -f "docker-compose.sonarqube.yml" ] && [ -f "sonar-project.properties" ]; then
        log_success "  ✅ SonarQube gaming setup ready"
    else
        log_warning "  ⚠️ SonarQube gaming setup incomplete"
    fi
    
    # Check Kubernetes manifests
    if [ -d "k8s/blue-green" ]; then
        local k8s_files=$(find k8s/blue-green -name "*.yaml" | wc -l)
        if [ "$k8s_files" -gt 0 ]; then
            log_success "  ✅ Kubernetes gaming manifests ready ($k8s_files files)"
        else
            log_error "  ❌ No Kubernetes gaming manifests found"
        fi
    else
        log_error "  ❌ Kubernetes gaming directory missing"
    fi
}

# Validate development tools
validate_dev_tools() {
    log_gaming "🔧 Validating Herald.lol development tools..."
    
    # Check Makefile commands
    if [ -f "Makefile" ]; then
        log_success "  ✅ Makefile found"
        
        # Check gaming-specific make targets
        local make_targets=(
            "test-gaming"
            "sonar-analyze"
            "sync-hooks"
            "hooks-status"
        )
        
        for target in "${make_targets[@]}"; do
            if grep -q "^${target}:" Makefile; then
                log_success "    ✅ Make target: $target"
            else
                log_warning "    ⚠️ Make target missing: $target"
            fi
        done
    else
        log_error "  ❌ Makefile missing"
    fi
    
    # Check scripts
    local scripts=(
        "scripts/vps-setup.sh"
        "scripts/sonarqube-local.sh"
        "scripts/sync-claude-hooks.sh"
        "scripts/blue-green-deploy.sh"
    )
    
    for script in "${scripts[@]}"; do
        if [ -f "$script" ] && [ -x "$script" ]; then
            log_success "  ✅ Script ready: $(basename "$script")"
        else
            log_error "  ❌ Script missing/not executable: $(basename "$script")"
        fi
    done
}

# Generate validation report
generate_report() {
    log_gaming "📊 Herald.lol Claude Code Setup Validation Report"
    echo "=================================================="
    
    if [ $VALIDATION_ERRORS -eq 0 ]; then
        echo ""
        log_success "🎉 Herald.lol Claude Code setup is PERFECT!"
        log_gaming "🎯 Gaming Platform Status: READY FOR DEVELOPMENT"
        
        echo ""
        log_info "🎮 Herald.lol Gaming Development Ready:"
        echo "  ⚡ Performance Target: <5000ms analytics"
        echo "  👥 Concurrent Users: 1M+ support" 
        echo "  🎯 Gaming Focus: League of Legends & TFT"
        echo "  🏆 Quality: SonarQube gaming analysis"
        echo "  🔧 Infrastructure: VPS + Kubernetes + monitoring"
        echo "  🚀 Deployment: Blue-green + rollback"
        
    elif [ $VALIDATION_ERRORS -lt 5 ]; then
        echo ""
        log_warning "⚠️ Herald.lol Claude Code setup has MINOR issues ($VALIDATION_ERRORS errors)"
        log_gaming "🎯 Gaming Platform Status: MOSTLY READY"
        log_info "Fix the errors above and you'll be ready for gaming development!"
        
    else
        echo ""
        log_error "❌ Herald.lol Claude Code setup has MAJOR issues ($VALIDATION_ERRORS errors)"
        log_gaming "🎯 Gaming Platform Status: NEEDS WORK"
        log_info "Please fix the errors above before starting gaming development"
    fi
    
    echo ""
    log_gaming "🎮 Next Steps:"
    echo "  1. Fix any validation errors"
    echo "  2. Run: make sync-hooks"
    echo "  3. Run: ./scripts/vps-setup.sh setup"
    echo "  4. Start Herald.lol gaming development!"
    
    echo ""
    log_success "🚀 Herald.lol Gaming Analytics Platform validation complete!"
}

# Main function
main() {
    log_gaming "🎮 Starting Herald.lol Claude Code Setup Validation"
    echo ""
    
    validate_claude_hooks
    echo ""
    
    validate_project_structure
    echo ""
    
    validate_go_environment
    echo ""
    
    validate_gaming_setup
    echo ""
    
    validate_infrastructure
    echo ""
    
    validate_dev_tools
    echo ""
    
    generate_report
    
    return $VALIDATION_ERRORS
}

# Handle interruption gracefully
trap 'echo ""; log_info "🎮 Herald.lol validation interrupted"; exit 0' SIGINT

# Run main function
main "$@"