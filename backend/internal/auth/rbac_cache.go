package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Herald.lol Gaming Analytics - RBAC Permission Cache
// Redis-based caching for user permissions to improve performance

// RedisPermissionCache implements PermissionCache using Redis
type RedisPermissionCache struct {
	redisClient RedisClient
	keyPrefix   string
}

// NewRedisPermissionCache creates new Redis permission cache
func NewRedisPermissionCache(redisClient RedisClient) *RedisPermissionCache {
	return &RedisPermissionCache{
		redisClient: redisClient,
		keyPrefix:   "herald:gaming:rbac:permissions:",
	}
}

// GetUserPermissions retrieves cached user permissions from Redis
func (c *RedisPermissionCache) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	key := c.keyPrefix + userID

	data, err := c.redisClient.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get cached gaming user permissions: %w", err)
	}

	var permissions []string
	if err := json.Unmarshal([]byte(data), &permissions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached gaming permissions: %w", err)
	}

	return permissions, nil
}

// SetUserPermissions caches user permissions in Redis
func (c *RedisPermissionCache) SetUserPermissions(ctx context.Context, userID string, permissions []string, ttl time.Duration) error {
	key := c.keyPrefix + userID

	data, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal gaming user permissions: %w", err)
	}

	if err := c.redisClient.Set(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("failed to cache gaming user permissions: %w", err)
	}

	return nil
}

// InvalidateUserPermissions removes user permissions from cache
func (c *RedisPermissionCache) InvalidateUserPermissions(ctx context.Context, userID string) error {
	key := c.keyPrefix + userID

	if err := c.redisClient.Del(ctx, key); err != nil {
		return fmt.Errorf("failed to invalidate gaming user permissions: %w", err)
	}

	return nil
}

// InvalidateRolePermissions removes all cached permissions for users with specific role
func (c *RedisPermissionCache) InvalidateRolePermissions(ctx context.Context, roleID string) error {
	// This would require tracking which users have which roles
	// For simplicity, we'll implement a pattern-based invalidation

	pattern := c.keyPrefix + "*"
	var cursor uint64

	for {
		keys, nextCursor, err := c.redisClient.Scan(ctx, cursor, pattern, 100)
		if err != nil {
			return fmt.Errorf("failed to scan gaming permission cache: %w", err)
		}

		// For each user, we need to check if they have the role
		// In a production system, you might maintain a separate index
		// For now, we'll invalidate all cached permissions (inefficient but safe)
		for _, key := range keys {
			c.redisClient.Del(ctx, key)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}

// Enhanced Permission Cache with metadata

// CachedUserPermissions represents cached user permissions with metadata
type CachedUserPermissions struct {
	UserID       string                 `json:"user_id"`
	Permissions  []string               `json:"permissions"`
	Roles        []string               `json:"roles"`      // Role IDs
	TeamRoles    []string               `json:"team_roles"` // Team role IDs
	Metadata     map[string]interface{} `json:"metadata"`   // Additional metadata
	CachedAt     time.Time              `json:"cached_at"`
	LastVerified time.Time              `json:"last_verified"` // Last permission verification
	TTL          time.Duration          `json:"ttl"`
}

// EnhancedRedisPermissionCache provides enhanced caching with metadata
type EnhancedRedisPermissionCache struct {
	*RedisPermissionCache
	userRoleKeyPrefix string
	roleKeyPrefix     string
}

// NewEnhancedRedisPermissionCache creates enhanced Redis permission cache
func NewEnhancedRedisPermissionCache(redisClient RedisClient) *EnhancedRedisPermissionCache {
	return &EnhancedRedisPermissionCache{
		RedisPermissionCache: NewRedisPermissionCache(redisClient),
		userRoleKeyPrefix:    "herald:gaming:rbac:user_roles:",
		roleKeyPrefix:        "herald:gaming:rbac:role_permissions:",
	}
}

// GetEnhancedUserPermissions retrieves cached user permissions with metadata
func (c *EnhancedRedisPermissionCache) GetEnhancedUserPermissions(ctx context.Context, userID string) (*CachedUserPermissions, error) {
	key := c.keyPrefix + userID

	data, err := c.redisClient.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get enhanced gaming user permissions: %w", err)
	}

	var cached CachedUserPermissions
	if err := json.Unmarshal([]byte(data), &cached); err != nil {
		return nil, fmt.Errorf("failed to unmarshal enhanced gaming permissions: %w", err)
	}

	return &cached, nil
}

// SetEnhancedUserPermissions caches user permissions with metadata
func (c *EnhancedRedisPermissionCache) SetEnhancedUserPermissions(ctx context.Context, cached *CachedUserPermissions, ttl time.Duration) error {
	key := c.keyPrefix + cached.UserID

	cached.CachedAt = time.Now()
	cached.TTL = ttl

	data, err := json.Marshal(cached)
	if err != nil {
		return fmt.Errorf("failed to marshal enhanced gaming permissions: %w", err)
	}

	if err := c.redisClient.Set(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("failed to cache enhanced gaming permissions: %w", err)
	}

	return nil
}

// CacheUserRoles caches user role assignments
func (c *EnhancedRedisPermissionCache) CacheUserRoles(ctx context.Context, userID string, roles []string, ttl time.Duration) error {
	key := c.userRoleKeyPrefix + userID

	roleData := map[string]interface{}{
		"user_id":   userID,
		"roles":     roles,
		"cached_at": time.Now(),
	}

	data, err := json.Marshal(roleData)
	if err != nil {
		return fmt.Errorf("failed to marshal gaming user roles: %w", err)
	}

	if err := c.redisClient.Set(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("failed to cache gaming user roles: %w", err)
	}

	return nil
}

// GetCachedUserRoles retrieves cached user roles
func (c *EnhancedRedisPermissionCache) GetCachedUserRoles(ctx context.Context, userID string) ([]string, error) {
	key := c.userRoleKeyPrefix + userID

	data, err := c.redisClient.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get cached gaming user roles: %w", err)
	}

	var roleData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &roleData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gaming user roles: %w", err)
	}

	rolesInterface, ok := roleData["roles"]
	if !ok {
		return []string{}, nil
	}

	rolesSlice, ok := rolesInterface.([]interface{})
	if !ok {
		return []string{}, nil
	}

	var roles []string
	for _, role := range rolesSlice {
		if roleStr, ok := role.(string); ok {
			roles = append(roles, roleStr)
		}
	}

	return roles, nil
}

