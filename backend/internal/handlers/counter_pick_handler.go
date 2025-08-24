// Counter Pick Handler for Herald.lol
package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/herald-lol/herald/backend/internal/services"
)

type CounterPickHandler struct {
	service *services.CounterPickService
}

func NewCounterPickHandler(service *services.CounterPickService) *CounterPickHandler {
	return &CounterPickHandler{
		service: service,
	}
}

func (h *CounterPickHandler) RegisterRoutes(rg *gin.RouterGroup) {
	counterPicks := rg.Group("/counter-picks")
	{
		// Single target counter analysis
		counterPicks.POST("/analyze", h.AnalyzeCounterPicks)
		counterPicks.GET("/suggestions/:champion/:role", h.GetCounterSuggestions)
		counterPicks.GET("/matchup/:champion1/:champion2", h.GetMatchupAnalysis)

		// Multi-target counter analysis
		counterPicks.POST("/multi-target", h.AnalyzeMultiTargetCounters)
		counterPicks.POST("/team-counters", h.GetTeamCounterStrategies)

		// Lane and phase specific
		counterPicks.GET("/lane-counters/:champion/:role", h.GetLaneCounters)
		counterPicks.GET("/teamfight-counters/:champion", h.GetTeamFightCounters)
		counterPicks.GET("/item-counters/:champion", h.GetItemCounters)
		counterPicks.GET("/playstyle-counters/:champion", h.GetPlayStyleCounters)

		// Meta and statistical data
		counterPicks.GET("/meta-counters", h.GetMetaCounters)
		counterPicks.GET("/winrate-data/:champion/:target", h.GetWinRateData)
		counterPicks.GET("/popularity/:champion", h.GetCounterPopularity)

		// User-specific features
		counterPicks.GET("/personalized/:summoner_id", h.GetPersonalizedCounters)
		counterPicks.POST("/favorites", h.AddToFavorites)
		counterPicks.GET("/favorites/:summoner_id", h.GetFavoriteCounters)
		counterPicks.DELETE("/favorites/:id", h.RemoveFromFavorites)

		// Historical data and learning
		counterPicks.GET("/history/:summoner_id", h.GetCounterHistory)
		counterPicks.POST("/feedback", h.SubmitCounterFeedback)
		counterPicks.GET("/performance/:summoner_id/:champion", h.GetCounterPerformance)

		// Ban strategy integration
		counterPicks.POST("/ban-strategy", h.GetCounterBanStrategy)
		counterPicks.GET("/threat-assessment", h.AssessThreatLevel)
	}
}

