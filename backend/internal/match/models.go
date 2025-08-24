package match

import (
	"time"

	"github.com/herald-lol/herald/backend/internal/riot"
)

// Herald.lol Gaming Analytics - Match Analysis Models
// Data structures for match analysis requests and results

// MatchAnalysisRequest contains parameters for match analysis
type MatchAnalysisRequest struct {
	Match                   *riot.Match `json:"match"`
	PlayerPUUID             string      `json:"player_puuid"`
	AnalysisDepth           string      `json:"analysis_depth"` // "basic", "standard", "detailed", "professional"
	IncludePhaseAnalysis    bool        `json:"include_phase_analysis"`
	IncludeKeyMoments       bool        `json:"include_key_moments"`
	IncludeTeamAnalysis     bool        `json:"include_team_analysis"`
	IncludeOpponentAnalysis bool        `json:"include_opponent_analysis"`
	CompareWithAverage      bool        `json:"compare_with_average"`
	FocusAreas              []string    `json:"focus_areas"` // "farming", "fighting", "vision", "objectives"
}

// MatchAnalysisResult contains comprehensive match analysis
type MatchAnalysisResult struct {
	MatchID               string                 `json:"match_id"`
	PlayerPUUID           string                 `json:"player_puuid"`
	MatchInfo             *MatchInfo             `json:"match_info"`
	Performance           *PerformanceAnalysis   `json:"performance"`
	PhaseAnalysis         *GamePhaseAnalysis     `json:"phase_analysis,omitempty"`
	KeyMoments            []*KeyMoment           `json:"key_moments,omitempty"`
	TeamAnalysis          *TeamAnalysis          `json:"team_analysis,omitempty"`
	OpponentAnalysis      *OpponentAnalysis      `json:"opponent_analysis,omitempty"`
	Insights              *MatchInsights         `json:"insights"`
	LearningOpportunities []*LearningOpportunity `json:"learning_opportunities"`
	OverallRating         float64                `json:"overall_rating"`    // 0-100
	PerformanceGrade      string                 `json:"performance_grade"` // "S+", "S", "A+", "A", "B+", "B", "C+", "C", "D"
	AnalyzedAt            time.Time              `json:"analyzed_at"`
}

// MatchInfo contains basic match information
type MatchInfo struct {
	GameMode     string    `json:"game_mode"`
	QueueType    string    `json:"queue_type"`
	GameDuration int       `json:"game_duration"`
	GameVersion  string    `json:"game_version"`
	Champion     string    `json:"champion"`
	Role         string    `json:"role"`
	Result       string    `json:"result"` // "Victory", "Defeat"
	KDA          float64   `json:"kda"`
	Score        int       `json:"score"` // Composite score
	PlayedAt     time.Time `json:"played_at"`
}

// PerformanceAnalysis contains detailed performance metrics
type PerformanceAnalysis struct {
	// Combat Stats
	Kills             int     `json:"kills"`
	Deaths            int     `json:"deaths"`
	Assists           int     `json:"assists"`
	KDA               float64 `json:"kda"`
	KillParticipation float64 `json:"kill_participation"`

	// Farming Stats
	CreepScore  int     `json:"creep_score"`
	CSPerMinute float64 `json:"cs_per_minute"`
	CSRating    string  `json:"cs_rating"`    // "Excellent", "Good", "Average", "Poor"
	CSAdvantage int     `json:"cs_advantage"` // vs lane opponent

	// Economic Stats
	Gold           int     `json:"gold"`
	GoldPerMinute  float64 `json:"gold_per_minute"`
	GoldEfficiency float64 `json:"gold_efficiency"`
	GoldAdvantage  int     `json:"gold_advantage"`

	// Combat Stats
	Damage           int     `json:"damage"`
	DamagePerMinute  float64 `json:"damage_per_minute"`
	DamageShare      float64 `json:"damage_share"`      // % of team damage
	DamageEfficiency float64 `json:"damage_efficiency"` // Damage per gold

	// Vision Stats
	VisionScore     int     `json:"vision_score"`
	VisionPerMinute float64 `json:"vision_per_minute"`
	VisionRating    string  `json:"vision_rating"`
	WardsPlaced     int     `json:"wards_placed"`
	WardsKilled     int     `json:"wards_killed"`
	ControlWards    int     `json:"control_wars_bought"`

	// Objective Stats
	ObjectiveStats *ObjectiveStats `json:"objective_stats"`

	// Special Achievements
	MultiKills *MultiKillStats `json:"multi_kills"`
	FirstBlood bool            `json:"first_blood"`
	SoloKills  int             `json:"solo_kills"`

	// Survival Stats
	LongestLiving int     `json:"longest_living"` // Seconds
	TimeDead      int     `json:"time_dead"`      // Seconds
	SurvivalRate  float64 `json:"survival_rate"`  // % of game alive

	// Overall Performance
	OverallRating     float64 `json:"overall_rating"`     // 0-100
	PerformanceLevel  string  `json:"performance_level"`  // "Outstanding", "Great", "Good", "Average", "Poor"
	RoleEffectiveness float64 `json:"role_effectiveness"` // How well they played their role
}

