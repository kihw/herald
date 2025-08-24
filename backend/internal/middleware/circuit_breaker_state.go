package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Herald.lol Gaming Analytics - Circuit Breaker State Management
// Manages circuit breaker states and transitions for gaming services

// CircuitBreakerStateManager manages circuit breaker states
type CircuitBreakerStateManager struct {
	redisClient *redis.Client
	config      *CircuitBreakerConfig
}

// ServiceHealthChecker checks service health
type ServiceHealthChecker struct {
	redisClient *redis.Client
	config      *CircuitBreakerConfig
}

// FallbackCache manages cached responses for fallback
type FallbackCache struct {
	redisClient *redis.Client
	config      *CircuitBreakerConfig
}

// GetCircuitInfo retrieves current circuit breaker information
func (sm *CircuitBreakerStateManager) GetCircuitInfo(ctx context.Context, serviceName string) (*CircuitBreakerInfo, error) {
	key := fmt.Sprintf("herald:circuit_breaker:%s", serviceName)
	result, err := sm.redisClient.HGetAll(ctx, key).Result()
	if err != nil || len(result) == 0 {
		// Initialize new circuit breaker
		return sm.initializeCircuit(ctx, serviceName), nil
	}

	info := &CircuitBreakerInfo{
		ServiceName: serviceName,
		State:       StateClosed, // Default
	}

	// Parse stored values
	if state, exists := result["state"]; exists {
		info.State = CircuitBreakerState(state)
	}

	if failureCountStr, exists := result["failure_count"]; exists {
		if count, err := strconv.ParseInt(failureCountStr, 10, 64); err == nil {
			info.FailureCount = count
		}
	}

	if successCountStr, exists := result["success_count"]; exists {
		if count, err := strconv.ParseInt(successCountStr, 10, 64); err == nil {
			info.SuccessCount = count
		}
	}

	if consecutiveFailuresStr, exists := result["consecutive_failures"]; exists {
		if count, err := strconv.ParseInt(consecutiveFailuresStr, 10, 64); err == nil {
			info.ConsecutiveFailures = count
		}
	}

	if lastFailureStr, exists := result["last_failure"]; exists && lastFailureStr != "" {
		if timestamp, err := strconv.ParseInt(lastFailureStr, 10, 64); err == nil {
			lastFailure := time.Unix(timestamp, 0)
			info.LastFailure = &lastFailure
		}
	}

	if lastSuccessStr, exists := result["last_success"]; exists && lastSuccessStr != "" {
		if timestamp, err := strconv.ParseInt(lastSuccessStr, 10, 64); err == nil {
			lastSuccess := time.Unix(timestamp, 0)
			info.LastSuccess = &lastSuccess
		}
	}

	if stateChangedAtStr, exists := result["state_changed_at"]; exists {
		if timestamp, err := strconv.ParseInt(stateChangedAtStr, 10, 64); err == nil {
			info.StateChangedAt = time.Unix(timestamp, 0)
		}
	}

	if nextRetryAtStr, exists := result["next_retry_at"]; exists && nextRetryAtStr != "" {
		if timestamp, err := strconv.ParseInt(nextRetryAtStr, 10, 64); err == nil {
			nextRetryAt := time.Unix(timestamp, 0)
			info.NextRetryAt = &nextRetryAt
		}
	}

	if healthStatus, exists := result["health_status"]; exists {
		info.HealthStatus = healthStatus
	}

	if errorRateStr, exists := result["error_rate"]; exists {
		if rate, err := strconv.ParseFloat(errorRateStr, 64); err == nil {
			info.ErrorRate = rate
		}
	}

	return info, nil
}

// initializeCircuit initializes a new circuit breaker
func (sm *CircuitBreakerStateManager) initializeCircuit(ctx context.Context, serviceName string) *CircuitBreakerInfo {
	now := time.Now()

	info := &CircuitBreakerInfo{
		ServiceName:         serviceName,
		State:               StateClosed,
		FailureCount:        0,
		SuccessCount:        0,
		ConsecutiveFailures: 0,
		StateChangedAt:      now,
		HealthStatus:        "unknown",
		ErrorRate:           0.0,
	}

	sm.saveCircuitInfo(ctx, info)
	return info
}

