package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// Herald.lol Gaming Analytics - Gaming-Specific Rate Limiter
// Advanced rate limiting for gaming platform with tier-based limits

// GamingRateLimiter provides gaming-specific rate limiting functionality
type GamingRateLimiter struct {
	redis        *redis.Client
	config       *GamingRateLimitConfig
	ipLimiter    *IPRateLimiter
	apiLimiter   *APIRateLimiter
	gameRateLimiter *GameSpecificRateLimiter
}

// GamingRateLimitConfig contains rate limiting configuration
type GamingRateLimitConfig struct {
	// Subscription tier limits
	FreeTierLimits     *TierLimits `json:"free_tier"`
	PremiumTierLimits  *TierLimits `json:"premium_tier"`
	ProTierLimits      *TierLimits `json:"pro_tier"`
	EnterpriseLimits   *TierLimits `json:"enterprise_tier"`
	
	// Gaming-specific limits
	RiotAPILimits      *GameAPILimits `json:"riot_api_limits"`
	AnalyticsLimits    *AnalyticsRateLimits `json:"analytics_limits"`
	ExportLimits       *ExportRateLimits `json:"export_limits"`
	
	// DDoS protection
	DDoSProtection     *DDoSProtectionConfig `json:"ddos_protection"`
	
	// IP-based limits
	IPLimits          *IPRateLimitConfig `json:"ip_limits"`
}

// TierLimits defines rate limits per subscription tier
type TierLimits struct {
	RequestsPerMinute    int           `json:"requests_per_minute"`
	RequestsPerHour      int           `json:"requests_per_hour"`
	RequestsPerDay       int           `json:"requests_per_day"`
	AnalyticsPerMinute   int           `json:"analytics_per_minute"`
	ExportsPerDay        int           `json:"exports_per_day"`
	BurstLimit           int           `json:"burst_limit"`
	WindowDuration       time.Duration `json:"window_duration"`
}

// GameAPILimits defines gaming API specific limits
type GameAPILimits struct {
	RiotPersonalLimit    int `json:"riot_personal_limit"`    // 100 req/2min
	RiotProductionLimit  int `json:"riot_production_limit"`  // Variable
	MatchDataPerHour     int `json:"match_data_per_hour"`
	SummonerDataPerHour  int `json:"summoner_data_per_hour"`
	RankedDataPerHour    int `json:"ranked_data_per_hour"`
}

// AnalyticsRateLimits defines analytics-specific rate limits
type AnalyticsRateLimits struct {
	BasicAnalyticsPerMin    int `json:"basic_analytics_per_min"`
	AdvancedAnalyticsPerMin int `json:"advanced_analytics_per_min"`
	RealTimeAnalyticsPerMin int `json:"realtime_analytics_per_min"`
	TeamAnalyticsPerMin     int `json:"team_analytics_per_min"`
	ComparisonAnalyticsPerMin int `json:"comparison_analytics_per_min"`
}

// ExportRateLimits defines data export rate limits
type ExportRateLimits struct {
	JSONExportsPerDay  int `json:"json_exports_per_day"`
	CSVExportsPerDay   int `json:"csv_exports_per_day"`
	PDFExportsPerDay   int `json:"pdf_exports_per_day"`
	ExcelExportsPerDay int `json:"excel_exports_per_day"`
	MaxExportSizeMB    int `json:"max_export_size_mb"`
}

// DDoSProtectionConfig defines DDoS protection settings
type DDoSProtectionConfig struct {
	Enabled              bool          `json:"enabled"`
	RequestThreshold     int           `json:"request_threshold"`
	WindowDuration       time.Duration `json:"window_duration"`
	BlockDuration        time.Duration `json:"block_duration"`
	SuspiciousPatterns   []string      `json:"suspicious_patterns"`
	GeoBlocking          bool          `json:"geo_blocking"`
	BlockedCountries     []string      `json:"blocked_countries"`
}

// IPRateLimitConfig defines IP-based rate limiting
type IPRateLimitConfig struct {
	RequestsPerMinute  int           `json:"requests_per_minute"`
	RequestsPerHour    int           `json:"requests_per_hour"`
	WhitelistedIPs     []string      `json:"whitelisted_ips"`
	BlacklistedIPs     []string      `json:"blacklisted_ips"`
	TrustedProxies     []string      `json:"trusted_proxies"`
}

// RateLimitResult contains rate limiting decision
type RateLimitResult struct {
	Allowed           bool          `json:"allowed"`
	Remaining         int           `json:"remaining"`
	Reset             time.Time     `json:"reset"`
	RetryAfter        time.Duration `json:"retry_after,omitempty"`
	Tier              string        `json:"tier"`
	LimitType         string        `json:"limit_type"`
	RequestsThisMinute int          `json:"requests_this_minute"`
	RequestsThisHour  int           `json:"requests_this_hour"`
	RequestsThisDay   int           `json:"requests_this_day"`
}

