package models

import (
	"time"
)

// GoldAnalysis represents comprehensive gold efficiency analysis results
type GoldAnalysis struct {
	ID                       string    `json:"id" db:"id"`
	PlayerID                 string    `json:"player_id" db:"player_id"`
	Champion                 string    `json:"champion,omitempty" db:"champion"`
	Position                 string    `json:"position,omitempty" db:"position"`
	TimeRange                string    `json:"time_range" db:"time_range"`
	
	// Core Gold Metrics
	AverageGoldEarned        float64   `json:"average_gold_earned" db:"average_gold_earned"`
	AverageGoldPerMinute     float64   `json:"average_gold_per_minute" db:"average_gold_per_minute"`
	GoldEfficiencyScore      float64   `json:"gold_efficiency_score" db:"gold_efficiency_score"`
	EconomyRating            string    `json:"economy_rating" db:"economy_rating"`
	
	// Gold Sources Analysis (JSON stored)
	GoldSources              string    `json:"gold_sources" db:"gold_sources"`
	
	// Item Efficiency (JSON stored)
	ItemEfficiency           string    `json:"item_efficiency" db:"item_efficiency"`
	
	// Spending Patterns (JSON stored)
	SpendingPatterns         string    `json:"spending_patterns" db:"spending_patterns"`
	
	// Game Phase Analysis (JSON stored)
	EarlyGameGold            string    `json:"early_game_gold" db:"early_game_gold"`
	MidGameGold              string    `json:"mid_game_gold" db:"mid_game_gold"`
	LateGameGold             string    `json:"late_game_gold" db:"late_game_gold"`
	
	// Comparative Analysis (JSON stored)
	RoleBenchmark            string    `json:"role_benchmark" db:"role_benchmark"`
	RankBenchmark            string    `json:"rank_benchmark" db:"rank_benchmark"`
	GlobalBenchmark          string    `json:"global_benchmark" db:"global_benchmark"`
	
	// Performance Impact
	GoldAdvantageWinRate     float64   `json:"gold_advantage_win_rate" db:"gold_advantage_win_rate"`
	GoldDisadvantageWinRate  float64   `json:"gold_disadvantage_win_rate" db:"gold_disadvantage_win_rate"`
	GoldImpactScore          float64   `json:"gold_impact_score" db:"gold_impact_score"`
	
	// Trend Analysis
	TrendDirection           string    `json:"trend_direction" db:"trend_direction"`
	TrendSlope               float64   `json:"trend_slope" db:"trend_slope"`
	TrendConfidence          float64   `json:"trend_confidence" db:"trend_confidence"`
	TrendData                string    `json:"trend_data" db:"trend_data"` // JSON stored
	
	// Economy Optimization (JSON stored)
	IncomeOptimization       string    `json:"income_optimization" db:"income_optimization"`
	SpendingOptimization     string    `json:"spending_optimization" db:"spending_optimization"`
	
	// Insights and Recommendations (JSON stored)
	StrengthAreas            string    `json:"strength_areas" db:"strength_areas"`
	ImprovementAreas         string    `json:"improvement_areas" db:"improvement_areas"`
	Recommendations          string    `json:"recommendations" db:"recommendations"`
	
	// Match Performance (JSON stored)
	RecentMatches            string    `json:"recent_matches" db:"recent_matches"`
	
	// Metadata
	GeneratedAt              time.Time `json:"generated_at" db:"generated_at"`
	LastUpdated              time.Time `json:"last_updated" db:"last_updated"`
}

