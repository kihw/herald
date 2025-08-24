// Coaching Models for Herald.lol
package models

import (
	"gorm.io/gorm"
	"time"
)

// CoachingInsight represents a coaching analysis and insights result
type CoachingInsight struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	SummonerID  string    `gorm:"not null;index" json:"summonerId"`
	InsightType string    `gorm:"not null;index" json:"insightType"` // match_analysis, skill_development, strategic, tactical, mental
	InsightData string    `gorm:"type:text" json:"insightData"`      // JSON stored as text
	Confidence  float64   `gorm:"not null" json:"confidence"`
	CreatedAt   time.Time `gorm:"not null" json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CoachingTip represents individual coaching tips
type CoachingTip struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TipID      string    `gorm:"not null;unique;index" json:"tipId"`
	SummonerID string    `gorm:"not null;index" json:"summonerId"`
	Category   string    `gorm:"not null;index" json:"category"` // mechanical, tactical, strategic, mental, champion_specific
	Type       string    `gorm:"not null" json:"type"`           // quick_tip, deep_insight, warning, opportunity
	Title      string    `gorm:"not null" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	Context    string    `gorm:"type:text" json:"context"` // JSON stored as text
	Relevance  float64   `json:"relevance"`                // 0-100 how relevant to player
	Actionable bool      `json:"actionable"`
	Difficulty string    `json:"difficulty"`                   // easy, moderate, hard
	Expected   string    `gorm:"type:text" json:"expected"`    // JSON stored as text (ExpectedOutcome)
	Related    string    `gorm:"type:text" json:"related"`     // JSON array of related tip IDs
	Status     string    `gorm:"default:active" json:"status"` // active, applied, dismissed, archived
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// CoachingTipFeedback represents user feedback on coaching tips
type CoachingTipFeedback struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TipID      uint      `gorm:"not null;index" json:"tipId"`
	SummonerID string    `gorm:"not null;index" json:"summonerId"`
	Helpful    bool      `json:"helpful"`
	Applied    bool      `json:"applied"`
	Effective  bool      `json:"effective"`
	Comments   string    `gorm:"type:text" json:"comments"`
	Rating     float64   `json:"rating"` // 1-10
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	// Foreign keys
	Tip  CoachingTip `gorm:"foreignKey:TipID" json:"tip,omitempty"`
	User User        `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// ImprovementPlan represents structured improvement plans
type ImprovementPlan struct {
	ID             string    `gorm:"primaryKey" json:"id"`
	SummonerID     string    `gorm:"not null;index" json:"summonerId"`
	PlanType       string    `gorm:"not null;index" json:"planType"` // skill_development, rank_climb, champion_mastery
	Title          string    `gorm:"not null" json:"title"`
	Description    string    `gorm:"type:text" json:"description"`
	Duration       string    `json:"duration"`                        // 4_weeks, 8_weeks, 12_weeks, custom
	Status         string    `gorm:"default:active" json:"status"`    // active, completed, paused, abandoned
	Progress       float64   `json:"progress"`                        // 0-100
	PlanData       string    `gorm:"type:text" json:"planData"`       // JSON stored as text
	MainObjectives string    `gorm:"type:text" json:"mainObjectives"` // JSON array
	DailyRoutine   string    `gorm:"type:text" json:"dailyRoutine"`   // JSON object
	WeeklyGoals    string    `gorm:"type:text" json:"weeklyGoals"`    // JSON array
	Checkpoints    string    `gorm:"type:text" json:"checkpoints"`    // JSON array
	SuccessMetrics string    `gorm:"type:text" json:"successMetrics"` // JSON array
	StartedAt      time.Time `json:"startedAt"`
	CompletedAt    time.Time `json:"completedAt,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// PracticeRoutine represents personalized practice routines
