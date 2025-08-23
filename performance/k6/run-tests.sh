#!/bin/bash

# Herald.lol Performance Testing Runner
# Execute k6 performance tests for gaming analytics platform

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
FRONTEND_URL="${FRONTEND_URL:-http://localhost:3000}"
RESULTS_DIR="./results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Create results directory
mkdir -p "$RESULTS_DIR"

echo -e "${BLUE}🎮 Herald.lol Performance Testing Suite${NC}"
echo -e "${BLUE}=======================================${NC}"
echo ""
echo -e "API Base URL: ${GREEN}$API_BASE_URL${NC}"
echo -e "Frontend URL: ${GREEN}$FRONTEND_URL${NC}"
echo -e "Results Directory: ${GREEN}$RESULTS_DIR${NC}"
echo ""

# Function to run k6 test
run_k6_test() {
    local test_file="$1"
    local test_name="$2"
    local output_file="$RESULTS_DIR/${test_name}_${TIMESTAMP}.json"
    local html_file="$RESULTS_DIR/${test_name}_${TIMESTAMP}.html"
    
    echo -e "${YELLOW}🚀 Running $test_name...${NC}"
    
    # Run k6 test with JSON output
    k6 run \
        --env API_BASE_URL="$API_BASE_URL" \
        --env FRONTEND_URL="$FRONTEND_URL" \
        --out json="$output_file" \
        "$test_file"
    
    local exit_code=$?
    
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✅ $test_name completed successfully${NC}"
    else
        echo -e "${RED}❌ $test_name failed with exit code $exit_code${NC}"
    fi
    
    # Generate HTML report if k6-reporter is available
    if command -v k6-reporter >/dev/null 2>&1; then
        echo -e "${BLUE}📊 Generating HTML report...${NC}"
        k6-reporter "$output_file" "$html_file"
    fi
    
    echo ""
    return $exit_code
}

# Function to check system health before tests
check_system_health() {
    echo -e "${BLUE}🏥 Checking system health...${NC}"
    
    # Check API health
    if curl -f -s "$API_BASE_URL/health" >/dev/null; then
        echo -e "${GREEN}✅ API health check passed${NC}"
    else
        echo -e "${RED}❌ API health check failed${NC}"
        echo -e "${YELLOW}⚠️  Continuing with tests anyway...${NC}"
    fi
    
    # Check frontend
    if curl -f -s "$FRONTEND_URL" >/dev/null; then
        echo -e "${GREEN}✅ Frontend health check passed${NC}"
    else
        echo -e "${RED}❌ Frontend health check failed${NC}"
        echo -e "${YELLOW}⚠️  Continuing with tests anyway...${NC}"
    fi
    
    echo ""
}

# Function to display test menu
show_menu() {
    echo -e "${BLUE}🎯 Select Herald.lol Performance Test:${NC}"
    echo ""
    echo -e "1) ${GREEN}Gaming Analytics Load Test${NC} - Standard load testing"
    echo -e "2) ${YELLOW}Gaming Stress Test${NC} - Extreme load conditions" 
    echo -e "3) ${RED}Gaming Spike Test${NC} - Sudden load spikes"
    echo -e "4) ${BLUE}Run All Tests${NC} - Complete test suite"
    echo -e "5) ${GREEN}Quick Health Check${NC} - Basic functionality"
    echo -e "6) Exit"
    echo ""
    echo -n "Enter your choice (1-6): "
}

# Function for quick health check
quick_health_check() {
    echo -e "${BLUE}🏥 Running Quick Health Check...${NC}"
    
    # Simple k6 health check script
    cat > /tmp/health_check.js << 'EOF'
import http from 'k6/http';
import { check } from 'k6';

export let options = {
  vus: 10,
  duration: '30s',
};

export default function() {
  const baseUrl = __ENV.API_BASE_URL || 'http://localhost:8080';
  
  const response = http.get(`${baseUrl}/health`);
  
  check(response, {
    'Health check successful': (r) => r.status === 200,
    'Response time < 1000ms': (r) => r.timings.duration < 1000,
  });
}
EOF
    
    k6 run --env API_BASE_URL="$API_BASE_URL" /tmp/health_check.js
    rm -f /tmp/health_check.js
}

