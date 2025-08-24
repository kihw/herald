package summoner

import "time"

// Herald.lol Gaming Analytics - Summoner Service Configuration
// Configuration settings and defaults for summoner analytics service

// DefaultSummonerServiceConfig returns default configuration for summoner service
func DefaultSummonerServiceConfig() *SummonerServiceConfig {
	return &SummonerServiceConfig{
		// Analysis settings
		MaxMatchesAnalyzed:   100,
		DefaultAnalysisDepth: 20,
		CacheAnalyticsTTL:    30 * time.Minute,
		CacheSummonerTTL:     15 * time.Minute,

		// Performance settings
		AnalysisTimeout:       2 * time.Minute, // <5s target but allow buffer for complex analysis
		MaxConcurrentAnalysis: 10,
		EnableProgressTrack:   true,

		// Feature flags
		EnableLiveGameAnalysis: true,
		EnablePredictions:      true,
		EnableComparisons:      true,
		EnableCoaching:         true,

		// Queue preferences (prioritized by importance)
		PrioritizedQueues: []int{
			420, // Ranked Solo/Duo (highest priority)
			440, // Ranked Flex 5v5
			430, // Normal Blind Pick
			400, // Normal Draft Pick
			450, // ARAM
		},

		// Queue weights for analysis
		QueueWeights: map[int]float64{
			420: 1.0, // Ranked Solo/Duo - full weight
			440: 0.8, // Ranked Flex - high weight
			430: 0.4, // Normal Blind - medium weight
			400: 0.5, // Normal Draft - medium weight
			450: 0.3, // ARAM - lower weight
		},

		// Analysis limits by subscription tier
		TierAnalysisLimits: map[string]int{
			"free":       10,  // 10 matches for free users
			"premium":    25,  // 25 matches for premium users
			"pro":        50,  // 50 matches for pro users
			"enterprise": 100, // 100 matches for enterprise users
		},
	}
}

// GetAnalysisProfileConfig returns analysis configuration by profile type
func GetAnalysisProfileConfig() map[string]*AnalysisProfileConfig {
	return map[string]*AnalysisProfileConfig{
		"basic": {
			Name:               "Basic Analysis",
			MatchLimit:         10,
			IncludeLiveGame:    false,
			IncludePredictions: false,
			IncludeComparisons: false,
			IncludeCoaching:    false,
			ProcessingPriority: 1,
			CacheDuration:      60 * time.Minute,
			EstimatedTime:      "5-10 seconds",
			Features: []string{
				"Core metrics calculation",
				"Basic champion statistics",
				"Recent performance summary",
			},
		},
		"standard": {
			Name:               "Standard Analysis",
			MatchLimit:         25,
			IncludeLiveGame:    true,
			IncludePredictions: false,
			IncludeComparisons: true,
			IncludeCoaching:    false,
			ProcessingPriority: 2,
			CacheDuration:      30 * time.Minute,
			EstimatedTime:      "10-20 seconds",
			Features: []string{
				"Detailed metrics analysis",
				"Role performance breakdown",
				"Champion mastery analysis",
				"Live game information",
				"Basic comparisons",
			},
		},
		"detailed": {
			Name:               "Detailed Analysis",
			MatchLimit:         50,
			IncludeLiveGame:    true,
			IncludePredictions: true,
			IncludeComparisons: true,
			IncludeCoaching:    true,
			ProcessingPriority: 3,
			CacheDuration:      15 * time.Minute,
			EstimatedTime:      "20-45 seconds",
			Features: []string{
				"Comprehensive performance analysis",
				"AI-powered insights and coaching",
				"Performance predictions",
				"Detailed trend analysis",
				"Live game predictions",
				"Advanced comparisons",
			},
		},
		"professional": {
			Name:               "Professional Analysis",
			MatchLimit:         100,
			IncludeLiveGame:    true,
			IncludePredictions: true,
			IncludeComparisons: true,
			IncludeCoaching:    true,
			ProcessingPriority: 5,
			CacheDuration:      10 * time.Minute,
			EstimatedTime:      "45-90 seconds",
			Features: []string{
				"Complete professional-grade analysis",
				"Advanced AI coaching and insights",
				"Detailed performance forecasting",
				"Team composition analysis",
				"Meta analysis and recommendations",
				"Professional training plans",
				"Skill gap analysis",
			},
		},
	}
}

