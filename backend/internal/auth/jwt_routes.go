package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Herald.lol Gaming Analytics - JWT Management Routes
// Advanced JWT token management endpoints with short expiration

// SetupJWTManagementRoutes sets up JWT token management routes
func SetupJWTManagementRoutes(router *gin.Engine, jwtManager *GamingJWTManager, middleware *GamingAuthMiddleware) {
	// Create JWT management group
	jwt := router.Group("/auth/jwt")
	jwt.Use(middleware.CORSForGaming())
	jwt.Use(middleware.GamingSecurityHeaders())
	jwt.Use(middleware.RateLimitByGamingTier())
	
	// Token generation and refresh endpoints
	setupTokenGenerationRoutes(jwt, jwtManager, middleware)
	
	// Token validation and introspection
	setupTokenValidationRoutes(jwt, jwtManager, middleware)
	
	// Token revocation and blacklisting
	setupTokenRevocationRoutes(jwt, jwtManager, middleware)
	
	// Gaming-specific token endpoints
	setupGamingTokenRoutes(jwt, jwtManager, middleware)
	
	// Token management utilities
	setupTokenUtilityRoutes(jwt, jwtManager, middleware)
}

// Token generation and refresh endpoints
func setupTokenGenerationRoutes(jwt *gin.RouterGroup, jwtManager *GamingJWTManager, middleware *GamingAuthMiddleware) {
	tokens := jwt.Group("/tokens")
	
	// Generate new gaming token pair (for authenticated users)
	tokens.POST("/generate", middleware.RequireGamingAuth(), func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			return
		}
		
		// Extract device info from request
		deviceInfo := &DeviceInfo{
			Platform:      c.GetHeader("X-Gaming-Platform"),
			UserAgent:     c.GetHeader("User-Agent"),
			IPAddress:     c.ClientIP(),
			GamingClient:  c.GetHeader("X-Gaming-Client"),
			ClientVersion: c.GetHeader("X-Gaming-Client-Version"),
		}
		
		// Extract gaming context from request
		gamingContext := &GamingContext{
			CurrentRegion:     c.GetHeader("X-Gaming-Region"),
			SessionType:       c.GetHeader("X-Gaming-Session-Type"),
			GamingPreferences: make(map[string]string),
		}
		
		// Parse gaming preferences from form data
		if preferredAnalytics := c.PostForm("preferred_analytics"); preferredAnalytics != "" {
			// Would parse comma-separated list
			gamingContext.PreferredAnalytics = []string{preferredAnalytics}
		}
		
		// Generate token pair
		tokenPair, err := jwtManager.GenerateGamingTokenPair(c.Request.Context(), user, deviceInfo, gamingContext)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate gaming token pair",
				"details": err.Error(),
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"tokens": tokenPair,
			"gaming_platform": "herald-lol",
			"user_id": user.ID,
			"subscription_tier": user.GamingProfile.SubscriptionTier,
		})
	})
	
	// Refresh gaming access token
	tokens.POST("/refresh", func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
			DeviceInfo   struct {
				Platform      string `json:"platform"`
				UserAgent     string `json:"user_agent"`
				IPAddress     string `json:"ip_address"`
				GamingClient  string `json:"gaming_client"`
				ClientVersion string `json:"client_version"`
				Fingerprint   string `json:"fingerprint"`
			} `json:"device_info"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid gaming token refresh request",
				"details": err.Error(),
			})
			return
		}
		
		// Convert request device info
		deviceInfo := &DeviceInfo{
			Platform:      req.DeviceInfo.Platform,
			UserAgent:     req.DeviceInfo.UserAgent,
			IPAddress:     req.DeviceInfo.IPAddress,
			GamingClient:  req.DeviceInfo.GamingClient,
			ClientVersion: req.DeviceInfo.ClientVersion,
			Fingerprint:   req.DeviceInfo.Fingerprint,
		}
		
		// If device info not provided, extract from headers
		if deviceInfo.IPAddress == "" {
			deviceInfo.IPAddress = c.ClientIP()
		}
		if deviceInfo.UserAgent == "" {
			deviceInfo.UserAgent = c.GetHeader("User-Agent")
		}
		if deviceInfo.Platform == "" {
			deviceInfo.Platform = c.GetHeader("X-Gaming-Platform")
		}
		
		// Refresh token
		tokenPair, err := jwtManager.RefreshGamingToken(c.Request.Context(), req.RefreshToken, deviceInfo)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Failed to refresh gaming token",
				"details": err.Error(),
				"action": "re_authenticate",
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"tokens": tokenPair,
			"gaming_platform": "herald-lol",
		})
	})
	
	// Batch refresh multiple tokens
	tokens.POST("/refresh/batch", middleware.RequireGamingAuth(), func(c *gin.Context) {
		var req struct {
			RefreshTokens []string `json:"refresh_tokens" binding:"required"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid batch refresh request",
			})
			return
		}
		
		// Extract device info
		deviceInfo := &DeviceInfo{
			Platform:     c.GetHeader("X-Gaming-Platform"),
			UserAgent:    c.GetHeader("User-Agent"),
			IPAddress:    c.ClientIP(),
			GamingClient: c.GetHeader("X-Gaming-Client"),
		}
		
		var results []gin.H
		var successCount int
		
		for _, refreshToken := range req.RefreshTokens {
			tokenPair, err := jwtManager.RefreshGamingToken(c.Request.Context(), refreshToken, deviceInfo)
			if err != nil {
				results = append(results, gin.H{
					"refresh_token": refreshToken[:20] + "...", // Partial token for identification
					"success": false,
					"error": err.Error(),
				})
			} else {
				results = append(results, gin.H{
					"refresh_token": refreshToken[:20] + "...",
					"success": true,
					"tokens": tokenPair,
				})
				successCount++
			}
		}
		
		c.JSON(http.StatusOK, gin.H{
			"batch_refresh": true,
			"total_tokens": len(req.RefreshTokens),
			"successful_refreshes": successCount,
			"results": results,
			"gaming_platform": "herald-lol",
		})
	})
}

