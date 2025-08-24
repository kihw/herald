package ratelimit

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Herald.lol Gaming Analytics - Game-Specific Rate Limiter
// Rate limiting for gaming operations with context-aware limits

// GameSpecificRateLimiter handles game-specific rate limiting
type GameSpecificRateLimiter struct {
	redis  *redis.Client
	config *GamingRateLimitConfig
}

// GameRateLimitResult contains game-specific rate limiting result
type GameRateLimitResult struct {
	Allowed    bool          `json:"allowed"`
	Remaining  int           `json:"remaining"`
	Reset      time.Time     `json:"reset"`
	RetryAfter time.Duration `json:"retry_after,omitempty"`
	GameType   string        `json:"game_type"`
	Operation  string        `json:"operation"`
	LimitType  string        `json:"limit_type"`
	Region     string        `json:"region,omitempty"`
}

// GameRateLimitRequest contains game-specific request information
type GameRateLimitRequest struct {
	UserID           string `json:"user_id"`
	GameType         string `json:"game_type"` // LoL, TFT
	Operation        string `json:"operation"` // analytics, match_data, export
	Region           string `json:"region,omitempty"`
	SubscriptionTier string `json:"subscription_tier"`
	Endpoint         string `json:"endpoint"`
	IsRealTime       bool   `json:"is_real_time"`
	DataSize         string `json:"data_size"` // small, medium, large
}

// NewGameSpecificRateLimiter creates new game-specific rate limiter
func NewGameSpecificRateLimiter(redis *redis.Client, config *GamingRateLimitConfig) *GameSpecificRateLimiter {
	return &GameSpecificRateLimiter{
		redis:  redis,
		config: config,
	}
}

// CheckGameLimits checks game-specific rate limits
func (g *GameSpecificRateLimiter) CheckGameLimits(ctx context.Context, request *GamingRateLimitRequest) (*RateLimitResult, error) {
	// Convert to game-specific request
	gameRequest := &GameRateLimitRequest{
		UserID:           request.UserID,
		GameType:         request.GameType,
		Operation:        g.getOperationType(request.Endpoint),
		Region:           request.Region,
		SubscriptionTier: request.SubscriptionTier,
		Endpoint:         request.Endpoint,
		IsRealTime:       g.isRealTimeEndpoint(request.Endpoint),
		DataSize:         g.getDataSize(request.Endpoint),
	}

	// Check analytics-specific limits
	if gameRequest.Operation == "analytics" {
		return g.checkAnalyticsLimits(ctx, gameRequest)
	}

	// Check export-specific limits
	if gameRequest.Operation == "export" {
		return g.checkExportLimits(ctx, gameRequest)
	}

	// Check real-time operation limits
	if gameRequest.IsRealTime {
		return g.checkRealTimeLimits(ctx, gameRequest)
	}

	// Check region-specific limits
	if gameRequest.Region != "" {
		return g.checkRegionLimits(ctx, gameRequest)
	}

	// Default: allow request
	return &RateLimitResult{
		Allowed:   true,
		LimitType: "no_game_limit",
	}, nil
}