// AnalysisProfileConfig contains configuration for analysis profiles
type AnalysisProfileConfig struct {
	Name               string        `json:"name"`
	MatchLimit         int           `json:"match_limit"`
	IncludeLiveGame    bool          `json:"include_live_game"`
	IncludePredictions bool          `json:"include_predictions"`
	IncludeComparisons bool          `json:"include_comparisons"`
	IncludeCoaching    bool          `json:"include_coaching"`
	ProcessingPriority int           `json:"processing_priority"`
	CacheDuration      time.Duration `json:"cache_duration"`
	EstimatedTime      string        `json:"estimated_time"`
	Features           []string      `json:"features"`
}

// GetSubscriptionTierFeatures returns available features by subscription tier
func GetSubscriptionTierFeatures() map[string]*TierFeatures {
	return map[string]*TierFeatures{
		"free": {
			Tier:               "Free",
			MaxAnalysisPerDay:  5,
			MaxMatchesAnalyzed: 10,
			AnalysisProfiles:   []string{"basic"},
			Features: TierFeatureFlags{
				LiveGameAnalysis:   false,
				AIInsights:         false,
				PredictiveAnalysis: false,
				Comparisons:        false,
				CoachingAdvice:     false,
				TrendAnalysis:      false,
				ExportData:         false,
				PrioritySupport:    false,
			},
			RateLimit: &TierRateLimit{
				RequestsPerMinute:  10,
				RequestsPerHour:    60,
				ConcurrentRequests: 2,
			},
		},
		"premium": {
			Tier:               "Premium",
			MaxAnalysisPerDay:  25,
			MaxMatchesAnalyzed: 25,
			AnalysisProfiles:   []string{"basic", "standard"},
			Features: TierFeatureFlags{
				LiveGameAnalysis:   true,
				AIInsights:         true,
				PredictiveAnalysis: false,
				Comparisons:        true,
				CoachingAdvice:     false,
				TrendAnalysis:      true,
				ExportData:         true,
				PrioritySupport:    false,
			},
			RateLimit: &TierRateLimit{
				RequestsPerMinute:  25,
				RequestsPerHour:    150,
				ConcurrentRequests: 3,
			},
		},
		"pro": {
			Tier:               "Pro",
			MaxAnalysisPerDay:  100,
			MaxMatchesAnalyzed: 50,
			AnalysisProfiles:   []string{"basic", "standard", "detailed"},
			Features: TierFeatureFlags{
				LiveGameAnalysis:   true,
				AIInsights:         true,
				PredictiveAnalysis: true,
				Comparisons:        true,
				CoachingAdvice:     true,
				TrendAnalysis:      true,
				ExportData:         true,
				PrioritySupport:    true,
			},
			RateLimit: &TierRateLimit{
				RequestsPerMinute:  50,
				RequestsPerHour:    500,
				ConcurrentRequests: 5,
			},
		},
		"enterprise": {
			Tier:               "Enterprise",
			MaxAnalysisPerDay:  -1, // Unlimited
			MaxMatchesAnalyzed: 100,
			AnalysisProfiles:   []string{"basic", "standard", "detailed", "professional"},
			Features: TierFeatureFlags{
				LiveGameAnalysis:   true,
				AIInsights:         true,
				PredictiveAnalysis: true,
				Comparisons:        true,
				CoachingAdvice:     true,
				TrendAnalysis:      true,
				ExportData:         true,
				PrioritySupport:    true,
			},
			RateLimit: &TierRateLimit{
				RequestsPerMinute:  100,
				RequestsPerHour:    1000,
				ConcurrentRequests: 10,
			},
		},
	}
}

// TierFeatures contains features available for each subscription tier
type TierFeatures struct {
	Tier               string           `json:"tier"`
	MaxAnalysisPerDay  int              `json:"max_analysis_per_day"` // -1 for unlimited
	MaxMatchesAnalyzed int              `json:"max_matches_analyzed"`
	AnalysisProfiles   []string         `json:"analysis_profiles"`
	Features           TierFeatureFlags `json:"features"`
	RateLimit          *TierRateLimit   `json:"rate_limit"`
}

