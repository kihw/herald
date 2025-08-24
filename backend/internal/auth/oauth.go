package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// Herald.lol Gaming Analytics - OAuth 2.0/OpenID Connect Implementation
// Comprehensive authentication for gaming platform users

// OAuthProvider represents supported OAuth providers
type OAuthProvider string

const (
	ProviderGoogle  OAuthProvider = "google"
	ProviderDiscord OAuthProvider = "discord"
	ProviderTwitch  OAuthProvider = "twitch"
	ProviderRiot    OAuthProvider = "riot"
	ProviderGitHub  OAuthProvider = "github"
)

// GamingUserInfo represents user information from OAuth providers
type GamingUserInfo struct {
	ID            string            `json:"id"`
	Email         string            `json:"email"`
	Name          string            `json:"name"`
	Username      string            `json:"username,omitempty"`
	Avatar        string            `json:"avatar,omitempty"`
	Provider      OAuthProvider     `json:"provider"`
	ProviderID    string            `json:"provider_id"`
	GamingProfile *GamingProfile    `json:"gaming_profile,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// GamingProfile represents gaming-specific user profile
type GamingProfile struct {
	SummonerName     string            `json:"summoner_name,omitempty"`
	Region           string            `json:"region,omitempty"`
	Rank             string            `json:"rank,omitempty"`
	MainChampions    []string          `json:"main_champions,omitempty"`
	PreferredRoles   []string          `json:"preferred_roles,omitempty"`
	DiscordUsername  string            `json:"discord_username,omitempty"`
	TwitchUsername   string            `json:"twitch_username,omitempty"`
	GamingGoals      []string          `json:"gaming_goals,omitempty"`
	SubscriptionTier string            `json:"subscription_tier,omitempty"`
	Preferences      map[string]string `json:"preferences,omitempty"`
}

// OAuthState represents OAuth state for CSRF protection
type OAuthState struct {
	State       string            `json:"state"`
	Provider    OAuthProvider     `json:"provider"`
	RedirectURL string            `json:"redirect_url,omitempty"`
	GamingData  map[string]string `json:"gaming_data,omitempty"`
	ExpiresAt   time.Time         `json:"expires_at"`
}

// GamingOAuthConfig holds OAuth configuration for gaming platform
type GamingOAuthConfig struct {
	GoogleConfig  *oauth2.Config
	DiscordConfig *oauth2.Config
	TwitchConfig  *oauth2.Config
	RiotConfig    *oauth2.Config
	GitHubConfig  *oauth2.Config

	JWTSecret         []byte
	JWTExpiration     time.Duration
	RefreshExpiration time.Duration

	DB              *gorm.DB
	StateStore      StateStore
	UserStore       UserStore
	GamingAnalytics GamingAnalyticsService
}

// StateStore interface for OAuth state management
type StateStore interface {
	StoreState(ctx context.Context, state string, oauthState *OAuthState) error
	GetState(ctx context.Context, state string) (*OAuthState, error)
	DeleteState(ctx context.Context, state string) error
	CleanupExpiredStates(ctx context.Context) error
}

// UserStore interface for user management
type UserStore interface {
	GetUserByProviderID(ctx context.Context, provider OAuthProvider, providerID string) (*GamingUserInfo, error)
	GetUserByEmail(ctx context.Context, email string) (*GamingUserInfo, error)
	CreateUser(ctx context.Context, user *GamingUserInfo) error
	UpdateUser(ctx context.Context, user *GamingUserInfo) error
	UpdateGamingProfile(ctx context.Context, userID string, profile *GamingProfile) error
}

// GamingAnalyticsService interface for gaming analytics integration
type GamingAnalyticsService interface {
	TrackUserLogin(ctx context.Context, userID string, provider OAuthProvider, metadata map[string]string)
	TrackUserRegistration(ctx context.Context, userID string, provider OAuthProvider, profile *GamingProfile)
}

// NewGamingOAuthConfig creates new OAuth configuration for gaming platform
func NewGamingOAuthConfig(
	db *gorm.DB,
	stateStore StateStore,
	userStore UserStore,
	gamingAnalytics GamingAnalyticsService,
	jwtSecret []byte,
) *GamingOAuthConfig {
	return &GamingOAuthConfig{
		DB:                db,
		StateStore:        stateStore,
		UserStore:         userStore,
		GamingAnalytics:   gamingAnalytics,
		JWTSecret:         jwtSecret,
		JWTExpiration:     15 * time.Minute,   // Gaming session duration
		RefreshExpiration: 7 * 24 * time.Hour, // Gaming refresh token duration

		GoogleConfig: &oauth2.Config{
			ClientID:     getEnvVar("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnvVar("GOOGLE_CLIENT_SECRET", ""),
			Endpoint:     google.Endpoint,
			RedirectURL:  getEnvVar("GOOGLE_REDIRECT_URL", ""),
			Scopes:       []string{"openid", "profile", "email"},
		},

		DiscordConfig: &oauth2.Config{
			ClientID:     getEnvVar("DISCORD_CLIENT_ID", ""),
			ClientSecret: getEnvVar("DISCORD_CLIENT_SECRET", ""),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/api/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
			RedirectURL: getEnvVar("DISCORD_REDIRECT_URL", ""),
			Scopes:      []string{"identify", "email"},
		},

		TwitchConfig: &oauth2.Config{
			ClientID:     getEnvVar("TWITCH_CLIENT_ID", ""),
			ClientSecret: getEnvVar("TWITCH_CLIENT_SECRET", ""),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://id.twitch.tv/oauth2/authorize",
				TokenURL: "https://id.twitch.tv/oauth2/token",
			},
			RedirectURL: getEnvVar("TWITCH_REDIRECT_URL", ""),
			Scopes:      []string{"openid", "user:read:email"},
		},

		RiotConfig: &oauth2.Config{
			ClientID:     getEnvVar("RIOT_CLIENT_ID", ""),
			ClientSecret: getEnvVar("RIOT_CLIENT_SECRET", ""),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://auth.riotgames.com/oauth2/authorize",
				TokenURL: "https://auth.riotgames.com/token",
			},
			RedirectURL: getEnvVar("RIOT_REDIRECT_URL", ""),
			Scopes:      []string{"openid", "cpid", "ppid"},
		},

		GitHubConfig: &oauth2.Config{
			ClientID:     getEnvVar("GITHUB_CLIENT_ID", ""),
			ClientSecret: getEnvVar("GITHUB_CLIENT_SECRET", ""),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
			RedirectURL: getEnvVar("GITHUB_REDIRECT_URL", ""),
			Scopes:      []string{"user:email", "read:user"},
		},
	}
}

// StartGamingOAuth initiates OAuth flow for gaming platform
func (oauth *GamingOAuthConfig) StartGamingOAuth(c *gin.Context) {
	provider := OAuthProvider(c.Param("provider"))

	// Validate gaming OAuth provider
	if !oauth.isValidProvider(provider) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":               "Invalid gaming OAuth provider",
			"supported_providers": []string{"google", "discord", "twitch", "riot", "github"},
		})
		return
	}

	// Generate secure state for CSRF protection
	state, err := oauth.generateSecureState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate OAuth state",
		})
		return
	}

	// Store OAuth state with gaming metadata
	oauthState := &OAuthState{
		State:       state,
		Provider:    provider,
		RedirectURL: c.Query("redirect_url"),
		GamingData: map[string]string{
			"client_ip":    c.ClientIP(),
			"user_agent":   c.GetHeader("User-Agent"),
			"gaming_flow":  c.Query("gaming_flow"),
			"utm_source":   c.Query("utm_source"),
			"utm_campaign": c.Query("utm_campaign"),
		},
		ExpiresAt: time.Now().Add(10 * time.Minute), // Gaming OAuth timeout
	}

	ctx := context.Background()
	if err := oauth.StateStore.StoreState(ctx, state, oauthState); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store OAuth state",
		})
		return
	}

	// Get OAuth config for gaming provider
	config := oauth.getProviderConfig(provider)
	if config == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gaming provider configuration not found",
		})
		return
	}

	// Generate authorization URL for gaming OAuth
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	// Track gaming OAuth initiation
	go oauth.GamingAnalytics.TrackUserLogin(ctx, "", provider, map[string]string{
		"action":      "oauth_start",
		"provider":    string(provider),
		"client_ip":   c.ClientIP(),
		"user_agent":  c.GetHeader("User-Agent"),
		"gaming_flow": c.Query("gaming_flow"),
	})

	c.JSON(http.StatusOK, gin.H{
		"auth_url":        authURL,
		"state":           state,
		"provider":        provider,
		"gaming_platform": "herald-lol",
		"expires_in":      600, // 10 minutes
	})
}

// HandleGamingOAuthCallback handles OAuth callback for gaming platform
func (oauth *GamingOAuthConfig) HandleGamingOAuthCallback(c *gin.Context) {
	provider := OAuthProvider(c.Param("provider"))
	state := c.Query("state")
	code := c.Query("code")
	errorParam := c.Query("error")

	ctx := context.Background()

	// Handle gaming OAuth errors
	if errorParam != "" {
		oauth.handleOAuthError(c, provider, errorParam, c.Query("error_description"))
		return
	}

	// Validate gaming OAuth state
	oauthState, err := oauth.StateStore.GetState(ctx, state)
	if err != nil || oauthState == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "Invalid or expired gaming OAuth state",
			"gaming_support": "Please restart the gaming authentication process",
		})
		return
	}

	// Clean up OAuth state
	defer oauth.StateStore.DeleteState(ctx, state)

	// Validate gaming provider consistency
	if oauthState.Provider != provider {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gaming provider mismatch in OAuth flow",
		})
		return
	}

	// Exchange code for token
	config := oauth.getProviderConfig(provider)
	token, err := config.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to exchange gaming OAuth code",
			"provider": provider,
		})
		return
	}

	// Get gaming user info from provider
	userInfo, err := oauth.getUserInfoFromProvider(ctx, provider, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to get gaming user info",
			"provider": provider,
		})
		return
	}

	// Find or create gaming user
	existingUser, err := oauth.UserStore.GetUserByProviderID(ctx, provider, userInfo.ProviderID)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to query gaming user",
		})
		return
	}

	var finalUser *GamingUserInfo

	if existingUser != nil {
		// Update existing gaming user
		existingUser.UpdatedAt = time.Now()
		existingUser.Metadata = mergeMetadata(existingUser.Metadata, oauthState.GamingData)

		if err := oauth.UserStore.UpdateUser(ctx, existingUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update gaming user",
			})
			return
		}

		finalUser = existingUser

		// Track gaming user login
		go oauth.GamingAnalytics.TrackUserLogin(ctx, finalUser.ID, provider, oauthState.GamingData)
	} else {
		// Create new gaming user
		userInfo.CreatedAt = time.Now()
		userInfo.UpdatedAt = time.Now()
		userInfo.Metadata = oauthState.GamingData
		userInfo.GamingProfile = oauth.initializeGamingProfile(provider, userInfo)

		if err := oauth.UserStore.CreateUser(ctx, userInfo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create gaming user",
			})
			return
		}

		finalUser = userInfo

		// Track gaming user registration
		go oauth.GamingAnalytics.TrackUserRegistration(ctx, finalUser.ID, provider, finalUser.GamingProfile)
	}

	// Generate gaming JWT tokens
	accessToken, err := oauth.generateJWT(finalUser, oauth.JWTExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate gaming access token",
		})
		return
	}

	refreshToken, err := oauth.generateJWT(finalUser, oauth.RefreshExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate gaming refresh token",
		})
		return
	}

	// Set secure gaming cookies
	oauth.setGamingTokenCookies(c, accessToken, refreshToken)

	// Determine redirect URL for gaming platform
	redirectURL := oauthState.RedirectURL
	if redirectURL == "" {
		redirectURL = getEnvVar("GAMING_DEFAULT_REDIRECT", "https://herald.lol/dashboard")
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"user":            finalUser,
		"access_token":    accessToken,
		"refresh_token":   refreshToken,
		"expires_in":      int(oauth.JWTExpiration.Seconds()),
		"redirect_url":    redirectURL,
		"gaming_platform": "herald-lol",
		"provider":        provider,
	})
}

// RefreshGamingToken refreshes gaming JWT token
func (oauth *GamingOAuthConfig) RefreshGamingToken(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		refreshToken = c.Query("refresh_token")
	}

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gaming refresh token required",
		})
		return
	}

	// Remove Bearer prefix if present
	if strings.HasPrefix(refreshToken, "Bearer ") {
		refreshToken = refreshToken[7:]
	}

	// Parse and validate gaming refresh token
	claims, err := oauth.parseJWT(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid gaming refresh token",
		})
		return
	}

	// Get gaming user
	ctx := context.Background()
	user, err := oauth.UserStore.GetUserByProviderID(ctx, OAuthProvider(claims.Provider), claims.ProviderID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Gaming user not found",
		})
		return
	}

	// Generate new gaming access token
	accessToken, err := oauth.generateJWT(user, oauth.JWTExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate new gaming access token",
		})
		return
	}

	// Optionally rotate refresh token for enhanced security
	var newRefreshToken string
	if c.Query("rotate") == "true" {
		newRefreshToken, err = oauth.generateJWT(user, oauth.RefreshExpiration)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate new gaming refresh token",
			})
			return
		}
	}

	response := gin.H{
		"access_token":    accessToken,
		"expires_in":      int(oauth.JWTExpiration.Seconds()),
		"gaming_platform": "herald-lol",
		"token_type":      "Bearer",
	}

	if newRefreshToken != "" {
		response["refresh_token"] = newRefreshToken
	}

	c.JSON(http.StatusOK, response)
}

// Gaming JWT Claims structure
type GamingJWTClaims struct {
	UserID            string            `json:"user_id"`
	Email             string            `json:"email"`
	Name              string            `json:"name"`
	Provider          string            `json:"provider"`
	ProviderID        string            `json:"provider_id"`
	SubscriptionTier  string            `json:"subscription_tier"`
	GamingPermissions []string          `json:"gaming_permissions"`
	GamingMetadata    map[string]string `json:"gaming_metadata"`
	jwt.RegisteredClaims
}

// Helper functions

func (oauth *GamingOAuthConfig) isValidProvider(provider OAuthProvider) bool {
	switch provider {
	case ProviderGoogle, ProviderDiscord, ProviderTwitch, ProviderRiot, ProviderGitHub:
		return true
	default:
		return false
	}
}

func (oauth *GamingOAuthConfig) generateSecureState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (oauth *GamingOAuthConfig) getProviderConfig(provider OAuthProvider) *oauth2.Config {
	switch provider {
	case ProviderGoogle:
		return oauth.GoogleConfig
	case ProviderDiscord:
		return oauth.DiscordConfig
	case ProviderTwitch:
		return oauth.TwitchConfig
	case ProviderRiot:
		return oauth.RiotConfig
	case ProviderGitHub:
		return oauth.GitHubConfig
	default:
		return nil
	}
}

func (oauth *GamingOAuthConfig) handleOAuthError(c *gin.Context, provider OAuthProvider, errorCode, errorDescription string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":             errorCode,
		"error_description": errorDescription,
		"provider":          provider,
		"gaming_support":    "Please try again or contact Herald.lol support",
	})
}

func (oauth *GamingOAuthConfig) initializeGamingProfile(provider OAuthProvider, userInfo *GamingUserInfo) *GamingProfile {
	profile := &GamingProfile{
		SubscriptionTier: "free", // Default gaming tier
		Preferences:      make(map[string]string),
	}

	// Provider-specific gaming profile initialization
	switch provider {
	case ProviderDiscord:
		profile.DiscordUsername = userInfo.Username
	case ProviderTwitch:
		profile.TwitchUsername = userInfo.Username
	case ProviderRiot:
		// Initialize with Riot gaming data if available
		profile.Region = "na1" // Default region
	}

	return profile
}

func (oauth *GamingOAuthConfig) generateJWT(user *GamingUserInfo, expiration time.Duration) (string, error) {
	now := time.Now()

	// Determine gaming permissions based on subscription tier
	permissions := oauth.getGamingPermissions(user.GamingProfile.SubscriptionTier)

	claims := GamingJWTClaims{
		UserID:            user.ID,
		Email:             user.Email,
		Name:              user.Name,
		Provider:          string(user.Provider),
		ProviderID:        user.ProviderID,
		SubscriptionTier:  user.GamingProfile.SubscriptionTier,
		GamingPermissions: permissions,
		GamingMetadata:    user.Metadata,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "herald.lol",
			Subject:   user.ID,
			Audience:  []string{"herald-gaming-api"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(oauth.JWTSecret)
}

func (oauth *GamingOAuthConfig) parseJWT(tokenString string) (*GamingJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &GamingJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected gaming token signing method: %v", token.Header["alg"])
		}
		return oauth.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*GamingJWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid gaming JWT token")
}

func (oauth *GamingOAuthConfig) getGamingPermissions(subscriptionTier string) []string {
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

func (oauth *GamingOAuthConfig) setGamingTokenCookies(c *gin.Context, accessToken, refreshToken string) {
	// Set secure gaming access token cookie
	c.SetCookie(
		"herald_access_token",
		accessToken,
		int(oauth.JWTExpiration.Seconds()),
		"/",
		getEnvVar("GAMING_COOKIE_DOMAIN", ".herald.lol"),
		true, // Secure
		true, // HttpOnly
	)

	// Set secure gaming refresh token cookie
	c.SetCookie(
		"herald_refresh_token",
		refreshToken,
		int(oauth.RefreshExpiration.Seconds()),
		"/",
		getEnvVar("GAMING_COOKIE_DOMAIN", ".herald.lol"),
		true, // Secure
		true, // HttpOnly
	)
}

func mergeMetadata(existing, new map[string]string) map[string]string {
	result := make(map[string]string)

	// Copy existing metadata
	for k, v := range existing {
		result[k] = v
	}

	// Override with new metadata
	for k, v := range new {
		result[k] = v
	}

	return result
}

func getEnvVar(key, defaultValue string) string {
	// This would typically use os.Getenv() in a real implementation
	// For now, returning the default value
	return defaultValue
}
