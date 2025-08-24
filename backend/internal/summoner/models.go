package summoner

import (
	"time"

	"github.com/herald-lol/herald/backend/internal/analytics"
	// "github.com/herald-lol/herald/backend/internal/riot"
)

// Herald.lol Gaming Analytics - Summoner Service Models
// Data structures for summoner analytics requests and responses

// SummonerAnalysisRequest contains parameters for summoner analysis
type SummonerAnalysisRequest struct {
	RequestID        string `json:"request_id"`
	Region           string `json:"region"`
	SummonerName     string `json:"summoner_name"`
	SubscriptionTier string `json:"subscription_tier"` // "free", "premium", "pro", "enterprise"
	AnalysisType     string `json:"analysis_type"`     // "basic", "standard", "detailed", "professional"
	TimeFrame        string `json:"time_frame"`        // "recent", "season", "all"
	MatchCount       int    `json:"match_count"`       // Override default match count
	UseCache         bool   `json:"use_cache"`         // Use cached results if available
	IncludeLiveGame  bool   `json:"include_live_game"` // Include current live game analysis
	IncludeInsights  bool   `json:"include_insights"`  // Include AI-generated insights
}

// SummonerAnalysisResponse contains complete summoner analysis
type SummonerAnalysisResponse struct {
	RequestID        string                    `json:"request_id"`
	Region           string                    `json:"region"`
	Summoner         *SummonerInfo             `json:"summoner"`
	RankedInfo       []*RankedInfo             `json:"ranked_info"`
	Analytics        *analytics.PlayerAnalysis `json:"analytics"`
	ChampionMastery  []*ChampionMasteryInfo    `json:"champion_mastery"`
	LiveGame         *LiveGameInfo             `json:"live_game,omitempty"`
	Recommendations  *RecommendationSummary    `json:"recommendations,omitempty"`
	ProcessingTimeMs int                       `json:"processing_time_ms"`
	StartedAt        time.Time                 `json:"started_at"`
	CompletedAt      time.Time                 `json:"completed_at"`
	CacheHit         bool                      `json:"cache_hit"`
	DataFreshness    string                    `json:"data_freshness"` // "live", "cached", "mixed"
}

// SummonerInfo contains basic summoner information
type SummonerInfo struct {
	ID            string    `json:"id"`
	PUUID         string    `json:"puuid"`
	Name          string    `json:"name"`
	Level         int       `json:"level"`
	ProfileIconID int       `json:"profile_icon_id"`
	AccountID     string    `json:"account_id,omitempty"`
	LastUpdated   time.Time `json:"last_updated"`
}

// RankedInfo contains ranked queue information
type RankedInfo struct {
	QueueType      string          `json:"queue_type"`
	Tier           string          `json:"tier"`
	Rank           string          `json:"rank"`
	LeaguePoints   int             `json:"league_points"`
	Wins           int             `json:"wins"`
	Losses         int             `json:"losses"`
	WinRate        float64         `json:"win_rate"`
	HotStreak      bool            `json:"hot_streak"`
	Veteran        bool            `json:"veteran"`
	FreshBlood     bool            `json:"fresh_blood"`
	Inactive       bool            `json:"inactive"`
	MiniSeries     *MiniSeriesInfo `json:"mini_series,omitempty"`
	PeakThisSeason *RankSnapshot   `json:"peak_this_season,omitempty"`
	Progression    []RankSnapshot  `json:"progression,omitempty"`
}

// MiniSeriesInfo contains promotional series information
type MiniSeriesInfo struct {
	Target   int    `json:"target"`
	Wins     int    `json:"wins"`
	Losses   int    `json:"losses"`
	Progress string `json:"progress"`
}

// RankSnapshot represents rank at a specific time
type RankSnapshot struct {
	Timestamp    time.Time `json:"timestamp"`
	Tier         string    `json:"tier"`
	Rank         string    `json:"rank"`
	LeaguePoints int       `json:"league_points"`
	Games        int       `json:"games"`
	WinRate      float64   `json:"win_rate"`
}

