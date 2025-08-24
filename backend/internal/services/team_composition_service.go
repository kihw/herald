package services

import (
	"fmt"
	"math"
	"sort"
	"time"

	"gorm.io/gorm"
)

// TeamCompositionService provides advanced team composition optimization and analysis
type TeamCompositionService struct {
	db                     *gorm.DB
	analyticsService       *AnalyticsService
	metaService            *MetaAnalyticsService
	matchPredictionService *MatchPredictionService
}

// NewTeamCompositionService creates a new team composition service
func NewTeamCompositionService(db *gorm.DB, analyticsService *AnalyticsService, metaService *MetaAnalyticsService, matchPredictionService *MatchPredictionService) *TeamCompositionService {
	return &TeamCompositionService{
		db:                     db,
		analyticsService:       analyticsService,
		metaService:            metaService,
		matchPredictionService: matchPredictionService,
	}
}

// TeamComposition represents an optimized team composition
type TeamComposition struct {
	ID              string `json:"id" gorm:"primaryKey"`
	CompositionName string `json:"composition_name"`
	CompositionType string `json:"composition_type"` // team_fight, split_push, poke, etc.
	Tier            string `json:"tier"`             // S+, S, A+, A, B+, B, C

	// Champion Assignments
	TopLaner ChampionAssignment `json:"top_laner" gorm:"embedded;embeddedPrefix:top_"`
	Jungler  ChampionAssignment `json:"jungler" gorm:"embedded;embeddedPrefix:jungle_"`
	MidLaner ChampionAssignment `json:"mid_laner" gorm:"embedded;embeddedPrefix:mid_"`
	ADCarry  ChampionAssignment `json:"ad_carry" gorm:"embedded;embeddedPrefix:adc_"`
	Support  ChampionAssignment `json:"support" gorm:"embedded;embeddedPrefix:support_"`

	// Composition Analysis
	OverallRating     float64 `json:"overall_rating"`     // 0-100
	SynergyScore      float64 `json:"synergy_score"`      // 0-100
	CounterResistance float64 `json:"counter_resistance"` // 0-100
	FlexibilityScore  float64 `json:"flexibility_score"`  // 0-100

	// Performance Metrics
	WinRateData    WinRateAnalysis    `json:"win_rate_data" gorm:"embedded"`
	ScalingProfile CompositionScaling `json:"scaling_profile" gorm:"embedded"`
	PowerSpikes    []PowerSpikeData   `json:"power_spikes" gorm:"type:text"`

	// Strategic Analysis
	WinConditions []CompositionWinCondition `json:"win_conditions" gorm:"type:text"`
	Weaknesses    []CompositionWeakness     `json:"weaknesses" gorm:"type:text"`
	CounterComps  []string                  `json:"counter_comps" gorm:"type:text"`
	SynergyComps  []string                  `json:"synergy_comps" gorm:"type:text"`

	// Meta Context
	MetaRelevance  float64 `json:"meta_relevance"`  // 0-100
	PatchStability float64 `json:"patch_stability"` // 0-100
	TrendDirection string  `json:"trend_direction"` // rising, stable, declining

	// Usage Statistics
	PickRate     float64 `json:"pick_rate"`      // 0-100
	BanRate      float64 `json:"ban_rate"`       // 0-100
	ProPlayUsage float64 `json:"pro_play_usage"` // 0-100

	// Optimization Context
	OptimizedFor     []string `json:"optimized_for" gorm:"type:text"`     // rank, role, playstyle
	RecommendedRanks []string `json:"recommended_ranks" gorm:"type:text"` // bronze, silver, gold, etc.
	DifficultyRating string   `json:"difficulty_rating"`                  // easy, medium, hard, expert

	// Metadata
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastAnalyzed time.Time `json:"last_analyzed"`
}

// ChampionAssignment represents a champion assigned to a specific role
type ChampionAssignment struct {
	ChampionName string `json:"champion_name"`
	ChampionID   int    `json:"champion_id"`
	Role         string `json:"role"`
	Priority     int    `json:"priority"` // 1 = highest priority

	// Champion Analysis
	RoleEfficiency float64 `json:"role_efficiency"` // 0-100
	MetaStrength   float64 `json:"meta_strength"`   // 0-100
	SynergyContrib float64 `json:"synergy_contrib"` // 0-100
	CounterResist  float64 `json:"counter_resist"`  // 0-100

	// Alternative Options
	Alternatives []AlternativeChampion `json:"alternatives" gorm:"type:text"`
	FlexOptions  []FlexOption          `json:"flex_options" gorm:"type:text"`

	// Performance Data
	WinRate        float64 `json:"win_rate"`  // 0-100
	PickRate       float64 `json:"pick_rate"` // 0-100
	BanRate        float64 `json:"ban_rate"`  // 0-100
	AvgKDA         float64 `json:"avg_kda"`
	AvgDamageShare float64 `json:"avg_damage_share"` // 0-100
}

// AlternativeChampion represents alternative champion options for a role
type AlternativeChampion struct {
	ChampionName    string  `json:"champion_name"`
	SimilarityScore float64 `json:"similarity_score"` // 0-100
	StrengthRating  float64 `json:"strength_rating"`  // 0-100
	SynergyScore    float64 `json:"synergy_score"`    // 0-100
	Reasoning       string  `json:"reasoning"`
}

