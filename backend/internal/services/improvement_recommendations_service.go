package services

import (
	"fmt"
	"math"
	"sort"
	"time"

	"gorm.io/gorm"
)

// ImprovementRecommendationsService provides personalized coaching and improvement suggestions
type ImprovementRecommendationsService struct {
	db                *gorm.DB
	analyticsService  *AnalyticsService
	predictiveService *PredictiveAnalyticsService
}

// NewImprovementRecommendationsService creates a new improvement recommendations service
func NewImprovementRecommendationsService(db *gorm.DB, analyticsService *AnalyticsService, predictiveService *PredictiveAnalyticsService) *ImprovementRecommendationsService {
	return &ImprovementRecommendationsService{
		db:                db,
		analyticsService:  analyticsService,
		predictiveService: predictiveService,
	}
}

// ImprovementRecommendation represents a personalized improvement suggestion
type ImprovementRecommendation struct {
	ID               string  `json:"id" gorm:"primaryKey"`
	SummonerID       string  `json:"summoner_id" gorm:"index"`
	Category         string  `json:"category"` // mechanical, macro, mental, champion_specific, etc.
	Priority         string  `json:"priority"` // critical, high, medium, low
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	ImpactScore      float64 `json:"impact_score"`        // 0-100
	DifficultyLevel  string  `json:"difficulty_level"`    // easy, medium, hard, expert
	TimeToSeeResults int     `json:"time_to_see_results"` // days
	EstimatedROI     float64 `json:"estimated_roi"`       // expected improvement %

	// Detailed Action Plan
	ActionPlan ImprovementActionPlan `json:"action_plan" gorm:"embedded"`

	// Progress Tracking
	ProgressTracking ProgressTrackingData `json:"progress_tracking" gorm:"embedded"`

	// Context and Reasoning
	RecommendationContext RecommendationContext `json:"recommendation_context" gorm:"embedded"`

	// Metadata
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ValidUntil time.Time `json:"valid_until"`
	Status     string    `json:"status"` // active, completed, dismissed, expired
}

// ImprovementActionPlan contains the detailed steps to implement the recommendation
type ImprovementActionPlan struct {
	PrimaryObjective    string                 `json:"primary_objective"`
	SecondaryObjectives []string               `json:"secondary_objectives" gorm:"type:text"`
	ActionSteps         []ActionStep           `json:"action_steps" gorm:"type:text"`
	PracticeExercises   []PracticeExercise     `json:"practice_exercises" gorm:"type:text"`
	Resources           []LearningResource     `json:"resources" gorm:"type:text"`
	Milestones          []ImprovementMilestone `json:"milestones" gorm:"type:text"`
	SuccessMetrics      []string               `json:"success_metrics" gorm:"type:text"`
}

// ActionStep represents a specific step in the improvement plan
type ActionStep struct {
	StepNumber    int      `json:"step_number"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Duration      string   `json:"duration"`  // e.g., "15 minutes daily"
	Frequency     string   `json:"frequency"` // e.g., "daily", "3x per week"
	Prerequisites []string `json:"prerequisites"`
	Tools         []string `json:"tools"` // practice tool, replay analysis, etc.
}

// PracticeExercise represents a specific practice routine
type PracticeExercise struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Duration     int      `json:"duration"` // minutes
	Difficulty   string   `json:"difficulty"`
	Focus        []string `json:"focus"` // what skills this targets
	Instructions []string `json:"instructions"`
	Variations   []string `json:"variations"` // different ways to practice
}

// LearningResource represents educational content
type LearningResource struct {
	Type        string `json:"type"` // video, guide, tool, coach
	Title       string `json:"title"`
	URL         string `json:"url,omitempty"`
	Description string `json:"description"`
	Duration    string `json:"duration,omitempty"`
	Difficulty  string `json:"difficulty"`
}

// ImprovementMilestone represents measurable goals
type ImprovementMilestone struct {
	MilestoneNumber int      `json:"milestone_number"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	TargetMetrics   []string `json:"target_metrics"`
	TimeFrame       string   `json:"time_frame"`
	RewardValue     int      `json:"reward_value"` // potential LP/rank improvement
}

