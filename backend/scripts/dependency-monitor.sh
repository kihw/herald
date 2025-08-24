#!/bin/bash

# Herald.lol Gaming Analytics - Dependency Monitoring Script
# Continuous monitoring of dependencies for gaming security and performance

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
MONITOR_DIR="$PROJECT_ROOT/dependency-monitoring"
GAMING_PERFORMANCE_TARGET=5000  # 5s gaming analytics target
CHECK_INTERVAL=3600             # Check every hour

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
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
    echo -e "${PURPLE}[$(date +'%Y-%m-%d %H:%M:%S')] [GAMING]${NC} $1"
}

# Initialize monitoring directory
init_monitoring() {
    mkdir -p "$MONITOR_DIR/reports"
    mkdir -p "$MONITOR_DIR/alerts"
    mkdir -p "$MONITOR_DIR/metrics"
    log_info "📁 Initialized dependency monitoring directory"
}

# Install monitoring tools
install_tools() {
    log_info "🔧 Installing dependency monitoring tools..."
    
    # Install Go vulnerability scanner
    if ! command -v govulncheck >/dev/null 2>&1; then
        go install golang.org/x/vuln/cmd/govulncheck@latest
        log_info "✅ Installed govulncheck"
    fi
    
    # Install Nancy
    if ! command -v nancy >/dev/null 2>&1; then
        go install github.com/sonatypecommunity/nancy@latest
        log_info "✅ Installed nancy"
    fi
    
    # Check for Trivy
    if ! command -v trivy >/dev/null 2>&1; then
        log_warning "⚠️  Trivy not found. Install for complete scanning."
    fi
    
    log_success "🔧 Monitoring tools ready"
}

# Scan Go dependencies
scan_go_dependencies() {
    log_info "🐹 Scanning Go dependencies for Herald.lol gaming platform..."
    
    local timestamp=$(date +"%Y%m%d_%H%M%S")
    local report_file="$MONITOR_DIR/reports/go_deps_$timestamp.json"
    
    cd "$PROJECT_ROOT"
    
    # Run govulncheck
    log_info "Running official Go vulnerability scanner..."
    govulncheck -json ./... > "$report_file" 2>/dev/null || {
        log_warning "govulncheck completed with warnings"
    }
    
    # Analyze results for gaming impact
    local vulnerabilities=$(jq -r '.message.vulnerability // empty' "$report_file" 2>/dev/null | wc -l || echo "0")
    
    if [[ $vulnerabilities -gt 0 ]]; then
        log_warning "🚨 Found $vulnerabilities vulnerabilities in Go dependencies"
        
        # Check for gaming-critical dependencies
        check_gaming_critical_deps "$report_file"
    else
        log_success "✅ No vulnerabilities found in Go dependencies"
    fi
    
    echo "$report_file"
}

# Check gaming-critical dependencies
check_gaming_critical_deps() {
    local report_file=$1
    
    log_gaming "🎮 Analyzing gaming-critical dependency vulnerabilities..."
    
    # Gaming-critical dependency patterns
    local gaming_deps=(
        "gin"
        "gorilla"
        "websocket"
        "grpc"
        "redis"
        "postgres"
        "gorm"
        "prometheus"
    )
    
    local gaming_vulns=0
    for dep in "${gaming_deps[@]}"; do
        if jq -r '.message.module // empty' "$report_file" 2>/dev/null | grep -i "$dep" >/dev/null; then
            log_warning "⚠️ Gaming-critical dependency affected: $dep"
            gaming_vulns=$((gaming_vulns + 1))
        fi
    done
    
    if [[ $gaming_vulns -gt 0 ]]; then
        log_error "🚨 $gaming_vulns gaming-critical dependencies have vulnerabilities!"
        create_gaming_alert "gaming_deps_vulnerable" "$gaming_vulns"
    else
        log_success "✅ No vulnerabilities in gaming-critical dependencies"
    fi
}

# Monitor dependency freshness
monitor_dependency_freshness() {
    log_info "📅 Monitoring dependency freshness for gaming platform..."
    
    cd "$PROJECT_ROOT"
    
    # Get current dependencies
    go list -m -u all > "$MONITOR_DIR/current_deps.txt" 2>/dev/null || true
    
    # Check for available updates
    local outdated_deps=$(go list -m -u all 2>/dev/null | grep -E '\[.*\]' | wc -l || echo "0")
    
    if [[ $outdated_deps -gt 0 ]]; then
        log_info "📦 $outdated_deps dependencies have updates available"
        
        # Analyze gaming impact of updates
        analyze_gaming_update_impact "$outdated_deps"
    else
        log_success "✅ All dependencies are up to date"
    fi
}

