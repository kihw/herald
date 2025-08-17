# 🚀 Production Deployment Guide

## 📋 Prerequisites

### System Requirements

- **Docker Desktop** (Latest version)
- **Docker Compose** v2.0+
- **Minimum 4GB RAM** (8GB recommended)
- **10GB free disk space**
- **Windows 10/11** or **Linux/macOS**

### Network Requirements

- **Port 80**: HTTP traffic (can be changed)
- **Port 443**: HTTPS traffic (optional)
- **Port 5432**: PostgreSQL (internal)
- **Port 6379**: Redis (internal)
- **Port 8001**: Backend API (internal)

## 🔧 Pre-Deployment Setup

### 1. Environment Configuration

```bash
# Copy environment template
cp .env.production .env

# Edit configuration
nano .env  # Linux/macOS
notepad .env  # Windows
```

**Required configurations:**

- `RIOT_API_KEY`: Your Riot Games API key
- `POSTGRES_PASSWORD`: Strong database password
- `REDIS_PASSWORD`: Strong Redis password
- `JWT_SECRET`: 32+ character secret key

### 2. SSL Configuration (Optional)

```bash
# Create SSL directory
mkdir -p nginx/ssl

# Copy your SSL certificates
cp your-cert.pem nginx/ssl/cert.pem
cp your-key.pem nginx/ssl/key.pem

# Update .env
SSL_ENABLED=true
```

## 🚀 Deployment Commands

### Quick Deploy (Recommended)

```powershell
# Windows PowerShell
.\deploy.ps1 -Build -Deploy

# Alternative: Manual Docker Compose
docker-compose -f docker-compose.prod.yml up -d --build
```

### Step-by-Step Deployment

**1. Build Services**

```bash
docker-compose -f docker-compose.prod.yml build
```

**2. Start Infrastructure**

```bash
docker-compose -f docker-compose.prod.yml up -d postgres redis
```

**3. Wait for Database Ready**

```bash
# Check database health
docker-compose -f docker-compose.prod.yml exec postgres pg_isready -U lol_user
```

**4. Deploy Application**

```bash
docker-compose -f docker-compose.prod.yml up -d backend frontend nginx
```

## 📊 Service Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Production Stack                        │
└─────────────────────────────────────────────────────────────┘

┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Nginx       │    │    Frontend     │    │    Backend      │
│   (Port 80)     │◄──►│   (React App)   │◄──►│   (Go API)      │
│   Load Balancer │    │   Static Files  │    │   Business Logic│
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                                              │
         │              ┌─────────────────┐             │
         │              │   PostgreSQL    │◄────────────┘
         │              │   (Database)    │
         │              └─────────────────┘
         │
         │              ┌─────────────────┐
         └──────────────►│      Redis      │
                        │     (Cache)     │
                        └─────────────────┘
```

## 🔍 Health Monitoring

### Service Health Checks

```bash
# Check all services
docker-compose -f docker-compose.prod.yml ps

# Check specific service health
docker-compose -f docker-compose.prod.yml exec backend curl -f http://localhost:8001/api/health

# Check application endpoints
curl http://localhost/health          # Frontend
curl http://localhost/api/health      # Backend (via proxy)
```

### Application URLs

- **Frontend**: http://localhost
- **API**: http://localhost/api
- **Direct Backend**: http://localhost:8001 (if needed)

## 📱 Management Commands

### Using Deploy Script (Windows)

```powershell
# Show status
.\deploy.ps1 -Status

# View logs
.\deploy.ps1 -Logs

# View specific service logs
.\deploy.ps1 -Logs -Service backend

# Stop all services
.\deploy.ps1 -Stop

# Clean up resources
.\deploy.ps1 -Clean

# Rebuild and redeploy
.\deploy.ps1 -Build -Deploy
```

### Using Docker Compose Directly

```bash
# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Restart service
docker-compose -f docker-compose.prod.yml restart backend

# Scale service (if needed)
docker-compose -f docker-compose.prod.yml up -d --scale backend=2

# Execute commands in container
docker-compose -f docker-compose.prod.yml exec backend sh
```

## 🔐 Security Configuration

### Firewall Rules

```bash
# Allow HTTP traffic
ufw allow 80/tcp

# Allow HTTPS traffic (if SSL enabled)
ufw allow 443/tcp

# Block direct access to internal ports
ufw deny 5432/tcp  # PostgreSQL
ufw deny 6379/tcp  # Redis
ufw deny 8001/tcp  # Backend (optional)
```

### Environment Security

- **Never commit** `.env` file to version control
- **Use strong passwords** (minimum 16 characters)
- **Rotate secrets** regularly (every 90 days)
- **Enable SSL/TLS** for production traffic
- **Configure CORS** properly in `.env`

## 🔄 Updates & Maintenance

### Application Updates

```bash
# 1. Pull latest code
git pull origin main

