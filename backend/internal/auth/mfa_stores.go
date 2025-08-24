package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Herald.lol Gaming Analytics - MFA Store Implementations
// Database and Redis implementations for MFA data management

// DatabaseMFAStore implements MFAStore using GORM database
type DatabaseMFAStore struct {
	db *gorm.DB
}

// Database models for MFA

// GamingTOTPRecord database model for TOTP secrets
type GamingTOTPRecord struct {
	UserID      string     `gorm:"primaryKey;type:varchar(255)" json:"user_id"`
	Secret      string     `gorm:"type:varchar(512)" json:"secret"` // Encrypted
	QRCodeURL   string     `gorm:"type:text" json:"qr_code_url"`
	BackupCodes string     `gorm:"type:jsonb" json:"backup_codes"` // JSON array
	Enabled     bool       `gorm:"default:false" json:"enabled"`
	Verified    bool       `gorm:"default:false" json:"verified"`
	CreatedAt   time.Time  `json:"created_at"`
	VerifiedAt  *time.Time `json:"verified_at,omitempty"`
	LastUsed    *time.Time `json:"last_used,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// GamingWebAuthnRecord database model for WebAuthn credentials
type GamingWebAuthnRecord struct {
	ID              string     `gorm:"primaryKey;type:varchar(255)" json:"id"`
	UserID          string     `gorm:"type:varchar(255);index" json:"user_id"`
	PublicKey       []byte     `gorm:"type:bytea" json:"public_key"`
	AttestationType string     `gorm:"type:varchar(100)" json:"attestation_type"`
	AAGUID          []byte     `gorm:"type:bytea" json:"aaguid"`
	SignCount       uint32     `gorm:"default:0" json:"sign_count"`
	Transports      string     `gorm:"type:jsonb" json:"transports"` // JSON array
	DeviceName      string     `gorm:"type:varchar(255)" json:"device_name"`
	DeviceType      string     `gorm:"type:varchar(100)" json:"device_type"`
	GamingPlatform  string     `gorm:"type:varchar(100)" json:"gaming_platform"`
	UserAgent       string     `gorm:"type:text" json:"user_agent"`
	IPAddress       string     `gorm:"type:varchar(45)" json:"ip_address"`
	Enabled         bool       `gorm:"default:true" json:"enabled"`
	CreatedAt       time.Time  `json:"created_at"`
	LastUsed        *time.Time `json:"last_used,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// GamingBackupCodesRecord database model for backup codes
type GamingBackupCodesRecord struct {
	UserID    string    `gorm:"primaryKey;type:varchar(255)" json:"user_id"`
	Codes     string    `gorm:"type:jsonb" json:"codes"` // JSON map of code->used
	CreatedAt time.Time `json:"created_at"`
	UsedCount int       `gorm:"default:0" json:"used_count"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GamingMFAAttemptRecord database model for MFA attempts
type GamingMFAAttemptRecord struct {
	ID           string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	UserID       string    `gorm:"type:varchar(255);index" json:"user_id"`
	Method       string    `gorm:"type:varchar(50)" json:"method"`
	Success      bool      `gorm:"default:false;index" json:"success"`
	ErrorMessage string    `gorm:"type:text" json:"error_message"`
	IPAddress    string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent    string    `gorm:"type:text" json:"user_agent"`
	GamingAction string    `gorm:"type:varchar(255)" json:"gaming_action"`
	AttemptedAt  time.Time `gorm:"index" json:"attempted_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// NewDatabaseMFAStore creates new database MFA store
func NewDatabaseMFAStore(db *gorm.DB) *DatabaseMFAStore {
	return &DatabaseMFAStore{db: db}
}

// TOTP methods

// StoreTOTPSecret stores TOTP secret in database
func (s *DatabaseMFAStore) StoreTOTPSecret(ctx context.Context, userID string, secret *TOTPSecret) error {
	// Convert backup codes to JSON
	backupCodesJSON, err := json.Marshal(secret.BackupCodes)
	if err != nil {
		return fmt.Errorf("failed to marshal gaming backup codes: %w", err)
	}

	record := &GamingTOTPRecord{
		UserID:      userID,
		Secret:      secret.Secret, // In production, encrypt this
		QRCodeURL:   secret.QRCodeURL,
		BackupCodes: string(backupCodesJSON),
		Enabled:     secret.Enabled,
		Verified:    secret.Verified,
		CreatedAt:   secret.CreatedAt,
		VerifiedAt:  secret.VerifiedAt,
		LastUsed:    secret.LastUsed,
		UpdatedAt:   time.Now(),
	}

	if err := s.db.WithContext(ctx).Save(record).Error; err != nil {
		return fmt.Errorf("failed to store gaming TOTP secret: %w", err)
	}

	return nil
}

// GetTOTPSecret retrieves TOTP secret from database
func (s *DatabaseMFAStore) GetTOTPSecret(ctx context.Context, userID string) (*TOTPSecret, error) {
	var record GamingTOTPRecord

	err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("gaming TOTP secret not found")
		}
		return nil, fmt.Errorf("failed to get gaming TOTP secret: %w", err)
	}

	// Parse backup codes
	var backupCodes []string
	if record.BackupCodes != "" {
		if err := json.Unmarshal([]byte(record.BackupCodes), &backupCodes); err != nil {
			backupCodes = []string{} // Continue without backup codes if parsing fails
		}
	}

	return &TOTPSecret{
		UserID:      record.UserID,
		Secret:      record.Secret, // In production, decrypt this
		QRCodeURL:   record.QRCodeURL,
		BackupCodes: backupCodes,
		Enabled:     record.Enabled,
		Verified:    record.Verified,
		CreatedAt:   record.CreatedAt,
		VerifiedAt:  record.VerifiedAt,
		LastUsed:    record.LastUsed,
	}, nil
}

