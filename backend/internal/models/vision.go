package models

import (
	"time"
)

// WardPlacement represents a ward placement event in a match
type WardPlacement struct {
	ID              string    `json:"id" db:"id"`
	MatchID         string    `json:"match_id" db:"match_id"`
	PlayerID        string    `json:"player_id" db:"player_id"`
	
	// Ward Details
	WardType        string    `json:"ward_type" db:"ward_type"` // "YELLOW", "CONTROL", "BLUE_TRINKET", "FARSIGHT"
	WardID          int       `json:"ward_id" db:"ward_id"`
	
	// Position on Map
	X               int       `json:"x" db:"x"`
	Y               int       `json:"y" db:"y"`
	MapSide         string    `json:"map_side" db:"map_side"` // "BLUE", "RED"
	
	// Timing Information
	Timestamp       int       `json:"timestamp" db:"timestamp"` // Game time in seconds
	GamePhase       string    `json:"game_phase" db:"game_phase"` // "early", "mid", "late"
	
	// Strategic Information
	Zone            string    `json:"zone" db:"zone"` // "jungle", "river", "lane", "baron", "dragon"
	Strategic       bool      `json:"strategic" db:"strategic"` // High-value location
	
	// Ward Lifecycle
	Duration        int       `json:"duration" db:"duration"` // How long ward lasted (seconds)
	KilledBy        string    `json:"killed_by" db:"killed_by"` // Player ID who killed it
	KilledAt        int       `json:"killed_at" db:"killed_at"` // Game time when killed
	Expired         bool      `json:"expired" db:"expired"` // True if ward expired naturally
	
	// Context
	TeamGold        int       `json:"team_gold" db:"team_gold"` // Team's total gold when placed
	ObjectiveState  string    `json:"objective_state" db:"objective_state"` // "dragon_spawning", "baron_up", etc.
	
	// Metadata
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// WardKill represents a ward destruction event
type WardKill struct {
	ID              string    `json:"id" db:"id"`
	MatchID         string    `json:"match_id" db:"match_id"`
	KillerID        string    `json:"killer_id" db:"killer_id"`
	
	// Ward Information
	WardType        string    `json:"ward_type" db:"ward_type"`
	WardOwnerID     string    `json:"ward_owner_id" db:"ward_owner_id"`
	
	// Position
	X               int       `json:"x" db:"x"`
	Y               int       `json:"y" db:"y"`
	Zone            string    `json:"zone" db:"zone"`
	
	// Timing
	Timestamp       int       `json:"timestamp" db:"timestamp"`
	GamePhase       string    `json:"game_phase" db:"game_phase"`
	
	// Strategic Value
	Strategic       bool      `json:"strategic" db:"strategic"`
	VisionDenied    int       `json:"vision_denied" db:"vision_denied"` // Estimated vision value denied
	
	// Context
	NearbyAllies    int       `json:"nearby_allies" db:"nearby_allies"`
	NearbyEnemies   int       `json:"nearby_enemies" db:"nearby_enemies"`
	SafetyLevel     string    `json:"safety_level" db:"safety_level"` // "safe", "risky", "dangerous"
	
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// VisionEvent represents any vision-related event during a match
type VisionEvent struct {
	ID              string    `json:"id" db:"id"`
	MatchID         string    `json:"match_id" db:"match_id"`
	PlayerID        string    `json:"player_id" db:"player_id"`
	
	// Event Type
	EventType       string    `json:"event_type" db:"event_type"` // "WARD_PLACED", "WARD_KILLED", "STEALTH_WARD_PLACED", etc.
	
	// Position and Timing
	X               int       `json:"x" db:"x"`
	Y               int       `json:"y" db:"y"`
	Timestamp       int       `json:"timestamp" db:"timestamp"`
	
	// Additional Data
	WardType        string    `json:"ward_type" db:"ward_type"`
	ItemID          int       `json:"item_id" db:"item_id"`
	
	// Strategic Context
	Zone            string    `json:"zone" db:"zone"`
	Strategic       bool      `json:"strategic" db:"strategic"`
	ObjectivePhase  string    `json:"objective_phase" db:"objective_phase"`
	
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// VisionHeatmapData represents aggregated heatmap data for visualization
type VisionHeatmapData struct {
	ID              string    `json:"id" db:"id"`
	PlayerID        string    `json:"player_id" db:"player_id"`
	TimeRange       string    `json:"time_range" db:"time_range"`
	
	// Heatmap Configuration
	WardType        string    `json:"ward_type" db:"ward_type"`
	MapSide         string    `json:"map_side" db:"map_side"`
	Position        string    `json:"position" db:"position"`
	Champion        string    `json:"champion" db:"champion"`
	
	// Heatmap Points (JSON stored)
	DataPoints      string    `json:"data_points" db:"data_points"` // JSON array of HeatmapPoint
	ZoneIntensity   string    `json:"zone_intensity" db:"zone_intensity"` // JSON map of zone->frequency
	
	// Statistics
	TotalPlacements int       `json:"total_placements" db:"total_placements"`
	Coverage        float64   `json:"coverage" db:"coverage"`
	EfficiencyScore float64   `json:"efficiency_score" db:"efficiency_score"`
	
	// Metadata
	GeneratedAt     time.Time `json:"generated_at" db:"generated_at"`
	LastUpdated     time.Time `json:"last_updated" db:"last_updated"`
}

// PlayerVisionStats represents aggregated vision statistics for a player
type PlayerVisionStats struct {
	ID              string    `json:"id" db:"id"`
	PlayerID        string    `json:"player_id" db:"player_id"`
	TimeRange       string    `json:"time_range" db:"time_range"`
	Position        string    `json:"position" db:"position"`
	
	// Basic Vision Metrics
	TotalMatches         int       `json:"total_matches" db:"total_matches"`
	AverageVisionScore   float64   `json:"average_vision_score" db:"average_vision_score"`
	MedianVisionScore    float64   `json:"median_vision_score" db:"median_vision_score"`
	BestVisionScore      float64   `json:"best_vision_score" db:"best_vision_score"`
	WorstVisionScore     float64   `json:"worst_vision_score" db:"worst_vision_score"`
	VisionScoreStdDev    float64   `json:"vision_score_std_dev" db:"vision_score_std_dev"`
	
	// Ward Statistics
	AverageWardsPlaced   float64   `json:"average_wards_placed" db:"average_wards_placed"`
	AverageWardsKilled   float64   `json:"average_wards_killed" db:"average_wards_killed"`
	ControlWardsPerGame  float64   `json:"control_wards_per_game" db:"control_wards_per_game"`
	WardEfficiency       float64   `json:"ward_efficiency" db:"ward_efficiency"`
	
	// Performance by Phase
	EarlyGameVision      float64   `json:"early_game_vision" db:"early_game_vision"`
	MidGameVision        float64   `json:"mid_game_vision" db:"mid_game_vision"`
	LateGameVision       float64   `json:"late_game_vision" db:"late_game_vision"`
	
	// Impact Metrics
	VisionImpactScore    float64   `json:"vision_impact_score" db:"vision_impact_score"`
	MapControlScore      float64   `json:"map_control_score" db:"map_control_score"`
	ObjectiveVision      float64   `json:"objective_vision" db:"objective_vision"`
	
	// Comparative Analysis
	RolePercentile       float64   `json:"role_percentile" db:"role_percentile"`
	RankPercentile       float64   `json:"rank_percentile" db:"rank_percentile"`
	GlobalPercentile     float64   `json:"global_percentile" db:"global_percentile"`
	
	// Trend Analysis
	TrendDirection       string    `json:"trend_direction" db:"trend_direction"`
	TrendSlope           float64   `json:"trend_slope" db:"trend_slope"`
	TrendConfidence      float64   `json:"trend_confidence" db:"trend_confidence"`
	
	// Metadata
	LastCalculated       time.Time `json:"last_calculated" db:"last_calculated"`
	LastUpdated          time.Time `json:"last_updated" db:"last_updated"`
}

// VisionBenchmarkData represents benchmark data for vision comparison
type VisionBenchmarkData struct {
	ID                   string    `json:"id" db:"id"`
	
	// Benchmark Scope
	BenchmarkType        string    `json:"benchmark_type" db:"benchmark_type"` // "role", "rank", "global", "champion"
	FilterValue          string    `json:"filter_value" db:"filter_value"` // "ADC", "GOLD", "Jinx", etc.
	Region               string    `json:"region" db:"region"`
	
	// Sample Information
	TotalPlayers         int       `json:"total_players" db:"total_players"`
	TotalMatches         int       `json:"total_matches" db:"total_matches"`
	
	// Vision Score Benchmarks
	AverageVisionScore   float64   `json:"average_vision_score" db:"average_vision_score"`
	MedianVisionScore    float64   `json:"median_vision_score" db:"median_vision_score"`
	Top10PercentVision   float64   `json:"top_10_percent_vision" db:"top_10_percent_vision"`
	Top25PercentVision   float64   `json:"top_25_percent_vision" db:"top_25_percent_vision"`
	Bottom25PercentVision float64  `json:"bottom_25_percent_vision" db:"bottom_25_percent_vision"`
	
	// Ward Benchmarks
	AverageWardsPlaced   float64   `json:"average_wards_placed" db:"average_wards_placed"`
	AverageWardsKilled   float64   `json:"average_wards_killed" db:"average_wards_killed"`
	AverageControlWards  float64   `json:"average_control_wards" db:"average_control_wards"`
	AverageWardEfficiency float64  `json:"average_ward_efficiency" db:"average_ward_efficiency"`
	
	// Phase-specific Benchmarks
	EarlyGameVisionAvg   float64   `json:"early_game_vision_avg" db:"early_game_vision_avg"`
	MidGameVisionAvg     float64   `json:"mid_game_vision_avg" db:"mid_game_vision_avg"`
	LateGameVisionAvg    float64   `json:"late_game_vision_avg" db:"late_game_vision_avg"`
	
	// Win Rate Correlation
	VisionWinCorrelation float64   `json:"vision_win_correlation" db:"vision_win_correlation"`
	HighVisionWinRate    float64   `json:"high_vision_win_rate" db:"high_vision_win_rate"`
	LowVisionWinRate     float64   `json:"low_vision_win_rate" db:"low_vision_win_rate"`
	
	// Metadata
	LastCalculated       time.Time `json:"last_calculated" db:"last_calculated"`
	ValidUntil           time.Time `json:"valid_until" db:"valid_until"`
}

// MapZoneStats represents statistics for specific map zones
type MapZoneStats struct {
	ID                   string    `json:"id" db:"id"`
	PlayerID             string    `json:"player_id" db:"player_id"`
	ZoneName             string    `json:"zone_name" db:"zone_name"`
	TimeRange            string    `json:"time_range" db:"time_range"`
	
	// Zone Activity
	WardsPlaced          int       `json:"wards_placed" db:"wards_placed"`
	ControlWardsPlaced   int       `json:"control_wards_placed" db:"control_wards_placed"`
	WardsKilled          int       `json:"wards_killed" db:"wards_killed"`
	
	// Efficiency in Zone
	ZoneEfficiency       float64   `json:"zone_efficiency" db:"zone_efficiency"`
	StrategicValue       float64   `json:"strategic_value" db:"strategic_value"`
	
	// Win Rate Impact
	WinRateWithVision    float64   `json:"win_rate_with_vision" db:"win_rate_with_vision"`
	WinRateWithoutVision float64   `json:"win_rate_without_vision" db:"win_rate_without_vision"`
	VisionImpact         float64   `json:"vision_impact" db:"vision_impact"`
	
	// Timing Patterns
	EarlyGameActivity    int       `json:"early_game_activity" db:"early_game_activity"`
	MidGameActivity      int       `json:"mid_game_activity" db:"mid_game_activity"`
	LateGameActivity     int       `json:"late_game_activity" db:"late_game_activity"`
	
	LastUpdated          time.Time `json:"last_updated" db:"last_updated"`
}

// VisionRecommendationRule represents rules for generating vision recommendations
type VisionRecommendationRule struct {
	ID              string    `json:"id" db:"id"`
	
	// Rule Configuration
	RuleName        string    `json:"rule_name" db:"rule_name"`
	Category        string    `json:"category" db:"category"` // "warding", "dewarding", "positioning", "timing"
	Priority        string    `json:"priority" db:"priority"` // "high", "medium", "low"
	
	// Conditions
	MinMatches      int       `json:"min_matches" db:"min_matches"`
	VisionThreshold float64   `json:"vision_threshold" db:"vision_threshold"`
	PositionFilter  string    `json:"position_filter" db:"position_filter"`
	RankFilter      string    `json:"rank_filter" db:"rank_filter"`
	
	// Recommendation Content
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description" db:"description"`
	ImpactText      string    `json:"impact_text" db:"impact_text"`
	
	// Game Phase Applicability
	EarlyGame       bool      `json:"early_game" db:"early_game"`
	MidGame         bool      `json:"mid_game" db:"mid_game"`
	LateGame        bool      `json:"late_game" db:"late_game"`
	
	// Visual Aids
	VisualAidURL    string    `json:"visual_aid_url" db:"visual_aid_url"`
	MapOverlay      string    `json:"map_overlay" db:"map_overlay"`
	
	// Metadata
	Active          bool      `json:"active" db:"active"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// VisionInsight represents generated insights about vision performance
type VisionInsight struct {
	ID              string    `json:"id" db:"id"`
	PlayerID        string    `json:"player_id" db:"player_id"`
	
	// Insight Details
	InsightType     string    `json:"insight_type" db:"insight_type"` // "strength", "weakness", "trend", "comparison"
	Category        string    `json:"category" db:"category"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description" db:"description"`
	
	// Supporting Data
	MetricValue     float64   `json:"metric_value" db:"metric_value"`
	BenchmarkValue  float64   `json:"benchmark_value" db:"benchmark_value"`
	DifferencePercent float64 `json:"difference_percent" db:"difference_percent"`
	
	// Confidence and Impact
	Confidence      float64   `json:"confidence" db:"confidence"` // 0-1
	ImpactLevel     string    `json:"impact_level" db:"impact_level"` // "high", "medium", "low"
	
	// Time Relevance
	TimeRange       string    `json:"time_range" db:"time_range"`
	RecentMatches   int       `json:"recent_matches" db:"recent_matches"`
	
	// Actionability
	Actionable      bool      `json:"actionable" db:"actionable"`
	RelatedZones    string    `json:"related_zones" db:"related_zones"` // JSON array
	RelatedPhases   string    `json:"related_phases" db:"related_phases"` // JSON array
	
	// Metadata
	GeneratedAt     time.Time `json:"generated_at" db:"generated_at"`
	ExpiresAt       time.Time `json:"expires_at" db:"expires_at"`
	Viewed          bool      `json:"viewed" db:"viewed"`
}

// TableNames for GORM
func (WardPlacement) TableName() string {
	return "ward_placements"
}

func (WardKill) TableName() string {
	return "ward_kills"
}

func (VisionEvent) TableName() string {
	return "vision_events"
}

func (VisionHeatmapData) TableName() string {
	return "vision_heatmap_data"
}

func (PlayerVisionStats) TableName() string {
	return "player_vision_stats"
}

func (VisionBenchmarkData) TableName() string {
	return "vision_benchmark_data"
}

func (MapZoneStats) TableName() string {
	return "map_zone_stats"
}

func (VisionRecommendationRule) TableName() string {
	return "vision_recommendation_rules"
}

func (VisionInsight) TableName() string {
	return "vision_insights"
}