package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Herald.lol Gaming Analytics - RBAC Store Implementations
// Database implementations for Role-Based Access Control

// DatabaseRBACStore implements RBACStore using GORM database
type DatabaseRBACStore struct {
	db *gorm.DB
}

// Database models for RBAC

// GamingRoleRecord database model for roles
type GamingRoleRecord struct {
	ID            string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name          string    `gorm:"uniqueIndex;type:varchar(255)" json:"name"`
	DisplayName   string    `gorm:"type:varchar(255)" json:"display_name"`
	Description   string    `gorm:"type:text" json:"description"`
	Type          string    `gorm:"type:varchar(50)" json:"type"`
	Category      string    `gorm:"type:varchar(100);index" json:"category"`
	Level         int       `gorm:"default:1;index" json:"level"`
	ParentRoleID  *string   `gorm:"type:varchar(255);index" json:"parent_role_id,omitempty"`
	IsSystem      bool      `gorm:"default:false;index" json:"is_system"`
	IsActive      bool      `gorm:"default:true;index" json:"is_active"`
	GamingContext string    `gorm:"type:jsonb" json:"gaming_context"` // JSON
	Metadata      string    `gorm:"type:jsonb" json:"metadata"`       // JSON
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedBy     string    `gorm:"type:varchar(255)" json:"created_by"`
}

// GamingPermissionRecord database model for permissions
type GamingPermissionRecord struct {
	ID               string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name             string    `gorm:"uniqueIndex;type:varchar(255)" json:"name"`
	DisplayName      string    `gorm:"type:varchar(255)" json:"display_name"`
	Description      string    `gorm:"type:text" json:"description"`
	Category         string    `gorm:"type:varchar(100);index" json:"category"`
	Resource         string    `gorm:"type:varchar(100);index" json:"resource"`
	Action           string    `gorm:"type:varchar(100);index" json:"action"`
	Scope            string    `gorm:"type:varchar(50);index" json:"scope"`
	RequiresMFA      bool      `gorm:"default:false;index" json:"requires_mfa"`
	SubscriptionTier string    `gorm:"type:varchar(50);index" json:"subscription_tier"`
	GamingContext    string    `gorm:"type:jsonb" json:"gaming_context"` // JSON
	IsSystem         bool      `gorm:"default:false;index" json:"is_system"`
	IsActive         bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedBy        string    `gorm:"type:varchar(255)" json:"created_by"`
}

