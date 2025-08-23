package middleware

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// Herald.lol Gaming Analytics - Adaptive Rate Limiting
// Dynamic rate limiting that adjusts based on system load and traffic patterns

// AdaptiveLimitManager manages dynamic rate limiting
type AdaptiveLimitManager struct {
	redisClient *redis.Client
	config      *GamingRateLimitConfig
	metrics     *SystemMetrics
	predictor   *TrafficPredictor
}

// SystemMetrics tracks system performance metrics
type SystemMetrics struct {
	redisClient *redis.Client
}

// TrafficPredictor predicts traffic patterns
type TrafficPredictor struct {
	redisClient *redis.Client
}

// AdaptiveLimitState represents current adaptive limiting state
type AdaptiveLimitState struct {
	BaseLimit       int     `json:"base_limit"`
	CurrentLimit    int     `json:"current_limit"`
	AdjustmentRatio float64 `json:"adjustment_ratio"`
	SystemLoad      float64 `json:"system_load"`
	TrafficScore    float64 `json:"traffic_score"`
	LastUpdated     time.Time `json:"last_updated"`
	Reason          string  `json:"reason"`
}

// MetricsSnapshot represents system metrics at a point in time
type MetricsSnapshot struct {
	Timestamp         time.Time `json:"timestamp"`
	CPUUsage          float64   `json:"cpu_usage"`
	MemoryUsage       float64   `json:"memory_usage"`
	ActiveConnections int64     `json:"active_connections"`
	RequestsPerSecond float64   `json:"requests_per_second"`
	ErrorRate         float64   `json:"error_rate"`
	ResponseTime      float64   `json:"response_time_ms"`
	QueueDepth        int64     `json:"queue_depth"`
}

// TrafficPattern represents detected traffic patterns
type TrafficPattern struct {
	Type        string    `json:"type"`
	Intensity   float64   `json:"intensity"`
	Duration    time.Duration `json:"duration"`
	Prediction  string    `json:"prediction"`
	Confidence  float64   `json:"confidence"`
	DetectedAt  time.Time `json:"detected_at"`
}

// NewAdaptiveLimitManager creates new adaptive rate limit manager
func NewAdaptiveLimitManager(redisClient *redis.Client, config *GamingRateLimitConfig) *AdaptiveLimitManager {
	return &AdaptiveLimitManager{
		redisClient: redisClient,
		config:      config,
		metrics: &SystemMetrics{
			redisClient: redisClient,
		},
		predictor: &TrafficPredictor{
			redisClient: redisClient,
		},
	}
}

// GetAdaptiveLimit returns dynamically adjusted rate limit
func (alm *AdaptiveLimitManager) GetAdaptiveLimit(ctx context.Context, baseLimit int) int {
	// Get current adaptive state
	state, err := alm.getCurrentState(ctx, baseLimit)
	if err != nil {
		return baseLimit // Fallback to base limit on error
	}
	
	// Check if state needs updating
	if time.Since(state.LastUpdated) > 30*time.Second {
		state = alm.updateAdaptiveState(ctx, baseLimit)
	}
	
	return state.CurrentLimit
}

// getCurrentState retrieves current adaptive limiting state
func (alm *AdaptiveLimitManager) getCurrentState(ctx context.Context, baseLimit int) (*AdaptiveLimitState, error) {
	key := "herald:adaptive_limits:state"
	result, err := alm.redisClient.HGetAll(ctx, key).Result()
	if err != nil || len(result) == 0 {
		// Initialize new state
		return alm.initializeState(ctx, baseLimit), nil
	}
	
	state := &AdaptiveLimitState{
		BaseLimit: baseLimit,
	}
	
	// Parse stored values
	if currentLimitStr, exists := result["current_limit"]; exists {
		if limit, err := strconv.Atoi(currentLimitStr); err == nil {
			state.CurrentLimit = limit
		}
	}
	
	if adjustmentRatioStr, exists := result["adjustment_ratio"]; exists {
		if ratio, err := strconv.ParseFloat(adjustmentRatioStr, 64); err == nil {
			state.AdjustmentRatio = ratio
		}
	}
	
	if systemLoadStr, exists := result["system_load"]; exists {
		if load, err := strconv.ParseFloat(systemLoadStr, 64); err == nil {
			state.SystemLoad = load
		}
	}
	
	if trafficScoreStr, exists := result["traffic_score"]; exists {
		if score, err := strconv.ParseFloat(trafficScoreStr, 64); err == nil {
			state.TrafficScore = score
		}
	}
	
	if lastUpdatedStr, exists := result["last_updated"]; exists {
		if timestamp, err := strconv.ParseInt(lastUpdatedStr, 10, 64); err == nil {
			state.LastUpdated = time.Unix(timestamp, 0)
		}
	}
	
	if reason, exists := result["reason"]; exists {
		state.Reason = reason
	}
	
	return state, nil
}

