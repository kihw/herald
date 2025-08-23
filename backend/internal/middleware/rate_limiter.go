package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// Herald.lol Gaming Analytics - Advanced Rate Limiting & DDoS Protection
// Comprehensive rate limiting system for gaming platform APIs

// GamingRateLimiter provides gaming-specific rate limiting
type GamingRateLimiter struct {
	redisClient    *redis.Client
	config         *GamingRateLimitConfig
	ddosProtector  *DDoSProtector
	adaptiveLimits *AdaptiveLimitManager
}

// GamingRateLimitConfig holds rate limiting configuration for gaming platform
type GamingRateLimitConfig struct {
	// Global limits
	GlobalRPM         int           `json:"global_rpm"`         // Requests per minute globally
	GlobalBurst       int           `json:"global_burst"`       // Burst capacity
	GlobalWindow      time.Duration `json:"global_window"`      // Time window
	
	// Per-user limits by subscription tier
	FreeTierRPM       int           `json:"free_tier_rpm"`
	PremiumTierRPM    int           `json:"premium_tier_rpm"`
	ProTierRPM        int           `json:"pro_tier_rpm"`
	EnterpriseTierRPM int           `json:"enterprise_tier_rpm"`
	
	// Gaming-specific limits
	AnalyticsRPM      int           `json:"analytics_rpm"`      // Analytics API calls
	RiotAPIProxyRPM   int           `json:"riot_api_proxy_rpm"` // Riot API proxy calls
	ExportRPM         int           `json:"export_rpm"`         // Data export limits
	
	// Geographic limits
	RegionLimits      map[string]int `json:"region_limits"`     // Per-region limits
	
	// API endpoint specific limits
	EndpointLimits    map[string]int `json:"endpoint_limits"`   // Per-endpoint limits
	
	// DDoS protection settings
	EnableDDoSProtection bool          `json:"enable_ddos_protection"`
	SuspiciousThreshold  int           `json:"suspicious_threshold"`
	BlockDuration       time.Duration  `json:"block_duration"`
	
	// Adaptive limiting
	EnableAdaptive      bool          `json:"enable_adaptive"`
	LoadThreshold       float64       `json:"load_threshold"`
	AdaptiveMultiplier  float64       `json:"adaptive_multiplier"`
}

// RateLimitResult represents the result of a rate limit check
type RateLimitResult struct {
	Allowed       bool          `json:"allowed"`
	Limit         int           `json:"limit"`
	Remaining     int           `json:"remaining"`
	ResetTime     time.Time     `json:"reset_time"`
	RetryAfter    time.Duration `json:"retry_after"`
	Tier          string        `json:"tier"`
	Reason        string        `json:"reason,omitempty"`
	BlockedUntil  *time.Time    `json:"blocked_until,omitempty"`
}

