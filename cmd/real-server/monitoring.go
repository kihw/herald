package main

import (
	"log"
	"runtime"
	"sync"
	"time"
)

// SystemMetrics represents comprehensive system monitoring data
type SystemMetrics struct {
	// Server information
	StartTime     time.Time `json:"start_time"`
	Uptime        float64   `json:"uptime_hours"`
	Version       string    `json:"version"`
	Environment   string    `json:"environment"`
	
	// Runtime metrics
	GoVersion     string  `json:"go_version"`
	NumGoroutines int     `json:"num_goroutines"`
	NumCPU        int     `json:"num_cpu"`
	
	// Memory metrics
	MemoryUsage   MemoryMetrics `json:"memory_usage"`
	
	// Database metrics
	DatabaseStats DatabaseMetrics `json:"database_stats"`
	
	// Cache metrics
	CacheStats    CacheStats `json:"cache_stats"`
	
	// WebSocket metrics
	WebSocketStats WebSocketStats `json:"websocket_stats"`
	
	// API metrics
	APIMetrics    APIMetrics `json:"api_metrics"`
	
	// Performance metrics
	Performance   PerformanceMetrics `json:"performance"`
	
	// Health status
	HealthStatus  HealthStatus `json:"health_status"`
}

// MemoryMetrics tracks memory usage
type MemoryMetrics struct {
	AllocatedMB     float64 `json:"allocated_mb"`
	TotalAllocMB    float64 `json:"total_alloc_mb"`
	SystemMB        float64 `json:"system_mb"`
	GCCycles        uint32  `json:"gc_cycles"`
	HeapObjectCount uint64  `json:"heap_objects"`
	StackMB         float64 `json:"stack_mb"`
}

// DatabaseMetrics tracks database performance
type DatabaseMetrics struct {
	TotalUsers      int                 `json:"total_users"`
	TotalMatches    int                 `json:"total_matches"`
	DatabaseSizeMB  float64            `json:"database_size_mb"`
	ActiveConnections int              `json:"active_connections"`
	QueryMetrics    map[string]QueryStat `json:"query_metrics"`
	LastBackup      *time.Time         `json:"last_backup"`
}

// QueryStat tracks individual query performance
type QueryStat struct {
	Count           int64         `json:"count"`
	TotalDuration   time.Duration `json:"total_duration_ms"`
	AverageDuration time.Duration `json:"average_duration_ms"`
	MaxDuration     time.Duration `json:"max_duration_ms"`
	ErrorCount      int64         `json:"error_count"`
}

// APIMetrics tracks API endpoint performance
type APIMetrics struct {
	TotalRequests     int64                    `json:"total_requests"`
	RequestsPerMinute float64                  `json:"requests_per_minute"`
	AverageLatencyMS  float64                  `json:"average_latency_ms"`
	ErrorRate         float64                  `json:"error_rate"`
	EndpointStats     map[string]EndpointStat  `json:"endpoint_stats"`
	StatusCodeCounts  map[int]int64            `json:"status_code_counts"`
}

// EndpointStat tracks individual endpoint performance
type EndpointStat struct {
	Count            int64         `json:"count"`
	AverageLatencyMS float64       `json:"average_latency_ms"`
	MaxLatencyMS     float64       `json:"max_latency_ms"`
	ErrorCount       int64         `json:"error_count"`
	LastAccessed     time.Time     `json:"last_accessed"`
}

// PerformanceMetrics tracks overall system performance
type PerformanceMetrics struct {
	CPU           CPUMetrics    `json:"cpu"`
	Network       NetworkMetrics `json:"network"`
	RiotAPI       RiotAPIMetrics `json:"riot_api"`
	ResponseTimes ResponseTimeMetrics `json:"response_times"`
}

