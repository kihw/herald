# Herald.lol Gaming Analytics - Infrastructure Documentation

🎮 **Production-ready AWS infrastructure for Herald gaming analytics platform**

## 🏗️ Architecture Overview

Herald.lol utilizes a modern, cloud-native architecture optimized for gaming analytics workloads:

- **Compute**: Amazon EKS with auto-scaling node groups
- **Database**: Aurora PostgreSQL cluster with read replicas across 3 AZs
- **Cache**: ElastiCache Redis cluster for gaming sessions
- **Storage**: S3 for assets and backups
- **CDN**: CloudFront for global content delivery
- **Service Mesh**: Istio for secure service communication
- **Monitoring**: Prometheus + Grafana + Jaeger tracing

## 🚀 Quick Start

### Prerequisites

Ensure you have the following tools installed:

```bash
# Required tools
aws-cli (v2.x)
terraform (v1.0+)
kubectl (v1.28+)
helm (v3.x)
istioctl (v1.19+)
```

### AWS Credentials Setup

```bash
# Configure AWS credentials
aws configure

# Verify access
aws sts get-caller-identity
```

### Deploy Infrastructure

```bash
# Deploy production infrastructure
./scripts/deploy-infrastructure.sh -e production -r us-east-1

# Dry run to see what will be created
./scripts/deploy-infrastructure.sh --dry-run

# Deploy with auto-confirmation
./scripts/deploy-infrastructure.sh -y
```

## 📋 Infrastructure Components

### 🔧 Core Services

#### EKS Cluster Configuration
- **Cluster Name**: `herald-gaming-cluster`
- **Version**: Kubernetes 1.28
- **Node Groups**: 
  - Gaming Primary: c5.2xlarge (3-100 nodes)
  - Analytics Processing: r5.4xlarge (0-20 nodes, Spot)
- **Add-ons**: CoreDNS, VPC-CNI, EBS CSI Driver
- **Networking**: Custom VPC with private/public subnets

#### Database Architecture
- **Primary**: Aurora PostgreSQL 15.4
- **Configuration**: Multi-AZ cluster with 3 instances
- **Performance**: r6g.2xlarge primary, r6g.xlarge readers
- **Backup**: 30-day retention, point-in-time recovery
- **Security**: Encryption at rest, VPC isolation

#### Redis Cache Cluster
- **Type**: ElastiCache Redis 7.x in cluster mode
- **Configuration**: 3 node groups, 1 replica each
- **Instance**: cache.r7g.2xlarge
- **Features**: Encryption, automatic failover, backup

### 🌐 Networking & Security

#### VPC Configuration
```
VPC CIDR: 10.0.0.0/16
├── Public Subnets:  10.0.101.0/24, 10.0.102.0/24, 10.0.103.0/24
└── Private Subnets: 10.0.1.0/24,   10.0.2.0/24,   10.0.3.0/24
```

#### Security Groups
- **RDS**: Port 5432 from EKS nodes only
- **ElastiCache**: Port 6379 from EKS nodes only
- **EKS**: Managed security groups with gaming optimizations

#### Istio Service Mesh
- **mTLS**: Strict mode for all gaming services
- **Circuit Breaker**: Gaming-optimized thresholds
- **Load Balancing**: LEAST_CONN for consistent latency
- **Tracing**: Jaeger integration for request tracking

## 🎮 Gaming-Specific Optimizations

### Performance Tuning

#### Network Optimizations
```bash
# Applied to all gaming nodes
net.core.rmem_max = 134217728
net.core.wmem_max = 134217728  
net.ipv4.tcp_rmem = 4096 87380 134217728
net.ipv4.tcp_wmem = 4096 65536 134217728
```

#### Database Parameters
- **max_connections**: 2000 (high concurrency)
- **effective_cache_size**: 24GB
- **work_mem**: 64MB per operation
- **random_page_cost**: 1.1 (SSD optimized)

#### Redis Configuration
- **Cluster Mode**: Enabled for horizontal scaling
- **Connection Pool**: 100 max connections per service
- **Memory Policy**: allkeys-lru for gaming data

### Gaming Workload Labels
```yaml
# Node selectors for gaming workloads
WorkloadType: gaming-analytics
NodeType: high-memory
gaming-platform: herald
environment: production
```

## 📊 Monitoring & Observability

### Prometheus Metrics
- **Infrastructure**: Node, pod, service metrics
- **Gaming**: Custom metrics for KDA, CS/min, Vision Score  
- **Performance**: Request latency, throughput, error rates

### Grafana Dashboards
- **Gaming Analytics**: Real-time gaming performance metrics
- **Infrastructure Health**: Cluster, database, cache status
- **Application Performance**: Service response times, errors

### Jaeger Tracing  
- **Service Map**: Visual representation of service calls
- **Request Tracing**: End-to-end request tracking
- **Performance Analysis**: Bottleneck identification

## 🔒 Security & Compliance

### Encryption
- **At Rest**: All storage encrypted (RDS, ElastiCache, S3, EBS)
- **In Transit**: TLS 1.3 for all communications
- **Service Mesh**: mTLS between all gaming services

