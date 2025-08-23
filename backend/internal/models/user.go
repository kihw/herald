package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a Herald.lol user
type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Authentication
	Email        string `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"-" gorm:"not null"`
	EmailVerified bool  `json:"email_verified" gorm:"default:false"`
	
	// Profile Information
	Username     string    `json:"username" gorm:"uniqueIndex;not null"`
	DisplayName  string    `json:"display_name"`
	Avatar       string    `json:"avatar"`
	Bio          string    `json:"bio"`
	Timezone     string    `json:"timezone" gorm:"default:'UTC'"`
	Language     string    `json:"language" gorm:"default:'en'"`
	
	// Account Status
	IsActive   bool      `json:"is_active" gorm:"default:true"`
	IsPremium  bool      `json:"is_premium" gorm:"default:false"`
	LastLogin  time.Time `json:"last_login"`
	LoginCount int       `json:"login_count" gorm:"default:0"`
	
	// Gaming Profiles
	RiotAccounts []RiotAccount `json:"riot_accounts" gorm:"foreignKey:UserID"`
	
	// Preferences
	Preferences UserPreferences `json:"preferences" gorm:"foreignKey:UserID"`
	
	// Subscriptions and Settings
	Subscription *Subscription `json:"subscription,omitempty" gorm:"foreignKey:UserID"`
	
	// Analytics
	TotalMatches     int       `json:"total_matches" gorm:"default:0"`
	LastSyncAt       time.Time `json:"last_sync_at"`
	FavoriteChampion string    `json:"favorite_champion"`
	MainRole         string    `json:"main_role"`
	CurrentRank      string    `json:"current_rank"`
}

// RiotAccount represents a linked Riot Games account
type RiotAccount struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	UserID uuid.UUID `json:"user_id" gorm:"not null;index"`
	User   User      `json:"-" gorm:"foreignKey:UserID"`
	
	// Riot Account Info
	PUUID        string `json:"puuid" gorm:"uniqueIndex;not null"`
	SummonerName string `json:"summoner_name" gorm:"not null"`
	TagLine      string `json:"tag_line" gorm:"not null"`
	SummonerID   string `json:"summoner_id"`
	AccountID    string `json:"account_id"`
	
	// Server and Region
	Region   string `json:"region" gorm:"not null"` // na1, euw1, etc.
	Platform string `json:"platform" gorm:"not null"` // americas, europe, asia
	
	// Account Status
	IsVerified bool      `json:"is_verified" gorm:"default:false"`
	IsPrimary  bool      `json:"is_primary" gorm:"default:false"`
	LastSyncAt time.Time `json:"last_sync_at"`
	
	// Current Rankings
	SoloQueueRank    string `json:"solo_queue_rank"`
	FlexQueueRank    string `json:"flex_queue_rank"`
	TFTRank          string `json:"tft_rank"`
	ArenaRank        string `json:"arena_rank"`
	
	// Profile Info
	SummonerLevel int    `json:"summoner_level"`
	ProfileIcon   int    `json:"profile_icon"`
	
	// Statistics
	TotalMasteryScore int `json:"total_mastery_score"`
}

