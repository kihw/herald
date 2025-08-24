package services

import "time"

// MatchAnalysisInsight represents insights from match analysis
type MatchAnalysisInsight struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Impact      string    `json:"impact"`
	Priority    int       `json:"priority"`
	Timestamp   time.Time `json:"timestamp"`
}

// TiltManagementPlan represents tilt management strategies
type TiltManagementPlan struct {
	ID                   string            `json:"id"`
	TriggerPatterns      []string          `json:"trigger_patterns"`
	PreventionTechniques []string          `json:"prevention_techniques"`
	RecoveryStrategies   []string          `json:"recovery_strategies"`
	BreakRecommendations map[string]string `json:"break_recommendations"`
	MindsetExercises     []MindsetExercise `json:"mindset_exercises"`
	ProgressTracking     ProgressMetrics   `json:"progress_tracking"`
	CreatedAt            time.Time         `json:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at"`
}

// ConfidencePlan represents confidence building strategies
type ConfidencePlan struct {
	ID                   string               `json:"id"`
	CurrentLevel         float64              `json:"current_level"`
	TargetLevel          float64              `json:"target_level"`
	BuildingExercises    []ConfidenceExercise `json:"building_exercises"`
	SuccessTracking      []Achievement        `json:"success_tracking"`
	AffirmationPractices []string             `json:"affirmation_practices"`
	GoalMilestones       []Milestone          `json:"goal_milestones"`
	ProgressMetrics      ProgressMetrics      `json:"progress_metrics"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

// FocusTraining represents focus improvement training
type FocusTraining struct {
	ID                  string          `json:"id"`
	ConcentrationLevel  float64         `json:"concentration_level"`
	AttentionSpan       time.Duration   `json:"attention_span"`
	DistractionFactors  []string        `json:"distraction_factors"`
	FocusExercises      []FocusExercise `json:"focus_exercises"`
	MeditationPractices []string        `json:"meditation_practices"`
	EnvironmentTips     []string        `json:"environment_tips"`
	ProgressTracking    ProgressMetrics `json:"progress_tracking"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

// MotivationSystem represents motivation enhancement system
type MotivationSystem struct {
	ID                   string          `json:"id"`
	MotivationLevel      float64         `json:"motivation_level"`
	PersonalGoals        []PersonalGoal  `json:"personal_goals"`
	RewardSystem         RewardSystem    `json:"reward_system"`
	ChallengeProgression []Challenge     `json:"challenge_progression"`
	SocialMotivation     SocialFactors   `json:"social_motivation"`
	ProgressTracking     ProgressMetrics `json:"progress_tracking"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

// StressManagement represents stress management strategies
type StressManagement struct {
	ID                   string                `json:"id"`
	StressLevel          float64               `json:"stress_level"`
	StressTriggers       []string              `json:"stress_triggers"`
	RelaxationTechniques []RelaxationTechnique `json:"relaxation_techniques"`
	BreathingExercises   []BreathingExercise   `json:"breathing_exercises"`
	PhysicalExercises    []PhysicalExercise    `json:"physical_exercises"`
	TimeManagement       TimeManagementPlan    `json:"time_management"`
	ProgressTracking     ProgressMetrics       `json:"progress_tracking"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
}

// MindfulnessProgram represents mindfulness training program
type MindfulnessProgram struct {
	ID                   string                `json:"id"`
	CurrentLevel         int                   `json:"current_level"`
	PracticeStreak       int                   `json:"practice_streak"`
	MeditationSessions   []MeditationSession   `json:"meditation_sessions"`
	MindfulnessExercises []MindfulnessExercise `json:"mindfulness_exercises"`
	AwarenessActivities  []AwarenessActivity   `json:"awareness_activities"`
	DailyPractices       []DailyPractice       `json:"daily_practices"`
	ProgressTracking     ProgressMetrics       `json:"progress_tracking"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
}

// Supporting types
type MindsetExercise struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Duration    time.Duration `json:"duration"`
	Frequency   string        `json:"frequency"`
}

type ProgressMetrics struct {
	StartDate       time.Time `json:"start_date"`
	CurrentProgress float64   `json:"current_progress"`
	TargetProgress  float64   `json:"target_progress"`
	WeeklyGoals     []string  `json:"weekly_goals"`
	Milestones      []string  `json:"milestones"`
}

type ConfidenceExercise struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Steps     []string `json:"steps"`
	Frequency string   `json:"frequency"`
}

type Achievement struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Points      int       `json:"points"`
	AchievedAt  time.Time `json:"achieved_at"`
}

type Milestone struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Target    float64   `json:"target"`
	Current   float64   `json:"current"`
	Deadline  time.Time `json:"deadline"`
	Completed bool      `json:"completed"`
}

type FocusExercise struct {
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Duration    time.Duration `json:"duration"`
	Description string        `json:"description"`
}

type PersonalGoal struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Deadline    time.Time `json:"deadline"`
	Progress    float64   `json:"progress"`
}

type RewardSystem struct {
	Points          int      `json:"points"`
	Level           int      `json:"level"`
	UnlockedRewards []string `json:"unlocked_rewards"`
	NextRewards     []string `json:"next_rewards"`
	BadgeCollection []Badge  `json:"badge_collection"`
}

type Badge struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	EarnedAt    time.Time `json:"earned_at"`
}

type Challenge struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Difficulty  int     `json:"difficulty"`
	Reward      int     `json:"reward"`
	Progress    float64 `json:"progress"`
	Completed   bool    `json:"completed"`
}

type SocialFactors struct {
	FriendsSupport    []string `json:"friends_support"`
	TeamMotivation    []string `json:"team_motivation"`
	CommunityGoals    []string `json:"community_goals"`
	CompetitiveSpirit bool     `json:"competitive_spirit"`
}

type RelaxationTechnique struct {
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Duration     time.Duration `json:"duration"`
	Instructions []string      `json:"instructions"`
}

type BreathingExercise struct {
	Name     string   `json:"name"`
	Pattern  string   `json:"pattern"`
	Duration int      `json:"duration"`
	Steps    []string `json:"steps"`
}

type PhysicalExercise struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Duration     int      `json:"duration"`
	Instructions []string `json:"instructions"`
}

type TimeManagementPlan struct {
	DailySchedule  map[string]string `json:"daily_schedule"`
	BreakIntervals []string          `json:"break_intervals"`
	PriorityTasks  []string          `json:"priority_tasks"`
	TimeBlocking   bool              `json:"time_blocking"`
}

type MeditationSession struct {
	ID          string        `json:"id"`
	Type        string        `json:"type"`
	Duration    time.Duration `json:"duration"`
	Theme       string        `json:"theme"`
	CompletedAt time.Time     `json:"completed_at"`
}

type MindfulnessExercise struct {
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	Instructions []string `json:"instructions"`
	Duration     int      `json:"duration"`
}

type AwarenessActivity struct {
	Name        string `json:"name"`
	Focus       string `json:"focus"`
	Description string `json:"description"`
	Frequency   string `json:"frequency"`
}

type DailyPractice struct {
	Name      string `json:"name"`
	Time      string `json:"time"`
	Duration  int    `json:"duration"`
	Completed bool   `json:"completed"`
}

// Additional coaching types for Herald.lol gaming platform
type ChampionCoachingTip struct {
	Champion   string   `json:"champion"`
	Role       string   `json:"role"`
	Tips       []string `json:"tips"`
	Priority   int      `json:"priority"`
	Difficulty string   `json:"difficulty"`
}

type MetaAdaptationGuidance struct {
	MetaPatch       string   `json:"meta_patch"`
	Adaptations     []string `json:"adaptations"`
	ChampionShifts  []string `json:"champion_shifts"`
	StrategyChanges []string `json:"strategy_changes"`
}

type PerformanceGoal struct {
	Metric     string  `json:"metric"`
	Current    float64 `json:"current"`
	Target     float64 `json:"target"`
	Timeline   string  `json:"timeline"`
	Difficulty string  `json:"difficulty"`
}

type CoachingSchedule struct {
	SessionType string    `json:"session_type"`
	Duration    int       `json:"duration"`
	Frequency   string    `json:"frequency"`
	Goals       []string  `json:"goals"`
	StartTime   time.Time `json:"start_time"`
}

type ProgressTracking struct {
	PlayerID     string             `json:"player_id"`
	Goals        []PerformanceGoal  `json:"goals"`
	Achievements []string           `json:"achievements"`
	LastUpdated  time.Time          `json:"last_updated"`
	Progress     map[string]float64 `json:"progress"`
}
