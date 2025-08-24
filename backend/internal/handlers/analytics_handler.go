package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herald-lol/herald/backend/internal/services"
)

// AnalyticsHandler handles analytics-related requests
type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// KDAAnalysisRequest represents request for KDA analysis
type KDAAnalysisRequest struct {
	TimeRange string `form:"time_range" json:"time_range" binding:"required"` // "7d", "30d", "90d"
	Champion  string `form:"champion" json:"champion"`                        // Optional champion filter
}

// CSAnalysisRequest represents request for CS analysis
type CSAnalysisRequest struct {
	TimeRange string `form:"time_range" json:"time_range" binding:"required"`
	Position  string `form:"position" json:"position"` // Optional position filter
	Champion  string `form:"champion" json:"champion"` // Optional champion filter
}

// PerformanceComparisonRequest represents request for performance comparison
type PerformanceComparisonRequest struct {
	TimeRange string `form:"time_range" json:"time_range" binding:"required"`
}

// GetKDAAnalysis godoc
// @Summary Get comprehensive KDA analysis
// @Description Analyzes player's KDA performance with trends, comparisons, and insights
// @Tags analytics
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string true "Time range (7d, 30d, 90d)"
// @Param champion query string false "Champion name filter"
// @Success 200 {object} services.KDAAnalysis
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/analytics/{player_id}/kda [get]
func (ah *AnalyticsHandler) GetKDAAnalysis(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	var req KDAAnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Validate time range
	if !isValidTimeRange(req.TimeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Perform KDA analysis
	analysis, err := ah.analyticsService.AnalyzeKDA(c.Request.Context(), playerID, req.TimeRange, req.Champion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze KDA data",
		})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetCSAnalysis godoc
// @Summary Get comprehensive CS (Creep Score) analysis
// @Description Analyzes player's farming performance with benchmarks and recommendations
// @Tags analytics
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string true "Time range (7d, 30d, 90d)"
// @Param position query string false "Position filter (TOP, JUNGLE, MID, ADC, SUPPORT)"
// @Param champion query string false "Champion name filter"
// @Success 200 {object} services.CSAnalysis
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/analytics/{player_id}/cs [get]
func (ah *AnalyticsHandler) GetCSAnalysis(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	var req CSAnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Validate time range
	if !isValidTimeRange(req.TimeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Validate position if provided
	if req.Position != "" && !isValidPosition(req.Position) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid position. Use: TOP, JUNGLE, MID, ADC, or SUPPORT",
		})
		return
	}

	// Perform CS analysis
	analysis, err := ah.analyticsService.AnalyzeCS(c.Request.Context(), playerID, req.TimeRange, req.Position, req.Champion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze CS data",
		})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetPerformanceComparison godoc
// @Summary Get comprehensive performance comparison
// @Description Compares player performance against rank, role, and global benchmarks
// @Tags analytics
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string true "Time range (7d, 30d, 90d)"
// @Success 200 {object} services.PerformanceComparison
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/analytics/{player_id}/comparison [get]
func (ah *AnalyticsHandler) GetPerformanceComparison(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	var req PerformanceComparisonRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Validate time range
	if !isValidTimeRange(req.TimeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Perform performance comparison
	comparison, err := ah.analyticsService.ComparePerformance(c.Request.Context(), playerID, req.TimeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "comparison_error",
			Message: "Failed to generate performance comparison",
		})
		return
	}

	c.JSON(http.StatusOK, comparison)
}

