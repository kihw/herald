package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// CacheItem represents a cached item with expiration
type CacheItem struct {
	Data       interface{} `json:"data"`
	Expiration int64       `json:"expiration"`
	CreatedAt  int64       `json:"created_at"`
	HitCount   int64       `json:"hit_count"`
}

// Cache represents an intelligent in-memory cache
type Cache struct {
	items map[string]*CacheItem
	mutex sync.RWMutex
	stats CacheStats
}

// CacheStats tracks cache performance metrics
type CacheStats struct {
	Hits         int64   `json:"hits"`
	Misses       int64   `json:"misses"`
	Evictions    int64   `json:"evictions"`
	Size         int     `json:"size"`
	HitRatio     float64 `json:"hit_ratio"`
	LastCleanup  int64   `json:"last_cleanup"`
	TotalMemory  int64   `json:"total_memory_bytes"`
}

// NewCache creates a new intelligent cache
func NewCache() *Cache {
	cache := &Cache{
		items: make(map[string]*CacheItem),
		stats: CacheStats{
			LastCleanup: time.Now().Unix(),
		},
	}
	
	// Start background cleanup goroutine
	go cache.startCleanupRoutine()
	
	log.Println("🧠 Intelligent cache system initialized")
	return cache
}

// Set stores an item in cache with TTL
func (c *Cache) Set(key string, value interface{}, ttlSeconds int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	expiration := time.Now().Add(time.Duration(ttlSeconds) * time.Second).Unix()
	
	c.items[key] = &CacheItem{
		Data:       value,
		Expiration: expiration,
		CreatedAt:  time.Now().Unix(),
		HitCount:   0,
	}
	
	c.stats.Size = len(c.items)
}

// Get retrieves an item from cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, exists := c.items[key]
	if !exists {
		c.stats.Misses++
		c.updateHitRatio()
		return nil, false
	}
	
	// Check if item has expired
	if time.Now().Unix() > item.Expiration {
		c.stats.Misses++
		c.updateHitRatio()
		// Don't delete here to avoid lock issues, cleanup will handle it
		return nil, false
	}
	
	// Update hit statistics
	item.HitCount++
	c.stats.Hits++
	c.updateHitRatio()
	
	return item.Data, true
}

// Delete removes an item from cache
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	delete(c.items, key)
	c.stats.Size = len(c.items)
}

// Clear removes all items from cache
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.items = make(map[string]*CacheItem)
	c.stats.Size = 0
	c.stats.Evictions++
	
	log.Println("🧹 Cache cleared")
}

// GetStats returns current cache statistics
func (c *Cache) GetStats() CacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	// Calculate memory usage estimate
	memoryUsage := int64(0)
	for _, item := range c.items {
		if data, err := json.Marshal(item); err == nil {
			memoryUsage += int64(len(data))
		}
	}
	
	stats := c.stats
	stats.Size = len(c.items)
	stats.TotalMemory = memoryUsage
	
	return stats
}

// updateHitRatio calculates the current hit ratio
func (c *Cache) updateHitRatio() {
	total := c.stats.Hits + c.stats.Misses
	if total > 0 {
		c.stats.HitRatio = float64(c.stats.Hits) / float64(total)
	}
}

// startCleanupRoutine runs background cleanup every 5 minutes
func (c *Cache) startCleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			c.cleanup()
		}
	}
}

// cleanup removes expired items
func (c *Cache) cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	now := time.Now().Unix()
	initialSize := len(c.items)
	removed := 0
	
	for key, item := range c.items {
		if now > item.Expiration {
			delete(c.items, key)
			removed++
		}
	}
	
	c.stats.Size = len(c.items)
	c.stats.Evictions += int64(removed)
	c.stats.LastCleanup = now
	
	if removed > 0 {
		log.Printf("🧹 Cache cleanup: removed %d expired items (%d → %d)", 
			removed, initialSize, len(c.items))
	}
}

// SmartCache provides intelligent caching strategies for different data types
type SmartCache struct {
	cache *Cache
}

