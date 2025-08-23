package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// Herald.lol Gaming Analytics - Circuit Breaker Pattern
// Protects gaming platform from cascading failures and provides graceful degradation

// GamingCircuitBreaker implements circuit breaker pattern for gaming services
type GamingCircuitBreaker struct {
	redisClient    *redis.Client
	config         *CircuitBreakerConfig
	stateManager   *CircuitBreakerStateManager
	healthChecker  *ServiceHealthChecker
	fallbackCache  *FallbackCache
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	// Failure thresholds
	FailureThreshold        int           `json:"failure_threshold"`         // Failures to trip circuit
	SuccessThreshold        int           `json:"success_threshold"`         // Successes to close circuit
	ConsecutiveFailures     int           `json:"consecutive_failures"`      // Consecutive failures threshold
	
	// Time windows
	TimeWindow              time.Duration `json:"time_window"`               // Evaluation window
	HalfOpenTimeout         time.Duration `json:"half_open_timeout"`         // Half-open state timeout
	OpenTimeout             time.Duration `json:"open_timeout"`              // Time to stay in open state
	
	// Health check configuration
	HealthCheckInterval     time.Duration `json:"health_check_interval"`     // Health check frequency
	HealthCheckTimeout      time.Duration `json:"health_check_timeout"`      // Health check timeout
	
	// Gaming-specific settings
	GamingServicePriority   map[string]int `json:"gaming_service_priority"`  // Service priority levels
	AnalyticsServiceFallback bool          `json:"analytics_service_fallback"` // Enable analytics fallback
	RiotAPICircuitBreaker   bool          `json:"riot_api_circuit_breaker"`  // Enable Riot API protection
	
	// Response configuration
	FallbackEnabled         bool          `json:"fallback_enabled"`          // Enable fallback responses
	CachedResponseTTL       time.Duration `json:"cached_response_ttl"`       // Cache TTL for fallback
	GracefulDegradation     bool          `json:"graceful_degradation"`      // Enable graceful degradation
}

// CircuitBreakerState represents circuit breaker states
type CircuitBreakerState string

const (
	StateClosed    CircuitBreakerState = "closed"    // Normal operation
	StateOpen      CircuitBreakerState = "open"      // Circuit is open, requests fail fast
	StateHalfOpen  CircuitBreakerState = "half_open" // Testing if service recovered
)

// CircuitBreakerInfo holds circuit breaker status information
type CircuitBreakerInfo struct {
	ServiceName         string              `json:"service_name"`
	State              CircuitBreakerState  `json:"state"`
	FailureCount       int64               `json:"failure_count"`
	SuccessCount       int64               `json:"success_count"`
	ConsecutiveFailures int64              `json:"consecutive_failures"`
	LastFailure        *time.Time          `json:"last_failure,omitempty"`
	LastSuccess        *time.Time          `json:"last_success,omitempty"`
	StateChangedAt     time.Time           `json:"state_changed_at"`
	NextRetryAt        *time.Time          `json:"next_retry_at,omitempty"`
	HealthStatus       string              `json:"health_status"`
	ErrorRate          float64             `json:"error_rate"`
}

// FallbackResponse represents cached fallback data
type FallbackResponse struct {
	Data      interface{} `json:"data"`
	CachedAt  time.Time   `json:"cached_at"`
	ExpiresAt time.Time   `json:"expires_at"`
	Source    string      `json:"source"`
	Stale     bool        `json:"stale"`
}