// CPUMetrics tracks CPU usage
type CPUMetrics struct {
	UsagePercent    float64 `json:"usage_percent"`
	LoadAverage1m   float64 `json:"load_average_1m"`
	LoadAverage5m   float64 `json:"load_average_5m"`
	LoadAverage15m  float64 `json:"load_average_15m"`
}

// NetworkMetrics tracks network usage
type NetworkMetrics struct {
	BytesReceivedTotal int64   `json:"bytes_received_total"`
	BytesSentTotal     int64   `json:"bytes_sent_total"`
	ConnectionsActive  int     `json:"connections_active"`
	ConnectionsTotal   int64   `json:"connections_total"`
	BandwidthUsageMbps float64 `json:"bandwidth_usage_mbps"`
}

// RiotAPIMetrics tracks Riot API usage
type RiotAPIMetrics struct {
	RequestsTotal     int64         `json:"requests_total"`
	RequestsPerMinute float64       `json:"requests_per_minute"`
	AverageLatencyMS  float64       `json:"average_latency_ms"`
	ErrorRate         float64       `json:"error_rate"`
	RateLimit         RateLimitInfo `json:"rate_limit"`
	LastError         *time.Time    `json:"last_error"`
}

// RateLimitInfo tracks Riot API rate limiting
type RateLimitInfo struct {
	PersonalLimit   int `json:"personal_limit"`
	PersonalUsed    int `json:"personal_used"`
	MethodLimit     int `json:"method_limit"`
	MethodUsed      int `json:"method_used"`
	RetryAfter      int `json:"retry_after_seconds"`
}

// ResponseTimeMetrics tracks response time percentiles
type ResponseTimeMetrics struct {
	P50MS  float64 `json:"p50_ms"`
	P90MS  float64 `json:"p90_ms"`
	P95MS  float64 `json:"p95_ms"`
	P99MS  float64 `json:"p99_ms"`
	P999MS float64 `json:"p999_ms"`
}

// HealthStatus represents overall system health
type HealthStatus struct {
	Status           string               `json:"status"` // "healthy", "degraded", "unhealthy"
	DatabaseHealth   string               `json:"database_health"`
	CacheHealth      string               `json:"cache_health"`
	WebSocketHealth  string               `json:"websocket_health"`
	RiotAPIHealth    string               `json:"riot_api_health"`
	Issues           []HealthIssue        `json:"issues"`
	LastHealthCheck  time.Time            `json:"last_health_check"`
	HealthScore      float64              `json:"health_score"` // 0-100
}

