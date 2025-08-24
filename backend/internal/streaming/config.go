package streaming

import (
	"time"
)

// Herald.lol Gaming Analytics - Real-time Streaming Configuration
// Configuration settings for real-time gaming data streaming service

// StreamingConfig contains main streaming service configuration
type StreamingConfig struct {
	// WebSocket settings
	EnableWebSocket      bool          `json:"enable_websocket"`
	WebSocketReadBuffer  int           `json:"websocket_read_buffer"`
	WebSocketWriteBuffer int           `json:"websocket_write_buffer"`
	MaxMessageSize       int64         `json:"max_message_size"`
	PingPeriod           time.Duration `json:"ping_period"`
	PongWait             time.Duration `json:"pong_wait"`
	WriteWait            time.Duration `json:"write_wait"`

	// Connection management
	MaxConnections       int           `json:"max_connections"`
	MaxChannelsPerClient int           `json:"max_channels_per_client"`
	ClientTimeout        time.Duration `json:"client_timeout"`
	ClientUpdateInterval time.Duration `json:"client_update_interval"`
	ConnectionRateLimit  int           `json:"connection_rate_limit"` // per minute

	// Live match streaming
	EnableLiveMatches       bool          `json:"enable_live_matches"`
	LiveMatchUpdateInterval time.Duration `json:"live_match_update_interval"`
	LiveMatchScanInterval   time.Duration `json:"live_match_scan_interval"`
	LiveMatchTTL            time.Duration `json:"live_match_ttl"`
	MaxLiveMatches          int           `json:"max_live_matches"`

	// Player updates
	EnablePlayerUpdates  bool          `json:"enable_player_updates"`
	PlayerUpdateInterval time.Duration `json:"player_update_interval"`
	PlayerStatusCacheTTL time.Duration `json:"player_status_cache_ttl"`

	// Analytics streaming
	EnableAnalyticsStream   bool          `json:"enable_analytics_stream"`
	AnalyticsUpdateInterval time.Duration `json:"analytics_update_interval"`
	TrendAnalysisInterval   time.Duration `json:"trend_analysis_interval"`

	// Notifications
	EnableNotifications     bool          `json:"enable_notifications"`
	NotificationTTL         time.Duration `json:"notification_ttl"`
	MaxNotificationsPerUser int           `json:"max_notifications_per_user"`

	// Event processing
	EventQueueSize         int           `json:"event_queue_size"`
	EventWorkers           int           `json:"event_workers"`
	EventProcessingTimeout time.Duration `json:"event_processing_timeout"`
	MaxEventRetries        int           `json:"max_event_retries"`

	// Performance settings
	ChannelTTL            time.Duration `json:"channel_ttl"`
	CleanupInterval       time.Duration `json:"cleanup_interval"`
	MetricsUpdateInterval time.Duration `json:"metrics_update_interval"`

	// Gaming-specific settings
	RiotAPIUpdateFrequency time.Duration `json:"riot_api_update_frequency"`
	MatchHistoryDepth      int           `json:"match_history_depth"`
	RealtimeLatencyTarget  time.Duration `json:"realtime_latency_target"`

	// Security settings
	RequireAuth         bool     `json:"require_auth"`
	AllowedOrigins      []string `json:"allowed_origins"`
	RateLimitingEnabled bool     `json:"rate_limiting_enabled"`

	// Subscription-based limits
	SubscriptionConfigs map[string]*SubscriptionConfig `json:"subscription_configs"`
}

