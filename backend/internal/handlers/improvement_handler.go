package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herald-lol/herald/backend/internal/services"
)

type ImprovementHandler struct {
	improvementService *services.ImprovementRecommendationsService
}

func NewImprovementHandler(improvementService *services.ImprovementRecommendationsService) *ImprovementHandler {
	return &ImprovementHandler{
		improvementService: improvementService,
	}
}

// RegisterRoutes registers all improvement recommendation routes
func (h *ImprovementHandler) RegisterRoutes(r *gin.RouterGroup) {
	improvement := r.Group("/improvement")
	{
		// Core Recommendation Endpoints
		improvement.GET("/recommendations/:summoner_id", h.GetPersonalizedRecommendations)
		improvement.GET("/recommendations/:summoner_id/active", h.GetActiveRecommendations)
		improvement.GET("/recommendation/:recommendation_id", h.GetRecommendationDetails)
		improvement.PUT("/recommendation/:recommendation_id/progress", h.UpdateRecommendationProgress)
		improvement.POST("/recommendation/:recommendation_id/complete", h.CompleteRecommendation)
		
		// Analysis and Insights
		improvement.GET("/analysis/:summoner_id", h.GetPlayerAnalysis)
		improvement.GET("/insights/:summoner_id", h.GetImprovementInsights)
		improvement.GET("/progress/:summoner_id", h.GetOverallProgress)
		
		// Specialized Recommendations
		improvement.GET("/quick-wins/:summoner_id", h.GetQuickWins)
		improvement.GET("/coaching-plan/:summoner_id", h.GetCoachingPlan)
	}
}

// GetPersonalizedRecommendations handles personalized recommendation requests
func (h *ImprovementHandler) GetPersonalizedRecommendations(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summoner_id is required",
		})
		return
	}

	// Parse query parameters
	options := services.RecommendationOptions{
		MaxRecommendations: 10, // default
	}

	if focusCategory := c.Query("focus_category"); focusCategory != "" {
		options.FocusCategory = focusCategory
	}

	if difficultyFilter := c.Query("difficulty_filter"); difficultyFilter != "" {
		options.DifficultyFilter = difficultyFilter
	}

	if timeConstraintStr := c.Query("time_constraint"); timeConstraintStr != "" {
		if timeConstraint, err := strconv.Atoi(timeConstraintStr); err == nil {
			options.TimeConstraint = timeConstraint
		}
	}

	if maxRecsStr := c.Query("max_recommendations"); maxRecsStr != "" {
		if maxRecs, err := strconv.Atoi(maxRecsStr); err == nil && maxRecs > 0 {
			options.MaxRecommendations = maxRecs
		}
	}

	options.IncludeAlternatives = c.Query("include_alternatives") == "true"

	recommendations, err := h.improvementService.GetPersonalizedRecommendations(summonerID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate personalized recommendations",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summoner_id":     summonerID,
		"recommendations": recommendations,
		"options":         options,
		"generated_at":    "2024-01-15T10:30:00Z",
	})
}

// GetActiveRecommendations handles active recommendations requests
func (h *ImprovementHandler) GetActiveRecommendations(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summoner_id is required",
		})
		return
	}

	recommendations, err := h.improvementService.GetActiveRecommendations(summonerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get active recommendations",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summoner_id":     summonerID,
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// GetRecommendationDetails handles individual recommendation detail requests
func (h *ImprovementHandler) GetRecommendationDetails(c *gin.Context) {
	recommendationID := c.Param("recommendation_id")
	if recommendationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "recommendation_id is required",
		})
		return
	}

	// Get recommendation progress
	progress, err := h.improvementService.GetRecommendationProgress(recommendationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get recommendation details",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendation_id": recommendationID,
		"progress":          progress,
	})
}

// UpdateRecommendationProgress handles progress update requests
func (h *ImprovementHandler) UpdateRecommendationProgress(c *gin.Context) {
	recommendationID := c.Param("recommendation_id")
	if recommendationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "recommendation_id is required",
		})
		return
	}

	var progressUpdate services.ProgressTrackingData
	if err := c.ShouldBindJSON(&progressUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid progress data",
			"details": err.Error(),
		})
		return
	}

	err := h.improvementService.UpdateRecommendationProgress(recommendationID, progressUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update recommendation progress",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendation_id": recommendationID,
		"status":            "progress_updated",
		"message":           "Progress updated successfully",
	})
}

