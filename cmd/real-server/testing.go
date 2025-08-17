package main

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

// TestResult represents the result of a system test
type TestResult struct {
	TestName    string        `json:"test_name"`
	Status      string        `json:"status"` // "passed", "failed", "warning"
	Duration    time.Duration `json:"duration_ms"`
	Message     string        `json:"message"`
	Details     interface{}   `json:"details,omitempty"`
	Timestamp   time.Time     `json:"timestamp"`
	Category    string        `json:"category"` // "database", "cache", "websocket", "api", "performance"
}

// TestSuite manages and executes comprehensive system tests
type TestSuite struct {
	results []TestResult
	mutex   sync.RWMutex
}

// Global test suite
var testSuite *TestSuite

// InitializeTesting initializes the testing system
func InitializeTesting() {
	testSuite = &TestSuite{
		results: make([]TestResult, 0),
	}
	log.Println("🧪 Testing system initialized")
}

// RunAllTests executes comprehensive system tests
func (ts *TestSuite) RunAllTests() []TestResult {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	
	ts.results = make([]TestResult, 0)
	
	log.Println("🧪 Starting comprehensive system tests...")
	
	// Database tests
	ts.runDatabaseTests()
	
	// Cache tests
	ts.runCacheTests()
	
	// WebSocket tests
	ts.runWebSocketTests()
	
	// Monitoring tests
	ts.runMonitoringTests()
	
	// Performance tests
	ts.runPerformanceTests()
	
	// Integration tests
	ts.runIntegrationTests()
	
	log.Printf("🧪 Testing completed: %d tests executed", len(ts.results))
	return ts.results
}

// runDatabaseTests tests database functionality
func (ts *TestSuite) runDatabaseTests() {
	log.Println("🗄️ Running database tests...")
	
	// Test database connection
	start := time.Now()
	err := database.Ping()
	duration := time.Since(start)
	
	if err != nil {
		ts.addResult(TestResult{
			TestName:  "Database Connection",
			Status:    "failed",
			Duration:  duration,
			Message:   fmt.Sprintf("Database ping failed: %v", err),
			Category:  "database",
			Timestamp: time.Now(),
		})
		return
	}
	
	ts.addResult(TestResult{
		TestName:  "Database Connection",
		Status:    "passed",
		Duration:  duration,
		Message:   "Database connection successful",
		Category:  "database",
		Timestamp: time.Now(),
	})
	
	// Test table existence
	ts.testTableExistence()
	
	// Test CRUD operations
	ts.testCRUDOperations()
	
	// Test query performance
	ts.testQueryPerformance()
}