// NewGamingRateLimiter creates new gaming rate limiter
func NewGamingRateLimiter(redisClient *redis.Client, config *GamingRateLimitConfig) *GamingRateLimiter {
	// Set gaming-specific defaults
	if config.GlobalRPM == 0 {
		config.GlobalRPM = 100000 // 100k RPM global limit
	}
	if config.GlobalBurst == 0 {
		config.GlobalBurst = 1000
	}
	if config.GlobalWindow == 0 {
		config.GlobalWindow = time.Minute
	}
	
	// Gaming subscription tier defaults
	if config.FreeTierRPM == 0 {
		config.FreeTierRPM = 60 // 1 per second for free users
	}
	if config.PremiumTierRPM == 0 {
		config.PremiumTierRPM = 300 // 5 per second for premium
	}
	if config.ProTierRPM == 0 {
		config.ProTierRPM = 1200 // 20 per second for pro
	}
	if config.EnterpriseTierRPM == 0 {
		config.EnterpriseTierRPM = 6000 // 100 per second for enterprise
	}
	
	// Gaming-specific API limits
	if config.AnalyticsRPM == 0 {
		config.AnalyticsRPM = 180 // 3 per second for analytics
	}
	if config.RiotAPIProxyRPM == 0 {
		config.RiotAPIProxyRPM = 100 // Respect Riot API limits
	}
	if config.ExportRPM == 0 {
		config.ExportRPM = 10 // 10 exports per minute
	}
	
	// DDoS protection defaults
	if config.SuspiciousThreshold == 0 {
		config.SuspiciousThreshold = 1000 // 1000 requests in window
	}
	if config.BlockDuration == 0 {
		config.BlockDuration = 15 * time.Minute
	}
	
	// Initialize region limits if not set
	if config.RegionLimits == nil {
		config.RegionLimits = map[string]int{
			"NA": 30000,  // North America
			"EUW": 25000, // Europe West
			"EUNE": 15000, // Europe Northeast
			"KR": 20000,  // Korea
			"JP": 10000,  // Japan
			"BR": 8000,   // Brazil
			"LAN": 5000,  // Latin America North
			"LAS": 5000,  // Latin America South
			"OCE": 3000,  // Oceania
			"RU": 7000,   // Russia
			"TR": 6000,   // Turkey
		}
	}
	
	// Initialize endpoint-specific limits
	if config.EndpointLimits == nil {
		config.EndpointLimits = map[string]int{
			"/api/gaming/analytics/export":     10,   // Export limits
			"/api/gaming/analytics/bulk":       30,   // Bulk analytics
			"/api/gaming/riot/summoner":        120,  // Summoner lookups
			"/api/gaming/riot/match":           180,  // Match data
			"/api/gaming/riot/league":          60,   // League data
			"/api/gaming/teams/create":         5,    // Team creation
			"/api/gaming/subscription/change":  3,    // Subscription changes
			"/api/gaming/admin/users":          60,   // Admin user operations
		}
	}
	
	limiter := &GamingRateLimiter{
		redisClient: redisClient,
		config:      config,
		ddosProtector: NewDDoSProtector(redisClient, config),
		adaptiveLimits: NewAdaptiveLimitManager(redisClient, config),
	}
	
	return limiter
}

// RateLimit middleware for gaming platform
func (rl *GamingRateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract client identifier
		clientID := rl.getClientIdentifier(c)
		
		// Check if client is blocked by DDoS protection
		if rl.config.EnableDDoSProtection {
			if blocked, blockInfo := rl.ddosProtector.IsBlocked(c.Request.Context(), clientID); blocked {
				c.Header("X-Gaming-Rate-Limit-Blocked", "true")
				c.Header("X-Gaming-Block-Reason", blockInfo.Reason)
				if blockInfo.BlockedUntil != nil {
					c.Header("Retry-After", strconv.Itoa(int(time.Until(*blockInfo.BlockedUntil).Seconds())))
				}
				
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Gaming API access temporarily blocked",
					"reason": blockInfo.Reason,
					"blocked_until": blockInfo.BlockedUntil,
					"gaming_platform": "herald-lol",
					"support_contact": "support@herald.lol",
				})
				c.Abort()
				return
			}
		}
		
		// Perform rate limit check
		result := rl.checkRateLimit(c.Request.Context(), c, clientID)
		
		// Set rate limit headers
		rl.setRateLimitHeaders(c, result)
		
		if !result.Allowed {
			// Log rate limit hit for analytics
			rl.logRateLimitHit(c, clientID, result)
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Gaming API rate limit exceeded",
				"limit": result.Limit,
				"remaining": result.Remaining,
				"reset_time": result.ResetTime,
				"retry_after": result.RetryAfter.Seconds(),
				"tier": result.Tier,
				"reason": result.Reason,
				"gaming_platform": "herald-lol",
				"upgrade_info": rl.getUpgradeInfo(result.Tier),
			})
			c.Abort()
			return
		}
		
		// Continue to next middleware
		c.Next()
	}
}

