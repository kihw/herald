package models

import (
	"time"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// GroupUser représente un utilisateur étendu pour le système de groupes
type GroupUser struct {
	ID          int       `json:"id" db:"id"`
	GoogleID    string    `json:"google_id,omitempty" db:"google_id"`
	Email       string    `json:"email" db:"email"`
	Name        string    `json:"name" db:"name"`
	Picture     string    `json:"picture,omitempty" db:"picture"`
	RiotID      string    `json:"riot_id,omitempty" db:"riot_id"`
	RiotTag     string    `json:"riot_tag,omitempty" db:"riot_tag"`
	Region      string    `json:"region" db:"region"`
	Rank        string    `json:"rank,omitempty" db:"rank"`
	LP          *int      `json:"lp,omitempty" db:"lp"`
	MMR         *int      `json:"mmr,omitempty" db:"mmr"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	LastActive  time.Time `json:"last_active" db:"last_active"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	Preferences UserPreferences `json:"preferences" db:"preferences"`
}

// UserPreferences contient les préférences utilisateur
type UserPreferences struct {
	Theme               string   `json:"theme"`                // dark, light
	NotificationsEmail  bool     `json:"notifications_email"`
	NotificationsBrowser bool    `json:"notifications_browser"`
	PrivacyLevel        string   `json:"privacy_level"`        // public, friends, private
	FavoriteChampions   []int    `json:"favorite_champions"`
	PreferredRoles      []string `json:"preferred_roles"`
	StatsVisibility     string   `json:"stats_visibility"`     // public, friends, private
}

// Value implémente l'interface driver.Valuer pour UserPreferences
func (up UserPreferences) Value() (driver.Value, error) {
	return json.Marshal(up)
}

// Scan implémente l'interface sql.Scanner pour UserPreferences
func (up *UserPreferences) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, up)
}

// Group représente un groupe d'amis
type Group struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	OwnerID     int       `json:"owner_id" db:"owner_id"`
	Owner       *User     `json:"owner,omitempty"`
	Privacy     string    `json:"privacy" db:"privacy"`           // public, private, invite_only
	InviteCode  string    `json:"invite_code,omitempty" db:"invite_code"`
	Settings    GroupSettings `json:"settings" db:"settings"`
	MemberCount int       `json:"member_count" db:"member_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Members     []GroupMember `json:"members,omitempty"`
	Stats       *GroupStats   `json:"stats,omitempty"`
}

// GroupSettings contient les paramètres du groupe
type GroupSettings struct {
	AllowInvitations    bool     `json:"allow_invitations"`
	AutoAcceptFriends   bool     `json:"auto_accept_friends"`
	ShowMemberStats     bool     `json:"show_member_stats"`
	AllowedRegions      []string `json:"allowed_regions"`
	MinRankRequirement  string   `json:"min_rank_requirement,omitempty"`
	ComparisonFeatures  []string `json:"comparison_features"`
	NotificationLevel   string   `json:"notification_level"`    // all, mentions, none
}

// Value implémente l'interface driver.Valuer pour GroupSettings
func (gs GroupSettings) Value() (driver.Value, error) {
	return json.Marshal(gs)
}

// Scan implémente l'interface sql.Scanner pour GroupSettings
func (gs *GroupSettings) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, gs)
}

// GroupMember représente un membre d'un groupe
type GroupMember struct {
	ID         int       `json:"id" db:"id"`
	GroupID    int       `json:"group_id" db:"group_id"`
	UserID     int       `json:"user_id" db:"user_id"`
	User       *User     `json:"user,omitempty"`
	Role       string    `json:"role" db:"role"`           // owner, admin, member
	Status     string    `json:"status" db:"status"`       // active, pending, banned
	JoinedAt   time.Time `json:"joined_at" db:"joined_at"`
	Nickname   string    `json:"nickname,omitempty" db:"nickname"`
	Permissions GroupMemberPermissions `json:"permissions" db:"permissions"`
}

// GroupMemberPermissions définit les permissions d'un membre
type GroupMemberPermissions struct {
	CanInvite      bool `json:"can_invite"`
	CanKick        bool `json:"can_kick"`
	CanEditGroup   bool `json:"can_edit_group"`
	CanViewStats   bool `json:"can_view_stats"`
	CanCreateComps bool `json:"can_create_comparisons"`
}

// Value implémente l'interface driver.Valuer pour GroupMemberPermissions
func (gmp GroupMemberPermissions) Value() (driver.Value, error) {
	return json.Marshal(gmp)
}

// Scan implémente l'interface sql.Scanner pour GroupMemberPermissions
func (gmp *GroupMemberPermissions) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, gmp)
}

// GroupInvite représente une invitation à un groupe
type GroupInvite struct {
	ID        int       `json:"id" db:"id"`
	GroupID   int       `json:"group_id" db:"group_id"`
	Group     *Group    `json:"group,omitempty"`
	InviterID int       `json:"inviter_id" db:"inviter_id"`
	Inviter   *User     `json:"inviter,omitempty"`
	InviteeID *int      `json:"invitee_id,omitempty" db:"invitee_id"`
	Invitee   *User     `json:"invitee,omitempty"`
	Email     string    `json:"email,omitempty" db:"email"`
	Status    string    `json:"status" db:"status"`         // pending, accepted, declined, expired
	Message   string    `json:"message,omitempty" db:"message"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// GroupStats contient les statistiques d'un groupe