// AnalyzeCounterPicks analyzes counter picks for a specific target champion
func (h *CounterPickHandler) AnalyzeCounterPicks(c *gin.Context) {
	var request struct {
		TargetChampion     string   `json:"targetChampion" binding:"required"`
		TargetRole         string   `json:"targetRole" binding:"required"`
		GameMode           string   `json:"gameMode"`
		PlayerChampionPool []string `json:"playerChampionPool"`
		PlayerRank         string   `json:"playerRank"`
		Preferences        struct {
			PrioritizeLane      bool `json:"prioritizeLane"`
			PrioritizeTeamfight bool `json:"prioritizeTeamfight"`
			PrioritizeMeta      bool `json:"prioritizeMeta"`
			PrioritizeComfort   bool `json:"prioritizeComfort"`
		} `json:"preferences"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default game mode
	if request.GameMode == "" {
		request.GameMode = "ranked"
	}

	analysis, err := h.service.AnalyzeCounterPicks(
		c.Request.Context(),
		request.TargetChampion,
		request.TargetRole,
		request.GameMode,
		request.PlayerChampionPool,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetCounterSuggestions gets counter suggestions for a specific champion and role
func (h *CounterPickHandler) GetCounterSuggestions(c *gin.Context) {
	champion := c.Param("champion")
	role := c.Param("role")

	if champion == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Champion and role are required"})
		return
	}

	// Get query parameters
	gameMode := c.DefaultQuery("gameMode", "ranked")
	limitStr := c.DefaultQuery("limit", "15")
	minStrength := c.DefaultQuery("minStrength", "60")
	playerChampions := c.QueryArray("playerChampions")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 15
	}

	minStrengthFloat, err := strconv.ParseFloat(minStrength, 64)
	if err != nil {
		minStrengthFloat = 60.0
	}

	analysis, err := h.service.AnalyzeCounterPicks(
		c.Request.Context(),
		champion,
		role,
		gameMode,
		playerChampions,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter and limit results
	var filteredCounters []interface{}
	for _, counter := range analysis.CounterPicks {
		if counter.CounterStrength >= minStrengthFloat && len(filteredCounters) < limit {
			filteredCounters = append(filteredCounters, counter)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"targetChampion": champion,
		"targetRole":     role,
		"gameMode":       gameMode,
		"suggestions":    filteredCounters,
		"metaContext":    analysis.MetaContext,
		"confidence":     analysis.Confidence,
	})
}

// GetMatchupAnalysis gets detailed matchup analysis between two champions
func (h *CounterPickHandler) GetMatchupAnalysis(c *gin.Context) {
	champion1 := c.Param("champion1")
	champion2 := c.Param("champion2")

	if champion1 == "" || champion2 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both champions are required"})
		return
	}

	// Get query parameters
	role := c.DefaultQuery("role", "mid")
	gameMode := c.DefaultQuery("gameMode", "ranked")

	// Analyze both directions
	analysis1, err := h.service.AnalyzeCounterPicks(
		c.Request.Context(),
		champion2, // champion2 is the target
		role,
		gameMode,
		[]string{champion1}, // champion1 is the counter
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	analysis2, err := h.service.AnalyzeCounterPicks(
		c.Request.Context(),
		champion1, // champion1 is the target
		role,
		gameMode,
		[]string{champion2}, // champion2 is the counter
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Find specific matchup data
	var matchup1, matchup2 interface{}
	if len(analysis1.CounterPicks) > 0 {
		matchup1 = analysis1.CounterPicks[0]
	}
	if len(analysis2.CounterPicks) > 0 {
		matchup2 = analysis2.CounterPicks[0]
	}

	c.JSON(http.StatusOK, gin.H{
		"champion1":            champion1,
		"champion2":            champion2,
		"role":                 role,
		"champion1VsChampion2": matchup1,
		"champion2VsChampion1": matchup2,
		"lanePhases":           analysis1.LaneCounters,
		"teamFightData":        analysis1.TeamFightCounters,
		"itemCounters":         analysis1.ItemCounters,
	})
}

// AnalyzeMultiTargetCounters analyzes counters for multiple target champions
func (h *CounterPickHandler) AnalyzeMultiTargetCounters(c *gin.Context) {
	var request struct {
		TargetChampions []struct {
			Champion    string  `json:"champion" binding:"required"`
			Role        string  `json:"role" binding:"required"`
			ThreatLevel string  `json:"threatLevel"` // low, medium, high, critical
			Priority    float64 `json:"priority"`    // 0-100
		} `json:"targetChampions" binding:"required"`
		PlayerChampionPool []string `json:"playerChampionPool"`
		GameMode           string   `json:"gameMode"`
		Strategy           string   `json:"strategy"` // universal, specific, hybrid
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default values
	if request.GameMode == "" {
		request.GameMode = "ranked"
	}
	if request.Strategy == "" {
		request.Strategy = "hybrid"
	}

	// Convert request format to service format
	var targetChampions []services.TargetChampionData
	for _, target := range request.TargetChampions {
		threatLevel := target.ThreatLevel
		if threatLevel == "" {
			threatLevel = "medium"
		}

		priority := target.Priority
		if priority == 0 {
			priority = 50
		}

		targetChampions = append(targetChampions, services.TargetChampionData{
			Champion:    target.Champion,
			Role:        target.Role,
			ThreatLevel: threatLevel,
			Priority:    priority,
			Reasons:     []string{fmt.Sprintf("Priority target: %s", target.Champion)},
		})
	}

	analysis, err := h.service.AnalyzeMultiTargetCounters(
		c.Request.Context(),
		targetChampions,
		request.GameMode,
		request.PlayerChampionPool,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetTeamCounterStrategies gets team-based counter strategies
func (h *CounterPickHandler) GetTeamCounterStrategies(c *gin.Context) {
	var request struct {
		EnemyTeam []struct {
			Champion string `json:"champion" binding:"required"`
			Role     string `json:"role" binding:"required"`
		} `json:"enemyTeam" binding:"required"`
		OurTeam []struct {
			Champion string `json:"champion"`
			Role     string `json:"role" binding:"required"`
		} `json:"ourTeam"`
		GameMode string `json:"gameMode"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.GameMode == "" {
		request.GameMode = "ranked"
	}

	// Convert enemy team to target champions
	var targetChampions []services.TargetChampionData
	for _, enemy := range request.EnemyTeam {
		targetChampions = append(targetChampions, services.TargetChampionData{
			Champion:    enemy.Champion,
			Role:        enemy.Role,
			ThreatLevel: "medium",
			Priority:    70,
			Reasons:     []string{"Enemy team member"},
		})
	}

	analysis, err := h.service.AnalyzeMultiTargetCounters(
		c.Request.Context(),
		targetChampions,
		request.GameMode,
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"enemyTeam":          request.EnemyTeam,
		"ourTeam":            request.OurTeam,
		"teamStrategies":     analysis.TeamCounters,
		"universalCounters":  analysis.UniversalCounters,
		"banRecommendations": analysis.BanRecommendations,
		"overallStrategy":    analysis.OverallStrategy,
		"confidence":         analysis.Confidence,
	})
}

