#!/bin/bash

# Herald.lol Gaming Analytics - Infrastructure Monitoring Script
# Real-time monitoring and drift detection for gaming platform

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
TERRAFORM_DIR="$PROJECT_ROOT/terraform/drift-detection"
DRIFT_SCRIPT_DIR="$SCRIPT_DIR/drift-detection"
AWS_REGION="${AWS_REGION:-us-east-1}"
GAMING_PERFORMANCE_TARGET=5000  # 5s gaming analytics target

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
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
    echo -e "${PURPLE}[$(date +'%Y-%m-%d %H:%M:%S')] [GAMING]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_info "🔍 Checking prerequisites for Herald.lol infrastructure monitoring..."
    
    local missing_tools=()
    
    # Check AWS CLI
    if ! command -v aws >/dev/null 2>&1; then
        missing_tools+=("aws-cli")
    fi
    
    # Check Terraform
    if ! command -v terraform >/dev/null 2>&1; then
        missing_tools+=("terraform")
    fi
    
    # Check jq
    if ! command -v jq >/dev/null 2>&1; then
        missing_tools+=("jq")
    fi
    
    # Check Python
    if ! command -v python3 >/dev/null 2>&1; then
        missing_tools+=("python3")
    fi
    
    if [ ${#missing_tools[@]} -gt 0 ]; then
        log_error "❌ Missing required tools: ${missing_tools[*]}"
        log_info "Please install missing tools and try again"
        return 1
    fi
    
    # Check AWS credentials
    if ! aws sts get-caller-identity >/dev/null 2>&1; then
        log_error "❌ AWS credentials not configured"
        log_info "Configure AWS credentials: aws configure"
        return 1
    fi
    
    log_success "✅ All prerequisites met"
    return 0
}

# Deploy drift detection infrastructure
deploy_infrastructure() {
    log_info "🚀 Deploying Herald.lol gaming infrastructure drift detection..."
    
    if [ ! -d "$TERRAFORM_DIR" ]; then
        log_error "❌ Terraform directory not found: $TERRAFORM_DIR"
        return 1
    fi
    
    cd "$TERRAFORM_DIR"
    
    # Initialize Terraform
    log_info "🏗️ Initializing Terraform..."
    terraform init
    
    # Plan deployment
    log_info "📋 Planning gaming infrastructure deployment..."
    terraform plan \
        -var="gaming_environment=production" \
        -var="gaming_performance_target_ms=$GAMING_PERFORMANCE_TARGET" \
        -var="notification_email=${NOTIFICATION_EMAIL:-admin@herald.lol}" \
        -var="slack_webhook_url=${SLACK_WEBHOOK_URL:-}" \
        -out=tfplan
    
    # Apply if plan looks good
    read -p "🎮 Deploy Herald.lol gaming infrastructure monitoring? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        log_info "🚀 Deploying gaming infrastructure..."
        terraform apply tfplan
        log_success "✅ Gaming infrastructure deployed successfully"
    else
        log_info "📋 Deployment cancelled"
        return 1
    fi
    
    cd "$SCRIPT_DIR"
}

# Build and deploy Lambda function
deploy_lambda() {
    log_info "📦 Building and deploying gaming drift detection Lambda..."
    
    if [ ! -f "$DRIFT_SCRIPT_DIR/build_lambda.sh" ]; then
        log_error "❌ Lambda build script not found"
        return 1
    fi
    
    # Build Lambda package
    cd "$DRIFT_SCRIPT_DIR"
    ./build_lambda.sh
    
    # Deploy if package exists
    if [ -f "drift_detection.zip" ]; then
        log_info "🚀 Deploying gaming drift detection Lambda..."
        
        # Get Lambda function name from Terraform output
        cd "$TERRAFORM_DIR"
        local function_name=$(terraform output -raw drift_detection_function_arn | sed 's/.*://')
        
        if [ -n "$function_name" ]; then
            aws lambda update-function-code \
                --function-name "$function_name" \
                --zip-file "fileb://$DRIFT_SCRIPT_DIR/drift_detection.zip" \
                --region "$AWS_REGION"
            
            log_success "✅ Gaming drift detection Lambda deployed"
        else
            log_error "❌ Could not get Lambda function name from Terraform"
            return 1
        fi
    else
        log_error "❌ Lambda package not found"
        return 1
    fi
    
    cd "$SCRIPT_DIR"
}

# Monitor gaming infrastructure
monitor_infrastructure() {
    log_info "📊 Starting Herald.lol gaming infrastructure monitoring..."
    log_gaming "🎮 Performance target: <${GAMING_PERFORMANCE_TARGET}ms analytics"
    
    # Get Terraform outputs
    cd "$TERRAFORM_DIR"
    
    local sns_topic_arn=$(terraform output -raw sns_topic_arn 2>/dev/null || echo "")
    local dashboard_url=$(terraform output -raw dashboard_url 2>/dev/null || echo "")
    local snapshots_bucket=$(terraform output -raw snapshots_bucket 2>/dev/null || echo "")
    
    cd "$SCRIPT_DIR"
    
    echo ""
    log_info "🎮 Herald.lol Gaming Infrastructure Status"
    echo "=========================================="
    
    # Check EKS clusters
    log_info "🏗️ Checking gaming EKS clusters..."
    local eks_clusters=$(aws eks list-clusters --region "$AWS_REGION" --query 'clusters[?contains(@, `herald`) || contains(@, `gaming`)]' --output text 2>/dev/null || echo "")
    if [ -n "$eks_clusters" ]; then
        echo "  ✅ Gaming EKS clusters: $eks_clusters"
    else
        echo "  ⚠️ No gaming EKS clusters found"
    fi
    
    # Check RDS instances
    log_info "💾 Checking gaming database instances..."
    local rds_count=$(aws rds describe-db-instances --region "$AWS_REGION" --query 'DBInstances[?contains(DBInstanceIdentifier, `herald`) || contains(DBInstanceIdentifier, `gaming`)] | length(@)' --output text 2>/dev/null || echo "0")
    echo "  📊 Gaming RDS instances: $rds_count"
    
    # Check ElastiCache
    log_info "⚡ Checking gaming cache clusters..."
    local cache_count=$(aws elasticache describe-cache-clusters --region "$AWS_REGION" --query 'CacheClusters[?contains(CacheClusterId, `herald`) || contains(CacheClusterId, `gaming`)] | length(@)' --output text 2>/dev/null || echo "0")
    echo "  🚀 Gaming cache clusters: $cache_count"
    
    # Check Load Balancers
    log_info "🌐 Checking gaming load balancers..."
    local elb_count=$(aws elbv2 describe-load-balancers --region "$AWS_REGION" --query 'LoadBalancers[?contains(LoadBalancerName, `herald`) || contains(LoadBalancerName, `gaming`)] | length(@)' --output text 2>/dev/null || echo "0")
    echo "  ⚖️ Gaming load balancers: $elb_count"
    
    echo ""
    log_info "📈 Gaming Infrastructure Monitoring Links"
    if [ -n "$dashboard_url" ]; then
        echo "  🎮 Gaming Dashboard: $dashboard_url"
    fi
    if [ -n "$sns_topic_arn" ]; then
        echo "  📧 Alert Topic: $sns_topic_arn"
    fi
    if [ -n "$snapshots_bucket" ]; then
        echo "  📸 Snapshots Bucket: s3://$snapshots_bucket"
    fi
    
    echo ""
}

# Test drift detection
test_drift_detection() {
    log_info "🧪 Testing Herald.lol gaming drift detection..."
    
    # Get Lambda function name
    cd "$TERRAFORM_DIR"
    local function_arn=$(terraform output -raw drift_detection_function_arn 2>/dev/null || echo "")
    
    if [ -z "$function_arn" ]; then
        log_error "❌ Lambda function not deployed"
        return 1
    fi
    
    local function_name=$(echo "$function_arn" | sed 's/.*://')
    
    # Create test event
    local test_event='{"test": true, "gaming_platform": "Herald.lol"}'
    
    log_info "🚀 Invoking gaming drift detection Lambda..."
    local result=$(aws lambda invoke \
        --function-name "$function_name" \
        --payload "$test_event" \
        --region "$AWS_REGION" \
        /tmp/lambda_output.json 2>&1)
    
    if [ $? -eq 0 ]; then
        log_success "✅ Gaming drift detection test successful"
        
        # Show response
        if [ -f "/tmp/lambda_output.json" ]; then
            log_info "📋 Lambda response:"
            jq . /tmp/lambda_output.json 2>/dev/null || cat /tmp/lambda_output.json
            rm -f /tmp/lambda_output.json
        fi
    else
        log_error "❌ Gaming drift detection test failed: $result"
        return 1
    fi
    
    cd "$SCRIPT_DIR"
}

# Show real-time dashboard
show_dashboard() {
    while true; do
        clear
        echo "🎮 Herald.lol Gaming Infrastructure Dashboard"
        echo "============================================="
        echo "⚡ Performance Target: <${GAMING_PERFORMANCE_TARGET}ms"
        echo "🕐 Last Update: $(date)"
        echo ""
        
        # AWS Resources Summary
        log_info "📊 AWS Gaming Resources"
        echo "======================="
        
        # EKS
        local eks_clusters=$(aws eks list-clusters --region "$AWS_REGION" --query 'clusters[?contains(@, `herald`) || contains(@, `gaming`)]' --output table 2>/dev/null || echo "No clusters")
        echo "🏗️ EKS Clusters:"
        echo "$eks_clusters" | head -10
        echo ""
        
        # RDS
        log_info "💾 Gaming Databases"
        aws rds describe-db-instances --region "$AWS_REGION" \
            --query 'DBInstances[?contains(DBInstanceIdentifier, `herald`) || contains(DBInstanceIdentifier, `gaming`)].{Name:DBInstanceIdentifier,Status:DBInstanceStatus,Engine:Engine,Class:DBInstanceClass}' \
            --output table 2>/dev/null | head -10 || echo "No gaming databases found"
        echo ""
        
        # ElastiCache
        log_info "⚡ Gaming Cache Clusters"
        aws elasticache describe-cache-clusters --region "$AWS_REGION" \
            --query 'CacheClusters[?contains(CacheClusterId, `herald`) || contains(CacheClusterId, `gaming`)].{ID:CacheClusterId,Status:CacheClusterStatus,Engine:Engine,NodeType:CacheNodeType}' \
            --output table 2>/dev/null | head -10 || echo "No gaming cache clusters found"
        echo ""
        
        # Recent Lambda executions
        if command -v aws >/dev/null 2>&1; then
            log_info "🔍 Recent Drift Detections"
            aws logs describe-log-groups --region "$AWS_REGION" \
                --log-group-name-prefix "/aws/lambda/herald-gaming-drift" \
                --query 'logGroups[0].{LogGroup:logGroupName,LastEvent:lastEventTime}' \
                --output table 2>/dev/null | head -5 || echo "No recent executions"
        fi
        
        echo ""
        echo "Press Ctrl+C to exit, or wait 30s for refresh..."
        sleep 30
    done
}

# Clean up infrastructure
cleanup_infrastructure() {
    log_warning "🧹 Cleaning up Herald.lol gaming infrastructure monitoring..."
    
    read -p "⚠️ This will destroy ALL gaming infrastructure monitoring. Continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Cleanup cancelled"
        return 0
    fi
    
    cd "$TERRAFORM_DIR"
    
    # Destroy infrastructure
    log_warning "💥 Destroying gaming infrastructure monitoring..."
    terraform destroy -auto-approve
    
    log_success "✅ Gaming infrastructure monitoring cleaned up"
    
    cd "$SCRIPT_DIR"
}

# Usage
usage() {
    cat << EOF
🎮 Herald.lol Gaming Infrastructure Monitor

Usage: $0 [COMMAND]

COMMANDS:
    deploy          Deploy gaming infrastructure monitoring
    deploy-lambda   Build and deploy drift detection Lambda
    monitor         Show gaming infrastructure status
    test            Test drift detection functionality
    dashboard       Show real-time gaming dashboard
    cleanup         Clean up all monitoring infrastructure
    -h, --help      Show this help message

ENVIRONMENT VARIABLES:
    NOTIFICATION_EMAIL    Email for gaming alerts
    SLACK_WEBHOOK_URL     Slack webhook for gaming notifications
    AWS_REGION           AWS region (default: us-east-1)

EXAMPLES:
    # Deploy complete gaming monitoring
    NOTIFICATION_EMAIL=admin@herald.lol $0 deploy

    # Show real-time gaming dashboard  
    $0 dashboard

    # Test drift detection
    $0 test

🎯 Gaming Focus: Monitor Herald.lol infrastructure for <${GAMING_PERFORMANCE_TARGET}ms target
🎮 Platform: Optimized for gaming analytics and real-time performance
📊 Monitoring: EKS, RDS, ElastiCache, Load Balancers, Auto Scaling
EOF
}

# Main function
main() {
    case "${1:-}" in
        deploy)
            check_prerequisites
            deploy_infrastructure
            deploy_lambda
            monitor_infrastructure
            ;;
        deploy-lambda)
            check_prerequisites
            deploy_lambda
            ;;
        monitor)
            check_prerequisites
            monitor_infrastructure
            ;;
        test)
            check_prerequisites
            test_drift_detection
            ;;
        dashboard)
            check_prerequisites
            show_dashboard
            ;;
        cleanup)
            check_prerequisites
            cleanup_infrastructure
            ;;
        -h|--help)
            usage
            ;;
        *)
            log_info "🎮 Herald.lol Gaming Infrastructure Monitor"
            echo ""
            usage
            exit 1
            ;;
    esac
}

# Handle Ctrl+C gracefully
trap 'echo ""; log_info "Gaming infrastructure monitoring stopped"; exit 0' SIGINT

# Run main function
main "$@"