package ratelimit

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Herald.lol Gaming Analytics - IP-based Rate Limiter
// Advanced IP rate limiting with geographic and pattern-based blocking

// IPRateLimiter handles IP-based rate limiting
type IPRateLimiter struct {
	redis  *redis.Client
	config *IPRateLimitConfig
}

// IPRateLimitResult contains IP rate limiting result
type IPRateLimitResult struct {
	Allowed    bool          `json:"allowed"`
	Remaining  int           `json:"remaining"`
	Reset      time.Time     `json:"reset"`
	RetryAfter time.Duration `json:"retry_after,omitempty"`
	Reason     string        `json:"reason,omitempty"`
	ClientIP   string        `json:"client_ip"`
}

// NewIPRateLimiter creates new IP rate limiter
func NewIPRateLimiter(redis *redis.Client, config *IPRateLimitConfig) *IPRateLimiter {
	return &IPRateLimiter{
		redis:  redis,
		config: config,
	}
}

// CheckIPLimit checks IP-based rate limits
func (i *IPRateLimiter) CheckIPLimit(ctx context.Context, clientIP string) (*RateLimitResult, error) {
	// Clean and validate IP
	cleanIP := i.cleanIP(clientIP)
	if cleanIP == "" {
		return &RateLimitResult{
			Allowed:   false,
			LimitType: "invalid_ip",
		}, nil
	}

	// Check if IP is blacklisted
	if i.isBlacklisted(cleanIP) {
		return &RateLimitResult{
			Allowed:   false,
			LimitType: "blacklisted_ip",
		}, nil
	}

	// Check if IP is whitelisted (skip other checks)
	if i.isWhitelisted(cleanIP) {
		return &RateLimitResult{
			Allowed:   true,
			LimitType: "whitelisted_ip",
		}, nil
	}

	// Check rate limits
	now := time.Now()

	// Minute-based limiting
	minuteResult, err := i.checkIPWindow(ctx, cleanIP, "minute", i.config.RequestsPerMinute, time.Minute, now)
	if err != nil {
		return nil, fmt.Errorf("failed to check IP minute limit: %w", err)
	}
	if !minuteResult.Allowed {
		return minuteResult, nil
	}

	// Hour-based limiting
	hourResult, err := i.checkIPWindow(ctx, cleanIP, "hour", i.config.RequestsPerHour, time.Hour, now)
	if err != nil {
		return nil, fmt.Errorf("failed to check IP hour limit: %w", err)
	}
	if !hourResult.Allowed {
		return hourResult, nil
	}

	return minuteResult, nil
}

// checkIPWindow checks rate limit for specific time window
func (i *IPRateLimiter) checkIPWindow(ctx context.Context, ip, window string, limit int, duration time.Duration, now time.Time) (*RateLimitResult, error) {
	var windowKey string
	var windowStart int64

	switch window {
	case "minute":
		windowStart = now.Unix() / 60
		windowKey = fmt.Sprintf("ip:%s:minute:%d", ip, windowStart)
	case "hour":
		windowStart = now.Unix() / 3600
		windowKey = fmt.Sprintf("ip:%s:hour:%d", ip, windowStart)
	default:
		return nil, fmt.Errorf("unknown window type: %s", window)
	}

	pipe := i.redis.Pipeline()
	pipe.Incr(ctx, windowKey)
	pipe.Expire(ctx, windowKey, duration)
	results, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check IP %s window: %w", window, err)
	}

	count := results[0].(*redis.IntCmd).Val()

	if int(count) > limit {
		return &RateLimitResult{
			Allowed:    false,
			Remaining:  0,
			Reset:      now.Add(duration),
			RetryAfter: duration,
			LimitType:  fmt.Sprintf("ip_%s_limit", window),
		}, nil
	}

	return &RateLimitResult{
		Allowed:   true,
		Remaining: limit - int(count),
		Reset:     now.Add(duration),
		LimitType: fmt.Sprintf("ip_%s_limit", window),
	}, nil
}

// cleanIP extracts and cleans IP address
func (i *IPRateLimiter) cleanIP(clientIP string) string {
	// Handle X-Forwarded-For header
	if strings.Contains(clientIP, ",") {
		ips := strings.Split(clientIP, ",")
		clientIP = strings.TrimSpace(ips[0])
	}

	// Parse IP to ensure validity
	ip := net.ParseIP(clientIP)
	if ip == nil {
		return ""
	}

	return ip.String()
}

