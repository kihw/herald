// Skill Progression Models for Herald.lol
package models

import (
	"gorm.io/gorm"
	"time"
)

// SkillProgressionAnalysis represents a skill progression analysis result
type SkillProgressionAnalysis struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	SummonerID    string    `gorm:"not null;index" json:"summonerId"`
	AnalysisType  string    `gorm:"not null;index" json:"analysisType"` // overall, champion, role, skill
	AnalysisData  string    `gorm:"type:text" json:"analysisData"`      // JSON stored as text
	OverallRating float64   `gorm:"not null" json:"overallRating"`
	Confidence    float64   `gorm:"not null" json:"confidence"`
	CreatedAt     time.Time `gorm:"not null" json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// SkillCategoryTracking tracks individual skill categories over time
type SkillCategoryTracking struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SummonerID  string    `gorm:"not null;index" json:"summonerId"`
	Category    string    `gorm:"not null;index" json:"category"` // mechanical, tactical, strategic, mental, champion_specific
	Rating      float64   `gorm:"not null" json:"rating"`         // 0-100
	Percentile  float64   `json:"percentile"`                     // compared to similar rank
	Improvement float64   `json:"improvement"`                    // change from previous measurement
	Confidence  float64   `json:"confidence"`
	MeasuredAt  time.Time `gorm:"not null;index" json:"measuredAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// SkillSubcategoryTracking tracks subcategories within each skill category
type SkillSubcategoryTracking struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CategoryID  uint      `gorm:"not null;index" json:"categoryId"`
	SummonerID  string    `gorm:"not null;index" json:"summonerId"`
	Subcategory string    `gorm:"not null" json:"subcategory"`
	Value       float64   `gorm:"not null" json:"value"`
	Target      float64   `json:"target"`
	Weight      float64   `json:"weight"` // importance to overall category
	Improvement float64   `json:"improvement"`
	Unit        string    `json:"unit"`
	Description string    `json:"description"`
	MeasuredAt  time.Time `gorm:"not null;index" json:"measuredAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Foreign keys
	Category SkillCategoryTracking `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// RankProgressionHistory tracks rank changes over time
type RankProgressionHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SummonerID  string    `gorm:"not null;index" json:"summonerId"`
	Rank        string    `gorm:"not null" json:"rank"`
	Division    string    `json:"division"` // I, II, III, IV
	LP          int       `json:"lp"`
	MMR         int       `json:"mmr"`         // estimated
	Change      int       `json:"change"`      // LP change
	MatchResult string    `json:"matchResult"` // win, loss
	Performance float64   `json:"performance"` // 0-100 performance rating
	Season      string    `gorm:"not null;index" json:"season"`
	GameMode    string    `gorm:"not null;index" json:"gameMode"`
	RecordedAt  time.Time `gorm:"not null;index" json:"recordedAt"`
	CreatedAt   time.Time `json:"createdAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// ChampionMasteryProgression tracks champion-specific skill development
type ChampionMasteryProgression struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	SummonerID        string    `gorm:"not null;index" json:"summonerId"`
	Champion          string    `gorm:"not null;index" json:"champion"`
	Role              string    `gorm:"not null;index" json:"role"`
	MasteryLevel      int       `json:"masteryLevel"`
	MasteryPoints     int       `json:"masteryPoints"`
	GamesPlayed       int       `json:"gamesPlayed"`
	WinRate           float64   `json:"winRate"`
	PerformanceRating float64   `json:"performanceRating"` // 0-100
	ProgressionStage  string    `json:"progressionStage"`  // learning, improving, mastering, expert
	MechanicsRating   float64   `json:"mechanicsRating"`
	KnowledgeRating   float64   `json:"knowledgeRating"`
	DecisionRating    float64   `json:"decisionRating"`
	ImprovementRate   float64   `json:"improvementRate"` // per week
	LastPlayed        time.Time `json:"lastPlayed"`
	MeasuredAt        time.Time `gorm:"not null;index" json:"measuredAt"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// CoreSkillMeasurement tracks fundamental gaming skills