// GetLaneCounters gets lane-specific counter information
func (h *CounterPickHandler) GetLaneCounters(c *gin.Context) {
	champion := c.Param("champion")
	role := c.Param("role")

	if champion == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Champion and role are required"})
		return
	}

	gameMode := c.DefaultQuery("gameMode", "ranked")

	analysis, err := h.service.AnalyzeCounterPicks(
		c.Request.Context(),
		champion,
		role,
		gameMode,
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"targetChampion": champion,
		"targetRole":     role,
		"laneCounters":   analysis.LaneCounters,
		"confidence":     analysis.Confidence,
	})
}

// GetTeamFightCounters gets team fight counter information
func (h *CounterPickHandler) GetTeamFightCounters(c *gin.Context) {
	champion := c.Param("champion")

	if champion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Champion is required"})
		return
	}

	gameMode := c.DefaultQuery("gameMode", "ranked")
	role := c.DefaultQuery("role", "mid")

	analysis, err := h.service.AnalyzeCounterPicks(
		c.Request.Context(),
		champion,
		role,
		gameMode,
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"targetChampion":    champion,
		"teamFightCounters": analysis.TeamFightCounters,
		"confidence":        analysis.Confidence,
	})
}

// GetItemCounters gets item-based counter recommendations
func (h *CounterPickHandler) GetItemCounters(c *gin.Context) {
	champion := c.Param("champion")

	if champion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Champion is required"})
		return
	}

	gameMode := c.DefaultQuery("gameMode", "ranked")
	role := c.DefaultQuery("role", "mid")

	analysis, err := h.service.AnalyzeCounterPicks(
		c.Request.Context(),
		champion,
		role,
		gameMode,
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"targetChampion": champion,
		"itemCounters":   analysis.ItemCounters,
		"confidence":     analysis.Confidence,
	})
}

// GetPlayStyleCounters gets playstyle-based counter strategies
func (h *CounterPickHandler) GetPlayStyleCounters(c *gin.Context) {
	champion := c.Param("champion")

	if champion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Champion is required"})
		return
	}

	gameMode := c.DefaultQuery("gameMode", "ranked")
	role := c.DefaultQuery("role", "mid")

	analysis, err := h.service.AnalyzeCounterPicks(
		c.Request.Context(),
		champion,
		role,
		gameMode,
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"targetChampion":    champion,
		"playStyleCounters": analysis.PlayStyleCounters,
		"confidence":        analysis.Confidence,
	})
}