// GetPlayerStats godoc
// @Summary Get aggregated player statistics
// @Description Returns comprehensive player statistics for the specified time period
// @Tags analytics
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param days query int false "Number of days (default: 30)"
// @Success 200 {object} models.PlayerStats
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/analytics/{player_id}/stats [get]
func (ah *AnalyticsHandler) GetPlayerStats(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	// Parse days parameter
	days := 30 // default
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
			days = d
		} else {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid days parameter. Must be between 1 and 365",
			})
			return
		}
	}

	// Get player statistics
	stats, err := ah.analyticsService.GetPlayerStats(c.Request.Context(), playerID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "stats_error",
			Message: "Failed to get player statistics",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetChampionStats godoc
// @Summary Get champion-specific statistics
// @Description Returns performance statistics for a specific champion
// @Tags analytics
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param champion_id path int true "Champion ID"
// @Param days query int false "Number of days (default: 30)"
// @Success 200 {object} models.ChampionStats
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/analytics/{player_id}/champion/{champion_id} [get]
func (ah *AnalyticsHandler) GetChampionStats(c *gin.Context) {
	playerID := c.Param("player_id")
	championIDStr := c.Param("champion_id")

	if playerID == "" || championIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID and Champion ID are required",
		})
		return
	}

	championID, err := strconv.Atoi(championIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid champion ID",
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

	// Get champion statistics
	stats, err := ah.analyticsService.GetChampionStats(c.Request.Context(), playerID, championID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "stats_error",
			Message: "Failed to get champion statistics",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetPerformanceTrends godoc
// @Summary Get performance trends over time
// @Description Returns performance trends and patterns for visualization
// @Tags analytics
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param metric query string true "Metric type (kda, cs, vision, damage, winrate)"
// @Param period query string false "Time period (daily, weekly, monthly) - default: daily"
// @Param days query int false "Number of days (default: 30)"
// @Success 200 {object} []services.KDATrendPoint
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/analytics/{player_id}/trends [get]
func (ah *AnalyticsHandler) GetPerformanceTrends(c *gin.Context) {
	playerID := c.Param("player_id")
	metric := c.Query("metric")

	if playerID == "" || metric == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID and metric are required",
		})
		return
	}

	// Validate metric
	validMetrics := map[string]bool{
		"kda":     true,
		"cs":      true,
		"vision":  true,
		"damage":  true,
		"winrate": true,
	}

	if !validMetrics[metric] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid metric. Use: kda, cs, vision, damage, or winrate",
		})
		return
	}

	period := c.DefaultQuery("period", "daily")
	validPeriods := map[string]bool{
		"daily":   true,
		"weekly":  true,
		"monthly": true,
	}

	if !validPeriods[period] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid period. Use: daily, weekly, or monthly",
		})
		return
	}

	// Parse days parameter
	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	// Get performance trends
	trends, err := ah.analyticsService.GetPerformanceTrends(c.Request.Context(), playerID, metric, period, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "trends_error",
			Message: "Failed to get performance trends",
		})
		return
	}

	c.JSON(http.StatusOK, trends)
}

// GetBenchmarks godoc
// @Summary Get performance benchmarks
// @Description Returns benchmark data for comparison (rank, role, global averages)
// @Tags analytics
// @Accept json
// @Produce json
// @Param tier query string false "Tier filter (IRON, BRONZE, SILVER, etc.)"
// @Param role query string false "Role filter (TOP, JUNGLE, MID, ADC, SUPPORT)"
// @Param champion query string false "Champion name filter"
// @Success 200 {object} services.PerformanceBenchmark
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/analytics/benchmarks [get]
func (ah *AnalyticsHandler) GetBenchmarks(c *gin.Context) {
	tier := c.Query("tier")
	role := c.Query("role")
	champion := c.Query("champion")

	// Validate tier if provided
	if tier != "" && !isValidTier(tier) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid tier",
		})
		return
	}

	// Validate role if provided
	if role != "" && !isValidPosition(role) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid role",
		})
		return
	}

	// Get benchmarks
	benchmarks, err := ah.analyticsService.GetBenchmarks(c.Request.Context(), tier, role, champion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "benchmarks_error",
			Message: "Failed to get performance benchmarks",
		})
		return
	}

	c.JSON(http.StatusOK, benchmarks)
}

// Validation helper functions

func isValidTimeRange(timeRange string) bool {
	validRanges := map[string]bool{
		"7d":  true,
		"30d": true,
		"90d": true,
	}
	return validRanges[timeRange]
}

func isValidPosition(position string) bool {
	validPositions := map[string]bool{
		"TOP":     true,
		"JUNGLE":  true,
		"MID":     true,
		"ADC":     true,
		"SUPPORT": true,
	}
	return validPositions[position]
}

func isValidTier(tier string) bool {
	validTiers := map[string]bool{
		"IRON":        true,
		"BRONZE":      true,
		"SILVER":      true,
		"GOLD":        true,
		"PLATINUM":    true,
		"EMERALD":     true,
		"DIAMOND":     true,
		"MASTER":      true,
		"GRANDMASTER": true,
		"CHALLENGER":  true,
	}
	return validTiers[tier]
}

// Register routes
func (ah *AnalyticsHandler) RegisterRoutes(router *gin.RouterGroup) {
	analytics := router.Group("/analytics")
	{
		// Player-specific analytics
		analytics.GET("/:player_id/kda", ah.GetKDAAnalysis)
		analytics.GET("/:player_id/cs", ah.GetCSAnalysis)
		analytics.GET("/:player_id/comparison", ah.GetPerformanceComparison)
		analytics.GET("/:player_id/stats", ah.GetPlayerStats)
		analytics.GET("/:player_id/trends", ah.GetPerformanceTrends)
		analytics.GET("/:player_id/champion/:champion_id", ah.GetChampionStats)

		// Global analytics
		analytics.GET("/benchmarks", ah.GetBenchmarks)
	}
}
