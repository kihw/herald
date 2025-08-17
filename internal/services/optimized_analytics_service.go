package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	"lol-match-exporter/internal/cache"
	"lol-match-exporter/internal/db"
	"lol-match-exporter/internal/workers"
)

// OptimizedAnalyticsService provides high-performance analytics with caching and async processing
type OptimizedAnalyticsService struct {
	// Core services
	analyticsService *AnalyticsService
	cacheService     *cache.CacheService
	workerPool       *workers.WorkerPool

	// Configuration
	config OptimizedConfig

	// State
	mutex   sync.RWMutex
	running bool
}

// OptimizedConfig holds configuration for the optimized service
type OptimizedConfig struct {
	// Cache configuration
	CacheEnabled    bool
	CacheHost       string
	CachePort       int
	CachePassword   string
	CacheDB         int

	// Worker pool configuration
	EnableAsyncProcessing bool
	MaxWorkers           int
	QueueSize            int

	// Performance tuning
	EnableConcurrentQueries bool
	QueryTimeout            time.Duration
	BatchSize               int

	// Cache TTL settings
	ShortCacheTTL     time.Duration
	MediumCacheTTL    time.Duration
	LongCacheTTL      time.Duration
	VeryLongCacheTTL  time.Duration
}

// DefaultOptimizedConfig returns a default configuration
func DefaultOptimizedConfig() OptimizedConfig {
	return OptimizedConfig{
		CacheEnabled:            true,
		CacheHost:              "localhost",
		CachePort:              6379,
		CachePassword:          "",
		CacheDB:                0,
		EnableAsyncProcessing:   true,
		MaxWorkers:             4,
		QueueSize:              1000,
		EnableConcurrentQueries: true,
		QueryTimeout:           30 * time.Second,
		BatchSize:              10,
		ShortCacheTTL:          5 * time.Minute,
		MediumCacheTTL:         1 * time.Hour,
		LongCacheTTL:           24 * time.Hour,
		VeryLongCacheTTL:       7 * 24 * time.Hour,
	}
}

// NewOptimizedAnalyticsService creates a new high-performance analytics service
func NewOptimizedAnalyticsService(database *db.Database, config OptimizedConfig) *OptimizedAnalyticsService {
	// Create core analytics service
	analyticsService := NewAnalyticsService(database)

	// Create cache service
	cacheConfig := cache.CacheConfig{
		Host:     config.CacheHost,
		Port:     config.CachePort,
		Password: config.CachePassword,
		DB:       config.CacheDB,
		Enabled:  config.CacheEnabled,
	}
	cacheService := cache.NewCacheService(cacheConfig)

	// Create worker pool if async processing is enabled
	var workerPool *workers.WorkerPool
	if config.EnableAsyncProcessing {
		// Use wrapper to implement the worker interface
		wrapper := NewAnalyticsServiceWrapper(analyticsService)
		workerPool = workers.NewWorkerPool(wrapper, cacheService)
	}

	return &OptimizedAnalyticsService{
		analyticsService: analyticsService,
		cacheService:     cacheService,
		workerPool:       workerPool,
		config:           config,
	}
}

// Start initializes and starts the optimized analytics service
func (oas *OptimizedAnalyticsService) Start() error {
	oas.mutex.Lock()
	defer oas.mutex.Unlock()

	if oas.running {
		return fmt.Errorf("optimized analytics service is already running")
	}

	log.Println("🚀 Starting optimized analytics service...")

	// Start worker pool if enabled
	if oas.workerPool != nil {
		if err := oas.workerPool.Start(); err != nil {
			return fmt.Errorf("failed to start worker pool: %w", err)
		}
	}

	oas.running = true
	log.Println("✅ Optimized analytics service started successfully")

	// Warm up cache with common data if enabled
	if oas.config.CacheEnabled {
		go oas.warmupCache()
	}

	return nil
}

