package models

import (
	"time"
)

// MatchData represents processed match data for analytics
type MatchData struct {
	ID           string `json:"id" db:"id"`
	MatchID      string `json:"match_id" db:"match_id"`
	PlayerID     string `json:"player_id" db:"player_id"`
	SummonerName string `json:"summoner_name" db:"summoner_name"`

	// Game Info
	GameMode     string    `json:"game_mode" db:"game_mode"`
	QueueType    string    `json:"queue_type" db:"queue_type"`
	GameDuration int       `json:"game_duration" db:"game_duration"` // in seconds
	GameVersion  string    `json:"game_version" db:"game_version"`
	Date         time.Time `json:"date" db:"created_at"`

	// Champion Info
	ChampionID   int    `json:"champion_id" db:"champion_id"`
	ChampionName string `json:"champion_name" db:"champion_name"`
	Position     string `json:"position" db:"position"`
	TeamID       int    `json:"team_id" db:"team_id"`

	// Core Performance Metrics
	Kills   int     `json:"kills" db:"kills"`
	Deaths  int     `json:"deaths" db:"deaths"`
	Assists int     `json:"assists" db:"assists"`
	KDA     float64 `json:"kda" db:"kda"`

	// Farming Metrics
	TotalCS              int     `json:"total_cs" db:"total_cs"`
	CSPerMinute          float64 `json:"cs_per_minute" db:"cs_per_minute"`
	NeutralMinionsKilled int     `json:"neutral_minions_killed" db:"neutral_minions_killed"`
	EnemyJungleCS        int     `json:"enemy_jungle_cs" db:"enemy_jungle_cs"`
	AllyJungleCS         int     `json:"ally_jungle_cs" db:"ally_jungle_cs"`

	// Vision Metrics
	VisionScore        int `json:"vision_score" db:"vision_score"`
	WardsPlaced        int `json:"wards_placed" db:"wards_placed"`
	WardsKilled        int `json:"wards_killed" db:"wards_killed"`
	ControlWardsPlaced int `json:"control_wards_placed" db:"control_wards_placed"`
	VisionWardsBought  int `json:"vision_wards_bought" db:"vision_wards_bought"`

	// Damage Metrics
	TotalDamageDealt       int     `json:"total_damage_dealt" db:"total_damage_dealt"`
	TotalDamageToChampions int     `json:"total_damage_to_champions" db:"total_damage_to_champions"`
	PhysicalDamageDealt    int     `json:"physical_damage_dealt" db:"physical_damage_dealt"`
	MagicDamageDealt       int     `json:"magic_damage_dealt" db:"magic_damage_dealt"`
	TrueDamageDealt        int     `json:"true_damage_dealt" db:"true_damage_dealt"`
	DamageShare            float64 `json:"damage_share" db:"damage_share"`

	// Damage Taken
	TotalDamageTaken    int `json:"total_damage_taken" db:"total_damage_taken"`
	PhysicalDamageTaken int `json:"physical_damage_taken" db:"physical_damage_taken"`
	MagicDamageTaken    int `json:"magic_damage_taken" db:"magic_damage_taken"`
	TrueDamageTaken     int `json:"true_damage_taken" db:"true_damage_taken"`
	DamageSelfMitigated int `json:"damage_self_mitigated" db:"damage_self_mitigated"`

	// Gold & Economy
	GoldEarned     int     `json:"gold_earned" db:"gold_earned"`
	GoldSpent      int     `json:"gold_spent" db:"gold_spent"`
	GoldPerMinute  float64 `json:"gold_per_minute" db:"gold_per_minute"`
	GoldEfficiency float64 `json:"gold_efficiency" db:"gold_efficiency"`

	// Items
	Item0   int   `json:"item_0" db:"item_0"`
	Item1   int   `json:"item_1" db:"item_1"`
	Item2   int   `json:"item_2" db:"item_2"`
	Item3   int   `json:"item_3" db:"item_3"`
	Item4   int   `json:"item_4" db:"item_4"`
	Item5   int   `json:"item_5" db:"item_5"`
	Trinket int   `json:"trinket" db:"trinket"`
	Items   []int `json:"items"`

	// Objectives
	DragonKills    int `json:"dragon_kills" db:"dragon_kills"`
	BaronKills     int `json:"baron_kills" db:"baron_kills"`
	TurretKills    int `json:"turret_kills" db:"turret_kills"`
	InhibitorKills int `json:"inhibitor_kills" db:"inhibitor_kills"`

	// Performance Indicators
	FirstBloodKill   bool `json:"first_blood_kill" db:"first_blood_kill"`
	FirstBloodAssist bool `json:"first_blood_assist" db:"first_blood_assist"`
	FirstTowerKill   bool `json:"first_tower_kill" db:"first_tower_kill"`
	FirstTowerAssist bool `json:"first_tower_assist" db:"first_tower_assist"`

	// Advanced Metrics
	LongestTimeSpentLiving int `json:"longest_time_spent_living" db:"longest_time_spent_living"`
	LargestKillingSpree    int `json:"largest_killing_spree" db:"largest_killing_spree"`
	LargestMultiKill       int `json:"largest_multi_kill" db:"largest_multi_kill"`
	DoubleKills            int `json:"double_kills" db:"double_kills"`
	TripleKills            int `json:"triple_kills" db:"triple_kills"`
	QuadraKills            int `json:"quadra_kills" db:"quadra_kills"`
	PentaKills             int `json:"penta_kills" db:"penta_kills"`

	// Team Fighting
	TeamfightParticipation float64 `json:"teamfight_participation" db:"teamfight_participation"`

	// Game Result
	Win            bool `json:"win" db:"win"`
	GameEndedEarly bool `json:"game_ended_early" db:"game_ended_early"`
	GameWasRemade  bool `json:"game_was_remade" db:"game_was_remade"`

	// Timestamps for different phases
	EarlyGameEnd int `json:"early_game_end" db:"early_game_end"` // Usually 15 minutes
	MidGameEnd   int `json:"mid_game_end" db:"mid_game_end"`     // Usually 25 minutes

	// Calculated fields
	KillParticipation float64 `json:"kill_participation"`
	DeathShare        float64 `json:"death_share"`
	CSAtTenMinutes    int     `json:"cs_at_ten_minutes"`
	CSAt15Minutes     int     `json:"cs_at_fifteen_minutes"`
	GoldDiffAt15      int     `json:"gold_diff_at_fifteen"`
	XPDiffAt15        int     `json:"xp_diff_at_fifteen"`
}