// NewGamingRateLimiter creates new gaming rate limiter
func NewGamingRateLimiter(redis *redis.Client, config *GamingRateLimitConfig) *GamingRateLimiter {
	return &GamingRateLimiter{
		redis:           redis,
		config:          config,
		ipLimiter:       NewIPRateLimiter(redis, config.IPLimits),
		apiLimiter:      NewAPIRateLimiter(redis),
		gameRateLimiter: NewGameSpecificRateLimiter(redis, config),
	}
}

// CheckGamingRateLimit performs comprehensive gaming rate limit check
func (g *GamingRateLimiter) CheckGamingRateLimit(ctx context.Context, request *GamingRateLimitRequest) (*RateLimitResult, error) {
	// 1. Check IP-based limits first (DDoS protection)
	if g.config.DDoSProtection.Enabled {
		if blocked, err := g.checkDDoSProtection(ctx, request); err != nil {
			return nil, err
		} else if blocked {
			return &RateLimitResult{
				Allowed:    false,
				LimitType:  "ddos_protection",
				RetryAfter: g.config.DDoSProtection.BlockDuration,
			}, nil
		}
	}
	
	// 2. Check IP rate limits
	ipResult, err := g.ipLimiter.CheckIPLimit(ctx, request.ClientIP)
	if err != nil {
		return nil, fmt.Errorf("failed to check IP limit: %w", err)
	}
	if !ipResult.Allowed {
		return ipResult, nil
	}
	
	// 3. Check subscription tier limits
	tierResult, err := g.checkTierLimits(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to check tier limits: %w", err)
	}
	if !tierResult.Allowed {
		return tierResult, nil
	}
	
	// 4. Check gaming-specific limits
	gamingResult, err := g.gameRateLimiter.CheckGameLimits(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to check gaming limits: %w", err)
	}
	if !gamingResult.Allowed {
		return gamingResult, nil
	}
	
	// 5. Update counters if all checks pass
	if err := g.updateCounters(ctx, request); err != nil {
		return nil, fmt.Errorf("failed to update counters: %w", err)
	}
	
	return &RateLimitResult{
		Allowed:   true,
		Remaining: tierResult.Remaining,
		Reset:     tierResult.Reset,
		Tier:      request.SubscriptionTier,
		LimitType: "tier_limit",
	}, nil
}

// checkDDoSProtection checks for DDoS attack patterns
func (g *GamingRateLimiter) checkDDoSProtection(ctx context.Context, request *GamingRateLimitRequest) (bool, error) {
	key := fmt.Sprintf("ddos:ip:%s", request.ClientIP)
	
	// Check request rate within window
	pipe := g.redis.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, g.config.DDoSProtection.WindowDuration)
	results, err := pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check DDoS protection: %w", err)
	}
	
	count := results[0].(*redis.IntCmd).Val()
	
	// Block if threshold exceeded
	if int(count) > g.config.DDoSProtection.RequestThreshold {
		// Add to blocked IPs
		blockKey := fmt.Sprintf("blocked:ip:%s", request.ClientIP)
		g.redis.Set(ctx, blockKey, "ddos_blocked", g.config.DDoSProtection.BlockDuration)
		
		return true, nil
	}
	
	// Check if IP is already blocked
	blockKey := fmt.Sprintf("blocked:ip:%s", request.ClientIP)
	blocked, err := g.redis.Exists(ctx, blockKey).Result()
	if err != nil {
		return false, err
	}
	
	return blocked > 0, nil
}