// GoldTransaction represents individual gold transactions in matches
type GoldTransaction struct {
	ID                       string    `json:"id" db:"id"`
	MatchID                  string    `json:"match_id" db:"match_id"`
	PlayerID                 string    `json:"player_id" db:"player_id"`
	
	// Transaction Details
	TransactionType          string    `json:"transaction_type" db:"transaction_type"` // "EARNED", "SPENT"
	Source                   string    `json:"source" db:"source"` // "MINION", "CHAMPION", "MONSTER", "ITEM", "PASSIVE"
	Amount                   int       `json:"amount" db:"amount"`
	
	// Context Information
	Timestamp                int       `json:"timestamp" db:"timestamp"` // Game time in seconds
	GamePhase                string    `json:"game_phase" db:"game_phase"` // "early", "mid", "late"
	X                        int       `json:"x" db:"x"`
	Y                        int       `json:"y" db:"y"`
	
	// Source Details
	SourceID                 string    `json:"source_id" db:"source_id"` // Champion ID, Monster ID, etc.
	SourceName               string    `json:"source_name" db:"source_name"`
	
	// Item Purchase Details (if applicable)
	ItemID                   int       `json:"item_id" db:"item_id"`
	ItemName                 string    `json:"item_name" db:"item_name"`
	ItemCost                 int       `json:"item_cost" db:"item_cost"`
	
	// Economic Context
	TotalGoldBefore          int       `json:"total_gold_before" db:"total_gold_before"`
	TotalGoldAfter           int       `json:"total_gold_after" db:"total_gold_after"`
	TeamGold                 int       `json:"team_gold" db:"team_gold"`
	EnemyTeamGold            int       `json:"enemy_team_gold" db:"enemy_team_gold"`
	
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
}

// ItemPurchase represents item purchase events with efficiency analysis
type ItemPurchase struct {
	ID                       string    `json:"id" db:"id"`
	MatchID                  string    `json:"match_id" db:"match_id"`
	PlayerID                 string    `json:"player_id" db:"player_id"`
	
	// Item Details
	ItemID                   int       `json:"item_id" db:"item_id"`
	ItemName                 string    `json:"item_name" db:"item_name"`
	ItemCost                 int       `json:"item_cost" db:"item_cost"`
	ItemCategory             string    `json:"item_category" db:"item_category"` // "damage", "defensive", "utility"
	
	// Purchase Context
	Timestamp                int       `json:"timestamp" db:"timestamp"`
	GamePhase                string    `json:"game_phase" db:"game_phase"`
	PurchaseOrder            int       `json:"purchase_order" db:"purchase_order"` // 1st, 2nd, 3rd item, etc.
	
	// Economic Context
	GoldSpent                int       `json:"gold_spent" db:"gold_spent"`
	GoldRemaining            int       `json:"gold_remaining" db:"gold_remaining"`
	BackTiming               bool      `json:"back_timing" db:"back_timing"` // Was this a good time to back?
	
	// Efficiency Metrics
	ValueEfficiency          float64   `json:"value_efficiency" db:"value_efficiency"` // 0-100 score
	TimingEfficiency         float64   `json:"timing_efficiency" db:"timing_efficiency"` // 0-100 score
	SituationalEfficiency    float64   `json:"situational_efficiency" db:"situational_efficiency"` // 0-100 score
	
	// Performance Impact
	DamageGained             int       `json:"damage_gained" db:"damage_gained"` // Estimated damage increase
	DefenseGained            int       `json:"defense_gained" db:"defense_gained"` // Estimated defense increase
	UtilityGained            int       `json:"utility_gained" db:"utility_gained"` // Estimated utility increase
	
	// Strategic Context
	CounterBuild             bool      `json:"counter_build" db:"counter_build"` // Was this a counter to enemy team?
	PowerSpike               bool      `json:"power_spike" db:"power_spike"` // Did this complete a power spike?
	OptimalChoice            bool      `json:"optimal_choice" db:"optimal_choice"` // Was this the optimal item choice?
	
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
}

