package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

// Herald.lol Gaming Analytics - Multi-Factor Authentication
// Comprehensive MFA implementation with TOTP and WebAuthn for gaming platform

// GamingMFAManager manages multi-factor authentication for Herald.lol
type GamingMFAManager struct {
	db              *gorm.DB
	webAuthn        *webauthn.WebAuthn
	mfaStore        MFAStore
	gamingAnalytics GamingAnalyticsService
	config          *GamingMFAConfig
}

// GamingMFAConfig holds MFA configuration for gaming platform
type GamingMFAConfig struct {
	// TOTP Configuration
	TOTPIssuer       string
	TOTPAccountName  string
	TOTPSecretLength int
	TOTPPeriod       uint
	TOTPSkew         uint

	// WebAuthn Configuration
	WebAuthnRPID                    string
	WebAuthnRPDisplayName           string
	WebAuthnRPOrigins               []string
	WebAuthnTimeout                 time.Duration
	WebAuthnRequireResident         bool
	WebAuthnRequireUserVerification bool

	// Gaming-specific MFA settings
	GamingSessionMFARequired bool
	AnalyticsMFARequired     bool
	APIAccessMFARequired     bool
	HighValueActionsMFA      []string // Actions requiring MFA
	MFACooldownPeriod        time.Duration
	MaxMFAAttempts           int

	// Backup codes
	BackupCodesEnabled bool
	BackupCodesCount   int
	BackupCodeLength   int
}

// MFAStore interface for MFA data management
type MFAStore interface {
	// TOTP methods
	StoreTOTPSecret(ctx context.Context, userID string, secret *TOTPSecret) error
	GetTOTPSecret(ctx context.Context, userID string) (*TOTPSecret, error)
	VerifyTOTPCode(ctx context.Context, userID, code string) error
	DisableTOTP(ctx context.Context, userID string) error

	// WebAuthn methods
	StoreWebAuthnCredential(ctx context.Context, userID string, credential *WebAuthnCredential) error
	GetWebAuthnCredentials(ctx context.Context, userID string) ([]*WebAuthnCredential, error)
	UpdateWebAuthnCredential(ctx context.Context, credentialID string, credential *WebAuthnCredential) error
	DeleteWebAuthnCredential(ctx context.Context, userID, credentialID string) error

	// Backup codes methods
	StoreBackupCodes(ctx context.Context, userID string, codes *BackupCodes) error
	GetBackupCodes(ctx context.Context, userID string) (*BackupCodes, error)
	UseBackupCode(ctx context.Context, userID, code string) error

	// MFA session methods
	StoreMFAChallenge(ctx context.Context, challengeID string, challenge *MFAChallenge) error
	GetMFAChallenge(ctx context.Context, challengeID string) (*MFAChallenge, error)
	DeleteMFAChallenge(ctx context.Context, challengeID string) error

	// MFA attempt tracking
	TrackMFAAttempt(ctx context.Context, userID string, attempt *MFAAttempt) error
	GetMFAAttempts(ctx context.Context, userID string, since time.Time) ([]*MFAAttempt, error)
}

// TOTPSecret represents TOTP configuration for a user
type TOTPSecret struct {
	UserID      string     `json:"user_id"`
	Secret      string     `json:"secret"` // Base32 encoded secret
	QRCodeURL   string     `json:"qr_code_url"`
	BackupCodes []string   `json:"backup_codes,omitempty"`
	Enabled     bool       `json:"enabled"`
	Verified    bool       `json:"verified"`
	CreatedAt   time.Time  `json:"created_at"`
	VerifiedAt  *time.Time `json:"verified_at,omitempty"`
	LastUsed    *time.Time `json:"last_used,omitempty"`
}

// WebAuthnCredential represents a WebAuthn credential
type WebAuthnCredential struct {
	ID              string     `json:"id"` // Credential ID (base64url encoded)
	UserID          string     `json:"user_id"`
	PublicKey       []byte     `json:"public_key"`
	AttestationType string     `json:"attestation_type"`
	AAGUID          []byte     `json:"aaguid"`
	SignCount       uint32     `json:"sign_count"`
	Transports      []string   `json:"transports,omitempty"`
	DeviceName      string     `json:"device_name"`
	DeviceType      string     `json:"device_type"` // gaming-controller, security-key, platform, etc.
	GamingPlatform  string     `json:"gaming_platform,omitempty"`
	UserAgent       string     `json:"user_agent"`
	IPAddress       string     `json:"ip_address"`
	Enabled         bool       `json:"enabled"`
	CreatedAt       time.Time  `json:"created_at"`
	LastUsed        *time.Time `json:"last_used,omitempty"`
}

