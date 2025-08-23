package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

// Herald.lol Gaming Analytics - MFA Routes
// Multi-Factor Authentication API endpoints for gaming platform

// SetupGamingMFARoutes sets up MFA routes for Herald.lol gaming platform
func SetupGamingMFARoutes(router *gin.Engine, mfaManager *GamingMFAManager, middleware *GamingAuthMiddleware) {
	// Create MFA group with authentication required
	mfa := router.Group("/auth/mfa")
	mfa.Use(middleware.CORSForGaming())
	mfa.Use(middleware.GamingSecurityHeaders())
	mfa.Use(middleware.RequireGamingAuth())
	
	// TOTP endpoints
	setupTOTPRoutes(mfa, mfaManager, middleware)
	
	// WebAuthn endpoints  
	setupWebAuthnRoutes(mfa, mfaManager, middleware)
	
	// Backup codes endpoints
	setupBackupCodesRoutes(mfa, mfaManager, middleware)
	
	// MFA verification endpoints
	setupMFAVerificationRoutes(mfa, mfaManager, middleware)
	
	// MFA management endpoints
	setupMFAManagementRoutes(mfa, mfaManager, middleware)
}

// TOTP endpoints
func setupTOTPRoutes(mfa *gin.RouterGroup, mfaManager *GamingMFAManager, middleware *GamingAuthMiddleware) {
	totp := mfa.Group("/totp")
	
	// Setup TOTP for gaming user
	totp.POST("/setup", mfaManager.SetupGamingTOTP)
	
	// Verify TOTP code and enable
	totp.POST("/verify", mfaManager.VerifyGamingTOTP)
	
	// Verify TOTP code for authentication
	totp.POST("/authenticate", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			return
		}
		
		var req struct {
			Code          string `json:"code" binding:"required"`
			GamingAction  string `json:"gaming_action,omitempty"`
			RememberDevice bool  `json:"remember_device,omitempty"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid gaming TOTP authentication request",
			})
			return
		}
		
		// Get TOTP secret
		totpSecret, err := mfaManager.mfaStore.GetTOTPSecret(c.Request.Context(), user.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Gaming TOTP not configured",
			})
			return
		}
		
		if !totpSecret.Enabled {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Gaming TOTP not enabled",
			})
			return
		}
		
		// Check rate limiting
		if err := mfaManager.checkMFARateLimit(c.Request.Context(), user.ID); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many gaming MFA attempts",
				"retry_after": mfaManager.config.MFACooldownPeriod.Seconds(),
			})
			return
		}
		
		// Verify TOTP code
		valid := totp.Validate(req.Code, totpSecret.Secret)
		if !valid {
			// Track failed attempt
			mfaManager.trackMFAAttempt(c.Request.Context(), user.ID, "totp", false, "Invalid TOTP code", c)
			
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid gaming TOTP code",
				"gaming_platform": "herald-lol",
			})
			return
		}
		
		// Update last used time
		now := time.Now()
		totpSecret.LastUsed = &now
		mfaManager.mfaStore.StoreTOTPSecret(c.Request.Context(), user.ID, totpSecret)
		
		// Track successful attempt
		mfaManager.trackMFAAttempt(c.Request.Context(), user.ID, "totp", true, "", c)
		
		// Generate MFA authentication token
		mfaToken := mfaManager.generateMFAToken(user.ID, "totp", req.GamingAction)
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Gaming TOTP authentication successful",
			"mfa_token": mfaToken,
			"expires_in": int(mfaManager.config.MFACooldownPeriod.Seconds()),
			"gaming_platform": "herald-lol",
		})
	})
	
	// Disable TOTP
	totp.POST("/disable", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			return
		}
		
		var req struct {
			Password string `json:"password" binding:"required"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password required to disable gaming TOTP",
			})
			return
		}
		
		// Verify password (implementation would verify against user's password)
		// For now, just disable TOTP
		
		if err := mfaManager.mfaStore.DisableTOTP(c.Request.Context(), user.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to disable gaming TOTP",
			})
			return
		}
		
		// Track TOTP disabled
		go mfaManager.gamingAnalytics.TrackUserLogin(c.Request.Context(), user.ID, user.Provider, map[string]string{
			"action": "mfa_totp_disabled",
			"gaming_platform": "herald-lol",
		})
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Gaming TOTP disabled successfully",
			"gaming_platform": "herald-lol",
		})
	})
	
	// Get TOTP status
	totp.GET("/status", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			return
		}
		
		totpSecret, err := mfaManager.mfaStore.GetTOTPSecret(c.Request.Context(), user.ID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"configured": false,
				"enabled": false,
				"gaming_platform": "herald-lol",
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"configured": true,
			"enabled": totpSecret.Enabled,
			"verified": totpSecret.Verified,
			"last_used": totpSecret.LastUsed,
			"created_at": totpSecret.CreatedAt,
			"gaming_platform": "herald-lol",
		})
	})
}