// ProgressTrackingData tracks the user's progress on recommendations
type ProgressTrackingData struct {
	OverallProgress    float64               `json:"overall_progress"` // 0-100
	CompletedSteps     []int                 `json:"completed_steps" gorm:"type:text"`
	CurrentMilestone   int                   `json:"current_milestone"`
	MilestoneProgress  []MilestoneProgress   `json:"milestone_progress" gorm:"type:text"`
	WeeklyProgress     []WeeklyProgressData  `json:"weekly_progress" gorm:"type:text"`
	PerformanceImpact  PerformanceImpactData `json:"performance_impact" gorm:"embedded"`
	LastProgressUpdate time.Time             `json:"last_progress_update"`
}

// MilestoneProgress tracks progress on individual milestones
type MilestoneProgress struct {
	MilestoneNumber int                `json:"milestone_number"`
	Progress        float64            `json:"progress"` // 0-100
	StartedAt       time.Time          `json:"started_at"`
	CompletedAt     *time.Time         `json:"completed_at,omitempty"`
	CurrentMetrics  map[string]float64 `json:"current_metrics"`
	TargetMetrics   map[string]float64 `json:"target_metrics"`
}

// WeeklyProgressData tracks weekly improvement metrics
type WeeklyProgressData struct {
	Week                int                `json:"week"`
	StartDate           time.Time          `json:"start_date"`
	PracticeTimeMinutes int                `json:"practice_time_minutes"`
	GamesPlayed         int                `json:"games_played"`
	SkillImprovement    map[string]float64 `json:"skill_improvement"` // skill -> improvement delta
	RankProgress        float64            `json:"rank_progress"`     // LP change
	ConsistencyScore    float64            `json:"consistency_score"` // 0-100
}

// PerformanceImpactData measures the actual impact of following recommendations
type PerformanceImpactData struct {
	BaselineMetrics    map[string]float64 `json:"baseline_metrics"`
	CurrentMetrics     map[string]float64 `json:"current_metrics"`
	ImprovementDeltas  map[string]float64 `json:"improvement_deltas"`
	ROIActual          float64            `json:"roi_actual"`          // actual improvement %
	ROIPredicted       float64            `json:"roi_predicted"`       // predicted improvement %
	ConfidenceInterval float64            `json:"confidence_interval"` // prediction accuracy
}

// RecommendationContext provides context for why this recommendation was made
type RecommendationContext struct {
	TriggeringFactors      []string                    `json:"triggering_factors"` // what led to this recommendation
	DataSources            []string                    `json:"data_sources"`       // recent games, long-term trends, etc.
	AnalysisDepth          string                      `json:"analysis_depth"`     // surface, moderate, deep
	ConfidenceScore        float64                     `json:"confidence_score"`   // 0-100
	AlternativeOptions     []AlternativeRecommendation `json:"alternative_options" gorm:"type:text"`
	PersonalizationFactors PersonalizationData         `json:"personalization_factors" gorm:"embedded"`
}

// AlternativeRecommendation represents other options the user could consider
type AlternativeRecommendation struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImpactScore float64 `json:"impact_score"`
	Difficulty  string  `json:"difficulty"`
	Reason      string  `json:"reason"` // why this wasn't chosen as primary
}

// PersonalizationData contains factors used to personalize the recommendation
type PersonalizationData struct {
	PlayStyle         string            `json:"play_style"`
	LearningStyle     string            `json:"learning_style"`  // visual, hands-on, analytical
	TimeCommitment    string            `json:"time_commitment"` // casual, moderate, intensive
	CurrentRank       string            `json:"current_rank"`
	MainRole          string            `json:"main_role"`
	ChampionPool      []string          `json:"champion_pool" gorm:"type:text"`
	WeakestAreas      []string          `json:"weakest_areas" gorm:"type:text"`
	StrengthAreas     []string          `json:"strength_areas" gorm:"type:text"`
	RecentPerformance RecentPerfSummary `json:"recent_performance" gorm:"embedded"`
	Goals             []string          `json:"goals" gorm:"type:text"`
	Preferences       UserPreferences   `json:"preferences" gorm:"embedded"`
}

