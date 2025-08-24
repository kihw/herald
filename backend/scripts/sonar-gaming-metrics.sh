#!/bin/bash

# Herald.lol Gaming Analytics - SonarQube Metrics Collection
# Automated collection and analysis of gaming-specific code quality metrics

set -euo pipefail

# Gaming configuration
SONARQUBE_URL="${SONARQUBE_URL:-http://sonarqube-herald:9000}"
SONARQUBE_TOKEN="${SONARQUBE_TOKEN:-squ_herald_gaming_token}"
GAMING_PROJECT_KEY="${GAMING_PROJECT_KEY:-herald-gaming-analytics}"
PERFORMANCE_TARGET_MS="${PERFORMANCE_TARGET_MS:-5000}"
METRICS_DATA_DIR="${METRICS_DATA_DIR:-/var/lib/gaming-metrics}"

# Colors for gaming console output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
GOLD='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] [GAMING-METRICS]${NC} $1"
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

# Initialize metrics collection
init_metrics() {
    log_gaming "🎮 Initializing Herald.lol gaming metrics collection..."
    
    # Create metrics directory
    mkdir -p "$METRICS_DATA_DIR"
    
    # Check SonarQube availability
    local max_retries=10
    local retry_count=0
    
    while [ $retry_count -lt $max_retries ]; do
        if curl -f "$SONARQUBE_URL/api/system/status" >/dev/null 2>&1; then
            log_success "✅ SonarQube is accessible for gaming metrics"
            break
        fi
        
        retry_count=$((retry_count + 1))
        log_warning "⏳ Waiting for SonarQube... (attempt $retry_count/$max_retries)"
        sleep 30
    done
    
    if [ $retry_count -eq $max_retries ]; then
        log_error "❌ SonarQube not accessible after $max_retries attempts"
        return 1
    fi
}