// GoldBenchmark represents benchmark data for gold efficiency comparison
type GoldBenchmark struct {
	ID                       string    `json:"id" db:"id"`
	
	// Benchmark Scope
	BenchmarkType            string    `json:"benchmark_type" db:"benchmark_type"` // "role", "rank", "global", "champion"
	FilterValue              string    `json:"filter_value" db:"filter_value"` // "ADC", "GOLD", "Jinx", etc.
	Region                   string    `json:"region" db:"region"`
	
	// Sample Information
	TotalPlayers             int       `json:"total_players" db:"total_players"`
	TotalMatches             int       `json:"total_matches" db:"total_matches"`
	
	// Gold Generation Benchmarks
	AverageGoldPerMinute     float64   `json:"average_gold_per_minute" db:"average_gold_per_minute"`
	MedianGoldPerMinute      float64   `json:"median_gold_per_minute" db:"median_gold_per_minute"`
	Top10PercentGPM          float64   `json:"top_10_percent_gpm" db:"top_10_percent_gpm"`
	Top25PercentGPM          float64   `json:"top_25_percent_gpm" db:"top_25_percent_gpm"`
	Bottom25PercentGPM       float64   `json:"bottom_25_percent_gpm" db:"bottom_25_percent_gpm"`
	
	// Efficiency Benchmarks
	AverageEfficiencyScore   float64   `json:"average_efficiency_score" db:"average_efficiency_score"`
	MedianEfficiencyScore    float64   `json:"median_efficiency_score" db:"median_efficiency_score"`
	Top10PercentEfficiency   float64   `json:"top_10_percent_efficiency" db:"top_10_percent_efficiency"`
	
	// Gold Sources Benchmarks
	AverageFarmingPercent    float64   `json:"average_farming_percent" db:"average_farming_percent"`
	AverageKillsPercent      float64   `json:"average_kills_percent" db:"average_kills_percent"`
	AverageObjectivePercent  float64   `json:"average_objective_percent" db:"average_objective_percent"`
	
	// Item Efficiency Benchmarks
	AverageItemsCompleted    float64   `json:"average_items_completed" db:"average_items_completed"`
	AverageFirstItemTiming   float64   `json:"average_first_item_timing" db:"average_first_item_timing"`
	AverageCoreItemsTiming   float64   `json:"average_core_items_timing" db:"average_core_items_timing"`
	
	// Phase-specific Benchmarks
	EarlyGameGPMAvg          float64   `json:"early_game_gpm_avg" db:"early_game_gpm_avg"`
	MidGameGPMAvg            float64   `json:"mid_game_gpm_avg" db:"mid_game_gpm_avg"`
	LateGameGPMAvg           float64   `json:"late_game_gpm_avg" db:"late_game_gpm_avg"`
	
	// Win Rate Correlation
	GoldWinCorrelation       float64   `json:"gold_win_correlation" db:"gold_win_correlation"`
	HighGoldWinRate          float64   `json:"high_gold_win_rate" db:"high_gold_win_rate"`
	LowGoldWinRate           float64   `json:"low_gold_win_rate" db:"low_gold_win_rate"`
	
	// Metadata
	LastCalculated           time.Time `json:"last_calculated" db:"last_calculated"`
	ValidUntil               time.Time `json:"valid_until" db:"valid_until"`
}

// PlayerGoldStats represents aggregated gold statistics for a player
type PlayerGoldStats struct {
	ID                       string    `json:"id" db:"id"`
	PlayerID                 string    `json:"player_id" db:"player_id"`
	TimeRange                string    `json:"time_range" db:"time_range"`
	Position                 string    `json:"position" db:"position"`
	Champion                 string    `json:"champion" db:"champion"`
	
	// Basic Statistics
	TotalMatches             int       `json:"total_matches" db:"total_matches"`
	AverageGoldEarned        float64   `json:"average_gold_earned" db:"average_gold_earned"`
	MedianGoldEarned         float64   `json:"median_gold_earned" db:"median_gold_earned"`
	BestGoldEarned           int64     `json:"best_gold_earned" db:"best_gold_earned"`
	WorstGoldEarned          int64     `json:"worst_gold_earned" db:"worst_gold_earned"`
	GoldStdDev               float64   `json:"gold_std_dev" db:"gold_std_dev"`
	
	// Gold Per Minute Stats
	AverageGPM               float64   `json:"average_gpm" db:"average_gpm"`
	BestGPM                  float64   `json:"best_gpm" db:"best_gpm"`
	GPMStdDev                float64   `json:"gpm_std_dev" db:"gpm_std_dev"`
	
	// Gold Sources Distribution
	FarmingGoldPercent       float64   `json:"farming_gold_percent" db:"farming_gold_percent"`
	KillsGoldPercent         float64   `json:"kills_gold_percent" db:"kills_gold_percent"`
	ObjectiveGoldPercent     float64   `json:"objective_gold_percent" db:"objective_gold_percent"`
	PassiveGoldPercent       float64   `json:"passive_gold_percent" db:"passive_gold_percent"`
	
	// Performance by Phase
	EarlyGameGPM             float64   `json:"early_game_gpm" db:"early_game_gpm"`
	MidGameGPM               float64   `json:"mid_game_gpm" db:"mid_game_gpm"`
	LateGameGPM              float64   `json:"late_game_gpm" db:"late_game_gpm"`
	
	// Efficiency Metrics
	GoldEfficiencyScore      float64   `json:"gold_efficiency_score" db:"gold_efficiency_score"`
	ItemEfficiencyScore      float64   `json:"item_efficiency_score" db:"item_efficiency_score"`
	SpendingEfficiencyScore  float64   `json:"spending_efficiency_score" db:"spending_efficiency_score"`
	
	// Impact Metrics
	GoldImpactScore          float64   `json:"gold_impact_score" db:"gold_impact_score"`
	EconomyConsistencyScore  float64   `json:"economy_consistency_score" db:"economy_consistency_score"`
	
	// Comparative Analysis
	RolePercentile           float64   `json:"role_percentile" db:"role_percentile"`
	RankPercentile           float64   `json:"rank_percentile" db:"rank_percentile"`
	GlobalPercentile         float64   `json:"global_percentile" db:"global_percentile"`
	
	// Win Rate Analysis
	HighGoldWinRate          float64   `json:"high_gold_win_rate" db:"high_gold_win_rate"`
	LowGoldWinRate           float64   `json:"low_gold_win_rate" db:"low_gold_win_rate"`
	
	// Trend Analysis
	TrendDirection           string    `json:"trend_direction" db:"trend_direction"`
	TrendSlope               float64   `json:"trend_slope" db:"trend_slope"`
	TrendConfidence          float64   `json:"trend_confidence" db:"trend_confidence"`
	
	// Metadata
	LastCalculated           time.Time `json:"last_calculated" db:"last_calculated"`
	LastUpdated              time.Time `json:"last_updated" db:"last_updated"`
}