// RecentPerfSummary summarizes recent performance trends
type RecentPerfSummary struct {
	Games             int      `json:"games"`
	WinRate           float64  `json:"win_rate"`
	AverageKDA        float64  `json:"average_kda"`
	ConsistencyRating float64  `json:"consistency_rating"`
	ImprovementTrend  string   `json:"improvement_trend"` // improving, stable, declining
	ProblemAreas      []string `json:"problem_areas" gorm:"type:text"`
}

// UserPreferences contains user-specific preferences for recommendations
type UserPreferences struct {
	PreferredDifficulty string   `json:"preferred_difficulty"`
	FocusAreas          []string `json:"focus_areas" gorm:"type:text"`
	AvoidanceAreas      []string `json:"avoidance_areas" gorm:"type:text"`
	PracticeStyle       string   `json:"practice_style"`     // structured, flexible, mixed
	FeedbackFrequency   string   `json:"feedback_frequency"` // daily, weekly, monthly
}

// GetPersonalizedRecommendations generates comprehensive improvement recommendations
func (s *ImprovementRecommendationsService) GetPersonalizedRecommendations(summonerID string, options RecommendationOptions) ([]*ImprovementRecommendation, error) {
	// Analyze current performance and identify improvement areas
	analysis, err := s.analyzePlayerPerformance(summonerID)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze player performance: %w", err)
	}

	// Generate personalized recommendations based on analysis
	recommendations := s.generateRecommendations(analysis, options)

	// Rank and filter recommendations based on impact and feasibility
	rankedRecommendations := s.rankRecommendations(recommendations, analysis)

	// Apply user preferences and constraints
	filteredRecommendations := s.applyUserPreferences(rankedRecommendations, analysis.PersonalizationData)

	// Limit results based on options
	if options.MaxRecommendations > 0 && len(filteredRecommendations) > options.MaxRecommendations {
		filteredRecommendations = filteredRecommendations[:options.MaxRecommendations]
	}

	return filteredRecommendations, nil
}

// RecommendationOptions configures how recommendations are generated
type RecommendationOptions struct {
	FocusCategory       string   `json:"focus_category,omitempty"`    // mechanical, macro, mental, etc.
	DifficultyFilter    string   `json:"difficulty_filter,omitempty"` // easy, medium, hard, expert
	TimeConstraint      int      `json:"time_constraint,omitempty"`   // max time per day in minutes
	MaxRecommendations  int      `json:"max_recommendations"`
	IncludeAlternatives bool     `json:"include_alternatives"`
	PriorityAreas       []string `json:"priority_areas,omitempty"`
}

// PlayerAnalysisResult contains comprehensive player analysis
type PlayerAnalysisResult struct {
	SummonerID             string                  `json:"summoner_id"`
	AnalysisDate           time.Time               `json:"analysis_date"`
	OverallRating          float64                 `json:"overall_rating"`        // 0-100
	SkillBreakdown         map[string]float64      `json:"skill_breakdown"`       // skill -> rating
	ImprovementPotential   map[string]float64      `json:"improvement_potential"` // skill -> potential gain
	CriticalWeaknesses     []CriticalWeakness      `json:"critical_weaknesses"`
	UnderutilizedStrengths []UnderutilizedStrength `json:"underutilized_strengths"`
	PersonalizationData    PersonalizationData     `json:"personalization_data"`
	RecentTrends           RecentTrendAnalysis     `json:"recent_trends"`
	CompetitiveBenchmark   CompetitiveBenchmark    `json:"competitive_benchmark"`
}

// CriticalWeakness represents a major area needing improvement
type CriticalWeakness struct {
	Area            string   `json:"area"`
	Severity        string   `json:"severity"`          // critical, high, medium, low
	ImpactOnWinRate float64  `json:"impact_on_winrate"` // estimated WR improvement if fixed
	Frequency       float64  `json:"frequency"`         // how often this issue occurs
	RootCauses      []string `json:"root_causes"`
	QuickWins       []string `json:"quick_wins"` // easy improvements
}

// UnderutilizedStrength represents strengths that could be leveraged more
type UnderutilizedStrength struct {
	Strength         string   `json:"strength"`
	CurrentUsage     float64  `json:"current_usage"` // 0-100
	OptimalUsage     float64  `json:"optimal_usage"` // 0-100
	LeverageStrategy []string `json:"leverage_strategy"`
	PotentialGain    float64  `json:"potential_gain"` // estimated improvement
}

