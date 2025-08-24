package streaming

import (
	"time"

	"github.com/gorilla/websocket"
)

// Herald.lol Gaming Analytics - Real-time Streaming Models
// Data structures for real-time gaming data streaming

// Connection and Client Models

// ClientConnection represents a WebSocket connection from a client
type ClientConnection struct {
	ID           string          `json:"id"`
	UserID       string          `json:"user_id"`
	PlayerPUUID  string          `json:"player_puuid"`
	Connection   *websocket.Conn `json:"-"`
	Channels     map[string]bool `json:"channels"`
	LastPing     time.Time       `json:"last_ping"`
	Connected    bool            `json:"connected"`
	JoinedAt     time.Time       `json:"joined_at"`
	MessageCount int             `json:"message_count"`
	UserAgent    string          `json:"user_agent,omitempty"`
	IPAddress    string          `json:"ip_address,omitempty"`
}

// StreamChannel represents a streaming channel with subscribers
type StreamChannel struct {
	Name         string                       `json:"name"`
	Type         string                       `json:"type"` // live_match, player, analytics, notification
	Subscribers  map[string]*ClientConnection `json:"subscribers"`
	CreatedAt    time.Time                    `json:"created_at"`
	MessageCount int                          `json:"message_count"`
	LastMessage  time.Time                    `json:"last_message"`
}

// Message Models

// StreamMessage represents a message sent to clients via WebSocket
type StreamMessage struct {
	Type      string      `json:"type"`
	Channel   string      `json:"channel,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
	MessageID string      `json:"message_id,omitempty"`
}

// ClientMessage represents a message received from a client
type ClientMessage struct {
	Type      string                 `json:"type"`
	Channel   string                 `json:"channel,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
}

// Live Match Streaming Models

// LiveMatchTracker tracks a live match for real-time updates
type LiveMatchTracker struct {
	MatchID      string                      `json:"match_id"`
	GameID       string                      `json:"game_id"`
	StartTime    time.Time                   `json:"start_time"`
	GameTime     int                         `json:"game_time"`     // seconds
	CurrentState string                      `json:"current_state"` // in_progress, paused, finished
	Participants map[string]*LiveParticipant `json:"participants"`
	Teams        map[string]*LiveTeam        `json:"teams"`
	Subscribers  map[string]bool             `json:"subscribers"`
	Events       []*LiveMatchEvent           `json:"events"`
	LastUpdate   time.Time                   `json:"last_update"`
	Region       string                      `json:"region"`
	Queue        string                      `json:"queue"`
	GameMode     string                      `json:"game_mode"`
}

// LiveParticipant represents a participant in a live match
type LiveParticipant struct {
	PlayerPUUID  string       `json:"player_puuid"`
	SummonerName string       `json:"summoner_name"`
	ChampionName string       `json:"champion_name"`
	ChampionID   int          `json:"champion_id"`
	TeamID       int          `json:"team_id"`
	Position     string       `json:"position"`
	Spells       []int        `json:"spells"`
	Runes        *LiveRuneSet `json:"runes"`

	// Real-time stats
	Level      int         `json:"level"`
	Kills      int         `json:"kills"`
	Deaths     int         `json:"deaths"`
	Assists    int         `json:"assists"`
	CS         int         `json:"cs"`
	Gold       int         `json:"gold"`
	Items      []int       `json:"items"`
	Position2D *Position2D `json:"position_2d,omitempty"`

	// Performance metrics
	KDA               float64 `json:"kda"`
	CSPerMinute       float64 `json:"cs_per_minute"`
	GoldPerMinute     float64 `json:"gold_per_minute"`
	KillParticipation float64 `json:"kill_participation"`

	// Streaming metadata
	LastUpdate  time.Time `json:"last_update"`
	UpdateCount int       `json:"update_count"`
}

// LiveTeam represents a team in a live match
type LiveTeam struct {
	TeamID  int      `json:"team_id"`
	Side    string   `json:"side"`    // blue, red
	Players []string `json:"players"` // Player PUUIDs

	// Team stats
	Kills   int `json:"kills"`
	Deaths  int `json:"deaths"`
	Assists int `json:"assists"`
	Gold    int `json:"gold"`
	Towers  int `json:"towers"`
	Dragons int `json:"dragons"`
	Barons  int `json:"barons"`
	Heralds int `json:"heralds"`

	// Objectives
	DragonSoul  string `json:"dragon_soul,omitempty"`
	ElderDragon bool   `json:"elder_dragon"`
	Baron       bool   `json:"baron"`

	// Performance
	TeamKDA        float64 `json:"team_kda"`
	GoldLead       int     `json:"gold_lead"`
	ObjectiveScore int     `json:"objective_score"`
}

