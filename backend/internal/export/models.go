package export

import (
	"time"

	"github.com/herald-lol/herald/backend/internal/analytics"
	"github.com/herald-lol/herald/backend/internal/match"
)

// Herald.lol Gaming Analytics - Export Data Models
// Data structures for gaming data export service

// Export Request Models

// PlayerExportRequest contains parameters for exporting player analytics data
type PlayerExportRequest struct {
	PlayerPUUID        string   `json:"player_puuid" validate:"required"`
	SummonerName       string   `json:"summoner_name"`
	Region             string   `json:"region" validate:"required"`
	Format             string   `json:"format" validate:"required,oneof=csv json xlsx pdf charts"`
	TimeRange          string   `json:"time_range" validate:"required"`
	GameModes          []string `json:"game_modes"`
	MatchIDs           []string `json:"match_ids"`
	IncludeDetails     bool     `json:"include_details"`
	IncludeCharts      bool     `json:"include_charts"`
	IncludeComparisons bool     `json:"include_comparisons"`
	CustomFields       []string `json:"custom_fields"`
	CompressionLevel   string   `json:"compression_level,omitempty"`
	EncryptionEnabled  bool     `json:"encryption_enabled,omitempty"`

	// Export options
	ExportOptions *ExportOptions `json:"export_options,omitempty"`
}

// MatchExportRequest contains parameters for exporting match analysis
type MatchExportRequest struct {
	MatchID           string   `json:"match_id" validate:"required"`
	PlayerPUUID       string   `json:"player_puuid" validate:"required"`
	Format            string   `json:"format" validate:"required,oneof=csv json xlsx pdf charts"`
	AnalysisDepth     string   `json:"analysis_depth,omitempty"`
	IncludeTimeline   bool     `json:"include_timeline"`
	IncludeHeatmaps   bool     `json:"include_heatmaps"`
	IncludeComparison bool     `json:"include_comparison"`
	ComparisonTargets []string `json:"comparison_targets"`

	// Export options
	ExportOptions *ExportOptions `json:"export_options,omitempty"`
}

// TeamExportRequest contains parameters for exporting team analytics
type TeamExportRequest struct {
	TeamName           string   `json:"team_name" validate:"required"`
	PlayerPUUIDs       []string `json:"player_puuids" validate:"required,min=2,max=5"`
	Format             string   `json:"format" validate:"required,oneof=csv json xlsx pdf charts"`
	TimeRange          string   `json:"time_range" validate:"required"`
	GameModes          []string `json:"game_modes"`
	SharedMatchIDs     []string `json:"shared_match_ids"`
	IncludeIndividual  bool     `json:"include_individual"`
	IncludeTeamMetrics bool     `json:"include_team_metrics"`
	IncludeSynergy     bool     `json:"include_synergy"`

	// Export options
	ExportOptions *ExportOptions `json:"export_options,omitempty"`
}

// ChampionExportRequest contains parameters for exporting champion analytics
type ChampionExportRequest struct {
	PlayerPUUID        string   `json:"player_puuid" validate:"required"`
	ChampionName       string   `json:"champion_name" validate:"required"`
	Format             string   `json:"format" validate:"required,oneof=csv json xlsx pdf charts"`
	TimeRange          string   `json:"time_range" validate:"required"`
	GameModes          []string `json:"game_modes"`
	IncludeBuildPaths  bool     `json:"include_build_paths"`
	IncludeMatchups    bool     `json:"include_matchups"`
	IncludeProgression bool     `json:"include_progression"`
	ComparisonPlayers  []string `json:"comparison_players"`

	// Export options
	ExportOptions *ExportOptions `json:"export_options,omitempty"`
}

// CustomReportRequest contains parameters for custom reports
type CustomReportRequest struct {
	ReportName   string                 `json:"report_name" validate:"required"`
	ReportType   string                 `json:"report_type" validate:"required"`
	Description  string                 `json:"description"`
	Format       string                 `json:"format" validate:"required,oneof=csv json xlsx pdf charts"`
	Query        string                 `json:"query"`
	Parameters   map[string]interface{} `json:"parameters"`
	Columns      []string               `json:"columns"`
	Filters      []ReportFilter         `json:"filters"`
	Aggregations []ReportAggregation    `json:"aggregations"`
	TimeRange    string                 `json:"time_range"`
	GameModes    []string               `json:"game_modes"`

	// Export options
	ExportOptions *ExportOptions `json:"export_options,omitempty"`
}

