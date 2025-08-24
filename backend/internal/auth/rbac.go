package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Herald.lol Gaming Analytics - Role-Based Access Control
// Comprehensive RBAC system for gaming platform with hierarchical roles and fine-grained permissions

// GamingRBACManager manages role-based access control for Herald.lol
type GamingRBACManager struct {
	db              *gorm.DB
	rbacStore       RBACStore
	gamingAnalytics GamingAnalyticsService
	config          *GamingRBACConfig
	roleHierarchy   *RoleHierarchy
	permissionCache PermissionCache
}

// GamingRBACConfig holds RBAC configuration for gaming platform
type GamingRBACConfig struct {
	// Role configuration
	DefaultRole         string
	AdminRole           string
	EnableRoleHierarchy bool
	EnableInheritance   bool

	// Permission configuration
	EnablePermissionCache bool
	CacheTTL              time.Duration
	MaxRolesPerUser       int
	MaxPermissionsPerRole int

	// Gaming-specific settings
	GamingRolePrefix     string
	TeamRoleEnabled      bool
	DynamicRolesEnabled  bool
	SubscriptionRoleSync bool

	// Audit settings
	EnableAuditLog     bool
	AuditRetentionDays int
}

// RBACStore interface for RBAC data management
type RBACStore interface {
	// Role methods
	CreateRole(ctx context.Context, role *GamingRole) error
	GetRole(ctx context.Context, roleID string) (*GamingRole, error)
	GetRoleByName(ctx context.Context, name string) (*GamingRole, error)
	UpdateRole(ctx context.Context, role *GamingRole) error
	DeleteRole(ctx context.Context, roleID string) error
	ListRoles(ctx context.Context, filters *RoleFilters) ([]*GamingRole, error)

	// Permission methods
	CreatePermission(ctx context.Context, permission *GamingPermission) error
	GetPermission(ctx context.Context, permissionID string) (*GamingPermission, error)
	GetPermissionByName(ctx context.Context, name string) (*GamingPermission, error)
	ListPermissions(ctx context.Context, filters *PermissionFilters) ([]*GamingPermission, error)

	// Role-Permission mapping
	AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]*GamingPermission, error)

	// User-Role mapping
	AssignRoleToUser(ctx context.Context, userID, roleID string, assignment *RoleAssignment) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]*UserRole, error)
	GetUsersWithRole(ctx context.Context, roleID string) ([]*UserRole, error)

	// Team-Role mapping (gaming-specific)
	AssignRoleToTeam(ctx context.Context, teamID, roleID string, assignment *TeamRoleAssignment) error
	RemoveRoleFromTeam(ctx context.Context, teamID, roleID string) error
	GetTeamRoles(ctx context.Context, teamID string) ([]*TeamRole, error)
	GetUserTeamRoles(ctx context.Context, userID string) ([]*TeamRole, error)

	// Audit methods
	LogRoleAction(ctx context.Context, action *RoleAuditLog) error
	GetRoleAuditLog(ctx context.Context, filters *AuditFilters) ([]*RoleAuditLog, error)
}

// PermissionCache interface for caching user permissions
type PermissionCache interface {
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	SetUserPermissions(ctx context.Context, userID string, permissions []string, ttl time.Duration) error
	InvalidateUserPermissions(ctx context.Context, userID string) error
	InvalidateRolePermissions(ctx context.Context, roleID string) error
}

// Core RBAC data structures

