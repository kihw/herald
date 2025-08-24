package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herald-lol/herald/backend/internal/services"
)

// GoldHandler handles gold efficiency analytics requests
type GoldHandler struct {
	goldService *services.GoldAnalyticsService
}

// NewGoldHandler creates a new gold handler
func NewGoldHandler(goldService *services.GoldAnalyticsService) *GoldHandler {
	return &GoldHandler{
		goldService: goldService,
	}
}

// GoldAnalysisRequest represents request for gold analysis
type GoldAnalysisRequest struct {
	TimeRange string `form:"time_range" json:"time_range" binding:"required"`
	Champion  string `form:"champion" json:"champion"`
	Position  string `form:"position" json:"position"`
	GameMode  string `form:"game_mode" json:"game_mode"`
}

// GoldOptimizationRequest represents request for gold optimization analysis
type GoldOptimizationRequest struct {
	TimeRange  string `form:"time_range" json:"time_range" binding:"required"`
	FocusArea  string `form:"focus_area" json:"focus_area"` // "income", "spending", "efficiency"
	Difficulty string `form:"difficulty" json:"difficulty"` // "easy", "medium", "hard"
}

// GetGoldAnalysis godoc
// @Summary Get comprehensive gold efficiency analysis
// @Description Analyzes player's gold generation, spending efficiency, and economic performance
// @Tags gold
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string true "Time range (7d, 30d, 90d)"
// @Param champion query string false "Champion name filter"
// @Param position query string false "Position filter"
// @Param game_mode query string false "Game mode filter"
// @Success 200 {object} services.GoldAnalysis
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/gold/{player_id}/analysis [get]
func (gh *GoldHandler) GetGoldAnalysis(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	var req GoldAnalysisRequest
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

	// Perform gold analysis
	analysis, err := gh.goldService.AnalyzeGold(c.Request.Context(), playerID, req.TimeRange, req.Champion, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze gold data",
		})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetGoldSources godoc
// @Summary Get gold income source breakdown
// @Description Analyzes where player's gold comes from (farming, kills, objectives, etc.)
// @Tags gold
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string true "Time range (7d, 30d, 90d)"
// @Param source_type query string false "Source type filter (farming, kills, objectives, passive)"
// @Success 200 {object} services.GoldSourcesData
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/gold/{player_id}/sources [get]
func (gh *GoldHandler) GetGoldSources(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	timeRange := c.DefaultQuery("time_range", "30d")
	sourceType := c.Query("source_type")

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Validate source type if provided
	if sourceType != "" {
		validSourceTypes := map[string]bool{
			"farming":    true,
			"kills":      true,
			"objectives": true,
			"passive":    true,
			"items":      true,
		}
		if !validSourceTypes[sourceType] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid source type. Use: farming, kills, objectives, passive, or items",
			})
			return
		}
	}

	// Get gold analysis
	analysis, err := gh.goldService.AnalyzeGold(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze gold sources",
		})
		return
	}

	c.JSON(http.StatusOK, analysis.GoldSources)
}

// GetItemEfficiency godoc
// @Summary Get item purchase and utilization efficiency
// @Description Analyzes how efficiently player purchases and uses items
// @Tags gold
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string false "Time range (default: 30d)"
// @Param item_category query string false "Item category filter (damage, defensive, utility)"
// @Success 200 {object} services.ItemEfficiencyData
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/gold/{player_id}/item-efficiency [get]
func (gh *GoldHandler) GetItemEfficiency(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	timeRange := c.DefaultQuery("time_range", "30d")
	itemCategory := c.Query("item_category")

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Validate item category if provided
	if itemCategory != "" {
		validCategories := map[string]bool{
			"damage":    true,
			"defensive": true,
			"utility":   true,
		}
		if !validCategories[itemCategory] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid item category. Use: damage, defensive, or utility",
			})
			return
		}
	}

	// Get gold analysis
	analysis, err := gh.goldService.AnalyzeGold(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze item efficiency",
		})
		return
	}

	c.JSON(http.StatusOK, analysis.ItemEfficiency)
}

