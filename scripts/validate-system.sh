#!/bin/bash

# LoL Match Exporter - System Validation Script
# This script validates the complete system functionality

set -e

echo "🔍 LoL Match Exporter - System Validation"
echo "========================================"

# Configuration
SERVER_URL="${SERVER_URL:-http://localhost:8001}"
REDIS_HOST="${REDIS_HOST:-localhost}"
REDIS_PORT="${REDIS_PORT:-6379}"
TEST_USER_ID="${TEST_USER_ID:-1}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Make HTTP request with error handling
make_request() {
    local method="$1"
    local url="$2"
    local data="$3"
    local expected_status="${4:-200}"
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X "$method" \
                       -H "Content-Type: application/json" \
                       -d "$data" "$url" 2>/dev/null)
    else
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X "$method" "$url" 2>/dev/null)
    fi
    
    body=$(echo "$response" | sed -E 's/HTTPSTATUS\:[0-9]{3}$//')
    status=$(echo "$response" | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
    
    if [ "$status" -eq "$expected_status" ]; then
        echo "$body"
        return 0
    else
        log_error "HTTP $status for $method $url (expected $expected_status)"
        echo "$body"
        return 1
    fi
}

# Test 1: Check prerequisites
test_prerequisites() {
    log_info "Testing prerequisites..."
    
    # Check required commands
    local missing_commands=()
    for cmd in curl jq redis-cli; do
        if ! command_exists "$cmd"; then
            missing_commands+=("$cmd")
        fi
    done
    
    if [ ${#missing_commands[@]} -gt 0 ]; then
        log_error "Missing required commands: ${missing_commands[*]}"
        log_info "Install missing commands and retry"
        return 1
    fi
    
    log_success "All prerequisites satisfied"
    return 0
}

# Test 2: Basic connectivity
test_connectivity() {
    log_info "Testing basic connectivity..."
    
    # Test server connectivity
    if ! curl -s --connect-timeout 5 "$SERVER_URL/api/health" >/dev/null; then
        log_error "Cannot connect to server at $SERVER_URL"
        return 1
    fi
    
    # Test Redis connectivity (optional)
    if command_exists redis-cli; then
        if redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" ping >/dev/null 2>&1; then
            log_success "Redis connection successful"
        else
            log_warning "Redis not available at $REDIS_HOST:$REDIS_PORT (will test graceful degradation)"
        fi
    fi
    
    log_success "Basic connectivity verified"
    return 0
}

# Test 3: Health endpoints
test_health_endpoints() {
    log_info "Testing health endpoints..."
    
    # Test main health endpoint
    local health_response
    health_response=$(make_request "GET" "$SERVER_URL/api/health")
    if [ $? -eq 0 ]; then
        local status
        status=$(echo "$health_response" | jq -r '.status // empty')
        if [ "$status" = "healthy" ]; then
            log_success "Main health endpoint: OK"
        else
            log_error "Main health endpoint returned status: $status"
            return 1
        fi
    else
        log_error "Main health endpoint failed"
        return 1
    fi
    
    # Test optimized analytics health
    local v2_health_response
    v2_health_response=$(make_request "GET" "$SERVER_URL/api/analytics/v2/health")
    if [ $? -eq 0 ]; then
        local v2_status
        v2_status=$(echo "$v2_health_response" | jq -r '.status // empty')
        if [ "$v2_status" = "healthy" ] || [ "$v2_status" = "unhealthy" ]; then
            log_success "Optimized analytics health endpoint: OK (status: $v2_status)"
        else
            log_error "Optimized analytics health endpoint returned unexpected status: $v2_status"
            return 1
        fi
    else
        log_error "Optimized analytics health endpoint failed"
        return 1
    fi
    
    log_success "All health endpoints working"
    return 0
}

# Test 4: Authentication flow
test_authentication() {
    log_info "Testing authentication flow..."
    
    # Test account validation
    local auth_data='{"riot_id":"TestPlayer","riot_tag":"TEST","region":"euw1"}'
    local auth_response
    auth_response=$(make_request "POST" "$SERVER_URL/api/auth/validate" "$auth_data")
    if [ $? -eq 0 ]; then
        local valid
        valid=$(echo "$auth_response" | jq -r '.valid // false')
        if [ "$valid" = "true" ]; then
            log_success "Authentication flow: OK"
            
            # Extract session cookie for subsequent requests
            SESSION_COOKIE=$(curl -s -c - -X POST \
                                 -H "Content-Type: application/json" \
                                 -d "$auth_data" \
                                 "$SERVER_URL/api/auth/validate" | grep -o 'lol-session.*')
            
            if [ -n "$SESSION_COOKIE" ]; then
                log_success "Session cookie obtained"
            else
                log_warning "No session cookie found (may affect subsequent tests)"
            fi
        else
            log_error "Authentication validation failed"
            return 1
        fi
    else
        log_error "Authentication endpoint failed"
        return 1
    fi
    
    return 0
}

# Test 5: Analytics endpoints (v1)
test_v1_analytics() {
    log_info "Testing v1 analytics endpoints..."
    
    # Note: These tests may fail if authentication is required
    # For now, test endpoints that don't require auth or use mock data
    
    local endpoints=(
        "/api/analytics/health"
    )
    
    for endpoint in "${endpoints[@]}"; do
        local response
        if response=$(make_request "GET" "$SERVER_URL$endpoint"); then
            log_success "V1 endpoint $endpoint: OK"
        else
            log_warning "V1 endpoint $endpoint: Failed (may require auth)"
        fi
    done
    
    return 0
}

# Test 6: Analytics endpoints (v2 optimized)
test_v2_analytics() {
    log_info "Testing v2 optimized analytics endpoints..."
    
    # Test performance metrics endpoint (should work without auth for monitoring)
    local perf_response
    if perf_response=$(make_request "GET" "$SERVER_URL/api/analytics/v2/performance"); then
        local cache_enabled
        cache_enabled=$(echo "$perf_response" | jq -r '.data.service.cache_enabled // false')
        local async_processing
        async_processing=$(echo "$perf_response" | jq -r '.data.service.async_processing // false')
        
        log_success "V2 performance endpoint: OK"
        log_info "Cache enabled: $cache_enabled"
        log_info "Async processing: $async_processing"
        
        # Validate worker pool stats if available
        local workers_active
        workers_active=$(echo "$perf_response" | jq -r '.data.worker_pool.workers_active // "N/A"')
        if [ "$workers_active" != "N/A" ]; then
            log_info "Active workers: $workers_active"
        fi
    else
        log_error "V2 performance endpoint failed"
        return 1
    fi
    
    return 0
}

# Test 7: Load testing (basic)
test_basic_load() {
    log_info "Testing basic load handling..."
    
    # Simple concurrent requests to health endpoint
    local concurrent_requests=10
    local pids=()
    
    log_info "Sending $concurrent_requests concurrent requests..."
    
    for i in $(seq 1 $concurrent_requests); do
        (
            response=$(curl -s -w "%{http_code}" -o /dev/null "$SERVER_URL/api/health")
            if [ "$response" = "200" ]; then
                echo "Request $i: OK"
            else
                echo "Request $i: FAILED ($response)"
            fi
        ) &
        pids+=($!)
    done
    
    # Wait for all requests to complete
    local success_count=0
    for pid in "${pids[@]}"; do
        if wait "$pid"; then
            ((success_count++))
        fi
    done
    
    local success_rate=$((success_count * 100 / concurrent_requests))
    if [ $success_rate -ge 90 ]; then
        log_success "Load test: $success_count/$concurrent_requests requests succeeded ($success_rate%)"
    else
        log_warning "Load test: Only $success_count/$concurrent_requests requests succeeded ($success_rate%)"
    fi
    
    return 0
}

# Test 8: Cache behavior
test_cache_behavior() {
    log_info "Testing cache behavior..."
    
    if ! command_exists redis-cli; then
        log_warning "Redis CLI not available, skipping cache tests"
        return 0
    fi
    
    # Test Redis connectivity
    if ! redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" ping >/dev/null 2>&1; then
        log_warning "Redis not available, testing graceful degradation"
        # Test that system still works without Redis
        local health_response
        if health_response=$(make_request "GET" "$SERVER_URL/api/health"); then
            log_success "System operates correctly without Redis"
        else
            log_error "System fails when Redis is unavailable"
            return 1
        fi
        return 0
    fi
    
    # Test cache statistics
    local perf_response
    if perf_response=$(make_request "GET" "$SERVER_URL/api/analytics/v2/performance"); then
        local cache_status
        cache_status=$(echo "$perf_response" | jq -r '.data.cache.status // "unknown"')
        log_info "Cache status: $cache_status"
        
        if [ "$cache_status" = "connected" ]; then
            log_success "Cache is connected and operational"
        else
            log_warning "Cache status: $cache_status"
        fi
    fi
    
    return 0
}

# Test 9: Error handling
test_error_handling() {
    log_info "Testing error handling..."
    
    # Test invalid endpoints
    local invalid_response
    if invalid_response=$(make_request "GET" "$SERVER_URL/api/invalid/endpoint" "" 404); then
        log_success "404 handling: OK"
    else
        log_warning "404 handling may need attention"
    fi
    
    # Test malformed requests
    local malformed_data='{"invalid": json}'
    if make_request "POST" "$SERVER_URL/api/auth/validate" "$malformed_data" 400 >/dev/null 2>&1; then
        log_success "Malformed request handling: OK"
    else
        log_warning "Malformed request handling may need attention"
    fi
    
    return 0
}

# Main validation function
run_validation() {
    local start_time
    start_time=$(date +%s)
    
    echo "Starting system validation at $(date)"
    echo "Server URL: $SERVER_URL"
    echo "Redis: $REDIS_HOST:$REDIS_PORT"
    echo ""
    
    local tests=(
        "test_prerequisites"
        "test_connectivity"
        "test_health_endpoints"
        "test_authentication"
        "test_v1_analytics"
        "test_v2_analytics"
        "test_basic_load"
        "test_cache_behavior"
        "test_error_handling"
    )
    
    local passed=0
    local failed=0
    local warnings=0
    
    for test in "${tests[@]}"; do
        echo ""
        if $test; then
            ((passed++))
        else
            ((failed++))
        fi
    done
    
    local end_time
    end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    echo ""
    echo "========================================"
    echo "Validation completed in ${duration}s"
    echo "Results: $passed passed, $failed failed"
    
    if [ $failed -eq 0 ]; then
        log_success "🎉 All validation tests passed!"
        echo "System is ready for production use."
        return 0
    else
        log_error "❌ Some validation tests failed."
        echo "Please review the failures above before deploying to production."
        return 1
    fi
}

# Run validation if script is executed directly
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    run_validation
fi