// RecentTrendAnalysis analyzes recent performance trends
type RecentTrendAnalysis struct {
	TrendPeriodDays       int                  `json:"trend_period_days"`
	OverallTrend          string               `json:"overall_trend"` // improving, stable, declining
	SkillTrends           map[string]string    `json:"skill_trends"`  // skill -> trend
	ConsistencyTrend      string               `json:"consistency_trend"`
	PerformanceVolatility float64              `json:"performance_volatility"`
	RecentBreakthroughs   []RecentBreakthrough `json:"recent_breakthroughs"`
	RecentStruggles       []RecentStruggle     `json:"recent_struggles"`
}

// RecentBreakthrough represents recent positive developments
type RecentBreakthrough struct {
	Area           string    `json:"area"`
	Description    string    `json:"description"`
	Impact         float64   `json:"impact"`
	Date           time.Time `json:"date"`
	Sustainability string    `json:"sustainability"` // high, medium, low
}

// RecentStruggle represents recent performance issues
type RecentStruggle struct {
	Area        string  `json:"area"`
	Description string  `json:"description"`
	Frequency   float64 `json:"frequency"` // how often it happens
	Severity    string  `json:"severity"`
	Pattern     string  `json:"pattern"` // situational, consistent, random
	Trend       string  `json:"trend"`   // worsening, stable, improving
}

// CompetitiveBenchmark compares player to others at similar skill level
type CompetitiveBenchmark struct {
	RankTier                 string             `json:"rank_tier"`
	RegionalPercentile       float64            `json:"regional_percentile"`
	SkillPercentiles         map[string]float64 `json:"skill_percentiles"` // skill -> percentile
	StrongerThanPeers        []string           `json:"stronger_than_peers"`
	WeakerThanPeers          []string           `json:"weaker_than_peers"`
	CompetitiveAdvantages    []string           `json:"competitive_advantages"`
	CompetitiveDisadvantages []string           `json:"competitive_disadvantages"`
}

