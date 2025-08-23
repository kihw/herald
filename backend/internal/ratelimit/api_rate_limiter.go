package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Herald.lol Gaming Analytics - API Rate Limiter
// Specific rate limiting for API endpoints and gaming operations

// APIRateLimiter handles API-specific rate limiting
type APIRateLimiter struct {
	redis *redis.Client
}

// APIRateLimitResult contains API rate limiting result
type APIRateLimitResult struct {
	Allowed    bool          `json:"allowed"`
	Remaining  int           `json:"remaining"`
	Reset      time.Time     `json:"reset"`
	RetryAfter time.Duration `json:"retry_after,omitempty"`
	Endpoint   string        `json:"endpoint"`
	LimitType  string        `json:"limit_type"`
}

// NewAPIRateLimiter creates new API rate limiter
func NewAPIRateLimiter(redis *redis.Client) *APIRateLimiter {
	return &APIRateLimiter{
		redis: redis,
	}
}

// CheckAPILimit checks API endpoint specific limits
func (a *APIRateLimiter) CheckAPILimit(ctx context.Context, request *APIRateLimitRequest) (*APIRateLimitResult, error) {
	// Get endpoint configuration
	endpointConfig := a.getEndpointConfig(request.Endpoint)
	if endpointConfig == nil {
		// No specific limits for this endpoint
		return &APIRateLimitResult{
			Allowed:   true,
			Endpoint:  request.Endpoint,
			LimitType: "no_limit",
		}, nil
	}
	
	now := time.Now()
	
	// Check per-user endpoint limit
	userEndpointKey := fmt.Sprintf("api:user:%s:endpoint:%s:minute:%d", 
		request.UserID, request.Endpoint, now.Unix()/60)
	
	pipe := a.redis.Pipeline()
	pipe.Incr(ctx, userEndpointKey)
	pipe.Expire(ctx, userEndpointKey, time.Minute)
	results, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check API endpoint limit: %w", err)
	}
	
	count := results[0].(*redis.IntCmd).Val()
	
	if int(count) > endpointConfig.RequestsPerMinute {
		return &APIRateLimitResult{
			Allowed:    false,
			Remaining:  0,
			Reset:      now.Add(time.Minute),
			RetryAfter: time.Minute,
			Endpoint:   request.Endpoint,
			LimitType:  "endpoint_limit",
		}, nil
	}
	
	return &APIRateLimitResult{
		Allowed:   true,
		Remaining: endpointConfig.RequestsPerMinute - int(count),
		Reset:     now.Add(time.Minute),
		Endpoint:  request.Endpoint,
		LimitType: "endpoint_limit",
	}, nil
}

// APIRateLimitRequest contains API rate limit request info
type APIRateLimitRequest struct {
	UserID           string `json:"user_id"`
	Endpoint         string `json:"endpoint"`
	Method           string `json:"method"`
	SubscriptionTier string `json:"subscription_tier"`
}

// EndpointConfig defines rate limits for specific endpoints
type EndpointConfig struct {
	Endpoint          string `json:"endpoint"`
	RequestsPerMinute int    `json:"requests_per_minute"`
	RequestsPerHour   int    `json:"requests_per_hour"`
	RequiresAuth      bool   `json:"requires_auth"`
	MinTier           string `json:"min_tier"`
	IsExpensive       bool   `json:"is_expensive"`
}

