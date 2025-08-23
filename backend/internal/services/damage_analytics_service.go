package services

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/herald/internal/models"
)

// DamageAnalyticsService handles damage and team contribution analysis
type DamageAnalyticsService struct {
	analyticsService *AnalyticsService
	matchService     *MatchService
}

// DamageAnalysis represents comprehensive damage analysis
type DamageAnalysis struct {
	PlayerID          string                    `json:"player_id"`
	Champion          string                    `json:"champion,omitempty"`
	Position          string                    `json:"position,omitempty"`
	TimeRange         string                    `json:"time_range"`
	
	// Core Damage Metrics
	TotalDamageDealt     float64               `json:"total_damage_dealt"`
	AverageDamagePerGame float64               `json:"average_damage_per_game"`
	DamagePerMinute      float64               `json:"damage_per_minute"`
	DamageShare          float64               `json:"damage_share_percent"`
	
	// Damage Type Breakdown
	PhysicalDamageShare  float64               `json:"physical_damage_share"`
	MagicDamageShare     float64               `json:"magic_damage_share"`
	TrueDamageShare      float64               `json:"true_damage_share"`
	
	// Damage Efficiency
	DamagePerGold        float64               `json:"damage_per_gold"`
	DamagePerItem        float64               `json:"damage_per_item"`
	DamageToChampions    float64               `json:"damage_to_champions_percent"`
	
	// Team Contribution Analysis
	TeamContribution     TeamContributionData  `json:"team_contribution"`
	
	// Comparative Analysis
	RoleBenchmark        DamageBenchmark       `json:"role_benchmark"`
	RankBenchmark        DamageBenchmark       `json:"rank_benchmark"`
	ChampionBenchmark    DamageBenchmark       `json:"champion_benchmark"`
	
	// Performance Analysis
	DamageConsistency    float64               `json:"damage_consistency"` // Lower std dev is better
	DamageReliability    string                `json:"damage_reliability"` // "consistent", "variable", "inconsistent"
	CarryPotential       float64               `json:"carry_potential"`    // 0-100 score
	
	// Game Phase Analysis
	EarlyGameDamage      GamePhaseDamage       `json:"early_game_damage"`
	MidGameDamage        GamePhaseDamage       `json:"mid_game_damage"`
	LateGameDamage       GamePhaseDamage       `json:"late_game_damage"`
	
	// Situational Analysis
	WinningGamesDamage   DamageSnapshot        `json:"winning_games_damage"`
	LosingGamesDamage    DamageSnapshot        `json:"losing_games_damage"`
	CloseGamesDamage     DamageSnapshot        `json:"close_games_damage"`
	
	// Trend Analysis
	TrendDirection       string                `json:"trend_direction"`
	TrendSlope           float64               `json:"trend_slope"`
	TrendConfidence      float64               `json:"trend_confidence"`
	
	// Insights and Recommendations
	StrengthAreas        []string              `json:"strength_areas"`
	ImprovementAreas     []string              `json:"improvement_areas"`
	DamageRecommendations []DamageRecommendation `json:"recommendations"`
	
	// Historical Data
	TrendData            []DamageTrendPoint    `json:"trend_data"`
	RecentMatches        []MatchDamageData     `json:"recent_matches"`
}

// TeamContributionData represents player's contribution to team success
type TeamContributionData struct {
	// Damage Contribution
	TeamDamageShare         float64   `json:"team_damage_share"`
	DamageCarryRate         float64   `json:"damage_carry_rate"`      // % of games where top damage
	DamageContributionScore float64   `json:"damage_contribution_score"` // 0-100
	
	// Kill Participation
	KillParticipation       float64   `json:"kill_participation"`
	FirstBloodInvolvement   float64   `json:"first_blood_involvement"`
	TeamfightPresence       float64   `json:"teamfight_presence"`
	
	// Objective Control
	DragonDamageShare       float64   `json:"dragon_damage_share"`
	BaronDamageShare        float64   `json:"baron_damage_share"`
	TowerDamageShare        float64   `json:"tower_damage_share"`
	
	// Support Metrics (for all roles)
	DamageAmplification     float64   `json:"damage_amplification"`  // Damage enabled for others
	DamageReduction         float64   `json:"damage_reduction"`      // Damage prevented/healed
	
	// Team Impact Score
	OverallTeamImpact       float64   `json:"overall_team_impact"`   // 0-100 combined score
	TeamImpactRanking       string    `json:"team_impact_ranking"`   // "excellent", "good", "average", "poor"
}

