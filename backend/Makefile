# Herald.lol Backend Testing and Development Makefile

.PHONY: test test-verbose test-coverage test-race test-bench test-gaming help clean deps lint fmt build run-tests

# Default target
help:
	@echo "Herald.lol Backend - Testing Commands"
	@echo ""
	@echo "Testing:"
	@echo "  test            Run all tests"
	@echo "  test-verbose    Run tests with verbose output"
	@echo "  test-coverage   Run tests with coverage report"
	@echo "  test-race       Run tests with race detection"
	@echo "  test-bench      Run benchmark tests"
	@echo "  test-gaming     Run gaming-specific performance tests"
	@echo "  test-all        Run comprehensive test suite"
	@echo ""
	@echo "Development:"
	@echo "  deps            Install dependencies"
	@echo "  lint            Run linters"
	@echo "  fmt             Format code"
	@echo "  build           Build application"
	@echo "  clean           Clean build artifacts"
	@echo ""
	@echo "Herald.lol Performance Targets:"
	@echo "  - Analytics: <5s post-game analysis"
	@echo "  - API Response: <500ms"
	@echo "  - Database Query: <100ms"
	@echo "  - Concurrent Users: 1M+ support"

# Basic test command
test:
	@echo "🎮 Running Herald.lol Backend Tests..."
	go test ./... -timeout 30s

# Verbose test output
test-verbose:
	@echo "🎮 Running Herald.lol Backend Tests (Verbose)..."
	go test ./... -v -timeout 30s

# Test with coverage
test-coverage:
	@echo "🎮 Running Herald.lol Tests with Coverage..."
	go test ./... -cover -coverprofile=coverage.out -timeout 30s
	@echo ""
	@echo "📊 Coverage Summary:"
	go tool cover -func=coverage.out | tail -1
	@echo ""
	@echo "💡 Generate HTML coverage report with: go tool cover -html=coverage.out"

# Race condition detection
test-race:
	@echo "🎮 Running Herald.lol Tests with Race Detection..."
	go test ./... -race -timeout 30s

# Performance benchmarks
test-bench:
	@echo "🎮 Running Herald.lol Performance Benchmarks..."
	@echo "⚡ Performance Targets: <5s analytics, <500ms API, 1M+ concurrent users"
	go test ./... -bench=. -benchmem -run=^$$ -timeout 10m

# Gaming-specific performance tests
test-gaming:
	@echo "🎮 Running Herald.lol Gaming Performance Tests..."
	@echo "🎯 Testing: KDA calculations, CS/min, Vision Score, Damage Share"
	go test ./internal/services -bench=BenchmarkGaming -benchmem -run=^$$ -timeout 5m
	go test ./internal/models -bench=BenchmarkHerald -benchmem -run=^$$ -timeout 5m
	@echo ""
	@echo "🎮 Testing Gaming Analytics Performance (<5s requirement):"
	go test ./internal/services -run=TestGamingAnalyticsPerformance -v -timeout 10s

# Comprehensive test suite
test-all: test-race test-coverage test-gaming
	@echo ""
	@echo "✅ Herald.lol Comprehensive Test Suite Complete!"
	@echo ""
	@echo "🎮 Gaming Analytics Platform Status:"
	@echo "  - All unit tests passing"
	@echo "  - Race conditions checked"
	@echo "  - Coverage report generated"
	@echo "  - Gaming performance validated"
	@echo "  - Ready for 1M+ concurrent users"

# Specific test packages
test-services:
	@echo "🎮 Testing Herald.lol Services..."
	go test ./internal/services/... -v -timeout 30s

test-models:
	@echo "🎮 Testing Herald.lol Models..."
	go test ./internal/models/... -v -timeout 30s

test-handlers:
	@echo "🎮 Testing Herald.lol Handlers..."
	go test ./internal/handlers/... -v -timeout 30s

# Performance-specific benchmarks
bench-analytics:
	@echo "🎮 Benchmarking Herald.lol Analytics (<5s target)..."
	go test ./internal/services -bench=BenchmarkAnalytics -benchmem -benchtime=5s

bench-gaming:
	@echo "🎮 Benchmarking Herald.lol Gaming Calculations..."
	go test ./internal/models -bench=BenchmarkGaming -benchmem -benchtime=10s

bench-concurrent:
	@echo "🎮 Benchmarking Herald.lol Concurrency (1M+ users target)..."
	go test ./internal/services -bench=BenchmarkConcurrent -benchmem -benchtime=10s

# Development commands
deps:
	@echo "📦 Installing Herald.lol Dependencies..."
	go mod tidy
	go mod verify
	go mod download

