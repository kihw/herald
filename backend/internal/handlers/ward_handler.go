package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herald/internal/services"
)

// WardHandler handles ward placement and map control analytics requests
type WardHandler struct {
	wardService *services.WardAnalyticsService
}

// NewWardHandler creates a new ward handler
func NewWardHandler(wardService *services.WardAnalyticsService) *WardHandler {
	return &WardHandler{
		wardService: wardService,
	}
}

// WardAnalysisRequest represents request for ward analysis
type WardAnalysisRequest struct {
	TimeRange string `form:"time_range" json:"time_range" binding:"required"`
	Champion  string `form:"champion" json:"champion"`
	Position  string `form:"position" json:"position"`
	GameMode  string `form:"game_mode" json:"game_mode"`
}

// MapControlRequest represents request for map control analysis
type MapControlRequest struct {
	TimeRange string `form:"time_range" json:"time_range" binding:"required"`
	Zone      string `form:"zone" json:"zone"`      // Specific zone filter
	Metric    string `form:"metric" json:"metric"`  // "control", "coverage", "efficiency"
}

// GetWardAnalysis godoc
// @Summary Get comprehensive ward placement and map control analysis
// @Description Analyzes player's ward placement patterns, map control, and vision impact
// @Tags ward
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string true "Time range (7d, 30d, 90d)"
// @Param champion query string false "Champion name filter"
// @Param position query string false "Position filter"
// @Param game_mode query string false "Game mode filter"
// @Success 200 {object} services.WardAnalysis
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/ward/{player_id}/analysis [get]
func (wh *WardHandler) GetWardAnalysis(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	var req WardAnalysisRequest
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
			"RANKED_SOLO_5x5":   true,
			"RANKED_FLEX_SR":    true,
			"NORMAL_DRAFT":      true,
			"ARAM":              true,
		}
		if !validGameModes[req.GameMode] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid game mode",
			})
			return
		}
	}

	// Perform ward analysis
	analysis, err := wh.wardService.AnalyzeWards(c.Request.Context(), playerID, req.TimeRange, req.Champion, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze ward data",
		})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetMapControl godoc
// @Summary Get map control analysis
// @Description Analyzes player's map control score and territory coverage
// @Tags ward
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string true "Time range (7d, 30d, 90d)"
// @Param zone query string false "Zone filter (river, jungle, dragon, baron)"
// @Param metric query string false "Metric type (control, coverage, efficiency)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/ward/{player_id}/map-control [get]
func (wh *WardHandler) GetMapControl(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	var req MapControlRequest
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

	// Validate zone if provided
	if req.Zone != "" {
		validZones := map[string]bool{
			"river":      true,
			"jungle":     true,
			"dragon":     true,
			"baron":      true,
			"objectives": true,
		}
		if !validZones[req.Zone] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid zone. Use: river, jungle, dragon, baron, or objectives",
			})
			return
		}
	}

	// Validate metric if provided
	if req.Metric != "" {
		validMetrics := map[string]bool{
			"control":    true,
			"coverage":   true,
			"efficiency": true,
		}
		if !validMetrics[req.Metric] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid metric. Use: control, coverage, or efficiency",
			})
			return
		}
	}

	// Get ward analysis
	analysis, err := wh.wardService.AnalyzeWards(c.Request.Context(), playerID, req.TimeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze map control",
		})
		return
	}

	// Create map control response
	mapControl := map[string]interface{}{
		"player_id":           playerID,
		"time_range":          req.TimeRange,
		"map_control_score":   analysis.MapControlScore,
		"territory_controlled": analysis.TerritoryControlled,
		"strategic_coverage":  analysis.StrategicCoverage,
		"zone_control":        analysis.ZoneControl,
		"river_control":       analysis.RiverControl,
		"jungle_control":      analysis.JungleControl,
	}

	// Filter by specific zone if requested
	if req.Zone != "" {
		switch req.Zone {
		case "river":
			mapControl["focus_data"] = analysis.RiverControl
		case "jungle":
			mapControl["focus_data"] = analysis.JungleControl
		case "dragon":
			if dragonControl, exists := analysis.ZoneControl["Dragon Pit"]; exists {
				mapControl["focus_data"] = dragonControl
			}
		case "baron":
			if baronControl, exists := analysis.ZoneControl["Baron Pit"]; exists {
				mapControl["focus_data"] = baronControl
			}
		}
	}

	c.JSON(http.StatusOK, mapControl)
}

