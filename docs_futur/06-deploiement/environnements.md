# Environnements et Stratégie de Déploiement

## Vue d'Ensemble des Environnements

Herald.lol adopte une **stratégie d'environnements multi-niveaux** qui garantit la qualité, la stabilité et la sécurité à travers tout le cycle de développement. Cette approche permet un développement agile tout en maintenant des standards de production rigoureux.

## Architecture Multi-Environnements

### Environnement de Développement (Development)

#### Configuration Locale
- **Infrastructure** : Docker Compose local avec services essentiels
- **Base de Données** : SQLite pour rapidité et PostgreSQL pour tests
- **Cache** : Redis local single-node
- **Storage** : Stockage local avec MinIO S3-compatible

#### Services de Développement
```yaml
# docker-compose.dev.yml
version: '3.8'
services:
  herald-backend:
    build: .
    environment:
      - ENV=development
      - DB_TYPE=sqlite
      - REDIS_URL=redis://redis:6379
      - RIOT_API_KEY=${RIOT_API_KEY_DEV}
    volumes:
      - ./data:/app/data
      
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
      
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: herald_dev
    volumes:
      - postgres_data:/var/lib/postgresql/data
```

#### Caractéristiques Développement
- **Hot Reloading** : Rechargement automatique frontend et backend
- **Debug Mode** : Logs détaillés et profiling activé
- **Mock APIs** : Services externes mockés pour tests
- **Seed Data** : Données de test pré-chargées

### Environnement de Test (Testing/Staging)

#### Infrastructure Cloud Staging
- **Platform** : AWS avec configuration similaire production
- **Scaling** : Version réduite de l'infrastructure production
- **Monitoring** : Monitoring complet avec alerting
- **Security** : Configuration sécurité identique production

#### Configuration Staging
```yaml
# terraform/environments/staging/main.tf
module "herald_staging" {
  source = "../../modules/herald"
  
  environment = "staging"
  instance_count = 2
  db_instance_class = "db.t3.medium"
  redis_node_type = "cache.t3.medium"
  
  # Reduced capacity for cost optimization
  min_capacity = 1
  max_capacity = 5
  
  tags = {
    Environment = "staging"
    Purpose = "testing"
  }
}
```

#### Testing Features
- **Automated Testing** : Suite complète tests automatisés
- **Performance Testing** : Tests charge et stress
- **Security Testing** : Scans vulnérabilités automatisés
- **Integration Testing** : Tests intégration APIs externes

### Environnement de Production (Production)

#### Infrastructure Production Haute Disponibilité
- **Multi-AZ Deployment** : Déploiement multi-zones disponibilité
- **Auto-Scaling** : Auto-scaling basé métriques performance
- **Load Balancing** : Load balancers avec health checks
- **CDN Global** : Distribution contenu mondiale

#### Configuration Production
```yaml
# terraform/environments/production/main.tf
module "herald_production" {
  source = "../../modules/herald"
  
  environment = "production"
  instance_count = 6
  db_instance_class = "db.r5.2xlarge"
  redis_node_type = "cache.r5.2xlarge"
  
  # High availability configuration
  multi_az = true
  backup_retention = 30
  
  # Auto-scaling configuration
  min_capacity = 3
  max_capacity = 20
  target_cpu_utilization = 70
  
  tags = {
    Environment = "production"
    Purpose = "live-service"
  }
}
```

#### Production Features
- **99.9% Uptime SLA** : Garantie disponibilité service
- **Real-Time Monitoring** : Monitoring temps réel 24/7
- **Disaster Recovery** : Plan reprise activité complet
- **Security Hardening** : Durcissement sécurité maximal

## Infrastructure as Code (IaC)

### Terraform Infrastructure Management

#### Module Architecture
```hcl
# modules/herald/main.tf
module "networking" {
  source = "./modules/networking"
  environment = var.environment
  vpc_cidr = var.vpc_cidr
}

module "compute" {
  source = "./modules/compute"
  environment = var.environment
  vpc_id = module.networking.vpc_id
  subnet_ids = module.networking.private_subnet_ids
}

module "database" {
  source = "./modules/database"
  environment = var.environment
  vpc_id = module.networking.vpc_id
  subnet_ids = module.networking.db_subnet_ids
}

module "cache" {
  source = "./modules/cache"
  environment = var.environment
  vpc_id = module.networking.vpc_id
  subnet_ids = module.networking.cache_subnet_ids
}
```