// NewGamingCircuitBreaker creates new circuit breaker for gaming platform
func NewGamingCircuitBreaker(redisClient *redis.Client, config *CircuitBreakerConfig) *GamingCircuitBreaker {
	// Set gaming-specific defaults
	if config.FailureThreshold == 0 {
		config.FailureThreshold = 10 // 10 failures to trip circuit
	}
	if config.SuccessThreshold == 0 {
		config.SuccessThreshold = 3 // 3 successes to close circuit
	}
	if config.ConsecutiveFailures == 0 {
		config.ConsecutiveFailures = 5 // 5 consecutive failures
	}
	if config.TimeWindow == 0 {
		config.TimeWindow = time.Minute // 1 minute evaluation window
	}
	if config.HalfOpenTimeout == 0 {
		config.HalfOpenTimeout = 30 * time.Second
	}
	if config.OpenTimeout == 0 {
		config.OpenTimeout = 60 * time.Second // Stay open for 1 minute
	}
	if config.HealthCheckInterval == 0 {
		config.HealthCheckInterval = 30 * time.Second
	}
	if config.HealthCheckTimeout == 0 {
		config.HealthCheckTimeout = 5 * time.Second
	}
	if config.CachedResponseTTL == 0 {
		config.CachedResponseTTL = 10 * time.Minute
	}
	
	// Initialize service priorities if not set
	if config.GamingServicePriority == nil {
		config.GamingServicePriority = map[string]int{
			"analytics":     10, // Highest priority
			"riot_api":      9,  // High priority
			"user_profile":  8,  // High priority
			"team_data":     7,  // Medium-high priority
			"match_history": 6,  // Medium priority
			"leaderboard":   5,  // Medium priority
			"notifications": 3,  // Low priority
			"social":        2,  // Low priority
			"cosmetics":     1,  // Lowest priority
		}
	}
	
	cb := &GamingCircuitBreaker{
		redisClient: redisClient,
		config:      config,
		stateManager: &CircuitBreakerStateManager{
			redisClient: redisClient,
			config:      config,
		},
		healthChecker: &ServiceHealthChecker{
			redisClient: redisClient,
			config:      config,
		},
		fallbackCache: &FallbackCache{
			redisClient: redisClient,
			config:      config,
		},
	}
	
	// Start background health checking
	go cb.startHealthCheckRoutine()
	
	return cb
}

// CircuitBreakerMiddleware creates circuit breaker middleware
func (cb *GamingCircuitBreaker) CircuitBreakerMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check circuit breaker state
		info, err := cb.stateManager.GetCircuitInfo(c.Request.Context(), serviceName)
		if err != nil {
			c.Next() // Continue on error
			return
		}
		
		// Handle different circuit states
		switch info.State {
		case StateOpen:
			// Circuit is open - fail fast or return fallback
			cb.handleOpenCircuit(c, serviceName, info)
			return
			
		case StateHalfOpen:
			// Half-open state - allow limited requests through
			if !cb.allowHalfOpenRequest(c, serviceName, info) {
				cb.handleOpenCircuit(c, serviceName, info)
				return
			}
			
		case StateClosed:
			// Normal operation - continue
		}
		
		// Set circuit breaker context
		c.Set("circuit_breaker_service", serviceName)
		c.Set("circuit_breaker_info", info)
		
		// Execute request with monitoring
		cb.executeWithMonitoring(c, serviceName)
	}
}

// GamingAnalyticsCircuitBreaker specialized circuit breaker for analytics
func (cb *GamingCircuitBreaker) GamingAnalyticsCircuitBreaker() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceName := "analytics"
		
		// Check if analytics service is healthy
		info, _ := cb.stateManager.GetCircuitInfo(c.Request.Context(), serviceName)
		
		if info.State == StateOpen {
			// Serve from cache if available
			if cached := cb.fallbackCache.GetCachedAnalytics(c.Request.Context(), c.Request.URL.Path); cached != nil {
				c.Header("X-Gaming-Fallback", "cached_analytics")
				c.Header("X-Gaming-Cache-Age", strconv.Itoa(int(time.Since(cached.CachedAt).Seconds())))
				
				c.JSON(http.StatusOK, gin.H{
					"data":           cached.Data,
					"cached":         true,
					"cached_at":      cached.CachedAt,
					"gaming_warning": "Analytics service temporarily unavailable - serving cached data",
					"gaming_platform": "herald-lol",
				})
				c.Abort()
				return
			}
			
			// No cache available - return degraded response
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":            "Gaming analytics temporarily unavailable",
				"service":          "analytics",
				"circuit_state":    info.State,
				"retry_after":      info.NextRetryAt,
				"gaming_platform":  "herald-lol",
				"degraded_info":    "Basic gaming stats may be available through alternative endpoints",
			})
			c.Abort()
			return
		}
		
		cb.executeWithMonitoring(c, serviceName)
	}
}

