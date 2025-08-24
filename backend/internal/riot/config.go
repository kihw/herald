package riot

import (
	"time"
)

// Herald.lol Gaming Analytics - Riot API Configuration
// Configuration structures and defaults for Riot Games API integration

// RiotAPIConfig contains global Riot API configuration
type RiotAPIConfig struct {
	// API Keys
	PersonalAPIKey   string `json:"personal_api_key"`
	ProductionAPIKey string `json:"production_api_key,omitempty"`

	// Rate Limiting
	UsePersonalKey    bool `json:"use_personal_key"` // true = 100req/2min, false = production limits
	RequestsPerMinute int  `json:"requests_per_minute"`
	BurstLimit        int  `json:"burst_limit"`

	// Caching
	CacheEnabled     bool          `json:"cache_enabled"`
	SummonerCacheTTL time.Duration `json:"summoner_cache_ttl"`
	MatchCacheTTL    time.Duration `json:"match_cache_ttl"`
	RankedCacheTTL   time.Duration `json:"ranked_cache_ttl"`
	MasteryCacheTTL  time.Duration `json:"mastery_cache_ttl"`

	// Request Settings
	RequestTimeout time.Duration `json:"request_timeout"`
	MaxRetries     int           `json:"max_retries"`
	RetryBackoff   time.Duration `json:"retry_backoff"`

	// Client Settings
	UserAgent         string `json:"user_agent"`
	EnableCompression bool   `json:"enable_compression"`

	// Regional Settings
	DefaultRegion     string            `json:"default_region"`
	RegionalEndpoints map[string]string `json:"regional_endpoints"`

	// Analytics Settings
	EnableAnalytics   bool          `json:"enable_analytics"`
	AnalyticsDepth    int           `json:"analytics_depth"` // Number of matches to analyze
	CacheAnalytics    bool          `json:"cache_analytics"`
	AnalyticsCacheTTL time.Duration `json:"analytics_cache_ttl"`

	// Compliance Settings
	RespectRiotLimits bool `json:"respect_riot_limits"`
	LogAPIRequests    bool `json:"log_api_requests"`
	MonitorUsage      bool `json:"monitor_usage"`
}

// RegionConfig contains region-specific configuration
type RegionConfig struct {
	RegionCode          string  `json:"region_code"`           // NA1, EUW1, etc.
	DisplayName         string  `json:"display_name"`          // North America, Europe West, etc.
	BaseURL             string  `json:"base_url"`              // Regional API base URL
	Timezone            string  `json:"timezone"`              // Region timezone
	Language            string  `json:"language"`              // Primary language
	Enabled             bool    `json:"enabled"`               // Whether region is enabled
	RateLimitMultiplier float64 `json:"rate_limit_multiplier"` // Regional rate limit adjustment
}

// QueueConfig contains queue-specific configuration
type QueueConfig struct {
	QueueID     int    `json:"queue_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"` // "Ranked", "Normal", "ARAM", etc.
	IsRanked    bool   `json:"is_ranked"`
	TeamSize    int    `json:"team_size"`
	Enabled     bool   `json:"enabled"`
	Priority    int    `json:"priority"` // Higher priority queues get more analysis
}

// ChampionConfig contains champion-specific configuration
type ChampionConfig struct {
	ChampionID    int      `json:"champion_id"`
	Name          string   `json:"name"`
	Title         string   `json:"title"`
	Roles         []string `json:"roles"`
	Tags          []string `json:"tags"`
	Difficulty    int      `json:"difficulty"`
	Enabled       bool     `json:"enabled"`
	AnalyticsTags []string `json:"analytics_tags"` // Tags for analytics categorization
}