// GetDefaultStreamingConfig returns default configuration for Herald.lol streaming
func GetDefaultStreamingConfig() *StreamingConfig {
	return &StreamingConfig{
		// WebSocket settings optimized for gaming
		EnableWebSocket:      true,
		WebSocketReadBuffer:  4096,
		WebSocketWriteBuffer: 4096,
		MaxMessageSize:       32768, // 32KB for detailed gaming data
		PingPeriod:           54 * time.Second,
		PongWait:             60 * time.Second,
		WriteWait:            10 * time.Second,

		// Connection management for 1M+ concurrent users
		MaxConnections:       1000000, // 1M concurrent connections target
		MaxChannelsPerClient: 50,
		ClientTimeout:        5 * time.Minute,
		ClientUpdateInterval: 30 * time.Second,
		ConnectionRateLimit:  1000, // connections per minute per IP

		// Live match streaming for real-time gaming
		EnableLiveMatches:       true,
		LiveMatchUpdateInterval: 5 * time.Second,  // 5s updates for live matches
		LiveMatchScanInterval:   30 * time.Second, // Scan for new matches every 30s
		LiveMatchTTL:            2 * time.Hour,    // Keep inactive matches for 2 hours
		MaxLiveMatches:          100000,           // Support 100k concurrent live matches

		// Player updates for gaming status
		EnablePlayerUpdates:  true,
		PlayerUpdateInterval: 15 * time.Second, // 15s player status updates
		PlayerStatusCacheTTL: 5 * time.Minute,  // Cache status for 5 minutes

		// Analytics streaming for performance insights
		EnableAnalyticsStream:   true,
		AnalyticsUpdateInterval: 10 * time.Second, // 10s analytics updates
		TrendAnalysisInterval:   1 * time.Minute,  // Trend analysis every minute

		// Gaming notifications
		EnableNotifications:     true,
		NotificationTTL:         24 * time.Hour, // Keep notifications for 24 hours
		MaxNotificationsPerUser: 100,

		// High-throughput event processing
		EventQueueSize:         100000, // 100k event buffer
		EventWorkers:           20,     // 20 parallel workers
		EventProcessingTimeout: 30 * time.Second,
		MaxEventRetries:        3,

		// Performance optimized for gaming platform
		ChannelTTL:            1 * time.Hour,
		CleanupInterval:       15 * time.Minute,
		MetricsUpdateInterval: 30 * time.Second,

		// Gaming-specific settings
		RiotAPIUpdateFrequency: 10 * time.Second,       // Riot API polling frequency
		MatchHistoryDepth:      20,                     // Last 20 matches for context
		RealtimeLatencyTarget:  500 * time.Millisecond, // <500ms latency target

		// Security for gaming platform
		RequireAuth:         true,
		AllowedOrigins:      []string{"https://herald.lol", "https://app.herald.lol", "https://cdn.herald.lol"},
		RateLimitingEnabled: true,

		// Subscription configurations
		SubscriptionConfigs: GetSubscriptionStreamingConfig(),
	}
}

// GetStreamingConfigByEnvironment returns environment-specific configuration
func GetStreamingConfigByEnvironment(env string) *StreamingConfig {
	baseConfig := GetDefaultStreamingConfig()

	switch env {
	case "development":
		// Development settings
		baseConfig.MaxConnections = 1000
		baseConfig.MaxLiveMatches = 100
		baseConfig.EventQueueSize = 1000
		baseConfig.EventWorkers = 2
		baseConfig.LiveMatchUpdateInterval = 10 * time.Second
		baseConfig.AllowedOrigins = []string{"*"} // Allow all origins in dev

	case "staging":
		// Staging settings
		baseConfig.MaxConnections = 10000
		baseConfig.MaxLiveMatches = 1000
		baseConfig.EventQueueSize = 10000
		baseConfig.EventWorkers = 5
		baseConfig.LiveMatchUpdateInterval = 8 * time.Second

	case "production":
		// Production settings (use defaults optimized for 1M+ users)
		// Additional production-specific optimizations
		baseConfig.EventWorkers = 50       // More workers for production load
		baseConfig.EventQueueSize = 500000 // Larger queue for peak traffic

	default:
		// Use default configuration
	}

	return baseConfig
}

// Performance Targets for Herald.lol Gaming Platform

// StreamingPerformanceTargets contains performance targets for real-time streaming
type StreamingPerformanceTargets struct {
	// Latency targets
	MessageDeliveryLatency    time.Duration `json:"message_delivery_latency"`     // <500ms
	LiveMatchUpdateLatency    time.Duration `json:"live_match_update_latency"`    // <1s
	PlayerStatusUpdateLatency time.Duration `json:"player_status_update_latency"` // <2s
	AnalyticsUpdateLatency    time.Duration `json:"analytics_update_latency"`     // <3s

	// Throughput targets
	MessagesPerSecond    int `json:"messages_per_second"`    // 100k+ msg/s
	ConnectionsPerSecond int `json:"connections_per_second"` // 1k+ conn/s
	EventsPerSecond      int `json:"events_per_second"`      // 50k+ events/s

	// Reliability targets
	MessageDeliveryRate   float64 `json:"message_delivery_rate"`   // 99.9%
	ConnectionSuccessRate float64 `json:"connection_success_rate"` // 99.95%
	EventProcessingRate   float64 `json:"event_processing_rate"`   // 99.8%

	// Scalability targets
	ConcurrentConnections int     `json:"concurrent_connections"`  // 1M+
	ConcurrentLiveMatches int     `json:"concurrent_live_matches"` // 100k+
	PeakTrafficMultiplier float64 `json:"peak_traffic_multiplier"` // 3x normal load

	// Gaming-specific targets
	LiveMatchAccuracy     float64       `json:"live_match_accuracy"`      // 99.5%
	PlayerStatusAccuracy  float64       `json:"player_status_accuracy"`   // 99%
	RealTimeDataFreshness time.Duration `json:"real_time_data_freshness"` // <5s

	// Resource efficiency targets
	CPUUsageLimit         float64 `json:"cpu_usage_limit"`         // <70%
	MemoryUsageLimit      int64   `json:"memory_usage_limit"`      // <8GB per instance
	NetworkBandwidthLimit int64   `json:"network_bandwidth_limit"` // <1Gbps per instance
}

