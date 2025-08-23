package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herald/internal/services"
)

// MetaHandler handles meta analysis and tier list requests
type MetaHandler struct {
	metaService *services.MetaAnalyticsService
}

// NewMetaHandler creates a new meta handler
func NewMetaHandler(metaService *services.MetaAnalyticsService) *MetaHandler {
	return &MetaHandler{
		metaService: metaService,
	}
}

// MetaAnalysisRequest represents request for meta analysis
type MetaAnalysisRequest struct {
	Patch     string `form:"patch" json:"patch" binding:"required"`
	Region    string `form:"region" json:"region"`
	Rank      string `form:"rank" json:"rank"`
	TimeRange string `form:"time_range" json:"time_range"`
}

// TierListRequest represents request for tier list
type TierListRequest struct {
	Patch  string `form:"patch" json:"patch" binding:"required"`
	Region string `form:"region" json:"region"`
	Rank   string `form:"rank" json:"rank"`
	Role   string `form:"role" json:"role"`
}

// ChampionMetaRequest represents request for champion meta stats
type ChampionMetaRequest struct {
	Champion string `form:"champion" json:"champion" binding:"required"`
	Patch    string `form:"patch" json:"patch" binding:"required"`
	Region   string `form:"region" json:"region"`
	Rank     string `form:"rank" json:"rank"`
}

// MetaTrendsRequest represents request for meta trends
type MetaTrendsRequest struct {
	Patch    string `form:"patch" json:"patch" binding:"required"`
	Region   string `form:"region" json:"region"`
	Rank     string `form:"rank" json:"rank"`
	Category string `form:"category" json:"category"`
}

// GetMetaAnalysis godoc
// @Summary Get comprehensive meta analysis
// @Description Provides complete meta analysis including tier lists, trends, and predictions
// @Tags meta
// @Accept json
// @Produce json
// @Param patch query string true "Patch version (e.g., 14.1)"
// @Param region query string false "Region filter (default: all)"
// @Param rank query string false "Rank filter (default: all)"
// @Param time_range query string false "Time range (default: 7d)"
// @Success 200 {object} services.MetaAnalysis
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/meta/analysis [get]
func (mh *MetaHandler) GetMetaAnalysis(c *gin.Context) {
	var req MetaAnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Validate patch format
	if req.Patch == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Patch version is required (e.g., 14.1)",
		})
		return
	}

	// Set defaults
	if req.Region == "" {
		req.Region = "all"
	}
	if req.Rank == "" {
		req.Rank = "all"
	}
	if req.TimeRange == "" {
		req.TimeRange = "7d"
	}

	// Validate time range
	if !isValidTimeRange(req.TimeRange) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid time range. Use: 7d, 14d, or 30d",
		})
		return
	}

	// Validate region if provided
	if req.Region != "all" {
		validRegions := map[string]bool{
			"na1":   true,
			"euw1":  true,
			"eun1":  true,
			"kr":    true,
			"jp1":   true,
			"br1":   true,
			"la1":   true,
			"la2":   true,
			"oc1":   true,
			"tr1":   true,
			"ru":    true,
		}
		if !validRegions[req.Region] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid region",
			})
			return
		}
	}

	// Validate rank if provided
	if req.Rank != "all" {
		validRanks := map[string]bool{
			"iron":     true,
			"bronze":   true,
			"silver":   true,
			"gold":     true,
			"platinum": true,
			"diamond":  true,
			"master":   true,
			"grandmaster": true,
			"challenger":  true,
		}
		if !validRanks[req.Rank] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid rank",
			})
			return
		}
	}

	// Perform meta analysis
	analysis, err := mh.metaService.AnalyzeMeta(c.Request.Context(), req.Patch, req.Region, req.Rank, req.TimeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to analyze meta data",
		})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetTierList godoc
// @Summary Get champion tier list
// @Description Returns tier list for champions based on current meta
// @Tags meta
// @Accept json
// @Produce json
// @Param patch query string true "Patch version (e.g., 14.1)"
// @Param region query string false "Region filter (default: all)"
// @Param rank query string false "Rank filter (default: all)"
// @Param role query string false "Role filter (TOP, JUNGLE, MID, ADC, SUPPORT)"
// @Success 200 {object} services.ChampionTierList
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/meta/tier-list [get]
func (mh *MetaHandler) GetTierList(c *gin.Context) {
	var req TierListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Validate patch
	if req.Patch == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Patch version is required",
		})
		return
	}

	// Set defaults
	if req.Region == "" {
		req.Region = "all"
	}
	if req.Rank == "" {
		req.Rank = "all"
	}
	if req.Role == "" {
		req.Role = "ALL"
	}

	// Validate role if provided
	if req.Role != "ALL" && !isValidPosition(req.Role) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid role. Use: TOP, JUNGLE, MID, ADC, SUPPORT, or ALL",
		})
		return
	}

	// Get tier list
	tierList, err := mh.metaService.GetTierList(c.Request.Context(), req.Patch, req.Region, req.Rank, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to get tier list",
		})
		return
	}

	c.JSON(http.StatusOK, tierList)
}