// TierFeatureFlags contains feature flags for subscription tiers
type TierFeatureFlags struct {
	LiveGameAnalysis   bool `json:"live_game_analysis"`
	AIInsights         bool `json:"ai_insights"`
	PredictiveAnalysis bool `json:"predictive_analysis"`
	Comparisons        bool `json:"comparisons"`
	CoachingAdvice     bool `json:"coaching_advice"`
	TrendAnalysis      bool `json:"trend_analysis"`
	ExportData         bool `json:"export_data"`
	PrioritySupport    bool `json:"priority_support"`
}

// TierRateLimit contains rate limiting for subscription tiers
type TierRateLimit struct {
	RequestsPerMinute  int `json:"requests_per_minute"`
	RequestsPerHour    int `json:"requests_per_hour"`
	ConcurrentRequests int `json:"concurrent_requests"`
}

// GetRegionLatencyConfig returns expected latency for different regions
func GetRegionLatencyConfig() map[string]*RegionLatency {
	return map[string]*RegionLatency{
		"NA1": {
			Region:              "NA1",
			ExpectedLatencyMs:   150,
			RiotAPILatencyMs:    100,
			ProcessingLatencyMs: 50,
			ReliabilityScore:    0.99,
		},
		"EUW1": {
			Region:              "EUW1",
			ExpectedLatencyMs:   180,
			RiotAPILatencyMs:    120,
			ProcessingLatencyMs: 60,
			ReliabilityScore:    0.98,
		},
		"KR": {
			Region:              "KR",
			ExpectedLatencyMs:   250,
			RiotAPILatencyMs:    200,
			ProcessingLatencyMs: 50,
			ReliabilityScore:    0.97,
		},
		"JP1": {
			Region:              "JP1",
			ExpectedLatencyMs:   220,
			RiotAPILatencyMs:    170,
			ProcessingLatencyMs: 50,
			ReliabilityScore:    0.98,
		},
	}
}

// RegionLatency contains latency expectations for regions
type RegionLatency struct {
	Region              string  `json:"region"`
	ExpectedLatencyMs   int     `json:"expected_latency_ms"`
	RiotAPILatencyMs    int     `json:"riot_api_latency_ms"`
	ProcessingLatencyMs int     `json:"processing_latency_ms"`
	ReliabilityScore    float64 `json:"reliability_score"`
}

// GetCacheStrategy returns caching strategy configuration
func GetCacheStrategy() *CacheStrategyConfig {
	return &CacheStrategyConfig{
		SummonerInfo: &CacheConfig{
			TTL:          15 * time.Minute,
			RefreshAfter: 10 * time.Minute,
			MaxSize:      10000,
		},
		RankedInfo: &CacheConfig{
			TTL:          5 * time.Minute,
			RefreshAfter: 3 * time.Minute,
			MaxSize:      10000,
		},
		MatchData: &CacheConfig{
			TTL:          24 * time.Hour, // Matches don't change
			RefreshAfter: 12 * time.Hour,
			MaxSize:      50000,
		},
		Analytics: &CacheConfig{
			TTL:          30 * time.Minute,
			RefreshAfter: 20 * time.Minute,
			MaxSize:      5000,
		},
		LiveGame: &CacheConfig{
			TTL:          30 * time.Second, // Very short for live data
			RefreshAfter: 15 * time.Second,
			MaxSize:      1000,
		},
		ChampionMastery: &CacheConfig{
			TTL:          30 * time.Minute,
			RefreshAfter: 20 * time.Minute,
			MaxSize:      10000,
		},
	}
}

// CacheStrategyConfig contains caching configuration
type CacheStrategyConfig struct {
	SummonerInfo    *CacheConfig `json:"summoner_info"`
	RankedInfo      *CacheConfig `json:"ranked_info"`
	MatchData       *CacheConfig `json:"match_data"`
	Analytics       *CacheConfig `json:"analytics"`
	LiveGame        *CacheConfig `json:"live_game"`
	ChampionMastery *CacheConfig `json:"champion_mastery"`
}

// CacheConfig contains individual cache configuration
type CacheConfig struct {
	TTL          time.Duration `json:"ttl"`           // Time to live
	RefreshAfter time.Duration `json:"refresh_after"` // Background refresh threshold
	MaxSize      int           `json:"max_size"`      // Maximum cache entries
}

