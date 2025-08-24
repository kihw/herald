package riot

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Herald.lol Gaming Analytics - Riot API Rate Limiter
// Compliant rate limiting for Riot Games API (100 requests per 2 minutes for personal keys)

// RiotRateLimiter handles Riot API specific rate limiting
type RiotRateLimiter struct {
	redis          *redis.Client
	requestsPerMin int
	burstLimit     int
}

// RiotRateLimiterStats contains rate limiter statistics
type RiotRateLimiterStats struct {
	RequestsThisMinute int `json:"requests_this_minute"`
	RequestsLast2Min   int `json:"requests_last_2_min"`
	RequestsToday      int `json:"requests_today"`
	RateLimitHits      int `json:"rate_limit_hits"`
	BurstLimitHits     int `json:"burst_limit_hits"`
}

// NewRiotRateLimiter creates new Riot API rate limiter
func NewRiotRateLimiter(redis *redis.Client, requestsPerMin, burstLimit int) *RiotRateLimiter {
	return &RiotRateLimiter{
		redis:          redis,
		requestsPerMin: requestsPerMin,
		burstLimit:     burstLimit,
	}
}

// CheckRateLimit checks if request is allowed under Riot API limits
func (r *RiotRateLimiter) CheckRateLimit(ctx context.Context) (bool, time.Duration, error) {
	now := time.Now()

	// Check Riot's 100 requests per 2 minutes limit (personal dev key)
	allowed, waitTime, err := r.checkRiotPersonalLimit(ctx, now)
	if err != nil {
		return false, 0, err
	}
	if !allowed {
		// Increment rate limit hits
		r.incrementRateLimitHits(ctx)
		return false, waitTime, nil
	}

	// Check burst limit (prevent too many requests in short time)
	allowed, waitTime, err = r.checkBurstLimit(ctx, now)
	if err != nil {
		return false, 0, err
	}
	if !allowed {
		r.incrementBurstLimitHits(ctx)
		return false, waitTime, nil
	}

	// Request is allowed, increment counters
	if err := r.incrementCounters(ctx, now); err != nil {
		return false, 0, fmt.Errorf("failed to increment counters: %w", err)
	}

	return true, 0, nil
}

// checkRiotPersonalLimit checks Riot's personal development key limit (100/2min)
func (r *RiotRateLimiter) checkRiotPersonalLimit(ctx context.Context, now time.Time) (bool, time.Duration, error) {
	// Use 2-minute sliding window
	twoMinWindow := now.Unix() / 120 // 120 seconds = 2 minutes
	key := fmt.Sprintf("riot:personal_limit:2min:%d", twoMinWindow)

	// Get current count
	count, err := r.redis.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return false, 0, fmt.Errorf("failed to get rate limit count: %w", err)
	}

	// Check if we're at the limit
	if count >= 100 {
		// Calculate wait time until next 2-minute window
		nextWindow := (twoMinWindow + 1) * 120
		waitTime := time.Until(time.Unix(nextWindow, 0))
		return false, waitTime, nil
	}

	return true, 0, nil
}

// checkBurstLimit checks short-term burst limit
func (r *RiotRateLimiter) checkBurstLimit(ctx context.Context, now time.Time) (bool, time.Duration, error) {
	// Use 10-second burst window
	burstWindow := now.Unix() / 10 // 10-second windows
	key := fmt.Sprintf("riot:burst_limit:10s:%d", burstWindow)

	count, err := r.redis.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return false, 0, fmt.Errorf("failed to get burst limit count: %w", err)
	}

	if count >= r.burstLimit {
		// Calculate wait time until next 10-second window
		nextWindow := (burstWindow + 1) * 10
		waitTime := time.Until(time.Unix(nextWindow, 0))
		return false, waitTime, nil
	}

	return true, 0, nil
}

// incrementCounters increments all rate limiting counters
func (r *RiotRateLimiter) incrementCounters(ctx context.Context, now time.Time) error {
	pipe := r.redis.Pipeline()

	// 2-minute window (Riot limit)
	twoMinKey := fmt.Sprintf("riot:personal_limit:2min:%d", now.Unix()/120)
	pipe.Incr(ctx, twoMinKey)
	pipe.Expire(ctx, twoMinKey, 3*time.Minute) // Keep a bit longer than window

	// 10-second burst window
	burstKey := fmt.Sprintf("riot:burst_limit:10s:%d", now.Unix()/10)
	pipe.Incr(ctx, burstKey)
	pipe.Expire(ctx, burstKey, 30*time.Second)

	// Daily counter (for statistics)
	dailyKey := fmt.Sprintf("riot:daily_requests:%d", now.Unix()/86400)
	pipe.Incr(ctx, dailyKey)
	pipe.Expire(ctx, dailyKey, 48*time.Hour)

	// Minute counter (for statistics)
	minuteKey := fmt.Sprintf("riot:minute_requests:%d", now.Unix()/60)
	pipe.Incr(ctx, minuteKey)
	pipe.Expire(ctx, minuteKey, 5*time.Minute)

	_, err := pipe.Exec(ctx)
	return err
}