// RiotAPICircuitBreaker specialized circuit breaker for Riot API
func (cb *GamingCircuitBreaker) RiotAPICircuitBreaker() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceName := "riot_api"
		
		if !cb.config.RiotAPICircuitBreaker {
			c.Next()
			return
		}
		
		info, _ := cb.stateManager.GetCircuitInfo(c.Request.Context(), serviceName)
		
		if info.State == StateOpen {
			// Riot API is down - check for cached data
			cacheKey := fmt.Sprintf("riot_api:%s", c.Request.URL.Path)
			if cached := cb.fallbackCache.GetCachedData(c.Request.Context(), cacheKey); cached != nil {
				c.Header("X-Gaming-Riot-API-Fallback", "cached")
				c.Header("X-Gaming-Cache-Age", strconv.Itoa(int(time.Since(cached.CachedAt).Seconds())))
				
				// Serve cached Riot API data
				c.JSON(http.StatusOK, gin.H{
					"data":           cached.Data,
					"cached":         true,
					"stale":          cached.Stale,
					"riot_api_note":  "Riot API temporarily unavailable - serving recent data",
					"gaming_platform": "herald-lol",
				})
				c.Abort()
				return
			}
			
			// No cache - return error with Riot-specific messaging
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":           "Riot Games API temporarily unavailable",
				"service":         "riot_api",
				"circuit_state":   info.State,
				"retry_after":     info.NextRetryAt,
				"gaming_platform": "herald-lol",
				"riot_status":     "https://status.riotgames.com/",
				"alternative":     "Try again in a few minutes or check Riot's service status",
			})
			c.Abort()
			return
		}
		
		cb.executeWithMonitoring(c, serviceName)
	}
}

// executeWithMonitoring executes request with circuit breaker monitoring
func (cb *GamingCircuitBreaker) executeWithMonitoring(c *gin.Context, serviceName string) {
	startTime := time.Now()
	
	// Execute the actual request
	c.Next()
	
	duration := time.Since(startTime)
	statusCode := c.Writer.Status()
	
	// Record request result
	success := statusCode >= 200 && statusCode < 500 // Client errors don't count as service failures
	cb.recordRequestResult(c.Request.Context(), serviceName, success, duration, statusCode)
	
	// Cache successful responses for fallback
	if success && cb.config.FallbackEnabled {
		cb.cacheSuccessfulResponse(c, serviceName)
	}
}

// handleOpenCircuit handles requests when circuit is open
func (cb *GamingCircuitBreaker) handleOpenCircuit(c *gin.Context, serviceName string, info *CircuitBreakerInfo) {
	// Set circuit breaker headers
	c.Header("X-Gaming-Circuit-Breaker", "open")
	c.Header("X-Gaming-Service", serviceName)
	c.Header("X-Gaming-Failure-Count", strconv.Itoa(int(info.FailureCount)))
	
	if info.NextRetryAt != nil {
		c.Header("Retry-After", strconv.Itoa(int(time.Until(*info.NextRetryAt).Seconds())))
	}
	
	// Try to serve fallback response
	if cb.config.FallbackEnabled {
		if fallback := cb.getFallbackResponse(c, serviceName); fallback != nil {
			c.Header("X-Gaming-Fallback", "cached")
			c.JSON(http.StatusOK, fallback)
			c.Abort()
			return
		}
	}
	
	// No fallback available - return error
	priority := cb.getServicePriority(serviceName)
	
	c.JSON(http.StatusServiceUnavailable, gin.H{
		"error":           fmt.Sprintf("Gaming service '%s' temporarily unavailable", serviceName),
		"service":         serviceName,
		"circuit_state":   string(info.State),
		"priority":        priority,
		"failure_count":   info.FailureCount,
		"error_rate":      info.ErrorRate,
		"retry_after":     info.NextRetryAt,
		"gaming_platform": "herald-lol",
		"support_info":    cb.getSupportInfo(serviceName),
	})
	c.Abort()
}

// allowHalfOpenRequest determines if request should be allowed in half-open state
func (cb *GamingCircuitBreaker) allowHalfOpenRequest(c *gin.Context, serviceName string, info *CircuitBreakerInfo) bool {
	// Allow only one request at a time in half-open state
	key := fmt.Sprintf("herald:circuit_breaker:half_open_lock:%s", serviceName)
	
	// Try to acquire lock for half-open request
	acquired, err := cb.redisClient.SetNX(c.Request.Context(), key, time.Now().Unix(), cb.config.HalfOpenTimeout).Result()
	if err != nil || !acquired {
		return false // Another request is already testing the service
	}
	
	// Set cleanup for the lock
	c.Header("X-Gaming-Circuit-Half-Open", "testing")
	
	return true
}

// recordRequestResult records the result of a request for circuit breaker evaluation
func (cb *GamingCircuitBreaker) recordRequestResult(ctx context.Context, serviceName string, success bool, duration time.Duration, statusCode int) {
	now := time.Now()
	
	// Update counters
	if success {
		cb.stateManager.RecordSuccess(ctx, serviceName, now, duration)
	} else {
		cb.stateManager.RecordFailure(ctx, serviceName, now, duration, statusCode)
	}
	
	// Check if circuit state should change
	cb.stateManager.EvaluateCircuitState(ctx, serviceName)
}