type GroupStats struct {
	ID                int       `json:"id" db:"id"`
	GroupID           int       `json:"group_id" db:"group_id"`
	TotalMembers      int       `json:"total_members" db:"total_members"`
	ActiveMembers     int       `json:"active_members" db:"active_members"`
	AverageRank       string    `json:"average_rank,omitempty" db:"average_rank"`
	AverageMMR        *float64  `json:"average_mmr,omitempty" db:"average_mmr"`
	TopChampions      ChampionStatList `json:"top_champions" db:"top_champions"`
	PopularRoles      RoleStatList     `json:"popular_roles" db:"popular_roles"`
	WinRateComparison WinRateMap `json:"winrate_comparison" db:"winrate_comparison"`
	LastUpdated       time.Time `json:"last_updated" db:"last_updated"`
}

// ChampionStat représente une statistique de champion
type ChampionStat struct {
	ChampionID   int     `json:"champion_id"`
	ChampionName string  `json:"champion_name"`
	PlayCount    int     `json:"play_count"`
	WinRate      float64 `json:"win_rate"`
	AvgKDA       float64 `json:"avg_kda"`
}

// RoleStat représente une statistique de rôle
type RoleStat struct {
	Role      string  `json:"role"`
	PlayCount int     `json:"play_count"`
	WinRate   float64 `json:"win_rate"`
}

// ChampionStatList est un alias pour []ChampionStat avec des méthodes customisées
type ChampionStatList []ChampionStat

// Value implémente l'interface driver.Valuer pour ChampionStatList
func (cs ChampionStatList) Value() (driver.Value, error) {
	return json.Marshal(cs)
}

// Scan implémente l'interface sql.Scanner pour ChampionStatList
func (cs *ChampionStatList) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, cs)
}

// RoleStatList est un alias pour []RoleStat avec des méthodes customisées
type RoleStatList []RoleStat

// Value implémente l'interface driver.Valuer pour RoleStatList
func (rs RoleStatList) Value() (driver.Value, error) {
	return json.Marshal(rs)
}

// Scan implémente l'interface sql.Scanner pour RoleStatList
func (rs *RoleStatList) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, rs)
}

// WinRateMap est un alias pour map[string]float64 avec des méthodes customisées
type WinRateMap map[string]float64

// Value implémente l'interface driver.Valuer pour WinRateMap
func (wr WinRateMap) Value() (driver.Value, error) {
	return json.Marshal(wr)
}

// Scan implémente l'interface sql.Scanner pour WinRateMap
func (wr *WinRateMap) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, wr)
}

// GroupComparison représente une comparaison entre membres du groupe
type GroupComparison struct {
	ID          int       `json:"id" db:"id"`
	GroupID     int       `json:"group_id" db:"group_id"`
	Group       *Group    `json:"group,omitempty"`
	CreatorID   int       `json:"creator_id" db:"creator_id"`
	Creator     *User     `json:"creator,omitempty"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CompareType string    `json:"compare_type" db:"compare_type"`   // champions, roles, performance, trends
	Parameters  ComparisonParameters `json:"parameters" db:"parameters"`
	Results     ComparisonResults    `json:"results,omitempty" db:"results"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ComparisonParameters définit les paramètres de comparaison
type ComparisonParameters struct {
	MemberIDs    []int    `json:"member_ids"`
	TimeRange    string   `json:"time_range"`      // last_week, last_month, season, all_time
	Champions    []int    `json:"champions,omitempty"`
	Roles        []string `json:"roles,omitempty"`
	GameModes    []int    `json:"game_modes,omitempty"`
	Metrics      []string `json:"metrics"`         // winrate, kda, cs, damage, vision, etc.
	MinGames     int      `json:"min_games"`
}

// Value implémente l'interface driver.Valuer pour ComparisonParameters
func (cp ComparisonParameters) Value() (driver.Value, error) {
	return json.Marshal(cp)
}

// Scan implémente l'interface sql.Scanner pour ComparisonParameters
func (cp *ComparisonParameters) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, cp)
}