// GetWardPlacementPatterns godoc
// @Summary Get ward placement pattern analysis
// @Description Analyzes player's ward placement patterns and optimality
// @Tags ward
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string false "Time range (default: 30d)"
// @Param pattern_type query string false "Pattern type (placement, timing, optimization)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/ward/{player_id}/patterns [get]
func (wh *WardHandler) GetWardPlacementPatterns(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	timeRange := c.DefaultQuery("time_range", "30d")
	patternType := c.Query("pattern_type")

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Validate pattern type if provided
	if patternType != "" {
		validPatternTypes := map[string]bool{
			"placement":    true,
			"timing":       true,
			"optimization": true,
		}
		if !validPatternTypes[patternType] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid pattern type. Use: placement, timing, or optimization",
			})
			return
		}
	}

	// Get ward analysis
	analysis, err := wh.wardService.AnalyzeWards(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze ward patterns",
		})
		return
	}

	// Create patterns response
	patterns := map[string]interface{}{
		"player_id":            playerID,
		"time_range":           timeRange,
		"placement_patterns":   analysis.PlacementPatterns,
		"optimal_placements":   analysis.OptimalPlacements,
		"placement_timing":     analysis.PlacementTiming,
		"placement_optimization": analysis.PlacementOptimization,
	}

	// Filter by specific pattern type if requested
	if patternType != "" {
		switch patternType {
		case "placement":
			patterns["focus_data"] = analysis.PlacementPatterns
		case "timing":
			patterns["focus_data"] = analysis.PlacementTiming
		case "optimization":
			patterns["focus_data"] = analysis.PlacementOptimization
		}
	}

	c.JSON(http.StatusOK, patterns)
}

// GetWardTypeAnalysis godoc
// @Summary Get ward type specific analysis
// @Description Analyzes performance by ward type (yellow, control, blue trinket)
// @Tags ward
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param ward_type query string false "Ward type filter (YELLOW, CONTROL, BLUE_TRINKET)"
// @Param time_range query string false "Time range (default: 30d)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/ward/{player_id}/types [get]
func (wh *WardHandler) GetWardTypeAnalysis(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	wardType := c.Query("ward_type")
	timeRange := c.DefaultQuery("time_range", "30d")

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Validate ward type if provided
	if wardType != "" {
		validWardTypes := map[string]bool{
			"YELLOW":       true,
			"CONTROL":      true,
			"BLUE_TRINKET": true,
		}
		if !validWardTypes[wardType] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid ward type. Use: YELLOW, CONTROL, or BLUE_TRINKET",
			})
			return
		}
	}

	// Get ward analysis
	analysis, err := wh.wardService.AnalyzeWards(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze ward types",
		})
		return
	}

	// Create ward type response
	wardTypeAnalysis := map[string]interface{}{
		"player_id":             playerID,
		"time_range":            timeRange,
		"yellow_wards_analysis": analysis.YellowWardsAnalysis,
		"control_wards_analysis": analysis.ControlWardsAnalysis,
		"blue_ward_analysis":    analysis.BlueWardAnalysis,
	}

	// Filter by specific ward type if requested
	if wardType != "" {
		switch wardType {
		case "YELLOW":
			wardTypeAnalysis["focus_analysis"] = analysis.YellowWardsAnalysis
		case "CONTROL":
			wardTypeAnalysis["focus_analysis"] = analysis.ControlWardsAnalysis
		case "BLUE_TRINKET":
			wardTypeAnalysis["focus_analysis"] = analysis.BlueWardAnalysis
		}
	}

	c.JSON(http.StatusOK, wardTypeAnalysis)
}

// GetWardTrends godoc
// @Summary Get ward performance trends
// @Description Returns ward placement and map control trends over time
// @Tags ward
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param metric query string true "Metric type (wards_placed, wards_killed, map_control, efficiency)"
// @Param period query string false "Period (daily, weekly) - default: daily"
// @Param days query int false "Number of days (default: 30)"
// @Success 200 {object} []services.WardTrendPoint
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/ward/{player_id}/trends [get]
func (wh *WardHandler) GetWardTrends(c *gin.Context) {
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
		"wards_placed":    true,
		"wards_killed":    true,
		"map_control":     true,
		"efficiency":      true,
	}
	if !validMetrics[metric] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid metric. Use: wards_placed, wards_killed, map_control, or efficiency",
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

	// Get ward analysis to extract trend data
	timeRange := "30d"
	if days <= 7 {
		timeRange = "7d"
	} else if days <= 90 {
		timeRange = "90d"
	}

	analysis, err := wh.wardService.AnalyzeWards(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to get ward trend data",
		})
		return
	}

	c.JSON(http.StatusOK, analysis.TrendData)
}

// GetWardClearing godoc
// @Summary Get ward clearing analysis
// @Description Analyzes player's ward clearing patterns and counter-warding
// @Tags ward
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string false "Time range (default: 30d)"
// @Param clearing_type query string false "Clearing type (proactive, reactive, opportunistic)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/ward/{player_id}/clearing [get]
func (wh *WardHandler) GetWardClearing(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	timeRange := c.DefaultQuery("time_range", "30d")
	clearingType := c.Query("clearing_type")

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Validate clearing type if provided
	if clearingType != "" {
		validClearingTypes := map[string]bool{
			"proactive":     true,
			"reactive":      true,
			"opportunistic": true,
		}
		if !validClearingTypes[clearingType] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid clearing type. Use: proactive, reactive, or opportunistic",
			})
			return
		}
	}

	// Get ward analysis
	analysis, err := wh.wardService.AnalyzeWards(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze ward clearing",
		})
		return
	}

	// Create ward clearing response
	wardClearing := map[string]interface{}{
		"player_id":               playerID,
		"time_range":              timeRange,
		"ward_clearing_patterns":  analysis.WardClearingPatterns,
		"counter_warding_score":   analysis.CounterWardingScore,
		"clearing_optimization":   analysis.ClearingOptimization,
		"vision_denied_score":     analysis.VisionDeniedScore,
	}

	c.JSON(http.StatusOK, wardClearing)
}