// GetChampionMeta godoc
// @Summary Get champion meta statistics
// @Description Returns detailed meta statistics for a specific champion
// @Tags meta
// @Accept json
// @Produce json
// @Param champion query string true "Champion name"
// @Param patch query string true "Patch version (e.g., 14.1)"
// @Param region query string false "Region filter (default: all)"
// @Param rank query string false "Rank filter (default: all)"
// @Success 200 {object} services.ChampionMetaStats
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/meta/champion [get]
func (mh *MetaHandler) GetChampionMeta(c *gin.Context) {
	var req ChampionMetaRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Validate required fields
	if req.Champion == "" || req.Patch == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Champion name and patch version are required",
		})
		return
	}

	// Set defaults
	if req.Region == "" {
		req.Region = "all"
	}
	if req.Rank == "" {
		req.Rank = "all"
	}

	// Get champion meta stats
	stats, err := mh.metaService.GetChampionMetaStats(c.Request.Context(), req.Champion, req.Patch, req.Region, req.Rank)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to get champion meta statistics",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetMetaTrends godoc
// @Summary Get meta trends analysis
// @Description Returns current meta trends and strategic shifts
// @Tags meta
// @Accept json
// @Produce json
// @Param patch query string true "Patch version (e.g., 14.1)"
// @Param region query string false "Region filter (default: all)"
// @Param rank query string false "Rank filter (default: all)"
// @Param category query string false "Trend category (strategies, champions, items)"
// @Success 200 {object} services.MetaTrendsData
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/meta/trends [get]
func (mh *MetaHandler) GetMetaTrends(c *gin.Context) {
	var req MetaTrendsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Validate patch
	if req.Patch == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Patch version is required",
		})
		return
	}

	// Set defaults
	if req.Region == "" {
		req.Region = "all"
	}
	if req.Rank == "" {
		req.Rank = "all"
	}

	// Validate category if provided
	if req.Category != "" {
		validCategories := map[string]bool{
			"strategies": true,
			"champions":  true,
			"items":      true,
			"objectives": true,
		}
		if !validCategories[req.Category] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid category. Use: strategies, champions, items, or objectives",
			})
			return
		}
	}

	// Get meta trends
	trends, err := mh.metaService.GetMetaTrends(c.Request.Context(), req.Patch, req.Region, req.Rank, req.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to get meta trends",
		})
		return
	}

	c.JSON(http.StatusOK, trends)
}

// GetBanAnalysis godoc
// @Summary Get ban phase analysis
// @Description Returns ban phase statistics and recommendations
// @Tags meta
// @Accept json
// @Produce json
// @Param patch query string true "Patch version (e.g., 14.1)"
// @Param region query string false "Region filter (default: all)"
// @Param rank query string false "Rank filter (default: all)"
// @Param ban_type query string false "Ban type filter (power_bans, target_bans, flex_bans)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/meta/bans [get]
func (mh *MetaHandler) GetBanAnalysis(c *gin.Context) {
	patch := c.Query("patch")
	region := c.DefaultQuery("region", "all")
	rank := c.DefaultQuery("rank", "all")
	banType := c.Query("ban_type")

	if patch == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Patch version is required",
		})
		return
	}

	// Validate ban type if provided
	if banType != "" {
		validBanTypes := map[string]bool{
			"power_bans":  true,
			"target_bans": true,
			"flex_bans":   true,
		}
		if !validBanTypes[banType] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid ban type. Use: power_bans, target_bans, or flex_bans",
			})
			return
		}
	}

	// Get full meta analysis to extract ban data
	analysis, err := mh.metaService.AnalyzeMeta(c.Request.Context(), patch, region, rank, "7d")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to get ban analysis",
		})
		return
	}

	banAnalysis := map[string]interface{}{
		"patch":                patch,
		"region":               region,
		"rank":                 rank,
		"top_banned_champions": analysis.BanAnalysis.TopBannedChampions,
		"ban_strategies":       analysis.BanAnalysis.BanStrategies,
		"role_targeting":       analysis.BanAnalysis.RoleTargeting,
		"power_bans":           analysis.BanAnalysis.PowerBans,
	}

	c.JSON(http.StatusOK, banAnalysis)
}

