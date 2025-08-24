package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Herald.lol Gaming Analytics - Authentication Middleware
// Comprehensive JWT and OAuth middleware for gaming platform

// GamingAuthMiddleware provides authentication middleware for Herald.lol
type GamingAuthMiddleware struct {
	config    *GamingOAuthConfig
	userStore *DatabaseUserStore
}

// NewGamingAuthMiddleware creates new authentication middleware for gaming platform
func NewGamingAuthMiddleware(config *GamingOAuthConfig, userStore *DatabaseUserStore) *GamingAuthMiddleware {
	return &GamingAuthMiddleware{
		config:    config,
		userStore: userStore,
	}
}

// RequireGamingAuth middleware that requires valid gaming JWT token
func (m *GamingAuthMiddleware) RequireGamingAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract gaming token from multiple sources
		token := m.extractGamingToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":           "Gaming authentication required",
				"gaming_platform": "herald-lol",
				"auth_endpoints": gin.H{
					"login":   "/auth/oauth/:provider",
					"refresh": "/auth/refresh",
				},
			})
			c.Abort()
			return
		}

		// Parse and validate gaming JWT token
		claims, err := m.config.parseJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":           "Invalid gaming authentication token",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}

		// Check token expiration
		if time.Now().After(claims.ExpiresAt.Time) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":           "Gaming authentication token expired",
				"gaming_platform": "herald-lol",
				"action":          "refresh_token",
			})
			c.Abort()
			return
		}

		// Get full gaming user info
		ctx := context.Background()
		user, err := m.userStore.GetUserByProviderID(ctx, OAuthProvider(claims.Provider), claims.ProviderID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":           "Gaming user not found",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}

		// Check if gaming user is active
		if !user.Metadata["is_active"] == "true" && user.Metadata["is_active"] != "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "Gaming account is disabled",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}

		// Update session activity if available
		if sessionToken := c.GetHeader("X-Session-Token"); sessionToken != "" {
			go m.userStore.UpdateGamingSessionActivity(ctx, sessionToken)
		}

		// Set gaming user context
		c.Set("gaming_user", user)
		c.Set("gaming_claims", claims)
		c.Set("gaming_user_id", claims.UserID)
		c.Set("gaming_subscription_tier", claims.SubscriptionTier)
		c.Set("gaming_permissions", claims.GamingPermissions)

		// Add gaming headers for downstream services
		c.Header("X-Gaming-User-ID", claims.UserID)
		c.Header("X-Gaming-Tier", claims.SubscriptionTier)
		c.Header("X-Gaming-Platform", "herald-lol")

		c.Next()
	}
}

// RequireGamingPermission middleware that requires specific gaming permission
func (m *GamingAuthMiddleware) RequireGamingPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First require authentication
		m.RequireGamingAuth()(c)
		if c.IsAborted() {
			return
		}

		// Get gaming permissions from context
		permissionsInterface, exists := c.Get("gaming_permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "Gaming permissions not found",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}

		permissions, ok := permissionsInterface.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "Invalid gaming permissions format",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}

		// Check if user has required gaming permission
		hasPermission := false
		for _, perm := range permissions {
			if perm == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":               "Insufficient gaming permissions",
				"required_permission": permission,
				"gaming_platform":     "herald-lol",
				"upgrade_info":        m.getUpgradeInfo(permission),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireGamingTier middleware that requires minimum gaming subscription tier
func (m *GamingAuthMiddleware) RequireGamingTier(minTier string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First require authentication
		m.RequireGamingAuth()(c)
		if c.IsAborted() {
			return
		}

		// Get gaming tier from context
		tier, exists := c.Get("gaming_subscription_tier")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "Gaming subscription tier not found",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}

		userTier, ok := tier.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "Invalid gaming subscription tier",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}

		// Check if user meets minimum gaming tier requirement
		if !m.meetsTierRequirement(userTier, minTier) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "Insufficient gaming subscription tier",
				"current_tier":    userTier,
				"required_tier":   minTier,
				"gaming_platform": "herald-lol",
				"upgrade_url":     "https://herald.lol/upgrade",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalGamingAuth middleware that optionally extracts gaming user if token is present
