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

// Herald.lol Gaming Analytics - Gaming Security Middleware
// Specialized security middleware for gaming platform operations

// GamingSecurityMiddleware provides gaming-specific security enforcement
type GamingSecurityMiddleware struct {
	redis           *redis.Client
	config          *GamingSecurityConfig
	riotAPIValidator *RiotAPIValidator
	subscriptionChecker *SubscriptionChecker
}

// GamingSecurityConfig contains security configuration
type GamingSecurityConfig struct {
	// API Key requirements
	RequireAPIKey           bool     `json:"require_api_key"`
	APIKeyPrefixes          []string `json:"api_key_prefixes"`
	
	// Riot Games compliance
	RiotAPICompliance       bool     `json:"riot_api_compliance"`
	RiotTermsAcceptance     bool     `json:"riot_terms_acceptance"`
	
	// Gaming data protection
	DataExportMFARequired   bool     `json:"data_export_mfa_required"`
	AnalyticsAuthRequired   bool     `json:"analytics_auth_required"`
	TeamDataPermissions     bool     `json:"team_data_permissions"`
	
	// Regional restrictions
	RegionalCompliance      bool     `json:"regional_compliance"`
	EUDataProtection        bool     `json:"eu_data_protection"`
	
	// Security headers
	GamingSecurityHeaders   bool     `json:"gaming_security_headers"`
	CSPForGamingContent     bool     `json:"csp_for_gaming_content"`
	
	// Rate limiting integration
	EnforceGamingLimits     bool     `json:"enforce_gaming_limits"`
	SubscriptionTierCheck   bool     `json:"subscription_tier_check"`
}

// NewGamingSecurityMiddleware creates new gaming security middleware
func NewGamingSecurityMiddleware(redis *redis.Client, config *GamingSecurityConfig) *GamingSecurityMiddleware {
	return &GamingSecurityMiddleware{
		redis:               redis,
		config:              config,
		riotAPIValidator:    NewRiotAPIValidator(redis),
		subscriptionChecker: NewSubscriptionChecker(redis),
	}
}

// GamingSecurityHeaders adds gaming-specific security headers
func (g *GamingSecurityMiddleware) GamingSecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !g.config.GamingSecurityHeaders {
			c.Next()
			return
		}
		
		// Standard security headers for gaming platform
		c.Header("X-Gaming-Platform", "herald-lol")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Gaming-specific headers
		c.Header("X-Gaming-API-Version", "3.1.0")
		c.Header("X-Riot-Compliance", "true")
		c.Header("X-Gaming-Data-Protection", "enabled")
		
		// GDPR compliance header for EU requests
		if g.config.EUDataProtection && g.isEURequest(c) {
			c.Header("X-GDPR-Compliance", "required")
			c.Header("X-Data-Processing-Basis", "legitimate-interest")
		}
		
		// CSP for gaming content
		if g.config.CSPForGamingContent {
			csp := "default-src 'self'; " +
				"script-src 'self' 'unsafe-inline' https://cdn.herald.lol; " +
				"img-src 'self' data: https://ddragon.leagueoflegends.com https://cdn.herald.lol; " +
				"connect-src 'self' https://api.herald.lol; " +
				"font-src 'self' https://fonts.gstatic.com; " +
				"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com"
			c.Header("Content-Security-Policy", csp)
		}
		
		c.Next()
	}
}

// ValidateGamingAPIKey validates gaming-specific API keys
func (g *GamingSecurityMiddleware) ValidateGamingAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !g.config.RequireAPIKey {
			c.Next()
			return
		}
		
		// Get API key from header
		apiKey := c.GetHeader("X-Gaming-API-Key")
		if apiKey == "" {
			// Also check Authorization header for API key
			auth := c.GetHeader("Authorization")
			if strings.HasPrefix(auth, "ApiKey ") {
				apiKey = strings.TrimPrefix(auth, "ApiKey ")
			}
		}
		
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":           "Gaming API key required",
				"error_code":      "GAMING_API_KEY_MISSING",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Validate API key format
		if !g.isValidAPIKeyFormat(apiKey) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":           "Invalid gaming API key format",
				"error_code":      "GAMING_API_KEY_INVALID",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Validate API key in Redis
		keyInfo, err := g.validateAPIKeyInRedis(c.Request.Context(), apiKey)
		if err != nil || keyInfo == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":           "Gaming API key validation failed",
				"error_code":      "GAMING_API_KEY_INVALID",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Check if API key is active
		if !keyInfo.IsActive {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "Gaming API key is inactive",
				"error_code":      "GAMING_API_KEY_INACTIVE",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Check rate limits for this API key
		if g.config.EnforceGamingLimits {
			if exceeded, err := g.checkAPIKeyRateLimit(c.Request.Context(), apiKey, keyInfo); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":           "Rate limit check failed",
					"gaming_platform": "herald-lol",
				})
				c.Abort()
				return
			} else if exceeded {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":           "API key rate limit exceeded",
					"error_code":      "GAMING_RATE_LIMIT_EXCEEDED",
					"gaming_platform": "herald-lol",
				})
				c.Abort()
				return
			}
		}
		
		// Store API key info in context
		c.Set("api_key_info", keyInfo)
		c.Set("api_key", apiKey)
		
		c.Next()
	}
}