#### Environment-Specific Variables
```hcl
# environments/production/terraform.tfvars
environment = "production"
region = "eu-west-1"

# Networking
vpc_cidr = "10.0.0.0/16"
availability_zones = ["eu-west-1a", "eu-west-1b", "eu-west-1c"]

# Compute
instance_type = "c5.2xlarge"
min_size = 3
max_size = 20
desired_capacity = 6

# Database
db_instance_class = "db.r5.2xlarge"
db_allocated_storage = 1000
backup_retention_period = 30
```

### Kubernetes Deployment

#### Namespace Isolation
```yaml
# k8s/namespaces/production.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: herald-production
  labels:
    environment: production
    app: herald
---
apiVersion: v1
kind: ResourceQuota
metadata:
  name: herald-quota
  namespace: herald-production
spec:
  hard:
    requests.cpu: "20"
    requests.memory: 40Gi
    limits.cpu: "40"
    limits.memory: 80Gi
```

#### Application Deployment
```yaml
# k8s/apps/herald-backend.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: herald-backend
  namespace: herald-production
spec:
  replicas: 6
  selector:
    matchLabels:
      app: herald-backend
  template:
    metadata:
      labels:
        app: herald-backend
    spec:
      containers:
      - name: backend
        image: herald/backend:v2.1.0
        ports:
        - containerPort: 8000
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: herald-secrets
              key: database-url
        resources:
          requests:
            memory: "2Gi"
            cpu: "1"
          limits:
            memory: "4Gi"
            cpu: "2"
        livenessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8000
          initialDelaySeconds: 5
          periodSeconds: 5
```

## CI/CD Pipeline

### GitHub Actions Workflow

#### Multi-Environment Pipeline
```yaml
# .github/workflows/deploy.yml
name: Deploy Herald.lol

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v4
      with:
        go-version: '1.23'
    - name: Run Tests
      run: |
        go test ./...
        cd web && npm test
    
  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Build Docker Images
      run: |
        docker build -t herald/backend:${{ github.sha }} .
        docker build -t herald/frontend:${{ github.sha }} ./web
    
  deploy-staging:
    needs: build
    if: github.ref == 'refs/heads/develop'
    runs-on: ubuntu-latest
    environment: staging
    steps:
    - name: Deploy to Staging
      run: |
        kubectl apply -f k8s/staging/
        kubectl set image deployment/herald-backend herald-backend=herald/backend:${{ github.sha }}
    
  deploy-production:
    needs: build
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment: production
    steps:
    - name: Deploy to Production
      run: |
        kubectl apply -f k8s/production/
        kubectl rollout restart deployment/herald-backend
```

#### Environment Protection Rules
- **Staging** : Déploiement automatique sur develop branch
- **Production** : Approbation manuelle requise + reviews
- **Rollback** : Capacité rollback automatique si health checks échouent
- **Blue-Green** : Déploiement blue-green pour zero-downtime

### ArgoCD GitOps

#### Application Configuration
```yaml
# argocd/herald-production.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: herald-production
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/herald/kubernetes-manifests
    targetRevision: main
    path: production
  destination:
    server: https://kubernetes.default.svc
    namespace: herald-production
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m
```

#### Progressive Delivery
- **Canary Deployments** : Déploiement canary avec Argo Rollouts
- **Blue-Green Strategy** : Stratégie blue-green pour updates majeures
- **Traffic Splitting** : Division trafic progressive 10% -> 50% -> 100%
- **Automatic Rollback** : Rollback automatique si métriques dégradées

## Configuration Management

### Environment-Specific Configuration

#### ConfigMaps Kubernetes
```yaml
# k8s/config/production-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: herald-config
  namespace: herald-production
data:
  app.yaml: |
    server:
      port: 8000
      read_timeout: 30s
      write_timeout: 30s
    
    database:
      max_connections: 100
      connection_timeout: 30s
    
    cache:
      default_ttl: 3600
      max_memory: 2gb
    
    analytics:
      batch_size: 1000
      processing_interval: 60s
```