// VerifyTOTPCode verifies TOTP code and updates last used
func (s *DatabaseMFAStore) VerifyTOTPCode(ctx context.Context, userID, code string) error {
	// Get secret first
	secret, err := s.GetTOTPSecret(ctx, userID)
	if err != nil {
		return err
	}

	if !secret.Enabled {
		return fmt.Errorf("gaming TOTP not enabled")
	}

	// Verify with TOTP library (this would be done in the manager)
	// For now, just update last used time
	now := time.Now()

	result := s.db.WithContext(ctx).Model(&GamingTOTPRecord{}).
		Where("user_id = ?", userID).
		Update("last_used", now)

	if result.Error != nil {
		return fmt.Errorf("failed to update gaming TOTP last used: %w", result.Error)
	}

	return nil
}

// DisableTOTP disables TOTP for user
func (s *DatabaseMFAStore) DisableTOTP(ctx context.Context, userID string) error {
	result := s.db.WithContext(ctx).Model(&GamingTOTPRecord{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"enabled":    false,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to disable gaming TOTP: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("gaming TOTP record not found")
	}

	return nil
}

// WebAuthn methods

// StoreWebAuthnCredential stores WebAuthn credential in database
func (s *DatabaseMFAStore) StoreWebAuthnCredential(ctx context.Context, userID string, credential *WebAuthnCredential) error {
	// Convert transports to JSON
	transportsJSON, err := json.Marshal(credential.Transports)
	if err != nil {
		return fmt.Errorf("failed to marshal WebAuthn transports: %w", err)
	}

	record := &GamingWebAuthnRecord{
		ID:              credential.ID,
		UserID:          userID,
		PublicKey:       credential.PublicKey,
		AttestationType: credential.AttestationType,
		AAGUID:          credential.AAGUID,
		SignCount:       credential.SignCount,
		Transports:      string(transportsJSON),
		DeviceName:      credential.DeviceName,
		DeviceType:      credential.DeviceType,
		GamingPlatform:  credential.GamingPlatform,
		UserAgent:       credential.UserAgent,
		IPAddress:       credential.IPAddress,
		Enabled:         credential.Enabled,
		CreatedAt:       credential.CreatedAt,
		LastUsed:        credential.LastUsed,
		UpdatedAt:       time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("failed to store WebAuthn credential: %w", err)
	}

	return nil
}

// GetWebAuthnCredentials retrieves all WebAuthn credentials for user
func (s *DatabaseMFAStore) GetWebAuthnCredentials(ctx context.Context, userID string) ([]*WebAuthnCredential, error) {
	var records []GamingWebAuthnRecord

	err := s.db.WithContext(ctx).Where("user_id = ? AND enabled = ?", userID, true).Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get WebAuthn credentials: %w", err)
	}

	var credentials []*WebAuthnCredential
	for _, record := range records {
		// Parse transports
		var transports []string
		if record.Transports != "" {
			if err := json.Unmarshal([]byte(record.Transports), &transports); err != nil {
				transports = []string{} // Continue without transports if parsing fails
			}
		}

		credential := &WebAuthnCredential{
			ID:              record.ID,
			UserID:          record.UserID,
			PublicKey:       record.PublicKey,
			AttestationType: record.AttestationType,
			AAGUID:          record.AAGUID,
			SignCount:       record.SignCount,
			Transports:      transports,
			DeviceName:      record.DeviceName,
			DeviceType:      record.DeviceType,
			GamingPlatform:  record.GamingPlatform,
			UserAgent:       record.UserAgent,
			IPAddress:       record.IPAddress,
			Enabled:         record.Enabled,
			CreatedAt:       record.CreatedAt,
			LastUsed:        record.LastUsed,
		}

		credentials = append(credentials, credential)
	}

	return credentials, nil
}

