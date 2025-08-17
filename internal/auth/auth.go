package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Session represents an authenticated user session
type Session struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// AuthService handles authentication operations
type AuthService struct {
	sessions      map[string]*Session
	sessionsMutex sync.RWMutex
	tokenExpiry   time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService() *AuthService {
	service := &AuthService{
		sessions:    make(map[string]*Session),
		tokenExpiry: time.Hour * 24, // 24 hours
	}
	
	// Start cleanup routine
	go service.cleanupExpiredSessions()
	
	return service
}

// HashPassword creates a bcrypt hash of the password
func (as *AuthService) HashPassword(password string) (string, error) {
	cost := 12 // Strong cost for security
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword checks if password matches hash
func (as *AuthService) VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// CreateSession creates a new session for a user
func (as *AuthService) CreateSession(userID int, username, email string) (*Session, error) {
	// Generate random token
	token, err := GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	session := &Session{
		UserID:    userID,
		Username:  username,
		Email:     email,
		Token:     token,
		ExpiresAt: time.Now().Add(as.tokenExpiry),
		CreatedAt: time.Now(),
	}

	as.sessionsMutex.Lock()
	as.sessions[token] = session
	as.sessionsMutex.Unlock()

	return session, nil
}

// ValidateSession validates a session token and returns the session
func (as *AuthService) ValidateSession(token string) (*Session, error) {
	as.sessionsMutex.RLock()
	session, exists := as.sessions[token]
	as.sessionsMutex.RUnlock()

	if !exists {
		return nil, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		as.deleteSession(token)
		return nil, errors.New("session expired")
	}

	return session, nil
}

// DeleteSession removes a session
func (as *AuthService) DeleteSession(token string) {
	as.deleteSession(token)
}

func (as *AuthService) deleteSession(token string) {
	as.sessionsMutex.Lock()
	delete(as.sessions, token)
	as.sessionsMutex.Unlock()
}

// RefreshSession extends the expiry of an existing session
func (as *AuthService) RefreshSession(token string) (*Session, error) {
	as.sessionsMutex.Lock()
	defer as.sessionsMutex.Unlock()

	session, exists := as.sessions[token]
	if !exists {
		return nil, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		delete(as.sessions, token)
		return nil, errors.New("session expired")
	}

	// Extend expiry
	session.ExpiresAt = time.Now().Add(as.tokenExpiry)
	return session, nil
}

// cleanupExpiredSessions periodically removes expired sessions
func (as *AuthService) cleanupExpiredSessions() {
	ticker := time.NewTicker(time.Hour) // Cleanup every hour
	defer ticker.Stop()

	for range ticker.C {
		as.sessionsMutex.Lock()
		now := time.Now()
		for token, session := range as.sessions {
			if now.After(session.ExpiresAt) {
				delete(as.sessions, token)
			}
		}
		as.sessionsMutex.Unlock()
	}
}

// GenerateRandomString generates a cryptographically secure random string
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// PasswordStrength checks password strength
func (as *AuthService) PasswordStrength(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	
	// Could add more sophisticated checks here
	// - uppercase/lowercase mix
	// - numbers
	// - special characters
	// - common password detection
	
	return nil
}

// GenerateTokens creates JWT-like tokens for OAuth authentication
func (as *AuthService) GenerateTokens(userID int) (accessToken, refreshToken string, err error) {
	// Generate access token (shorter lived)
	accessToken, err = GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}
	
	// Generate refresh token (longer lived)
	refreshToken, err = GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}
	
	// Store session with access token
	session := &Session{
		UserID:    userID,
		Token:     accessToken,
		ExpiresAt: time.Now().Add(time.Hour), // 1 hour for access token
		CreatedAt: time.Now(),
	}
	
	as.sessionsMutex.Lock()
	as.sessions[accessToken] = session
	as.sessionsMutex.Unlock()
	
	return accessToken, refreshToken, nil
}
