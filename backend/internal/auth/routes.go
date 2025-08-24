package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Herald.lol Gaming Analytics - Authentication Routes
// RESTful API endpoints for gaming OAuth and authentication

// GamingAuthRoutes sets up authentication routes for Herald.lol gaming platform
func SetupGamingAuthRoutes(router *gin.Engine, authConfig *GamingOAuthConfig, middleware *GamingAuthMiddleware) {
	// Create auth group with gaming middleware
	auth := router.Group("/auth")
	auth.Use(middleware.CORSForGaming())
	auth.Use(middleware.GamingSecurityHeaders())

	// OAuth 2.0/OpenID Connect routes
	setupOAuthRoutes(auth, authConfig)

	// JWT token management routes
	setupTokenRoutes(auth, authConfig, middleware)

	// Gaming user profile routes
	setupProfileRoutes(auth, authConfig, middleware)

	// Gaming session management routes
	setupSessionRoutes(auth, authConfig, middleware)

	// Gaming authentication utilities
	setupUtilityRoutes(auth, authConfig, middleware)
}

// OAuth 2.0/OpenID Connect routes
func setupOAuthRoutes(auth *gin.RouterGroup, authConfig *GamingOAuthConfig) {
	oauth := auth.Group("/oauth")

	// Start OAuth flow for gaming providers
	oauth.GET("/:provider", authConfig.StartGamingOAuth)

	// Handle OAuth callbacks for gaming providers
	oauth.GET("/:provider/callback", authConfig.HandleGamingOAuthCallback)

	// List supported gaming OAuth providers
	oauth.GET("/providers", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"supported_providers": []gin.H{
				{
					"name":            "google",
					"display_name":    "Google",
					"icon":            "https://cdn.herald.lol/icons/google.svg",
					"description":     "Sign in with Google",
					"gaming_features": []string{"profile_sync"},
				},
				{
					"name":            "discord",
					"display_name":    "Discord",
					"icon":            "https://cdn.herald.lol/icons/discord.svg",
					"description":     "Sign in with Discord",
					"gaming_features": []string{"rich_presence", "social_features"},
				},
				{
					"name":            "twitch",
					"display_name":    "Twitch",
					"icon":            "https://cdn.herald.lol/icons/twitch.svg",
					"description":     "Sign in with Twitch",
					"gaming_features": []string{"streaming_integration", "clips_analysis"},
				},
				{
					"name":            "riot",
					"display_name":    "Riot Games",
					"icon":            "https://cdn.herald.lol/icons/riot.svg",
					"description":     "Sign in with Riot Games",
					"gaming_features": []string{"lol_integration", "tft_integration", "direct_match_import"},
				},
				{
					"name":            "github",
					"display_name":    "GitHub",
					"icon":            "https://cdn.herald.lol/icons/github.svg",
					"description":     "Sign in with GitHub",
					"gaming_features": []string{"developer_features"},
				},
			},
			"gaming_platform": "herald-lol",
		})
	})
}

// JWT token management routes
func setupTokenRoutes(auth *gin.RouterGroup, authConfig *GamingOAuthConfig, middleware *GamingAuthMiddleware) {
	tokens := auth.Group("/tokens")

	// Refresh gaming access token
	tokens.POST("/refresh", authConfig.RefreshGamingToken)

	// Validate gaming token
	tokens.POST("/validate", func(c *gin.Context) {
		token := c.PostForm("token")
		if token == "" {
			token = c.GetHeader("Authorization")
			if token != "" && len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}
		}

		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":           "Gaming token required",
				"gaming_platform": "herald-lol",
			})
			return
		}

		// Parse and validate gaming token
		claims, err := authConfig.parseJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"valid":           false,
				"error":           "Invalid gaming token",
				"gaming_platform": "herald-lol",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"valid":             true,
			"user_id":           claims.UserID,
			"email":             claims.Email,
			"name":              claims.Name,
			"provider":          claims.Provider,
			"subscription_tier": claims.SubscriptionTier,
			"permissions":       claims.GamingPermissions,
			"expires_at":        claims.ExpiresAt.Time,
			"gaming_platform":   "herald-lol",
		})
	})

	// Revoke gaming token (logout)
	tokens.POST("/revoke", middleware.RequireGamingAuth(), func(c *gin.Context) {
		// Get session token if available
		sessionToken := c.GetHeader("X-Session-Token")

		if sessionToken != "" {
			// Invalidate session in database
			userStore := middleware.userStore
			userStore.InvalidateGamingSession(c.Request.Context(), sessionToken)
		}

		// Clear gaming cookies
		c.SetCookie("herald_access_token", "", -1, "/", "", true, true)
		c.SetCookie("herald_refresh_token", "", -1, "/", "", true, true)

		c.JSON(http.StatusOK, gin.H{
			"success":         true,
			"message":         "Gaming token revoked successfully",
			"gaming_platform": "herald-lol",
		})
	})

	// Get gaming token info
	tokens.GET("/info", middleware.RequireGamingAuth(), func(c *gin.Context) {
		claims, exists := GetGamingClaims(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming claims not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id":           claims.UserID,
			"email":             claims.Email,
			"name":              claims.Name,
			"provider":          claims.Provider,
			"subscription_tier": claims.SubscriptionTier,
			"permissions":       claims.GamingPermissions,
			"issued_at":         claims.IssuedAt.Time,
			"expires_at":        claims.ExpiresAt.Time,
			"gaming_platform":   "herald-lol",
		})
	})
}