// FlexOption represents flexible pick options
type FlexOption struct {
	AlternativeRole string  `json:"alternative_role"`
	FlexStrength    float64 `json:"flex_strength"`   // 0-100
	RoleEfficiency  float64 `json:"role_efficiency"` // 0-100
	StrategicValue  string  `json:"strategic_value"`
}

// WinRateAnalysis contains comprehensive win rate data
type WinRateAnalysis struct {
	OverallWinRate  float64          `json:"overall_win_rate"` // 0-100
	EarlyGameWR     float64          `json:"early_game_wr"`    // 0-100
	MidGameWR       float64          `json:"mid_game_wr"`      // 0-100
	LateGameWR      float64          `json:"late_game_wr"`     // 0-100
	RankBreakdown   []RankWinRate    `json:"rank_breakdown" gorm:"type:text"`
	RegionVariance  []RegionWinRate  `json:"region_variance" gorm:"type:text"`
	MatchupWinRates []MatchupWinRate `json:"matchup_winrates" gorm:"type:text"`
	SampleSize      int              `json:"sample_size"`
	Confidence      float64          `json:"confidence"` // 0-100
}

// RankWinRate represents win rate by rank tier
type RankWinRate struct {
	Rank       string  `json:"rank"`
	WinRate    float64 `json:"win_rate"` // 0-100
	SampleSize int     `json:"sample_size"`
	Confidence float64 `json:"confidence"` // 0-100
}

// RegionWinRate represents win rate by region
type RegionWinRate struct {
	Region          string  `json:"region"`
	WinRate         float64 `json:"win_rate"` // 0-100
	SampleSize      int     `json:"sample_size"`
	PopularityIndex float64 `json:"popularity_index"` // 0-100
}

// MatchupWinRate represents win rate against specific compositions
type MatchupWinRate struct {
	OpponentType  string  `json:"opponent_type"`
	WinRate       float64 `json:"win_rate"` // 0-100
	SampleSize    int     `json:"sample_size"`
	MatchupRating string  `json:"matchup_rating"` // favored, even, unfavored
}

// CompositionScaling represents how the composition scales over time
type CompositionScaling struct {
	EarlyGameStrength float64 `json:"early_game_strength"` // 0-100 (0-15 min)
	MidGameStrength   float64 `json:"mid_game_strength"`   // 0-100 (15-25 min)
	LateGameStrength  float64 `json:"late_game_strength"`  // 0-100 (25+ min)
	ScalingCurve      string  `json:"scaling_curve"`       // linear, exponential, plateau, spike
	PeakTiming        int     `json:"peak_timing"`         // minute when comp peaks
	FalloffPoint      int     `json:"falloff_point"`       // minute when comp starts declining
	ConsistencyRating float64 `json:"consistency_rating"`  // 0-100
}

// PowerSpikeData represents significant power spikes

// CompositionWinCondition represents paths to victory
type CompositionWinCondition struct {
	Condition    string   `json:"condition"`
	Probability  float64  `json:"probability"` // 0-100
	GamePhase    string   `json:"game_phase"`  // early, mid, late, any
	Requirements []string `json:"requirements"`
	EnabledBy    []string `json:"enabled_by"`
	CounteredBy  []string `json:"countered_by"`
	SuccessRate  float64  `json:"success_rate"` // 0-100
}

// CompositionWeakness represents vulnerabilities
type CompositionWeakness struct {
	Weakness          string   `json:"weakness"`
	Severity          string   `json:"severity"`           // low, medium, high, critical
	ExploitRate       float64  `json:"exploit_rate"`       // 0-100
	GamePhase         string   `json:"game_phase"`         // when this weakness is most exploitable
	Mitigation        []string `json:"mitigation"`         // how to minimize this weakness
	CounterStrategies []string `json:"counter_strategies"` // how opponents exploit this
}

// CompositionOptimizationRequest represents a request for composition optimization
type CompositionOptimizationRequest struct {
	OptimizationType string                   `json:"optimization_type"` // maximize_win_rate, meta_strength, synergy, etc.
	Constraints      OptimizationConstraints  `json:"constraints"`
	PlayerData       []PlayerOptimizationData `json:"player_data"`
	MetaContext      MetaOptimizationContext  `json:"meta_context"`
	Preferences      OptimizationPreferences  `json:"preferences"`
}

// OptimizationConstraints defines limitations for optimization
type OptimizationConstraints struct {
	BannedChampions   []string                  `json:"banned_champions"`
	RequiredChampions []ChampionRoleRequirement `json:"required_champions"`
	ForbiddenRoles    []string                  `json:"forbidden_roles"`    // roles player can't/won't play
	ChampionPool      map[string][]string       `json:"champion_pool"`      // role -> available champions
	MaxDifficulty     string                    `json:"max_difficulty"`     // easy, medium, hard, expert
	MinMetaRelevance  float64                   `json:"min_meta_relevance"` // 0-100
}

// ChampionRoleRequirement specifies required champion-role combinations
type ChampionRoleRequirement struct {
	ChampionName string `json:"champion_name"`
	Role         string `json:"role"`
	Priority     int    `json:"priority"` // 1 = must have, 2 = strongly preferred, etc.
}

// PlayerOptimizationData contains player-specific data for optimization
type PlayerOptimizationData struct {
	PlayerID          string                `json:"player_id"`
	Role              string                `json:"role"`
	ChampionMastery   []ChampionMasteryData `json:"champion_mastery"`
	PlayStyle         string                `json:"play_style"`       // aggressive, passive, supportive, etc.
	SkillLevel        float64               `json:"skill_level"`      // 0-100
	RoleFlexibility   []string              `json:"role_flexibility"` // other roles they can play
	RecentPerformance RecentPerformanceData `json:"recent_performance"`
}