// updateAdaptiveState updates adaptive limiting state based on current conditions
func (alm *AdaptiveLimitManager) updateAdaptiveState(ctx context.Context, baseLimit int) *AdaptiveLimitState {
	// Collect current system metrics
	metrics := alm.metrics.CollectMetrics(ctx)
	
	// Analyze traffic patterns
	patterns := alm.predictor.AnalyzeTrafficPatterns(ctx)
	
	// Calculate system load score
	systemLoad := alm.calculateSystemLoad(metrics)
	
	// Calculate traffic complexity score
	trafficScore := alm.calculateTrafficScore(patterns)
	
	// Determine adjustment ratio
	adjustmentRatio, reason := alm.calculateAdjustmentRatio(systemLoad, trafficScore, patterns)
	
	// Apply adjustment
	currentLimit := int(float64(baseLimit) * adjustmentRatio)
	
	// Apply bounds
	minLimit := int(float64(baseLimit) * 0.1)  // Never go below 10% of base
	maxLimit := int(float64(baseLimit) * 3.0)  // Never exceed 300% of base
	
	if currentLimit < minLimit {
		currentLimit = minLimit
		reason = fmt.Sprintf("%s (limited to minimum)", reason)
	}
	if currentLimit > maxLimit {
		currentLimit = maxLimit
		reason = fmt.Sprintf("%s (limited to maximum)", reason)
	}
	
	state := &AdaptiveLimitState{
		BaseLimit:       baseLimit,
		CurrentLimit:    currentLimit,
		AdjustmentRatio: adjustmentRatio,
		SystemLoad:      systemLoad,
		TrafficScore:    trafficScore,
		LastUpdated:     time.Now(),
		Reason:          reason,
	}
	
	// Store updated state
	alm.storeState(ctx, state)
	
	// Log adjustment for monitoring
	alm.logAdjustment(ctx, state)
	
	return state
}

// initializeState creates initial adaptive state
func (alm *AdaptiveLimitManager) initializeState(ctx context.Context, baseLimit int) *AdaptiveLimitState {
	state := &AdaptiveLimitState{
		BaseLimit:       baseLimit,
		CurrentLimit:    baseLimit,
		AdjustmentRatio: 1.0,
		SystemLoad:      0.5,
		TrafficScore:    0.5,
		LastUpdated:     time.Now(),
		Reason:          "initialized",
	}
	
	alm.storeState(ctx, state)
	return state
}

// calculateSystemLoad calculates overall system load score (0.0 to 1.0)
func (alm *AdaptiveLimitManager) calculateSystemLoad(metrics *MetricsSnapshot) float64 {
	// Weight different metrics
	cpuWeight := 0.3
	memoryWeight := 0.2
	connectionWeight := 0.2
	errorWeight := 0.15
	responseTimeWeight := 0.15
	
	// Normalize metrics to 0-1 scale
	cpuScore := math.Min(metrics.CPUUsage/100.0, 1.0)
	memoryScore := math.Min(metrics.MemoryUsage/100.0, 1.0)
	
	// Connection load (normalized against expected capacity)
	connectionScore := math.Min(float64(metrics.ActiveConnections)/10000.0, 1.0)
	
	// Error rate score
	errorScore := math.Min(metrics.ErrorRate/10.0, 1.0) // 10% error rate = max score
	
	// Response time score (normalize against target of 100ms)
	responseTimeScore := math.Min(metrics.ResponseTime/500.0, 1.0) // 500ms = max score
	
	// Calculate weighted average
	systemLoad := (cpuScore*cpuWeight +
		memoryScore*memoryWeight +
		connectionScore*connectionWeight +
		errorScore*errorWeight +
		responseTimeScore*responseTimeWeight)
	
	return systemLoad
}