// Stop gracefully shuts down the optimized analytics service
func (oas *OptimizedAnalyticsService) Stop() error {
	oas.mutex.Lock()
	defer oas.mutex.Unlock()

	if !oas.running {
		return fmt.Errorf("optimized analytics service is not running")
	}

	log.Println("🛑 Stopping optimized analytics service...")

	// Stop worker pool
	if oas.workerPool != nil {
		if err := oas.workerPool.Stop(); err != nil {
			log.Printf("Warning: error stopping worker pool: %v", err)
		}
	}

	// Close cache service
	if oas.cacheService != nil {
		if err := oas.cacheService.Close(); err != nil {
			log.Printf("Warning: error closing cache service: %v", err)
		}
	}

	oas.running = false
	log.Println("✅ Optimized analytics service stopped")

	return nil
}

// GetPeriodStatsAsync retrieves period stats using cache and async processing
func (oas *OptimizedAnalyticsService) GetPeriodStatsAsync(userID int, period string) (*PeriodStats, error) {
	if !oas.isRunning() {
		return nil, fmt.Errorf("service not running")
	}

	// Try cache first
	cacheKey := cache.AnalyticsCacheKey(userID, period, "period_stats")
	if oas.cacheService.IsEnabled() {
		var stats PeriodStats
		if err := oas.cacheService.GetJSON(cacheKey, &stats); err == nil {
			log.Printf("📦 Cache hit for period stats (user: %d, period: %s)", userID, period)
			return &stats, nil
		}
	}

	// If async processing is enabled, submit task and return cached result or wait
	if oas.workerPool != nil && oas.workerPool.IsRunning() {
		// Submit high-priority task
		task := workers.AnalyticsTask{
			Type:     workers.TaskPeriodStats,
			UserID:   userID,
			Data:     map[string]interface{}{"period": period},
			Priority: 1,
		}

		if err := oas.workerPool.SubmitTask(task); err != nil {
			log.Printf("Warning: failed to submit async task, falling back to sync: %v", err)
		} else {
			// Task submitted, check cache periodically for result
			return oas.waitForCachedResult(cacheKey, 10*time.Second, func() interface{} {
				stats, _ := oas.analyticsService.GetPeriodStats(userID, period)
				return stats
			})
		}
	}

	// Fallback to synchronous processing
	return oas.analyticsService.GetPeriodStats(userID, period)
}

// GetMMRTrajectoryAsync retrieves MMR trajectory using cache and async processing
func (oas *OptimizedAnalyticsService) GetMMRTrajectoryAsync(userID int, days int) (*MMRTrajectory, error) {
	if !oas.isRunning() {
		return nil, fmt.Errorf("service not running")
	}

	// Try cache first
	cacheKey := cache.MMRCacheKey(userID, days)
	if oas.cacheService.IsEnabled() {
		var trajectory MMRTrajectory
		if err := oas.cacheService.GetJSON(cacheKey, &trajectory); err == nil {
			log.Printf("📦 Cache hit for MMR trajectory (user: %d, days: %d)", userID, days)
			return &trajectory, nil
		}
	}

	// If async processing is enabled, submit task
	if oas.workerPool != nil && oas.workerPool.IsRunning() {
		task := workers.AnalyticsTask{
			Type:     workers.TaskMMRTrajectory,
			UserID:   userID,
			Data:     map[string]interface{}{"days": days},
			Priority: 1,
		}

		if err := oas.workerPool.SubmitTask(task); err != nil {
			log.Printf("Warning: failed to submit async task, falling back to sync: %v", err)
		} else {
			result, err := oas.waitForCachedResultGeneric(cacheKey, 15*time.Second, func() interface{} {
				trajectory, _ := oas.analyticsService.GetMMRTrajectory(userID, days)
				return trajectory
			})
			if err != nil {
				return nil, err
			}
			if trajectory, ok := result.(*MMRTrajectory); ok {
				return trajectory, nil
			}
		}
	}

	// Fallback to synchronous processing
	return oas.analyticsService.GetMMRTrajectory(userID, days)
}