// Gaming user profile routes
func setupProfileRoutes(auth *gin.RouterGroup, authConfig *GamingOAuthConfig, middleware *GamingAuthMiddleware) {
	profile := auth.Group("/profile")
	profile.Use(middleware.RequireGamingAuth())

	// Get gaming user profile
	profile.GET("", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found in context",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":            user,
			"gaming_platform": "herald-lol",
		})
	})

	// Update gaming user profile
	profile.PUT("", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found in context",
			})
			return
		}

		var updateData struct {
			Name     string `json:"name,omitempty"`
			Username string `json:"username,omitempty"`
			Avatar   string `json:"avatar,omitempty"`
		}

		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid gaming profile data",
				"details": err.Error(),
			})
			return
		}

		// Update gaming user profile
		if updateData.Name != "" {
			user.Name = updateData.Name
		}
		if updateData.Username != "" {
			user.Username = updateData.Username
		}
		if updateData.Avatar != "" {
			user.Avatar = updateData.Avatar
		}

		// Save to database
		userStore := middleware.userStore
		if err := userStore.UpdateUser(c.Request.Context(), user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update gaming user profile",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":         true,
			"user":            user,
			"gaming_platform": "herald-lol",
		})
	})

	// Get gaming profile
	profile.GET("/gaming", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found in context",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"gaming_profile":  user.GamingProfile,
			"gaming_platform": "herald-lol",
		})
	})

	// Update gaming profile
	profile.PUT("/gaming", func(c *gin.Context) {
		userID, exists := GetGamingUserID(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user ID not found in context",
			})
			return
		}

		var gamingProfile GamingProfile
		if err := c.ShouldBindJSON(&gamingProfile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid gaming profile data",
				"details": err.Error(),
			})
			return
		}

		// Validate gaming profile data
		if err := validateGamingProfile(&gamingProfile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid gaming profile",
				"details": err.Error(),
			})
			return
		}

		// Update gaming profile
		userStore := middleware.userStore
		if err := userStore.UpdateGamingProfile(c.Request.Context(), userID, &gamingProfile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update gaming profile",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":         true,
			"gaming_profile":  gamingProfile,
			"gaming_platform": "herald-lol",
		})
	})

	// Link gaming accounts
	profile.POST("/link/:provider", func(c *gin.Context) {
		provider := c.Param("provider")

		// This would initiate linking another gaming account
		// Implementation would be similar to OAuth flow but for account linking

		c.JSON(http.StatusOK, gin.H{
			"message":         "Gaming account linking initiated",
			"provider":        provider,
			"gaming_platform": "herald-lol",
			"status":          "not_implemented", // TODO: Implement account linking
		})
	})
}

// Gaming session management routes
func setupSessionRoutes(auth *gin.RouterGroup, authConfig *GamingOAuthConfig, middleware *GamingAuthMiddleware) {
	sessions := auth.Group("/sessions")
	sessions.Use(middleware.RequireGamingAuth())

	// Get active gaming sessions
	sessions.GET("", func(c *gin.Context) {
		// This would return active gaming sessions for the user
		c.JSON(http.StatusOK, gin.H{
			"sessions":        []gin.H{}, // TODO: Implement session listing
			"gaming_platform": "herald-lol",
		})
	})

	// Terminate gaming session
	sessions.DELETE("/:sessionId", func(c *gin.Context) {
		sessionId := c.Param("sessionId")

		// This would terminate a specific gaming session
		userStore := middleware.userStore
		if err := userStore.InvalidateGamingSession(c.Request.Context(), sessionId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to terminate gaming session",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":         true,
			"message":         "Gaming session terminated",
			"gaming_platform": "herald-lol",
		})
	})

	// Terminate all gaming sessions
	sessions.DELETE("", func(c *gin.Context) {
		// This would terminate all gaming sessions for the user
		c.JSON(http.StatusOK, gin.H{
			"success":         true,
			"message":         "All gaming sessions terminated",
			"gaming_platform": "herald-lol",
			"status":          "not_implemented", // TODO: Implement bulk session termination
		})
	})
}