func (m *GamingAuthMiddleware) OptionalGamingAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract gaming token
		token := m.extractGamingToken(c)
		if token == "" {
			// No token provided, continue without authentication
			c.Set("gaming_authenticated", false)
			c.Next()
			return
		}

		// Parse and validate gaming JWT token
		claims, err := m.config.parseJWT(token)
		if err != nil || time.Now().After(claims.ExpiresAt.Time) {
			// Invalid or expired token, continue without authentication
			c.Set("gaming_authenticated", false)
			c.Next()
			return
		}

		// Get gaming user info
		ctx := context.Background()
		user, err := m.userStore.GetUserByProviderID(ctx, OAuthProvider(claims.Provider), claims.ProviderID)
		if err != nil {
			// User not found, continue without authentication
			c.Set("gaming_authenticated", false)
			c.Next()
			return
		}

		// Set gaming user context
		c.Set("gaming_authenticated", true)
		c.Set("gaming_user", user)
		c.Set("gaming_claims", claims)
		c.Set("gaming_user_id", claims.UserID)
		c.Set("gaming_subscription_tier", claims.SubscriptionTier)
		c.Set("gaming_permissions", claims.GamingPermissions)

		// Add gaming headers
		c.Header("X-Gaming-User-ID", claims.UserID)
		c.Header("X-Gaming-Tier", claims.SubscriptionTier)
		c.Header("X-Gaming-Platform", "herald-lol")

		c.Next()
	}
}

// RateLimitByGamingTier middleware that applies rate limiting based on gaming subscription tier
func (m *GamingAuthMiddleware) RateLimitByGamingTier() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get gaming tier (requires OptionalGamingAuth or RequireGamingAuth to run first)
		tier := "free" // Default to free tier

		if tierInterface, exists := c.Get("gaming_subscription_tier"); exists {
			if tierString, ok := tierInterface.(string); ok {
				tier = tierString
			}
		}

		// Set rate limiting headers based on gaming tier
		rateLimits := m.getGamingRateLimits(tier)
		c.Header("X-Gaming-RateLimit-Limit", rateLimits["limit"])
		c.Header("X-Gaming-RateLimit-Window", rateLimits["window"])
		c.Header("X-Gaming-RateLimit-Tier", tier)

		// The actual rate limiting would be handled by Kong or other middleware
		// This middleware just sets the appropriate headers

		c.Next()
	}
}

// CORSForGaming middleware that handles CORS for gaming platform
func (m *GamingAuthMiddleware) CORSForGaming() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Gaming platform allowed origins
		allowedOrigins := []string{
			"https://herald.lol",
			"https://www.herald.lol",
			"https://app.herald.lol",
			"https://staging.herald.lol",
			"http://localhost:3000", // Development
			"http://localhost:3001", // Development
		}

		// Check if origin is allowed for gaming platform
		originAllowed := false
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				originAllowed = true
				break
			}
		}

		if originAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Session-Token, X-Gaming-Client, X-Herald-Version")
		c.Header("Access-Control-Expose-Headers", "X-Gaming-User-ID, X-Gaming-Tier, X-Gaming-Platform, X-Gaming-RateLimit-Limit, X-Gaming-RateLimit-Remaining")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// GamingSecurityHeaders middleware that adds security headers for gaming platform
func (m *GamingAuthMiddleware) GamingSecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Gaming platform security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("X-Gaming-Platform", "herald-lol")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' https://cdn.herald.lol; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; img-src 'self' data: https:; font-src 'self' https://fonts.gstatic.com")

		// Gaming API specific headers
		c.Header("X-Robots-Tag", "noindex, nofollow")
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.Next()
	}
}

// Helper methods

// extractGamingToken extracts JWT token from various sources
func (m *GamingAuthMiddleware) extractGamingToken(c *gin.Context) string {
	// 1. Check Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// Handle "Bearer <token>" format
		if strings.HasPrefix(authHeader, "Bearer ") {
			return authHeader[7:]
		}
		// Handle direct token
		return authHeader
	}

	// 2. Check X-Auth-Token header (gaming specific)
	if token := c.GetHeader("X-Auth-Token"); token != "" {
		return token
	}

	// 3. Check query parameter
	if token := c.Query("token"); token != "" {
		return token
	}

	// 4. Check gaming access token cookie
	if token, err := c.Cookie("herald_access_token"); err == nil && token != "" {
		return token
	}

	return ""
}

