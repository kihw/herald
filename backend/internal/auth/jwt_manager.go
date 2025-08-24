package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Herald.lol Gaming Analytics - JWT Token Manager
// Advanced JWT token management with short expiration, rotation, and security

// GamingJWTManager handles JWT token lifecycle for Herald.lol gaming platform
type GamingJWTManager struct {
	config            *GamingJWTConfig
	refreshTokenStore RefreshTokenStore
	blacklistStore    TokenBlacklistStore
	rotationStore     TokenRotationStore
	userStore         *DatabaseUserStore
	gamingAnalytics   GamingAnalyticsService
}

// GamingJWTConfig holds JWT configuration for gaming platform
type GamingJWTConfig struct {
	AccessTokenSecret  []byte
	RefreshTokenSecret []byte
	AccessTokenTTL     time.Duration // Short expiration for security
	RefreshTokenTTL    time.Duration
	RefreshRotationTTL time.Duration

	// Gaming-specific configuration
	GamingTokenTTL    time.Duration // Ultra-short for gaming analytics
	AnalyticsTokenTTL time.Duration // For analytics-specific operations

	// Security configuration
	EnableTokenRotation bool
	EnableBlacklist     bool
	MaxRefreshAttempts  int
	TokenVersioning     bool

	// Performance configuration
	TokenCacheEnabled bool
	TokenCacheTTL     time.Duration

	Issuer   string
	Audience []string
}

// RefreshTokenStore interface for refresh token management
type RefreshTokenStore interface {
	StoreRefreshToken(ctx context.Context, tokenID string, refreshToken *GamingRefreshToken) error
	GetRefreshToken(ctx context.Context, tokenID string) (*GamingRefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID string) error
	CleanupExpiredTokens(ctx context.Context) error
	GetUserRefreshTokens(ctx context.Context, userID string) ([]*GamingRefreshToken, error)
	RevokeAllUserTokens(ctx context.Context, userID string) error
}

// TokenBlacklistStore interface for token blacklisting
type TokenBlacklistStore interface {
	BlacklistToken(ctx context.Context, tokenID string, expiresAt time.Time) error
	IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error)
	CleanupExpiredBlacklist(ctx context.Context) error
}

// TokenRotationStore interface for token rotation tracking
type TokenRotationStore interface {
	TrackRotation(ctx context.Context, oldTokenID, newTokenID string, userID string) error
	GetRotationChain(ctx context.Context, tokenID string) ([]*TokenRotation, error)
	IsRotationValid(ctx context.Context, tokenID string) (bool, error)
}

