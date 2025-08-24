package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Herald.lol Gaming Analytics - Riot Games API Client
// Compliant Riot Games API client with rate limiting and caching

// RiotClient handles all Riot Games API interactions
type RiotClient struct {
	httpClient  *http.Client
	redis       *redis.Client
	config      *RiotClientConfig
	rateLimiter *RiotRateLimiter
}

// RiotClientConfig contains client configuration
type RiotClientConfig struct {
	APIKey           string        `json:"api_key"`
	BaseURL          string        `json:"base_url"`
	Timeout          time.Duration `json:"timeout"`
	MaxRetries       int           `json:"max_retries"`
	CacheEnabled     bool          `json:"cache_enabled"`
	CacheTTL         time.Duration `json:"cache_ttl"`
	RateLimitEnabled bool          `json:"rate_limit_enabled"`
	UserAgent        string        `json:"user_agent"`
	RequestsPerMin   int           `json:"requests_per_min"` // Riot API limit
	BurstLimit       int           `json:"burst_limit"`      // Short burst limit
}

// NewRiotClient creates new Riot Games API client
func NewRiotClient(redis *redis.Client, config *RiotClientConfig) *RiotClient {
	httpClient := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: false,
		},
	}

	return &RiotClient{
		httpClient:  httpClient,
		redis:       redis,
		config:      config,
		rateLimiter: NewRiotRateLimiter(redis, config.RequestsPerMin, config.BurstLimit),
	}
}

// GetSummonerByName retrieves summoner information by summoner name
func (r *RiotClient) GetSummonerByName(ctx context.Context, region, summonerName string) (*Summoner, error) {
	endpoint := fmt.Sprintf("/lol/summoner/v4/summoners/by-name/%s", summonerName)
	cacheKey := fmt.Sprintf("summoner:%s:%s", region, strings.ToLower(summonerName))

	// Try cache first
	if r.config.CacheEnabled {
		if cached, err := r.getCachedResponse(ctx, cacheKey); err == nil && cached != "" {
			var summoner Summoner
			if json.Unmarshal([]byte(cached), &summoner) == nil {
				return &summoner, nil
			}
		}
	}

	// Make API request
	response, err := r.makeRiotAPIRequest(ctx, region, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get summoner by name: %w", err)
	}

	var summoner Summoner
	if err := json.Unmarshal(response, &summoner); err != nil {
		return nil, fmt.Errorf("failed to unmarshal summoner response: %w", err)
	}

	// Cache the response
	if r.config.CacheEnabled {
		r.cacheResponse(ctx, cacheKey, response, r.config.CacheTTL)
	}

	return &summoner, nil
}

// GetSummonerByPUUID retrieves summoner by PUUID
func (r *RiotClient) GetSummonerByPUUID(ctx context.Context, region, puuid string) (*Summoner, error) {
	endpoint := fmt.Sprintf("/lol/summoner/v4/summoners/by-puuid/%s", puuid)
	cacheKey := fmt.Sprintf("summoner:puuid:%s:%s", region, puuid)

	if r.config.CacheEnabled {
		if cached, err := r.getCachedResponse(ctx, cacheKey); err == nil && cached != "" {
			var summoner Summoner
			if json.Unmarshal([]byte(cached), &summoner) == nil {
				return &summoner, nil
			}
		}
	}

	response, err := r.makeRiotAPIRequest(ctx, region, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get summoner by PUUID: %w", err)
	}

	var summoner Summoner
	if err := json.Unmarshal(response, &summoner); err != nil {
		return nil, fmt.Errorf("failed to unmarshal summoner response: %w", err)
	}

	if r.config.CacheEnabled {
		r.cacheResponse(ctx, cacheKey, response, r.config.CacheTTL)
	}

	return &summoner, nil
}

// GetRankedInfo retrieves ranked information for summoner
func (r *RiotClient) GetRankedInfo(ctx context.Context, region, summonerID string) ([]RankedEntry, error) {
	endpoint := fmt.Sprintf("/lol/league/v4/entries/by-summoner/%s", summonerID)
	cacheKey := fmt.Sprintf("ranked:%s:%s", region, summonerID)

	if r.config.CacheEnabled {
		if cached, err := r.getCachedResponse(ctx, cacheKey); err == nil && cached != "" {
			var entries []RankedEntry
			if json.Unmarshal([]byte(cached), &entries) == nil {
				return entries, nil
			}
		}
	}

	response, err := r.makeRiotAPIRequest(ctx, region, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get ranked info: %w", err)
	}

	var entries []RankedEntry
	if err := json.Unmarshal(response, &entries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ranked response: %w", err)
	}

	if r.config.CacheEnabled {
		r.cacheResponse(ctx, cacheKey, response, 10*time.Minute) // Shorter cache for ranked data
	}

	return entries, nil
}

