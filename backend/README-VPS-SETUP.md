# 🎮 Herald.lol VPS Development Environment

## 📋 Overview

Complete VPS setup for Herald.lol Gaming Analytics Platform development environment. This implements Phase 1 Q1 2025 Infrastructure Foundation with all required gaming services.

## 🏗️ Architecture Stack

```
Herald.lol VPS Gaming Development Stack
├── 🐳 Docker Compose Orchestration
├── 🗄️ PostgreSQL 15 (Gaming-optimized)
├── ⚡ Redis 7 (Gaming cache + sessions)
├── 📊 InfluxDB 2.7 (Gaming time-series metrics)
├── 📈 Prometheus (Gaming monitoring)
├── 📊 Grafana (Gaming dashboards)
├── 🔍 ELK Stack (Gaming log analytics)
├── 🔒 HashiCorp Vault (Gaming secrets)
├── 🌐 NGINX (Gaming reverse proxy)
└── 🎮 Herald API (Gaming backend)
```

## 🚀 Quick Start

### Prerequisites Check

```bash
# Check current environment
./scripts/vps-setup.sh status

# View setup options
./scripts/vps-setup.sh --help
```

### Complete Setup

```bash
# Complete Herald.lol gaming environment setup
./scripts/vps-setup.sh setup

# This will:
# 1. Install Docker + Docker Compose (if needed)
# 2. Create gaming directory structure
# 3. Generate gaming configurations
# 4. Start all Herald.lol services
# 5. Show gaming environment status
```

### Manual Setup Steps

```bash
# 1. Create configuration files
./scripts/vps-setup.sh setup

# 2. Start gaming infrastructure
docker-compose -f docker-compose.herald-dev.yml up -d

# 3. Check gaming services status
./scripts/vps-setup.sh status

# 4. View gaming logs
./scripts/vps-setup.sh logs
```

## 🎯 Gaming Services

### Core Gaming Infrastructure

| Service | Port | Purpose | Gaming Focus |
|---------|------|---------|--------------|
| **Herald API** | 8080 | Gaming backend | LoL/TFT analytics |
| **PostgreSQL** | 5432 | Gaming database | Player/match data |
| **Redis** | 6379 | Gaming cache | Performance optimization |
| **InfluxDB** | 8086 | Gaming metrics | Time-series analytics |

### Gaming Monitoring Stack

| Service | Port | Credentials | Purpose |
|---------|------|-------------|---------|
| **Grafana** | 3001 | admin:herald_gaming_grafana_2025 | Gaming dashboards |
| **Prometheus** | 9090 | - | Gaming metrics collection |
| **Kibana** | 5601 | - | Gaming log visualization |
| **Elasticsearch** | 9200 | - | Gaming log storage |

### Gaming Security & Ops

| Service | Port | Credentials | Purpose |
|---------|------|-------------|---------|
| **Vault** | 8200 | token:herald-gaming-vault-root-token | Gaming secrets |
| **NGINX** | 80/443 | - | Gaming reverse proxy |

## 🎮 Gaming URLs

### Development Access

```bash
# Gaming API
curl http://localhost/api/health

# Gaming Dashboards
open http://localhost/grafana/
# Login: admin / herald_gaming_grafana_2025

# Gaming Logs  
open http://localhost/kibana/

# Gaming Metrics
open http://localhost:9090

# Gaming Vault
open http://localhost:8200
# Token: herald-gaming-vault-root-token
```

### Database Connections

```bash
# PostgreSQL Gaming Database
psql -h localhost -p 5432 -U herald_dev -d herald_gaming_dev
# Password: herald_gaming_dev_2025

# Redis Gaming Cache
redis-cli -h localhost -p 6379

# InfluxDB Gaming Metrics
influx -host localhost -port 8086
# Org: herald-gaming
# Token: herald-gaming-analytics-token-dev
```

## ⚙️ Configuration

### Environment Variables

```bash
# Create/edit .env file
cat > .env << 'EOF'
# Herald.lol Gaming Configuration
RIOT_API_KEY=RGAPI-your-riot-api-key-here
RIOT_REGION=euw1
GAMING_PERFORMANCE_TARGET_MS=5000
GAMING_MAX_CONCURRENT_USERS=1000000
EOF
```

### Gaming Database Schema

The setup automatically creates:
- **gaming** schema: Players, matches, champions
- **analytics** schema: Performance metrics, statistics
- **riot_api** schema: API data integration
- **Indexes**: Optimized for gaming queries (<100ms)
- **Views**: Pre-calculated gaming statistics

### Gaming Monitoring

Pre-configured gaming alerts:
- **Analytics Response Time** > 5000ms
- **Database Slow Queries** > 1000ms  
- **Memory Usage** > 90%
- **Gaming API Errors** > 5%

## 🔧 Development Workflow

### Daily Development

```bash
# Start gaming development day
./scripts/vps-setup.sh start

# Check gaming services
./scripts/vps-setup.sh status

# View gaming logs  
./scripts/vps-setup.sh logs herald-api

# Test gaming features
make test-gaming

# Stop at end of day
./scripts/vps-setup.sh stop
```

### Gaming Hot Reload

