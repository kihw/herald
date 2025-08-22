package models

import (
	"time"
	"github.com/lib/pq"
)

// User represents a registered user validated via Riot ID
type User struct {
	ID              int       `json:"id" db:"id"`
	RiotID          string    `json:"riot_id" db:"riot_id"`           // Riot ID (gameName)
	RiotTag         string    `json:"riot_tag" db:"riot_tag"`         // Riot Tag (tagLine)
	RiotPUUID       string    `json:"riot_puuid" db:"riot_puuid"`     // Primary identifier from Riot
	SummonerID      *string   `json:"summoner_id,omitempty" db:"summoner_id"`
	AccountID       *string   `json:"account_id,omitempty" db:"account_id"`
	ProfileIconID   int       `json:"profile_icon_id" db:"profile_icon_id"`
	SummonerLevel   int       `json:"summoner_level" db:"summoner_level"`
	Region          string    `json:"region" db:"region"`             // Primary region
	IsValidated     bool      `json:"is_validated" db:"is_validated"` // Account validation status
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	LastSync        *time.Time `json:"last_sync,omitempty" db:"last_sync"`
}

// UserSettings represents user configuration for match synchronization
type UserSettings struct {
	ID                int           `json:"id" db:"id"`
	UserID            int           `json:"user_id" db:"user_id"`
	Platform          string        `json:"platform" db:"platform"`
	QueueTypes        pq.Int64Array `json:"queue_types" db:"queue_types"`
	Language          string        `json:"language" db:"language"`
	IncludeTimeline   bool          `json:"include_timeline" db:"include_timeline"`
	IncludeAllData    bool          `json:"include_all_data" db:"include_all_data"`
	LightMode         bool          `json:"light_mode" db:"light_mode"`
	AutoSyncEnabled   bool          `json:"auto_sync_enabled" db:"auto_sync_enabled"`
	SyncFrequencyHours int          `json:"sync_frequency_hours" db:"sync_frequency_hours"`
	CreatedAt         time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at" db:"updated_at"`
}