// getEndpointConfig returns configuration for specific endpoint
func (a *APIRateLimiter) getEndpointConfig(endpoint string) *EndpointConfig {
	// Gaming Analytics endpoints with specific limits
	endpointConfigs := map[string]*EndpointConfig{
		// Analytics endpoints (expensive operations)
		"/api/v1/gaming/analytics/summoner/:region/:summonerName": {
			Endpoint:          endpoint,
			RequestsPerMinute: 30,
			RequestsPerHour:   1800,
			RequiresAuth:      true,
			MinTier:          "free",
			IsExpensive:      true,
		},
		"/api/v1/gaming/analytics/summoner/:region/:summonerName/export": {
			Endpoint:          endpoint,
			RequestsPerMinute: 2,
			RequestsPerHour:   120,
			RequiresAuth:      true,
			MinTier:          "premium",
			IsExpensive:      true,
		},
		"/api/v1/gaming/analytics/summoner/:region/:summonerName/trends": {
			Endpoint:          endpoint,
			RequestsPerMinute: 20,
			RequestsPerHour:   1200,
			RequiresAuth:      true,
			MinTier:          "free",
			IsExpensive:      true,
		},
		"/api/v1/gaming/analytics/summoner/:region/:summonerName/compare": {
			Endpoint:          endpoint,
			RequestsPerMinute: 15,
			RequestsPerHour:   900,
			RequiresAuth:      true,
			MinTier:          "premium",
			IsExpensive:      true,
		},
		
		// Match data endpoints
		"/api/v1/gaming/matches/:region/:matchId": {
			Endpoint:          endpoint,
			RequestsPerMinute: 60,
			RequestsPerHour:   3600,
			RequiresAuth:      true,
			MinTier:          "free",
			IsExpensive:      false,
		},
		"/api/v1/gaming/matches/:region/:matchId/analyze": {
			Endpoint:          endpoint,
			RequestsPerMinute: 10,
			RequestsPerHour:   600,
			RequiresAuth:      true,
			MinTier:          "premium",
			IsExpensive:      true,
		},
		"/api/v1/gaming/matches/:region/summoner/:summonerName/recent": {
			Endpoint:          endpoint,
			RequestsPerMinute: 30,
			RequestsPerHour:   1800,
			RequiresAuth:      true,
			MinTier:          "free",
			IsExpensive:      false,
		},
		"/api/v1/gaming/matches/:region/summoner/:summonerName/live": {
			Endpoint:          endpoint,
			RequestsPerMinute: 60,
			RequestsPerHour:   3600,
			RequiresAuth:      true,
			MinTier:          "free",
			IsExpensive:      false,
		},
		
		// Team management endpoints
		"/api/v1/gaming/teams": {
			Endpoint:          endpoint,
			RequestsPerMinute: 60,
			RequestsPerHour:   3600,
			RequiresAuth:      true,
			MinTier:          "free",
			IsExpensive:      false,
		},
		"/api/v1/gaming/teams/:teamId/analytics": {
			Endpoint:          endpoint,
			RequestsPerMinute: 20,
			RequestsPerHour:   1200,
			RequiresAuth:      true,
			MinTier:          "premium",
			IsExpensive:      true,
		},
		
		// Riot API proxy endpoints (comply with Riot limits)
		"/api/v1/riot/summoner/:region/by-name/:summonerName": {
			Endpoint:          endpoint,
			RequestsPerMinute: 50, // Stay under Riot's 100/2min limit
			RequestsPerHour:   3000,
			RequiresAuth:      true,
			MinTier:          "free",
			IsExpensive:      false,
		},
		"/api/v1/riot/league/:region/by-summoner/:summonerName": {
			Endpoint:          endpoint,
			RequestsPerMinute: 50,
			RequestsPerHour:   3000,
			RequiresAuth:      true,
			MinTier:          "free",
			IsExpensive:      false,
		},
		
		// Gaming insights endpoints (AI-powered, expensive)
		"/api/v1/gaming/insights/summoner/:region/:summonerName": {
			Endpoint:          endpoint,
			RequestsPerMinute: 10,
			RequestsPerHour:   600,
			RequiresAuth:      true,
			MinTier:          "premium",
			IsExpensive:      true,
		},
		"/api/v1/gaming/insights/summoner/:region/:summonerName/coaching": {
			Endpoint:          endpoint,
			RequestsPerMinute: 5,
			RequestsPerHour:   300,
			RequiresAuth:      true,
			MinTier:          "pro",
			IsExpensive:      true,
		},
		
		// Authentication endpoints
		"/api/v1/auth/login": {
			Endpoint:          endpoint,
			RequestsPerMinute: 10,
			RequestsPerHour:   60,
			RequiresAuth:      false,
			MinTier:          "free",
			IsExpensive:      false,
		},
		"/api/v1/auth/register": {
			Endpoint:          endpoint,
			RequestsPerMinute: 5,
			RequestsPerHour:   30,
			RequiresAuth:      false,
			MinTier:          "free",
			IsExpensive:      false,
		},
		
		// System endpoints (unrestricted)
		"/api/v1/health": {
			Endpoint:          endpoint,
			RequestsPerMinute: 1000,
			RequestsPerHour:   60000,
			RequiresAuth:      false,
			MinTier:          "free",
			IsExpensive:      false,
		},
	}
	
	return endpointConfigs[endpoint]
}