The Herald API service includes hot reload for development:
- **File changes** → Automatic rebuild
- **Code updates** → Instant restart
- **Config changes** → Auto-reload
- **Database changes** → Migration detection

### Gaming Performance Testing

```bash
# Run gaming benchmarks
make perf-test

# Validate gaming performance targets
make validate-performance

# Gaming load testing preparation
make load-test
```

## 📊 Gaming Monitoring

### Grafana Gaming Dashboards

Pre-configured gaming dashboards:
1. **Herald.lol Overview** - Gaming platform health
2. **Gaming Performance** - Response times, throughput  
3. **Player Analytics** - User behavior, engagement
4. **Riot API Integration** - Rate limits, errors
5. **Database Performance** - Query times, connections

### Gaming Metrics Collection

Automatic collection of gaming metrics:
- **API Performance**: Response times, error rates
- **Database**: Query performance, connections
- **Cache**: Hit rates, memory usage  
- **Gaming Logic**: KDA calculations, match processing
- **User Behavior**: Session times, feature usage

### Gaming Log Analysis

ELK stack configured for gaming logs:
- **Application Logs**: Herald API, services
- **Access Logs**: NGINX, API requests
- **Database Logs**: PostgreSQL slow queries
- **Gaming Events**: Match processing, player actions

## 🔒 Gaming Security

### Secrets Management

All gaming secrets managed via Vault:
```bash
# Set Riot API key
vault kv put secret/herald/riot api_key="RGAPI-your-key"

# Set database passwords  
vault kv put secret/herald/db password="secure-password"

# Gaming JWT secrets
vault kv put secret/herald/auth jwt_secret="gaming-jwt-secret"
```

### Gaming Data Protection

- **Player Data**: GDPR compliant encryption
- **API Keys**: Vault-managed rotation
- **Passwords**: Bcrypt hashing
- **Sessions**: Secure Redis storage
- **Communications**: TLS 1.3 encryption

## 🧹 Maintenance

### Gaming Environment Cleanup

```bash
# Stop gaming services
./scripts/vps-setup.sh stop

# Complete gaming cleanup (removes data)
./scripts/vps-setup.sh cleanup

# Docker gaming cleanup
docker system prune -f
```

### Gaming Backup Strategy

Automated backups configured:
- **Database**: Daily PostgreSQL dumps
- **Metrics**: InfluxDB snapshots  
- **Configurations**: Git versioning
- **Logs**: 30-day retention
- **Volumes**: Weekly snapshots

### Gaming Health Monitoring

Built-in health checks:
- **API Health**: /health endpoint
- **Database**: Connection pooling monitoring
- **Cache**: Redis ping checks
- **Services**: Docker health checks
- **Gaming Logic**: Performance validation

## 🔧 Troubleshooting

### Common Gaming Issues

```bash
# Gaming services won't start
./scripts/vps-setup.sh logs
docker-compose -f docker-compose.herald-dev.yml ps

# Gaming database connection issues
./scripts/vps-setup.sh logs herald-postgres
psql -h localhost -p 5432 -U herald_dev -d herald_gaming_dev

# Gaming API not responding
./scripts/vps-setup.sh logs herald-api
curl -v http://localhost/api/health

# Gaming performance issues
./scripts/vps-setup.sh logs herald-prometheus
open http://localhost:9090
```

### Gaming Performance Tuning

```bash
# Increase PostgreSQL gaming performance
# Edit: monitoring/prometheus/prometheus.yml
shared_buffers = '512MB'
effective_cache_size = '2GB'

# Redis gaming optimization  
# Edit: docker-compose.herald-dev.yml
maxmemory 1gb
maxmemory-policy allkeys-lru

# Gaming API performance
# Edit: .env
GAMING_PERFORMANCE_TARGET_MS=3000
```

## 📚 Next Steps

### Gaming Development Ready

After successful setup:

1. **Update Riot API Key** in .env file
2. **Test Gaming Features**: `make test-gaming`  
3. **Explore Gaming Dashboards**: http://localhost/grafana/
4. **Review Gaming Logs**: http://localhost/kibana/
5. **Start Gaming Development**: Begin Phase 1 features

### Phase 1 Q1 2025 Infrastructure ✅

Completed infrastructure foundation:
- [x] ✅ VPS development environment (Docker Compose)
- [x] ✅ Database architecture (PostgreSQL + InfluxDB)  
- [x] ✅ Monitoring stack (Prometheus + Grafana)
- [x] ✅ Logging infrastructure (ELK stack)
- [x] ✅ Secrets management (HashiCorp Vault)
- [x] ✅ Distributed caching (Redis cluster)
- [x] ✅ Reverse proxy (NGINX)

### Ready for Next Phase

Infrastructure prepared for:
- **Cloud Migration**: AWS/GCP deployment ready
- **Kubernetes**: Container orchestration ready  
- **Auto-scaling**: Performance monitoring ready
- **Production**: Security and monitoring ready

---

## 🎯 Herald.lol Gaming Performance Targets

> **Analytics Response**: <5000ms  
> **Concurrent Users**: 1M+ supported  
> **Uptime**: 99.9% target  
> **Gaming Focus**: League of Legends & TFT analytics  
> **Infrastructure**: Cloud-native, scalable, secure

**🎮 Herald.lol Gaming Development Environment Ready!**