# 🎮 Herald.lol Gaming Analytics - API Gateway

Kong API Gateway configuration optimized for Herald.lol gaming analytics platform with focus on performance, scalability, and gaming-specific requirements.

## 🎯 Gaming Performance Targets

- **Analytics Load Time:** <5 seconds
- **UI Response Time:** <2 seconds  
- **System Uptime:** 99.9%
- **Concurrent Users:** 1M+ support
- **API Response Time:** <1 second

## 🏗️ Architecture Overview

```
Gaming Clients → Kong API Gateway → Herald.lol Backend Services
     ↓               ↓                        ↓
  Rate Limiting   Authentication         Gaming Analytics
  Load Balancing    Monitoring            Riot API Integration  
  Caching          Security               Match Analysis
  Transformations  Logging                Team Composition
```

## 🚀 Quick Start

### 1. Environment Setup

```bash
# Copy environment configuration
cp gateway/.env.kong.example gateway/.env.kong

# Update configuration for your environment
nano gateway/.env.kong
```

### 2. Start Gaming API Gateway

```bash
# Start Kong with gaming services
docker-compose -f gateway/docker-compose.kong.yml up -d

# Run gaming setup script
./gateway/scripts/kong-setup.sh
```

### 3. Verify Gaming Setup

```bash
# Check Kong admin
curl http://localhost:8001

# Check gaming endpoints
curl http://localhost:8000/api/analytics/health
curl http://localhost:8000/api/riot/health
```

## 🎮 Gaming Services Configuration

### Gaming Analytics Service
- **Endpoint:** `/api/analytics`
- **Rate Limit:** 1000/min, 10000/hour
- **Timeout:** 5s (Herald.lol requirement)
- **Caching:** 10 minutes for analysis results

### Riot API Integration
- **Endpoint:** `/api/riot`
- **Rate Limit:** 95/2min, 550/10min (Riot compliance)
- **Caching:** 5 minutes for match data
- **Circuit Breaker:** Automatic failover

### Match Analysis Service
- **Endpoint:** `/api/matches`
- **Rate Limit:** 1000/min
- **Performance SLA:** <5s requirement
- **Caching:** 10 minutes for match analysis

### Team Composition Service
- **Endpoint:** `/api/team-composition`
- **Features:** Champion synergy optimization
- **Rate Limit:** 500/min
- **Caching:** 30 minutes for compositions

### Counter-Picks Service
- **Endpoint:** `/api/counter-picks`
- **Features:** Champion counter analysis
- **Rate Limit:** 500/min
- **Caching:** 30 minutes for counter data

## 🔐 Gaming Authentication & Security

### Multi-Tier User Support

| Tier | Rate Limit | Features |
|------|-----------|----------|
| Free | 100/min | Basic analytics |
| Premium | 500/min | Advanced analytics |
| Pro | 2000/min | Professional tools |
| Enterprise | 10000/min | Full platform access |

### Security Features
- **SSL/TLS:** TLS 1.2+ encryption
- **CORS:** Gaming-specific origins
- **Security Headers:** XSS, CSRF, clickjacking protection
- **API Key Authentication:** Multi-tier access control
- **Rate Limiting:** Distributed Redis-based limiting

## 📊 Gaming Monitoring & Metrics

### Kong Manager UI
- **URL:** http://localhost:1337
- **Features:** Real-time gaming metrics
- **Dashboards:** Performance, errors, usage

### Datadog Integration
- **Metrics:** Gaming performance, errors, usage
- **Logs:** Access logs, error logs
- **Alerts:** Performance degradation, error spikes

### Prometheus Metrics
- **Gaming Analytics:** Request counts, latencies
- **Riot API:** Rate limit usage, failures
- **System Health:** Uptime, response times

## 🎯 Gaming Performance Optimization

### Caching Strategy
```
Gaming Analytics: 10min TTL
Riot API Data: 5min TTL  
Match Analysis: 10min TTL
Team Compositions: 30min TTL
Counter-picks: 30min TTL
```

### Load Balancing
- **Algorithm:** Round-robin with health checks
- **Health Checks:** HTTP /health endpoint
- **Failover:** Automatic unhealthy backend removal