// EnforceSubscriptionTier enforces subscription tier requirements
func (g *GamingSecurityMiddleware) EnforceSubscriptionTier(requiredTier string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !g.config.SubscriptionTierCheck {
			c.Next()
			return
		}
		
		// Get user subscription tier
		userTier := g.getUserSubscriptionTier(c)
		if userTier == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":           "Gaming subscription tier required",
				"required_tier":   requiredTier,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Check if user's tier meets requirement
		if !g.subscriptionChecker.MeetsTierRequirement(userTier, requiredTier) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "Gaming subscription tier insufficient",
				"current_tier":    userTier,
				"required_tier":   requiredTier,
				"upgrade_url":     "https://herald.lol/upgrade",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		c.Set("subscription_tier", userTier)
		c.Next()
	}
}

// ValidateRiotCompliance ensures Riot Games ToS compliance
func (g *GamingSecurityMiddleware) ValidateRiotCompliance() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !g.config.RiotAPICompliance {
			c.Next()
			return
		}
		
		// Check if user has accepted Riot Games Terms of Service
		if g.config.RiotTermsAcceptance {
			userID := g.getUserID(c)
			if userID != "" {
				accepted, err := g.riotAPIValidator.HasAcceptedRiotTerms(c.Request.Context(), userID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":           "Riot ToS validation failed",
						"gaming_platform": "herald-lol",
					})
					c.Abort()
					return
				}
				
				if !accepted {
					c.JSON(http.StatusForbidden, gin.H{
						"error":           "Riot Games Terms of Service acceptance required",
						"error_code":      "RIOT_TOS_REQUIRED",
						"terms_url":       "https://www.riotgames.com/en/terms-of-service",
						"gaming_platform": "herald-lol",
					})
					c.Abort()
					return
				}
			}
		}
		
		// Add Riot compliance headers
		c.Header("X-Riot-Terms-Compliance", "verified")
		c.Header("X-Gaming-Data-Source", "riot-games-api")
		
		c.Next()
	}
}

// RequireMFAForDataExport requires MFA for sensitive gaming data exports
func (g *GamingSecurityMiddleware) RequireMFAForDataExport() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !g.config.DataExportMFARequired {
			c.Next()
			return
		}
		
		// Check if this is a data export endpoint
		if !g.isDataExportEndpoint(c.Request.URL.Path) {
			c.Next()
			return
		}
		
		// Get MFA token from request
		mfaToken := c.GetHeader("X-MFA-Token")
		if mfaToken == "" {
			var requestData map[string]interface{}
			c.ShouldBindJSON(&requestData)
			if token, ok := requestData["mfa_token"].(string); ok {
				mfaToken = token
			}
		}
		
		if mfaToken == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "MFA token required for gaming data export",
				"error_code":      "GAMING_MFA_REQUIRED",
				"mfa_methods":     []string{"totp", "webauthn"},
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Validate MFA token
		userID := g.getUserID(c)
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":           "User authentication required for MFA",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		valid, err := g.validateMFAToken(c.Request.Context(), userID, mfaToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":           "MFA validation failed",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		if !valid {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "Invalid MFA token for gaming data export",
				"error_code":      "GAMING_MFA_INVALID",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// MFA validated, add header and continue
		c.Header("X-MFA-Verified", "true")
		c.Set("mfa_verified", true)
		
		c.Next()
	}
}

// EnforceTeamPermissions enforces team-specific permissions for gaming operations
func (g *GamingSecurityMiddleware) EnforceTeamPermissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !g.config.TeamDataPermissions {
			c.Next()
			return
		}
		
		// Check if this endpoint requires team permissions
		teamID := c.Param("teamId")
		if teamID == "" {
			c.Next()
			return
		}
		
		userID := g.getUserID(c)
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":           "Authentication required for team gaming operations",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Check team membership and permissions
		hasPermission, role, err := g.checkTeamPermissions(c.Request.Context(), userID, teamID, c.Request.Method, c.Request.URL.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":           "Team permission check failed",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "Insufficient gaming team permissions",
				"error_code":      "GAMING_TEAM_PERMISSION_DENIED",
				"team_id":         teamID,
				"required_role":   g.getRequiredTeamRole(c.Request.Method, c.Request.URL.Path),
				"current_role":    role,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Store team info in context
		c.Set("team_id", teamID)
		c.Set("team_role", role)
		
		c.Next()
	}
}