// meetsTierRequirement checks if user tier meets minimum requirement
func (m *GamingAuthMiddleware) meetsTierRequirement(userTier, minTier string) bool {
	tierLevels := map[string]int{
		"free":       0,
		"premium":    1,
		"pro":        2,
		"enterprise": 3,
	}

	userLevel, userExists := tierLevels[userTier]
	minLevel, minExists := tierLevels[minTier]

	if !userExists || !minExists {
		return false
	}

	return userLevel >= minLevel
}

// getUpgradeInfo returns upgrade information for insufficient permissions
func (m *GamingAuthMiddleware) getUpgradeInfo(permission string) map[string]string {
	upgradeInfo := map[string]string{
		"upgrade_url": "https://herald.lol/upgrade",
		"contact":     "support@herald.lol",
	}

	// Map permissions to required tiers
	permissionTiers := map[string]string{
		"analytics:advanced": "premium",
		"api:extended":       "premium",
		"coaching:premium":   "pro",
		"team:management":    "enterprise",
		"api:unlimited":      "enterprise",
	}

	if requiredTier, exists := permissionTiers[permission]; exists {
		upgradeInfo["required_tier"] = requiredTier
		upgradeInfo["upgrade_url"] = "https://herald.lol/upgrade?tier=" + requiredTier
	}

	return upgradeInfo
}

// getGamingRateLimits returns rate limits for gaming tiers
func (m *GamingAuthMiddleware) getGamingRateLimits(tier string) map[string]string {
	limits := map[string]map[string]string{
		"free": {
			"limit":  "100",
			"window": "60", // per minute
		},
		"premium": {
			"limit":  "500",
			"window": "60",
		},
		"pro": {
			"limit":  "2000",
			"window": "60",
		},
		"enterprise": {
			"limit":  "10000",
			"window": "60",
		},
	}

	if tierLimits, exists := limits[tier]; exists {
		return tierLimits
	}

	return limits["free"] // Default to free tier limits
}

// Gaming Context Helpers

// GetGamingUser extracts gaming user from Gin context
func GetGamingUser(c *gin.Context) (*GamingUserInfo, bool) {
	if user, exists := c.Get("gaming_user"); exists {
		if gamingUser, ok := user.(*GamingUserInfo); ok {
			return gamingUser, true
		}
	}
	return nil, false
}

// GetGamingClaims extracts gaming JWT claims from Gin context
func GetGamingClaims(c *gin.Context) (*GamingJWTClaims, bool) {
	if claims, exists := c.Get("gaming_claims"); exists {
		if gamingClaims, ok := claims.(*GamingJWTClaims); ok {
			return gamingClaims, true
		}
	}
	return nil, false
}

// GetGamingUserID extracts gaming user ID from Gin context
func GetGamingUserID(c *gin.Context) (string, bool) {
	if userID, exists := c.Get("gaming_user_id"); exists {
		if userIDString, ok := userID.(string); ok {
			return userIDString, true
		}
	}
	return "", false
}

// IsGamingAuthenticated checks if request is authenticated for gaming platform
func IsGamingAuthenticated(c *gin.Context) bool {
	if authenticated, exists := c.Get("gaming_authenticated"); exists {
		if auth, ok := authenticated.(bool); ok {
			return auth
		}
	}
	return false
}

// GetGamingSubscriptionTier extracts gaming subscription tier from Gin context
func GetGamingSubscriptionTier(c *gin.Context) (string, bool) {
	if tier, exists := c.Get("gaming_subscription_tier"); exists {
		if tierString, ok := tier.(string); ok {
			return tierString, true
		}
	}
	return "free", false
}

// HasGamingPermission checks if user has specific gaming permission
func HasGamingPermission(c *gin.Context, permission string) bool {
	if permissions, exists := c.Get("gaming_permissions"); exists {
		if permSlice, ok := permissions.([]string); ok {
			for _, perm := range permSlice {
				if perm == permission {
					return true
				}
			}
		}
	}
	return false
}