// checkAnalyticsLimits checks analytics-specific rate limits
func (g *GameSpecificRateLimiter) checkAnalyticsLimits(ctx context.Context, request *GameRateLimitRequest) (*RateLimitResult, error) {
	now := time.Now()

	// Get tier-specific analytics limits
	tierLimits := g.getTierLimits(request.SubscriptionTier)
	analyticsLimits := g.config.AnalyticsLimits

	var limit int
	var limitType string

	// Determine specific analytics limit based on endpoint complexity
	switch {
	case strings.Contains(request.Endpoint, "compare") || strings.Contains(request.Endpoint, "insights"):
		limit = analyticsLimits.ComparisonAnalyticsPerMin
		limitType = "comparison_analytics"
	case strings.Contains(request.Endpoint, "trends"):
		limit = analyticsLimits.AdvancedAnalyticsPerMin
		limitType = "advanced_analytics"
	case strings.Contains(request.Endpoint, "team"):
		limit = analyticsLimits.TeamAnalyticsPerMin
		limitType = "team_analytics"
	case request.IsRealTime:
		limit = analyticsLimits.RealTimeAnalyticsPerMin
		limitType = "realtime_analytics"
	default:
		limit = analyticsLimits.BasicAnalyticsPerMin
		limitType = "basic_analytics"
	}

	// Apply tier multipliers
	if request.SubscriptionTier == "premium" {
		limit = int(float64(limit) * 2.0)
	} else if request.SubscriptionTier == "pro" {
		limit = int(float64(limit) * 5.0)
	} else if request.SubscriptionTier == "enterprise" {
		limit = int(float64(limit) * 10.0)
	}

	// Check analytics limit
	analyticsKey := fmt.Sprintf("analytics:%s:%s:minute:%d",
		request.UserID, limitType, now.Unix()/60)

	pipe := g.redis.Pipeline()
	pipe.Incr(ctx, analyticsKey)
	pipe.Expire(ctx, analyticsKey, time.Minute)
	results, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check analytics limit: %w", err)
	}

	count := results[0].(*redis.IntCmd).Val()

	if int(count) > limit {
		return &RateLimitResult{
			Allowed:    false,
			Remaining:  0,
			Reset:      now.Add(time.Minute),
			RetryAfter: time.Minute,
			LimitType:  fmt.Sprintf("%s_limit", limitType),
		}, nil
	}

	return &RateLimitResult{
		Allowed:   true,
		Remaining: limit - int(count),
		Reset:     now.Add(time.Minute),
		LimitType: fmt.Sprintf("%s_limit", limitType),
	}, nil
}

// checkExportLimits checks data export rate limits
func (g *GameSpecificRateLimiter) checkExportLimits(ctx context.Context, request *GameRateLimitRequest) (*RateLimitResult, error) {
	now := time.Now()
	exportLimits := g.config.ExportLimits

	// Get tier-specific export limits
	var dailyLimit int
	switch request.SubscriptionTier {
	case "free":
		dailyLimit = 1
	case "premium":
		dailyLimit = exportLimits.JSONExportsPerDay
	case "pro":
		dailyLimit = exportLimits.JSONExportsPerDay * 5
	case "enterprise":
		dailyLimit = -1 // Unlimited
	default:
		dailyLimit = 1
	}

	if dailyLimit == -1 {
		return &RateLimitResult{
			Allowed:   true,
			LimitType: "unlimited_exports",
		}, nil
	}

	// Check daily export limit
	exportKey := fmt.Sprintf("export:%s:day:%d",
		request.UserID, now.Unix()/86400)

	pipe := g.redis.Pipeline()
	pipe.Incr(ctx, exportKey)
	pipe.Expire(ctx, exportKey, 24*time.Hour)
	results, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check export limit: %w", err)
	}

	count := results[0].(*redis.IntCmd).Val()

	if int(count) > dailyLimit {
		return &RateLimitResult{
			Allowed:    false,
			Remaining:  0,
			Reset:      now.Add(24 * time.Hour),
			RetryAfter: 24 * time.Hour,
			LimitType:  "export_daily_limit",
		}, nil
	}

	return &RateLimitResult{
		Allowed:   true,
		Remaining: dailyLimit - int(count),
		Reset:     now.Add(24 * time.Hour),
		LimitType: "export_daily_limit",
	}, nil
}

// checkRealTimeLimits checks real-time operation limits
func (g *GameSpecificRateLimiter) checkRealTimeLimits(ctx context.Context, request *GameRateLimitRequest) (*RateLimitResult, error) {
	now := time.Now()
	analyticsLimits := g.config.AnalyticsLimits

	var limit int
	switch request.SubscriptionTier {
	case "free":
		limit = analyticsLimits.RealTimeAnalyticsPerMin / 2 // Reduced for free tier
	case "premium":
		limit = analyticsLimits.RealTimeAnalyticsPerMin
	case "pro":
		limit = analyticsLimits.RealTimeAnalyticsPerMin * 3
	case "enterprise":
		limit = analyticsLimits.RealTimeAnalyticsPerMin * 10
	default:
		limit = analyticsLimits.RealTimeAnalyticsPerMin / 2
	}

	realtimeKey := fmt.Sprintf("realtime:%s:minute:%d",
		request.UserID, now.Unix()/60)

	pipe := g.redis.Pipeline()
	pipe.Incr(ctx, realtimeKey)
	pipe.Expire(ctx, realtimeKey, time.Minute)
	results, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check real-time limit: %w", err)
	}

	count := results[0].(*redis.IntCmd).Val()

	if int(count) > limit {
		return &RateLimitResult{
			Allowed:    false,
			Remaining:  0,
			Reset:      now.Add(time.Minute),
			RetryAfter: time.Minute,
			LimitType:  "realtime_limit",
		}, nil
	}

	return &RateLimitResult{
		Allowed:   true,
		Remaining: limit - int(count),
		Reset:     now.Add(time.Minute),
		LimitType: "realtime_limit",
	}, nil
}