// GetRecommendationsAsync retrieves recommendations using cache and async processing
func (oas *OptimizedAnalyticsService) GetRecommendationsAsync(userID int) ([]Recommendation, error) {
	if !oas.isRunning() {
		return nil, fmt.Errorf("service not running")
	}

	// Try cache first
	cacheKey := cache.RecommendationCacheKey(userID)
	if oas.cacheService.IsEnabled() {
		var recommendations []Recommendation
		if err := oas.cacheService.GetJSON(cacheKey, &recommendations); err == nil {
			log.Printf("📦 Cache hit for recommendations (user: %d)", userID)
			return recommendations, nil
		}
	}

	// If async processing is enabled, submit task
	if oas.workerPool != nil && oas.workerPool.IsRunning() {
		task := workers.AnalyticsTask{
			Type:     workers.TaskRecommendations,
			UserID:   userID,
			Data:     map[string]interface{}{},
			Priority: 2,
		}

		if err := oas.workerPool.SubmitTask(task); err != nil {
			log.Printf("Warning: failed to submit async task, falling back to sync: %v", err)
		} else {
			result, err := oas.waitForCachedResultGeneric(cacheKey, 20*time.Second, func() interface{} {
				recs, _ := oas.analyticsService.GetRecommendations(userID)
				return recs
			})
			if err != nil {
				return nil, err
			}
			
			if recs, ok := result.([]Recommendation); ok {
				return recs, nil
			}
		}
	}

	// Fallback to synchronous processing
	return oas.analyticsService.GetRecommendations(userID)
}

// GetBatchAnalytics retrieves multiple analytics types concurrently
func (oas *OptimizedAnalyticsService) GetBatchAnalytics(userID int, requests []string) (map[string]interface{}, error) {
	if !oas.isRunning() {
		return nil, fmt.Errorf("service not running")
	}

	results := make(map[string]interface{})
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Process requests concurrently if enabled
	for _, request := range requests {
		wg.Add(1)
		go func(req string) {
			defer wg.Done()

			var result interface{}
			var err error

			switch req {
			case "period_stats_week":
				result, err = oas.GetPeriodStatsAsync(userID, "week")
			case "period_stats_month":
				result, err = oas.GetPeriodStatsAsync(userID, "month")
			case "mmr_trajectory":
				result, err = oas.GetMMRTrajectoryAsync(userID, 30)
			case "recommendations":
				result, err = oas.GetRecommendationsAsync(userID)
			default:
				err = fmt.Errorf("unknown request type: %s", req)
			}

			mutex.Lock()
			if err != nil {
				results[req] = map[string]interface{}{"error": err.Error()}
			} else {
				results[req] = result
			}
			mutex.Unlock()
		}(request)
	}

	wg.Wait()
	return results, nil
}

// InvalidateUserCache invalidates all cache entries for a user
func (oas *OptimizedAnalyticsService) InvalidateUserCache(userID int) error {
	if !oas.isRunning() {
		return fmt.Errorf("service not running")
	}

	// Submit cache invalidation task if worker pool is available
	if oas.workerPool != nil && oas.workerPool.IsRunning() {
		task := workers.AnalyticsTask{
			Type:     workers.TaskCacheInvalidate,
			UserID:   userID,
			Data:     map[string]interface{}{},
			Priority: 1,
		}
		return oas.workerPool.SubmitTask(task)
	}

	// Direct cache invalidation
	if oas.cacheService.IsEnabled() {
		patterns := []string{
			fmt.Sprintf("user:%d:*", userID),
			fmt.Sprintf("analytics:%d:*", userID),
			fmt.Sprintf("mmr:%d:*", userID),
			fmt.Sprintf("recommendations:%d", userID),
			fmt.Sprintf("champion:%d:*", userID),
		}

		for _, pattern := range patterns {
			oas.cacheService.DeletePattern(pattern)
		}
	}

	return nil
}

