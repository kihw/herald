#!/bin/bash

# Herald.lol Gaming Analytics Platform - Infrastructure Deployment Script
# Automated deployment of AWS infrastructure using Terraform

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
TERRAFORM_DIR="${PROJECT_ROOT}/backend/terraform"
K8S_DIR="${PROJECT_ROOT}/backend/k8s"

# Default values
ENVIRONMENT="production"
AWS_REGION="us-east-1"
CLUSTER_NAME="herald-gaming-cluster"
DRY_RUN=false
SKIP_CONFIRMATION=false
DESTROY_MODE=false

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
    
    🎮 Gaming Analytics Platform - Infrastructure Deployment 🎮
EOF
    echo -e "${NC}"
}

show_help() {
    cat << EOF
Herald.lol Infrastructure Deployment Script

Usage: $0 [OPTIONS]

OPTIONS:
    -e, --environment ENV       Deployment environment (default: production)
    -r, --region REGION        AWS region (default: us-east-1)
    -c, --cluster-name NAME    EKS cluster name (default: herald-gaming-cluster)
    -d, --dry-run             Show what would be deployed without executing
    -y, --yes                 Skip confirmation prompts
    -D, --destroy             Destroy infrastructure (USE WITH CAUTION!)
    -h, --help                Show this help message

EXAMPLES:
    # Deploy production infrastructure
    $0 -e production -r us-east-1

    # Dry run deployment
    $0 --dry-run

    # Deploy to staging with auto-confirmation
    $0 -e staging -y

    # Destroy infrastructure (BE CAREFUL!)
    $0 --destroy -e staging
EOF
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if required tools are installed
    local missing_tools=()
    
    if ! command -v terraform &> /dev/null; then
        missing_tools+=("terraform")
    fi
    
    if ! command -v aws &> /dev/null; then
        missing_tools+=("aws-cli")
    fi
    
    if ! command -v kubectl &> /dev/null; then
        missing_tools+=("kubectl")
    fi
    
    if ! command -v helm &> /dev/null; then
        missing_tools+=("helm")
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_error "Missing required tools: ${missing_tools[*]}"
        log_error "Please install missing tools and retry"
        exit 1
    fi
    
    # Check AWS credentials
    if ! aws sts get-caller-identity &> /dev/null; then
        log_error "AWS credentials not configured or invalid"
        log_error "Please run: aws configure"
        exit 1
    fi
    
    # Check Terraform directory exists
    if [ ! -d "${TERRAFORM_DIR}" ]; then
        log_error "Terraform directory not found: ${TERRAFORM_DIR}"
        exit 1
    fi
    
    log_success "All prerequisites satisfied"
}

validate_aws_permissions() {
    log_info "Validating AWS permissions..."
    
    local required_permissions=(
        "eks:CreateCluster"
        "ec2:CreateVpc"
        "rds:CreateDBCluster"
        "elasticache:CreateReplicationGroup"
        "s3:CreateBucket"
        "iam:CreateRole"
        "cloudformation:CreateStack"
    )
    
    # Simple validation by checking if we can list resources
    if ! aws eks list-clusters --region "${AWS_REGION}" &> /dev/null; then
        log_warning "May not have sufficient EKS permissions"
    fi
    
    if ! aws ec2 describe-vpcs --region "${AWS_REGION}" &> /dev/null; then
        log_warning "May not have sufficient EC2 permissions"
    fi
    
    log_success "AWS permissions validation completed"
}

init_terraform() {
    log_info "Initializing Terraform..."
    
    cd "${TERRAFORM_DIR}"
    
    # Initialize Terraform
    terraform init -input=false
    
    # Validate configuration
    terraform validate
    
    log_success "Terraform initialized successfully"
}

plan_deployment() {
    log_info "Planning deployment..."
    
    cd "${TERRAFORM_DIR}"
    
    # Create terraform plan
    terraform plan \
        -var="environment=${ENVIRONMENT}" \
        -var="aws_region=${AWS_REGION}" \
        -var="cluster_name=${CLUSTER_NAME}" \
        -out="herald-gaming.tfplan" \
        -input=false
    
    log_success "Deployment plan created: herald-gaming.tfplan"
}

apply_deployment() {
    log_info "Applying deployment..."
    
    cd "${TERRAFORM_DIR}"
    
    if [ "${DRY_RUN}" = true ]; then
        log_info "DRY RUN MODE - Would apply: herald-gaming.tfplan"
        terraform show herald-gaming.tfplan
        return
    fi
    
    # Apply the plan
    terraform apply \
        -input=false \
        -auto-approve \
        herald-gaming.tfplan
    
    log_success "Infrastructure deployment completed"
}

configure_kubectl() {
    log_info "Configuring kubectl for Herald gaming cluster..."
    
    # Update kubeconfig
    aws eks update-kubeconfig \
        --region "${AWS_REGION}" \
        --name "${CLUSTER_NAME}" \
        --alias "herald-gaming-${ENVIRONMENT}"
    
    # Test connection
    if kubectl cluster-info &> /dev/null; then
        log_success "kubectl configured successfully"
        kubectl get nodes
    else
        log_error "Failed to configure kubectl"
        exit 1
    fi
}

