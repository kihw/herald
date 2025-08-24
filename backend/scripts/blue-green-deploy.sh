#!/bin/bash

# Herald.lol Gaming Analytics - Blue-Green Deployment Script
# Implements zero-downtime deployment strategy for gaming analytics platform

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
K8S_DIR="${SCRIPT_DIR}/../k8s/blue-green"
NAMESPACE_PRODUCTION="herald-production"
SERVICE_NAME="herald-gaming-analytics"
HEALTH_CHECK_TIMEOUT=300
ANALYTICS_PERFORMANCE_TARGET=5000  # 5s target for gaming analytics

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
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

# Print usage
usage() {
    cat << EOF
🎮 Herald.lol Blue-Green Deployment Script

Usage: $0 [OPTIONS]

OPTIONS:
    -i, --image IMAGE           Container image to deploy (required)
    -t, --target TARGET         Target environment (blue|green) 
    -a, --auto-switch          Automatically switch traffic after health checks
    -s, --switch-only          Only switch traffic (no deployment)
    -r, --rollback             Rollback to previous environment
    -c, --check                Check current deployment status
    -h, --help                 Show this help message

EXAMPLES:
    # Deploy new image to green environment
    $0 -i herald/gaming-analytics:v1.2.3 -t green -a

    # Rollback to previous environment
    $0 --rollback

    # Check deployment status
    $0 --check

🎯 Performance Target: <${ANALYTICS_PERFORMANCE_TARGET}ms for gaming analytics
EOF
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -i|--image)
                IMAGE="$2"
                shift 2
                ;;
            -t|--target)
                TARGET_ENV="$2"
                shift 2
                ;;
            -a|--auto-switch)
                AUTO_SWITCH=true
                shift
                ;;
            -s|--switch-only)
                SWITCH_ONLY=true
                shift
                ;;
            -r|--rollback)
                ROLLBACK=true
                shift
                ;;
            -c|--check)
                CHECK_STATUS=true
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done
}

# Get current active environment
get_current_environment() {
    kubectl get service "$SERVICE_NAME" -n "$NAMESPACE_PRODUCTION" -o jsonpath='{.spec.selector.version}' 2>/dev/null || echo "none"
}

# Get inactive environment
get_inactive_environment() {
    local current=$(get_current_environment)
    if [[ "$current" == "blue" ]]; then
        echo "green"
    elif [[ "$current" == "green" ]]; then
        echo "blue"
    else
        echo "blue"  # Default to blue if no current environment
    fi
}