// ObjectiveStats contains objective participation statistics
type ObjectiveStats struct {
	DragonKills       int     `json:"dragon_kills"`
	BaronKills        int     `json:"baron_kills"`
	TurretKills       int     `json:"turret_kills"`
	TurretDamage      int     `json:"turret_damage"`
	InhibitorKills    int     `json:"inhibitor_kills"`
	ObjectiveParticip float64 `json:"objective_participation"` // % of team objectives
	ObjectiveControl  float64 `json:"objective_control"`       // Impact on objectives
}

// MultiKillStats contains multi-kill statistics
type MultiKillStats struct {
	DoubleKills  int `json:"double_kills"`
	TripleKills  int `json:"triple_kills"`
	QuadraKills  int `json:"quadra_kills"`
	PentaKills   int `json:"penta_kills"`
	LargestSpree int `json:"largest_spree"`
}

// GamePhaseAnalysis breaks down performance by game phases
type GamePhaseAnalysis struct {
	LanePhase        *PhasePerformance `json:"lane_phase"`
	MidGame          *PhasePerformance `json:"mid_game"`
	LateGame         *PhasePerformance `json:"late_game"`
	StrongestPhase   string            `json:"strongest_phase"`
	WeakestPhase     string            `json:"weakest_phase"`
	PhaseConsistency float64           `json:"phase_consistency"` // 0-100
}

// PhasePerformance contains performance data for a game phase
type PhasePerformance struct {
	Phase     string `json:"phase"`
	StartTime int    `json:"start_time"`
	EndTime   int    `json:"end_time"`
	Duration  int    `json:"duration"`

	// Performance in this phase
	Kills   int     `json:"kills"`
	Deaths  int     `json:"deaths"`
	Assists int     `json:"assists"`
	KDA     float64 `json:"kda"`

	// Economic performance
	GoldEarned int `json:"gold_earned"`
	CSGained   int `json:"cs_gained"`

	// Combat performance
	DamageDealt int `json:"damage_dealt"`

	// Key events in this phase
	KeyEvents []string `json:"key_events"`

	// Phase rating
	PhaseRating float64 `json:"phase_rating"` // 0-100
	PhaseGrade  string  `json:"phase_grade"`
	Impact      string  `json:"impact"` // "High", "Medium", "Low"

	// Improvement suggestions
	Improvements []string `json:"improvements"`
}

// KeyMoment represents a significant moment in the match
type KeyMoment struct {
	Type           string `json:"type"`       // "First Blood", "Multi Kill", "Objective", etc.
	Timestamp      int    `json:"timestamp"`  // Game time in seconds
	GamePhase      string `json:"game_phase"` // "Lane", "Mid", "Late"
	Impact         string `json:"impact"`     // "Very Positive", "Positive", "Neutral", "Negative", "Very Negative"
	Importance     int    `json:"importance"` // 1-10 scale
	Description    string `json:"description"`
	LearningPoint  string `json:"learning_point"`
	Context        string `json:"context"`
	Recommendation string `json:"recommendation"`
}

// TeamAnalysis contains team performance analysis
type TeamAnalysis struct {
	TeamID            int     `json:"team_id"`
	TeamSize          int     `json:"team_size"`
	TeamKDA           float64 `json:"team_kda"`
	TotalDamage       int     `json:"total_damage"`
	TotalGold         int     `json:"total_gold"`
	AverageVision     float64 `json:"average_vision"`
	ObjectivesSecured int     `json:"objectives_secured"`

	// Player's role in team
	PlayerContribution *PlayerTeamContribution `json:"player_contribution"`

	// Team synergy
	SynergyRating      float64 `json:"synergy_rating"`  // 0-100
	TeamplayRating     float64 `json:"teamplay_rating"` // 0-100
	CommunicationScore float64 `json:"communication_score"`

	// Team composition analysis
	Composition *TeamComposition `json:"composition"`

	// Win conditions
	WinConditions   []string `json:"win_conditions"`
	ExecutionRating float64  `json:"execution_rating"` // How well team executed
}