// ChampionMasteryInfo contains champion mastery information
type ChampionMasteryInfo struct {
	ChampionID     int       `json:"champion_id"`
	ChampionName   string    `json:"champion_name"`
	ChampionLevel  int       `json:"champion_level"`
	ChampionPoints int       `json:"champion_points"`
	LastPlayTime   time.Time `json:"last_play_time"`
	ChestGranted   bool      `json:"chest_granted"`
	TokensEarned   int       `json:"tokens_earned"`
	PointsToNext   int64     `json:"points_to_next"`
	ProgressToNext float64   `json:"progress_to_next"`
}

// LiveGameInfo contains current live game information
type LiveGameInfo struct {
	GameID       int64                `json:"game_id"`
	GameType     string               `json:"game_type"`
	GameMode     string               `json:"game_mode"`
	GameLength   int64                `json:"game_length"`
	MapID        int                  `json:"map_id"`
	Participants []*LiveParticipant   `json:"participants"`
	PlayerTeam   *LiveTeamInfo        `json:"player_team"`
	EnemyTeam    *LiveTeamInfo        `json:"enemy_team"`
	Predictions  *LiveGamePredictions `json:"predictions,omitempty"`
	Analysis     *LiveGameAnalysis    `json:"analysis,omitempty"`
}

// LiveParticipant contains live game participant info
type LiveParticipant struct {
	SummonerName  string  `json:"summoner_name"`
	ChampionName  string  `json:"champion_name"`
	ChampionID    int64   `json:"champion_id"`
	TeamID        int     `json:"team_id"`
	Spell1        int64   `json:"spell1"`
	Spell2        int64   `json:"spell2"`
	ProfileIconID int64   `json:"profile_icon_id"`
	Rank          string  `json:"rank,omitempty"`
	WinRate       float64 `json:"win_rate,omitempty"`
	KDA           float64 `json:"kda,omitempty"`
	IsBot         bool    `json:"is_bot"`
}

// LiveTeamInfo contains live game team information
type LiveTeamInfo struct {
	TeamID         int                `json:"team_id"`
	Participants   []*LiveParticipant `json:"participants"`
	Bans           []ChampionBan      `json:"bans"`
	TeamStrength   float64            `json:"team_strength"`
	WinProbability float64            `json:"win_probability"`
	Composition    *TeamComposition   `json:"composition"`
}

// ChampionBan contains champion ban information
type ChampionBan struct {
	ChampionID int `json:"champion_id"`
	PickTurn   int `json:"pick_turn"`
}

// TeamComposition analyzes team composition
type TeamComposition struct {
	CompositionType   string   `json:"composition_type"` // "Poke", "Teamfight", "Split Push", etc.
	StrengthPhases    []string `json:"strength_phases"`  // "Early", "Mid", "Late"
	WeakPhases        []string `json:"weak_phases"`
	WinConditions     []string `json:"win_conditions"`
	ThreatsToWatch    []string `json:"threats_to_watch"`
	StrengthRating    float64  `json:"strength_rating"`    // 0-100
	SynergyRating     float64  `json:"synergy_rating"`     // 0-100
	FlexibilityRating float64  `json:"flexibility_rating"` // 0-100
}

// LiveGamePredictions contains match outcome predictions
type LiveGamePredictions struct {
	WinProbability      float64            `json:"win_probability"`
	ConfidenceLevel     float64            `json:"confidence_level"`
	KeyFactors          []PredictionFactor `json:"key_factors"`
	GameDurationPredict string             `json:"game_duration_predict"`
	MVPPrediction       string             `json:"mvp_prediction"`
	FirstObjectives     map[string]float64 `json:"first_objectives"` // Probability of getting first tower, dragon, etc.
}

// PredictionFactor represents a factor in match prediction
type PredictionFactor struct {
	Factor      string  `json:"factor"`
	Impact      float64 `json:"impact"`     // -1 to 1
	Confidence  float64 `json:"confidence"` // 0 to 1
	Description string  `json:"description"`
}