// RecordSuccess records a successful request
func (sm *CircuitBreakerStateManager) RecordSuccess(ctx context.Context, serviceName string, timestamp time.Time, duration time.Duration) {
	key := fmt.Sprintf("herald:circuit_breaker:%s", serviceName)

	// Atomic operations to update success metrics
	pipe := sm.redisClient.TxPipeline()

	pipe.HIncrBy(ctx, key, "success_count", 1)
	pipe.HSet(ctx, key, "last_success", timestamp.Unix())
	pipe.HSet(ctx, key, "consecutive_failures", 0) // Reset consecutive failures
	pipe.HSet(ctx, key, "health_status", "healthy")

	// Record response time
	responseTimeKey := fmt.Sprintf("herald:circuit_breaker:response_times:%s", serviceName)
	pipe.ZAdd(ctx, responseTimeKey, &redis.Z{
		Score:  float64(timestamp.Unix()),
		Member: duration.Milliseconds(),
	})
	pipe.Expire(ctx, responseTimeKey, sm.config.TimeWindow*2)

	pipe.Exec(ctx)

	// Update error rate
	sm.updateErrorRate(ctx, serviceName)
}

// RecordFailure records a failed request
func (sm *CircuitBreakerStateManager) RecordFailure(ctx context.Context, serviceName string, timestamp time.Time, duration time.Duration, statusCode int) {
	key := fmt.Sprintf("herald:circuit_breaker:%s", serviceName)

	// Atomic operations to update failure metrics
	pipe := sm.redisClient.TxPipeline()

	pipe.HIncrBy(ctx, key, "failure_count", 1)
	pipe.HIncrBy(ctx, key, "consecutive_failures", 1)
	pipe.HSet(ctx, key, "last_failure", timestamp.Unix())
	pipe.HSet(ctx, key, "health_status", "degraded")

	// Record failure details
	failureKey := fmt.Sprintf("herald:circuit_breaker:failures:%s", serviceName)
	failureData := map[string]interface{}{
		"timestamp":   timestamp.Unix(),
		"duration_ms": duration.Milliseconds(),
		"status_code": statusCode,
	}
	pipe.ZAdd(ctx, failureKey, &redis.Z{
		Score:  float64(timestamp.Unix()),
		Member: failureData,
	})
	pipe.Expire(ctx, failureKey, sm.config.TimeWindow*2)

	pipe.Exec(ctx)

	// Update error rate
	sm.updateErrorRate(ctx, serviceName)
}

// EvaluateCircuitState evaluates and potentially changes circuit state
func (sm *CircuitBreakerStateManager) EvaluateCircuitState(ctx context.Context, serviceName string) {
	info, err := sm.GetCircuitInfo(ctx, serviceName)
	if err != nil {
		return
	}

	now := time.Now()
	shouldTransition, newState, reason := sm.shouldTransitionState(info, now)

	if shouldTransition {
		sm.transitionToState(ctx, serviceName, info, newState, reason)
	}
}

// shouldTransitionState determines if circuit state should change
func (sm *CircuitBreakerStateManager) shouldTransitionState(info *CircuitBreakerInfo, now time.Time) (bool, CircuitBreakerState, string) {
	switch info.State {
	case StateClosed:
		// Check if we should open the circuit
		if info.ConsecutiveFailures >= int64(sm.config.ConsecutiveFailures) {
			return true, StateOpen, "consecutive failures threshold exceeded"
		}

		if info.ErrorRate > 0.5 && info.FailureCount >= int64(sm.config.FailureThreshold) {
			return true, StateOpen, "error rate and failure count thresholds exceeded"
		}

	case StateOpen:
		// Check if we should try half-open
		if info.NextRetryAt != nil && now.After(*info.NextRetryAt) {
			return true, StateHalfOpen, "retry timeout reached"
		}

	case StateHalfOpen:
		// Check if we should close or reopen
		if info.ConsecutiveFailures > 0 {
			return true, StateOpen, "failure in half-open state"
		}

		if info.SuccessCount >= int64(sm.config.SuccessThreshold) {
			return true, StateClosed, "success threshold met in half-open state"
		}

		// Check for half-open timeout
		if now.Sub(info.StateChangedAt) > sm.config.HalfOpenTimeout {
			return true, StateOpen, "half-open timeout exceeded"
		}
	}

	return false, info.State, ""
}