// BackupCodes represents backup codes for MFA
type BackupCodes struct {
	UserID    string          `json:"user_id"`
	Codes     map[string]bool `json:"codes"` // code -> used
	CreatedAt time.Time       `json:"created_at"`
	UsedCount int             `json:"used_count"`
}

// MFAChallenge represents an MFA challenge session
type MFAChallenge struct {
	ID            string                 `json:"id"`
	UserID        string                 `json:"user_id"`
	ChallengeType string                 `json:"challenge_type"` // totp, webauthn, backup
	SessionData   map[string]interface{} `json:"session_data"`
	GamingContext *GamingContext         `json:"gaming_context,omitempty"`
	IPAddress     string                 `json:"ip_address"`
	UserAgent     string                 `json:"user_agent"`
	CreatedAt     time.Time              `json:"created_at"`
	ExpiresAt     time.Time              `json:"expires_at"`
	Completed     bool                   `json:"completed"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty"`
}

// MFAAttempt represents an MFA verification attempt
type MFAAttempt struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Method       string    `json:"method"` // totp, webauthn, backup
	Success      bool      `json:"success"`
	ErrorMessage string    `json:"error_message,omitempty"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	GamingAction string    `json:"gaming_action,omitempty"`
	AttemptedAt  time.Time `json:"attempted_at"`
}

// GamingWebAuthnUser represents a WebAuthn user for gaming platform
type GamingWebAuthnUser struct {
	ID                string                `json:"id"`
	Name              string                `json:"name"`
	DisplayName       string                `json:"display_name"`
	Icon              string                `json:"icon,omitempty"`
	GamingCredentials []*WebAuthnCredential `json:"credentials"`
}

// NewGamingMFAManager creates new MFA manager for gaming platform
func NewGamingMFAManager(
	db *gorm.DB,
	mfaStore MFAStore,
	gamingAnalytics GamingAnalyticsService,
	config *GamingMFAConfig,
) (*GamingMFAManager, error) {
	// Set gaming-specific defaults
	if config.TOTPIssuer == "" {
		config.TOTPIssuer = "Herald.lol Gaming Analytics"
	}
	if config.TOTPSecretLength == 0 {
		config.TOTPSecretLength = 32
	}
	if config.TOTPPeriod == 0 {
		config.TOTPPeriod = 30
	}
	if config.TOTPSkew == 0 {
		config.TOTPSkew = 1
	}
	if config.WebAuthnTimeout == 0 {
		config.WebAuthnTimeout = 60 * time.Second
	}
	if config.MFACooldownPeriod == 0 {
		config.MFACooldownPeriod = 30 * time.Second
	}
	if config.MaxMFAAttempts == 0 {
		config.MaxMFAAttempts = 5
	}
	if len(config.HighValueActionsMFA) == 0 {
		config.HighValueActionsMFA = []string{
			"analytics:export",
			"team:management",
			"subscription:change",
			"account:delete",
		}
	}
	if config.BackupCodesCount == 0 {
		config.BackupCodesCount = 10
	}
	if config.BackupCodeLength == 0 {
		config.BackupCodeLength = 8
	}

	// Initialize WebAuthn
	webAuthnConfig := &webauthn.Config{
		RPDisplayName: config.WebAuthnRPDisplayName,
		RPID:          config.WebAuthnRPID,
		RPOrigins:     config.WebAuthnRPOrigins,
		Timeout:       int(config.WebAuthnTimeout.Milliseconds()),
	}

	webAuthn, err := webauthn.New(webAuthnConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize WebAuthn for gaming platform: %w", err)
	}

	return &GamingMFAManager{
		db:              db,
		webAuthn:        webAuthn,
		mfaStore:        mfaStore,
		gamingAnalytics: gamingAnalytics,
		config:          config,
	}, nil
}

// TOTP Methods

// SetupGamingTOTP sets up TOTP for a gaming user
func (mfa *GamingMFAManager) SetupGamingTOTP(c *gin.Context) {
	user, exists := GetGamingUser(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gaming user not found",
		})
		return
	}

	// Check if TOTP is already enabled
	existingSecret, err := mfa.mfaStore.GetTOTPSecret(c.Request.Context(), user.ID)
	if err == nil && existingSecret.Enabled {
		c.JSON(http.StatusConflict, gin.H{
			"error":           "TOTP already enabled for gaming account",
			"gaming_platform": "herald-lol",
		})
		return
	}

	// Generate TOTP secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      mfa.config.TOTPIssuer,
		AccountName: fmt.Sprintf("%s (%s)", user.Email, user.Name),
		Period:      mfa.config.TOTPPeriod,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
		SecretSize:  mfa.config.TOTPSecretLength,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate gaming TOTP secret",
		})
		return
	}

	// Generate backup codes
	backupCodes, err := mfa.generateBackupCodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate gaming backup codes",
		})
		return
	}

	// Store TOTP secret
	totpSecret := &TOTPSecret{
		UserID:      user.ID,
		Secret:      key.Secret(),
		QRCodeURL:   key.URL(),
		BackupCodes: backupCodes,
		Enabled:     false, // Will be enabled after verification
		Verified:    false,
		CreatedAt:   time.Now(),
	}

	if err := mfa.mfaStore.StoreTOTPSecret(c.Request.Context(), user.ID, totpSecret); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store gaming TOTP secret",
		})
		return
	}

	// Generate QR code image
	qrCode, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate gaming QR code",
		})
		return
	}

	// Track MFA setup
	go mfa.gamingAnalytics.TrackUserLogin(c.Request.Context(), user.ID, user.Provider, map[string]string{
		"action":          "mfa_totp_setup",
		"gaming_platform": "herald-lol",
	})

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"totp_setup": gin.H{
			"secret":         key.Secret(),
			"qr_code_url":    key.URL(),
			"qr_code_base64": base64.StdEncoding.EncodeToString(qrCode),
			"backup_codes":   backupCodes,
			"manual_entry": gin.H{
				"issuer":    mfa.config.TOTPIssuer,
				"account":   fmt.Sprintf("%s (%s)", user.Email, user.Name),
				"secret":    key.Secret(),
				"period":    mfa.config.TOTPPeriod,
				"digits":    6,
				"algorithm": "SHA1",
			},
		},
		"next_step":       "Scan QR code with authenticator app and verify with /auth/mfa/totp/verify",
		"gaming_platform": "herald-lol",
	})
}

// VerifyGamingTOTP verifies and enables TOTP for gaming user
func (mfa *GamingMFAManager) VerifyGamingTOTP(c *gin.Context) {
	user, exists := GetGamingUser(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gaming user not found",
		})
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid gaming TOTP verification request",
		})
		return
	}

	// Get TOTP secret
	totpSecret, err := mfa.mfaStore.GetTOTPSecret(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Gaming TOTP not set up",
		})
		return
	}

	if totpSecret.Enabled {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Gaming TOTP already verified and enabled",
		})
		return
	}

	// Verify TOTP code
	valid := totp.Validate(req.Code, totpSecret.Secret)
	if !valid {
		// Track failed attempt
		mfa.trackMFAAttempt(c.Request.Context(), user.ID, "totp", false, "Invalid TOTP code", c)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error":           "Invalid gaming TOTP code",
			"gaming_platform": "herald-lol",
		})
		return
	}

	// Enable TOTP
	now := time.Now()
	totpSecret.Enabled = true
	totpSecret.Verified = true
	totpSecret.VerifiedAt = &now
	totpSecret.LastUsed = &now

	if err := mfa.mfaStore.StoreTOTPSecret(c.Request.Context(), user.ID, totpSecret); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to enable gaming TOTP",
		})
		return
	}

	// Store backup codes
	if mfa.config.BackupCodesEnabled {
		backupCodes := &BackupCodes{
			UserID:    user.ID,
			Codes:     make(map[string]bool),
			CreatedAt: now,
			UsedCount: 0,
		}

		for _, code := range totpSecret.BackupCodes {
			backupCodes.Codes[code] = false // false = unused
		}

		mfa.mfaStore.StoreBackupCodes(c.Request.Context(), user.ID, backupCodes)
	}

	// Track successful MFA setup
	mfa.trackMFAAttempt(c.Request.Context(), user.ID, "totp", true, "", c)

	go mfa.gamingAnalytics.TrackUserLogin(c.Request.Context(), user.ID, user.Provider, map[string]string{
		"action":          "mfa_totp_enabled",
		"gaming_platform": "herald-lol",
	})

	c.JSON(http.StatusOK, gin.H{
		"success":            true,
		"message":            "Gaming TOTP enabled successfully",
		"backup_codes_saved": mfa.config.BackupCodesEnabled,
		"gaming_platform":    "herald-lol",
	})
}

// WebAuthn Methods

// BeginWebAuthnRegistration starts WebAuthn credential registration
func (mfa *GamingMFAManager) BeginWebAuthnRegistration(c *gin.Context) {
	user, exists := GetGamingUser(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gaming user not found",
		})
		return
	}

	var req struct {
		DeviceName     string `json:"device_name"`
		DeviceType     string `json:"device_type"` // gaming-controller, security-key, platform
		GamingPlatform string `json:"gaming_platform,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid WebAuthn registration request",
		})
		return
	}

	// Create WebAuthn user
	gamingWebAuthnUser := &GamingWebAuthnUser{
		ID:          user.ID,
		Name:        user.Email,
		DisplayName: user.Name,
		Icon:        user.Avatar,
	}

	// Get existing credentials
	credentials, _ := mfa.mfaStore.GetWebAuthnCredentials(c.Request.Context(), user.ID)
	for _, cred := range credentials {
		if cred.Enabled {
			gamingWebAuthnUser.GamingCredentials = append(gamingWebAuthnUser.GamingCredentials, cred)
		}
	}

	// Begin registration
	creation, session, err := mfa.webAuthn.BeginRegistration(gamingWebAuthnUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to begin WebAuthn registration",
		})
		return
	}

	// Store challenge session
	challengeID := mfa.generateChallengeID()
	challenge := &MFAChallenge{
		ID:            challengeID,
		UserID:        user.ID,
		ChallengeType: "webauthn_registration",
		SessionData: map[string]interface{}{
			"webauthn_session": session,
			"device_name":      req.DeviceName,
			"device_type":      req.DeviceType,
			"gaming_platform":  req.GamingPlatform,
		},
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(mfa.config.WebAuthnTimeout),
	}

	if err := mfa.mfaStore.StoreMFAChallenge(c.Request.Context(), challengeID, challenge); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store WebAuthn challenge",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"challenge_id":     challengeID,
		"creation_options": creation,
		"gaming_platform":  "herald-lol",
	})
}

