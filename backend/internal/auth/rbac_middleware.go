package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Herald.lol Gaming Analytics - RBAC Middleware
// Role-Based Access Control middleware for gaming platform endpoints

// GamingRBACMiddleware provides RBAC enforcement for gaming platform
type GamingRBACMiddleware struct {
	rbacManager *GamingRBACManager
}

// NewGamingRBACMiddleware creates new RBAC middleware for gaming platform
func NewGamingRBACMiddleware(rbacManager *GamingRBACManager) *GamingRBACMiddleware {
	return &GamingRBACMiddleware{
		rbacManager: rbacManager,
	}
}

// RequireGamingPermission middleware that requires specific gaming permission
func (m *GamingRBACMiddleware) RequireGamingPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required for gaming permission",
				"permission": permission,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Check if user has required permission
		hasPermission, err := m.rbacManager.HasPermission(c.Request.Context(), user.ID, permission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check gaming permission",
				"permission": permission,
			})
			c.Abort()
			return
		}
		
		if !hasPermission {
			// Get user's available permissions for helpful error message
			userPermissions, _ := m.rbacManager.GetUserPermissions(c.Request.Context(), user.ID)
			
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient gaming permissions",
				"required_permission": permission,
				"available_permissions": userPermissions,
				"gaming_platform": "herald-lol",
				"upgrade_info": m.getUpgradeInfo(permission),
			})
			c.Abort()
			return
		}
		
		// Set permission context
		c.Set("gaming_permission_verified", true)
		c.Set("gaming_required_permission", permission)
		
		// Add permission headers
		c.Header("X-Gaming-Permission-Verified", "true")
		c.Header("X-Gaming-Required-Permission", permission)
		
		c.Next()
	}
}

// RequireGamingRole middleware that requires specific gaming role
func (m *GamingRBACMiddleware) RequireGamingRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required for gaming role",
				"role": roleName,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Get user roles
		userRoles, err := m.rbacManager.rbacStore.GetUserRoles(c.Request.Context(), user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get gaming user roles",
			})
			c.Abort()
			return
		}
		
		// Check if user has required role
		hasRole := false
		for _, userRole := range userRoles {
			if !userRole.IsActive {
				continue
			}
			
			// Get role details
			role, err := m.rbacManager.rbacStore.GetRole(c.Request.Context(), userRole.RoleID)
			if err != nil {
				continue
			}
			
			if role.Name == roleName {
				hasRole = true
				break
			}
		}
		
		if !hasRole {
			// Get user's available roles for helpful error message
			var availableRoles []string
			for _, userRole := range userRoles {
				if role, err := m.rbacManager.rbacStore.GetRole(c.Request.Context(), userRole.RoleID); err == nil {
					availableRoles = append(availableRoles, role.Name)
				}
			}
			
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient gaming role",
				"required_role": roleName,
				"available_roles": availableRoles,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Set role context
		c.Set("gaming_role_verified", true)
		c.Set("gaming_required_role", roleName)
		
		c.Next()
	}
}

// RequireGamingTeamRole middleware for team-specific role requirements
func (m *GamingRBACMiddleware) RequireGamingTeamRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required for gaming team role",
				"role": roleName,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Get team ID from URL parameter
		teamID := c.Param("teamId")
		if teamID == "" {
			teamID = c.GetHeader("X-Gaming-Team-ID")
		}
		
		if teamID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Team ID required for gaming team role verification",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Get user team roles
		teamRoles, err := m.rbacManager.rbacStore.GetUserTeamRoles(c.Request.Context(), user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get gaming user team roles",
			})
			c.Abort()
			return
		}
		
		// Check if user has required team role
		hasTeamRole := false
		for _, teamRole := range teamRoles {
			if !teamRole.IsActive || teamRole.TeamID != teamID {
				continue
			}
			
			// Get role details
			role, err := m.rbacManager.rbacStore.GetRole(c.Request.Context(), teamRole.RoleID)
			if err != nil {
				continue
			}
			
			if role.Name == roleName {
				hasTeamRole = true
				c.Set("gaming_team_role", teamRole)
				break
			}
		}
		
		if !hasTeamRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient gaming team role",
				"required_role": roleName,
				"team_id": teamID,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Set team context
		c.Set("gaming_team_role_verified", true)
		c.Set("gaming_team_id", teamID)
		c.Set("gaming_required_team_role", roleName)
		
		c.Next()
	}
}

