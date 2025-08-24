// Team Composition Handler for Herald.lol
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/herald-lol/herald/backend/internal/services"
)

type TeamCompositionHandler struct {
	service *services.TeamCompositionService
}

func NewTeamCompositionHandler(service *services.TeamCompositionService) *TeamCompositionHandler {
	return &TeamCompositionHandler{
		service: service,
	}
}

func (h *TeamCompositionHandler) RegisterRoutes(rg *gin.RouterGroup) {
	composition := rg.Group("/team-composition")
	{
		composition.POST("/optimize", h.OptimizeTeamComposition)
		composition.POST("/analyze", h.AnalyzeComposition)
		composition.GET("/suggestions/:summoner_id", h.GetCompositionSuggestions)
		composition.POST("/validate", h.ValidateComposition)
		composition.POST("/compare", h.CompareCompositions)
		composition.GET("/meta-compositions", h.GetMetaCompositions)
		composition.POST("/synergy-analysis", h.AnalyzeTeamSynergy)
		composition.POST("/counter-analysis", h.AnalyzeCounters)
		composition.GET("/role-recommendations/:summoner_id/:role", h.GetRoleRecommendations)
		composition.POST("/draft-optimization", h.OptimizeDraftPicks)
		composition.GET("/player-comfort/:summoner_id", h.GetPlayerComfortPicks)
		composition.POST("/win-condition-analysis", h.AnalyzeWinConditions)
		composition.POST("/scaling-analysis", h.AnalyzeScaling)
		composition.GET("/champion-pools/:summoner_id", h.GetChampionPools)
		composition.POST("/ban-strategy", h.GetBanStrategy)
	}
}