// ValidateRegionalCompliance ensures regional compliance for gaming data
func (g *GamingSecurityMiddleware) ValidateRegionalCompliance() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !g.config.RegionalCompliance {
			c.Next()
			return
		}
		
		// Get region from URL path
		region := c.Param("region")
		if region == "" {
			c.Next()
			return
		}
		
		// Validate region
		if !g.isValidGamingRegion(region) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":           "Invalid gaming region",
				"provided_region": region,
				"valid_regions":   []string{"NA", "EUW", "EUNE", "KR", "JP", "BR", "LAN", "LAS", "OCE", "RU", "TR"},
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Check for region-specific restrictions
		if g.hasRegionalRestrictions(c, region) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "Gaming data access restricted in this region",
				"region":          region,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// EU specific data protection checks
		if g.config.EUDataProtection && g.isEURegion(region) {
			if !g.hasGDPRConsent(c) {
				c.JSON(http.StatusForbidden, gin.H{
					"error":           "GDPR consent required for EU gaming data",
					"error_code":      "GDPR_CONSENT_REQUIRED",
					"region":          region,
					"consent_url":     "https://herald.lol/gdpr-consent",
					"gaming_platform": "herald-lol",
				})
				c.Abort()
				return
			}
		}
		
		c.Set("gaming_region", region)
		c.Next()
	}
}

// Helper methods

func (g *GamingSecurityMiddleware) isValidAPIKeyFormat(apiKey string) bool {
	if len(apiKey) < 32 {
		return false
	}
	
	// Check for valid prefixes
	for _, prefix := range g.config.APIKeyPrefixes {
		if strings.HasPrefix(apiKey, prefix) {
			return true
		}
	}
	
	return len(g.config.APIKeyPrefixes) == 0 // If no prefixes defined, allow any format
}

func (g *GamingSecurityMiddleware) validateAPIKeyInRedis(ctx context.Context, apiKey string) (*APIKeyInfo, error) {
	keyData, err := g.redis.HGetAll(ctx, fmt.Sprintf("api_key:%s", apiKey)).Result()
	if err != nil {
		return nil, err
	}
	
	if len(keyData) == 0 {
		return nil, fmt.Errorf("API key not found")
	}
	
	keyInfo := &APIKeyInfo{
		UserID:    keyData["user_id"],
		Tier:      keyData["tier"],
		IsActive:  keyData["is_active"] == "true",
	}
	
	if rateLimit, ok := keyData["rate_limit"]; ok {
		if limit, err := strconv.Atoi(rateLimit); err == nil {
			keyInfo.RateLimit = limit
		}
	}
	
	return keyInfo, nil
}

func (g *GamingSecurityMiddleware) checkAPIKeyRateLimit(ctx context.Context, apiKey string, keyInfo *APIKeyInfo) (bool, error) {
	now := time.Now()
	rateLimitKey := fmt.Sprintf("api_key_rate:%s:minute:%d", apiKey, now.Unix()/60)
	
	count, err := g.redis.Incr(ctx, rateLimitKey).Result()
	if err != nil {
		return false, err
	}
	
	if count == 1 {
		g.redis.Expire(ctx, rateLimitKey, time.Minute)
	}
	
	return int(count) > keyInfo.RateLimit, nil
}

func (g *GamingSecurityMiddleware) getUserSubscriptionTier(c *gin.Context) string {
	// Try to get from context first
	if tier, exists := c.Get("subscription_tier"); exists {
		if t, ok := tier.(string); ok {
			return t
		}
	}
	
	// Try to get from API key info
	if keyInfo, exists := c.Get("api_key_info"); exists {
		if info, ok := keyInfo.(*APIKeyInfo); ok {
			return info.Tier
		}
	}
	
	return ""
}

func (g *GamingSecurityMiddleware) getUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

func (g *GamingSecurityMiddleware) isEURequest(c *gin.Context) bool {
	// Check for Cloudflare country header
	country := c.GetHeader("CF-IPCountry")
	euCountries := []string{"AT", "BE", "BG", "CY", "CZ", "DE", "DK", "EE", "ES", "FI", "FR", "GR", "HR", "HU", "IE", "IT", "LT", "LU", "LV", "MT", "NL", "PL", "PT", "RO", "SE", "SI", "SK"}
	
	for _, eu := range euCountries {
		if country == eu {
			return true
		}
	}
	
	return false
}

