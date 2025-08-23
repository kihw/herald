// Counter Pick Models for Herald.lol
package models

import (
	"time"
	"gorm.io/gorm"
)

// CounterPickAnalysis represents a counter-pick analysis result
type CounterPickAnalysis struct {
	ID             string    `gorm:"primaryKey" json:"id"`
	TargetChampion string    `gorm:"not null;index" json:"targetChampion"`
	TargetRole     string    `gorm:"not null;index" json:"targetRole"`
	GameMode       string    `gorm:"not null;index" json:"gameMode"`
	AnalysisData   string    `gorm:"type:text" json:"analysisData"` // JSON stored as text
	Confidence     float64   `gorm:"not null" json:"confidence"`
	CreatedAt      time.Time `gorm:"not null" json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// CounterPickSuggestion represents individual counter pick suggestions
type CounterPickSuggestion struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	AnalysisID          string    `gorm:"not null;index" json:"analysisId"`
	Champion            string    `gorm:"not null;index" json:"champion"`
	TargetChampion      string    `gorm:"not null;index" json:"targetChampion"`
	CounterStrength     float64   `gorm:"not null" json:"counterStrength"`
	WinRateAdvantage    float64   `json:"winRateAdvantage"`
	LaneAdvantage       float64   `json:"laneAdvantage"`
	TeamFightAdvantage  float64   `json:"teamFightAdvantage"`
	ScalingAdvantage    float64   `json:"scalingAdvantage"`
	MetaFit             float64   `json:"metaFit"`
	PlayerComfort       float64   `json:"playerComfort"`
	BanPriority         float64   `json:"banPriority"`
	Flexibility         float64   `json:"flexibility"`
	SafetyRating        float64   `json:"safetyRating"`
	MatchupDifficulty   string    `json:"matchupDifficulty"`
	CounterReasons      string    `gorm:"type:text" json:"counterReasons"` // JSON array as text
	PlayingTips         string    `gorm:"type:text" json:"playingTips"`    // JSON array as text
	ItemRecommendations string    `gorm:"type:text" json:"itemRecommendations"` // JSON array as text
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`

	// Foreign key
	Analysis CounterPickAnalysis `gorm:"foreignKey:AnalysisID;references:ID" json:"analysis,omitempty"`
}

// LaneCounterData represents lane-specific counter information
type LaneCounterData struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	AnalysisID        string    `gorm:"not null;index" json:"analysisId"`
	Phase             string    `gorm:"not null" json:"phase"` // early, mid, late
	Advantage         float64   `json:"advantage"`             // -100 to 100
	AllInPotential    float64   `json:"allInPotential"`
	RoamingPotential  float64   `json:"roamingPotential"`
	ScalingComparison string    `json:"scalingComparison"`
	KeyFactors        string    `gorm:"type:text" json:"keyFactors"`       // JSON array as text
	PlayingTips       string    `gorm:"type:text" json:"playingTips"`      // JSON array as text
	WardingTips       string    `gorm:"type:text" json:"wardingTips"`      // JSON array as text
	TradingPatterns   string    `gorm:"type:text" json:"tradingPatterns"`  // JSON array as text
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`

	// Foreign key
	Analysis CounterPickAnalysis `gorm:"foreignKey:AnalysisID;references:ID" json:"analysis,omitempty"`
}

// TeamFightCounterData represents team fight counter information
type TeamFightCounterData struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	AnalysisID       string    `gorm:"not null;index" json:"analysisId"`
	CounterType      string    `gorm:"not null" json:"counterType"` // engage, disengage, peel, burst, etc.
	Effectiveness    float64   `json:"effectiveness"`               // 0-100
	Positioning      string    `gorm:"type:text" json:"positioning"`      // JSON array as text
	ComboCounters    string    `gorm:"type:text" json:"comboCounters"`    // JSON array as text
	TeamCoordination string    `gorm:"type:text" json:"teamCoordination"` // JSON array as text
	ObjectiveControl string    `gorm:"type:text" json:"objectiveControl"` // JSON array as text
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	// Foreign key
	Analysis CounterPickAnalysis `gorm:"foreignKey:AnalysisID;references:ID" json:"analysis,omitempty"`
}

// ItemCounterData represents item-based counter strategies
type ItemCounterData struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	AnalysisID        string    `gorm:"not null;index" json:"analysisId"`
	ItemName          string    `gorm:"not null" json:"itemName"`
	CounterType       string    `gorm:"not null" json:"counterType"` // defensive, offensive, utility
	Effectiveness     float64   `json:"effectiveness"`               // 0-100
	BuildPriority     int       `json:"buildPriority"`               // 1-6
	SituationalUse    string    `json:"situationalUse"`
	CostEffectiveness float64   `json:"costEffectiveness"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`

	// Foreign key
	Analysis CounterPickAnalysis `gorm:"foreignKey:AnalysisID;references:ID" json:"analysis,omitempty"`
}

// PlayStyleCounterData represents playstyle-specific counter strategies
type PlayStyleCounterData struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	AnalysisID      string    `gorm:"not null;index" json:"analysisId"`
	TargetPlayStyle string    `gorm:"not null" json:"targetPlayStyle"`
	CounterStrategy string    `gorm:"not null" json:"counterStrategy"`
	RiskLevel       string    `json:"riskLevel"` // low, medium, high
	KeyPrinciples   string    `gorm:"type:text" json:"keyPrinciples"` // JSON array as text
	Timing          string    `gorm:"type:text" json:"timing"`        // JSON array as text
	TeamSupport     string    `gorm:"type:text" json:"teamSupport"`   // JSON array as text
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Foreign key
	Analysis CounterPickAnalysis `gorm:"foreignKey:AnalysisID;references:ID" json:"analysis,omitempty"`
}

// MultiTargetCounterAnalysis represents analysis for countering multiple champions
type MultiTargetCounterAnalysis struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	GameMode     string    `gorm:"not null;index" json:"gameMode"`
	AnalysisData string    `gorm:"type:text" json:"analysisData"` // JSON stored as text
	Confidence   float64   `gorm:"not null" json:"confidence"`
	CreatedAt    time.Time `gorm:"not null" json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// UniversalCounterSuggestion represents champions that counter multiple targets