// transitionToState transitions circuit to new state
func (sm *CircuitBreakerStateManager) transitionToState(ctx context.Context, serviceName string, info *CircuitBreakerInfo, newState CircuitBreakerState, reason string) {
	now := time.Now()
	key := fmt.Sprintf("herald:circuit_breaker:%s", serviceName)

	// Update state information
	updates := map[string]interface{}{
		"state":             string(newState),
		"state_changed_at":  now.Unix(),
		"transition_reason": reason,
	}

	// Set next retry time for open state
	if newState == StateOpen {
		nextRetryAt := now.Add(sm.config.OpenTimeout)
		updates["next_retry_at"] = nextRetryAt.Unix()
	} else {
		updates["next_retry_at"] = ""
	}

	// Reset counters for closed state
	if newState == StateClosed {
		updates["failure_count"] = 0
		updates["success_count"] = 0
		updates["consecutive_failures"] = 0
		updates["error_rate"] = 0.0
		updates["health_status"] = "healthy"
	}

	sm.redisClient.HMSet(ctx, key, updates)

	// Log state transition
	sm.logStateTransition(ctx, serviceName, info.State, newState, reason)

	// Send alerts for critical state changes
	if newState == StateOpen {
		sm.sendServiceDownAlert(ctx, serviceName, reason)
	} else if info.State == StateOpen && newState == StateClosed {
		sm.sendServiceRecoveredAlert(ctx, serviceName)
	}
}

// TransitionToHalfOpen manually transitions circuit to half-open (for health checks)
func (sm *CircuitBreakerStateManager) TransitionToHalfOpen(ctx context.Context, serviceName string) {
	info, err := sm.GetCircuitInfo(ctx, serviceName)
	if err != nil || info.State != StateOpen {
		return
	}

	sm.transitionToState(ctx, serviceName, info, StateHalfOpen, "health check triggered")
}

// updateErrorRate calculates and updates error rate
func (sm *CircuitBreakerStateManager) updateErrorRate(ctx context.Context, serviceName string) {
	now := time.Now()
	windowStart := now.Add(-sm.config.TimeWindow)

	// Get failure and success counts in time window
	failureKey := fmt.Sprintf("herald:circuit_breaker:failures:%s", serviceName)
	successKey := fmt.Sprintf("herald:circuit_breaker:response_times:%s", serviceName)

	failureCount, _ := sm.redisClient.ZCount(ctx, failureKey, strconv.FormatInt(windowStart.Unix(), 10), "+inf").Result()
	successCount, _ := sm.redisClient.ZCount(ctx, successKey, strconv.FormatInt(windowStart.Unix(), 10), "+inf").Result()

	totalRequests := failureCount + successCount
	var errorRate float64

	if totalRequests > 0 {
		errorRate = float64(failureCount) / float64(totalRequests)
	}

	// Update error rate
	key := fmt.Sprintf("herald:circuit_breaker:%s", serviceName)
	sm.redisClient.HSet(ctx, key, "error_rate", errorRate)
}

// saveCircuitInfo saves circuit breaker info to Redis
func (sm *CircuitBreakerStateManager) saveCircuitInfo(ctx context.Context, info *CircuitBreakerInfo) {
	key := fmt.Sprintf("herald:circuit_breaker:%s", info.ServiceName)

	data := map[string]interface{}{
		"service_name":         info.ServiceName,
		"state":                string(info.State),
		"failure_count":        info.FailureCount,
		"success_count":        info.SuccessCount,
		"consecutive_failures": info.ConsecutiveFailures,
		"state_changed_at":     info.StateChangedAt.Unix(),
		"health_status":        info.HealthStatus,
		"error_rate":           info.ErrorRate,
	}

	if info.LastFailure != nil {
		data["last_failure"] = info.LastFailure.Unix()
	}

	if info.LastSuccess != nil {
		data["last_success"] = info.LastSuccess.Unix()
	}

	if info.NextRetryAt != nil {
		data["next_retry_at"] = info.NextRetryAt.Unix()
	}

	sm.redisClient.HMSet(ctx, key, data)
	sm.redisClient.Expire(ctx, key, 24*time.Hour) // Keep circuit state for 24 hours
}