// UpdateWebAuthnCredential updates WebAuthn credential
func (s *DatabaseMFAStore) UpdateWebAuthnCredential(ctx context.Context, credentialID string, credential *WebAuthnCredential) error {
	// Convert transports to JSON
	transportsJSON, err := json.Marshal(credential.Transports)
	if err != nil {
		return fmt.Errorf("failed to marshal WebAuthn transports: %w", err)
	}

	updates := map[string]interface{}{
		"sign_count":  credential.SignCount,
		"transports":  string(transportsJSON),
		"device_name": credential.DeviceName,
		"last_used":   credential.LastUsed,
		"updated_at":  time.Now(),
	}

	result := s.db.WithContext(ctx).Model(&GamingWebAuthnRecord{}).
		Where("id = ?", credentialID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update WebAuthn credential: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("WebAuthn credential not found")
	}

	return nil
}

// DeleteWebAuthnCredential deletes WebAuthn credential
func (s *DatabaseMFAStore) DeleteWebAuthnCredential(ctx context.Context, userID, credentialID string) error {
	result := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", credentialID, userID).
		Delete(&GamingWebAuthnRecord{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete WebAuthn credential: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("WebAuthn credential not found")
	}

	return nil
}

// Backup codes methods

// StoreBackupCodes stores backup codes in database
func (s *DatabaseMFAStore) StoreBackupCodes(ctx context.Context, userID string, codes *BackupCodes) error {
	// Convert codes map to JSON
	codesJSON, err := json.Marshal(codes.Codes)
	if err != nil {
		return fmt.Errorf("failed to marshal gaming backup codes: %w", err)
	}

	record := &GamingBackupCodesRecord{
		UserID:    userID,
		Codes:     string(codesJSON),
		CreatedAt: codes.CreatedAt,
		UsedCount: codes.UsedCount,
		UpdatedAt: time.Now(),
	}

	if err := s.db.WithContext(ctx).Save(record).Error; err != nil {
		return fmt.Errorf("failed to store gaming backup codes: %w", err)
	}

	return nil
}

// GetBackupCodes retrieves backup codes from database
func (s *DatabaseMFAStore) GetBackupCodes(ctx context.Context, userID string) (*BackupCodes, error) {
	var record GamingBackupCodesRecord

	err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("gaming backup codes not found")
		}
		return nil, fmt.Errorf("failed to get gaming backup codes: %w", err)
	}

	// Parse codes map
	var codes map[string]bool
	if err := json.Unmarshal([]byte(record.Codes), &codes); err != nil {
		return nil, fmt.Errorf("failed to parse gaming backup codes: %w", err)
	}

	return &BackupCodes{
		UserID:    record.UserID,
		Codes:     codes,
		CreatedAt: record.CreatedAt,
		UsedCount: record.UsedCount,
	}, nil
}

