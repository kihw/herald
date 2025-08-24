package services

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/herald-lol/herald/backend/internal/models"
)

// GoldAnalyticsService handles gold efficiency and economy analytics
type GoldAnalyticsService struct {
	analyticsService *AnalyticsService
}

// GoldAnalysis represents comprehensive gold efficiency analysis
type GoldAnalysis struct {
	PlayerID  string `json:"player_id"`
	Champion  string `json:"champion,omitempty"`
	Position  string `json:"position,omitempty"`
	TimeRange string `json:"time_range"`

	// Core Gold Metrics
	AverageGoldEarned    float64 `json:"average_gold_earned"`
	AverageGoldPerMinute float64 `json:"average_gold_per_minute"`
	GoldEfficiencyScore  float64 `json:"gold_efficiency_score"` // 0-100 scale
	EconomyRating        string  `json:"economy_rating"`        // "excellent", "good", "average", "poor"

	// Gold Sources Analysis
	GoldSources GoldSourcesData `json:"gold_sources"`

	// Spending Efficiency
	ItemEfficiency   ItemEfficiencyData   `json:"item_efficiency"`
	SpendingPatterns SpendingPatternsData `json:"spending_patterns"`

	// Game Phase Analysis
	EarlyGameGold GoldPhaseData `json:"early_game_gold"`
	MidGameGold   GoldPhaseData `json:"mid_game_gold"`
	LateGameGold  GoldPhaseData `json:"late_game_gold"`

	// Comparative Analysis
	RoleBenchmark   GoldBenchmark `json:"role_benchmark"`
	RankBenchmark   GoldBenchmark `json:"rank_benchmark"`
	GlobalBenchmark GoldBenchmark `json:"global_benchmark"`

	// Performance Impact
	GoldAdvantageWinRate    float64 `json:"gold_advantage_win_rate"`
	GoldDisadvantageWinRate float64 `json:"gold_disadvantage_win_rate"`
	GoldImpactScore         float64 `json:"gold_impact_score"`

	// Trend Analysis
	TrendDirection  string           `json:"trend_direction"`
	TrendSlope      float64          `json:"trend_slope"`
	TrendConfidence float64          `json:"trend_confidence"`
	TrendData       []GoldTrendPoint `json:"trend_data"`

	// Economy Optimization
	IncomeOptimization   IncomeOptimizationData   `json:"income_optimization"`
	SpendingOptimization SpendingOptimizationData `json:"spending_optimization"`

	// Insights and Recommendations
	StrengthAreas    []string             `json:"strength_areas"`
	ImprovementAreas []string             `json:"improvement_areas"`
	Recommendations  []GoldRecommendation `json:"recommendations"`

	// Match Performance
	RecentMatches []MatchGoldData `json:"recent_matches"`
}

// GoldSourcesData represents breakdown of gold income sources
type GoldSourcesData struct {
	FarmingGold      float64 `json:"farming_gold"` // CS + monsters
	FarmingPercent   float64 `json:"farming_percent"`
	KillsGold        float64 `json:"kills_gold"` // Champion kills
	KillsPercent     float64 `json:"kills_percent"`
	AssistsGold      float64 `json:"assists_gold"` // Assist gold
	AssistsPercent   float64 `json:"assists_percent"`
	ObjectiveGold    float64 `json:"objective_gold"` // Dragons, Baron, turrets
	ObjectivePercent float64 `json:"objective_percent"`
	PassiveGold      float64 `json:"passive_gold"` // Base income
	PassivePercent   float64 `json:"passive_percent"`
	ItemsGold        float64 `json:"items_gold"` // Item actives, selling items
	ItemsPercent     float64 `json:"items_percent"`

	// Efficiency Metrics
	CSGoldPerMinute    float64 `json:"cs_gold_per_minute"`
	KillGoldEfficiency float64 `json:"kill_gold_efficiency"` // Gold per kill vs deaths lost
	ObjectiveGoldShare float64 `json:"objective_gold_share"` // Team objective gold %
}

// ItemEfficiencyData represents item purchase and utilization efficiency
type ItemEfficiencyData struct {
	AverageItemsCompleted int     `json:"average_items_completed"`
	ItemCompletionSpeed   float64 `json:"item_completion_speed"` // Items per minute
	GoldSpentOnItems      float64 `json:"gold_spent_on_items"`
	ItemValueEfficiency   float64 `json:"item_value_efficiency"` // Value gained vs gold spent

	// Item Categories
	DamageItemsPercent    float64 `json:"damage_items_percent"`
	DefensiveItemsPercent float64 `json:"defensive_items_percent"`
	UtilityItemsPercent   float64 `json:"utility_items_percent"`

	// Power Spikes
	FirstItemTiming float64 `json:"first_item_timing"` // Minutes to first item
	CoreItemsTiming float64 `json:"core_items_timing"` // Minutes to core build
	SixItemsTiming  float64 `json:"six_items_timing"`  // Minutes to full build

	// Optimization Metrics
	OptimalItemOrder       bool    `json:"optimal_item_order"`
	CounterBuildEfficiency float64 `json:"counter_build_efficiency"` // Adapting to enemy team
	ComponentUtilization   float64 `json:"component_utilization"`    // Using item components effectively
}