// GamingRefreshToken represents a gaming refresh token
type GamingRefreshToken struct {
	ID               string            `json:"id"`
	UserID           string            `json:"user_id"`
	TokenHash        string            `json:"token_hash"`
	DeviceInfo       *DeviceInfo       `json:"device_info"`
	GamingContext    *GamingContext    `json:"gaming_context"`
	Metadata         map[string]string `json:"metadata"`
	IssuedAt         time.Time         `json:"issued_at"`
	ExpiresAt        time.Time         `json:"expires_at"`
	LastUsedAt       *time.Time        `json:"last_used_at,omitempty"`
	UsageCount       int64             `json:"usage_count"`
	IsRevoked        bool              `json:"is_revoked"`
	RevokedAt        *time.Time        `json:"revoked_at,omitempty"`
	RevokedReason    string            `json:"revoked_reason,omitempty"`
	RotationCount    int               `json:"rotation_count"`
	RotationParentID string            `json:"rotation_parent_id,omitempty"`
	Version          int               `json:"version"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

// DeviceInfo represents device/client information
type DeviceInfo struct {
	DeviceID      string `json:"device_id"`
	Platform      string `json:"platform"`
	UserAgent     string `json:"user_agent"`
	IPAddress     string `json:"ip_address"`
	GamingClient  string `json:"gaming_client,omitempty"`
	ClientVersion string `json:"client_version,omitempty"`
	Fingerprint   string `json:"fingerprint,omitempty"`
}

// GamingContext represents gaming-specific context
type GamingContext struct {
	LastGameSession    *time.Time        `json:"last_game_session,omitempty"`
	CurrentRegion      string            `json:"current_region,omitempty"`
	PreferredAnalytics []string          `json:"preferred_analytics,omitempty"`
	GamingPreferences  map[string]string `json:"gaming_preferences,omitempty"`
	SessionType        string            `json:"session_type,omitempty"` // web, mobile, api, analytics
}

// TokenRotation represents token rotation history
type TokenRotation struct {
	ID         string    `json:"id"`
	OldTokenID string    `json:"old_token_id"`
	NewTokenID string    `json:"new_token_id"`
	UserID     string    `json:"user_id"`
	RotatedAt  time.Time `json:"rotated_at"`
	Reason     string    `json:"reason"`
	IPAddress  string    `json:"ip_address"`
}

// Enhanced JWT Claims for gaming platform
type EnhancedGamingJWTClaims struct {
	// Standard claims
	jwt.RegisteredClaims

	// Gaming user claims
	UserID     string `json:"uid"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Provider   string `json:"provider"`
	ProviderID string `json:"provider_id"`

	// Gaming subscription and permissions
	SubscriptionTier   string     `json:"sub_tier"`
	GamingPermissions  []string   `json:"gaming_perms"`
	SubscriptionExpiry *time.Time `json:"sub_exp,omitempty"`

	// Gaming context
	GamingRegion       string     `json:"gaming_region,omitempty"`
	LastGameActivity   *time.Time `json:"last_game,omitempty"`
	PreferredAnalytics []string   `json:"pref_analytics,omitempty"`

	// Token management
	TokenID      string `json:"jti"`        // JWT ID for tracking
	TokenType    string `json:"token_type"` // access, refresh, gaming, analytics
	TokenVersion int    `json:"token_version"`

	// Security
	DeviceFingerprint string `json:"device_fp,omitempty"`
	SessionID         string `json:"sid,omitempty"`
	IPAddress         string `json:"ip,omitempty"`

	// Gaming-specific metadata
	GamingMetadata map[string]string `json:"gaming_meta,omitempty"`
}

// NewGamingJWTManager creates new JWT manager for gaming platform
func NewGamingJWTManager(
	config *GamingJWTConfig,
	refreshTokenStore RefreshTokenStore,
	blacklistStore TokenBlacklistStore,
	rotationStore TokenRotationStore,
	userStore *DatabaseUserStore,
	gamingAnalytics GamingAnalyticsService,
) *GamingJWTManager {
	// Set gaming-specific defaults
	if config.AccessTokenTTL == 0 {
		config.AccessTokenTTL = 15 * time.Minute // Short expiration for security
	}
	if config.RefreshTokenTTL == 0 {
		config.RefreshTokenTTL = 7 * 24 * time.Hour // 7 days
	}
	if config.GamingTokenTTL == 0 {
		config.GamingTokenTTL = 5 * time.Minute // Ultra-short for gaming operations
	}
	if config.AnalyticsTokenTTL == 0 {
		config.AnalyticsTokenTTL = 30 * time.Minute // Analytics-specific tokens
	}
	if config.RefreshRotationTTL == 0 {
		config.RefreshRotationTTL = 30 * 24 * time.Hour // 30 days
	}

	return &GamingJWTManager{
		config:            config,
		refreshTokenStore: refreshTokenStore,
		blacklistStore:    blacklistStore,
		rotationStore:     rotationStore,
		userStore:         userStore,
		gamingAnalytics:   gamingAnalytics,
	}
}