// GetPickAnalysis godoc
// @Summary Get pick phase analysis
// @Description Returns pick phase statistics and recommendations
// @Tags meta
// @Accept json
// @Produce json
// @Param patch query string true "Patch version (e.g., 14.1)"
// @Param region query string false "Region filter (default: all)"
// @Param rank query string false "Rank filter (default: all)"
// @Param pick_type query string false "Pick type filter (blind_pick, flex_pick, counter_pick)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/meta/picks [get]
func (mh *MetaHandler) GetPickAnalysis(c *gin.Context) {
	patch := c.Query("patch")
	region := c.DefaultQuery("region", "all")
	rank := c.DefaultQuery("rank", "all")
	pickType := c.Query("pick_type")

	if patch == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Patch version is required",
		})
		return
	}

	// Validate pick type if provided
	if pickType != "" {
		validPickTypes := map[string]bool{
			"blind_pick":   true,
			"flex_pick":    true,
			"counter_pick": true,
		}
		if !validPickTypes[pickType] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid pick type. Use: blind_pick, flex_pick, or counter_pick",
			})
			return
		}
	}

	// Get full meta analysis to extract pick data
	analysis, err := mh.metaService.AnalyzeMeta(c.Request.Context(), patch, region, rank, "7d")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to get pick analysis",
		})
		return
	}

	pickAnalysis := map[string]interface{}{
		"patch":               patch,
		"region":              region,
		"rank":                rank,
		"blind_pick_safe":     analysis.PickAnalysis.BlindPickSafe,
		"flex_picks":          analysis.PickAnalysis.FlexPicks,
		"counter_picks":       analysis.PickAnalysis.CounterPicks,
		"first_pick_priority": analysis.PickAnalysis.FirstPickPriority,
		"last_pick_options":   analysis.PickAnalysis.LastPickOptions,
	}

	c.JSON(http.StatusOK, pickAnalysis)
}

// GetMetaPredictions godoc
// @Summary Get meta predictions
// @Description Returns predictions for upcoming meta shifts and champion changes
// @Tags meta
// @Accept json
// @Produce json
// @Param patch query string true "Current patch version (e.g., 14.1)"
// @Param prediction_type query string false "Prediction type (champions, strategies, items)"
// @Param confidence_threshold query float64 false "Minimum confidence threshold (0.0-1.0)"
// @Success 200 {object} services.MetaPredictions
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/meta/predictions [get]
func (mh *MetaHandler) GetMetaPredictions(c *gin.Context) {
	patch := c.Query("patch")
	predictionType := c.Query("prediction_type")
	confidenceThresholdStr := c.Query("confidence_threshold")

	if patch == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Patch version is required",
		})
		return
	}

	// Validate prediction type if provided
	if predictionType != "" {
		validPredictionTypes := map[string]bool{
			"champions":   true,
			"strategies":  true,
			"items":       true,
		}
		if !validPredictionTypes[predictionType] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid prediction type. Use: champions, strategies, or items",
			})
			return
		}
	}

	// Parse confidence threshold
	var confidenceThreshold float64 = 0.0
	if confidenceThresholdStr != "" {
		threshold, err := strconv.ParseFloat(confidenceThresholdStr, 64)
		if err != nil || threshold < 0.0 || threshold > 1.0 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid confidence threshold. Must be between 0.0 and 1.0",
			})
			return
		}
		confidenceThreshold = threshold
	}

	// Get full meta analysis to extract predictions
	analysis, err := mh.metaService.AnalyzeMeta(c.Request.Context(), patch, "all", "all", "7d")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to get meta predictions",
		})
		return
	}

	predictions := analysis.Predictions

	// Filter predictions by confidence threshold
	if confidenceThreshold > 0.0 {
		filteredTierPredictions := []services.ChampionTierPrediction{}
		for _, pred := range predictions.NextPatchPredictions {
			if pred.Confidence >= confidenceThreshold {
				filteredTierPredictions = append(filteredTierPredictions, pred)
			}
		}
		predictions.NextPatchPredictions = filteredTierPredictions

		filteredEmergingPredictions := []services.EmergingPrediction{}
		for _, pred := range predictions.EmergingChampions {
			if pred.Confidence >= confidenceThreshold {
				filteredEmergingPredictions = append(filteredEmergingPredictions, pred)
			}
		}
		predictions.EmergingChampions = filteredEmergingPredictions
	}

	c.JSON(http.StatusOK, predictions)
}

