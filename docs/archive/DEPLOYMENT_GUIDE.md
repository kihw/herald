# 🚀 LoL Match Exporter - Production Deployment Guide

## 📋 Overview

This guide covers the complete deployment process for the optimized LoL Match Exporter with Go native analytics, Redis caching, and goroutine-based worker pools.

## 🏗️ Production Architecture

```
Internet → Load Balancer → Go Application Servers → Redis Cluster
                                    ↓
                              PostgreSQL Database
                                    ↓
                              Monitoring Stack
```

## 📦 Pre-Production Checklist

### ✅ Infrastructure Requirements

**Minimum Requirements:**
- **CPU**: 4 cores (8 recommended for high traffic)
- **RAM**: 8GB (16GB recommended)
- **Storage**: 100GB SSD (for logs and temporary data)
- **Network**: 1Gbps bandwidth

**Software Dependencies:**
- **Go**: 1.21+ (for compilation)
- **Redis**: 7.0+ (for caching)
- **PostgreSQL**: 14+ (for persistent data)
- **Docker**: 20.10+ (for containerized deployment)

### ✅ Environment Setup

**1. Environment Variables**
```bash
# Core Application
PORT=8001
GIN_MODE=release
PROJECT_DIR=/app

# Database Configuration
DB_HOST=postgres.production.internal
DB_PORT=5432
DB_NAME=lol_match_exporter
DB_USER=app_user
DB_PASSWORD=secure_password
DB_SSL_MODE=require

# Redis Configuration
REDIS_HOST=redis.production.internal
REDIS_PORT=6379
REDIS_PASSWORD=redis_secure_password
REDIS_DB=0

# Analytics Configuration
ANALYTICS_CACHE_ENABLED=true
ANALYTICS_MAX_WORKERS=8
ANALYTICS_QUEUE_SIZE=1000
ANALYTICS_QUERY_TIMEOUT=30s

# Cache TTL Settings
CACHE_SHORT_TTL=5m
CACHE_MEDIUM_TTL=1h
CACHE_LONG_TTL=24h
CACHE_VERY_LONG_TTL=168h

# Monitoring
METRICS_ENABLED=true
METRICS_PORT=9090
LOG_LEVEL=info
```

**2. Build Configuration**
```bash
# Production build with optimizations
go build -ldflags="-s -w" -o analytics-server ./cmd/analytics-server

# Create minimal Docker image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY analytics-server .
COPY web/dist ./web/dist
EXPOSE 8001
CMD ["./analytics-server"]
```

## 🐳 Docker Deployment

### Docker Compose Production

```yaml
version: '3.8'

services:
  analytics-server:
    build: .
    ports:
      - "8001:8001"
    environment:
      - GIN_MODE=release
      - REDIS_HOST=redis
      - DB_HOST=postgres
    depends_on:
      - redis
      - postgres
    deploy:
      replicas: 3
      resources:
        limits:
          memory: 2G
          cpus: '1.0'
        reservations:
          memory: 1G
          cpus: '0.5'
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8001/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    deploy:
      resources:
        limits:
          memory: 4G
          cpus: '2.0'

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    deploy:
      resources:
        limits:
          memory: 8G
          cpus: '4.0'

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/ssl/certs
    depends_on:
      - analytics-server

volumes:
  redis_data:
  postgres_data:
```

## ⚙️ Configuration Tuning

### Redis Optimization

```redis
# redis.conf production settings
maxmemory 4gb
maxmemory-policy allkeys-lru
save 900 1
save 300 10
save 60 10000
tcp-keepalive 60
timeout 300
```

### Go Application Tuning

```go
// Production configuration
config := services.OptimizedConfig{
    // Cache settings
    CacheEnabled:     true,
    CacheHost:        os.Getenv("REDIS_HOST"),
    CachePort:        6379,
    CachePassword:    os.Getenv("REDIS_PASSWORD"),
    
    // Worker pool optimization
    EnableAsyncProcessing:   true,
    MaxWorkers:             8,  // 2x CPU cores
    QueueSize:              1000,
    EnableConcurrentQueries: true,
    QueryTimeout:           30 * time.Second,
    
    // Cache TTL optimization
    ShortCacheTTL:    5 * time.Minute,    // Real-time data
    MediumCacheTTL:   1 * time.Hour,      // Analytics results
    LongCacheTTL:     24 * time.Hour,     // User profiles
    VeryLongCacheTTL: 7 * 24 * time.Hour, // Historical data
}
```

## 📊 Performance Monitoring

### Health Check Endpoints

```bash
# Application health
curl http://localhost:8001/api/health

# Optimized analytics health
curl http://localhost:8001/api/analytics/v2/health

# Performance metrics
curl http://localhost:8001/api/analytics/v2/performance
```

