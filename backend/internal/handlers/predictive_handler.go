package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herald-lol/backend/internal/services"
)

type PredictiveHandler struct {
	predictiveService *services.PredictiveAnalyticsService
}

func NewPredictiveHandler(predictiveService *services.PredictiveAnalyticsService) *PredictiveHandler {
	return &PredictiveHandler{
		predictiveService: predictiveService,
	}
}

// RegisterRoutes registers all predictive analytics routes
func (h *PredictiveHandler) RegisterRoutes(r *gin.RouterGroup) {
	predictive := r.Group("/predictive")
	{
		// Performance Predictions
		predictive.GET("/performance/:summoner_id", h.GetPerformancePrediction)
		predictive.GET("/rank-progression/:summoner_id", h.GetRankProgressionPrediction)
		predictive.GET("/skill-development/:summoner_id", h.GetSkillDevelopmentForecast)
		
		// Champion Recommendations
		predictive.GET("/champion-recommendations/:summoner_id", h.GetChampionRecommendations)
		predictive.GET("/meta-adaptation/:summoner_id", h.GetMetaAdaptationForecast)
		
		// Team Analytics
		predictive.GET("/team-performance/:team_id", h.GetTeamPerformancePrediction)
		predictive.GET("/team-synergy/:team_id", h.GetTeamSynergyAnalysis)
		
		// Career Analytics
		predictive.GET("/career-trajectory/:summoner_id", h.GetCareerTrajectoryForecast)
		predictive.GET("/player-potential/:summoner_id", h.GetPlayerPotentialAssessment)
	}
}

// GetPerformancePrediction handles performance prediction requests
func (h *PredictiveHandler) GetPerformancePrediction(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summoner_id is required",
		})
		return
	}

	// Parse query parameters
	horizon := c.DefaultQuery("horizon", "short_term")
	champion := c.Query("champion")
	role := c.Query("role")
	patch := c.Query("patch")

	prediction, err := h.predictiveService.GetPerformancePrediction(
		summonerID,
		horizon,
		champion,
		role,
		patch,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate performance prediction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summoner_id": summonerID,
		"prediction": prediction,
	})
}

// GetRankProgressionPrediction handles rank progression prediction requests
func (h *PredictiveHandler) GetRankProgressionPrediction(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summoner_id is required",
		})
		return
	}

	// Parse query parameters
	targetRank := c.Query("target_rank")
	timeframeStr := c.DefaultQuery("timeframe_days", "30")
	
	timeframeDays, err := strconv.Atoi(timeframeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid timeframe_days parameter",
		})
		return
	}

	prediction, err := h.predictiveService.GetRankProgressionPrediction(
		summonerID,
		targetRank,
		timeframeDays,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate rank progression prediction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summoner_id": summonerID,
		"prediction": prediction,
	})
}

// GetSkillDevelopmentForecast handles skill development forecast requests
func (h *PredictiveHandler) GetSkillDevelopmentForecast(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summoner_id is required",
		})
		return
	}

	// Parse query parameters
	skillCategory := c.Query("skill_category")
	forecastPeriodStr := c.DefaultQuery("forecast_period_days", "90")
	
	forecastPeriod, err := strconv.Atoi(forecastPeriodStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid forecast_period_days parameter",
		})
		return
	}

	forecast, err := h.predictiveService.GetSkillDevelopmentForecast(
		summonerID,
		skillCategory,
		forecastPeriod,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate skill development forecast",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summoner_id": summonerID,
		"forecast": forecast,
	})
}

// GetChampionRecommendations handles champion recommendation requests
func (h *PredictiveHandler) GetChampionRecommendations(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summoner_id is required",
		})
		return
	}

	// Parse query parameters
	role := c.Query("role")
	playstyle := c.Query("playstyle")
	metaFocus := c.DefaultQuery("meta_focus", "current")
	limitStr := c.DefaultQuery("limit", "10")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit parameter",
		})
		return
	}

	recommendations, err := h.predictiveService.GetChampionRecommendations(
		summonerID,
		role,
		playstyle,
		metaFocus,
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate champion recommendations",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summoner_id": summonerID,
		"recommendations": recommendations,
	})
}

// GetMetaAdaptationForecast handles meta adaptation forecast requests
func (h *PredictiveHandler) GetMetaAdaptationForecast(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summoner_id is required",
		})
		return
	}

	// Parse query parameters
	targetPatch := c.Query("target_patch")
	adaptationStyle := c.DefaultQuery("adaptation_style", "gradual")

	forecast, err := h.predictiveService.GetMetaAdaptationForecast(
		summonerID,
		targetPatch,
		adaptationStyle,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate meta adaptation forecast",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summoner_id": summonerID,
		"forecast": forecast,
	})
}

// GetTeamPerformancePrediction handles team performance prediction requests
func (h *PredictiveHandler) GetTeamPerformancePrediction(c *gin.Context) {
	teamID := c.Param("team_id")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "team_id is required",
		})
		return
	}

	// Parse query parameters
	gameType := c.DefaultQuery("game_type", "ranked")
	predictionScope := c.DefaultQuery("prediction_scope", "next_games")
	gamesCountStr := c.DefaultQuery("games_count", "10")
	
	gamesCount, err := strconv.Atoi(gamesCountStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid games_count parameter",
		})
		return
	}

	prediction, err := h.predictiveService.GetTeamPerformancePrediction(
		teamID,
		gameType,
		predictionScope,
		gamesCount,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate team performance prediction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"team_id": teamID,
		"prediction": prediction,
	})
}

// GetTeamSynergyAnalysis handles team synergy analysis requests
func (h *PredictiveHandler) GetTeamSynergyAnalysis(c *gin.Context) {
	teamID := c.Param("team_id")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "team_id is required",
		})
		return
	}

	// Parse query parameters
	analysisDepth := c.DefaultQuery("analysis_depth", "comprehensive")
	includeRecommendations := c.DefaultQuery("include_recommendations", "true") == "true"

	analysis, err := h.predictiveService.GetTeamSynergyAnalysis(
		teamID,
		analysisDepth,
		includeRecommendations,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate team synergy analysis",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"team_id": teamID,
		"analysis": analysis,
	})
}

// GetCareerTrajectoryForecast handles career trajectory forecast requests
func (h *PredictiveHandler) GetCareerTrajectoryForecast(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summoner_id is required",
		})
		return
	}

	// Parse query parameters
	forecastHorizon := c.DefaultQuery("forecast_horizon", "long_term")
	careerGoals := c.Query("career_goals")
	analysisDepth := c.DefaultQuery("analysis_depth", "detailed")

	forecast, err := h.predictiveService.GetCareerTrajectoryForecast(
		summonerID,
		forecastHorizon,
		careerGoals,
		analysisDepth,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate career trajectory forecast",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summoner_id": summonerID,
		"forecast": forecast,
	})
}

// GetPlayerPotentialAssessment handles player potential assessment requests
func (h *PredictiveHandler) GetPlayerPotentialAssessment(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summoner_id is required",
		})
		return
	}

	// Parse query parameters
	assessmentType := c.DefaultQuery("assessment_type", "comprehensive")
	includeRecommendations := c.DefaultQuery("include_recommendations", "true") == "true"
	competitiveLevel := c.Query("competitive_level")

	assessment, err := h.predictiveService.GetPlayerPotentialAssessment(
		summonerID,
		assessmentType,
		includeRecommendations,
		competitiveLevel,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate player potential assessment",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summoner_id": summonerID,
		"assessment": assessment,
	})
}