// GamingRole represents a role in the gaming platform
type GamingRole struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	DisplayName   string                 `json:"display_name"`
	Description   string                 `json:"description"`
	Type          RoleType               `json:"type"`
	Category      string                 `json:"category"` // gaming, admin, team, subscription
	Level         int                    `json:"level"`    // Hierarchy level
	ParentRoleID  *string                `json:"parent_role_id,omitempty"`
	IsSystem      bool                   `json:"is_system"`
	IsActive      bool                   `json:"is_active"`
	GamingContext map[string]interface{} `json:"gaming_context,omitempty"`
	Metadata      map[string]string      `json:"metadata,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	CreatedBy     string                 `json:"created_by"`

	// Related data
	Permissions []*GamingPermission `json:"permissions,omitempty"`
	ChildRoles  []*GamingRole       `json:"child_roles,omitempty"`
	UserCount   int                 `json:"user_count,omitempty"`
}

// GamingPermission represents a permission in the gaming platform
type GamingPermission struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"` // e.g., gaming:analytics:view
	DisplayName      string                 `json:"display_name"`
	Description      string                 `json:"description"`
	Category         string                 `json:"category"` // gaming, analytics, team, admin
	Resource         string                 `json:"resource"` // analytics, teams, users, etc.
	Action           string                 `json:"action"`   // view, create, update, delete, manage
	Scope            PermissionScope        `json:"scope"`    // self, team, organization, global
	RequiresMFA      bool                   `json:"requires_mfa"`
	SubscriptionTier string                 `json:"subscription_tier,omitempty"` // free, premium, pro, enterprise
	GamingContext    map[string]interface{} `json:"gaming_context,omitempty"`
	IsSystem         bool                   `json:"is_system"`
	IsActive         bool                   `json:"is_active"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	CreatedBy        string                 `json:"created_by"`
}

// UserRole represents a role assignment to a user
type UserRole struct {
	ID         string                 `json:"id"`
	UserID     string                 `json:"user_id"`
	RoleID     string                 `json:"role_id"`
	AssignedBy string                 `json:"assigned_by"`
	AssignedAt time.Time              `json:"assigned_at"`
	ExpiresAt  *time.Time             `json:"expires_at,omitempty"`
	IsActive   bool                   `json:"is_active"`
	Scope      string                 `json:"scope,omitempty"`   // team:123, org:456, global
	Context    map[string]interface{} `json:"context,omitempty"` // Additional context

	// Related data
	Role *GamingRole     `json:"role,omitempty"`
	User *GamingUserInfo `json:"user,omitempty"`
}

// TeamRole represents gaming team-specific role assignments
type TeamRole struct {
	ID         string     `json:"id"`
	TeamID     string     `json:"team_id"`
	UserID     string     `json:"user_id"`
	RoleID     string     `json:"role_id"`
	Position   string     `json:"position,omitempty"` // captain, analyst, coach, player
	AssignedBy string     `json:"assigned_by"`
	AssignedAt time.Time  `json:"assigned_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	IsActive   bool       `json:"is_active"`

	// Gaming-specific fields
	GameRole string   `json:"game_role,omitempty"` // support, adc, mid, jungle, top
	Champion []string `json:"champion,omitempty"`  // Main champions
	Region   string   `json:"region,omitempty"`

	// Related data
	Role *GamingRole     `json:"role,omitempty"`
	User *GamingUserInfo `json:"user,omitempty"`
	Team *GamingTeam     `json:"team,omitempty"`
}

// RoleAssignment represents role assignment details
type RoleAssignment struct {
	AssignedBy string                 `json:"assigned_by"`
	Reason     string                 `json:"reason,omitempty"`
	ExpiresAt  *time.Time             `json:"expires_at,omitempty"`
	Scope      string                 `json:"scope,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`
}