// Export Response Models

// ExportResult contains the result of an export operation
type ExportResult struct {
	ExportID    string          `json:"export_id"`
	Format      string          `json:"format"`
	FileSize    int             `json:"file_size"`
	Status      string          `json:"status"`
	DownloadURL string          `json:"download_url"`
	CreatedAt   time.Time       `json:"created_at"`
	ExpiresAt   time.Time       `json:"expires_at"`
	Metadata    *ExportMetadata `json:"metadata,omitempty"`
}

// ExportStatus contains the status of an export job
type ExportStatus struct {
	ExportID     string    `json:"export_id"`
	Status       string    `json:"status"`   // pending, processing, completed, failed, expired
	Progress     int       `json:"progress"` // 0-100
	FileSize     int       `json:"file_size"`
	DownloadURL  string    `json:"download_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	ErrorMessage string    `json:"error_message,omitempty"`
}

// ExportSummary contains a summary of an export
type ExportSummary struct {
	ExportID    string    `json:"export_id"`
	Format      string    `json:"format"`
	Status      string    `json:"status"`
	FileSize    int       `json:"file_size"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	DataType    string    `json:"data_type"`
	Description string    `json:"description"`
}

// ExportMetadata contains metadata about an export
type ExportMetadata struct {
	PlayerPUUID  string `json:"player_puuid,omitempty"`
	MatchID      string `json:"match_id,omitempty"`
	TeamName     string `json:"team_name,omitempty"`
	ChampionName string `json:"champion_name,omitempty"`
	ReportName   string `json:"report_name,omitempty"`
	TimeRange    string `json:"time_range,omitempty"`
	DataPoints   int    `json:"data_points"`
	Compressed   bool   `json:"compressed"`
	Encrypted    bool   `json:"encrypted"`
	CustomQuery  string `json:"custom_query,omitempty"`
}

// Export Data Models

// PlayerExportData contains comprehensive player data for export
type PlayerExportData struct {
	PlayerInfo *PlayerInfo                   `json:"player_info"`
	Summary    *analytics.PerformanceSummary `json:"summary"`
	Matches    []*MatchExportData            `json:"matches"`
	TimeRange  string                        `json:"time_range"`
	ExportedAt time.Time                     `json:"exported_at"`
	TotalGames int                           `json:"total_games"`

	// Additional analytics data
	ChampionStats map[string]*ChampionStats `json:"champion_stats,omitempty"`
	RoleStats     map[string]*RoleStats     `json:"role_stats,omitempty"`
	TrendAnalysis *TrendAnalysis            `json:"trend_analysis,omitempty"`
	Achievements  []*Achievement            `json:"achievements,omitempty"`
}

// MatchExportData contains detailed match data for export
type MatchExportData struct {
	MatchID               string                       `json:"match_id"`
	Champion              string                       `json:"champion"`
	Role                  string                       `json:"role"`
	Result                string                       `json:"result"`
	Duration              int                          `json:"duration"`
	Performance           *match.PerformanceAnalysis   `json:"performance"`
	PhaseAnalysis         *match.GamePhaseAnalysis     `json:"phase_analysis,omitempty"`
	KeyMoments            []*match.KeyMoment           `json:"key_moments,omitempty"`
	TeamAnalysis          *match.TeamAnalysis          `json:"team_analysis,omitempty"`
	Insights              *match.MatchInsights         `json:"insights,omitempty"`
	LearningOpportunities []*match.LearningOpportunity `json:"learning_opportunities,omitempty"`
	OverallRating         float64                      `json:"overall_rating"`

	// Timeline data (optional)
	Timeline *MatchTimeline `json:"timeline,omitempty"`
	Heatmaps *MatchHeatmaps `json:"heatmaps,omitempty"`
}