// PlayerStats represents aggregated player statistics
type PlayerStats struct {
	PlayerID     string    `json:"player_id" db:"player_id"`
	SummonerName string    `json:"summoner_name" db:"summoner_name"`
	Region       string    `json:"region" db:"region"`
	LastUpdated  time.Time `json:"last_updated" db:"last_updated"`

	// Time Range
	TimeRange    string `json:"time_range" db:"time_range"`
	TotalMatches int    `json:"total_matches" db:"total_matches"`

	// Win Rate
	Wins    int     `json:"wins" db:"wins"`
	Losses  int     `json:"losses" db:"losses"`
	WinRate float64 `json:"win_rate" db:"win_rate"`

	// KDA Stats
	TotalKills   int     `json:"total_kills" db:"total_kills"`
	TotalDeaths  int     `json:"total_deaths" db:"total_deaths"`
	TotalAssists int     `json:"total_assists" db:"total_assists"`
	AverageKDA   float64 `json:"average_kda" db:"average_kda"`
	BestKDA      float64 `json:"best_kda" db:"best_kda"`

	// Farm Stats
	AverageCS       float64 `json:"average_cs" db:"average_cs"`
	AverageCSPerMin float64 `json:"average_cs_per_minute" db:"average_cs_per_minute"`
	BestCSPerMin    float64 `json:"best_cs_per_minute" db:"best_cs_per_minute"`

	// Vision Stats
	AverageVisionScore float64 `json:"average_vision_score" db:"average_vision_score"`
	AverageWardsPlaced float64 `json:"average_wards_placed" db:"average_wards_placed"`
	AverageWardsKilled float64 `json:"average_wards_killed" db:"average_wards_killed"`

	// Damage Stats
	AverageDamageShare float64 `json:"average_damage_share" db:"average_damage_share"`
	AverageDamageDealt float64 `json:"average_damage_dealt" db:"average_damage_dealt"`

	// Gold Stats
	AverageGoldPerMin     float64 `json:"average_gold_per_minute" db:"average_gold_per_minute"`
	AverageGoldEfficiency float64 `json:"average_gold_efficiency" db:"average_gold_efficiency"`

	// Champion Performance
	MainChampion         string  `json:"main_champion" db:"main_champion"`
	MainChampionWinRate  float64 `json:"main_champion_win_rate" db:"main_champion_win_rate"`
	TotalChampionsPlayed int     `json:"total_champions_played" db:"total_champions_played"`

	// Role Performance
	MainRole         string         `json:"main_role" db:"main_role"`
	RoleDistribution map[string]int `json:"role_distribution"`

	// Current Rank
	CurrentTier  string `json:"current_tier" db:"current_tier"`
	CurrentRank  string `json:"current_rank" db:"current_rank"`
	LeaguePoints int    `json:"league_points" db:"league_points"`

	// Performance Trends
	RecentForm       string  `json:"recent_form"`       // "improving", "declining", "stable"
	PerformanceScore float64 `json:"performance_score"` // 0-100
}