// TeamRoleAssignment represents team role assignment details
type TeamRoleAssignment struct {
	AssignedBy string     `json:"assigned_by"`
	Position   string     `json:"position,omitempty"`
	GameRole   string     `json:"game_role,omitempty"`
	Champion   []string   `json:"champion,omitempty"`
	Reason     string     `json:"reason,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

// RoleAuditLog represents RBAC audit log entry
type RoleAuditLog struct {
	ID           string                 `json:"id"`
	Action       string                 `json:"action"` // role_created, permission_assigned, user_role_assigned, etc.
	ActorID      string                 `json:"actor_id"`
	ActorType    string                 `json:"actor_type"` // user, system, admin
	TargetID     string                 `json:"target_id"`
	TargetType   string                 `json:"target_type"` // user, role, permission, team
	RoleID       *string                `json:"role_id,omitempty"`
	PermissionID *string                `json:"permission_id,omitempty"`
	Changes      map[string]interface{} `json:"changes,omitempty"`
	IPAddress    string                 `json:"ip_address"`
	UserAgent    string                 `json:"user_agent"`
	Timestamp    time.Time              `json:"timestamp"`

	// Gaming context
	GamingAction string  `json:"gaming_action,omitempty"`
	TeamID       *string `json:"team_id,omitempty"`
}

// Enums and constants

type RoleType string

const (
	RoleTypeSystem       RoleType = "system"
	RoleTypeGaming       RoleType = "gaming"
	RoleTypeTeam         RoleType = "team"
	RoleTypeSubscription RoleType = "subscription"
	RoleTypeCustom       RoleType = "custom"
)

type PermissionScope string

const (
	ScopeSelf         PermissionScope = "self"
	ScopeTeam         PermissionScope = "team"
	ScopeOrganization PermissionScope = "organization"
	ScopeGlobal       PermissionScope = "global"
)

// Predefined gaming roles and permissions

var DefaultGamingRoles = []*GamingRole{
	{
		Name:        "gaming:user",
		DisplayName: "Gaming User",
		Description: "Basic gaming platform user",
		Type:        RoleTypeGaming,
		Category:    "gaming",
		Level:       1,
		IsSystem:    true,
		IsActive:    true,
	},
	{
		Name:        "gaming:premium",
		DisplayName: "Premium Gaming User",
		Description: "Premium gaming platform user with advanced features",
		Type:        RoleTypeGaming,
		Category:    "gaming",
		Level:       2,
		IsSystem:    true,
		IsActive:    true,
	},
	{
		Name:        "gaming:pro",
		DisplayName: "Pro Gaming User",
		Description: "Professional gaming user with team features",
		Type:        RoleTypeGaming,
		Category:    "gaming",
		Level:       3,
		IsSystem:    true,
		IsActive:    true,
	},
	{
		Name:        "team:player",
		DisplayName: "Team Player",
		Description: "Gaming team player",
		Type:        RoleTypeTeam,
		Category:    "team",
		Level:       1,
		IsSystem:    true,
		IsActive:    true,
	},
	{
		Name:        "team:captain",
		DisplayName: "Team Captain",
		Description: "Gaming team captain with management permissions",
		Type:        RoleTypeTeam,
		Category:    "team",
		Level:       2,
		IsSystem:    true,
		IsActive:    true,
	},
	{
		Name:        "team:coach",
		DisplayName: "Team Coach",
		Description: "Gaming team coach with analysis permissions",
		Type:        RoleTypeTeam,
		Category:    "team",
		Level:       2,
		IsSystem:    true,
		IsActive:    true,
	},
	{
		Name:        "team:manager",
		DisplayName: "Team Manager",
		Description: "Gaming team manager with full team permissions",
		Type:        RoleTypeTeam,
		Category:    "team",
		Level:       3,
		IsSystem:    true,
		IsActive:    true,
	},
	{
		Name:        "admin:gaming",
		DisplayName: "Gaming Administrator",
		Description: "Gaming platform administrator",
		Type:        RoleTypeSystem,
		Category:    "admin",
		Level:       10,
		IsSystem:    true,
		IsActive:    true,
	},
}

var DefaultGamingPermissions = []*GamingPermission{
	// Gaming Analytics Permissions
	{
		Name:        "gaming:analytics:view",
		DisplayName: "View Gaming Analytics",
		Description: "View basic gaming analytics and statistics",
		Category:    "gaming",
		Resource:    "analytics",
		Action:      "view",
		Scope:       ScopeSelf,
		IsSystem:    true,
		IsActive:    true,
	},
	{
		Name:             "gaming:analytics:advanced",
		DisplayName:      "Advanced Gaming Analytics",
		Description:      "Access advanced gaming analytics and insights",
		Category:         "gaming",
		Resource:         "analytics",
		Action:           "view",
		Scope:            ScopeSelf,
		SubscriptionTier: "premium",
		IsSystem:         true,
		IsActive:         true,
	},
	{
		Name:             "gaming:analytics:export",
		DisplayName:      "Export Gaming Analytics",
		Description:      "Export gaming analytics data",
		Category:         "gaming",
		Resource:         "analytics",
		Action:           "export",
		Scope:            ScopeSelf,
		RequiresMFA:      true,
		SubscriptionTier: "pro",
		IsSystem:         true,
		IsActive:         true,
	},

	// Team Management Permissions
	{
		Name:        "team:view",
		DisplayName: "View Team",
		Description: "View team information and basic statistics",
		Category:    "team",
		Resource:    "team",
		Action:      "view",
		Scope:       ScopeTeam,
		IsSystem:    true,
		IsActive:    true,
	},
	{
		Name:        "team:manage:players",
		DisplayName: "Manage Team Players",
		Description: "Add, remove, and manage team players",
		Category:    "team",
		Resource:    "team",
		Action:      "manage",
		Scope:       ScopeTeam,
		RequiresMFA: true,
		IsSystem:    true,
		IsActive:    true,
	},
	{
		Name:        "team:manage:settings",
		DisplayName: "Manage Team Settings",
		Description: "Modify team settings and configuration",
		Category:    "team",
		Resource:    "team",
		Action:      "manage",
		Scope:       ScopeTeam,
		RequiresMFA: true,
		IsSystem:    true,
		IsActive:    true,
	},

	// API Access Permissions
	{
		Name:        "api:basic",
		DisplayName: "Basic API Access",
		Description: "Basic API access with rate limiting",
		Category:    "api",
		Resource:    "api",
		Action:      "access",
		Scope:       ScopeSelf,
		IsSystem:    true,
		IsActive:    true,
	},
	{
		Name:             "api:extended",
		DisplayName:      "Extended API Access",
		Description:      "Extended API access with higher rate limits",
		Category:         "api",
		Resource:         "api",
		Action:           "access",
		Scope:            ScopeSelf,
		SubscriptionTier: "premium",
		IsSystem:         true,
		IsActive:         true,
	},
	{
		Name:             "api:unlimited",
		DisplayName:      "Unlimited API Access",
		Description:      "Unlimited API access for enterprise users",
		Category:         "api",
		Resource:         "api",
		Action:           "access",
		Scope:            ScopeGlobal,
		SubscriptionTier: "enterprise",
		IsSystem:         true,
		IsActive:         true,
	},

	// Admin Permissions
	{
		Name:        "admin:users:manage",
		DisplayName: "Manage Users",
		Description: "Manage gaming platform users",
		Category:    "admin",
		Resource:    "users",
		Action:      "manage",
		Scope:       ScopeGlobal,
		RequiresMFA: true,
		IsSystem:    true,
		IsActive:    true,
	},
	{
		Name:        "admin:roles:manage",
		DisplayName: "Manage Roles",
		Description: "Manage gaming platform roles and permissions",
		Category:    "admin",
		Resource:    "roles",
		Action:      "manage",
		Scope:       ScopeGlobal,
		RequiresMFA: true,
		IsSystem:    true,
		IsActive:    true,
	},
}

// Role hierarchy for inheritance
type RoleHierarchy struct {
	hierarchy map[string][]string // role_id -> child_role_ids
	levels    map[string]int      // role_id -> level
}

// Filters for queries

type RoleFilters struct {
	Type     *RoleType `json:"type,omitempty"`
	Category *string   `json:"category,omitempty"`
	IsActive *bool     `json:"is_active,omitempty"`
	IsSystem *bool     `json:"is_system,omitempty"`
	Level    *int      `json:"level,omitempty"`
}

type PermissionFilters struct {
	Category    *string          `json:"category,omitempty"`
	Resource    *string          `json:"resource,omitempty"`
	Action      *string          `json:"action,omitempty"`
	Scope       *PermissionScope `json:"scope,omitempty"`
	RequiresMFA *bool            `json:"requires_mfa,omitempty"`
	IsActive    *bool            `json:"is_active,omitempty"`
}

type AuditFilters struct {
	ActorID   *string    `json:"actor_id,omitempty"`
	TargetID  *string    `json:"target_id,omitempty"`
	Action    *string    `json:"action,omitempty"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	TeamID    *string    `json:"team_id,omitempty"`
}