// Token validation and introspection
func setupTokenValidationRoutes(jwt *gin.RouterGroup, jwtManager *GamingJWTManager, middleware *GamingAuthMiddleware) {
	validation := jwt.Group("/validate")
	
	// Validate gaming access token
	validation.POST("/access", func(c *gin.Context) {
		var req struct {
			Token string `json:"token" binding:"required"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			// Also check Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				req.Token = authHeader[7:]
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Gaming token required",
				})
				return
			}
		}
		
		// Parse token using OAuth config (reuse existing method)
		claims, err := jwtManager.config.parseJWT(req.Token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"valid": false,
				"error": "Invalid gaming token",
				"gaming_platform": "herald-lol",
			})
			return
		}
		
		// Check if token is blacklisted
		if jwtManager.config.EnableBlacklist && jwtManager.blacklistStore != nil {
			blacklisted, err := jwtManager.blacklistStore.IsTokenBlacklisted(c.Request.Context(), claims.TokenID)
			if err == nil && blacklisted {
				c.JSON(http.StatusUnauthorized, gin.H{
					"valid": false,
					"error": "Gaming token is blacklisted",
					"gaming_platform": "herald-lol",
				})
				return
			}
		}
		
		c.JSON(http.StatusOK, gin.H{
			"valid": true,
			"token_id": claims.TokenID,
			"user_id": claims.UserID,
			"email": claims.Email,
			"name": claims.Name,
			"provider": claims.Provider,
			"subscription_tier": claims.SubscriptionTier,
			"permissions": claims.GamingPermissions,
			"token_type": claims.TokenType,
			"expires_at": claims.ExpiresAt.Time,
			"issued_at": claims.IssuedAt.Time,
			"gaming_platform": "herald-lol",
		})
	})
	
	// Introspect gaming token (detailed information)
	validation.POST("/introspect", middleware.RequireGamingAuth(), func(c *gin.Context) {
		claims, exists := GetGamingClaims(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming claims not found",
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"active": true,
			"token_id": claims.TokenID,
			"token_type": claims.TokenType,
			"token_version": claims.TokenVersion,
			"user_id": claims.UserID,
			"email": claims.Email,
			"name": claims.Name,
			"provider": claims.Provider,
			"provider_id": claims.ProviderID,
			"subscription_tier": claims.SubscriptionTier,
			"subscription_expiry": claims.SubscriptionExpiry,
			"gaming_permissions": claims.GamingPermissions,
			"gaming_region": claims.GamingRegion,
			"last_game_activity": claims.LastGameActivity,
			"preferred_analytics": claims.PreferredAnalytics,
			"device_fingerprint": claims.DeviceFingerprint,
			"session_id": claims.SessionID,
			"ip_address": claims.IPAddress,
			"gaming_metadata": claims.GamingMetadata,
			"issued_at": claims.IssuedAt.Time,
			"expires_at": claims.ExpiresAt.Time,
			"not_before": claims.NotBefore.Time,
			"issuer": claims.Issuer,
			"audience": claims.Audience,
			"gaming_platform": "herald-lol",
		})
	})
	
	// Check token expiration
	validation.GET("/expiration/:tokenId", middleware.RequireGamingAuth(), func(c *gin.Context) {
		tokenID := c.Param("tokenId")
		claims, exists := GetGamingClaims(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming claims not found",
			})
			return
		}
		
		// Only allow checking own tokens or with admin permission
		if claims.TokenID != tokenID && !HasGamingPermission(c, "token:admin") {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Not authorized to check this gaming token",
			})
			return
		}
		
		now := time.Now()
		expiresAt := claims.ExpiresAt.Time
		
		c.JSON(http.StatusOK, gin.H{
			"token_id": tokenID,
			"expires_at": expiresAt,
			"is_expired": now.After(expiresAt),
			"time_until_expiry": int(time.Until(expiresAt).Seconds()),
			"gaming_platform": "herald-lol",
		})
	})
}

// Token revocation and blacklisting
func setupTokenRevocationRoutes(jwt *gin.RouterGroup, jwtManager *GamingJWTManager, middleware *GamingAuthMiddleware) {
	revocation := jwt.Group("/revoke")
	revocation.Use(middleware.RequireGamingAuth())
	
	// Revoke current gaming token (logout)
	revocation.POST("/current", func(c *gin.Context) {
		claims, exists := GetGamingClaims(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming claims not found",
			})
			return
		}
		
		// Blacklist current access token
		if jwtManager.config.EnableBlacklist && jwtManager.blacklistStore != nil {
			err := jwtManager.blacklistStore.BlacklistToken(c.Request.Context(), claims.TokenID, claims.ExpiresAt.Time)
			if err != nil {
				// Log error but continue
			}
		}
		
		// Clear gaming cookies
		domain := c.GetHeader("X-Gaming-Domain")
		if domain == "" {
			domain = ".herald.lol"
		}
		
		c.SetCookie("herald_access_token", "", -1, "/", domain, true, true)
		c.SetCookie("herald_refresh_token", "", -1, "/", domain, true, true)
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Gaming token revoked successfully",
			"token_id": claims.TokenID,
			"gaming_platform": "herald-lol",
		})
	})
	
	// Revoke specific gaming refresh token
	revocation.POST("/refresh/:tokenId", func(c *gin.Context) {
		tokenID := c.Param("tokenId")
		userID, exists := GetGamingUserID(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user ID not found",
			})
			return
		}
		
		// Get refresh token to verify ownership
		refreshToken, err := jwtManager.refreshTokenStore.GetRefreshToken(c.Request.Context(), tokenID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Gaming refresh token not found",
			})
			return
		}
		
		// Verify ownership
		if refreshToken.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Not authorized to revoke this gaming token",
			})
			return
		}
		
		// Revoke refresh token
		if err := jwtManager.refreshTokenStore.RevokeRefreshToken(c.Request.Context(), tokenID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to revoke gaming refresh token",
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Gaming refresh token revoked successfully",
			"token_id": tokenID,
			"gaming_platform": "herald-lol",
		})
	})
	
	// Revoke all gaming tokens for current user
	revocation.POST("/all", func(c *gin.Context) {
		userID, exists := GetGamingUserID(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user ID not found",
			})
			return
		}
		
		// Revoke all refresh tokens
		if err := jwtManager.refreshTokenStore.RevokeAllUserTokens(c.Request.Context(), userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to revoke all gaming tokens",
			})
			return
		}
		
		// Clear cookies
		domain := c.GetHeader("X-Gaming-Domain")
		if domain == "" {
			domain = ".herald.lol"
		}
		
		c.SetCookie("herald_access_token", "", -1, "/", domain, true, true)
		c.SetCookie("herald_refresh_token", "", -1, "/", domain, true, true)
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "All gaming tokens revoked successfully",
			"user_id": userID,
			"gaming_platform": "herald-lol",
		})
	})
}

// Gaming-specific token endpoints
func setupGamingTokenRoutes(jwt *gin.RouterGroup, jwtManager *GamingJWTManager, middleware *GamingAuthMiddleware) {
	gaming := jwt.Group("/gaming")
	gaming.Use(middleware.RequireGamingAuth())
	
	// Generate gaming analytics token
	gaming.POST("/analytics-token", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			return
		}
		
		var req struct {
			AnalyticsScope []string `json:"analytics_scope"`
			TTL           int      `json:"ttl_minutes,omitempty"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid gaming analytics token request",
			})
			return
		}
		
		// Validate analytics scope against user permissions
		userPermissions, _ := c.Get("gaming_permissions")
		if !validateAnalyticsScope(req.AnalyticsScope, userPermissions.([]string)) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions for requested gaming analytics scope",
				"requested_scope": req.AnalyticsScope,
				"user_permissions": userPermissions,
			})
			return
		}
		
		// Generate analytics token
		token, err := jwtManager.GenerateGamingAnalyticsToken(c.Request.Context(), user, req.AnalyticsScope)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate gaming analytics token",
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"analytics_token": token,
			"analytics_scope": req.AnalyticsScope,
			"expires_in": int(jwtManager.config.AnalyticsTokenTTL.Seconds()),
			"gaming_platform": "herald-lol",
		})
	})
	
	// List active gaming sessions
	gaming.GET("/sessions", func(c *gin.Context) {
		userID, exists := GetGamingUserID(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user ID not found",
			})
			return
		}
		
		// Get active refresh tokens (representing sessions)
		refreshTokens, err := jwtManager.refreshTokenStore.GetUserRefreshTokens(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get gaming sessions",
			})
			return
		}
		
		var sessions []gin.H
		for _, token := range refreshTokens {
			if !token.IsRevoked && time.Now().Before(token.ExpiresAt) {
				session := gin.H{
					"session_id": token.ID,
					"device_info": token.DeviceInfo,
					"gaming_context": token.GamingContext,
					"issued_at": token.IssuedAt,
					"expires_at": token.ExpiresAt,
					"last_used_at": token.LastUsedAt,
					"usage_count": token.UsageCount,
				}
				sessions = append(sessions, session)
			}
		}
		
		c.JSON(http.StatusOK, gin.H{
			"active_sessions": sessions,
			"session_count": len(sessions),
			"gaming_platform": "herald-lol",
		})
	})
}