// GetMatchHistory retrieves match history for summoner
func (r *RiotClient) GetMatchHistory(ctx context.Context, region, puuid string, start, count int) ([]string, error) {
	endpoint := fmt.Sprintf("/lol/match/v5/matches/by-puuid/%s/ids?start=%d&count=%d", puuid, start, count)
	cacheKey := fmt.Sprintf("matches:%s:%s:%d:%d", region, puuid, start, count)

	if r.config.CacheEnabled {
		if cached, err := r.getCachedResponse(ctx, cacheKey); err == nil && cached != "" {
			var matchIds []string
			if json.Unmarshal([]byte(cached), &matchIds) == nil {
				return matchIds, nil
			}
		}
	}

	response, err := r.makeRiotAPIRequest(ctx, region, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get match history: %w", err)
	}

	var matchIds []string
	if err := json.Unmarshal(response, &matchIds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal match history response: %w", err)
	}

	if r.config.CacheEnabled {
		r.cacheResponse(ctx, cacheKey, response, 5*time.Minute)
	}

	return matchIds, nil
}

// GetMatch retrieves detailed match information
func (r *RiotClient) GetMatch(ctx context.Context, region, matchID string) (*Match, error) {
	endpoint := fmt.Sprintf("/lol/match/v5/matches/%s", matchID)
	cacheKey := fmt.Sprintf("match:%s:%s", region, matchID)

	if r.config.CacheEnabled {
		if cached, err := r.getCachedResponse(ctx, cacheKey); err == nil && cached != "" {
			var match Match
			if json.Unmarshal([]byte(cached), &match) == nil {
				return &match, nil
			}
		}
	}

	response, err := r.makeRiotAPIRequest(ctx, region, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	var match Match
	if err := json.Unmarshal(response, &match); err != nil {
		return nil, fmt.Errorf("failed to unmarshal match response: %w", err)
	}

	if r.config.CacheEnabled {
		// Matches don't change, cache for longer
		r.cacheResponse(ctx, cacheKey, response, 24*time.Hour)
	}

	return &match, nil
}

// GetLiveGame retrieves current live game information
func (r *RiotClient) GetLiveGame(ctx context.Context, region, summonerID string) (*LiveGame, error) {
	endpoint := fmt.Sprintf("/lol/spectator/v4/active-games/by-summoner/%s", summonerID)

	// Don't cache live game data
	response, err := r.makeRiotAPIRequest(ctx, region, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get live game: %w", err)
	}

	var liveGame LiveGame
	if err := json.Unmarshal(response, &liveGame); err != nil {
		return nil, fmt.Errorf("failed to unmarshal live game response: %w", err)
	}

	return &liveGame, nil
}

// GetChampionMastery retrieves champion mastery for summoner
func (r *RiotClient) GetChampionMastery(ctx context.Context, region, summonerID string) ([]ChampionMastery, error) {
	endpoint := fmt.Sprintf("/lol/champion-mastery/v4/champion-masteries/by-summoner/%s", summonerID)
	cacheKey := fmt.Sprintf("mastery:%s:%s", region, summonerID)

	if r.config.CacheEnabled {
		if cached, err := r.getCachedResponse(ctx, cacheKey); err == nil && cached != "" {
			var masteries []ChampionMastery
			if json.Unmarshal([]byte(cached), &masteries) == nil {
				return masteries, nil
			}
		}
	}

	response, err := r.makeRiotAPIRequest(ctx, region, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get champion mastery: %w", err)
	}

	var masteries []ChampionMastery
	if err := json.Unmarshal(response, &masteries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal champion mastery response: %w", err)
	}

	if r.config.CacheEnabled {
		r.cacheResponse(ctx, cacheKey, response, 30*time.Minute)
	}

	return masteries, nil
}

// makeRiotAPIRequest makes HTTP request to Riot API with rate limiting
func (r *RiotClient) makeRiotAPIRequest(ctx context.Context, region, endpoint string) ([]byte, error) {
	// Check rate limits
	if r.config.RateLimitEnabled {
		allowed, waitTime, err := r.rateLimiter.CheckRateLimit(ctx)
		if err != nil {
			return nil, fmt.Errorf("rate limit check failed: %w", err)
		}
		if !allowed {
			return nil, fmt.Errorf("rate limit exceeded, retry after: %v", waitTime)
		}
	}

	// Build URL
	baseURL := r.getRegionalURL(region)
	url := baseURL + endpoint

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("X-Riot-Token", r.config.APIKey)
	req.Header.Set("User-Agent", r.config.UserAgent)
	req.Header.Set("Accept", "application/json")

	// Execute request with retries
	var response *http.Response
	var lastErr error

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		response, lastErr = r.httpClient.Do(req)
		if lastErr == nil && response.StatusCode < 500 {
			break
		}

		if response != nil {
			response.Body.Close()
		}

		if attempt < r.config.MaxRetries {
			// Exponential backoff
			backoff := time.Duration(1<<attempt) * time.Second
			time.Sleep(backoff)
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("request failed after %d attempts: %w", r.config.MaxRetries, lastErr)
	}
	defer response.Body.Close()

	// Handle response
	if response.StatusCode == 429 {
		// Rate limited by Riot
		return nil, fmt.Errorf("rate limited by Riot API")
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("API request failed with status: %d", response.StatusCode)
	}

	// Read response body
	body := make([]byte, 0, response.ContentLength)
	buffer := make([]byte, 4096)
	for {
		n, err := response.Body.Read(buffer)
		if n > 0 {
			body = append(body, buffer[:n]...)
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
	}

	return body, nil
}

// getRegionalURL returns the correct regional URL
func (r *RiotClient) getRegionalURL(region string) string {
	regionalMappings := map[string]string{
		"NA1":  "https://na1.api.riotgames.com",
		"EUW1": "https://euw1.api.riotgames.com",
		"EUN1": "https://eun1.api.riotgames.com",
		"KR":   "https://kr.api.riotgames.com",
		"JP1":  "https://jp1.api.riotgames.com",
		"BR1":  "https://br1.api.riotgames.com",
		"LA1":  "https://la1.api.riotgames.com",
		"LA2":  "https://la2.api.riotgames.com",
		"OC1":  "https://oc1.api.riotgames.com",
		"TR1":  "https://tr1.api.riotgames.com",
		"RU":   "https://ru.api.riotgames.com",
	}

	if url, exists := regionalMappings[region]; exists {
		return url
	}

	// Default to NA1 if region not found
	return regionalMappings["NA1"]
}

// getCachedResponse retrieves cached response from Redis
func (r *RiotClient) getCachedResponse(ctx context.Context, key string) (string, error) {
	cacheKey := fmt.Sprintf("riot:cache:%s", key)
	return r.redis.Get(ctx, cacheKey).Result()
}

// cacheResponse stores response in Redis cache
func (r *RiotClient) cacheResponse(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	cacheKey := fmt.Sprintf("riot:cache:%s", key)
	return r.redis.Set(ctx, cacheKey, string(data), ttl).Err()
}

// GetClientStats returns client usage statistics
func (r *RiotClient) GetClientStats(ctx context.Context) (*RiotClientStats, error) {
	stats := &RiotClientStats{}

	// Get rate limiter stats
	if rateLimiterStats, err := r.rateLimiter.GetStats(ctx); err == nil {
		stats.RequestsThisMinute = rateLimiterStats.RequestsThisMinute
		stats.RequestsToday = rateLimiterStats.RequestsToday
		stats.RateLimitHits = rateLimiterStats.RateLimitHits
	}

	// Get cache hit rate
	cachePattern := "riot:cache:*"
	keys, err := r.redis.Keys(ctx, cachePattern).Result()
	if err == nil {
		stats.CachedEntries = len(keys)
	}

	return stats, nil
}

// ClearCache clears all cached Riot API responses
func (r *RiotClient) ClearCache(ctx context.Context) error {
	pattern := "riot:cache:*"
	keys, err := r.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get cache keys: %w", err)
	}

	if len(keys) > 0 {
		return r.redis.Del(ctx, keys...).Err()
	}

	return nil
}

// DefaultRiotClientConfig returns default client configuration
func DefaultRiotClientConfig(apiKey string) *RiotClientConfig {
	return &RiotClientConfig{
		APIKey:           apiKey,
		BaseURL:          "https://na1.api.riotgames.com",
		Timeout:          30 * time.Second,
		MaxRetries:       3,
		CacheEnabled:     true,
		CacheTTL:         15 * time.Minute,
		RateLimitEnabled: true,
		UserAgent:        "Herald.lol/1.0 (Gaming Analytics Platform)",
		RequestsPerMin:   100, // Personal development key limit
		BurstLimit:       20,  // Short burst allowance
	}
}

// RiotClientStats contains client usage statistics
type RiotClientStats struct {
	RequestsThisMinute int `json:"requests_this_minute"`
	RequestsToday      int `json:"requests_today"`
	RateLimitHits      int `json:"rate_limit_hits"`
	CachedEntries      int `json:"cached_entries"`
}