// NewSmartCache creates a new smart cache with predefined strategies
func NewSmartCache() *SmartCache {
	return &SmartCache{
		cache: NewCache(),
	}
}

// Cache keys constants
const (
	CacheKeyUserStats       = "user_stats_%d"
	CacheKeyChampionStats   = "champion_stats_%d"
	CacheKeyGameModeStats   = "gamemode_stats_%d"
	CacheKeyRecommendations = "recommendations_%d"
	CacheKeyAnalysis        = "analysis_%d"
	CacheKeyMatches         = "matches_%d_%d_%d" // userID, page, limit
	CacheKeyTrends          = "trends_%d_%d"     // userID, days
)

// GetUserStats retrieves user stats with intelligent caching
func (sc *SmartCache) GetUserStats(userID int, fetcher func() (UserStats, error)) (UserStats, error) {
	key := fmt.Sprintf(CacheKeyUserStats, userID)
	
	if cached, hit := sc.cache.Get(key); hit {
		if stats, ok := cached.(UserStats); ok {
			log.Printf("📈 Cache HIT: User stats for user %d", userID)
			return stats, nil
		}
	}
	
	log.Printf("📉 Cache MISS: Fetching user stats for user %d", userID)
	stats, err := fetcher()
	if err != nil {
		return stats, err
	}
	
	// Cache for 5 minutes (stats change frequently)
	sc.cache.Set(key, stats, 300)
	return stats, nil
}

// GetChampionStats retrieves champion stats with caching
func (sc *SmartCache) GetChampionStats(userID int, fetcher func() ([]ChampionStats, error)) ([]ChampionStats, error) {
	key := fmt.Sprintf(CacheKeyChampionStats, userID)
	
	if cached, hit := sc.cache.Get(key); hit {
		if stats, ok := cached.([]ChampionStats); ok {
			log.Printf("📈 Cache HIT: Champion stats for user %d", userID)
			return stats, nil
		}
	}
	
	log.Printf("📉 Cache MISS: Fetching champion stats for user %d", userID)
	stats, err := fetcher()
	if err != nil {
		return stats, err
	}
	
	// Cache for 10 minutes (changes less frequently)
	sc.cache.Set(key, stats, 600)
	return stats, nil
}

// GetRecommendations retrieves AI recommendations with caching
func (sc *SmartCache) GetRecommendations(userID int, fetcher func() ([]Recommendation, error)) ([]Recommendation, error) {
	key := fmt.Sprintf(CacheKeyRecommendations, userID)
	
	if cached, hit := sc.cache.Get(key); hit {
		if recs, ok := cached.([]Recommendation); ok {
			log.Printf("📈 Cache HIT: Recommendations for user %d", userID)
			return recs, nil
		}
	}
	
	log.Printf("📉 Cache MISS: Generating recommendations for user %d", userID)
	recs, err := fetcher()
	if err != nil {
		return recs, err
	}
	
	// Cache for 15 minutes (AI recommendations are computationally expensive)
	sc.cache.Set(key, recs, 900)
	return recs, nil
}

// InvalidateUserCache removes all cached data for a specific user
func (sc *SmartCache) InvalidateUserCache(userID int) {
	keys := []string{
		fmt.Sprintf(CacheKeyUserStats, userID),
		fmt.Sprintf(CacheKeyChampionStats, userID),
		fmt.Sprintf(CacheKeyGameModeStats, userID),
		fmt.Sprintf(CacheKeyRecommendations, userID),
		fmt.Sprintf(CacheKeyAnalysis, userID),
	}
	
	for _, key := range keys {
		sc.cache.Delete(key)
	}
	
	log.Printf("🗑️ Invalidated cache for user %d", userID)
}

// GetCacheStats returns cache performance statistics
func (sc *SmartCache) GetCacheStats() CacheStats {
	return sc.cache.GetStats()
}

// Global smart cache instance
var smartCache *SmartCache

// InitializeCache initializes the global cache system
func InitializeCache() {
	smartCache = NewSmartCache()
}