// CheckRiotAPILimit checks Riot Games API proxy limits
func (a *APIRateLimiter) CheckRiotAPILimit(ctx context.Context, request *RiotAPIRequest) (*APIRateLimitResult, error) {
	now := time.Now()
	
	// Check personal development key limits (100 req/2min)
	if request.KeyType == "personal" {
		key := fmt.Sprintf("riot:personal:2min:%d", now.Unix()/120) // 2-minute window
		
		pipe := a.redis.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, 2*time.Minute)
		results, err := pipe.Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to check Riot API personal limit: %w", err)
		}
		
		count := results[0].(*redis.IntCmd).Val()
		
		if int(count) > 100 {
			return &APIRateLimitResult{
				Allowed:    false,
				Remaining:  0,
				Reset:      now.Add(2 * time.Minute),
				RetryAfter: 2 * time.Minute,
				Endpoint:   request.Endpoint,
				LimitType:  "riot_personal_limit",
			}, nil
		}
		
		return &APIRateLimitResult{
			Allowed:   true,
			Remaining: 100 - int(count),
			Reset:     now.Add(2 * time.Minute),
			Endpoint:  request.Endpoint,
			LimitType: "riot_personal_limit",
		}, nil
	}
	
	// Check production key limits (varies by endpoint)
	return a.checkRiotProductionLimits(ctx, request, now)
}

// checkRiotProductionLimits checks production Riot API limits
func (a *APIRateLimiter) checkRiotProductionLimits(ctx context.Context, request *RiotAPIRequest, now time.Time) (*APIRateLimitResult, error) {
	// Different endpoints have different Riot limits
	var limit int
	var window time.Duration
	
	switch {
	case request.Endpoint == "summoner":
		limit = 2000  // 2000 req/min for summoner endpoints
		window = time.Minute
	case request.Endpoint == "match":
		limit = 1000  // 1000 req/min for match endpoints
		window = time.Minute
	case request.Endpoint == "league":
		limit = 1500  // 1500 req/min for league endpoints
		window = time.Minute
	default:
		limit = 500   // Conservative default
		window = time.Minute
	}
	
	windowKey := fmt.Sprintf("riot:production:%s:minute:%d", request.Endpoint, now.Unix()/60)
	
	pipe := a.redis.Pipeline()
	pipe.Incr(ctx, windowKey)
	pipe.Expire(ctx, windowKey, window)
	results, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check Riot API production limit: %w", err)
	}
	
	count := results[0].(*redis.IntCmd).Val()
	
	if int(count) > limit {
		return &APIRateLimitResult{
			Allowed:    false,
			Remaining:  0,
			Reset:      now.Add(window),
			RetryAfter: window,
			Endpoint:   request.Endpoint,
			LimitType:  "riot_production_limit",
		}, nil
	}
	
	return &APIRateLimitResult{
		Allowed:   true,
		Remaining: limit - int(count),
		Reset:     now.Add(window),
		Endpoint:  request.Endpoint,
		LimitType: "riot_production_limit",
	}, nil
}

// RiotAPIRequest contains Riot API request information
type RiotAPIRequest struct {
	Endpoint string `json:"endpoint"` // summoner, match, league, etc.
	Region   string `json:"region"`
	KeyType  string `json:"key_type"` // personal, production
	UserID   string `json:"user_id"`
}