// UserPreferences represents user preferences and settings
type UserPreferences struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	UserID uuid.UUID `json:"user_id" gorm:"uniqueIndex;not null"`
	User   User      `json:"-" gorm:"foreignKey:UserID"`
	
	// UI Preferences
	Theme                string `json:"theme" gorm:"default:'dark'"` // dark, light, auto
	CompactMode          bool   `json:"compact_mode" gorm:"default:false"`
	ShowDetailedStats    bool   `json:"show_detailed_stats" gorm:"default:true"`
	DefaultTimeframe     string `json:"default_timeframe" gorm:"default:'7d'"` // 1d, 7d, 30d, season
	
	// Notification Settings
	EmailNotifications   bool `json:"email_notifications" gorm:"default:true"`
	PushNotifications    bool `json:"push_notifications" gorm:"default:true"`
	MatchNotifications   bool `json:"match_notifications" gorm:"default:true"`
	RankChangeNotifications bool `json:"rank_change_notifications" gorm:"default:true"`
	
	// Analytics Preferences  
	AutoSyncMatches      bool   `json:"auto_sync_matches" gorm:"default:true"`
	SyncInterval         int    `json:"sync_interval" gorm:"default:300"` // seconds
	IncludeNormalGames   bool   `json:"include_normal_games" gorm:"default:true"`
	IncludeARAMGames     bool   `json:"include_aram_games" gorm:"default:true"`
	FavoriteGameModes    string `json:"favorite_game_modes"` // JSON array
	
	// Privacy Settings
	PublicProfile        bool `json:"public_profile" gorm:"default:true"`
	ShowInLeaderboards   bool `json:"show_in_leaderboards" gorm:"default:true"`
	AllowDataExport      bool `json:"allow_data_export" gorm:"default:true"`
	
	// Coaching Preferences
	ReceiveAICoaching    bool   `json:"receive_ai_coaching" gorm:"default:true"`
	CoachingFocus        string `json:"coaching_focus"` // improvement_areas, champion_pool, macro, micro
	SkillLevel           string `json:"skill_level" gorm:"default:'intermediate'"` // beginner, intermediate, advanced, expert
	PreferredCoachingStyle string `json:"preferred_coaching_style" gorm:"default:'balanced'"` // gentle, direct, balanced
}

// Subscription represents a user's subscription plan
type Subscription struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	UserID uuid.UUID `json:"user_id" gorm:"uniqueIndex;not null"`
	User   User      `json:"-" gorm:"foreignKey:UserID"`
	
	// Subscription Details
	Plan         string    `json:"plan" gorm:"not null"` // free, premium, elite, enterprise
	Status       string    `json:"status" gorm:"not null"` // active, canceled, expired, trial
	StartedAt    time.Time `json:"started_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	TrialEndsAt  *time.Time `json:"trial_ends_at"`
	
	// Billing
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency" gorm:"default:'USD'"`
	Interval     string  `json:"interval" gorm:"default:'monthly'"` // monthly, yearly
	PaymentMethod string `json:"payment_method"`
	
	// Features
	MaxRiotAccounts    int  `json:"max_riot_accounts" gorm:"default:1"`
	UnlimitedAnalytics bool `json:"unlimited_analytics" gorm:"default:false"`
	AICoachingAccess   bool `json:"ai_coaching_access" gorm:"default:false"`
	AdvancedMetrics    bool `json:"advanced_metrics" gorm:"default:false"`
	DataExportAccess   bool `json:"data_export_access" gorm:"default:false"`
	PrioritySupport    bool `json:"priority_support" gorm:"default:false"`
}

// BeforeCreate hook to set UUID for new users
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for RiotAccount
func (ra *RiotAccount) BeforeCreate(tx *gorm.DB) error {
	if ra.ID == uuid.Nil {
		ra.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for UserPreferences
func (up *UserPreferences) BeforeCreate(tx *gorm.DB) error {
	if up.ID == uuid.Nil {
		up.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for Subscription
func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// GetFullDisplayName returns the full display name or username if display name is empty
func (u *User) GetFullDisplayName() string {
	if u.DisplayName != "" {
		return u.DisplayName
	}
	return u.Username
}

// GetPrimaryRiotAccount returns the primary Riot account
func (u *User) GetPrimaryRiotAccount() *RiotAccount {
	for _, account := range u.RiotAccounts {
		if account.IsPrimary {
			return &account
		}
	}
	if len(u.RiotAccounts) > 0 {
		return &u.RiotAccounts[0]
	}
	return nil
}

// HasValidSubscription checks if user has an active subscription
func (u *User) HasValidSubscription() bool {
	if u.Subscription == nil {
		return false
	}
	return u.Subscription.Status == "active" && u.Subscription.ExpiresAt.After(time.Now())
}

// GetSubscriptionPlan returns the current subscription plan
func (u *User) GetSubscriptionPlan() string {
	if u.Subscription == nil {
		return "free"
	}
	return u.Subscription.Plan
}

// CanAddRiotAccount checks if user can add more Riot accounts
func (u *User) CanAddRiotAccount() bool {
	if u.Subscription == nil {
		return len(u.RiotAccounts) < 1 // Free tier: 1 account
	}
	return len(u.RiotAccounts) < u.Subscription.MaxRiotAccounts
}