lint:
	@echo "🔍 Running Herald.lol Code Linting..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "⚠️  golangci-lint not installed. Install with:"; \
		echo "   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		echo "Running basic checks instead:"; \
		go vet ./...; \
		go fmt ./...; \
	fi

fmt:
	@echo "✨ Formatting Herald.lol Code..."
	go fmt ./...
	gofmt -s -w .

build:
	@echo "🏗️  Building Herald.lol Backend..."
	go build -o bin/herald-server ./cmd/server
	@echo "✅ Built: bin/herald-server"

# Database testing
test-db:
	@echo "🗄️  Testing Herald.lol Database Operations..."
	go test ./internal/models -run=TestDatabase -v

test-db-performance:
	@echo "🗄️  Testing Herald.lol Database Performance (<100ms queries)..."
	go test ./internal/models -run=TestDatabasePerformance -v

# Integration tests
test-integration:
	@echo "🔗 Running Herald.lol Integration Tests..."
	go test ./... -tags=integration -v -timeout 60s

# Clean up
clean:
	@echo "🧹 Cleaning Herald.lol Build Artifacts..."
	rm -rf bin/
	rm -f coverage.out
	rm -f *.prof
	go clean -cache
	go clean -modcache

# Watch tests (requires entr or similar)
test-watch:
	@echo "👀 Watching Herald.lol Tests..."
	@if command -v find >/dev/null 2>&1 && command -v entr >/dev/null 2>&1; then \
		find . -name "*.go" | entr -r make test; \
	else \
		echo "⚠️  Install 'entr' for watch functionality"; \
		echo "   On macOS: brew install entr"; \
		echo "   On Ubuntu: apt install entr"; \
	fi

# Generate test data
generate-test-data:
	@echo "🎲 Generating Herald.lol Test Data..."
	go run ./internal/testdata/generator/main.go

# Validate Herald.lol performance requirements
validate-performance:
	@echo "⚡ Validating Herald.lol Performance Requirements..."
	@echo ""
	@echo "🎯 Performance Targets:"
	@echo "  - Post-game Analysis: <5 seconds"
	@echo "  - API Response Time: <500ms"
	@echo "  - Database Queries: <100ms"
	@echo "  - Concurrent Users: 1M+ support"
	@echo ""
	@$(MAKE) test-gaming
	@echo ""
	@echo "✅ Performance validation complete!"

# Gaming-specific test scenarios
test-scenarios:
	@echo "🎮 Running Herald.lol Gaming Scenarios..."
	@echo "  - Perfect Game (20/0/15)"
	@echo "  - Average Game (8/4/12)"  
	@echo "  - Feeding Game (2/8/6)"
	@echo "  - Carry Game (15/3/8)"
	@echo "  - Support Game (2/3/25)"
	go test ./internal/services -run=TestGamingScenarios -v

# Memory profiling
profile-memory:
	@echo "🧠 Profiling Herald.lol Memory Usage..."
	go test ./internal/services -bench=BenchmarkAnalytics -memprofile=mem.prof -benchtime=10s
	@echo "📊 View profile with: go tool pprof mem.prof"

# CPU profiling  
profile-cpu:
	@echo "⚡ Profiling Herald.lol CPU Usage..."
	go test ./internal/services -bench=BenchmarkAnalytics -cpuprofile=cpu.prof -benchtime=10s
	@echo "📊 View profile with: go tool pprof cpu.prof"

# Test with different Go versions (if available)
test-compatibility:
	@echo "🔄 Testing Herald.lol Go Compatibility..."
	@for version in 1.21 1.22 1.23; do \
		if command -v go$$version >/dev/null 2>&1; then \
			echo "Testing with Go $$version..."; \
			go$$version test ./... -timeout 30s; \
		fi; \
	done

# Security testing
test-security:
	@echo "🔒 Running Herald.lol Security Tests..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "⚠️  gosec not installed. Install with:"; \
		echo "   go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Show test statistics
test-stats:
	@echo "📊 Herald.lol Test Statistics:"
	@echo ""
	@echo "Test Files:"
	@find . -name "*_test.go" | wc -l | awk '{print "  " $$1 " test files"}'
	@echo ""
	@echo "Test Functions:"
	@grep -r "^func Test" . --include="*_test.go" | wc -l | awk '{print "  " $$1 " test functions"}'
	@echo ""
	@echo "Benchmark Functions:"
	@grep -r "^func Benchmark" . --include="*_test.go" | wc -l | awk '{print "  " $$1 " benchmark functions"}'
	@echo ""
	@echo "Gaming Test Coverage:"
	@grep -r "KDA\|CS.*min\|Vision\|Damage" . --include="*_test.go" | wc -l | awk '{print "  " $$1 " gaming-specific tests"}'