// logStateTransition logs circuit breaker state transitions
func (sm *CircuitBreakerStateManager) logStateTransition(ctx context.Context, serviceName string, oldState, newState CircuitBreakerState, reason string) {
	logEntry := map[string]interface{}{
		"service":   serviceName,
		"old_state": string(oldState),
		"new_state": string(newState),
		"reason":    reason,
		"timestamp": time.Now().Unix(),
		"platform":  "herald-lol",
	}

	logKey := fmt.Sprintf("herald:circuit_breaker:transitions:%s", time.Now().Format("2006-01-02"))
	sm.redisClient.LPush(ctx, logKey, logEntry)
	sm.redisClient.Expire(ctx, logKey, 7*24*time.Hour) // Keep logs for 7 days
}

// sendServiceDownAlert sends alert when service goes down
func (sm *CircuitBreakerStateManager) sendServiceDownAlert(ctx context.Context, serviceName, reason string) {
	alert := map[string]interface{}{
		"type":            "service_down",
		"service":         serviceName,
		"reason":          reason,
		"severity":        sm.getServiceSeverity(serviceName),
		"timestamp":       time.Now().Unix(),
		"platform":        "herald-lol",
		"action_required": sm.getActionRequired(serviceName),
	}

	alertKey := "herald:circuit_breaker:alerts:service_down"
	sm.redisClient.LPush(ctx, alertKey, alert)
	sm.redisClient.LTrim(ctx, alertKey, 0, 99) // Keep last 100 alerts
}

// sendServiceRecoveredAlert sends alert when service recovers
func (sm *CircuitBreakerStateManager) sendServiceRecoveredAlert(ctx context.Context, serviceName string) {
	alert := map[string]interface{}{
		"type":      "service_recovered",
		"service":   serviceName,
		"severity":  "info",
		"timestamp": time.Now().Unix(),
		"platform":  "herald-lol",
	}

	alertKey := "herald:circuit_breaker:alerts:service_recovered"
	sm.redisClient.LPush(ctx, alertKey, alert)
	sm.redisClient.LTrim(ctx, alertKey, 0, 99) // Keep last 100 alerts
}

// CheckServiceHealth performs health check for a service
func (hc *ServiceHealthChecker) CheckServiceHealth(ctx context.Context, serviceName string) bool {
	// Get recent error patterns
	errorRate := hc.getRecentErrorRate(ctx, serviceName)
	responseTime := hc.getAverageResponseTime(ctx, serviceName)

	// Health criteria
	maxErrorRate := 0.1       // 10% error rate threshold
	maxResponseTime := 5000.0 // 5 seconds response time threshold

	if errorRate > maxErrorRate {
		hc.recordHealthCheck(ctx, serviceName, false, fmt.Sprintf("High error rate: %.2f%%", errorRate*100))
		return false
	}

	if responseTime > maxResponseTime {
		hc.recordHealthCheck(ctx, serviceName, false, fmt.Sprintf("Slow response: %.0fms", responseTime))
		return false
	}

	// Additional gaming-specific health checks
	switch serviceName {
	case "riot_api":
		return hc.checkRiotAPIHealth(ctx)
	case "analytics":
		return hc.checkAnalyticsHealth(ctx)
	default:
		hc.recordHealthCheck(ctx, serviceName, true, "Basic health check passed")
		return true
	}
}

// getRecentErrorRate gets recent error rate for service
func (hc *ServiceHealthChecker) getRecentErrorRate(ctx context.Context, serviceName string) float64 {
	now := time.Now()
	windowStart := now.Add(-hc.config.TimeWindow)

	failureKey := fmt.Sprintf("herald:circuit_breaker:failures:%s", serviceName)
	successKey := fmt.Sprintf("herald:circuit_breaker:response_times:%s", serviceName)

	failureCount, _ := hc.redisClient.ZCount(ctx, failureKey, strconv.FormatInt(windowStart.Unix(), 10), "+inf").Result()
	successCount, _ := hc.redisClient.ZCount(ctx, successKey, strconv.FormatInt(windowStart.Unix(), 10), "+inf").Result()

	totalRequests := failureCount + successCount
	if totalRequests == 0 {
		return 0.0
	}

	return float64(failureCount) / float64(totalRequests)
}