// GetSpendingPatterns godoc
// @Summary Get gold spending behavior analysis
// @Description Analyzes player's gold spending patterns and back timing
// @Tags gold
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string false "Time range (default: 30d)"
// @Param spending_type query string false "Spending type filter (items, wards, consumables)"
// @Success 200 {object} services.SpendingPatternsData
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/gold/{player_id}/spending [get]
func (gh *GoldHandler) GetSpendingPatterns(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	timeRange := c.DefaultQuery("time_range", "30d")
	spendingType := c.Query("spending_type")

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Validate spending type if provided
	if spendingType != "" {
		validSpendingTypes := map[string]bool{
			"items":       true,
			"wards":       true,
			"consumables": true,
		}
		if !validSpendingTypes[spendingType] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid spending type. Use: items, wards, or consumables",
			})
			return
		}
	}

	// Get gold analysis
	analysis, err := gh.goldService.AnalyzeGold(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze spending patterns",
		})
		return
	}

	c.JSON(http.StatusOK, analysis.SpendingPatterns)
}

// GetGoldTrends godoc
// @Summary Get gold performance trends
// @Description Returns gold efficiency and GPM trends over time
// @Tags gold
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param metric query string true "Metric type (gold_per_minute, efficiency, farming_efficiency, spending_efficiency)"
// @Param period query string false "Period (daily, weekly) - default: daily"
// @Param days query int false "Number of days (default: 30)"
// @Success 200 {object} []services.GoldTrendPoint
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/gold/{player_id}/trends [get]
func (gh *GoldHandler) GetGoldTrends(c *gin.Context) {
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
		"gold_per_minute":     true,
		"efficiency":          true,
		"farming_efficiency":  true,
		"spending_efficiency": true,
	}
	if !validMetrics[metric] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid metric. Use: gold_per_minute, efficiency, farming_efficiency, or spending_efficiency",
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

	// Get gold analysis to extract trend data
	timeRange := "30d"
	if days <= 7 {
		timeRange = "7d"
	} else if days <= 90 {
		timeRange = "90d"
	}

	analysis, err := gh.goldService.AnalyzeGold(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to get gold trend data",
		})
		return
	}

	c.JSON(http.StatusOK, analysis.TrendData)
}

// GetGoldComparison godoc
// @Summary Compare gold performance with benchmarks
// @Description Compares player's gold metrics against role/rank/global benchmarks
// @Tags gold
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param benchmark_type query string true "Benchmark type (role, rank, global)"
// @Param time_range query string false "Time range (default: 30d)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/gold/{player_id}/comparison [get]
func (gh *GoldHandler) GetGoldComparison(c *gin.Context) {
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

	// Get gold analysis
	analysis, err := gh.goldService.AnalyzeGold(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze gold data",
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
				"gold_per_minute":   analysis.AverageGoldPerMinute,
				"gold_efficiency":   analysis.GoldEfficiencyScore,
				"economy_rating":    analysis.EconomyRating,
				"gold_impact_score": analysis.GoldImpactScore,
			},
			"benchmark_metrics": analysis.RoleBenchmark,
			"percentile":        analysis.RoleBenchmark.PlayerPercentile,
		}
	case "rank":
		comparison = map[string]interface{}{
			"benchmark_type": "rank",
			"player_metrics": map[string]interface{}{
				"gold_per_minute":   analysis.AverageGoldPerMinute,
				"gold_efficiency":   analysis.GoldEfficiencyScore,
				"economy_rating":    analysis.EconomyRating,
				"gold_impact_score": analysis.GoldImpactScore,
			},
			"benchmark_metrics": analysis.RankBenchmark,
			"percentile":        analysis.RankBenchmark.PlayerPercentile,
		}
	case "global":
		comparison = map[string]interface{}{
			"benchmark_type": "global",
			"player_metrics": map[string]interface{}{
				"gold_per_minute":   analysis.AverageGoldPerMinute,
				"gold_efficiency":   analysis.GoldEfficiencyScore,
				"economy_rating":    analysis.EconomyRating,
				"gold_impact_score": analysis.GoldImpactScore,
			},
			"benchmark_metrics": analysis.GlobalBenchmark,
			"percentile":        analysis.GlobalBenchmark.PlayerPercentile,
		}
	}

	c.JSON(http.StatusOK, comparison)
}