# Function to generate summary report
generate_summary_report() {
    local summary_file="$RESULTS_DIR/test_summary_${TIMESTAMP}.md"
    
    echo -e "${BLUE}📊 Generating test summary report...${NC}"
    
    cat > "$summary_file" << EOF
# Herald.lol Performance Testing Summary

**Test Date:** $(date)
**API Base URL:** $API_BASE_URL  
**Frontend URL:** $FRONTEND_URL

## Performance Requirements
- **Analytics Response Time:** <5 seconds
- **UI Load Time:** <2 seconds  
- **Uptime Target:** 99.9%
- **Concurrent Users:** 1M+ support

## Test Results

### Gaming Analytics Load Test
$(if [ -f "$RESULTS_DIR/gaming-load_${TIMESTAMP}.json" ]; then
    echo "✅ Completed - Check detailed results in JSON file"
else
    echo "❌ Not run or failed"
fi)

### Gaming Stress Test  
$(if [ -f "$RESULTS_DIR/gaming-stress_${TIMESTAMP}.json" ]; then
    echo "✅ Completed - Check detailed results in JSON file"
else
    echo "❌ Not run or failed" 
fi)

### Gaming Spike Test
$(if [ -f "$RESULTS_DIR/gaming-spike_${TIMESTAMP}.json" ]; then
    echo "✅ Completed - Check detailed results in JSON file"
else
    echo "❌ Not run or failed"
fi)

## Files Generated
- Test results: $RESULTS_DIR/*_${TIMESTAMP}.json
- HTML reports: $RESULTS_DIR/*_${TIMESTAMP}.html (if available)
- Summary report: $summary_file

## Next Steps
1. Review detailed JSON results for metrics analysis
2. Check HTML reports for visual performance insights  
3. Address any performance bottlenecks identified
4. Re-run tests after optimizations

---
*Generated by Herald.lol Performance Testing Suite*
EOF

    echo -e "${GREEN}✅ Summary report generated: $summary_file${NC}"
}

# Main execution
main() {
    check_system_health
    
    if [ $# -eq 0 ]; then
        # Interactive mode
        while true; do
            show_menu
            read -r choice
            
            case $choice in
                1)
                    run_k6_test "gaming-analytics-load.js" "gaming-load"
                    ;;
                2)
                    echo -e "${YELLOW}⚠️  This will apply extreme load to your system. Continue? (y/N)${NC}"
                    read -r confirm
                    if [[ $confirm =~ ^[Yy]$ ]]; then
                        run_k6_test "gaming-stress-test.js" "gaming-stress"
                    fi
                    ;;
                3)
                    run_k6_test "gaming-spike-test.js" "gaming-spike"
                    ;;
                4)
                    echo -e "${BLUE}🚀 Running complete Herald.lol test suite...${NC}"
                    run_k6_test "gaming-analytics-load.js" "gaming-load"
                    run_k6_test "gaming-spike-test.js" "gaming-spike"
                    
                    echo -e "${YELLOW}⚠️  Run stress test? This applies extreme load (y/N)${NC}"
                    read -r confirm
                    if [[ $confirm =~ ^[Yy]$ ]]; then
                        run_k6_test "gaming-stress-test.js" "gaming-stress"
                    fi
                    
                    generate_summary_report
                    ;;
                5)
                    quick_health_check
                    ;;
                6)
                    echo -e "${GREEN}🎮 Herald.lol Performance Testing Complete!${NC}"
                    break
                    ;;
                *)
                    echo -e "${RED}Invalid choice. Please select 1-6.${NC}"
                    ;;
            esac
        done
    else
        # Command line mode
        case "$1" in
            "load")
                run_k6_test "gaming-analytics-load.js" "gaming-load"
                ;;
            "stress") 
                run_k6_test "gaming-stress-test.js" "gaming-stress"
                ;;
            "spike")
                run_k6_test "gaming-spike-test.js" "gaming-spike"
                ;;
            "all")
                run_k6_test "gaming-analytics-load.js" "gaming-load"
                run_k6_test "gaming-spike-test.js" "gaming-spike" 
                run_k6_test "gaming-stress-test.js" "gaming-stress"
                generate_summary_report
                ;;
            "health")
                quick_health_check
                ;;
            *)
                echo -e "${RED}Usage: $0 [load|stress|spike|all|health]${NC}"
                exit 1
                ;;
        esac
    fi
}

# Run main function
main "$@"