// getAverageResponseTime gets average response time for service
func (hc *ServiceHealthChecker) getAverageResponseTime(ctx context.Context, serviceName string) float64 {
	now := time.Now()
	windowStart := now.Add(-hc.config.TimeWindow)

	responseTimeKey := fmt.Sprintf("herald:circuit_breaker:response_times:%s", serviceName)

	// Get recent response times
	results, err := hc.redisClient.ZRangeByScore(ctx, responseTimeKey, &redis.ZRangeBy{
		Min: strconv.FormatInt(windowStart.Unix(), 10),
		Max: "+inf",
	}).Result()

	if err != nil || len(results) == 0 {
		return 0.0
	}

	var totalTime float64
	for _, result := range results {
		if time, err := strconv.ParseFloat(result, 64); err == nil {
			totalTime += time
		}
	}

	return totalTime / float64(len(results))
}

// checkRiotAPIHealth performs Riot API specific health check
func (hc *ServiceHealthChecker) checkRiotAPIHealth(ctx context.Context) bool {
	// Check if we're hitting rate limits
	rateLimitKey := "herald:rate_limit:riot_api_errors"
	rateLimitErrors, _ := hc.redisClient.Get(ctx, rateLimitKey).Int64()

	if rateLimitErrors > 10 { // Too many rate limit errors
		hc.recordHealthCheck(ctx, "riot_api", false, "Rate limit errors detected")
		return false
	}

	hc.recordHealthCheck(ctx, "riot_api", true, "Riot API health check passed")
	return true
}

// checkAnalyticsHealth performs analytics service health check
func (hc *ServiceHealthChecker) checkAnalyticsHealth(ctx context.Context) bool {
	// Check if analytics queries are completing
	analyticsQueueKey := "herald:analytics:processing_queue"
	queueSize, _ := hc.redisClient.LLen(ctx, analyticsQueueKey).Result()

	if queueSize > 1000 { // Analytics queue is backing up
		hc.recordHealthCheck(ctx, "analytics", false, "Analytics queue backing up")
		return false
	}

	hc.recordHealthCheck(ctx, "analytics", true, "Analytics health check passed")
	return true
}

// recordHealthCheck records health check result
func (hc *ServiceHealthChecker) recordHealthCheck(ctx context.Context, serviceName string, healthy bool, message string) {
	healthRecord := map[string]interface{}{
		"service":   serviceName,
		"healthy":   healthy,
		"message":   message,
		"timestamp": time.Now().Unix(),
		"platform":  "herald-lol",
	}

	healthKey := fmt.Sprintf("herald:circuit_breaker:health_checks:%s", serviceName)
	hc.redisClient.LPush(ctx, healthKey, healthRecord)
	hc.redisClient.LTrim(ctx, healthKey, 0, 49) // Keep last 50 health checks
	hc.redisClient.Expire(ctx, healthKey, 24*time.Hour)
}

// CacheData caches data for fallback responses
func (fc *FallbackCache) CacheData(ctx context.Context, key string, data *FallbackResponse) {
	cacheKey := fmt.Sprintf("herald:circuit_breaker:fallback:%s", key)

	serialized, err := json.Marshal(data)
	if err != nil {
		return
	}

	fc.redisClient.Set(ctx, cacheKey, serialized, fc.config.CachedResponseTTL)
}

// GetCachedData retrieves cached data for fallback
func (fc *FallbackCache) GetCachedData(ctx context.Context, key string) *FallbackResponse {
	cacheKey := fmt.Sprintf("herald:circuit_breaker:fallback:%s", key)

	data, err := fc.redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil
	}

	var cached FallbackResponse
	if err := json.Unmarshal([]byte(data), &cached); err != nil {
		return nil
	}

	// Check if data is stale
	cached.Stale = time.Now().After(cached.ExpiresAt)

	return &cached
}

// GetCachedAnalytics retrieves cached analytics data
func (fc *FallbackCache) GetCachedAnalytics(ctx context.Context, endpoint string) *FallbackResponse {
	return fc.GetCachedData(ctx, fmt.Sprintf("analytics:%s", endpoint))
}

// Helper methods

func (sm *CircuitBreakerStateManager) getServiceSeverity(serviceName string) string {
	switch serviceName {
	case "analytics", "riot_api":
		return "high"
	case "user_profile", "team_data":
		return "medium"
	default:
		return "low"
	}
}

func (sm *CircuitBreakerStateManager) getActionRequired(serviceName string) string {
	switch serviceName {
	case "analytics":
		return "Check analytics database and processing queue"
	case "riot_api":
		return "Verify Riot API status and rate limiting"
	case "user_profile":
		return "Check user service database connectivity"
	default:
		return "Investigate service health and dependencies"
	}
}