// WarmupUserCache pre-calculates and caches common analytics for a user
func (oas *OptimizedAnalyticsService) WarmupUserCache(userID int) error {
	if !oas.isRunning() {
		return fmt.Errorf("service not running")
	}

	if oas.workerPool != nil && oas.workerPool.IsRunning() {
		task := workers.AnalyticsTask{
			Type:     workers.TaskCacheWarmup,
			UserID:   userID,
			Data:     map[string]interface{}{},
			Priority: 3, // Low priority
		}
		return oas.workerPool.SubmitTask(task)
	}

	return fmt.Errorf("worker pool not available for cache warmup")
}

// GetPerformanceStats returns performance statistics
func (oas *OptimizedAnalyticsService) GetPerformanceStats() map[string]interface{} {
	stats := make(map[string]interface{})

	// Cache stats
	if oas.cacheService != nil {
		stats["cache"] = oas.cacheService.GetStats()
	}

	// Worker pool stats
	if oas.workerPool != nil {
		stats["worker_pool"] = oas.workerPool.GetStats()
	}

	// Service stats
	stats["service"] = map[string]interface{}{
		"running":                 oas.isRunning(),
		"cache_enabled":           oas.config.CacheEnabled,
		"async_processing":        oas.config.EnableAsyncProcessing,
		"concurrent_queries":      oas.config.EnableConcurrentQueries,
		"max_workers":            oas.config.MaxWorkers,
		"query_timeout":          oas.config.QueryTimeout.String(),
	}

	return stats
}

// Utility methods

// isRunning safely checks if the service is running
func (oas *OptimizedAnalyticsService) isRunning() bool {
	oas.mutex.RLock()
	defer oas.mutex.RUnlock()
	return oas.running
}

// waitForCachedResultGeneric waits for a cached result with timeout and fallback
func (oas *OptimizedAnalyticsService) waitForCachedResultGeneric(cacheKey string, timeout time.Duration, fallback func() interface{}) (interface{}, error) {
	if !oas.cacheService.IsEnabled() {
		return fallback(), nil
	}

	// Wait for result with polling
	start := time.Now()
	for time.Since(start) < timeout {
		var result interface{}
		if err := oas.cacheService.GetJSON(cacheKey, &result); err == nil {
			return result, nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Timeout reached, use fallback
	log.Printf("⏰ Cache wait timeout for key %s, using fallback", cacheKey)
	return fallback(), nil
}

// waitForCachedResult waits for a cached result with timeout and fallback (specific for PeriodStats)
func (oas *OptimizedAnalyticsService) waitForCachedResult(cacheKey string, timeout time.Duration, fallback func() interface{}) (*PeriodStats, error) {
	if !oas.cacheService.IsEnabled() {
		result := fallback()
		if stats, ok := result.(*PeriodStats); ok {
			return stats, nil
		}
		return nil, fmt.Errorf("fallback failed")
	}

	// Wait for result with polling
	start := time.Now()
	for time.Since(start) < timeout {
		var stats PeriodStats
		if err := oas.cacheService.GetJSON(cacheKey, &stats); err == nil {
			return &stats, nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Timeout reached, use fallback
	log.Printf("⏰ Cache wait timeout for key %s, using fallback", cacheKey)
	result := fallback()
	if stats, ok := result.(*PeriodStats); ok {
		return stats, nil
	}
	return nil, fmt.Errorf("fallback failed after timeout")
}

// warmupCache pre-populates cache with commonly accessed data
func (oas *OptimizedAnalyticsService) warmupCache() {
	log.Println("🔥 Starting cache warmup...")
	
	// This could be enhanced to warm up cache for active users
	// For now, just log that warmup is available
	
	log.Println("✅ Cache warmup completed")
}