// RequireAnyGamingPermission middleware that requires any of the specified permissions
func (m *GamingRBACMiddleware) RequireAnyGamingPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required for gaming permissions",
				"required_permissions": permissions,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Check if user has any of the required permissions
		hasAnyPermission := false
		var grantedPermission string
		
		for _, permission := range permissions {
			hasPermission, err := m.rbacManager.HasPermission(c.Request.Context(), user.ID, permission)
			if err != nil {
				continue // Skip permission check errors
			}
			
			if hasPermission {
				hasAnyPermission = true
				grantedPermission = permission
				break
			}
		}
		
		if !hasAnyPermission {
			userPermissions, _ := m.rbacManager.GetUserPermissions(c.Request.Context(), user.ID)
			
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient gaming permissions",
				"required_any_of": permissions,
				"available_permissions": userPermissions,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Set permission context
		c.Set("gaming_permission_verified", true)
		c.Set("gaming_granted_permission", grantedPermission)
		c.Set("gaming_required_permissions", permissions)
		
		c.Next()
	}
}

// RequireAllGamingPermissions middleware that requires all specified permissions
func (m *GamingRBACMiddleware) RequireAllGamingPermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required for gaming permissions",
				"required_permissions": permissions,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Check if user has all required permissions
		var missingPermissions []string
		var grantedPermissions []string
		
		for _, permission := range permissions {
			hasPermission, err := m.rbacManager.HasPermission(c.Request.Context(), user.ID, permission)
			if err != nil || !hasPermission {
				missingPermissions = append(missingPermissions, permission)
			} else {
				grantedPermissions = append(grantedPermissions, permission)
			}
		}
		
		if len(missingPermissions) > 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient gaming permissions",
				"missing_permissions": missingPermissions,
				"granted_permissions": grantedPermissions,
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Set permission context
		c.Set("gaming_permissions_verified", true)
		c.Set("gaming_granted_permissions", grantedPermissions)
		
		c.Next()
	}
}

// RequireSubscriptionTier middleware that requires minimum subscription tier
func (m *GamingRBACMiddleware) RequireSubscriptionTier(minTier string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required for gaming subscription check",
				"gaming_platform": "herald-lol",
			})
			c.Abort()
			return
		}
		
		// Get user subscription tier
		userTier := "free" // Default tier
		if user.GamingProfile != nil && user.GamingProfile.SubscriptionTier != "" {
			userTier = user.GamingProfile.SubscriptionTier
		}
		
		// Check if user meets minimum tier requirement
		if !m.meetsTierRequirement(userTier, minTier) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient gaming subscription tier",
				"current_tier": userTier,
				"required_tier": minTier,
				"gaming_platform": "herald-lol",
				"upgrade_url": "https://herald.lol/upgrade?tier=" + minTier,
			})
			c.Abort()
			return
		}
		
		// Set subscription context
		c.Set("gaming_subscription_verified", true)
		c.Set("gaming_subscription_tier", userTier)
		
		c.Next()
	}
}

// EnrichGamingContext middleware that adds RBAC context to requests
func (m *GamingRBACMiddleware) EnrichGamingContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := GetGamingUser(c)
		if !exists {
			c.Next()
			return
		}
		
		ctx := context.Background()
		
		// Get user permissions (async to avoid blocking)
		go func() {
			permissions, err := m.rbacManager.GetUserPermissions(ctx, user.ID)
			if err == nil {
				c.Set("gaming_user_permissions", permissions)
			}
		}()
		
		// Get user roles
		userRoles, err := m.rbacManager.rbacStore.GetUserRoles(ctx, user.ID)
		if err == nil {
			var roleNames []string
			for _, userRole := range userRoles {
				if role, err := m.rbacManager.rbacStore.GetRole(ctx, userRole.RoleID); err == nil {
					roleNames = append(roleNames, role.Name)
				}
			}
			c.Set("gaming_user_roles", roleNames)
		}
		
		// Get team roles
		teamRoles, err := m.rbacManager.rbacStore.GetUserTeamRoles(ctx, user.ID)
		if err == nil {
			c.Set("gaming_user_team_roles", teamRoles)
		}
		
		c.Next()
	}
}

// AdminOnly middleware that restricts access to gaming administrators
func (m *GamingRBACMiddleware) AdminOnly() gin.HandlerFunc {
	return m.RequireAnyGamingPermission(
		"admin:users:manage",
		"admin:roles:manage",
		"admin:gaming:manage",
		"admin:system:manage",
	)
}

// TeamManagerOnly middleware for team management operations
func (m *GamingRBACMiddleware) TeamManagerOnly() gin.HandlerFunc {
	return m.RequireAnyGamingPermission(
		"team:manage:settings",
		"team:manage:players",
		"admin:teams:manage",
	)
}

// PremiumOnly middleware for premium features
func (m *GamingRBACMiddleware) PremiumOnly() gin.HandlerFunc {
	return m.RequireSubscriptionTier("premium")
}

// ProOnly middleware for professional features
func (m *GamingRBACMiddleware) ProOnly() gin.HandlerFunc {
	return m.RequireSubscriptionTier("pro")
}

// EnterpriseOnly middleware for enterprise features
func (m *GamingRBACMiddleware) EnterpriseOnly() gin.HandlerFunc {
	return m.RequireSubscriptionTier("enterprise")
}

// Helper methods