// GamingTeam represents a gaming team (simplified)
type GamingTeam struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Region      string    `json:"region"`
	GameType    string    `json:"game_type"`
	CreatedAt   time.Time `json:"created_at"`
}

// NewGamingRBACManager creates new RBAC manager for gaming platform
func NewGamingRBACManager(
	db *gorm.DB,
	rbacStore RBACStore,
	gamingAnalytics GamingAnalyticsService,
	permissionCache PermissionCache,
	config *GamingRBACConfig,
) *GamingRBACManager {
	// Set gaming-specific defaults
	if config.DefaultRole == "" {
		config.DefaultRole = "gaming:user"
	}
	if config.AdminRole == "" {
		config.AdminRole = "admin:gaming"
	}
	if config.GamingRolePrefix == "" {
		config.GamingRolePrefix = "gaming:"
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = 15 * time.Minute
	}
	if config.MaxRolesPerUser == 0 {
		config.MaxRolesPerUser = 10
	}
	if config.MaxPermissionsPerRole == 0 {
		config.MaxPermissionsPerRole = 100
	}
	if config.AuditRetentionDays == 0 {
		config.AuditRetentionDays = 90
	}

	manager := &GamingRBACManager{
		db:              db,
		rbacStore:       rbacStore,
		gamingAnalytics: gamingAnalytics,
		permissionCache: permissionCache,
		config:          config,
		roleHierarchy: &RoleHierarchy{
			hierarchy: make(map[string][]string),
			levels:    make(map[string]int),
		},
	}

	return manager
}

// Core RBAC operations

