package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herald-lol/herald/backend/internal/services"
)

// DamageHandler handles damage analytics requests
type DamageHandler struct {
	damageService *services.DamageAnalyticsService
}

// NewDamageHandler creates a new damage handler
func NewDamageHandler(damageService *services.DamageAnalyticsService) *DamageHandler {
	return &DamageHandler{
		damageService: damageService,
	}
}

// DamageAnalysisRequest represents request for damage analysis
type DamageAnalysisRequest struct {
	TimeRange string `form:"time_range" json:"time_range" binding:"required"`
	Champion  string `form:"champion" json:"champion"`
	Position  string `form:"position" json:"position"`
	GameMode  string `form:"game_mode" json:"game_mode"`
}

// TeamContributionRequest represents request for team contribution analysis
type TeamContributionRequest struct {
	TimeRange string `form:"time_range" json:"time_range" binding:"required"`
	TeamRole  string `form:"team_role" json:"team_role"`
	Season    string `form:"season" json:"season"`
}

// GetDamageAnalysis godoc
// @Summary Get comprehensive damage analysis
// @Description Analyzes player's damage output with team contribution metrics and carry potential
// @Tags damage
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string true "Time range (7d, 30d, 90d)"
// @Param champion query string false "Champion name filter"
// @Param position query string false "Position filter"
// @Param game_mode query string false "Game mode filter"
// @Success 200 {object} services.DamageAnalysis
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/damage/{player_id}/analysis [get]
func (dh *DamageHandler) GetDamageAnalysis(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	var req DamageAnalysisRequest
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

	// Validate game mode if provided
	if req.GameMode != "" {
		validGameModes := map[string]bool{
			"RANKED_SOLO_5x5": true,
			"RANKED_FLEX_SR":  true,
			"NORMAL_DRAFT":    true,
			"ARAM":            true,
			"RANKED_TFT":      true,
		}
		if !validGameModes[req.GameMode] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid game mode",
			})
			return
		}
	}

	// Perform damage analysis
	analysis, err := dh.damageService.AnalyzeDamage(c.Request.Context(), playerID, req.TimeRange, req.Champion, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze damage data",
		})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetTeamContribution godoc
// @Summary Get team contribution metrics
// @Description Analyzes player's contribution to team success through damage metrics
// @Tags damage
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string true "Time range (7d, 30d, 90d)"
// @Param team_role query string false "Team role filter (carry, support, tank, etc.)"
// @Param season query string false "Season filter"
// @Success 200 {object} services.TeamContributionData
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/damage/{player_id}/contribution [get]
func (dh *DamageHandler) GetTeamContribution(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	var req TeamContributionRequest
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

	// Perform damage analysis first to get team contribution
	analysis, err := dh.damageService.AnalyzeDamage(c.Request.Context(), playerID, req.TimeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze damage data for team contribution",
		})
		return
	}

	c.JSON(http.StatusOK, analysis.TeamContribution)
}

// GetDamageDistribution godoc
// @Summary Get damage distribution analysis
// @Description Analyzes damage distribution across different targets and phases
// @Tags damage
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string false "Time range (default: 30d)"
// @Param target_type query string false "Target type filter (champions, structures, monsters)"
// @Success 200 {object} services.DamageDistribution
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/damage/{player_id}/distribution [get]
func (dh *DamageHandler) GetDamageDistribution(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	timeRange := c.DefaultQuery("time_range", "30d")
	targetType := c.Query("target_type")

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Validate target type if provided
	if targetType != "" {
		validTargetTypes := map[string]bool{
			"champions":  true,
			"structures": true,
			"monsters":   true,
			"objectives": true,
		}
		if !validTargetTypes[targetType] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid target type. Use: champions, structures, monsters, or objectives",
			})
			return
		}
	}

	// Get damage analysis
	analysis, err := dh.damageService.AnalyzeDamage(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze damage distribution",
		})
		return
	}

	c.JSON(http.StatusOK, analysis.DamageDistribution)
}

// GetDamageTrends godoc
// @Summary Get damage performance trends
// @Description Returns damage output and efficiency trends over time
// @Tags damage
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param metric query string true "Metric type (damage_per_minute, damage_share, carry_potential, efficiency)"
// @Param period query string false "Period (daily, weekly) - default: daily"
// @Param days query int false "Number of days (default: 30)"
// @Success 200 {object} []services.DamageTrendPoint
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/damage/{player_id}/trends [get]
func (dh *DamageHandler) GetDamageTrends(c *gin.Context) {
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
		"damage_per_minute": true,
		"damage_share":      true,
		"carry_potential":   true,
		"efficiency":        true,
	}
	if !validMetrics[metric] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid metric. Use: damage_per_minute, damage_share, carry_potential, or efficiency",
		})
		return
	}

	period := c.DefaultQuery("period", "daily")
	validPeriods := map[string]bool{"daily": true, "weekly": true}
	if !validPeriods[period] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid period. Use: daily or weekly",
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

	// Get damage analysis to extract trend data
	timeRange := "30d"
	if days <= 7 {
		timeRange = "7d"
	} else if days <= 90 {
		timeRange = "90d"
	}

	analysis, err := dh.damageService.AnalyzeDamage(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to get damage trend data",
		})
		return
	}

	c.JSON(http.StatusOK, analysis.TrendData)
}