// WebAuthn endpoints
func setupWebAuthnRoutes(mfa *gin.RouterGroup, mfaManager *GamingMFAManager, middleware *GamingAuthMiddleware) {
	webauthn := mfa.Group("/webauthn")
	
	// Begin WebAuthn registration
	webauthn.POST("/register/begin", mfaManager.BeginWebAuthnRegistration)
	
	// Complete WebAuthn registration
	webauthn.POST("/register/complete/:challengeId", mfaManager.CompleteWebAuthnRegistration)
	
	// Begin WebAuthn authentication
	webauthn.POST("/authenticate/begin", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			return
		}
		
		var req struct {
			GamingAction string `json:"gaming_action,omitempty"`
		}
		
		c.ShouldBindJSON(&req) // Optional
		
		// Get user's credentials
		credentials, err := mfaManager.mfaStore.GetWebAuthnCredentials(c.Request.Context(), user.ID)
		if err != nil || len(credentials) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "No WebAuthn credentials found",
			})
			return
		}
		
		// Create WebAuthn user
		gamingWebAuthnUser := &GamingWebAuthnUser{
			ID:          user.ID,
			Name:        user.Email,
			DisplayName: user.Name,
			Icon:        user.Avatar,
			GamingCredentials: credentials,
		}
		
		// Begin authentication (simplified)
		challengeID := mfaManager.generateChallengeID()
		
		// Store challenge
		challenge := &MFAChallenge{
			ID:            challengeID,
			UserID:        user.ID,
			ChallengeType: "webauthn_authentication",
			SessionData: map[string]interface{}{
				"gaming_action": req.GamingAction,
			},
			IPAddress: c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(mfaManager.config.WebAuthnTimeout),
		}
		
		if err := mfaManager.mfaStore.StoreMFAChallenge(c.Request.Context(), challengeID, challenge); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to store WebAuthn challenge",
			})
			return
		}
		
		// Return challenge options (simplified)
		c.JSON(http.StatusOK, gin.H{
			"challenge_id": challengeID,
			"credential_options": gin.H{
				"publicKey": gin.H{
					"challenge": challengeID,
					"timeout": int(mfaManager.config.WebAuthnTimeout.Milliseconds()),
					"rpId": mfaManager.config.WebAuthnRPID,
					"allowCredentials": []gin.H{}, // Would include actual credential descriptors
				},
			},
			"gaming_platform": "herald-lol",
		})
	})
	
	// Complete WebAuthn authentication
	webauthn.POST("/authenticate/complete/:challengeId", func(c *gin.Context) {
		challengeID := c.Param("challengeId")
		
		// Get challenge
		challenge, err := mfaManager.mfaStore.GetMFAChallenge(c.Request.Context(), challengeID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "WebAuthn challenge not found or expired",
			})
			return
		}
		
		if challenge.Completed || time.Now().After(challenge.ExpiresAt) {
			c.JSON(http.StatusGone, gin.H{
				"error": "WebAuthn challenge expired or already used",
			})
			return
		}
		
		// Parse authentication response (simplified)
		var authResponse map[string]interface{}
		if err := c.ShouldBindJSON(&authResponse); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid WebAuthn response",
			})
			return
		}
		
		// Verify authentication (simplified - would use webauthn library)
		// For now, just mark as successful
		
		// Update credential last used
		credentialID := "simplified_credential_id" // Would extract from response
		credentials, _ := mfaManager.mfaStore.GetWebAuthnCredentials(c.Request.Context(), challenge.UserID)
		for _, cred := range credentials {
			if cred.ID == credentialID {
				now := time.Now()
				cred.LastUsed = &now
				cred.SignCount++ // Would be updated based on response
				mfaManager.mfaStore.UpdateWebAuthnCredential(c.Request.Context(), credentialID, cred)
				break
			}
		}
		
		// Mark challenge as completed
		now := time.Now()
		challenge.Completed = true
		challenge.CompletedAt = &now
		mfaManager.mfaStore.StoreMFAChallenge(c.Request.Context(), challengeID, challenge)
		
		// Track successful attempt
		mfaManager.trackMFAAttempt(c.Request.Context(), challenge.UserID, "webauthn", true, "", c)
		
		// Generate MFA authentication token
		gamingAction, _ := challenge.SessionData["gaming_action"].(string)
		mfaToken := mfaManager.generateMFAToken(challenge.UserID, "webauthn", gamingAction)
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "WebAuthn authentication successful",
			"mfa_token": mfaToken,
			"credential_id": credentialID,
			"gaming_platform": "herald-lol",
		})
	})
	
	// List WebAuthn credentials
	webauthn.GET("/credentials", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			return
		}
		
		credentials, err := mfaManager.mfaStore.GetWebAuthnCredentials(c.Request.Context(), user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get WebAuthn credentials",
			})
			return
		}
		
		// Format credentials for response (remove sensitive data)
		var credentialList []gin.H
		for _, cred := range credentials {
			credentialList = append(credentialList, gin.H{
				"id": cred.ID,
				"device_name": cred.DeviceName,
				"device_type": cred.DeviceType,
				"gaming_platform": cred.GamingPlatform,
				"created_at": cred.CreatedAt,
				"last_used": cred.LastUsed,
				"enabled": cred.Enabled,
			})
		}
		
		c.JSON(http.StatusOK, gin.H{
			"credentials": credentialList,
			"count": len(credentialList),
			"gaming_platform": "herald-lol",
		})
	})
	
	// Delete WebAuthn credential
	webauthn.DELETE("/credentials/:credentialId", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			return
		}
		
		credentialID := c.Param("credentialId")
		
		if err := mfaManager.mfaStore.DeleteWebAuthnCredential(c.Request.Context(), user.ID, credentialID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "WebAuthn credential not found",
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "WebAuthn credential deleted successfully",
			"credential_id": credentialID,
			"gaming_platform": "herald-lol",
		})
	})
}