// HasPermission checks if user has specific permission
func (rbac *GamingRBACManager) HasPermission(ctx context.Context, userID, permission string) (bool, error) {
	// Check cache first
	if rbac.config.EnablePermissionCache {
		cachedPermissions, err := rbac.permissionCache.GetUserPermissions(ctx, userID)
		if err == nil {
			for _, perm := range cachedPermissions {
				if perm == permission {
					return true, nil
				}
			}
			return false, nil
		}
	}

	// Get user permissions from database
	permissions, err := rbac.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// Check for exact match or wildcard
	for _, perm := range permissions {
		if perm == permission || rbac.matchesWildcard(perm, permission) {
			return true, nil
		}
	}

	return false, nil
}

// GetUserPermissions retrieves all permissions for a user
func (rbac *GamingRBACManager) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	// Get user roles
	userRoles, err := rbac.rbacStore.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	permissionMap := make(map[string]bool)

	// Collect permissions from all roles
	for _, userRole := range userRoles {
		if !userRole.IsActive {
			continue
		}

		// Check if role has expired
		if userRole.ExpiresAt != nil && time.Now().After(*userRole.ExpiresAt) {
			continue
		}

		// Get role permissions
		rolePermissions, err := rbac.rbacStore.GetRolePermissions(ctx, userRole.RoleID)
		if err != nil {
			continue // Skip this role if we can't get permissions
		}

		for _, permission := range rolePermissions {
			if permission.IsActive {
				permissionMap[permission.Name] = true
			}
		}

		// Include inherited permissions if hierarchy is enabled
		if rbac.config.EnableInheritance {
			inheritedPermissions, err := rbac.getInheritedPermissions(ctx, userRole.RoleID)
			if err == nil {
				for _, permission := range inheritedPermissions {
					permissionMap[permission] = true
				}
			}
		}
	}

	// Get team-based permissions
	teamRoles, err := rbac.rbacStore.GetUserTeamRoles(ctx, userID)
	if err == nil {
		for _, teamRole := range teamRoles {
			if !teamRole.IsActive {
				continue
			}

			// Check if team role has expired
			if teamRole.ExpiresAt != nil && time.Now().After(*teamRole.ExpiresAt) {
				continue
			}

			rolePermissions, err := rbac.rbacStore.GetRolePermissions(ctx, teamRole.RoleID)
			if err != nil {
				continue
			}

			for _, permission := range rolePermissions {
				if permission.IsActive {
					// Add team context to permission
					teamPermission := fmt.Sprintf("%s:team:%s", permission.Name, teamRole.TeamID)
					permissionMap[teamPermission] = true
					// Also add the base permission
					permissionMap[permission.Name] = true
				}
			}
		}
	}

	// Convert map to slice
	var permissions []string
	for permission := range permissionMap {
		permissions = append(permissions, permission)
	}

	// Cache permissions
	if rbac.config.EnablePermissionCache {
		rbac.permissionCache.SetUserPermissions(ctx, userID, permissions, rbac.config.CacheTTL)
	}

	return permissions, nil
}

// Helper methods

func (rbac *GamingRBACManager) matchesWildcard(pattern, permission string) bool {
	// Simple wildcard matching for permissions like "gaming:*" or "team:*:manage"
	if !strings.Contains(pattern, "*") {
		return false
	}

	patternParts := strings.Split(pattern, ":")
	permissionParts := strings.Split(permission, ":")

	if len(patternParts) != len(permissionParts) {
		return false
	}

	for i, part := range patternParts {
		if part != "*" && part != permissionParts[i] {
			return false
		}
	}

	return true
}

func (rbac *GamingRBACManager) getInheritedPermissions(ctx context.Context, roleID string) ([]string, error) {
	// Get parent role
	role, err := rbac.rbacStore.GetRole(ctx, roleID)
	if err != nil || role.ParentRoleID == nil {
		return []string{}, nil
	}

	// Get parent permissions
	parentPermissions, err := rbac.rbacStore.GetRolePermissions(ctx, *role.ParentRoleID)
	if err != nil {
		return []string{}, nil
	}

	var permissions []string
	for _, permission := range parentPermissions {
		if permission.IsActive {
			permissions = append(permissions, permission.Name)
		}
	}

	// Recursively get inherited permissions
	inheritedPermissions, err := rbac.getInheritedPermissions(ctx, *role.ParentRoleID)
	if err == nil {
		permissions = append(permissions, inheritedPermissions...)
	}

	return permissions, nil
}