// checkTierLimits checks subscription tier-based limits
func (g *GamingRateLimiter) checkTierLimits(ctx context.Context, request *GamingRateLimitRequest) (*RateLimitResult, error) {
	var limits *TierLimits
	
	// Get limits based on subscription tier
	switch request.SubscriptionTier {
	case "free":
		limits = g.config.FreeTierLimits
	case "premium":
		limits = g.config.PremiumTierLimits
	case "pro":
		limits = g.config.ProTierLimits
	case "enterprise":
		limits = g.config.EnterpriseLimits
	default:
		limits = g.config.FreeTierLimits // Default to free tier
	}
	
	// Check multiple time windows
	now := time.Now()
	
	// Minute window
	minuteKey := fmt.Sprintf("tier:%s:user:%s:minute:%d", request.SubscriptionTier, request.UserID, now.Unix()/60)
	minuteCount, err := g.redis.Incr(ctx, minuteKey).Result()
	if err != nil {
		return nil, err
	}
	g.redis.Expire(ctx, minuteKey, time.Minute)
	
	if int(minuteCount) > limits.RequestsPerMinute {
		return &RateLimitResult{
			Allowed:    false,
			Remaining:  0,
			Reset:      now.Add(time.Minute),
			RetryAfter: time.Minute,
			Tier:       request.SubscriptionTier,
			LimitType:  "tier_minute_limit",
		}, nil
	}
	
	// Hour window
	hourKey := fmt.Sprintf("tier:%s:user:%s:hour:%d", request.SubscriptionTier, request.UserID, now.Unix()/3600)
	hourCount, err := g.redis.Incr(ctx, hourKey).Result()
	if err != nil {
		return nil, err
	}
	g.redis.Expire(ctx, hourKey, time.Hour)
	
	if int(hourCount) > limits.RequestsPerHour {
		return &RateLimitResult{
			Allowed:    false,
			Remaining:  0,
			Reset:      now.Add(time.Hour),
			RetryAfter: time.Hour,
			Tier:       request.SubscriptionTier,
			LimitType:  "tier_hour_limit",
		}, nil
	}
	
	// Day window
	dayKey := fmt.Sprintf("tier:%s:user:%s:day:%d", request.SubscriptionTier, request.UserID, now.Unix()/86400)
	dayCount, err := g.redis.Incr(ctx, dayKey).Result()
	if err != nil {
		return nil, err
	}
	g.redis.Expire(ctx, dayKey, 24*time.Hour)
	
	if int(dayCount) > limits.RequestsPerDay {
		return &RateLimitResult{
			Allowed:    false,
			Remaining:  0,
			Reset:      now.Add(24 * time.Hour),
			RetryAfter: 24 * time.Hour,
			Tier:       request.SubscriptionTier,
			LimitType:  "tier_day_limit",
		}, nil
	}
	
	return &RateLimitResult{
		Allowed:            true,
		Remaining:          limits.RequestsPerMinute - int(minuteCount),
		Reset:              now.Add(time.Minute),
		Tier:               request.SubscriptionTier,
		LimitType:          "tier_limit",
		RequestsThisMinute: int(minuteCount),
		RequestsThisHour:   int(hourCount),
		RequestsThisDay:    int(dayCount),
	}, nil
}

// updateCounters updates rate limit counters
func (g *GamingRateLimiter) updateCounters(ctx context.Context, request *GamingRateLimitRequest) error {
	pipe := g.redis.Pipeline()
	now := time.Now()
	
	// Update user counters
	userMinuteKey := fmt.Sprintf("user:%s:minute:%d", request.UserID, now.Unix()/60)
	userHourKey := fmt.Sprintf("user:%s:hour:%d", request.UserID, now.Unix()/3600)
	userDayKey := fmt.Sprintf("user:%s:day:%d", request.UserID, now.Unix()/86400)
	
	pipe.Incr(ctx, userMinuteKey)
	pipe.Expire(ctx, userMinuteKey, time.Minute)
	pipe.Incr(ctx, userHourKey)
	pipe.Expire(ctx, userHourKey, time.Hour)
	pipe.Incr(ctx, userDayKey)
	pipe.Expire(ctx, userDayKey, 24*time.Hour)
	
	// Update endpoint-specific counters
	if request.Endpoint != "" {
		endpointKey := fmt.Sprintf("endpoint:%s:minute:%d", request.Endpoint, now.Unix()/60)
		pipe.Incr(ctx, endpointKey)
		pipe.Expire(ctx, endpointKey, time.Minute)
	}
	
	_, err := pipe.Exec(ctx)
	return err
}

// GamingRateLimitRequest contains request information for rate limiting
type GamingRateLimitRequest struct {
	UserID           string `json:"user_id"`
	ClientIP         string `json:"client_ip"`
	SubscriptionTier string `json:"subscription_tier"`
	Endpoint         string `json:"endpoint"`
	Method           string `json:"method"`
	APIKey           string `json:"api_key"`
	UserAgent        string `json:"user_agent"`
	Region           string `json:"region,omitempty"`
	GameType         string `json:"game_type,omitempty"` // LoL, TFT, etc.
}