// GetGoldOptimization godoc
// @Summary Get gold efficiency optimization suggestions
// @Description Provides personalized recommendations for improving gold generation and spending
// @Tags gold
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string false "Time range (default: 30d)"
// @Param focus_area query string false "Focus area (income, spending, efficiency)"
// @Param difficulty query string false "Implementation difficulty (easy, medium, hard)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/gold/{player_id}/optimization [get]
func (gh *GoldHandler) GetGoldOptimization(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	timeRange := c.DefaultQuery("time_range", "30d")
	focusArea := c.Query("focus_area")
	difficulty := c.Query("difficulty")

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
			"income":     true,
			"spending":   true,
			"efficiency": true,
		}
		if !validFocusAreas[focusArea] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid focus area. Use: income, spending, or efficiency",
			})
			return
		}
	}

	// Validate difficulty if provided
	if difficulty != "" {
		validDifficulties := map[string]bool{
			"easy":   true,
			"medium": true,
			"hard":   true,
		}
		if !validDifficulties[difficulty] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid difficulty. Use: easy, medium, or hard",
			})
			return
		}
	}

	// Get gold analysis
	analysis, err := gh.goldService.AnalyzeGold(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze gold optimization",
		})
		return
	}

	// Filter recommendations based on focus area and difficulty
	recommendations := analysis.Recommendations
	if difficulty != "" {
		filtered := []services.GoldRecommendation{}
		for _, rec := range recommendations {
			if rec.ImplementationDifficulty == difficulty {
				filtered = append(filtered, rec)
			}
		}
		recommendations = filtered
	}

	// Create optimization response
	optimization := map[string]interface{}{
		"player_id":  playerID,
		"time_range": timeRange,
		"current_performance": map[string]interface{}{
			"gold_per_minute": analysis.AverageGoldPerMinute,
			"gold_efficiency": analysis.GoldEfficiencyScore,
			"economy_rating":  analysis.EconomyRating,
		},
		"income_optimization":   analysis.IncomeOptimization,
		"spending_optimization": analysis.SpendingOptimization,
		"recommendations":       recommendations,
		"improvement_potential": map[string]interface{}{
			"expected_gpm_increase": analysis.IncomeOptimization.ExpectedGPMIncrease,
			"expected_wr_increase":  analysis.IncomeOptimization.ExpectedWinRateIncrease,
		},
	}

	c.JSON(http.StatusOK, optimization)
}

// GetGoldPhaseAnalysis godoc
// @Summary Get gold efficiency by game phase
// @Description Analyzes gold performance across early, mid, and late game phases
// @Tags gold
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string false "Time range (default: 30d)"
// @Param phase query string false "Phase filter (early, mid, late)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/gold/{player_id}/phases [get]
func (gh *GoldHandler) GetGoldPhaseAnalysis(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	timeRange := c.DefaultQuery("time_range", "30d")
	phase := c.Query("phase")

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Validate phase if provided
	if phase != "" {
		validPhases := map[string]bool{
			"early": true,
			"mid":   true,
			"late":  true,
		}
		if !validPhases[phase] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid phase. Use: early, mid, or late",
			})
			return
		}
	}

	// Get gold analysis
	analysis, err := gh.goldService.AnalyzeGold(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze gold phase data",
		})
		return
	}

	// Create phase analysis response
	phaseAnalysis := map[string]interface{}{
		"player_id":  playerID,
		"time_range": timeRange,
		"early_game": analysis.EarlyGameGold,
		"mid_game":   analysis.MidGameGold,
		"late_game":  analysis.LateGameGold,
	}

	// Filter by specific phase if requested
	if phase != "" {
		switch phase {
		case "early":
			phaseAnalysis["phase_data"] = analysis.EarlyGameGold
		case "mid":
			phaseAnalysis["phase_data"] = analysis.MidGameGold
		case "late":
			phaseAnalysis["phase_data"] = analysis.LateGameGold
		}
	}

	c.JSON(http.StatusOK, phaseAnalysis)
}

// Register routes for gold analytics
func (gh *GoldHandler) RegisterRoutes(router *gin.RouterGroup) {
	gold := router.Group("/gold")
	{
		gold.GET("/:player_id/analysis", gh.GetGoldAnalysis)
		gold.GET("/:player_id/sources", gh.GetGoldSources)
		gold.GET("/:player_id/item-efficiency", gh.GetItemEfficiency)
		gold.GET("/:player_id/spending", gh.GetSpendingPatterns)
		gold.GET("/:player_id/trends", gh.GetGoldTrends)
		gold.GET("/:player_id/comparison", gh.GetGoldComparison)
		gold.GET("/:player_id/optimization", gh.GetGoldOptimization)
		gold.GET("/:player_id/phases", gh.GetGoldPhaseAnalysis)
	}
}