// GetPerformanceTargets returns performance targets for the service
func GetPerformanceTargets() *PerformanceTargets {
	return &PerformanceTargets{
		AnalysisTimeTarget: map[string]time.Duration{
			"basic":        3 * time.Second,
			"standard":     10 * time.Second,
			"detailed":     30 * time.Second,
			"professional": 60 * time.Second,
		},
		ConcurrentUserTarget: 1000,
		ThroughputTarget: map[string]int{
			"requests_per_second": 100,
			"analysis_per_minute": 300,
		},
		AvailabilityTarget: 0.999, // 99.9%
		ErrorRateTarget:    0.001, // 0.1%
		CacheHitRateTarget: 0.80,  // 80%
	}
}

// PerformanceTargets contains service performance targets
type PerformanceTargets struct {
	AnalysisTimeTarget   map[string]time.Duration `json:"analysis_time_target"`
	ConcurrentUserTarget int                      `json:"concurrent_user_target"`
	ThroughputTarget     map[string]int           `json:"throughput_target"`
	AvailabilityTarget   float64                  `json:"availability_target"`
	ErrorRateTarget      float64                  `json:"error_rate_target"`
	CacheHitRateTarget   float64                  `json:"cache_hit_rate_target"`
}

// GetMonitoringConfig returns monitoring configuration
func GetMonitoringConfig() *MonitoringConfig {
	return &MonitoringConfig{
		MetricsEnabled:      true,
		TracingEnabled:      true,
		LoggingLevel:        "info",
		HealthCheckInterval: 30 * time.Second,
		AlertThresholds: &AlertThresholds{
			HighLatency:     5 * time.Second,
			HighErrorRate:   0.05, // 5%
			LowCacheHitRate: 0.50, // 50%
			HighMemoryUsage: 0.85, // 85%
			HighCPUUsage:    0.80, // 80%
		},
		MetricsRetention: 7 * 24 * time.Hour, // 7 days
	}
}

// MonitoringConfig contains monitoring configuration
type MonitoringConfig struct {
	MetricsEnabled      bool             `json:"metrics_enabled"`
	TracingEnabled      bool             `json:"tracing_enabled"`
	LoggingLevel        string           `json:"logging_level"`
	HealthCheckInterval time.Duration    `json:"health_check_interval"`
	AlertThresholds     *AlertThresholds `json:"alert_thresholds"`
	MetricsRetention    time.Duration    `json:"metrics_retention"`
}

// AlertThresholds contains alerting thresholds
type AlertThresholds struct {
	HighLatency     time.Duration `json:"high_latency"`
	HighErrorRate   float64       `json:"high_error_rate"`
	LowCacheHitRate float64       `json:"low_cache_hit_rate"`
	HighMemoryUsage float64       `json:"high_memory_usage"`
	HighCPUUsage    float64       `json:"high_cpu_usage"`
}

// GetFeatureFlags returns current feature flags
func GetFeatureFlags() *FeatureFlags {
	return &FeatureFlags{
		EnableBetaFeatures:      false,
		EnableExperimentalAI:    false,
		EnableAdvancedAnalytics: true,
		EnableRealTimeUpdates:   true,
		EnableComparisons:       true,
		EnablePredictions:       true,
		EnableCoaching:          true,
		EnableTeamAnalysis:      true,
		EnableMetaAnalysis:      true,
		EnableAutoRefresh:       true,
	}
}

// FeatureFlags contains global feature flags
type FeatureFlags struct {
	EnableBetaFeatures      bool `json:"enable_beta_features"`
	EnableExperimentalAI    bool `json:"enable_experimental_ai"`
	EnableAdvancedAnalytics bool `json:"enable_advanced_analytics"`
	EnableRealTimeUpdates   bool `json:"enable_real_time_updates"`
	EnableComparisons       bool `json:"enable_comparisons"`
	EnablePredictions       bool `json:"enable_predictions"`
	EnableCoaching          bool `json:"enable_coaching"`
	EnableTeamAnalysis      bool `json:"enable_team_analysis"`
	EnableMetaAnalysis      bool `json:"enable_meta_analysis"`
	EnableAutoRefresh       bool `json:"enable_auto_refresh"`
}
