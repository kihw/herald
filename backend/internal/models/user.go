package models

import (
	"gorm.io/gorm"
	"time"
)

// User represents a Herald.lol user
type User struct {
	ID           string     `gorm:"primaryKey" json:"id"`
	Email        string     `gorm:"uniqueIndex;not null" json:"email"`
	Username     string     `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string     `gorm:"not null" json:"-"`
	DisplayName  string     `json:"displayName"`
	RiotID       string     `json:"riotId"`
	Region       string     `json:"region"`
	SummonerID   string     `gorm:"index" json:"summonerId"`
	PUUID        string     `gorm:"index" json:"puuid"`
	IsVerified   bool       `gorm:"default:false" json:"isVerified"`
	IsActive     bool       `gorm:"default:true" json:"isActive"`
	LastLoginAt  *time.Time `json:"lastLoginAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`

	// Gaming-specific fields
	PreferredRole string `json:"preferredRole"`
	MainChampions string `gorm:"type:text" json:"mainChampions"` // JSON array
	CurrentRank   string `json:"currentRank"`
	PeakRank      string `json:"peakRank"`
}

// UserPreferences represents user gaming preferences
type UserPreferences struct {
	ID                      uint      `gorm:"primaryKey" json:"id"`
	UserID                  string    `gorm:"not null;index" json:"userId"`
	Theme                   string    `gorm:"default:dark" json:"theme"`
	Language                string    `gorm:"default:en" json:"language"`
	Region                  string    `json:"region"`
	Timezone                string    `json:"timezone"`
	EmailNotifications      bool      `gorm:"default:true" json:"emailNotifications"`
	PushNotifications       bool      `gorm:"default:true" json:"pushNotifications"`
	MatchAlerts             bool      `gorm:"default:true" json:"matchAlerts"`
	AnalyticsSharing        bool      `gorm:"default:false" json:"analyticsSharing"`
	PrivacyMode             bool      `gorm:"default:false" json:"privacyMode"`
	AutoSyncMatches         bool      `gorm:"default:true" json:"autoSyncMatches"`
	CoachingRecommendations bool      `gorm:"default:true" json:"coachingRecommendations"`
	CreatedAt               time.Time `json:"createdAt"`
	UpdatedAt               time.Time `json:"updatedAt"`
}

// BeforeCreate hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.DisplayName == "" {
		u.DisplayName = u.Username
	}
	return nil
}

// BeforeUpdate hook
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// GetDisplayName returns the display name or username
func (u *User) GetDisplayName() string {
	if u.DisplayName != "" {
		return u.DisplayName
	}
	return u.Username
}

// GetMainChampions returns the main champions as a slice
func (u *User) GetMainChampions() []string {
	// TODO: Parse JSON string to slice
	return []string{}
}

// Table names
func (User) TableName() string {
	return "users"
}

func (UserPreferences) TableName() string {
	return "user_preferences"
}