// SpendingPatternsData represents gold spending behavior analysis
type SpendingPatternsData struct {
	ControlWardsPercent   float64               `json:"control_wards_percent"` // % gold on control wards
	ConsumablesPercent    float64               `json:"consumables_percent"`   // % gold on potions
	BackTiming            BackTimingData        `json:"back_timing"`
	GoldEfficiencyByPhase []PhaseGoldEfficiency `json:"gold_efficiency_by_phase"`

	// Shopping Behavior
	AverageShoppingTime float64 `json:"average_shopping_time"` // Seconds in shop per back
	OptimalBackTiming   float64 `json:"optimal_back_timing"`   // % of backs at good timings
	EmergencyBacks      int     `json:"emergency_backs"`       // Forced backs due to low HP
}

// BackTimingData represents recall timing analysis
type BackTimingData struct {
	AverageBackTiming float64 `json:"average_back_timing"` // Minutes between backs
	OptimalBacks      int     `json:"optimal_backs"`       // Good timing backs
	SuboptimalBacks   int     `json:"suboptimal_backs"`    // Poor timing backs
	ForcedbBacks      int     `json:"forced_backs"`        // Emergency backs
	GoldPerBack       float64 `json:"gold_per_back"`       // Average gold spent per back
}

// PhaseGoldEfficiency represents gold efficiency by game phase
type PhaseGoldEfficiency struct {
	Phase              string  `json:"phase"` // "early", "mid", "late"
	GoldPerMinute      float64 `json:"gold_per_minute"`
	SpendingEfficiency float64 `json:"spending_efficiency"` // 0-100
	IncomeEfficiency   float64 `json:"income_efficiency"`   // 0-100
	EconomyScore       float64 `json:"economy_score"`       // Combined score
}

// GoldPhaseData represents gold metrics for a game phase
type GoldPhaseData struct {
	Phase                  string  `json:"phase"`
	AverageGoldPerMinute   float64 `json:"average_gold_per_minute"`
	GoldAdvantage          float64 `json:"gold_advantage"`          // vs enemy laner
	FarmingEfficiency      float64 `json:"farming_efficiency"`      // % of available CS gold
	KillParticipation      float64 `json:"kill_participation"`      // % of team's kill gold
	ObjectiveParticipation float64 `json:"objective_participation"` // % of team's objective gold
	SpendingScore          float64 `json:"spending_score"`          // 0-100 efficiency
	EfficiencyRating       string  `json:"efficiency_rating"`
}

// GoldBenchmark represents gold performance benchmarks
type GoldBenchmark struct {
	Category             string  `json:"category"`
	AverageGoldPerMinute float64 `json:"average_gold_per_minute"`
	Top10Percent         float64 `json:"top_10_percent"`
	Top25Percent         float64 `json:"top_25_percent"`
	Median               float64 `json:"median"`
	PlayerPercentile     float64 `json:"player_percentile"`
	EfficiencyAverage    float64 `json:"efficiency_average"`
}

// IncomeOptimizationData represents income improvement opportunities
type IncomeOptimizationData struct {
	CSImprovementPotential        float64 `json:"cs_improvement_potential"`        // Additional GPM from better CS
	KillParticipationPotential    float64 `json:"kp_improvement_potential"`        // Additional GPM from better KP
	ObjectiveImprovementPotential float64 `json:"objective_improvement_potential"` // Additional GPM from objectives

	// Specific Recommendations
	EarlyFarmingSuggestions    []string `json:"early_farming_suggestions"`
	MidGamePositionSuggestions []string `json:"mid_game_position_suggestions"`
	LateGameFocusSuggestions   []string `json:"late_game_focus_suggestions"`

	// Potential Impact
	ExpectedGPMIncrease     float64 `json:"expected_gpm_increase"`
	ExpectedWinRateIncrease float64 `json:"expected_win_rate_increase"`
}

// SpendingOptimizationData represents spending efficiency improvements
type SpendingOptimizationData struct {
	ItemOrderOptimization     []string `json:"item_order_optimization"`
	BackTimingOptimization    []string `json:"back_timing_optimization"`
	GoldAllocationSuggestions []string `json:"gold_allocation_suggestions"`

	// Component Management
	ComponentBuyingTips []string `json:"component_buying_tips"`
	PowerSpikeTiming    []string `json:"power_spike_timing"`

	// Economic Priorities
	EarlyGamePriorities []string `json:"early_game_priorities"`
	MidGamePriorities   []string `json:"mid_game_priorities"`
	LateGamePriorities  []string `json:"late_game_priorities"`
}

