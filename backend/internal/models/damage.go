package models

import (
	"time"
)

// DamageAnalysis represents comprehensive damage analysis results
type DamageAnalysis struct {
	ID                   string                 `json:"id" db:"id"`
	PlayerID             string                 `json:"player_id" db:"player_id"`
	Champion             string                 `json:"champion,omitempty" db:"champion"`
	Position             string                 `json:"position,omitempty" db:"position"`
	TimeRange            string                 `json:"time_range" db:"time_range"`
	
	// Core Damage Metrics
	DamageShare          float64                `json:"damage_share" db:"damage_share"`
	DamagePerMinute      float64                `json:"damage_per_minute" db:"damage_per_minute"`
	TotalDamage          int64                  `json:"total_damage" db:"total_damage"`
	PhysicalDamagePercent float64               `json:"physical_damage_percent" db:"physical_damage_percent"`
	MagicDamagePercent   float64                `json:"magic_damage_percent" db:"magic_damage_percent"`
	TrueDamagePercent    float64                `json:"true_damage_percent" db:"true_damage_percent"`
	
	// Performance Metrics
	CarryPotential       float64                `json:"carry_potential" db:"carry_potential"`
	EfficiencyRating     string                 `json:"efficiency_rating" db:"efficiency_rating"`
	DamageConsistency    float64                `json:"damage_consistency" db:"damage_consistency"`
	
	// Team Contribution (JSON stored)
	TeamContribution     string                 `json:"team_contribution" db:"team_contribution"`
	
	// Distribution Analysis (JSON stored)
	DamageDistribution   string                 `json:"damage_distribution" db:"damage_distribution"`
	
	// Game Phase Analysis (JSON stored)
	GamePhaseAnalysis    string                 `json:"game_phase_analysis" db:"game_phase_analysis"`
	
	// Win Rate Analysis
	HighDamageWinRate    float64                `json:"high_damage_win_rate" db:"high_damage_win_rate"`
	LowDamageWinRate     float64                `json:"low_damage_win_rate" db:"low_damage_win_rate"`
	
	// Benchmark Comparisons (JSON stored)
	RoleBenchmark        string                 `json:"role_benchmark" db:"role_benchmark"`
	RankBenchmark        string                 `json:"rank_benchmark" db:"rank_benchmark"`
	GlobalBenchmark      string                 `json:"global_benchmark" db:"global_benchmark"`
	
	// Trend Analysis (JSON stored)
	TrendData            string                 `json:"trend_data" db:"trend_data"`
	
	// Recommendations (JSON stored)
	Recommendations      string                 `json:"recommendations" db:"recommendations"`
	
	// Metadata
	GeneratedAt          time.Time              `json:"generated_at" db:"generated_at"`
	LastUpdated          time.Time              `json:"last_updated" db:"last_updated"`
}

// DamageEvent represents individual damage events in matches
type DamageEvent struct {
	ID                   string    `json:"id" db:"id"`
	MatchID              string    `json:"match_id" db:"match_id"`
	PlayerID             string    `json:"player_id" db:"player_id"`
	
	// Event Details
	EventType            string    `json:"event_type" db:"event_type"` // "CHAMPION_DAMAGE", "STRUCTURE_DAMAGE", "MONSTER_DAMAGE"
	DamageAmount         int       `json:"damage_amount" db:"damage_amount"`
	DamageType           string    `json:"damage_type" db:"damage_type"` // "PHYSICAL", "MAGICAL", "TRUE"
	
	// Target Information
	TargetType           string    `json:"target_type" db:"target_type"` // "CHAMPION", "MINION", "MONSTER", "STRUCTURE"
	TargetID             string    `json:"target_id" db:"target_id"`
	TargetName           string    `json:"target_name" db:"target_name"`
	
	// Timing and Position
	Timestamp            int       `json:"timestamp" db:"timestamp"` // Game time in seconds
	GamePhase            string    `json:"game_phase" db:"game_phase"` // "early", "mid", "late"
	X                    int       `json:"x" db:"x"`
	Y                    int       `json:"y" db:"y"`
	
	// Context Information
	ItemsUsed            string    `json:"items_used" db:"items_used"` // JSON array of item IDs
	AbilityUsed          string    `json:"ability_used" db:"ability_used"` // Ability that caused damage
	IsCritical           bool      `json:"is_critical" db:"is_critical"`
	
	// Team Fight Context
	IsTeamFight          bool      `json:"is_team_fight" db:"is_team_fight"`
	NearbyAllies         int       `json:"nearby_allies" db:"nearby_allies"`
	NearbyEnemies        int       `json:"nearby_enemies" db:"nearby_enemies"`
	
	// Outcome Context
	KilledTarget         bool      `json:"killed_target" db:"killed_target"`
	AssistOnKill         bool      `json:"assist_on_kill" db:"assist_on_kill"`
	
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
}