// Gaming authentication utilities
func setupUtilityRoutes(auth *gin.RouterGroup, authConfig *GamingOAuthConfig, middleware *GamingAuthMiddleware) {
	utils := auth.Group("/utils")

	// Get gaming authentication status
	utils.GET("/status", middleware.OptionalGamingAuth(), func(c *gin.Context) {
		authenticated := IsGamingAuthenticated(c)

		response := gin.H{
			"authenticated":   authenticated,
			"gaming_platform": "herald-lol",
		}

		if authenticated {
			user, _ := GetGamingUser(c)
			tier, _ := GetGamingSubscriptionTier(c)

			response["user_id"] = user.ID
			response["email"] = user.Email
			response["name"] = user.Name
			response["subscription_tier"] = tier
			response["provider"] = string(user.Provider)
		}

		c.JSON(http.StatusOK, response)
	})

	// Gaming platform health check
	utils.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":          "healthy",
			"service":         "herald-gaming-auth",
			"version":         "1.0.0",
			"gaming_platform": "herald-lol",
			"timestamp":       c.GetHeader("Date"),
		})
	})

	// Gaming subscription tiers info
	utils.GET("/tiers", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"tiers": []gin.H{
				{
					"name":         "free",
					"display_name": "Free",
					"price":        "$0/month",
					"features": []string{
						"Basic analytics",
						"100 API requests/minute",
						"Match history (30 days)",
					},
					"gaming_permissions": []string{"analytics:basic", "api:limited"},
				},
				{
					"name":         "premium",
					"display_name": "Premium",
					"price":        "$9.99/month",
					"features": []string{
						"Advanced analytics",
						"500 API requests/minute",
						"Match history (90 days)",
						"Team composition optimizer",
						"Counter-pick analysis",
					},
					"gaming_permissions": []string{"analytics:advanced", "api:standard", "coaching:basic"},
				},
				{
					"name":         "pro",
					"display_name": "Pro",
					"price":        "$29.99/month",
					"features": []string{
						"Professional analytics",
						"2000 API requests/minute",
						"Unlimited match history",
						"AI coaching insights",
						"Team management tools",
						"Custom analytics exports",
					},
					"gaming_permissions": []string{"analytics:advanced", "api:extended", "coaching:premium", "team:basic"},
				},
				{
					"name":         "enterprise",
					"display_name": "Enterprise",
					"price":        "Custom pricing",
					"features": []string{
						"Enterprise analytics",
						"10000 API requests/minute",
						"Full data access",
						"Custom integrations",
						"Dedicated support",
						"Team management suite",
					},
					"gaming_permissions": []string{"analytics:advanced", "api:unlimited", "coaching:premium", "team:management", "export:all"},
				},
			},
			"gaming_platform": "herald-lol",
		})
	})

	// Gaming permissions info
	utils.GET("/permissions", middleware.RequireGamingAuth(), func(c *gin.Context) {
		permissions, _ := c.Get("gaming_permissions")
		tier, _ := GetGamingSubscriptionTier(c)

		c.JSON(http.StatusOK, gin.H{
			"current_permissions": permissions,
			"subscription_tier":   tier,
			"permission_descriptions": gin.H{
				"analytics:basic":    "Access basic gaming analytics",
				"analytics:advanced": "Access advanced gaming analytics and insights",
				"api:limited":        "Limited API access (100 req/min)",
				"api:standard":       "Standard API access (500 req/min)",
				"api:extended":       "Extended API access (2000 req/min)",
				"api:unlimited":      "Unlimited API access (10000 req/min)",
				"coaching:basic":     "Basic coaching recommendations",
				"coaching:premium":   "Premium AI coaching insights",
				"team:basic":         "Basic team management features",
				"team:management":    "Full team management suite",
				"export:basic":       "Basic data export capabilities",
				"export:all":         "Full data export and custom integrations",
			},
			"gaming_platform": "herald-lol",
		})
	})
}

// Helper functions

// validateGamingProfile validates gaming profile data
func validateGamingProfile(profile *GamingProfile) error {
	// Validate region
	validRegions := []string{"na1", "euw1", "eun1", "kr", "br1", "la1", "la2", "oc1", "ru", "tr1", "jp1"}
	if profile.Region != "" {
		validRegion := false
		for _, region := range validRegions {
			if profile.Region == region {
				validRegion = true
				break
			}
		}
		if !validRegion {
			return gin.H{"field": "region", "message": "Invalid region"}.(error)
		}
	}

	// Validate subscription tier
	validTiers := []string{"free", "premium", "pro", "enterprise"}
	if profile.SubscriptionTier != "" {
		validTier := false
		for _, tier := range validTiers {
			if profile.SubscriptionTier == tier {
				validTier = true
				break
			}
		}
		if !validTier {
			return gin.H{"field": "subscription_tier", "message": "Invalid subscription tier"}.(error)
		}
	}

	return nil
}