// GoldTrendPoint represents gold performance over time
type GoldTrendPoint struct {
	Date               time.Time `json:"date"`
	GoldPerMinute      float64   `json:"gold_per_minute"`
	GoldEfficiency     float64   `json:"gold_efficiency"`
	FarmingEfficiency  float64   `json:"farming_efficiency"`
	SpendingEfficiency float64   `json:"spending_efficiency"`
	MovingAverage      float64   `json:"moving_average"`
}

// GoldRecommendation represents actionable gold efficiency advice
type GoldRecommendation struct {
	Priority                 string   `json:"priority"` // "high", "medium", "low"
	Category                 string   `json:"category"` // "farming", "spending", "timing", "itemization"
	Title                    string   `json:"title"`
	Description              string   `json:"description"`
	Impact                   string   `json:"impact"`
	GamePhase                []string `json:"game_phase"`
	ExpectedGPMIncrease      float64  `json:"expected_gpm_increase"`
	ImplementationDifficulty string   `json:"implementation_difficulty"` // "easy", "medium", "hard"
}

// MatchGoldData represents gold performance in a specific match
type MatchGoldData struct {
	MatchID             string    `json:"match_id"`
	Champion            string    `json:"champion"`
	Position            string    `json:"position"`
	TotalGoldEarned     int       `json:"total_gold_earned"`
	GoldPerMinute       float64   `json:"gold_per_minute"`
	GoldEfficiencyScore float64   `json:"gold_efficiency_score"`
	FarmingGold         int       `json:"farming_gold"`
	KillGold            int       `json:"kill_gold"`
	ObjectiveGold       int       `json:"objective_gold"`
	ItemsCompleted      int       `json:"items_completed"`
	ControlWardsSpent   int       `json:"control_wards_spent"`
	GameDuration        int       `json:"game_duration"`
	Result              string    `json:"result"`
	Date                time.Time `json:"date"`
	GoldAdvantageAt15   int       `json:"gold_advantage_at_15"`
	TeamGoldShare       float64   `json:"team_gold_share"`
}

// NewGoldAnalyticsService creates a new gold analytics service
func NewGoldAnalyticsService(analyticsService *AnalyticsService) *GoldAnalyticsService {
	return &GoldAnalyticsService{
		analyticsService: analyticsService,
	}
}

// AnalyzeGold performs comprehensive gold efficiency analysis
func (gas *GoldAnalyticsService) AnalyzeGold(ctx context.Context, playerID string, timeRange string, champion string, position string) (*GoldAnalysis, error) {
	// Get match data with gold information
	matches, err := gas.getGoldMatches(ctx, playerID, timeRange, champion, position)
	if err != nil {
		return nil, fmt.Errorf("failed to get gold match data: %w", err)
	}

	if len(matches) == 0 {
		return &GoldAnalysis{
			PlayerID:  playerID,
			Champion:  champion,
			Position:  position,
			TimeRange: timeRange,
		}, nil
	}

	analysis := &GoldAnalysis{
		PlayerID:  playerID,
		Champion:  champion,
		Position:  position,
		TimeRange: timeRange,
	}

	// Calculate core gold metrics
	gas.calculateGoldBasics(analysis, matches)

	// Analyze gold sources
	gas.analyzeGoldSources(analysis, matches)

	// Analyze item efficiency
	gas.analyzeItemEfficiency(analysis, matches)

	// Analyze spending patterns
	gas.analyzeSpendingPatterns(analysis, matches)

	// Analyze by game phases
	gas.analyzeGoldPhases(analysis, matches)

	// Perform comparative analysis
	err = gas.performGoldComparison(ctx, analysis, position)
	if err != nil {
		fmt.Printf("Warning: failed to perform gold comparison: %v", err)
	}

	// Calculate gold impact score
	gas.calculateGoldImpactScore(analysis)

	// Perform trend analysis
	gas.analyzeGoldTrend(analysis, matches)

	// Generate optimization suggestions
	gas.generateOptimizationSuggestions(analysis)

	// Generate recommendations
	gas.generateGoldRecommendations(analysis)

	// Prepare trend data for visualization
	gas.generateGoldTrendData(analysis, matches)

	return analysis, nil
}