// PlayerTeamContribution shows player's contribution to team
type PlayerTeamContribution struct {
	KillParticipation  float64 `json:"kill_participation"`  // % of team kills/assists
	DamageShare        float64 `json:"damage_share"`        // % of team damage
	GoldShare          float64 `json:"gold_share"`          // % of team gold
	VisionContribution float64 `json:"vision_contribution"` // % of team vision
	ObjectiveShare     float64 `json:"objective_share"`     // % of objective participation
	Leadership         float64 `json:"leadership"`          // Leadership rating
	Supportiveness     float64 `json:"supportiveness"`      // How much they helped team
}

// TeamComposition analyzes team comp strengths/weaknesses
type TeamComposition struct {
	Type                string   `json:"type"`            // "Teamfight", "Poke", "Split Push", "Pick"
	StrengthPhases      []string `json:"strength_phases"` // When the comp is strong
	WeakPhases          []string `json:"weak_phases"`     // When the comp is weak
	WinConditions       []string `json:"win_conditions"`
	CounterStrategy     []string `json:"counter_strategy"`
	FlexibilityRating   float64  `json:"flexibility_rating"`   // How adaptable the comp is
	ExecutionDifficulty float64  `json:"execution_difficulty"` // How hard to execute
	SynergyRating       float64  `json:"synergy_rating"`
}

// OpponentAnalysis contains analysis of lane opponent (if available)
type OpponentAnalysis struct {
	OpponentChampion string `json:"opponent_champion"`
	OpponentRole     string `json:"opponent_role"`

	// Head-to-head stats
	KillsAgainst int `json:"kills_against"`
	DeathsTo     int `json:"deaths_to"`

	// Lane comparison
	CSAdvantage     int `json:"cs_advantage"`     // Player's CS - Opponent's CS
	GoldAdvantage   int `json:"gold_advantage"`   // Player's Gold - Opponent's Gold
	DamageAdvantage int `json:"damage_advantage"` // Player's Damage - Opponent's Damage

	// Laning outcome
	LaneResult        string  `json:"lane_result"`         // "Won", "Lost", "Even"
	LaningPhaseRating float64 `json:"laning_phase_rating"` // 0-100

	// Matchup analysis
	MatchupDifficulty     string  `json:"matchup_difficulty"`      // "Easy", "Medium", "Hard", "Extreme"
	ExpectedOutcome       string  `json:"expected_outcome"`        // What should have happened
	ActualOutcome         string  `json:"actual_outcome"`          // What actually happened
	PerformanceVsExpected float64 `json:"performance_vs_expected"` // How well vs expectation
}

// MatchInsights contains AI-generated insights about the match
type MatchInsights struct {
	Strengths         []string `json:"strengths"`
	Weaknesses        []string `json:"weaknesses"`
	KeyTakeaways      []string `json:"key_takeaways"`
	OverallAssessment string   `json:"overall_assessment"`

	// Specific insights
	LanePhaseFeedback string `json:"lane_phase_feedback"`
	TeamfightFeedback string `json:"teamfight_feedback"`
	ObjectiveFeedback string `json:"objective_feedback"`
	VisionFeedback    string `json:"vision_feedback"`

	// Improvement focus
	PrimaryFocus   string `json:"primary_focus"` // Most important thing to work on
	SecondaryFocus string `json:"secondary_focus"`

	// Performance context
	DifficultyContext string `json:"difficulty_context"` // How hard was this match
	MetaContext       string `json:"meta_context"`       // How this fits current meta
}

// LearningOpportunity represents specific learning opportunities
type LearningOpportunity struct {
	Category            string   `json:"category"` // "Positioning", "Farming", "Fighting", etc.
	Description         string   `json:"description"`
	Importance          string   `json:"importance"` // "High", "Medium", "Low"
	Difficulty          string   `json:"difficulty"` // "Easy", "Medium", "Hard"
	ActionSteps         []string `json:"action_steps"`
	PracticeMethod      string   `json:"practice_method"` // How to practice this
	ExpectedImprovement string   `json:"expected_improvement"`
	TimeToImprove       string   `json:"time_to_improve"` // "Days", "Weeks", "Months"
}