// ChampionMasteryData represents mastery on specific champions
type ChampionMasteryData struct {
	ChampionName      string  `json:"champion_name"`
	MasteryLevel      int     `json:"mastery_level"` // 1-7
	MasteryPoints     int     `json:"mastery_points"`
	GamesPlayed       int     `json:"games_played"`
	WinRate           float64 `json:"win_rate"`           // 0-100
	PerformanceRating float64 `json:"performance_rating"` // 0-100
	ComfortLevel      float64 `json:"comfort_level"`      // 0-100
}

// RecentPerformanceData represents recent performance trends
type RecentPerformanceData struct {
	RecentGames        int     `json:"recent_games"`
	WinRate            float64 `json:"win_rate"` // 0-100
	AvgKDA             float64 `json:"avg_kda"`
	ConsistencyRating  float64 `json:"consistency_rating"` // 0-100
	FormTrend          string  `json:"form_trend"`         // improving, stable, declining
	BestPerformingRole string  `json:"best_performing_role"`
}

// MetaOptimizationContext provides meta context for optimization
type MetaOptimizationContext struct {
	CurrentPatch  string `json:"current_patch"`
	Region        string `json:"region"`
	RankTier      string `json:"rank_tier"`
	GameMode      string `json:"game_mode"`      // ranked, tournament, etc.
	PriorityMeta  bool   `json:"priority_meta"`  // prioritize meta strength
	AvoidBans     bool   `json:"avoid_bans"`     // avoid high ban rate champions
	MetaStability string `json:"meta_stability"` // stable, volatile, new
}

// OptimizationPreferences contains user preferences for optimization
type OptimizationPreferences struct {
	PreferredStyle  string   `json:"preferred_style"`  // team_fight, split_push, poke, etc.
	AggressionLevel string   `json:"aggression_level"` // passive, balanced, aggressive
	ComplexityLevel string   `json:"complexity_level"` // simple, moderate, complex
	LearningGoals   []string `json:"learning_goals"`   // improve_teamwork, learn_new_champs, etc.
	TimeInvestment  string   `json:"time_investment"`  // low, medium, high
	RiskTolerance   string   `json:"risk_tolerance"`   // conservative, balanced, risky
}

// OptimizationResult represents the result of composition optimization
type OptimizationResult struct {
	OptimalCompositions  []TeamComposition            `json:"optimal_compositions"`
	OptimizationScore    float64                      `json:"optimization_score"` // 0-100
	AlternativeOptions   []TeamComposition            `json:"alternative_options"`
	OptimizationReport   OptimizationReport           `json:"optimization_report"`
	Recommendations      []OptimizationRecommendation `json:"recommendations"`
	ConstraintsSatisfied bool                         `json:"constraints_satisfied"`
}

// OptimizationReport provides detailed analysis of the optimization process
type OptimizationReport struct {
	OptimizationTime      float64            `json:"optimization_time"` // milliseconds
	CompositionsEvaluated int                `json:"compositions_evaluated"`
	OptimizationMethod    string             `json:"optimization_method"` // genetic_algorithm, brute_force, etc.
	ConvergenceScore      float64            `json:"convergence_score"`   // 0-100
	ConstraintsApplied    []ConstraintReport `json:"constraints_applied"`
	ImprovementAreas      []ImprovementArea  `json:"improvement_areas"`
	ConfidenceLevel       float64            `json:"confidence_level"` // 0-100
}

// ConstraintReport shows which constraints were applied
type ConstraintReport struct {
	ConstraintType string   `json:"constraint_type"`
	Applied        bool     `json:"applied"`
	Impact         string   `json:"impact"` // low, medium, high
	LimitedOptions []string `json:"limited_options"`
}

// ImprovementArea suggests areas for composition improvement
type ImprovementArea struct {
	Area                   string   `json:"area"`
	CurrentScore           float64  `json:"current_score"`   // 0-100
	PotentialScore         float64  `json:"potential_score"` // 0-100
	ImprovementSuggestions []string `json:"improvement_suggestions"`
	Priority               string   `json:"priority"` // low, medium, high
}

// OptimizationRecommendation provides specific recommendations
type OptimizationRecommendation struct {
	Type                string   `json:"type"` // champion_swap, role_change, strategy_adjust
	Description         string   `json:"description"`
	ExpectedImprovement float64  `json:"expected_improvement"` // 0-100
	Difficulty          string   `json:"difficulty"`           // easy, medium, hard
	Implementation      []string `json:"implementation"`       // steps to implement
	Priority            string   `json:"priority"`             // low, medium, high
}

// OptimizeComposition generates optimal team compositions based on constraints and preferences
func (s *TeamCompositionService) OptimizeComposition(request CompositionOptimizationRequest) (*OptimizationResult, error) {
	startTime := time.Now()

	// Initialize optimization result
	result := &OptimizationResult{
		ConstraintsSatisfied: true,
	}

	// Generate candidate compositions based on constraints
	candidateCompositions, err := s.generateCandidateCompositions(request)
	if err != nil {
		return nil, fmt.Errorf("failed to generate candidate compositions: %w", err)
	}

	// Evaluate and score all candidates
	scoredCompositions := s.evaluateCompositions(candidateCompositions, request)

	// Sort by optimization score
	sort.Slice(scoredCompositions, func(i, j int) bool {
		return scoredCompositions[i].OverallRating > scoredCompositions[j].OverallRating
	})

	// Select optimal compositions
	if len(scoredCompositions) > 0 {
		result.OptimalCompositions = scoredCompositions[:min(3, len(scoredCompositions))]
		result.OptimizationScore = scoredCompositions[0].OverallRating
	}

	// Select alternative options
	if len(scoredCompositions) > 3 {
		result.AlternativeOptions = scoredCompositions[3:min(8, len(scoredCompositions))]
	}

	// Generate optimization report
	result.OptimizationReport = s.generateOptimizationReport(request, scoredCompositions, startTime)

	// Generate recommendations
	result.Recommendations = s.generateOptimizationRecommendations(request, result.OptimalCompositions)

	return result, nil
}