// GetHeraldStreamingTargets returns Herald.lol specific streaming performance targets
func GetHeraldStreamingTargets() *StreamingPerformanceTargets {
	return &StreamingPerformanceTargets{
		// Ultra-low latency for competitive gaming
		MessageDeliveryLatency:    500 * time.Millisecond,
		LiveMatchUpdateLatency:    1 * time.Second,
		PlayerStatusUpdateLatency: 2 * time.Second,
		AnalyticsUpdateLatency:    3 * time.Second,

		// High throughput for gaming platform scale
		MessagesPerSecond:    100000,
		ConnectionsPerSecond: 1000,
		EventsPerSecond:      50000,

		// Gaming-grade reliability
		MessageDeliveryRate:   0.999,  // 99.9%
		ConnectionSuccessRate: 0.9995, // 99.95%
		EventProcessingRate:   0.998,  // 99.8%

		// Massive gaming community scalability
		ConcurrentConnections: 1000000, // 1M concurrent users
		ConcurrentLiveMatches: 100000,  // 100k live matches
		PeakTrafficMultiplier: 3.0,     // Handle 3x normal load

		// Gaming accuracy requirements
		LiveMatchAccuracy:     0.995, // 99.5%
		PlayerStatusAccuracy:  0.99,  // 99%
		RealTimeDataFreshness: 5 * time.Second,

		// Efficient resource usage
		CPUUsageLimit:         0.7,                    // 70% CPU max
		MemoryUsageLimit:      8 * 1024 * 1024 * 1024, // 8GB
		NetworkBandwidthLimit: 1000 * 1024 * 1024,     // 1Gbps
	}
}

// StreamingFeatureConfig contains feature-specific configurations
type StreamingFeatureConfig struct {
	// Live match features
	LiveMatchFeatures *LiveMatchFeatureConfig `json:"live_match_features"`

	// Player tracking features
	PlayerTrackingFeatures *PlayerTrackingFeatureConfig `json:"player_tracking_features"`

	// Analytics features
	AnalyticsFeatures *AnalyticsFeatureConfig `json:"analytics_features"`

	// Notification features
	NotificationFeatures *NotificationFeatureConfig `json:"notification_features"`
}

// LiveMatchFeatureConfig contains live match streaming configuration
type LiveMatchFeatureConfig struct {
	EnableRealTimeStats     bool `json:"enable_real_time_stats"`
	EnablePositionTracking  bool `json:"enable_position_tracking"`
	EnableEventPrediction   bool `json:"enable_event_prediction"`
	EnableTeamFightAnalysis bool `json:"enable_team_fight_analysis"`
	EnableObjectiveAlerts   bool `json:"enable_objective_alerts"`
	EnablePerformanceAlerts bool `json:"enable_performance_alerts"`

	// Update intervals
	StatsUpdateInterval    time.Duration `json:"stats_update_interval"`
	PositionUpdateInterval time.Duration `json:"position_update_interval"`
	EventCheckInterval     time.Duration `json:"event_check_interval"`

	// Thresholds
	SignificantEventThreshold float64 `json:"significant_event_threshold"`
	PerformanceAlertThreshold float64 `json:"performance_alert_threshold"`
}