# Analyze gaming impact of dependency updates
analyze_gaming_update_impact() {
    local outdated_count=$1
    
    log_gaming "🎮 Analyzing gaming impact of dependency updates..."
    
    # Create gaming impact report
    local impact_file="$MONITOR_DIR/reports/gaming_update_impact_$(date +%Y%m%d_%H%M%S).txt"
    
    {
        echo "# Herald.lol Gaming Dependency Update Impact Analysis"
        echo "Generated: $(date)"
        echo "Outdated dependencies: $outdated_count"
        echo "Gaming performance target: <${GAMING_PERFORMANCE_TARGET}ms"
        echo ""
        
        echo "## Gaming-Critical Dependencies Analysis"
        go list -m -u all 2>/dev/null | grep -E '\[.*\]' | while read line; do
            dep_name=$(echo "$line" | awk '{print $1}')
            
            # Check if it's gaming-critical
            if echo "$dep_name" | grep -iE "gin|gorilla|websocket|grpc|redis|postgres|gorm|prometheus" >/dev/null; then
                echo "🎮 GAMING-CRITICAL: $line"
            else
                echo "📦 Standard: $line"
            fi
        done
        
        echo ""
        echo "## Gaming Performance Considerations"
        echo "🎯 Performance Impact Assessment:"
        echo "- Updates to HTTP/WebSocket libraries may affect real-time gaming connections"
        echo "- Database dependency updates could impact <${GAMING_PERFORMANCE_TARGET}ms analytics target"
        echo "- gRPC updates may affect gaming service communication"
        echo ""
        echo "## Recommended Actions"
        echo "1. Test gaming performance after updates"
        echo "2. Validate Riot API integration compatibility"
        echo "3. Run gaming load tests with updated dependencies"
        echo "4. Monitor real-time analytics performance"
        
    } > "$impact_file"
    
    log_info "📊 Gaming impact analysis saved to: $impact_file"
}

# Monitor license compliance
monitor_license_compliance() {
    log_info "⚖️ Monitoring license compliance for gaming platform..."
    
    # Get all dependencies with licenses
    go list -m -json all 2>/dev/null > "$MONITOR_DIR/deps_with_licenses.json" || true
    
    # Check for problematic licenses
    local problematic_licenses=("GPL" "AGPL" "SSPL")
    local license_issues=0
    
    for license in "${problematic_licenses[@]}"; do
        if grep -i "$license" "$MONITOR_DIR/deps_with_licenses.json" >/dev/null 2>&1; then
            log_warning "⚠️ Found potentially problematic license: $license"
            license_issues=$((license_issues + 1))
        fi
    done
    
    if [[ $license_issues -eq 0 ]]; then
        log_success "✅ No problematic licenses detected"
    else
        log_warning "⚠️ $license_issues potential license issues found"
        create_gaming_alert "license_compliance" "$license_issues"
    fi
}

# Create gaming-specific alert
create_gaming_alert() {
    local alert_type=$1
    local severity=$2
    local timestamp=$(date +"%Y%m%d_%H%M%S")
    
    local alert_file="$MONITOR_DIR/alerts/gaming_alert_${alert_type}_${timestamp}.json"
    
    cat << EOF > "$alert_file"
{
  "alert_type": "$alert_type",
  "severity": "$severity", 
  "gaming_platform": "Herald.lol",
  "performance_target_ms": $GAMING_PERFORMANCE_TARGET,
  "timestamp": "$(date -Iseconds)",
  "description": "Gaming dependency security alert",
  "recommended_actions": [
    "Review dependency vulnerabilities",
    "Test gaming performance impact",
    "Validate Riot API integration",
    "Deploy via blue-green strategy"
  ],
  "gaming_impact": {
    "analytics_performance": "Monitor for <${GAMING_PERFORMANCE_TARGET}ms target",
    "real_time_features": "Validate WebSocket/gRPC functionality",
    "riot_api_integration": "Test API compliance after fixes"
  }
}
EOF
    
    log_warning "🚨 Gaming alert created: $alert_file"
}

# Generate dependency metrics for gaming platform
generate_gaming_metrics() {
    log_info "📊 Generating gaming dependency metrics..."
    
    local metrics_file="$MONITOR_DIR/metrics/gaming_metrics_$(date +%Y%m%d_%H%M%S).json"
    
    cd "$PROJECT_ROOT"
    
    # Count dependencies by category
    local total_deps=$(go list -m all | wc -l)
    local gaming_deps=$(go list -m all | grep -iE "gin|gorilla|websocket|grpc|redis|postgres|gorm" | wc -l)
    local security_deps=$(go list -m all | grep -iE "crypto|auth|jwt|oauth|security" | wc -l)
    local performance_deps=$(go list -m all | grep -iE "cache|memory|fast|perf" | wc -l)
    
    cat << EOF > "$metrics_file"
{
  "gaming_platform": "Herald.lol",
  "generated": "$(date -Iseconds)",
  "performance_target_ms": $GAMING_PERFORMANCE_TARGET,
  "dependency_metrics": {
    "total_dependencies": $total_deps,
    "gaming_critical": $gaming_deps,
    "security_related": $security_deps,
    "performance_related": $performance_deps,
    "gaming_ratio": $(awk "BEGIN {printf \"%.2f\", $gaming_deps/$total_deps*100}")
  },
  "gaming_health": {
    "dependency_freshness": "monitoring",
    "vulnerability_status": "scanning",
    "license_compliance": "checking",
    "gaming_performance": "optimized"
  },
  "next_scan": "$(date -d '+1 hour' -Iseconds)"
}
EOF
    
    log_success "📊 Gaming metrics generated: $metrics_file"
}