// isWhitelisted checks if IP is in whitelist
func (i *IPRateLimiter) isWhitelisted(ip string) bool {
	return i.isInIPList(ip, i.config.WhitelistedIPs)
}

// isBlacklisted checks if IP is in blacklist
func (i *IPRateLimiter) isBlacklisted(ip string) bool {
	return i.isInIPList(ip, i.config.BlacklistedIPs)
}

// isInIPList checks if IP is in the given list (supports CIDR)
func (i *IPRateLimiter) isInIPList(checkIP string, ipList []string) bool {
	ip := net.ParseIP(checkIP)
	if ip == nil {
		return false
	}

	for _, entry := range ipList {
		// Check if it's a CIDR range
		if strings.Contains(entry, "/") {
			_, cidr, err := net.ParseCIDR(entry)
			if err != nil {
				continue
			}
			if cidr.Contains(ip) {
				return true
			}
		} else {
			// Direct IP comparison
			if entry == checkIP {
				return true
			}
		}
	}

	return false
}

// BlockIP temporarily blocks an IP address
func (i *IPRateLimiter) BlockIP(ctx context.Context, ip string, duration time.Duration, reason string) error {
	blockKey := fmt.Sprintf("blocked:ip:%s", ip)
	err := i.redis.Set(ctx, blockKey, reason, duration).Err()
	if err != nil {
		return fmt.Errorf("failed to block IP %s: %w", ip, err)
	}

	// Log the block
	logKey := fmt.Sprintf("block_log:ip:%s:%d", ip, time.Now().Unix())
	i.redis.HSet(ctx, logKey, map[string]interface{}{
		"ip":        ip,
		"reason":    reason,
		"timestamp": time.Now().Unix(),
		"duration":  int(duration.Seconds()),
	})
	i.redis.Expire(ctx, logKey, 24*time.Hour) // Keep logs for 24 hours

	return nil
}

// UnblockIP removes IP from blocked list
func (i *IPRateLimiter) UnblockIP(ctx context.Context, ip string) error {
	blockKey := fmt.Sprintf("blocked:ip:%s", ip)
	err := i.redis.Del(ctx, blockKey).Err()
	if err != nil {
		return fmt.Errorf("failed to unblock IP %s: %w", ip, err)
	}
	return nil
}

// IsIPBlocked checks if IP is temporarily blocked
func (i *IPRateLimiter) IsIPBlocked(ctx context.Context, ip string) (bool, string, error) {
	blockKey := fmt.Sprintf("blocked:ip:%s", ip)
	reason, err := i.redis.Get(ctx, blockKey).Result()
	if err == redis.Nil {
		return false, "", nil
	}
	if err != nil {
		return false, "", fmt.Errorf("failed to check if IP %s is blocked: %w", ip, err)
	}
	return true, reason, nil
}

// GetIPStats retrieves current IP statistics
func (i *IPRateLimiter) GetIPStats(ctx context.Context, ip string) (*IPStats, error) {
	now := time.Now()

	// Get current minute and hour counts
	minuteKey := fmt.Sprintf("ip:%s:minute:%d", ip, now.Unix()/60)
	hourKey := fmt.Sprintf("ip:%s:hour:%d", ip, now.Unix()/3600)

	pipe := i.redis.Pipeline()
	minuteCmd := pipe.Get(ctx, minuteKey)
	hourCmd := pipe.Get(ctx, hourKey)
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to get IP stats: %w", err)
	}

	minuteCount := 0
	hourCount := 0

	if minuteVal, err := minuteCmd.Result(); err == nil {
		if val, err := redis.ParseInt(minuteVal); err == nil {
			minuteCount = int(val)
		}
	}

	if hourVal, err := hourCmd.Result(); err == nil {
		if val, err := redis.ParseInt(hourVal); err == nil {
			hourCount = int(val)
		}
	}

	// Check if blocked
	blocked, reason, err := i.IsIPBlocked(ctx, ip)
	if err != nil {
		return nil, err
	}

	return &IPStats{
		IP:                  ip,
		RequestsThisMinute:  minuteCount,
		RequestsThisHour:    hourCount,
		MinuteLimit:         i.config.RequestsPerMinute,
		HourLimit:           i.config.RequestsPerHour,
		IsBlocked:           blocked,
		BlockReason:         reason,
		IsWhitelisted:       i.isWhitelisted(ip),
		IsBlacklisted:       i.isBlacklisted(ip),
		RemainingThisMinute: i.config.RequestsPerMinute - minuteCount,
		RemainingThisHour:   i.config.RequestsPerHour - hourCount,
	}, nil
}