// PlayerTrackingFeatureConfig contains player tracking configuration
type PlayerTrackingFeatureConfig struct {
	EnableStatusTracking      bool `json:"enable_status_tracking"`
	EnableRankTracking        bool `json:"enable_rank_tracking"`
	EnableMatchTracking       bool `json:"enable_match_tracking"`
	EnableAchievementTracking bool `json:"enable_achievement_tracking"`
	EnableFriendUpdates       bool `json:"enable_friend_updates"`
	EnablePerformanceTracking bool `json:"enable_performance_tracking"`

	// Tracking intervals
	StatusCheckInterval      time.Duration `json:"status_check_interval"`
	RankCheckInterval        time.Duration `json:"rank_check_interval"`
	PerformanceCheckInterval time.Duration `json:"performance_check_interval"`

	// Privacy settings
	RespectPrivacySettings bool `json:"respect_privacy_settings"`
	AnonymizeData          bool `json:"anonymize_data"`
}

// AnalyticsFeatureConfig contains analytics streaming configuration
type AnalyticsFeatureConfig struct {
	EnableTrendAlerts         bool `json:"enable_trend_alerts"`
	EnableMilestoneAlerts     bool `json:"enable_milestone_alerts"`
	EnablePerformanceAlerts   bool `json:"enable_performance_alerts"`
	EnableMetaUpdates         bool `json:"enable_meta_updates"`
	EnablePredictiveAnalytics bool `json:"enable_predictive_analytics"`
	EnableComparisonAlerts    bool `json:"enable_comparison_alerts"`

	// Analysis intervals
	TrendAnalysisInterval  time.Duration `json:"trend_analysis_interval"`
	MetaAnalysisInterval   time.Duration `json:"meta_analysis_interval"`
	MilestoneCheckInterval time.Duration `json:"milestone_check_interval"`

	// Alert thresholds
	TrendSignificanceThreshold float64 `json:"trend_significance_threshold"`
	PerformanceChangeThreshold float64 `json:"performance_change_threshold"`
	MilestoneProgressThreshold float64 `json:"milestone_progress_threshold"`
}

// NotificationFeatureConfig contains notification system configuration
type NotificationFeatureConfig struct {
	EnablePushNotifications    bool `json:"enable_push_notifications"`
	EnableEmailNotifications   bool `json:"enable_email_notifications"`
	EnableInAppNotifications   bool `json:"enable_in_app_notifications"`
	EnableWebPushNotifications bool `json:"enable_web_push_notifications"`
	EnableCustomSounds         bool `json:"enable_custom_sounds"`
	EnableRichNotifications    bool `json:"enable_rich_notifications"`

	// Notification limits
	MaxNotificationsPerHour int `json:"max_notifications_per_hour"`
	MaxNotificationsPerDay  int `json:"max_notifications_per_day"`

	// Delivery settings
	BatchDelivery         bool          `json:"batch_delivery"`
	DeliveryRetryAttempts int           `json:"delivery_retry_attempts"`
	DeliveryTimeout       time.Duration `json:"delivery_timeout"`

	// Priority handling
	HighPriorityDeliveryTime   time.Duration `json:"high_priority_delivery_time"`
	NormalPriorityDeliveryTime time.Duration `json:"normal_priority_delivery_time"`
	LowPriorityDeliveryTime    time.Duration `json:"low_priority_delivery_time"`
}

// GetDefaultFeatureConfig returns default feature configuration for Herald.lol
func GetDefaultFeatureConfig() *StreamingFeatureConfig {
	return &StreamingFeatureConfig{
		LiveMatchFeatures: &LiveMatchFeatureConfig{
			EnableRealTimeStats:     true,
			EnablePositionTracking:  true,
			EnableEventPrediction:   true,
			EnableTeamFightAnalysis: true,
			EnableObjectiveAlerts:   true,
			EnablePerformanceAlerts: true,

			StatsUpdateInterval:    5 * time.Second,
			PositionUpdateInterval: 2 * time.Second,
			EventCheckInterval:     1 * time.Second,

			SignificantEventThreshold: 0.8,
			PerformanceAlertThreshold: 0.9,
		},

		PlayerTrackingFeatures: &PlayerTrackingFeatureConfig{
			EnableStatusTracking:      true,
			EnableRankTracking:        true,
			EnableMatchTracking:       true,
			EnableAchievementTracking: true,
			EnableFriendUpdates:       true,
			EnablePerformanceTracking: true,

			StatusCheckInterval:      30 * time.Second,
			RankCheckInterval:        5 * time.Minute,
			PerformanceCheckInterval: 1 * time.Minute,

			RespectPrivacySettings: true,
			AnonymizeData:          false,
		},

		AnalyticsFeatures: &AnalyticsFeatureConfig{
			EnableTrendAlerts:         true,
			EnableMilestoneAlerts:     true,
			EnablePerformanceAlerts:   true,
			EnableMetaUpdates:         true,
			EnablePredictiveAnalytics: true,
			EnableComparisonAlerts:    true,

			TrendAnalysisInterval:  1 * time.Minute,
			MetaAnalysisInterval:   15 * time.Minute,
			MilestoneCheckInterval: 5 * time.Minute,

			TrendSignificanceThreshold: 0.15, // 15% change
			PerformanceChangeThreshold: 0.10, // 10% change
			MilestoneProgressThreshold: 0.90, // 90% progress
		},

		NotificationFeatures: &NotificationFeatureConfig{
			EnablePushNotifications:    true,
			EnableEmailNotifications:   true,
			EnableInAppNotifications:   true,
			EnableWebPushNotifications: true,
			EnableCustomSounds:         true,
			EnableRichNotifications:    true,

			MaxNotificationsPerHour: 50,
			MaxNotificationsPerDay:  200,

			BatchDelivery:         true,
			DeliveryRetryAttempts: 3,
			DeliveryTimeout:       10 * time.Second,

			HighPriorityDeliveryTime:   1 * time.Second,
			NormalPriorityDeliveryTime: 5 * time.Second,
			LowPriorityDeliveryTime:    30 * time.Second,
		},
	}
}