// GenerateGamingTokenPair generates access and refresh token pair for gaming
func (jm *GamingJWTManager) GenerateGamingTokenPair(ctx context.Context, user *GamingUserInfo, deviceInfo *DeviceInfo, gamingContext *GamingContext) (*TokenPair, error) {
	now := time.Now()

	// Generate token IDs
	accessTokenID := jm.generateTokenID()
	refreshTokenID := jm.generateTokenID()

	// Create enhanced claims for access token
	accessClaims := &EnhancedGamingJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        accessTokenID,
			Subject:   user.ID,
			Issuer:    jm.config.Issuer,
			Audience:  jm.config.Audience,
			ExpiresAt: jwt.NewNumericDate(now.Add(jm.config.AccessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UserID:            user.ID,
		Email:             user.Email,
		Name:              user.Name,
		Provider:          string(user.Provider),
		ProviderID:        user.ProviderID,
		SubscriptionTier:  user.GamingProfile.SubscriptionTier,
		GamingPermissions: jm.getGamingPermissions(user.GamingProfile.SubscriptionTier),
		TokenID:           accessTokenID,
		TokenType:         "access",
		TokenVersion:      1,
		GamingMetadata:    user.Metadata,
	}

	// Add gaming context if provided
	if gamingContext != nil {
		accessClaims.GamingRegion = gamingContext.CurrentRegion
		accessClaims.LastGameActivity = gamingContext.LastGameSession
		accessClaims.PreferredAnalytics = gamingContext.PreferredAnalytics
	}

	// Add device context if provided
	if deviceInfo != nil {
		accessClaims.DeviceFingerprint = deviceInfo.Fingerprint
		accessClaims.IPAddress = deviceInfo.IPAddress
	}

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jm.config.AccessTokenSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign gaming access token: %w", err)
	}

	// Create refresh token claims
	refreshClaims := &EnhancedGamingJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        refreshTokenID,
			Subject:   user.ID,
			Issuer:    jm.config.Issuer,
			Audience:  jm.config.Audience,
			ExpiresAt: jwt.NewNumericDate(now.Add(jm.config.RefreshTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UserID:       user.ID,
		TokenID:      refreshTokenID,
		TokenType:    "refresh",
		TokenVersion: 1,
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jm.config.RefreshTokenSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign gaming refresh token: %w", err)
	}

	// Store refresh token
	gamingRefreshToken := &GamingRefreshToken{
		ID:            refreshTokenID,
		UserID:        user.ID,
		TokenHash:     jm.hashToken(refreshTokenString),
		DeviceInfo:    deviceInfo,
		GamingContext: gamingContext,
		Metadata:      user.Metadata,
		IssuedAt:      now,
		ExpiresAt:     now.Add(jm.config.RefreshTokenTTL),
		UsageCount:    0,
		IsRevoked:     false,
		Version:       1,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := jm.refreshTokenStore.StoreRefreshToken(ctx, refreshTokenID, gamingRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to store gaming refresh token: %w", err)
	}

	// Track token generation
	go jm.gamingAnalytics.TrackUserLogin(ctx, user.ID, user.Provider, map[string]string{
		"action":            "token_generated",
		"token_type":        "access_refresh_pair",
		"access_token_ttl":  jm.config.AccessTokenTTL.String(),
		"refresh_token_ttl": jm.config.RefreshTokenTTL.String(),
	})

	return &TokenPair{
		AccessToken:      accessTokenString,
		RefreshToken:     refreshTokenString,
		TokenType:        "Bearer",
		ExpiresIn:        int(jm.config.AccessTokenTTL.Seconds()),
		RefreshExpiresIn: int(jm.config.RefreshTokenTTL.Seconds()),
		AccessTokenID:    accessTokenID,
		RefreshTokenID:   refreshTokenID,
	}, nil
}

// RefreshGamingToken refreshes access token using refresh token
func (jm *GamingJWTManager) RefreshGamingToken(ctx context.Context, refreshTokenString string, deviceInfo *DeviceInfo) (*TokenPair, error) {
	// Parse refresh token
	refreshClaims, err := jm.parseRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid gaming refresh token: %w", err)
	}

	// Check if token is blacklisted
	if jm.config.EnableBlacklist {
		blacklisted, err := jm.blacklistStore.IsTokenBlacklisted(ctx, refreshClaims.TokenID)
		if err != nil {
			return nil, fmt.Errorf("failed to check token blacklist: %w", err)
		}
		if blacklisted {
			return nil, fmt.Errorf("gaming refresh token is blacklisted")
		}
	}

	// Get stored refresh token
	storedToken, err := jm.refreshTokenStore.GetRefreshToken(ctx, refreshClaims.TokenID)
	if err != nil {
		return nil, fmt.Errorf("gaming refresh token not found: %w", err)
	}

	// Validate stored token
	if storedToken.IsRevoked {
		return nil, fmt.Errorf("gaming refresh token is revoked")
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return nil, fmt.Errorf("gaming refresh token expired")
	}

	// Verify token hash
	if storedToken.TokenHash != jm.hashToken(refreshTokenString) {
		return nil, fmt.Errorf("gaming refresh token hash mismatch")
	}

	// Get user info
	user, err := jm.userStore.GetUserByProviderID(ctx, OAuthProvider(refreshClaims.Provider), refreshClaims.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("gaming user not found: %w", err)
	}

	// Check refresh attempt limits
	if jm.config.MaxRefreshAttempts > 0 && storedToken.UsageCount >= int64(jm.config.MaxRefreshAttempts) {
		// Revoke token due to excessive usage
		jm.revokeRefreshToken(ctx, refreshClaims.TokenID, "excessive_usage")
		return nil, fmt.Errorf("gaming refresh token usage limit exceeded")
	}

	now := time.Now()
	var newTokenPair *TokenPair

	if jm.config.EnableTokenRotation {
		// Token rotation: generate new refresh token
		newTokenPair, err = jm.GenerateGamingTokenPair(ctx, user, deviceInfo, storedToken.GamingContext)
		if err != nil {
			return nil, fmt.Errorf("failed to generate new gaming token pair: %w", err)
		}

		// Track rotation
		if jm.rotationStore != nil {
			jm.rotationStore.TrackRotation(ctx, refreshClaims.TokenID, newTokenPair.RefreshTokenID, user.ID)
		}

		// Revoke old refresh token
		jm.revokeRefreshToken(ctx, refreshClaims.TokenID, "rotated")

		// Track rotation event
		go jm.gamingAnalytics.TrackUserLogin(ctx, user.ID, user.Provider, map[string]string{
			"action":       "token_rotated",
			"old_token_id": refreshClaims.TokenID,
			"new_token_id": newTokenPair.RefreshTokenID,
			"usage_count":  fmt.Sprintf("%d", storedToken.UsageCount),
		})
	} else {
		// No rotation: generate new access token only
		accessTokenID := jm.generateTokenID()

		// Create new access token claims
		accessClaims := &EnhancedGamingJWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        accessTokenID,
				Subject:   user.ID,
				Issuer:    jm.config.Issuer,
				Audience:  jm.config.Audience,
				ExpiresAt: jwt.NewNumericDate(now.Add(jm.config.AccessTokenTTL)),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now),
			},
			UserID:            user.ID,
			Email:             user.Email,
			Name:              user.Name,
			Provider:          string(user.Provider),
			ProviderID:        user.ProviderID,
			SubscriptionTier:  user.GamingProfile.SubscriptionTier,
			GamingPermissions: jm.getGamingPermissions(user.GamingProfile.SubscriptionTier),
			TokenID:           accessTokenID,
			TokenType:         "access",
			TokenVersion:      storedToken.Version,
			GamingMetadata:    user.Metadata,
		}

		// Generate new access token
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
		accessTokenString, err := accessToken.SignedString(jm.config.AccessTokenSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to sign new gaming access token: %w", err)
		}

		// Update refresh token usage
		storedToken.LastUsedAt = &now
		storedToken.UsageCount++
		storedToken.UpdatedAt = now

		if err := jm.refreshTokenStore.StoreRefreshToken(ctx, refreshClaims.TokenID, storedToken); err != nil {
			return nil, fmt.Errorf("failed to update gaming refresh token: %w", err)
		}

		newTokenPair = &TokenPair{
			AccessToken:      accessTokenString,
			RefreshToken:     refreshTokenString, // Reuse existing refresh token
			TokenType:        "Bearer",
			ExpiresIn:        int(jm.config.AccessTokenTTL.Seconds()),
			RefreshExpiresIn: int(time.Until(storedToken.ExpiresAt).Seconds()),
			AccessTokenID:    accessTokenID,
			RefreshTokenID:   refreshClaims.TokenID,
		}

		// Track refresh event
		go jm.gamingAnalytics.TrackUserLogin(ctx, user.ID, user.Provider, map[string]string{
			"action":      "token_refreshed",
			"token_id":    refreshClaims.TokenID,
			"usage_count": fmt.Sprintf("%d", storedToken.UsageCount),
		})
	}

	return newTokenPair, nil
}