type PracticeRoutine struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	SummonerID    string    `gorm:"not null;index" json:"summonerId"`
	RoutineName   string    `gorm:"not null" json:"routineName"`
	RoutineType   string    `gorm:"not null;index" json:"routineType"` // daily, weekly, custom
	Duration      string    `json:"duration"`                          // total time commitment
	SkillFocus    string    `gorm:"type:text" json:"skillFocus"`       // JSON array
	Phases        string    `gorm:"type:text" json:"phases"`           // JSON array
	Equipment     string    `gorm:"type:text" json:"equipment"`        // JSON array
	Progression   string    `gorm:"type:text" json:"progression"`      // JSON object
	Alternatives  string    `gorm:"type:text" json:"alternatives"`     // JSON array
	Effectiveness float64   `json:"effectiveness"`                     // 0-100
	TimesUsed     int       `json:"timesUsed"`
	AverageRating float64   `json:"averageRating"`                // user ratings
	Status        string    `gorm:"default:active" json:"status"` // active, archived, favorite
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// PracticeSession represents individual practice sessions
type PracticeSession struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	RoutineID       string    `gorm:"index" json:"routineId"`
	SummonerID      string    `gorm:"not null;index" json:"summonerId"`
	SessionType     string    `gorm:"not null;index" json:"sessionType"` // cs_drill, mechanics, vod_review, theory, custom
	FocusAreas      string    `gorm:"type:text" json:"focusAreas"`       // JSON array
	Duration        int       `json:"duration"`                          // actual duration in minutes
	PlannedDuration int       `json:"plannedDuration"`                   // planned duration
	Quality         float64   `json:"quality"`                           // 1-10 subjective rating
	Goals           string    `gorm:"type:text" json:"goals"`            // JSON array
	Achievements    string    `gorm:"type:text" json:"achievements"`     // JSON array
	Notes           string    `gorm:"type:text" json:"notes"`
	ImprovementSeen bool      `json:"improvementSeen"`
	FollowUpNeeded  bool      `json:"followUpNeeded"`
	Effectiveness   float64   `json:"effectiveness"` // 0-100 how effective was the session
	StartedAt       time.Time `gorm:"not null;index" json:"startedAt"`
	CompletedAt     time.Time `json:"completedAt"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Foreign keys
	Routine PracticeRoutine `gorm:"foreignKey:RoutineID;references:ID" json:"routine,omitempty"`
	User    User            `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// TacticalAdvice represents specific tactical advice and tips
