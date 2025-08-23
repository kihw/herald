package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/herald/internal/services"
)

// ChampionHandler handles champion-specific analytics requests
type ChampionHandler struct {
	championService *services.ChampionAnalyticsService
}

// NewChampionHandler creates a new champion handler
func NewChampionHandler(championService *services.ChampionAnalyticsService) *ChampionHandler {
	return &ChampionHandler{
		championService: championService,
	}
}

// ChampionAnalysisRequest represents request for champion analysis
type ChampionAnalysisRequest struct {
	TimeRange string `form:"time_range" json:"time_range" binding:"required"`
	Champion  string `form:"champion" json:"champion" binding:"required"`
	Position  string `form:"position" json:"position"`
}

// ChampionMasteryRequest represents request for champion mastery ranking
type ChampionMasteryRequest struct {
	TimeRange string `form:"time_range" json:"time_range" binding:"required"`
	Limit     int    `form:"limit" json:"limit"`
}

// ChampionComparisonRequest represents request for champion comparison
type ChampionComparisonRequest struct {
	Champions []string `form:"champions" json:"champions" binding:"required"`
	TimeRange string   `form:"time_range" json:"time_range" binding:"required"`
}

// GetChampionAnalysis godoc
// @Summary Get comprehensive champion-specific performance analysis
// @Description Analyzes player's performance on a specific champion including mechanics, builds, matchups
// @Tags champion
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string true "Time range (7d, 30d, 90d)"
// @Param champion query string true "Champion name"
// @Param position query string false "Position filter"
// @Success 200 {object} services.ChampionAnalysis
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/champion/{player_id}/analysis [get]
func (ch *ChampionHandler) GetChampionAnalysis(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	var req ChampionAnalysisRequest
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

	// Validate champion name
	if req.Champion == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Champion name is required",
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

	// Perform champion analysis
	analysis, err := ch.championService.AnalyzeChampion(c.Request.Context(), playerID, req.Champion, req.TimeRange, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze champion performance",
		})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetChampionMastery godoc
// @Summary Get champion mastery ranking
// @Description Returns player's champion mastery ranking sorted by performance
// @Tags champion
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string true "Time range (7d, 30d, 90d)"
// @Param limit query int false "Limit number of champions returned (default: 10)"
// @Success 200 {object} []services.ChampionMasteryRanking
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/champion/{player_id}/mastery [get]
func (ch *ChampionHandler) GetChampionMastery(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	var req ChampionMasteryRequest
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

	// Set default limit
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 50 {
		req.Limit = 50 // Maximum limit
	}

	// Get champion mastery ranking
	rankings, err := ch.championService.GetChampionMasteryRanking(c.Request.Context(), playerID, req.TimeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to get champion mastery rankings",
		})
		return
	}

	// Apply limit
	if len(rankings) > req.Limit {
		rankings = rankings[:req.Limit]
	}

	c.JSON(http.StatusOK, rankings)
}

// GetChampionComparison godoc
// @Summary Compare performance across multiple champions
// @Description Returns comparative analysis of player's performance on multiple champions
// @Tags champion
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param champions query string true "Comma-separated list of champion names"
// @Param time_range query string true "Time range (7d, 30d, 90d)"
// @Success 200 {object} services.ChampionComparisonData
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/champion/{player_id}/comparison [get]
func (ch *ChampionHandler) GetChampionComparison(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	championsStr := c.Query("champions")
	timeRange := c.Query("time_range")

	if championsStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Champions parameter is required",
		})
		return
	}

	if timeRange == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Time range is required",
		})
		return
	}

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Parse champions list
	champions := strings.Split(championsStr, ",")
	if len(champions) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "At least one champion is required",
		})
		return
	}

	// Limit number of champions for comparison
	if len(champions) > 10 {
		champions = champions[:10]
	}

	// Trim whitespace from champion names
	for i, champion := range champions {
		champions[i] = strings.TrimSpace(champion)
		if champions[i] == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Champion names cannot be empty",
			})
			return
		}
	}

	// Get champion comparison
	comparison, err := ch.championService.GetChampionComparison(c.Request.Context(), playerID, champions, timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to get champion comparison",
		})
		return
	}

	c.JSON(http.StatusOK, comparison)
}

