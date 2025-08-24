#!/bin/bash

# Herald.lol Gaming Analytics - VPS Development Environment Setup
# Complete infrastructure setup for Phase 1 Q1 2025 Foundation

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Gaming colors for better terminal experience
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
GOLD='\033[1;33m'  # Herald.lol gaming gold
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
    echo -e "${GOLD}[$(date +'%Y-%m-%d %H:%M:%S')] [HERALD-GAMING]${NC} $1"
}

# Check if running on VPS/container
check_environment() {
    log_gaming "🎮 Checking Herald.lol VPS gaming environment..."
    
    if [ ! -f "/.dockerenv" ] && [ -z "${container:-}" ]; then
        log_warning "⚠️ Running on bare metal VPS - will adapt setup"
    else
        log_info "🐳 Detected containerized environment"
    fi
    
    # Check system resources
    local memory=$(free -m | awk 'NR==2{printf "%.0f", $2}')
    local disk=$(df -h / | awk 'NR==2{print $4}')
    
    log_info "💾 System Resources:"
    echo "  Memory: ${memory}MB"
    echo "  Disk Space: ${disk}"
    
    if [ "$memory" -lt 2048 ]; then
        log_warning "⚠️ Low memory detected - gaming performance may be limited"
    fi
}

# Install Docker if not available
install_docker() {
    log_gaming "🐳 Setting up Docker for Herald.lol gaming platform..."
    
    if command -v docker >/dev/null 2>&1; then
        log_success "✅ Docker already installed"
        docker --version
        return 0
    fi
    
    log_info "📦 Installing Docker for gaming development..."
    
    # Try different installation methods
    if command -v apk >/dev/null 2>&1; then
        # Alpine Linux
        log_info "🐧 Installing Docker on Alpine Linux..."
        apk update
        apk add --no-cache docker docker-compose
        rc-update add docker default
        service docker start
        
    elif command -v apt-get >/dev/null 2>&1; then
        # Debian/Ubuntu
        log_info "🐧 Installing Docker on Debian/Ubuntu..."
        apt-get update
        apt-get install -y ca-certificates curl gnupg lsb-release
        
        # Add Docker GPG key
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
        
        # Add repository
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
        
        # Install Docker
        apt-get update
        apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
        systemctl enable docker
        systemctl start docker
        
    elif command -v yum >/dev/null 2>&1; then
        # CentOS/RHEL
        log_info "🐧 Installing Docker on CentOS/RHEL..."
        yum install -y yum-utils
        yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
        yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
        systemctl enable docker
        systemctl start docker
        
    else
        log_error "❌ Unsupported system - manual Docker installation required"
        return 1
    fi
    
    # Add current user to docker group
    if [ -n "${USER:-}" ] && [ "$USER" != "root" ]; then
        usermod -aG docker "$USER"
        log_info "👤 Added $USER to docker group (restart required)"
    fi
    
    log_success "✅ Docker installed successfully for Herald.lol gaming"
}

# Install Docker Compose if needed
install_docker_compose() {
    log_info "🔧 Setting up Docker Compose for Herald.lol..."
    
    if docker compose version >/dev/null 2>&1; then
        log_success "✅ Docker Compose (plugin) available"
        return 0
    fi
    
    if command -v docker-compose >/dev/null 2>&1; then
        log_success "✅ Docker Compose (standalone) available"
        return 0
    fi
    
    log_info "📦 Installing Docker Compose..."
    
    # Install Docker Compose
    local compose_version="v2.21.0"
    local arch=$(uname -m)
    local url="https://github.com/docker/compose/releases/download/${compose_version}/docker-compose-linux-${arch}"
    
    curl -L "$url" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    
    # Verify installation
    if docker-compose --version >/dev/null 2>&1; then
        log_success "✅ Docker Compose installed successfully"
    else
        log_error "❌ Docker Compose installation failed"
        return 1
    fi
}