// Regional configuration for global Herald.lol deployment

// RegionalConfig contains region-specific streaming configuration
type RegionalConfig struct {
	Region                  string        `json:"region"`
	RiotAPIEndpoint         string        `json:"riot_api_endpoint"`
	StreamingServerEndpoint string        `json:"streaming_server_endpoint"`
	MaxConnections          int           `json:"max_connections"`
	LatencyTarget           time.Duration `json:"latency_target"`
	DataCenterLocations     []string      `json:"data_center_locations"`
	PrimaryLanguage         string        `json:"primary_language"`
	SupportedLanguages      []string      `json:"supported_languages"`
	ComplianceRequirements  []string      `json:"compliance_requirements"`
}

// GetRegionalConfigs returns region-specific configurations for global deployment
func GetRegionalConfigs() map[string]*RegionalConfig {
	return map[string]*RegionalConfig{
		"na": {
			Region:                  "North America",
			RiotAPIEndpoint:         "americas.api.riotgames.com",
			StreamingServerEndpoint: "stream-na.herald.lol",
			MaxConnections:          400000, // 400k concurrent for NA
			LatencyTarget:           300 * time.Millisecond,
			DataCenterLocations:     []string{"us-east-1", "us-west-2", "ca-central-1"},
			PrimaryLanguage:         "en-US",
			SupportedLanguages:      []string{"en-US", "es-ES", "pt-BR"},
			ComplianceRequirements:  []string{"COPPA", "CCPA"},
		},
		"euw": {
			Region:                  "Europe West",
			RiotAPIEndpoint:         "europe.api.riotgames.com",
			StreamingServerEndpoint: "stream-euw.herald.lol",
			MaxConnections:          300000, // 300k concurrent for EUW
			LatencyTarget:           250 * time.Millisecond,
			DataCenterLocations:     []string{"eu-west-1", "eu-central-1"},
			PrimaryLanguage:         "en-GB",
			SupportedLanguages:      []string{"en-GB", "de-DE", "fr-FR", "es-ES", "it-IT"},
			ComplianceRequirements:  []string{"GDPR", "ePrivacy"},
		},
		"kr": {
			Region:                  "Korea",
			RiotAPIEndpoint:         "asia.api.riotgames.com",
			StreamingServerEndpoint: "stream-kr.herald.lol",
			MaxConnections:          200000, // 200k concurrent for KR
			LatencyTarget:           200 * time.Millisecond,
			DataCenterLocations:     []string{"ap-northeast-2"},
			PrimaryLanguage:         "ko-KR",
			SupportedLanguages:      []string{"ko-KR"},
			ComplianceRequirements:  []string{"PIPA"},
		},
		"cn": {
			Region:                  "China",
			RiotAPIEndpoint:         "asia.api.riotgames.com",
			StreamingServerEndpoint: "stream-cn.herald.lol",
			MaxConnections:          150000,                 // 150k concurrent for CN
			LatencyTarget:           400 * time.Millisecond, // Higher latency due to restrictions
			DataCenterLocations:     []string{"cn-north-1", "cn-northwest-1"},
			PrimaryLanguage:         "zh-CN",
			SupportedLanguages:      []string{"zh-CN"},
			ComplianceRequirements:  []string{"Cybersecurity Law", "Data Security Law"},
		},
	}
}