// meetsTierRequirement checks if user tier meets minimum requirement
func (m *GamingRBACMiddleware) meetsTierRequirement(userTier, minTier string) bool {
	tierLevels := map[string]int{
		"free":       0,
		"premium":    1,
		"pro":        2,
		"enterprise": 3,
	}
	
	userLevel, userExists := tierLevels[userTier]
	minLevel, minExists := tierLevels[minTier]
	
	if !userExists || !minExists {
		return false
	}
	
	return userLevel >= minLevel
}

// getUpgradeInfo returns upgrade information for insufficient permissions
func (m *GamingRBACMiddleware) getUpgradeInfo(permission string) map[string]string {
	upgradeInfo := map[string]string{
		"upgrade_url": "https://herald.lol/upgrade",
		"contact":     "support@herald.lol",
	}
	
	// Map permissions to required tiers
	permissionTiers := map[string]string{
		"gaming:analytics:advanced": "premium",
		"gaming:analytics:export":   "pro",
		"api:extended":             "premium",
		"api:unlimited":            "enterprise",
		"team:manage:players":      "pro",
		"team:manage:settings":     "pro",
		"admin:users:manage":       "enterprise",
		"admin:roles:manage":       "enterprise",
	}
	
	if requiredTier, exists := permissionTiers[permission]; exists {
		upgradeInfo["required_tier"] = requiredTier
		upgradeInfo["upgrade_url"] = "https://herald.lol/upgrade?tier=" + requiredTier
	}
	
	return upgradeInfo
}

// RBAC Context Helpers

// GetGamingPermissions extracts gaming permissions from Gin context
func GetGamingPermissions(c *gin.Context) ([]string, bool) {
	if permissions, exists := c.Get("gaming_user_permissions"); exists {
		if permSlice, ok := permissions.([]string); ok {
			return permSlice, true
		}
	}
	return []string{}, false
}

// GetGamingRoles extracts gaming roles from Gin context
func GetGamingRoles(c *gin.Context) ([]string, bool) {
	if roles, exists := c.Get("gaming_user_roles"); exists {
		if roleSlice, ok := roles.([]string); ok {
			return roleSlice, true
		}
	}
	return []string{}, false
}

// GetGamingTeamRoles extracts gaming team roles from Gin context
func GetGamingTeamRoles(c *gin.Context) ([]*TeamRole, bool) {
	if teamRoles, exists := c.Get("gaming_user_team_roles"); exists {
		if teamRoleSlice, ok := teamRoles.([]*TeamRole); ok {
			return teamRoleSlice, true
		}
	}
	return []*TeamRole{}, false
}

// HasGamingPermissionInContext checks if permission was verified in context
func HasGamingPermissionInContext(c *gin.Context, permission string) bool {
	if verified, exists := c.Get("gaming_permission_verified"); exists {
		if verifiedBool, ok := verified.(bool); ok && verifiedBool {
			if requiredPerm, exists := c.Get("gaming_required_permission"); exists {
				if requiredPermStr, ok := requiredPerm.(string); ok {
					return requiredPermStr == permission
				}
			}
			
			// Check granted permissions for any/all permission middleware
			if grantedPerms, exists := c.Get("gaming_granted_permissions"); exists {
				if grantedSlice, ok := grantedPerms.([]string); ok {
					for _, granted := range grantedSlice {
						if granted == permission {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// GetVerifiedGamingTeamID gets the team ID that was verified in team role middleware
func GetVerifiedGamingTeamID(c *gin.Context) (string, bool) {
	if teamID, exists := c.Get("gaming_team_id"); exists {
		if teamIDStr, ok := teamID.(string); ok {
			return teamIDStr, true
		}
	}
	return "", false
}

// IsGamingAdmin checks if current user has admin permissions
func IsGamingAdmin(c *gin.Context) bool {
	permissions, exists := GetGamingPermissions(c)
	if !exists {
		return false
	}
	
	adminPermissions := []string{
		"admin:users:manage",
		"admin:roles:manage", 
		"admin:gaming:manage",
		"admin:system:manage",
	}
	
	for _, adminPerm := range adminPermissions {
		for _, userPerm := range permissions {
			if userPerm == adminPerm {
				return true
			}
		}
	}
	
	return false
}

// IsTeamManager checks if current user has team management permissions for specific team
func IsTeamManager(c *gin.Context, teamID string) bool {
	teamRoles, exists := GetGamingTeamRoles(c)
	if !exists {
		return false
	}
	
	managerRoles := []string{"team:manager", "team:captain"}
	
	for _, teamRole := range teamRoles {
		if teamRole.TeamID == teamID && teamRole.IsActive {
			// Get role name (this would require role lookup in real implementation)
			for _, managerRole := range managerRoles {
				if strings.Contains(teamRole.Position, "manager") || strings.Contains(teamRole.Position, "captain") {
					return true
				}
			}
		}
	}
	
	return false
}