// CacheRolePermissions caches permissions for a specific role
func (c *EnhancedRedisPermissionCache) CacheRolePermissions(ctx context.Context, roleID string, permissions []string, ttl time.Duration) error {
	key := c.roleKeyPrefix + roleID

	permData := map[string]interface{}{
		"role_id":     roleID,
		"permissions": permissions,
		"cached_at":   time.Now(),
	}

	data, err := json.Marshal(permData)
	if err != nil {
		return fmt.Errorf("failed to marshal gaming role permissions: %w", err)
	}

	if err := c.redisClient.Set(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("failed to cache gaming role permissions: %w", err)
	}

	return nil
}

// GetCachedRolePermissions retrieves cached role permissions
func (c *EnhancedRedisPermissionCache) GetCachedRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	key := c.roleKeyPrefix + roleID

	data, err := c.redisClient.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get cached gaming role permissions: %w", err)
	}

	var permData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &permData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gaming role permissions: %w", err)
	}

	permissionsInterface, ok := permData["permissions"]
	if !ok {
		return []string{}, nil
	}

	permissionsSlice, ok := permissionsInterface.([]interface{})
	if !ok {
		return []string{}, nil
	}

	var permissions []string
	for _, perm := range permissionsSlice {
		if permStr, ok := perm.(string); ok {
			permissions = append(permissions, permStr)
		}
	}

	return permissions, nil
}

// InvalidateAllUserPermissions invalidates all cached user permissions
func (c *EnhancedRedisPermissionCache) InvalidateAllUserPermissions(ctx context.Context) error {
	// Invalidate user permissions
	if err := c.invalidateByPattern(ctx, c.keyPrefix+"*"); err != nil {
		return fmt.Errorf("failed to invalidate user permissions: %w", err)
	}

	// Invalidate user roles
	if err := c.invalidateByPattern(ctx, c.userRoleKeyPrefix+"*"); err != nil {
		return fmt.Errorf("failed to invalidate user roles: %w", err)
	}

	return nil
}

// InvalidateRoleCache invalidates cached data for specific role
func (c *EnhancedRedisPermissionCache) InvalidateRoleCache(ctx context.Context, roleID string) error {
	// Invalidate role permissions
	roleKey := c.roleKeyPrefix + roleID
	if err := c.redisClient.Del(ctx, roleKey); err != nil {
		return fmt.Errorf("failed to invalidate gaming role permissions: %w", err)
	}

	// Invalidate all user permissions (since role changed)
	// In production, you might maintain an index of users with this role
	return c.InvalidateAllUserPermissions(ctx)
}