// TeamFightDamage represents damage dealt during team fights
type TeamFightDamage struct {
	ID                   string    `json:"id" db:"id"`
	MatchID              string    `json:"match_id" db:"match_id"`
	PlayerID             string    `json:"player_id" db:"player_id"`
	TeamFightID          string    `json:"team_fight_id" db:"team_fight_id"`
	
	// Team Fight Info
	StartTime            int       `json:"start_time" db:"start_time"`
	EndTime              int       `json:"end_time" db:"end_time"`
	Duration             int       `json:"duration" db:"duration"`
	
	// Damage Metrics
	TotalDamage          int       `json:"total_damage" db:"total_damage"`
	DamageToChampions    int       `json:"damage_to_champions" db:"damage_to_champions"`
	DamageShare          float64   `json:"damage_share" db:"damage_share"`
	DPS                  float64   `json:"dps" db:"dps"` // Damage per second
	
	// Team Fight Performance
	KillsInFight         int       `json:"kills_in_fight" db:"kills_in_fight"`
	AssistsInFight       int       `json:"assists_in_fight" db:"assists_in_fight"`
	DeathsInFight        int       `json:"deaths_in_fight" db:"deaths_in_fight"`
	
	// Positioning Metrics
	AverageDistanceToEnemies float64 `json:"avg_distance_to_enemies" db:"avg_distance_to_enemies"`
	FrontLineTime        float64   `json:"front_line_time" db:"front_line_time"` // Time spent in front line
	BackLineTime         float64   `json:"back_line_time" db:"back_line_time"`   // Time spent in back line
	
	// Impact Metrics
	CarryPerformance     float64   `json:"carry_performance" db:"carry_performance"` // 0-100 score
	ClutchFactor         float64   `json:"clutch_factor" db:"clutch_factor"`         // Damage when low HP
	
	// Team Fight Outcome
	TeamWon              bool      `json:"team_won" db:"team_won"`
	ObjectiveSecured     string    `json:"objective_secured" db:"objective_secured"` // "DRAGON", "BARON", etc.
	
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
}

// DamageBenchmarkData represents benchmark data for damage comparison
type DamageBenchmarkData struct {
	ID                   string    `json:"id" db:"id"`
	
	// Benchmark Scope
	BenchmarkType        string    `json:"benchmark_type" db:"benchmark_type"` // "role", "rank", "global", "champion"
	FilterValue          string    `json:"filter_value" db:"filter_value"`     // "ADC", "GOLD", "Jinx", etc.
	Region               string    `json:"region" db:"region"`
	
	// Sample Information
	TotalPlayers         int       `json:"total_players" db:"total_players"`
	TotalMatches         int       `json:"total_matches" db:"total_matches"`
	
	// Damage Benchmarks
	AverageDamageShare   float64   `json:"average_damage_share" db:"average_damage_share"`
	MedianDamageShare    float64   `json:"median_damage_share" db:"median_damage_share"`
	Top10PercentDamage   float64   `json:"top_10_percent_damage" db:"top_10_percent_damage"`
	Top25PercentDamage   float64   `json:"top_25_percent_damage" db:"top_25_percent_damage"`
	Bottom25PercentDamage float64  `json:"bottom_25_percent_damage" db:"bottom_25_percent_damage"`
	
	// DPM Benchmarks
	AverageDPM           float64   `json:"average_dpm" db:"average_dpm"`
	MedianDPM            float64   `json:"median_dpm" db:"median_dpm"`
	Top10PercentDPM      float64   `json:"top_10_percent_dpm" db:"top_10_percent_dpm"`
	
	// Carry Potential Benchmarks
	AverageCarryPotential float64  `json:"average_carry_potential" db:"average_carry_potential"`
	HighCarryWinRate     float64   `json:"high_carry_win_rate" db:"high_carry_win_rate"`
	LowCarryWinRate      float64   `json:"low_carry_win_rate" db:"low_carry_win_rate"`
	
	// Phase-specific Benchmarks
	EarlyGameDamageAvg   float64   `json:"early_game_damage_avg" db:"early_game_damage_avg"`
	MidGameDamageAvg     float64   `json:"mid_game_damage_avg" db:"mid_game_damage_avg"`
	LateGameDamageAvg    float64   `json:"late_game_damage_avg" db:"late_game_damage_avg"`
	
	// Win Rate Correlation
	DamageWinCorrelation float64   `json:"damage_win_correlation" db:"damage_win_correlation"`
	HighDamageWinRate    float64   `json:"high_damage_win_rate" db:"high_damage_win_rate"`
	LowDamageWinRate     float64   `json:"low_damage_win_rate" db:"low_damage_win_rate"`
	
	// Metadata
	LastCalculated       time.Time `json:"last_calculated" db:"last_calculated"`
	ValidUntil           time.Time `json:"valid_until" db:"valid_until"`
}

