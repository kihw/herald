package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"lol-match-exporter/internal/models"
	"lol-match-exporter/internal/services"
)

type ValidationHandler struct {
	UserValidationService *services.UserValidationService
	RiotValidationService *services.RiotValidationService
}

func NewValidationHandler(
	userValidationService *services.UserValidationService,
	riotValidationService *services.RiotValidationService,
) *ValidationHandler {
	return &ValidationHandler{
		UserValidationService: userValidationService,
		RiotValidationService: riotValidationService,
	}
}

// ValidateAccount validates a Riot account and creates/logs in the user
func (h *ValidationHandler) ValidateAccount(c *gin.Context) {
	var req models.RiotAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate required fields
	if req.RiotID == "" || req.RiotTag == "" || req.Region == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_fields",
			"message": "Riot ID, Tag, and Region are required",
		})
		return
	}

	// Validate and create/update user
	user, err := h.UserValidationService.ValidateAndCreateUser(req.RiotID, req.RiotTag, req.Region)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_failed",
			"message": "Failed to validate Riot account",
			"details": err.Error(),
		})
		return
	}

	// Create session
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("authenticated", true)
	session.Set("login_time", time.Now().Unix())
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "session_error",
			"message": "Failed to create session",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, models.RiotAccountValidationResponse{
		Valid: true,
		User:  user,
	})
}

// GetProfile returns the current user's profile
func (h *ValidationHandler) GetProfile(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "not_authenticated",
			"message": "User not authenticated",
		})
		return
	}

	user, err := h.UserValidationService.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "user_not_found",
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Logout logs out the current user
func (h *ValidationHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "logout_error",
			"message": "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// CheckSession checks if user has a valid session
func (h *ValidationHandler) CheckSession(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	authenticated := session.Get("authenticated")
	
	if userID != nil && authenticated == true {
		user, err := h.UserValidationService.GetUserByID(userID.(int))
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"authenticated": true,
				"user":          user,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": false,
	})
}

// GetSupportedRegions returns list of supported regions
func (h *ValidationHandler) GetSupportedRegions(c *gin.Context) {
	regions := h.RiotValidationService.GetSupportedRegions()
	c.JSON(http.StatusOK, gin.H{
		"regions": regions,
	})
}

// RequireAuth middleware to require authentication
func (h *ValidationHandler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		authenticated := session.Get("authenticated")
		
		if userID == nil || authenticated != true {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "authentication_required",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		// Store user ID in context for use in handlers
		c.Set("user_id", userID.(int))
		c.Next()
	}
}

// GetCurrentUserID helper to get current user ID from context
func GetCurrentUserID(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, gin.Error{Err: gin.Error{}.Err, Type: gin.ErrorTypePublic}
	}
	
	switch v := userID.(type) {
	case int:
		return v, nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, gin.Error{Err: gin.Error{}.Err, Type: gin.ErrorTypePublic}
	}
}