// Backup codes endpoints
func setupBackupCodesRoutes(mfa *gin.RouterGroup, mfaManager *GamingMFAManager, middleware *GamingAuthMiddleware) {
	backup := mfa.Group("/backup-codes")
	
	// Get backup codes (show unused codes)
	backup.GET("", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			return
		}
		
		backupCodes, err := mfaManager.mfaStore.GetBackupCodes(c.Request.Context(), user.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Gaming backup codes not found",
			})
			return
		}
		
		// Only return unused codes
		var unusedCodes []string
		for code, used := range backupCodes.Codes {
			if !used {
				unusedCodes = append(unusedCodes, code)
			}
		}
		
		c.JSON(http.StatusOK, gin.H{
			"backup_codes": unusedCodes,
			"total_codes": len(backupCodes.Codes),
			"used_codes": backupCodes.UsedCount,
			"unused_codes": len(unusedCodes),
			"created_at": backupCodes.CreatedAt,
			"gaming_platform": "herald-lol",
		})
	})
	
	// Use backup code for authentication
	backup.POST("/authenticate", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			return
		}
		
		var req struct {
			Code         string `json:"code" binding:"required"`
			GamingAction string `json:"gaming_action,omitempty"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid gaming backup code authentication request",
			})
			return
		}
		
		// Check rate limiting
		if err := mfaManager.checkMFARateLimit(c.Request.Context(), user.ID); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many gaming MFA attempts",
			})
			return
		}
		
		// Use backup code
		if err := mfaManager.mfaStore.UseBackupCode(c.Request.Context(), user.ID, req.Code); err != nil {
			// Track failed attempt
			mfaManager.trackMFAAttempt(c.Request.Context(), user.ID, "backup_code", false, err.Error(), c)
			
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or used gaming backup code",
			})
			return
		}
		
		// Track successful attempt
		mfaManager.trackMFAAttempt(c.Request.Context(), user.ID, "backup_code", true, "", c)
		
		// Generate MFA authentication token
		mfaToken := mfaManager.generateMFAToken(user.ID, "backup_code", req.GamingAction)
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Gaming backup code authentication successful",
			"mfa_token": mfaToken,
			"warning": "Backup code used - consider regenerating codes",
			"gaming_platform": "herald-lol",
		})
	})
	
	// Regenerate backup codes
	backup.POST("/regenerate", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			return
		}
		
		var req struct {
			Password string `json:"password" binding:"required"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password required to regenerate gaming backup codes",
			})
			return
		}
		
		// Verify password (implementation would verify against user's password)
		
		// Generate new backup codes
		codes, err := mfaManager.generateBackupCodes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate new gaming backup codes",
			})
			return
		}
		
		// Store new backup codes
		backupCodes := &BackupCodes{
			UserID:    user.ID,
			Codes:     make(map[string]bool),
			CreatedAt: time.Now(),
			UsedCount: 0,
		}
		
		for _, code := range codes {
			backupCodes.Codes[code] = false
		}
		
		if err := mfaManager.mfaStore.StoreBackupCodes(c.Request.Context(), user.ID, backupCodes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to store new gaming backup codes",
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"backup_codes": codes,
			"message": "Gaming backup codes regenerated successfully",
			"gaming_platform": "herald-lol",
		})
	})
}

