package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/herald-lol/herald/backend/internal/services"
)

type RiotHandler struct {
	riotService *services.RiotService
}

func NewRiotHandler(riotService *services.RiotService) *RiotHandler {
	return &RiotHandler{
		riotService: riotService,
	}
}

// LinkAccount links a Riot account to the current user
// @Summary Link Riot account
// @Description Link a Riot Games account to the current user
// @Tags riot
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body LinkAccountRequest true "Account linking details"
// @Success 201 {object} models.RiotAccount
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /riot/link [post]
func (h *RiotHandler) LinkAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in token",
		})
		return
	}

	var req struct {
		Region   string `json:"region" binding:"required"`
		GameName string `json:"game_name" binding:"required"`
		TagLine  string `json:"tag_line" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	// Validate region
	validRegions := []string{"na1", "euw1", "eun1", "kr", "jp1", "br1", "la1", "la2", "oc1", "tr1", "ru"}
	isValidRegion := false
	for _, region := range validRegions {
		if req.Region == region {
			isValidRegion = true
			break
		}
	}

	if !isValidRegion {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid region",
			Message: "Region must be one of: " + "na1, euw1, eun1, kr, jp1, br1, la1, la2, oc1, tr1, ru",
		})
		return
	}

	account, err := h.riotService.LinkRiotAccount(c.Request.Context(), userID.(uuid.UUID).String(), req.Region, req.GameName, req.TagLine)
	if err != nil {
		switch err {
		case services.ErrSummonerNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Summoner not found",
				Message: "No summoner found with that name and tag",
			})
		case services.ErrAPIKeyInvalid:
			c.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Error:   "Service unavailable",
				Message: "Riot API service is currently unavailable",
			})
		case services.ErrRateLimitExceeded:
			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Error:   "Rate limit exceeded",
				Message: "Too many requests, please try again later",
			})
		default:
			if err.Error() == "account is already linked" {
				c.JSON(http.StatusConflict, ErrorResponse{
					Error:   "Account already linked",
					Message: "This Riot account is already linked to another user",
				})
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Link failed",
					Message: "Failed to link Riot account",
				})
			}
		}
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetLinkedAccounts returns all linked Riot accounts for the current user
// @Summary Get linked accounts
// @Description Get all Riot accounts linked to the current user
// @Tags riot
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.RiotAccount
// @Failure 401 {object} ErrorResponse
// @Router /riot/accounts [get]
func (h *RiotHandler) GetLinkedAccounts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in token",
		})
		return
	}

	// This would be implemented in the RiotService
	// For now, return empty array
	c.JSON(http.StatusOK, []interface{}{})
}

// SyncMatches syncs recent matches for a Riot account
// @Summary Sync matches
// @Description Sync recent matches for the specified Riot account
// @Tags riot
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param account_id path string true "Riot Account ID"
// @Param request body SyncMatchesRequest true "Sync parameters"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /riot/accounts/{account_id}/sync [post]
func (h *RiotHandler) SyncMatches(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in token",
		})
		return
	}

	accountID := c.Param("account_id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Missing account ID",
			Message: "Riot account ID is required",
		})
		return
	}

	var req struct {
		Count int `json:"count"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		req.Count = 20 // Default count
	}

	// Validate count
	if req.Count <= 0 || req.Count > 100 {
		req.Count = 20
	}

	err := h.riotService.SyncMatchHistory(c.Request.Context(), userID.(uuid.UUID).String(), accountID, req.Count)
	if err != nil {
		switch err {
		case services.ErrRateLimitExceeded:
			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Error:   "Rate limit exceeded",
				Message: "Too many requests, please try again later",
			})
		case services.ErrAPIKeyInvalid:
			c.JSON(http.StatusServiceUnavailable, ErrorResponse{
				Error:   "Service unavailable",
				Message: "Riot API service is currently unavailable",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Sync failed",
				Message: "Failed to sync matches",
			})
		}
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Matches synced successfully",
	})
}