// GamingRateLimitMiddleware creates Gin middleware for gaming rate limiting
func (g *GamingRateLimiter) GamingRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract request information
		request := &GamingRateLimitRequest{
			UserID:           getUserID(c),
			ClientIP:         c.ClientIP(),
			SubscriptionTier: getSubscriptionTier(c),
			Endpoint:         c.FullPath(),
			Method:           c.Request.Method,
			APIKey:           c.GetHeader("X-Gaming-API-Key"),
			UserAgent:        c.GetHeader("User-Agent"),
			Region:           c.Param("region"),
			GameType:         c.Query("game_type"),
		}
		
		// Check rate limits
		result, err := g.CheckGamingRateLimit(c.Request.Context(), request)
		if err != nil {
			c.JSON(500, gin.H{
				"error":           "Gaming rate limit check failed",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Add rate limit headers
		c.Header("X-Gaming-Rate-Limit", strconv.Itoa(getTierLimit(request.SubscriptionTier)))
		c.Header("X-Gaming-Rate-Remaining", strconv.Itoa(result.Remaining))
		c.Header("X-Gaming-Rate-Reset", strconv.FormatInt(result.Reset.Unix(), 10))
		c.Header("X-Gaming-Rate-Tier", result.Tier)
		
		if !result.Allowed {
			c.Header("Retry-After", strconv.FormatInt(int64(result.RetryAfter.Seconds()), 10))
			
			c.JSON(429, gin.H{
				"error":               "Gaming rate limit exceeded",
				"limit_type":          result.LimitType,
				"retry_after_seconds": int(result.RetryAfter.Seconds()),
				"tier":                result.Tier,
				"gaming_platform":     "herald-lol",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// Helper functions

func getUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return "anonymous"
}

func getSubscriptionTier(c *gin.Context) string {
	if tier, exists := c.Get("subscription_tier"); exists {
		if t, ok := tier.(string); ok {
			return t
		}
	}
	return "free"
}

func getTierLimit(tier string) int {
	limits := map[string]int{
		"free":       60,
		"premium":    300,
		"pro":        1200,
		"enterprise": 6000,
	}
	
	if limit, exists := limits[tier]; exists {
		return limit
	}
	return limits["free"]
}

// DefaultGamingRateLimitConfig returns default configuration
func DefaultGamingRateLimitConfig() *GamingRateLimitConfig {
	return &GamingRateLimitConfig{
		FreeTierLimits: &TierLimits{
			RequestsPerMinute:  60,
			RequestsPerHour:    3600,
			RequestsPerDay:     50000,
			AnalyticsPerMinute: 30,
			ExportsPerDay:      1,
			BurstLimit:         10,
			WindowDuration:     time.Minute,
		},
		PremiumTierLimits: &TierLimits{
			RequestsPerMinute:  300,
			RequestsPerHour:    18000,
			RequestsPerDay:     250000,
			AnalyticsPerMinute: 100,
			ExportsPerDay:      5,
			BurstLimit:         50,
			WindowDuration:     time.Minute,
		},
		ProTierLimits: &TierLimits{
			RequestsPerMinute:  1200,
			RequestsPerHour:    72000,
			RequestsPerDay:     1000000,
			AnalyticsPerMinute: 500,
			ExportsPerDay:      25,
			BurstLimit:         200,
			WindowDuration:     time.Minute,
		},
		EnterpriseLimits: &TierLimits{
			RequestsPerMinute:  6000,
			RequestsPerHour:    360000,
			RequestsPerDay:     5000000,
			AnalyticsPerMinute: 2000,
			ExportsPerDay:      -1, // Unlimited
			BurstLimit:         1000,
			WindowDuration:     time.Minute,
		},
		RiotAPILimits: &GameAPILimits{
			RiotPersonalLimit:   100,
			RiotProductionLimit: 3000,
			MatchDataPerHour:    1000,
			SummonerDataPerHour: 2000,
			RankedDataPerHour:   1500,
		},
		AnalyticsLimits: &AnalyticsRateLimits{
			BasicAnalyticsPerMin:      30,
			AdvancedAnalyticsPerMin:   100,
			RealTimeAnalyticsPerMin:   50,
			TeamAnalyticsPerMin:       20,
			ComparisonAnalyticsPerMin: 10,
		},
		ExportLimits: &ExportRateLimits{
			JSONExportsPerDay:  5,
			CSVExportsPerDay:   3,
			PDFExportsPerDay:   2,
			ExcelExportsPerDay: 1,
			MaxExportSizeMB:    100,
		},
		DDoSProtection: &DDoSProtectionConfig{
			Enabled:          true,
			RequestThreshold: 1000,
			WindowDuration:   time.Minute,
			BlockDuration:    time.Hour,
			SuspiciousPatterns: []string{
				"bot", "crawler", "scanner", "attack",
			},
			GeoBlocking:      false,
			BlockedCountries: []string{},
		},
		IPLimits: &IPRateLimitConfig{
			RequestsPerMinute: 300,
			RequestsPerHour:   18000,
			WhitelistedIPs:    []string{},
			BlacklistedIPs:    []string{},
			TrustedProxies: []string{
				"10.0.0.0/8",
				"172.16.0.0/12",
				"192.168.0.0/16",
			},
		},
	}
}