// GamePhaseDamage represents damage metrics for a specific game phase
type GamePhaseDamage struct {
	Phase                   string    `json:"phase"` // "early", "mid", "late"
	AverageDamage           float64   `json:"average_damage"`
	DamageShare             float64   `json:"damage_share"`
	DamagePerMinute         float64   `json:"damage_per_minute"`
	ConsistencyRating       string    `json:"consistency_rating"`
	RelativePerformance     float64   `json:"relative_performance"` // vs role average
}

// DamageSnapshot represents damage performance in specific game contexts
type DamageSnapshot struct {
	Context                 string    `json:"context"` // "winning", "losing", "close"
	GameCount              int       `json:"game_count"`
	AverageDamage          float64   `json:"average_damage"`
	DamageShare            float64   `json:"damage_share"`
	PerformanceRating      string    `json:"performance_rating"`
	ImpactDifference       float64   `json:"impact_difference"` // vs overall average
}

// DamageBenchmark represents damage performance benchmarks
type DamageBenchmark struct {
	Category               string    `json:"category"` // "role", "rank", "champion"
	Filter                 string    `json:"filter"`   // "ADC", "GOLD", "Jinx"
	AverageDamageShare     float64   `json:"average_damage_share"`
	MedianDamageShare      float64   `json:"median_damage_share"`
	Top10PercentDamage     float64   `json:"top_10_percent_damage"`
	Top25PercentDamage     float64   `json:"top_25_percent_damage"`
	PlayerPercentile       float64   `json:"player_percentile"`
	ComparisonRating       string    `json:"comparison_rating"` // "excellent", "above_average", "average", "below_average", "poor"
}

// DamageRecommendation represents actionable damage improvement advice
type DamageRecommendation struct {
	Priority               string    `json:"priority"` // "high", "medium", "low"
	Category               string    `json:"category"` // "itemization", "positioning", "target_selection", "timing"
	Title                  string    `json:"title"`
	Description            string    `json:"description"`
	ExpectedImprovement    string    `json:"expected_improvement"`
	GamePhase              []string  `json:"game_phase"`
	ItemSuggestions        []string  `json:"item_suggestions,omitempty"`
	PositioningTips        []string  `json:"positioning_tips,omitempty"`
}

// DamageTrendPoint represents damage performance over time
type DamageTrendPoint struct {
	Date                   time.Time `json:"date"`
	DamageShare            float64   `json:"damage_share"`
	DamagePerMinute        float64   `json:"damage_per_minute"`
	TeamContribution       float64   `json:"team_contribution"`
	MovingAverage          float64   `json:"moving_average"`
	CarryPotential         float64   `json:"carry_potential"`
}

// MatchDamageData represents damage data from a specific match
type MatchDamageData struct {
	MatchID                string    `json:"match_id"`
	Champion               string    `json:"champion"`
	Position               string    `json:"position"`
	TotalDamage            int       `json:"total_damage"`
	DamageToChampions      int       `json:"damage_to_champions"`
	DamageShare            float64   `json:"damage_share"`
	DamagePerMinute        float64   `json:"damage_per_minute"`
	PhysicalDamage         int       `json:"physical_damage"`
	MagicDamage            int       `json:"magic_damage"`
	TrueDamage             int       `json:"true_damage"`
	GameDuration           int       `json:"game_duration"`
	Result                 string    `json:"result"`
	Date                   time.Time `json:"date"`
	TeamContribution       float64   `json:"team_contribution_score"`
	ItemBuild              []int     `json:"item_build"`
}