// ChampionStats represents champion-specific statistics
type ChampionStats struct {
	PlayerID     string `json:"player_id" db:"player_id"`
	ChampionID   int    `json:"champion_id" db:"champion_id"`
	ChampionName string `json:"champion_name" db:"champion_name"`

	// Basic Stats
	TotalMatches int     `json:"total_matches" db:"total_matches"`
	Wins         int     `json:"wins" db:"wins"`
	Losses       int     `json:"losses" db:"losses"`
	WinRate      float64 `json:"win_rate" db:"win_rate"`

	// Performance
	AverageKDA         float64 `json:"average_kda" db:"average_kda"`
	AverageCSPerMin    float64 `json:"average_cs_per_minute" db:"average_cs_per_minute"`
	AverageVisionScore float64 `json:"average_vision_score" db:"average_vision_score"`
	AverageDamageShare float64 `json:"average_damage_share" db:"average_damage_share"`

	// Mastery
	MasteryLevel  int `json:"mastery_level" db:"mastery_level"`
	MasteryPoints int `json:"mastery_points" db:"mastery_points"`

	// Last Played
	LastPlayed time.Time `json:"last_played" db:"last_played"`

	// Performance Grade
	Grade string `json:"grade"` // S+, S, S-, A+, A, A-, B+, B, B-, C+, C, D

	// Role specific stats
	PreferredRole     string                 `json:"preferred_role" db:"preferred_role"`
	RoleSpecificStats map[string]interface{} `json:"role_specific_stats"`
}

// MatchTimeline represents detailed match timeline data
type MatchTimeline struct {
	MatchID  string `json:"match_id" db:"match_id"`
	PlayerID string `json:"player_id" db:"player_id"`

	// Timeline events
	Events []TimelineEvent `json:"events"`

	// Performance by game phase
	EarlyGame GamePhaseStats `json:"early_game"`
	MidGame   GamePhaseStats `json:"mid_game"`
	LateGame  GamePhaseStats `json:"late_game"`

	// Gold and XP deltas
	GoldDeltas []DeltaPoint `json:"gold_deltas"`
	XPDeltas   []DeltaPoint `json:"xp_deltas"`
	CSDeltas   []DeltaPoint `json:"cs_deltas"`
}

// TimelineEvent represents an event in match timeline
type TimelineEvent struct {
	Timestamp int                    `json:"timestamp"`
	EventType string                 `json:"event_type"`
	PlayerID  string                 `json:"player_id"`
	Position  Position               `json:"position"`
	Data      map[string]interface{} `json:"data"`
}

// GamePhaseStats represents performance in a specific game phase
type GamePhaseStats struct {
	StartTime   int `json:"start_time"`
	EndTime     int `json:"end_time"`
	Kills       int `json:"kills"`
	Deaths      int `json:"deaths"`
	Assists     int `json:"assists"`
	CS          int `json:"cs"`
	Gold        int `json:"gold"`
	XP          int `json:"xp"`
	DamageDealt int `json:"damage_dealt"`
	DamageTaken int `json:"damage_taken"`
}

