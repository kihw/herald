package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herald-lol/herald/backend/internal/services"
)

// VisionHandler handles vision analytics requests
type VisionHandler struct {
	visionService *services.VisionAnalyticsService
}

// NewVisionHandler creates a new vision handler
func NewVisionHandler(visionService *services.VisionAnalyticsService) *VisionHandler {
	return &VisionHandler{
		visionService: visionService,
	}
}

// VisionAnalysisRequest represents request for vision analysis
type VisionAnalysisRequest struct {
	TimeRange string `form:"time_range" json:"time_range" binding:"required"`
	Champion  string `form:"champion" json:"champion"`
	Position  string `form:"position" json:"position"`
}

// HeatmapRequest represents request for heatmap generation
type HeatmapRequest struct {
	TimeRange string `form:"time_range" json:"time_range" binding:"required"`
	WardType  string `form:"ward_type" json:"ward_type" binding:"required"`
	MapSide   string `form:"map_side" json:"map_side"`
}

// GetVisionAnalysis godoc
// @Summary Get comprehensive vision analysis
// @Description Analyzes player's vision control performance with heatmaps and recommendations
// @Tags vision
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string true "Time range (7d, 30d, 90d)"
// @Param champion query string false "Champion name filter"
// @Param position query string false "Position filter"
// @Success 200 {object} services.VisionAnalysis
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/vision/{player_id}/analysis [get]
func (vh *VisionHandler) GetVisionAnalysis(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	var req VisionAnalysisRequest
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

	// Perform vision analysis
	analysis, err := vh.visionService.AnalyzeVision(c.Request.Context(), playerID, req.TimeRange, req.Champion, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze vision data",
		})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetVisionHeatmap godoc
// @Summary Generate vision heatmap
// @Description Creates heatmap visualization for ward placements and vision control
// @Tags vision
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string true "Time range (7d, 30d, 90d)"
// @Param ward_type query string true "Ward type (YELLOW, CONTROL, BLUE_TRINKET, ALL)"
// @Param map_side query string false "Map side filter (BLUE, RED, BOTH)"
// @Success 200 {object} services.HeatmapData
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/vision/{player_id}/heatmap [get]
func (vh *VisionHandler) GetVisionHeatmap(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	var req HeatmapRequest
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

	// Validate ward type
	validWardTypes := map[string]bool{
		"YELLOW":       true,
		"CONTROL":      true,
		"BLUE_TRINKET": true,
		"ALL":          true,
	}
	if !validWardTypes[req.WardType] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid ward type. Use: YELLOW, CONTROL, BLUE_TRINKET, or ALL",
		})
		return
	}

	// Validate map side if provided
	if req.MapSide != "" {
		validMapSides := map[string]bool{
			"BLUE": true,
			"RED":  true,
			"BOTH": true,
		}
		if !validMapSides[req.MapSide] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid map side. Use: BLUE, RED, or BOTH",
			})
			return
		}
	}

	// Generate heatmap
	heatmap, err := vh.visionService.GenerateVisionHeatmap(c.Request.Context(), playerID, req.TimeRange, req.WardType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "heatmap_error",
			Message: "Failed to generate vision heatmap",
		})
		return
	}

	c.JSON(http.StatusOK, heatmap)
}

// GetVisionRecommendations godoc
// @Summary Get vision improvement recommendations
// @Description Provides personalized recommendations for improving vision control
// @Tags vision
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param time_range query string false "Time range (default: 30d)"
// @Param priority query string false "Priority filter (high, medium, low)"
// @Param category query string false "Category filter (warding, dewarding, positioning, timing)"
// @Success 200 {object} []services.VisionRecommendation
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/vision/{player_id}/recommendations [get]
func (vh *VisionHandler) GetVisionRecommendations(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	timeRange := c.DefaultQuery("time_range", "30d")
	priority := c.Query("priority")
	category := c.Query("category")

	// Validate filters
	if timeRange != "" && !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	if priority != "" {
		validPriorities := map[string]bool{"high": true, "medium": true, "low": true}
		if !validPriorities[priority] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid priority. Use: high, medium, or low",
			})
			return
		}
	}

	if category != "" {
		validCategories := map[string]bool{
			"warding":     true,
			"dewarding":   true,
			"positioning": true,
			"timing":      true,
		}
		if !validCategories[category] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid category. Use: warding, dewarding, positioning, or timing",
			})
			return
		}
	}

	// Get vision analysis first to generate recommendations
	analysis, err := vh.visionService.AnalyzeVision(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze vision data for recommendations",
		})
		return
	}

	// Filter recommendations based on query parameters
	recommendations := analysis.Recommendations

	if priority != "" {
		filtered := []services.VisionRecommendation{}
		for _, rec := range recommendations {
			if rec.Priority == priority {
				filtered = append(filtered, rec)
			}
		}
		recommendations = filtered
	}

	if category != "" {
		filtered := []services.VisionRecommendation{}
		for _, rec := range recommendations {
			if rec.Category == category {
				filtered = append(filtered, rec)
			}
		}
		recommendations = filtered
	}

	c.JSON(http.StatusOK, recommendations)
}