// GetCacheStats returns cache statistics
func (c *EnhancedRedisPermissionCache) GetCacheStats(ctx context.Context) (*CacheStats, error) {
	stats := &CacheStats{
		Timestamp: time.Now(),
	}

	// Count user permission keys
	userPermKeys, err := c.countKeys(ctx, c.keyPrefix+"*")
	if err == nil {
		stats.UserPermissionKeys = userPermKeys
	}

	// Count user role keys
	userRoleKeys, err := c.countKeys(ctx, c.userRoleKeyPrefix+"*")
	if err == nil {
		stats.UserRoleKeys = userRoleKeys
	}

	// Count role permission keys
	rolePermKeys, err := c.countKeys(ctx, c.roleKeyPrefix+"*")
	if err == nil {
		stats.RolePermissionKeys = rolePermKeys
	}

	stats.TotalKeys = stats.UserPermissionKeys + stats.UserRoleKeys + stats.RolePermissionKeys

	return stats, nil
}

// CacheStats represents cache statistics
type CacheStats struct {
	UserPermissionKeys int       `json:"user_permission_keys"`
	UserRoleKeys       int       `json:"user_role_keys"`
	RolePermissionKeys int       `json:"role_permission_keys"`
	TotalKeys          int       `json:"total_keys"`
	Timestamp          time.Time `json:"timestamp"`
	HitRate            float64   `json:"hit_rate,omitempty"`  // Would be calculated over time
	MissRate           float64   `json:"miss_rate,omitempty"` // Would be calculated over time
}

// Helper methods

func (c *EnhancedRedisPermissionCache) invalidateByPattern(ctx context.Context, pattern string) error {
	var cursor uint64

	for {
		keys, nextCursor, err := c.redisClient.Scan(ctx, cursor, pattern, 100)
		if err != nil {
			return err
		}

		for _, key := range keys {
			c.redisClient.Del(ctx, key)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}

func (c *EnhancedRedisPermissionCache) countKeys(ctx context.Context, pattern string) (int, error) {
	count := 0
	var cursor uint64

	for {
		keys, nextCursor, err := c.redisClient.Scan(ctx, cursor, pattern, 100)
		if err != nil {
			return 0, err
		}

		count += len(keys)

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return count, nil
}

// Cache warming utilities

// WarmUserPermissionCache pre-loads permissions for active users
func (c *EnhancedRedisPermissionCache) WarmUserPermissionCache(ctx context.Context, rbacManager *GamingRBACManager, userIDs []string) error {
	for _, userID := range userIDs {
		// Get user permissions from RBAC manager
		permissions, err := rbacManager.GetUserPermissions(ctx, userID)
		if err != nil {
			continue // Skip failed users
		}

		// Get user roles
		userRoles, err := rbacManager.rbacStore.GetUserRoles(ctx, userID)
		if err != nil {
			continue
		}

		var roleIDs []string
		for _, userRole := range userRoles {
			roleIDs = append(roleIDs, userRole.RoleID)
		}

		// Cache enhanced permissions
		cached := &CachedUserPermissions{
			UserID:       userID,
			Permissions:  permissions,
			Roles:        roleIDs,
			LastVerified: time.Now(),
			Metadata: map[string]interface{}{
				"warmed": true,
			},
		}

		c.SetEnhancedUserPermissions(ctx, cached, rbacManager.config.CacheTTL)
	}

	return nil
}

// Batch cache operations

// BatchInvalidateUsers invalidates permissions for multiple users
func (c *EnhancedRedisPermissionCache) BatchInvalidateUsers(ctx context.Context, userIDs []string) error {
	for _, userID := range userIDs {
		if err := c.InvalidateUserPermissions(ctx, userID); err != nil {
			// Log error but continue with other users
		}
	}
	return nil
}

// BatchCacheUserPermissions caches permissions for multiple users
func (c *EnhancedRedisPermissionCache) BatchCacheUserPermissions(ctx context.Context, userPermissions map[string][]string, ttl time.Duration) error {
	for userID, permissions := range userPermissions {
		if err := c.SetUserPermissions(ctx, userID, permissions, ttl); err != nil {
			// Log error but continue with other users
		}
	}
	return nil
}