// LiveGameAnalysis contains real-time game analysis
type LiveGameAnalysis struct {
	CurrentPhase      string              `json:"current_phase"` // "Pick/Ban", "Early", "Mid", "Late"
	PlayerPerformance *LivePlayerStats    `json:"player_performance"`
	TeamComparison    *LiveTeamComparison `json:"team_comparison"`
	KeyMoments        []LiveKeyMoment     `json:"key_moments"`
	Recommendations   []string            `json:"recommendations"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

// LivePlayerStats contains live player performance stats
type LivePlayerStats struct {
	EstimatedKDA      float64 `json:"estimated_kda"`
	EstimatedCS       int     `json:"estimated_cs"`
	EstimatedGold     int     `json:"estimated_gold"`
	LaneState         string  `json:"lane_state"` // "winning", "losing", "even"
	ObjectivePresence float64 `json:"objective_presence"`
	WardingEfficiency float64 `json:"warding_efficiency"`
	PerformanceRating float64 `json:"performance_rating"` // 0-100
}

// LiveTeamComparison compares both teams in live game
type LiveTeamComparison struct {
	GoldAdvantage       int     `json:"gold_advantage"`
	ExperienceAdvantage int     `json:"experience_advantage"`
	DragonAdvantage     int     `json:"dragon_advantage"`
	TowerAdvantage      int     `json:"tower_advantage"`
	KillAdvantage       int     `json:"kill_advantage"`
	OverallAdvantage    string  `json:"overall_advantage"`  // "blue", "red", "even"
	AdvantageStrength   float64 `json:"advantage_strength"` // 0-1
}

// LiveKeyMoment represents significant moments in live game
type LiveKeyMoment struct {
	Timestamp   int    `json:"timestamp"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Impact      string `json:"impact"`   // "positive", "negative", "neutral"
	Severity    int    `json:"severity"` // 1-5
}

// RecommendationSummary contains coaching recommendations
type RecommendationSummary struct {
	ImmediateFocus      []string             `json:"immediate_focus"`
	ChampionPool        []string             `json:"champion_pool"`
	RoleFocus           string               `json:"role_focus"`
	SkillPriorities     []SkillPriority      `json:"skill_priorities"`
	TrainingPlan        *TrainingPlanSummary `json:"training_plan"`
	NextRankTarget      string               `json:"next_rank_target"`
	EstimatedTimeToRank string               `json:"estimated_time_to_rank"`
	ConfidenceLevel     float64              `json:"confidence_level"`
}

// SkillPriority represents a skill area to focus on
type SkillPriority struct {
	Skill        string  `json:"skill"`
	Priority     int     `json:"priority"`      // 1-5
	CurrentLevel float64 `json:"current_level"` // 0-100
	TargetLevel  float64 `json:"target_level"`  // 0-100
	Improvement  float64 `json:"improvement"`   // Points needed
}

// TrainingPlanSummary contains summary of training recommendations
type TrainingPlanSummary struct {
	Duration       string   `json:"duration"`
	DailyTime      string   `json:"daily_time"`
	KeyExercises   []string `json:"key_exercises"`
	Milestones     []string `json:"milestones"`
	SuccessMetrics []string `json:"success_metrics"`
}

// Comparison models

// SummonerComparisonRequest contains comparison request parameters
type SummonerComparisonRequest struct {
	Region           string `json:"region"`
	Summoner1Name    string `json:"summoner1_name"`
	Summoner2Name    string `json:"summoner2_name"`
	SubscriptionTier string `json:"subscription_tier"`
	ComparisonType   string `json:"comparison_type"` // "basic", "detailed", "head_to_head"
	TimeFrame        string `json:"time_frame"`
}

// SummonerComparisonResponse contains comparison results
type SummonerComparisonResponse struct {
	Summoner1   *SummonerInfo     `json:"summoner1"`
	Summoner2   *SummonerInfo     `json:"summoner2"`
	Comparison  *ComparisonResult `json:"comparison"`
	GeneratedAt time.Time         `json:"generated_at"`
}