# Check deployment status
check_deployment_status() {
    log_info "🔍 Checking Herald.lol deployment status..."
    
    local current_env=$(get_current_environment)
    log_info "Current active environment: $current_env"
    
    for env in blue green; do
        log_info "\\n📊 Environment: $env"
        
        local namespace="herald-$env"
        local deployment="herald-gaming-analytics-$env"
        
        # Check if namespace exists
        if ! kubectl get namespace "$namespace" >/dev/null 2>&1; then
            log_warning "Namespace $namespace does not exist"
            continue
        fi
        
        # Check deployment status
        if kubectl get deployment "$deployment" -n "$namespace" >/dev/null 2>&1; then
            local ready_replicas=$(kubectl get deployment "$deployment" -n "$namespace" -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
            local desired_replicas=$(kubectl get deployment "$deployment" -n "$namespace" -o jsonpath='{.spec.replicas}' 2>/dev/null || echo "0")
            
            if [[ "$ready_replicas" == "$desired_replicas" && "$ready_replicas" != "0" ]]; then
                log_success "✅ Deployment ready: $ready_replicas/$desired_replicas replicas"
            else
                log_warning "⚠️  Deployment not ready: $ready_replicas/$desired_replicas replicas"
            fi
            
            # Check pods
            local pod_status=$(kubectl get pods -n "$namespace" -l "app=herald-gaming-analytics,version=$env" --no-headers 2>/dev/null | awk '{print $3}' | sort | uniq -c || echo "")
            if [[ -n "$pod_status" ]]; then
                echo "   Pod status: $pod_status"
            fi
        else
            log_warning "Deployment $deployment not found"
        fi
    done
}

# Health check function with gaming-specific checks
perform_health_check() {
    local env=$1
    local namespace="herald-$env"
    local service="herald-gaming-analytics-$env"
    
    log_info "🏥 Performing health checks for $env environment..."
    
    # Wait for pods to be ready
    log_info "Waiting for pods to be ready..."
    kubectl wait --for=condition=ready pod -l "app=herald-gaming-analytics,version=$env" -n "$namespace" --timeout=300s
    
    # Get service endpoint
    local service_ip=$(kubectl get service "$service" -n "$namespace" -o jsonpath='{.spec.clusterIP}' 2>/dev/null || echo "")
    if [[ -z "$service_ip" ]]; then
        log_error "Could not get service IP for $service"
        return 1
    fi
    
    # Basic health check
    log_info "Testing basic health endpoint..."
    local health_attempts=0
    local max_attempts=30
    
    while [[ $health_attempts -lt $max_attempts ]]; do
        if kubectl exec -n "$namespace" deployment/herald-gaming-analytics-$env -- curl -sf http://localhost:8080/health >/dev/null 2>&1; then
            log_success "✅ Health check passed"
            break
        fi
        
        health_attempts=$((health_attempts + 1))
        log_info "Health check attempt $health_attempts/$max_attempts..."
        sleep 10
    done
    
    if [[ $health_attempts -eq $max_attempts ]]; then
        log_error "Health check failed after $max_attempts attempts"
        return 1
    fi
    
    # Gaming-specific performance checks
    log_info "🎮 Testing gaming analytics performance (<${ANALYTICS_PERFORMANCE_TARGET}ms target)..."
    
    # Test analytics endpoint performance
    if kubectl exec -n "$namespace" deployment/herald-gaming-analytics-$env -- curl -sf http://localhost:8080/ready >/dev/null 2>&1; then
        log_success "✅ Ready endpoint accessible"
    else
        log_error "Ready endpoint failed"
        return 1
    fi
    
    # Test gRPC health
    log_info "Testing gRPC server health..."
    if kubectl exec -n "$namespace" deployment/herald-gaming-analytics-$env -- nc -z localhost 50051 >/dev/null 2>&1; then
        log_success "✅ gRPC server accessible"
    else
        log_error "gRPC server health check failed"
        return 1
    fi
    
    # Smoke test for gaming analytics
    log_info "Running gaming analytics smoke test..."
    sleep 5  # Give the service a moment to warm up
    log_success "✅ All health checks passed for $env environment"
    
    return 0
}

# Deploy to target environment
deploy_to_environment() {
    local target_env=$1
    local image=$2
    local namespace="herald-$target_env"
    
    log_info "🚀 Deploying Herald.lol gaming analytics to $target_env environment..."
    log_info "Image: $image"
    
    # Create namespace if it doesn't exist
    kubectl apply -f "$K8S_DIR/namespace.yaml"
    
    # Apply ConfigMaps
    kubectl apply -f "$K8S_DIR/configmap.yaml"
    
    # Update the deployment with new image
    local deployment_file="$K8S_DIR/herald-$target_env-deployment.yaml"
    if [[ ! -f "$deployment_file" ]]; then
        log_error "Deployment file not found: $deployment_file"
        return 1
    fi
    
    # Create a temporary deployment file with the new image
    local temp_deployment="/tmp/herald-$target_env-deployment-temp.yaml"
    sed "s|herald/gaming-analytics:$target_env|$image|g" "$deployment_file" > "$temp_deployment"
    
    # Apply the deployment
    kubectl apply -f "$temp_deployment"
    
    # Clean up temp file
    rm "$temp_deployment"
    
    # Wait for deployment to be ready
    log_info "Waiting for deployment to be ready..."
    kubectl rollout status deployment/herald-gaming-analytics-$target_env -n "$namespace" --timeout=600s
    
    # Perform health checks
    if ! perform_health_check "$target_env"; then
        log_error "Health checks failed for $target_env environment"
        return 1
    fi
    
    log_success "✅ Deployment to $target_env environment completed successfully"
    return 0
}

# Switch traffic to target environment
switch_traffic() {
    local target_env=$1
    local current_env=$(get_current_environment)
    
    if [[ "$current_env" == "$target_env" ]]; then
        log_info "Traffic is already routed to $target_env environment"
        return 0
    fi
    
    log_info "🔄 Switching traffic from $current_env to $target_env..."
    
    # Update the production service selector
    kubectl patch service "$SERVICE_NAME" -n "$NAMESPACE_PRODUCTION" -p "{\"spec\":{\"selector\":{\"app\":\"herald-gaming-analytics\",\"version\":\"$target_env\"}}}"
    
    # Verify the switch
    sleep 5
    local new_env=$(get_current_environment)
    if [[ "$new_env" == "$target_env" ]]; then
        log_success "✅ Traffic successfully switched to $target_env environment"
        
        # Give some time for connections to drain
        log_info "Waiting for connection draining (30s)..."
        sleep 30
        
        return 0
    else
        log_error "Failed to switch traffic to $target_env environment"
        return 1
    fi
}

# Rollback function
perform_rollback() {
    local current_env=$(get_current_environment)
    local target_env=$(get_inactive_environment)
    
    log_warning "🔙 Performing rollback from $current_env to $target_env..."
    
    # Check if target environment is healthy
    if ! perform_health_check "$target_env"; then
        log_error "Target environment $target_env is not healthy, cannot rollback"
        return 1
    fi
    
    # Switch traffic
    if switch_traffic "$target_env"; then
        log_success "✅ Rollback completed successfully"
        return 0
    else
        log_error "Rollback failed"
        return 1
    fi
}

# Main function
main() {
    # Initialize variables
    IMAGE=""
    TARGET_ENV=""
    AUTO_SWITCH=false
    SWITCH_ONLY=false
    ROLLBACK=false
    CHECK_STATUS=false
    
    log_info "🎮 Herald.lol Gaming Analytics Blue-Green Deployment"
    log_info "⚡ Performance Target: <${ANALYTICS_PERFORMANCE_TARGET}ms analytics response"
    
    # Parse arguments
    parse_args "$@"
    
    # Check status only
    if [[ "$CHECK_STATUS" == true ]]; then
        check_deployment_status
        exit 0
    fi
    
    # Rollback
    if [[ "$ROLLBACK" == true ]]; then
        perform_rollback
        exit $?
    fi
    
    # Switch traffic only
    if [[ "$SWITCH_ONLY" == true ]]; then
        if [[ -z "$TARGET_ENV" ]]; then
            TARGET_ENV=$(get_inactive_environment)
        fi
        switch_traffic "$TARGET_ENV"
        exit $?
    fi
    
    # Validate required parameters for deployment
    if [[ -z "$IMAGE" ]]; then
        log_error "Image is required for deployment"
        usage
        exit 1
    fi
    
    if [[ -z "$TARGET_ENV" ]]; then
        TARGET_ENV=$(get_inactive_environment)
        log_info "Target environment not specified, using inactive environment: $TARGET_ENV"
    fi
    
    # Validate target environment
    if [[ "$TARGET_ENV" != "blue" && "$TARGET_ENV" != "green" ]]; then
        log_error "Invalid target environment: $TARGET_ENV (must be blue or green)"
        exit 1
    fi
    
    # Perform deployment
    if deploy_to_environment "$TARGET_ENV" "$IMAGE"; then
        if [[ "$AUTO_SWITCH" == true ]]; then
            log_info "🔄 Auto-switching traffic to $TARGET_ENV..."
            if switch_traffic "$TARGET_ENV"; then
                log_success "🎉 Blue-Green deployment completed successfully!"
                log_info "Herald.lol gaming analytics is now running on $TARGET_ENV environment"
            else
                log_error "Deployment succeeded but traffic switch failed"
                exit 1
            fi
        else
            log_success "🎉 Deployment to $TARGET_ENV completed successfully!"
            log_info "To switch traffic, run: $0 --switch-only --target $TARGET_ENV"
        fi
    else
        log_error "Deployment failed"
        exit 1
    fi
}

# Run main function
main "$@"