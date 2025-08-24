#!/bin/bash

# Herald.lol Gaming Analytics - Deployment Monitoring Script
# Monitors deployment health and performance during blue-green switches

set -euo pipefail

# Configuration
ANALYTICS_PERFORMANCE_TARGET=5000  # 5s target for gaming analytics
MONITORING_DURATION=300            # 5 minutes of monitoring after switch
CHECK_INTERVAL=10                  # Check every 10 seconds

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

# Get current environment
get_current_environment() {
    kubectl get service herald-gaming-analytics -n herald-production -o jsonpath='{.spec.selector.version}' 2>/dev/null || echo "none"
}

# Check deployment health
check_deployment_health() {
    local env=$1
    local namespace="herald-$env"
    
    # Check pod status
    local ready_pods=$(kubectl get pods -n "$namespace" -l "app=herald-gaming-analytics,version=$env" --field-selector=status.phase=Running --no-headers 2>/dev/null | wc -l || echo "0")
    local total_pods=$(kubectl get deployment "herald-gaming-analytics-$env" -n "$namespace" -o jsonpath='{.spec.replicas}' 2>/dev/null || echo "0")
    
    if [[ "$ready_pods" -eq "$total_pods" && "$ready_pods" -gt 0 ]]; then
        echo "✅"
    else
        echo "❌ ($ready_pods/$total_pods ready)"
    fi
}