// DeltaPoint represents a point-in-time comparison
type DeltaPoint struct {
	Timestamp int     `json:"timestamp"`
	Value     float64 `json:"value"`
	Opponent  float64 `json:"opponent"`
	Delta     float64 `json:"delta"`
}

// Position represents a map position
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// RankedStats represents ranked-specific statistics
type RankedStats struct {
	PlayerID  string `json:"player_id" db:"player_id"`
	QueueType string `json:"queue_type" db:"queue_type"` // RANKED_SOLO_5x5, RANKED_FLEX_SR

	// Current Rank
	Tier         string `json:"tier" db:"tier"`
	Rank         string `json:"rank" db:"rank"`
	LeaguePoints int    `json:"league_points" db:"league_points"`

	// Series Information
	InSeries     bool   `json:"in_series" db:"in_series"`
	SeriesWins   int    `json:"series_wins" db:"series_wins"`
	SeriesLosses int    `json:"series_losses" db:"series_losses"`
	SeriesTarget string `json:"series_target" db:"series_target"`

	// Season Stats
	Wins    int     `json:"wins" db:"wins"`
	Losses  int     `json:"losses" db:"losses"`
	WinRate float64 `json:"win_rate" db:"win_rate"`

	// LP Tracking
	LPGains       []LPChange `json:"lp_gains"`
	LPLosses      []LPChange `json:"lp_losses"`
	AverageLPGain float64    `json:"average_lp_gain"`
	AverageLPLoss float64    `json:"average_lp_loss"`

	// Rank Progression
	RankHistory           []RankChange `json:"rank_history"`
	HighestRankThisSeason string       `json:"highest_rank_this_season"`
	LowestRankThisSeason  string       `json:"lowest_rank_this_season"`

	// Performance in Rank
	PerformanceInTier float64 `json:"performance_in_tier"` // Compared to other players in same tier

	LastUpdated time.Time `json:"last_updated" db:"last_updated"`
}

// LPChange represents a change in League Points
type LPChange struct {
	MatchID    string    `json:"match_id"`
	Change     int       `json:"change"`
	PreviousLP int       `json:"previous_lp"`
	NewLP      int       `json:"new_lp"`
	Timestamp  time.Time `json:"timestamp"`
}

// RankChange represents a rank change event
type RankChange struct {
	PreviousTier string    `json:"previous_tier"`
	PreviousRank string    `json:"previous_rank"`
	NewTier      string    `json:"new_tier"`
	NewRank      string    `json:"new_rank"`
	ChangeType   string    `json:"change_type"` // "promotion", "demotion", "tier_change"
	MatchID      string    `json:"match_id"`
	Timestamp    time.Time `json:"timestamp"`
}

// PerformanceBenchmark represents performance benchmarks for comparison
type PerformanceBenchmark struct {
	Tier     string `json:"tier"`
	Role     string `json:"role"`
	Champion string `json:"champion,omitempty"`

	// Sample Size
	TotalPlayers int `json:"total_players"`
	TotalMatches int `json:"total_matches"`

	// Benchmark Stats
	AverageKDA      float64 `json:"average_kda"`
	MedianKDA       float64 `json:"median_kda"`
	Top10PercentKDA float64 `json:"top_10_percent_kda"`

	AverageCSPerMin      float64 `json:"average_cs_per_minute"`
	MedianCSPerMin       float64 `json:"median_cs_per_minute"`
	Top10PercentCSPerMin float64 `json:"top_10_percent_cs_per_minute"`

	AverageVisionScore float64 `json:"average_vision_score"`
	AverageDamageShare float64 `json:"average_damage_share"`
	AverageWinRate     float64 `json:"average_win_rate"`

	// Last Updated
	LastCalculated time.Time `json:"last_calculated"`
}

// TableNames for GORM
func (MatchData) TableName() string {
	return "match_data"
}

func (PlayerStats) TableName() string {
	return "player_stats"
}

func (ChampionStats) TableName() string {
	return "champion_stats"
}

func (MatchTimeline) TableName() string {
	return "match_timelines"
}

func (RankedStats) TableName() string {
	return "ranked_stats"
}

func (PerformanceBenchmark) TableName() string {
	return "performance_benchmarks"
}
