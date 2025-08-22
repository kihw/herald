package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"lol-match-exporter/internal/models"
	"lol-match-exporter/internal/services"
)

// GroupHandler gère les routes relatives aux groupes
type GroupHandler struct {
	groupService      *services.GroupService
	comparisonService *services.ComparisonService
	db                *sql.DB
}

// NewGroupHandler crée un nouveau handler pour les groupes
func NewGroupHandler(db *sql.DB, analyticsService *services.AnalyticsService) *GroupHandler {
	return &GroupHandler{
		groupService:      services.NewGroupService(db),
		comparisonService: services.NewComparisonService(db, analyticsService),
		db:                db,
	}
}

// CreateGroupRequest définit la structure de création d'un groupe
type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"` // public, private, invite_only
}

// CreateGroup crée un nouveau groupe
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pour l'instant, on simule un utilisateur connecté avec l'ID 1
	// Dans une vraie app, on récupérerait l'ID depuis la session/JWT
	userID := 1

	group, err := h.groupService.CreateGroup(userID, req.Name, req.Description, req.Privacy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"group":   group,
		"message": "Group created successfully",
	})
}

// GetGroup récupère un groupe par son ID
func (h *GroupHandler) GetGroup(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	group, err := h.groupService.GetGroupByID(groupID)
	if err != nil {
		if err.Error() == "group not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"group":   group,
	})
}

// GetUserGroups récupère tous les groupes d'un utilisateur
func (h *GroupHandler) GetUserGroups(c *gin.Context) {
	userID := 1 // Simulé pour l'instant

	groups, err := h.groupService.GetGroupsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"groups":  groups,
		"count":   len(groups),
	})
}

// SearchGroupsRequest définit la structure de recherche de groupes
type SearchGroupsRequest struct {
	Query string `form:"q" binding:"required"`
	Limit int    `form:"limit"`
}

// SearchGroups recherche des groupes publics
func (h *GroupHandler) SearchGroups(c *gin.Context) {
	var req SearchGroupsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}

	groups, err := h.groupService.SearchPublicGroups(req.Query, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"groups":  groups,
		"count":   len(groups),
		"query":   req.Query,
	})
}

// JoinGroupRequest définit la structure pour rejoindre un groupe
type JoinGroupRequest struct {
	InviteCode string `json:"invite_code" binding:"required"`
}

// JoinGroup permet de rejoindre un groupe via un code d'invitation
func (h *GroupHandler) JoinGroup(c *gin.Context) {
	var req JoinGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := 1 // Simulé

	err := h.groupService.JoinGroupByInviteCode(req.InviteCode, userID)
	if err != nil {
		if err.Error() == "invalid invite code" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid invite code"})
			return
		}
		if err.Error() == "user is already a member of this group" {
			c.JSON(http.StatusConflict, gin.H{"error": "You are already a member of this group"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully joined the group",
	})
}

// InviteToGroupRequest définit la structure pour inviter à un groupe
type InviteToGroupRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Message string `json:"message"`
}

// InviteToGroup invite un utilisateur à rejoindre un groupe
func (h *GroupHandler) InviteToGroup(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var req InviteToGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inviterID := 1 // Simulé

	invite, err := h.groupService.CreateGroupInvite(groupID, inviterID, req.Email, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"invite":  invite,
		"message": "Invitation sent successfully",
	})
}

// GetGroupMembers récupère les membres d'un groupe
func (h *GroupHandler) GetGroupMembers(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	members, err := h.groupService.GetGroupMembers(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"members": members,
		"count":   len(members),
	})
}

// RemoveMemberRequest définit la structure pour retirer un membre
type RemoveMemberRequest struct {
	UserID int `json:"user_id" binding:"required"`
}

// RemoveMember retire un membre d'un groupe
func (h *GroupHandler) RemoveMember(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var req RemoveMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.groupService.RemoveMemberFromGroup(groupID, req.UserID)
	if err != nil {
		if err.Error() == "cannot remove group owner" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot remove group owner"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Member removed successfully",
	})
}