// incrementRateLimitHits increments rate limit hit counter
func (r *RiotRateLimiter) incrementRateLimitHits(ctx context.Context) {
	now := time.Now()
	key := fmt.Sprintf("riot:rate_limit_hits:%d", now.Unix()/3600) // Hourly
	r.redis.Incr(ctx, key)
	r.redis.Expire(ctx, key, 25*time.Hour)
}

// incrementBurstLimitHits increments burst limit hit counter
func (r *RiotRateLimiter) incrementBurstLimitHits(ctx context.Context) {
	now := time.Now()
	key := fmt.Sprintf("riot:burst_limit_hits:%d", now.Unix()/3600) // Hourly
	r.redis.Incr(ctx, key)
	r.redis.Expire(ctx, key, 25*time.Hour)
}

// GetStats returns rate limiter statistics
func (r *RiotRateLimiter) GetStats(ctx context.Context) (*RiotRateLimiterStats, error) {
	now := time.Now()
	stats := &RiotRateLimiterStats{}

	// Get requests this minute
	minuteKey := fmt.Sprintf("riot:minute_requests:%d", now.Unix()/60)
	if count, err := r.redis.Get(ctx, minuteKey).Int(); err == nil {
		stats.RequestsThisMinute = count
	}

	// Get requests in last 2 minutes (current Riot limit window)
	twoMinKey := fmt.Sprintf("riot:personal_limit:2min:%d", now.Unix()/120)
	if count, err := r.redis.Get(ctx, twoMinKey).Int(); err == nil {
		stats.RequestsLast2Min = count
	}

	// Get requests today
	dailyKey := fmt.Sprintf("riot:daily_requests:%d", now.Unix()/86400)
	if count, err := r.redis.Get(ctx, dailyKey).Int(); err == nil {
		stats.RequestsToday = count
	}

	// Get rate limit hits this hour
	rateLimitKey := fmt.Sprintf("riot:rate_limit_hits:%d", now.Unix()/3600)
	if count, err := r.redis.Get(ctx, rateLimitKey).Int(); err == nil {
		stats.RateLimitHits = count
	}

	// Get burst limit hits this hour
	burstLimitKey := fmt.Sprintf("riot:burst_limit_hits:%d", now.Unix()/3600)
	if count, err := r.redis.Get(ctx, burstLimitKey).Int(); err == nil {
		stats.BurstLimitHits = count
	}

	return stats, nil
}

// GetRemainingRequests returns remaining requests in current 2-minute window
func (r *RiotRateLimiter) GetRemainingRequests(ctx context.Context) (int, error) {
	now := time.Now()
	twoMinKey := fmt.Sprintf("riot:personal_limit:2min:%d", now.Unix()/120)

	count, err := r.redis.Get(ctx, twoMinKey).Int()
	if err != nil && err != redis.Nil {
		return 0, err
	}

	remaining := 100 - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// WaitForRateLimit waits until rate limit resets
func (r *RiotRateLimiter) WaitForRateLimit(ctx context.Context) error {
	for {
		allowed, waitTime, err := r.CheckRateLimit(ctx)
		if err != nil {
			return err
		}
		if allowed {
			return nil
		}

		// Wait for the specified time or until context is cancelled
		select {
		case <-time.After(waitTime):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// ResetCounters resets all rate limiting counters (for testing)
func (r *RiotRateLimiter) ResetCounters(ctx context.Context) error {
	patterns := []string{
		"riot:personal_limit:*",
		"riot:burst_limit:*",
		"riot:daily_requests:*",
		"riot:minute_requests:*",
		"riot:rate_limit_hits:*",
		"riot:burst_limit_hits:*",
	}

	for _, pattern := range patterns {
		keys, err := r.redis.Keys(ctx, pattern).Result()
		if err != nil {
			continue
		}
		if len(keys) > 0 {
			r.redis.Del(ctx, keys...)
		}
	}

	return nil
}

// IsRateLimited checks if currently rate limited without making a request
func (r *RiotRateLimiter) IsRateLimited(ctx context.Context) (bool, time.Duration, error) {
	now := time.Now()

	// Check 2-minute limit
	twoMinKey := fmt.Sprintf("riot:personal_limit:2min:%d", now.Unix()/120)
	count, err := r.redis.Get(ctx, twoMinKey).Int()
	if err != nil && err != redis.Nil {
		return false, 0, err
	}

	if count >= 100 {
		nextWindow := ((now.Unix() / 120) + 1) * 120
		waitTime := time.Until(time.Unix(nextWindow, 0))
		return true, waitTime, nil
	}

	// Check burst limit
	burstKey := fmt.Sprintf("riot:burst_limit:10s:%d", now.Unix()/10)
	burstCount, err := r.redis.Get(ctx, burstKey).Int()
	if err != nil && err != redis.Nil {
		return false, 0, err
	}

	if burstCount >= r.burstLimit {
		nextWindow := ((now.Unix() / 10) + 1) * 10
		waitTime := time.Until(time.Unix(nextWindow, 0))
		return true, waitTime, nil
	}

	return false, 0, nil
}
