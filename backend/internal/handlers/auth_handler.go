package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/herald-lol/herald/backend/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.RegisterRequest true "Registration details"
// @Success 201 {object} services.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	response, err := h.authService.Register(req)
	if err != nil {
		switch err {
		case services.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "User already exists",
				Message: "A user with this email or username already exists",
			})
		case services.ErrWeakPassword:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Weak password",
				Message: "Password must be at least 6 characters long",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Registration failed",
				Message: "An error occurred during registration",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login handles user authentication
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.LoginRequest true "Login credentials"
// @Success 200 {object} services.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Invalid credentials",
				Message: "Email or password is incorrect",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Login failed",
				Message: "An error occurred during login",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} services.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		switch err {
		case services.ErrInvalidToken:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Invalid token",
				Message: "Refresh token is invalid or expired",
			})
		case services.ErrUserNotFound:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "User not found",
				Message: "User associated with token not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Token refresh failed",
				Message: "An error occurred while refreshing token",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// ChangePassword handles password change
// @Summary Change user password
// @Description Change the current user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Password change details"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in token",
		})
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	err := h.authService.ChangePassword(userID.(uuid.UUID), req.CurrentPassword, req.NewPassword)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid current password",
				Message: "Current password is incorrect",
			})
		case services.ErrWeakPassword:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Weak password",
				Message: "New password is too weak",
			})
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "User not found",
				Message: "User not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Password change failed",
				Message: "An error occurred while changing password",
			})
		}
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// ResetPassword handles password reset request
// @Summary Request password reset
// @Description Request a password reset email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Email for password reset"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	err := h.authService.ResetPassword(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Reset password failed",
			Message: "An error occurred while processing password reset",
		})
		return
	}

	// Always return success for security (don't reveal if email exists)
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "If an account with that email exists, a password reset email has been sent",
	})
}

// GetProfile returns the current user's profile
// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} ErrorResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not found in context",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Logout handles user logout (typically just returns success as JWT is stateless)
// @Summary Logout user
// @Description Logout the current user (client should discard tokens)
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a stateless JWT setup, logout is handled client-side by discarding tokens
	// In a more sophisticated setup, you might maintain a blacklist of tokens

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}

// AuthMiddleware validates JWT tokens and adds user to context
func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Missing authorization header",
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check for Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Invalid authorization header",
				Message: "Authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		user, err := h.authService.ValidateToken(tokenString)
		if err != nil {
			switch err {
			case services.ErrInvalidToken:
				c.JSON(http.StatusUnauthorized, ErrorResponse{
					Error:   "Invalid token",
					Message: "Token is invalid",
				})
			case services.ErrTokenExpired:
				c.JSON(http.StatusUnauthorized, ErrorResponse{
					Error:   "Token expired",
					Message: "Token has expired",
				})
			case services.ErrUserNotFound:
				c.JSON(http.StatusUnauthorized, ErrorResponse{
					Error:   "User not found",
					Message: "User associated with token not found",
				})
			default:
				c.JSON(http.StatusUnauthorized, ErrorResponse{
					Error:   "Authentication failed",
					Message: "Failed to authenticate user",
				})
			}
			c.Abort()
			return
		}

		// Add user to context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("username", user.Username)
		c.Set("email", user.Email)
		c.Set("is_premium", user.IsPremium)

		c.Next()
	})
}

// OptionalAuthMiddleware validates JWT tokens but doesn't require them
func (h *AuthHandler) OptionalAuthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			user, err := h.authService.ValidateToken(tokenString)
			if err == nil && user != nil {
				c.Set("user", user)
				c.Set("user_id", user.ID)
				c.Set("username", user.Username)
				c.Set("email", user.Email)
				c.Set("is_premium", user.IsPremium)
			}
		}

		c.Next()
	})
}

// Common response types
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