// ComparisonResult contains detailed comparison between summoners
type ComparisonResult struct {
	OverallWinner     string                       `json:"overall_winner"` // "summoner1", "summoner2", "tie"
	WinnerMargin      float64                      `json:"winner_margin"`  // 0-100
	MetricComparisons map[string]*MetricComparison `json:"metric_comparisons"`
	StrengthAreas     *ComparisonStrengths         `json:"strength_areas"`
	RankComparison    *RankComparison              `json:"rank_comparison"`
	ChampionOverlap   *ChampionOverlapAnalysis     `json:"champion_overlap"`
	HeadToHeadStats   *HeadToHeadStats             `json:"head_to_head_stats,omitempty"`
	Summary           string                       `json:"summary"`
}

// MetricComparison compares a specific metric between summoners
type MetricComparison struct {
	Metric       string  `json:"metric"`
	Summoner1    float64 `json:"summoner1"`
	Summoner2    float64 `json:"summoner2"`
	Winner       string  `json:"winner"`
	Difference   float64 `json:"difference"`
	Significance string  `json:"significance"` // "minor", "moderate", "major", "extreme"
}

// ComparisonStrengths shows what each summoner is better at
type ComparisonStrengths struct {
	Summoner1Strengths []string `json:"summoner1_strengths"`
	Summoner2Strengths []string `json:"summoner2_strengths"`
	SharedStrengths    []string `json:"shared_strengths"`
}

// RankComparison compares ranked performance
type RankComparison struct {
	CurrentRankWinner string  `json:"current_rank_winner"`
	Summoner1Rank     string  `json:"summoner1_rank"`
	Summoner2Rank     string  `json:"summoner2_rank"`
	RankDifference    int     `json:"rank_difference"` // Difference in rank tiers
	LPDifference      int     `json:"lp_difference"`
	WinRateComparison float64 `json:"win_rate_comparison"`
	ClimbingRate      string  `json:"climbing_rate"` // Who's climbing faster
}

// ChampionOverlapAnalysis analyzes champion pool overlap
type ChampionOverlapAnalysis struct {
	SharedChampions    []string `json:"shared_champions"`
	UniqueToSummoner1  []string `json:"unique_to_summoner1"`
	UniqueToSummoner2  []string `json:"unique_to_summoner2"`
	OverlapPercentage  float64  `json:"overlap_percentage"`
	PoolSimilarity     string   `json:"pool_similarity"` // "very_similar", "somewhat_similar", "different"
	BestSharedChampion string   `json:"best_shared_champion"`
}

// HeadToHeadStats contains direct competition statistics
type HeadToHeadStats struct {
	GamesPlayed       int       `json:"games_played"`
	Summoner1Wins     int       `json:"summoner1_wins"`
	Summoner2Wins     int       `json:"summoner2_wins"`
	HeadToHeadWinRate float64   `json:"head_to_head_win_rate"`
	LastMatchDate     time.Time `json:"last_match_date"`
	RecentForm        string    `json:"recent_form"` // "W-L-W-L-W" format
}

// Trends models

// SummonerTrendsRequest contains trend analysis parameters
type SummonerTrendsRequest struct {
	Region       string   `json:"region"`
	SummonerName string   `json:"summoner_name"`
	TimeWindows  []string `json:"time_windows"` // ["7d", "30d", "90d", "season"]
}

// SummonerTrendsResponse contains trend analysis results
type SummonerTrendsResponse struct {
	SummonerName string                  `json:"summoner_name"`
	Region       string                  `json:"region"`
	Trends       map[string]*TrendPeriod `json:"trends"`
	AnalyzedAt   time.Time               `json:"analyzed_at"`
}

// TrendPeriod contains trend data for a specific time period
type TrendPeriod struct {
	Period        string             `json:"period"`
	StartDate     time.Time          `json:"start_date"`
	EndDate       time.Time          `json:"end_date"`
	GamesPlayed   int                `json:"games_played"`
	Metrics       *TrendMetrics      `json:"metrics"`
	Improvements  []TrendImprovement `json:"improvements"`
	Declines      []TrendDecline     `json:"declines"`
	Achievements  []TrendAchievement `json:"achievements"`
	TrendStrength string             `json:"trend_strength"` // "strong_positive", "positive", "stable", "negative", "strong_negative"
}