// TeamExportData contains team analytics data for export
type TeamExportData struct {
	TeamName    string              `json:"team_name"`
	Players     []*PlayerExportData `json:"players"`
	TeamMetrics *TeamMetrics        `json:"team_metrics"`
	TimeRange   string              `json:"time_range"`
	ExportedAt  time.Time           `json:"exported_at"`

	// Team-specific analytics
	SynergyAnalysis    *SynergyAnalysis `json:"synergy_analysis,omitempty"`
	RoleEfficiency     *RoleEfficiency  `json:"role_efficiency,omitempty"`
	CommunicationScore float64          `json:"communication_score,omitempty"`
}

// ChampionExportData contains champion-specific analytics for export
type ChampionExportData struct {
	ChampionName       string                        `json:"champion_name"`
	PlayerPUUID        string                        `json:"player_puuid"`
	PerformanceHistory []*ChampionPerformanceHistory `json:"performance_history"`
	Statistics         *ChampionStatistics           `json:"statistics"`
	Trends             *ChampionTrends               `json:"trends"`
	Comparisons        *ChampionComparisons          `json:"comparisons"`
	Recommendations    []*ChampionRecommendation     `json:"recommendations"`
	TimeRange          string                        `json:"time_range"`
	ExportedAt         time.Time                     `json:"exported_at"`

	// Build and matchup data
	BuildPaths      []*BuildPath         `json:"build_paths,omitempty"`
	Matchups        []*ChampionMatchup   `json:"matchups,omitempty"`
	ProgressionData *ChampionProgression `json:"progression_data,omitempty"`
}

// CustomReportData contains custom report data for export
type CustomReportData struct {
	ReportName   string                   `json:"report_name"`
	Description  string                   `json:"description"`
	Parameters   map[string]interface{}   `json:"parameters"`
	Columns      []string                 `json:"columns"`
	DataRows     []map[string]interface{} `json:"data_rows"`
	Filters      []ReportFilter           `json:"filters"`
	Aggregations []ReportAggregation      `json:"aggregations"`
	GeneratedAt  time.Time                `json:"generated_at"`

	// Report metadata
	TotalRows   int           `json:"total_rows"`
	QueryTime   time.Duration `json:"query_time"`
	DataSources []string      `json:"data_sources"`
}

// Supporting Data Models

// PlayerInfo contains basic player information
type PlayerInfo struct {
	PUUID        string `json:"puuid"`
	SummonerName string `json:"summoner_name"`
	Region       string `json:"region"`
	Rank         string `json:"rank"`
	LP           int    `json:"lp"`
	Level        int    `json:"level"`
	ProfileIcon  int    `json:"profile_icon"`
}

// ChampionStats contains statistics for a specific champion
type ChampionStats struct {
	ChampionName  string  `json:"champion_name"`
	GamesPlayed   int     `json:"games_played"`
	WinRate       float64 `json:"win_rate"`
	AverageKDA    float64 `json:"average_kda"`
	AverageCS     float64 `json:"average_cs"`
	AverageDamage int     `json:"average_damage"`
	PlayRate      float64 `json:"play_rate"`
	Performance   string  `json:"performance"` // Excellent, Good, Average, Poor
}

// RoleStats contains statistics for a specific role
type RoleStats struct {
	Role          string  `json:"role"`
	GamesPlayed   int     `json:"games_played"`
	WinRate       float64 `json:"win_rate"`
	AverageKDA    float64 `json:"average_kda"`
	AverageCS     float64 `json:"average_cs"`
	AverageVision float64 `json:"average_vision"`
	PlayRate      float64 `json:"play_rate"`
	Performance   string  `json:"performance"`
}

// TrendAnalysis contains trend analysis data
type TrendAnalysis struct {
	PerformanceTrend string            `json:"performance_trend"` // Improving, Stable, Declining
	WinRateTrend     string            `json:"win_rate_trend"`
	KDATrend         string            `json:"kda_trend"`
	CSPerMinTrend    string            `json:"cs_per_min_trend"`
	TrendData        []*TrendDataPoint `json:"trend_data"`
	Consistency      float64           `json:"consistency"` // 0-100
	Volatility       float64           `json:"volatility"`
}

