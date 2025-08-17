package models

import (
	"time"
)

// TrendDirection represents performance trend direction
type TrendDirection string

const (
	TrendImproving TrendDirection = "improving"
	TrendDeclining TrendDirection = "declining"
	TrendStable    TrendDirection = "stable"
)

// PerformanceMetrics contains comprehensive champion/role performance data
type PerformanceMetrics struct {
	GamesPlayed      int             `json:"games_played"`
	Wins             int             `json:"wins"`
	Losses           int             `json:"losses"`
	WinRate          float64         `json:"win_rate"`
	AvgKills         float64         `json:"avg_kills"`
	AvgDeaths        float64         `json:"avg_deaths"`
	AvgAssists       float64         `json:"avg_assists"`
	AvgKDA           float64         `json:"avg_kda"`
	AvgCSPerMin      float64         `json:"avg_cs_per_min"`
	AvgGoldPerMin    float64         `json:"avg_gold_per_min"`
	AvgDamagePerMin  float64         `json:"avg_damage_per_min"`
	AvgVisionScore   float64         `json:"avg_vision_score"`
	PerformanceScore float64         `json:"performance_score"`
	TrendDirection   TrendDirection  `json:"trend_direction"`
}

// PeriodStats contains statistics for a specific time period
type PeriodStats struct {
	Period          string                    `json:"period"`
	TotalGames      int                       `json:"total_games"`
	WinRate         float64                   `json:"win_rate"`
	AvgKDA          float64                   `json:"avg_kda"`
	BestRole        string                    `json:"best_role"`
	WorstRole       string                    `json:"worst_role"`
	TopChampions    []ChampionPerformance     `json:"top_champions"`
	RolePerformance map[string]PerformanceMetrics `json:"role_performance"`
	RecentTrend     string                    `json:"recent_trend"`
	Suggestions     []string                  `json:"suggestions"`
}

// ChampionPerformance represents champion performance data
type ChampionPerformance struct {
	ChampionID       int     `json:"champion_id"`
	ChampionName     string  `json:"champion_name"`
	Games            int     `json:"games"`
	WinRate          float64 `json:"win_rate"`
	PerformanceScore float64 `json:"performance_score"`
	AvgKDA           float64 `json:"avg_kda"`
}

// ChampionMasteryAnalysis contains champion mastery analysis
type ChampionMasteryAnalysis struct {
	ChampionID          int                      `json:"champion_id"`
	ChampionName        string                   `json:"champion_name"`
	GamesPlayed         int                      `json:"games_played"`
	PerformanceMetrics  PerformanceMetrics       `json:"performance_metrics"`
	MasteryScore        float64                  `json:"mastery_score"`
	BestGame            GameSummary              `json:"best_game"`
	WorstGame           GameSummary              `json:"worst_game"`
	ImprovementSuggestions []string              `json:"improvement_suggestions"`
	SkillProgression    SkillProgression         `json:"skill_progression"`
}

// GameSummary represents a summary of a single game
type GameSummary struct {
	MatchID          string    `json:"match_id"`
	Date             time.Time `json:"date"`
	Champion         string    `json:"champion"`
	Role             string    `json:"role"`
	KDA              string    `json:"kda"`
	CS               int       `json:"cs"`
	Gold             int       `json:"gold"`
	Damage           int       `json:"damage"`
	VisionScore      int       `json:"vision_score"`
	GameDuration     int       `json:"game_duration"`
	Win              bool      `json:"win"`
	PerformanceScore float64   `json:"performance_score"`
}

// SkillProgression represents skill progression over time
type SkillProgression struct {
	ProgressionData     []ProgressionPeriod `json:"progression_data"`
	OverallImprovement  float64             `json:"overall_improvement"`
	Trend               string              `json:"trend"`
	InsufficientData    bool                `json:"insufficient_data,omitempty"`
}

// ProgressionPeriod represents progression in a specific period
type ProgressionPeriod struct {
	Period           int     `json:"period"`
	Games            int     `json:"games"`
	WinRate          float64 `json:"win_rate"`
	AvgKDA           float64 `json:"avg_kda"`
	PerformanceScore float64 `json:"performance_score"`
}

// ImprovementSuggestion represents a suggestion for improvement
type ImprovementSuggestion struct {
	Type                string  `json:"type"`
	Title               string  `json:"title"`
	Description         string  `json:"description"`
	Priority            int     `json:"priority"`
	ExpectedImprovement string  `json:"expected_improvement"`
	ActionItems         []string `json:"action_items"`
	ChampionID          *int    `json:"champion_id,omitempty"`
	Role                *string `json:"role,omitempty"`
	TimePeriod          string  `json:"time_period"`
}

// PerformanceTrends contains detailed performance trends over time
type PerformanceTrends struct {
	DailyTrend         TrendMetrics `json:"daily_trend"`
	WeeklyTrend        TrendMetrics `json:"weekly_trend"`
	MonthlyTrend       TrendMetrics `json:"monthly_trend"`
	SeasonalTrend      TrendMetrics `json:"seasonal_trend"`
	ImprovementVelocity float64     `json:"improvement_velocity"`
	ConsistencyScore    float64     `json:"consistency_score"`
	PeakPerformance     PeakPerformancePeriod `json:"peak_performance"`
}

// TrendMetrics represents metrics for a specific trend period
type TrendMetrics struct {
	Trend    string  `json:"trend"`
	Games    int     `json:"games"`
	WinRate  float64 `json:"win_rate"`
	Wins     int     `json:"wins"`
	Losses   int     `json:"losses"`
}

// PeakPerformancePeriod represents the period of peak performance
type PeakPerformancePeriod struct {
	Period      string  `json:"period"`
	Performance float64 `json:"performance"`
	Games       int     `json:"games"`
}

// ChampionStats represents champion statistics stored in database
type ChampionStats struct {
	ID                  int       `json:"id"`
	UserID              int       `json:"user_id"`
	ChampionID          int       `json:"champion_id"`
	ChampionName        string    `json:"champion_name"`
	Role                string    `json:"role"`
	Season              string    `json:"season"`
	TimePeriod          string    `json:"time_period"`
	GamesPlayed         int       `json:"games_played"`
	Wins                int       `json:"wins"`
	Losses              int       `json:"losses"`
	WinRate             float64   `json:"win_rate"`
	AvgKills            float64   `json:"avg_kills"`
	AvgDeaths           float64   `json:"avg_deaths"`
	AvgAssists          float64   `json:"avg_assists"`
	AvgKDA              float64   `json:"avg_kda"`
	AvgCSPerMin         float64   `json:"avg_cs_per_min"`
	AvgGoldPerMin       float64   `json:"avg_gold_per_min"`
	AvgDamagePerMin     float64   `json:"avg_damage_per_min"`
	AvgVisionScore      float64   `json:"avg_vision_score"`
	PerformanceScore    float64   `json:"performance_score"`
	TrendDirection      string    `json:"trend_direction"`
	LastPlayed          time.Time `json:"last_played"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}