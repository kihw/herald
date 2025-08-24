package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Match represents a League of Legends match
type Match struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Riot API Match Data
	MatchID    string `json:"match_id" gorm:"uniqueIndex;not null"` // Riot match ID
	GameID     int64  `json:"game_id"`
	PlatformID string `json:"platform_id" gorm:"not null"`

	// Game Information
	GameMode string `json:"game_mode" gorm:"not null"` // CLASSIC, ARAM, etc.
	GameType string `json:"game_type" gorm:"not null"` // MATCHED_GAME, etc.
	QueueID  int    `json:"queue_id" gorm:"not null"`  // 420 (Ranked Solo), 440 (Ranked Flex), etc.
	MapID    int    `json:"map_id" gorm:"not null"`    // 11 (Summoner's Rift), etc.

	// Timing
	GameStartTimestamp int64  `json:"game_start_timestamp"`
	GameEndTimestamp   int64  `json:"game_end_timestamp"`
	GameDuration       int    `json:"game_duration"` // seconds
	GameVersion        string `json:"game_version"`

	// Match Participants
	Participants []MatchParticipant `json:"participants" gorm:"foreignKey:MatchID"`

	// Match Status
	IsProcessed bool      `json:"is_processed" gorm:"default:false"`
	ProcessedAt time.Time `json:"processed_at"`
	IsAnalyzed  bool      `json:"is_analyzed" gorm:"default:false"`
	AnalyzedAt  time.Time `json:"analyzed_at"`

	// Game Outcome
	WinningTeam int `json:"winning_team"` // 100 (Blue) or 200 (Red)
}

// MatchParticipant represents a participant in a League of Legends match
type MatchParticipant struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	MatchID uuid.UUID `json:"match_id" gorm:"not null;index"`
	Match   Match     `json:"-" gorm:"foreignKey:MatchID"`

	// Player Information
	PUUID        string `json:"puuid" gorm:"not null;index"`
	SummonerName string `json:"summoner_name"`
	SummonerID   string `json:"summoner_id"`

	// Game Position
	ParticipantID int    `json:"participant_id"` // 1-10
	TeamID        int    `json:"team_id"`        // 100 (Blue) or 200 (Red)
	TeamPosition  string `json:"team_position"`  // TOP, JUNGLE, MIDDLE, BOTTOM, UTILITY

	// Champion Information
	ChampionID   int    `json:"champion_id" gorm:"not null"`
	ChampionName string `json:"champion_name" gorm:"not null"`

	// Summoner Spells
	Spell1ID int `json:"spell1_id"`
	Spell2ID int `json:"spell2_id"`

	// Runes
	PrimaryRuneStyle int `json:"primary_rune_style"`
	SubRuneStyle     int `json:"sub_rune_style"`

	// Core Stats
	Kills         int  `json:"kills"`
	Deaths        int  `json:"deaths"`
	Assists       int  `json:"assists"`
	ChampionLevel int  `json:"champion_level"`
	Won           bool `json:"won"`

	// Combat Stats
	TotalDamageDealt            int `json:"total_damage_dealt"`
	TotalDamageDealtToChampions int `json:"total_damage_dealt_to_champions"`
	TotalDamageTaken            int `json:"total_damage_taken"`
	TotalHeal                   int `json:"total_heal"`
	TotalHealsOnTeammates       int `json:"total_heals_on_teammates"`
	DamageDealtToObjectives     int `json:"damage_dealt_to_objectives"`
	DamageDealtToTurrets        int `json:"damage_dealt_to_turrets"`

	// Economy
	GoldEarned  int     `json:"gold_earned"`
	GoldSpent   int     `json:"gold_spent"`
	TotalCS     int     `json:"total_cs"` // totalMinionsKilled + neutralMinionsKilled
	CSPerMinute float64 `json:"cs_per_minute"`

	// Vision
	VisionScore             int `json:"vision_score"`
	WardsPlaced             int `json:"wards_placed"`
	WardsKilled             int `json:"wards_killed"`
	ControlWardsPlaced      int `json:"control_wards_placed"`
	VisionWardsBoughtInGame int `json:"vision_wards_bought_in_game"`

	// Items
	Item0 int `json:"item0"`
	Item1 int `json:"item1"`
	Item2 int `json:"item2"`
	Item3 int `json:"item3"`
	Item4 int `json:"item4"`
	Item5 int `json:"item5"`
	Item6 int `json:"item6"` // Trinket

	// Performance Metrics (calculated by Herald.lol)
	KDA               float64 `json:"kda"`
	KillParticipation float64 `json:"kill_participation"` // (kills + assists) / team kills
	DamageShare       float64 `json:"damage_share"`       // damage dealt to champions / team total
	GoldShare         float64 `json:"gold_share"`         // gold earned / team total
	PerformanceScore  float64 `json:"performance_score"`  // Herald.lol performance rating

	// Objectives
	TurretKills    int `json:"turret_kills"`
	InhibitorKills int `json:"inhibitor_kills"`
	DragonKills    int `json:"dragon_kills"`
	BaronKills     int `json:"baron_kills"`

	// Player Behavior
	FirstBloodKill      bool `json:"first_blood_kill"`
	FirstBloodAssist    bool `json:"first_blood_assist"`
	LargestKillingSpree int  `json:"largest_killing_spree"`
	LargestMultiKill    int  `json:"largest_multi_kill"`

	// Game Events
	EarlyGamePerformance float64 `json:"early_game_performance"` // 0-15 min
	MidGamePerformance   float64 `json:"mid_game_performance"`   // 15-25 min
	LateGamePerformance  float64 `json:"late_game_performance"`  // 25+ min

	// Position-Specific Metrics
	JungleCS         int     `json:"jungle_cs"`          // For junglers
	LaneCS           int     `json:"lane_cs"`            // For laners
	SupportItemQuest bool    `json:"support_item_quest"` // For supports
	RoamingScore     float64 `json:"roaming_score"`      // Map presence
	TeamfightScore   float64 `json:"teamfight_score"`    // Teamfight contribution

	// Advanced Analytics
	EconomicEfficiency    float64 `json:"economic_efficiency"`    // Gold efficiency
	VisionEfficiency      float64 `json:"vision_efficiency"`      // Vision score per minute
	ObjectiveContribution float64 `json:"objective_contribution"` // Objective damage and participation
}