// GetChampionPowerSpikes godoc
// @Summary Get champion power spike analysis
// @Description Analyzes when the player's champion performance peaks during games
// @Tags champion
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param champion query string true "Champion name"
// @Param time_range query string false "Time range (default: 30d)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/champion/{player_id}/power-spikes [get]
func (ch *ChampionHandler) GetChampionPowerSpikes(c *gin.Context) {
	playerID := c.Param("player_id")
	champion := c.Query("champion")
	timeRange := c.DefaultQuery("time_range", "30d")

	if playerID == "" || champion == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID and champion are required",
		})
		return
	}

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Get full champion analysis (we'll extract power spikes from it)
	analysis, err := ch.championService.AnalyzeChampion(c.Request.Context(), playerID, champion, timeRange, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze champion power spikes",
		})
		return
	}

	// Create power spikes response
	powerSpikes := map[string]interface{}{
		"player_id":    playerID,
		"champion":     champion,
		"time_range":   timeRange,
		"power_spikes": analysis.PowerSpikes,
		"carry_potential": analysis.CarryPotential,
		"scaling_analysis": map[string]interface{}{
			"early_game_rating": analysis.LanePhasePerformance.PhaseRating,
			"mid_game_rating":   analysis.MidGamePerformance.PhaseRating,
			"late_game_rating":  analysis.LateGamePerformance.PhaseRating,
		},
	}

	c.JSON(http.StatusOK, powerSpikes)
}

// GetChampionMatchups godoc
// @Summary Get champion matchup analysis
// @Description Analyzes player's performance against specific opponent champions
// @Tags champion
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param champion query string true "Champion name"
// @Param opponent query string false "Opponent champion filter"
// @Param time_range query string false "Time range (default: 30d)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/champion/{player_id}/matchups [get]
func (ch *ChampionHandler) GetChampionMatchups(c *gin.Context) {
	playerID := c.Param("player_id")
	champion := c.Query("champion")
	opponent := c.Query("opponent")
	timeRange := c.DefaultQuery("time_range", "30d")

	if playerID == "" || champion == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID and champion are required",
		})
		return
	}

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Get full champion analysis
	analysis, err := ch.championService.AnalyzeChampion(c.Request.Context(), playerID, champion, timeRange, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze champion matchups",
		})
		return
	}

	// Create matchup response
	matchups := map[string]interface{}{
		"player_id":           playerID,
		"champion":            champion,
		"time_range":          timeRange,
		"matchup_performance": analysis.MatchupPerformance,
		"strength_matchups":   analysis.StrengthMatchups,
		"weakness_matchups":   analysis.WeaknessMatchups,
	}

	// Filter by specific opponent if provided
	if opponent != "" {
		var specificMatchup *services.MatchupData
		
		// Look in strength matchups
		for _, matchup := range analysis.StrengthMatchups {
			if strings.EqualFold(matchup.OpponentChampion, opponent) {
				specificMatchup = &matchup
				break
			}
		}
		
		// Look in weakness matchups if not found
		if specificMatchup == nil {
			for _, matchup := range analysis.WeaknessMatchups {
				if strings.EqualFold(matchup.OpponentChampion, opponent) {
					specificMatchup = &matchup
					break
				}
			}
		}
		
		if specificMatchup != nil {
			matchups["specific_matchup"] = specificMatchup
		} else {
			matchups["specific_matchup"] = map[string]string{
				"message": "No data available for the specified matchup",
			}
		}
	}

	c.JSON(http.StatusOK, matchups)
}

// GetChampionBuilds godoc
// @Summary Get champion build analysis
// @Description Analyzes player's item builds and rune setups for a specific champion
// @Tags champion
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param champion query string true "Champion name"
// @Param build_type query string false "Build type (items, runes, skills)"
// @Param time_range query string false "Time range (default: 30d)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/champion/{player_id}/builds [get]
func (ch *ChampionHandler) GetChampionBuilds(c *gin.Context) {
	playerID := c.Param("player_id")
	champion := c.Query("champion")
	buildType := c.Query("build_type")
	timeRange := c.DefaultQuery("time_range", "30d")

	if playerID == "" || champion == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID and champion are required",
		})
		return
	}

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Validate build type if provided
	if buildType != "" {
		validBuildTypes := map[string]bool{
			"items":  true,
			"runes":  true,
			"skills": true,
		}
		if !validBuildTypes[buildType] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid build type. Use: items, runes, or skills",
			})
			return
		}
	}

	// Get full champion analysis
	analysis, err := ch.championService.AnalyzeChampion(c.Request.Context(), playerID, champion, timeRange, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze champion builds",
		})
		return
	}

	// Create builds response
	builds := map[string]interface{}{
		"player_id":  playerID,
		"champion":   champion,
		"time_range": timeRange,
		"item_builds": analysis.ItemBuilds,
		"skill_order": analysis.SkillOrder,
		"rune_optimization": analysis.RuneOptimization,
	}

	// Filter by specific build type if provided
	if buildType != "" {
		switch buildType {
		case "items":
			builds["focus_data"] = analysis.ItemBuilds
		case "runes":
			builds["focus_data"] = analysis.RuneOptimization
		case "skills":
			builds["focus_data"] = analysis.SkillOrder
		}
	}

	c.JSON(http.StatusOK, builds)
}