// PlayerDamageStats represents aggregated damage statistics for a player
type PlayerDamageStats struct {
	ID                   string    `json:"id" db:"id"`
	PlayerID             string    `json:"player_id" db:"player_id"`
	TimeRange            string    `json:"time_range" db:"time_range"`
	Position             string    `json:"position" db:"position"`
	Champion             string    `json:"champion" db:"champion"`
	
	// Basic Statistics
	TotalMatches         int       `json:"total_matches" db:"total_matches"`
	AverageDamage        float64   `json:"average_damage" db:"average_damage"`
	MedianDamage         float64   `json:"median_damage" db:"median_damage"`
	BestDamage           int64     `json:"best_damage" db:"best_damage"`
	WorstDamage          int64     `json:"worst_damage" db:"worst_damage"`
	DamageStdDev         float64   `json:"damage_std_dev" db:"damage_std_dev"`
	
	// Damage Per Minute Stats
	AverageDPM           float64   `json:"average_dpm" db:"average_dpm"`
	BestDPM              float64   `json:"best_dpm" db:"best_dpm"`
	DPMStdDev            float64   `json:"dpm_std_dev" db:"dpm_std_dev"`
	
	// Damage Share Stats
	AverageDamageShare   float64   `json:"average_damage_share" db:"average_damage_share"`
	BestDamageShare      float64   `json:"best_damage_share" db:"best_damage_share"`
	DamageShareStdDev    float64   `json:"damage_share_std_dev" db:"damage_share_std_dev"`
	
	// Damage Type Distribution
	PhysicalDamagePercent float64  `json:"physical_damage_percent" db:"physical_damage_percent"`
	MagicDamagePercent   float64   `json:"magic_damage_percent" db:"magic_damage_percent"`
	TrueDamagePercent    float64   `json:"true_damage_percent" db:"true_damage_percent"`
	
	// Performance by Phase
	EarlyGameDamage      float64   `json:"early_game_damage" db:"early_game_damage"`
	MidGameDamage        float64   `json:"mid_game_damage" db:"mid_game_damage"`
	LateGameDamage       float64   `json:"late_game_damage" db:"late_game_damage"`
	
	// Team Fight Performance
	TeamFightDamageShare float64   `json:"team_fight_damage_share" db:"team_fight_damage_share"`
	TeamFightKP          float64   `json:"team_fight_kp" db:"team_fight_kp"`
	
	// Impact Metrics
	CarryPotential       float64   `json:"carry_potential" db:"carry_potential"`
	DamageImpactScore    float64   `json:"damage_impact_score" db:"damage_impact_score"`
	ConsistencyScore     float64   `json:"consistency_score" db:"consistency_score"`
	
	// Comparative Analysis
	RolePercentile       float64   `json:"role_percentile" db:"role_percentile"`
	RankPercentile       float64   `json:"rank_percentile" db:"rank_percentile"`
	GlobalPercentile     float64   `json:"global_percentile" db:"global_percentile"`
	
	// Win Rate Analysis
	HighDamageWinRate    float64   `json:"high_damage_win_rate" db:"high_damage_win_rate"`
	LowDamageWinRate     float64   `json:"low_damage_win_rate" db:"low_damage_win_rate"`
	
	// Trend Analysis
	TrendDirection       string    `json:"trend_direction" db:"trend_direction"`
	TrendSlope           float64   `json:"trend_slope" db:"trend_slope"`
	TrendConfidence      float64   `json:"trend_confidence" db:"trend_confidence"`
	
	// Metadata
	LastCalculated       time.Time `json:"last_calculated" db:"last_calculated"`
	LastUpdated          time.Time `json:"last_updated" db:"last_updated"`
}