# Create required directories
create_directories() {
    log_gaming "📁 Creating Herald.lol gaming directory structure..."
    
    cd "$PROJECT_ROOT"
    
    # Create monitoring directories
    mkdir -p monitoring/{prometheus,grafana,logstash}
    mkdir -p monitoring/prometheus/{rules,targets}
    mkdir -p monitoring/grafana/{dashboards,datasources}
    mkdir -p monitoring/logstash/{pipeline,config}
    
    # Create nginx configuration
    mkdir -p nginx/conf.d
    
    # Create vault configuration
    mkdir -p vault/config
    
    # Create SQL initialization
    mkdir -p sql/init
    
    # Create logs directory
    mkdir -p logs
    
    log_success "✅ Herald.lol gaming directories created"
}

# Create configuration files
create_configs() {
    log_gaming "⚙️ Creating Herald.lol gaming configurations..."
    
    cd "$PROJECT_ROOT"
    
    # Prometheus configuration
    cat > monitoring/prometheus/prometheus.yml << 'EOF'
# Herald.lol Gaming Analytics - Prometheus Configuration
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    gaming_platform: "herald-lol"
    environment: "development"

rule_files:
  - "rules/*.yml"

alerting:
  alertmanagers:
    - static_configs:
        - targets: []

scrape_configs:
  # Herald.lol API metrics
  - job_name: 'herald-api'
    static_configs:
      - targets: ['herald-api:9091']
    scrape_interval: 5s
    metrics_path: /metrics
    
  # PostgreSQL metrics
  - job_name: 'postgres'
    static_configs:
      - targets: ['herald-postgres:9187']
    scrape_interval: 10s
    
  # Redis metrics
  - job_name: 'redis'
    static_configs:
      - targets: ['herald-redis:9121']
    scrape_interval: 10s
    
  # NGINX metrics
  - job_name: 'nginx'
    static_configs:
      - targets: ['herald-nginx:9113']
    scrape_interval: 10s
    
  # Node exporter
  - job_name: 'node'
    static_configs:
      - targets: ['localhost:9100']
    scrape_interval: 15s
EOF

    # Gaming performance rules
    cat > monitoring/prometheus/rules/gaming-rules.yml << 'EOF'
groups:
  - name: herald.gaming.performance
    rules:
      # Gaming response time alert
      - alert: GamingAnalyticsResponseTimeTooHigh
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket{job="herald-api"}[5m])) > 5
        for: 2m
        labels:
          severity: critical
          gaming: performance
        annotations:
          summary: "Herald.lol analytics response time > 5s"
          description: "Gaming analytics taking {{ $value }}s (target: <5s)"
          
      # Database performance
      - alert: GamingDatabaseSlowQueries
        expr: rate(postgres_slow_queries_total[5m]) > 0.1
        for: 1m
        labels:
          severity: warning
          gaming: database
        annotations:
          summary: "Herald.lol database slow queries detected"
          
      # Memory usage
      - alert: GamingHighMemoryUsage
        expr: (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) > 0.9
        for: 5m
        labels:
          severity: warning
          gaming: resources
        annotations:
          summary: "Herald.lol high memory usage: {{ $value }}%"
EOF

    # Grafana datasources
    cat > monitoring/grafana/datasources/datasources.yml << 'EOF'
apiVersion: 1

datasources:
  # Prometheus for Herald.lol gaming metrics
  - name: Herald Gaming Prometheus
    type: prometheus
    access: proxy
    url: http://herald-prometheus:9090
    isDefault: true
    jsonData:
      timeInterval: "5s"
      
  # InfluxDB for Herald.lol gaming time-series
  - name: Herald Gaming InfluxDB
    type: influxdb
    access: proxy
    url: http://herald-influxdb:8086
    database: gaming-metrics
    jsonData:
      organization: herald-gaming
      defaultBucket: gaming-metrics
      version: Flux
    secureJsonData:
      token: herald-gaming-analytics-token-dev
      
  # PostgreSQL for Herald.lol gaming data
  - name: Herald Gaming Database
    type: postgres
    url: herald-postgres:5432
    database: herald_gaming_dev
    user: herald_dev
    secureJsonData:
      password: herald_gaming_dev_2025
EOF

    # NGINX configuration
    cat > nginx/nginx.conf << 'EOF'