// Token management utilities
func setupTokenUtilityRoutes(jwt *gin.RouterGroup, jwtManager *GamingJWTManager, middleware *GamingAuthMiddleware) {
	utils := jwt.Group("/utils")
	
	// Gaming JWT configuration info
	utils.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"access_token_ttl": jwtManager.config.AccessTokenTTL.String(),
			"refresh_token_ttl": jwtManager.config.RefreshTokenTTL.String(),
			"gaming_token_ttl": jwtManager.config.GamingTokenTTL.String(),
			"analytics_token_ttl": jwtManager.config.AnalyticsTokenTTL.String(),
			"token_rotation_enabled": jwtManager.config.EnableTokenRotation,
			"blacklist_enabled": jwtManager.config.EnableBlacklist,
			"max_refresh_attempts": jwtManager.config.MaxRefreshAttempts,
			"token_versioning": jwtManager.config.TokenVersioning,
			"issuer": jwtManager.config.Issuer,
			"audience": jwtManager.config.Audience,
			"gaming_platform": "herald-lol",
		})
	})
	
	// Gaming token statistics (admin only)
	utils.GET("/stats", middleware.RequireGamingPermission("token:admin"), func(c *gin.Context) {
		// This would return token usage statistics
		c.JSON(http.StatusOK, gin.H{
			"stats": gin.H{
				"total_active_tokens": 0,      // Would be calculated
				"tokens_issued_today": 0,      // Would be calculated  
				"tokens_refreshed_today": 0,   // Would be calculated
				"blacklisted_tokens": 0,       // Would be calculated
			},
			"gaming_platform": "herald-lol",
			"status": "not_implemented", // TODO: Implement token statistics
		})
	})
	
	// JWT health check
	utils.GET("/health", func(c *gin.Context) {
		// Test token generation
		testClaims := &EnhancedGamingJWTClaims{
			UserID:    "test",
			TokenID:   "health_check",
			TokenType: "test",
		}
		
		testToken := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
		_, err := testToken.SignedString(jwtManager.config.AccessTokenSecret)
		
		healthy := err == nil
		
		c.JSON(http.StatusOK, gin.H{
			"service": "herald-jwt-manager",
			"status": map[string]bool{
				"healthy": healthy,
				"jwt_signing": err == nil,
			},
			"config": gin.H{
				"access_token_ttl": jwtManager.config.AccessTokenTTL.String(),
				"refresh_token_ttl": jwtManager.config.RefreshTokenTTL.String(),
			},
			"gaming_platform": "herald-lol",
			"timestamp": time.Now(),
		})
	})
}

// Helper functions

// validateAnalyticsScope validates requested analytics scope against user permissions
func validateAnalyticsScope(requestedScope, userPermissions []string) bool {
	for _, scope := range requestedScope {
		hasPermission := false
		
		// Check if user has permission for this scope
		for _, permission := range userPermissions {
			if permission == scope || permission == "analytics:advanced" {
				hasPermission = true
				break
			}
		}
		
		if !hasPermission {
			return false
		}
	}
	
	return true
}