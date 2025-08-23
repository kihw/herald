package services

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	
	"github.com/herald-lol/backend/internal/config"
	"github.com/herald-lol/backend/internal/models"
)

type AuthService struct {
	db     *gorm.DB
	config *config.Config
}

type AuthClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	IsPremium bool     `json:"is_premium"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required,min=3,max=20"`
	Password    string `json:"password" binding:"required,min=6"`
	DisplayName string `json:"display_name,omitempty"`
}

type AuthResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         models.User `json:"user"`
	ExpiresIn    int         `json:"expires_in"`
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrWeakPassword       = errors.New("password is too weak")
)

func NewAuthService(db *gorm.DB, config *config.Config) *AuthService {
	return &AuthService{
		db:     db,
		config: config,
	}
}

// Register creates a new user account
func (s *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		return nil, ErrUserAlreadyExists
	}

	// Validate password strength
	if err := s.validatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		DisplayName:  req.DisplayName,
		IsActive:     true,
		IsPremium:    false,
		LoginCount:   0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create user
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create default preferences
	preferences := models.UserPreferences{
		UserID:                      user.ID,
		Theme:                       "dark",
		CompactMode:                 false,
		ShowDetailedStats:           true,
		DefaultTimeframe:           "7d",
		EmailNotifications:         true,
		PushNotifications:          true,
		MatchNotifications:         true,
		RankChangeNotifications:    true,
		AutoSyncMatches:            true,
		SyncInterval:               300,
		IncludeNormalGames:         true,
		IncludeARAMGames:           true,
		PublicProfile:              true,
		ShowInLeaderboards:         true,
		AllowDataExport:            true,
		ReceiveAICoaching:          true,
		SkillLevel:                 "intermediate",
		PreferredCoachingStyle:     "balanced",
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
	}

	if err := tx.Create(&preferences).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create free subscription
	subscription := models.Subscription{
		UserID:             user.ID,
		Plan:               "free",
		Status:             "active",
		StartedAt:          time.Now(),
		ExpiresAt:          time.Now().AddDate(100, 0, 0), // Free plan never expires
		Amount:             0,
		Currency:           "USD",
		Interval:           "monthly",
		MaxRiotAccounts:    1,
		UnlimitedAnalytics: false,
		AICoachingAccess:   false,
		AdvancedMetrics:    false,
		DataExportAccess:   false,
		PrioritySupport:    false,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := tx.Create(&subscription).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	// Load user with relationships
	if err := s.db.Preload("RiotAccounts").Preload("Preferences").Preload("Subscription").First(&user, user.ID).Error; err != nil {
		return nil, err
	}

	// Generate tokens
	token, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
		ExpiresIn:    int(s.config.JWT.Expiration.Seconds()),
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	var user models.User

	// Find user by email
	if err := s.db.Preload("RiotAccounts").Preload("Preferences").Preload("Subscription").Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update login stats
	user.LastLogin = time.Now()
	user.LoginCount++
	s.db.Save(&user)

	// Generate tokens
	token, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
		ExpiresIn:    int(s.config.JWT.Expiration.Seconds()),
	}, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	// Get user from database
	var user models.User
	if err := s.db.Preload("RiotAccounts").Preload("Preferences").Preload("Subscription").First(&user, claims.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	// Parse refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Get user
	var user models.User
	if err := s.db.Preload("RiotAccounts").Preload("Preferences").Preload("Subscription").First(&user, claims.UserID).Error; err != nil {
		return nil, ErrUserNotFound
	}

	// Generate new tokens
	newToken, newRefreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
		User:         user,
		ExpiresIn:    int(s.config.JWT.Expiration.Seconds()),
	}, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(userID uuid.UUID, currentPassword, newPassword string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return ErrUserNotFound
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// Validate new password strength
	if err := s.validatePasswordStrength(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()
	
	return s.db.Save(&user).Error
}

// ResetPassword initiates password reset (would typically send email)
func (s *AuthService) ResetPassword(email string) error {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		// Don't reveal if email exists or not for security
		return nil
	}

	// TODO: Generate reset token and send email
	// For now, just log that a reset was requested
	
	return nil
}

// generateTokens generates access and refresh tokens for a user
func (s *AuthService) generateTokens(user models.User) (string, string, error) {
	now := time.Now()
	
	// Access token claims
	accessClaims := AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		IsPremium: user.IsPremium,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.JWT.Expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   user.ID.String(),
			Issuer:    "herald.lol",
		},
	}

	// Refresh token claims (longer expiration)
	refreshClaims := AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		IsPremium: user.IsPremium,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.JWT.Expiration * 7)), // 7x longer
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   user.ID.String(),
			Issuer:    "herald.lol",
		},
	}

	// Generate tokens
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// validatePasswordStrength validates password strength
func (s *AuthService) validatePasswordStrength(password string) error {
	if len(password) < 6 {
		return ErrWeakPassword
	}
	
	// TODO: Add more password strength validation
	// - Check for common passwords
	// - Require mix of letters, numbers, symbols
	// - Check against breach databases
	
	return nil
}