// DamageInsight represents AI-generated insights about damage performance
type DamageInsight struct {
	ID                   string    `json:"id" db:"id"`
	PlayerID             string    `json:"player_id" db:"player_id"`
	
	// Insight Details
	InsightType          string    `json:"insight_type" db:"insight_type"` // "strength", "weakness", "trend", "comparison"
	Category             string    `json:"category" db:"category"`         // "damage_output", "consistency", "positioning"
	Title                string    `json:"title" db:"title"`
	Description          string    `json:"description" db:"description"`
	
	// Supporting Data
	MetricValue          float64   `json:"metric_value" db:"metric_value"`
	BenchmarkValue       float64   `json:"benchmark_value" db:"benchmark_value"`
	DifferencePercent    float64   `json:"difference_percent" db:"difference_percent"`
	
	// Confidence and Impact
	Confidence           float64   `json:"confidence" db:"confidence"`     // 0-1
	ImpactLevel          string    `json:"impact_level" db:"impact_level"` // "high", "medium", "low"
	
	// Time Relevance
	TimeRange            string    `json:"time_range" db:"time_range"`
	RecentMatches        int       `json:"recent_matches" db:"recent_matches"`
	
	// Actionability
	Actionable           bool      `json:"actionable" db:"actionable"`
	RelatedChampions     string    `json:"related_champions" db:"related_champions"` // JSON array
	RelatedPhases        string    `json:"related_phases" db:"related_phases"`       // JSON array
	
	// Recommendations
	RecommendedActions   string    `json:"recommended_actions" db:"recommended_actions"` // JSON array
	
	// Metadata
	GeneratedAt          time.Time `json:"generated_at" db:"generated_at"`
	ExpiresAt            time.Time `json:"expires_at" db:"expires_at"`
	Viewed               bool      `json:"viewed" db:"viewed"`
}

// DamageOptimization represents optimization suggestions for damage performance
type DamageOptimization struct {
	ID                   string    `json:"id" db:"id"`
	PlayerID             string    `json:"player_id" db:"player_id"`
	
	// Optimization Target
	Champion             string    `json:"champion" db:"champion"`
	Position             string    `json:"position" db:"position"`
	GamePhase            string    `json:"game_phase" db:"game_phase"`
	
	// Current Performance
	CurrentDamageShare   float64   `json:"current_damage_share" db:"current_damage_share"`
	CurrentDPM           float64   `json:"current_dpm" db:"current_dpm"`
	CurrentCarryPotential float64  `json:"current_carry_potential" db:"current_carry_potential"`
	
	// Optimization Goals
	TargetDamageShare    float64   `json:"target_damage_share" db:"target_damage_share"`
	TargetDPM            float64   `json:"target_dpm" db:"target_dpm"`
	TargetCarryPotential float64   `json:"target_carry_potential" db:"target_carry_potential"`
	
	// Improvement Areas (JSON stored)
	ItemBuildSuggestions string    `json:"item_build_suggestions" db:"item_build_suggestions"`
	AbilityOrderSuggestions string `json:"ability_order_suggestions" db:"ability_order_suggestions"`
	PositioningSuggestions string  `json:"positioning_suggestions" db:"positioning_suggestions"`
	TimingSuggestions    string    `json:"timing_suggestions" db:"timing_suggestions"`
	
	// Expected Impact
	ExpectedWinRateGain  float64   `json:"expected_win_rate_gain" db:"expected_win_rate_gain"`
	ImplementationDifficulty string `json:"implementation_difficulty" db:"implementation_difficulty"` // "easy", "medium", "hard"
	
	// Metadata
	GeneratedAt          time.Time `json:"generated_at" db:"generated_at"`
	LastUpdated          time.Time `json:"last_updated" db:"last_updated"`
	Priority             string    `json:"priority" db:"priority"` // "high", "medium", "low"
}

// TableNames for GORM
func (DamageAnalysis) TableName() string {
	return "damage_analysis"
}

func (DamageEvent) TableName() string {
	return "damage_events"
}

func (TeamFightDamage) TableName() string {
	return "team_fight_damage"
}

func (DamageBenchmarkData) TableName() string {
	return "damage_benchmark_data"
}

func (PlayerDamageStats) TableName() string {
	return "player_damage_stats"
}

func (DamageInsight) TableName() string {
	return "damage_insights"
}

func (DamageOptimization) TableName() string {
	return "damage_optimizations"
}