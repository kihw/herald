#!/bin/bash

# Herald.lol Go Testing Helper Script

set -e

echo "🧪 Running Herald.lol Go Tests..."

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
TIMEOUT=30s
COVERAGE_THRESHOLD=80

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Change to backend directory
cd "$(dirname "$0")"

# Check if Go is available (in Docker context)
if ! command -v go &> /dev/null; then
    print_error "Go is not available. Run this script inside the Docker container."
    exit 1
fi

print_status "Starting Go test suite for Herald.lol..."

# Clean test cache
print_status "Cleaning test cache..."
go clean -testcache

# Run unit tests with coverage
print_status "Running unit tests with coverage..."
go test ./... -v -race -timeout=$TIMEOUT -coverprofile=coverage.out -covermode=atomic

# Check if tests passed
if [ $? -eq 0 ]; then
    print_success "All unit tests passed!"
else
    print_error "Some unit tests failed!"
    exit 1
fi

# Generate coverage report
print_status "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

# Check coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
COVERAGE_INT=${COVERAGE%.*}

echo "📊 Test Coverage: ${COVERAGE}%"

if [ "$COVERAGE_INT" -ge "$COVERAGE_THRESHOLD" ]; then
    print_success "Coverage threshold met (${COVERAGE}% >= ${COVERAGE_THRESHOLD}%)"
else
    print_warning "Coverage below threshold (${COVERAGE}% < ${COVERAGE_THRESHOLD}%)"
fi

# Run benchmarks
print_status "Running benchmark tests..."
go test ./... -bench=. -benchmem -run=^$ > benchmarks.txt

if [ -f benchmarks.txt ]; then
    print_success "Benchmark results saved to benchmarks.txt"
    echo "📊 Benchmark Summary:"
    grep -E "(Benchmark|PASS)" benchmarks.txt | head -10
fi

# Run race condition detection
print_status "Running race condition detection..."
go test ./... -race -short

# Check for common Go issues
print_status "Running additional checks..."

# Check for potential nil pointer dereferences
print_status "Checking for potential issues with go vet..."
go vet ./...

# Check code formatting
print_status "Checking code formatting..."
UNFORMATTED=$(gofmt -l .)
if [ -n "$UNFORMATTED" ]; then
    print_warning "The following files need formatting:"
    echo "$UNFORMATTED"
else
    print_success "All Go files are properly formatted"
fi

# Generate test summary
echo ""
echo "🎯 Herald.lol Go Test Summary:"
echo "================================"
echo "✅ Unit Tests: PASSED"
echo "📊 Coverage: ${COVERAGE}%"
echo "🏁 Race Detection: PASSED"
echo "📝 Code Formatting: $([ -n "$UNFORMATTED" ] && echo "NEEDS WORK" || echo "PASSED")"
echo "🔍 Vet Analysis: PASSED"

print_success "Go test suite completed successfully!"

# Instructions for viewing results
echo ""
echo "📋 View Results:"
echo "  - Coverage Report: open coverage.html in browser"
echo "  - Benchmark Results: cat benchmarks.txt"
echo "  - Run specific tests: go test -v ./path/to/package -run TestName"
echo "  - Run benchmarks only: go test -bench=. -run=^$ ./..."

exit 0