type UniversalCounterSuggestion struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	AnalysisID      string    `gorm:"not null;index" json:"analysisId"`
	Champion        string    `gorm:"not null;index" json:"champion"`
	AverageStrength float64   `json:"averageStrength"`
	Versatility     float64   `json:"versatility"`
	CountersTargets string    `gorm:"type:text" json:"countersTargets"` // JSON array as text
	RecommendReasons string   `gorm:"type:text" json:"recommendReasons"` // JSON array as text
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Foreign key
	Analysis MultiTargetCounterAnalysis `gorm:"foreignKey:AnalysisID;references:ID" json:"analysis,omitempty"`
}

// SpecificCounterSuggestion represents specialized counters for specific champions
type SpecificCounterSuggestion struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	AnalysisID       string    `gorm:"not null;index" json:"analysisId"`
	Champion         string    `gorm:"not null;index" json:"champion"`
	PrimaryTarget    string    `gorm:"not null" json:"primaryTarget"`
	SecondaryTargets string    `gorm:"type:text" json:"secondaryTargets"` // JSON array as text
	CounterStrength  float64   `json:"counterStrength"`
	Specialization   string    `json:"specialization"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	// Foreign key
	Analysis MultiTargetCounterAnalysis `gorm:"foreignKey:AnalysisID;references:ID" json:"analysis,omitempty"`
}

// TeamCounterStrategy represents team-based counter strategies
type TeamCounterStrategy struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	AnalysisID        string    `gorm:"not null;index" json:"analysisId"`
	Strategy          string    `gorm:"not null" json:"strategy"`
	Effectiveness     float64   `json:"effectiveness"`
	Complexity        string    `json:"complexity"` // simple, moderate, complex
	Description       string    `gorm:"type:text" json:"description"`
	RequiredChampions string    `gorm:"type:text" json:"requiredChampions"` // JSON array as text
	Execution         string    `gorm:"type:text" json:"execution"`         // JSON array as text
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`

	// Foreign key
	Analysis MultiTargetCounterAnalysis `gorm:"foreignKey:AnalysisID;references:ID" json:"analysis,omitempty"`
}