type TacticalAdvice struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	AdviceID   string    `gorm:"not null;unique;index" json:"adviceId"`
	SummonerID string    `gorm:"not null;index" json:"summonerId"`
	Category   string    `gorm:"not null;index" json:"category"` // laning, teamfighting, positioning, vision, etc.
	Situation  string    `gorm:"not null" json:"situation"`
	Problem    string    `gorm:"type:text" json:"problem"`
	Solution   string    `gorm:"type:text" json:"solution"`
	Reasoning  string    `gorm:"type:text" json:"reasoning"`
	Examples   string    `gorm:"type:text" json:"examples"` // JSON array of PracticalExample
	Difficulty string    `json:"difficulty"`                // easy, moderate, hard
	Impact     float64   `json:"impact"`                    // 0-100
	Frequency  string    `json:"frequency"`                 // how often this situation occurs
	Urgency    string    `json:"urgency"`                   // critical, high, medium, low
	Related    string    `gorm:"type:text" json:"related"`  // JSON array of related advice IDs
	Applied    bool      `json:"applied"`                   // has user applied this advice
	Helpful    bool      `json:"helpful"`                   // user feedback
	Rating     float64   `json:"rating"`                    // user rating 1-10
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// StrategicGuidance represents high-level strategic guidance
type StrategicGuidance struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	GuidanceID   string    `gorm:"not null;unique;index" json:"guidanceId"`
	SummonerID   string    `gorm:"not null;index" json:"summonerId"`
	StrategyType string    `gorm:"not null;index" json:"strategyType"` // macro, draft, adaptation, win_conditions
	Title        string    `gorm:"not null" json:"title"`
	Overview     string    `gorm:"type:text" json:"overview"`
	Principles   string    `gorm:"type:text" json:"principles"`  // JSON array of StrategicPrinciple
	Application  string    `gorm:"type:text" json:"application"` // JSON object StrategyApplication
	Counters     string    `gorm:"type:text" json:"counters"`    // JSON array of StrategyCounter
	Mastery      string    `gorm:"type:text" json:"mastery"`     // JSON object MasteryProgression
	Advanced     string    `gorm:"type:text" json:"advanced"`    // JSON array of AdvancedConcept
	Difficulty   string    `json:"difficulty"`                   // beginner, intermediate, advanced, expert
	Relevance    float64   `json:"relevance"`                    // 0-100 how relevant to current player level
	Mastered     bool      `json:"mastered"`                     // has player mastered this
	InProgress   bool      `json:"inProgress"`                   // currently working on this
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// MentalCoachingPlan represents mental coaching and mindset development
type MentalCoachingPlan struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	SummonerID         string    `gorm:"not null;index" json:"summonerId"`
	PlanType           string    `gorm:"not null;index" json:"planType"`      // tilt_management, confidence_building, focus_training, stress_management
	CurrentMentalState string    `gorm:"type:text" json:"currentMentalState"` // JSON object
	Goals              string    `gorm:"type:text" json:"goals"`              // JSON array
	Techniques         string    `gorm:"type:text" json:"techniques"`         // JSON array
	DailyPractices     string    `gorm:"type:text" json:"dailyPractices"`     // JSON array
	TiltTriggers       string    `gorm:"type:text" json:"tiltTriggers"`       // JSON array
	CopingStrategies   string    `gorm:"type:text" json:"copingStrategies"`   // JSON array
	ProgressMetrics    string    `gorm:"type:text" json:"progressMetrics"`    // JSON array
	Status             string    `gorm:"default:active" json:"status"`        // active, completed, paused
	Progress           float64   `json:"progress"`                            // 0-100
	Effectiveness      float64   `json:"effectiveness"`                       // user-reported effectiveness
	StartedAt          time.Time `json:"startedAt"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// PerformanceGoal represents specific performance goals
type PerformanceGoal struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	SummonerID      string    `gorm:"not null;index" json:"summonerId"`
	GoalType        string    `gorm:"not null;index" json:"goalType"` // rank, skill_metric, champion_mastery, habit
	Title           string    `gorm:"not null" json:"title"`
	Description     string    `gorm:"type:text" json:"description"`
	Target          string    `gorm:"not null" json:"target"`       // target value
	Current         string    `json:"current"`                      // current value
	Measurement     string    `json:"measurement"`                  // how to measure progress
	Timeline        string    `json:"timeline"`                     // target completion date
	Priority        string    `json:"priority"`                     // critical, high, medium, low
	Status          string    `gorm:"default:active" json:"status"` // active, completed, paused, failed
	Progress        float64   `json:"progress"`                     // 0-100
	Milestones      string    `gorm:"type:text" json:"milestones"`  // JSON array
	Strategies      string    `gorm:"type:text" json:"strategies"`  // JSON array
	Blockers        string    `gorm:"type:text" json:"blockers"`    // JSON array
	Support         string    `gorm:"type:text" json:"support"`     // JSON array
	Achieved        bool      `json:"achieved"`
	AchievementDate time.Time `json:"achievementDate,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// CoachingSchedule represents personalized coaching schedules
