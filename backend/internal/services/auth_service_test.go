package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/herald-lol/herald/backend/internal/config"
	"github.com/herald-lol/herald/backend/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Run migrations
	err = db.AutoMigrate(
		&models.User{},
		&models.RiotAccount{},
		&models.UserPreferences{},
		&models.Subscription{},
	)
	require.NoError(t, err)

	return db
}

func setupTestConfig() *config.Config {
	return &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-jwt-secret-key",
			Expiration: 24 * time.Hour,
		},
		Server: config.ServerConfig{
			Environment: "test",
		},
	}
}

func TestAuthService_Register(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()
	authService := NewAuthService(db, cfg)

	tests := []struct {
		name    string
		request RegisterRequest
		wantErr bool
		errType error
	}{
		{
			name: "successful registration",
			request: RegisterRequest{
				Email:       "test@herald.lol",
				Username:    "testuser",
				Password:    "password123",
				DisplayName: "Test User",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			request: RegisterRequest{
				Email:    "test@herald.lol",
				Username: "anotheruser",
				Password: "password123",
			},
			wantErr: true,
			errType: ErrUserAlreadyExists,
		},
		{
			name: "duplicate username",
			request: RegisterRequest{
				Email:    "another@herald.lol",
				Username: "testuser",
				Password: "password123",
			},
			wantErr: true,
			errType: ErrUserAlreadyExists,
		},
		{
			name: "weak password",
			request: RegisterRequest{
				Email:    "weak@herald.lol",
				Username: "weakuser",
				Password: "123", // Too short
			},
			wantErr: true,
			errType: ErrWeakPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Register(tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, response)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.Token)
				assert.NotEmpty(t, response.RefreshToken)
				assert.Equal(t, tt.request.Email, response.User.Email)
				assert.Equal(t, tt.request.Username, response.User.Username)
				assert.True(t, response.User.IsActive)
				assert.False(t, response.User.IsPremium)

				// Verify user preferences were created
				assert.NotNil(t, response.User.Preferences)
				assert.Equal(t, "dark", response.User.Preferences.Theme)

				// Verify subscription was created
				assert.NotNil(t, response.User.Subscription)
				assert.Equal(t, "free", response.User.Subscription.Plan)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()
	authService := NewAuthService(db, cfg)

	// Create a test user first
	registerReq := RegisterRequest{
		Email:    "login@herald.lol",
		Username: "loginuser",
		Password: "password123",
	}
	_, err := authService.Register(registerReq)
	require.NoError(t, err)

	tests := []struct {
		name    string
		request LoginRequest
		wantErr bool
		errType error
	}{
		{
			name: "successful login",
			request: LoginRequest{
				Email:    "login@herald.lol",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "wrong password",
			request: LoginRequest{
				Email:    "login@herald.lol",
				Password: "wrongpassword",
			},
			wantErr: true,
			errType: ErrInvalidCredentials,
		},
		{
			name: "non-existent user",
			request: LoginRequest{
				Email:    "nonexistent@herald.lol",
				Password: "password123",
			},
			wantErr: true,
			errType: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Login(tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, response)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.Token)
				assert.NotEmpty(t, response.RefreshToken)
				assert.Equal(t, tt.request.Email, response.User.Email)

				// Verify login count was incremented
				assert.Equal(t, 1, response.User.LoginCount)
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()
	authService := NewAuthService(db, cfg)

	// Create a test user and get a token
	registerReq := RegisterRequest{
		Email:    "token@herald.lol",
		Username: "tokenuser",
		Password: "password123",
	}
	response, err := authService.Register(registerReq)
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
		errType error
	}{
		{
			name:    "valid token",
			token:   response.Token,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid-token",
			wantErr: true,
			errType: ErrInvalidToken,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
			errType: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := authService.ValidateToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, response.User.Email, user.Email)
			}
		})
	}
}

func TestAuthService_ChangePassword(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()
	authService := NewAuthService(db, cfg)

	// Create a test user
	registerReq := RegisterRequest{
		Email:    "changepw@herald.lol",
		Username: "changepwuser",
		Password: "password123",
	}
	response, err := authService.Register(registerReq)
	require.NoError(t, err)

	tests := []struct {
		name            string
		userID          uuid.UUID
		currentPassword string
		newPassword     string
		wantErr         bool
		errType         error
	}{
		{
			name:            "successful password change",
			userID:          response.User.ID,
			currentPassword: "password123",
			newPassword:     "newpassword123",
			wantErr:         false,
		},
		{
			name:            "wrong current password",
			userID:          response.User.ID,
			currentPassword: "wrongpassword",
			newPassword:     "newpassword123",
			wantErr:         true,
			errType:         ErrInvalidCredentials,
		},
		{
			name:            "weak new password",
			userID:          response.User.ID,
			currentPassword: "password123",
			newPassword:     "123", // Too short
			wantErr:         true,
			errType:         ErrWeakPassword,
		},
		{
			name:            "non-existent user",
			userID:          uuid.New(),
			currentPassword: "password123",
			newPassword:     "newpassword123",
			wantErr:         true,
			errType:         ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.ChangePassword(tt.userID, tt.currentPassword, tt.newPassword)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)

				// Verify we can login with new password
				loginReq := LoginRequest{
					Email:    "changepw@herald.lol",
					Password: tt.newPassword,
				}
				loginResponse, err := authService.Login(loginReq)
				assert.NoError(t, err)
				assert.NotNil(t, loginResponse)
			}
		})
	}
}

// Benchmark tests
func BenchmarkAuthService_Register(b *testing.B) {
	db := setupTestDB(&testing.T{})
	cfg := setupTestConfig()
	authService := NewAuthService(db, cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		request := RegisterRequest{
			Email:    "bench" + string(rune(i)) + "@herald.lol",
			Username: "benchuser" + string(rune(i)),
			Password: "password123",
		}
		_, _ = authService.Register(request)
	}
}

func BenchmarkAuthService_Login(b *testing.B) {
	db := setupTestDB(&testing.T{})
	cfg := setupTestConfig()
	authService := NewAuthService(db, cfg)

	// Create test user
	registerReq := RegisterRequest{
		Email:    "benchmark@herald.lol",
		Username: "benchuser",
		Password: "password123",
	}
	_, err := authService.Register(registerReq)
	if err != nil {
		b.Fatal(err)
	}

	loginReq := LoginRequest{
		Email:    "benchmark@herald.lol",
		Password: "password123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = authService.Login(loginReq)
	}
}

func BenchmarkAuthService_ValidateToken(b *testing.B) {
	db := setupTestDB(&testing.T{})
	cfg := setupTestConfig()
	authService := NewAuthService(db, cfg)

	// Create test user and get token
	registerReq := RegisterRequest{
		Email:    "validate@herald.lol",
		Username: "validateuser",
		Password: "password123",
	}
	response, err := authService.Register(registerReq)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = authService.ValidateToken(response.Token)
	}
}