// GetVisionTrends godoc
// @Summary Get vision performance trends
// @Description Returns vision score and ward efficiency trends over time
// @Tags vision
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param metric query string true "Metric type (vision_score, wards_placed, wards_killed, efficiency)"
// @Param period query string false "Period (daily, weekly) - default: daily"
// @Param days query int false "Number of days (default: 30)"
// @Success 200 {object} []services.VisionTrendPoint
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/vision/{player_id}/trends [get]
func (vh *VisionHandler) GetVisionTrends(c *gin.Context) {
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
		"vision_score": true,
		"wards_placed": true,
		"wards_killed": true,
		"efficiency":   true,
	}
	if !validMetrics[metric] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid metric. Use: vision_score, wards_placed, wards_killed, or efficiency",
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

	// Get vision analysis to extract trend data
	timeRange := "30d"
	if days <= 7 {
		timeRange = "7d"
	} else if days <= 90 {
		timeRange = "90d"
	}

	analysis, err := vh.visionService.AnalyzeVision(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to get vision trend data",
		})
		return
	}

	c.JSON(http.StatusOK, analysis.TrendData)
}

// GetVisionZoneAnalysis godoc
// @Summary Get zone-specific vision analysis
// @Description Analyzes vision performance in specific map zones
// @Tags vision
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param zone query string false "Zone filter (jungle, river, dragon, baron, etc.)"
// @Param time_range query string false "Time range (default: 30d)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/vision/{player_id}/zones [get]
func (vh *VisionHandler) GetVisionZoneAnalysis(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Player ID is required",
		})
		return
	}

	zone := c.Query("zone")
	timeRange := c.DefaultQuery("time_range", "30d")

	// Validate time range
	if !isValidTimeRange(timeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 30d, or 90d",
		})
		return
	}

	// Get vision analysis
	analysis, err := vh.visionService.AnalyzeVision(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze zone vision data",
		})
		return
	}

	// Extract zone-specific data from heatmaps
	zoneData := map[string]interface{}{
		"player_id":  playerID,
		"time_range": timeRange,
		"zones":      analysis.WardHeatmaps,
	}

	// Filter by specific zone if requested
	if zone != "" {
		// Filter heatmap data for specific zone
		filteredData := map[string]interface{}{
			"zone":            zone,
			"intensity":       0,
			"coverage":        0.0,
			"strategic_value": 0.0,
		}

		// Extract zone-specific intensity from heatmaps
		if intensity, exists := analysis.WardHeatmaps.YellowWards.Intensity[zone]; exists {
			filteredData["intensity"] = intensity
		}
		if intensity, exists := analysis.WardHeatmaps.ControlWards.Intensity[zone]; exists {
			filteredData["control_ward_intensity"] = intensity
		}
		if intensity, exists := analysis.WardHeatmaps.WardKills.Intensity[zone]; exists {
			filteredData["ward_kill_intensity"] = intensity
		}

		zoneData["zone_analysis"] = filteredData
	}

	c.JSON(http.StatusOK, zoneData)
}

// GetVisionComparison godoc
// @Summary Compare vision performance with benchmarks
// @Description Compares player's vision metrics against role/rank/global benchmarks
// @Tags vision
// @Accept json
// @Produce json
// @Param player_id path string true "Player ID"
// @Param benchmark_type query string true "Benchmark type (role, rank, global)"
// @Param time_range query string false "Time range (default: 30d)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/vision/{player_id}/comparison [get]
func (vh *VisionHandler) GetVisionComparison(c *gin.Context) {
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

	// Get vision analysis
	analysis, err := vh.visionService.AnalyzeVision(c.Request.Context(), playerID, timeRange, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze vision data",
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
				"vision_score":    analysis.AverageVisionScore,
				"wards_placed":    analysis.AverageWardsPlaced,
				"wards_killed":    analysis.AverageWardsKilled,
				"ward_efficiency": analysis.WardEfficiency,
			},
			"benchmark_metrics": analysis.RoleBenchmark,
			"percentile":        analysis.RoleBenchmark.PlayerPercentile,
		}
	case "rank":
		comparison = map[string]interface{}{
			"benchmark_type": "rank",
			"player_metrics": map[string]interface{}{
				"vision_score":    analysis.AverageVisionScore,
				"wards_placed":    analysis.AverageWardsPlaced,
				"wards_killed":    analysis.AverageWardsKilled,
				"ward_efficiency": analysis.WardEfficiency,
			},
			"benchmark_metrics": analysis.RankBenchmark,
			"percentile":        analysis.RankBenchmark.PlayerPercentile,
		}
	case "global":
		comparison = map[string]interface{}{
			"benchmark_type": "global",
			"player_metrics": map[string]interface{}{
				"vision_score":    analysis.AverageVisionScore,
				"wards_placed":    analysis.AverageWardsPlaced,
				"wards_killed":    analysis.AverageWardsKilled,
				"ward_efficiency": analysis.WardEfficiency,
			},
			"benchmark_metrics": analysis.GlobalBenchmark,
			"percentile":        analysis.GlobalBenchmark.PlayerPercentile,
		}
	}

	c.JSON(http.StatusOK, comparison)
}

// Register routes for vision analytics
func (vh *VisionHandler) RegisterRoutes(router *gin.RouterGroup) {
	vision := router.Group("/vision")
	{
		vision.GET("/:player_id/analysis", vh.GetVisionAnalysis)
		vision.GET("/:player_id/heatmap", vh.GetVisionHeatmap)
		vision.GET("/:player_id/recommendations", vh.GetVisionRecommendations)
		vision.GET("/:player_id/trends", vh.GetVisionTrends)
		vision.GET("/:player_id/zones", vh.GetVisionZoneAnalysis)
		vision.GET("/:player_id/comparison", vh.GetVisionComparison)
	}
}
