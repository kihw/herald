#!/bin/bash

# LoL Match Exporter - Performance Benchmark Script
# This script runs comprehensive performance tests and generates reports

set -e

echo "📊 LoL Match Exporter - Performance Benchmark"
echo "============================================="

# Configuration
SERVER_URL="${SERVER_URL:-http://localhost:8001}"
BENCHMARK_DURATION="${BENCHMARK_DURATION:-60}"
CONCURRENT_USERS="${CONCURRENT_USERS:-50}"
RAMP_UP_TIME="${RAMP_UP_TIME:-10}"
OUTPUT_DIR="${OUTPUT_DIR:-./benchmark-results}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Helper functions
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

log_benchmark() {
    echo -e "${CYAN}[BENCHMARK]${NC} $1"
}

# Check dependencies
check_dependencies() {
    log_info "Checking dependencies..."
    
    local missing_deps=()
    for cmd in curl jq ab hey wrk; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            missing_deps+=("$cmd")
        fi
    done
    
    if [ ${#missing_deps[@]} -gt 0 ]; then
        log_warning "Missing optional tools: ${missing_deps[*]}"
        log_info "Some benchmarks may be skipped. Install tools for complete testing."
    fi
    
    # Check curl (required)
    if ! command -v curl >/dev/null 2>&1; then
        log_error "curl is required but not installed"
        exit 1
    fi
    
    log_success "Dependencies checked"
}

# Setup benchmark environment
setup_benchmark() {
    log_info "Setting up benchmark environment..."
    
    # Create output directory
    mkdir -p "$OUTPUT_DIR"
    
    # Create timestamp for this benchmark run
    TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
    REPORT_FILE="$OUTPUT_DIR/benchmark_report_$TIMESTAMP.txt"
    
    # Test server connectivity
    if ! curl -s --connect-timeout 5 "$SERVER_URL/api/health" >/dev/null; then
        log_error "Cannot connect to server at $SERVER_URL"
        exit 1
    fi
    
    log_success "Benchmark environment ready"
}

# Benchmark 1: Basic health check performance
benchmark_health_check() {
    log_benchmark "Running health check performance test..."
    
    local endpoint="$SERVER_URL/api/health"
    local requests=1000
    local concurrency=10
    
    if command -v ab >/dev/null 2>&1; then
        log_info "Using Apache Bench (ab) for health check test"
        ab -n $requests -c $concurrency -g "$OUTPUT_DIR/health_check_gnuplot.dat" "$endpoint" > "$OUTPUT_DIR/health_check_ab.txt" 2>&1
        
        # Extract key metrics
        local rps
        rps=$(grep "Requests per second" "$OUTPUT_DIR/health_check_ab.txt" | awk '{print $4}')
        local mean_time
        mean_time=$(grep "Time per request.*mean" "$OUTPUT_DIR/health_check_ab.txt" | head -1 | awk '{print $4}')
        
        log_success "Health check: $rps req/sec, $mean_time ms average"
        
        echo "Health Check Performance:" >> "$REPORT_FILE"
        echo "  Requests per second: $rps" >> "$REPORT_FILE"
        echo "  Average response time: $mean_time ms" >> "$REPORT_FILE"
        echo "" >> "$REPORT_FILE"
    else
        log_warning "Apache Bench (ab) not available, skipping detailed health check test"
    fi
}

# Benchmark 2: Analytics v2 endpoints
benchmark_analytics_v2() {
    log_benchmark "Running analytics v2 performance test..."
    
    local endpoints=(
        "$SERVER_URL/api/analytics/v2/health"
        "$SERVER_URL/api/analytics/v2/performance"
    )
    
    for endpoint in "${endpoints[@]}"; do
        local endpoint_name
        endpoint_name=$(basename "$endpoint")
        
        log_info "Testing endpoint: $endpoint_name"
        
        if command -v ab >/dev/null 2>&1; then
            ab -n 500 -c 25 "$endpoint" > "$OUTPUT_DIR/analytics_v2_${endpoint_name}_ab.txt" 2>&1
            
            local rps
            rps=$(grep "Requests per second" "$OUTPUT_DIR/analytics_v2_${endpoint_name}_ab.txt" | awk '{print $4}')
            local p95
            p95=$(grep "95%" "$OUTPUT_DIR/analytics_v2_${endpoint_name}_ab.txt" | awk '{print $2}')
            
            log_success "$endpoint_name: $rps req/sec, 95th percentile: $p95 ms"
            
            echo "Analytics V2 - $endpoint_name:" >> "$REPORT_FILE"
            echo "  Requests per second: $rps" >> "$REPORT_FILE"
            echo "  95th percentile: $p95 ms" >> "$REPORT_FILE"
            echo "" >> "$REPORT_FILE"
        fi
    done
}

# Benchmark 3: Concurrent user simulation
benchmark_concurrent_users() {
    log_benchmark "Running concurrent user simulation..."
    
    local endpoint="$SERVER_URL/api/health"
    
    if command -v hey >/dev/null 2>&1; then
        log_info "Using hey for concurrent user simulation"
        hey -n 5000 -c "$CONCURRENT_USERS" -o csv "$endpoint" > "$OUTPUT_DIR/concurrent_users_hey.csv" 2>&1
        
        # Parse hey output for summary
        local total_time
        total_time=$(hey -n 5000 -c "$CONCURRENT_USERS" "$endpoint" 2>&1 | grep "Total time:" | awk '{print $3}')
        local rps
        rps=$(hey -n 5000 -c "$CONCURRENT_USERS" "$endpoint" 2>&1 | grep "Requests/sec:" | awk '{print $2}')
        
        log_success "Concurrent users ($CONCURRENT_USERS): $rps req/sec, total time: $total_time"
        
        echo "Concurrent User Simulation ($CONCURRENT_USERS users):" >> "$REPORT_FILE"
        echo "  Requests per second: $rps" >> "$REPORT_FILE"
        echo "  Total time: $total_time" >> "$REPORT_FILE"
        echo "" >> "$REPORT_FILE"
        
    elif command -v ab >/dev/null 2>&1; then
        log_info "Using Apache Bench for concurrent user simulation"
        ab -n 5000 -c "$CONCURRENT_USERS" "$endpoint" > "$OUTPUT_DIR/concurrent_users_ab.txt" 2>&1
        
        local rps
        rps=$(grep "Requests per second" "$OUTPUT_DIR/concurrent_users_ab.txt" | awk '{print $4}')
        local failed
        failed=$(grep "Failed requests" "$OUTPUT_DIR/concurrent_users_ab.txt" | awk '{print $3}')
        
        log_success "Concurrent users ($CONCURRENT_USERS): $rps req/sec, $failed failed"
        
        echo "Concurrent User Simulation ($CONCURRENT_USERS users):" >> "$REPORT_FILE"
        echo "  Requests per second: $rps" >> "$REPORT_FILE"
        echo "  Failed requests: $failed" >> "$REPORT_FILE"
        echo "" >> "$REPORT_FILE"
    else
        log_warning "No suitable tool for concurrent user simulation"
    fi
}

# Benchmark 4: Sustained load test
benchmark_sustained_load() {
    log_benchmark "Running sustained load test ($BENCHMARK_DURATION seconds)..."
    
    local endpoint="$SERVER_URL/api/health"
    
    if command -v wrk >/dev/null 2>&1; then
        log_info "Using wrk for sustained load test"
        wrk -t4 -c20 -d"${BENCHMARK_DURATION}s" --timeout 10s "$endpoint" > "$OUTPUT_DIR/sustained_load_wrk.txt" 2>&1
        
        # Extract metrics from wrk output
        local rps
        rps=$(grep "Requests/sec:" "$OUTPUT_DIR/sustained_load_wrk.txt" | awk '{print $2}')
        local latency_avg
        latency_avg=$(grep "Latency" "$OUTPUT_DIR/sustained_load_wrk.txt" | awk '{print $2}')
        local latency_p99
        latency_p99=$(grep "99%" "$OUTPUT_DIR/sustained_load_wrk.txt" | awk '{print $2}')
        
        log_success "Sustained load: $rps req/sec, avg latency: $latency_avg, 99th percentile: $latency_p99"
        
        echo "Sustained Load Test (${BENCHMARK_DURATION}s):" >> "$REPORT_FILE"
        echo "  Requests per second: $rps" >> "$REPORT_FILE"
        echo "  Average latency: $latency_avg" >> "$REPORT_FILE"
        echo "  99th percentile latency: $latency_p99" >> "$REPORT_FILE"
        echo "" >> "$REPORT_FILE"
        
    elif command -v ab >/dev/null 2>&1; then
        log_info "Using Apache Bench for sustained load test"
        # Calculate requests for duration
        local total_requests=$((BENCHMARK_DURATION * 20))  # 20 req/sec baseline
        ab -n $total_requests -c 20 "$endpoint" > "$OUTPUT_DIR/sustained_load_ab.txt" 2>&1
        
        local rps
        rps=$(grep "Requests per second" "$OUTPUT_DIR/sustained_load_ab.txt" | awk '{print $4}')
        local time_per_request
        time_per_request=$(grep "Time per request.*mean" "$OUTPUT_DIR/sustained_load_ab.txt" | head -1 | awk '{print $4}')
        
        log_success "Sustained load: $rps req/sec, avg time: $time_per_request ms"
        
        echo "Sustained Load Test (${BENCHMARK_DURATION}s simulation):" >> "$REPORT_FILE"
        echo "  Requests per second: $rps" >> "$REPORT_FILE"
        echo "  Average time per request: $time_per_request ms" >> "$REPORT_FILE"
        echo "" >> "$REPORT_FILE"
    else
        log_warning "No suitable tool for sustained load test"
    fi
}

# Benchmark 5: Memory and resource usage
benchmark_resource_usage() {
    log_benchmark "Monitoring resource usage during load..."
    
    # Get initial system stats
    local initial_performance
    initial_performance=$(curl -s "$SERVER_URL/api/analytics/v2/performance" 2>/dev/null || echo '{}')
    
    log_info "Initial system state captured"
    
    # Run a moderate load while monitoring
    if command -v ab >/dev/null 2>&1; then
        log_info "Running background load for resource monitoring"
        ab -n 2000 -c 10 "$SERVER_URL/api/health" > /dev/null 2>&1 &
        local load_pid=$!
        
        # Monitor for 30 seconds
        local monitor_duration=30
        local monitor_interval=5
        local monitor_count=$((monitor_duration / monitor_interval))
        
        echo "Resource Usage Monitoring:" >> "$REPORT_FILE"
        echo "  Timestamp | Workers Active | Tasks Processed | Cache Status" >> "$REPORT_FILE"
        
        for i in $(seq 1 $monitor_count); do
            sleep $monitor_interval
            
            local current_performance
            current_performance=$(curl -s "$SERVER_URL/api/analytics/v2/performance" 2>/dev/null || echo '{}')
            
            local workers_active
            workers_active=$(echo "$current_performance" | jq -r '.data.worker_pool.workers_active // "N/A"' 2>/dev/null || echo "N/A")
            local tasks_processed
            tasks_processed=$(echo "$current_performance" | jq -r '.data.worker_pool.tasks_processed // "N/A"' 2>/dev/null || echo "N/A")
            local cache_status
            cache_status=$(echo "$current_performance" | jq -r '.data.cache.status // "N/A"' 2>/dev/null || echo "N/A")
            
            local timestamp
            timestamp=$(date '+%H:%M:%S')
            
            echo "  $timestamp | $workers_active | $tasks_processed | $cache_status" >> "$REPORT_FILE"
        done
        
        # Stop background load
        kill $load_pid 2>/dev/null || true
        wait $load_pid 2>/dev/null || true
        
        log_success "Resource monitoring completed"
    else
        log_warning "Cannot run resource monitoring without load testing tools"
    fi
    
    echo "" >> "$REPORT_FILE"
}

# Benchmark 6: Cache performance test
benchmark_cache_performance() {
    log_benchmark "Testing cache performance..."
    
    # Test cache hit/miss patterns
    local endpoint="$SERVER_URL/api/analytics/v2/performance"
    
    # First request (cache miss)
    log_info "Testing cache miss scenario"
    local start_time
    start_time=$(date +%s%3N)
    curl -s "$endpoint" > /dev/null
    local end_time
    end_time=$(date +%s%3N)
    local first_request_time=$((end_time - start_time))
    
    # Second request (cache hit)
    log_info "Testing cache hit scenario"
    start_time=$(date +%s%3N)
    curl -s "$endpoint" > /dev/null
    end_time=$(date +%s%3N)
    local second_request_time=$((end_time - start_time))
    
    # Calculate cache improvement
    local improvement=0
    if [ $first_request_time -gt 0 ]; then
        improvement=$(echo "scale=2; (($first_request_time - $second_request_time) * 100) / $first_request_time" | bc -l 2>/dev/null || echo "N/A")
    fi
    
    log_success "Cache miss: ${first_request_time}ms, Cache hit: ${second_request_time}ms, Improvement: ${improvement}%"
    
    echo "Cache Performance Test:" >> "$REPORT_FILE"
    echo "  First request (cache miss): ${first_request_time}ms" >> "$REPORT_FILE"
    echo "  Second request (cache hit): ${second_request_time}ms" >> "$REPORT_FILE"
    echo "  Performance improvement: ${improvement}%" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
}

# Generate summary report
generate_summary() {
    log_info "Generating performance summary..."
    
    # Add system information to report
    {
        echo "========================================"
        echo "PERFORMANCE BENCHMARK SUMMARY"
        echo "========================================"
        echo "Timestamp: $(date)"
        echo "Server URL: $SERVER_URL"
        echo "Duration: ${BENCHMARK_DURATION}s"
        echo "Concurrent Users: $CONCURRENT_USERS"
        echo "Output Directory: $OUTPUT_DIR"
        echo ""
        echo "System Information:"
        echo "  OS: $(uname -s)"
        echo "  Architecture: $(uname -m)"
        if command -v nproc >/dev/null 2>&1; then
            echo "  CPU Cores: $(nproc)"
        fi
        echo ""
    } | cat - "$REPORT_FILE" > temp_report && mv temp_report "$REPORT_FILE"
    
    # Add footer
    {
        echo ""
        echo "========================================"
        echo "BENCHMARK COMPLETE"
        echo "========================================"
        echo "Full report available at: $REPORT_FILE"
        echo "Raw data files in: $OUTPUT_DIR"
    } >> "$REPORT_FILE"
    
    log_success "Performance report generated: $REPORT_FILE"
}

# Main benchmark execution
run_benchmark() {
    local start_time
    start_time=$(date +%s)
    
    echo "Starting performance benchmark at $(date)"
    echo "Configuration:"
    echo "  Server: $SERVER_URL"
    echo "  Duration: ${BENCHMARK_DURATION}s"
    echo "  Concurrent Users: $CONCURRENT_USERS"
    echo "  Output: $OUTPUT_DIR"
    echo ""
    
    # Initialize report file
    echo "PERFORMANCE BENCHMARK RESULTS" > "$REPORT_FILE"
    echo "Generated: $(date)" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    
    # Run benchmark tests
    benchmark_health_check
    benchmark_analytics_v2
    benchmark_concurrent_users
    benchmark_sustained_load
    benchmark_resource_usage
    benchmark_cache_performance
    
    # Generate final report
    generate_summary
    
    local end_time
    end_time=$(date +%s)
    local total_duration=$((end_time - start_time))
    
    echo ""
    echo "============================================"
    echo "🎉 Benchmark completed in ${total_duration}s"
    echo "📊 Report: $REPORT_FILE"
    echo "📁 Data: $OUTPUT_DIR"
    echo "============================================"
    
    # Display key findings
    if [ -f "$REPORT_FILE" ]; then
        echo ""
        echo "Key Performance Metrics:"
        grep -E "(Requests per second|Average latency|95th percentile)" "$REPORT_FILE" | head -5
    fi
}

# Main execution
check_dependencies
setup_benchmark

if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    run_benchmark
fi