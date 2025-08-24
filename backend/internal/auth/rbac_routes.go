package auth

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Herald.lol Gaming Analytics - RBAC Routes
// HTTP handlers for Role-Based Access Control operations

// RBACRoutes defines RBAC route handlers
type RBACRoutes struct {
	rbac      *GamingRBACManager
	validator *validator.Validate
}

// NewRBACRoutes creates new RBAC routes handler
func NewRBACRoutes(rbac *GamingRBACManager) *RBACRoutes {
	return &RBACRoutes{
		rbac:      rbac,
		validator: validator.New(),
	}
}

// RegisterRoutes registers all RBAC routes
func (r *RBACRoutes) RegisterRoutes(group *gin.RouterGroup) {
	// Role management routes
	roles := group.Group("/roles")
	{
		roles.GET("", r.ListRoles)
		roles.GET("/:id", r.GetRole)
		roles.POST("", r.CreateRole)
		roles.PUT("/:id", r.UpdateRole)
		roles.DELETE("/:id", r.DeleteRole)
		roles.GET("/:id/permissions", r.GetRolePermissions)
		roles.POST("/:id/permissions/:permissionId", r.AssignPermissionToRole)
		roles.DELETE("/:id/permissions/:permissionId", r.RemovePermissionFromRole)
	}

	// Permission management routes
	permissions := group.Group("/permissions")
	{
		permissions.GET("", r.ListPermissions)
		permissions.GET("/:id", r.GetPermission)
		permissions.POST("", r.CreatePermission)
		permissions.PUT("/:id", r.UpdatePermission)
		permissions.DELETE("/:id", r.DeletePermission)
	}

	// User role assignment routes
	users := group.Group("/users")
	{
		users.GET("/:id/roles", r.GetUserRoles)
		users.POST("/:id/roles/:roleId", r.AssignRoleToUser)
		users.DELETE("/:id/roles/:roleId", r.RemoveRoleFromUser)
		users.GET("/:id/permissions", r.GetUserPermissions)
		users.GET("/:id/teams/:teamId/roles", r.GetUserTeamRoles)
	}

	// Team role management routes
	teams := group.Group("/teams")
	{
		teams.GET("/:id/roles", r.GetTeamRoles)
		teams.GET("/:id/members", r.GetTeamMembers)
		teams.POST("/:id/members/:userId/roles/:roleId", r.AssignTeamRole)
		teams.DELETE("/:id/members/:userId/roles/:roleId", r.RemoveTeamRole)
	}

	// Audit and monitoring routes
	audit := group.Group("/audit")
	{
		audit.GET("/logs", r.GetAuditLogs)
		audit.GET("/roles/:id/history", r.GetRoleHistory)
		audit.GET("/users/:id/history", r.GetUserRoleHistory)
	}
}

// Role management handlers

// ListRoles handles GET /api/v1/rbac/roles
func (r *RBACRoutes) ListRoles(c *gin.Context) {
	filters := &RoleFilters{}

	// Parse query parameters
	if roleType := c.Query("type"); roleType != "" {
		filters.Type = &roleType
	}
	if category := c.Query("category"); category != "" {
		filters.Category = &category
	}
	if isActive := c.Query("active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			filters.IsActive = &active
		}
	}
	if isSystem := c.Query("system"); isSystem != "" {
		if system, err := strconv.ParseBool(isSystem); err == nil {
			filters.IsSystem = &system
		}
	}
	if level := c.Query("level"); level != "" {
		if lvl, err := strconv.Atoi(level); err == nil {
			filters.Level = &lvl
		}
	}

	roles, err := r.rbac.ListRoles(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":           "Failed to list gaming roles",
			"details":         err.Error(),
			"gaming_platform": "herald-lol",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles":           roles,
		"total":           len(roles),
		"gaming_platform": "herald-lol",
	})
}

// GetRole handles GET /api/v1/rbac/roles/:id
func (r *RBACRoutes) GetRole(c *gin.Context) {
	roleID := c.Param("id")

	role, err := r.rbac.GetRole(c.Request.Context(), roleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":           "Gaming role not found",
			"details":         err.Error(),
			"gaming_platform": "herald-lol",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"role":            role,
		"gaming_platform": "herald-lol",
	})
}

// CreateRoleRequest represents role creation request
type CreateRoleRequest struct {
	Name         string                 `json:"name" validate:"required,min=3,max=100"`
	DisplayName  string                 `json:"display_name" validate:"required,min=3,max=100"`
	Description  string                 `json:"description" validate:"required,max=500"`
	Type         string                 `json:"type" validate:"required,oneof=gaming team admin system"`
	Category     string                 `json:"category" validate:"required,max=100"`
	Level        int                    `json:"level" validate:"min=1,max=100"`
	ParentRoleID *string                `json:"parent_role_id,omitempty"`
	Context      map[string]interface{} `json:"gaming_context,omitempty"`
	Metadata     map[string]string      `json:"metadata,omitempty"`
}

// CreateRole handles POST /api/v1/rbac/roles
func (r *RBACRoutes) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":           "Invalid gaming role data",
			"details":         err.Error(),
			"gaming_platform": "herald-lol",
		})
		return
	}

	if err := r.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":           "Gaming role validation failed",
			"details":         err.Error(),
			"gaming_platform": "herald-lol",
		})
		return
	}

	// Get current user info from context
	userID := GetUserIDFromContext(c)

	role := &GamingRole{
		Name:          req.Name,
		DisplayName:   req.DisplayName,
		Description:   req.Description,
		Type:          RoleType(req.Type),
		Category:      req.Category,
		Level:         req.Level,
		ParentRoleID:  req.ParentRoleID,
		GamingContext: req.Context,
		Metadata:      req.Metadata,
		IsSystem:      false,
		IsActive:      true,
		CreatedBy:     userID,
	}

	if err := r.rbac.CreateRole(c.Request.Context(), role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":           "Failed to create gaming role",
			"details":         err.Error(),
			"gaming_platform": "herald-lol",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "Gaming role created successfully",
		"role":            role,
		"gaming_platform": "herald-lol",
	})
}

// GetUserPermissions handles GET /api/v1/rbac/users/:id/permissions
func (r *RBACRoutes) GetUserPermissions(c *gin.Context) {
	userID := c.Param("id")

	permissions, err := r.rbac.GetUserPermissions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":           "Failed to get user gaming permissions",
			"details":         err.Error(),
			"gaming_platform": "herald-lol",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":         userID,
		"permissions":     permissions,
		"total":           len(permissions),
		"gaming_platform": "herald-lol",
	})
}

// Placeholder handlers for other endpoints (shortened for brevity)
func (r *RBACRoutes) UpdateRole(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) DeleteRole(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) ListPermissions(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) GetPermission(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) CreatePermission(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) UpdatePermission(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) DeletePermission(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) GetRolePermissions(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) AssignPermissionToRole(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) RemovePermissionFromRole(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) GetUserRoles(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) AssignRoleToUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) RemoveRoleFromUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) GetUserTeamRoles(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) GetTeamRoles(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) GetTeamMembers(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) AssignTeamRole(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) RemoveTeamRole(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) GetAuditLogs(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) GetRoleHistory(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

func (r *RBACRoutes) GetUserRoleHistory(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented", "gaming_platform": "herald-lol"})
}

// Helper function to get user ID from context
func GetUserIDFromContext(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return "system"
	}

	if id, ok := userID.(string); ok {
		return id
	}

	return "system"
}
