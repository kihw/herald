#!/bin/bash

# Herald.lol Gaming Analytics - Lambda Build Script
# Build and package drift detection Lambda for gaming infrastructure

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="$SCRIPT_DIR/build"
PACKAGE_NAME="drift_detection.zip"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Clean previous build
clean_build() {
    log_info "🧹 Cleaning previous build..."
    rm -rf "$BUILD_DIR"
    rm -f "$SCRIPT_DIR/$PACKAGE_NAME"
    mkdir -p "$BUILD_DIR"
}

# Install Python dependencies
install_dependencies() {
    log_info "📦 Installing Python dependencies for gaming drift detection..."
    
    # Create requirements.txt for Lambda
    cat << 'EOF' > "$BUILD_DIR/requirements.txt"
boto3==1.34.0
requests==2.31.0
botocore==1.34.0
urllib3==1.26.18
EOF

    # Install dependencies
    pip install -r "$BUILD_DIR/requirements.txt" -t "$BUILD_DIR" --no-deps
    
    log_success "✅ Dependencies installed"
}

# Copy Lambda function
copy_function() {
    log_info "📄 Copying gaming drift detection function..."
    
    # Copy main function
    cp "$SCRIPT_DIR/drift_detector.py" "$BUILD_DIR/main.py"
    
    # Create Lambda handler wrapper
    cat << 'EOF' > "$BUILD_DIR/lambda_function.py"
"""
Herald.lol Gaming Analytics - Lambda Handler
Entry point for AWS Lambda drift detection
"""

from main import handler

def lambda_handler(event, context):
    """AWS Lambda entry point"""
    return handler(event, context)
EOF
    
    log_success "✅ Function copied"
}

# Create configuration files
create_config() {
    log_info "⚙️ Creating gaming configuration files..."
    
    # Create gaming-specific configuration
    cat << 'EOF' > "$BUILD_DIR/gaming_config.json"
{
  "platform": "Herald.lol",
  "gaming_components": {
    "critical": ["eks_clusters", "elasticache", "load_balancers"],
    "performance": ["rds_instances", "auto_scaling_groups", "ec2_instances"],
    "monitoring": ["cloudwatch", "sns", "lambda"]
  },
  "performance_targets": {
    "analytics_response_ms": 5000,
    "concurrent_users": 1000000,
    "uptime_percentage": 99.9
  },
  "gaming_tags": {
    "required": ["Project", "Environment", "Component"],
    "gaming_identifiers": ["herald", "gaming", "analytics", "riot"]
  },
  "alert_thresholds": {
    "critical_score": 7,
    "performance_score": 5,
    "moderate_score": 3
  }
}
EOF

    # Create deployment metadata
    cat << EOF > "$BUILD_DIR/deployment_info.json"
{
  "build_timestamp": "$(date -Iseconds)",
  "platform": "Herald.lol Gaming Analytics",
  "version": "1.0.0",
  "environment": "production",
  "gaming_optimized": true,
  "performance_target_ms": 5000,
  "build_system": "$(uname -s)/$(uname -m)",
  "python_version": "$(python3 --version)"
}
EOF
    
    log_success "✅ Configuration created"
}

# Optimize for Lambda
optimize_package() {
    log_info "⚡ Optimizing package for gaming performance..."
    
    cd "$BUILD_DIR"
    
    # Remove unnecessary files to reduce cold start
    find . -name "*.pyc" -delete
    find . -name "__pycache__" -type d -exec rm -rf {} + 2>/dev/null || true
    find . -name "*.pyo" -delete
    find . -name "tests" -type d -exec rm -rf {} + 2>/dev/null || true
    find . -name "test_*" -delete
    find . -name "*_test.py" -delete
    
    # Remove documentation and examples
    find . -name "*.md" -delete
    find . -name "*.rst" -delete
    find . -name "*.txt" -delete 2>/dev/null || true
    find . -name "examples" -type d -exec rm -rf {} + 2>/dev/null || true
    find . -name "docs" -type d -exec rm -rf {} + 2>/dev/null || true
    
    # Keep only requirements.txt
    echo "boto3==1.34.0" > requirements.txt
    echo "requests==2.31.0" >> requirements.txt
    
    log_success "✅ Package optimized for gaming performance"
}

# Create deployment package
create_package() {
    log_info "📦 Creating deployment package for gaming infrastructure monitoring..."
    
    cd "$BUILD_DIR"
    
    # Create the ZIP package
    zip -r "../$PACKAGE_NAME" . -q
    
    cd "$SCRIPT_DIR"
    
    # Verify package
    local package_size=$(du -h "$PACKAGE_NAME" | cut -f1)
    log_success "✅ Gaming drift detection package created: $PACKAGE_NAME ($package_size)"
    
    # Show package contents summary
    log_info "📋 Package contents:"
    unzip -l "$PACKAGE_NAME" | head -20
    echo "..."
}