// cacheSuccessfulResponse caches successful response for fallback
func (cb *GamingCircuitBreaker) cacheSuccessfulResponse(c *gin.Context, serviceName string) {
	// Only cache GET requests
	if c.Request.Method != "GET" {
		return
	}
	
	// Get response data (simplified - in real implementation would capture response body)
	cacheKey := fmt.Sprintf("%s:%s", serviceName, c.Request.URL.Path)
	
	// Cache placeholder data (in real implementation, would cache actual response)
	fallback := &FallbackResponse{
		Data:      gin.H{"status": "cached", "service": serviceName},
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(cb.config.CachedResponseTTL),
		Source:    "circuit_breaker_cache",
		Stale:     false,
	}
	
	cb.fallbackCache.CacheData(c.Request.Context(), cacheKey, fallback)
}

// getFallbackResponse gets fallback response for a service
func (cb *GamingCircuitBreaker) getFallbackResponse(c *gin.Context, serviceName string) interface{} {
	cacheKey := fmt.Sprintf("%s:%s", serviceName, c.Request.URL.Path)
	cached := cb.fallbackCache.GetCachedData(c.Request.Context(), cacheKey)
	
	if cached != nil {
		return gin.H{
			"data":              cached.Data,
			"cached":            true,
			"cached_at":         cached.CachedAt,
			"stale":             cached.Stale,
			"gaming_fallback":   "Circuit breaker serving cached data",
			"gaming_platform":   "herald-lol",
		}
	}
	
	// Return default degraded response for specific services
	return cb.getDefaultDegradedResponse(serviceName)
}

// getDefaultDegradedResponse returns default degraded response
func (cb *GamingCircuitBreaker) getDefaultDegradedResponse(serviceName string) interface{} {
	switch serviceName {
	case "analytics":
		return gin.H{
			"error":           "Analytics service temporarily unavailable",
			"degraded_mode":   true,
			"basic_stats":     gin.H{"status": "limited_data_available"},
			"gaming_platform": "herald-lol",
		}
	case "riot_api":
		return gin.H{
			"error":           "Riot API temporarily unavailable",
			"riot_status":     "https://status.riotgames.com/",
			"gaming_platform": "herald-lol",
		}
	case "leaderboard":
		return gin.H{
			"error":           "Leaderboard temporarily unavailable",
			"cached_rankings": "Check back in a few minutes for updated rankings",
			"gaming_platform": "herald-lol",
		}
	default:
		return gin.H{
			"error":           fmt.Sprintf("Service '%s' temporarily unavailable", serviceName),
			"gaming_platform": "herald-lol",
		}
	}
}

// getServicePriority gets service priority level
func (cb *GamingCircuitBreaker) getServicePriority(serviceName string) int {
	if priority, exists := cb.config.GamingServicePriority[serviceName]; exists {
		return priority
	}
	return 5 // Default medium priority
}

// getSupportInfo returns support information for service failures
func (cb *GamingCircuitBreaker) getSupportInfo(serviceName string) gin.H {
	info := gin.H{
		"status_page": "https://status.herald.lol",
		"support":     "support@herald.lol",
	}
	
	switch serviceName {
	case "riot_api":
		info["riot_status"] = "https://status.riotgames.com/"
		info["note"] = "This may be due to Riot Games API maintenance"
	case "analytics":
		info["note"] = "Your gaming data is safe - analytics will resume once service is restored"
	}
	
	return info
}

// startHealthCheckRoutine starts background health checking
func (cb *GamingCircuitBreaker) startHealthCheckRoutine() {
	ticker := time.NewTicker(cb.config.HealthCheckInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		ctx := context.Background()
		cb.performHealthChecks(ctx)
	}
}

// performHealthChecks performs health checks on all services
func (cb *GamingCircuitBreaker) performHealthChecks(ctx context.Context) {
	services := []string{"analytics", "riot_api", "user_profile", "team_data", "match_history"}
	
	for _, service := range services {
		info, err := cb.stateManager.GetCircuitInfo(ctx, service)
		if err != nil {
			continue
		}
		
		// Only check health for open circuits
		if info.State == StateOpen {
			healthy := cb.healthChecker.CheckServiceHealth(ctx, service)
			if healthy {
				// Service appears to be healthy - transition to half-open
				cb.stateManager.TransitionToHalfOpen(ctx, service)
			}
		}
	}
}