// calculateGoldBasics calculates fundamental gold metrics
func (gas *GoldAnalyticsService) calculateGoldBasics(analysis *GoldAnalysis, matches []models.MatchData) {
	if len(matches) == 0 {
		return
	}

	totalGold := 0.0
	totalMinutes := 0.0
	efficiencyScores := make([]float64, 0, len(matches))

	for _, match := range matches {
		goldEarned := float64(match.GoldEarned)
		minutes := float64(match.GameDuration) / 60.0

		totalGold += goldEarned
		totalMinutes += minutes

		// Calculate efficiency score for this match (0-100)
		gpm := goldEarned / minutes
		efficiencyScore := gas.calculateMatchGoldEfficiency(match)
		efficiencyScores = append(efficiencyScores, efficiencyScore)
	}

	analysis.AverageGoldEarned = totalGold / float64(len(matches))
	analysis.AverageGoldPerMinute = totalGold / totalMinutes
	analysis.GoldEfficiencyScore = gas.calculateMean(efficiencyScores)
	analysis.EconomyRating = gas.rateEconomy(analysis.GoldEfficiencyScore)
}

// analyzeGoldSources breaks down gold income sources
func (gas *GoldAnalyticsService) analyzeGoldSources(analysis *GoldAnalysis, matches []models.MatchData) {
	totalGold := 0.0
	totalFarming := 0.0
	totalKills := 0.0
	totalAssists := 0.0
	totalObjectives := 0.0
	totalPassive := 0.0

	for _, match := range matches {
		goldEarned := float64(match.GoldEarned)
		totalGold += goldEarned

		// Estimate gold sources based on match data
		farmingGold := float64(match.CS)*20.0 + float64(match.NeutralMinionsKilled)*30.0 // Rough estimates
		killsGold := float64(match.Kills) * 300.0                                        // Average kill gold
		assistsGold := float64(match.Assists) * 150.0                                    // Average assist gold

		totalFarming += farmingGold
		totalKills += killsGold
		totalAssists += assistsGold

		// Estimate other sources
		objectiveGold := goldEarned * 0.15               // Rough estimate for objectives
		passiveGold := float64(match.GameDuration) * 1.5 // Base gold per second

		totalObjectives += objectiveGold
		totalPassive += passiveGold
	}

	analysis.GoldSources = GoldSourcesData{
		FarmingGold:      totalFarming / float64(len(matches)),
		FarmingPercent:   (totalFarming / totalGold) * 100,
		KillsGold:        totalKills / float64(len(matches)),
		KillsPercent:     (totalKills / totalGold) * 100,
		AssistsGold:      totalAssists / float64(len(matches)),
		AssistsPercent:   (totalAssists / totalGold) * 100,
		ObjectiveGold:    totalObjectives / float64(len(matches)),
		ObjectivePercent: (totalObjectives / totalGold) * 100,
		PassiveGold:      totalPassive / float64(len(matches)),
		PassivePercent:   (totalPassive / totalGold) * 100,

		// Calculate efficiency metrics
		CSGoldPerMinute:    (totalFarming / totalGold) * analysis.AverageGoldPerMinute,
		KillGoldEfficiency: gas.calculateKillGoldEfficiency(matches),
		ObjectiveGoldShare: 25.0, // Placeholder - would need team data
	}
}

// analyzeItemEfficiency analyzes item purchase and utilization
func (gas *GoldAnalyticsService) analyzeItemEfficiency(analysis *GoldAnalysis, matches []models.MatchData) {
	totalItems := 0.0
	totalGameTime := 0.0
	firstItemTimings := make([]float64, 0)

	for _, match := range matches {
		totalItems += float64(match.ItemsPurchased)
		gameMinutes := float64(match.GameDuration) / 60.0
		totalGameTime += gameMinutes

		// Estimate first item timing (rough calculation)
		if match.GoldEarned > 3000 { // Assume first item around 3000+ gold
			firstItemTimings = append(firstItemTimings, gameMinutes*0.2) // Rough estimate
		}
	}

	analysis.ItemEfficiency = ItemEfficiencyData{
		AverageItemsCompleted: int(totalItems / float64(len(matches))),
		ItemCompletionSpeed:   totalItems / totalGameTime,
		GoldSpentOnItems:      analysis.AverageGoldEarned * 0.85, // Most gold goes to items
		ItemValueEfficiency:   75.0,                              // Placeholder efficiency score

		// Rough estimates for item categories
		DamageItemsPercent:    60.0,
		DefensiveItemsPercent: 25.0,
		UtilityItemsPercent:   15.0,

		// Timing estimates
		FirstItemTiming: gas.calculateMean(firstItemTimings),
		CoreItemsTiming: 20.0, // Minutes to core items
		SixItemsTiming:  35.0, // Minutes to full build

		// Optimization metrics
		OptimalItemOrder:       true,
		CounterBuildEfficiency: 70.0,
		ComponentUtilization:   80.0,
	}
}