// TrendDataPoint represents a data point in trend analysis
type TrendDataPoint struct {
	Date        time.Time `json:"date"`
	WinRate     float64   `json:"win_rate"`
	KDA         float64   `json:"kda"`
	CSPerMin    float64   `json:"cs_per_min"`
	VisionScore float64   `json:"vision_score"`
	Rating      float64   `json:"rating"`
}

// Achievement represents a gaming achievement
type Achievement struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Rarity      string    `json:"rarity"` // Common, Rare, Epic, Legendary
	UnlockedAt  time.Time `json:"unlocked_at"`
	Category    string    `json:"category"`
	Points      int       `json:"points"`
}

// TeamMetrics contains team-level performance metrics
type TeamMetrics struct {
	TeamWinRate         float64            `json:"team_win_rate"`
	AverageTeamKDA      float64            `json:"average_team_kda"`
	AverageGameDuration int                `json:"average_game_duration"`
	ObjectiveControl    float64            `json:"objective_control"`
	TeamFightRating     float64            `json:"team_fight_rating"`
	MacroPlay           float64            `json:"macro_play"`
	PlayerSynergy       map[string]float64 `json:"player_synergy"` // Player pairs
}

// SynergyAnalysis contains team synergy analysis
type SynergyAnalysis struct {
	OverallSynergy     float64                   `json:"overall_synergy"`
	PlayerCombinations map[string]*PlayerSynergy `json:"player_combinations"`
	BestDuos           []*PlayerDuo              `json:"best_duos"`
	CommunicationScore float64                   `json:"communication_score"`
	CoordinationScore  float64                   `json:"coordination_score"`
}

// PlayerSynergy represents synergy between two players
type PlayerSynergy struct {
	Player1         string  `json:"player1"`
	Player2         string  `json:"player2"`
	SynergyScore    float64 `json:"synergy_score"`
	GamesPlayed     int     `json:"games_played"`
	WinRateTogether float64 `json:"win_rate_together"`
	KDAImprovement  float64 `json:"kda_improvement"`
}

// PlayerDuo represents a strong player duo
type PlayerDuo struct {
	Player1     string  `json:"player1"`
	Player2     string  `json:"player2"`
	Synergy     float64 `json:"synergy"`
	WinRate     float64 `json:"win_rate"`
	GamesPlayed int     `json:"games_played"`
}

// RoleEfficiency contains role efficiency analysis
type RoleEfficiency struct {
	TopLane     *RoleEfficiencyData `json:"top_lane"`
	Jungle      *RoleEfficiencyData `json:"jungle"`
	MidLane     *RoleEfficiencyData `json:"mid_lane"`
	BotLane     *RoleEfficiencyData `json:"bot_lane"`
	Support     *RoleEfficiencyData `json:"support"`
	OverallTeam float64             `json:"overall_team"`
}

// RoleEfficiencyData contains efficiency data for a role
type RoleEfficiencyData struct {
	PlayerName       string   `json:"player_name"`
	Efficiency       float64  `json:"efficiency"` // 0-100
	StrengthAreas    []string `json:"strength_areas"`
	ImprovementAreas []string `json:"improvement_areas"`
	RoleAlignment    float64  `json:"role_alignment"` // How well suited for role
}

// ChampionPerformanceHistory contains historical performance data
type ChampionPerformanceHistory struct {
	Date             time.Time `json:"date"`
	MatchID          string    `json:"match_id"`
	Result           string    `json:"result"`
	KDA              float64   `json:"kda"`
	CSPerMin         float64   `json:"cs_per_min"`
	DamageShare      float64   `json:"damage_share"`
	VisionScore      int       `json:"vision_score"`
	GameDuration     int       `json:"game_duration"`
	Rating           float64   `json:"rating"`
	OpponentChampion string    `json:"opponent_champion"`
}