// CompleteRecommendation handles recommendation completion requests
func (h *ImprovementHandler) CompleteRecommendation(c *gin.Context) {
	recommendationID := c.Param("recommendation_id")
	if recommendationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "recommendation_id is required",
		})
		return
	}

	err := h.improvementService.CompleteRecommendation(recommendationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to complete recommendation",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendation_id": recommendationID,
		"status":            "completed",
		"message":           "Recommendation completed successfully",
	})
}

// GetPlayerAnalysis handles player analysis requests
func (h *ImprovementHandler) GetPlayerAnalysis(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summoner_id is required",
		})
		return
	}

	// Mock analysis response - in real implementation this would call the service
	analysis := gin.H{
		"summoner_id":    summonerID,
		"analysis_date":  "2024-01-15T10:30:00Z",
		"overall_rating": 72.5,
		"skill_breakdown": gin.H{
			"mechanical_skill":  68.0,
			"game_knowledge":    75.0,
			"map_awareness":     65.0,
			"team_fighting":     78.0,
			"laning":            70.0,
			"objective_control": 73.0,
			"vision_control":    60.0,
			"positioning":       72.0,
			"decision_making":   74.0,
			"mental_resilience": 69.0,
		},
		"improvement_potential": gin.H{
			"vision_control":    20.0,
			"map_awareness":     18.0,
			"mechanical_skill":  15.0,
			"positioning":       12.0,
			"mental_resilience": 12.0,
		},
		"critical_weaknesses": []gin.H{
			{
				"area":               "vision_control",
				"severity":           "high",
				"impact_on_winrate":  8.5,
				"frequency":          85.0,
				"root_causes":        []string{"infrequent ward placement", "poor ward positioning", "not clearing enemy vision"},
				"quick_wins":         []string{"buy more control wards", "ward before major objectives", "use trinket on cooldown"},
			},
		},
		"competitive_benchmark": gin.H{
			"rank_tier":            "Gold II",
			"regional_percentile":  68.5,
			"stronger_than_peers":  []string{"team_fighting", "decision_making"},
			"weaker_than_peers":    []string{"vision_control", "map_awareness"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"analysis": analysis,
	})
}

// GetImprovementInsights handles improvement insights requests
func (h *ImprovementHandler) GetImprovementInsights(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summoner_id is required",
		})
		return
	}

	insights := gin.H{
		"summoner_id": summonerID,
		"insights": []gin.H{
			{
				"type":        "opportunity",
				"priority":    "high",
				"title":       "Vision Control Quick Win",
				"description": "Improving ward placement can increase win rate by 8-10% with minimal effort",
				"impact":      "high",
				"difficulty":  "easy",
				"timeframe":   "1-2 weeks",
				"action_steps": []string{
					"Purchase control ward every back",
					"Ward objectives 30 seconds before spawn",
					"Use trinket ward on cooldown",
				},
			},
			{
				"type":        "trend",
				"priority":    "medium",
				"title":       "Consistent Improvement Pattern",
				"description": "Your recent games show steady improvement in team fighting",
				"impact":      "medium",
				"trend":       "positive",
				"continue_doing": []string{
					"Look for team fight opportunities",
					"Focus on positioning in fights",
					"Follow up on team engages",
				},
			},
			{
				"type":        "warning",
				"priority":    "medium",
				"title":       "Positioning Regression",
				"description": "Recent decline in positioning - address before it impacts climb",
				"impact":      "medium",
				"difficulty":  "medium",
				"suggested_focus": []string{
					"Review positioning mistakes in replays",
					"Practice safe farming positions",
					"Work on team fight positioning",
				},
			},
		],
		"overall_trajectory": gin.H{
			"direction":         "improving",
			"consistency":       "stable",
			"key_focus_areas":   []string{"vision_control", "positioning"},
			"strength_areas":    []string{"team_fighting", "decision_making"},
			"predicted_outcome": "Platinum within 2-3 months with focused improvement",
		},
	}

	c.JSON(http.StatusOK, insights)
}