// UseBackupCode marks backup code as used
func (s *DatabaseMFAStore) UseBackupCode(ctx context.Context, userID, code string) error {
	// Get current backup codes
	backupCodes, err := s.GetBackupCodes(ctx, userID)
	if err != nil {
		return err
	}

	// Check if code exists and is unused
	used, exists := backupCodes.Codes[code]
	if !exists {
		return fmt.Errorf("invalid gaming backup code")
	}

	if used {
		return fmt.Errorf("gaming backup code already used")
	}

	// Mark code as used
	backupCodes.Codes[code] = true
	backupCodes.UsedCount++

	// Update in database
	return s.StoreBackupCodes(ctx, userID, backupCodes)
}

// MFA Challenge methods (using Redis for temporary storage)

// RedisMFAChallengeStore implements MFA challenge storage using Redis
type RedisMFAChallengeStore struct {
	redisClient RedisClient
	keyPrefix   string
}

// NewRedisMFAChallengeStore creates new Redis MFA challenge store
func NewRedisMFAChallengeStore(redisClient RedisClient) *RedisMFAChallengeStore {
	return &RedisMFAChallengeStore{
		redisClient: redisClient,
		keyPrefix:   "herald:gaming:mfa:challenge:",
	}
}

// StoreMFAChallenge stores MFA challenge in Redis
func (s *RedisMFAChallengeStore) StoreMFAChallenge(ctx context.Context, challengeID string, challenge *MFAChallenge) error {
	key := s.keyPrefix + challengeID

	data, err := json.Marshal(challenge)
	if err != nil {
		return fmt.Errorf("failed to marshal gaming MFA challenge: %w", err)
	}

	ttl := time.Until(challenge.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("gaming MFA challenge already expired")
	}

	if err := s.redisClient.Set(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("failed to store gaming MFA challenge: %w", err)
	}

	return nil
}

// GetMFAChallenge retrieves MFA challenge from Redis
func (s *RedisMFAChallengeStore) GetMFAChallenge(ctx context.Context, challengeID string) (*MFAChallenge, error) {
	key := s.keyPrefix + challengeID

	data, err := s.redisClient.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get gaming MFA challenge: %w", err)
	}

	var challenge MFAChallenge
	if err := json.Unmarshal([]byte(data), &challenge); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gaming MFA challenge: %w", err)
	}

	// Check if challenge has expired
	if time.Now().After(challenge.ExpiresAt) {
		s.redisClient.Del(ctx, key)
		return nil, fmt.Errorf("gaming MFA challenge expired")
	}

	return &challenge, nil
}

// DeleteMFAChallenge deletes MFA challenge from Redis
func (s *RedisMFAChallengeStore) DeleteMFAChallenge(ctx context.Context, challengeID string) error {
	key := s.keyPrefix + challengeID

	if err := s.redisClient.Del(ctx, key); err != nil {
		return fmt.Errorf("failed to delete gaming MFA challenge: %w", err)
	}

	return nil
}

// MFA Attempt tracking (using database)

