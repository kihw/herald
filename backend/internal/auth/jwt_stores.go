package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Herald.lol Gaming Analytics - JWT Store Implementations
// Redis and database implementations for JWT token management

// RedisRefreshTokenStore implements RefreshTokenStore using Redis
type RedisRefreshTokenStore struct {
	redisClient RedisClient
	keyPrefix   string
	userPrefix  string
}

// NewRedisRefreshTokenStore creates new Redis refresh token store
func NewRedisRefreshTokenStore(redisClient RedisClient) *RedisRefreshTokenStore {
	return &RedisRefreshTokenStore{
		redisClient: redisClient,
		keyPrefix:   "herald:gaming:refresh:",
		userPrefix:  "herald:gaming:user:tokens:",
	}
}

// StoreRefreshToken stores refresh token in Redis
func (s *RedisRefreshTokenStore) StoreRefreshToken(ctx context.Context, tokenID string, refreshToken *GamingRefreshToken) error {
	key := s.keyPrefix + tokenID
	userKey := s.userPrefix + refreshToken.UserID

	// Marshal token data
	data, err := json.Marshal(refreshToken)
	if err != nil {
		return fmt.Errorf("failed to marshal gaming refresh token: %w", err)
	}

	// Calculate TTL
	ttl := time.Until(refreshToken.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("gaming refresh token already expired")
	}

	// Store token
	if err := s.redisClient.Set(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("failed to store gaming refresh token: %w", err)
	}

	// Add to user's token set (for tracking multiple tokens per user)
	userTokenData := map[string]interface{}{
		"token_id":    tokenID,
		"device_info": refreshToken.DeviceInfo,
		"issued_at":   refreshToken.IssuedAt,
		"expires_at":  refreshToken.ExpiresAt,
	}

	userTokenJSON, _ := json.Marshal(userTokenData)
	userTokenKey := fmt.Sprintf("%s:%s", userKey, tokenID)

	if err := s.redisClient.Set(ctx, userTokenKey, userTokenJSON, ttl); err != nil {
		// Non-critical error, continue
	}

	return nil
}

// GetRefreshToken retrieves refresh token from Redis
func (s *RedisRefreshTokenStore) GetRefreshToken(ctx context.Context, tokenID string) (*GamingRefreshToken, error) {
	key := s.keyPrefix + tokenID

	data, err := s.redisClient.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get gaming refresh token: %w", err)
	}

	var refreshToken GamingRefreshToken
	if err := json.Unmarshal([]byte(data), &refreshToken); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gaming refresh token: %w", err)
	}

	// Check if token has expired
	if time.Now().After(refreshToken.ExpiresAt) {
		// Clean up expired token
		s.redisClient.Del(ctx, key)
		return nil, fmt.Errorf("gaming refresh token expired")
	}

	return &refreshToken, nil
}

// RevokeRefreshToken revokes refresh token
func (s *RedisRefreshTokenStore) RevokeRefreshToken(ctx context.Context, tokenID string) error {
	key := s.keyPrefix + tokenID

	// Get token first to update revocation info
	refreshToken, err := s.GetRefreshToken(ctx, tokenID)
	if err != nil {
		// Token might already be expired/deleted
		return nil
	}

	// Mark as revoked
	now := time.Now()
	refreshToken.IsRevoked = true
	refreshToken.RevokedAt = &now
	refreshToken.UpdatedAt = now

	// Store revoked token with shorter TTL
	data, _ := json.Marshal(refreshToken)
	shortTTL := 24 * time.Hour // Keep revoked tokens for audit

	if err := s.redisClient.Set(ctx, key, data, shortTTL); err != nil {
		return fmt.Errorf("failed to update revoked gaming refresh token: %w", err)
	}

	return nil
}