// generateCandidateCompositions creates candidate compositions based on constraints
func (s *TeamCompositionService) generateCandidateCompositions(request CompositionOptimizationRequest) ([]TeamComposition, error) {
	var compositions []TeamComposition

	// Get available champions for each role based on constraints
	availableChampions := s.getAvailableChampions(request.Constraints, request.PlayerData)

	// Generate compositions based on different strategies
	strategies := []string{"meta_optimal", "synergy_focused", "balanced", "comfort_picks"}

	for _, strategy := range strategies {
		strategyCompositions := s.generateCompositionsForStrategy(strategy, availableChampions, request)
		compositions = append(compositions, strategyCompositions...)
	}

	// Add predefined meta compositions if they satisfy constraints
	metaCompositions := s.getMetaCompositions(request.MetaContext)
	for _, comp := range metaCompositions {
		if s.satisfiesConstraints(comp, request.Constraints) {
			compositions = append(compositions, comp)
		}
	}

	return compositions, nil
}

// getAvailableChampions returns available champions for each role based on constraints
func (s *TeamCompositionService) getAvailableChampions(constraints OptimizationConstraints, playerData []PlayerOptimizationData) map[string][]string {
	availableChampions := make(map[string][]string)

	// Initialize with basic champion pools or use provided pools
	if len(constraints.ChampionPool) > 0 {
		availableChampions = constraints.ChampionPool
	} else {
		// Use default champion pools
		availableChampions = map[string][]string{
			"TOP":     {"Garen", "Darius", "Fiora", "Camille", "Malphite", "Ornn", "Aatrox", "Jax"},
			"JUNGLE":  {"Graves", "Kindred", "Lee Sin", "Elise", "Amumu", "Rammus", "Kha'Zix", "Evelynn"},
			"MID":     {"Yasuo", "Zed", "Orianna", "Syndra", "Ahri", "LeBlanc", "Lux", "Malzahar"},
			"ADC":     {"Jinx", "Kai'Sa", "Ezreal", "Caitlyn", "Jhin", "Ashe", "Vayne", "Aphelios"},
			"SUPPORT": {"Thresh", "Leona", "Nautilus", "Soraka", "Lulu", "Braum", "Pyke", "Blitzcrank"},
		}
	}

	// Remove banned champions
	for role, champions := range availableChampions {
		var filtered []string
		for _, champion := range champions {
			banned := false
			for _, bannedChamp := range constraints.BannedChampions {
				if champion == bannedChamp {
					banned = true
					break
				}
			}
			if !banned {
				filtered = append(filtered, champion)
			}
		}
		availableChampions[role] = filtered
	}

	// Filter by player mastery and preferences if available
	if len(playerData) > 0 {
		for _, player := range playerData {
			if champions, exists := availableChampions[player.Role]; exists {
				var filteredByMastery []string
				for _, champion := range champions {
					// Check if player has reasonable mastery on this champion
					hasGoodMastery := s.playerHasGoodMastery(player, champion)
					if hasGoodMastery || s.championIsEasyToLearn(champion) {
						filteredByMastery = append(filteredByMastery, champion)
					}
				}
				if len(filteredByMastery) > 0 {
					availableChampions[player.Role] = filteredByMastery
				}
			}
		}
	}

	return availableChampions
}

// generateCompositionsForStrategy generates compositions based on a specific strategy
func (s *TeamCompositionService) generateCompositionsForStrategy(strategy string, availableChampions map[string][]string, request CompositionOptimizationRequest) []TeamComposition {
	var compositions []TeamComposition

	switch strategy {
	case "meta_optimal":
		compositions = s.generateMetaOptimalCompositions(availableChampions, request.MetaContext)
	case "synergy_focused":
		compositions = s.generateSynergyFocusedCompositions(availableChampions)
	case "balanced":
		compositions = s.generateBalancedCompositions(availableChampions)
	case "comfort_picks":
		compositions = s.generateComfortCompositions(availableChampions, request.PlayerData)
	}

	return compositions
}