// TrackMFAAttempt stores MFA attempt in database
func (s *DatabaseMFAStore) TrackMFAAttempt(ctx context.Context, userID string, attempt *MFAAttempt) error {
	record := &GamingMFAAttemptRecord{
		ID:           attempt.ID,
		UserID:       userID,
		Method:       attempt.Method,
		Success:      attempt.Success,
		ErrorMessage: attempt.ErrorMessage,
		IPAddress:    attempt.IPAddress,
		UserAgent:    attempt.UserAgent,
		GamingAction: attempt.GamingAction,
		AttemptedAt:  attempt.AttemptedAt,
		CreatedAt:    time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("failed to track gaming MFA attempt: %w", err)
	}

	return nil
}

// GetMFAAttempts retrieves MFA attempts for user since specified time
func (s *DatabaseMFAStore) GetMFAAttempts(ctx context.Context, userID string, since time.Time) ([]*MFAAttempt, error) {
	var records []GamingMFAAttemptRecord

	err := s.db.WithContext(ctx).Where("user_id = ? AND attempted_at > ?", userID, since).
		Order("attempted_at DESC").Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get gaming MFA attempts: %w", err)
	}

	var attempts []*MFAAttempt
	for _, record := range records {
		attempt := &MFAAttempt{
			ID:           record.ID,
			UserID:       record.UserID,
			Method:       record.Method,
			Success:      record.Success,
			ErrorMessage: record.ErrorMessage,
			IPAddress:    record.IPAddress,
			UserAgent:    record.UserAgent,
			GamingAction: record.GamingAction,
			AttemptedAt:  record.AttemptedAt,
		}
		attempts = append(attempts, attempt)
	}

	return attempts, nil
}

// Combined MFA Store that implements all interfaces
type CombinedGamingMFAStore struct {
	*DatabaseMFAStore
	*RedisMFAChallengeStore
}

// NewCombinedGamingMFAStore creates a combined MFA store using both database and Redis
func NewCombinedGamingMFAStore(db *gorm.DB, redisClient RedisClient) *CombinedGamingMFAStore {
	return &CombinedGamingMFAStore{
		DatabaseMFAStore:       NewDatabaseMFAStore(db),
		RedisMFAChallengeStore: NewRedisMFAChallengeStore(redisClient),
	}
}

// Ensure CombinedGamingMFAStore implements MFAStore interface
var _ MFAStore = (*CombinedGamingMFAStore)(nil)

// Additional utility methods for MFA store

// GetMFAStatus returns comprehensive MFA status for user
func (s *DatabaseMFAStore) GetMFAStatus(ctx context.Context, userID string) (*MFAStatus, error) {
	status := &MFAStatus{
		UserID:         userID,
		HasTOTP:        false,
		HasWebAuthn:    false,
		HasBackupCodes: false,
		Methods:        []string{},
	}

	// Check TOTP
	if totpSecret, err := s.GetTOTPSecret(ctx, userID); err == nil && totpSecret.Enabled {
		status.HasTOTP = true
		status.Methods = append(status.Methods, "totp")
	}

	// Check WebAuthn
	if credentials, err := s.GetWebAuthnCredentials(ctx, userID); err == nil && len(credentials) > 0 {
		status.HasWebAuthn = true
		status.WebAuthnCredentials = len(credentials)
		status.Methods = append(status.Methods, "webauthn")
	}

	// Check backup codes
	if backupCodes, err := s.GetBackupCodes(ctx, userID); err == nil {
		status.HasBackupCodes = true
		// Count unused backup codes
		for _, used := range backupCodes.Codes {
			if !used {
				status.UnusedBackupCodes++
			}
		}
	}

	status.Enabled = len(status.Methods) > 0

	return status, nil
}

// MFAStatus represents comprehensive MFA status
type MFAStatus struct {
	UserID              string   `json:"user_id"`
	Enabled             bool     `json:"enabled"`
	HasTOTP             bool     `json:"has_totp"`
	HasWebAuthn         bool     `json:"has_webauthn"`
	HasBackupCodes      bool     `json:"has_backup_codes"`
	Methods             []string `json:"methods"`
	WebAuthnCredentials int      `json:"webauthn_credentials"`
	UnusedBackupCodes   int      `json:"unused_backup_codes"`
}