// analyzePlayerPerformance conducts comprehensive player analysis
func (s *ImprovementRecommendationsService) analyzePlayerPerformance(summonerID string) (*PlayerAnalysisResult, error) {
	// This would integrate with all other analytics services
	// For now, return mock analysis data

	analysis := &PlayerAnalysisResult{
		SummonerID:    summonerID,
		AnalysisDate:  time.Now(),
		OverallRating: 72.5,
		SkillBreakdown: map[string]float64{
			"mechanical_skill":  68.0,
			"game_knowledge":    75.0,
			"map_awareness":     65.0,
			"team_fighting":     78.0,
			"laning":            70.0,
			"objective_control": 73.0,
			"vision_control":    60.0,
			"positioning":       72.0,
			"decision_making":   74.0,
			"mental_resilience": 69.0,
		},
		ImprovementPotential: map[string]float64{
			"vision_control":    20.0, // highest potential
			"map_awareness":     18.0,
			"mechanical_skill":  15.0,
			"positioning":       12.0,
			"mental_resilience": 12.0,
			"laning":            10.0,
			"objective_control": 8.0,
			"decision_making":   7.0,
			"game_knowledge":    5.0,
			"team_fighting":     3.0, // already strong
		},
		CriticalWeaknesses: []CriticalWeakness{
			{
				Area:            "vision_control",
				Severity:        "high",
				ImpactOnWinRate: 8.5,
				Frequency:       85.0,
				RootCauses:      []string{"infrequent ward placement", "poor ward positioning", "not clearing enemy vision"},
				QuickWins:       []string{"buy more control wards", "ward before major objectives", "use trinket on cooldown"},
			},
			{
				Area:            "map_awareness",
				Severity:        "medium",
				ImpactOnWinRate: 6.2,
				Frequency:       70.0,
				RootCauses:      []string{"tunnel vision during farming", "not tracking enemy jungler", "poor minimap usage"},
				QuickWins:       []string{"look at minimap every 5 seconds", "ping missing enemies", "ward river bushes"},
			},
		},
		UnderutilizedStrengths: []UnderutilizedStrength{
			{
				Strength:         "team_fighting",
				CurrentUsage:     65.0,
				OptimalUsage:     85.0,
				LeverageStrategy: []string{"engage more team fights", "position more aggressively", "follow up on team plays"},
				PotentialGain:    5.5,
			},
		},
		PersonalizationData: PersonalizationData{
			PlayStyle:      "teamfight_oriented",
			LearningStyle:  "hands_on",
			TimeCommitment: "moderate",
			CurrentRank:    "Gold II",
			MainRole:       "ADC",
			ChampionPool:   []string{"Jinx", "Kai'Sa", "Ezreal"},
			WeakestAreas:   []string{"vision_control", "map_awareness", "positioning"},
			StrengthAreas:  []string{"team_fighting", "decision_making", "game_knowledge"},
			RecentPerformance: RecentPerfSummary{
				Games:             20,
				WinRate:           58.0,
				AverageKDA:        2.1,
				ConsistencyRating: 72.0,
				ImprovementTrend:  "stable",
				ProblemAreas:      []string{"vision_control", "early_game_positioning"},
			},
			Goals: []string{"reach_platinum", "improve_consistency", "better_teamwork"},
			Preferences: UserPreferences{
				PreferredDifficulty: "medium",
				FocusAreas:          []string{"vision_control", "positioning"},
				PracticeStyle:       "structured",
				FeedbackFrequency:   "weekly",
			},
		},
		RecentTrends: RecentTrendAnalysis{
			TrendPeriodDays: 14,
			OverallTrend:    "improving",
			SkillTrends: map[string]string{
				"vision_control": "improving",
				"team_fighting":  "stable",
				"positioning":    "declining",
			},
			ConsistencyTrend:      "stable",
			PerformanceVolatility: 15.2,
		},
		CompetitiveBenchmark: CompetitiveBenchmark{
			RankTier:           "Gold II",
			RegionalPercentile: 68.5,
			SkillPercentiles: map[string]float64{
				"team_fighting":   82.0, // above average
				"decision_making": 75.0,
				"game_knowledge":  73.0,
				"vision_control":  45.0, // below average
				"map_awareness":   52.0,
			},
			StrongerThanPeers: []string{"team_fighting", "decision_making"},
			WeakerThanPeers:   []string{"vision_control", "map_awareness"},
		},
	}

	return analysis, nil
}

// generateRecommendations creates personalized improvement recommendations
func (s *ImprovementRecommendationsService) generateRecommendations(analysis *PlayerAnalysisResult, options RecommendationOptions) []*ImprovementRecommendation {
	var recommendations []*ImprovementRecommendation

	// Generate recommendations for critical weaknesses
	for _, weakness := range analysis.CriticalWeaknesses {
		rec := s.createWeaknessRecommendation(analysis, weakness)
		recommendations = append(recommendations, rec)
	}

	// Generate recommendations for underutilized strengths
	for _, strength := range analysis.UnderutilizedStrengths {
		rec := s.createStrengthLeverageRecommendation(analysis, strength)
		recommendations = append(recommendations, rec)
	}

	// Generate recommendations based on recent trends
	trendRecs := s.createTrendBasedRecommendations(analysis)
	recommendations = append(recommendations, trendRecs...)

	// Generate general improvement recommendations
	generalRecs := s.createGeneralImprovementRecommendations(analysis)
	recommendations = append(recommendations, generalRecs...)

	return recommendations
}