// GenerateGamingAnalyticsToken generates short-lived token for gaming analytics
func (jm *GamingJWTManager) GenerateGamingAnalyticsToken(ctx context.Context, user *GamingUserInfo, analyticsScope []string) (string, error) {
	now := time.Now()
	tokenID := jm.generateTokenID()

	claims := &EnhancedGamingJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			Subject:   user.ID,
			Issuer:    jm.config.Issuer,
			Audience:  []string{"herald-gaming-analytics"},
			ExpiresAt: jwt.NewNumericDate(now.Add(jm.config.AnalyticsTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UserID:             user.ID,
		SubscriptionTier:   user.GamingProfile.SubscriptionTier,
		GamingPermissions:  analyticsScope,
		TokenID:            tokenID,
		TokenType:          "analytics",
		TokenVersion:       1,
		PreferredAnalytics: analyticsScope,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jm.config.AccessTokenSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign gaming analytics token: %w", err)
	}

	// Track analytics token generation
	go jm.gamingAnalytics.TrackUserLogin(ctx, user.ID, user.Provider, map[string]string{
		"action":          "analytics_token_generated",
		"analytics_scope": fmt.Sprintf("%v", analyticsScope),
		"token_ttl":       jm.config.AnalyticsTokenTTL.String(),
	})

	return tokenString, nil
}

// TokenPair represents access and refresh token pair
type TokenPair struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	AccessTokenID    string `json:"access_token_id,omitempty"`
	RefreshTokenID   string `json:"refresh_token_id,omitempty"`
}