// GetOverallProgress handles overall progress tracking requests
func (h *ImprovementHandler) GetOverallProgress(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summoner_id is required",
		})
		return
	}

	progress := gin.H{
		"summoner_id": summonerID,
		"overall_progress": gin.H{
			"improvement_score":    75.2,
			"active_recommendations": 4,
			"completed_recommendations": 2,
			"total_recommendations": 6,
			"streak_days": 12,
			"consistency_rating": 82.0,
		},
		"recent_achievements": []gin.H{
			{
				"achievement": "Vision Control Milestone",
				"description": "Achieved 50+ vision score average",
				"date":        "2024-01-10",
				"impact":      "6.5% win rate improvement",
			},
			{
				"achievement": "Consistency Streak",
				"description": "7 days of following improvement plan",
				"date":        "2024-01-08",
				"impact":      "Improved mental resilience",
			},
		},
		"skill_progress": gin.H{
			"vision_control": gin.H{
				"baseline":    60.0,
				"current":     72.5,
				"target":      80.0,
				"progress":    62.5, // percentage to target
				"trend":       "improving",
			},
			"positioning": gin.H{
				"baseline":    72.0,
				"current":     69.0,
				"target":      78.0,
				"progress":    -50.0, // negative progress
				"trend":       "declining",
			},
		},
		"weekly_summary": gin.H{
			"week_number":       3,
			"practice_minutes":  180,
			"games_played":     15,
			"recommendations_worked_on": 3,
			"improvement_areas_focused": []string{"vision_control", "map_awareness"},
			"win_rate_change": 8.5,
			"rank_progress":   125, // LP gained
		},
	}

	c.JSON(http.StatusOK, progress)
}

// GetQuickWins handles quick win recommendations requests
func (h *ImprovementHandler) GetQuickWins(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summoner_id is required",
		})
		return
	}

	quickWins := gin.H{
		"summoner_id": summonerID,
		"quick_wins": []gin.H{
			{
				"title":              "Buy Control Wards Every Back",
				"description":        "Simple habit that dramatically improves vision control",
				"expected_impact":    "5-8% win rate increase",
				"time_to_implement":  "immediate",
				"difficulty":         "very_easy",
				"roi_score":          95.0,
				"instructions": []string{
					"Always have 75g saved for control wards",
					"Place control wards in key bushes",
					"Replace destroyed control wards immediately",
				},
			},
			{
				"title":              "Look at Minimap Every 5 Seconds",
				"description":        "Develop map awareness through conscious habit",
				"expected_impact":    "3-5% win rate increase",
				"time_to_implement":  "3-5 days",
				"difficulty":         "easy",
				"roi_score":          85.0,
				"habit_formation": gin.H{
					"trigger":     "after every CS",
					"action":      "glance at minimap",
					"reward":      "better map awareness",
					"tracking":    "count map looks per minute",
				},
			},
			{
				"title":              "Ward Before Objectives Spawn",
				"description":        "Place vision 30-60 seconds before dragon/baron",
				"expected_impact":    "4-6% win rate increase",
				"time_to_implement":  "immediate",
				"difficulty":         "easy",
				"roi_score":          88.0,
				"timing_guide": gin.H{
					"dragon":      "30 seconds before spawn",
					"baron":       "45 seconds before spawn",
					"rift_herald": "30 seconds before spawn",
				},
			},
		},
		"implementation_plan": gin.H{
			"week_1": []string{"Focus on control ward purchases", "Begin minimap habit formation"},
			"week_2": []string{"Master objective warding timing", "Consolidate all habits"},
			"expected_cumulative_impact": "12-18% win rate improvement",
			"success_metrics": []string{
				"Control ward purchases per game > 3",
				"Map awareness checks per minute > 8",
				"Objective vision control > 70%",
			},
		},
	}

	c.JSON(http.StatusOK, quickWins)
}