// GetMetaCounters gets current meta counter information
func (h *CounterPickHandler) GetMetaCounters(c *gin.Context) {
	// Get query parameters
	role := c.Query("role")
	gameMode := c.DefaultQuery("gameMode", "ranked")
	rank := c.DefaultQuery("rank", "all")
	region := c.DefaultQuery("region", "global")
	limitStr := c.DefaultQuery("limit", "20")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	// This would integrate with meta service to get current strong picks
	// and then analyze their counters
	strongPicks := []string{"Yasuo", "Zed", "Akali", "Katarina", "Ahri"} // Mock data

	var metaCounters []interface{}
	for _, champion := range strongPicks {
		targetRole := role
		if targetRole == "" {
			targetRole = "mid" // Default for most meta champions
		}

		analysis, err := h.service.AnalyzeCounterPicks(
			c.Request.Context(),
			champion,
			targetRole,
			gameMode,
			nil,
		)
		if err != nil {
			continue
		}

		// Get top 3 counters for each meta champion
		topCounters := analysis.CounterPicks
		if len(topCounters) > 3 {
			topCounters = topCounters[:3]
		}

		metaCounters = append(metaCounters, gin.H{
			"targetChampion": champion,
			"targetRole":     targetRole,
			"topCounters":    topCounters,
			"metaContext":    analysis.MetaContext,
		})

		if len(metaCounters) >= limit {
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"metaCounters": metaCounters,
		"gameMode":     gameMode,
		"rank":         rank,
		"region":       region,
	})
}

// GetPersonalizedCounters gets personalized counter recommendations
func (h *CounterPickHandler) GetPersonalizedCounters(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Summoner ID is required"})
		return
	}

	// Get query parameters
	targetChampion := c.Query("targetChampion")
	targetRole := c.Query("targetRole")
	gameMode := c.DefaultQuery("gameMode", "ranked")

	if targetChampion == "" || targetRole == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target champion and role are required"})
		return
	}

	// TODO: Get player's champion pool from database
	playerChampionPool := []string{} // Would be populated from player data

	analysis, err := h.service.AnalyzeCounterPicks(
		c.Request.Context(),
		targetChampion,
		targetRole,
		gameMode,
		playerChampionPool,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summonerID":           summonerID,
		"targetChampion":       targetChampion,
		"targetRole":           targetRole,
		"personalizedCounters": analysis.CounterPicks,
		"recommendations": gin.H{
			"comfort": "Based on your most played champions",
			"meta":    "Current patch recommendations",
			"winrate": "Highest success rate matchups",
		},
		"confidence": analysis.Confidence,
	})
}

// AddToFavorites adds a counter pick to user favorites
func (h *CounterPickHandler) AddToFavorites(c *gin.Context) {
	var request struct {
		SummonerID     string  `json:"summonerId" binding:"required"`
		Champion       string  `json:"champion" binding:"required"`
		TargetChampion string  `json:"targetChampion" binding:"required"`
		Role           string  `json:"role" binding:"required"`
		Notes          string  `json:"notes"`
		PersonalRating float64 `json:"personalRating"` // 1-10
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Save to database
	// This would create a CounterPickFavorites record

	c.JSON(http.StatusOK, gin.H{
		"message": "Counter pick added to favorites",
		"data": gin.H{
			"summonerId":     request.SummonerID,
			"champion":       request.Champion,
			"targetChampion": request.TargetChampion,
			"role":           request.Role,
			"notes":          request.Notes,
			"personalRating": request.PersonalRating,
		},
	})
}

// GetFavoriteCounters gets user's favorite counter picks
func (h *CounterPickHandler) GetFavoriteCounters(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Summoner ID is required"})
		return
	}

	// TODO: Query from database
	// This would fetch CounterPickFavorites records

	c.JSON(http.StatusOK, gin.H{
		"summonerID": summonerID,
		"favorites": []gin.H{
			{
				"id":             1,
				"champion":       "Malphite",
				"targetChampion": "Yasuo",
				"role":           "top",
				"notes":          "Rock solid counter, easy lane",
				"personalRating": 9.0,
				"timesUsed":      5,
				"successRate":    80.0,
			},
		},
	})
}