// GetDamageComparison godoc
// @Summary Compare damage performance with benchmarks
// @Description Compares player's damage metrics against role/rank/global benchmarks
// @Tags damage
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param benchmark_type query string true "Benchmark type (role, rank, global)"
// @Param time_range query string false "Time range (default: 30d)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/damage/{player_id}/comparison [get]
func (dh *DamageHandler) GetDamageComparison(c *gin.Context) {
	playerID := c.Param("player_id")
	benchmarkType := c.Query("benchmark_type")

	if playerID == "" || benchmarkType == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID and benchmark type are required",
		})
		return
	}

	// Validate benchmark type
	validBenchmarks := map[string]bool{
		"role":   true,
		"rank":   true,
		"global": true,
	}
	if !validBenchmarks[benchmarkType] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid benchmark type. Use: role, rank, or global",
		})
		return
	}

	timeRange := c.DefaultQuery("time_range", "30d")
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Get damage analysis
	analysis, err := dh.damageService.AnalyzeDamage(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze damage data",
		})
		return
	}

	// Extract relevant benchmark comparison
	var comparison map[string]interface{}

	switch benchmarkType {
	case "role":
		comparison = map[string]interface{}{
			"benchmark_type": "role",
			"player_metrics": map[string]interface{}{
				"damage_share":      analysis.DamageShare,
				"damage_per_minute": analysis.DamagePerMinute,
				"carry_potential":   analysis.CarryPotential,
				"efficiency_rating": analysis.EfficiencyRating,
			},
			"benchmark_metrics": analysis.RoleBenchmark,
			"percentile":        analysis.RoleBenchmark.PlayerPercentile,
		}
	case "rank":
		comparison = map[string]interface{}{
			"benchmark_type": "rank",
			"player_metrics": map[string]interface{}{
				"damage_share":      analysis.DamageShare,
				"damage_per_minute": analysis.DamagePerMinute,
				"carry_potential":   analysis.CarryPotential,
				"efficiency_rating": analysis.EfficiencyRating,
			},
			"benchmark_metrics": analysis.RankBenchmark,
			"percentile":        analysis.RankBenchmark.PlayerPercentile,
		}
	case "global":
		comparison = map[string]interface{}{
			"benchmark_type": "global",
			"player_metrics": map[string]interface{}{
				"damage_share":      analysis.DamageShare,
				"damage_per_minute": analysis.DamagePerMinute,
				"carry_potential":   analysis.CarryPotential,
				"efficiency_rating": analysis.EfficiencyRating,
			},
			"benchmark_metrics": analysis.GlobalBenchmark,
			"percentile":        analysis.GlobalBenchmark.PlayerPercentile,
		}
	}

	c.JSON(http.StatusOK, comparison)
}

// GetCarryPotentialAnalysis godoc
// @Summary Get carry potential analysis
// @Description Analyzes player's ability to carry games based on damage output and team impact
// @Tags damage
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string false "Time range (default: 30d)"
// @Param win_condition query string false "Win condition analysis (damage_carry, utility_carry, mixed)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/damage/{player_id}/carry-potential [get]
func (dh *DamageHandler) GetCarryPotentialAnalysis(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	timeRange := c.DefaultQuery("time_range", "30d")
	winCondition := c.Query("win_condition")

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Validate win condition if provided
	if winCondition != "" {
		validConditions := map[string]bool{
			"damage_carry":  true,
			"utility_carry": true,
			"mixed":         true,
		}
		if !validConditions[winCondition] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid win condition. Use: damage_carry, utility_carry, or mixed",
			})
			return
		}
	}

	// Get damage analysis
	analysis, err := dh.damageService.AnalyzeDamage(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze carry potential",
		})
		return
	}

	// Create carry potential response
	carryAnalysis := map[string]interface{}{
		"player_id":         playerID,
		"time_range":        timeRange,
		"carry_potential":   analysis.CarryPotential,
		"damage_share":      analysis.DamageShare,
		"consistency_score": analysis.DamageConsistency,
		"team_contribution": analysis.TeamContribution,
		"game_phase_analysis": map[string]interface{}{
			"early_game_carry": analysis.GamePhaseAnalysis.EarlyGame.CarryPotential,
			"mid_game_carry":   analysis.GamePhaseAnalysis.MidGame.CarryPotential,
			"late_game_carry":  analysis.GamePhaseAnalysis.LateGame.CarryPotential,
		},
		"win_rate_correlation": map[string]interface{}{
			"high_damage_games": analysis.HighDamageWinRate,
			"low_damage_games":  analysis.LowDamageWinRate,
			"impact_score":      analysis.CarryPotential,
		},
		"recommendations": analysis.Recommendations,
	}

	c.JSON(http.StatusOK, carryAnalysis)
}

// Register routes for damage analytics
func (dh *DamageHandler) RegisterRoutes(router *gin.RouterGroup) {
	damage := router.Group("/damage")
	{
		damage.GET("/:player_id/analysis", dh.GetDamageAnalysis)
		damage.GET("/:player_id/contribution", dh.GetTeamContribution)
		damage.GET("/:player_id/distribution", dh.GetDamageDistribution)
		damage.GET("/:player_id/trends", dh.GetDamageTrends)
		damage.GET("/:player_id/comparison", dh.GetDamageComparison)
		damage.GET("/:player_id/carry-potential", dh.GetCarryPotentialAnalysis)
	}
}