// DamageProfile represents a player's damage dealing profile
type DamageProfile struct {
	PlayerID               string                 `json:"player_id"`
	ProfileType            string                 `json:"profile_type"` // "burst", "sustained", "mixed", "utility"
	DamagePattern          DamagePattern          `json:"damage_pattern"`
	OptimalConditions      []string               `json:"optimal_conditions"`
	WeakConditions         []string               `json:"weak_conditions"`
	RecommendedItems       []ItemRecommendation   `json:"recommended_items"`
	PlaystyleAdvice        []string               `json:"playstyle_advice"`
}

// DamagePattern represents how a player typically deals damage
type DamagePattern struct {
	BurstPotential         float64   `json:"burst_potential"`      // 0-100
	SustainedDamage        float64   `json:"sustained_damage"`     // 0-100
	EarlyGameImpact        float64   `json:"early_game_impact"`    // 0-100
	LateGameScaling        float64   `json:"late_game_scaling"`    // 0-100
	TeamfightEffectiveness float64   `json:"teamfight_effectiveness"` // 0-100
	SkirmishPotential      float64   `json:"skirmish_potential"`   // 0-100
}

// ItemRecommendation represents recommended items for damage optimization
type ItemRecommendation struct {
	ItemID                 int       `json:"item_id"`
	ItemName               string    `json:"item_name"`
	ReasonCode             string    `json:"reason_code"` // "damage_increase", "survivability", "utility"
	ExpectedImpact         string    `json:"expected_impact"`
	GamePhase              string    `json:"game_phase"` // "early", "mid", "late", "situational"
	Priority               int       `json:"priority"` // 1-5
}

// NewDamageAnalyticsService creates a new damage analytics service
func NewDamageAnalyticsService(analyticsService *AnalyticsService, matchService *MatchService) *DamageAnalyticsService {
	return &DamageAnalyticsService{
		analyticsService: analyticsService,
		matchService:     matchService,
	}
}

// AnalyzeDamage performs comprehensive damage analysis
func (das *DamageAnalyticsService) AnalyzeDamage(ctx context.Context, playerID string, timeRange string, champion string, position string) (*DamageAnalysis, error) {
	// Get match data with damage information
	matches, err := das.getDamageMatches(ctx, playerID, timeRange, champion, position)
	if err != nil {
		return nil, fmt.Errorf("failed to get damage match data: %w", err)
	}

	if len(matches) == 0 {
		return &DamageAnalysis{
			PlayerID:  playerID,
			Champion:  champion,
			Position:  position,
			TimeRange: timeRange,
		}, nil
	}

	analysis := &DamageAnalysis{
		PlayerID:  playerID,
		Champion:  champion,
		Position:  position,
		TimeRange: timeRange,
	}

	// Calculate basic damage statistics
	das.calculateDamageBasics(analysis, matches)

	// Analyze damage by game phases
	das.analyzeDamagePhases(analysis, matches)

	// Calculate team contribution metrics
	das.calculateTeamContribution(analysis, matches)

	// Perform comparative analysis
	err = das.performDamageComparison(ctx, analysis, position, champion)
	if err != nil {
		fmt.Printf("Warning: failed to perform damage comparison: %v", err)
	}

	// Calculate carry potential and consistency
	das.calculateCarryMetrics(analysis, matches)

	// Analyze situational performance
	das.analyzeSituationalDamage(analysis, matches)

	// Perform trend analysis
	das.analyzeDamageTrend(analysis, matches)

	// Generate insights and recommendations
	das.generateDamageRecommendations(analysis)

	// Prepare trend data for visualization
	das.generateDamageTrendData(analysis, matches)

	// Cache results
	das.cacheDamageAnalysis(ctx, analysis)

	return analysis, nil
}

// AnalyzeDamageProfile creates a comprehensive damage profile for a player
func (das *DamageAnalyticsService) AnalyzeDamageProfile(ctx context.Context, playerID string, timeRange string) (*DamageProfile, error) {
	matches, err := das.getDamageMatches(ctx, playerID, timeRange, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get match data for damage profile: %w", err)
	}

	if len(matches) == 0 {
		return &DamageProfile{
			PlayerID:    playerID,
			ProfileType: "insufficient_data",
		}, nil
	}

	profile := &DamageProfile{
		PlayerID: playerID,
	}

	// Analyze damage patterns
	das.analyzeDamagePatterns(profile, matches)

	// Determine profile type
	das.classifyDamageProfile(profile)

	// Generate recommendations
	das.generateProfileRecommendations(profile)

	return profile, nil
}

