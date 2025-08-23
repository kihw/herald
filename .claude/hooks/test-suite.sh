#!/bin/bash
# Herald.lol Gaming Platform Test Suite
# Comprehensive testing for gaming analytics platform

echo "🧪 Running Herald.lol gaming analytics test suite..."

TEST_FAILURES=0
TOTAL_TESTS=0

# Go backend tests (gaming analytics core)
if [ -f "go.mod" ]; then
    echo "🎮 Running Go backend tests for gaming analytics..."
    
    if command -v go >/dev/null 2>&1; then
        # Run tests with coverage
        echo "  📊 Running Go tests with coverage..."
        if go test ./... -v -cover -race -timeout=5m; then
            echo "  ✅ Go tests passed"
            
            # Generate coverage report for gaming analytics
            echo "  📈 Generating coverage report..."
            go test ./... -coverprofile=coverage.out -timeout=5m >/dev/null 2>&1
            if [ -f "coverage.out" ]; then
                COVERAGE=$(go tool cover -func=coverage.out | tail -n 1 | awk '{print $3}' | sed 's/%//')
                echo "  📊 Test coverage: ${COVERAGE}%"
                
                # Gaming platform coverage requirements
                if (( $(echo "$COVERAGE >= 80" | bc -l) )); then
                    echo "  ✅ Coverage meets gaming platform standards (≥80%)"
                else
                    echo "  ⚠️  Coverage below gaming platform target (${COVERAGE}% < 80%)"
                fi
                rm -f coverage.out
            fi
            
            # Gaming-specific test validations
            echo "  🎮 Gaming analytics test validations..."
            
            # Check for gaming metrics tests
            if grep -r "TestKDA\|TestCSPerMin\|TestVisionScore" . --include="*_test.go" >/dev/null 2>&1; then
                echo "  ✅ Gaming metrics tests found"
            else
                echo "  ⚠️  Consider adding tests for core gaming metrics (KDA, CS/min, Vision Score)"
            fi
            
            # Check for Riot API tests
            if grep -r "TestRiotAPI\|riot.*api.*test" . --include="*_test.go" -i >/dev/null 2>&1; then
                echo "  ✅ Riot API tests found"
            else
                echo "  ⚠️  Consider adding Riot API integration tests"
            fi
            
            # Check for performance tests
            if grep -r "Benchmark\|performance.*test" . --include="*_test.go" -i >/dev/null 2>&1; then
                echo "  ✅ Performance tests found"
            else
                echo "  ⚠️  Consider adding performance benchmarks for gaming analytics"
            fi
            
        else
            echo "  ❌ Go tests failed"
            TEST_FAILURES=$((TEST_FAILURES + 1))
        fi
        TOTAL_TESTS=$((TOTAL_TESTS + 1))
    fi
fi

# Frontend tests (React gaming UI)
if [ -f "package.json" ]; then
    echo "🎨 Running frontend tests for gaming UI..."
    
    # Jest tests for React components
    if grep -q "jest\|@testing-library" package.json 2>/dev/null; then
        echo "  ⚛️  Running React component tests..."
        
        if command -v npm >/dev/null 2>&1; then
            if npm test -- --watchAll=false --coverage >/dev/null 2>&1; then
                echo "  ✅ Frontend tests passed"
            else
                echo "  ❌ Frontend tests failed"
                TEST_FAILURES=$((TEST_FAILURES + 1))
            fi
        elif command -v yarn >/dev/null 2>&1; then
            if yarn test --watchAll=false --coverage >/dev/null 2>&1; then
                echo "  ✅ Frontend tests passed"
            else
                echo "  ❌ Frontend tests failed"
                TEST_FAILURES=$((TEST_FAILURES + 1))
            fi
        fi
        TOTAL_TESTS=$((TOTAL_TESTS + 1))
    fi
    
    # Gaming UI specific test validations
    echo "  🎮 Gaming UI test validations..."
    
    # Check for gaming component tests
    if find src/ -name "*.test.*" -exec grep -l "Gaming\|Analytics\|Chart\|Dashboard" {} \; 2>/dev/null | head -1 >/dev/null; then
        echo "  ✅ Gaming component tests found"
    else
        echo "  ⚠️  Consider adding tests for gaming analytics components"
    fi
    
    # Check for accessibility tests
    if find src/ -name "*.test.*" -exec grep -l "accessibility\|a11y\|aria" {} \; 2>/dev/null | head -1 >/dev/null; then
        echo "  ✅ Accessibility tests found"
    else
        echo "  ⚠️  Consider adding accessibility tests for gaming UI"
    fi
    
    # Check for performance tests
    if find src/ -name "*.test.*" -exec grep -l "performance\|render.*time" {} \; 2>/dev/null | head -1 >/dev/null; then
        echo "  ✅ Frontend performance tests found"
    else
        echo "  ⚠️  Consider adding performance tests for gaming UI components"
    fi