type CoreSkillMeasurement struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	SummonerID    string    `gorm:"not null;index" json:"summonerId"`
	SkillType     string    `gorm:"not null;index" json:"skillType"` // mechanical, game_knowledge, strategic, mental, communication, adaptability, leadership
	Rating        float64   `gorm:"not null" json:"rating"`          // 0-100
	Percentile    float64   `json:"percentile"`
	Components    string    `gorm:"type:text" json:"components"` // JSON array of skill components
	Trend         string    `json:"trend"`                       // improving, stable, declining
	TrendStrength float64   `json:"trendStrength"`
	NextLevel     string    `json:"nextLevel"`
	Blockers      string    `gorm:"type:text" json:"blockers"`  // JSON array
	Catalysts     string    `gorm:"type:text" json:"catalysts"` // JSON array
	MeasuredAt    time.Time `gorm:"not null;index" json:"measuredAt"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// LearningCurveData tracks learning patterns and efficiency
type LearningCurveData struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	SummonerID        string    `gorm:"not null;index" json:"summonerId"`
	CurveType         string    `json:"curveType"`     // linear, exponential, plateau, sigmoid
	LearningPhase     string    `json:"learningPhase"` // beginner, intermediate, advanced, expert
	LearningRate      float64   `json:"learningRate"`  // skill points per hour
	RetentionRate     float64   `json:"retentionRate"`
	TransferRate      float64   `json:"transferRate"`      // skill transfer between champions
	FocusQuality      float64   `json:"focusQuality"`      // quality of practice
	ImprovementFactor float64   `json:"improvementFactor"` // learning multiplier
	InPlateau         bool      `json:"inPlateau"`
	PlateauDuration   int       `json:"plateauDuration"` // weeks
	PlateauLevel      float64   `json:"plateauLevel"`
	BreakoutChance    float64   `json:"breakoutChance"`
	OptimalHours      float64   `json:"optimalHours"`   // hours per week
	OptimalSession    int       `json:"optimalSession"` // minutes per session
	MeasuredAt        time.Time `gorm:"not null;index" json:"measuredAt"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// SkillMilestone represents skill achievement milestones
type SkillMilestone struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	SummonerID      string    `gorm:"not null;index" json:"summonerId"`
	MilestoneID     string    `gorm:"not null;unique;index" json:"milestoneId"`
	Name            string    `gorm:"not null" json:"name"`
	Category        string    `gorm:"not null;index" json:"category"`
	Description     string    `gorm:"type:text" json:"description"`
	Achieved        bool      `gorm:"default:false" json:"achieved"`
	Progress        float64   `json:"progress"`                      // 0-100
	Requirements    string    `gorm:"type:text" json:"requirements"` // JSON array
	Reward          string    `gorm:"type:text" json:"reward"`       // JSON object
	Difficulty      string    `json:"difficulty"`                    // easy, moderate, hard, extreme
	EstimatedTime   string    `json:"estimatedTime"`
	Tips            string    `gorm:"type:text" json:"tips"` // JSON array
	AchievementDate time.Time `json:"achievementDate,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// ProgressionPrediction stores skill progression predictions
type ProgressionPrediction struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	SummonerID     string    `gorm:"not null;index" json:"summonerId"`
	PredictionType string    `gorm:"not null;index" json:"predictionType"` // rank, skill, champion_mastery
	PredictedValue string    `json:"predictedValue"`                       // predicted rank or rating
	Confidence     float64   `json:"confidence"`                           // 0-100
	TimeFrame      string    `json:"timeFrame"`                            // 1 week, 1 month, 1 season
	KeyFactors     string    `gorm:"type:text" json:"keyFactors"`          // JSON array
	Requirements   string    `gorm:"type:text" json:"requirements"`        // JSON array
	Likelihood     float64   `json:"likelihood"`                           // 0-100
	Assumptions    string    `gorm:"type:text" json:"assumptions"`         // JSON array
	ActualResult   string    `json:"actualResult,omitempty"`               // for accuracy tracking
	Accuracy       float64   `json:"accuracy,omitempty"`                   // how accurate was prediction
	PredictedAt    time.Time `gorm:"not null;index" json:"predictedAt"`
	ValidUntil     time.Time `json:"validUntil"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// PotentialAssessment stores long-term potential analysis