// GetTeamContribution calculates comprehensive team contribution metrics
func (das *DamageAnalyticsService) GetTeamContribution(ctx context.Context, playerID string, timeRange string) (*TeamContributionData, error) {
	matches, err := das.getDamageMatches(ctx, playerID, timeRange, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get match data: %w", err)
	}

	contribution := &TeamContributionData{}
	das.calculateTeamContributionDetails(contribution, matches)
	
	return contribution, nil
}

// Helper functions for damage analysis

func (das *DamageAnalyticsService) calculateDamageBasics(analysis *DamageAnalysis, matches []models.MatchData) {
	if len(matches) == 0 {
		return
	}

	totalDamage := 0.0
	totalTime := 0.0
	damageShares := make([]float64, 0, len(matches))
	damageValues := make([]float64, 0, len(matches))

	physicalDamageTotal := 0.0
	magicDamageTotal := 0.0
	trueDamageTotal := 0.0
	totalGold := 0.0

	for _, match := range matches {
		damage := float64(match.TotalDamageToChampions)
		gameTime := float64(match.GameDuration) / 60.0 // Convert to minutes
		
		totalDamage += damage
		totalTime += gameTime
		damageShares = append(damageShares, match.DamageShare)
		damageValues = append(damageValues, damage)

		physicalDamageTotal += float64(match.PhysicalDamageDealt)
		magicDamageTotal += float64(match.MagicDamageDealt)
		trueDamageTotal += float64(match.TrueDamageDealt)
		totalGold += float64(match.GoldEarned)
	}

	// Calculate averages and metrics
	analysis.TotalDamageDealt = totalDamage
	analysis.AverageDamagePerGame = totalDamage / float64(len(matches))
	if totalTime > 0 {
		analysis.DamagePerMinute = totalDamage / totalTime
	}
	analysis.DamageShare = das.calculateMean(damageShares)

	// Damage type breakdown
	totalDamageAllTypes := physicalDamageTotal + magicDamageTotal + trueDamageTotal
	if totalDamageAllTypes > 0 {
		analysis.PhysicalDamageShare = (physicalDamageTotal / totalDamageAllTypes) * 100
		analysis.MagicDamageShare = (magicDamageTotal / totalDamageAllTypes) * 100
		analysis.TrueDamageShare = (trueDamageTotal / totalDamageAllTypes) * 100
	}

	// Efficiency metrics
	if totalGold > 0 {
		analysis.DamagePerGold = totalDamage / totalGold
	}

	// Calculate consistency
	stdDev := das.calculateStandardDeviation(damageValues)
	if analysis.AverageDamagePerGame > 0 {
		analysis.DamageConsistency = 100 - ((stdDev / analysis.AverageDamagePerGame) * 100)
		
		// Classify reliability
		coefficientOfVariation := stdDev / analysis.AverageDamagePerGame
		if coefficientOfVariation < 0.2 {
			analysis.DamageReliability = "consistent"
		} else if coefficientOfVariation < 0.4 {
			analysis.DamageReliability = "variable"
		} else {
			analysis.DamageReliability = "inconsistent"
		}
	}
}

