package models

import "time"

// MetagameSnapshot représente un snapshot des données métagame à un moment donné
type MetagameSnapshot struct {
	ID        int       `json:"id" db:"id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Patch     string    `json:"patch" db:"patch"`
	Region    string    `json:"region" db:"region"`
	Tier      string    `json:"tier" db:"tier"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ChampionMetrics représente les métriques d'un champion dans le métagame
type ChampionMetrics struct {
	ID             int       `json:"id" db:"id"`
	SnapshotID     int       `json:"snapshot_id" db:"snapshot_id"`
	ChampionID     int       `json:"champion_id" db:"champion_id"`
	ChampionName   string    `json:"champion_name" db:"champion_name"`
	Role           string    `json:"role" db:"role"`
	PickRate       float64   `json:"pick_rate" db:"pick_rate"`
	BanRate        float64   `json:"ban_rate" db:"ban_rate"`
	WinRate        float64   `json:"win_rate" db:"win_rate"`
	Presence       float64   `json:"presence" db:"presence"`
	GamesPlayed    int       `json:"games_played" db:"games_played"`
	Wins           int       `json:"wins" db:"wins"`
	Losses         int       `json:"losses" db:"losses"`
	AvgKDA         float64   `json:"avg_kda" db:"avg_kda"`
	AvgDamage      float64   `json:"avg_damage" db:"avg_damage"`
	AvgGold        float64   `json:"avg_gold" db:"avg_gold"`
	AvgCS          float64   `json:"avg_cs" db:"avg_cs"`
	TierScore      float64   `json:"tier_score" db:"tier_score"`
	TrendDirection string    `json:"trend_direction" db:"trend_direction"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// ItemMetrics représente les métriques d'un item dans le métagame
type ItemMetrics struct {
	ID              int       `json:"id" db:"id"`
	SnapshotID      int       `json:"snapshot_id" db:"snapshot_id"`
	ItemID          int       `json:"item_id" db:"item_id"`
	ItemName        string    `json:"item_name" db:"item_name"`
	PickRate        float64   `json:"pick_rate" db:"pick_rate"`
	WinRate         float64   `json:"win_rate" db:"win_rate"`
	Role            string    `json:"role" db:"role"`
	ChampionID      int       `json:"champion_id" db:"champion_id"`
	Position        string    `json:"position" db:"position"` // core, situational, boots, etc.
	AverageGameTime float64   `json:"average_game_time" db:"average_game_time"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// MetagameTrend représente une tendance dans le métagame
type MetagameTrend struct {
	ID            int       `json:"id" db:"id"`
	Type          string    `json:"type" db:"type"` // champion, item, strategy
	EntityID      int       `json:"entity_id" db:"entity_id"`
	EntityName    string    `json:"entity_name" db:"entity_name"`
	TrendType     string    `json:"trend_type" db:"trend_type"` // rising, falling, stable
	ChangePercent float64   `json:"change_percent" db:"change_percent"`
	Timeframe     string    `json:"timeframe" db:"timeframe"` // daily, weekly, monthly
	Confidence    float64   `json:"confidence" db:"confidence"`
	Description   string    `json:"description" db:"description"`
	StartDate     time.Time `json:"start_date" db:"start_date"`
	EndDate       time.Time `json:"end_date" db:"end_date"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// MetagameStats représente les statistiques globales du métagame
type MetagameStats struct {
	ID                 int       `json:"id" db:"id"`
	SnapshotID         int       `json:"snapshot_id" db:"snapshot_id"`
	TotalGames         int       `json:"total_games" db:"total_games"`
	UniqueChampions    int       `json:"unique_champions" db:"unique_champions"`
	AvgGameDuration    float64   `json:"avg_game_duration" db:"avg_game_duration"`
	MostPickedChampion string    `json:"most_picked_champion" db:"most_picked_champion"`
	MostBannedChampion string    `json:"most_banned_champion" db:"most_banned_champion"`
	HighestWinRate     float64   `json:"highest_win_rate" db:"highest_win_rate"`
	LowestWinRate      float64   `json:"lowest_win_rate" db:"lowest_win_rate"`
	DiversityIndex     float64   `json:"diversity_index" db:"diversity_index"`
	PowerLevel         string    `json:"power_level" db:"power_level"` // low, medium, high
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}