// IPStats contains IP address statistics
type IPStats struct {
	IP                  string `json:"ip"`
	RequestsThisMinute  int    `json:"requests_this_minute"`
	RequestsThisHour    int    `json:"requests_this_hour"`
	MinuteLimit         int    `json:"minute_limit"`
	HourLimit           int    `json:"hour_limit"`
	IsBlocked           bool   `json:"is_blocked"`
	BlockReason         string `json:"block_reason,omitempty"`
	IsWhitelisted       bool   `json:"is_whitelisted"`
	IsBlacklisted       bool   `json:"is_blacklisted"`
	RemainingThisMinute int    `json:"remaining_this_minute"`
	RemainingThisHour   int    `json:"remaining_this_hour"`
}

// GetBlockedIPs returns list of currently blocked IPs
func (i *IPRateLimiter) GetBlockedIPs(ctx context.Context) ([]BlockedIPInfo, error) {
	pattern := "blocked:ip:*"
	keys, err := i.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get blocked IPs: %w", err)
	}

	var blockedIPs []BlockedIPInfo
	for _, key := range keys {
		// Extract IP from key
		parts := strings.Split(key, ":")
		if len(parts) < 3 {
			continue
		}
		ip := parts[2]

		// Get block reason and TTL
		pipe := i.redis.Pipeline()
		reasonCmd := pipe.Get(ctx, key)
		ttlCmd := pipe.TTL(ctx, key)
		_, err := pipe.Exec(ctx)
		if err != nil {
			continue
		}

		reason, _ := reasonCmd.Result()
		ttl, _ := ttlCmd.Result()

		blockedIPs = append(blockedIPs, BlockedIPInfo{
			IP:        ip,
			Reason:    reason,
			ExpiresIn: ttl,
			BlockedAt: time.Now().Add(-ttl),
		})
	}

	return blockedIPs, nil
}

// BlockedIPInfo contains information about blocked IP
type BlockedIPInfo struct {
	IP        string        `json:"ip"`
	Reason    string        `json:"reason"`
	ExpiresIn time.Duration `json:"expires_in"`
	BlockedAt time.Time     `json:"blocked_at"`
}

// ClearIPCounters clears all counters for an IP
func (i *IPRateLimiter) ClearIPCounters(ctx context.Context, ip string) error {
	now := time.Now()

	// Clear current windows
	keys := []string{
		fmt.Sprintf("ip:%s:minute:%d", ip, now.Unix()/60),
		fmt.Sprintf("ip:%s:hour:%d", ip, now.Unix()/3600),
	}

	// Also clear previous windows
	keys = append(keys,
		fmt.Sprintf("ip:%s:minute:%d", ip, (now.Unix()/60)-1),
		fmt.Sprintf("ip:%s:hour:%d", ip, (now.Unix()/3600)-1),
	)

	err := i.redis.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("failed to clear IP counters: %w", err)
	}

	return nil
}

// UpdateIPLimits dynamically updates IP rate limits
func (i *IPRateLimiter) UpdateIPLimits(ctx context.Context, config *IPRateLimitConfig) error {
	// Validate new configuration
	if config.RequestsPerMinute <= 0 || config.RequestsPerHour <= 0 {
		return fmt.Errorf("invalid rate limits: minute=%d, hour=%d",
			config.RequestsPerMinute, config.RequestsPerHour)
	}

	// Update configuration
	i.config = config

	// Store configuration in Redis for persistence
	configKey := "ip_rate_limit_config"
	configData := map[string]interface{}{
		"requests_per_minute": config.RequestsPerMinute,
		"requests_per_hour":   config.RequestsPerHour,
		"updated_at":          time.Now().Unix(),
	}

	err := i.redis.HSet(ctx, configKey, configData).Err()
	if err != nil {
		return fmt.Errorf("failed to store IP rate limit config: %w", err)
	}

	return nil
}