// CreateComparisonRequest définit la structure pour créer une comparaison
type CreateComparisonRequest struct {
	Name        string                          `json:"name" binding:"required"`
	Description string                          `json:"description"`
	CompareType string                          `json:"compare_type" binding:"required"`
	Parameters  models.ComparisonParameters     `json:"parameters" binding:"required"`
}

// CreateComparison crée une nouvelle comparaison dans un groupe
func (h *GroupHandler) CreateComparison(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var req CreateComparisonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	creatorID := 1 // Simulé

	comparison, err := h.comparisonService.CreateComparison(
		groupID, creatorID, req.Name, req.Description, req.CompareType, req.Parameters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":    true,
		"comparison": comparison,
		"message":    "Comparison created successfully",
	})
}

// GetGroupComparisons récupère les comparaisons d'un groupe
func (h *GroupHandler) GetGroupComparisons(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	comparisons, err := h.comparisonService.GetGroupComparisons(groupID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"comparisons": comparisons,
		"count":       len(comparisons),
	})
}

// GetComparison récupère une comparaison spécifique
func (h *GroupHandler) GetComparison(c *gin.Context) {
	comparisonID, err := strconv.Atoi(c.Param("comparisonId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comparison ID"})
		return
	}

	comparison, err := h.comparisonService.GetComparisonByID(comparisonID)
	if err != nil {
		if err.Error() == "comparison not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comparison not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"comparison": comparison,
	})
}

// RegenerateComparison régénère les résultats d'une comparaison
func (h *GroupHandler) RegenerateComparison(c *gin.Context) {
	comparisonID, err := strconv.Atoi(c.Param("comparisonId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comparison ID"})
		return
	}

	err = h.comparisonService.GenerateComparisonResults(comparisonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Récupérer la comparaison mise à jour
	comparison, err := h.comparisonService.GetComparisonByID(comparisonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"comparison": comparison,
		"message":    "Comparison results regenerated successfully",
	})
}

// GetGroupStats récupère les statistiques d'un groupe
func (h *GroupHandler) GetGroupStats(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Pour l'instant, retourner des statistiques simulées
	// Dans une vraie app, on calculerait les vraies stats
	stats := &models.GroupStats{
		ID:              1,
		GroupID:         groupID,
		TotalMembers:    5,
		ActiveMembers:   4,
		AverageRank:     "Gold II",
		AverageMMR:      nil,
		TopChampions:    []models.ChampionStat{
			{ChampionID: 1, ChampionName: "Jinx", PlayCount: 45, WinRate: 64.4, AvgKDA: 2.3},
			{ChampionID: 2, ChampionName: "Yasuo", PlayCount: 38, WinRate: 52.6, AvgKDA: 1.8},
			{ChampionID: 3, ChampionName: "Thresh", PlayCount: 42, WinRate: 58.9, AvgKDA: 1.9},
		},
		PopularRoles:    []models.RoleStat{
			{Role: "ADC", PlayCount: 78, WinRate: 61.2},
			{Role: "MID", PlayCount: 65, WinRate: 56.8},
			{Role: "SUPPORT", PlayCount: 58, WinRate: 59.3},
		},
		WinRateComparison: map[string]float64{
			"current_month": 58.5,
			"last_month":    54.2,
			"season":        56.8,
		},
		LastUpdated: time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"stats":   stats,
	})
}

// UpdateGroupSettingsRequest définit la structure pour mettre à jour les paramètres
type UpdateGroupSettingsRequest struct {
	Settings models.GroupSettings `json:"settings" binding:"required"`
}

// UpdateGroupSettings met à jour les paramètres d'un groupe
func (h *GroupHandler) UpdateGroupSettings(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var req UpdateGroupSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.groupService.UpdateGroupSettings(groupID, req.Settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Group settings updated successfully",
	})
}