// GetCoachingPlan handles comprehensive coaching plan requests
func (h *ImprovementHandler) GetCoachingPlan(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summoner_id is required",
		})
		return
	}

	// Parse query parameters
	duration := c.DefaultQuery("duration", "30") // days
	intensityStr := c.DefaultQuery("intensity", "moderate")
	
	coachingPlan := gin.H{
		"summoner_id": summonerID,
		"plan_details": gin.H{
			"duration_days":    duration,
			"intensity_level":  intensityStr,
			"primary_goal":     "Reach Platinum rank",
			"secondary_goals":  []string{"Improve consistency", "Master vision control", "Better positioning"},
			"estimated_outcome": "70% probability of reaching Platinum within 30 days",
		},
		"weekly_breakdown": []gin.H{
			{
				"week":           1,
				"focus_theme":    "Foundation & Quick Wins",
				"primary_skills": []string{"vision_control", "map_awareness"},
				"daily_tasks": gin.H{
					"practice_time":     20, // minutes
					"games_minimum":     3,
					"specific_focus":    []string{"Control ward usage", "Minimap checks"},
					"success_metrics":   []string{"Vision score > 40", "Deaths < 6 per game"},
				},
				"milestone": gin.H{
					"title":       "Vision Control Foundation",
					"description": "Establish consistent warding habits",
					"target":      "Average vision score 45+",
				},
			},
			{
				"week":           2,
				"focus_theme":    "Positioning & Safety",
				"primary_skills": []string{"positioning", "safety", "laning"},
				"daily_tasks": gin.H{
					"practice_time":     25,
					"games_minimum":     3,
					"specific_focus":    []string{"Safe farming positions", "Team fight positioning"},
					"success_metrics":   []string{"Deaths < 5 per game", "Kill participation > 60%"},
				},
				"milestone": gin.H{
					"title":       "Positioning Mastery",
					"description": "Reduce deaths through better positioning",
					"target":      "Average deaths < 5 per game",
				},
			},
			{
				"week":           3,
				"focus_theme":    "Team Fighting & Macro",
				"primary_skills": []string{"team_fighting", "objective_control", "decision_making"},
				"daily_tasks": gin.H{
					"practice_time":     30,
					"games_minimum":     4,
					"specific_focus":    []string{"Team fight positioning", "Objective timing"},
					"success_metrics":   []string{"Damage dealt > 20k per game", "Objective control participation > 70%"},
				},
				"milestone": gin.H{
					"title":       "Team Impact Excellence",
					"description": "Maximize impact in team scenarios",
					"target":      "Consistent positive KDA and high damage",
				},
			},
			{
				"week":           4,
				"focus_theme":    "Consistency & Climbing",
				"primary_skills": []string{"mental_game", "consistency", "adaptation"},
				"daily_tasks": gin.H{
					"practice_time":     25,
					"games_minimum":     5,
					"specific_focus":    []string{"Consistent performance", "Adapting to game state"},
					"success_metrics":   []string{"Win rate > 60%", "Performance consistency > 80%"},
				},
				"milestone": gin.H{
					"title":       "Ranking Consistency",
					"description": "Maintain high performance for climbing",
					"target":      "Consistent positive LP gains",
				},
			},
		},
		"daily_routine": gin.H{
			"warm_up": gin.H{
				"duration_minutes": 10,
				"activities":       []string{"Practice tool mechanics", "Review previous game mistakes"},
			},
			"focused_practice": gin.H{
				"duration_minutes": 15,
				"activities":       []string{"Specific skill drills", "Replay analysis"},
			},
			"ranked_games": gin.H{
				"minimum_games": 3,
				"maximum_games": 6,
				"focus_mindset":  []string{"Apply learned skills", "Stay positive", "Focus on improvement not just wins"},
			},
			"review_session": gin.H{
				"duration_minutes": 10,
				"activities":       []string{"Quick game review", "Update progress tracking", "Plan next session"},
			},
		},
		"progress_tracking": gin.H{
			"daily_metrics": []string{
				"Vision score average",
				"Deaths per game",
				"Kill participation %",
				"Damage dealt average",
				"LP change",
			},
			"weekly_review": []string{
				"Overall win rate",
				"Skill improvement ratings",
				"Milestone achievement",
				"Consistency metrics",
				"Rank progress",
			},
			"success_indicators": []string{
				"Consistent upward LP trend",
				"Decreasing death count",
				"Increasing vision score",
				"Better KDA consistency",
				"Improved game impact",
			},
		},
	}

	c.JSON(http.StatusOK, coachingPlan)
}