deploy_istio() {
    log_info "Deploying Istio service mesh..."
    
    if [ "${DRY_RUN}" = true ]; then
        log_info "DRY RUN MODE - Would deploy Istio service mesh"
        return
    fi
    
    # Install Istio
    if ! command -v istioctl &> /dev/null; then
        log_warning "istioctl not found, installing..."
        curl -L https://istio.io/downloadIstio | sh -
        export PATH="$PWD/istio-*/bin:$PATH"
    fi
    
    # Apply Istio configuration
    kubectl apply -f "${K8S_DIR}/istio/istio-install.yaml"
    
    # Wait for Istio to be ready
    kubectl wait --for=condition=Ready pods --all -n istio-system --timeout=300s
    
    log_success "Istio service mesh deployed successfully"
}

show_deployment_summary() {
    log_info "Deployment Summary:"
    
    cd "${TERRAFORM_DIR}"
    
    echo -e "${GREEN}================================${NC}"
    echo -e "${GREEN}Herald.lol Gaming Infrastructure${NC}"
    echo -e "${GREEN}================================${NC}"
    
    # Get outputs
    local cluster_endpoint=$(terraform output -raw cluster_endpoint 2>/dev/null || echo "N/A")
    local rds_endpoint=$(terraform output -raw rds_cluster_endpoint 2>/dev/null || echo "N/A")
    local redis_endpoint=$(terraform output -raw redis_configuration_endpoint 2>/dev/null || echo "N/A")
    local cdn_domain=$(terraform output -raw cloudfront_domain_name 2>/dev/null || echo "N/A")
    
    echo "Environment: ${ENVIRONMENT}"
    echo "Region: ${AWS_REGION}"
    echo "Cluster: ${CLUSTER_NAME}"
    echo ""
    echo "Endpoints:"
    echo "  EKS Cluster: ${cluster_endpoint}"
    echo "  Database: ${rds_endpoint}"
    echo "  Redis: ${redis_endpoint}"
    echo "  CDN: ${cdn_domain}"
    echo ""
    echo "kubectl configuration:"
    echo "  aws eks update-kubeconfig --region ${AWS_REGION} --name ${CLUSTER_NAME}"
    echo -e "${GREEN}================================${NC}"
}

cleanup_deployment() {
    log_warning "Destroying Herald gaming infrastructure..."
    
    if [ "${SKIP_CONFIRMATION}" = false ]; then
        echo -e "${RED}WARNING: This will destroy all Herald gaming infrastructure!${NC}"
        read -p "Are you absolutely sure? Type 'YES' to continue: " confirmation
        
        if [ "${confirmation}" != "YES" ]; then
            log_info "Destruction cancelled"
            exit 0
        fi
    fi
    
    cd "${TERRAFORM_DIR}"
    
    terraform destroy \
        -var="environment=${ENVIRONMENT}" \
        -var="aws_region=${AWS_REGION}" \
        -var="cluster_name=${CLUSTER_NAME}" \
        -auto-approve
    
    log_success "Infrastructure destroyed"
}

main() {
    show_banner
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -e|--environment)
                ENVIRONMENT="$2"
                shift 2
                ;;
            -r|--region)
                AWS_REGION="$2"
                shift 2
                ;;
            -c|--cluster-name)
                CLUSTER_NAME="$2"
                shift 2
                ;;
            -d|--dry-run)
                DRY_RUN=true
                shift
                ;;
            -y|--yes)
                SKIP_CONFIRMATION=true
                shift
                ;;
            -D|--destroy)
                DESTROY_MODE=true
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
    
    log_info "Starting Herald gaming infrastructure deployment"
    log_info "Environment: ${ENVIRONMENT}"
    log_info "Region: ${AWS_REGION}"
    log_info "Cluster: ${CLUSTER_NAME}"
    
    if [ "${DESTROY_MODE}" = true ]; then
        cleanup_deployment
        exit 0
    fi
    
    # Confirmation
    if [ "${SKIP_CONFIRMATION}" = false ] && [ "${DRY_RUN}" = false ]; then
        echo -e "${YELLOW}This will deploy Herald gaming infrastructure to AWS${NC}"
        read -p "Continue? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Deployment cancelled"
            exit 0
        fi
    fi
    
    # Execute deployment steps
    check_prerequisites
    validate_aws_permissions
    init_terraform
    plan_deployment
    apply_deployment
    
    if [ "${DRY_RUN}" = false ]; then
        configure_kubectl
        deploy_istio
        show_deployment_summary
        
        log_success "Herald gaming infrastructure deployment completed! 🎮"
        log_info "You can now deploy Herald gaming services to the cluster"
    else
        log_info "Dry run completed - no resources were created"
    fi
}

# Trap to cleanup on exit
trap 'log_error "Deployment interrupted"' INT TERM

# Run main function
main "$@"