// TFTMatch represents a Teamfight Tactics match
type TFTMatch struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Riot API Match Data
	MatchID     string `json:"match_id" gorm:"uniqueIndex;not null"`
	DataVersion string `json:"data_version"`

	// Game Information
	GameDateTime int64  `json:"game_datetime"`
	GameLength   int    `json:"game_length"`
	GameVersion  string `json:"game_version"`
	QueueID      int    `json:"queue_id"`
	TFTSetNumber int    `json:"tft_set_number"`
	TFTGameType  string `json:"tft_game_type"`

	// Participants
	TFTParticipants []TFTParticipant `json:"participants" gorm:"foreignKey:TFTMatchID"`

	// Processing Status
	IsProcessed bool      `json:"is_processed" gorm:"default:false"`
	ProcessedAt time.Time `json:"processed_at"`
	IsAnalyzed  bool      `json:"is_analyzed" gorm:"default:false"`
	AnalyzedAt  time.Time `json:"analyzed_at"`
}

// TFTParticipant represents a participant in a TFT match
type TFTParticipant struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	TFTMatchID uuid.UUID `json:"tft_match_id" gorm:"not null;index"`
	TFTMatch   TFTMatch  `json:"-" gorm:"foreignKey:TFTMatchID"`

	// Player Information
	PUUID string `json:"puuid" gorm:"not null;index"`

	// Match Results
	Placement            int `json:"placement"` // 1-8
	Level                int `json:"level"`
	LastRound            int `json:"last_round"`
	PlayersEliminated    int `json:"players_eliminated"`
	TotalDamageToPlayers int `json:"total_damage_to_players"`

	// Economy
	GoldLeft       int `json:"gold_left"`
	TimeEliminated int `json:"time_eliminated"`

	// Composition
	Units    []TFTUnit    `json:"units" gorm:"foreignKey:TFTParticipantID;constraint:OnDelete:CASCADE"`
	Traits   []TFTTrait   `json:"traits" gorm:"foreignKey:TFTParticipantID;constraint:OnDelete:CASCADE"`
	Augments []TFTAugment `json:"augments" gorm:"foreignKey:TFTParticipantID;constraint:OnDelete:CASCADE"`

	// Performance Analytics (calculated by Herald.lol)
	EconomyScore     float64 `json:"economy_score"`
	PositioningScore float64 `json:"positioning_score"`
	CompositionScore float64 `json:"composition_score"`
	PerformanceScore float64 `json:"performance_score"`
}

// TFTUnit represents a unit in a TFT composition
type TFTUnit struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TFTParticipantID uuid.UUID `json:"tft_participant_id" gorm:"not null"`

	CharacterID string `json:"character_id"`
	ItemNames   string `json:"item_names"` // JSON array of items
	Name        string `json:"name"`
	Rarity      int    `json:"rarity"`
	Tier        int    `json:"tier"`
}

// TFTTrait represents an active trait in a TFT composition
type TFTTrait struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TFTParticipantID uuid.UUID `json:"tft_participant_id" gorm:"not null"`

	Name        string `json:"name"`
	NumUnits    int    `json:"num_units"`
	Style       int    `json:"style"` // Trait activation level
	TierCurrent int    `json:"tier_current"`
	TierTotal   int    `json:"tier_total"`
}

// TFTAugment represents an augment in TFT
type TFTAugment struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TFTParticipantID uuid.UUID `json:"tft_participant_id" gorm:"not null"`

	Name string `json:"name"`
}

// BeforeCreate hooks
func (m *Match) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (mp *MatchParticipant) BeforeCreate(tx *gorm.DB) error {
	if mp.ID == uuid.Nil {
		mp.ID = uuid.New()
	}
	return nil
}

func (tm *TFTMatch) BeforeCreate(tx *gorm.DB) error {
	if tm.ID == uuid.Nil {
		tm.ID = uuid.New()
	}
	return nil
}

func (tp *TFTParticipant) BeforeCreate(tx *gorm.DB) error {
	if tp.ID == uuid.Nil {
		tp.ID = uuid.New()
	}
	return nil
}

// Helper methods
func (mp *MatchParticipant) CalculateKDA() float64 {
	if mp.Deaths == 0 {
		return float64(mp.Kills + mp.Assists)
	}
	return float64(mp.Kills+mp.Assists) / float64(mp.Deaths)
}

func (mp *MatchParticipant) CalculateCSPerMinute(gameDurationSeconds int) float64 {
	if gameDurationSeconds == 0 {
		return 0
	}
	minutes := float64(gameDurationSeconds) / 60.0
	return float64(mp.TotalCS) / minutes
}

func (mp *MatchParticipant) IsWin() bool {
	return mp.Won
}

func (tp *TFTParticipant) IsTop4() bool {
	return tp.Placement <= 4
}

func (tp *TFTParticipant) IsWin() bool {
	return tp.Placement == 1
}