// ChampionStatistics contains statistical data for champion performance
type ChampionStatistics struct {
	TotalGames          int     `json:"total_games"`
	WinRate             float64 `json:"win_rate"`
	AverageKDA          float64 `json:"average_kda"`
	AverageCSPerMin     float64 `json:"average_cs_per_min"`
	AverageDamageShare  float64 `json:"average_damage_share"`
	AverageVisionScore  float64 `json:"average_vision_score"`
	AverageGameDuration int     `json:"average_game_duration"`
	BestPerformance     float64 `json:"best_performance"`
	WorstPerformance    float64 `json:"worst_performance"`
	ConsistencyScore    float64 `json:"consistency_score"`

	// Rank-specific stats
	RankStats map[string]*ChampionRankStats `json:"rank_stats"`
}

// ChampionRankStats contains champion stats by rank
type ChampionRankStats struct {
	Rank        string  `json:"rank"`
	Games       int     `json:"games"`
	WinRate     float64 `json:"win_rate"`
	AverageKDA  float64 `json:"average_kda"`
	Performance string  `json:"performance"`
}

// ChampionTrends contains trend analysis for champion
type ChampionTrends struct {
	WinRateTrend      string    `json:"win_rate_trend"`
	KDATrend          string    `json:"kda_trend"`
	PerformanceTrend  string    `json:"performance_trend"`
	RecentImprovement bool      `json:"recent_improvement"`
	LearningCurve     string    `json:"learning_curve"` // Steep, Moderate, Flat
	MasteryLevel      string    `json:"mastery_level"`  // Beginner, Intermediate, Advanced, Master
	TrendStartDate    time.Time `json:"trend_start_date"`
}

// ChampionComparisons contains comparison data
type ChampionComparisons struct {
	VsAveragePlayer *ChampionComparisonData `json:"vs_average_player"`
	VsSimilarRank   *ChampionComparisonData `json:"vs_similar_rank"`
	VsHigherRank    *ChampionComparisonData `json:"vs_higher_rank"`
	PersonalBest    *ChampionComparisonData `json:"personal_best"`
}

// ChampionComparisonData contains specific comparison metrics
type ChampionComparisonData struct {
	WinRateDiff     float64 `json:"win_rate_diff"`
	KDADiff         float64 `json:"kda_diff"`
	CSPerMinDiff    float64 `json:"cs_per_min_diff"`
	DamageShareDiff float64 `json:"damage_share_diff"`
	VisionScoreDiff float64 `json:"vision_score_diff"`
	OverallDiff     float64 `json:"overall_diff"`
	Percentile      float64 `json:"percentile"` // 0-100
}

// ChampionRecommendation contains improvement recommendations
type ChampionRecommendation struct {
	Category       string   `json:"category"`
	Priority       string   `json:"priority"` // High, Medium, Low
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	ActionSteps    []string `json:"action_steps"`
	ExpectedImpact float64  `json:"expected_impact"` // Expected rating improvement
	TimeFrame      string   `json:"time_frame"`
}

// BuildPath contains champion build information
type BuildPath struct {
	BuildID      string             `json:"build_id"`
	ItemSequence []string           `json:"item_sequence"`
	WinRate      float64            `json:"win_rate"`
	GamesPlayed  int                `json:"games_played"`
	AverageCost  int                `json:"average_cost"`
	PowerSpikes  []*BuildPowerSpike `json:"power_spikes"`
	Situational  bool               `json:"situational"`
	Description  string             `json:"description"`
}

// BuildPowerSpike represents a power spike in build path
type BuildPowerSpike struct {
	ItemName       string  `json:"item_name"`
	CompletionTime int     `json:"completion_time"` // In minutes
	PowerIncrease  float64 `json:"power_increase"`
	Description    string  `json:"description"`
}

// ChampionMatchup contains matchup information
type ChampionMatchup struct {
	OpponentChampion string   `json:"opponent_champion"`
	GamesPlayed      int      `json:"games_played"`
	WinRate          float64  `json:"win_rate"`
	AverageKDA       float64  `json:"average_kda"`
	LaneDominance    float64  `json:"lane_dominance"` // -1 to 1
	Difficulty       string   `json:"difficulty"`     // Easy, Medium, Hard
	Tips             []string `json:"tips"`
	CounterPicks     []string `json:"counter_picks"`
}