// Helper methods

func (jm *GamingJWTManager) generateTokenID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func (jm *GamingJWTManager) hashToken(token string) string {
	// In production, use crypto/sha256 for proper hashing
	return fmt.Sprintf("hash_%s", token[:20]) // Simplified for example
}

func (jm *GamingJWTManager) parseRefreshToken(tokenString string) (*EnhancedGamingJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &EnhancedGamingJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected gaming refresh token signing method: %v", token.Header["alg"])
		}
		return jm.config.RefreshTokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*EnhancedGamingJWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid gaming refresh token")
}

func (jm *GamingJWTManager) revokeRefreshToken(ctx context.Context, tokenID, reason string) {
	// Implementation to revoke refresh token
	go func() {
		if err := jm.refreshTokenStore.RevokeRefreshToken(ctx, tokenID); err != nil {
			// Log error
		}
	}()
}

func (jm *GamingJWTManager) getGamingPermissions(subscriptionTier string) []string {
	switch subscriptionTier {
	case "enterprise":
		return []string{"analytics:advanced", "api:unlimited", "coaching:premium", "team:management", "export:all"}
	case "pro":
		return []string{"analytics:advanced", "api:extended", "coaching:premium", "team:basic", "export:basic"}
	case "premium":
		return []string{"analytics:advanced", "api:standard", "coaching:basic", "export:basic"}
	case "free":
		fallthrough
	default:
		return []string{"analytics:basic", "api:limited"}
	}
}