func (das *DamageAnalyticsService) analyzeDamagePhases(analysis *DamageAnalysis, matches []models.MatchData) {
	// This would analyze damage by game phases using timeline data
	// For now, we'll use simplified calculations based on game duration patterns

	earlyDamage := make([]float64, 0)
	midDamage := make([]float64, 0)
	lateDamage := make([]float64, 0)

	for _, match := range matches {
		damage := float64(match.TotalDamageToChampions)
		duration := match.GameDuration

		// Classify games by duration and analyze damage patterns
		if duration <= 1200 { // Games <= 20 minutes (early game focused)
			earlyDamage = append(earlyDamage, damage)
		} else if duration <= 2100 { // Games 20-35 minutes (mid game)
			midDamage = append(midDamage, damage)
		} else { // Games > 35 minutes (late game)
			lateDamage = append(lateDamage, damage)
		}
	}

	// Calculate phase-specific metrics
	analysis.EarlyGameDamage = GamePhaseDamage{
		Phase:               "early",
		AverageDamage:       das.calculateMean(earlyDamage),
		ConsistencyRating:   das.rateConsistency(das.calculateStandardDeviation(earlyDamage), das.calculateMean(earlyDamage)),
		RelativePerformance: das.calculateRelativePerformance(das.calculateMean(earlyDamage), analysis.AverageDamagePerGame),
	}

	analysis.MidGameDamage = GamePhaseDamage{
		Phase:               "mid",
		AverageDamage:       das.calculateMean(midDamage),
		ConsistencyRating:   das.rateConsistency(das.calculateStandardDeviation(midDamage), das.calculateMean(midDamage)),
		RelativePerformance: das.calculateRelativePerformance(das.calculateMean(midDamage), analysis.AverageDamagePerGame),
	}

	analysis.LateGameDamage = GamePhaseDamage{
		Phase:               "late",
		AverageDamage:       das.calculateMean(lateDamage),
		ConsistencyRating:   das.rateConsistency(das.calculateStandardDeviation(lateDamage), das.calculateMean(lateDamage)),
		RelativePerformance: das.calculateRelativePerformance(das.calculateMean(lateDamage), analysis.AverageDamagePerGame),
	}
}

func (das *DamageAnalyticsService) calculateTeamContribution(analysis *DamageAnalysis, matches []models.MatchData) {
	if len(matches) == 0 {
		analysis.TeamContribution = TeamContributionData{}
		return
	}

	contribution := &TeamContributionData{}
	
	// Calculate damage share metrics
	damageShares := make([]float64, 0, len(matches))
	killParticipations := make([]float64, 0, len(matches))
	carryGames := 0

	for _, match := range matches {
		damageShares = append(damageShares, match.DamageShare)
		killParticipations = append(killParticipations, match.TeamfightParticipation)
		
		// Consider it a "carry" game if damage share > 35% or top damage dealer
		if match.DamageShare > 35.0 {
			carryGames++
		}
	}

	contribution.TeamDamageShare = das.calculateMean(damageShares)
	contribution.KillParticipation = das.calculateMean(killParticipations)
	contribution.DamageCarryRate = (float64(carryGames) / float64(len(matches))) * 100

	// Calculate overall contribution score (0-100)
	scoreComponents := []float64{
		math.Min(contribution.TeamDamageShare * 2, 50),  // Max 50 points from damage share
		math.Min(contribution.KillParticipation, 30),     // Max 30 points from kill participation
		math.Min(contribution.DamageCarryRate * 0.2, 20), // Max 20 points from carry rate
	}

	contribution.DamageContributionScore = das.calculateSum(scoreComponents)
	contribution.OverallTeamImpact = contribution.DamageContributionScore

	// Rate team impact
	if contribution.OverallTeamImpact >= 80 {
		contribution.TeamImpactRanking = "excellent"
	} else if contribution.OverallTeamImpact >= 65 {
		contribution.TeamImpactRanking = "good"
	} else if contribution.OverallTeamImpact >= 45 {
		contribution.TeamImpactRanking = "average"
	} else {
		contribution.TeamImpactRanking = "poor"
	}

	analysis.TeamContribution = *contribution
}

func (das *DamageAnalyticsService) calculateCarryMetrics(analysis *DamageAnalysis, matches []models.MatchData) {
	if len(matches) == 0 {
		return
	}

	// Calculate carry potential based on multiple factors
	carryScore := 0.0

	// Factor 1: Average damage share (0-40 points)
	damageShareScore := math.Min(analysis.DamageShare * 1.2, 40)
	carryScore += damageShareScore

	// Factor 2: Consistency (0-25 points)
	consistencyScore := math.Min(analysis.DamageConsistency * 0.25, 25)
	carryScore += consistencyScore

	// Factor 3: Damage per minute relative to role (0-35 points)
	// This would typically compare to role benchmarks
	dpmScore := math.Min((analysis.DamagePerMinute / 1000) * 35, 35)
	carryScore += dpmScore

	analysis.CarryPotential = math.Min(carryScore, 100)
}