// GetSummonerInfo gets basic summoner information
// @Summary Get summoner info
// @Description Get summoner information by name and tag
// @Tags riot
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GetSummonerRequest true "Summoner lookup details"
// @Success 200 {object} services.Summoner
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /riot/summoner [post]
func (h *RiotHandler) GetSummonerInfo(c *gin.Context) {
	var req struct {
		Region   string `json:"region" binding:"required"`
		GameName string `json:"game_name" binding:"required"`
		TagLine  string `json:"tag_line" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	// Get account first
	account, err := h.riotService.GetAccountByRiotID(c.Request.Context(), req.Region, req.GameName, req.TagLine)
	if err != nil {
		switch err {
		case services.ErrSummonerNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Summoner not found",
				Message: "No summoner found with that name and tag",
			})
		case services.ErrRateLimitExceeded:
			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Error:   "Rate limit exceeded",
				Message: "Too many requests, please try again later",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Lookup failed",
				Message: "Failed to lookup summoner",
			})
		}
		return
	}

	// Get summoner info
	summoner, err := h.riotService.GetSummonerByPUUID(c.Request.Context(), req.Region, account.PUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Lookup failed",
			Message: "Failed to get summoner information",
		})
		return
	}

	c.JSON(http.StatusOK, summoner)
}

// GetRankedInfo gets ranked information for a summoner
// @Summary Get ranked info
// @Description Get ranked information by summoner ID
// @Tags riot
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GetRankedInfoRequest true "Ranked lookup details"
// @Success 200 {array} services.LeagueEntry
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /riot/ranked [post]
func (h *RiotHandler) GetRankedInfo(c *gin.Context) {
	var req struct {
		Region     string `json:"region" binding:"required"`
		SummonerID string `json:"summoner_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	entries, err := h.riotService.GetLeagueEntries(c.Request.Context(), req.Region, req.SummonerID)
	if err != nil {
		switch err {
		case services.ErrSummonerNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Summoner not found",
				Message: "No ranked information found for that summoner",
			})
		case services.ErrRateLimitExceeded:
			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Error:   "Rate limit exceeded",
				Message: "Too many requests, please try again later",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Lookup failed",
				Message: "Failed to get ranked information",
			})
		}
		return
	}

	c.JSON(http.StatusOK, entries)
}

// GetMatchHistory gets match history for a summoner
// @Summary Get match history
// @Description Get recent match history for a summoner
// @Tags riot
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param puuid query string true "Player PUUID"
// @Param region query string true "Region"
// @Param count query int false "Number of matches (default: 20, max: 100)"
// @Success 200 {object} services.MatchHistory
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /riot/matches [get]
func (h *RiotHandler) GetMatchHistory(c *gin.Context) {
	puuid := c.Query("puuid")
	region := c.Query("region")
	countStr := c.Query("count")

	if puuid == "" || region == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Missing parameters",
			Message: "PUUID and region are required",
		})
		return
	}

	count := 20 // Default
	if countStr != "" {
		if parsed, err := strconv.Atoi(countStr); err == nil && parsed > 0 && parsed <= 100 {
			count = parsed
		}
	}

	matchHistory, err := h.riotService.GetMatchHistory(c.Request.Context(), region, puuid, count)
	if err != nil {
		switch err {
		case services.ErrRateLimitExceeded:
			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Error:   "Rate limit exceeded",
				Message: "Too many requests, please try again later",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Lookup failed",
				Message: "Failed to get match history",
			})
		}
		return
	}

	c.JSON(http.StatusOK, matchHistory)
}

// GetMatchDetails gets detailed information about a specific match
// @Summary Get match details
// @Description Get detailed information about a specific match
// @Tags riot
// @Produce json
// @Security BearerAuth
// @Param match_id path string true "Match ID"
// @Param region query string true "Region"
// @Success 200 {object} services.MatchDetails
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /riot/matches/{match_id} [get]
func (h *RiotHandler) GetMatchDetails(c *gin.Context) {
	matchID := c.Param("match_id")
	region := c.Query("region")

	if matchID == "" || region == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Missing parameters",
			Message: "Match ID and region are required",
		})
		return
	}

	matchDetails, err := h.riotService.GetMatchDetails(c.Request.Context(), region, matchID)
	if err != nil {
		switch err {
		case services.ErrMatchNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Match not found",
				Message: "No match found with that ID",
			})
		case services.ErrRateLimitExceeded:
			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Error:   "Rate limit exceeded",
				Message: "Too many requests, please try again later",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Lookup failed",
				Message: "Failed to get match details",
			})
		}
		return
	}

	c.JSON(http.StatusOK, matchDetails)
}

// GetRateLimitStatus gets current rate limit status
// @Summary Get rate limit status
// @Description Get current rate limit status for different regions
// @Tags riot
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /riot/rate-limit [get]
func (h *RiotHandler) GetRateLimitStatus(c *gin.Context) {
	// This would show rate limit status per region
	// For now, return basic status
	status := map[string]interface{}{
		"status": "operational",
		"regions": map[string]interface{}{
			"na1":  map[string]interface{}{"available": true, "requests_remaining": "unknown"},
			"euw1": map[string]interface{}{"available": true, "requests_remaining": "unknown"},
			"kr":   map[string]interface{}{"available": true, "requests_remaining": "unknown"},
		},
	}

	c.JSON(http.StatusOK, status)
}
