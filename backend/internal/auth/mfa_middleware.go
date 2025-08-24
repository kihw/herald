package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Herald.lol Gaming Analytics - MFA Middleware
// Multi-Factor Authentication middleware for protecting gaming actions

// GamingMFAMiddleware provides MFA enforcement for gaming platform
type GamingMFAMiddleware struct {
	mfaManager *GamingMFAManager
}

// NewGamingMFAMiddleware creates new MFA middleware for gaming platform
func NewGamingMFAMiddleware(mfaManager *GamingMFAManager) *GamingMFAMiddleware {
	return &GamingMFAMiddleware{
		mfaManager: mfaManager,
	}
}

// RequireGamingMFA middleware that requires MFA for specific actions
func (m *GamingMFAMiddleware) RequireGamingMFA(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			c.Abort()
			return
		}

		// Check if action requires MFA
		requiresMFA := m.mfaManager.actionRequiresMFA(action)
		if !requiresMFA {
			c.Next()
			return
		}

		// Get user's MFA status
		ctx := context.Background()
		status, err := m.mfaManager.mfaStore.(*CombinedGamingMFAStore).GetMFAStatus(ctx, user.ID)
		if err != nil || status == nil || !status.Enabled {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "MFA required for this gaming action",
				"action":          action,
				"mfa_configured":  false,
				"setup_url":       "/auth/mfa/totp/setup",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}

		// Check for MFA token in request
		mfaToken := m.extractMFAToken(c)
		if mfaToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":       "MFA authentication required for this gaming action",
				"action":      action,
				"mfa_methods": status.Methods,
				"mfa_endpoints": gin.H{
					"totp":         "/auth/mfa/totp/authenticate",
					"webauthn":     "/auth/mfa/webauthn/authenticate/begin",
					"backup_codes": "/auth/mfa/backup-codes/authenticate",
				},
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}

		// Verify MFA token
		if !m.mfaManager.verifyMFAToken(mfaToken, action) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":           "Invalid or expired gaming MFA token",
				"action":          action,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}

		// Set MFA context
		c.Set("mfa_verified", true)
		c.Set("mfa_action", action)
		c.Set("mfa_token", mfaToken)

		// Add MFA headers
		c.Header("X-Gaming-MFA-Verified", "true")
		c.Header("X-Gaming-MFA-Action", action)

		c.Next()
	}
}

// RequireGamingMFAForHighValue middleware for high-value gaming actions
func (m *GamingMFAMiddleware) RequireGamingMFAForHighValue() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Determine action from request
		action := m.determineGamingAction(c)

		// Apply MFA middleware for the action
		m.RequireGamingMFA(action)(c)
	}
}

// EnforceGamingSessionMFA middleware that enforces MFA for gaming sessions
func (m *GamingMFAMiddleware) EnforceGamingSessionMFA() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.mfaManager.config.GamingSessionMFARequired {
			c.Next()
			return
		}

		user, exists := GetGamingUser(c)
		if !exists {
			c.Next()
			return
		}

		// Check if user has MFA enabled
		ctx := context.Background()
		status, err := m.mfaManager.mfaStore.(*CombinedGamingMFAStore).GetMFAStatus(ctx, user.ID)
		if err != nil || status == nil || !status.Enabled {
			// If gaming session MFA is required but not configured, require setup
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "MFA setup required for gaming sessions",
				"gaming_platform": "herald-lol",
				"setup_required":  true,
				"setup_url":       "/auth/mfa/totp/setup",
			})
			c.Abort()
			return
		}

		// Check for recent MFA verification in session
		if !m.hasRecentMFAVerification(c) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":           "MFA verification required for gaming session",
				"gaming_platform": "herald-lol",
				"mfa_required":    true,
				"mfa_methods":     status.Methods,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalGamingMFA middleware that optionally checks MFA
func (m *GamingMFAMiddleware) OptionalGamingMFA() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.Set("mfa_configured", false)
			c.Next()
			return
		}

		// Get MFA status
		ctx := context.Background()
		status, err := m.mfaManager.mfaStore.(*CombinedGamingMFAStore).GetMFAStatus(ctx, user.ID)

		mfaConfigured := err == nil && status != nil && status.Enabled
		c.Set("mfa_configured", mfaConfigured)

		if mfaConfigured {
			c.Set("mfa_methods", status.Methods)
			c.Set("mfa_status", status)
		}

		// Check for MFA token
		mfaToken := m.extractMFAToken(c)
		if mfaToken != "" {
			c.Set("mfa_token_present", true)
			c.Set("mfa_verified", m.mfaManager.verifyMFAToken(mfaToken, ""))
		} else {
			c.Set("mfa_token_present", false)
			c.Set("mfa_verified", false)
		}

		c.Next()
	}
}

// GamingAnalyticsMFA middleware for analytics operations
func (m *GamingMFAMiddleware) GamingAnalyticsMFA() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.mfaManager.config.AnalyticsMFARequired {
			c.Next()
			return
		}

		// Determine analytics action
		action := "analytics:" + m.getAnalyticsOperation(c)

		// Apply MFA requirement
		m.RequireGamingMFA(action)(c)
	}
}

// GamingAPIMFA middleware for API access
func (m *GamingMFAMiddleware) GamingAPIMFA() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.mfaManager.config.APIAccessMFARequired {
			c.Next()
			return
		}

		// Check if this is an API request requiring MFA
		if m.isHighPrivilegeAPIRequest(c) {
			action := "api:high_privilege"
			m.RequireGamingMFA(action)(c)
			return
		}

		c.Next()
	}
}