// createWeaknessRecommendation creates a recommendation to address a critical weakness
func (s *ImprovementRecommendationsService) createWeaknessRecommendation(analysis *PlayerAnalysisResult, weakness CriticalWeakness) *ImprovementRecommendation {
	rec := &ImprovementRecommendation{
		ID:               fmt.Sprintf("weakness_%s_%s", analysis.SummonerID, weakness.Area),
		SummonerID:       analysis.SummonerID,
		Category:         weakness.Area,
		Priority:         weakness.Severity,
		Title:            fmt.Sprintf("Improve %s", formatSkillName(weakness.Area)),
		Description:      fmt.Sprintf("Address critical weakness in %s to improve win rate by up to %.1f%%", weakness.Area, weakness.ImpactOnWinRate),
		ImpactScore:      weakness.ImpactOnWinRate * 10, // convert to 0-100 scale
		DifficultyLevel:  s.calculateDifficultyLevel(weakness.Area, analysis.PersonalizationData),
		TimeToSeeResults: s.estimateTimeToResults(weakness.Area, weakness.Severity),
		EstimatedROI:     weakness.ImpactOnWinRate,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		ValidUntil:       time.Now().AddDate(0, 0, 30), // 30 days validity
		Status:           "active",
	}

	// Create detailed action plan
	rec.ActionPlan = s.createActionPlan(weakness.Area, analysis.PersonalizationData)

	// Initialize progress tracking
	rec.ProgressTracking = ProgressTrackingData{
		OverallProgress:    0.0,
		CompletedSteps:     []int{},
		CurrentMilestone:   1,
		LastProgressUpdate: time.Now(),
		PerformanceImpact: PerformanceImpactData{
			BaselineMetrics: map[string]float64{
				weakness.Area: analysis.SkillBreakdown[weakness.Area],
			},
			ROIPredicted: weakness.ImpactOnWinRate,
		},
	}

	// Set recommendation context
	rec.RecommendationContext = RecommendationContext{
		TriggeringFactors:      weakness.RootCauses,
		DataSources:            []string{"recent_matches", "skill_analysis", "competitive_benchmark"},
		AnalysisDepth:          "deep",
		ConfidenceScore:        85.0,
		PersonalizationFactors: analysis.PersonalizationData,
	}

	return rec
}

// createStrengthLeverageRecommendation creates a recommendation to better utilize strengths
func (s *ImprovementRecommendationsService) createStrengthLeverageRecommendation(analysis *PlayerAnalysisResult, strength UnderutilizedStrength) *ImprovementRecommendation {
	rec := &ImprovementRecommendation{
		ID:               fmt.Sprintf("strength_%s_%s", analysis.SummonerID, strength.Strength),
		SummonerID:       analysis.SummonerID,
		Category:         "strength_leverage",
		Priority:         "medium",
		Title:            fmt.Sprintf("Better Leverage Your %s", formatSkillName(strength.Strength)),
		Description:      fmt.Sprintf("You're strong at %s but not fully utilizing it. Potential %.1f%% improvement.", strength.Strength, strength.PotentialGain),
		ImpactScore:      strength.PotentialGain * 10,
		DifficultyLevel:  "medium",
		TimeToSeeResults: 7, // usually quicker to leverage existing strengths
		EstimatedROI:     strength.PotentialGain,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		ValidUntil:       time.Now().AddDate(0, 0, 21), // 21 days validity
		Status:           "active",
	}

	// Create action plan for leveraging strength
	rec.ActionPlan = ImprovementActionPlan{
		PrimaryObjective:    fmt.Sprintf("Increase utilization of %s from %.1f%% to %.1f%%", strength.Strength, strength.CurrentUsage, strength.OptimalUsage),
		SecondaryObjectives: []string{"Identify more opportunities", "Build confidence in strength", "Develop situational awareness"},
		ActionSteps: []ActionStep{
			{
				StepNumber:  1,
				Title:       "Identify Opportunities",
				Description: "Recognize situations where you can leverage this strength",
				Duration:    "10 minutes per game",
				Frequency:   "every game",
				Tools:       []string{"mental_checklist", "replay_review"},
			},
			{
				StepNumber:  2,
				Title:       "Practice Active Usage",
				Description: "Consciously apply this strength more often",
				Duration:    "entire game",
				Frequency:   "daily",
				Tools:       []string{"in_game_reminders", "buddy_system"},
			},
		},
		SuccessMetrics: []string{
			fmt.Sprintf("Increase %s utilization to %.1f%%", strength.Strength, strength.OptimalUsage),
			"Consistent positive impact in games",
			"Improved confidence in ability",
		},
	}

	return rec
}

// Helper functions

func formatSkillName(skill string) string {
	switch skill {
	case "vision_control":
		return "Vision Control"
	case "map_awareness":
		return "Map Awareness"
	case "team_fighting":
		return "Team Fighting"
	case "mechanical_skill":
		return "Mechanical Skill"
	case "positioning":
		return "Positioning"
	case "decision_making":
		return "Decision Making"
	default:
		return skill
	}
}

