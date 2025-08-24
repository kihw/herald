package analytics

import (
	"time"

	"github.com/herald-lol/herald/backend/internal/riot"
)

// Herald.lol Gaming Analytics - Analytics Data Models
// Data structures for gaming analytics results and requests

// PlayerAnalysisRequest contains data needed for player analysis
type PlayerAnalysisRequest struct {
	SummonerID   string        `json:"summoner_id"`
	SummonerName string        `json:"summoner_name"`
	PlayerPUUID  string        `json:"player_puuid"`
	Region       string        `json:"region"`
	CurrentRank  string        `json:"current_rank"`
	Matches      []*riot.Match `json:"matches"`
	TimeFrame    string        `json:"time_frame"`    // "recent", "season", "all"
	AnalysisType string        `json:"analysis_type"` // "basic", "detailed", "pro"
}

// PlayerAnalysis contains comprehensive player analysis results
type PlayerAnalysis struct {
	SummonerID       string                      `json:"summoner_id"`
	SummonerName     string                      `json:"summoner_name"`
	Region           string                      `json:"region"`
	AnalyzedAt       time.Time                   `json:"analyzed_at"`
	TotalMatches     int                         `json:"total_matches"`
	PerformanceScore float64                     `json:"performance_score"` // 0-100
	CoreMetrics      *CoreMetrics                `json:"core_metrics"`
	RoleMetrics      map[string]*RolePerformance `json:"role_metrics"`
	ChampionMetrics  []*ChampionPerformance      `json:"champion_metrics"`
	Trends           *TrendAnalysis              `json:"trends"`
	Insights         *GameInsights               `json:"insights"`
	SkillAssessment  *SkillAssessment            `json:"skill_assessment"`
	Recommendations  *ImprovementRecommendations `json:"recommendations"`
	CompetitiveData  *CompetitiveAnalysis        `json:"competitive_data,omitempty"`
}

// CoreMetrics contains fundamental gaming performance metrics
type CoreMetrics struct {
	// KDA Metrics
	AverageKDA     float64 `json:"average_kda"`
	AverageKills   float64 `json:"average_kills"`
	AverageDeaths  float64 `json:"average_deaths"`
	AverageAssists float64 `json:"average_assists"`

	// Economic Metrics
	AverageCS      float64 `json:"average_cs"`
	CSPerMinute    float64 `json:"cs_per_minute"`
	AverageGold    int     `json:"average_gold"`
	GoldEfficiency float64 `json:"gold_efficiency"` // 0-2.0 scale

	// Combat Metrics
	AverageDamage int     `json:"average_damage"`
	DamageShare   float64 `json:"damage_share"` // 0-1.0 team damage share

	// Utility Metrics
	AverageVision float64 `json:"average_vision"`

	// Overall Metrics
	WinRate float64 `json:"win_rate"` // 0-1.0

	// Advanced Metrics
	KillParticipation float64 `json:"kill_participation"`
	ObjectiveControl  float64 `json:"objective_control"`
	LaneEfficiency    float64 `json:"lane_efficiency"`
	TeamfightPresence float64 `json:"teamfight_presence"`
}

// RolePerformance contains role-specific performance metrics
type RolePerformance struct {
	Role              string  `json:"role"`
	GamesPlayed       int     `json:"games_played"`
	WinRate           float64 `json:"win_rate"`
	AverageKDA        float64 `json:"average_kda"`
	AverageCS         float64 `json:"average_cs"`
	CSPerMinute       float64 `json:"cs_per_minute"`
	AverageDamage     int     `json:"average_damage"`
	AverageGold       int     `json:"average_gold"`
	AverageVision     float64 `json:"average_vision"`
	PerformanceRating float64 `json:"performance_rating"` // 0-100
	RoleRank          string  `json:"role_rank"`          // Estimated rank for this role
	Consistency       float64 `json:"consistency"`        // Performance consistency 0-1
	ImpactRating      float64 `json:"impact_rating"`      // Game impact 0-100
}