### Access Control
- **RBAC**: Role-based access for Kubernetes resources
- **IAM**: Least privilege AWS permissions
- **Network**: Private subnets, security group restrictions

### Compliance
- **GDPR**: Data protection for EU gaming data
- **Riot ToS**: Compliant API usage and data handling
- **Audit**: CloudTrail logging, VPC Flow Logs

## 🚀 Scaling & Performance

### Auto-Scaling Configuration
```yaml
# Gaming primary nodes
minSize: 3
maxSize: 100
targetCPU: 70%
targetMemory: 80%

# Analytics processing nodes  
minSize: 0
maxSize: 20
scaleUpCooldown: 60s
scaleDownCooldown: 300s
```

### Performance Targets
- **API Response**: <500ms average
- **Database Queries**: <100ms average
- **Cache Access**: <1ms average
- **Gaming Analytics**: <5s post-game analysis

### Load Testing
```bash
# Performance testing with k6
k6 run --vus 1000 --duration 10m scripts/load-test-gaming.js

# Database stress testing
pgbench -h ${RDS_ENDPOINT} -U herald_admin -c 100 -j 4 -T 300
```

## 🛠️ Operational Procedures

### Daily Operations

#### Health Checks
```bash
# Check cluster health
kubectl get nodes
kubectl get pods --all-namespaces

# Check database health
aws rds describe-db-clusters --db-cluster-identifier herald-gaming-cluster

# Check Redis health
aws elasticache describe-replication-groups --replication-group-id herald-gaming-redis
```

#### Monitoring Alerts
- **Node CPU**: >80% for 5 minutes
- **Memory Usage**: >85% for 3 minutes
- **Database Connections**: >1800 active connections
- **Gaming API Errors**: >5% error rate for 2 minutes

### Backup & Recovery

#### Database Backups
- **Automated**: Daily backups with 30-day retention
- **Point-in-Time**: Recovery to any second within retention
- **Cross-Region**: Backup replication for disaster recovery

#### Disaster Recovery
```bash
# RDS failover (automated)
aws rds failover-db-cluster --db-cluster-identifier herald-gaming-cluster

# EKS node replacement (automated via ASG)
# Pods automatically rescheduled on healthy nodes
```

## 🔧 Troubleshooting

### Common Issues

#### Gaming Service Connectivity
```bash
# Check service mesh status  
istioctl proxy-status
kubectl get vs,dr,gw -n herald-gaming

# Test service-to-service communication
kubectl exec -n herald-gaming deployment/gaming-api -- curl -v http://analytics-service:8080/health
```

#### Database Performance
```bash
# Check active queries
psql -h ${RDS_ENDPOINT} -U herald_admin -c "SELECT * FROM pg_stat_activity WHERE state = 'active';"

# Monitor connection pool
kubectl logs -n herald-gaming deployment/gaming-api | grep "pool"
```

#### Cache Issues
```bash
# Redis cluster info
redis-cli -h ${REDIS_ENDPOINT} cluster info
redis-cli -h ${REDIS_ENDPOINT} cluster nodes

# Check cache hit ratio
redis-cli -h ${REDIS_ENDPOINT} info stats | grep cache_hit
```

### Performance Optimization

#### Gaming Latency Issues
1. Check Istio circuit breaker status
2. Verify node resource availability  
3. Review database query performance
4. Analyze cache hit ratios

#### High CPU Usage
1. Scale node groups horizontally
2. Check for memory leaks in gaming services
3. Review analytics processing efficiency
4. Optimize database queries

## 📚 Additional Resources

### Documentation Links
- [AWS EKS Best Practices](https://aws.github.io/aws-eks-best-practices/)
- [Istio Gaming Configuration](https://istio.io/latest/docs/tasks/traffic-management/)
- [PostgreSQL Gaming Optimization](https://wiki.postgresql.org/wiki/Performance_Optimization)
- [Redis Gaming Cache Patterns](https://redis.io/docs/manual/patterns/)

### Gaming-Specific Resources
- [Riot Games API Documentation](https://developer.riotgames.com/)
- [League of Legends Data Dragon](https://developer.riotgames.com/docs/lol#data-dragon)
- [Gaming Analytics Best Practices](https://www.herald.lol/docs/analytics)

### Support Contacts
- **Infrastructure**: infrastructure@herald.lol
- **Gaming Services**: gaming@herald.lol  
- **Security**: security@herald.lol
- **Emergency**: +1-XXX-XXX-XXXX (24/7 on-call)

---

## 🎮 Gaming Analytics Architecture Summary

Herald.lol provides a robust, scalable infrastructure specifically optimized for gaming analytics workloads:

- **🏎️ Performance**: Sub-5-second gaming analytics processing
- **🌍 Global Scale**: Multi-region deployment with edge optimization
- **🔒 Security**: End-to-end encryption with gaming data compliance
- **📊 Observability**: Comprehensive monitoring and tracing
- **⚡ Reliability**: 99.9% uptime with automated failover
- **🎯 Gaming-First**: Optimized for League of Legends and TFT analytics

**Ready to revolutionize gaming analytics! 🎮**