// HealthIssue represents a system health issue
type HealthIssue struct {
	Severity    string    `json:"severity"` // "critical", "warning", "info"
	Component   string    `json:"component"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
	Resolved    bool      `json:"resolved"`
}

// SystemMonitor manages all system metrics
type SystemMonitor struct {
	metrics     SystemMetrics
	startTime   time.Time
	mutex       sync.RWMutex
	
	// Metric collection intervals
	updateInterval time.Duration
	
	// Historical data (for trends)
	historySize    int
	metricHistory  []SystemMetrics
}

// Global system monitor
var systemMonitor *SystemMonitor

// InitializeMonitoring initializes the system monitoring
func InitializeMonitoring() {
	systemMonitor = &SystemMonitor{
		startTime:      time.Now(),
		updateInterval: 30 * time.Second,
		historySize:    100, // Keep last 100 metrics snapshots
		metricHistory:  make([]SystemMetrics, 0, 100),
		metrics: SystemMetrics{
			StartTime:   time.Now(),
			Version:     "1.0.0",
			Environment: "development",
			GoVersion:   runtime.Version(),
			NumCPU:      runtime.NumCPU(),
			APIMetrics: APIMetrics{
				EndpointStats:    make(map[string]EndpointStat),
				StatusCodeCounts: make(map[int]int64),
			},
			DatabaseStats: DatabaseMetrics{
				QueryMetrics: make(map[string]QueryStat),
			},
			HealthStatus: HealthStatus{
				Status:  "healthy",
				Issues:  make([]HealthIssue, 0),
				HealthScore: 100.0,
			},
		},
	}
	
	// Start background monitoring
	go systemMonitor.startMonitoring()
	
	log.Println("📊 System monitoring initialized")
}

// startMonitoring runs the monitoring loop
func (sm *SystemMonitor) startMonitoring() {
	ticker := time.NewTicker(sm.updateInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			sm.updateMetrics()
		}
	}
}

// updateMetrics collects and updates all system metrics
func (sm *SystemMonitor) updateMetrics() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	// Update uptime
	sm.metrics.Uptime = time.Since(sm.startTime).Hours()
	
	// Update runtime metrics
	sm.metrics.NumGoroutines = runtime.NumGoroutine()
	
	// Update memory metrics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	sm.metrics.MemoryUsage = MemoryMetrics{
		AllocatedMB:     float64(memStats.Alloc) / 1024 / 1024,
		TotalAllocMB:    float64(memStats.TotalAlloc) / 1024 / 1024,
		SystemMB:        float64(memStats.Sys) / 1024 / 1024,
		GCCycles:        memStats.NumGC,
		HeapObjectCount: memStats.HeapObjects,
		StackMB:         float64(memStats.StackSys) / 1024 / 1024,
	}
	
	// Update database metrics (if database is available)
	if database != nil {
		sm.updateDatabaseMetrics()
	}
	
	// Update cache metrics (if cache is available)
	if smartCache != nil {
		sm.metrics.CacheStats = smartCache.GetCacheStats()
	}
	
	// Update WebSocket metrics (if WebSocket is available)
	if wsHub != nil {
		sm.metrics.WebSocketStats = wsHub.GetStats()
	}
	
	// Update health status
	sm.updateHealthStatus()
	
	// Add to history
	sm.addToHistory()
	
	log.Printf("📊 Metrics updated - Goroutines: %d, Memory: %.1f MB, Health: %s", 
		sm.metrics.NumGoroutines, sm.metrics.MemoryUsage.AllocatedMB, sm.metrics.HealthStatus.Status)
}

// updateDatabaseMetrics updates database-related metrics
func (sm *SystemMonitor) updateDatabaseMetrics() {
	// Count total users
	var userCount int
	database.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	
	// Count total matches
	var matchCount int
	database.QueryRow("SELECT COUNT(*) FROM matches").Scan(&matchCount)
	
	sm.metrics.DatabaseStats.TotalUsers = userCount
	sm.metrics.DatabaseStats.TotalMatches = matchCount
	
	// Estimate database size (simplified)
	sm.metrics.DatabaseStats.DatabaseSizeMB = float64(sm.metrics.DatabaseStats.TotalMatches) * 0.001 // Rough estimate
	sm.metrics.DatabaseStats.ActiveConnections = 1 // SQLite is single connection
}

// updateHealthStatus evaluates and updates system health
func (sm *SystemMonitor) updateHealthStatus() {
	var issues []HealthIssue
	healthScore := 100.0
	
	// Check memory usage
	if sm.metrics.MemoryUsage.AllocatedMB > 500 {
		issues = append(issues, HealthIssue{
			Severity:  "warning",
			Component: "memory",
			Message:   "High memory usage detected",
			Timestamp: time.Now(),
		})
		healthScore -= 10
	}
	
	// Check goroutine count
	if sm.metrics.NumGoroutines > 100 {
		issues = append(issues, HealthIssue{
			Severity:  "warning",
			Component: "runtime",
			Message:   "High number of goroutines",
			Timestamp: time.Now(),
		})
		healthScore -= 5
	}
	
	// Check cache hit ratio
	if sm.metrics.CacheStats.HitRatio < 0.5 && sm.metrics.CacheStats.Hits+sm.metrics.CacheStats.Misses > 100 {
		issues = append(issues, HealthIssue{
			Severity:  "warning",
			Component: "cache",
			Message:   "Low cache hit ratio",
			Timestamp: time.Now(),
		})
		healthScore -= 15
	}
	
	// Determine overall status
	status := "healthy"
	if healthScore < 90 {
		status = "degraded"
	}
	if healthScore < 70 {
		status = "unhealthy"
	}
	
	sm.metrics.HealthStatus = HealthStatus{
		Status:           status,
		DatabaseHealth:   "healthy",
		CacheHealth:      "healthy",
		WebSocketHealth:  "healthy",
		RiotAPIHealth:    "healthy",
		Issues:           issues,
		LastHealthCheck:  time.Now(),
		HealthScore:      healthScore,
	}
}

// addToHistory adds current metrics to historical data
func (sm *SystemMonitor) addToHistory() {
	// Create a copy of current metrics
	snapshot := sm.metrics
	
	// Add to history
	if len(sm.metricHistory) >= sm.historySize {
		// Remove oldest entry
		sm.metricHistory = sm.metricHistory[1:]
	}
	sm.metricHistory = append(sm.metricHistory, snapshot)
}

// GetMetrics returns current system metrics
func (sm *SystemMonitor) GetMetrics() SystemMetrics {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	return sm.metrics
}

// GetMetricsHistory returns historical metrics data
func (sm *SystemMonitor) GetMetricsHistory() []SystemMetrics {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	// Return a copy of the history
	history := make([]SystemMetrics, len(sm.metricHistory))
	copy(history, sm.metricHistory)
	return history
}

// RecordAPIRequest records an API request for metrics
func (sm *SystemMonitor) RecordAPIRequest(endpoint string, statusCode int, latencyMS float64) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	// Update total requests
	sm.metrics.APIMetrics.TotalRequests++
	
	// Update status code counts
	sm.metrics.APIMetrics.StatusCodeCounts[statusCode]++
	
	// Update endpoint stats
	stat := sm.metrics.APIMetrics.EndpointStats[endpoint]
	stat.Count++
	stat.LastAccessed = time.Now()
	
	// Calculate average latency
	if stat.Count == 1 {
		stat.AverageLatencyMS = latencyMS
		stat.MaxLatencyMS = latencyMS
	} else {
		stat.AverageLatencyMS = (stat.AverageLatencyMS*float64(stat.Count-1) + latencyMS) / float64(stat.Count)
		if latencyMS > stat.MaxLatencyMS {
			stat.MaxLatencyMS = latencyMS
		}
	}
	
	if statusCode >= 400 {
		stat.ErrorCount++
	}
	
	sm.metrics.APIMetrics.EndpointStats[endpoint] = stat
	
	// Calculate overall error rate
	totalErrors := int64(0)
	for code, count := range sm.metrics.APIMetrics.StatusCodeCounts {
		if code >= 400 {
			totalErrors += count
		}
	}
	sm.metrics.APIMetrics.ErrorRate = float64(totalErrors) / float64(sm.metrics.APIMetrics.TotalRequests) * 100
}

// RecordDatabaseQuery records a database query for metrics
func (sm *SystemMonitor) RecordDatabaseQuery(queryType string, duration time.Duration, err error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	stat := sm.metrics.DatabaseStats.QueryMetrics[queryType]
	stat.Count++
	stat.TotalDuration += duration
	stat.AverageDuration = stat.TotalDuration / time.Duration(stat.Count)
	
	if duration > stat.MaxDuration {
		stat.MaxDuration = duration
	}
	
	if err != nil {
		stat.ErrorCount++
	}
	
	sm.metrics.DatabaseStats.QueryMetrics[queryType] = stat
}

// GetHealthStatus returns current health status
func (sm *SystemMonitor) GetHealthStatus() HealthStatus {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	return sm.metrics.HealthStatus
}

// IsHealthy returns true if the system is healthy
func (sm *SystemMonitor) IsHealthy() bool {
	status := sm.GetHealthStatus()
	return status.Status == "healthy"
}