// GamingAnalyticsRateLimit specialized rate limiting for analytics endpoints
func (rl *GamingRateLimiter) GamingAnalyticsRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := rl.getClientIdentifier(c)
		
		// Check analytics-specific rate limits
		result := rl.checkAnalyticsRateLimit(c.Request.Context(), c, clientID)
		rl.setRateLimitHeaders(c, result)
		
		if !result.Allowed {
			c.Header("X-Gaming-Analytics-Limit-Exceeded", "true")
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Gaming analytics rate limit exceeded",
				"limit": result.Limit,
				"remaining": result.Remaining,
				"reset_time": result.ResetTime,
				"analytics_limits": gin.H{
					"basic": rl.config.AnalyticsRPM,
					"premium": rl.config.AnalyticsRPM * 2,
					"pro": rl.config.AnalyticsRPM * 5,
				},
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// RiotAPIProxyRateLimit rate limiting for Riot API proxy
func (rl *GamingRateLimiter) RiotAPIProxyRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := rl.getClientIdentifier(c)
		
		// Check Riot API proxy limits
		result := rl.checkRiotAPIRateLimit(c.Request.Context(), c, clientID)
		rl.setRateLimitHeaders(c, result)
		
		if !result.Allowed {
			c.Header("X-Gaming-Riot-API-Limit-Exceeded", "true")
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Riot API proxy rate limit exceeded",
				"limit": result.Limit,
				"remaining": result.Remaining,
				"reset_time": result.ResetTime,
				"riot_api_info": gin.H{
					"note": "Rate limits protect against Riot Games API violations",
					"upgrade": "Consider upgrading for higher limits",
				},
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// PerRegionRateLimit region-based rate limiting
func (rl *GamingRateLimiter) PerRegionRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		region := rl.extractRegion(c)
		if region == "" {
			c.Next()
			return
		}
		
		regionLimit, exists := rl.config.RegionLimits[region]
		if !exists {
			c.Next()
			return
		}
		
		// Check region-specific limits
		key := fmt.Sprintf("herald:rate_limit:region:%s", region)
		result := rl.checkLimit(c.Request.Context(), key, regionLimit, rl.config.GlobalWindow)
		
		c.Header("X-Gaming-Region", region)
		c.Header("X-Gaming-Region-Limit", strconv.Itoa(regionLimit))
		c.Header("X-Gaming-Region-Remaining", strconv.Itoa(result.Remaining))
		
		if !result.Allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Regional gaming API rate limit exceeded",
				"region": region,
				"limit": regionLimit,
				"remaining": result.Remaining,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// checkRateLimit performs comprehensive rate limit checking
func (rl *GamingRateLimiter) checkRateLimit(ctx context.Context, c *gin.Context, clientID string) *RateLimitResult {
	// Determine user's subscription tier
	tier := rl.getUserTier(c)
	limit := rl.getTierLimit(tier)
	
	// Get endpoint-specific limit if applicable
	endpoint := c.Request.URL.Path
	if endpointLimit, exists := rl.config.EndpointLimits[endpoint]; exists {
		if endpointLimit < limit {
			limit = endpointLimit // Use the more restrictive limit
		}
	}
	
	// Apply adaptive limiting if enabled
	if rl.config.EnableAdaptive {
		adaptiveLimit := rl.adaptiveLimits.GetAdaptiveLimit(ctx, limit)
		if adaptiveLimit < limit {
			limit = adaptiveLimit
		}
	}
	
	// Check the limit
	key := fmt.Sprintf("herald:rate_limit:user:%s", clientID)
	result := rl.checkLimit(ctx, key, limit, rl.config.GlobalWindow)
	result.Tier = tier
	
	// Record request for DDoS detection
	if rl.config.EnableDDoSProtection {
		rl.ddosProtector.RecordRequest(ctx, clientID, c.ClientIP())
	}
	
	return result
}

// checkAnalyticsRateLimit checks analytics-specific rate limits
func (rl *GamingRateLimiter) checkAnalyticsRateLimit(ctx context.Context, c *gin.Context, clientID string) *RateLimitResult {
	tier := rl.getUserTier(c)
	
	// Analytics limits based on tier
	var limit int
	switch tier {
	case "premium":
		limit = rl.config.AnalyticsRPM * 2
	case "pro":
		limit = rl.config.AnalyticsRPM * 5
	case "enterprise":
		limit = rl.config.AnalyticsRPM * 10
	default:
		limit = rl.config.AnalyticsRPM
	}
	
	key := fmt.Sprintf("herald:rate_limit:analytics:%s", clientID)
	result := rl.checkLimit(ctx, key, limit, rl.config.GlobalWindow)
	result.Tier = tier
	result.Reason = "analytics"
	
	return result
}

// checkRiotAPIRateLimit checks Riot API proxy rate limits
func (rl *GamingRateLimiter) checkRiotAPIRateLimit(ctx context.Context, c *gin.Context, clientID string) *RateLimitResult {
	tier := rl.getUserTier(c)
	
	// Riot API limits are more conservative
	var limit int
	switch tier {
	case "premium":
		limit = rl.config.RiotAPIProxyRPM + 50
	case "pro":
		limit = rl.config.RiotAPIProxyRPM + 200
	case "enterprise":
		limit = rl.config.RiotAPIProxyRPM + 500
	default:
		limit = rl.config.RiotAPIProxyRPM
	}
	
	key := fmt.Sprintf("herald:rate_limit:riot_api:%s", clientID)
	result := rl.checkLimit(ctx, key, limit, rl.config.GlobalWindow)
	result.Tier = tier
	result.Reason = "riot_api_proxy"
	
	return result
}

// checkLimit performs the actual rate limit check using Redis
func (rl *GamingRateLimiter) checkLimit(ctx context.Context, key string, limit int, window time.Duration) *RateLimitResult {
	pipe := rl.redisClient.TxPipeline()
	
	// Current timestamp
	now := time.Now()
	windowStart := now.Add(-window)
	
	// Remove old entries
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart.UnixNano(), 10))
	
	// Count current requests
	countCmd := pipe.ZCard(ctx, key)
	
	// Add current request
	pipe.ZAdd(ctx, key, &redis.Z{
		Score:  float64(now.UnixNano()),
		Member: fmt.Sprintf("%d-%d", now.UnixNano(), rl.generateRequestID()),
	})
	
	// Set expiration
	pipe.Expire(ctx, key, window+time.Minute)
	
	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		// On Redis error, allow the request but log the error
		return &RateLimitResult{
			Allowed:   true,
			Limit:     limit,
			Remaining: limit - 1,
			ResetTime: now.Add(window),
			Tier:      "unknown",
		}
	}
	
	currentCount := int(countCmd.Val())
	remaining := limit - currentCount
	if remaining < 0 {
		remaining = 0
	}
	
	// Calculate reset time
	resetTime := now.Add(window)
	
	// Calculate retry after
	retryAfter := time.Duration(0)
	if currentCount >= limit {
		retryAfter = window
	}
	
	return &RateLimitResult{
		Allowed:    currentCount <= limit,
		Limit:      limit,
		Remaining:  remaining,
		ResetTime:  resetTime,
		RetryAfter: retryAfter,
	}
}

