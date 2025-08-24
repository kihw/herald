package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Herald.lol Gaming Analytics - Auth Store Implementations
// Database and cache implementations for gaming OAuth state and user management

// RedisStateStore implements StateStore using Redis for gaming OAuth state
type RedisStateStore struct {
	redisClient RedisClient
	keyPrefix   string
}

// RedisClient interface for Redis operations
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)
}

// NewRedisStateStore creates new Redis state store for gaming OAuth
func NewRedisStateStore(redisClient RedisClient, keyPrefix string) *RedisStateStore {
	if keyPrefix == "" {
		keyPrefix = "herald:gaming:oauth:state:"
	}

	return &RedisStateStore{
		redisClient: redisClient,
		keyPrefix:   keyPrefix,
	}
}

// StoreState stores OAuth state in Redis with gaming metadata
func (s *RedisStateStore) StoreState(ctx context.Context, state string, oauthState *OAuthState) error {
	key := s.keyPrefix + state

	data, err := json.Marshal(oauthState)
	if err != nil {
		return fmt.Errorf("failed to marshal gaming OAuth state: %w", err)
	}

	expiration := time.Until(oauthState.ExpiresAt)
	if expiration <= 0 {
		expiration = 10 * time.Minute // Default gaming OAuth timeout
	}

	if err := s.redisClient.Set(ctx, key, data, expiration); err != nil {
		return fmt.Errorf("failed to store gaming OAuth state in Redis: %w", err)
	}

	return nil
}

// GetState retrieves OAuth state from Redis
func (s *RedisStateStore) GetState(ctx context.Context, state string) (*OAuthState, error) {
	key := s.keyPrefix + state

	data, err := s.redisClient.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get gaming OAuth state from Redis: %w", err)
	}

	var oauthState OAuthState
	if err := json.Unmarshal([]byte(data), &oauthState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gaming OAuth state: %w", err)
	}

	// Check if state has expired
	if time.Now().After(oauthState.ExpiresAt) {
		// Clean up expired state
		s.redisClient.Del(ctx, key)
		return nil, fmt.Errorf("gaming OAuth state has expired")
	}

	return &oauthState, nil
}

// DeleteState removes OAuth state from Redis
func (s *RedisStateStore) DeleteState(ctx context.Context, state string) error {
	key := s.keyPrefix + state

	if err := s.redisClient.Del(ctx, key); err != nil {
		return fmt.Errorf("failed to delete gaming OAuth state from Redis: %w", err)
	}

	return nil
}