type CoachingSchedule struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SummonerID   string    `gorm:"not null;index" json:"summonerId"`
	ScheduleType string    `gorm:"not null" json:"scheduleType"` // daily, weekly, monthly
	Activities   string    `gorm:"type:text" json:"activities"`  // JSON array
	TimeSlots    string    `gorm:"type:text" json:"timeSlots"`   // JSON array
	Reminders    string    `gorm:"type:text" json:"reminders"`   // JSON array
	Flexibility  string    `gorm:"type:text" json:"flexibility"` // JSON object
	Adherence    float64   `json:"adherence"`                    // 0-100 how well they follow schedule
	Adjustments  string    `gorm:"type:text" json:"adjustments"` // JSON array of recent adjustments
	Status       string    `gorm:"default:active" json:"status"` // active, paused, archived
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// ProgressTracking represents coaching progress tracking
type ProgressTracking struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	SummonerID        string    `gorm:"not null;index" json:"summonerId"`
	TrackingPeriod    string    `gorm:"not null;index" json:"trackingPeriod"` // daily, weekly, monthly
	Metrics           string    `gorm:"type:text" json:"metrics"`             // JSON object
	Improvements      string    `gorm:"type:text" json:"improvements"`        // JSON array
	Setbacks          string    `gorm:"type:text" json:"setbacks"`            // JSON array
	Breakthroughs     string    `gorm:"type:text" json:"breakthroughs"`       // JSON array
	OverallProgress   float64   `json:"overallProgress"`                      // 0-100
	MotivationLevel   float64   `json:"motivationLevel"`                      // 0-100
	EngagementLevel   float64   `json:"engagementLevel"`                      // 0-100
	SatisfactionLevel float64   `json:"satisfactionLevel"`                    // 0-100
	Notes             string    `gorm:"type:text" json:"notes"`
	RecordedAt        time.Time `gorm:"not null;index" json:"recordedAt"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// MatchAnalysisInsight represents insights from match analysis
type MatchAnalysisInsight struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SummonerID  string    `gorm:"not null;index" json:"summonerId"`
	MatchID     string    `gorm:"not null;index" json:"matchId"`
	InsightType string    `gorm:"not null;index" json:"insightType"` // tactical, strategic, mechanical, mental
	Category    string    `json:"category"`                          // laning, teamfighting, decision_making, etc.
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Timestamp   int       `json:"timestamp"` // game timestamp in seconds
	Severity    string    `json:"severity"`  // critical, major, minor, positive
	Impact      float64   `json:"impact"`    // -100 to 100
	Actionable  bool      `json:"actionable"`
	Advice      string    `gorm:"type:text" json:"advice"`
	Context     string    `gorm:"type:text" json:"context"` // JSON object with game context
	Reviewed    bool      `json:"reviewed"`                 // has player reviewed this
	Applied     bool      `json:"applied"`                  // has player applied the advice
	Rating      float64   `json:"rating"`                   // user rating of insight quality
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// ChampionCoachingTip represents champion-specific coaching
type ChampionCoachingTip struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SummonerID  string    `gorm:"not null;index" json:"summonerId"`
	Champion    string    `gorm:"not null;index" json:"champion"`
	Role        string    `gorm:"not null;index" json:"role"`
	TipCategory string    `gorm:"not null" json:"tipCategory"` // mechanics, combos, builds, matchups, positioning
	Title       string    `gorm:"not null" json:"title"`
	Content     string    `gorm:"type:text" json:"content"`
	Difficulty  string    `json:"difficulty"`                 // beginner, intermediate, advanced, expert
	Mastery     float64   `json:"mastery"`                    // 0-100 current mastery of this tip
	Priority    string    `json:"priority"`                   // high, medium, low
	Practiced   bool      `json:"practiced"`                  // has player practiced this
	Mastered    bool      `json:"mastered"`                   // has player mastered this
	Examples    string    `gorm:"type:text" json:"examples"`  // JSON array
	Resources   string    `gorm:"type:text" json:"resources"` // JSON array
	Related     string    `gorm:"type:text" json:"related"`   // JSON array of related tip IDs
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Foreign key
	User User `gorm:"foreignKey:SummonerID;references:SummonerID" json:"user,omitempty"`
}

// GORM Hooks
func (c *CoachingInsight) BeforeCreate(tx *gorm.DB) error {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	return nil
}

func (c *CoachingTip) BeforeCreate(tx *gorm.DB) error {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	return nil
}

func (p *PracticeSession) BeforeCreate(tx *gorm.DB) error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	if p.StartedAt.IsZero() {
		p.StartedAt = time.Now()
	}
	return nil
}

// Table names
func (CoachingInsight) TableName() string {
	return "coaching_insights"
}

func (CoachingTip) TableName() string {
	return "coaching_tips"
}

func (CoachingTipFeedback) TableName() string {
	return "coaching_tip_feedback"
}

func (ImprovementPlan) TableName() string {
	return "improvement_plans"
}

func (PracticeRoutine) TableName() string {
	return "practice_routines"
}

func (PracticeSession) TableName() string {
	return "practice_sessions"
}

func (TacticalAdvice) TableName() string {
	return "tactical_advice"
}

func (StrategicGuidance) TableName() string {
	return "strategic_guidance"
}

func (MentalCoachingPlan) TableName() string {
	return "mental_coaching_plans"
}

func (PerformanceGoal) TableName() string {
	return "performance_goals"
}

func (CoachingSchedule) TableName() string {
	return "coaching_schedules"
}

func (ProgressTracking) TableName() string {
	return "progress_tracking"
}

func (MatchAnalysisInsight) TableName() string {
	return "match_analysis_insights"
}

func (ChampionCoachingTip) TableName() string {
	return "champion_coaching_tips"
}