// ChampionPerformance contains champion-specific performance metrics
type ChampionPerformance struct {
	ChampionName     string    `json:"champion_name"`
	ChampionID       int       `json:"champion_id"`
	GamesPlayed      int       `json:"games_played"`
	WinRate          float64   `json:"win_rate"`
	AverageKDA       float64   `json:"average_kda"`
	AverageCS        float64   `json:"average_cs"`
	CSPerMinute      float64   `json:"cs_per_minute"`
	MasteryLevel     int       `json:"mastery_level"`
	PerformanceTrend string    `json:"performance_trend"` // "improving", "declining", "stable"
	LastPlayed       time.Time `json:"last_played"`
	RecentForm       string    `json:"recent_form"`   // "WWLWL" recent game results
	BanRate          float64   `json:"ban_rate"`      // How often it gets banned
	PickPriority     int       `json:"pick_priority"` // 1-5 recommendation priority
	PowerSpike       string    `json:"power_spike"`   // "early", "mid", "late"
	Difficulty       int       `json:"difficulty"`    // 1-5 champion difficulty
}

// TrendAnalysis contains performance trend analysis
type TrendAnalysis struct {
	WinRateTrend     string  `json:"win_rate_trend"` // "improving", "declining", "stable"
	KDATrend         string  `json:"kda_trend"`
	CSPerMinTrend    string  `json:"cs_per_min_trend"`
	VisionTrend      string  `json:"vision_trend"`
	DamageTrend      string  `json:"damage_trend"`
	PerformanceTrend string  `json:"performance_trend"` // Overall trend
	RecentWinRate    float64 `json:"recent_win_rate"`   // Last 20 games
	TrendConfidence  float64 `json:"trend_confidence"`  // 0-1 confidence score
	TrendPeriod      string  `json:"trend_period"`      // Time period analyzed
	TrendStrength    string  `json:"trend_strength"`    // "strong", "moderate", "weak"
	PeakPerformance  float64 `json:"peak_performance"`  // Best recent performance score
	ConsistencyScore float64 `json:"consistency_score"` // Performance consistency
}

// GameInsights contains AI-generated insights and recommendations
type GameInsights struct {
	StrengthAreas     []string `json:"strength_areas"`
	ImprovementAreas  []string `json:"improvement_areas"`
	PlaystyleProfile  string   `json:"playstyle_profile"` // "Aggressive", "Passive", "Balanced"
	RecommendedChamps []string `json:"recommended_champs"`
	CoachingTips      []string `json:"coaching_tips"`
	NextGoals         []string `json:"next_goals"`
	SkillLevel        string   `json:"skill_level"`     // Estimated skill level
	Confidence        float64  `json:"confidence"`      // Confidence in insights 0-1
	KeyFocus          string   `json:"key_focus"`       // Primary area to focus on
	PlaytimeAdvice    string   `json:"playtime_advice"` // When to play for best results
	MetaAlignment     string   `json:"meta_alignment"`  // How well aligned with current meta
}

// SkillAssessment contains detailed skill assessment
type SkillAssessment struct {
	OverallSkill    string             `json:"overall_skill"`    // "Bronze", "Silver", etc.
	SkillAreas      map[string]float64 `json:"skill_areas"`      // Individual skill ratings
	MechanicalSkill float64            `json:"mechanical_skill"` // 0-100
	GameKnowledge   float64            `json:"game_knowledge"`   // 0-100
	DecisionMaking  float64            `json:"decision_making"`  // 0-100
	Positioning     float64            `json:"positioning"`      // 0-100
	Teamwork        float64            `json:"teamwork"`         // 0-100
	Adaptability    float64            `json:"adaptability"`     // 0-100
	Consistency     float64            `json:"consistency"`      // 0-100
	ImprovementRate float64            `json:"improvement_rate"` // Rate of improvement
	SkillCeiling    string             `json:"skill_ceiling"`    // Estimated potential
	LearningStyle   string             `json:"learning_style"`   // "Visual", "Practice", "Analysis"
}

// ImprovementRecommendations contains specific improvement suggestions
type ImprovementRecommendations struct {
	ImmediateFocus  []string      `json:"immediate_focus"`  // Most important improvements
	ShortTermGoals  []string      `json:"short_term_goals"` // 1-2 week goals
	LongTermGoals   []string      `json:"long_term_goals"`  // 1+ month goals
	PracticeRoutine []string      `json:"practice_routine"` // Specific practice suggestions
	ChampionPool    []string      `json:"champion_pool"`    // Recommended champions
	RoleFocus       string        `json:"role_focus"`       // Primary role recommendation
	TrainingPlan    *TrainingPlan `json:"training_plan"`    // Detailed training plan
	ResourceLinks   []string      `json:"resource_links"`   // Educational resources
	CoachingAreas   []string      `json:"coaching_areas"`   // Areas that need coaching
	MentorshipValue string        `json:"mentorship_value"` // "High", "Medium", "Low"
}