// LiveRuneSet represents rune configuration
type LiveRuneSet struct {
	PrimaryTree   int   `json:"primary_tree"`
	SecondaryTree int   `json:"secondary_tree"`
	Runes         []int `json:"runes"`
	StatRunes     []int `json:"stat_runes"`
}

// Position2D represents a 2D position on the map
type Position2D struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Live Match Event Models

// LiveMatchEvent represents an event that occurred in a live match
type LiveMatchEvent struct {
	EventID      string                 `json:"event_id"`
	Type         string                 `json:"type"`
	GameTime     int                    `json:"game_time"`
	Description  string                 `json:"description"`
	Participants []string               `json:"participants"` // Player PUUIDs involved
	TeamID       int                    `json:"team_id,omitempty"`
	Position     *Position2D            `json:"position,omitempty"`
	Impact       string                 `json:"impact"` // low, medium, high, critical
	Details      map[string]interface{} `json:"details,omitempty"`
	Timestamp    time.Time              `json:"timestamp"`
}

// Player Update Models

// PlayerUpdate represents a real-time player update
type PlayerUpdate struct {
	PlayerPUUID string                 `json:"player_puuid"`
	UpdateType  string                 `json:"update_type"` // status, rank, match_start, match_end
	Timestamp   time.Time              `json:"timestamp"`
	Data        map[string]interface{} `json:"data"`

	// Specific update types
	StatusUpdate      *PlayerStatusUpdate      `json:"status_update,omitempty"`
	RankUpdate        *PlayerRankUpdate        `json:"rank_update,omitempty"`
	MatchUpdate       *PlayerMatchUpdate       `json:"match_update,omitempty"`
	AchievementUpdate *PlayerAchievementUpdate `json:"achievement_update,omitempty"`
}

// PlayerStatusUpdate represents a player status change
type PlayerStatusUpdate struct {
	OnlineStatus    string    `json:"online_status"` // online, in_game, away, offline
	CurrentActivity string    `json:"current_activity"`
	LastSeen        time.Time `json:"last_seen"`
	GameClient      string    `json:"game_client,omitempty"`
}

// PlayerRankUpdate represents a rank change
type PlayerRankUpdate struct {
	Queue       string       `json:"queue"`
	OldRank     string       `json:"old_rank"`
	NewRank     string       `json:"new_rank"`
	OldLP       int          `json:"old_lp"`
	NewLP       int          `json:"new_lp"`
	LPChange    int          `json:"lp_change"`
	PromoSeries *PromoSeries `json:"promo_series,omitempty"`
}

// PlayerMatchUpdate represents match start/end updates
type PlayerMatchUpdate struct {
	MatchID     string                   `json:"match_id"`
	Status      string                   `json:"status"` // started, ended
	Champion    string                   `json:"champion,omitempty"`
	Role        string                   `json:"role,omitempty"`
	Queue       string                   `json:"queue"`
	Result      string                   `json:"result,omitempty"` // victory, defeat
	Duration    int                      `json:"duration,omitempty"`
	Performance *MatchPerformanceSummary `json:"performance,omitempty"`
}

// PlayerAchievementUpdate represents new achievement unlocked
type PlayerAchievementUpdate struct {
	AchievementID string    `json:"achievement_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	Rarity        string    `json:"rarity"`
	Points        int       `json:"points"`
	UnlockedAt    time.Time `json:"unlocked_at"`
}

// PromoSeries represents ranked promotion series
type PromoSeries struct {
	Target   string   `json:"target"`
	Wins     int      `json:"wins"`
	Losses   int      `json:"losses"`
	Progress []string `json:"progress"` // W, L, N (not played)
}

// MatchPerformanceSummary represents a brief match performance summary
type MatchPerformanceSummary struct {
	KDA     float64 `json:"kda"`
	Kills   int     `json:"kills"`
	Deaths  int     `json:"deaths"`
	Assists int     `json:"assists"`
	CS      int     `json:"cs"`
	Damage  int     `json:"damage"`
	Vision  int     `json:"vision"`
	Rating  float64 `json:"rating"`
}

// Notification Models

// Notification represents a real-time notification
type Notification struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"` // achievement, rank_up, friend_online, match_ready, etc.
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Icon      string                 `json:"icon,omitempty"`
	Priority  string                 `json:"priority"` // low, normal, high, urgent
	Category  string                 `json:"category"` // gaming, social, system, promotion
	Data      map[string]interface{} `json:"data,omitempty"`
	ActionURL string                 `json:"action_url,omitempty"`
	ExpiresAt *time.Time             `json:"expires_at,omitempty"`
	CreatedAt time.Time              `json:"created_at"`

	// Targeting
	UserIDs      []string `json:"user_ids,omitempty"`
	PlayerPUUIDs []string `json:"player_puuids,omitempty"`
	Regions      []string `json:"regions,omitempty"`

	// Delivery tracking
	DeliveredTo map[string]time.Time `json:"delivered_to,omitempty"`
	ReadBy      map[string]time.Time `json:"read_by,omitempty"`
}