# Collect gaming-specific metrics
collect_gaming_metrics() {
    log_gaming "📊 Collecting Herald.lol gaming code quality metrics..."
    
    local timestamp=$(date +%s)
    local metrics_file="$METRICS_DATA_DIR/gaming-metrics-$timestamp.json"
    
    # Get basic project metrics
    local response=$(curl -s -H "Authorization: Bearer $SONARQUBE_TOKEN" \
        "$SONARQUBE_URL/api/measures/component?component=$GAMING_PROJECT_KEY&metricKeys=ncloc,complexity,coverage,duplicated_lines_density,violations,bugs,vulnerabilities,code_smells,security_hotspots")
    
    if [ $? -ne 0 ]; then
        log_error "❌ Failed to fetch basic gaming metrics"
        return 1
    fi
    
    # Parse and enhance with gaming context
    echo "$response" | jq --arg timestamp "$timestamp" --arg target "$PERFORMANCE_TARGET_MS" '
    {
        "timestamp": ($timestamp | tonumber),
        "gaming_project": "herald-gaming-analytics",
        "performance_target_ms": ($target | tonumber),
        "metrics": {
            "code_lines": (.component.measures[] | select(.metric == "ncloc") | .value // "0" | tonumber),
            "complexity": (.component.measures[] | select(.metric == "complexity") | .value // "0" | tonumber),
            "coverage_percent": (.component.measures[] | select(.metric == "coverage") | .value // "0" | tonumber),
            "duplication_percent": (.component.measures[] | select(.metric == "duplicated_lines_density") | .value // "0" | tonumber),
            "total_issues": (.component.measures[] | select(.metric == "violations") | .value // "0" | tonumber),
            "bugs": (.component.measures[] | select(.metric == "bugs") | .value // "0" | tonumber),
            "vulnerabilities": (.component.measures[] | select(.metric == "vulnerabilities") | .value // "0" | tonumber),
            "code_smells": (.component.measures[] | select(.metric == "code_smells") | .value // "0" | tonumber),
            "security_hotspots": (.component.measures[] | select(.metric == "security_hotspots") | .value // "0" | tonumber)
        }
    }' > "$metrics_file"
    
    log_success "✅ Gaming metrics collected: $metrics_file"
}

# Analyze gaming performance impact
analyze_gaming_impact() {
    log_gaming "🎯 Analyzing gaming performance impact..."
    
    local latest_metrics=$(ls -t "$METRICS_DATA_DIR"/gaming-metrics-*.json | head -1)
    
    if [ ! -f "$latest_metrics" ]; then
        log_error "❌ No gaming metrics found"
        return 1
    fi
    
    # Extract key gaming metrics
    local code_lines=$(jq -r '.metrics.code_lines' "$latest_metrics")
    local complexity=$(jq -r '.metrics.complexity' "$latest_metrics")
    local coverage=$(jq -r '.metrics.coverage_percent' "$latest_metrics")
    local bugs=$(jq -r '.metrics.bugs' "$latest_metrics")
    local vulnerabilities=$(jq -r '.metrics.vulnerabilities' "$latest_metrics")
    
    log_info "🎮 Herald.lol Gaming Code Quality Report"
    echo "========================================="
    echo "📊 Lines of Code: $code_lines"
    echo "🔧 Complexity: $complexity"
    echo "🛡️ Coverage: $coverage%"
    echo "🐛 Bugs: $bugs"
    echo "🔒 Vulnerabilities: $vulnerabilities"
    echo ""
    
    # Gaming-specific impact analysis
    local gaming_quality_score=100
    
    # Penalize for gaming performance risks
    if [ "$bugs" -gt 0 ]; then
        gaming_quality_score=$((gaming_quality_score - bugs * 5))
        log_warning "⚠️ Bugs detected: May impact gaming performance"
    fi
    
    if [ "$vulnerabilities" -gt 0 ]; then
        gaming_quality_score=$((gaming_quality_score - vulnerabilities * 10))
        log_warning "🔒 Security vulnerabilities: May compromise gaming data"
    fi
    
    if [ "$(echo "$coverage < 70" | bc -l 2>/dev/null || echo "1")" -eq 1 ]; then
        gaming_quality_score=$((gaming_quality_score - 10))
        log_warning "🧪 Low test coverage: Gaming reliability risk"
    fi
    
    # Ensure minimum score
    if [ "$gaming_quality_score" -lt 0 ]; then
        gaming_quality_score=0
    fi
    
    echo "🎯 Gaming Quality Score: $gaming_quality_score/100"
    
    # Performance impact assessment
    if [ "$gaming_quality_score" -ge 80 ]; then
        log_success "✅ Excellent gaming code quality - <${PERFORMANCE_TARGET_MS}ms target achievable"
    elif [ "$gaming_quality_score" -ge 60 ]; then
        log_warning "⚠️ Good gaming code quality - Monitor performance closely"
    else
        log_error "❌ Gaming code quality needs improvement - Performance target at risk"
    fi
}

# Get gaming-specific issues
get_gaming_issues() {
    log_gaming "🔍 Analyzing gaming-specific code issues..."
    
    # Get issues from SonarQube
    local issues_response=$(curl -s -H "Authorization: Bearer $SONARQUBE_TOKEN" \
        "$SONARQUBE_URL/api/issues/search?componentKeys=$GAMING_PROJECT_KEY&types=BUG,VULNERABILITY,CODE_SMELL&ps=100")
    
    if [ $? -ne 0 ]; then
        log_error "❌ Failed to fetch gaming issues"
        return 1
    fi
    
    # Analyze for gaming-specific patterns
    local gaming_issues_file="$METRICS_DATA_DIR/gaming-issues-$(date +%s).json"
    echo "$issues_response" | jq '
    {
        "total_issues": .total,
        "gaming_categories": {
            "performance_critical": [.issues[] | select(.message | test("performance|timeout|slow|cache"))],
            "riot_api_related": [.issues[] | select(.message | test("api|http|request|rate.?limit"))],
            "gaming_logic": [.issues[] | select(.component | test("gaming|match|player|analytics"))],
            "security_gaming": [.issues[] | select(.type == "VULNERABILITY" and (.component | test("gaming|player|riot")))]
        }
    }' > "$gaming_issues_file"
    
    # Report gaming-specific findings
    local perf_issues=$(jq '.gaming_categories.performance_critical | length' "$gaming_issues_file")
    local api_issues=$(jq '.gaming_categories.riot_api_related | length' "$gaming_issues_file")
    local gaming_logic_issues=$(jq '.gaming_categories.gaming_logic | length' "$gaming_issues_file")
    local security_gaming_issues=$(jq '.gaming_categories.security_gaming | length' "$gaming_issues_file")
    
    echo ""
    log_info "🎮 Gaming-Specific Issues Analysis"
    echo "=================================="
    echo "⚡ Performance Critical: $perf_issues issues"
    echo "🔗 Riot API Related: $api_issues issues"
    echo "🎯 Gaming Logic: $gaming_logic_issues issues"
    echo "🔒 Gaming Security: $security_gaming_issues issues"
    
    # Recommendations based on gaming issues
    if [ "$perf_issues" -gt 0 ]; then
        log_warning "⚠️ Performance issues detected - May impact <${PERFORMANCE_TARGET_MS}ms target"
    fi
    
    if [ "$api_issues" -gt 0 ]; then
        log_warning "⚠️ Riot API issues detected - Check rate limiting and error handling"
    fi
    
    if [ "$security_gaming_issues" -gt 0 ]; then
        log_error "🔒 Gaming security issues - Player data at risk"
    fi
}

# Generate gaming dashboard data
generate_dashboard_data() {
    log_gaming "📈 Generating Herald.lol gaming dashboard data..."
    
    local dashboard_file="$METRICS_DATA_DIR/gaming-dashboard-$(date +%s).json"
    
    # Combine metrics from recent collections
    jq -s '
    {
        "herald_gaming_dashboard": {
            "last_updated": now,
            "performance_target_ms": '"$PERFORMANCE_TARGET_MS"',
            "metrics_history": .,
            "gaming_trends": {
                "code_growth": [.[] | .metrics.code_lines],
                "quality_trend": [.[] | (100 - (.metrics.bugs + .metrics.vulnerabilities) * 2)],
                "coverage_trend": [.[] | .metrics.coverage_percent]
            },
            "gaming_alerts": {
                "performance_risk": (.[0].metrics.bugs > 5),
                "security_risk": (.[0].metrics.vulnerabilities > 0),
                "coverage_low": (.[0].metrics.coverage_percent < 70)
            }
        }
    }' "$METRICS_DATA_DIR"/gaming-metrics-*.json > "$dashboard_file" 2>/dev/null || {
        echo '{"herald_gaming_dashboard": {"error": "No metrics data available"}}' > "$dashboard_file"
    }
    
    log_success "✅ Gaming dashboard data: $dashboard_file"
}

# Send gaming metrics to external monitoring
send_metrics() {
    log_gaming "📡 Sending Herald.lol gaming metrics..."
    
    # Send to CloudWatch (if AWS CLI available)
    if command -v aws >/dev/null 2>&1; then
        local latest_metrics=$(ls -t "$METRICS_DATA_DIR"/gaming-metrics-*.json | head -1)
        
        if [ -f "$latest_metrics" ]; then
            local bugs=$(jq -r '.metrics.bugs' "$latest_metrics")
            local vulnerabilities=$(jq -r '.metrics.vulnerabilities' "$latest_metrics")
            local coverage=$(jq -r '.metrics.coverage_percent' "$latest_metrics")
            
            # Send gaming-specific CloudWatch metrics
            aws cloudwatch put-metric-data \
                --namespace "Herald/Gaming/CodeQuality" \
                --metric-data MetricName=GamingBugs,Value="$bugs",Unit=Count \
                --metric-data MetricName=GamingVulnerabilities,Value="$vulnerabilities",Unit=Count \
                --metric-data MetricName=GamingCoverage,Value="$coverage",Unit=Percent \
                2>/dev/null && log_success "✅ Metrics sent to CloudWatch" || log_warning "⚠️ CloudWatch metrics failed"
        fi
    fi
}

# Main metrics collection workflow
main() {
    log_gaming "🚀 Starting Herald.lol gaming metrics collection workflow..."
    
    init_metrics
    collect_gaming_metrics
    analyze_gaming_impact
    get_gaming_issues
    generate_dashboard_data
    send_metrics
    
    log_success "✅ Herald.lol gaming metrics collection completed!"
    log_info "📊 Next collection in 1 hour"
}

# Handle interruptions gracefully
trap 'echo ""; log_info "🎮 Gaming metrics collection interrupted"; exit 0' SIGINT SIGTERM

# Run main workflow
main "$@"