// DefaultRiotAPIConfig returns default API configuration
func DefaultRiotAPIConfig() *RiotAPIConfig {
	return &RiotAPIConfig{
		// Rate Limiting (Personal Development Key limits)
		UsePersonalKey:    true,
		RequestsPerMinute: 50, // Conservative limit under 100/2min
		BurstLimit:        20, // Allow short bursts

		// Caching
		CacheEnabled:     true,
		SummonerCacheTTL: 15 * time.Minute,
		MatchCacheTTL:    24 * time.Hour,   // Matches don't change
		RankedCacheTTL:   10 * time.Minute, // Ranked data changes frequently
		MasteryCacheTTL:  30 * time.Minute,

		// Request Settings
		RequestTimeout: 30 * time.Second,
		MaxRetries:     3,
		RetryBackoff:   2 * time.Second,

		// Client Settings
		UserAgent:         "Herald.lol/1.0 (Gaming Analytics Platform)",
		EnableCompression: true,

		// Regional Settings
		DefaultRegion: "NA1",
		RegionalEndpoints: map[string]string{
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
		},

		// Analytics Settings
		EnableAnalytics:   true,
		AnalyticsDepth:    20, // Analyze last 20 matches
		CacheAnalytics:    true,
		AnalyticsCacheTTL: 30 * time.Minute,

		// Compliance Settings
		RespectRiotLimits: true,
		LogAPIRequests:    true,
		MonitorUsage:      true,
	}
}

// GetSupportedRegions returns list of supported League of Legends regions
func GetSupportedRegions() []RegionConfig {
	return []RegionConfig{
		{
			RegionCode:          "NA1",
			DisplayName:         "North America",
			BaseURL:             "https://na1.api.riotgames.com",
			Timezone:            "America/Los_Angeles",
			Language:            "en_US",
			Enabled:             true,
			RateLimitMultiplier: 1.0,
		},
		{
			RegionCode:          "EUW1",
			DisplayName:         "Europe West",
			BaseURL:             "https://euw1.api.riotgames.com",
			Timezone:            "Europe/London",
			Language:            "en_GB",
			Enabled:             true,
			RateLimitMultiplier: 1.0,
		},
		{
			RegionCode:          "EUN1",
			DisplayName:         "Europe Nordic & East",
			BaseURL:             "https://eun1.api.riotgames.com",
			Timezone:            "Europe/Stockholm",
			Language:            "en_GB",
			Enabled:             true,
			RateLimitMultiplier: 1.0,
		},
		{
			RegionCode:          "KR",
			DisplayName:         "Korea",
			BaseURL:             "https://kr.api.riotgames.com",
			Timezone:            "Asia/Seoul",
			Language:            "ko_KR",
			Enabled:             true,
			RateLimitMultiplier: 0.8, // Higher traffic region
		},
		{
			RegionCode:          "JP1",
			DisplayName:         "Japan",
			BaseURL:             "https://jp1.api.riotgames.com",
			Timezone:            "Asia/Tokyo",
			Language:            "ja_JP",
			Enabled:             true,
			RateLimitMultiplier: 0.9,
		},
		{
			RegionCode:          "BR1",
			DisplayName:         "Brazil",
			BaseURL:             "https://br1.api.riotgames.com",
			Timezone:            "America/Sao_Paulo",
			Language:            "pt_BR",
			Enabled:             true,
			RateLimitMultiplier: 1.1,
		},
		{
			RegionCode:          "LA1",
			DisplayName:         "Latin America North",
			BaseURL:             "https://la1.api.riotgames.com",
			Timezone:            "America/Mexico_City",
			Language:            "es_MX",
			Enabled:             true,
			RateLimitMultiplier: 1.2,
		},
		{
			RegionCode:          "LA2",
			DisplayName:         "Latin America South",
			BaseURL:             "https://la2.api.riotgames.com",
			Timezone:            "America/Santiago",
			Language:            "es_CL",
			Enabled:             true,
			RateLimitMultiplier: 1.2,
		},
		{
			RegionCode:          "OC1",
			DisplayName:         "Oceania",
			BaseURL:             "https://oc1.api.riotgames.com",
			Timezone:            "Australia/Sydney",
			Language:            "en_AU",
			Enabled:             true,
			RateLimitMultiplier: 1.3, // Lower traffic region
		},
		{
			RegionCode:          "TR1",
			DisplayName:         "Turkey",
			BaseURL:             "https://tr1.api.riotgames.com",
			Timezone:            "Europe/Istanbul",
			Language:            "tr_TR",
			Enabled:             true,
			RateLimitMultiplier: 1.1,
		},
		{
			RegionCode:          "RU",
			DisplayName:         "Russia",
			BaseURL:             "https://ru.api.riotgames.com",
			Timezone:            "Europe/Moscow",
			Language:            "ru_RU",
			Enabled:             true,
			RateLimitMultiplier: 1.0,
		},
	}
}