// Analytics Streaming Models

// AnalyticsUpdate represents a real-time analytics update
type AnalyticsUpdate struct {
	Type      string      `json:"type"`     // performance_change, trend_alert, milestone, comparison
	Category  string      `json:"category"` // player, champion, meta, rank
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`

	// Specific update types
	TrendAlert        *TrendAlert              `json:"trend_alert,omitempty"`
	Milestone         *MilestoneUpdate         `json:"milestone,omitempty"`
	PerformanceChange *PerformanceChangeUpdate `json:"performance_change,omitempty"`
	MetaUpdate        *MetaAnalyticsUpdate     `json:"meta_update,omitempty"`
}

// TrendAlert represents a trend-based alert
type TrendAlert struct {
	PlayerPUUID string  `json:"player_puuid"`
	Metric      string  `json:"metric"`    // winrate, kda, cs_per_min, etc.
	Direction   string  `json:"direction"` // increasing, decreasing
	Magnitude   string  `json:"magnitude"` // slight, moderate, significant, dramatic
	OldValue    float64 `json:"old_value"`
	NewValue    float64 `json:"new_value"`
	Change      float64 `json:"change"`
	TimeSpan    string  `json:"time_span"`
	Confidence  float64 `json:"confidence"` // 0-1
}

// MilestoneUpdate represents a milestone achievement in analytics
type MilestoneUpdate struct {
	PlayerPUUID   string                 `json:"player_puuid"`
	MilestoneType string                 `json:"milestone_type"`
	Description   string                 `json:"description"`
	Value         float64                `json:"value"`
	Target        float64                `json:"target"`
	Progress      float64                `json:"progress"` // 0-1
	Rarity        string                 `json:"rarity"`
	Category      string                 `json:"category"`
	Details       map[string]interface{} `json:"details,omitempty"`
}

// PerformanceChangeUpdate represents a significant performance change
type PerformanceChangeUpdate struct {
	PlayerPUUID       string  `json:"player_puuid"`
	Champion          string  `json:"champion,omitempty"`
	Role              string  `json:"role,omitempty"`
	PerformanceMetric string  `json:"performance_metric"`
	OldRating         float64 `json:"old_rating"`
	NewRating         float64 `json:"new_rating"`
	Change            float64 `json:"change"`
	ChangeType        string  `json:"change_type"`  // improvement, decline, plateau
	Significance      string  `json:"significance"` // minor, moderate, major
	TimeFrame         string  `json:"time_frame"`
	SampleSize        int     `json:"sample_size"`
}

// MetaAnalyticsUpdate represents meta game analytics updates
type MetaAnalyticsUpdate struct {
	Region     string `json:"region"`
	Queue      string `json:"queue"`
	Patch      string `json:"patch"`
	UpdateType string `json:"update_type"` // champion_tier_change, item_popularity, role_meta

	ChampionUpdates []ChampionMetaUpdate `json:"champion_updates,omitempty"`
	ItemUpdates     []ItemMetaUpdate     `json:"item_updates,omitempty"`
	RoleUpdates     []RoleMetaUpdate     `json:"role_updates,omitempty"`
}

// ChampionMetaUpdate represents champion meta changes
type ChampionMetaUpdate struct {
	ChampionName   string  `json:"champion_name"`
	Role           string  `json:"role"`
	OldTier        string  `json:"old_tier"`
	NewTier        string  `json:"new_tier"`
	OldWinRate     float64 `json:"old_win_rate"`
	NewWinRate     float64 `json:"new_win_rate"`
	OldPickRate    float64 `json:"old_pick_rate"`
	NewPickRate    float64 `json:"new_pick_rate"`
	TrendDirection string  `json:"trend_direction"` // rising, falling, stable
}

// ItemMetaUpdate represents item meta changes
type ItemMetaUpdate struct {
	ItemName      string   `json:"item_name"`
	OldPopularity float64  `json:"old_popularity"`
	NewPopularity float64  `json:"new_popularity"`
	WinRateImpact float64  `json:"win_rate_impact"`
	Champions     []string `json:"champions"` // Champions commonly building this item
}

// RoleMetaUpdate represents role meta changes
type RoleMetaUpdate struct {
	Role              string   `json:"role"`
	PopularChampions  []string `json:"popular_champions"`
	EmergingPicks     []string `json:"emerging_picks"`
	DeciningPicks     []string `json:"declining_picks"`
	AverageGameLength int      `json:"average_game_length"`
	KeyItems          []string `json:"key_items"`
}

// Event Processing Models

// StreamEvent represents an event in the streaming system
type StreamEvent struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Source      string      `json:"source"`
	MatchID     string      `json:"match_id,omitempty"`
	PlayerPUUID string      `json:"player_puuid,omitempty"`
	Data        interface{} `json:"data"`
	Timestamp   time.Time   `json:"timestamp"`
	Priority    int         `json:"priority"` // 1-10, higher is more important
	Processed   bool        `json:"processed"`
	ProcessedAt *time.Time  `json:"processed_at,omitempty"`
	Error       string      `json:"error,omitempty"`
}

// EventProcessor interface for processing different types of streaming events
type EventProcessor interface {
	Process(event *StreamEvent) error
	GetEventType() string
	GetPriority() int
}

// Statistics and Monitoring Models

// StreamingStats represents current streaming service statistics
type StreamingStats struct {
	ConnectedClients  int           `json:"connected_clients"`
	ActiveChannels    int           `json:"active_channels"`
	LiveMatches       int           `json:"live_matches"`
	EventsProcessed   int64         `json:"events_processed"`
	MessagesDelivered int64         `json:"messages_delivered"`
	DroppedEvents     int64         `json:"dropped_events"`
	FailedEvents      int64         `json:"failed_events"`
	AverageLatency    time.Duration `json:"average_latency"`
	PeakConcurrent    int           `json:"peak_concurrent"`
	Uptime            time.Duration `json:"uptime"`

	// Performance metrics
	CPUUsage        float64 `json:"cpu_usage"`
	MemoryUsage     int64   `json:"memory_usage"` // bytes
	NetworkBytesIn  int64   `json:"network_bytes_in"`
	NetworkBytesOut int64   `json:"network_bytes_out"`

	// Gaming-specific metrics
	LiveMatchUpdates  int64 `json:"live_match_updates"`
	PlayerUpdates     int64 `json:"player_updates"`
	AnalyticsUpdates  int64 `json:"analytics_updates"`
	NotificationsSent int64 `json:"notifications_sent"`
}

// StreamingMetrics handles metrics collection for the streaming service
type StreamingMetrics struct {
	StartTime         time.Time `json:"start_time"`
	ConnectionCount   int64     `json:"connection_count"`
	EventsProcessed   int64     `json:"events_processed"`
	MessagesDelivered int64     `json:"messages_delivered"`
	DroppedEvents     int64     `json:"dropped_events"`
	FailedEvents      int64     `json:"failed_events"`

	// Latency tracking
	LatencySum     time.Duration `json:"latency_sum"`
	LatencyCount   int64         `json:"latency_count"`
	AverageLatency time.Duration `json:"average_latency"`

	// Peak tracking
	PeakConnections int `json:"peak_connections"`
	PeakChannels    int `json:"peak_channels"`
	PeakLiveMatches int `json:"peak_live_matches"`

	mutex sync.RWMutex `json:"-"`
}

// NewStreamingMetrics creates a new metrics collector
func NewStreamingMetrics() *StreamingMetrics {
	return &StreamingMetrics{
		StartTime: time.Now(),
	}
}

// IncrementConnections increments the connection count
func (m *StreamingMetrics) IncrementConnections() {
	m.mutex.Lock()
	m.ConnectionCount++
	if int(m.ConnectionCount) > m.PeakConnections {
		m.PeakConnections = int(m.ConnectionCount)
	}
	m.mutex.Unlock()
}

// DecrementConnections decrements the connection count
func (m *StreamingMetrics) DecrementConnections() {
	m.mutex.Lock()
	if m.ConnectionCount > 0 {
		m.ConnectionCount--
	}
	m.mutex.Unlock()
}

// IncrementEvents increments the events processed count
func (m *StreamingMetrics) IncrementEvents() {
	m.mutex.Lock()
	m.EventsProcessed++
	m.mutex.Unlock()
}

// IncrementProcessedEvents increments successfully processed events
func (m *StreamingMetrics) IncrementProcessedEvents() {
	m.IncrementEvents()
}

// IncrementDroppedEvents increments dropped events count
func (m *StreamingMetrics) IncrementDroppedEvents() {
	m.mutex.Lock()
	m.DroppedEvents++
	m.mutex.Unlock()
}

// IncrementFailedEvents increments failed events count
func (m *StreamingMetrics) IncrementFailedEvents() {
	m.mutex.Lock()
	m.FailedEvents++
	m.mutex.Unlock()
}

// IncrementMessages increments messages delivered count
func (m *StreamingMetrics) IncrementMessages(count int) {
	m.mutex.Lock()
	m.MessagesDelivered += int64(count)
	m.mutex.Unlock()
}

// RecordLatency records message latency
func (m *StreamingMetrics) RecordLatency(latency time.Duration) {
	m.mutex.Lock()
	m.LatencySum += latency
	m.LatencyCount++
	m.AverageLatency = m.LatencySum / time.Duration(m.LatencyCount)
	m.mutex.Unlock()
}

// UpdateStats updates current statistics
func (m *StreamingMetrics) UpdateStats(clients, channels, liveMatches int) {
	m.mutex.Lock()
	if clients > m.PeakConnections {
		m.PeakConnections = clients
	}
	if channels > m.PeakChannels {
		m.PeakChannels = channels
	}
	if liveMatches > m.PeakLiveMatches {
		m.PeakLiveMatches = liveMatches
	}
	m.mutex.Unlock()
}

// Configuration Models

// SubscriptionConfig represents subscription-based streaming limits
type SubscriptionConfig struct {
	MaxConnections        int  `json:"max_connections"`
	MaxChannelsPerClient  int  `json:"max_channels_per_client"`
	MessageRateLimit      int  `json:"message_rate_limit"` // per minute
	MaxMessageSize        int  `json:"max_message_size"`   // bytes
	LiveMatchAccess       bool `json:"live_match_access"`
	AnalyticsStreamAccess bool `json:"analytics_stream_access"`
	PriorityDelivery      bool `json:"priority_delivery"`
	HistoricalData        bool `json:"historical_data"`
	CustomNotifications   bool `json:"custom_notifications"`
}

// GetSubscriptionStreamingConfig returns streaming limits by subscription tier
func GetSubscriptionStreamingConfig() map[string]*SubscriptionConfig {
	return map[string]*SubscriptionConfig{
		"free": {
			MaxConnections:        1,
			MaxChannelsPerClient:  5,
			MessageRateLimit:      60,   // 1 per second
			MaxMessageSize:        1024, // 1KB
			LiveMatchAccess:       false,
			AnalyticsStreamAccess: false,
			PriorityDelivery:      false,
			HistoricalData:        false,
			CustomNotifications:   false,
		},
		"premium": {
			MaxConnections:        3,
			MaxChannelsPerClient:  15,
			MessageRateLimit:      300,  // 5 per second
			MaxMessageSize:        5120, // 5KB
			LiveMatchAccess:       true,
			AnalyticsStreamAccess: true,
			PriorityDelivery:      false,
			HistoricalData:        true,
			CustomNotifications:   false,
		},
		"pro": {
			MaxConnections:        10,
			MaxChannelsPerClient:  50,
			MessageRateLimit:      1800,  // 30 per second
			MaxMessageSize:        10240, // 10KB
			LiveMatchAccess:       true,
			AnalyticsStreamAccess: true,
			PriorityDelivery:      true,
			HistoricalData:        true,
			CustomNotifications:   true,
		},
		"enterprise": {
			MaxConnections:        -1,    // unlimited
			MaxChannelsPerClient:  -1,    // unlimited
			MessageRateLimit:      -1,    // unlimited
			MaxMessageSize:        51200, // 50KB
			LiveMatchAccess:       true,
			AnalyticsStreamAccess: true,
			PriorityDelivery:      true,
			HistoricalData:        true,
			CustomNotifications:   true,
		},
	}
}