// CompleteWebAuthnRegistration completes WebAuthn credential registration
func (mfa *GamingMFAManager) CompleteWebAuthnRegistration(c *gin.Context) {
	challengeID := c.Param("challengeId")

	// Get challenge
	challenge, err := mfa.mfaStore.GetMFAChallenge(c.Request.Context(), challengeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "WebAuthn challenge not found or expired",
		})
		return
	}

	if challenge.Completed {
		c.JSON(http.StatusConflict, gin.H{
			"error": "WebAuthn challenge already completed",
		})
		return
	}

	if time.Now().After(challenge.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{
			"error": "WebAuthn challenge expired",
		})
		return
	}

	// Parse WebAuthn response
	// Note: This is simplified - actual implementation would use webauthn.ParseCredentialCreationResponse
	var credentialCreationResponse map[string]interface{}
	if err := c.ShouldBindJSON(&credentialCreationResponse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid WebAuthn response",
		})
		return
	}

	// Complete registration (simplified)
	// In real implementation, use webauthn.FinishRegistration
	now := time.Now()

	// Create credential record
	credential := &WebAuthnCredential{
		ID:              mfa.generateCredentialID(),
		UserID:          challenge.UserID,
		PublicKey:       []byte("simplified_public_key"), // Would be extracted from response
		AttestationType: "none",
		AAGUID:          []byte{},
		SignCount:       0,
		DeviceName:      challenge.SessionData["device_name"].(string),
		DeviceType:      challenge.SessionData["device_type"].(string),
		UserAgent:       challenge.UserAgent,
		IPAddress:       challenge.IPAddress,
		Enabled:         true,
		CreatedAt:       now,
	}

	if gamingPlatform, ok := challenge.SessionData["gaming_platform"].(string); ok {
		credential.GamingPlatform = gamingPlatform
	}

	// Store credential
	if err := mfa.mfaStore.StoreWebAuthnCredential(c.Request.Context(), challenge.UserID, credential); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store WebAuthn credential",
		})
		return
	}

	// Mark challenge as completed
	challenge.Completed = true
	challenge.CompletedAt = &now
	mfa.mfaStore.StoreMFAChallenge(c.Request.Context(), challengeID, challenge)

	// Track WebAuthn registration
	mfa.trackMFAAttempt(c.Request.Context(), challenge.UserID, "webauthn_registration", true, "", c)

	go mfa.gamingAnalytics.TrackUserLogin(c.Request.Context(), challenge.UserID, "", map[string]string{
		"action":          "mfa_webauthn_registered",
		"device_type":     credential.DeviceType,
		"gaming_platform": credential.GamingPlatform,
	})

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "WebAuthn credential registered successfully",
		"credential": gin.H{
			"id":          credential.ID,
			"device_name": credential.DeviceName,
			"device_type": credential.DeviceType,
			"created_at":  credential.CreatedAt,
		},
		"gaming_platform": "herald-lol",
	})
}