// ComparisonResults contient les résultats d'une comparaison
type ComparisonResults struct {
	Summary       ComparisonSummary        `json:"summary"`
	MemberStats   map[string]interface{}   `json:"member_stats"`
	Charts        []ChartData              `json:"charts"`
	Rankings      []MemberRanking          `json:"rankings"`
	Insights      []string                 `json:"insights"`
	GeneratedAt   time.Time                `json:"generated_at"`
}

// ComparisonSummary résume les résultats de comparaison
type ComparisonSummary struct {
	TopPerformer    string  `json:"top_performer"`
	BestMetric      string  `json:"best_metric"`
	AverageWinRate  float64 `json:"average_win_rate"`
	TotalGamesCompared int  `json:"total_games_compared"`
	TimeSpan        string  `json:"time_span"`
}

// ChartData représente des données pour un graphique
type ChartData struct {
	Type      string                 `json:"type"`        // bar, line, radar, pie
	Title     string                 `json:"title"`
	Labels    []string               `json:"labels"`
	Datasets  []ChartDataset         `json:"datasets"`
	Options   map[string]interface{} `json:"options,omitempty"`
}

// ChartDataset représente un dataset pour un graphique
type ChartDataset struct {
	Label           string    `json:"label"`
	Data            []float64 `json:"data"`
	BackgroundColor []string  `json:"background_color,omitempty"`
	BorderColor     []string  `json:"border_color,omitempty"`
}

// MemberRanking représente le classement d'un membre
type MemberRanking struct {
	UserID      int     `json:"user_id"`
	Username    string  `json:"username"`
	Rank        int     `json:"rank"`
	Score       float64 `json:"score"`
	Metric      string  `json:"metric"`
	Change      string  `json:"change"`    // up, down, same
}

// Value implémente l'interface driver.Valuer pour ComparisonResults
func (cr ComparisonResults) Value() (driver.Value, error) {
	return json.Marshal(cr)
}

// Scan implémente l'interface sql.Scanner pour ComparisonResults
func (cr *ComparisonResults) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, cr)
}

// GetDefaultUserPreferences retourne les préférences par défaut
func GetDefaultUserPreferences() UserPreferences {
	return UserPreferences{
		Theme:                "dark",
		NotificationsEmail:   true,
		NotificationsBrowser: false,
		PrivacyLevel:         "friends",
		FavoriteChampions:    []int{},
		PreferredRoles:       []string{},
		StatsVisibility:      "friends",
	}
}

// GetDefaultGroupSettings retourne les paramètres par défaut pour un groupe
func GetDefaultGroupSettings() GroupSettings {
	return GroupSettings{
		AllowInvitations:   true,
		AutoAcceptFriends:  false,
		ShowMemberStats:    true,
		AllowedRegions:     []string{"euw1", "na1", "eun1"},
		ComparisonFeatures: []string{"champions", "roles", "performance"},
		NotificationLevel:  "mentions",
	}
}

// GetMemberPermissions retourne les permissions par défaut selon le rôle
func GetMemberPermissions(role string) GroupMemberPermissions {
	switch role {
	case "owner":
		return GroupMemberPermissions{
			CanInvite:      true,
			CanKick:        true,
			CanEditGroup:   true,
			CanViewStats:   true,
			CanCreateComps: true,
		}
	case "admin":
		return GroupMemberPermissions{
			CanInvite:      true,
			CanKick:        true,
			CanEditGroup:   false,
			CanViewStats:   true,
			CanCreateComps: true,
		}
	default: // member
		return GroupMemberPermissions{
			CanInvite:      false,
			CanKick:        false,
			CanEditGroup:   false,
			CanViewStats:   true,
			CanCreateComps: true,
		}
	}
}