// GetSupportedQueues returns list of supported ranked queues
func GetSupportedQueues() []QueueConfig {
	return []QueueConfig{
		{
			QueueID:     420,
			Name:        "Ranked Solo/Duo",
			Description: "5v5 Ranked Solo/Duo games on Summoner's Rift",
			Category:    "Ranked",
			IsRanked:    true,
			TeamSize:    5,
			Enabled:     true,
			Priority:    10, // Highest priority
		},
		{
			QueueID:     440,
			Name:        "Ranked Flex 5v5",
			Description: "5v5 Ranked Flex games on Summoner's Rift",
			Category:    "Ranked",
			IsRanked:    true,
			TeamSize:    5,
			Enabled:     true,
			Priority:    8,
		},
		{
			QueueID:     430,
			Name:        "Normal Blind Pick",
			Description: "5v5 Normal Blind Pick games on Summoner's Rift",
			Category:    "Normal",
			IsRanked:    false,
			TeamSize:    5,
			Enabled:     true,
			Priority:    5,
		},
		{
			QueueID:     400,
			Name:        "Normal Draft Pick",
			Description: "5v5 Normal Draft Pick games on Summoner's Rift",
			Category:    "Normal",
			IsRanked:    false,
			TeamSize:    5,
			Enabled:     true,
			Priority:    6,
		},
		{
			QueueID:     450,
			Name:        "ARAM",
			Description: "5v5 ARAM games on Howling Abyss",
			Category:    "ARAM",
			IsRanked:    false,
			TeamSize:    5,
			Enabled:     true,
			Priority:    3,
		},
		{
			QueueID:     1700,
			Name:        "Arena",
			Description: "2v2v2v2 Arena games",
			Category:    "Arena",
			IsRanked:    false,
			TeamSize:    2,
			Enabled:     true,
			Priority:    2,
		},
	}
}

// GetChampionRoles returns standard champion roles
func GetChampionRoles() []string {
	return []string{
		"TOP",
		"JUNGLE",
		"MIDDLE",
		"BOTTOM",
		"UTILITY", // Support
	}
}

// GetRankedTiers returns ranked tier hierarchy
func GetRankedTiers() []string {
	return []string{
		"UNRANKED",
		"IRON",
		"BRONZE",
		"SILVER",
		"GOLD",
		"PLATINUM",
		"EMERALD",
		"DIAMOND",
		"MASTER",
		"GRANDMASTER",
		"CHALLENGER",
	}
}

// GetRankedRanks returns ranked rank hierarchy within tiers
func GetRankedRanks() []string {
	return []string{
		"IV",
		"III",
		"II",
		"I",
	}
}

// ValidateRegion checks if region code is supported
func ValidateRegion(region string) bool {
	supportedRegions := GetSupportedRegions()
	for _, r := range supportedRegions {
		if r.RegionCode == region && r.Enabled {
			return true
		}
	}
	return false
}

// ValidateQueueID checks if queue ID is supported
func ValidateQueueID(queueID int) bool {
	supportedQueues := GetSupportedQueues()
	for _, q := range supportedQueues {
		if q.QueueID == queueID && q.Enabled {
			return true
		}
	}
	return false
}

// GetRegionConfig returns configuration for specific region
func GetRegionConfig(regionCode string) *RegionConfig {
	supportedRegions := GetSupportedRegions()
	for _, r := range supportedRegions {
		if r.RegionCode == regionCode {
			return &r
		}
	}
	return nil
}

// GetQueueConfig returns configuration for specific queue
func GetQueueConfig(queueID int) *QueueConfig {
	supportedQueues := GetSupportedQueues()
	for _, q := range supportedQueues {
		if q.QueueID == queueID {
			return &q
		}
	}
	return nil
}

// IsRankedQueue checks if queue ID is a ranked queue
func IsRankedQueue(queueID int) bool {
	queueConfig := GetQueueConfig(queueID)
	return queueConfig != nil && queueConfig.IsRanked
}

// GetQueuePriority returns priority level for queue (higher = more important)
func GetQueuePriority(queueID int) int {
	queueConfig := GetQueueConfig(queueID)
	if queueConfig != nil {
		return queueConfig.Priority
	}
	return 1 // Default low priority
}
