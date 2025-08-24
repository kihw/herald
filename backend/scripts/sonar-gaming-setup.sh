#!/bin/bash

# Herald.lol Gaming Analytics - SonarQube Setup Script
# Complete setup and configuration for gaming platform code quality

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
SONAR_HOST="${SONAR_HOST:-http://localhost:9000}"
SONAR_PROJECT_KEY="herald-gaming-analytics"
GAMING_PERFORMANCE_TARGET=5000
GAMING_CONCURRENT_USERS=1000000

# Colors for gaming theme
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
GOLD='\033[0;33m'
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

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_gaming() {
    echo -e "${GOLD}[GAMING]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_info "🔍 Checking prerequisites for Herald.lol SonarQube setup..."
    
    local missing_tools=()
    
    # Check Docker and Docker Compose
    if ! command -v docker >/dev/null 2>&1; then
        missing_tools+=("docker")
    fi
    
    if ! command -v docker-compose >/dev/null 2>&1 && ! docker compose version >/dev/null 2>&1; then
        missing_tools+=("docker-compose")
    fi
    
    # Check curl for API calls
    if ! command -v curl >/dev/null 2>&1; then
        missing_tools+=("curl")
    fi
    
    # Check jq for JSON processing
    if ! command -v jq >/dev/null 2>&1; then
        missing_tools+=("jq")
    fi
    
    if [ ${#missing_tools[@]} -gt 0 ]; then
        log_error "❌ Missing required tools: ${missing_tools[*]}"
        log_info "Please install missing tools and try again"
        return 1
    fi
    
    log_success "✅ All prerequisites met"
    return 0
}

# Prepare gaming environment
prepare_gaming_environment() {
    log_info "🎮 Preparing Herald.lol gaming environment..."
    
    cd "$PROJECT_ROOT"
    
    # Create required directories
    log_info "📁 Creating gaming-specific directories..."
    mkdir -p volumes/sonarqube/{data,logs,extensions,database}
    mkdir -p sonarqube/{config,db-init,quality-profiles}
    
    # Set permissions for SonarQube
    if [ "$(id -u)" = "0" ]; then
        log_warning "⚠️ Running as root, adjusting permissions..."
        chown -R 999:999 volumes/sonarqube/
    else
        log_info "📝 Setting up volume permissions..."
        # Ensure directories are writable
        chmod -R 755 volumes/sonarqube/
    fi
    
    # Create gaming-specific SonarQube configuration
    create_gaming_sonar_config
    
    log_success "✅ Gaming environment prepared"
}

# Create gaming-specific SonarQube configuration
create_gaming_sonar_config() {
    log_gaming "🎮 Creating Herald.lol gaming-specific SonarQube configuration..."
    
    # Create sonar.properties for gaming optimization
    cat << 'EOF' > sonarqube/config/sonar.properties
# Herald.lol Gaming Analytics - SonarQube Configuration
# Optimized for gaming platform performance

# Gaming database configuration
sonar.jdbc.maxActive=60
sonar.jdbc.maxIdle=5
sonar.jdbc.minIdle=2
sonar.jdbc.maxWait=5000
sonar.jdbc.minEvictableIdleTimeMillis=600000
sonar.jdbc.timeBetweenEvictionRunsMillis=30000

# Gaming web server configuration
sonar.web.javaOpts=-Xmx2G -Xms1G -XX:+UseG1GC -XX:MaxGCPauseMillis=200
sonar.web.http.maxThreads=50
sonar.web.http.minThreads=5
sonar.web.http.acceptCount=25

# Gaming compute engine configuration
sonar.ce.javaOpts=-Xmx2G -Xms512M -XX:+UseG1GC -XX:MaxGCPauseMillis=200
sonar.ce.workerCount=2

# Gaming Elasticsearch configuration
sonar.search.javaOpts=-Xmx1G -Xms512M -XX:+UseG1GC -XX:MaxGCPauseMillis=200

# Gaming platform security
sonar.forceAuthentication=true
sonar.security.realm=sonar

# Gaming performance settings
sonar.web.sessionTimeoutInMinutes=480
sonar.dbcleaner.hoursBeforeKeepingOnlyOneSnapshotByDay=24
sonar.dbcleaner.weeksBeforeDeletingAllSnapshots=4

# Gaming analytics retention
sonar.dbcleaner.audit.weeksBeforeDeletion=4

# Gaming platform logging
sonar.log.level=INFO
sonar.path.logs=logs

# Gaming telemetry (disabled for performance)
sonar.telemetry.enable=false
EOF

    # Create gaming-specific database initialization
    cat << 'EOF' > sonarqube/db-init/01-gaming-setup.sql
-- Herald.lol Gaming Analytics - Database Setup
-- Optimized for gaming platform performance

-- Gaming performance optimizations
ALTER SYSTEM SET shared_preload_libraries = 'pg_stat_statements';
ALTER SYSTEM SET max_connections = 200;
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET work_mem = '16MB';
ALTER SYSTEM SET maintenance_work_mem = '128MB';
ALTER SYSTEM SET random_page_cost = 1.1;
ALTER SYSTEM SET effective_io_concurrency = 200;

-- Gaming analytics specific indexes
-- (These will be created by SonarQube, but we can prepare the database)

-- Gaming platform metadata
COMMENT ON DATABASE herald_sonar IS 'Herald.lol Gaming Analytics Platform - SonarQube Database';

SELECT pg_reload_conf();
EOF

    log_success "✅ Gaming SonarQube configuration created"
}

# Start SonarQube with gaming configuration
start_sonarqube() {
    log_info "🚀 Starting Herald.lol SonarQube gaming setup..."
    
    cd "$PROJECT_ROOT"
    
    # Create environment file if it doesn't exist
    if [ ! -f ".env" ]; then
        cat << 'EOF' > .env
# Herald.lol Gaming Analytics - SonarQube Environment
SONAR_DB_PASSWORD=herald_gaming_sonar_2024
EOF
        log_info "📝 Created .env file with gaming defaults"
    fi
    
    # Start SonarQube services
    log_info "🐳 Starting gaming SonarQube containers..."
    
    # Use docker-compose or docker compose based on availability
    if command -v docker-compose >/dev/null 2>&1; then
        docker-compose -f docker-compose.sonarqube.yml up -d
    else
        docker compose -f docker-compose.sonarqube.yml up -d
    fi
    
    # Wait for SonarQube to be ready
    log_info "⏳ Waiting for Herald.lol SonarQube to be ready..."
    wait_for_sonarqube
    
    log_success "✅ SonarQube started successfully"
}

# Wait for SonarQube to be ready
wait_for_sonarqube() {
    local max_attempts=60
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "$SONAR_HOST/api/system/status" | grep -q "UP"; then
            log_success "✅ SonarQube is ready!"
            return 0
        fi
        
        echo -n "."
        sleep 5
        attempt=$((attempt + 1))
    done
    
    log_error "❌ SonarQube failed to start within $(($max_attempts * 5)) seconds"
    return 1
}

# Configure gaming project
configure_gaming_project() {
    log_gaming "🎮 Configuring Herald.lol gaming project in SonarQube..."
    
    # Get admin token (using default admin credentials)
    local admin_token=$(get_or_create_admin_token)
    
    if [ -z "$admin_token" ]; then
        log_error "❌ Failed to get admin token"
        return 1
    fi
    
    # Create gaming project
    log_info "📊 Creating gaming analytics project..."
    curl -s -u "admin:admin" -X POST \
        "$SONAR_HOST/api/projects/create" \
        -d "project=$SONAR_PROJECT_KEY" \
        -d "name=Herald.lol Gaming Analytics Platform" \
        -d "visibility=private" >/dev/null || log_warning "Project may already exist"
    
    # Set gaming project settings
    log_gaming "⚙️ Configuring gaming-specific project settings..."
    
    # Gaming performance targets
    curl -s -u "admin:admin" -X POST \
        "$SONAR_HOST/api/settings/set" \
        -d "component=$SONAR_PROJECT_KEY" \
        -d "key=sonar.gaming.performance.target" \
        -d "value=$GAMING_PERFORMANCE_TARGET" >/dev/null
    
    curl -s -u "admin:admin" -X POST \
        "$SONAR_HOST/api/settings/set" \
        -d "component=$SONAR_PROJECT_KEY" \
        -d "key=sonar.gaming.concurrent.users" \
        -d "value=$GAMING_CONCURRENT_USERS" >/dev/null
    
    # Gaming-specific quality profile
    setup_gaming_quality_profiles "$admin_token"
    
    # Create gaming quality gate
    setup_gaming_quality_gate "$admin_token"
    
    log_success "✅ Gaming project configured"
}

# Get or create admin token
get_or_create_admin_token() {
    # Try to create a token for gaming analytics
    local token_response=$(curl -s -u "admin:admin" -X POST \
        "$SONAR_HOST/api/user_tokens/generate" \
        -d "name=herald-gaming-token" \
        -d "type=USER_TOKEN" 2>/dev/null || echo "")
    
    if [ -n "$token_response" ]; then
        echo "$token_response" | jq -r '.token' 2>/dev/null || echo ""
    else
        # If token creation fails, we'll use basic auth
        echo ""
    fi
}

# Setup gaming quality profiles
setup_gaming_quality_profiles() {
    local token=$1
    
    log_gaming "🎯 Setting up Herald.lol gaming quality profiles..."
    
    # Create Go quality profile for gaming
    create_go_gaming_profile
    
    # Create TypeScript quality profile for gaming
    create_typescript_gaming_profile
    
    log_success "✅ Gaming quality profiles created"
}

# Create Go quality profile for gaming
create_go_gaming_profile() {
    log_info "🐹 Creating Go quality profile for gaming platform..."
    
    # Export and modify Sonar way profile for gaming
    cat << 'EOF' > sonarqube/quality-profiles/herald-gaming-go.xml
<?xml version="1.0" encoding="UTF-8"?>
<profile>
  <name>Herald Gaming Go Profile</name>
  <language>go</language>
  <parent>Sonar way</parent>
  <rules>
    <!-- Gaming performance optimizations -->
    <rule>
      <repositoryKey>go</repositoryKey>
      <key>S3776</key>
      <priority>INFO</priority>
      <parameters>
        <parameter>
          <key>threshold</key>
          <value>20</value> <!-- Higher threshold for gaming analytics complexity -->
        </parameter>
      </parameters>
    </rule>
    
    <!-- Gaming configuration flexibility -->
    <rule>
      <repositoryKey>go</repositoryKey>
      <key>S107</key>
      <priority>INFO</priority>
      <parameters>
        <parameter>
          <key>max</key>
          <value>10</value> <!-- More parameters allowed for gaming config -->
        </parameter>
      </parameters>
    </rule>
    
    <!-- Gaming error handling -->
    <rule>
      <repositoryKey>go</repositoryKey>
      <key>S1764</key>
      <priority>MAJOR</priority>
    </rule>
    
    <!-- Gaming security rules -->
    <rule>
      <repositoryKey>go</repositoryKey>
      <key>S2068</key>
      <priority>BLOCKER</priority>
    </rule>
  </rules>
</profile>
EOF

    log_success "✅ Go gaming quality profile created"
}

# Create TypeScript quality profile for gaming
create_typescript_gaming_profile() {
    log_info "📝 Creating TypeScript quality profile for gaming platform..."
    
    cat << 'EOF' > sonarqube/quality-profiles/herald-gaming-typescript.xml
<?xml version="1.0" encoding="UTF-8"?>
<profile>
  <name>Herald Gaming TypeScript Profile</name>
  <language>ts</language>
  <parent>Sonar way</parent>
  <rules>
    <!-- Gaming frontend performance -->
    <rule>
      <repositoryKey>typescript</repositoryKey>
      <key>S3776</key>
      <priority>INFO</priority>
      <parameters>
        <parameter>
          <key>threshold</key>
          <value>20</value> <!-- Higher threshold for gaming UI complexity -->
        </parameter>
      </parameters>
    </rule>
    
    <!-- Gaming React patterns -->
    <rule>
      <repositoryKey>typescript</repositoryKey>
      <key>S6426</key>
      <priority>MAJOR</priority>
    </rule>
    
    <!-- Gaming security -->
    <rule>
      <repositoryKey>typescript</repositoryKey>
      <key>S2068</key>
      <priority>BLOCKER</priority>
    </rule>
  </rules>
</profile>
EOF

    log_success "✅ TypeScript gaming quality profile created"
}

# Setup gaming quality gate
setup_gaming_quality_gate() {
    local token=$1
    
    log_gaming "🎯 Setting up Herald.lol gaming quality gate..."
    
    # Create gaming quality gate
    local gate_response=$(curl -s -u "admin:admin" -X POST \
        "$SONAR_HOST/api/qualitygates/create" \
        -d "name=Herald Gaming Gate" 2>/dev/null || echo "")
    
    if [ -n "$gate_response" ]; then
        local gate_id=$(echo "$gate_response" | jq -r '.id' 2>/dev/null || echo "")
        
        if [ -n "$gate_id" ] && [ "$gate_id" != "null" ]; then
            # Gaming-specific quality conditions
            setup_gaming_quality_conditions "$gate_id"
            
            # Set as default for gaming project
            curl -s -u "admin:admin" -X POST \
                "$SONAR_HOST/api/qualitygates/select" \
                -d "projectKey=$SONAR_PROJECT_KEY" \
                -d "gateId=$gate_id" >/dev/null
            
            log_success "✅ Gaming quality gate created and set as default"
        else
            log_warning "⚠️ Could not extract quality gate ID, using default"
        fi
    else
        log_warning "⚠️ Could not create quality gate, using default"
    fi
}

# Setup gaming quality conditions
setup_gaming_quality_conditions() {
    local gate_id=$1
    
    log_info "⚙️ Setting up gaming quality conditions..."
    
    # Gaming performance: Coverage should be at least 70% for gaming reliability
    curl -s -u "admin:admin" -X POST \
        "$SONAR_HOST/api/qualitygates/create_condition" \
        -d "gateId=$gate_id" \
        -d "metric=coverage" \
        -d "op=LT" \
        -d "error=70" >/dev/null
    
    # Gaming quality: No new bugs for gaming stability
    curl -s -u "admin:admin" -X POST \
        "$SONAR_HOST/api/qualitygates/create_condition" \
        -d "gateId=$gate_id" \
        -d "metric=new_bugs" \
        -d "op=GT" \
        -d "error=0" >/dev/null
    
    # Gaming security: No new vulnerabilities for gaming data protection
    curl -s -u "admin:admin" -X POST \
        "$SONAR_HOST/api/qualitygates/create_condition" \
        -d "gateId=$gate_id" \
        -d "metric=new_vulnerabilities" \
        -d "op=GT" \
        -d "error=0" >/dev/null
    
    # Gaming maintainability: Technical debt should be manageable
    curl -s -u "admin:admin" -X POST \
        "$SONAR_HOST/api/qualitygates/create_condition" \
        -d "gateId=$gate_id" \
        -d "metric=sqale_rating" \
        -d "op=GT" \
        -d "error=1" >/dev/null
    
    log_success "✅ Gaming quality conditions configured"
}

# Run initial gaming analysis
run_initial_analysis() {
    log_gaming "🔍 Running initial Herald.lol gaming code analysis..."
    
    cd "$PROJECT_ROOT"
    
    # Check if sonar-scanner is available
    if command -v sonar-scanner >/dev/null 2>&1; then
        log_info "🔍 Running SonarQube scanner..."
        sonar-scanner \
            -Dsonar.host.url="$SONAR_HOST" \
            -Dsonar.login="admin" \
            -Dsonar.password="admin"
    else
        log_info "🐳 Running SonarQube scanner via Docker..."
        
        # Use docker-compose or docker compose
        if command -v docker-compose >/dev/null 2>&1; then
            docker-compose -f docker-compose.sonarqube.yml run --rm sonar-scanner-herald
        else
            docker compose -f docker-compose.sonarqube.yml run --rm sonar-scanner-herald
        fi
    fi
    
    log_success "✅ Initial gaming analysis completed"
}

# Display gaming dashboard
show_gaming_dashboard() {
    log_gaming "🎮 Herald.lol Gaming Code Quality Dashboard"
    echo "========================================"
    echo "⚡ Performance Target: <${GAMING_PERFORMANCE_TARGET}ms"
    echo "👥 Concurrent Users: ${GAMING_CONCURRENT_USERS:,}"
    echo "🕐 SonarQube URL: $SONAR_HOST"
    echo ""
    
    log_info "📊 Gaming Project Information"
    echo "Project Key: $SONAR_PROJECT_KEY"
    echo "Project Name: Herald.lol Gaming Analytics Platform"
    echo ""
    
    log_info "🎯 Gaming Quality Targets"
    echo "• Code Coverage: ≥70% (gaming reliability)"
    echo "• New Bugs: 0 (gaming stability)"
    echo "• Vulnerabilities: 0 (gaming security)"
    echo "• Maintainability: A rating (gaming performance)"
    echo ""
    
    log_info "🔗 Gaming Access Links"
    echo "• Dashboard: $SONAR_HOST/dashboard?id=$SONAR_PROJECT_KEY"
    echo "• Issues: $SONAR_HOST/project/issues?id=$SONAR_PROJECT_KEY"
    echo "• Coverage: $SONAR_HOST/component_measures?id=$SONAR_PROJECT_KEY&metric=coverage"
    echo "• Security: $SONAR_HOST/project/security_hotspots?id=$SONAR_PROJECT_KEY"
    echo ""
    
    log_success "🎮 Herald.lol SonarQube setup completed!"
}

# Cleanup SonarQube
cleanup_sonarqube() {
    log_warning "🧹 Cleaning up Herald.lol SonarQube..."
    
    read -p "⚠️ This will stop and remove all SonarQube containers and data. Continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Cleanup cancelled"
        return 0
    fi
    
    cd "$PROJECT_ROOT"
    
    # Stop containers
    if command -v docker-compose >/dev/null 2>&1; then
        docker-compose -f docker-compose.sonarqube.yml down -v
    else
        docker compose -f docker-compose.sonarqube.yml down -v
    fi
    
    # Remove volumes
    log_warning "🗑️ Removing gaming SonarQube data..."
    rm -rf volumes/sonarqube/
    
    log_success "✅ Gaming SonarQube cleanup completed"
}

# Usage
usage() {
    cat << EOF
🎮 Herald.lol Gaming SonarQube Setup

Usage: $0 [COMMAND]

COMMANDS:
    setup           Complete gaming SonarQube setup
    start           Start SonarQube containers
    configure       Configure gaming project
    analyze         Run code analysis
    dashboard       Show gaming dashboard info
    cleanup         Clean up all SonarQube data
    -h, --help      Show this help message

ENVIRONMENT VARIABLES:
    SONAR_HOST      SonarQube URL (default: http://localhost:9000)

EXAMPLES:
    # Complete gaming setup
    $0 setup

    # Just start containers
    $0 start

    # Run gaming analysis
    $0 analyze

🎯 Gaming Focus: Code quality for Herald.lol gaming analytics
⚡ Performance: Optimized for <${GAMING_PERFORMANCE_TARGET}ms gaming target
👥 Scale: Support for ${GAMING_CONCURRENT_USERS:,}+ concurrent gaming users
EOF
}

# Main function
main() {
    case "${1:-}" in
        setup)
            check_prerequisites
            prepare_gaming_environment
            start_sonarqube
            configure_gaming_project
            show_gaming_dashboard
            ;;
        start)
            check_prerequisites
            start_sonarqube
            ;;
        configure)
            configure_gaming_project
            ;;
        analyze)
            run_initial_analysis
            ;;
        dashboard)
            show_gaming_dashboard
            ;;
        cleanup)
            cleanup_sonarqube
            ;;
        -h|--help)
            usage
            ;;
        *)
            log_gaming "🎮 Herald.lol Gaming SonarQube Setup"
            echo ""
            usage
            exit 1
            ;;
    esac
}

# Handle Ctrl+C gracefully
trap 'echo ""; log_info "Gaming SonarQube setup interrupted"; exit 0' SIGINT

# Run main function
main "$@"