# Herald.lol Gaming Analytics - NGINX Configuration
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log notice;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
    use epoll;
    multi_accept on;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;
    
    # Gaming optimized logging
    log_format gaming '$remote_addr - $remote_user [$time_local] "$request" '
                     '$status $body_bytes_sent "$http_referer" '
                     '"$http_user_agent" "$http_x_forwarded_for" '
                     'rt=$request_time uct="$upstream_connect_time" '
                     'uht="$upstream_header_time" urt="$upstream_response_time"';
    
    access_log /var/log/nginx/access.log gaming;
    
    # Performance optimizations for gaming
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    
    # Gzip for gaming assets
    gzip on;
    gzip_vary on;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/json
        application/javascript
        application/xml+rss
        application/atom+xml
        image/svg+xml;
    
    include /etc/nginx/conf.d/*.conf;
}
EOF

    cat > nginx/conf.d/herald-gaming.conf << 'EOF'
# Herald.lol Gaming Analytics - Site Configuration
upstream herald_api {
    server herald-api:8080;
    keepalive 32;
}

upstream herald_grafana {
    server herald-grafana:3000;
}

upstream herald_kibana {
    server herald-kibana:5601;
}

server {
    listen 80;
    server_name localhost herald.local;
    
    # Security headers for gaming platform
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Referrer-Policy strict-origin-when-cross-origin;
    
    # Herald.lol API
    location /api/ {
        proxy_pass http://herald_api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Gaming performance optimizations
        proxy_connect_timeout 5s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
        
        # CORS for gaming frontend
        add_header Access-Control-Allow-Origin *;
        add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS";
        add_header Access-Control-Allow-Headers "Content-Type, Authorization";
    }
    
    # Grafana gaming dashboards
    location /grafana/ {
        proxy_pass http://herald_grafana/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # Kibana gaming logs
    location /kibana/ {
        proxy_pass http://herald_kibana/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # Health check
    location /health {
        access_log off;
        return 200 "Herald.lol Gaming Platform OK\n";
        add_header Content-Type text/plain;
    }
}
EOF

    # Database initialization
    cat > sql/init/01-herald-gaming-init.sql << 'EOF'
-- Herald.lol Gaming Analytics - Database Initialization
-- PostgreSQL setup for gaming platform development

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Gaming performance optimizations
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET work_mem = '16MB';
ALTER SYSTEM SET maintenance_work_mem = '128MB';
ALTER SYSTEM SET max_connections = '200';
ALTER SYSTEM SET log_min_duration_statement = '1000';

-- Gaming-specific schemas
CREATE SCHEMA IF NOT EXISTS gaming;
CREATE SCHEMA IF NOT EXISTS analytics;
CREATE SCHEMA IF NOT EXISTS riot_api;

-- Grant permissions
GRANT ALL PRIVILEGES ON SCHEMA gaming TO herald_dev;
GRANT ALL PRIVILEGES ON SCHEMA analytics TO herald_dev;
GRANT ALL PRIVILEGES ON SCHEMA riot_api TO herald_dev;

-- Initial gaming tables
CREATE TABLE IF NOT EXISTS gaming.players (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    riot_puuid VARCHAR(78) UNIQUE NOT NULL,
    summoner_name VARCHAR(100) NOT NULL,
    region VARCHAR(10) NOT NULL,
    tier VARCHAR(20),
    rank_value VARCHAR(5),
    league_points INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gaming.matches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    riot_match_id VARCHAR(20) UNIQUE NOT NULL,
    game_creation TIMESTAMP NOT NULL,
    game_duration INTEGER NOT NULL,
    queue_id INTEGER NOT NULL,
    region VARCHAR(10) NOT NULL,
    processed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS analytics.player_performance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    player_id UUID REFERENCES gaming.players(id),
    match_id UUID REFERENCES gaming.matches(id),
    champion_id INTEGER NOT NULL,
    kills INTEGER DEFAULT 0,
    deaths INTEGER DEFAULT 0,
    assists INTEGER DEFAULT 0,
    cs_per_minute DECIMAL(5,2),
    vision_score INTEGER DEFAULT 0,
    damage_dealt INTEGER DEFAULT 0,
    gold_earned INTEGER DEFAULT 0,
    game_won BOOLEAN DEFAULT FALSE,
    performance_score DECIMAL(5,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Gaming indexes for performance
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_players_riot_puuid ON gaming.players(riot_puuid);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_matches_riot_id ON gaming.matches(riot_match_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_performance_player_id ON analytics.player_performance(player_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_performance_match_id ON analytics.player_performance(match_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_performance_created_at ON analytics.player_performance(created_at DESC);

-- Gaming triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_players_updated_at BEFORE UPDATE ON gaming.players
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Gaming analytics views
CREATE OR REPLACE VIEW analytics.player_stats_summary AS
SELECT 
    p.riot_puuid,
    p.summoner_name,
    COUNT(pp.id) as total_games,
    AVG(pp.kills) as avg_kills,
    AVG(pp.deaths) as avg_deaths,
    AVG(pp.assists) as avg_assists,
    AVG(pp.cs_per_minute) as avg_cs_per_minute,
    AVG(pp.vision_score) as avg_vision_score,
    ROUND(SUM(CASE WHEN pp.game_won THEN 1 ELSE 0 END)::DECIMAL / COUNT(pp.id) * 100, 2) as win_rate,
    AVG(pp.performance_score) as avg_performance_score
FROM gaming.players p
LEFT JOIN analytics.player_performance pp ON p.id = pp.player_id
GROUP BY p.id, p.riot_puuid, p.summoner_name;

-- Gaming performance tracking
INSERT INTO gaming.players (riot_puuid, summoner_name, region, tier, rank_value) VALUES
('dev-test-player-1-puuid', 'TestPlayer1', 'EUW1', 'GOLD', 'II')
ON CONFLICT (riot_puuid) DO NOTHING;

-- Analyze tables for query optimization
ANALYZE gaming.players;
ANALYZE gaming.matches;
ANALYZE analytics.player_performance;
EOF

    log_success "✅ Herald.lol gaming configurations created"
}

# Set up environment file
create_env_file() {
    log_info "📝 Creating Herald.lol gaming environment file..."
    
    cd "$PROJECT_ROOT"
    
    if [ ! -f ".env" ]; then
        cat > .env << 'EOF'
# Herald.lol Gaming Analytics - Development Environment
COMPOSE_PROJECT_NAME=herald-gaming-dev

# Riot API Configuration (set your actual API key)
RIOT_API_KEY=RGAPI-your-riot-api-key-here
RIOT_REGION=euw1

# Gaming Performance Targets
GAMING_PERFORMANCE_TARGET_MS=5000
GAMING_MAX_CONCURRENT_USERS=1000000

# Database Settings
POSTGRES_DB=herald_gaming_dev
POSTGRES_USER=herald_dev
POSTGRES_PASSWORD=herald_gaming_dev_2025

# Redis Settings
REDIS_PASSWORD=

# InfluxDB Settings
INFLUXDB_ADMIN_PASSWORD=herald_gaming_metrics_2025
INFLUXDB_TOKEN=herald-gaming-analytics-token-dev

# Grafana Settings
GRAFANA_ADMIN_PASSWORD=herald_gaming_grafana_2025

# Development Settings
DEBUG=true
LOG_LEVEL=info
EOF
        log_info "📝 Created .env file - please update RIOT_API_KEY"
    else
        log_info "📝 .env file already exists"
    fi
}

# Build and start Herald.lol services
start_services() {
    log_gaming "🚀 Starting Herald.lol gaming development environment..."
    
    cd "$PROJECT_ROOT"
    
    # Make sure Docker is running
    if ! docker info >/dev/null 2>&1; then
        log_error "❌ Docker is not running"
        return 1
    fi
    
    # Start infrastructure services first
    log_info "🗄️ Starting Herald.lol gaming infrastructure..."
    docker-compose -f docker-compose.herald-dev.yml up -d \
        herald-postgres \
        herald-redis \
        herald-influxdb \
        herald-vault \
        herald-elasticsearch
    
    # Wait for databases to be ready
    log_info "⏳ Waiting for Herald.lol gaming databases to be ready..."
    sleep 30
    
    # Start monitoring stack
    log_info "📊 Starting Herald.lol gaming monitoring..."
    docker-compose -f docker-compose.herald-dev.yml up -d \
        herald-prometheus \
        herald-grafana \
        herald-logstash \
        herald-kibana
    
    # Wait for monitoring to be ready
    sleep 15
    
    # Start API and proxy
    log_info "🌐 Starting Herald.lol gaming API and proxy..."
    docker-compose -f docker-compose.herald-dev.yml up -d \
        herald-nginx
    
    # Show status
    show_status
}

# Show Herald.lol gaming environment status
show_status() {
    log_gaming "🎮 Herald.lol Gaming Development Environment Status"
    echo "================================================="
    
    cd "$PROJECT_ROOT"
    
    echo ""
    log_info "🌐 Gaming Services URLs:"
    echo "  🎯 Herald.lol API:      http://localhost/api/"
    echo "  📊 Gaming Dashboards:   http://localhost/grafana/ (admin:herald_gaming_grafana_2025)"
    echo "  📝 Gaming Logs:         http://localhost/kibana/"
    echo "  📈 Gaming Metrics:      http://localhost:9090 (Prometheus)"
    echo "  🔒 Gaming Vault:        http://localhost:8200 (token:herald-gaming-vault-root-token)"
    echo "  🗄️ Gaming Database:     localhost:5432 (herald_dev:herald_gaming_dev_2025)"
    echo "  ⚡ Gaming Cache:        localhost:6379"
    echo "  📊 Gaming Metrics DB:   http://localhost:8086"
    
    echo ""
    log_info "🐳 Gaming Container Status:"
    docker-compose -f docker-compose.herald-dev.yml ps
    
    echo ""
    log_info "🎯 Gaming Performance Targets:"
    echo "  ⚡ Analytics Response: <5000ms"
    echo "  👥 Concurrent Users: 1M+"
    echo "  🏆 Uptime: 99.9%"
    echo "  🎮 Gaming Focus: League of Legends analytics"
    
    echo ""
    log_success "✅ Herald.lol Gaming Development Environment Ready!"
    log_info "📚 Next steps:"
    echo "  1. Update RIOT_API_KEY in .env file"
    echo "  2. Run: make test-gaming"
    echo "  3. Start coding Herald.lol gaming features!"
}

# Cleanup function
cleanup_services() {
    log_warning "🧹 Cleaning up Herald.lol gaming environment..."
    
    cd "$PROJECT_ROOT"
    
    docker-compose -f docker-compose.herald-dev.yml down -v
    docker system prune -f
    
    log_success "✅ Herald.lol gaming environment cleaned up"
}

# Main function
main() {
    case "${1:-setup}" in
        setup)
            log_gaming "🎮 Setting up Herald.lol Gaming Development Environment"
            check_environment
            install_docker
            install_docker_compose
            create_directories
            create_configs
            create_env_file
            start_services
            ;;
        start)
            start_services
            ;;
        stop)
            cd "$PROJECT_ROOT"
            docker-compose -f docker-compose.herald-dev.yml stop
            ;;
        restart)
            cd "$PROJECT_ROOT"
            docker-compose -f docker-compose.herald-dev.yml restart
            ;;
        status)
            show_status
            ;;
        logs)
            cd "$PROJECT_ROOT"
            docker-compose -f docker-compose.herald-dev.yml logs -f "${2:-}"
            ;;
        cleanup)
            cleanup_services
            ;;
        --help|-h)
            echo "Herald.lol Gaming Analytics VPS Setup"
            echo ""
            echo "Usage: $0 [command]"
            echo ""
            echo "Commands:"
            echo "  setup    - Complete Herald.lol gaming setup (default)"
            echo "  start    - Start Herald.lol gaming services"
            echo "  stop     - Stop Herald.lol gaming services"
            echo "  restart  - Restart Herald.lol gaming services"
            echo "  status   - Show Herald.lol gaming status"
            echo "  logs     - Show Herald.lol gaming logs"
            echo "  cleanup  - Clean up Herald.lol gaming environment"
            echo ""
            echo "🎮 Herald.lol Gaming Analytics Platform"
            echo "⚡ Performance Target: <5000ms"
            echo "👥 Concurrent Users: 1M+"
            ;;
        *)
            log_error "❌ Unknown command: $1"
            echo "Use '$0 --help' for usage information"
            exit 1
            ;;
    esac
}

# Handle interruption gracefully
trap 'echo ""; log_info "🎮 Herald.lol setup interrupted"; exit 0' SIGINT SIGTERM

# Run main function
main "$@"