// checkRegionLimits checks region-specific limits
func (g *GameSpecificRateLimiter) checkRegionLimits(ctx context.Context, request *GameRateLimitRequest) (*RateLimitResult, error) {
	now := time.Now()

	// Some regions might have specific limits due to Riot API constraints
	var regionMultiplier float64
	switch request.Region {
	case "KR", "JP": // High traffic regions
		regionMultiplier = 0.8
	case "OCE", "TR": // Lower traffic regions
		regionMultiplier = 1.2
	default:
		regionMultiplier = 1.0
	}

	tierLimits := g.getTierLimits(request.SubscriptionTier)
	baseLimit := tierLimits.RequestsPerMinute
	adjustedLimit := int(float64(baseLimit) * regionMultiplier)

	regionKey := fmt.Sprintf("region:%s:%s:minute:%d",
		request.Region, request.UserID, now.Unix()/60)

	pipe := g.redis.Pipeline()
	pipe.Incr(ctx, regionKey)
	pipe.Expire(ctx, regionKey, time.Minute)
	results, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check region limit: %w", err)
	}

	count := results[0].(*redis.IntCmd).Val()

	if int(count) > adjustedLimit {
		return &RateLimitResult{
			Allowed:    false,
			Remaining:  0,
			Reset:      now.Add(time.Minute),
			RetryAfter: time.Minute,
			LimitType:  "region_limit",
		}, nil
	}

	return &RateLimitResult{
		Allowed:   true,
		Remaining: adjustedLimit - int(count),
		Reset:     now.Add(time.Minute),
		LimitType: "region_limit",
	}, nil
}

// Helper functions

func (g *GameSpecificRateLimiter) getOperationType(endpoint string) string {
	switch {
	case strings.Contains(endpoint, "analytics"):
		return "analytics"
	case strings.Contains(endpoint, "export"):
		return "export"
	case strings.Contains(endpoint, "matches"):
		return "match_data"
	case strings.Contains(endpoint, "insights"):
		return "insights"
	case strings.Contains(endpoint, "teams"):
		return "team_data"
	default:
		return "general"
	}
}

func (g *GameSpecificRateLimiter) isRealTimeEndpoint(endpoint string) bool {
	realTimeEndpoints := []string{
		"/live",
		"/current",
		"/realtime",
		"/streaming",
	}

	for _, pattern := range realTimeEndpoints {
		if strings.Contains(endpoint, pattern) {
			return true
		}
	}

	return false
}

func (g *GameSpecificRateLimiter) getDataSize(endpoint string) string {
	switch {
	case strings.Contains(endpoint, "export"):
		return "large"
	case strings.Contains(endpoint, "analytics"):
		return "medium"
	case strings.Contains(endpoint, "trends"):
		return "medium"
	case strings.Contains(endpoint, "compare"):
		return "large"
	default:
		return "small"
	}
}

func (g *GameSpecificRateLimiter) getTierLimits(tier string) *TierLimits {
	switch tier {
	case "free":
		return g.config.FreeTierLimits
	case "premium":
		return g.config.PremiumTierLimits
	case "pro":
		return g.config.ProTierLimits
	case "enterprise":
		return g.config.EnterpriseLimits
	default:
		return g.config.FreeTierLimits
	}
}