// generateMetaOptimalCompositions generates compositions optimized for current meta
func (s *TeamCompositionService) generateMetaOptimalCompositions(availableChampions map[string][]string, metaContext MetaOptimizationContext) []TeamComposition {
	var compositions []TeamComposition

	// Get current meta champions for each role
	metaChampions := s.getCurrentMetaChampions(metaContext)

	// Generate compositions using meta champions
	topChamps := s.intersectChampions(availableChampions["TOP"], metaChampions["TOP"])
	jungleChamps := s.intersectChampions(availableChampions["JUNGLE"], metaChampions["JUNGLE"])
	midChamps := s.intersectChampions(availableChampions["MID"], metaChampions["MID"])
	adcChamps := s.intersectChampions(availableChampions["ADC"], metaChampions["ADC"])
	supportChamps := s.intersectChampions(availableChampions["SUPPORT"], metaChampions["SUPPORT"])

	// Generate top meta combinations (limit to avoid explosion)
	maxCombinations := 10
	combinations := 0

	for _, top := range topChamps[:min(3, len(topChamps))] {
		for _, jungle := range jungleChamps[:min(3, len(jungleChamps))] {
			for _, mid := range midChamps[:min(3, len(midChamps))] {
				for _, adc := range adcChamps[:min(2, len(adcChamps))] {
					for _, support := range supportChamps[:min(2, len(supportChamps))] {
						if combinations >= maxCombinations {
							break
						}

						comp := s.createComposition(top, jungle, mid, adc, support, "meta_optimal")
						compositions = append(compositions, comp)
						combinations++
					}
				}
			}
		}
	}

	return compositions
}

// generateSynergyFocusedCompositions generates compositions optimized for champion synergies
func (s *TeamCompositionService) generateSynergyFocusedCompositions(availableChampions map[string][]string) []TeamComposition {
	var compositions []TeamComposition

	// Define synergy groups (champions that work well together)
	synergyGroups := []map[string][]string{
		// Engage composition
		{
			"TOP":     {"Malphite", "Ornn", "Garen"},
			"JUNGLE":  {"Ammu", "Rammus", "Sejuani"},
			"MID":     {"Orianna", "Yasuo", "Malzahar"},
			"ADC":     {"Jinx", "Kai'Sa", "Aphelios"},
			"SUPPORT": {"Leona", "Nautilus", "Thresh"},
		},
		// Protect the carry composition
		{
			"TOP":     {"Ornn", "Malphite", "Shen"},
			"JUNGLE":  {"Graves", "Kindred", "Sejuani"},
			"MID":     {"Lulu", "Orianna", "Karma"},
			"ADC":     {"Jinx", "Kog'Maw", "Vayne"},
			"SUPPORT": {"Lulu", "Janna", "Soraka"},
		},
		// Poke composition
		{
			"TOP":     {"Jayce", "Gnar", "Kennen"},
			"JUNGLE":  {"Nidalee", "Graves", "Elise"},
			"MID":     {"Zoe", "Xerath", "Vel'Koz"},
			"ADC":     {"Ezreal", "Varus", "Jhin"},
			"SUPPORT": {"Xerath", "Vel'Koz", "Brand"},
		},
	}

	// Generate compositions for each synergy group
	for _, group := range synergyGroups {
		topOptions := s.intersectChampions(availableChampions["TOP"], group["TOP"])
		jungleOptions := s.intersectChampions(availableChampions["JUNGLE"], group["JUNGLE"])
		midOptions := s.intersectChampions(availableChampions["MID"], group["MID"])
		adcOptions := s.intersectChampions(availableChampions["ADC"], group["ADC"])
		supportOptions := s.intersectChampions(availableChampions["SUPPORT"], group["SUPPORT"])

		// Generate best combinations from this synergy group
		if len(topOptions) > 0 && len(jungleOptions) > 0 && len(midOptions) > 0 &&
			len(adcOptions) > 0 && len(supportOptions) > 0 {

			// Take best option from each role for this synergy group
			comp := s.createComposition(
				topOptions[0], jungleOptions[0], midOptions[0],
				adcOptions[0], supportOptions[0], "synergy_focused")
			compositions = append(compositions, comp)
		}
	}

	return compositions
}

// generateBalancedCompositions generates well-rounded compositions
func (s *TeamCompositionService) generateBalancedCompositions(availableChampions map[string][]string) []TeamComposition {
	var compositions []TeamComposition

	// Define balanced champion archetypes for each role
	balancedPicks := map[string][]string{
		"TOP":     {"Garen", "Malphite", "Darius", "Ornn"},   // Mix of tank and damage
		"JUNGLE":  {"Graves", "Ammu", "Lee Sin", "Rammus"},   // Mix of carry and utility
		"MID":     {"Orianna", "Ahri", "Malzahar", "Syndra"}, // Mix of control and burst
		"ADC":     {"Jinx", "Ezreal", "Ashe", "Caitlyn"},     // Mix of scaling and utility
		"SUPPORT": {"Thresh", "Leona", "Soraka", "Braum"},    // Mix of engage and protection
	}

	// Filter by available champions
	filteredPicks := make(map[string][]string)
	for role, champions := range balancedPicks {
		filteredPicks[role] = s.intersectChampions(availableChampions[role], champions)
	}

	// Generate a few balanced compositions
	for i := 0; i < 3; i++ {
		if s.hasEnoughChampions(filteredPicks) {
			comp := s.createComposition(
				filteredPicks["TOP"][i%len(filteredPicks["TOP"])],
				filteredPicks["JUNGLE"][i%len(filteredPicks["JUNGLE"])],
				filteredPicks["MID"][i%len(filteredPicks["MID"])],
				filteredPicks["ADC"][i%len(filteredPicks["ADC"])],
				filteredPicks["SUPPORT"][i%len(filteredPicks["SUPPORT"])],
				"balanced")
			compositions = append(compositions, comp)
		}
	}

	return compositions
}