// GamingRolePermissionRecord database model for role-permission mapping
type GamingRolePermissionRecord struct {
	ID           string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	RoleID       string    `gorm:"type:varchar(255);index;uniqueIndex:role_permission_unique" json:"role_id"`
	PermissionID string    `gorm:"type:varchar(255);index;uniqueIndex:role_permission_unique" json:"permission_id"`
	AssignedBy   string    `gorm:"type:varchar(255)" json:"assigned_by"`
	AssignedAt   time.Time `json:"assigned_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// GamingUserRoleRecord database model for user-role mapping
type GamingUserRoleRecord struct {
	ID         string     `gorm:"primaryKey;type:varchar(255)" json:"id"`
	UserID     string     `gorm:"type:varchar(255);index;uniqueIndex:user_role_unique" json:"user_id"`
	RoleID     string     `gorm:"type:varchar(255);index;uniqueIndex:user_role_unique" json:"role_id"`
	AssignedBy string     `gorm:"type:varchar(255)" json:"assigned_by"`
	AssignedAt time.Time  `json:"assigned_at"`
	ExpiresAt  *time.Time `gorm:"index" json:"expires_at,omitempty"`
	IsActive   bool       `gorm:"default:true;index" json:"is_active"`
	Scope      string     `gorm:"type:varchar(255)" json:"scope"`
	Context    string     `gorm:"type:jsonb" json:"context"` // JSON
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// GamingTeamRoleRecord database model for team-role mapping
type GamingTeamRoleRecord struct {
	ID         string     `gorm:"primaryKey;type:varchar(255)" json:"id"`
	TeamID     string     `gorm:"type:varchar(255);index" json:"team_id"`
	UserID     string     `gorm:"type:varchar(255);index" json:"user_id"`
	RoleID     string     `gorm:"type:varchar(255);index" json:"role_id"`
	Position   string     `gorm:"type:varchar(100)" json:"position"`
	AssignedBy string     `gorm:"type:varchar(255)" json:"assigned_by"`
	AssignedAt time.Time  `json:"assigned_at"`
	ExpiresAt  *time.Time `gorm:"index" json:"expires_at,omitempty"`
	IsActive   bool       `gorm:"default:true;index" json:"is_active"`
	GameRole   string     `gorm:"type:varchar(100)" json:"game_role"`
	Champion   string     `gorm:"type:jsonb" json:"champion"` // JSON array
	Region     string     `gorm:"type:varchar(50)" json:"region"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// GamingRoleAuditLogRecord database model for role audit log
type GamingRoleAuditLogRecord struct {
	ID           string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Action       string    `gorm:"type:varchar(100);index" json:"action"`
	ActorID      string    `gorm:"type:varchar(255);index" json:"actor_id"`
	ActorType    string    `gorm:"type:varchar(50)" json:"actor_type"`
	TargetID     string    `gorm:"type:varchar(255);index" json:"target_id"`
	TargetType   string    `gorm:"type:varchar(50)" json:"target_type"`
	RoleID       *string   `gorm:"type:varchar(255);index" json:"role_id,omitempty"`
	PermissionID *string   `gorm:"type:varchar(255);index" json:"permission_id,omitempty"`
	Changes      string    `gorm:"type:jsonb" json:"changes"` // JSON
	IPAddress    string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent    string    `gorm:"type:text" json:"user_agent"`
	Timestamp    time.Time `gorm:"index" json:"timestamp"`
	GamingAction string    `gorm:"type:varchar(255)" json:"gaming_action"`
	TeamID       *string   `gorm:"type:varchar(255);index" json:"team_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// NewDatabaseRBACStore creates new database RBAC store
func NewDatabaseRBACStore(db *gorm.DB) *DatabaseRBACStore {
	return &DatabaseRBACStore{db: db}
}

// Role methods

// CreateRole creates a new role
func (s *DatabaseRBACStore) CreateRole(ctx context.Context, role *GamingRole) error {
	// Convert to database record
	record, err := s.roleToRecord(role)
	if err != nil {
		return fmt.Errorf("failed to convert role to record: %w", err)
	}

	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	if err := s.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("failed to create gaming role: %w", err)
	}

	// Update role with generated ID
	role.ID = record.ID
	role.CreatedAt = record.CreatedAt
	role.UpdatedAt = record.UpdatedAt

	return nil
}

// GetRole retrieves role by ID
func (s *DatabaseRBACStore) GetRole(ctx context.Context, roleID string) (*GamingRole, error) {
	var record GamingRoleRecord

	err := s.db.WithContext(ctx).Where("id = ?", roleID).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("gaming role not found")
		}
		return nil, fmt.Errorf("failed to get gaming role: %w", err)
	}

	return s.recordToRole(&record)
}

// GetRoleByName retrieves role by name
func (s *DatabaseRBACStore) GetRoleByName(ctx context.Context, name string) (*GamingRole, error) {
	var record GamingRoleRecord

	err := s.db.WithContext(ctx).Where("name = ?", name).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("gaming role not found")
		}
		return nil, fmt.Errorf("failed to get gaming role by name: %w", err)
	}

	return s.recordToRole(&record)
}

// UpdateRole updates existing role
func (s *DatabaseRBACStore) UpdateRole(ctx context.Context, role *GamingRole) error {
	record, err := s.roleToRecord(role)
	if err != nil {
		return fmt.Errorf("failed to convert role to record: %w", err)
	}

	record.UpdatedAt = time.Now()

	if err := s.db.WithContext(ctx).Save(record).Error; err != nil {
		return fmt.Errorf("failed to update gaming role: %w", err)
	}

	role.UpdatedAt = record.UpdatedAt
	return nil
}

// DeleteRole deletes role (soft delete by marking inactive)
func (s *DatabaseRBACStore) DeleteRole(ctx context.Context, roleID string) error {
	result := s.db.WithContext(ctx).Model(&GamingRoleRecord{}).
		Where("id = ?", roleID).
		Update("is_active", false)

	if result.Error != nil {
		return fmt.Errorf("failed to delete gaming role: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("gaming role not found")
	}

	return nil
}

// ListRoles lists roles with filters
func (s *DatabaseRBACStore) ListRoles(ctx context.Context, filters *RoleFilters) ([]*GamingRole, error) {
	query := s.db.WithContext(ctx).Model(&GamingRoleRecord{})

	if filters != nil {
		if filters.Type != nil {
			query = query.Where("type = ?", *filters.Type)
		}
		if filters.Category != nil {
			query = query.Where("category = ?", *filters.Category)
		}
		if filters.IsActive != nil {
			query = query.Where("is_active = ?", *filters.IsActive)
		}
		if filters.IsSystem != nil {
			query = query.Where("is_system = ?", *filters.IsSystem)
		}
		if filters.Level != nil {
			query = query.Where("level = ?", *filters.Level)
		}
	}

	var records []GamingRoleRecord
	if err := query.Order("level ASC, name ASC").Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to list gaming roles: %w", err)
	}

	var roles []*GamingRole
	for _, record := range records {
		role, err := s.recordToRole(&record)
		if err != nil {
			continue // Skip invalid records
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// Permission methods

// CreatePermission creates a new permission
func (s *DatabaseRBACStore) CreatePermission(ctx context.Context, permission *GamingPermission) error {
	record, err := s.permissionToRecord(permission)
	if err != nil {
		return fmt.Errorf("failed to convert permission to record: %w", err)
	}

	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	if err := s.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("failed to create gaming permission: %w", err)
	}

	permission.ID = record.ID
	permission.CreatedAt = record.CreatedAt
	permission.UpdatedAt = record.UpdatedAt

	return nil
}

// GetPermission retrieves permission by ID
func (s *DatabaseRBACStore) GetPermission(ctx context.Context, permissionID string) (*GamingPermission, error) {
	var record GamingPermissionRecord

	err := s.db.WithContext(ctx).Where("id = ?", permissionID).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("gaming permission not found")
		}
		return nil, fmt.Errorf("failed to get gaming permission: %w", err)
	}

	return s.recordToPermission(&record)
}

// GetPermissionByName retrieves permission by name
func (s *DatabaseRBACStore) GetPermissionByName(ctx context.Context, name string) (*GamingPermission, error) {
	var record GamingPermissionRecord

	err := s.db.WithContext(ctx).Where("name = ?", name).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("gaming permission not found")
		}
		return nil, fmt.Errorf("failed to get gaming permission by name: %w", err)
	}

	return s.recordToPermission(&record)
}

// ListPermissions lists permissions with filters
func (s *DatabaseRBACStore) ListPermissions(ctx context.Context, filters *PermissionFilters) ([]*GamingPermission, error) {
	query := s.db.WithContext(ctx).Model(&GamingPermissionRecord{})

	if filters != nil {
		if filters.Category != nil {
			query = query.Where("category = ?", *filters.Category)
		}
		if filters.Resource != nil {
			query = query.Where("resource = ?", *filters.Resource)
		}
		if filters.Action != nil {
			query = query.Where("action = ?", *filters.Action)
		}
		if filters.Scope != nil {
			query = query.Where("scope = ?", *filters.Scope)
		}
		if filters.RequiresMFA != nil {
			query = query.Where("requires_mfa = ?", *filters.RequiresMFA)
		}
		if filters.IsActive != nil {
			query = query.Where("is_active = ?", *filters.IsActive)
		}
	}

	var records []GamingPermissionRecord
	if err := query.Order("category ASC, resource ASC, action ASC").Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to list gaming permissions: %w", err)
	}

	var permissions []*GamingPermission
	for _, record := range records {
		permission, err := s.recordToPermission(&record)
		if err != nil {
			continue // Skip invalid records
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// Role-Permission mapping methods

// AssignPermissionToRole assigns permission to role
func (s *DatabaseRBACStore) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	record := &GamingRolePermissionRecord{
		ID:           fmt.Sprintf("%s_%s", roleID, permissionID),
		RoleID:       roleID,
		PermissionID: permissionID,
		AssignedBy:   "system", // Would be set by caller
		AssignedAt:   time.Now(),
		CreatedAt:    time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("failed to assign permission to gaming role: %w", err)
	}

	return nil
}

// RemovePermissionFromRole removes permission from role
func (s *DatabaseRBACStore) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	result := s.db.WithContext(ctx).Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&GamingRolePermissionRecord{})

	if result.Error != nil {
		return fmt.Errorf("failed to remove permission from gaming role: %w", result.Error)
	}

	return nil
}

// GetRolePermissions gets all permissions for a role
func (s *DatabaseRBACStore) GetRolePermissions(ctx context.Context, roleID string) ([]*GamingPermission, error) {
	var records []GamingPermissionRecord

	err := s.db.WithContext(ctx).
		Table("gaming_permission_records").
		Joins("JOIN gaming_role_permission_records ON gaming_permission_records.id = gaming_role_permission_records.permission_id").
		Where("gaming_role_permission_records.role_id = ?", roleID).
		Find(&records).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	var permissions []*GamingPermission
	for _, record := range records {
		permission, err := s.recordToPermission(&record)
		if err != nil {
			continue // Skip invalid records
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// User-Role mapping methods

// AssignRoleToUser assigns role to user
func (s *DatabaseRBACStore) AssignRoleToUser(ctx context.Context, userID, roleID string, assignment *RoleAssignment) error {
	record := &GamingUserRoleRecord{
		ID:         fmt.Sprintf("%s_%s_%d", userID, roleID, time.Now().UnixNano()),
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: assignment.AssignedBy,
		AssignedAt: time.Now(),
		ExpiresAt:  assignment.ExpiresAt,
		IsActive:   true,
		Scope:      assignment.Scope,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Marshal context to JSON
	if assignment.Context != nil {
		contextJSON, err := json.Marshal(assignment.Context)
		if err == nil {
			record.Context = string(contextJSON)
		}
	}

	if err := s.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("failed to assign role to gaming user: %w", err)
	}

	return nil
}

// RemoveRoleFromUser removes role from user
func (s *DatabaseRBACStore) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	result := s.db.WithContext(ctx).Model(&GamingUserRoleRecord{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Update("is_active", false)

	if result.Error != nil {
		return fmt.Errorf("failed to remove role from gaming user: %w", result.Error)
	}

	return nil
}

// GetUserRoles gets all roles for a user
func (s *DatabaseRBACStore) GetUserRoles(ctx context.Context, userID string) ([]*UserRole, error) {
	var records []GamingUserRoleRecord

	err := s.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).
		Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	var userRoles []*UserRole
	for _, record := range records {
		userRole, err := s.recordToUserRole(&record)
		if err != nil {
			continue // Skip invalid records
		}
		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
}

// GetUsersWithRole gets all users with specific role
func (s *DatabaseRBACStore) GetUsersWithRole(ctx context.Context, roleID string) ([]*UserRole, error) {
	var records []GamingUserRoleRecord

	err := s.db.WithContext(ctx).Where("role_id = ? AND is_active = ?", roleID, true).
		Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get users with role: %w", err)
	}

	var userRoles []*UserRole
	for _, record := range records {
		userRole, err := s.recordToUserRole(&record)
		if err != nil {
			continue // Skip invalid records
		}
		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
}

// Team-Role mapping methods

// AssignRoleToTeam assigns role to team member
func (s *DatabaseRBACStore) AssignRoleToTeam(ctx context.Context, teamID, roleID string, assignment *TeamRoleAssignment) error {
	record := &GamingTeamRoleRecord{
		ID:         fmt.Sprintf("%s_%s_%d", teamID, roleID, time.Now().UnixNano()),
		TeamID:     teamID,
		RoleID:     roleID,
		Position:   assignment.Position,
		AssignedBy: assignment.AssignedBy,
		AssignedAt: time.Now(),
		ExpiresAt:  assignment.ExpiresAt,
		IsActive:   true,
		GameRole:   assignment.GameRole,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Marshal champion list to JSON
	if assignment.Champion != nil {
		championJSON, err := json.Marshal(assignment.Champion)
		if err == nil {
			record.Champion = string(championJSON)
		}
	}

	if err := s.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("failed to assign role to gaming team: %w", err)
	}

	return nil
}

// RemoveRoleFromTeam removes role from team
func (s *DatabaseRBACStore) RemoveRoleFromTeam(ctx context.Context, teamID, roleID string) error {
	result := s.db.WithContext(ctx).Model(&GamingTeamRoleRecord{}).
		Where("team_id = ? AND role_id = ?", teamID, roleID).
		Update("is_active", false)

	if result.Error != nil {
		return fmt.Errorf("failed to remove role from gaming team: %w", result.Error)
	}

	return nil
}

// GetTeamRoles gets all roles for a team
func (s *DatabaseRBACStore) GetTeamRoles(ctx context.Context, teamID string) ([]*TeamRole, error) {
	var records []GamingTeamRoleRecord

	err := s.db.WithContext(ctx).Where("team_id = ? AND is_active = ?", teamID, true).
		Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get team roles: %w", err)
	}

	var teamRoles []*TeamRole
	for _, record := range records {
		teamRole, err := s.recordToTeamRole(&record)
		if err != nil {
			continue // Skip invalid records
		}
		teamRoles = append(teamRoles, teamRole)
	}

	return teamRoles, nil
}

// GetUserTeamRoles gets all team roles for a user
func (s *DatabaseRBACStore) GetUserTeamRoles(ctx context.Context, userID string) ([]*TeamRole, error) {
	var records []GamingTeamRoleRecord

	err := s.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).
		Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user team roles: %w", err)
	}

	var teamRoles []*TeamRole
	for _, record := range records {
		teamRole, err := s.recordToTeamRole(&record)
		if err != nil {
			continue // Skip invalid records
		}
		teamRoles = append(teamRoles, teamRole)
	}

	return teamRoles, nil
}

// Audit methods

// LogRoleAction logs RBAC action for audit
func (s *DatabaseRBACStore) LogRoleAction(ctx context.Context, action *RoleAuditLog) error {
	record := &GamingRoleAuditLogRecord{
		ID:           action.ID,
		Action:       action.Action,
		ActorID:      action.ActorID,
		ActorType:    action.ActorType,
		TargetID:     action.TargetID,
		TargetType:   action.TargetType,
		RoleID:       action.RoleID,
		PermissionID: action.PermissionID,
		IPAddress:    action.IPAddress,
		UserAgent:    action.UserAgent,
		Timestamp:    action.Timestamp,
		GamingAction: action.GamingAction,
		TeamID:       action.TeamID,
		CreatedAt:    time.Now(),
	}

	// Marshal changes to JSON
	if action.Changes != nil {
		changesJSON, err := json.Marshal(action.Changes)
		if err == nil {
			record.Changes = string(changesJSON)
		}
	}

	if err := s.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("failed to log gaming role action: %w", err)
	}

	return nil
}

// GetRoleAuditLog gets role audit log with filters
func (s *DatabaseRBACStore) GetRoleAuditLog(ctx context.Context, filters *AuditFilters) ([]*RoleAuditLog, error) {
	query := s.db.WithContext(ctx).Model(&GamingRoleAuditLogRecord{})

	if filters != nil {
		if filters.ActorID != nil {
			query = query.Where("actor_id = ?", *filters.ActorID)
		}
		if filters.TargetID != nil {
			query = query.Where("target_id = ?", *filters.TargetID)
		}
		if filters.Action != nil {
			query = query.Where("action = ?", *filters.Action)
		}
		if filters.StartTime != nil {
			query = query.Where("timestamp >= ?", *filters.StartTime)
		}
		if filters.EndTime != nil {
			query = query.Where("timestamp <= ?", *filters.EndTime)
		}
		if filters.TeamID != nil {
			query = query.Where("team_id = ?", *filters.TeamID)
		}
	}

	var records []GamingRoleAuditLogRecord
	if err := query.Order("timestamp DESC").Limit(1000).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to get gaming role audit log: %w", err)
	}

	var auditLogs []*RoleAuditLog
	for _, record := range records {
		auditLog, err := s.recordToAuditLog(&record)
		if err != nil {
			continue // Skip invalid records
		}
		auditLogs = append(auditLogs, auditLog)
	}

	return auditLogs, nil
}

// Helper methods for record conversion

func (s *DatabaseRBACStore) roleToRecord(role *GamingRole) (*GamingRoleRecord, error) {
	record := &GamingRoleRecord{
		ID:           role.ID,
		Name:         role.Name,
		DisplayName:  role.DisplayName,
		Description:  role.Description,
		Type:         string(role.Type),
		Category:     role.Category,
		Level:        role.Level,
		ParentRoleID: role.ParentRoleID,
		IsSystem:     role.IsSystem,
		IsActive:     role.IsActive,
		CreatedBy:    role.CreatedBy,
	}

	// Marshal JSON fields
	if role.GamingContext != nil {
		if contextJSON, err := json.Marshal(role.GamingContext); err == nil {
			record.GamingContext = string(contextJSON)
		}
	}

	if role.Metadata != nil {
		if metadataJSON, err := json.Marshal(role.Metadata); err == nil {
			record.Metadata = string(metadataJSON)
		}
	}

	return record, nil
}

func (s *DatabaseRBACStore) recordToRole(record *GamingRoleRecord) (*GamingRole, error) {
	role := &GamingRole{
		ID:           record.ID,
		Name:         record.Name,
		DisplayName:  record.DisplayName,
		Description:  record.Description,
		Type:         RoleType(record.Type),
		Category:     record.Category,
		Level:        record.Level,
		ParentRoleID: record.ParentRoleID,
		IsSystem:     record.IsSystem,
		IsActive:     record.IsActive,
		CreatedAt:    record.CreatedAt,
		UpdatedAt:    record.UpdatedAt,
		CreatedBy:    record.CreatedBy,
	}

	// Unmarshal JSON fields
	if record.GamingContext != "" {
		var context map[string]interface{}
		if err := json.Unmarshal([]byte(record.GamingContext), &context); err == nil {
			role.GamingContext = context
		}
	}

	if record.Metadata != "" {
		var metadata map[string]string
		if err := json.Unmarshal([]byte(record.Metadata), &metadata); err == nil {
			role.Metadata = metadata
		}
	}

	return role, nil
}

func (s *DatabaseRBACStore) permissionToRecord(permission *GamingPermission) (*GamingPermissionRecord, error) {
	record := &GamingPermissionRecord{
		ID:               permission.ID,
		Name:             permission.Name,
		DisplayName:      permission.DisplayName,
		Description:      permission.Description,
		Category:         permission.Category,
		Resource:         permission.Resource,
		Action:           permission.Action,
		Scope:            string(permission.Scope),
		RequiresMFA:      permission.RequiresMFA,
		SubscriptionTier: permission.SubscriptionTier,
		IsSystem:         permission.IsSystem,
		IsActive:         permission.IsActive,
		CreatedBy:        permission.CreatedBy,
	}

	// Marshal JSON fields
	if permission.GamingContext != nil {
		if contextJSON, err := json.Marshal(permission.GamingContext); err == nil {
			record.GamingContext = string(contextJSON)
		}
	}

	return record, nil
}

func (s *DatabaseRBACStore) recordToPermission(record *GamingPermissionRecord) (*GamingPermission, error) {
	permission := &GamingPermission{
		ID:               record.ID,
		Name:             record.Name,
		DisplayName:      record.DisplayName,
		Description:      record.Description,
		Category:         record.Category,
		Resource:         record.Resource,
		Action:           record.Action,
		Scope:            PermissionScope(record.Scope),
		RequiresMFA:      record.RequiresMFA,
		SubscriptionTier: record.SubscriptionTier,
		IsSystem:         record.IsSystem,
		IsActive:         record.IsActive,
		CreatedAt:        record.CreatedAt,
		UpdatedAt:        record.UpdatedAt,
		CreatedBy:        record.CreatedBy,
	}

	// Unmarshal JSON fields
	if record.GamingContext != "" {
		var context map[string]interface{}
		if err := json.Unmarshal([]byte(record.GamingContext), &context); err == nil {
			permission.GamingContext = context
		}
	}

	return permission, nil
}

func (s *DatabaseRBACStore) recordToUserRole(record *GamingUserRoleRecord) (*UserRole, error) {
	userRole := &UserRole{
		ID:         record.ID,
		UserID:     record.UserID,
		RoleID:     record.RoleID,
		AssignedBy: record.AssignedBy,
		AssignedAt: record.AssignedAt,
		ExpiresAt:  record.ExpiresAt,
		IsActive:   record.IsActive,
		Scope:      record.Scope,
	}

	// Unmarshal context
	if record.Context != "" {
		var context map[string]interface{}
		if err := json.Unmarshal([]byte(record.Context), &context); err == nil {
			userRole.Context = context
		}
	}

	return userRole, nil
}

func (s *DatabaseRBACStore) recordToTeamRole(record *GamingTeamRoleRecord) (*TeamRole, error) {
	teamRole := &TeamRole{
		ID:         record.ID,
		TeamID:     record.TeamID,
		UserID:     record.UserID,
		RoleID:     record.RoleID,
		Position:   record.Position,
		AssignedBy: record.AssignedBy,
		AssignedAt: record.AssignedAt,
		ExpiresAt:  record.ExpiresAt,
		IsActive:   record.IsActive,
		GameRole:   record.GameRole,
		Region:     record.Region,
	}

	// Unmarshal champion list
	if record.Champion != "" {
		var champion []string
		if err := json.Unmarshal([]byte(record.Champion), &champion); err == nil {
			teamRole.Champion = champion
		}
	}

	return teamRole, nil
}

func (s *DatabaseRBACStore) recordToAuditLog(record *GamingRoleAuditLogRecord) (*RoleAuditLog, error) {
	auditLog := &RoleAuditLog{
		ID:           record.ID,
		Action:       record.Action,
		ActorID:      record.ActorID,
		ActorType:    record.ActorType,
		TargetID:     record.TargetID,
		TargetType:   record.TargetType,
		RoleID:       record.RoleID,
		PermissionID: record.PermissionID,
		IPAddress:    record.IPAddress,
		UserAgent:    record.UserAgent,
		Timestamp:    record.Timestamp,
		GamingAction: record.GamingAction,
		TeamID:       record.TeamID,
	}

	// Unmarshal changes
	if record.Changes != "" {
		var changes map[string]interface{}
		if err := json.Unmarshal([]byte(record.Changes), &changes); err == nil {
			auditLog.Changes = changes
		}
	}

	return auditLog, nil
}
