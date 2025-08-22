#!/bin/bash

# Test script for Herald.lol API endpoints
# Tests both authentication and group management functionality

echo "🚀 Starting Herald.lol API endpoint tests..."

# Configuration
BASE_URL="http://localhost:8000"
TEST_USER_EMAIL="test@herald-lol.com"
TEST_RIOT_ID="TestUser"
TEST_RIOT_TAG="EUW"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    
    case $status in
        "SUCCESS")
            echo -e "${GREEN}✅ $message${NC}"
            ;;
        "ERROR")
            echo -e "${RED}❌ $message${NC}"
            ;;
        "INFO")
            echo -e "${BLUE}ℹ️  $message${NC}"
            ;;
        "WARNING")
            echo -e "${YELLOW}⚠️  $message${NC}"
            ;;
    esac
}

# Function to test endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4
    local description=$5
    local auth_header=$6
    
    echo
    print_status "INFO" "Testing: $description"
    echo "Endpoint: $method $BASE_URL$endpoint"
    
    if [ -n "$auth_header" ]; then
        if [ -n "$data" ]; then
            response=$(curl -s -w "\n%{http_code}" -X $method \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $auth_header" \
                -d "$data" \
                "$BASE_URL$endpoint")
        else
            response=$(curl -s -w "\n%{http_code}" -X $method \
                -H "Authorization: Bearer $auth_header" \
                "$BASE_URL$endpoint")
        fi
    else
        if [ -n "$data" ]; then
            response=$(curl -s -w "\n%{http_code}" -X $method \
                -H "Content-Type: application/json" \
                -d "$data" \
                "$BASE_URL$endpoint")
        else
            response=$(curl -s -w "\n%{http_code}" -X $method \
                "$BASE_URL$endpoint")
        fi
    fi
    
    # Extract status code (last line) and body (everything else)
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$status_code" = "$expected_status" ]; then
        print_status "SUCCESS" "Status: $status_code (Expected: $expected_status)"
        if [ -n "$body" ] && [ "$body" != "null" ]; then
            echo "Response: $body" | jq . 2>/dev/null || echo "Response: $body"
        fi
    else
        print_status "ERROR" "Status: $status_code (Expected: $expected_status)"
        echo "Response: $body"
        return 1
    fi
    
    return 0
}

# Check if server is running
print_status "INFO" "Checking if server is running at $BASE_URL"
if ! curl -s "$BASE_URL/health" > /dev/null 2>&1; then
    print_status "ERROR" "Server is not running at $BASE_URL"
    print_status "INFO" "Please start the server first with: cd /home/debian/herald && docker-compose up"
    exit 1
fi

print_status "SUCCESS" "Server is running"

# Test counter
total_tests=0
passed_tests=0

# Helper function to increment test counter
run_test() {
    total_tests=$((total_tests + 1))
    if test_endpoint "$@"; then
        passed_tests=$((passed_tests + 1))
    fi
}

echo
echo "=================================="
echo "🧪 AUTHENTICATION ENDPOINTS TESTS"
echo "=================================="

# Test health check
run_test "GET" "/health" "" "200" "Health Check"

# Test OAuth endpoints (these might need actual OAuth setup)
run_test "GET" "/api/auth/google" "" "200" "Google OAuth URL Generation"

# Test invalid auth
run_test "GET" "/api/user/profile" "" "401" "Unauthorized Profile Access"

echo
echo "============================="
echo "👥 GROUP ENDPOINTS TESTS"
echo "============================="

# Mock auth token for testing (in real scenario, this would come from OAuth)
AUTH_TOKEN="mock-jwt-token-for-testing"

# Test group creation (will fail without proper auth, but tests endpoint)
GROUP_DATA='{
    "name": "Test Group",
    "description": "A test group for League players",
    "privacy": "public"
}'

run_test "POST" "/api/groups" "$GROUP_DATA" "401" "Create Group (Unauthorized)" ""

# Test fetching user groups (will fail without proper auth)
run_test "GET" "/api/groups/my" "" "401" "Get User Groups (Unauthorized)" ""

# Test joining group by invite code
JOIN_DATA='{"invite_code": "TEST123"}'
run_test "POST" "/api/groups/join" "$JOIN_DATA" "401" "Join Group by Invite (Unauthorized)" ""

# Test getting group details
run_test "GET" "/api/groups/1" "" "401" "Get Group Details (Unauthorized)" ""

# Test group statistics
run_test "GET" "/api/groups/1/stats" "" "401" "Get Group Statistics (Unauthorized)" ""

echo
echo "================================"
echo "📊 COMPARISON ENDPOINTS TESTS"
echo "================================"

# Test comparison creation
COMPARISON_DATA='{
    "name": "Test Comparison",
    "description": "Comparing player performance",
    "compare_type": "champions",
    "parameters": {
        "member_ids": [1, 2, 3],
        "time_range": "30d",
        "metrics": ["winrate", "kda", "cs"],
        "min_games": 5
    }
}'

run_test "POST" "/api/groups/1/comparisons" "$COMPARISON_DATA" "401" "Create Comparison (Unauthorized)" ""

# Test getting comparisons
run_test "GET" "/api/groups/1/comparisons" "" "401" "Get Group Comparisons (Unauthorized)" ""

# Test getting specific comparison
run_test "GET" "/api/groups/1/comparisons/1" "" "401" "Get Comparison Details (Unauthorized)" ""

echo
echo "=========================="
echo "🔧 UTILITY ENDPOINTS TESTS"
echo "=========================="

# Test public endpoints that might not require auth
run_test "GET" "/api/public/champions" "" "404" "Get Champions List (Not Implemented)" ""
run_test "GET" "/api/public/regions" "" "404" "Get Regions List (Not Implemented)" ""

echo
echo "========================="
echo "📈 PERFORMANCE TESTS"
echo "========================="

print_status "INFO" "Testing endpoint response times..."

# Test multiple requests to check performance
for i in {1..5}; do
    start_time=$(date +%s%N)
    curl -s "$BASE_URL/health" > /dev/null
    end_time=$(date +%s%N)
    duration=$(( (end_time - start_time) / 1000000 )) # Convert to milliseconds
    
    if [ $duration -lt 100 ]; then
        print_status "SUCCESS" "Health check #$i: ${duration}ms (Fast)"
    elif [ $duration -lt 500 ]; then
        print_status "WARNING" "Health check #$i: ${duration}ms (Acceptable)"
    else
        print_status "ERROR" "Health check #$i: ${duration}ms (Slow)"
    fi
done

echo
echo "========================"
echo "📝 TEST SUMMARY"
echo "========================"

echo "Total tests run: $total_tests"
echo "Tests passed: $passed_tests"
echo "Tests failed: $((total_tests - passed_tests))"

if [ $passed_tests -eq $total_tests ]; then
    print_status "SUCCESS" "All tests passed! 🎉"
    exit 0
elif [ $passed_tests -gt $((total_tests / 2)) ]; then
    print_status "WARNING" "Most tests passed, but some failures detected"
    exit 1
else
    print_status "ERROR" "Many tests failed - check server configuration"
    exit 1
fi

echo
echo "========================"
echo "💡 NEXT STEPS"
echo "========================"

echo "1. Set up proper OAuth authentication"
echo "2. Configure Riot API integration"
echo "3. Implement missing public endpoints"
echo "4. Add comprehensive error handling"
echo "5. Set up database with test data"

print_status "INFO" "Test script completed"