// Match Series Analysis Models

// MatchSeriesRequest contains parameters for analyzing multiple matches
type MatchSeriesRequest struct {
	Matches            []*riot.Match `json:"matches"`
	PlayerPUUID        string        `json:"player_puuid"`
	AnalysisType       string        `json:"analysis_type"` // "trend", "pattern", "consistency"
	TimeFrame          string        `json:"time_frame"`    // "recent", "week", "month"
	FocusMetrics       []string      `json:"focus_metrics"` // Which metrics to focus on
	IncludeComparisons bool          `json:"include_comparisons"`
}

// MatchSeriesAnalysis contains analysis of multiple matches
type MatchSeriesAnalysis struct {
	PlayerPUUID        string                 `json:"player_puuid"`
	TotalMatches       int                    `json:"total_matches"`
	SuccessfulAnalyses int                    `json:"successful_analyses"`
	AnalysisType       string                 `json:"analysis_type"`
	TimeFrame          string                 `json:"time_frame"`
	MatchAnalyses      []*MatchAnalysisResult `json:"match_analyses"`

	// Series-level insights
	SeriesInsights      *SeriesInsights      `json:"series_insights"`
	PerformancePatterns *PerformancePatterns `json:"performance_patterns"`
	ImprovementAreas    []SeriesImprovement  `json:"improvement_areas"`
	ConsistencyMetrics  *ConsistencyMetrics  `json:"consistency_metrics"`

	AnalyzedAt time.Time `json:"analyzed_at"`
}

// SeriesInsights contains insights across multiple matches
type SeriesInsights struct {
	OverallTrend          string   `json:"overall_trend"` // "Improving", "Declining", "Stable"
	BestMatch             string   `json:"best_match"`    // Match ID of best performance
	WorstMatch            string   `json:"worst_match"`   // Match ID of worst performance
	ConsistentStrengths   []string `json:"consistent_strengths"`
	ConsistentWeaknesses  []string `json:"consistent_weaknesses"`
	PerformanceVolatility float64  `json:"performance_volatility"` // How much performance varies

	// Win/Loss patterns
	WinStreak      int      `json:"win_streak"`
	LossStreak     int      `json:"loss_streak"`
	WinConditions  []string `json:"win_conditions"`  // What leads to wins
	LossConditions []string `json:"loss_conditions"` // What leads to losses
}

// PerformancePatterns identifies patterns in performance
type PerformancePatterns struct {
	ChampionPatterns   map[string]*ChampionPattern `json:"champion_patterns"`
	RolePatterns       map[string]*RolePattern     `json:"role_patterns"`
	TimePatterns       *TimePattern                `json:"time_patterns"`
	MatchupPatterns    *MatchupPattern             `json:"matchup_patterns"`
	GameLengthPatterns *GameLengthPattern          `json:"game_length_patterns"`
}

// ChampionPattern shows performance patterns on specific champions
type ChampionPattern struct {
	Champion           string   `json:"champion"`
	GamesPlayed        int      `json:"games_played"`
	WinRate            float64  `json:"win_rate"`
	AveragePerformance float64  `json:"average_performance"`
	Consistency        float64  `json:"consistency"`
	TrendDirection     string   `json:"trend_direction"`
	KeyStrengths       []string `json:"key_strengths"`
	ImprovementAreas   []string `json:"improvement_areas"`
}

// RolePattern shows performance patterns by role
type RolePattern struct {
	Role               string  `json:"role"`
	GamesPlayed        int     `json:"games_played"`
	AveragePerformance float64 `json:"average_performance"`
	RelativeStrength   string  `json:"relative_strength"` // vs other roles
	Specialization     float64 `json:"specialization"`    // How specialized in this role
}

// TimePattern analyzes performance by time of day/day of week
type TimePattern struct {
	BestTimeOfDay      string  `json:"best_time_of_day"`
	WorstTimeOfDay     string  `json:"worst_time_of_day"`
	TimeConsistency    float64 `json:"time_consistency"`
	WeekdayPerformance float64 `json:"weekday_performance"`
	WeekendPerformance float64 `json:"weekend_performance"`
}

// MatchupPattern analyzes performance against different matchups
type MatchupPattern struct {
	EasiestMatchups     []string `json:"easiest_matchups"`
	HardestMatchups     []string `json:"hardest_matchups"`
	MatchupAdaptability float64  `json:"matchup_adaptability"`
	CounterplaySkill    float64  `json:"counterplay_skill"`
}