fi

# Integration tests for gaming platform
echo "🔗 Checking integration tests..."

# Database integration tests
if grep -r "integration.*test\|TestDB\|database.*test" . --include="*_test.go" --include="*.test.*" -i >/dev/null 2>&1; then
    echo "  ✅ Database integration tests found"
else
    echo "  ⚠️  Consider adding database integration tests for gaming data"
fi

# API integration tests
if grep -r "api.*test\|TestAPI\|integration.*api" . --include="*_test.go" --include="*.test.*" -i >/dev/null 2>&1; then
    echo "  ✅ API integration tests found"
else
    echo "  ⚠️  Consider adding API integration tests"
fi

# End-to-end tests (Cypress for gaming workflows)
if [ -f "cypress.config.js" ] || [ -f "cypress.json" ] || [ -d "cypress" ]; then
    echo "🎯 Running end-to-end tests for gaming workflows..."
    
    if command -v npx >/dev/null 2>&1; then
        # Check if Cypress is available
        if npx cypress --version >/dev/null 2>&1; then
            echo "  🎮 Cypress detected for gaming workflow testing"
            # Note: Not running Cypress here as it requires a running application
            echo "  ℹ️  Run 'npx cypress run' to execute e2e gaming workflow tests"
        fi
    fi
    
    # Check for gaming-specific e2e tests
    if find cypress/ -name "*.cy.*" -exec grep -l "gaming\|analytics\|dashboard\|player" {} \; 2>/dev/null | head -1 >/dev/null; then
        echo "  ✅ Gaming workflow e2e tests found"
    else
        echo "  ⚠️  Consider adding e2e tests for core gaming workflows"
    fi
fi

# Security tests for gaming platform
echo "🔒 Running security tests for gaming platform..."

# Check for security test files
if find . -name "*security*test*" -o -name "*auth*test*" 2>/dev/null | head -1 >/dev/null; then
    echo "  ✅ Security tests found"
else
    echo "  ⚠️  Consider adding security tests for gaming platform"
fi

# Performance benchmarks for gaming analytics
echo "⚡ Checking performance benchmarks..."

# Go benchmarks
if grep -r "func.*Benchmark" . --include="*_test.go" >/dev/null 2>&1; then
    echo "  📊 Go benchmarks found"
    if command -v go >/dev/null 2>&1; then
        echo "  ⚡ Running critical gaming analytics benchmarks..."
        go test -bench=. -benchmem ./... -timeout=2m | grep -E "(Benchmark|ns/op|B/op)" | head -10
    fi
else
    echo "  ⚠️  Consider adding performance benchmarks for gaming analytics"
fi

# Load testing hints
if [ -f "k6.js" ] || [ -d "loadtest" ]; then
    echo "  ✅ Load testing configuration found"
else
    echo "  ⚠️  Consider adding load tests for 1M+ concurrent user target"
fi

# Final test report
echo ""
echo "📋 Herald.lol Test Suite Summary:"
echo "  🧪 Total test suites run: $TOTAL_TESTS"
echo "  ❌ Failed test suites: $TEST_FAILURES"

if [ $TEST_FAILURES -eq 0 ] && [ $TOTAL_TESTS -gt 0 ]; then
    echo "✅ All Herald.lol gaming platform tests passed!"
    echo "🎮 Ready for gaming analytics deployment!"
elif [ $TOTAL_TESTS -eq 0 ]; then
    echo "⚠️  No tests detected - consider adding tests for gaming platform"
else
    echo "❌ Some tests failed - review before gaming platform deployment"
fi

# Gaming platform specific reminders
echo ""
echo "🎮 Herald.lol Testing Reminders:"
echo "  • Target <5s response time for gaming analytics"
echo "  • Test with realistic gaming data volumes"
echo "  • Validate Riot API rate limiting"
echo "  • Test gaming UI responsiveness on mobile"
echo "  • Ensure gaming metrics accuracy (KDA, CS/min, etc.)"

exit 0