// calculateTrafficScore calculates traffic complexity score (0.0 to 1.0)
func (alm *AdaptiveLimitManager) calculateTrafficScore(patterns []*TrafficPattern) float64 {
	if len(patterns) == 0 {
		return 0.3 // Default moderate traffic score
	}
	
	var totalScore float64
	var totalWeight float64
	
	for _, pattern := range patterns {
		weight := pattern.Confidence
		var patternScore float64
		
		switch pattern.Type {
		case "burst":
			patternScore = 0.8 + (pattern.Intensity * 0.2) // Burst traffic is challenging
		case "sustained_high":
			patternScore = 0.7 + (pattern.Intensity * 0.3) // Sustained high load
		case "flash_crowd":
			patternScore = 0.9 // Flash crowds are very challenging
		case "gaming_event":
			patternScore = 0.6 + (pattern.Intensity * 0.4) // Gaming events vary
		case "normal":
			patternScore = 0.3 // Normal traffic
		case "low":
			patternScore = 0.1 // Low traffic
		default:
			patternScore = 0.5 // Unknown pattern
		}
		
		totalScore += patternScore * weight
		totalWeight += weight
	}
	
	if totalWeight > 0 {
		return totalScore / totalWeight
	}
	
	return 0.3 // Default score
}

// calculateAdjustmentRatio determines how to adjust the rate limit
func (alm *AdaptiveLimitManager) calculateAdjustmentRatio(systemLoad, trafficScore float64, patterns []*TrafficPattern) (float64, string) {
	// Base adjustment calculation
	loadFactor := 1.0 - systemLoad      // Higher load = lower limits
	trafficFactor := 1.0 - trafficScore // Complex traffic = lower limits
	
	// Combine factors
	baseFactor := (loadFactor + trafficFactor) / 2.0
	
	// Apply gaming-specific adjustments
	gamingFactor := 1.0
	var reason string
	
	// Check for gaming-specific patterns
	for _, pattern := range patterns {
		switch pattern.Type {
		case "gaming_event":
			// During gaming events, be more permissive initially
			gamingFactor *= 1.2
			reason = "gaming event detected - increased limits"
		case "flash_crowd":
			// Flash crowds need quick throttling
			gamingFactor *= 0.6
			reason = "flash crowd detected - reduced limits"
		case "analytics_burst":
			// Analytics bursts are expected, allow some headroom
			gamingFactor *= 1.1
			reason = "analytics burst - slight increase"
		case "riot_api_spike":
			// Riot API spikes need careful management
			gamingFactor *= 0.8
			reason = "Riot API spike - protective reduction"
		}
	}
	
	// Time-based adjustments for gaming platform
	hour := time.Now().Hour()
	var timeBasedFactor float64
	
	if hour >= 18 && hour <= 23 { // Peak gaming hours (6 PM - 11 PM)
		timeBasedFactor = 1.3
		if reason == "" {
			reason = "peak gaming hours - increased capacity"
		}
	} else if hour >= 2 && hour <= 6 { // Low activity hours (2 AM - 6 AM)
		timeBasedFactor = 0.7
		if reason == "" {
			reason = "low activity hours - reduced capacity"
		}
	} else {
		timeBasedFactor = 1.0
	}
	
	// Calculate final adjustment ratio
	adjustmentRatio := baseFactor * gamingFactor * timeBasedFactor
	
	// Add load-based emergency adjustments
	if systemLoad > 0.9 {
		adjustmentRatio *= 0.5
		reason = "critical system load - emergency reduction"
	} else if systemLoad > 0.8 {
		adjustmentRatio *= 0.7
		reason = "high system load - protective reduction"
	} else if systemLoad < 0.3 && trafficScore < 0.4 {
		adjustmentRatio *= 1.5
		reason = "low system load - increased capacity"
	}
	
	// Ensure reasonable bounds
	if adjustmentRatio < 0.1 {
		adjustmentRatio = 0.1
		reason = "minimum safety limit applied"
	}
	if adjustmentRatio > 3.0 {
		adjustmentRatio = 3.0
		reason = "maximum expansion limit applied"
	}
	
	if reason == "" {
		reason = fmt.Sprintf("adaptive adjustment (load: %.2f, traffic: %.2f)", systemLoad, trafficScore)
	}
	
	return adjustmentRatio, reason
}

