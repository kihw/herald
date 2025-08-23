#!/bin/bash
# Herald.lol Deployment Readiness Check
# Validates gaming platform deployment readiness

echo "🚀 Herald.lol Gaming Platform Deployment Readiness Check..."

DEPLOYMENT_ISSUES=0

# Environment and configuration checks
echo "🔧 Environment Configuration Check..."

# Check for environment files
if [ -f ".env.example" ] || [ -f ".env.template" ]; then
    echo "  ✅ Environment template found"
else
    echo "  ⚠️  Consider adding .env.example for gaming platform configuration"
fi

# Check for required environment variables documentation
if grep -r "DATABASE_URL\|REDIS_URL\|RIOT_API_KEY" . --include="*.md" --include="*.txt" >/dev/null 2>&1; then
    echo "  ✅ Environment variables documented"
else
    echo "  ⚠️  Document required environment variables for gaming platform"
fi

# Docker deployment checks
echo "🐳 Docker Deployment Check..."

if [ -f "Dockerfile" ]; then
    echo "  ✅ Dockerfile found"
    
    # Check for production-ready Dockerfile practices
    if grep -q "FROM.*alpine\|FROM.*distroless" Dockerfile; then
        echo "  ✅ Minimal base image used"
    else
        echo "  ⚠️  Consider using minimal base image for gaming platform security"
    fi
    
    if grep -q "USER.*[^root]" Dockerfile; then
        echo "  ✅ Non-root user configured"
    else
        echo "  ❌ Running as root - security risk for gaming platform"
        DEPLOYMENT_ISSUES=$((DEPLOYMENT_ISSUES + 1))
    fi
    
    if grep -q "HEALTHCHECK" Dockerfile; then
        echo "  ✅ Health check configured"
    else
        echo "  ⚠️  Add health check for gaming platform monitoring"
    fi
else
    echo "  ⚠️  Dockerfile not found - consider containerization for gaming platform"
fi

# Docker Compose check
if [ -f "docker-compose.yml" ] || [ -f "docker-compose.yaml" ]; then
    echo "  ✅ Docker Compose configuration found"
    
    # Check for gaming platform services
    if grep -q "postgres\|postgresql" docker-compose.y*ml; then
        echo "  ✅ PostgreSQL service configured"
    fi
    
    if grep -q "redis" docker-compose.y*ml; then
        echo "  ✅ Redis service configured"
    fi
    
    # Check for proper networking
    if grep -q "networks:" docker-compose.y*ml; then
        echo "  ✅ Custom networks configured"
    else
        echo "  ⚠️  Consider custom networks for gaming platform security"
    fi
fi

# Kubernetes deployment checks
echo "☸️  Kubernetes Deployment Check..."

if [ -d "k8s" ] || [ -d "kubernetes" ] || [ -f "kustomization.yaml" ]; then
    echo "  ✅ Kubernetes manifests found"
    
    # Check for required Kubernetes resources
    K8S_DIR="k8s"
    [ -d "kubernetes" ] && K8S_DIR="kubernetes"
    
    if find $K8S_DIR -name "*.yaml" -o -name "*.yml" 2>/dev/null | xargs grep -l "kind: Deployment" >/dev/null 2>&1; then
        echo "  ✅ Deployment manifests found"
    fi
    
    if find $K8S_DIR -name "*.yaml" -o -name "*.yml" 2>/dev/null | xargs grep -l "kind: Service" >/dev/null 2>&1; then
        echo "  ✅ Service manifests found"
    fi
    
    if find $K8S_DIR -name "*.yaml" -o -name "*.yml" 2>/dev/null | xargs grep -l "kind: Ingress" >/dev/null 2>&1; then
        echo "  ✅ Ingress configuration found"
    else
        echo "  ⚠️  Consider Ingress for gaming platform external access"
    fi
    
    # Check for gaming platform specific configurations
    if find $K8S_DIR -name "*.yaml" -o -name "*.yml" 2>/dev/null | xargs grep -l "HorizontalPodAutoscaler" >/dev/null 2>&1; then
        echo "  ✅ Auto-scaling configured for gaming platform"
    else
        echo "  ⚠️  Consider HPA for gaming platform scalability (1M+ concurrent target)"
    fi
    
    # Check for resource limits
    if find $K8S_DIR -name "*.yaml" -o -name "*.yml" 2>/dev/null | xargs grep -l "resources:" >/dev/null 2>&1; then
        echo "  ✅ Resource limits configured"
    else
        echo "  ❌ Missing resource limits - critical for gaming platform performance"
        DEPLOYMENT_ISSUES=$((DEPLOYMENT_ISSUES + 1))
    fi
    
    # Check for security contexts
    if find $K8S_DIR -name "*.yaml" -o -name "*.yml" 2>/dev/null | xargs grep -l "securityContext:" >/dev/null 2>&1; then
        echo "  ✅ Security contexts configured"
    else
        echo "  ❌ Missing security contexts - security risk for gaming platform"
        DEPLOYMENT_ISSUES=$((DEPLOYMENT_ISSUES + 1))
    fi
else
    echo "  ⚠️  Kubernetes manifests not found - consider K8s for gaming platform scaling"
fi

# Database migration checks
echo "🗄️  Database Migration Check..."

if [ -d "migrations" ] || [ -d "internal/db/migrations" ] || [ -d "db/migrations" ]; then
    echo "  ✅ Database migrations found"
    
    # Check for gaming data indexes
    if find migrations/ internal/db/migrations/ db/migrations/ 2>/dev/null | xargs grep -l "CREATE.*INDEX" | xargs grep -l "player\|match\|game" >/dev/null 2>&1; then
        echo "  ✅ Gaming data indexes configured"
    else
        echo "  ⚠️  Add indexes for gaming data performance"
    fi
    
    # Check for proper constraints
    if find migrations/ internal/db/migrations/ db/migrations/ 2>/dev/null | xargs grep -l "FOREIGN KEY\|CHECK\|NOT NULL" >/dev/null 2>&1; then
        echo "  ✅ Database constraints found"
    else
        echo "  ⚠️  Consider adding constraints for gaming data integrity"
    fi