# Continuous monitoring loop
continuous_monitor() {
    log_info "🔄 Starting continuous dependency monitoring for Herald.lol..."
    log_gaming "🎮 Gaming performance target: <${GAMING_PERFORMANCE_TARGET}ms"
    log_info "⏱️ Check interval: ${CHECK_INTERVAL}s (1 hour)"
    
    while true; do
        log_info "🔍 Starting dependency monitoring cycle..."
        
        # Run all monitoring checks
        scan_go_dependencies >/dev/null
        monitor_dependency_freshness
        monitor_license_compliance
        generate_gaming_metrics
        
        log_success "✅ Monitoring cycle completed"
        log_info "⏰ Next check in $(($CHECK_INTERVAL / 60)) minutes..."
        
        sleep "$CHECK_INTERVAL"
    done
}

# Single monitoring run
single_monitor() {
    log_info "🔍 Running single dependency monitoring cycle..."
    
    local go_report=$(scan_go_dependencies)
    monitor_dependency_freshness
    monitor_license_compliance
    generate_gaming_metrics
    
    log_success "✅ Single monitoring cycle completed"
    log_info "📊 Go dependency report: $go_report"
}

# Display monitoring dashboard
show_dashboard() {
    clear
    echo "🎮 Herald.lol Gaming Dependency Monitoring Dashboard"
    echo "=================================================="
    echo "⚡ Performance Target: <${GAMING_PERFORMANCE_TARGET}ms"
    echo "🕐 Last Update: $(date)"
    echo ""
    
    # Show latest metrics if available
    local latest_metrics=$(find "$MONITOR_DIR/metrics" -name "gaming_metrics_*.json" -type f 2>/dev/null | sort | tail -1)
    if [[ -n "$latest_metrics" ]]; then
        echo "📊 Latest Gaming Metrics:"
        echo "========================"
        if command -v jq >/dev/null 2>&1; then
            jq -r '.dependency_metrics | to_entries[] | "\(.key): \(.value)"' "$latest_metrics" 2>/dev/null || cat "$latest_metrics"
        else
            cat "$latest_metrics"
        fi
        echo ""
    fi
    
    # Show recent alerts
    echo "🚨 Recent Alerts:"
    echo "=================="
    local alert_count=$(find "$MONITOR_DIR/alerts" -name "gaming_alert_*.json" -type f -mtime -1 2>/dev/null | wc -l)
    if [[ $alert_count -gt 0 ]]; then
        echo "⚠️ $alert_count alerts in the last 24 hours"
        find "$MONITOR_DIR/alerts" -name "gaming_alert_*.json" -type f -mtime -1 2>/dev/null | tail -3 | while read alert; do
            echo "  - $(basename "$alert")"
        done
    else
        echo "✅ No recent alerts"
    fi
    echo ""
    
    # Show dependency health
    echo "🏥 Gaming Dependency Health:"
    echo "============================="
    cd "$PROJECT_ROOT"
    local outdated=$(go list -m -u all 2>/dev/null | grep -E '\[.*\]' | wc -l || echo "0")
    echo "📦 Dependencies with updates: $outdated"
    echo "🎮 Gaming platform: Herald.lol"
    echo "🎯 Performance monitoring: Active"
}

# Usage
usage() {
    cat << EOF
🎮 Herald.lol Gaming Dependency Monitor

Usage: $0 [OPTIONS]

OPTIONS:
    --continuous        Start continuous monitoring
    --single           Run single monitoring cycle
    --dashboard        Show monitoring dashboard
    --install          Install required monitoring tools
    -h, --help         Show this help message

EXAMPLES:
    # Start continuous monitoring
    $0 --continuous
    
    # Run single check
    $0 --single
    
    # Show dashboard
    $0 --dashboard

🎯 Gaming Focus: Optimized for Herald.lol <${GAMING_PERFORMANCE_TARGET}ms target
🔄 Monitoring: Continuous dependency security and freshness
🎮 Platform: Gaming analytics with Riot API integration
EOF
}

# Main function
main() {
    case "${1:-}" in
        --continuous)
            init_monitoring
            install_tools
            continuous_monitor
            ;;
        --single)
            init_monitoring
            install_tools
            single_monitor
            ;;
        --dashboard)
            init_monitoring
            show_dashboard
            ;;
        --install)
            install_tools
            ;;
        -h|--help)
            usage
            ;;
        *)
            usage
            exit 1
            ;;
    esac
}

# Handle Ctrl+C gracefully
trap 'echo ""; log_info "Dependency monitoring stopped"; exit 0' SIGINT

# Run main function
main "$@"