func (s *ImprovementRecommendationsService) calculateDifficultyLevel(area string, personalData PersonalizationData) string {
	// Consider user's learning style, time commitment, and current skill level
	if personalData.PreferredDifficulty == "easy" {
		return "easy"
	}

	// Some areas are inherently more difficult
	hardAreas := []string{"decision_making", "macro_play", "shot_calling"}
	for _, hardArea := range hardAreas {
		if area == hardArea {
			return "hard"
		}
	}

	return "medium"
}

func (s *ImprovementRecommendationsService) estimateTimeToResults(area string, severity string) int {
	baseTime := map[string]int{
		"vision_control":   7,  // quick wins possible
		"map_awareness":    14, // habit formation
		"positioning":      21, // muscle memory
		"mechanical_skill": 28, // practice intensive
		"decision_making":  35, // experience needed
	}

	days := baseTime[area]
	if days == 0 {
		days = 14 // default
	}

	// Adjust based on severity
	switch severity {
	case "critical":
		days = int(float64(days) * 1.5) // takes longer for critical issues
	case "high":
		days = int(float64(days) * 1.2)
	}

	return days
}

func (s *ImprovementRecommendationsService) createActionPlan(area string, personalData PersonalizationData) ImprovementActionPlan {
	// This would create detailed, personalized action plans based on the weakness area
	// For brevity, returning a template plan

	plan := ImprovementActionPlan{
		PrimaryObjective: fmt.Sprintf("Significantly improve %s performance", area),
		ActionSteps: []ActionStep{
			{
				StepNumber:  1,
				Title:       "Foundation Assessment",
				Description: "Evaluate current level and identify specific gaps",
				Duration:    "30 minutes",
				Frequency:   "one-time",
				Tools:       []string{"self_assessment", "replay_analysis"},
			},
			{
				StepNumber:  2,
				Title:       "Targeted Practice",
				Description: "Focus on specific weaknesses through deliberate practice",
				Duration:    "20 minutes daily",
				Frequency:   "daily",
				Tools:       []string{"practice_tool", "custom_games"},
			},
		},
		Milestones: []ImprovementMilestone{
			{
				MilestoneNumber: 1,
				Title:           "Foundation Established",
				Description:     "Basic improvements visible in gameplay",
				TimeFrame:       "1 week",
				RewardValue:     50, // LP equivalent
			},
		},
		SuccessMetrics: []string{
			"Measurable improvement in relevant statistics",
			"Positive trend in match performance",
			"Increased consistency",
		},
	}

	return plan
}

func (s *ImprovementRecommendationsService) createTrendBasedRecommendations(analysis *PlayerAnalysisResult) []*ImprovementRecommendation {
	// Create recommendations based on recent trends
	var recommendations []*ImprovementRecommendation

	// Example: if positioning is declining, create a recommendation to address it
	if analysis.RecentTrends.SkillTrends["positioning"] == "declining" {
		rec := &ImprovementRecommendation{
			ID:               fmt.Sprintf("trend_%s_positioning", analysis.SummonerID),
			SummonerID:       analysis.SummonerID,
			Category:         "positioning",
			Priority:         "high",
			Title:            "Address Declining Positioning",
			Description:      "Your positioning has been declining recently - let's get it back on track",
			ImpactScore:      70.0,
			DifficultyLevel:  "medium",
			TimeToSeeResults: 10,
			EstimatedROI:     6.0,
			CreatedAt:        time.Now(),
			Status:           "active",
		}
		recommendations = append(recommendations, rec)
	}

	return recommendations
}

func (s *ImprovementRecommendationsService) createGeneralImprovementRecommendations(analysis *PlayerAnalysisResult) []*ImprovementRecommendation {
	// Create general recommendations based on rank and role
	var recommendations []*ImprovementRecommendation

	// Add mental resilience recommendation for all players
	rec := &ImprovementRecommendation{
		ID:               fmt.Sprintf("general_%s_mental", analysis.SummonerID),
		SummonerID:       analysis.SummonerID,
		Category:         "mental_game",
		Priority:         "medium",
		Title:            "Strengthen Mental Resilience",
		Description:      "Develop better mental game and tilt resistance for more consistent performance",
		ImpactScore:      50.0,
		DifficultyLevel:  "medium",
		TimeToSeeResults: 14,
		EstimatedROI:     4.0,
		CreatedAt:        time.Now(),
		Status:           "active",
	}

	recommendations = append(recommendations, rec)

	return recommendations
}