// Helper methods

func (rl *GamingRateLimiter) getClientIdentifier(c *gin.Context) string {
	// Try to get user ID first
	if userID := rl.getUserID(c); userID != "" {
		return fmt.Sprintf("user:%s", userID)
	}
	
	// Try API key
	if apiKey := c.GetHeader("X-Gaming-API-Key"); apiKey != "" {
		return fmt.Sprintf("api_key:%s", apiKey)
	}
	
	// Fall back to IP address
	return fmt.Sprintf("ip:%s", c.ClientIP())
}

func (rl *GamingRateLimiter) getUserID(c *gin.Context) string {
	// Try to extract user ID from JWT token or session
	if userIDValue, exists := c.Get("gaming_user_id"); exists {
		if userID, ok := userIDValue.(string); ok {
			return userID
		}
	}
	
	// Try Authorization header
	if auth := c.GetHeader("Authorization"); strings.HasPrefix(auth, "Bearer ") {
		// Extract user ID from JWT token (simplified)
		// In real implementation, you would decode and validate the JWT
		return "jwt_user" // Placeholder
	}
	
	return ""
}

func (rl *GamingRateLimiter) getUserTier(c *gin.Context) string {
	// Try to get tier from context (set by auth middleware)
	if tierValue, exists := c.Get("gaming_subscription_tier"); exists {
		if tier, ok := tierValue.(string); ok {
			return tier
		}
	}
	
	// Check headers
	if tier := c.GetHeader("X-Gaming-Subscription-Tier"); tier != "" {
		return tier
	}
	
	return "free" // Default tier
}