// MFARecoveryGuard middleware that protects MFA recovery operations
func (m *GamingMFAMiddleware) MFARecoveryGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Require additional verification for MFA recovery
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			c.Abort()
			return
		}

		// Check for password confirmation or other recovery proof
		recoveryToken := c.GetHeader("X-Gaming-Recovery-Token")
		if recoveryToken == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "Recovery verification required for gaming MFA operations",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}

		// Verify recovery token (implementation would verify password/email/etc)
		if !m.verifyRecoveryToken(user.ID, recoveryToken) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":           "Invalid gaming recovery verification",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}

		c.Set("recovery_verified", true)
		c.Next()
	}
}

// Helper methods

// extractMFAToken extracts MFA token from request
func (m *GamingMFAMiddleware) extractMFAToken(c *gin.Context) string {
	// Check X-Gaming-MFA-Token header
	if token := c.GetHeader("X-Gaming-MFA-Token"); token != "" {
		return token
	}

	// Check X-MFA-Token header
	if token := c.GetHeader("X-MFA-Token"); token != "" {
		return token
	}

	// Check query parameter
	if token := c.Query("mfa_token"); token != "" {
		return token
	}

	// Check form data
	if token := c.PostForm("mfa_token"); token != "" {
		return token
	}

	return ""
}

// determineGamingAction determines the gaming action from request
func (m *GamingMFAMiddleware) determineGamingAction(c *gin.Context) string {
	// Check X-Gaming-Action header
	if action := c.GetHeader("X-Gaming-Action"); action != "" {
		return action
	}

	// Determine from request path and method
	path := c.Request.URL.Path
	method := c.Request.Method

	// Gaming-specific action mapping
	switch {
	case strings.Contains(path, "/analytics/export"):
		return "analytics:export"
	case strings.Contains(path, "/team/") && method == "POST":
		return "team:management"
	case strings.Contains(path, "/subscription/"):
		return "subscription:change"
	case strings.Contains(path, "/account/delete"):
		return "account:delete"
	case strings.Contains(path, "/api/admin/"):
		return "admin:access"
	default:
		return "unknown:action"
	}
}

// getAnalyticsOperation determines analytics operation type
func (m *GamingMFAMiddleware) getAnalyticsOperation(c *gin.Context) string {
	path := c.Request.URL.Path

	switch {
	case strings.Contains(path, "/export"):
		return "export"
	case strings.Contains(path, "/advanced"):
		return "advanced"
	case strings.Contains(path, "/bulk"):
		return "bulk"
	case strings.Contains(path, "/raw"):
		return "raw"
	default:
		return "basic"
	}
}

// isHighPrivilegeAPIRequest checks if request requires high privilege
func (m *GamingMFAMiddleware) isHighPrivilegeAPIRequest(c *gin.Context) bool {
	path := c.Request.URL.Path
	method := c.Request.Method

	// High privilege API operations
	highPrivilegePaths := []string{
		"/api/admin/",
		"/api/users/",
		"/api/teams/manage",
		"/api/subscription/",
		"/api/billing/",
		"/api/export/bulk",
	}

	for _, privilegePath := range highPrivilegePaths {
		if strings.Contains(path, privilegePath) {
			return true
		}
	}

	// All DELETE operations on user data
	if method == "DELETE" && strings.Contains(path, "/api/") {
		return true
	}

	return false
}

// hasRecentMFAVerification checks for recent MFA verification in session
func (m *GamingMFAMiddleware) hasRecentMFAVerification(c *gin.Context) bool {
	// Check session for recent MFA verification
	// This would typically check a session store or JWT claims

	// Check for MFA session cookie
	if mfaSession, err := c.Cookie("herald_mfa_session"); err == nil && mfaSession != "" {
		// Verify MFA session token (simplified)
		return len(mfaSession) > 20
	}

	// Check for MFA token in Authorization header
	mfaToken := m.extractMFAToken(c)
	if mfaToken != "" {
		return m.mfaManager.verifyMFAToken(mfaToken, "session")
	}

	return false
}

// verifyRecoveryToken verifies recovery token for MFA operations
func (m *GamingMFAMiddleware) verifyRecoveryToken(userID, token string) bool {
	// In real implementation, this would verify:
	// - Password confirmation
	// - Email verification link
	// - SMS verification code
	// - Admin override token
	// - Emergency recovery code

	// Simplified verification
	return len(token) > 10 && strings.HasPrefix(token, "recovery_")
}

// Gaming MFA Context Helpers

// GetMFAVerified checks if current request has MFA verification
func GetMFAVerified(c *gin.Context) bool {
	if verified, exists := c.Get("mfa_verified"); exists {
		if verifiedBool, ok := verified.(bool); ok {
			return verifiedBool
		}
	}
	return false
}

// GetMFAAction gets the MFA action from context
func GetMFAAction(c *gin.Context) (string, bool) {
	if action, exists := c.Get("mfa_action"); exists {
		if actionString, ok := action.(string); ok {
			return actionString, true
		}
	}
	return "", false
}

// GetMFAConfigured checks if user has MFA configured
func GetMFAConfigured(c *gin.Context) bool {
	if configured, exists := c.Get("mfa_configured"); exists {
		if configuredBool, ok := configured.(bool); ok {
			return configuredBool
		}
	}
	return false
}

// GetMFAMethods gets available MFA methods from context
func GetMFAMethods(c *gin.Context) ([]string, bool) {
	if methods, exists := c.Get("mfa_methods"); exists {
		if methodsSlice, ok := methods.([]string); ok {
			return methodsSlice, true
		}
	}
	return []string{}, false
}

// GetMFAStatus gets MFA status from context
func GetMFAStatus(c *gin.Context) (*MFAStatus, bool) {
	if status, exists := c.Get("mfa_status"); exists {
		if statusStruct, ok := status.(*MFAStatus); ok {
			return statusStruct, true
		}
	}
	return nil, false
}