### Connection Pooling
- **Keep-alive:** 60s timeout
- **Pool Size:** 32 connections per upstream
- **Max Requests:** 1000 per connection

## 🛠️ Development & Testing

### Local Development

```bash
# Start development environment
docker-compose -f gateway/docker-compose.kong.yml up

# Apply gaming configuration
./gateway/scripts/kong-setup.sh

# Test gaming endpoints
curl -H "apikey: herald-gaming-test-key-123" \
     http://localhost:8000/api/analytics/health
```

### Gaming API Testing

```bash
# Test gaming analytics endpoint
curl -X GET "http://localhost:8000/api/analytics/kda?summoner=test" \
     -H "apikey: herald-gaming-test-key-123"

# Test Riot API endpoint
curl -X GET "http://localhost:8000/api/riot/summoner/test" \
     -H "apikey: herald-gaming-test-key-123"

# Test match analysis
curl -X GET "http://localhost:8000/api/matches/analysis/12345" \
     -H "apikey: herald-gaming-test-key-123"
```

## 🔧 Configuration Management

### Kong Configuration Files
- `kong/kong.yml` - Declarative gaming configuration
- `kong/nginx-custom.conf` - Gaming performance optimization
- `kong/redis.conf` - Gaming cache configuration

### Gaming Plugin Configuration
- **Rate Limiting:** Redis-based distributed limiting
- **Caching:** Gaming-optimized TTL settings
- **Monitoring:** Prometheus + Datadog integration
- **Security:** CORS, headers, authentication

### Environment Variables
```bash
# Gaming performance
KONG_NGINX_WORKER_PROCESSES=auto
KONG_MEM_CACHE_SIZE=256m

# Gaming database
KONG_PG_PASSWORD=secure_password

# Gaming Redis
REDIS_PASSWORD=redis_password
```

## 🚨 Troubleshooting

### Common Gaming Issues

**Kong not starting:**
```bash
# Check database connection
docker-compose -f gateway/docker-compose.kong.yml logs kong-database

# Check Kong logs
docker-compose -f gateway/docker-compose.kong.yml logs kong
```

**Gaming endpoints not responding:**
```bash
# Check service health
curl http://localhost:8001/services

# Check route configuration
curl http://localhost:8001/routes

# Test backend connectivity
curl http://herald-backend:8080/health
```

**Rate limiting issues:**
```bash
# Check Redis connection
docker-compose -f gateway/docker-compose.kong.yml logs herald-redis

# Check rate limit status
curl http://localhost:8001/plugins | jq '.data[] | select(.name=="rate-limiting-advanced")'
```

## 📚 Gaming API Documentation

### Gaming Analytics Endpoints
- `GET /api/analytics/kda` - KDA analysis
- `GET /api/analytics/cs` - CS/min calculations  
- `GET /api/analytics/vision` - Vision score analysis
- `GET /api/analytics/damage` - Damage analysis
- `GET /api/analytics/gold` - Gold efficiency

### Riot API Endpoints
- `GET /api/riot/summoner/{name}` - Summoner data
- `GET /api/riot/matches/{id}` - Match details
- `GET /api/riot/league/{id}` - League information

### Team Composition Endpoints
- `GET /api/team-composition/optimize` - Optimize team comp
- `GET /api/team-composition/synergy` - Champion synergy
- `POST /api/team-composition/analyze` - Analyze composition

## 🔗 Related Documentation

- [Herald.lol Gaming Analytics](../CLAUDE.md)
- [Gaming Backend Services](../backend/README.md)
- [Gaming Frontend Application](../frontend/README.md)
- [Gaming Testing Strategy](../TESTING_STRATEGY.md)
- [Kong Official Documentation](https://docs.konghq.com/)

## 🎮 Gaming Support

For Herald.lol gaming platform support:
- **Issues:** Create GitHub issue with `gaming` label
- **Performance:** Check gaming metrics dashboard
- **Security:** Review gaming security logs
- **API:** Test gaming endpoints with provided keys

---

**Herald.lol Gaming Analytics Platform**  
Performance Target: <5s analytics, <2s UI, 99.9% uptime, 1M+ concurrent users