// Match represents a League of Legends match
type Match struct {
	ID               int             `json:"id" db:"id"`
	MatchID          string          `json:"match_id" db:"match_id"`
	Platform         string          `json:"platform" db:"platform"`
	GameCreation     int64           `json:"game_creation" db:"game_creation"`
	GameDuration     int             `json:"game_duration" db:"game_duration"`
	GameEndTimestamp *int64          `json:"game_end_timestamp,omitempty" db:"game_end_timestamp"`
	GameMode         *string         `json:"game_mode,omitempty" db:"game_mode"`
	GameType         *string         `json:"game_type,omitempty" db:"game_type"`
	GameVersion      *string         `json:"game_version,omitempty" db:"game_version"`
	MapID            *int            `json:"map_id,omitempty" db:"map_id"`
	QueueID          *int            `json:"queue_id,omitempty" db:"queue_id"`
	SeasonID         *int            `json:"season_id,omitempty" db:"season_id"`
	TournamentCode   *string         `json:"tournament_code,omitempty" db:"tournament_code"`
	DataVersion      *string         `json:"data_version,omitempty" db:"data_version"`
	RawData          *map[string]interface{} `json:"raw_data,omitempty" db:"raw_data"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
}

// MatchParticipant represents a user's participation in a match
type MatchParticipant struct {
	ID                          int                    `json:"id" db:"id"`
	MatchID                     int                    `json:"match_id" db:"match_id"`
	UserID                      int                    `json:"user_id" db:"user_id"`
	ParticipantID               int                    `json:"participant_id" db:"participant_id"`
	TeamID                      int                    `json:"team_id" db:"team_id"`
	ChampionID                  int                    `json:"champion_id" db:"champion_id"`
	ChampionName                *string                `json:"champion_name,omitempty" db:"champion_name"`
	ChampionLevel               *int                   `json:"champion_level,omitempty" db:"champion_level"`
	Kills                       int                    `json:"kills" db:"kills"`
	Deaths                      int                    `json:"deaths" db:"deaths"`
	Assists                     int                    `json:"assists" db:"assists"`
	TotalDamageDealt            int                    `json:"total_damage_dealt" db:"total_damage_dealt"`
	TotalDamageDealtToChampions int                    `json:"total_damage_dealt_to_champions" db:"total_damage_dealt_to_champions"`
	TotalDamageTaken            int                    `json:"total_damage_taken" db:"total_damage_taken"`
	GoldEarned                  int                    `json:"gold_earned" db:"gold_earned"`
	TotalMinionsKilled          int                    `json:"total_minions_killed" db:"total_minions_killed"`
	VisionScore                 int                    `json:"vision_score" db:"vision_score"`
	Item0                       int                    `json:"item0" db:"item0"`
	Item1                       int                    `json:"item1" db:"item1"`
	Item2                       int                    `json:"item2" db:"item2"`
	Item3                       int                    `json:"item3" db:"item3"`
	Item4                       int                    `json:"item4" db:"item4"`
	Item5                       int                    `json:"item5" db:"item5"`
	Item6                       int                    `json:"item6" db:"item6"`
	Win                         bool                   `json:"win" db:"win"`
	DetailedStats              *map[string]interface{} `json:"detailed_stats,omitempty" db:"detailed_stats"`
	CreatedAt                   time.Time              `json:"created_at" db:"created_at"`
}

// SyncJob represents a synchronization job
type SyncJob struct {
	ID                 int        `json:"id" db:"id"`
	UserID             int        `json:"user_id" db:"user_id"`
	JobType            string     `json:"job_type" db:"job_type"`
	Status             string     `json:"status" db:"status"`
	StartedAt          *time.Time `json:"started_at,omitempty" db:"started_at"`
	CompletedAt        *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	MatchesProcessed   int        `json:"matches_processed" db:"matches_processed"`
	MatchesNew         int        `json:"matches_new" db:"matches_new"`
	MatchesUpdated     int        `json:"matches_updated" db:"matches_updated"`
	ErrorMessage       *string    `json:"error_message,omitempty" db:"error_message"`
	LastMatchTimestamp *int64     `json:"last_match_timestamp,omitempty" db:"last_match_timestamp"`
}

// SystemConfig represents system configuration
type SystemConfig struct {
	ID          int       `json:"id" db:"id"`
	Key         string    `json:"key" db:"key"`
	Value       string    `json:"value" db:"value"`
	Description *string   `json:"description,omitempty" db:"description"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Request/Response DTOs

// Riot Account Validation DTOs
type RiotAccountRequest struct {
	RiotID string `json:"riot_id" binding:"required"`
	RiotTag string `json:"riot_tag" binding:"required"`
	Region string `json:"region" binding:"required"`
}

type RiotAccountValidationResponse struct {
	Valid         bool   `json:"valid"`
	User          *User  `json:"user,omitempty"`
	ErrorMessage  string `json:"error_message,omitempty"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	SessionID    string `json:"session_id"`
	ExpiresIn    int    `json:"expires_in"`
}

// Riot API response structures
type RiotAccountResponse struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type RiotSummonerResponse struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	PUUID         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

type SyncRequest struct {
	Force bool `json:"force"` // Force sync even if cooldown not reached
}

type SyncResponse struct {
	JobID   int    `json:"job_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type DashboardStats struct {
	TotalMatches    int     `json:"total_matches"`
	WinRate         float64 `json:"win_rate"`
	AverageKDA      float64 `json:"average_kda"`
	FavoriteChampion string  `json:"favorite_champion"`
	LastSyncAt      *time.Time `json:"last_sync_at"`
	NextSyncAt      *time.Time `json:"next_sync_at"`
}

type MatchSummary struct {
	Match       Match            `json:"match"`
	Participant MatchParticipant `json:"participant"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
}