// GetCounterBanStrategy gets ban strategy based on counter analysis
func (h *CounterPickHandler) GetCounterBanStrategy(c *gin.Context) {
	var request struct {
		ThreateningChampions []string `json:"threateningChampions" binding:"required"`
		PlayerChampionPool   []string `json:"playerChampionPool"`
		BanPhase             string   `json:"banPhase"` // first, second
		GameMode             string   `json:"gameMode"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.GameMode == "" {
		request.GameMode = "ranked"
	}
	if request.BanPhase == "" {
		request.BanPhase = "first"
	}

	// Convert threatening champions to target data
	var targetChampions []services.TargetChampionData
	for _, champion := range request.ThreateningChampions {
		targetChampions = append(targetChampions, services.TargetChampionData{
			Champion:    champion,
			Role:        "unknown", // Would need role detection
			ThreatLevel: "high",
			Priority:    80,
			Reasons:     []string{"Threatening to team composition"},
		})
	}

	analysis, err := h.service.AnalyzeMultiTargetCounters(
		c.Request.Context(),
		targetChampions,
		request.GameMode,
		request.PlayerChampionPool,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"banStrategy":       analysis.BanRecommendations,
		"universalCounters": analysis.UniversalCounters,
		"specificCounters":  analysis.SpecificCounters,
		"recommendation":    analysis.OverallStrategy,
		"confidence":        analysis.Confidence,
	})
}

// Helper function to get current user ID (mock implementation)
func getCurrentUserID(c *gin.Context) uint {
	// This would extract user ID from JWT token or session
	return 1 // Mock user ID
}

// Additional utility endpoints
func (h *CounterPickHandler) RemoveFromFavorites(c *gin.Context) {
	id := c.Param("id")
	// TODO: Delete from database
	c.JSON(http.StatusOK, gin.H{"message": "Favorite removed", "id": id})
}

func (h *CounterPickHandler) GetCounterHistory(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	// TODO: Query counter pick history from database
	c.JSON(http.StatusOK, gin.H{
		"summonerID": summonerID,
		"history":    []gin.H{},
	})
}

func (h *CounterPickHandler) SubmitCounterFeedback(c *gin.Context) {
	var request struct {
		Champion       string  `json:"champion" binding:"required"`
		TargetChampion string  `json:"targetChampion" binding:"required"`
		MatchID        string  `json:"matchId" binding:"required"`
		Result         string  `json:"result" binding:"required"` // win/loss
		Performance    float64 `json:"performance"`               // 1-10
		Feedback       string  `json:"feedback"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Store feedback for model improvement
	c.JSON(http.StatusOK, gin.H{"message": "Feedback submitted successfully"})
}

func (h *CounterPickHandler) GetCounterPerformance(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	champion := c.Param("champion")
	// TODO: Calculate performance metrics
	c.JSON(http.StatusOK, gin.H{
		"summonerID": summonerID,
		"champion":   champion,
		"performance": gin.H{
			"overallWinRate": 65.5,
			"counterWinRate": 58.2,
			"averageKDA":     2.1,
			"gamesPlayed":    23,
		},
	})
}

func (h *CounterPickHandler) GetWinRateData(c *gin.Context) {
	champion := c.Param("champion")
	target := c.Param("target")
	// TODO: Get statistical win rate data
	c.JSON(http.StatusOK, gin.H{
		"champion":         champion,
		"target":           target,
		"winRate":          58.3,
		"sampleSize":       1247,
		"confidence":       85.2,
		"laneWinRate":      62.1,
		"teamFightWinRate": 55.7,
	})
}

func (h *CounterPickHandler) GetCounterPopularity(c *gin.Context) {
	champion := c.Param("champion")
	// TODO: Get popularity and usage statistics
	c.JSON(http.StatusOK, gin.H{
		"champion":    champion,
		"pickRate":    23.4,
		"banRate":     15.2,
		"counterRate": 8.7, // How often it's picked as a counter
	})
}

func (h *CounterPickHandler) AssessThreatLevel(c *gin.Context) {
	enemyChampions := c.QueryArray("champions")
	// TODO: Assess threat level of enemy team
	c.JSON(http.StatusOK, gin.H{
		"champions": enemyChampions,
		"threatAssessment": gin.H{
			"overall":    "high",
			"early":      "medium",
			"mid":        "high",
			"late":       "critical",
			"priorities": []string{"Ban Yasuo", "Counter Zed", "Focus Jinx"},
		},
	})
}