# Generate deployment documentation
generate_docs() {
    log_info "📚 Generating deployment documentation..."
    
    cat << EOF > "$SCRIPT_DIR/DEPLOYMENT.md"
# Herald.lol Gaming Infrastructure Drift Detection - Deployment Guide

## 🎮 Package Information
- **Platform:** Herald.lol Gaming Analytics
- **Build Date:** $(date)
- **Package:** $PACKAGE_NAME
- **Size:** $(du -h "$SCRIPT_DIR/$PACKAGE_NAME" | cut -f1)

## 🚀 Deployment Instructions

### 1. Terraform Deployment
\`\`\`bash
cd terraform/drift-detection
terraform init
terraform plan -var="notification_email=your-email@example.com"
terraform apply
\`\`\`

### 2. Manual Lambda Deployment
\`\`\`bash
aws lambda update-function-code \\
  --function-name herald-gaming-drift-detection \\
  --zip-file fileb://$PACKAGE_NAME
\`\`\`

### 3. Environment Variables
Set these environment variables in Lambda:
- \`GAMING_ENVIRONMENT\`: production/blue/green
- \`GAMING_PERFORMANCE_TARGET_MS\`: 5000
- \`GAMING_CONCURRENT_USERS\`: 1000000
- \`SNS_TOPIC_ARN\`: SNS topic ARN for alerts
- \`S3_SNAPSHOTS_BUCKET\`: S3 bucket for snapshots
- \`SLACK_WEBHOOK_URL\`: Slack webhook (optional)

## 🎯 Gaming Performance Configuration

### Lambda Configuration
- **Runtime:** Python 3.11
- **Timeout:** 300 seconds (5 minutes)
- **Memory:** 256 MB
- **Architecture:** x86_64

### Triggers
- **EventBridge Rule:** Every 4 hours
- **Manual Execution:** For immediate checks

## 🎮 Gaming-Specific Features

### Infrastructure Components Monitored
- **EKS Clusters:** Gaming container orchestration
- **RDS Instances:** Gaming analytics database
- **ElastiCache:** Gaming performance cache
- **Load Balancers:** Gaming traffic distribution
- **Auto Scaling Groups:** Gaming workload scaling
- **EC2 Instances:** Gaming compute resources

### Gaming Impact Analysis
- **Critical Impact:** Changes affecting core gaming functionality
- **Performance Impact:** Changes affecting <5s analytics target
- **Gaming-Specific:** Changes to gaming-tagged resources

### Alert Levels
- **Critical (7-10):** Immediate action required
- **Performance (5-6):** Monitor gaming performance
- **Moderate (3-4):** Review when convenient
- **Low (0-2):** Informational

## 📊 Monitoring

### CloudWatch Dashboard
Access the gaming infrastructure dashboard:
https://console.aws.amazon.com/cloudwatch/home#dashboards:name=Herald-Gaming-Infrastructure-Drift

### Metrics
- Lambda execution duration
- Drift detection frequency
- Gaming risk scores
- Alert notifications

## 🎮 Gaming Performance Considerations

### Cold Start Optimization
- Package size minimized for faster cold starts
- Dependencies optimized for gaming workloads
- Configuration cached for performance

### Gaming-Specific Monitoring
- EKS cluster version tracking
- Database performance class monitoring
- Cache node type validation
- Auto scaling capacity tracking

## 🔧 Troubleshooting

### Common Issues
1. **Permission Errors:** Verify IAM role has required permissions
2. **S3 Access:** Ensure snapshots bucket exists and is accessible
3. **SNS Failures:** Verify SNS topic ARN and permissions
4. **Slack Alerts:** Check webhook URL and network connectivity

### Gaming-Specific Issues
1. **False Positives:** Adjust gaming tag filters
2. **Performance Impact:** Review resource type changes
3. **Gaming Component Detection:** Update gaming identifier patterns

---
**Herald.lol Gaming Analytics Platform**  
*Infrastructure monitoring optimized for gaming performance*
EOF

    log_success "✅ Deployment documentation generated"
}

# Test package locally
test_package() {
    log_info "🧪 Testing gaming drift detection package..."
    
    # Extract to temporary directory for testing
    local test_dir="$BUILD_DIR/test"
    mkdir -p "$test_dir"
    cd "$test_dir"
    unzip -q "../../$PACKAGE_NAME"
    
    # Basic syntax check
    if python3 -m py_compile main.py; then
        log_success "✅ Python syntax validation passed"
    else
        log_warning "⚠️ Python syntax validation failed"
        return 1
    fi
    
    # Check imports
    if python3 -c "import main; print('✅ Import test passed')"; then
        log_success "✅ Import test passed"
    else
        log_warning "⚠️ Import test failed"
        return 1
    fi
    
    cd "$SCRIPT_DIR"
    rm -rf "$test_dir"
    
    log_success "✅ Gaming drift detection package testing completed"
}

# Main build function
main() {
    echo "🎮 Herald.lol Gaming Infrastructure Drift Detection - Lambda Build"
    echo "================================================================"
    echo "⚡ Performance Target: <5000ms gaming analytics"
    echo "👥 Concurrent Users: 1M+ supported"
    echo ""
    
    clean_build
    install_dependencies
    copy_function
    create_config
    optimize_package
    create_package
    test_package
    generate_docs
    
    echo ""
    log_success "🎉 Gaming drift detection Lambda package build completed!"
    echo ""
    echo "📦 Package: $SCRIPT_DIR/$PACKAGE_NAME"
    echo "📚 Docs: $SCRIPT_DIR/DEPLOYMENT.md"
    echo ""
    echo "🚀 Next Steps:"
    echo "1. Deploy using Terraform: cd terraform/drift-detection && terraform apply"
    echo "2. Or manually upload: aws lambda update-function-code ..."
    echo "3. Configure environment variables"
    echo "4. Test with gaming infrastructure"
    echo ""
    echo "🎮 Ready to monitor Herald.lol gaming infrastructure!"
}

# Run main function
main "$@"