package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lol-match-exporter/internal/models"
	"lol-match-exporter/internal/services"
)

type DashboardHandler struct {
	UserValidationService *services.UserValidationService
	// TODO: Add RiotService and SyncService when ready
}

func NewDashboardHandler(userValidationService *services.UserValidationService) *DashboardHandler {
	return &DashboardHandler{
		UserValidationService: userValidationService,
	}
}

// GetStats returns dashboard statistics for the authenticated user
func (h *DashboardHandler) GetStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user_not_found",
			"message": "User ID not found in context",
		})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_user_id",
			"message": "Invalid user ID format",
		})
		return
	}

	// Get user information
	user, err := h.UserValidationService.GetUserByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "user_not_found",
			"message": "User not found",
		})
		return
	}

	// For now, return basic stats - will be expanded with match data later
	stats := map[string]interface{}{
		"user": user,
		"stats": map[string]interface{}{
			"total_matches":     0,
			"win_rate":         0.0,
			"average_kda":      0.0,
			"favorite_champion": "Unknown",
			"last_sync_at":     user.LastSync,
			"next_sync_at":     nil,
		},
		"recent_matches": []interface{}{},
		"summary": map[string]interface{}{
			"account_validated": user.IsValidated,
			"region":           user.Region,
			"summoner_level":   user.SummonerLevel,
		},
	}

	c.JSON(http.StatusOK, stats)
}

// GetMatches returns match history for the authenticated user
func (h *DashboardHandler) GetMatches(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user_not_found",
			"message": "User ID not found in context",
		})
		return
	}

	// Parse pagination parameters
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// For now, return empty matches - will be implemented with match synchronization
	matches := models.PaginatedResponse{
		Data:       []interface{}{},
		Page:       page,
		PerPage:    limit,
		Total:      0,
		TotalPages: 0,
	}

	c.JSON(http.StatusOK, matches)
}

// SyncMatches initiates match synchronization for the authenticated user
func (h *DashboardHandler) SyncMatches(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user_not_found",
			"message": "User ID not found in context",
		})
		return
	}

	_, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_user_id",
			"message": "Invalid user ID format",
		})
		return
	}

	// For now, return a placeholder response
	// TODO: Implement actual match synchronization with Riot API
	response := models.SyncResponse{
		JobID:   0,
		Status:  "placeholder",
		Message: "Match synchronization not yet implemented. Coming soon!",
	}

	c.JSON(http.StatusOK, response)
}

// GetSettings returns user settings
func (h *DashboardHandler) GetSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user_not_found",
			"message": "User ID not found in context",
		})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_user_id",
			"message": "Invalid user ID format",
		})
		return
	}

	// Get user information
	user, err := h.UserValidationService.GetUserByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "user_not_found",
			"message": "User not found",
		})
		return
	}

	// Return default settings - will be expanded later
	settings := map[string]interface{}{
		"user_id":               user.ID,
		"platform":              user.Region,
		"queue_types":           []int{420, 440}, // Ranked Solo/Duo, Ranked Flex
		"language":              "en_US",
		"include_timeline":      true,
		"include_all_data":      true,
		"light_mode":            true,
		"auto_sync_enabled":     true,
		"sync_frequency_hours":  24,
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateSettings updates user settings
func (h *DashboardHandler) UpdateSettings(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user_not_found",
			"message": "User ID not found in context",
		})
		return
	}

	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid settings format",
			"details": err.Error(),
		})
		return
	}

	// For now, just return success - will implement actual settings update later
	c.JSON(http.StatusOK, gin.H{
		"message": "Settings updated successfully",
		"settings": settings,
	})
}