// TrendMetrics contains trending metrics
type TrendMetrics struct {
	WinRate     *MetricTrend `json:"win_rate"`
	KDA         *MetricTrend `json:"kda"`
	CSPerMinute *MetricTrend `json:"cs_per_minute"`
	VisionScore *MetricTrend `json:"vision_score"`
	DamageShare *MetricTrend `json:"damage_share"`
}

// MetricTrend represents a trending metric
type MetricTrend struct {
	StartValue    float64 `json:"start_value"`
	EndValue      float64 `json:"end_value"`
	Change        float64 `json:"change"`
	PercentChange float64 `json:"percent_change"`
	Direction     string  `json:"direction"`    // "up", "down", "stable"
	Significance  string  `json:"significance"` // "minor", "moderate", "major"
}

// TrendImprovement represents an area of improvement
type TrendImprovement struct {
	Area        string  `json:"area"`
	Improvement float64 `json:"improvement"`
	Description string  `json:"description"`
}

// TrendDecline represents an area of decline
type TrendDecline struct {
	Area        string  `json:"area"`
	Decline     float64 `json:"decline"`
	Description string  `json:"description"`
	Concern     string  `json:"concern"` // "low", "medium", "high"
}

// TrendAchievement represents achievements during the period
type TrendAchievement struct {
	Achievement  string    `json:"achievement"`
	Date         time.Time `json:"date"`
	Significance string    `json:"significance"`
}

// Insights models

// SummonerInsightsRequest contains insights request parameters
type SummonerInsightsRequest struct {
	Region           string   `json:"region"`
	SummonerName     string   `json:"summoner_name"`
	SubscriptionTier string   `json:"subscription_tier"`
	InsightTypes     []string `json:"insight_types"` // ["performance", "coaching", "predictions", "comparisons"]
}

// SummonerInsightsResponse contains AI-generated insights
type SummonerInsightsResponse struct {
	SummonerName string            `json:"summoner_name"`
	Region       string            `json:"region"`
	Insights     *EnhancedInsights `json:"insights"`
	GeneratedAt  time.Time         `json:"generated_at"`
}

// EnhancedInsights contains comprehensive AI insights
type EnhancedInsights struct {
	PerformanceInsights *PerformanceInsights `json:"performance_insights"`
	CoachingInsights    *CoachingInsights    `json:"coaching_insights"`
	PredictiveInsights  *PredictiveInsights  `json:"predictive_insights"`
	MetaInsights        *MetaInsights        `json:"meta_insights"`
	PersonalizedAdvice  []PersonalizedTip    `json:"personalized_advice"`
}

// PerformanceInsights provides performance-specific insights
type PerformanceInsights struct {
	CurrentForm         string            `json:"current_form"`
	PerformancePattern  string            `json:"performance_pattern"`
	ConsistencyRating   float64           `json:"consistency_rating"`
	ClutchFactor        float64           `json:"clutch_factor"`
	AdaptabilityScore   float64           `json:"adaptability_score"`
	KeyPerformanceAreas []PerformanceArea `json:"key_performance_areas"`
	ImprovementVelocity float64           `json:"improvement_velocity"`
}

// PerformanceArea represents a specific performance area
type PerformanceArea struct {
	Area      string  `json:"area"`
	Score     float64 `json:"score"`
	Rank      string  `json:"rank"` // Relative ranking
	Trend     string  `json:"trend"`
	Potential float64 `json:"potential"` // Improvement potential
}

// CoachingInsights provides coaching-specific insights
type CoachingInsights struct {
	LearningStyle      string           `json:"learning_style"`
	MotivationFactors  []string         `json:"motivation_factors"`
	LearningObstacles  []string         `json:"learning_obstacles"`
	OptimalPlayTimes   []string         `json:"optimal_play_times"`
	TiltTriggers       []string         `json:"tilt_triggers"`
	SuccessPatterns    []SuccessPattern `json:"success_patterns"`
	RecommendedRoutine string           `json:"recommended_routine"`
}

// SuccessPattern identifies patterns in successful games
type SuccessPattern struct {
	Pattern     string  `json:"pattern"`
	Frequency   float64 `json:"frequency"`
	Impact      float64 `json:"impact"`
	Description string  `json:"description"`
}