type PotentialAssessment struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	SummonerID        string    `gorm:"not null;index" json:"summonerId"`
	OverallPotential  float64   `json:"overallPotential"` // 0-100
	PeakPrediction    string    `json:"peakPrediction"`   // predicted peak rank
	NaturalTalent     float64   `json:"naturalTalent"`
	WorkEthic         float64   `json:"workEthic"`
	LearningSpeed     float64   `json:"learningSpeed"`
	Consistency       float64   `json:"consistency"`
	Adaptability      float64   `json:"adaptability"`
	CompetitiveDrive  float64   `json:"competitiveDrive"`
	TalentProfile     string    `json:"talentProfile"`                      // prodigy, grinder, balanced, late_bloomer
	LimitingFactors   string    `gorm:"type:text" json:"limitingFactors"`   // JSON array
	StrengthAreas     string    `gorm:"type:text" json:"strengthAreas"`     // JSON array
	UntappedPotential string    `gorm:"type:text" json:"untappedPotential"` // JSON array
	Recommendations   string    `gorm:"type:text" json:"recommendations"`   // JSON array
	Confidence        float64   `json:"confidence"`
	AssessedAt        time.Time `gorm:"not null;index" json:"assessedAt"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// ProgressionRecommendation stores personalized improvement recommendations
type ProgressionRecommendation struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	SummonerID       string    `gorm:"not null;index" json:"summonerId"`
	RecommendationID string    `gorm:"not null;unique;index" json:"recommendationId"`
	Type             string    `gorm:"not null;index" json:"type"`     // practice, champion, playstyle, mindset
	Priority         string    `gorm:"not null;index" json:"priority"` // critical, high, medium, low
	Title            string    `gorm:"not null" json:"title"`
	Description      string    `gorm:"type:text" json:"description"`
	ImpactRating     float64   `json:"impactRating"` // 0-100
	Difficulty       string    `json:"difficulty"`   // easy, moderate, hard
	TimeCommitment   string    `json:"timeCommitment"`
	ActionSteps      string    `gorm:"type:text" json:"actionSteps"` // JSON array
	Success          string    `gorm:"type:text" json:"success"`     // JSON object with success metrics
	Related          string    `gorm:"type:text" json:"related"`     // JSON array of related recommendation IDs
	Status           string    `gorm:"default:active" json:"status"` // active, completed, dismissed, paused
	Progress         float64   `json:"progress"`                     // 0-100
	StartedAt        time.Time `json:"startedAt,omitempty"`
	CompletedAt      time.Time `json:"completedAt,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// SkillBreakthrough records significant skill improvements
type SkillBreakthrough struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SummonerID   string    `gorm:"not null;index" json:"summonerId"`
	SkillArea    string    `gorm:"not null;index" json:"skillArea"`
	ImpactSize   float64   `json:"impactSize"` // magnitude of improvement
	Trigger      string    `json:"trigger"`    // what caused the breakthrough
	Description  string    `gorm:"type:text" json:"description"`
	Lessons      string    `gorm:"type:text" json:"lessons"` // JSON array
	BeforeRating float64   `json:"beforeRating"`
	AfterRating  float64   `json:"afterRating"`
	Duration     int       `json:"duration"`  // days it took
	Validated    bool      `json:"validated"` // confirmed by subsequent performance
	OccurredAt   time.Time `gorm:"not null;index" json:"occurredAt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// SkillPracticeSession tracks dedicated practice activities
type SkillPracticeSession struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	SummonerID      string    `gorm:"not null;index" json:"summonerId"`
	SessionType     string    `gorm:"not null;index" json:"sessionType"` // cs_drill, mechanics, vod_review, theory
	FocusAreas      string    `gorm:"type:text" json:"focusAreas"`       // JSON array
	Duration        int       `json:"duration"`                          // minutes
	Quality         float64   `json:"quality"`                           // 1-10 subjective rating
	Goals           string    `gorm:"type:text" json:"goals"`            // JSON array
	Achievements    string    `gorm:"type:text" json:"achievements"`     // JSON array
	Notes           string    `gorm:"type:text" json:"notes"`
	ImprovementSeen bool      `json:"improvementSeen"`
	FollowUpNeeded  bool      `json:"followUpNeeded"`
	StartedAt       time.Time `gorm:"not null;index" json:"startedAt"`
	CompletedAt     time.Time `json:"completedAt"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// SkillGoal represents player-set skill improvement goals