// GoldInsight represents AI-generated insights about gold performance
type GoldInsight struct {
	ID                       string    `json:"id" db:"id"`
	PlayerID                 string    `json:"player_id" db:"player_id"`
	
	// Insight Details
	InsightType              string    `json:"insight_type" db:"insight_type"` // "strength", "weakness", "trend", "comparison"
	Category                 string    `json:"category" db:"category"` // "gold_generation", "spending_efficiency", "farming"
	Title                    string    `json:"title" db:"title"`
	Description              string    `json:"description" db:"description"`
	
	// Supporting Data
	MetricValue              float64   `json:"metric_value" db:"metric_value"`
	BenchmarkValue           float64   `json:"benchmark_value" db:"benchmark_value"`
	DifferencePercent        float64   `json:"difference_percent" db:"difference_percent"`
	
	// Confidence and Impact
	Confidence               float64   `json:"confidence" db:"confidence"` // 0-1
	ImpactLevel              string    `json:"impact_level" db:"impact_level"` // "high", "medium", "low"
	
	// Time Relevance
	TimeRange                string    `json:"time_range" db:"time_range"`
	RecentMatches            int       `json:"recent_matches" db:"recent_matches"`
	
	// Actionability
	Actionable               bool      `json:"actionable" db:"actionable"`
	RelatedPhases            string    `json:"related_phases" db:"related_phases"` // JSON array
	RelatedChampions         string    `json:"related_champions" db:"related_champions"` // JSON array
	
	// Recommendations
	RecommendedActions       string    `json:"recommended_actions" db:"recommended_actions"` // JSON array
	ExpectedImprovement      float64   `json:"expected_improvement" db:"expected_improvement"`
	
	// Metadata
	GeneratedAt              time.Time `json:"generated_at" db:"generated_at"`
	ExpiresAt                time.Time `json:"expires_at" db:"expires_at"`
	Viewed                   bool      `json:"viewed" db:"viewed"`
}