// GetEndpointStats retrieves statistics for specific endpoint
func (a *APIRateLimiter) GetEndpointStats(ctx context.Context, endpoint string, timeRange string) (*EndpointStats, error) {
	now := time.Now()
	var keys []string
	var duration time.Duration
	
	switch timeRange {
	case "minute":
		keys = append(keys, fmt.Sprintf("endpoint:%s:minute:%d", endpoint, now.Unix()/60))
		duration = time.Minute
	case "hour":
		for i := 0; i < 60; i++ {
			keys = append(keys, fmt.Sprintf("endpoint:%s:minute:%d", endpoint, (now.Unix()/60)-int64(i)))
		}
		duration = time.Hour
	case "day":
		for i := 0; i < 24; i++ {
			keys = append(keys, fmt.Sprintf("endpoint:%s:hour:%d", endpoint, (now.Unix()/3600)-int64(i)))
		}
		duration = 24 * time.Hour
	default:
		return nil, fmt.Errorf("invalid time range: %s", timeRange)
	}
	
	// Get counts for all keys
	pipe := a.redis.Pipeline()
	for _, key := range keys {
		pipe.Get(ctx, key)
	}
	results, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoint stats: %w", err)
	}
	
	totalRequests := 0
	for _, result := range results {
		if cmd, ok := result.(*redis.StringCmd); ok {
			if val, err := cmd.Result(); err == nil {
				if count, err := redis.ParseInt(val); err == nil {
					totalRequests += int(count)
				}
			}
		}
	}
	
	config := a.getEndpointConfig(endpoint)
	
	return &EndpointStats{
		Endpoint:         endpoint,
		TotalRequests:    totalRequests,
		TimeRange:        timeRange,
		Duration:         duration,
		RequestsPerMinute: totalRequests / int(duration.Minutes()),
		ConfiguredLimit:  config,
		IsAtLimit:        config != nil && totalRequests >= config.RequestsPerMinute*int(duration.Minutes()),
	}, nil
}

// EndpointStats contains endpoint usage statistics
type EndpointStats struct {
	Endpoint          string          `json:"endpoint"`
	TotalRequests     int             `json:"total_requests"`
	TimeRange         string          `json:"time_range"`
	Duration          time.Duration   `json:"duration"`
	RequestsPerMinute int             `json:"requests_per_minute"`
	ConfiguredLimit   *EndpointConfig `json:"configured_limit"`
	IsAtLimit         bool            `json:"is_at_limit"`
}

// GetTopEndpoints returns most used endpoints
func (a *APIRateLimiter) GetTopEndpoints(ctx context.Context, limit int) ([]EndpointUsage, error) {
	// This would require additional tracking of endpoint usage
	// For now, return configured endpoints with their limits
	
	endpoints := []EndpointUsage{
		{Endpoint: "/api/v1/gaming/analytics/summoner/:region/:summonerName", RequestsPerHour: 1800, IsExpensive: true},
		{Endpoint: "/api/v1/gaming/matches/:region/:matchId", RequestsPerHour: 3600, IsExpensive: false},
		{Endpoint: "/api/v1/riot/summoner/:region/by-name/:summonerName", RequestsPerHour: 3000, IsExpensive: false},
		{Endpoint: "/api/v1/gaming/teams/:teamId/analytics", RequestsPerHour: 1200, IsExpensive: true},
		{Endpoint: "/api/v1/gaming/insights/summoner/:region/:summonerName", RequestsPerHour: 600, IsExpensive: true},
	}
	
	if len(endpoints) > limit {
		endpoints = endpoints[:limit]
	}
	
	return endpoints, nil
}

// EndpointUsage contains endpoint usage information
type EndpointUsage struct {
	Endpoint         string `json:"endpoint"`
	RequestsPerHour  int    `json:"requests_per_hour"`
	IsExpensive      bool   `json:"is_expensive"`
}

// ClearEndpointCounters clears all counters for an endpoint
func (a *APIRateLimiter) ClearEndpointCounters(ctx context.Context, endpoint string) error {
	now := time.Now()
	
	// Clear current and previous windows
	keys := []string{
		fmt.Sprintf("endpoint:%s:minute:%d", endpoint, now.Unix()/60),
		fmt.Sprintf("endpoint:%s:minute:%d", endpoint, (now.Unix()/60)-1),
		fmt.Sprintf("endpoint:%s:hour:%d", endpoint, now.Unix()/3600),
		fmt.Sprintf("endpoint:%s:hour:%d", endpoint, (now.Unix()/3600)-1),
	}
	
	// Also clear user-specific counters
	pattern := fmt.Sprintf("api:user:*:endpoint:%s:*", endpoint)
	userKeys, err := a.redis.Keys(ctx, pattern).Result()
	if err == nil {
		keys = append(keys, userKeys...)
	}
	
	if len(keys) > 0 {
		err := a.redis.Del(ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("failed to clear endpoint counters: %w", err)
		}
	}
	
	return nil
}