// TrainingPlan contains structured training recommendations
type TrainingPlan struct {
	Duration          string             `json:"duration"`           // "2 weeks", "1 month"
	DailyTimeCommit   string             `json:"daily_time_commit"`  // "30 minutes", "1 hour"
	WeeklyGoals       []string           `json:"weekly_goals"`       // Goals per week
	PracticeExercises []PracticeExercise `json:"practice_exercises"` // Specific exercises
	MilestoneChecks   []string           `json:"milestone_checks"`   // Progress checkpoints
	SuccessMetrics    []string           `json:"success_metrics"`    // How to measure success
}

// PracticeExercise represents a specific training exercise
type PracticeExercise struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Duration    string   `json:"duration"`
	Difficulty  int      `json:"difficulty"` // 1-5
	Priority    int      `json:"priority"`   // 1-5
	Category    string   `json:"category"`   // "Mechanical", "Strategic", etc.
	Tools       []string `json:"tools"`      // Required tools/modes
	Frequency   string   `json:"frequency"`  // "Daily", "3x/week", etc.
}

// CompetitiveAnalysis contains competitive/ranked specific analysis
type CompetitiveAnalysis struct {
	CurrentRank     string         `json:"current_rank"`
	PeakRank        string         `json:"peak_rank"`
	RankProgression []RankSnapshot `json:"rank_progression"`
	LPGainAverage   int            `json:"lp_gain_average"`
	LPLossAverage   int            `json:"lp_loss_average"`
	PromoSuccess    float64        `json:"promo_success"`     // Promotion success rate
	RankStability   string         `json:"rank_stability"`    // "Stable", "Climbing", "Declining"
	SeasonGoal      string         `json:"season_goal"`       // Predicted season end rank
	ClimbingRate    string         `json:"climbing_rate"`     // "Fast", "Steady", "Slow"
	TiltFactor      float64        `json:"tilt_factor"`       // 0-1 tendency to tilt
	OptimalPlayTime []string       `json:"optimal_play_time"` // Best times to play
	DodgeRecommend  []string       `json:"dodge_recommend"`   // When to dodge
}

// RankSnapshot represents rank at a point in time
type RankSnapshot struct {
	Timestamp    time.Time `json:"timestamp"`
	Tier         string    `json:"tier"`
	Rank         string    `json:"rank"`
	LeaguePoints int       `json:"league_points"`
	Games        int       `json:"games"`
	WinRate      float64   `json:"win_rate"`
}

// SkillGap represents the gap between current and target performance
type SkillGap struct {
	KDAGap         float64 `json:"kda_gap"`
	CSGap          float64 `json:"cs_gap"`
	VisionGap      float64 `json:"vision_gap"`
	DamageGap      float64 `json:"damage_gap"`
	WinRateGap     float64 `json:"win_rate_gap"`
	OverallGap     float64 `json:"overall_gap"`     // 0-100
	EstimatedGames int     `json:"estimated_games"` // Games needed to close gap
	TimeEstimate   string  `json:"time_estimate"`   // "2 weeks", "1 month"
	Difficulty     string  `json:"difficulty"`      // "Easy", "Moderate", "Hard"
}

// MatchAnalysisResult contains detailed match analysis
type MatchAnalysisResult struct {
	MatchID          string  `json:"match_id"`
	GameMode         string  `json:"game_mode"`
	Champion         string  `json:"champion"`
	Role             string  `json:"role"`
	Duration         int     `json:"duration"`
	Outcome          string  `json:"outcome"`           // "Victory", "Defeat"
	PerformanceScore float64 `json:"performance_score"` // 0-100

	// Detailed metrics
	KDA         float64 `json:"kda"`
	CS          int     `json:"cs"`
	CSPerMinute float64 `json:"cs_per_minute"`
	Gold        int     `json:"gold"`
	Damage      int     `json:"damage"`
	DamageShare float64 `json:"damage_share"`
	Vision      int     `json:"vision"`

	// Game flow analysis
	LanePhase *PhaseAnalysis `json:"lane_phase"`
	MidGame   *PhaseAnalysis `json:"mid_game"`
	LateGame  *PhaseAnalysis `json:"late_game"`

	// Advanced insights
	KeyMoments     []KeyMoment `json:"key_moments"`
	Strengths      []string    `json:"strengths"`
	Weaknesses     []string    `json:"weaknesses"`
	LearningPoints []string    `json:"learning_points"`
	MVPRating      float64     `json:"mvp_rating"` // 0-10

	PlayedAt time.Time `json:"played_at"`
}