func (das *DamageAnalyticsService) analyzeSituationalDamage(analysis *DamageAnalysis, matches []models.MatchData) {
	winningGames := make([]float64, 0)
	losingGames := make([]float64, 0)
	closeGames := make([]float64, 0)

	winningShares := make([]float64, 0)
	losingShares := make([]float64, 0)
	closeShares := make([]float64, 0)

	for _, match := range matches {
		damage := float64(match.TotalDamageToChampions)
		share := match.DamageShare

		// Classify game type (simplified)
		if match.Win {
			winningGames = append(winningGames, damage)
			winningShares = append(winningShares, share)
		} else {
			losingGames = append(losingGames, damage)
			losingShares = append(losingShares, share)
		}

		// Close games (duration-based heuristic)
		if match.GameDuration >= 1800 { // Games >= 30 minutes
			closeGames = append(closeGames, damage)
			closeShares = append(closeShares, share)
		}
	}

	// Calculate situational snapshots
	analysis.WinningGamesDamage = DamageSnapshot{
		Context:           "winning",
		GameCount:        len(winningGames),
		AverageDamage:    das.calculateMean(winningGames),
		DamageShare:      das.calculateMean(winningShares),
		PerformanceRating: das.ratePerformance(das.calculateMean(winningGames), analysis.AverageDamagePerGame),
		ImpactDifference:  das.calculateMean(winningGames) - analysis.AverageDamagePerGame,
	}

	analysis.LosingGamesDamage = DamageSnapshot{
		Context:           "losing",
		GameCount:        len(losingGames),
		AverageDamage:    das.calculateMean(losingGames),
		DamageShare:      das.calculateMean(losingShares),
		PerformanceRating: das.ratePerformance(das.calculateMean(losingGames), analysis.AverageDamagePerGame),
		ImpactDifference:  das.calculateMean(losingGames) - analysis.AverageDamagePerGame,
	}

	analysis.CloseGamesDamage = DamageSnapshot{
		Context:           "close",
		GameCount:        len(closeGames),
		AverageDamage:    das.calculateMean(closeGames),
		DamageShare:      das.calculateMean(closeShares),
		PerformanceRating: das.ratePerformance(das.calculateMean(closeGames), analysis.AverageDamagePerGame),
		ImpactDifference:  das.calculateMean(closeGames) - analysis.AverageDamagePerGame,
	}
}