func (rl *GamingRateLimiter) getTierLimit(tier string) int {
	switch tier {
	case "premium":
		return rl.config.PremiumTierRPM
	case "pro":
		return rl.config.ProTierRPM
	case "enterprise":
		return rl.config.EnterpriseTierRPM
	default:
		return rl.config.FreeTierRPM
	}
}

func (rl *GamingRateLimiter) extractRegion(c *gin.Context) string {
	// Try header first
	if region := c.GetHeader("X-Gaming-Region"); region != "" {
		return strings.ToUpper(region)
	}
	
	// Try query parameter
	if region := c.Query("region"); region != "" {
		return strings.ToUpper(region)
	}
	
	// Try to extract from URL path
	path := c.Request.URL.Path
	if strings.Contains(path, "/region/") {
		parts := strings.Split(path, "/region/")
		if len(parts) > 1 {
			regionPart := strings.Split(parts[1], "/")[0]
			return strings.ToUpper(regionPart)
		}
	}
	
	return ""
}

func (rl *GamingRateLimiter) setRateLimitHeaders(c *gin.Context, result *RateLimitResult) {
	c.Header("X-Gaming-Rate-Limit", strconv.Itoa(result.Limit))
	c.Header("X-Gaming-Rate-Remaining", strconv.Itoa(result.Remaining))
	c.Header("X-Gaming-Rate-Reset", strconv.FormatInt(result.ResetTime.Unix(), 10))
	c.Header("X-Gaming-Rate-Tier", result.Tier)
	
	if !result.Allowed {
		c.Header("Retry-After", strconv.Itoa(int(result.RetryAfter.Seconds())))
	}
}

func (rl *GamingRateLimiter) logRateLimitHit(c *gin.Context, clientID string, result *RateLimitResult) {
	// Log to Redis for analytics (simplified)
	logEntry := map[string]interface{}{
		"client_id":  clientID,
		"ip":         c.ClientIP(),
		"path":       c.Request.URL.Path,
		"method":     c.Request.Method,
		"tier":       result.Tier,
		"limit":      result.Limit,
		"remaining":  result.Remaining,
		"timestamp":  time.Now().Unix(),
		"user_agent": c.GetHeader("User-Agent"),
	}
	
	key := fmt.Sprintf("herald:rate_limit:hits:%s", time.Now().Format("2006-01-02"))
	rl.redisClient.LPush(context.Background(), key, logEntry)
	rl.redisClient.Expire(context.Background(), key, 7*24*time.Hour) // Keep logs for 7 days
}

func (rl *GamingRateLimiter) getUpgradeInfo(currentTier string) gin.H {
	upgrades := gin.H{
		"upgrade_url": "https://herald.lol/upgrade",
		"contact": "sales@herald.lol",
	}
	
	switch currentTier {
	case "free":
		upgrades["recommended"] = "premium"
		upgrades["benefit"] = "5x higher rate limits + advanced analytics"
	case "premium":
		upgrades["recommended"] = "pro"
		upgrades["benefit"] = "Team features + 4x higher limits"
	case "pro":
		upgrades["recommended"] = "enterprise"
		upgrades["benefit"] = "Unlimited API + dedicated support"
	}
	
	return upgrades
}

func (rl *GamingRateLimiter) generateRequestID() int64 {
	return time.Now().UnixNano()
}