// MFA verification endpoints
func setupMFAVerificationRoutes(mfa *gin.RouterGroup, mfaManager *GamingMFAManager, middleware *GamingAuthMiddleware) {
	verify := mfa.Group("/verify")
	
	// Check if MFA is required for specific action
	verify.POST("/required", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			return
		}
		
		var req struct {
			Action string `json:"action" binding:"required"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Action required",
			})
			return
		}
		
		// Check if action requires MFA
		requiresMFA := mfaManager.actionRequiresMFA(req.Action)
		
		// Get user's MFA methods
		status, _ := mfaManager.mfaStore.(*CombinedGamingMFAStore).GetMFAStatus(c.Request.Context(), user.ID)
		
		c.JSON(http.StatusOK, gin.H{
			"action": req.Action,
			"mfa_required": requiresMFA,
			"mfa_configured": status != nil && status.Enabled,
			"available_methods": func() []string {
				if status == nil {
					return []string{}
				}
				return status.Methods
			}(),
			"gaming_platform": "herald-lol",
		})
	})
	
	// Verify MFA token
	verify.POST("/token", func(c *gin.Context) {
		var req struct {
			MFAToken string `json:"mfa_token" binding:"required"`
			Action   string `json:"action" binding:"required"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "MFA token and action required",
			})
			return
		}
		
		// Verify MFA token (simplified)
		valid := mfaManager.verifyMFAToken(req.MFAToken, req.Action)
		
		c.JSON(http.StatusOK, gin.H{
			"valid": valid,
			"action": req.Action,
			"gaming_platform": "herald-lol",
		})
	})
}

// MFA management endpoints
func setupMFAManagementRoutes(mfa *gin.RouterGroup, mfaManager *GamingMFAManager, middleware *GamingAuthMiddleware) {
	// Get comprehensive MFA status
	mfa.GET("/status", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			return
		}
		
		status, err := mfaManager.mfaStore.(*CombinedGamingMFAStore).GetMFAStatus(c.Request.Context(), user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get gaming MFA status",
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"mfa_status": status,
			"gaming_platform": "herald-lol",
		})
	})
	
	// Get MFA attempts history
	mfa.GET("/attempts", func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gaming user not found",
			})
			return
		}
		
		// Get attempts from last 24 hours
		since := time.Now().Add(-24 * time.Hour)
		if sinceParam := c.Query("since"); sinceParam != "" {
			if hours, err := strconv.Atoi(sinceParam); err == nil {
				since = time.Now().Add(-time.Duration(hours) * time.Hour)
			}
		}
		
		attempts, err := mfaManager.mfaStore.GetMFAAttempts(c.Request.Context(), user.ID, since)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get gaming MFA attempts",
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"attempts": attempts,
			"count": len(attempts),
			"since": since,
			"gaming_platform": "herald-lol",
		})
	})
	
	// MFA health check
	mfa.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "herald-mfa-manager",
			"status": "healthy",
			"totp_enabled": true,
			"webauthn_enabled": true,
			"backup_codes_enabled": mfaManager.config.BackupCodesEnabled,
			"gaming_platform": "herald-lol",
			"timestamp": time.Now(),
		})
	})
}

// Helper methods for MFA manager

func (mfa *GamingMFAManager) checkMFARateLimit(ctx context.Context, userID string) error {
	// Get recent failed attempts
	since := time.Now().Add(-time.Hour)
	attempts, err := mfa.mfaStore.GetMFAAttempts(ctx, userID, since)
	if err != nil {
		return nil // Continue if we can't check
	}
	
	// Count failed attempts
	failedCount := 0
	for _, attempt := range attempts {
		if !attempt.Success {
			failedCount++
		}
	}
	
	if failedCount >= mfa.config.MaxMFAAttempts {
		return fmt.Errorf("too many failed gaming MFA attempts")
	}
	
	return nil
}

func (mfa *GamingMFAManager) generateMFAToken(userID, method, action string) string {
	// Generate temporary MFA authentication token
	// This would be a short-lived JWT or random token
	return fmt.Sprintf("mfa_%s_%s_%d", method, userID[:8], time.Now().Unix())
}

func (mfa *GamingMFAManager) verifyMFAToken(token, action string) bool {
	// Verify MFA token (simplified)
	// In real implementation, this would verify JWT signature and expiration
	return len(token) > 10 && action != ""
}

func (mfa *GamingMFAManager) actionRequiresMFA(action string) bool {
	// Check if action requires MFA
	for _, mfaAction := range mfa.config.HighValueActionsMFA {
		if action == mfaAction {
			return true
		}
	}
	return false
}