# 2. Rebuild and redeploy
.\deploy.ps1 -Build -Deploy

# 3. Verify deployment
.\deploy.ps1 -Status
```

### Database Migrations

```bash
# Run migrations (when available)
docker-compose -f docker-compose.prod.yml exec backend ./server --migrate

# Backup before migrations
docker-compose -f docker-compose.prod.yml exec postgres pg_dump -U lol_user lol_match_db > backup.sql
```

### Data Backup

```bash
# Database backup
docker-compose -f docker-compose.prod.yml exec postgres pg_dump -U lol_user -h localhost lol_match_db | gzip > backup_$(date +%Y%m%d_%H%M%S).sql.gz

# Redis backup
docker-compose -f docker-compose.prod.yml exec redis redis-cli --rdb dump.rdb
```

## 🐛 Troubleshooting

### Common Issues

**Service Won't Start**

```bash
# Check service logs
docker-compose -f docker-compose.prod.yml logs service_name

# Check resource usage
docker stats

# Restart service
docker-compose -f docker-compose.prod.yml restart service_name
```

**Database Connection Issues**

```bash
# Check database status
docker-compose -f docker-compose.prod.yml exec postgres pg_isready -U lol_user

# Check database logs
docker-compose -f docker-compose.prod.yml logs postgres

# Recreate database (WARNING: Data loss)
docker-compose -f docker-compose.prod.yml down -v
docker-compose -f docker-compose.prod.yml up -d postgres
```

**High Memory Usage**

```bash
# Check container resource usage
docker stats

# Restart high-memory services
docker-compose -f docker-compose.prod.yml restart backend

# Clear Redis cache
docker-compose -f docker-compose.prod.yml exec redis redis-cli FLUSHALL
```

### Performance Tuning

**Database Optimization**

```bash
# Check database performance
docker-compose -f docker-compose.prod.yml exec postgres psql -U lol_user -d lol_match_db -c "SELECT * FROM pg_stat_activity;"

# Analyze query performance
docker-compose -f docker-compose.prod.yml exec postgres psql -U lol_user -d lol_match_db -c "EXPLAIN ANALYZE SELECT * FROM matches LIMIT 10;"
```

**Redis Cache Tuning**

```bash
# Check Redis memory usage
docker-compose -f docker-compose.prod.yml exec redis redis-cli INFO memory

# Check cache hit rates
docker-compose -f docker-compose.prod.yml exec redis redis-cli INFO stats
```

## 📊 Production Monitoring

### Metrics Collection

The application includes built-in performance monitoring:

- **Render time tracking**
- **Memory usage monitoring**
- **API response times**
- **Cache hit rates**
- **Error rates**

### Log Aggregation

```bash
# Centralized logging
docker-compose -f docker-compose.prod.yml logs -f | grep ERROR

# Export logs for analysis
docker-compose -f docker-compose.prod.yml logs --since 1h > recent_logs.txt
```

### Health Monitoring Script

```powershell
# Create monitoring script (monitor.ps1)
while ($true) {
    $status = .\deploy.ps1 -Status
    Write-Host "$(Get-Date): $status"
    Start-Sleep -Seconds 300  # Check every 5 minutes
}
```

## 🎯 Production Checklist

Before going live, ensure:

- [ ] **Environment configured** (`.env` file complete)
- [ ] **SSL certificates** installed (if using HTTPS)
- [ ] **Firewall rules** configured properly
- [ ] **Database backups** scheduled
- [ ] **Monitoring** set up
- [ ] **Domain DNS** pointing to server (if applicable)
- [ ] **Load testing** completed
- [ ] **Security audit** passed
- [ ] **Documentation** updated
- [ ] **Team access** configured

## 🆘 Emergency Procedures

### Quick Rollback

```bash
# Stop current deployment
.\deploy.ps1 -Stop

# Deploy previous version
git checkout previous_working_tag
.\deploy.ps1 -Build -Deploy
```

### Emergency Stop

```bash
# Stop all services immediately
docker-compose -f docker-compose.prod.yml down

# Force stop if needed
docker stop $(docker ps -aq)
```

### Data Recovery

```bash
# Restore database from backup
docker-compose -f docker-compose.prod.yml exec -T postgres psql -U lol_user -d lol_match_db < backup.sql
```

---

## 🎉 Success!

Your LoL Match Exporter is now running in production mode!

- **Frontend**: http://localhost
- **API Documentation**: http://localhost/api/health
- **Monitoring**: Available via deploy script status

For support, check the logs and troubleshooting section above. 🚀