// generateComfortCompositions generates compositions based on player comfort
func (s *TeamCompositionService) generateComfortCompositions(availableChampions map[string][]string, playerData []PlayerOptimizationData) []TeamComposition {
	var compositions []TeamComposition

	if len(playerData) == 0 {
		return compositions
	}

	// Create player comfort map
	comfortChampions := make(map[string][]string)
	for _, player := range playerData {
		var comfortable []string
		for _, mastery := range player.ChampionMastery {
			if mastery.ComfortLevel >= 70.0 { // Only include champions player is comfortable with
				comfortable = append(comfortable, mastery.ChampionName)
			}
		}
		if len(comfortable) > 0 {
			comfortChampions[player.Role] = s.intersectChampions(availableChampions[player.Role], comfortable)
		}
	}

	// Generate compositions using comfort picks
	if s.hasEnoughChampions(comfortChampions) {
		comp := s.createComposition(
			s.getBestComfortPick(comfortChampions["TOP"]),
			s.getBestComfortPick(comfortChampions["JUNGLE"]),
			s.getBestComfortPick(comfortChampions["MID"]),
			s.getBestComfortPick(comfortChampions["ADC"]),
			s.getBestComfortPick(comfortChampions["SUPPORT"]),
			"comfort_picks")
		compositions = append(compositions, comp)
	}

	return compositions
}

// evaluateCompositions evaluates and scores composition candidates
func (s *TeamCompositionService) evaluateCompositions(compositions []TeamComposition, request CompositionOptimizationRequest) []TeamComposition {
	for i := range compositions {
		score := s.calculateCompositionScore(&compositions[i], request)
		compositions[i].OverallRating = score

		// Add additional analysis
		s.analyzeCompositionSynergy(&compositions[i])
		s.analyzeCompositionScaling(&compositions[i])
		s.analyzeCompositionWeaknesses(&compositions[i])
		s.addMetaContext(&compositions[i], request.MetaContext)
	}

	return compositions
}

// calculateCompositionScore calculates overall score for a composition
func (s *TeamCompositionService) calculateCompositionScore(comp *TeamComposition, request CompositionOptimizationRequest) float64 {
	var totalScore float64
	weights := s.getOptimizationWeights(request.OptimizationType)

	// Meta strength score (0-100)
	metaScore := s.calculateMetaScore(comp, request.MetaContext)
	totalScore += metaScore * weights.MetaWeight

	// Synergy score (0-100)
	synergyScore := s.calculateSynergyScore(comp)
	totalScore += synergyScore * weights.SynergyWeight

	// Player comfort score (0-100)
	comfortScore := s.calculateComfortScore(comp, request.PlayerData)
	totalScore += comfortScore * weights.ComfortWeight

	// Win rate score (0-100)
	winRateScore := s.calculateWinRateScore(comp)
	totalScore += winRateScore * weights.WinRateWeight

	// Difficulty penalty (easier compositions get slight bonus)
	difficultyScore := s.calculateDifficultyScore(comp)
	totalScore += difficultyScore * weights.DifficultyWeight

	// Flexibility bonus
	flexibilityScore := s.calculateFlexibilityScore(comp)
	totalScore += flexibilityScore * weights.FlexibilityWeight

	return totalScore
}

// OptimizationWeights defines weights for different optimization criteria
type OptimizationWeights struct {
	MetaWeight        float64
	SynergyWeight     float64
	ComfortWeight     float64
	WinRateWeight     float64
	DifficultyWeight  float64
	FlexibilityWeight float64
}

// getOptimizationWeights returns weights based on optimization type
func (s *TeamCompositionService) getOptimizationWeights(optimizationType string) OptimizationWeights {
	switch optimizationType {
	case "maximize_win_rate":
		return OptimizationWeights{
			MetaWeight:        0.3,
			SynergyWeight:     0.2,
			ComfortWeight:     0.1,
			WinRateWeight:     0.4,
			DifficultyWeight:  0.0,
			FlexibilityWeight: 0.0,
		}
	case "meta_strength":
		return OptimizationWeights{
			MetaWeight:        0.5,
			SynergyWeight:     0.2,
			ComfortWeight:     0.1,
			WinRateWeight:     0.2,
			DifficultyWeight:  0.0,
			FlexibilityWeight: 0.0,
		}
	case "synergy":
		return OptimizationWeights{
			MetaWeight:        0.2,
			SynergyWeight:     0.4,
			ComfortWeight:     0.2,
			WinRateWeight:     0.2,
			DifficultyWeight:  0.0,
			FlexibilityWeight: 0.0,
		}
	case "comfort":
		return OptimizationWeights{
			MetaWeight:        0.15,
			SynergyWeight:     0.15,
			ComfortWeight:     0.5,
			WinRateWeight:     0.15,
			DifficultyWeight:  0.05,
			FlexibilityWeight: 0.0,
		}
	case "balanced":
	default:
		return OptimizationWeights{
			MetaWeight:        0.25,
			SynergyWeight:     0.25,
			ComfortWeight:     0.2,
			WinRateWeight:     0.2,
			DifficultyWeight:  0.05,
			FlexibilityWeight: 0.05,
		}
	}
}

// Helper functions for composition analysis

func (s *TeamCompositionService) calculateMetaScore(comp *TeamComposition, metaContext MetaOptimizationContext) float64 {
	// This would analyze current meta strength of champions
	// For now, return mock score
	return 75.0 + float64((time.Now().UnixNano() % 20)) - 10.0
}

func (s *TeamCompositionService) calculateSynergyScore(comp *TeamComposition) float64 {
	// This would analyze champion synergies
	// For now, return mock score based on composition type
	baseScore := 70.0

	// Team fight compositions typically have good synergy
	if comp.CompositionType == "team_fight" {
		baseScore += 10.0
	}

	variance := float64((time.Now().UnixNano() % 20)) - 10.0
	return math.Max(0, math.Min(100, baseScore+variance))
}