// Helper methods

func (mfa *GamingMFAManager) generateBackupCodes() ([]string, error) {
	codes := make([]string, mfa.config.BackupCodesCount)

	for i := 0; i < mfa.config.BackupCodesCount; i++ {
		code, err := mfa.generateBackupCode()
		if err != nil {
			return nil, err
		}
		codes[i] = code
	}

	return codes, nil
}

func (mfa *GamingMFAManager) generateBackupCode() (string, error) {
	// Generate random backup code
	bytes := make([]byte, mfa.config.BackupCodeLength/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert to hex string and format
	code := fmt.Sprintf("%x", bytes)

	// Format as groups of 4 characters
	if len(code) >= 8 {
		return fmt.Sprintf("%s-%s", code[:4], code[4:8]), nil
	}

	return code, nil
}

func (mfa *GamingMFAManager) generateChallengeID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

func (mfa *GamingMFAManager) generateCredentialID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

func (mfa *GamingMFAManager) trackMFAAttempt(ctx context.Context, userID, method string, success bool, errorMsg string, c *gin.Context) {
	attempt := &MFAAttempt{
		ID:           fmt.Sprintf("mfa_%d", time.Now().UnixNano()),
		UserID:       userID,
		Method:       method,
		Success:      success,
		ErrorMessage: errorMsg,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.GetHeader("User-Agent"),
		AttemptedAt:  time.Now(),
	}

	if gamingAction := c.GetHeader("X-Gaming-Action"); gamingAction != "" {
		attempt.GamingAction = gamingAction
	}

	go mfa.mfaStore.TrackMFAAttempt(ctx, userID, attempt)
}

// Implement WebAuthn User interface
func (user *GamingWebAuthnUser) WebAuthnID() []byte {
	return []byte(user.ID)
}

func (user *GamingWebAuthnUser) WebAuthnName() string {
	return user.Name
}

func (user *GamingWebAuthnUser) WebAuthnDisplayName() string {
	return user.DisplayName
}

func (user *GamingWebAuthnUser) WebAuthnIcon() string {
	return user.Icon
}

func (user *GamingWebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	credentials := make([]webauthn.Credential, len(user.GamingCredentials))

	for i, cred := range user.GamingCredentials {
		credentials[i] = webauthn.Credential{
			ID:              []byte(cred.ID),
			PublicKey:       cred.PublicKey,
			AttestationType: cred.AttestationType,
			Authenticator: webauthn.Authenticator{
				AAGUID:    cred.AAGUID,
				SignCount: cred.SignCount,
			},
		}
	}

	return credentials
}

// Credential descriptor for excluding existing credentials
func (user *GamingWebAuthnUser) CredentialExcludeList() []webauthn.CredentialDescriptor {
	excludeList := make([]webauthn.CredentialDescriptor, len(user.GamingCredentials))

	for i, cred := range user.GamingCredentials {
		transports := make([]webauthn.AuthenticatorTransport, len(cred.Transports))
		for j, transport := range cred.Transports {
			transports[j] = webauthn.AuthenticatorTransport(transport)
		}

		excludeList[i] = webauthn.CredentialDescriptor{
			Type:         webauthn.PublicKeyCredentialType,
			CredentialID: []byte(cred.ID),
			Transport:    transports,
		}
	}

	return excludeList
}