func (das *DamageAnalyticsService) generateDamageRecommendations(analysis *DamageAnalysis) {
	recommendations := make([]DamageRecommendation, 0)

	// Damage share analysis
	if analysis.DamageShare < 20 {
		recommendations = append(recommendations, DamageRecommendation{
			Priority:            "high",
			Category:           "positioning",
			Title:              "Increase Damage Output",
			Description:        fmt.Sprintf("Your damage share (%.1f%%) is below role expectations. Focus on safer positioning to deal more consistent damage.", analysis.DamageShare),
			ExpectedImprovement: "5-10% increase in damage share can improve win rate by 15-20%",
			GamePhase:          []string{"mid", "late"},
			PositioningTips:    []string{"Stay behind frontline", "Use max range for abilities", "Position for teamfight objectives"},
		})
	}

	// Consistency analysis
	if analysis.DamageReliability == "inconsistent" {
		recommendations = append(recommendations, DamageRecommendation{
			Priority:            "medium",
			Category:           "consistency",
			Title:              "Improve Damage Consistency",
			Description:        "Your damage output varies significantly between games. Work on consistent farming and positioning.",
			ExpectedImprovement: "Better consistency leads to more predictable team performance",
			GamePhase:          []string{"early", "mid", "late"},
		})
	}

	// Carry potential analysis
	if analysis.CarryPotential < 50 {
		recommendations = append(recommendations, DamageRecommendation{
			Priority:            "high",
			Category:           "itemization",
			Title:              "Optimize Item Builds for Damage",
			Description:        fmt.Sprintf("Your carry potential (%.1f) suggests room for itemization improvement.", analysis.CarryPotential),
			ExpectedImprovement: "Optimized builds can increase damage output by 20-30%",
			GamePhase:          []string{"mid", "late"},
			ItemSuggestions:    []string{"Focus on core damage items first", "Consider situational damage items", "Avoid excessive defensive items early"},
		})
	}

	// Phase-specific recommendations
	if analysis.LateGameDamage.RelativePerformance < 0.8 {
		recommendations = append(recommendations, DamageRecommendation{
			Priority:            "medium",
			Category:           "scaling",
			Title:              "Improve Late Game Scaling",
			Description:        "Your late game damage falls off compared to your overall performance. Focus on scaling builds and positioning.",
			ExpectedImprovement: "Better scaling increases late game carry potential",
			GamePhase:          []string{"late"},
		})
	}

	analysis.DamageRecommendations = recommendations

	// Identify strength and improvement areas
	analysis.StrengthAreas = []string{}
	analysis.ImprovementAreas = []string{}

	if analysis.DamageShare > 30 {
		analysis.StrengthAreas = append(analysis.StrengthAreas, "High Damage Output")
	} else if analysis.DamageShare < 20 {
		analysis.ImprovementAreas = append(analysis.ImprovementAreas, "Damage Output")
	}

	if analysis.DamageReliability == "consistent" {
		analysis.StrengthAreas = append(analysis.StrengthAreas, "Consistent Performance")
	} else if analysis.DamageReliability == "inconsistent" {
		analysis.ImprovementAreas = append(analysis.ImprovementAreas, "Performance Consistency")
	}

	if analysis.TeamContribution.OverallTeamImpact > 70 {
		analysis.StrengthAreas = append(analysis.StrengthAreas, "Team Contribution")
	} else if analysis.TeamContribution.OverallTeamImpact < 45 {
		analysis.ImprovementAreas = append(analysis.ImprovementAreas, "Team Impact")
	}
}

// Utility functions

func (das *DamageAnalyticsService) calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (das *DamageAnalyticsService) calculateSum(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum
}

func (das *DamageAnalyticsService) calculateStandardDeviation(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	mean := das.calculateMean(values)
	sumSquaredDiffs := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiffs += diff * diff
	}
	variance := sumSquaredDiffs / float64(len(values))
	return math.Sqrt(variance)
}

func (das *DamageAnalyticsService) rateConsistency(stdDev, mean float64) string {
	if mean == 0 {
		return "insufficient_data"
	}
	cv := stdDev / mean
	if cv < 0.2 {
		return "excellent"
	} else if cv < 0.3 {
		return "good"
	} else if cv < 0.4 {
		return "average"
	} else {
		return "poor"
	}
}

func (das *DamageAnalyticsService) calculateRelativePerformance(phaseAvg, overallAvg float64) float64 {
	if overallAvg == 0 {
		return 1.0
	}
	return phaseAvg / overallAvg
}

func (das *DamageAnalyticsService) ratePerformance(value, baseline float64) string {
	if baseline == 0 {
		return "unknown"
	}
	ratio := value / baseline
	if ratio >= 1.2 {
		return "excellent"
	} else if ratio >= 1.05 {
		return "good"
	} else if ratio >= 0.95 {
		return "average"
	} else if ratio >= 0.8 {
		return "below_average"
	} else {
		return "poor"
	}
}

// Placeholder functions for data access
func (das *DamageAnalyticsService) getDamageMatches(ctx context.Context, playerID string, timeRange string, champion string, position string) ([]models.MatchData, error) {
	// This would query the database for matches with damage data
	return []models.MatchData{}, nil
}

func (das *DamageAnalyticsService) performDamageComparison(ctx context.Context, analysis *DamageAnalysis, position, champion string) error {
	// This would query benchmark data from database
	analysis.RoleBenchmark = DamageBenchmark{
		Category:           "role",
		Filter:             position,
		AverageDamageShare: 25.0, // Placeholder
		PlayerPercentile:   das.calculatePercentile(analysis.DamageShare, 25.0),
		ComparisonRating:   das.rateComparison(analysis.DamageShare, 25.0),
	}
	return nil
}