// CollectMetrics collects current system metrics
func (sm *SystemMetrics) CollectMetrics(ctx context.Context) *MetricsSnapshot {
	now := time.Now()
	
	// In a real implementation, these would collect actual metrics
	// For now, we'll simulate with Redis-based metrics
	
	metrics := &MetricsSnapshot{
		Timestamp: now,
	}
	
	// Get CPU usage (simulated from request patterns)
	cpuKey := "herald:metrics:cpu_usage"
	if cpuStr, err := sm.redisClient.Get(ctx, cpuKey).Result(); err == nil {
		if cpu, err := strconv.ParseFloat(cpuStr, 64); err == nil {
			metrics.CPUUsage = cpu
		}
	} else {
		metrics.CPUUsage = 50.0 // Default moderate CPU usage
	}
	
	// Get memory usage (simulated)
	memKey := "herald:metrics:memory_usage"
	if memStr, err := sm.redisClient.Get(ctx, memKey).Result(); err == nil {
		if mem, err := strconv.ParseFloat(memStr, 64); err == nil {
			metrics.MemoryUsage = mem
		}
	} else {
		metrics.MemoryUsage = 60.0 // Default moderate memory usage
	}
	
	// Get active connections from rate limiter data
	connPattern := "herald:rate_limit:*"
	keys, _ := sm.redisClient.Keys(ctx, connPattern).Result()
	metrics.ActiveConnections = int64(len(keys))
	
	// Calculate requests per second from recent activity
	rpsKey := "herald:metrics:requests_per_second"
	if rpsStr, err := sm.redisClient.Get(ctx, rpsKey).Result(); err == nil {
		if rps, err := strconv.ParseFloat(rpsStr, 64); err == nil {
			metrics.RequestsPerSecond = rps
		}
	} else {
		metrics.RequestsPerSecond = float64(metrics.ActiveConnections) / 60.0 // Estimate
	}
	
	// Get error rate from logs
	errorKey := "herald:metrics:error_rate"
	if errorStr, err := sm.redisClient.Get(ctx, errorKey).Result(); err == nil {
		if errorRate, err := strconv.ParseFloat(errorStr, 64); err == nil {
			metrics.ErrorRate = errorRate
		}
	} else {
		metrics.ErrorRate = 2.0 // Default 2% error rate
	}
	
	// Get response time metrics
	rtKey := "herald:metrics:response_time"
	if rtStr, err := sm.redisClient.Get(ctx, rtKey).Result(); err == nil {
		if rt, err := strconv.ParseFloat(rtStr, 64); err == nil {
			metrics.ResponseTime = rt
		}
	} else {
		metrics.ResponseTime = 150.0 // Default 150ms response time
	}
	
	// Simulate queue depth
	metrics.QueueDepth = int64(float64(metrics.ActiveConnections) * 0.1)
	
	return metrics
}

// AnalyzeTrafficPatterns analyzes current traffic patterns
func (tp *TrafficPredictor) AnalyzeTrafficPatterns(ctx context.Context) []*TrafficPattern {
	now := time.Now()
	var patterns []*TrafficPattern
	
	// Analyze request volume patterns
	if volumePattern := tp.analyzeVolumePattern(ctx); volumePattern != nil {
		patterns = append(patterns, volumePattern)
	}
	
	// Analyze gaming-specific patterns
	if gamingPattern := tp.analyzeGamingPatterns(ctx); gamingPattern != nil {
		patterns = append(patterns, gamingPattern)
	}
	
	// Analyze time-based patterns
	if timePattern := tp.analyzeTimePatterns(ctx, now); timePattern != nil {
		patterns = append(patterns, timePattern)
	}
	
	return patterns
}

// analyzeVolumePattern analyzes request volume patterns
func (tp *TrafficPredictor) analyzeVolumePattern(ctx context.Context) *TrafficPattern {
	// Get request counts for recent time windows
	windows := []time.Duration{1 * time.Minute, 5 * time.Minute, 15 * time.Minute}
	var counts []int64
	
	for _, window := range windows {
		key := fmt.Sprintf("herald:metrics:requests:%s", window.String())
		if count, err := tp.redisClient.Get(ctx, key).Int64(); err == nil {
			counts = append(counts, count)
		} else {
			counts = append(counts, 100) // Default value
		}
	}
	
	// Analyze pattern
	if len(counts) < 3 {
		return nil
	}
	
	// Check for burst pattern
	if counts[0] > counts[1]*2 && counts[1] > counts[2]*2 {
		return &TrafficPattern{
			Type:       "burst",
			Intensity:  float64(counts[0]) / float64(counts[2]),
			Duration:   time.Minute,
			Prediction: "short-term burst likely to continue",
			Confidence: 0.7,
			DetectedAt: time.Now(),
		}
	}
	
	// Check for sustained high load
	avgShort := (counts[0] + counts[1]) / 2
	if avgShort > counts[2]*1.5 {
		return &TrafficPattern{
			Type:       "sustained_high",
			Intensity:  float64(avgShort) / float64(counts[2]),
			Duration:   5 * time.Minute,
			Prediction: "sustained high load",
			Confidence: 0.6,
			DetectedAt: time.Now(),
		}
	}
	
	return &TrafficPattern{
		Type:       "normal",
		Intensity:  1.0,
		Duration:   time.Minute,
		Prediction: "normal traffic pattern",
		Confidence: 0.5,
		DetectedAt: time.Now(),
	}
}