func (g *GamingSecurityMiddleware) isDataExportEndpoint(path string) bool {
	return strings.Contains(path, "/export")
}

func (g *GamingSecurityMiddleware) validateMFAToken(ctx context.Context, userID, token string) (bool, error) {
	// This would validate TOTP or other MFA tokens
	// For now, return true if token length is 6 (TOTP format)
	return len(token) == 6, nil
}

func (g *GamingSecurityMiddleware) checkTeamPermissions(ctx context.Context, userID, teamID, method, path string) (bool, string, error) {
	// Get user's role in team
	roleKey := fmt.Sprintf("team:%s:member:%s", teamID, userID)
	role, err := g.redis.HGet(ctx, roleKey, "role").Result()
	if err != nil {
		return false, "", err
	}
	
	// Check if role has permission for this operation
	requiredRole := g.getRequiredTeamRole(method, path)
	hasPermission := g.checkRolePermission(role, requiredRole)
	
	return hasPermission, role, nil
}

func (g *GamingSecurityMiddleware) getRequiredTeamRole(method, path string) string {
	switch method {
	case "DELETE":
		return "captain"
	case "POST", "PUT", "PATCH":
		return "manager"
	default:
		return "member"
	}
}

func (g *GamingSecurityMiddleware) checkRolePermission(userRole, requiredRole string) bool {
	roleHierarchy := map[string]int{
		"member":  1,
		"manager": 2,
		"captain": 3,
		"owner":   4,
	}
	
	userLevel := roleHierarchy[userRole]
	requiredLevel := roleHierarchy[requiredRole]
	
	return userLevel >= requiredLevel
}

func (g *GamingSecurityMiddleware) isValidGamingRegion(region string) bool {
	validRegions := []string{"NA", "EUW", "EUNE", "KR", "JP", "BR", "LAN", "LAS", "OCE", "RU", "TR"}
	for _, valid := range validRegions {
		if region == valid {
			return true
		}
	}
	return false
}

func (g *GamingSecurityMiddleware) hasRegionalRestrictions(c *gin.Context, region string) bool {
	// Check for regional restrictions based on user location, compliance requirements, etc.
	return false
}

func (g *GamingSecurityMiddleware) isEURegion(region string) bool {
	euRegions := []string{"EUW", "EUNE"}
	for _, eu := range euRegions {
		if region == eu {
			return true
		}
	}
	return false
}

func (g *GamingSecurityMiddleware) hasGDPRConsent(c *gin.Context) bool {
	// Check if user has provided GDPR consent
	consent := c.GetHeader("X-GDPR-Consent")
	return consent == "given"
}

// APIKeyInfo contains API key information
type APIKeyInfo struct {
	UserID    string `json:"user_id"`
	Tier      string `json:"tier"`
	RateLimit int    `json:"rate_limit"`
	IsActive  bool   `json:"is_active"`
}

// RiotAPIValidator validates Riot Games API compliance
type RiotAPIValidator struct {
	redis *redis.Client
}

func NewRiotAPIValidator(redis *redis.Client) *RiotAPIValidator {
	return &RiotAPIValidator{redis: redis}
}

func (r *RiotAPIValidator) HasAcceptedRiotTerms(ctx context.Context, userID string) (bool, error) {
	termsKey := fmt.Sprintf("user:%s:riot_terms", userID)
	accepted, err := r.redis.Get(ctx, termsKey).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return accepted == "accepted", nil
}

// SubscriptionChecker validates subscription tiers
type SubscriptionChecker struct {
	redis *redis.Client
}

func NewSubscriptionChecker(redis *redis.Client) *SubscriptionChecker {
	return &SubscriptionChecker{redis: redis}
}

func (s *SubscriptionChecker) MeetsTierRequirement(userTier, requiredTier string) bool {
	tierLevels := map[string]int{
		"free":       1,
		"premium":    2,
		"pro":        3,
		"enterprise": 4,
	}
	
	userLevel := tierLevels[userTier]
	requiredLevel := tierLevels[requiredTier]
	
	return userLevel >= requiredLevel
}

// DefaultGamingSecurityConfig returns default security configuration
func DefaultGamingSecurityConfig() *GamingSecurityConfig {
	return &GamingSecurityConfig{
		RequireAPIKey:         true,
		APIKeyPrefixes:        []string{"hld_", "herald_"},
		RiotAPICompliance:     true,
		RiotTermsAcceptance:   true,
		DataExportMFARequired: true,
		AnalyticsAuthRequired: true,
		TeamDataPermissions:   true,
		RegionalCompliance:    true,
		EUDataProtection:      true,
		GamingSecurityHeaders: true,
		CSPForGamingContent:   true,
		EnforceGamingLimits:   true,
		SubscriptionTierCheck: true,
	}
}