// analyzeSpendingPatterns analyzes gold spending behavior
func (gas *GoldAnalyticsService) analyzeSpendingPatterns(analysis *GoldAnalysis, matches []models.MatchData) {
	analysis.SpendingPatterns = SpendingPatternsData{
		ControlWardsPercent: 5.0, // Typical % spent on control wards
		ConsumablesPercent:  3.0, // Typical % spent on consumables

		BackTiming: BackTimingData{
			AverageBackTiming: 8.0, // Minutes between backs
			OptimalBacks:      6,
			SuboptimalBacks:   3,
			ForcedbBacks:      2,
			GoldPerBack:       1200.0,
		},

		GoldEfficiencyByPhase: []PhaseGoldEfficiency{
			{
				Phase:              "early",
				GoldPerMinute:      analysis.AverageGoldPerMinute * 0.7,
				SpendingEfficiency: 80.0,
				IncomeEfficiency:   85.0,
				EconomyScore:       82.5,
			},
			{
				Phase:              "mid",
				GoldPerMinute:      analysis.AverageGoldPerMinute * 1.2,
				SpendingEfficiency: 85.0,
				IncomeEfficiency:   80.0,
				EconomyScore:       82.5,
			},
			{
				Phase:              "late",
				GoldPerMinute:      analysis.AverageGoldPerMinute * 0.9,
				SpendingEfficiency: 90.0,
				IncomeEfficiency:   70.0,
				EconomyScore:       80.0,
			},
		},

		AverageShoppingTime: 15.0,
		OptimalBackTiming:   70.0,
		EmergencyBacks:      2,
	}
}

// analyzeGoldPhases analyzes gold performance by game phases
func (gas *GoldAnalyticsService) analyzeGoldPhases(analysis *GoldAnalysis, matches []models.MatchData) {
	analysis.EarlyGameGold = GoldPhaseData{
		Phase:                  "early",
		AverageGoldPerMinute:   analysis.AverageGoldPerMinute * 0.7,
		GoldAdvantage:          50.0, // Average gold advantage
		FarmingEfficiency:      85.0, // % of available CS
		KillParticipation:      15.0, // % of team kills
		ObjectiveParticipation: 20.0, // % of objectives
		SpendingScore:          80.0,
		EfficiencyRating:       gas.rateEfficiency(80.0),
	}

	analysis.MidGameGold = GoldPhaseData{
		Phase:                  "mid",
		AverageGoldPerMinute:   analysis.AverageGoldPerMinute * 1.2,
		GoldAdvantage:          200.0,
		FarmingEfficiency:      75.0,
		KillParticipation:      25.0,
		ObjectiveParticipation: 30.0,
		SpendingScore:          85.0,
		EfficiencyRating:       gas.rateEfficiency(85.0),
	}

	analysis.LateGameGold = GoldPhaseData{
		Phase:                  "late",
		AverageGoldPerMinute:   analysis.AverageGoldPerMinute * 0.9,
		GoldAdvantage:          -100.0, // Often behind in late game
		FarmingEfficiency:      60.0,
		KillParticipation:      30.0,
		ObjectiveParticipation: 40.0,
		SpendingScore:          90.0,
		EfficiencyRating:       gas.rateEfficiency(90.0),
	}
}

// performGoldComparison compares against benchmarks
func (gas *GoldAnalyticsService) performGoldComparison(ctx context.Context, analysis *GoldAnalysis, position string) error {
	// Get role benchmark
	roleBenchmark, err := gas.getGoldBenchmark(ctx, "role", position, "")
	if err == nil {
		analysis.RoleBenchmark = *roleBenchmark
		analysis.RoleBenchmark.PlayerPercentile = gas.calculatePercentile(analysis.AverageGoldPerMinute, roleBenchmark.AverageGoldPerMinute)
	}

	// Get rank benchmark
	rankBenchmark, err := gas.getGoldBenchmark(ctx, "rank", "GOLD", "")
	if err == nil {
		analysis.RankBenchmark = *rankBenchmark
		analysis.RankBenchmark.PlayerPercentile = gas.calculatePercentile(analysis.AverageGoldPerMinute, rankBenchmark.AverageGoldPerMinute)
	}

	// Get global benchmark
	globalBenchmark, err := gas.getGoldBenchmark(ctx, "global", "", "")
	if err == nil {
		analysis.GlobalBenchmark = *globalBenchmark
		analysis.GlobalBenchmark.PlayerPercentile = gas.calculatePercentile(analysis.AverageGoldPerMinute, globalBenchmark.AverageGoldPerMinute)
	}

	return nil
}