// rankRecommendations sorts recommendations by priority and impact
func (s *ImprovementRecommendationsService) rankRecommendations(recommendations []*ImprovementRecommendation, analysis *PlayerAnalysisResult) []*ImprovementRecommendation {
	sort.Slice(recommendations, func(i, j int) bool {
		// Primary sort by priority
		priorityOrder := map[string]int{"critical": 4, "high": 3, "medium": 2, "low": 1}
		if priorityOrder[recommendations[i].Priority] != priorityOrder[recommendations[j].Priority] {
			return priorityOrder[recommendations[i].Priority] > priorityOrder[recommendations[j].Priority]
		}

		// Secondary sort by impact score
		return recommendations[i].ImpactScore > recommendations[j].ImpactScore
	})

	return recommendations
}

// applyUserPreferences filters and adjusts recommendations based on user preferences
func (s *ImprovementRecommendationsService) applyUserPreferences(recommendations []*ImprovementRecommendation, personalData PersonalizationData) []*ImprovementRecommendation {
	var filtered []*ImprovementRecommendation

	for _, rec := range recommendations {
		// Check if recommendation matches user's focus areas
		if len(personalData.Preferences.FocusAreas) > 0 {
			match := false
			for _, focus := range personalData.Preferences.FocusAreas {
				if rec.Category == focus {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		// Check if recommendation is in avoidance areas
		avoid := false
		for _, avoidArea := range personalData.Preferences.AvoidanceAreas {
			if rec.Category == avoidArea {
				avoid = true
				break
			}
		}
		if avoid {
			continue
		}

		// Check difficulty preference
		if personalData.Preferences.PreferredDifficulty != "" && rec.DifficultyLevel != personalData.Preferences.PreferredDifficulty {
			// Allow recommendations one level up or down from preferred difficulty
			difficultyOrder := map[string]int{"easy": 1, "medium": 2, "hard": 3, "expert": 4}
			userLevel := difficultyOrder[personalData.Preferences.PreferredDifficulty]
			recLevel := difficultyOrder[rec.DifficultyLevel]
			if math.Abs(float64(userLevel-recLevel)) > 1 {
				continue
			}
		}

		filtered = append(filtered, rec)
	}

	return filtered
}

// GetRecommendationProgress gets progress tracking data for a recommendation
func (s *ImprovementRecommendationsService) GetRecommendationProgress(recommendationID string) (*ProgressTrackingData, error) {
	var rec ImprovementRecommendation
	if err := s.db.Where("id = ?", recommendationID).First(&rec).Error; err != nil {
		return nil, fmt.Errorf("recommendation not found: %w", err)
	}

	return &rec.ProgressTracking, nil
}

// UpdateRecommendationProgress updates progress on a recommendation
func (s *ImprovementRecommendationsService) UpdateRecommendationProgress(recommendationID string, progress ProgressTrackingData) error {
	return s.db.Model(&ImprovementRecommendation{}).Where("id = ?", recommendationID).
		Updates(map[string]interface{}{
			"progress_tracking": progress,
			"updated_at":        time.Now(),
		}).Error
}

// GetActiveRecommendations gets all active recommendations for a summoner
func (s *ImprovementRecommendationsService) GetActiveRecommendations(summonerID string) ([]*ImprovementRecommendation, error) {
	var recommendations []*ImprovementRecommendation
	err := s.db.Where("summoner_id = ? AND status = ? AND valid_until > ?",
		summonerID, "active", time.Now()).
		Order("priority DESC, impact_score DESC").
		Find(&recommendations).Error

	return recommendations, err
}

// CompleteRecommendation marks a recommendation as completed
func (s *ImprovementRecommendationsService) CompleteRecommendation(recommendationID string) error {
	return s.db.Model(&ImprovementRecommendation{}).Where("id = ?", recommendationID).
		Updates(map[string]interface{}{
			"status":     "completed",
			"updated_at": time.Now(),
		}).Error
}