// GoldOptimization represents optimization suggestions for gold efficiency
type GoldOptimization struct {
	ID                       string    `json:"id" db:"id"`
	PlayerID                 string    `json:"player_id" db:"player_id"`
	
	// Optimization Target
	OptimizationType         string    `json:"optimization_type" db:"optimization_type"` // "income", "spending", "efficiency"
	Champion                 string    `json:"champion" db:"champion"`
	Position                 string    `json:"position" db:"position"`
	GamePhase                string    `json:"game_phase" db:"game_phase"`
	
	// Current Performance
	CurrentGPM               float64   `json:"current_gpm" db:"current_gpm"`
	CurrentEfficiencyScore   float64   `json:"current_efficiency_score" db:"current_efficiency_score"`
	CurrentEconomyRating     string    `json:"current_economy_rating" db:"current_economy_rating"`
	
	// Optimization Goals
	TargetGPM                float64   `json:"target_gpm" db:"target_gpm"`
	TargetEfficiencyScore    float64   `json:"target_efficiency_score" db:"target_efficiency_score"`
	TargetEconomyRating      string    `json:"target_economy_rating" db:"target_economy_rating"`
	
	// Improvement Strategies (JSON stored)
	IncomeStrategies         string    `json:"income_strategies" db:"income_strategies"`
	SpendingStrategies       string    `json:"spending_strategies" db:"spending_strategies"`
	EfficiencyStrategies     string    `json:"efficiency_strategies" db:"efficiency_strategies"`
	
	// Expected Impact
	ExpectedGPMIncrease      float64   `json:"expected_gpm_increase" db:"expected_gpm_increase"`
	ExpectedWinRateGain      float64   `json:"expected_win_rate_gain" db:"expected_win_rate_gain"`
	ImplementationDifficulty string    `json:"implementation_difficulty" db:"implementation_difficulty"` // "easy", "medium", "hard"
	
	// Timeline
	ExpectedTimeToImprovement int      `json:"expected_time_to_improvement" db:"expected_time_to_improvement"` // Days
	MilestoneTargets         string    `json:"milestone_targets" db:"milestone_targets"` // JSON array
	
	// Tracking
	ProgressTracking         string    `json:"progress_tracking" db:"progress_tracking"` // JSON object
	
	// Metadata
	GeneratedAt              time.Time `json:"generated_at" db:"generated_at"`
	LastUpdated              time.Time `json:"last_updated" db:"last_updated"`
	Priority                 string    `json:"priority" db:"priority"` // "high", "medium", "low"
	Active                   bool      `json:"active" db:"active"`
}

// BackEvent represents recall/back events with timing analysis
type BackEvent struct {
	ID                       string    `json:"id" db:"id"`
	MatchID                  string    `json:"match_id" db:"match_id"`
	PlayerID                 string    `json:"player_id" db:"player_id"`
	
	// Back Details
	Timestamp                int       `json:"timestamp" db:"timestamp"` // Game time when back started
	Duration                 int       `json:"duration" db:"duration"` // Time spent shopping
	GoldSpent                int       `json:"gold_spent" db:"gold_spent"`
	GoldRemaining            int       `json:"gold_remaining" db:"gold_remaining"`
	
	// Timing Analysis
	BackReason               string    `json:"back_reason" db:"back_reason"` // "low_hp", "low_mana", "items", "forced"
	TimingQuality            string    `json:"timing_quality" db:"timing_quality"` // "optimal", "good", "suboptimal", "forced"
	OpportunityCost          int       `json:"opportunity_cost" db:"opportunity_cost"` // Estimated gold/XP lost
	
	// Shopping Efficiency
	ItemsPurchased           string    `json:"items_purchased" db:"items_purchased"` // JSON array of item IDs
	ShoppingTime             int       `json:"shopping_time" db:"shopping_time"` // Seconds in shop
	OptimalPurchases         bool      `json:"optimal_purchases" db:"optimal_purchases"`
	
	// Context
	GamePhase                string    `json:"game_phase" db:"game_phase"`
	WaveState                string    `json:"wave_state" db:"wave_state"` // "pushing", "frozen", "neutral"
	ObjectiveState           string    `json:"objective_state" db:"objective_state"` // Upcoming objectives
	
	// Impact
	PowerSpikeAchieved       bool      `json:"power_spike_achieved" db:"power_spike_achieved"`
	GoldEfficiencyGain       float64   `json:"gold_efficiency_gain" db:"gold_efficiency_gain"`
	
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
}

// TableNames for GORM
func (GoldAnalysis) TableName() string {
	return "gold_analysis"
}

func (GoldTransaction) TableName() string {
	return "gold_transactions"
}

func (ItemPurchase) TableName() string {
	return "item_purchases"
}

func (GoldBenchmark) TableName() string {
	return "gold_benchmarks"
}

func (PlayerGoldStats) TableName() string {
	return "player_gold_stats"
}

func (GoldInsight) TableName() string {
	return "gold_insights"
}

func (GoldOptimization) TableName() string {
	return "gold_optimizations"
}

func (BackEvent) TableName() string {
	return "back_events"
}