// calculateGoldImpactScore calculates overall gold impact
func (gas *GoldAnalyticsService) calculateGoldImpactScore(analysis *GoldAnalysis) {
	baseScore := 50.0

	// Factor 1: Gold per minute vs role average
	if analysis.RoleBenchmark.AverageGoldPerMinute > 0 {
		gpmRatio := analysis.AverageGoldPerMinute / analysis.RoleBenchmark.AverageGoldPerMinute
		baseScore += (gpmRatio - 1.0) * 30.0
	}

	// Factor 2: Efficiency score
	efficiencyBonus := (analysis.GoldEfficiencyScore - 50.0) * 0.5
	baseScore += efficiencyBonus

	// Factor 3: Gold sources diversity
	sourcesBalance := gas.calculateSourcesBalance(analysis.GoldSources)
	baseScore += sourcesBalance * 10.0

	// Clamp between 0-100
	analysis.GoldImpactScore = math.Max(0.0, math.Min(100.0, baseScore))
}

// analyzeGoldTrend analyzes gold performance trends
func (gas *GoldAnalyticsService) analyzeGoldTrend(analysis *GoldAnalysis, matches []models.MatchData) {
	if len(matches) < 5 {
		analysis.TrendDirection = "insufficient_data"
		return
	}

	// Sort matches by date
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Date.Before(matches[j].Date)
	})

	// Extract GPM values
	gpmValues := make([]float64, len(matches))
	for i, match := range matches {
		gameMinutes := float64(match.GameDuration) / 60.0
		if gameMinutes > 0 {
			gpmValues[i] = float64(match.GoldEarned) / gameMinutes
		}
	}

	// Calculate trend using linear regression
	slope, confidence := gas.calculateLinearRegression(gpmValues)

	analysis.TrendSlope = slope
	analysis.TrendConfidence = confidence

	// Determine trend direction
	if slope > 5.0 && confidence > 0.6 {
		analysis.TrendDirection = "improving"
	} else if slope < -5.0 && confidence > 0.6 {
		analysis.TrendDirection = "declining"
	} else {
		analysis.TrendDirection = "stable"
	}
}

// generateOptimizationSuggestions creates income and spending optimization advice
func (gas *GoldAnalyticsService) generateOptimizationSuggestions(analysis *GoldAnalysis) {
	// Income optimization
	csImprovement := 50.0  // Potential additional GPM from better CS
	kpImprovement := 30.0  // Potential additional GPM from better kill participation
	objImprovement := 25.0 // Potential additional GPM from better objective control

	analysis.IncomeOptimization = IncomeOptimizationData{
		CSImprovementPotential:        csImprovement,
		KillParticipationPotential:    kpImprovement,
		ObjectiveImprovementPotential: objImprovement,

		EarlyFarmingSuggestions: []string{
			"Focus on last-hitting minions more accurately",
			"Take jungle camps when available",
			"Prioritize farming over trades in early game",
		},

		MidGamePositionSuggestions: []string{
			"Be present for dragon and herald fights",
			"Farm side lanes safely while team groups",
			"Take enemy jungle camps when ahead",
		},

		LateGameFocusSuggestions: []string{
			"Focus on team fights over farming",
			"Secure baron and elder dragon",
			"Farm efficiently between team objectives",
		},

		ExpectedGPMIncrease:     csImprovement + kpImprovement + objImprovement,
		ExpectedWinRateIncrease: 5.0,
	}

	// Spending optimization
	analysis.SpendingOptimization = SpendingOptimizationData{
		ItemOrderOptimization: []string{
			"Complete damage items before defensive ones",
			"Build components that help in lane first",
			"Consider enemy team comp in item choices",
		},

		BackTimingOptimization: []string{
			"Back when you have gold for meaningful items",
			"Don't back with small amounts of gold",
			"Time backs with wave management",
		},

		GoldAllocationSuggestions: []string{
			"Spend 5-8% of gold on control wards",
			"Buy consumables for sustain in lane",
			"Prioritize power spikes for objectives",
		},

		ComponentBuyingTips: []string{
			"Buy the most useful components first",
			"Consider component passive effects",
			"Build components that help current game state",
		},

		PowerSpikeTiming: []string{
			"Plan team fights around item completions",
			"Communicate power spikes to team",
			"Use advantage windows effectively",
		},

		EarlyGamePriorities: []string{"Farming tools", "Sustain items", "Basic components"},
		MidGamePriorities:   []string{"Core items", "Team fight items", "Control wards"},
		LateGamePriorities:  []string{"Situational items", "Defensive options", "Team utility"},
	}
}

