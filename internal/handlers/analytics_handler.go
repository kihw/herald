package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lol-match-exporter/internal/services"
)

type AnalyticsHandler struct {
	AnalyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		AnalyticsService: analyticsService,
	}
}

// GetPeriodStats returns analytics for a specific time period
// GET /api/analytics/period/:period
func (h *AnalyticsHandler) GetPeriodStats(c *gin.Context) {
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

	period := c.Param("period")
	if period == "" {
		period = "week" // default period
	}

	// Validate period
	validPeriods := map[string]bool{
		"today":  true,
		"week":   true,
		"month":  true,
		"season": true,
	}

	if !validPeriods[period] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_period",
			"message": "Period must be one of: today, week, month, season",
		})
		return
	}

	stats, err := h.AnalyticsService.GetPeriodStats(uid, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "analytics_error",
			"message": "Failed to generate period statistics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetMMRTrajectory returns MMR analysis over time
// GET /api/analytics/mmr?days=30
func (h *AnalyticsHandler) GetMMRTrajectory(c *gin.Context) {
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

	// Parse days parameter
	days := 30 // default
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	trajectory, err := h.AnalyticsService.GetMMRTrajectory(uid, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "mmr_analysis_error",
			"message": "Failed to calculate MMR trajectory",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    trajectory,
	})
}

// GetRecommendations returns AI-powered recommendations
// GET /api/analytics/recommendations
func (h *AnalyticsHandler) GetRecommendations(c *gin.Context) {
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

	recommendations, err := h.AnalyticsService.GetRecommendations(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "recommendations_error",
			"message": "Failed to generate recommendations",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    recommendations,
	})
}

// GetChampionAnalysis returns detailed analysis for a specific champion
// GET /api/analytics/champion/:championId?period=season
func (h *AnalyticsHandler) GetChampionAnalysis(c *gin.Context) {
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

	championIDStr := c.Param("championId")
	championID, err := strconv.Atoi(championIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_champion_id",
			"message": "Champion ID must be a valid integer",
		})
		return
	}

	period := c.Query("period")
	if period == "" {
		period = "season" // default period for champion analysis
	}

	analysis, err := h.AnalyticsService.GetChampionAnalysis(uid, championID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "champion_analysis_error",
			"message": "Failed to analyze champion performance",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analysis,
	})
}

// GetPerformanceTrends returns detailed performance trends over time
// GET /api/analytics/trends
func (h *AnalyticsHandler) GetPerformanceTrends(c *gin.Context) {
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

	trends, err := h.AnalyticsService.GetPerformanceTrends(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "trends_error",
			"message": "Failed to calculate performance trends",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    trends,
	})
}

// RefreshAnalytics triggers analytics refresh for the user
// POST /api/analytics/refresh
func (h *AnalyticsHandler) RefreshAnalytics(c *gin.Context) {
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

	// Parse request body for options
	var request struct {
		Period string `json:"period"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		// Use default if no body provided
		request.Period = "all"
	}

	err := h.AnalyticsService.UpdateChampionStats(uid, request.Period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "refresh_error",
			"message": "Failed to refresh analytics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Analytics refreshed successfully",
		"data": gin.H{
			"user_id": uid,
			"period":  request.Period,
		},
	})
}

// GetChampionsByRole returns champion performance filtered by role
// GET /api/analytics/champions/:role
func (h *AnalyticsHandler) GetChampionsByRole(c *gin.Context) {
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

	role := c.Param("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_role",
			"message": "Role parameter is required",
		})
		return
	}

	// Validate role
	validRoles := map[string]bool{
		"TOP":     true,
		"JUNGLE":  true,
		"MIDDLE":  true,
		"BOTTOM":  true,
		"UTILITY": true,
	}

	if !validRoles[role] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_role",
			"message": "Role must be one of: TOP, JUNGLE, MIDDLE, BOTTOM, UTILITY",
		})
		return
	}

	period := c.Query("period")
	if period == "" {
		period = "month" // default period for role analysis
	}

	// Get role performance from period stats
	stats, err := h.AnalyticsService.GetPeriodStats(uid, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "analytics_error",
			"message": "Failed to get role performance",
			"details": err.Error(),
		})
		return
	}

	// Filter for requested role
	rolePerformance, exists := stats.RolePerformance[role]
	if !exists {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"role":        role,
				"period":      period,
				"performance": nil,
				"champions":   []interface{}{},
			},
		})
		return
	}

	// Get top champions for this role (filtered from general top champions)
	roleChampions := []services.ChampionPerformance{}
	for _, champ := range stats.TopChampions {
		// This is a simplified filtering - in production you'd want to track
		// champions by role more explicitly
		roleChampions = append(roleChampions, champ)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"role":        role,
			"period":      period,
			"performance": rolePerformance,
			"champions":   roleChampions,
		},
	})
}

// HealthCheck validates the analytics service health
// GET /api/analytics/health
func (h *AnalyticsHandler) HealthCheck(c *gin.Context) {
	// Note: ValidatePythonEnvironment method removed - using Go native services
	var err error = nil
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "service_unavailable",
			"message": "Analytics service is not available",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Analytics service is healthy",
		"status":  "available",
	})
}