### Monitoring Stack Integration

**Prometheus Configuration (prometheus.yml)**
```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'lol-analytics'
    static_configs:
      - targets: ['analytics-server:9090']
    metrics_path: /metrics
    scrape_interval: 10s
```

**Grafana Dashboard Metrics**
- Request rate and latency
- Cache hit/miss ratios
- Worker pool utilization
- Memory and CPU usage
- Error rates and response codes

## 🔒 Security Configuration

### Application Security

```go
// Production security middleware
r.Use(gin.Recovery())
r.Use(middleware.RateLimiter(1000, time.Minute)) // 1000 req/min
r.Use(middleware.CORS(production_origins))
r.Use(middleware.RequestID())
r.Use(middleware.Logging())
```

### Network Security

```nginx
# nginx.conf security headers
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Strict-Transport-Security "max-age=31536000" always;

# Rate limiting
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
limit_req zone=api burst=20 nodelay;
```

## 🚀 Deployment Steps

### 1. Pre-deployment Preparation

```bash
# Clone and prepare
git clone <repository>
cd lol_match_exporter

# Set up environment
cp .env.production .env
# Edit .env with production values

# Build application
make build-production
```

### 2. Database Migration

```bash
# Run database migrations
./analytics-server migrate up

# Verify database schema
./analytics-server migrate status
```

### 3. Service Deployment

```bash
# Deploy with Docker Compose
docker-compose -f docker-compose.prod.yml up -d

# Verify deployment
docker-compose ps
docker-compose logs analytics-server
```

### 4. Post-deployment Validation

```bash
# Health checks
curl http://localhost:8001/api/health
curl http://localhost:8001/api/analytics/v2/health

# Performance test
./scripts/performance-test.sh

# Load test
./scripts/load-test.sh
```

## 📈 Performance Benchmarks

### Expected Performance Metrics

**Single Instance (4 CPU, 8GB RAM):**
- **Throughput**: 1000+ requests/second
- **Latency**: <100ms p95 for cached responses
- **Cache Hit Rate**: >85% for repeated queries
- **Worker Utilization**: 70-80% under normal load

**Scaled Deployment (3 instances):**
- **Throughput**: 3000+ requests/second
- **High Availability**: 99.9% uptime
- **Fault Tolerance**: Service degradation vs failure

### Load Testing Commands

```bash
# Basic load test
ab -n 10000 -c 100 http://localhost:8001/api/analytics/v2/health

# Analytics endpoint test
ab -n 1000 -c 50 -H "Authorization: Bearer <token>" \
   http://localhost:8001/api/analytics/v2/period/week

# Batch processing test
ab -n 500 -c 25 -p batch_request.json -T application/json \
   http://localhost:8001/api/analytics/v2/batch
```

## 🔧 Troubleshooting

### Common Issues

**High Memory Usage**
```bash
# Check Go memory stats
curl http://localhost:8001/debug/pprof/heap

# Redis memory analysis
redis-cli info memory
redis-cli --bigkeys
```

**Cache Performance Issues**
```bash
# Redis performance monitoring
redis-cli monitor
redis-cli info stats

# Application cache metrics
curl http://localhost:8001/api/analytics/v2/performance
```

**Worker Pool Saturation**
```bash
# Check worker pool stats
curl http://localhost:8001/api/analytics/v2/performance | jq '.worker_pool'

# Adjust worker count
export ANALYTICS_MAX_WORKERS=12
```

## 📋 Maintenance Procedures

### Regular Maintenance

**Daily:**
- Monitor error logs
- Check performance metrics
- Verify cache hit rates

**Weekly:**
- Database maintenance and backups
- Redis data cleanup
- Performance baseline review

**Monthly:**
- Security updates
- Capacity planning review
- Performance optimization

### Backup Procedures

```bash
# Database backup
pg_dump $DB_NAME > backup_$(date +%Y%m%d).sql

# Redis backup
redis-cli --rdb redis_backup_$(date +%Y%m%d).rdb

# Application configuration backup
tar -czf config_backup_$(date +%Y%m%d).tar.gz .env nginx.conf
```

## 🎯 Production Readiness Checklist

- [ ] Environment variables configured
- [ ] Database migrations applied
- [ ] Redis cluster deployed and tested
- [ ] Load balancer configured
- [ ] SSL certificates installed
- [ ] Monitoring stack deployed
- [ ] Backup procedures tested
- [ ] Disaster recovery plan documented
- [ ] Performance baselines established
- [ ] Security audit completed

---

🚀 **This deployment guide ensures a production-ready, scalable, and maintainable LoL Match Exporter deployment with enterprise-grade performance and reliability.**