// generateGoldRecommendations creates actionable gold efficiency recommendations
func (gas *GoldAnalyticsService) generateGoldRecommendations(analysis *GoldAnalysis) {
	recommendations := []GoldRecommendation{}

	// Check gold per minute vs benchmark
	if analysis.RoleBenchmark.PlayerPercentile < 50 {
		recommendations = append(recommendations, GoldRecommendation{
			Priority:                 "high",
			Category:                 "farming",
			Title:                    "Improve Gold Generation",
			Description:              fmt.Sprintf("Your GPM (%.0f) is below role average (%.0f). Focus on CS and map presence.", analysis.AverageGoldPerMinute, analysis.RoleBenchmark.AverageGoldPerMinute),
			Impact:                   "Increasing GPM by 50 can improve win rate by 8-12%",
			GamePhase:                []string{"early", "mid"},
			ExpectedGPMIncrease:      50.0,
			ImplementationDifficulty: "medium",
		})
	}

	// Check farming efficiency
	if analysis.GoldSources.FarmingPercent < 45 {
		recommendations = append(recommendations, GoldRecommendation{
			Priority:                 "high",
			Category:                 "farming",
			Title:                    "Increase Farming Focus",
			Description:              fmt.Sprintf("Only %.1f%% of your gold comes from farming. Aim for 50-60%% to improve consistency.", analysis.GoldSources.FarmingPercent),
			Impact:                   "Better farming provides reliable gold income",
			GamePhase:                []string{"early", "mid"},
			ExpectedGPMIncrease:      40.0,
			ImplementationDifficulty: "easy",
		})
	}

	// Check gold efficiency score
	if analysis.GoldEfficiencyScore < 60 {
		recommendations = append(recommendations, GoldRecommendation{
			Priority:                 "medium",
			Category:                 "spending",
			Title:                    "Optimize Gold Spending",
			Description:              fmt.Sprintf("Your gold efficiency (%.1f/100) suggests suboptimal spending. Focus on power spike timing.", analysis.GoldEfficiencyScore),
			Impact:                   "Better gold efficiency amplifies your gold advantage",
			GamePhase:                []string{"early", "mid", "late"},
			ExpectedGPMIncrease:      0.0, // Efficiency, not raw GPM
			ImplementationDifficulty: "medium",
		})
	}

	// Check kill participation
	if analysis.GoldSources.KillsPercent+analysis.GoldSources.AssistsPercent < 20 {
		recommendations = append(recommendations, GoldRecommendation{
			Priority:                 "medium",
			Category:                 "teamwork",
			Title:                    "Increase Kill Participation",
			Description:              "Low kill participation limits your gold income. Be more present for team fights and skirmishes.",
			Impact:                   "Higher kill participation provides burst gold income",
			GamePhase:                []string{"mid", "late"},
			ExpectedGPMIncrease:      30.0,
			ImplementationDifficulty: "medium",
		})
	}

	analysis.Recommendations = recommendations

	// Generate strength and improvement areas
	analysis.StrengthAreas = []string{}
	analysis.ImprovementAreas = []string{}

	if analysis.RoleBenchmark.PlayerPercentile > 75 {
		analysis.StrengthAreas = append(analysis.StrengthAreas, "Gold Generation")
	} else if analysis.RoleBenchmark.PlayerPercentile < 25 {
		analysis.ImprovementAreas = append(analysis.ImprovementAreas, "Gold Generation")
	}

	if analysis.GoldEfficiencyScore > 80 {
		analysis.StrengthAreas = append(analysis.StrengthAreas, "Gold Efficiency")
	} else if analysis.GoldEfficiencyScore < 60 {
		analysis.ImprovementAreas = append(analysis.ImprovementAreas, "Gold Efficiency")
	}

	if analysis.GoldSources.FarmingPercent > 55 {
		analysis.StrengthAreas = append(analysis.StrengthAreas, "Farming Focus")
	} else if analysis.GoldSources.FarmingPercent < 40 {
		analysis.ImprovementAreas = append(analysis.ImprovementAreas, "Farming Focus")
	}
}

// generateGoldTrendData creates trend visualization data
func (gas *GoldAnalyticsService) generateGoldTrendData(analysis *GoldAnalysis, matches []models.MatchData) {
	analysis.TrendData = []GoldTrendPoint{}

	// Sort matches by date
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Date.Before(matches[j].Date)
	})

	// Create trend points with moving averages
	windowSize := 5
	for i := range matches {
		gameMinutes := float64(matches[i].GameDuration) / 60.0
		var gpm float64
		if gameMinutes > 0 {
			gpm = float64(matches[i].GoldEarned) / gameMinutes
		}

		point := GoldTrendPoint{
			Date:               matches[i].Date,
			GoldPerMinute:      gpm,
			GoldEfficiency:     gas.calculateMatchGoldEfficiency(matches[i]),
			FarmingEfficiency:  float64(matches[i].CS) / gameMinutes * 10.0, // CS per 10 min
			SpendingEfficiency: 80.0,                                        // Placeholder
		}

		// Calculate moving average
		if i >= windowSize-1 {
			sum := 0.0
			for j := i - windowSize + 1; j <= i; j++ {
				gameMin := float64(matches[j].GameDuration) / 60.0
				if gameMin > 0 {
					sum += float64(matches[j].GoldEarned) / gameMin
				}
			}
			point.MovingAverage = sum / float64(windowSize)
		} else {
			point.MovingAverage = point.GoldPerMinute
		}

		analysis.TrendData = append(analysis.TrendData, point)
	}
}