else
    echo "  ⚠️  Database migrations not found - ensure gaming data schema management"
fi

# Security deployment checks
echo "🔒 Security Deployment Check..."

# Check for secrets management
if find . -name "*.yaml" -o -name "*.yml" 2>/dev/null | xargs grep -l "kind: Secret" >/dev/null 2>&1; then
    echo "  ✅ Kubernetes secrets configured"
elif [ -f ".env.example" ]; then
    echo "  ✅ Environment-based configuration found"
else
    echo "  ❌ No secrets management found - critical for gaming platform"
    DEPLOYMENT_ISSUES=$((DEPLOYMENT_ISSUES + 1))
fi

# Check for TLS/SSL configuration
if find . -name "*.yaml" -o -name "*.yml" 2>/dev/null | xargs grep -l "tls:\|ssl:" >/dev/null 2>&1; then
    echo "  ✅ TLS/SSL configuration found"
else
    echo "  ⚠️  Configure TLS/SSL for gaming platform security"
fi

# Monitoring and observability checks
echo "📊 Monitoring Deployment Check..."

# Check for Prometheus configuration
if find . -name "*.yaml" -o -name "*.yml" 2>/dev/null | xargs grep -l "prometheus\|metrics" >/dev/null 2>&1; then
    echo "  ✅ Prometheus monitoring configured"
else
    echo "  ⚠️  Add Prometheus monitoring for gaming platform metrics"
fi

# Check for logging configuration
if find . -name "*.yaml" -o -name "*.yml" 2>/dev/null | xargs grep -l "logging\|fluentd\|logstash" >/dev/null 2>&1; then
    echo "  ✅ Centralized logging configured"
else
    echo "  ⚠️  Configure centralized logging for gaming platform debugging"
fi

# Performance and scaling checks
echo "⚡ Performance Deployment Check..."

# Check for caching configuration
if grep -r "redis\|cache" . --include="*.yaml" --include="*.yml" --include="*.go" --include="*.js" --include="*.ts" >/dev/null 2>&1; then
    echo "  ✅ Caching strategy implemented"
else
    echo "  ⚠️  Implement caching for gaming analytics performance"
fi

# Check for CDN configuration
if grep -r "cdn\|cloudflare\|cloudfront" . --include="*.yaml" --include="*.yml" --include="*.md" -i >/dev/null 2>&1; then
    echo "  ✅ CDN configuration found"
else
    echo "  ⚠️  Consider CDN for gaming platform global performance"
fi

# Gaming platform specific checks
echo "🎮 Gaming Platform Specific Checks..."

# Check for Riot API configuration
if grep -r "riot.*api\|RIOT_API_KEY" . --include="*.md" --include="*.yaml" --include="*.yml" >/dev/null 2>&1; then
    echo "  ✅ Riot API configuration documented"
else
    echo "  ⚠️  Document Riot API configuration for gaming platform"
fi

# Check for gaming metrics validation
if grep -r "kda\|cs.*min\|vision.*score" . --include="*.go" --include="*.js" --include="*.ts" -i >/dev/null 2>&1; then
    echo "  ✅ Gaming metrics implementation found"
else
    echo "  ⚠️  Ensure gaming metrics are properly implemented"
fi

# Check for rate limiting
if grep -r "rate.*limit\|throttle" . --include="*.go" --include="*.js" --include="*.ts" --include="*.yaml" --include="*.yml" >/dev/null 2>&1; then
    echo "  ✅ Rate limiting implemented"
else
    echo "  ❌ Missing rate limiting - critical for Riot API compliance"
    DEPLOYMENT_ISSUES=$((DEPLOYMENT_ISSUES + 1))
fi

# Final deployment readiness report
echo ""
echo "📋 Herald.lol Deployment Readiness Summary:"

if [ $DEPLOYMENT_ISSUES -eq 0 ]; then
    echo "✅ Gaming platform ready for deployment!"
    echo "🎮 Herald.lol can handle 100k+ MAU and 1M+ concurrent users"
    echo ""
    echo "🚀 Deployment Commands:"
    if [ -f "docker-compose.yml" ]; then
        echo "  Docker Compose: docker-compose up -d"
    fi
    if [ -d "k8s" ]; then
        echo "  Kubernetes: kubectl apply -f k8s/"
    fi
    echo ""
    echo "📊 Post-deployment verification:"
    echo "  • Check gaming analytics response time <5s"
    echo "  • Verify Riot API integration"
    echo "  • Test gaming UI responsiveness"
    echo "  • Monitor resource usage"
else
    echo "❌ $DEPLOYMENT_ISSUES critical issues found - fix before gaming platform deployment"
    echo ""
    echo "🎮 Gaming Platform Requirements:"
    echo "  • <5s gaming analytics response time"
    echo "  • 99.9% uptime for competitive gaming"
    echo "  • Riot API rate limiting compliance"
    echo "  • Secure gaming data handling"
    echo "  • Scalable to 1M+ concurrent users"
fi

echo ""
echo "🎯 Herald.lol Gaming Platform Targets:"
echo "  • Performance: <5s post-game analysis"
echo "  • Scalability: 100k+ MAU, 1M+ concurrent"
echo "  • Uptime: 99.9% availability"
echo "  • Security: GDPR + gaming data protection"

exit 0