package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheService provides Redis caching functionality for analytics
type CacheService struct {
	client *redis.Client
	ctx    context.Context
}

// CacheConfig holds Redis configuration
type CacheConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	Enabled  bool
}

// NewCacheService creates a new Redis cache service
func NewCacheService(config CacheConfig) *CacheService {
	if !config.Enabled {
		log.Println("📦 Cache Redis désactivé")
		return &CacheService{
			client: nil,
			ctx:    context.Background(),
		}
	}

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// Test connection
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("⚠️  Redis connection failed: %v", err)
		log.Println("📦 Continuing without cache...")
		return &CacheService{
			client: nil,
			ctx:    ctx,
		}
	}

	log.Println("🚀 Redis cache service connected successfully")
	return &CacheService{
		client: rdb,
		ctx:    ctx,
	}
}

// IsEnabled returns true if Redis cache is enabled and connected
func (cs *CacheService) IsEnabled() bool {
	return cs.client != nil
}

// GetString retrieves a string value from cache
func (cs *CacheService) GetString(key string) (string, error) {
	if !cs.IsEnabled() {
		return "", fmt.Errorf("cache not enabled")
	}

	return cs.client.Get(cs.ctx, key).Result()
}

// SetString stores a string value in cache with TTL
func (cs *CacheService) SetString(key string, value string, ttl time.Duration) error {
	if !cs.IsEnabled() {
		return nil // Silently ignore if cache disabled
	}

	return cs.client.Set(cs.ctx, key, value, ttl).Err()
}

// GetJSON retrieves and unmarshals a JSON object from cache
func (cs *CacheService) GetJSON(key string, dest interface{}) error {
	if !cs.IsEnabled() {
		return fmt.Errorf("cache not enabled")
	}

	jsonStr, err := cs.client.Get(cs.ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(jsonStr), dest)
}

// SetJSON marshals and stores a JSON object in cache with TTL
func (cs *CacheService) SetJSON(key string, value interface{}, ttl time.Duration) error {
	if !cs.IsEnabled() {
		return nil // Silently ignore if cache disabled
	}

	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return cs.client.Set(cs.ctx, key, jsonBytes, ttl).Err()
}

// Delete removes a key from cache
func (cs *CacheService) Delete(key string) error {
	if !cs.IsEnabled() {
		return nil
	}

	return cs.client.Del(cs.ctx, key).Err()
}

// DeletePattern removes all keys matching a pattern
func (cs *CacheService) DeletePattern(pattern string) error {
	if !cs.IsEnabled() {
		return nil
	}

	keys, err := cs.client.Keys(cs.ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return cs.client.Del(cs.ctx, keys...).Err()
	}

	return nil
}

// Exists checks if a key exists in cache
func (cs *CacheService) Exists(key string) bool {
	if !cs.IsEnabled() {
		return false
	}

	count, err := cs.client.Exists(cs.ctx, key).Result()
	return err == nil && count > 0
}

// TTL returns the remaining TTL of a key
func (cs *CacheService) TTL(key string) (time.Duration, error) {
	if !cs.IsEnabled() {
		return 0, fmt.Errorf("cache not enabled")
	}

	return cs.client.TTL(cs.ctx, key).Result()
}

// GetStats returns cache statistics
func (cs *CacheService) GetStats() map[string]interface{} {
	if !cs.IsEnabled() {
		return map[string]interface{}{
			"enabled": false,
			"status":  "disabled",
		}
	}

	info, err := cs.client.Info(cs.ctx).Result()
	if err != nil {
		return map[string]interface{}{
			"enabled": true,
			"status":  "error",
			"error":   err.Error(),
		}
	}

	// Parse basic stats from Redis INFO
	return map[string]interface{}{
		"enabled":     true,
		"status":      "connected",
		"info_length": len(info),
	}
}

// FlushAll clears all keys from cache (use with caution!)
func (cs *CacheService) FlushAll() error {
	if !cs.IsEnabled() {
		return nil
	}

	return cs.client.FlushAll(cs.ctx).Err()
}

// Close closes the Redis connection
func (cs *CacheService) Close() error {
	if !cs.IsEnabled() {
		return nil
	}

	return cs.client.Close()
}

// Cache key generators for different types of data

// UserCacheKey generates cache keys for user-specific data
func UserCacheKey(userID int, dataType string) string {
	return fmt.Sprintf("user:%d:%s", userID, dataType)
}

// AnalyticsCacheKey generates cache keys for analytics data
func AnalyticsCacheKey(userID int, period string, dataType string) string {
	return fmt.Sprintf("analytics:%d:%s:%s", userID, period, dataType)
}

// MMRCacheKey generates cache keys for MMR data
func MMRCacheKey(userID int, days int) string {
	return fmt.Sprintf("mmr:%d:%d", userID, days)
}

// RecommendationCacheKey generates cache keys for recommendations
func RecommendationCacheKey(userID int) string {
	return fmt.Sprintf("recommendations:%d", userID)
}

// ChampionCacheKey generates cache keys for champion data
func ChampionCacheKey(userID int, championID int, period string) string {
	return fmt.Sprintf("champion:%d:%d:%s", userID, championID, period)
}

// GlobalCacheKey generates cache keys for global/shared data
func GlobalCacheKey(dataType string) string {
	return fmt.Sprintf("global:%s", dataType)
}

// Cache TTL constants
const (
	// Short-term cache (5 minutes) - for frequently changing data
	TTLShort = 5 * time.Minute

	// Medium-term cache (1 hour) - for analytics data
	TTLMedium = 1 * time.Hour

	// Long-term cache (24 hours) - for user stats and recommendations
	TTLLong = 24 * time.Hour

	// Very long-term cache (7 days) - for historical data and meta information
	TTLVeryLong = 7 * 24 * time.Hour
)