// ChampionProgression contains progression tracking
type ChampionProgression struct {
	MasteryPoints    int                     `json:"mastery_points"`
	MasteryLevel     int                     `json:"mastery_level"`
	GamesPlayed      int                     `json:"games_played"`
	SkillProgression *SkillProgression       `json:"skill_progression"`
	Milestones       []*ProgressionMilestone `json:"milestones"`
	NextGoals        []*ProgressionGoal      `json:"next_goals"`
}

// SkillProgression tracks skill development
type SkillProgression struct {
	Farming      *SkillLevel `json:"farming"`
	Trading      *SkillLevel `json:"trading"`
	Teamfighting *SkillLevel `json:"teamfighting"`
	Positioning  *SkillLevel `json:"positioning"`
	MacroPlay    *SkillLevel `json:"macro_play"`
	Overall      *SkillLevel `json:"overall"`
}

// SkillLevel represents proficiency in a skill
type SkillLevel struct {
	Level      int       `json:"level"`    // 1-10
	Progress   float64   `json:"progress"` // 0-100% to next level
	Trend      string    `json:"trend"`    // Improving, Stable, Declining
	LastUpdate time.Time `json:"last_update"`
}

// ProgressionMilestone represents an achieved milestone
type ProgressionMilestone struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AchievedAt  time.Time `json:"achieved_at"`
	Difficulty  string    `json:"difficulty"`
	Category    string    `json:"category"`
}

// ProgressionGoal represents a future goal
type ProgressionGoal struct {
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	TargetValue   float64       `json:"target_value"`
	CurrentValue  float64       `json:"current_value"`
	Progress      float64       `json:"progress"` // 0-100%
	EstimatedTime time.Duration `json:"estimated_time"`
	Priority      string        `json:"priority"`
}

// MatchTimeline contains timeline data for match
type MatchTimeline struct {
	Intervals     []*TimelineInterval `json:"intervals"`
	KeyEvents     []*TimelineEvent    `json:"key_events"`
	GoldGraphData []*GraphDataPoint   `json:"gold_graph_data"`
	XPGraphData   []*GraphDataPoint   `json:"xp_graph_data"`
}

// TimelineInterval represents a time interval in match
type TimelineInterval struct {
	Timestamp int              `json:"timestamp"`
	Gold      int              `json:"gold"`
	XP        int              `json:"xp"`
	CS        int              `json:"cs"`
	Level     int              `json:"level"`
	Position  *MapPosition     `json:"position"`
	Items     []string         `json:"items"`
	Events    []*TimelineEvent `json:"events"`
}

// TimelineEvent represents an event in match timeline
type TimelineEvent struct {
	Type         string       `json:"type"`
	Timestamp    int          `json:"timestamp"`
	Description  string       `json:"description"`
	Position     *MapPosition `json:"position,omitempty"`
	Participants []string     `json:"participants,omitempty"`
	Value        int          `json:"value,omitempty"`
}

// MapPosition represents a position on the map
type MapPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// GraphDataPoint represents a data point for graphs
type GraphDataPoint struct {
	Timestamp int     `json:"timestamp"`
	Value     float64 `json:"value"`
}

// MatchHeatmaps contains heatmap data for match
type MatchHeatmaps struct {
	DeathHeatmap    *HeatmapData `json:"death_heatmap"`
	WardHeatmap     *HeatmapData `json:"ward_heatmap"`
	PositionHeatmap *HeatmapData `json:"position_heatmap"`
	DamageHeatmap   *HeatmapData `json:"damage_heatmap"`
}

// HeatmapData contains heatmap visualization data
type HeatmapData struct {
	DataPoints []*HeatmapPoint `json:"data_points"`
	Intensity  [][]float64     `json:"intensity"`
	MaxValue   float64         `json:"max_value"`
	ColorScale []string        `json:"color_scale"`
}

// HeatmapPoint represents a point in heatmap
type HeatmapPoint struct {
	X         int     `json:"x"`
	Y         int     `json:"y"`
	Intensity float64 `json:"intensity"`
	Count     int     `json:"count"`
}

// Report Configuration Models

// ReportFilter represents a filter for custom reports
type ReportFilter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // equals, not_equals, greater_than, less_than, contains, in
	Value    interface{} `json:"value"`
}