// GetChampionCoaching godoc
// @Summary Get champion-specific coaching recommendations
// @Description Provides personalized coaching recommendations for improving champion performance
// @Tags champion
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param champion query string true "Champion name"
// @Param focus_area query string false "Focus area (mechanics, builds, matchups, positioning)"
// @Param time_range query string false "Time range (default: 30d)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/champion/{player_id}/coaching [get]
func (ch *ChampionHandler) GetChampionCoaching(c *gin.Context) {
	playerID := c.Param("player_id")
	champion := c.Query("champion")
	focusArea := c.Query("focus_area")
	timeRange := c.DefaultQuery("time_range", "30d")

	if playerID == "" || champion == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID and champion are required",
		})
		return
	}

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Validate focus area if provided
	if focusArea != "" {
		validFocusAreas := map[string]bool{
			"mechanics":    true,
			"builds":       true,
			"matchups":     true,
			"positioning":  true,
			"team_fighting": true,
		}
		if !validFocusAreas[focusArea] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid focus area. Use: mechanics, builds, matchups, positioning, or team_fighting",
			})
			return
		}
	}

	// Get full champion analysis
	analysis, err := ch.championService.AnalyzeChampion(c.Request.Context(), playerID, champion, timeRange, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to generate champion coaching recommendations",
		})
		return
	}

	// Create coaching response
	coaching := map[string]interface{}{
		"player_id":              playerID,
		"champion":               champion,
		"time_range":             timeRange,
		"current_performance": map[string]interface{}{
			"overall_rating":   analysis.OverallRating,
			"mechanics_score":  analysis.MechanicsScore,
			"game_knowledge_score": analysis.GameKnowledgeScore,
			"consistency_score": analysis.ConsistencyScore,
		},
		"core_strengths":         analysis.CoreStrengths,
		"improvement_areas":      analysis.ImprovementAreas,
		"playstyle_recommendations": analysis.PlayStyleRecommendations,
		"training_recommendations": analysis.TrainingRecommendations,
		"learning_curve":         analysis.LearningCurve,
	}

	c.JSON(http.StatusOK, coaching)
}

// GetChampionTrends godoc
// @Summary Get champion performance trends
// @Description Returns champion performance trends over time
// @Tags champion
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param champion query string true "Champion name"
// @Param metric query string false "Metric type (rating, winrate, kda, dpm)"
// @Param period query string false "Period (daily, weekly) - default: daily"
// @Param days query int false "Number of days (default: 30)"
// @Success 200 {object} []services.ChampionTrendPoint
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/champion/{player_id}/trends [get]
func (ch *ChampionHandler) GetChampionTrends(c *gin.Context) {
	playerID := c.Param("player_id")
	champion := c.Query("champion")
	metric := c.DefaultQuery("metric", "rating")
	
	if playerID == "" || champion == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID and champion are required",
		})
		return
	}

	// Validate metric
	validMetrics := map[string]bool{
		"rating":   true,
		"winrate":  true,
		"kda":      true,
		"dpm":      true,
		"cs":       true,
	}
	if !validMetrics[metric] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid metric. Use: rating, winrate, kda, dpm, or cs",
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

	// Determine time range based on days
	timeRange := "30d"
	if days <= 7 {
		timeRange = "7d"
	} else if days <= 90 {
		timeRange = "90d"
	}

	// Get champion analysis to extract trend data
	analysis, err := ch.championService.AnalyzeChampion(c.Request.Context(), playerID, champion, timeRange, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to get champion trend data",
		})
		return
	}

	c.JSON(http.StatusOK, analysis.TrendData)
}

// Register routes for champion analytics
func (ch *ChampionHandler) RegisterRoutes(router *gin.RouterGroup) {
	champion := router.Group("/champion")
	{
		champion.GET("/:player_id/analysis", ch.GetChampionAnalysis)
		champion.GET("/:player_id/mastery", ch.GetChampionMastery)
		champion.GET("/:player_id/comparison", ch.GetChampionComparison)
		champion.GET("/:player_id/power-spikes", ch.GetChampionPowerSpikes)
		champion.GET("/:player_id/matchups", ch.GetChampionMatchups)
		champion.GET("/:player_id/builds", ch.GetChampionBuilds)
		champion.GET("/:player_id/coaching", ch.GetChampionCoaching)
		champion.GET("/:player_id/trends", ch.GetChampionTrends)
	}
}