// GetWardOptimization godoc
// @Summary Get ward placement and clearing optimization suggestions
// @Description Provides personalized recommendations for improving ward efficiency
// @Tags ward
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string false "Time range (default: 30d)"
// @Param optimization_type query string false "Optimization type (placement, clearing, timing)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/ward/{player_id}/optimization [get]
func (wh *WardHandler) GetWardOptimization(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	timeRange := c.DefaultQuery("time_range", "30d")
	optimizationType := c.Query("optimization_type")

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Validate optimization type if provided
	if optimizationType != "" {
		validOptimizationTypes := map[string]bool{
			"placement": true,
			"clearing":  true,
			"timing":    true,
		}
		if !validOptimizationTypes[optimizationType] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid optimization type. Use: placement, clearing, or timing",
			})
			return
		}
	}

	// Get ward analysis
	analysis, err := wh.wardService.AnalyzeWards(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze ward optimization",
		})
		return
	}

	// Create optimization response
	optimization := map[string]interface{}{
		"player_id":               playerID,
		"time_range":              timeRange,
		"current_performance": map[string]interface{}{
			"map_control_score":     analysis.MapControlScore,
			"ward_efficiency":       analysis.WardEfficiency,
			"counter_warding_score": analysis.CounterWardingScore,
		},
		"placement_optimization": analysis.PlacementOptimization,
		"clearing_optimization":  analysis.ClearingOptimization,
		"recommendations":        analysis.Recommendations,
		"improvement_potential": map[string]interface{}{
			"expected_control_gain": analysis.PlacementOptimization.ExpectedControlGain,
			"expected_safety_gain":  analysis.PlacementOptimization.ExpectedSafetyGain,
			"expected_denial_gain":  analysis.ClearingOptimization.ExpectedDenialGain,
		},
	}

	c.JSON(http.StatusOK, optimization)
}

// GetObjectiveControl godoc
// @Summary Get objective-specific vision control analysis
// @Description Analyzes vision setup and control around objectives (Dragon, Baron, Herald)
// @Tags ward
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param objective query string false "Objective filter (dragon, baron, herald, elder)"
// @Param time_range query string false "Time range (default: 30d)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/ward/{player_id}/objectives [get]
func (wh *WardHandler) GetObjectiveControl(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	objective := c.Query("objective")
	timeRange := c.DefaultQuery("time_range", "30d")

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Validate objective if provided
	if objective != "" {
		validObjectives := map[string]bool{
			"dragon": true,
			"baron":  true,
			"herald": true,
			"elder":  true,
		}
		if !validObjectives[objective] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid objective. Use: dragon, baron, herald, or elder",
			})
			return
		}
	}

	// Get ward analysis
	analysis, err := wh.wardService.AnalyzeWards(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze objective control",
		})
		return
	}

	// Create objective control response
	objectiveControl := map[string]interface{}{
		"player_id":        playerID,
		"time_range":       timeRange,
		"objective_setup":  analysis.ObjectiveSetup,
		"safety_provided":  analysis.SafetyProvided,
		"strategic_coverage": analysis.StrategicCoverage,
	}

	// Add objective-specific data if requested
	if objective != "" {
		switch objective {
		case "dragon":
			objectiveControl["focus_score"] = analysis.ObjectiveSetup.DragonSetupScore
			objectiveControl["focus_coverage"] = analysis.StrategicCoverage.DragonPitCoverage
		case "baron":
			objectiveControl["focus_score"] = analysis.ObjectiveSetup.BaronSetupScore
			objectiveControl["focus_coverage"] = analysis.StrategicCoverage.BaronPitCoverage
		case "herald":
			objectiveControl["focus_score"] = analysis.ObjectiveSetup.HeraldSetupScore
		case "elder":
			objectiveControl["focus_score"] = analysis.ObjectiveSetup.ElderSetupScore
		}
	}

	c.JSON(http.StatusOK, objectiveControl)
}

// Register routes for ward analytics
func (wh *WardHandler) RegisterRoutes(router *gin.RouterGroup) {
	ward := router.Group("/ward")
	{
		ward.GET("/:player_id/analysis", wh.GetWardAnalysis)
		ward.GET("/:player_id/map-control", wh.GetMapControl)
		ward.GET("/:player_id/patterns", wh.GetWardPlacementPatterns)
		ward.GET("/:player_id/types", wh.GetWardTypeAnalysis)
		ward.GET("/:player_id/trends", wh.GetWardTrends)
		ward.GET("/:player_id/clearing", wh.GetWardClearing)
		ward.GET("/:player_id/optimization", wh.GetWardOptimization)
		ward.GET("/:player_id/objectives", wh.GetObjectiveControl)
	}
}