#### Secrets Management
```yaml
# k8s/secrets/production-secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: herald-secrets
  namespace: herald-production
type: Opaque
data:
  database-url: <base64-encoded-database-url>
  redis-url: <base64-encoded-redis-url>
  riot-api-key: <base64-encoded-riot-api-key>
  jwt-secret: <base64-encoded-jwt-secret>
```

### External Secrets Operator

#### AWS Secrets Manager Integration
```yaml
# k8s/external-secrets/production.yaml
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: aws-secrets-manager
  namespace: herald-production
spec:
  provider:
    aws:
      service: SecretsManager
      region: eu-west-1
      auth:
        jwt:
          serviceAccountRef:
            name: external-secrets-sa
---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: herald-external-secrets
  namespace: herald-production
spec:
  refreshInterval: 15s
  secretStoreRef:
    name: aws-secrets-manager
    kind: SecretStore
  target:
    name: herald-secrets
    creationPolicy: Owner
  data:
  - secretKey: riot-api-key
    remoteRef:
      key: herald/production/riot-api-key
```

## Monitoring et Observability par Environnement

### Environment-Specific Monitoring

#### Production Monitoring Stack
```yaml
# monitoring/production/prometheus.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s
      evaluation_interval: 15s
    
    rule_files:
    - "/etc/prometheus/rules/*.yml"
    
    scrape_configs:
    - job_name: 'herald-backend'
      kubernetes_sd_configs:
      - role: endpoints
        namespaces:
          names: ['herald-production']
      relabel_configs:
      - source_labels: [__meta_kubernetes_service_name]
        action: keep
        regex: herald-backend
    
    - job_name: 'herald-frontend'
      kubernetes_sd_configs:
      - role: endpoints
        namespaces:
          names: ['herald-production']
```

#### Alerting Rules Production
```yaml
# monitoring/alerts/production.yml
groups:
- name: herald.production
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
    for: 5m
    labels:
      severity: critical
      environment: production
    annotations:
      summary: "High error rate detected"
      
  - alert: HighLatency
    expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
    for: 5m
    labels:
      severity: warning
      environment: production
```

### Environment Health Checks

#### Comprehensive Health Endpoints
```go
// health/checker.go
type HealthChecker struct {
    db    *sql.DB
    redis *redis.Client
    env   string
}

func (h *HealthChecker) CheckHealth() HealthStatus {
    status := HealthStatus{
        Environment: h.env,
        Timestamp:   time.Now(),
        Status:      "healthy",
        Checks:      make(map[string]CheckResult),
    }
    
    // Database health
    if err := h.db.Ping(); err != nil {
        status.Checks["database"] = CheckResult{
            Status: "unhealthy",
            Error:  err.Error(),
        }
        status.Status = "degraded"
    }
    
    // Redis health
    if err := h.redis.Ping().Err(); err != nil {
        status.Checks["redis"] = CheckResult{
            Status: "unhealthy", 
            Error:  err.Error(),
        }
        status.Status = "degraded"
    }
    
    return status
}
```

## Security par Environnement

### Network Security

#### VPC et Network ACLs
```hcl
# terraform/modules/networking/security.tf
resource "aws_security_group" "herald_backend" {
  name_prefix = "${var.environment}-herald-backend"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 8000
    to_port     = 8000
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.main.cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${var.environment}-herald-backend-sg"
    Environment = var.environment
  }
}
```

#### Pod Security Standards
```yaml
# k8s/security/pod-security-policy.yaml
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: herald-restricted
spec:
  privileged: false
  allowPrivilegeEscalation: false
  requiredDropCapabilities:
  - ALL
  volumes:
  - 'configMap'
  - 'emptyDir'
  - 'projected'
  - 'secret'
  - 'downwardAPI'
  - 'persistentVolumeClaim'
  runAsUser:
    rule: 'MustRunAsNonRoot'
  seLinux:
    rule: 'RunAsAny'
  fsGroup:
    rule: 'RunAsAny'
```

Cette stratégie d'environnements robuste garantit que Herald.lol maintient la plus haute qualité et sécurité à travers tous les stades de développement et déploiement.