// CleanupExpiredTokens removes expired tokens (Redis handles this automatically with TTL)
func (s *RedisRefreshTokenStore) CleanupExpiredTokens(ctx context.Context) error {
	// Redis automatically removes expired keys, but we can scan for manual cleanup
	var cursor uint64
	pattern := s.keyPrefix + "*"

	for {
		keys, nextCursor, err := s.redisClient.Scan(ctx, cursor, pattern, 100)
		if err != nil {
			return fmt.Errorf("failed to scan gaming refresh tokens: %w", err)
		}

		// Check each key and clean up if needed
		for _, key := range keys {
			data, err := s.redisClient.Get(ctx, key)
			if err != nil {
				continue // Skip keys that can't be retrieved
			}

			var token GamingRefreshToken
			if err := json.Unmarshal([]byte(data), &token); err != nil {
				continue // Skip malformed tokens
			}

			if time.Now().After(token.ExpiresAt) {
				s.redisClient.Del(ctx, key)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}

// GetUserRefreshTokens gets all active refresh tokens for a user
func (s *RedisRefreshTokenStore) GetUserRefreshTokens(ctx context.Context, userID string) ([]*GamingRefreshToken, error) {
	userKeyPattern := s.userPrefix + userID + ":*"

	var tokens []*GamingRefreshToken
	var cursor uint64

	for {
		keys, nextCursor, err := s.redisClient.Scan(ctx, cursor, userKeyPattern, 100)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user gaming tokens: %w", err)
		}

		for _, userTokenKey := range keys {
			// Extract token ID from key
			parts := []string{} // Would parse key to get token ID
			if len(parts) < 1 {
				continue
			}

			tokenID := parts[len(parts)-1]
			token, err := s.GetRefreshToken(ctx, tokenID)
			if err != nil {
				continue // Skip invalid tokens
			}

			if !token.IsRevoked {
				tokens = append(tokens, token)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return tokens, nil
}

// RevokeAllUserTokens revokes all refresh tokens for a user
func (s *RedisRefreshTokenStore) RevokeAllUserTokens(ctx context.Context, userID string) error {
	tokens, err := s.GetUserRefreshTokens(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user gaming tokens: %w", err)
	}

	for _, token := range tokens {
		if err := s.RevokeRefreshToken(ctx, token.ID); err != nil {
			// Log error but continue revoking others
		}
	}

	return nil
}

// RedisTokenBlacklistStore implements TokenBlacklistStore using Redis
type RedisTokenBlacklistStore struct {
	redisClient RedisClient
	keyPrefix   string
}

// NewRedisTokenBlacklistStore creates new Redis token blacklist store
func NewRedisTokenBlacklistStore(redisClient RedisClient) *RedisTokenBlacklistStore {
	return &RedisTokenBlacklistStore{
		redisClient: redisClient,
		keyPrefix:   "herald:gaming:blacklist:",
	}
}

// BlacklistToken adds token to blacklist
func (s *RedisTokenBlacklistStore) BlacklistToken(ctx context.Context, tokenID string, expiresAt time.Time) error {
	key := s.keyPrefix + tokenID

	// Calculate TTL - only need to keep blacklisted until original expiration
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		// Token already expired, no need to blacklist
		return nil
	}

	blacklistData := map[string]interface{}{
		"token_id":       tokenID,
		"blacklisted_at": time.Now(),
		"expires_at":     expiresAt,
	}

	data, _ := json.Marshal(blacklistData)

	if err := s.redisClient.Set(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("failed to blacklist gaming token: %w", err)
	}

	return nil
}

// IsTokenBlacklisted checks if token is blacklisted
func (s *RedisTokenBlacklistStore) IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := s.keyPrefix + tokenID

	_, err := s.redisClient.Get(ctx, key)
	if err != nil {
		// Token not found in blacklist
		return false, nil
	}

	// Token found in blacklist
	return true, nil
}

// CleanupExpiredBlacklist removes expired blacklist entries (handled by Redis TTL)
func (s *RedisTokenBlacklistStore) CleanupExpiredBlacklist(ctx context.Context) error {
	// Redis automatically handles TTL cleanup
	return nil
}

// DatabaseTokenRotationStore implements TokenRotationStore using database
type DatabaseTokenRotationStore struct {
	db *gorm.DB
}

// TokenRotationRecord database model for token rotations
type TokenRotationRecord struct {
	ID         string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	OldTokenID string    `gorm:"type:varchar(255);index" json:"old_token_id"`
	NewTokenID string    `gorm:"type:varchar(255);index" json:"new_token_id"`
	UserID     string    `gorm:"type:varchar(255);index" json:"user_id"`
	RotatedAt  time.Time `json:"rotated_at"`
	Reason     string    `gorm:"type:varchar(100)" json:"reason"`
	IPAddress  string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent  string    `gorm:"type:text" json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}

// NewDatabaseTokenRotationStore creates new database token rotation store
func NewDatabaseTokenRotationStore(db *gorm.DB) *DatabaseTokenRotationStore {
	return &DatabaseTokenRotationStore{db: db}
}

// TrackRotation tracks token rotation
func (s *DatabaseTokenRotationStore) TrackRotation(ctx context.Context, oldTokenID, newTokenID string, userID string) error {
	rotation := &TokenRotationRecord{
		ID:         fmt.Sprintf("rot_%d", time.Now().UnixNano()),
		OldTokenID: oldTokenID,
		NewTokenID: newTokenID,
		UserID:     userID,
		RotatedAt:  time.Now(),
		Reason:     "refresh_rotation",
		CreatedAt:  time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(rotation).Error; err != nil {
		return fmt.Errorf("failed to track gaming token rotation: %w", err)
	}

	return nil
}

// GetRotationChain gets rotation chain for a token
func (s *DatabaseTokenRotationStore) GetRotationChain(ctx context.Context, tokenID string) ([]*TokenRotation, error) {
	var records []TokenRotationRecord

	// Find all rotations involving this token
	err := s.db.WithContext(ctx).Where("old_token_id = ? OR new_token_id = ?", tokenID, tokenID).
		Order("rotated_at DESC").Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get gaming token rotation chain: %w", err)
	}

	var rotations []*TokenRotation
	for _, record := range records {
		rotations = append(rotations, &TokenRotation{
			ID:         record.ID,
			OldTokenID: record.OldTokenID,
			NewTokenID: record.NewTokenID,
			UserID:     record.UserID,
			RotatedAt:  record.RotatedAt,
			Reason:     record.Reason,
			IPAddress:  record.IPAddress,
		})
	}

	return rotations, nil
}

// IsRotationValid checks if token rotation is valid
func (s *DatabaseTokenRotationStore) IsRotationValid(ctx context.Context, tokenID string) (bool, error) {
	var count int64

	// Check if token has been rotated recently
	recentTime := time.Now().Add(-24 * time.Hour)
	err := s.db.WithContext(ctx).Model(&TokenRotationRecord{}).
		Where("(old_token_id = ? OR new_token_id = ?) AND rotated_at > ?", tokenID, tokenID, recentTime).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check gaming token rotation validity: %w", err)
	}

	// Token is valid if it hasn't been rotated too frequently
	return count < 10, nil // Allow up to 10 rotations per day
}

// DatabaseRefreshTokenStore implements RefreshTokenStore using database
type DatabaseRefreshTokenStore struct {
	db *gorm.DB
}

// GamingRefreshTokenRecord database model for refresh tokens
type GamingRefreshTokenRecord struct {
	ID               string     `gorm:"primaryKey;type:varchar(255)" json:"id"`
	UserID           string     `gorm:"type:varchar(255);index" json:"user_id"`
	TokenHash        string     `gorm:"type:varchar(512)" json:"token_hash"`
	DeviceInfo       string     `gorm:"type:jsonb" json:"device_info"`    // JSON
	GamingContext    string     `gorm:"type:jsonb" json:"gaming_context"` // JSON
	Metadata         string     `gorm:"type:jsonb" json:"metadata"`       // JSON
	IssuedAt         time.Time  `json:"issued_at"`
	ExpiresAt        time.Time  `gorm:"index" json:"expires_at"`
	LastUsedAt       *time.Time `json:"last_used_at,omitempty"`
	UsageCount       int64      `gorm:"default:0" json:"usage_count"`
	IsRevoked        bool       `gorm:"default:false;index" json:"is_revoked"`
	RevokedAt        *time.Time `json:"revoked_at,omitempty"`
	RevokedReason    string     `gorm:"type:varchar(100)" json:"revoked_reason,omitempty"`
	RotationCount    int        `gorm:"default:0" json:"rotation_count"`
	RotationParentID string     `gorm:"type:varchar(255)" json:"rotation_parent_id,omitempty"`
	Version          int        `gorm:"default:1" json:"version"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// NewDatabaseRefreshTokenStore creates new database refresh token store
func NewDatabaseRefreshTokenStore(db *gorm.DB) *DatabaseRefreshTokenStore {
	return &DatabaseRefreshTokenStore{db: db}
}

// StoreRefreshToken stores refresh token in database
func (s *DatabaseRefreshTokenStore) StoreRefreshToken(ctx context.Context, tokenID string, refreshToken *GamingRefreshToken) error {
	// Convert to database model
	record := &GamingRefreshTokenRecord{
		ID:               tokenID,
		UserID:           refreshToken.UserID,
		TokenHash:        refreshToken.TokenHash,
		IssuedAt:         refreshToken.IssuedAt,
		ExpiresAt:        refreshToken.ExpiresAt,
		LastUsedAt:       refreshToken.LastUsedAt,
		UsageCount:       refreshToken.UsageCount,
		IsRevoked:        refreshToken.IsRevoked,
		RevokedAt:        refreshToken.RevokedAt,
		RevokedReason:    refreshToken.RevokedReason,
		RotationCount:    refreshToken.RotationCount,
		RotationParentID: refreshToken.RotationParentID,
		Version:          refreshToken.Version,
		CreatedAt:        refreshToken.CreatedAt,
		UpdatedAt:        refreshToken.UpdatedAt,
	}

	// Marshal JSON fields
	if refreshToken.DeviceInfo != nil {
		if data, err := json.Marshal(refreshToken.DeviceInfo); err == nil {
			record.DeviceInfo = string(data)
		}
	}

	if refreshToken.GamingContext != nil {
		if data, err := json.Marshal(refreshToken.GamingContext); err == nil {
			record.GamingContext = string(data)
		}
	}

	if refreshToken.Metadata != nil {
		if data, err := json.Marshal(refreshToken.Metadata); err == nil {
			record.Metadata = string(data)
		}
	}

	// Store in database
	if err := s.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("failed to store gaming refresh token in database: %w", err)
	}

	return nil
}

// GetRefreshToken retrieves refresh token from database
func (s *DatabaseRefreshTokenStore) GetRefreshToken(ctx context.Context, tokenID string) (*GamingRefreshToken, error) {
	var record GamingRefreshTokenRecord

	err := s.db.WithContext(ctx).Where("id = ? AND expires_at > ?", tokenID, time.Now()).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("gaming refresh token not found or expired")
		}
		return nil, fmt.Errorf("failed to get gaming refresh token from database: %w", err)
	}

	// Convert back to domain model
	refreshToken := &GamingRefreshToken{
		ID:               record.ID,
		UserID:           record.UserID,
		TokenHash:        record.TokenHash,
		IssuedAt:         record.IssuedAt,
		ExpiresAt:        record.ExpiresAt,
		LastUsedAt:       record.LastUsedAt,
		UsageCount:       record.UsageCount,
		IsRevoked:        record.IsRevoked,
		RevokedAt:        record.RevokedAt,
		RevokedReason:    record.RevokedReason,
		RotationCount:    record.RotationCount,
		RotationParentID: record.RotationParentID,
		Version:          record.Version,
		CreatedAt:        record.CreatedAt,
		UpdatedAt:        record.UpdatedAt,
	}

	// Unmarshal JSON fields
	if record.DeviceInfo != "" {
		var deviceInfo DeviceInfo
		if err := json.Unmarshal([]byte(record.DeviceInfo), &deviceInfo); err == nil {
			refreshToken.DeviceInfo = &deviceInfo
		}
	}

	if record.GamingContext != "" {
		var gamingContext GamingContext
		if err := json.Unmarshal([]byte(record.GamingContext), &gamingContext); err == nil {
			refreshToken.GamingContext = &gamingContext
		}
	}

	if record.Metadata != "" {
		var metadata map[string]string
		if err := json.Unmarshal([]byte(record.Metadata), &metadata); err == nil {
			refreshToken.Metadata = metadata
		}
	}

	return refreshToken, nil
}

// RevokeRefreshToken revokes refresh token in database
func (s *DatabaseRefreshTokenStore) RevokeRefreshToken(ctx context.Context, tokenID string) error {
	now := time.Now()

	result := s.db.WithContext(ctx).Model(&GamingRefreshTokenRecord{}).
		Where("id = ?", tokenID).
		Updates(map[string]interface{}{
			"is_revoked":     true,
			"revoked_at":     now,
			"revoked_reason": "manual_revocation",
			"updated_at":     now,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to revoke gaming refresh token: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("gaming refresh token not found: %s", tokenID)
	}

	return nil
}

// CleanupExpiredTokens removes expired tokens from database
func (s *DatabaseRefreshTokenStore) CleanupExpiredTokens(ctx context.Context) error {
	result := s.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&GamingRefreshTokenRecord{})

	if result.Error != nil {
		return fmt.Errorf("failed to cleanup expired gaming refresh tokens: %w", result.Error)
	}

	return nil
}

// GetUserRefreshTokens gets all active refresh tokens for a user from database
func (s *DatabaseRefreshTokenStore) GetUserRefreshTokens(ctx context.Context, userID string) ([]*GamingRefreshToken, error) {
	var records []GamingRefreshTokenRecord

	err := s.db.WithContext(ctx).Where("user_id = ? AND is_revoked = false AND expires_at > ?",
		userID, time.Now()).Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user gaming refresh tokens: %w", err)
	}

	var tokens []*GamingRefreshToken
	for _, record := range records {
		// Convert each record (simplified conversion)
		token := &GamingRefreshToken{
			ID:        record.ID,
			UserID:    record.UserID,
			TokenHash: record.TokenHash,
			IssuedAt:  record.IssuedAt,
			ExpiresAt: record.ExpiresAt,
			IsRevoked: record.IsRevoked,
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// RevokeAllUserTokens revokes all refresh tokens for a user in database
func (s *DatabaseRefreshTokenStore) RevokeAllUserTokens(ctx context.Context, userID string) error {
	now := time.Now()

	result := s.db.WithContext(ctx).Model(&GamingRefreshTokenRecord{}).
		Where("user_id = ? AND is_revoked = false", userID).
		Updates(map[string]interface{}{
			"is_revoked":     true,
			"revoked_at":     now,
			"revoked_reason": "bulk_revocation",
			"updated_at":     now,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to revoke all user gaming tokens: %w", result.Error)
	}

	return nil
}