func (das *DamageAnalyticsService) calculatePercentile(playerValue, benchmark float64) float64 {
	if benchmark == 0 {
		return 50.0
	}
	ratio := playerValue / benchmark
	if ratio >= 1.5 {
		return 90.0
	} else if ratio >= 1.2 {
		return 75.0
	} else if ratio >= 1.0 {
		return 60.0
	} else if ratio >= 0.8 {
		return 40.0
	} else {
		return 25.0
	}
}

func (das *DamageAnalyticsService) rateComparison(playerValue, benchmark float64) string {
	percentile := das.calculatePercentile(playerValue, benchmark)
	if percentile >= 80 {
		return "excellent"
	} else if percentile >= 60 {
		return "above_average"
	} else if percentile >= 40 {
		return "average"
	} else if percentile >= 25 {
		return "below_average"
	} else {
		return "poor"
	}
}

func (das *DamageAnalyticsService) analyzeDamageTrend(analysis *DamageAnalysis, matches []models.MatchData) {
	if len(matches) < 5 {
		analysis.TrendDirection = "insufficient_data"
		return
	}

	// Sort matches by date and calculate trend
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Date.Before(matches[j].Date)
	})

	damageShares := make([]float64, len(matches))
	for i, match := range matches {
		damageShares[i] = match.DamageShare
	}

	slope, confidence := das.calculateLinearRegression(damageShares)
	analysis.TrendSlope = slope
	analysis.TrendConfidence = confidence

	if slope > 1.0 && confidence > 0.6 {
		analysis.TrendDirection = "improving"
	} else if slope < -1.0 && confidence > 0.6 {
		analysis.TrendDirection = "declining"
	} else {
		analysis.TrendDirection = "stable"
	}
}

func (das *DamageAnalyticsService) calculateLinearRegression(values []float64) (slope, confidence float64) {
	// Implementation similar to other analytics services
	if len(values) < 2 {
		return 0, 0
	}
	
	// Linear regression calculation
	n := float64(len(values))
	sumX, sumY, sumXY, sumXX := 0.0, 0.0, 0.0, 0.0
	
	for i, y := range values {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}
	
	denominator := n*sumXX - sumX*sumX
	if denominator == 0 {
		return 0, 0
	}
	
	slope = (n*sumXY - sumX*sumY) / denominator
	
	// Calculate R-squared
	meanY := sumY / n
	ssTotal, ssRes := 0.0, 0.0
	
	for i, y := range values {
		predicted := slope*float64(i) + (sumY-slope*sumX)/n
		ssTotal += (y - meanY) * (y - meanY)
		ssRes += (y - predicted) * (y - predicted)
	}
	
	confidence = 1 - (ssRes / ssTotal)
	if confidence < 0 {
		confidence = 0
	}
	
	return slope, confidence
}

func (das *DamageAnalyticsService) generateDamageTrendData(analysis *DamageAnalysis, matches []models.MatchData) {
	// Generate trend visualization data
	analysis.TrendData = []DamageTrendPoint{}
	// Implementation would process match data into trend points
}

func (das *DamageAnalyticsService) cacheDamageAnalysis(ctx context.Context, analysis *DamageAnalysis) {
	// Cache analysis results
}

// Additional helper functions for damage profile analysis
func (das *DamageAnalyticsService) analyzeDamagePatterns(profile *DamageProfile, matches []models.MatchData) {
	// Analyze how the player typically deals damage
}

func (das *DamageAnalyticsService) classifyDamageProfile(profile *DamageProfile) {
	// Classify the player's damage profile type
}

func (das *DamageAnalyticsService) generateProfileRecommendations(profile *DamageProfile) {
	// Generate profile-specific recommendations
}

func (das *DamageAnalyticsService) calculateTeamContributionDetails(contribution *TeamContributionData, matches []models.MatchData) {
	// Calculate detailed team contribution metrics
}