// ReportAggregation represents an aggregation for custom reports
type ReportAggregation struct {
	Field    string `json:"field"`
	Function string `json:"function"` // sum, avg, count, min, max, distinct
	Alias    string `json:"alias,omitempty"`
}

// Export Configuration Models

// ExportOptions contains options for export customization
type ExportOptions struct {
	IncludeCharts      bool     `json:"include_charts"`
	ChartTypes         []string `json:"chart_types"`
	IncludeRawData     bool     `json:"include_raw_data"`
	IncludeSummary     bool     `json:"include_summary"`
	IncludeComparisons bool     `json:"include_comparisons"`
	CustomSections     []string `json:"custom_sections"`
	BrandingEnabled    bool     `json:"branding_enabled"`
	WatermarkEnabled   bool     `json:"watermark_enabled"`

	// Format-specific options
	CSVOptions   *CSVExportOptions   `json:"csv_options,omitempty"`
	JSONOptions  *JSONExportOptions  `json:"json_options,omitempty"`
	XLSXOptions  *XLSXExportOptions  `json:"xlsx_options,omitempty"`
	PDFOptions   *PDFExportOptions   `json:"pdf_options,omitempty"`
	ChartOptions *ChartExportOptions `json:"chart_options,omitempty"`
}

// Format-specific export options
type CSVExportOptions struct {
	Delimiter      string `json:"delimiter"`
	IncludeHeaders bool   `json:"include_headers"`
	DateFormat     string `json:"date_format"`
	NumberFormat   string `json:"number_format"`
}

type JSONExportOptions struct {
	Pretty        bool   `json:"pretty"`
	IncludeSchema bool   `json:"include_schema"`
	DateFormat    string `json:"date_format"`
}

type XLSXExportOptions struct {
	MultipleSheets   bool     `json:"multiple_sheets"`
	IncludeCharts    bool     `json:"include_charts"`
	AutoFitColumns   bool     `json:"auto_fit_columns"`
	HeaderFormatting bool     `json:"header_formatting"`
	SheetNames       []string `json:"sheet_names,omitempty"`
}

type PDFExportOptions struct {
	PageSize      string      `json:"page_size"`   // A4, Letter, etc.
	Orientation   string      `json:"orientation"` // portrait, landscape
	IncludeHeader bool        `json:"include_header"`
	IncludeFooter bool        `json:"include_footer"`
	FontSize      int         `json:"font_size"`
	Margins       *PDFMargins `json:"margins"`
}

type PDFMargins struct {
	Top    float64 `json:"top"`
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
	Right  float64 `json:"right"`
}

type ChartExportOptions struct {
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Interactive bool     `json:"interactive"`
	Theme       string   `json:"theme"`
	ColorScheme []string `json:"color_scheme"`
	ShowLegend  bool     `json:"show_legend"`
	ShowGrid    bool     `json:"show_grid"`
}

// ExportFormat describes a supported export format
type ExportFormat struct {
	Name        string   `json:"name"`
	Key         string   `json:"key"`
	Description string   `json:"description"`
	Extensions  []string `json:"extensions"`
	MimeType    string   `json:"mime_type"`
	Features    []string `json:"features"`
}

// Cache Models

// CachedExport represents a cached export result
type CachedExport struct {
	ExportID    string    `json:"export_id"`
	Format      string    `json:"format"`
	FileSize    int       `json:"file_size"`
	DownloadURL string    `json:"download_url"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// Storage Models

// StoredExportInfo contains information about stored exports
type StoredExportInfo struct {
	ExportID     string    `json:"export_id"`
	Status       string    `json:"status"`
	Progress     int       `json:"progress"`
	FileSize     int       `json:"file_size"`
	DownloadURL  string    `json:"download_url"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	ErrorMessage string    `json:"error_message,omitempty"`
}

// UserExport represents an export belonging to a user
type UserExport struct {
	ExportID    string    `json:"export_id"`
	UserID      string    `json:"user_id"`
	Format      string    `json:"format"`
	Status      string    `json:"status"`
	FileSize    int       `json:"file_size"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	DataType    string    `json:"data_type"`
	Description string    `json:"description"`
}