// GameLengthPattern analyzes performance by game length
type GameLengthPattern struct {
	ShortGamePerformance  float64 `json:"short_game_performance"`  // <20 min
	MediumGamePerformance float64 `json:"medium_game_performance"` // 20-35 min
	LongGamePerformance   float64 `json:"long_game_performance"`   // >35 min
	OptimalGameLength     int     `json:"optimal_game_length"`     // Best performance length
	Stamina               float64 `json:"stamina"`                 // Performance in long games
}

// SeriesImprovement represents improvement opportunities from series
type SeriesImprovement struct {
	Area            string   `json:"area"`
	Priority        string   `json:"priority"`         // "High", "Medium", "Low"
	Frequency       int      `json:"frequency"`        // How often this appears
	ImpactPotential float64  `json:"impact_potential"` // Potential rating improvement
	Difficulty      string   `json:"difficulty"`
	SpecificActions []string `json:"specific_actions"`
	MeasurableGoals []string `json:"measurable_goals"`
}

// ConsistencyMetrics measures performance consistency
type ConsistencyMetrics struct {
	OverallConsistency float64 `json:"overall_consistency"` // 0-100
	KDAConsistency     float64 `json:"kda_consistency"`
	CSConsistency      float64 `json:"cs_consistency"`
	VisionConsistency  float64 `json:"vision_consistency"`
	DamageConsistency  float64 `json:"damage_consistency"`

	PerformanceRange       float64 `json:"performance_range"` // Max - Min performance
	StandardDeviation      float64 `json:"standard_deviation"`
	CoefficientOfVariation float64 `json:"coefficient_of_variation"`

	// Clutch performance
	ClutchRating     float64 `json:"clutch_rating"`     // Performance in close games
	PressureHandling float64 `json:"pressure_handling"` // Performance under pressure
	Tilt             float64 `json:"tilt_resistance"`   // Resistance to tilting
}

// Match Comparison Models

// MatchComparisonRequest contains parameters for comparing two matches
type MatchComparisonRequest struct {
	Match1      *riot.Match `json:"match1"`
	Match2      *riot.Match `json:"match2"`
	PlayerPUUID string      `json:"player_puuid"`
}

// MatchComparisonResult contains comparison between two matches
type MatchComparisonResult struct {
	Match1Analysis        *MatchAnalysisResult    `json:"match1_analysis"`
	Match2Analysis        *MatchAnalysisResult    `json:"match2_analysis"`
	PerformanceComparison *PerformanceComparison  `json:"performance_comparison"`
	ImprovementAreas      []ComparisonImprovement `json:"improvement_areas"`
	Summary               string                  `json:"summary"`
	ComparedAt            time.Time               `json:"compared_at"`
}

// PerformanceComparison compares performance metrics between matches
type PerformanceComparison struct {
	KDAComparison    *MetricComparison `json:"kda_comparison"`
	CSComparison     *MetricComparison `json:"cs_comparison"`
	DamageComparison *MetricComparison `json:"damage_comparison"`
	VisionComparison *MetricComparison `json:"vision_comparison"`
	GoldComparison   *MetricComparison `json:"gold_comparison"`

	OverallImprovement string  `json:"overall_improvement"` // "Better", "Worse", "Similar"
	ImprovementScore   float64 `json:"improvement_score"`   // -100 to +100

	BetterAreas  []string `json:"better_areas"`
	WorseAreas   []string `json:"worse_areas"`
	SimilarAreas []string `json:"similar_areas"`
}

// MetricComparison compares a specific metric between matches
type MetricComparison struct {
	Metric        string  `json:"metric"`
	Match1Value   float64 `json:"match1_value"`
	Match2Value   float64 `json:"match2_value"`
	Change        float64 `json:"change"`
	PercentChange float64 `json:"percent_change"`
	Direction     string  `json:"direction"`    // "Improved", "Declined", "Same"
	Significance  string  `json:"significance"` // "Major", "Moderate", "Minor"
}

// ComparisonImprovement represents improvements identified through comparison
type ComparisonImprovement struct {
	Area           string `json:"area"`
	Improvement    string `json:"improvement"`
	Evidence       string `json:"evidence"`
	Recommendation string `json:"recommendation"`
	Priority       string `json:"priority"`
}