# Test gaming analytics performance
test_analytics_performance() {
    local env=$1
    local namespace="herald-$env"
    local service="herald-gaming-analytics-$env"
    
    # Get a pod to test against
    local pod=$(kubectl get pods -n "$namespace" -l "app=herald-gaming-analytics,version=$env" --field-selector=status.phase=Running -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
    
    if [[ -z "$pod" ]]; then
        echo "No running pods"
        return 1
    fi
    
    # Test health endpoint response time
    local start_time=$(date +%s%3N)
    if kubectl exec -n "$namespace" "$pod" -- curl -sf http://localhost:8080/health >/dev/null 2>&1; then
        local end_time=$(date +%s%3N)
        local response_time=$((end_time - start_time))
        
        if [[ $response_time -le $ANALYTICS_PERFORMANCE_TARGET ]]; then
            echo "✅ ${response_time}ms"
        else
            echo "⚠️ ${response_time}ms (>${ANALYTICS_PERFORMANCE_TARGET}ms target)"
        fi
    else
        echo "❌ Health check failed"
        return 1
    fi
}

# Test gRPC connectivity
test_grpc_health() {
    local env=$1
    local namespace="herald-$env"
    
    local pod=$(kubectl get pods -n "$namespace" -l "app=herald-gaming-analytics,version=$env" --field-selector=status.phase=Running -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
    
    if [[ -z "$pod" ]]; then
        echo "No running pods"
        return 1
    fi
    
    if kubectl exec -n "$namespace" "$pod" -- nc -z localhost 50051 >/dev/null 2>&1; then
        echo "✅ gRPC healthy"
    else
        echo "❌ gRPC failed"
        return 1
    fi
}

# Get resource usage
get_resource_usage() {
    local env=$1
    local namespace="herald-$env"
    
    # Get average CPU and memory usage
    local metrics=$(kubectl top pods -n "$namespace" -l "app=herald-gaming-analytics,version=$env" --no-headers 2>/dev/null | awk '{cpu+=$2; mem+=$3; count++} END {if(count>0) printf "CPU: %.0fm, Mem: %.0fMi", cpu/count, mem/count; else print "No metrics"}' || echo "Metrics unavailable")
    echo "$metrics"
}

# Monitor deployment
monitor_deployment() {
    local duration=${1:-$MONITORING_DURATION}
    local current_env=$(get_current_environment)
    
    if [[ "$current_env" == "none" ]]; then
        log_error "No active deployment found"
        return 1
    fi
    
    log_info "🎮 Monitoring Herald.lol gaming analytics deployment"
    log_info "📊 Active environment: $current_env"
    log_info "⏱️  Monitoring for ${duration}s (target: <${ANALYTICS_PERFORMANCE_TARGET}ms)"
    log_info "🔄 Check interval: ${CHECK_INTERVAL}s"
    echo ""
    
    local start_time=$(date +%s)
    local end_time=$((start_time + duration))
    local check_count=0
    local failed_checks=0
    local performance_violations=0
    
    # Print header
    printf "%-10s %-15s %-15s %-20s %-20s %-30s\n" "Time" "Health" "Performance" "gRPC" "Resources" "Status"
    printf "%-10s %-15s %-15s %-20s %-20s %-30s\n" "----------" "---------------" "---------------" "--------------------" "--------------------" "------------------------------"
    
    while [[ $(date +%s) -lt $end_time ]]; do
        local current_time=$(date +'%H:%M:%S')
        local health_status=$(check_deployment_health "$current_env")
        local performance_status=$(test_analytics_performance "$current_env")
        local grpc_status=$(test_grpc_health "$current_env")
        local resource_usage=$(get_resource_usage "$current_env")
        
        check_count=$((check_count + 1))
        
        # Determine overall status
        local overall_status="🟢 OK"
        if [[ "$health_status" == *"❌"* ]]; then
            overall_status="🔴 HEALTH FAIL"
            failed_checks=$((failed_checks + 1))
        elif [[ "$performance_status" == *"❌"* ]]; then
            overall_status="🔴 PERF FAIL"
            failed_checks=$((failed_checks + 1))
        elif [[ "$grpc_status" == *"❌"* ]]; then
            overall_status="🔴 GRPC FAIL"
            failed_checks=$((failed_checks + 1))
        elif [[ "$performance_status" == *"⚠️"* ]]; then
            overall_status="🟡 PERF SLOW"
            performance_violations=$((performance_violations + 1))
        fi
        
        printf "%-10s %-15s %-15s %-20s %-20s %-30s\n" "$current_time" "$health_status" "$performance_status" "$grpc_status" "$resource_usage" "$overall_status"
        
        sleep "$CHECK_INTERVAL"
    done
    
    echo ""
    log_info "📈 Monitoring Summary:"
    log_info "  Total checks: $check_count"
    log_info "  Failed checks: $failed_checks"
    log_info "  Performance violations: $performance_violations"
    log_info "  Success rate: $(awk "BEGIN {printf \"%.1f%%\", (($check_count-$failed_checks)/$check_count)*100}")"
    
    if [[ $failed_checks -eq 0 ]]; then
        log_success "✅ All health checks passed!"
    else
        log_warning "⚠️  $failed_checks health check failures detected"
    fi
    
    if [[ $performance_violations -eq 0 ]]; then
        log_success "🎯 All performance checks met <${ANALYTICS_PERFORMANCE_TARGET}ms target"
    else
        log_warning "⚠️  $performance_violations performance violations detected"
    fi
}

# Real-time monitoring
realtime_monitor() {
    log_info "🔄 Starting real-time monitoring (Press Ctrl+C to stop)"
    
    while true; do
        clear
        echo "🎮 Herald.lol Gaming Analytics - Real-time Monitor"
        echo "=================================================="
        echo "⚡ Performance Target: <${ANALYTICS_PERFORMANCE_TARGET}ms"
        echo "🕐 Last Update: $(date)"
        echo ""
        
        local current_env=$(get_current_environment)
        if [[ "$current_env" != "none" ]]; then
            echo "📊 Active Environment: $current_env"
            echo "🏥 Health: $(check_deployment_health "$current_env")"
            echo "⚡ Performance: $(test_analytics_performance "$current_env")"
            echo "🔌 gRPC: $(test_grpc_health "$current_env")"
            echo "💾 Resources: $(get_resource_usage "$current_env")"
            
            # Show both environments
            echo ""
            echo "🔍 Both Environments:"
            for env in blue green; do
                local health=$(check_deployment_health "$env")
                local resources=$(get_resource_usage "$env")
                local status_indicator="  "
                if [[ "$env" == "$current_env" ]]; then
                    status_indicator="🟢"
                fi
                echo "  $status_indicator $env: $health | $resources"
            done
        else
            echo "❌ No active deployment found"
        fi
        
        echo ""
        echo "Press Ctrl+C to stop monitoring"
        sleep 5
    done
}

# Usage
usage() {
    cat << EOF
🎮 Herald.lol Gaming Analytics Deployment Monitor

Usage: $0 [OPTIONS]

OPTIONS:
    -m, --monitor DURATION     Monitor deployment for specified duration (seconds)
    -r, --realtime            Real-time monitoring (interactive)
    -c, --check               Single health check
    -h, --help                Show this help message

EXAMPLES:
    # Monitor for 5 minutes after deployment
    $0 --monitor 300
    
    # Real-time monitoring
    $0 --realtime
    
    # Single health check
    $0 --check

🎯 Performance Target: <${ANALYTICS_PERFORMANCE_TARGET}ms for gaming analytics
EOF
}

# Main function
main() {
    case "${1:-}" in
        -m|--monitor)
            monitor_deployment "${2:-$MONITORING_DURATION}"
            ;;
        -r|--realtime)
            realtime_monitor
            ;;
        -c|--check)
            local current_env=$(get_current_environment)
            if [[ "$current_env" != "none" ]]; then
                echo "🎮 Herald.lol Gaming Analytics Health Check"
                echo "Active Environment: $current_env"
                echo "Health: $(check_deployment_health "$current_env")"
                echo "Performance: $(test_analytics_performance "$current_env")"
                echo "gRPC: $(test_grpc_health "$current_env")"
                echo "Resources: $(get_resource_usage "$current_env")"
            else
                log_error "No active deployment found"
                exit 1
            fi
            ;;
        -h|--help)
            usage
            ;;
        *)
            usage
            ;;
    esac
}

# Handle Ctrl+C gracefully
trap 'echo ""; log_info "Monitoring stopped"; exit 0' SIGINT

# Run main function
main "$@"