// CleanupExpiredStates removes expired OAuth states (periodic cleanup)
func (s *RedisStateStore) CleanupExpiredStates(ctx context.Context) error {
	var cursor uint64
	pattern := s.keyPrefix + "*"

	for {
		keys, nextCursor, err := s.redisClient.Scan(ctx, cursor, pattern, 100)
		if err != nil {
			return fmt.Errorf("failed to scan gaming OAuth states: %w", err)
		}

		// Check each key and delete expired ones
		for _, key := range keys {
			data, err := s.redisClient.Get(ctx, key)
			if err != nil {
				continue // Skip keys that can't be retrieved
			}

			var oauthState OAuthState
			if err := json.Unmarshal([]byte(data), &oauthState); err != nil {
				continue // Skip malformed states
			}

			if time.Now().After(oauthState.ExpiresAt) {
				s.redisClient.Del(ctx, key)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break // Scan complete
		}
	}

	return nil
}

// DatabaseUserStore implements UserStore using GORM for gaming users
type DatabaseUserStore struct {
	db *gorm.DB
}

// GamingUser database model for gaming users
type GamingUser struct {
	ID            string     `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Email         string     `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Name          string     `gorm:"type:varchar(255)" json:"name"`
	Username      string     `gorm:"type:varchar(255)" json:"username"`
	Avatar        string     `gorm:"type:text" json:"avatar"`
	Provider      string     `gorm:"type:varchar(50)" json:"provider"`
	ProviderID    string     `gorm:"type:varchar(255);uniqueIndex:idx_provider_id" json:"provider_id"`
	Metadata      string     `gorm:"type:jsonb" json:"metadata"`       // JSON metadata
	GamingProfile string     `gorm:"type:jsonb" json:"gaming_profile"` // JSON gaming profile
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	EmailVerified bool       `gorm:"default:false" json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Gaming-specific fields
	SubscriptionTier      string     `gorm:"type:varchar(50);default:'free'" json:"subscription_tier"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at"`
	GamingPermissions     string     `gorm:"type:jsonb" json:"gaming_permissions"` // JSON permissions array
	TotalPlaytime         int64      `gorm:"default:0" json:"total_playtime"`      // Minutes
	AnalyticsCount        int64      `gorm:"default:0" json:"analytics_count"`
	APIUsageCount         int64      `gorm:"default:0" json:"api_usage_count"`
	LastAnalyticsAt       *time.Time `json:"last_analytics_at"`
}

// GamingUserSession represents active gaming sessions
type GamingUserSession struct {
	ID               string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID           string    `gorm:"type:varchar(255);index" json:"user_id"`
	SessionToken     string    `gorm:"type:varchar(512);uniqueIndex" json:"session_token"`
	RefreshToken     string    `gorm:"type:varchar(512);uniqueIndex" json:"refresh_token"`
	IPAddress        string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent        string    `gorm:"type:text" json:"user_agent"`
	GamingClientInfo string    `gorm:"type:jsonb" json:"gaming_client_info"` // JSON client metadata
	ExpiresAt        time.Time `json:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	LastActivityAt   time.Time `json:"last_activity_at"`
	IsActive         bool      `gorm:"default:true" json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Gaming-specific session data
	GamingContext    string `gorm:"type:jsonb" json:"gaming_context"`    // Current gaming context
	AnalyticsSession string `gorm:"type:jsonb" json:"analytics_session"` // Analytics session data
}

// NewDatabaseUserStore creates new database user store for gaming
func NewDatabaseUserStore(db *gorm.DB) *DatabaseUserStore {
	return &DatabaseUserStore{db: db}
}

// GetUserByProviderID retrieves gaming user by provider and provider ID
func (s *DatabaseUserStore) GetUserByProviderID(ctx context.Context, provider OAuthProvider, providerID string) (*GamingUserInfo, error) {
	var dbUser GamingUser

	err := s.db.WithContext(ctx).Where("provider = ? AND provider_id = ?", string(provider), providerID).First(&dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get gaming user by provider ID: %w", err)
	}

	return s.dbUserToUserInfo(&dbUser)
}

// GetUserByEmail retrieves gaming user by email
func (s *DatabaseUserStore) GetUserByEmail(ctx context.Context, email string) (*GamingUserInfo, error) {
	var dbUser GamingUser

	err := s.db.WithContext(ctx).Where("email = ?", email).First(&dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get gaming user by email: %w", err)
	}

	return s.dbUserToUserInfo(&dbUser)
}

// CreateUser creates new gaming user
func (s *DatabaseUserStore) CreateUser(ctx context.Context, user *GamingUserInfo) error {
	dbUser, err := s.userInfoToDBUser(user)
	if err != nil {
		return fmt.Errorf("failed to convert gaming user info: %w", err)
	}

	// Set gaming-specific defaults
	dbUser.IsActive = true
	dbUser.EmailVerified = user.Email != ""
	dbUser.SubscriptionTier = "free"
	dbUser.CreatedAt = time.Now()
	dbUser.UpdatedAt = time.Now()

	if err := s.db.WithContext(ctx).Create(dbUser).Error; err != nil {
		return fmt.Errorf("failed to create gaming user: %w", err)
	}

	// Update the user info with database-generated values
	updatedUser, err := s.dbUserToUserInfo(dbUser)
	if err != nil {
		return fmt.Errorf("failed to convert created gaming user: %w", err)
	}

	*user = *updatedUser
	return nil
}

// UpdateUser updates existing gaming user
func (s *DatabaseUserStore) UpdateUser(ctx context.Context, user *GamingUserInfo) error {
	dbUser, err := s.userInfoToDBUser(user)
	if err != nil {
		return fmt.Errorf("failed to convert gaming user info: %w", err)
	}

	dbUser.UpdatedAt = time.Now()
	dbUser.LastLoginAt = &dbUser.UpdatedAt

	if err := s.db.WithContext(ctx).Save(dbUser).Error; err != nil {
		return fmt.Errorf("failed to update gaming user: %w", err)
	}

	return nil
}

// UpdateGamingProfile updates gaming-specific profile
func (s *DatabaseUserStore) UpdateGamingProfile(ctx context.Context, userID string, profile *GamingProfile) error {
	profileJSON, err := json.Marshal(profile)
	if err != nil {
		return fmt.Errorf("failed to marshal gaming profile: %w", err)
	}

	result := s.db.WithContext(ctx).Model(&GamingUser{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"gaming_profile": string(profileJSON),
			"updated_at":     time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update gaming profile: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("gaming user not found: %s", userID)
	}

	return nil
}

// Helper methods for data conversion

func (s *DatabaseUserStore) dbUserToUserInfo(dbUser *GamingUser) (*GamingUserInfo, error) {
	userInfo := &GamingUserInfo{
		ID:         dbUser.ID,
		Email:      dbUser.Email,
		Name:       dbUser.Name,
		Username:   dbUser.Username,
		Avatar:     dbUser.Avatar,
		Provider:   OAuthProvider(dbUser.Provider),
		ProviderID: dbUser.ProviderID,
		CreatedAt:  dbUser.CreatedAt,
		UpdatedAt:  dbUser.UpdatedAt,
	}

	// Parse metadata JSON
	if dbUser.Metadata != "" {
		var metadata map[string]string
		if err := json.Unmarshal([]byte(dbUser.Metadata), &metadata); err == nil {
			userInfo.Metadata = metadata
		}
	}

	// Parse gaming profile JSON
	if dbUser.GamingProfile != "" {
		var profile GamingProfile
		if err := json.Unmarshal([]byte(dbUser.GamingProfile), &profile); err == nil {
			userInfo.GamingProfile = &profile
		}
	}

	// Ensure gaming profile exists
	if userInfo.GamingProfile == nil {
		userInfo.GamingProfile = &GamingProfile{
			SubscriptionTier: dbUser.SubscriptionTier,
			Preferences:      make(map[string]string),
		}
	}

	return userInfo, nil
}

func (s *DatabaseUserStore) userInfoToDBUser(userInfo *GamingUserInfo) (*GamingUser, error) {
	dbUser := &GamingUser{
		ID:         userInfo.ID,
		Email:      userInfo.Email,
		Name:       userInfo.Name,
		Username:   userInfo.Username,
		Avatar:     userInfo.Avatar,
		Provider:   string(userInfo.Provider),
		ProviderID: userInfo.ProviderID,
		CreatedAt:  userInfo.CreatedAt,
		UpdatedAt:  userInfo.UpdatedAt,
	}

	// Marshal metadata to JSON
	if userInfo.Metadata != nil {
		metadataJSON, err := json.Marshal(userInfo.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal gaming user metadata: %w", err)
		}
		dbUser.Metadata = string(metadataJSON)
	}

	// Marshal gaming profile to JSON
	if userInfo.GamingProfile != nil {
		profileJSON, err := json.Marshal(userInfo.GamingProfile)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal gaming profile: %w", err)
		}
		dbUser.GamingProfile = string(profileJSON)
		dbUser.SubscriptionTier = userInfo.GamingProfile.SubscriptionTier
	}

	return dbUser, nil
}

// Analytics Store Implementation
type GamingAnalyticsStore struct {
	db *gorm.DB
}

// NewGamingAnalyticsStore creates new gaming analytics store
func NewGamingAnalyticsStore(db *gorm.DB) *GamingAnalyticsStore {
	return &GamingAnalyticsStore{db: db}
}

// TrackUserLogin tracks gaming user login events
func (s *GamingAnalyticsStore) TrackUserLogin(ctx context.Context, userID string, provider OAuthProvider, metadata map[string]string) {
	// Implementation for tracking login events
	// This would typically insert into an analytics table
	go func() {
		// Async analytics tracking
		loginEvent := map[string]interface{}{
			"event_type": "user_login",
			"user_id":    userID,
			"provider":   string(provider),
			"metadata":   metadata,
			"timestamp":  time.Now(),
		}

		// Store in analytics table or send to analytics service
		s.storeAnalyticsEvent(ctx, loginEvent)
	}()
}

// TrackUserRegistration tracks gaming user registration events
func (s *GamingAnalyticsStore) TrackUserRegistration(ctx context.Context, userID string, provider OAuthProvider, profile *GamingProfile) {
	// Implementation for tracking registration events
	go func() {
		registrationEvent := map[string]interface{}{
			"event_type":        "user_registration",
			"user_id":           userID,
			"provider":          string(provider),
			"subscription_tier": profile.SubscriptionTier,
			"gaming_profile":    profile,
			"timestamp":         time.Now(),
		}

		s.storeAnalyticsEvent(ctx, registrationEvent)
	}()
}

func (s *GamingAnalyticsStore) storeAnalyticsEvent(ctx context.Context, event map[string]interface{}) {
	// Implementation would store event in analytics database
	// For now, this is a placeholder
}

// Session Management

// CreateGamingSession creates new gaming session
func (s *DatabaseUserStore) CreateGamingSession(ctx context.Context, userID, sessionToken, refreshToken string, clientInfo map[string]interface{}) error {
	clientInfoJSON, _ := json.Marshal(clientInfo)

	session := &GamingUserSession{
		UserID:           userID,
		SessionToken:     sessionToken,
		RefreshToken:     refreshToken,
		GamingClientInfo: string(clientInfoJSON),
		ExpiresAt:        time.Now().Add(15 * time.Minute),   // Gaming session duration
		RefreshExpiresAt: time.Now().Add(7 * 24 * time.Hour), // Gaming refresh duration
		LastActivityAt:   time.Now(),
		IsActive:         true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(session).Error; err != nil {
		return fmt.Errorf("failed to create gaming session: %w", err)
	}

	return nil
}

// GetGamingSession retrieves gaming session by token
func (s *DatabaseUserStore) GetGamingSession(ctx context.Context, sessionToken string) (*GamingUserSession, error) {
	var session GamingUserSession

	err := s.db.WithContext(ctx).Where("session_token = ? AND is_active = ? AND expires_at > ?",
		sessionToken, true, time.Now()).First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get gaming session: %w", err)
	}

	return &session, nil
}

// UpdateGamingSessionActivity updates gaming session last activity
func (s *DatabaseUserStore) UpdateGamingSessionActivity(ctx context.Context, sessionToken string) error {
	result := s.db.WithContext(ctx).Model(&GamingUserSession{}).
		Where("session_token = ?", sessionToken).
		Update("last_activity_at", time.Now())

	if result.Error != nil {
		return fmt.Errorf("failed to update gaming session activity: %w", result.Error)
	}

	return nil
}

// InvalidateGamingSession invalidates gaming session
func (s *DatabaseUserStore) InvalidateGamingSession(ctx context.Context, sessionToken string) error {
	result := s.db.WithContext(ctx).Model(&GamingUserSession{}).
		Where("session_token = ?", sessionToken).
		Update("is_active", false)

	if result.Error != nil {
		return fmt.Errorf("failed to invalidate gaming session: %w", result.Error)
	}

	return nil
}