// GetGameStats retrieves game-specific statistics
func (g *GameSpecificRateLimiter) GetGameStats(ctx context.Context, userID, gameType string) (*GameStats, error) {
	now := time.Now()

	// Get current counts for different operations
	keys := []string{
		fmt.Sprintf("analytics:%s:basic_analytics:minute:%d", userID, now.Unix()/60),
		fmt.Sprintf("analytics:%s:advanced_analytics:minute:%d", userID, now.Unix()/60),
		fmt.Sprintf("analytics:%s:realtime_analytics:minute:%d", userID, now.Unix()/60),
		fmt.Sprintf("export:%s:day:%d", userID, now.Unix()/86400),
	}

	pipe := g.redis.Pipeline()
	for _, key := range keys {
		pipe.Get(ctx, key)
	}
	results, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get game stats: %w", err)
	}

	stats := &GameStats{
		UserID:   userID,
		GameType: gameType,
	}

	// Parse results
	if len(results) >= 4 {
		if val, err := results[0].(*redis.StringCmd).Result(); err == nil {
			if count, err := redis.ParseInt(val); err == nil {
				stats.BasicAnalyticsThisMinute = int(count)
			}
		}
		if val, err := results[1].(*redis.StringCmd).Result(); err == nil {
			if count, err := redis.ParseInt(val); err == nil {
				stats.AdvancedAnalyticsThisMinute = int(count)
			}
		}
		if val, err := results[2].(*redis.StringCmd).Result(); err == nil {
			if count, err := redis.ParseInt(val); err == nil {
				stats.RealtimeAnalyticsThisMinute = int(count)
			}
		}
		if val, err := results[3].(*redis.StringCmd).Result(); err == nil {
			if count, err := redis.ParseInt(val); err == nil {
				stats.ExportsToday = int(count)
			}
		}
	}

	return stats, nil
}

// GameStats contains game-specific usage statistics
type GameStats struct {
	UserID                      string `json:"user_id"`
	GameType                    string `json:"game_type"`
	BasicAnalyticsThisMinute    int    `json:"basic_analytics_this_minute"`
	AdvancedAnalyticsThisMinute int    `json:"advanced_analytics_this_minute"`
	RealtimeAnalyticsThisMinute int    `json:"realtime_analytics_this_minute"`
	ExportsToday                int    `json:"exports_today"`
}

// ResetGameCounters resets all game-specific counters for a user
func (g *GameSpecificRateLimiter) ResetGameCounters(ctx context.Context, userID string) error {
	now := time.Now()

	// Get all keys for this user
	patterns := []string{
		fmt.Sprintf("analytics:%s:*", userID),
		fmt.Sprintf("export:%s:*", userID),
		fmt.Sprintf("realtime:%s:*", userID),
		fmt.Sprintf("region:*:%s:*", userID),
	}

	var allKeys []string
	for _, pattern := range patterns {
		keys, err := g.redis.Keys(ctx, pattern).Result()
		if err == nil {
			allKeys = append(allKeys, keys...)
		}
	}

	if len(allKeys) > 0 {
		err := g.redis.Del(ctx, allKeys...).Err()
		if err != nil {
			return fmt.Errorf("failed to reset game counters: %w", err)
		}
	}

	return nil
}

// CheckBurstLimit checks burst limit for high-frequency requests
func (g *GameSpecificRateLimiter) CheckBurstLimit(ctx context.Context, userID, subscriptionTier string) (*RateLimitResult, error) {
	tierLimits := g.getTierLimits(subscriptionTier)
	now := time.Now()

	// Check requests in the last 10 seconds (burst window)
	burstWindow := 10 * time.Second
	burstKey := fmt.Sprintf("burst:%s:10s:%d", userID, now.Unix()/10)

	pipe := g.redis.Pipeline()
	pipe.Incr(ctx, burstKey)
	pipe.Expire(ctx, burstKey, burstWindow)
	results, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check burst limit: %w", err)
	}

	count := results[0].(*redis.IntCmd).Val()

	if int(count) > tierLimits.BurstLimit {
		return &RateLimitResult{
			Allowed:    false,
			Remaining:  0,
			Reset:      now.Add(burstWindow),
			RetryAfter: burstWindow,
			LimitType:  "burst_limit",
		}, nil
	}

	return &RateLimitResult{
		Allowed:   true,
		Remaining: tierLimits.BurstLimit - int(count),
		Reset:     now.Add(burstWindow),
		LimitType: "burst_limit",
	}, nil
}