func (s *TeamCompositionService) calculateComfortScore(comp *TeamComposition, playerData []PlayerOptimizationData) float64 {
	if len(playerData) == 0 {
		return 50.0 // neutral score if no player data
	}

	totalComfort := 0.0
	playersFound := 0

	champions := []string{
		comp.TopLaner.ChampionName,
		comp.Jungler.ChampionName,
		comp.MidLaner.ChampionName,
		comp.ADCarry.ChampionName,
		comp.Support.ChampionName,
	}

	for _, player := range playerData {
		var roleChampion string
		switch player.Role {
		case "TOP":
			roleChampion = comp.TopLaner.ChampionName
		case "JUNGLE":
			roleChampion = comp.Jungler.ChampionName
		case "MID":
			roleChampion = comp.MidLaner.ChampionName
		case "ADC":
			roleChampion = comp.ADCarry.ChampionName
		case "SUPPORT":
			roleChampion = comp.Support.ChampionName
		}

		if roleChampion != "" {
			comfort := s.getPlayerChampionComfort(player, roleChampion)
			totalComfort += comfort
			playersFound++
		}
	}

	if playersFound > 0 {
		return totalComfort / float64(playersFound)
	}

	return 50.0
}

func (s *TeamCompositionService) calculateWinRateScore(comp *TeamComposition) float64 {
	// This would look up actual win rates
	// For now, return mock data
	baseWinRate := 50.0
	variance := float64((time.Now().UnixNano() % 20)) - 10.0
	return math.Max(30, math.Min(70, baseWinRate+variance))
}

func (s *TeamCompositionService) calculateDifficultyScore(comp *TeamComposition) float64 {
	// Easier compositions get higher scores
	// This would analyze champion difficulty
	return 75.0 // mock score
}

func (s *TeamCompositionService) calculateFlexibilityScore(comp *TeamComposition) float64 {
	// This would analyze how flexible the composition is
	return 65.0 // mock score
}

// Helper functions

func (s *TeamCompositionService) createComposition(top, jungle, mid, adc, support, strategy string) TeamComposition {
	return TeamComposition{
		ID:              fmt.Sprintf("comp_%d", time.Now().UnixNano()),
		CompositionName: fmt.Sprintf("%s Composition", strategy),
		CompositionType: s.identifyCompositionType(top, jungle, mid, adc, support),
		TopLaner:        s.createChampionAssignment(top, "TOP"),
		Jungler:         s.createChampionAssignment(jungle, "JUNGLE"),
		MidLaner:        s.createChampionAssignment(mid, "MID"),
		ADCarry:         s.createChampionAssignment(adc, "ADC"),
		Support:         s.createChampionAssignment(support, "SUPPORT"),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		LastAnalyzed:    time.Now(),
	}
}

func (s *TeamCompositionService) createChampionAssignment(champion, role string) ChampionAssignment {
	return ChampionAssignment{
		ChampionName:   champion,
		Role:           role,
		Priority:       1,
		RoleEfficiency: 80.0,
		MetaStrength:   75.0,
		SynergyContrib: 70.0,
		CounterResist:  65.0,
		WinRate:        52.0,
		PickRate:       15.0,
		BanRate:        8.0,
		AvgKDA:         2.1,
		AvgDamageShare: 22.0,
	}
}

func (s *TeamCompositionService) identifyCompositionType(top, jungle, mid, adc, support string) string {
	// This would analyze the champions and determine composition type
	// For now, return common types
	types := []string{"team_fight", "split_push", "poke", "pick", "siege"}
	return types[int(time.Now().UnixNano())%len(types)]
}

func (s *TeamCompositionService) intersectChampions(available, desired []string) []string {
	var result []string
	for _, champ := range desired {
		for _, avail := range available {
			if champ == avail {
				result = append(result, champ)
				break
			}
		}
	}
	return result
}

func (s *TeamCompositionService) hasEnoughChampions(championMap map[string][]string) bool {
	requiredRoles := []string{"TOP", "JUNGLE", "MID", "ADC", "SUPPORT"}
	for _, role := range requiredRoles {
		if len(championMap[role]) == 0 {
			return false
		}
	}
	return true
}

func (s *TeamCompositionService) getBestComfortPick(champions []string) string {
	if len(champions) == 0 {
		return "Unknown"
	}
	return champions[0] // Return first (assumed best) option
}

func (s *TeamCompositionService) getCurrentMetaChampions(metaContext MetaOptimizationContext) map[string][]string {
	// This would fetch current meta champions from meta service
	// For now, return mock meta champions
	return map[string][]string{
		"TOP":     {"Aatrox", "Fiora", "Camille", "Ornn", "Malphite"},
		"JUNGLE":  {"Graves", "Lee Sin", "Elise", "Ammu", "Kindred"},
		"MID":     {"Yasuo", "Orianna", "Syndra", "LeBlanc", "Ahri"},
		"ADC":     {"Jinx", "Kai'Sa", "Ezreal", "Aphelios", "Caitlyn"},
		"SUPPORT": {"Thresh", "Leona", "Nautilus", "Lulu", "Braum"},
	}
}