// analyzeGamingPatterns analyzes gaming-specific traffic patterns
func (tp *TrafficPredictor) analyzeGamingPatterns(ctx context.Context) *TrafficPattern {
	// Check for analytics bursts
	analyticsKey := "herald:rate_limit:analytics:*"
	analyticsKeys, _ := tp.redisClient.Keys(ctx, analyticsKey).Result()
	
	if len(analyticsKeys) > 50 {
		return &TrafficPattern{
			Type:       "analytics_burst",
			Intensity:  float64(len(analyticsKeys)) / 50.0,
			Duration:   2 * time.Minute,
			Prediction: "analytics burst pattern",
			Confidence: 0.8,
			DetectedAt: time.Now(),
		}
	}
	
	// Check for Riot API spikes
	riotKey := "herald:rate_limit:riot_api:*"
	riotKeys, _ := tp.redisClient.Keys(ctx, riotKey).Result()
	
	if len(riotKeys) > 30 {
		return &TrafficPattern{
			Type:       "riot_api_spike",
			Intensity:  float64(len(riotKeys)) / 30.0,
			Duration:   3 * time.Minute,
			Prediction: "riot api usage spike",
			Confidence: 0.7,
			DetectedAt: time.Now(),
		}
	}
	
	return nil
}

// analyzeTimePatterns analyzes time-based patterns
func (tp *TrafficPredictor) analyzeTimePatterns(ctx context.Context, now time.Time) *TrafficPattern {
	hour := now.Hour()
	weekday := now.Weekday()
	
	// Gaming event times (simplified)
	if weekday == time.Saturday || weekday == time.Sunday {
		if hour >= 14 && hour <= 18 { // Weekend gaming hours
			return &TrafficPattern{
				Type:       "gaming_event",
				Intensity:  1.5,
				Duration:   4 * time.Hour,
				Prediction: "weekend gaming peak",
				Confidence: 0.9,
				DetectedAt: now,
			}
		}
	}
	
	// Weekday evening gaming
	if weekday >= time.Monday && weekday <= time.Friday {
		if hour >= 18 && hour <= 23 {
			return &TrafficPattern{
				Type:       "gaming_event",
				Intensity:  1.3,
				Duration:   5 * time.Hour,
				Prediction: "weekday evening gaming",
				Confidence: 0.8,
				DetectedAt: now,
			}
		}
	}
	
	return nil
}

// Helper methods

func (alm *AdaptiveLimitManager) storeState(ctx context.Context, state *AdaptiveLimitState) {
	key := "herald:adaptive_limits:state"
	data := map[string]interface{}{
		"current_limit":    state.CurrentLimit,
		"adjustment_ratio": state.AdjustmentRatio,
		"system_load":      state.SystemLoad,
		"traffic_score":    state.TrafficScore,
		"last_updated":     state.LastUpdated.Unix(),
		"reason":          state.Reason,
	}
	
	alm.redisClient.HMSet(ctx, key, data)
	alm.redisClient.Expire(ctx, key, 24*time.Hour)
}

func (alm *AdaptiveLimitManager) logAdjustment(ctx context.Context, state *AdaptiveLimitState) {
	logEntry := map[string]interface{}{
		"base_limit":       state.BaseLimit,
		"current_limit":    state.CurrentLimit,
		"adjustment_ratio": state.AdjustmentRatio,
		"system_load":      state.SystemLoad,
		"traffic_score":    state.TrafficScore,
		"reason":          state.Reason,
		"timestamp":       state.LastUpdated.Unix(),
		"platform":        "herald-lol",
	}
	
	logKey := fmt.Sprintf("herald:adaptive_limits:log:%s", time.Now().Format("2006-01-02"))
	alm.redisClient.LPush(ctx, logKey, logEntry)
	alm.redisClient.Expire(ctx, logKey, 7*24*time.Hour) // Keep logs for 7 days
}