// Helper functions

func (gas *GoldAnalyticsService) calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (gas *GoldAnalyticsService) calculateMatchGoldEfficiency(match models.MatchData) float64 {
	// Simplified efficiency calculation based on gold sources
	gameMinutes := float64(match.GameDuration) / 60.0
	if gameMinutes == 0 {
		return 50.0
	}

	gpm := float64(match.GoldEarned) / gameMinutes
	csPerMin := float64(match.CS) / gameMinutes

	// Base efficiency on GPM and CS rate
	baseScore := 50.0
	if gpm > 400 {
		baseScore += 20.0
	}
	if csPerMin > 6.0 {
		baseScore += 15.0
	}
	if match.Kills+match.Assists > match.Deaths*2 {
		baseScore += 15.0
	}

	return math.Max(0.0, math.Min(100.0, baseScore))
}

func (gas *GoldAnalyticsService) calculateKillGoldEfficiency(matches []models.MatchData) float64 {
	totalKillGold := 0.0
	totalDeathCost := 0.0

	for _, match := range matches {
		killGold := float64(match.Kills+match.Assists) * 225.0 // Average kill/assist gold
		deathCost := float64(match.Deaths) * 200.0             // Average death cost

		totalKillGold += killGold
		totalDeathCost += deathCost
	}

	if totalDeathCost == 0 {
		return 100.0
	}

	return (totalKillGold / (totalKillGold + totalDeathCost)) * 100.0
}

func (gas *GoldAnalyticsService) calculateSourcesBalance(sources GoldSourcesData) float64 {
	// Calculate how balanced gold sources are (diversity is good)
	percentages := []float64{
		sources.FarmingPercent,
		sources.KillsPercent + sources.AssistsPercent,
		sources.ObjectivePercent,
	}

	// Calculate variance - lower variance means better balance
	mean := (percentages[0] + percentages[1] + percentages[2]) / 3.0
	variance := 0.0
	for _, p := range percentages {
		variance += (p - mean) * (p - mean)
	}
	variance /= float64(len(percentages))

	// Convert to 0-1 scale where lower variance = higher score
	return 1.0 - (variance / 1000.0) // Normalize variance
}

func (gas *GoldAnalyticsService) calculateLinearRegression(values []float64) (slope, confidence float64) {
	if len(values) < 2 {
		return 0, 0
	}

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

func (gas *GoldAnalyticsService) calculatePercentile(playerValue, benchmarkAverage float64) float64 {
	if benchmarkAverage == 0 {
		return 50.0
	}
	ratio := playerValue / benchmarkAverage
	if ratio >= 1.5 {
		return 95.0
	} else if ratio >= 1.3 {
		return 80.0
	} else if ratio >= 1.1 {
		return 65.0
	} else if ratio >= 1.0 {
		return 55.0
	} else if ratio >= 0.9 {
		return 45.0
	} else if ratio >= 0.8 {
		return 30.0
	} else {
		return 15.0
	}
}

func (gas *GoldAnalyticsService) rateEconomy(score float64) string {
	switch {
	case score >= 85:
		return "excellent"
	case score >= 70:
		return "good"
	case score >= 55:
		return "average"
	default:
		return "poor"
	}
}

func (gas *GoldAnalyticsService) rateEfficiency(score float64) string {
	switch {
	case score >= 85:
		return "excellent"
	case score >= 70:
		return "good"
	case score >= 55:
		return "average"
	default:
		return "poor"
	}
}

// Placeholder functions for data access (would be implemented with real database queries)

func (gas *GoldAnalyticsService) getGoldMatches(ctx context.Context, playerID string, timeRange string, champion string, position string) ([]models.MatchData, error) {
	// This would query the database for matches with gold data
	return []models.MatchData{}, nil
}

func (gas *GoldAnalyticsService) getGoldBenchmark(ctx context.Context, category, filter, champion string) (*GoldBenchmark, error) {
	// This would query benchmark data from database
	return &GoldBenchmark{
		Category:             category,
		AverageGoldPerMinute: 450.0, // Placeholder
		Top10Percent:         600.0,
		Top25Percent:         520.0,
		Median:               430.0,
		EfficiencyAverage:    70.0,
	}, nil
}