func (s *TeamCompositionService) satisfiesConstraints(comp TeamComposition, constraints OptimizationConstraints) bool {
	// Check if composition satisfies all constraints
	champions := []string{
		comp.TopLaner.ChampionName,
		comp.Jungler.ChampionName,
		comp.MidLaner.ChampionName,
		comp.ADCarry.ChampionName,
		comp.Support.ChampionName,
	}

	// Check banned champions
	for _, champion := range champions {
		for _, banned := range constraints.BannedChampions {
			if champion == banned {
				return false
			}
		}
	}

	return true
}

func (s *TeamCompositionService) getMetaCompositions(metaContext MetaOptimizationContext) []TeamComposition {
	// This would return current meta compositions
	// For now, return empty slice
	return []TeamComposition{}
}

func (s *TeamCompositionService) playerHasGoodMastery(player PlayerOptimizationData, champion string) bool {
	for _, mastery := range player.ChampionMastery {
		if mastery.ChampionName == champion && mastery.ComfortLevel >= 60.0 {
			return true
		}
	}
	return false
}

func (s *TeamCompositionService) championIsEasyToLearn(champion string) bool {
	// Define easy-to-learn champions
	easyChampions := []string{"Garen", "Malphite", "Ammu", "Annie", "Ashe", "Soraka"}
	for _, easy := range easyChampions {
		if champion == easy {
			return true
		}
	}
	return false
}

func (s *TeamCompositionService) getPlayerChampionComfort(player PlayerOptimizationData, champion string) float64 {
	for _, mastery := range player.ChampionMastery {
		if mastery.ChampionName == champion {
			return mastery.ComfortLevel
		}
	}
	return 30.0 // low comfort for unknown champions
}

func (s *TeamCompositionService) analyzeCompositionSynergy(comp *TeamComposition) {
	// Analyze champion synergies and update synergy score
	comp.SynergyScore = s.calculateSynergyScore(comp)
}

func (s *TeamCompositionService) analyzeCompositionScaling(comp *TeamComposition) {
	// Analyze how the composition scales over time
	comp.ScalingProfile = CompositionScaling{
		EarlyGameStrength: 70.0 + float64((time.Now().UnixNano() % 20)) - 10.0,
		MidGameStrength:   75.0 + float64((time.Now().UnixNano() % 20)) - 10.0,
		LateGameStrength:  80.0 + float64((time.Now().UnixNano() % 20)) - 10.0,
		ScalingCurve:      "linear",
		PeakTiming:        25,
		FalloffPoint:      35,
		ConsistencyRating: 75.0,
	}
}

func (s *TeamCompositionService) analyzeCompositionWeaknesses(comp *TeamComposition) {
	// Identify composition weaknesses
	comp.Weaknesses = []CompositionWeakness{
		{
			Weakness:          "Early game vulnerability",
			Severity:          "medium",
			ExploitRate:       65.0,
			GamePhase:         "early",
			Mitigation:        []string{"Safe laning", "Ward coverage", "Jungle protection"},
			CounterStrategies: []string{"Early aggression", "Invades", "Lane pressure"},
		},
	}
}

func (s *TeamCompositionService) addMetaContext(comp *TeamComposition, metaContext MetaOptimizationContext) {
	// Add meta relevance and context
	comp.MetaRelevance = 80.0
	comp.PatchStability = 85.0
	comp.TrendDirection = "stable"
}

func (s *TeamCompositionService) generateOptimizationReport(request CompositionOptimizationRequest, compositions []TeamComposition, startTime time.Time) OptimizationReport {
	return OptimizationReport{
		OptimizationTime:      float64(time.Since(startTime).Milliseconds()),
		CompositionsEvaluated: len(compositions),
		OptimizationMethod:    "heuristic_search",
		ConvergenceScore:      85.0,
		ConfidenceLevel:       82.0,
	}
}

func (s *TeamCompositionService) generateOptimizationRecommendations(request CompositionOptimizationRequest, compositions []TeamComposition) []OptimizationRecommendation {
	var recommendations []OptimizationRecommendation

	if len(compositions) > 0 {
		recommendations = append(recommendations, OptimizationRecommendation{
			Type:                "strategy_focus",
			Description:         "Consider practicing team fight coordination to maximize this composition's potential",
			ExpectedImprovement: 15.0,
			Difficulty:          "medium",
			Implementation:      []string{"Practice 5v5 team fights", "Work on engage timing", "Coordinate ultimate usage"},
			Priority:            "high",
		})
	}

	return recommendations
}

// Utility functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetCompositionsByType retrieves compositions filtered by type and criteria
func (s *TeamCompositionService) GetCompositionsByType(compositionType string, tier string, limit int) ([]*TeamComposition, error) {
	var compositions []*TeamComposition

	query := s.db.Where("composition_type = ?", compositionType)

	if tier != "" {
		query = query.Where("tier = ?", tier)
	}

	err := query.Order("overall_rating DESC").
		Limit(limit).
		Find(&compositions).Error

	return compositions, err
}

// AnalyzeCounterCompositions analyzes compositions that counter a given composition
func (s *TeamCompositionService) AnalyzeCounterCompositions(targetComposition TeamComposition) ([]TeamComposition, error) {
	// This would analyze and return compositions that counter the target
	// For now, return empty slice
	return []TeamComposition{}, nil
}

// GetCompositionWinRate gets detailed win rate data for a composition
func (s *TeamCompositionService) GetCompositionWinRate(compositionID string) (*WinRateAnalysis, error) {
	var composition TeamComposition
	err := s.db.Where("id = ?", compositionID).First(&composition).Error
	if err != nil {
		return nil, err
	}

	return &composition.WinRateData, nil
}