// BanRecommendation represents ban strategy recommendations
type BanRecommendation struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	AnalysisID   string    `gorm:"not null;index" json:"analysisId"`
	Champion     string    `gorm:"not null;index" json:"champion"`
	Priority     float64   `json:"priority"` // 0-100
	Reasoning    string    `gorm:"type:text" json:"reasoning"`
	Impact       string    `gorm:"type:text" json:"impact"`
	Alternatives string    `gorm:"type:text" json:"alternatives"` // JSON array as text
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`

	// Foreign key
	Analysis MultiTargetCounterAnalysis `gorm:"foreignKey:AnalysisID;references:ID" json:"analysis,omitempty"`
}

// CounterPickMetrics represents performance metrics for counter picks
type CounterPickMetrics struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Champion        string    `gorm:"not null;index" json:"champion"`
	TargetChampion  string    `gorm:"not null;index" json:"targetChampion"`
	Role            string    `gorm:"not null;index" json:"role"`
	GameMode        string    `gorm:"not null;index" json:"gameMode"`
	Patch           string    `gorm:"not null;index" json:"patch"`
	SampleSize      int       `json:"sampleSize"`
	WinRate         float64   `json:"winRate"`
	LaneWinRate     float64   `json:"laneWinRate"`
	KDA             float64   `json:"kda"`
	DamageShare     float64   `json:"damageShare"`
	GoldDifferential float64  `json:"goldDifferential"`
	CSAdvantage     float64   `json:"csAdvantage"`
	VisionScore     float64   `json:"visionScore"`
	CounterStrength float64   `json:"counterStrength"` // Calculated overall strength
	Confidence      float64   `json:"confidence"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// CounterPickHistory represents historical counter pick performance
type CounterPickHistory struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null;index" json:"userId"`
	SummonerID     string    `gorm:"not null;index" json:"summonerId"`
	Champion       string    `gorm:"not null;index" json:"champion"`
	TargetChampion string    `gorm:"not null;index" json:"targetChampion"`
	Role           string    `gorm:"not null" json:"role"`
	GameMode       string    `gorm:"not null" json:"gameMode"`
	MatchID        string    `gorm:"not null;unique;index" json:"matchId"`
	Result         string    `gorm:"not null" json:"result"` // win, loss
	Performance    float64   `json:"performance"`            // 0-100 performance score
	PredictedStrength float64 `json:"predictedStrength"`    // What we predicted
	ActualStrength    float64 `json:"actualStrength"`       // What actually happened
	Accuracy          float64 `json:"accuracy"`             // How accurate our prediction was
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`

	// Foreign keys
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// CounterPickFavorites represents user's favorite counter picks
type CounterPickFavorites struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null;index" json:"userId"`
	SummonerID     string    `gorm:"not null;index" json:"summonerId"`
	Champion       string    `gorm:"not null;index" json:"champion"`
	TargetChampion string    `gorm:"not null;index" json:"targetChampion"`
	Role           string    `gorm:"not null" json:"role"`
	Notes          string    `gorm:"type:text" json:"notes"`
	PersonalRating float64   `json:"personalRating"` // User's personal rating 1-10
	TimesUsed      int       `json:"timesUsed"`
	SuccessRate    float64   `json:"successRate"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	// Foreign keys
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// GORM Hooks
func (c *CounterPickAnalysis) BeforeCreate(tx *gorm.DB) error {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	return nil
}

func (c *CounterPickSuggestion) BeforeCreate(tx *gorm.DB) error {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	return nil
}

func (m *MultiTargetCounterAnalysis) BeforeCreate(tx *gorm.DB) error {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return nil
}

// Table names
func (CounterPickAnalysis) TableName() string {
	return "counter_pick_analyses"
}

func (CounterPickSuggestion) TableName() string {
	return "counter_pick_suggestions"
}

func (LaneCounterData) TableName() string {
	return "lane_counter_data"
}

func (TeamFightCounterData) TableName() string {
	return "team_fight_counter_data"
}

func (ItemCounterData) TableName() string {
	return "item_counter_data"
}

func (PlayStyleCounterData) TableName() string {
	return "play_style_counter_data"
}

func (MultiTargetCounterAnalysis) TableName() string {
	return "multi_target_counter_analyses"
}

func (UniversalCounterSuggestion) TableName() string {
	return "universal_counter_suggestions"
}

func (SpecificCounterSuggestion) TableName() string {
	return "specific_counter_suggestions"
}

func (TeamCounterStrategy) TableName() string {
	return "team_counter_strategies"
}

func (BanRecommendation) TableName() string {
	return "ban_recommendations"
}

func (CounterPickMetrics) TableName() string {
	return "counter_pick_metrics"
}

func (CounterPickHistory) TableName() string {
	return "counter_pick_history"
}

func (CounterPickFavorites) TableName() string {
	return "counter_pick_favorites"
}