// PhaseAnalysis contains analysis of a specific game phase
type PhaseAnalysis struct {
	Phase             string   `json:"phase"`          // "Lane", "Mid", "Late"
	Duration          int      `json:"duration"`       // Phase duration in seconds
	CSAdvantage       int      `json:"cs_advantage"`   // CS difference
	GoldAdvantage     int      `json:"gold_advantage"` // Gold difference
	KillsInPhase      int      `json:"kills_in_phase"`
	DeathsInPhase     int      `json:"deaths_in_phase"`
	DamageInPhase     int      `json:"damage_in_phase"`
	PerformanceRating float64  `json:"performance_rating"` // 0-100
	KeyEvents         []string `json:"key_events"`
	Impact            string   `json:"impact"` // "High", "Medium", "Low"
}

// KeyMoment represents a significant moment in the game
type KeyMoment struct {
	Timestamp   int    `json:"timestamp"` // Game time in seconds
	Type        string `json:"type"`      // "Kill", "Death", "Objective", etc.
	Description string `json:"description"`
	Impact      string `json:"impact"`   // "Positive", "Negative", "Neutral"
	Learning    string `json:"learning"` // What can be learned
	Severity    int    `json:"severity"` // 1-5, importance level
}

// TeamAnalysis contains team-based analysis
type TeamAnalysis struct {
	TeamID          string       `json:"team_id"`
	TeamName        string       `json:"team_name"`
	Members         []TeamMember `json:"members"`
	TeamMetrics     *TeamMetrics `json:"team_metrics"`
	Synergy         *TeamSynergy `json:"synergy"`
	Recommendations []string     `json:"recommendations"`
	AnalyzedAt      time.Time    `json:"analyzed_at"`
}

// TeamMember represents a team member's role and performance
type TeamMember struct {
	SummonerName     string   `json:"summoner_name"`
	Role             string   `json:"role"`
	ChampionPool     []string `json:"champion_pool"`
	PerformanceScore float64  `json:"performance_score"`
	Consistency      float64  `json:"consistency"`
	TeamplayRating   float64  `json:"teamplay_rating"`
	Leadership       float64  `json:"leadership"`
}

// TeamMetrics contains team-level performance metrics
type TeamMetrics struct {
	TeamWinRate        float64 `json:"team_win_rate"`
	AvgGameDuration    int     `json:"avg_game_duration"`
	ObjectiveControl   float64 `json:"objective_control"`
	LaneSwapSuccess    float64 `json:"lane_swap_success"`
	TeamfightWinRate   float64 `json:"teamfight_win_rate"`
	CommunicationScore float64 `json:"communication_score"`
}

// TeamSynergy analyzes how well team members work together
type TeamSynergy struct {
	OverallSynergy     float64            `json:"overall_synergy"`
	RoleSynergy        map[string]float64 `json:"role_synergy"`
	ChampionSynergy    float64            `json:"champion_synergy"`
	PlaystyleMesh      string             `json:"playstyle_mesh"` // "Excellent", "Good", "Needs Work"
	CommunicationFit   float64            `json:"communication_fit"`
	StrategicAlignment float64            `json:"strategic_alignment"`
}

// AnalyticsResult represents the result of gaming analytics processing
type AnalyticsResult struct {
	PlayerID         string             `json:"player_id"`
	MatchID          string             `json:"match_id"`
	ProcessingTime   time.Duration      `json:"processing_time"`
	PerformanceScore float64            `json:"performance_score"`
	Metrics          map[string]float64 `json:"metrics"`
	Recommendations  []string           `json:"recommendations"`
	Errors           []string           `json:"errors"`
	Success          bool               `json:"success"`
}