// OptimizeTeamComposition optimizes team composition based on multiple criteria
func (h *TeamCompositionHandler) OptimizeTeamComposition(c *gin.Context) {
	var request struct {
		PlayerData []struct {
			SummonerID   string   `json:"summonerId" binding:"required"`
			Role         string   `json:"role" binding:"required"`
			ChampionPool []string `json:"championPool"`
			ComfortLevel int      `json:"comfortLevel"` // 1-10
			RecentGames  int      `json:"recentGames"`
		} `json:"playerData" binding:"required"`
		Strategy          string   `json:"strategy"` // meta_optimal, synergy_focused, balanced, comfort_picks
		BannedChampions   []string `json:"bannedChampions"`
		RequiredChampions []string `json:"requiredChampions"`
		GameMode          string   `json:"gameMode"`
		Constraints       struct {
			MaxNewChampions int  `json:"maxNewChampions"`
			RequireADC      bool `json:"requireADC"`
			RequireTank     bool `json:"requireTank"`
			PreferLateGame  bool `json:"preferLateGame"`
			PreferEarlyGame bool `json:"preferEarlyGame"`
		} `json:"constraints"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate strategy
	validStrategies := []string{"meta_optimal", "synergy_focused", "balanced", "comfort_picks"}
	strategyValid := false
	for _, validStrategy := range validStrategies {
		if request.Strategy == validStrategy {
			strategyValid = true
			break
		}
	}
	if !strategyValid {
		request.Strategy = "balanced" // Default strategy
	}

	result, err := h.service.OptimizeComposition(request.PlayerData, request.Strategy, request.BannedChampions, request.RequiredChampions, request.GameMode, request.Constraints)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AnalyzeComposition analyzes an existing team composition
func (h *TeamCompositionHandler) AnalyzeComposition(c *gin.Context) {
	var request struct {
		BlueTeam []struct {
			Champion string `json:"champion" binding:"required"`
			Role     string `json:"role" binding:"required"`
		} `json:"blueTeam" binding:"required"`
		RedTeam []struct {
			Champion string `json:"champion" binding:"required"`
			Role     string `json:"role" binding:"required"`
		} `json:"redTeam"`
		GameMode string `json:"gameMode"`
		Patch    string `json:"patch"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	analysis, err := h.service.AnalyzeComposition(request.BlueTeam, request.RedTeam, request.GameMode, request.Patch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetCompositionSuggestions gets composition suggestions for a player
func (h *TeamCompositionHandler) GetCompositionSuggestions(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Summoner ID is required"})
		return
	}

	// Get query parameters
	role := c.Query("role")
	gameMode := c.DefaultQuery("gameMode", "ranked")
	strategy := c.DefaultQuery("strategy", "balanced")
	limitStr := c.DefaultQuery("limit", "10")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	suggestions, err := h.service.GetCompositionSuggestions(summonerID, role, gameMode, strategy, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, suggestions)
}

// ValidateComposition validates a team composition
func (h *TeamCompositionHandler) ValidateComposition(c *gin.Context) {
	var request struct {
		Composition []struct {
			Champion string `json:"champion" binding:"required"`
			Role     string `json:"role" binding:"required"`
		} `json:"composition" binding:"required"`
		GameMode string `json:"gameMode"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validation, err := h.service.ValidateComposition(request.Composition, request.GameMode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, validation)
}

// CompareCompositions compares multiple team compositions
func (h *TeamCompositionHandler) CompareCompositions(c *gin.Context) {
	var request struct {
		Compositions []struct {
			Name      string `json:"name"`
			Champions []struct {
				Champion string `json:"champion" binding:"required"`
				Role     string `json:"role" binding:"required"`
			} `json:"champions" binding:"required"`
		} `json:"compositions" binding:"required"`
		GameMode string   `json:"gameMode"`
		Criteria []string `json:"criteria"` // synergy, scaling, teamfight, etc.
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comparison, err := h.service.CompareCompositions(request.Compositions, request.GameMode, request.Criteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comparison)
}

// GetMetaCompositions gets current meta compositions
func (h *TeamCompositionHandler) GetMetaCompositions(c *gin.Context) {
	// Get query parameters
	gameMode := c.DefaultQuery("gameMode", "ranked")
	rank := c.DefaultQuery("rank", "all")
	region := c.DefaultQuery("region", "global")
	patch := c.Query("patch")
	limitStr := c.DefaultQuery("limit", "20")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	compositions, err := h.service.GetMetaCompositions(gameMode, rank, region, patch, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, compositions)
}

// AnalyzeTeamSynergy analyzes synergy between champions
func (h *TeamCompositionHandler) AnalyzeTeamSynergy(c *gin.Context) {
	var request struct {
		Champions []struct {
			Champion string `json:"champion" binding:"required"`
			Role     string `json:"role" binding:"required"`
		} `json:"champions" binding:"required"`
		SynergyType string `json:"synergyType"` // teamfight, engage, protect, poke, etc.
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	synergy, err := h.service.AnalyzeTeamSynergy(request.Champions, request.SynergyType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, synergy)
}

// AnalyzeCounters analyzes counter-pick opportunities
func (h *TeamCompositionHandler) AnalyzeCounters(c *gin.Context) {
	var request struct {
		EnemyComposition []struct {
			Champion string `json:"champion" binding:"required"`
			Role     string `json:"role" binding:"required"`
		} `json:"enemyComposition" binding:"required"`
		AvailableChampions []string `json:"availableChampions"`
		TargetRole         string   `json:"targetRole"`
		CounterType        string   `json:"counterType"` // lane, teamfight, splitpush, etc.
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	counters, err := h.service.AnalyzeCounters(request.EnemyComposition, request.AvailableChampions, request.TargetRole, request.CounterType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, counters)
}

// GetRoleRecommendations gets champion recommendations for a specific role
func (h *TeamCompositionHandler) GetRoleRecommendations(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	role := c.Param("role")

	if summonerID == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Summoner ID and role are required"})
		return
	}

	// Get query parameters
	existingTeam := c.QueryArray("existing_champions")
	gameMode := c.DefaultQuery("gameMode", "ranked")
	strategy := c.DefaultQuery("strategy", "balanced")
	limitStr := c.DefaultQuery("limit", "15")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 15
	}

	recommendations, err := h.service.GetRoleRecommendations(summonerID, role, existingTeam, gameMode, strategy, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

// OptimizeDraftPicks optimizes picks during draft phase
func (h *TeamCompositionHandler) OptimizeDraftPicks(c *gin.Context) {
	var request struct {
		DraftState struct {
			BluePicks   []string `json:"bluePicks"`
			RedPicks    []string `json:"redPicks"`
			BlueBans    []string `json:"blueBans"`
			RedBans     []string `json:"redBans"`
			CurrentTurn string   `json:"currentTurn"` // blue_pick, red_pick, blue_ban, red_ban
			PickOrder   []string `json:"pickOrder"`   // remaining picks
		} `json:"draftState" binding:"required"`
		PlayerData []struct {
			SummonerID   string   `json:"summonerId" binding:"required"`
			Role         string   `json:"role" binding:"required"`
			ChampionPool []string `json:"championPool"`
		} `json:"playerData" binding:"required"`
		GameMode string `json:"gameMode"`
		Strategy string `json:"strategy"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	optimization, err := h.service.OptimizeDraftPicks(request.DraftState, request.PlayerData, request.GameMode, request.Strategy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, optimization)
}

// GetPlayerComfortPicks gets comfort picks for a player
func (h *TeamCompositionHandler) GetPlayerComfortPicks(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Summoner ID is required"})
		return
	}

	// Get query parameters
	role := c.Query("role")
	gameMode := c.DefaultQuery("gameMode", "ranked")
	recentGamesStr := c.DefaultQuery("recentGames", "50")
	limitStr := c.DefaultQuery("limit", "20")

	recentGames, err := strconv.Atoi(recentGamesStr)
	if err != nil {
		recentGames = 50
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	comfortPicks, err := h.service.GetPlayerComfortPicks(summonerID, role, gameMode, recentGames, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comfortPicks)
}

// AnalyzeWinConditions analyzes win conditions for a team composition
func (h *TeamCompositionHandler) AnalyzeWinConditions(c *gin.Context) {
	var request struct {
		TeamComposition []struct {
			Champion string `json:"champion" binding:"required"`
			Role     string `json:"role" binding:"required"`
		} `json:"teamComposition" binding:"required"`
		EnemyComposition []struct {
			Champion string `json:"champion" binding:"required"`
			Role     string `json:"role" binding:"required"`
		} `json:"enemyComposition"`
		GameMode string `json:"gameMode"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	winConditions, err := h.service.AnalyzeWinConditions(request.TeamComposition, request.EnemyComposition, request.GameMode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, winConditions)
}

// AnalyzeScaling analyzes scaling patterns for team composition
func (h *TeamCompositionHandler) AnalyzeScaling(c *gin.Context) {
	var request struct {
		TeamComposition []struct {
			Champion string `json:"champion" binding:"required"`
			Role     string `json:"role" binding:"required"`
		} `json:"teamComposition" binding:"required"`
		CompareAgainst []struct {
			Champion string `json:"champion" binding:"required"`
			Role     string `json:"role" binding:"required"`
		} `json:"compareAgainst"`
		GameMode string `json:"gameMode"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scaling, err := h.service.AnalyzeScaling(request.TeamComposition, request.CompareAgainst, request.GameMode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scaling)
}

// GetChampionPools gets champion pools for team members
func (h *TeamCompositionHandler) GetChampionPools(c *gin.Context) {
	summonerID := c.Param("summoner_id")
	if summonerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Summoner ID is required"})
		return
	}

	// Get query parameters
	role := c.Query("role")
	gameMode := c.DefaultQuery("gameMode", "ranked")
	poolType := c.DefaultQuery("poolType", "comfort") // comfort, meta, flex, wide
	recentGamesStr := c.DefaultQuery("recentGames", "100")

	recentGames, err := strconv.Atoi(recentGamesStr)
	if err != nil {
		recentGames = 100
	}

	championPools, err := h.service.GetChampionPools(summonerID, role, gameMode, poolType, recentGames)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, championPools)
}

// GetBanStrategy gets optimal ban strategy
func (h *TeamCompositionHandler) GetBanStrategy(c *gin.Context) {
	var request struct {
		PlayerData []struct {
			SummonerID string `json:"summonerId" binding:"required"`
			Role       string `json:"role" binding:"required"`
		} `json:"playerData" binding:"required"`
		EnemyData []struct {
			SummonerID string `json:"summonerId"`
			Role       string `json:"role" binding:"required"`
		} `json:"enemyData"`
		BanPhase     string   `json:"banPhase"` // first_ban, second_ban
		ExistingBans []string `json:"existingBans"`
		GameMode     string   `json:"gameMode"`
		Strategy     string   `json:"strategy"` // target_player, protect_comp, meta_deny
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	banStrategy, err := h.service.GetBanStrategy(request.PlayerData, request.EnemyData, request.BanPhase, request.ExistingBans, request.GameMode, request.Strategy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, banStrategy)
}