// testTableExistence verifies all required tables exist
func (ts *TestSuite) testTableExistence() {
	start := time.Now()
	
	tables := []string{"users", "matches"}
	missingTables := make([]string, 0)
	
	for _, table := range tables {
		var name string
		query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?`
		err := database.QueryRow(query, table).Scan(&name)
		if err != nil {
			missingTables = append(missingTables, table)
		}
	}
	
	duration := time.Since(start)
	
	if len(missingTables) > 0 {
		ts.addResult(TestResult{
			TestName:  "Table Existence",
			Status:    "failed",
			Duration:  duration,
			Message:   fmt.Sprintf("Missing tables: %v", missingTables),
			Details:   map[string]interface{}{"missing_tables": missingTables},
			Category:  "database",
			Timestamp: time.Now(),
		})
		return
	}
	
	ts.addResult(TestResult{
		TestName:  "Table Existence",
		Status:    "passed",
		Duration:  duration,
		Message:   "All required tables exist",
		Details:   map[string]interface{}{"tables": tables},
		Category:  "database",
		Timestamp: time.Now(),
	})
}

// testCRUDOperations tests basic database operations
func (ts *TestSuite) testCRUDOperations() {
	start := time.Now()
	
	// Test user insertion
	testUser := User{
		RiotID:        "TestUser",
		RiotTag:       "TEST",
		RiotPUUID:     "test-puuid-" + fmt.Sprintf("%d", time.Now().UnixNano()),
		Region:        "euw1",
		SummonerLevel: 100,
	}
	
	userID, err := database.UpsertUser(testUser)
	if err != nil {
		ts.addResult(TestResult{
			TestName:  "User CRUD Operations",
			Status:    "failed",
			Duration:  time.Since(start),
			Message:   fmt.Sprintf("Failed to insert user: %v", err),
			Category:  "database",
			Timestamp: time.Now(),
		})
		return
	}
	
	// Test match insertion
	testMatch := Match{
		ID:           "TEST_MATCH_" + fmt.Sprintf("%d", time.Now().UnixNano()),
		UserID:       userID,
		GameCreation: time.Now().Unix() * 1000,
		GameDuration: 1800,
		GameMode:     "CLASSIC",
		QueueID:      420,
		ChampionName: "Aatrox",
		ChampionID:   266,
		Kills:        10,
		Deaths:       3,
		Assists:      15,
		Win:          true,
	}
	
	err = database.SaveMatch(testMatch)
	if err != nil {
		ts.addResult(TestResult{
			TestName:  "Match CRUD Operations",
			Status:    "failed",
			Duration:  time.Since(start),
			Message:   fmt.Sprintf("Failed to insert match: %v", err),
			Category:  "database",
			Timestamp: time.Now(),
		})
		return
	}
	
	// Test data retrieval
	matches, err := database.GetMatchesByUser(userID, 10, 0)
	if err != nil || len(matches) == 0 {
		ts.addResult(TestResult{
			TestName:  "Data Retrieval",
			Status:    "failed",
			Duration:  time.Since(start),
			Message:   "Failed to retrieve inserted data",
			Category:  "database",
			Timestamp: time.Now(),
		})
		return
	}
	
	duration := time.Since(start)
	ts.addResult(TestResult{
		TestName:  "CRUD Operations",
		Status:    "passed",
		Duration:  duration,
		Message:   "All CRUD operations successful",
		Details: map[string]interface{}{
			"user_id":      userID,
			"match_count":  len(matches),
			"operations":   []string{"user_insert", "match_insert", "data_retrieval"},
		},
		Category:  "database",
		Timestamp: time.Now(),
	})
}

// testQueryPerformance tests database query performance
func (ts *TestSuite) testQueryPerformance() {
	start := time.Now()
	
	// Test multiple query types and measure performance
	queryTests := map[string]func() error{
		"user_count": func() error {
			var count int
			return database.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
		},
		"match_count": func() error {
			var count int
			return database.QueryRow("SELECT COUNT(*) FROM matches").Scan(&count)
		},
		"user_stats": func() error {
			if database != nil {
				_, err := database.GetUserStats(1)
				return err
			}
			return fmt.Errorf("database not available")
		},
	}
	
	results := make(map[string]time.Duration)
	failures := make([]string, 0)
	
	for testName, testFunc := range queryTests {
		queryStart := time.Now()
		err := testFunc()
		queryDuration := time.Since(queryStart)
		
		results[testName] = queryDuration
		
		if err != nil {
			failures = append(failures, testName)
		}
		
		// Warning if query takes too long
		if queryDuration > 100*time.Millisecond {
			ts.addResult(TestResult{
				TestName:  fmt.Sprintf("Query Performance - %s", testName),
				Status:    "warning",
				Duration:  queryDuration,
				Message:   "Query execution time exceeds recommended threshold",
				Category:  "database",
				Timestamp: time.Now(),
			})
		}
	}
	
	duration := time.Since(start)
	status := "passed"
	message := "All query performance tests passed"
	
	if len(failures) > 0 {
		status = "failed"
		message = fmt.Sprintf("Query failures: %v", failures)
	}
	
	ts.addResult(TestResult{
		TestName:  "Query Performance",
		Status:    status,
		Duration:  duration,
		Message:   message,
		Details: map[string]interface{}{
			"query_times": results,
			"failures":    failures,
		},
		Category:  "database",
		Timestamp: time.Now(),
	})
}

// runCacheTests tests cache functionality
func (ts *TestSuite) runCacheTests() {
	log.Println("🧠 Running cache tests...")
	
	if smartCache == nil {
		ts.addResult(TestResult{
			TestName:  "Cache Availability",
			Status:    "failed",
			Duration:  0,
			Message:   "Smart cache not initialized",
			Category:  "cache",
			Timestamp: time.Now(),
		})
		return
	}
	
	start := time.Now()
	
	// Test cache operations
	testKey := "test_cache_key"
	testValue := map[string]interface{}{
		"test_data": "cache_test_value",
		"timestamp": time.Now(),
	}
	
	// Test cache set
	smartCache.cache.Set(testKey, testValue, 60)
	
	// Test cache get
	_, hit := smartCache.cache.Get(testKey)
	if !hit {
		ts.addResult(TestResult{
			TestName:  "Cache Hit/Miss",
			Status:    "failed",
			Duration:  time.Since(start),
			Message:   "Cache miss on recently set value",
			Category:  "cache",
			Timestamp: time.Now(),
		})
		return
	}
	
	// Test cache statistics
	stats := smartCache.GetCacheStats()
	
	duration := time.Since(start)
	ts.addResult(TestResult{
		TestName:  "Cache Operations",
		Status:    "passed",
		Duration:  duration,
		Message:   "Cache operations successful",
		Details: map[string]interface{}{
			"hit_ratio":    stats.HitRatio,
			"cache_size":   stats.Size,
			"total_hits":   stats.Hits,
			"total_misses": stats.Misses,
		},
		Category:  "cache",
		Timestamp: time.Now(),
	})
	
	// Performance warning if hit ratio is low
	if stats.HitRatio < 0.7 && stats.Hits+stats.Misses > 50 {
		ts.addResult(TestResult{
			TestName:  "Cache Hit Ratio",
			Status:    "warning",
			Duration:  0,
			Message:   fmt.Sprintf("Low cache hit ratio: %.2f%%", stats.HitRatio*100),
			Category:  "cache",
			Timestamp: time.Now(),
		})
	}
}

// runWebSocketTests tests WebSocket functionality
func (ts *TestSuite) runWebSocketTests() {
	log.Println("🔌 Running WebSocket tests...")
	
	if wsHub == nil {
		ts.addResult(TestResult{
			TestName:  "WebSocket Availability",
			Status:    "failed",
			Duration:  0,
			Message:   "WebSocket hub not initialized",
			Category:  "websocket",
			Timestamp: time.Now(),
		})
		return
	}
	
	start := time.Now()
	
	// Test WebSocket statistics
	stats := wsHub.GetStats()
	
	duration := time.Since(start)
	ts.addResult(TestResult{
		TestName:  "WebSocket Statistics",
		Status:    "passed",
		Duration:  duration,
		Message:   "WebSocket system operational",
		Details: map[string]interface{}{
			"active_connections": stats.ActiveConnections,
			"total_connections":  stats.TotalConnections,
			"messages_sent":      stats.MessagesSent,
			"error_count":        stats.ErrorCount,
		},
		Category:  "websocket",
		Timestamp: time.Now(),
	})
	
	// Warning if error rate is high
	if stats.ErrorCount > 0 && stats.MessagesSent > 0 {
		errorRate := float64(stats.ErrorCount) / float64(stats.MessagesSent) * 100
		if errorRate > 5 {
			ts.addResult(TestResult{
				TestName:  "WebSocket Error Rate",
				Status:    "warning",
				Duration:  0,
				Message:   fmt.Sprintf("High WebSocket error rate: %.2f%%", errorRate),
				Category:  "websocket",
				Timestamp: time.Now(),
			})
		}
	}
}

// runMonitoringTests tests system monitoring
func (ts *TestSuite) runMonitoringTests() {
	log.Println("📊 Running monitoring tests...")
	
	if systemMonitor == nil {
		ts.addResult(TestResult{
			TestName:  "Monitoring Availability",
			Status:    "failed",
			Duration:  0,
			Message:   "System monitor not initialized",
			Category:  "monitoring",
			Timestamp: time.Now(),
		})
		return
	}
	
	start := time.Now()
	
	// Test health status
	health := systemMonitor.GetHealthStatus()
	
	// Test metrics collection
	metrics := systemMonitor.GetMetrics()
	
	duration := time.Since(start)
	
	status := "passed"
	message := "Monitoring system operational"
	
	if health.Status == "unhealthy" {
		status = "failed"
		message = "System health is unhealthy"
	} else if health.Status == "degraded" {
		status = "warning"
		message = "System health is degraded"
	}
	
	ts.addResult(TestResult{
		TestName:  "System Monitoring",
		Status:    status,
		Duration:  duration,
		Message:   message,
		Details: map[string]interface{}{
			"health_status":  health.Status,
			"health_score":   health.HealthScore,
			"memory_mb":      metrics.MemoryUsage.AllocatedMB,
			"goroutines":     metrics.NumGoroutines,
			"uptime_hours":   metrics.Uptime,
		},
		Category:  "monitoring",
		Timestamp: time.Now(),
	})
}

// runPerformanceTests tests system performance
func (ts *TestSuite) runPerformanceTests() {
	log.Println("⚡ Running performance tests...")
	
	start := time.Now()
	
	// Memory usage test
	metrics := systemMonitor.GetMetrics()
	memoryMB := metrics.MemoryUsage.AllocatedMB
	
	memoryStatus := "passed"
	memoryMessage := "Memory usage within normal range"
	
	if memoryMB > 100 {
		memoryStatus = "warning"
		memoryMessage = "High memory usage detected"
	}
	if memoryMB > 500 {
		memoryStatus = "failed"
		memoryMessage = "Critical memory usage"
	}
	
	ts.addResult(TestResult{
		TestName:  "Memory Usage",
		Status:    memoryStatus,
		Duration:  time.Since(start),
		Message:   fmt.Sprintf("%s: %.1f MB", memoryMessage, memoryMB),
		Details:   map[string]interface{}{"memory_mb": memoryMB},
		Category:  "performance",
		Timestamp: time.Now(),
	})
	
	// Goroutine count test
	goroutines := metrics.NumGoroutines
	goroutineStatus := "passed"
	goroutineMessage := "Goroutine count normal"
	
	if goroutines > 50 {
		goroutineStatus = "warning"
		goroutineMessage = "High goroutine count"
	}
	if goroutines > 200 {
		goroutineStatus = "failed"
		goroutineMessage = "Critical goroutine count"
	}
	
	ts.addResult(TestResult{
		TestName:  "Goroutine Count",
		Status:    goroutineStatus,
		Duration:  0,
		Message:   fmt.Sprintf("%s: %d", goroutineMessage, goroutines),
		Details:   map[string]interface{}{"goroutines": goroutines},
		Category:  "performance",
		Timestamp: time.Now(),
	})
}

// runIntegrationTests tests end-to-end functionality
func (ts *TestSuite) runIntegrationTests() {
	log.Println("🔄 Running integration tests...")
	
	start := time.Now()
	
	// Test full data flow: cache -> database -> response
	testCompleted := true
	integrationDetails := make(map[string]interface{})
	
	// Test 1: Database to Cache integration
	if smartCache != nil && database != nil {
		// Simulate getting user stats with caching
		fetcher := func() (UserStats, error) {
			return database.GetUserStats(1)
		}
		
		// First call (should miss cache)
		_, err := smartCache.GetUserStats(1, fetcher)
		if err != nil {
			integrationDetails["cache_db_error"] = err.Error()
			testCompleted = false
		}
		
		// Second call (should hit cache)
		_, err = smartCache.GetUserStats(1, fetcher)
		if err != nil {
			integrationDetails["cache_hit_error"] = err.Error()
			testCompleted = false
		}
	}
	
	// Test 2: Monitoring integration
	if systemMonitor != nil {
		health := systemMonitor.GetHealthStatus()
		integrationDetails["monitoring_health"] = health.Status
	}
	
	duration := time.Since(start)
	status := "passed"
	message := "All integration tests passed"
	
	if !testCompleted {
		status = "failed"
		message = "Integration test failures detected"
	}
	
	ts.addResult(TestResult{
		TestName:  "Integration Tests",
		Status:    status,
		Duration:  duration,
		Message:   message,
		Details:   integrationDetails,
		Category:  "integration",
		Timestamp: time.Now(),
	})
}

// addResult adds a test result to the suite
func (ts *TestSuite) addResult(result TestResult) {
	ts.results = append(ts.results, result)
	
	// Log result
	emoji := "✅"
	if result.Status == "failed" {
		emoji = "❌"
	} else if result.Status == "warning" {
		emoji = "⚠️"
	}
	
	log.Printf("🧪 %s %s [%s] %s (%.2fms)", 
		emoji, result.TestName, result.Category, result.Message, 
		float64(result.Duration.Nanoseconds())/1e6)
}

// GetTestResults returns all test results
func (ts *TestSuite) GetTestResults() []TestResult {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	
	results := make([]TestResult, len(ts.results))
	copy(results, ts.results)
	return results
}

// GetTestSummary returns a summary of test results
func (ts *TestSuite) GetTestSummary() map[string]interface{} {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	
	summary := map[string]interface{}{
		"total_tests": len(ts.results),
		"passed":      0,
		"failed":      0,
		"warnings":    0,
		"categories":  make(map[string]int),
		"avg_duration_ms": 0.0,
		"total_duration_ms": 0.0,
	}
	
	totalDuration := int64(0)
	
	for _, result := range ts.results {
		switch result.Status {
		case "passed":
			summary["passed"] = summary["passed"].(int) + 1
		case "failed":
			summary["failed"] = summary["failed"].(int) + 1
		case "warning":
			summary["warnings"] = summary["warnings"].(int) + 1
		}
		
		categories := summary["categories"].(map[string]int)
		categories[result.Category]++
		summary["categories"] = categories
		
		totalDuration += result.Duration.Nanoseconds()
	}
	
	if len(ts.results) > 0 {
		summary["avg_duration_ms"] = float64(totalDuration) / float64(len(ts.results)) / 1e6
	}
	summary["total_duration_ms"] = float64(totalDuration) / 1e6
	
	// Calculate success rate
	totalTests := float64(len(ts.results))
	if totalTests > 0 {
		successRate := float64(summary["passed"].(int)) / totalTests * 100
		summary["success_rate"] = math.Round(successRate*100) / 100
	}
	
	return summary
}

// RunStressTest performs stress testing on the system
func (ts *TestSuite) RunStressTest(duration time.Duration, concurrentRequests int) TestResult {
	log.Printf("💪 Starting stress test: %d concurrent requests for %v", 
		concurrentRequests, duration)
	
	start := time.Now()
	done := make(chan bool)
	results := make(chan bool, concurrentRequests*100)
	
	// Worker function
	worker := func() {
		for {
			select {
			case <-done:
				return
			default:
				// Simulate various operations
				if smartCache != nil {
					smartCache.cache.Get("stress_test_key")
				}
				if database != nil {
					var count int
					database.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
				}
				results <- true
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
	
	// Start workers
	for i := 0; i < concurrentRequests; i++ {
		go worker()
	}
	
	// Stop after duration
	go func() {
		time.Sleep(duration)
		close(done)
	}()
	
	// Count successful operations
	successCount := 0
	timeout := time.After(duration + 5*time.Second)
	
	for {
		select {
		case <-results:
			successCount++
		case <-timeout:
			goto finished
		}
		
		if time.Since(start) > duration {
			break
		}
	}
	
finished:
	testDuration := time.Since(start)
	opsPerSecond := float64(successCount) / testDuration.Seconds()
	
	status := "passed"
	message := "Stress test completed successfully"
	
	if opsPerSecond < 100 {
		status = "warning"
		message = "Low operations per second during stress test"
	}
	if opsPerSecond < 50 {
		status = "failed"
		message = "Critical performance during stress test"
	}
	
	result := TestResult{
		TestName:  "Stress Test",
		Status:    status,
		Duration:  testDuration,
		Message:   message,
		Details: map[string]interface{}{
			"operations_completed": successCount,
			"operations_per_second": math.Round(opsPerSecond*100) / 100,
			"concurrent_requests":   concurrentRequests,
			"test_duration_ms":      testDuration.Milliseconds(),
		},
		Category:  "performance",
		Timestamp: time.Now(),
	}
	
	log.Printf("💪 Stress test completed: %.1f ops/sec, %d total operations", 
		opsPerSecond, successCount)
	
	return result
}