// PredictiveInsights provides future performance predictions
type PredictiveInsights struct {
	RankPrediction      *RankPrediction      `json:"rank_prediction"`
	PerformanceForecast *PerformanceForecast `json:"performance_forecast"`
	SeasonGoals         []SeasonGoal         `json:"season_goals"`
	RiskFactors         []RiskFactor         `json:"risk_factors"`
}

// RankPrediction predicts future rank
type RankPrediction struct {
	PredictedRank   string   `json:"predicted_rank"`
	Confidence      float64  `json:"confidence"`
	TimeFrame       string   `json:"time_frame"`
	KeyFactors      []string `json:"key_factors"`
	RequiredChanges []string `json:"required_changes"`
}

// PerformanceForecast predicts performance trends
type PerformanceForecast struct {
	ExpectedTrend      string             `json:"expected_trend"`
	MetricForecasts    map[string]float64 `json:"metric_forecasts"`
	InfluencingFactors []string           `json:"influencing_factors"`
}

// SeasonGoal represents a seasonal objective
type SeasonGoal struct {
	Goal                string  `json:"goal"`
	Achievability       string  `json:"achievability"` // "easy", "moderate", "challenging", "difficult"
	RequiredImprovement float64 `json:"required_improvement"`
	EstimatedTime       string  `json:"estimated_time"`
}

// RiskFactor represents potential obstacles
type RiskFactor struct {
	Risk        string  `json:"risk"`
	Probability float64 `json:"probability"`
	Impact      string  `json:"impact"`
	Mitigation  string  `json:"mitigation"`
}

// MetaInsights provides meta game insights
type MetaInsights struct {
	MetaAlignment        string                `json:"meta_alignment"`
	MetaChampions        []string              `json:"meta_champions"`
	CounterPicks         []CounterPickInfo     `json:"counter_picks"`
	BuildRecommendations []BuildRecommendation `json:"build_recommendations"`
	StrategyTips         []StrategyTip         `json:"strategy_tips"`
}

// CounterPickInfo provides counter pick suggestions
type CounterPickInfo struct {
	AgainstChampion string  `json:"against_champion"`
	CounterPick     string  `json:"counter_pick"`
	Effectiveness   float64 `json:"effectiveness"`
	Reasoning       string  `json:"reasoning"`
}

// BuildRecommendation suggests item builds
type BuildRecommendation struct {
	Champion    string   `json:"champion"`
	Situation   string   `json:"situation"`
	CoreItems   []string `json:"core_items"`
	Situational []string `json:"situational"`
	Reasoning   string   `json:"reasoning"`
}

// StrategyTip provides strategic advice
type StrategyTip struct {
	Tip         string `json:"tip"`
	Category    string `json:"category"`
	Difficulty  int    `json:"difficulty"`
	Impact      string `json:"impact"`
	Application string `json:"application"`
}

// PersonalizedTip provides personalized advice
type PersonalizedTip struct {
	Tip         string   `json:"tip"`
	Priority    int      `json:"priority"`
	Category    string   `json:"category"`
	Reasoning   string   `json:"reasoning"`
	ActionSteps []string `json:"action_steps"`
}

// Progress tracking

// AnalysisProgressTracker tracks analysis progress for long-running operations
type AnalysisProgressTracker struct {
	RequestID     string    `json:"request_id"`
	CurrentStep   string    `json:"current_step"`
	Progress      int       `json:"progress"` // 0-100
	StartedAt     time.Time `json:"started_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	EstimatedTime string    `json:"estimated_time"`
	IsComplete    bool      `json:"is_complete"`
	Error         string    `json:"error,omitempty"`
}

// UpdateProgress updates the progress tracker
func (t *AnalysisProgressTracker) UpdateProgress(step string, progress int) {
	t.CurrentStep = step
	t.Progress = progress
	t.UpdatedAt = time.Now()
}

// Complete marks the tracker as complete
func (t *AnalysisProgressTracker) Complete() {
	t.IsComplete = true
	t.Progress = 100
	t.UpdatedAt = time.Now()
}

// SetError sets an error state
func (t *AnalysisProgressTracker) SetError(err error) {
	t.Error = err.Error()
	t.UpdatedAt = time.Now()
}