type SkillGoal struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	SummonerID      string    `gorm:"not null;index" json:"summonerId"`
	GoalType        string    `gorm:"not null;index" json:"goalType"` // rank, skill_rating, champion_mastery, custom
	Target          string    `gorm:"not null" json:"target"`         // target value (rank, rating, etc.)
	Current         string    `json:"current"`                        // current value
	Priority        string    `json:"priority"`                       // high, medium, low
	Deadline        time.Time `json:"deadline"`
	Description     string    `gorm:"type:text" json:"description"`
	Strategy        string    `gorm:"type:text" json:"strategy"`    // JSON array of strategies
	Milestones      string    `gorm:"type:text" json:"milestones"`  // JSON array
	Progress        float64   `json:"progress"`                     // 0-100
	Status          string    `gorm:"default:active" json:"status"` // active, completed, paused, abandoned
	Achieved        bool      `json:"achieved"`
	AchievementDate time.Time `json:"achievementDate,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// SkillBenchmark stores reference performance standards
type SkillBenchmark struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	SkillArea     string    `gorm:"not null;index" json:"skillArea"`
	Rank          string    `gorm:"not null;index" json:"rank"`
	Role          string    `gorm:"index" json:"role"`
	MetricName    string    `gorm:"not null" json:"metricName"`
	ExpectedValue float64   `json:"expectedValue"`
	MinValue      float64   `json:"minValue"`
	MaxValue      float64   `json:"maxValue"`
	Unit          string    `json:"unit"`
	Description   string    `json:"description"`
	SampleSize    int       `json:"sampleSize"`
	Confidence    float64   `json:"confidence"`
	Patch         string    `gorm:"index" json:"patch"`
	Region        string    `gorm:"index" json:"region"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedAt     time.Time `json:"createdAt"`
}

// GORM Hooks
func (s *SkillProgressionAnalysis) BeforeCreate(tx *gorm.DB) error {
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	return nil
}

func (s *SkillCategoryTracking) BeforeCreate(tx *gorm.DB) error {
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	if s.MeasuredAt.IsZero() {
		s.MeasuredAt = time.Now()
	}
	return nil
}

func (p *SkillPracticeSession) BeforeCreate(tx *gorm.DB) error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	if p.StartedAt.IsZero() {
		p.StartedAt = time.Now()
	}
	return nil
}

// Table names
func (SkillProgressionAnalysis) TableName() string {
	return "skill_progression_analyses"
}

func (SkillCategoryTracking) TableName() string {
	return "skill_category_tracking"
}

func (SkillSubcategoryTracking) TableName() string {
	return "skill_subcategory_tracking"
}

func (RankProgressionHistory) TableName() string {
	return "rank_progression_history"
}

func (ChampionMasteryProgression) TableName() string {
	return "champion_mastery_progression"
}

func (CoreSkillMeasurement) TableName() string {
	return "core_skill_measurements"
}

func (LearningCurveData) TableName() string {
	return "learning_curve_data"
}

func (SkillMilestone) TableName() string {
	return "skill_milestones"
}

func (ProgressionPrediction) TableName() string {
	return "progression_predictions"
}

func (PotentialAssessment) TableName() string {
	return "potential_assessments"
}

func (ProgressionRecommendation) TableName() string {
	return "progression_recommendations"
}

func (SkillBreakthrough) TableName() string {
	return "skill_breakthroughs"
}

func (SkillPracticeSession) TableName() string {
	return "skill_practice_sessions"
}

func (SkillGoal) TableName() string {
	return "skill_goals"
}

func (SkillBenchmark) TableName() string {
	return "skill_benchmarks"
}