// GetMetaRecommendations godoc
// @Summary Get meta-based recommendations
// @Description Returns personalized recommendations based on current meta
// @Tags meta
// @Accept json
// @Produce json
// @Param patch query string true "Patch version (e.g., 14.1)"
// @Param rank query string false "Target rank (default: all)"
// @Param role query string false "Role filter"
// @Param recommendation_type query string false "Type of recommendation (champion_pool, strategy, builds)"
// @Success 200 {object} []services.MetaRecommendation
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/meta/recommendations [get]
func (mh *MetaHandler) GetMetaRecommendations(c *gin.Context) {
	patch := c.Query("patch")
	rank := c.DefaultQuery("rank", "all")
	role := c.Query("role")
	recommendationType := c.Query("recommendation_type")

	if patch == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Patch version is required",
		})
		return
	}

	// Validate recommendation type if provided
	if recommendationType != "" {
		validRecommendationTypes := map[string]bool{
			"champion_pool": true,
			"strategy":      true,
			"builds":        true,
			"bans":          true,
		}
		if !validRecommendationTypes[recommendationType] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid recommendation type. Use: champion_pool, strategy, builds, or bans",
			})
			return
		}
	}

	// Validate role if provided
	if role != "" && !isValidPosition(role) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid role. Use: TOP, JUNGLE, MID, ADC, or SUPPORT",
		})
		return
	}

	// Get full meta analysis to extract recommendations
	analysis, err := mh.metaService.AnalyzeMeta(c.Request.Context(), patch, "all", rank, "7d")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "analysis_error",
			Message: "Failed to get meta recommendations",
		})
		return
	}

	recommendations := analysis.Recommendations

	// Filter recommendations by type if specified
	if recommendationType != "" {
		filteredRecommendations := []services.MetaRecommendation{}
		for _, rec := range recommendations {
			if rec.Type == recommendationType {
				filteredRecommendations = append(filteredRecommendations, rec)
			}
		}
		recommendations = filteredRecommendations
	}

	// Filter recommendations by target rank
	if rank != "all" {
		filteredRecommendations := []services.MetaRecommendation{}
		for _, rec := range recommendations {
			if rec.TargetRank == rank || rec.TargetRank == "All" {
				filteredRecommendations = append(filteredRecommendations, rec)
			}
		}
		recommendations = filteredRecommendations
	}

	c.JSON(http.StatusOK, recommendations)
}

// GetMetaHistory godoc
// @Summary Get meta history and evolution
// @Description Returns historical meta data showing evolution over patches
// @Tags meta
// @Accept json
// @Produce json
// @Param start_patch query string true "Start patch version (e.g., 14.1)"
// @Param end_patch query string false "End patch version (default: current)"
// @Param champion query string false "Champion filter"
// @Param metric query string false "Metric to track (win_rate, pick_rate, ban_rate, tier)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/v1/meta/history [get]
func (mh *MetaHandler) GetMetaHistory(c *gin.Context) {
	startPatch := c.Query("start_patch")
	endPatch := c.DefaultQuery("end_patch", "current")
	champion := c.Query("champion")
	metric := c.DefaultQuery("metric", "tier")

	if startPatch == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Start patch version is required",
		})
		return
	}

	// Validate metric
	validMetrics := map[string]bool{
		"win_rate":  true,
		"pick_rate": true,
		"ban_rate":  true,
		"tier":      true,
		"presence":  true,
	}
	if !validMetrics[metric] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid metric. Use: win_rate, pick_rate, ban_rate, tier, or presence",
		})
		return
	}

	// This would typically query historical meta data
	// For now, we'll return a simulated response
	history := map[string]interface{}{
		"start_patch": startPatch,
		"end_patch":   endPatch,
		"champion":    champion,
		"metric":      metric,
		"data_points": []map[string]interface{}{
			{
				"patch": startPatch,
				"value": 52.5,
				"tier":  "A",
			},
		},
		"trend_analysis": map[string]interface{}{
			"direction": "stable",
			"change":    0.5,
			"volatility": "low",
		},
	}

	c.JSON(http.StatusOK, history)
}

// Register routes for meta analytics
func (mh *MetaHandler) RegisterRoutes(router *gin.RouterGroup) {
	meta := router.Group("/meta")
	{
		meta.GET("/analysis", mh.GetMetaAnalysis)
		meta.GET("/tier-list", mh.GetTierList)
		meta.GET("/champion", mh.GetChampionMeta)
		meta.GET("/trends", mh.GetMetaTrends)
		meta.GET("/bans", mh.GetBanAnalysis)
		meta.GET("/picks", mh.GetPickAnalysis)
		meta.GET("/predictions", mh.GetMetaPredictions)
		meta.GET("/recommendations", mh.GetMetaRecommendations)
		meta.GET("/history", mh.GetMetaHistory)
	}
}