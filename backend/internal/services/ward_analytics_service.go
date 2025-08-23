package services

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/herald/internal/models"
)

// WardAnalyticsService handles ward placement and map control analytics
type WardAnalyticsService struct {
	visionService *VisionAnalyticsService
	mapService    *MapService
}

// WardAnalysis represents comprehensive ward placement and map control analysis
type WardAnalysis struct {
	PlayerID              string                    `json:"player_id"`
	Champion              string                    `json:"champion,omitempty"`
	Position              string                    `json:"position,omitempty"`
	TimeRange             string                    `json:"time_range"`
	
	// Core Ward Metrics
	AverageWardsPlaced    float64                   `json:"average_wards_placed"`
	AverageWardsKilled    float64                   `json:"average_wards_killed"`
	WardEfficiency        float64                   `json:"ward_efficiency"` // Wards killed / Wards placed
	WardSurvivalRate      float64                   `json:"ward_survival_rate"` // % wards that survive full duration
	
	// Map Control Analysis
	MapControlScore       float64                   `json:"map_control_score"` // 0-100 comprehensive score
	TerritoryControlled   float64                   `json:"territory_controlled"` // % of map under vision
	StrategicCoverage     StrategicCoverageData     `json:"strategic_coverage"`
	
	// Ward Placement Analysis
	PlacementPatterns     WardPlacementPatterns     `json:"placement_patterns"`
	OptimalPlacements     OptimalPlacementData      `json:"optimal_placements"`
	PlacementTiming       PlacementTimingData       `json:"placement_timing"`
	
	// Ward Types Analysis
	YellowWardsAnalysis   WardTypeAnalysis          `json:"yellow_wards_analysis"`
	ControlWardsAnalysis  WardTypeAnalysis          `json:"control_wards_analysis"`
	BlueWardAnalysis      WardTypeAnalysis          `json:"blue_ward_analysis"`
	
	// Clearing Analysis
	WardClearingPatterns  WardClearingData          `json:"ward_clearing_patterns"`
	CounterWardingScore   float64                   `json:"counter_warding_score"` // 0-100 score
	
	// Game Phase Analysis
	EarlyGameWards        WardPhaseData             `json:"early_game_wards"`
	MidGameWards          WardPhaseData             `json:"mid_game_wards"`
	LateGameWards         WardPhaseData             `json:"late_game_wards"`
	
	// Zone-Specific Analysis
	ZoneControl           map[string]ZoneControlData `json:"zone_control"`
	RiverControl          RiverControlData          `json:"river_control"`
	JungleControl         JungleControlData         `json:"jungle_control"`
	
	// Impact Analysis
	VisionDeniedScore     float64                   `json:"vision_denied_score"`
	ObjectiveSetup        ObjectiveSetupData        `json:"objective_setup"`
	SafetyProvided        SafetyData                `json:"safety_provided"`
	
	// Comparative Analysis
	RoleBenchmark         WardBenchmark             `json:"role_benchmark"`
	RankBenchmark         WardBenchmark             `json:"rank_benchmark"`
	GlobalBenchmark       WardBenchmark             `json:"global_benchmark"`
	
	// Performance Impact
	HighVisionWinRate     float64                   `json:"high_vision_win_rate"`
	LowVisionWinRate      float64                   `json:"low_vision_win_rate"`
	WardImpactScore       float64                   `json:"ward_impact_score"`
	
	// Trend Analysis
	TrendDirection        string                    `json:"trend_direction"`
	TrendSlope            float64                   `json:"trend_slope"`
	TrendConfidence       float64                   `json:"trend_confidence"`
	TrendData             []WardTrendPoint          `json:"trend_data"`
	
	// Optimization
	PlacementOptimization PlacementOptimizationData `json:"placement_optimization"`
	ClearingOptimization  ClearingOptimizationData  `json:"clearing_optimization"`
	
	// Insights and Recommendations
	StrengthAreas         []string                  `json:"strength_areas"`
	ImprovementAreas      []string                  `json:"improvement_areas"`
	Recommendations       []WardRecommendation      `json:"recommendations"`
	
	// Match Performance
	RecentMatches         []MatchWardData           `json:"recent_matches"`
}

// StrategicCoverageData represents coverage of strategic map areas
type StrategicCoverageData struct {
	DragonPitCoverage     float64                   `json:"dragon_pit_coverage"`     // % of time dragon pit is warded
	BaronPitCoverage      float64                   `json:"baron_pit_coverage"`      // % of time baron pit is warded
	RiverBrushesCoverage  float64                   `json:"river_brushes_coverage"`  // % of time river brushes are warded
	JungleEntryCoverage   float64                   `json:"jungle_entry_coverage"`   // % of jungle entrances covered
	LaneBrushesCoverage   float64                   `json:"lane_brushes_coverage"`   // % of lane brushes covered
	OverallStrategicScore float64                   `json:"overall_strategic_score"` // Weighted strategic coverage score
}

// WardPlacementPatterns represents analysis of ward placement patterns
type WardPlacementPatterns struct {
	AggregateScore        float64                   `json:"aggregate_score"`        // How well wards are clustered
	DiversityScore        float64                   `json:"diversity_score"`        // How diverse placement locations are
	PredictabilityScore   float64                   `json:"predictability_score"`   // How predictable placements are
	AdaptabilityScore     float64                   `json:"adaptability_score"`     // How well player adapts to game state
	OptimalityScore       float64                   `json:"optimality_score"`       // How often optimal spots are chosen
	
	// Common Placement Zones
	FavoriteZones         []ZonePlacementData       `json:"favorite_zones"`
	AvoidedZones          []ZonePlacementData       `json:"avoided_zones"`
	UnderUtilizedZones    []ZonePlacementData       `json:"under_utilized_zones"`
}

// OptimalPlacementData represents analysis of placement optimality
type OptimalPlacementData struct {
	OptimalPlacements     int                       `json:"optimal_placements"`     // Count of optimal placements
	SuboptimalPlacements  int                       `json:"suboptimal_placements"`  // Count of suboptimal placements
	WastedPlacements      int                       `json:"wasted_placements"`      // Count of wasted placements
	OptimalityRate        float64                   `json:"optimality_rate"`        // % of placements that are optimal
	
	// Specific Optimization Areas
	TimingOptimization    float64                   `json:"timing_optimization"`    // % improvement possible through timing
	LocationOptimization  float64                   `json:"location_optimization"`  // % improvement possible through location
	TypeOptimization      float64                   `json:"type_optimization"`      // % improvement possible through ward type
}

// PlacementTimingData represents ward placement timing analysis
type PlacementTimingData struct {
	PreObjectiveWarding   float64                   `json:"pre_objective_warding"`   // % wards placed before objectives
	ReactiveWarding       float64                   `json:"reactive_warding"`        // % wards placed reactively
	ProactiveWarding      float64                   `json:"proactive_warding"`       // % wards placed proactively
	EmergencyWarding      float64                   `json:"emergency_warding"`       // % wards placed in emergency
	
	// Timing Quality
	ExcellentTiming       int                       `json:"excellent_timing"`        // Count of perfectly timed wards
	GoodTiming            int                       `json:"good_timing"`             // Count of well-timed wards
	PoorTiming            int                       `json:"poor_timing"`             // Count of poorly timed wards
	TimingScore           float64                   `json:"timing_score"`            // Overall timing score 0-100
}

// WardTypeAnalysis represents analysis for specific ward types
type WardTypeAnalysis struct {
	WardType              string                    `json:"ward_type"`               // "YELLOW", "CONTROL", "BLUE_TRINKET"
	AveragePlaced         float64                   `json:"average_placed"`          // Average wards placed per game
	PlacementEfficiency   float64                   `json:"placement_efficiency"`    // How efficiently this ward type is used
	SurvivalRate          float64                   `json:"survival_rate"`           // % that survive full duration
	StrategicUsage        float64                   `json:"strategic_usage"`         // % placed in strategic locations
	TimingQuality         float64                   `json:"timing_quality"`          // Quality of timing for this ward type
	ImpactScore           float64                   `json:"impact_score"`            // Overall impact score for this ward type
	
	// Usage Patterns
	EarlyGameUsage        float64                   `json:"early_game_usage"`        // % used in early game
	MidGameUsage          float64                   `json:"mid_game_usage"`          // % used in mid game
	LateGameUsage         float64                   `json:"late_game_usage"`         // % used in late game
	
	// Optimization Potential
	OptimizationPotential float64                   `json:"optimization_potential"`  // % improvement possible
	RecommendedUsage      string                    `json:"recommended_usage"`       // Usage recommendation
}

// WardClearingData represents ward clearing analysis
type WardClearingData struct {
	TotalWardsCleared     int                       `json:"total_wards_cleared"`     // Total enemy wards cleared
	ClearingEfficiency    float64                   `json:"clearing_efficiency"`     // Wards cleared per opportunity
	StrategicClearing     float64                   `json:"strategic_clearing"`      // % cleared from strategic areas
	TimingQuality         float64                   `json:"timing_quality"`          // Quality of clearing timing
	SafetyScore           float64                   `json:"safety_score"`            // How safely clearing is performed
	
	// Clearing Patterns
	ProactiveClearing     float64                   `json:"proactive_clearing"`      // % clearing done proactively
	ReactiveClearing      float64                   `json:"reactive_clearing"`       // % clearing done reactively
	OpportunisticClearing float64                   `json:"opportunistic_clearing"`  // % clearing done opportunistically
	
	// Zone-Specific Clearing
	RiverClearing         float64                   `json:"river_clearing"`          // % of river wards cleared
	JungleClearing        float64                   `json:"jungle_clearing"`         // % of jungle wards cleared
	ObjectiveClearing     float64                   `json:"objective_clearing"`      // % of objective wards cleared
}

// WardPhaseData represents ward performance by game phase
type WardPhaseData struct {
	Phase                 string                    `json:"phase"`                   // "early", "mid", "late"
	WardsPlaced           float64                   `json:"wards_placed"`            // Average wards placed in this phase
	WardsKilled           float64                   `json:"wards_killed"`            // Average wards killed in this phase
	PlacementQuality      float64                   `json:"placement_quality"`       // Quality of placements 0-100
	StrategicFocus        float64                   `json:"strategic_focus"`         // % of strategic placements
	MapCoverage           float64                   `json:"map_coverage"`            // % of map covered by vision
	EfficiencyRating      string                    `json:"efficiency_rating"`       // "excellent", "good", "average", "poor"
}

// ZoneControlData represents control analysis for specific zones
type ZoneControlData struct {
	ZoneName              string                    `json:"zone_name"`
	ControlPercentage     float64                   `json:"control_percentage"`      // % of time zone is under player's vision
	ContestLevel          string                    `json:"contest_level"`           // "high", "medium", "low" - how contested
	StrategicValue        float64                   `json:"strategic_value"`         // Strategic importance 0-100
	WardsPlaced           int                       `json:"wards_placed"`            // Total wards placed in zone
	WardsCleared          int                       `json:"wards_cleared"`           // Total enemy wards cleared from zone
	ControlEfficiency     float64                   `json:"control_efficiency"`      // How efficiently zone is controlled
	RecommendedFocus      string                    `json:"recommended_focus"`       // Recommendation for this zone
}

// RiverControlData represents specific river control analysis
type RiverControlData struct {
	OverallRiverControl   float64                   `json:"overall_river_control"`   // % of river under control
	TopRiverControl       float64                   `json:"top_river_control"`       // % of top river controlled
	BottomRiverControl    float64                   `json:"bottom_river_control"`    // % of bottom river controlled
	ScuttleCrabControl    float64                   `json:"scuttle_crab_control"`    // % of scuttle crab areas controlled
	RiverBrushControl     float64                   `json:"river_brush_control"`     // % of river brushes controlled
	CrossingPoints        float64                   `json:"crossing_points"`         // % of river crossing points covered
	RiverScore            float64                   `json:"river_score"`             // Overall river control score 0-100
}

// JungleControlData represents jungle vision control analysis
type JungleControlData struct {
	OwnJungleControl      float64                   `json:"own_jungle_control"`      // % of own jungle under vision
	EnemyJungleControl    float64                   `json:"enemy_jungle_control"`    // % of enemy jungle under vision
	BuffControl           float64                   `json:"buff_control"`            // % of buff camps under vision
	CampTimerTracking     float64                   `json:"camp_timer_tracking"`     // % of jungle camps tracked
	InvasionDetection     float64                   `json:"invasion_detection"`      // % of invasions detected early
	CounterJungling       float64                   `json:"counter_jungling"`        // Counter-jungling vision setup
	JungleScore           float64                   `json:"jungle_score"`            // Overall jungle control score 0-100
}

// ObjectiveSetupData represents objective setup analysis
type ObjectiveSetupData struct {
	DragonSetupScore      float64                   `json:"dragon_setup_score"`      // Dragon vision setup quality 0-100
	BaronSetupScore       float64                   `json:"baron_setup_score"`       // Baron vision setup quality 0-100
	HeraldSetupScore      float64                   `json:"herald_setup_score"`      // Herald vision setup quality 0-100
	ElderSetupScore       float64                   `json:"elder_setup_score"`       // Elder dragon setup quality 0-100
	
	// Timing Analysis
	PreObjectiveSetup     float64                   `json:"pre_objective_setup"`     // % objectives with pre-setup
	SetupTiming           float64                   `json:"setup_timing"`            // Average setup time before objectives
	SetupEfficiency       float64                   `json:"setup_efficiency"`        // Setup efficiency score 0-100
	
	// Coverage Analysis
	ApproachCoverage      float64                   `json:"approach_coverage"`       // % of objective approaches covered
	EscapeCoverage        float64                   `json:"escape_coverage"`         // % of escape routes covered
	FlankCoverage         float64                   `json:"flank_coverage"`          // % of flank routes covered
}

// SafetyData represents safety provided by vision
type SafetyData struct {
	GankPrevention        float64                   `json:"gank_prevention"`         // % of potential ganks detected
	InvasionDetection     float64                   `json:"invasion_detection"`      // % of invasions detected
	RotationTracking      float64                   `json:"rotation_tracking"`       // % of enemy rotations tracked
	SafeFarmingProvided   float64                   `json:"safe_farming_provided"`   // % increase in safe farming time
	SafetyScore           float64                   `json:"safety_score"`            // Overall safety score 0-100
	
	// Risk Mitigation
	HighRiskDetection     float64                   `json:"high_risk_detection"`     // % of high-risk situations detected
	MediumRiskDetection   float64                   `json:"medium_risk_detection"`   // % of medium-risk situations detected
	RiskMitigationScore   float64                   `json:"risk_mitigation_score"`   // Overall risk mitigation 0-100
}

// WardBenchmark represents ward performance benchmarks
type WardBenchmark struct {
	Category              string                    `json:"category"`
	AverageWardsPlaced    float64                   `json:"average_wards_placed"`
	AverageWardsKilled    float64                   `json:"average_wards_killed"`
	AverageMapControl     float64                   `json:"average_map_control"`
	Top10Percent          float64                   `json:"top_10_percent"`
	Top25Percent          float64                   `json:"top_25_percent"`
	Median                float64                   `json:"median"`
	PlayerPercentile      float64                   `json:"player_percentile"`
}

// WardTrendPoint represents ward performance over time
type WardTrendPoint struct {
	Date                  time.Time                 `json:"date"`
	WardsPlaced           float64                   `json:"wards_placed"`
	WardsKilled           float64                   `json:"wards_killed"`
	MapControlScore       float64                   `json:"map_control_score"`
	WardEfficiency        float64                   `json:"ward_efficiency"`
	MovingAverage         float64                   `json:"moving_average"`
}

// PlacementOptimizationData represents placement optimization suggestions
type PlacementOptimizationData struct {
	OptimalSpots          []OptimalSpotData         `json:"optimal_spots"`
	TimingImprovements    []TimingImprovementData   `json:"timing_improvements"`
	TypeOptimizations     []TypeOptimizationData    `json:"type_optimizations"`
	
	// Expected Impact
	ExpectedControlGain   float64                   `json:"expected_control_gain"`    // Expected map control increase %
	ExpectedSafetyGain    float64                   `json:"expected_safety_gain"`     // Expected safety increase %
	ImplementationTips    []string                  `json:"implementation_tips"`
}

// ClearingOptimizationData represents clearing optimization suggestions
type ClearingOptimizationData struct {
	PriorityTargets       []PriorityTargetData      `json:"priority_targets"`
	ClearingOpportunities []ClearingOpportunityData `json:"clearing_opportunities"`
	SafetyClearingTips    []string                  `json:"safety_clearing_tips"`
	
	// Expected Impact
	ExpectedDenialGain    float64                   `json:"expected_denial_gain"`     // Expected vision denial increase %
	ExpectedSafetyGain    float64                   `json:"expected_safety_gain"`     // Expected clearing safety increase %
}

// Supporting data structures
type ZonePlacementData struct {
	Zone                  string                    `json:"zone"`
	PlacementCount        int                       `json:"placement_count"`
	PlacementPercentage   float64                   `json:"placement_percentage"`
	Effectiveness         float64                   `json:"effectiveness"`
}

type OptimalSpotData struct {
	Zone                  string                    `json:"zone"`
	Coordinates           [2]int                    `json:"coordinates"`
	StrategicValue        float64                   `json:"strategic_value"`
	CurrentUsage          float64                   `json:"current_usage"`
	RecommendedUsage      float64                   `json:"recommended_usage"`
	Reason                string                    `json:"reason"`
}

type TimingImprovementData struct {
	Situation             string                    `json:"situation"`
	CurrentTiming         float64                   `json:"current_timing"`
	OptimalTiming         float64                   `json:"optimal_timing"`
	ImprovementNeeded     float64                   `json:"improvement_needed"`
	Tips                  []string                  `json:"tips"`
}

type TypeOptimizationData struct {
	Situation             string                    `json:"situation"`
	CurrentType           string                    `json:"current_type"`
	RecommendedType       string                    `json:"recommended_type"`
	EfficiencyGain        float64                   `json:"efficiency_gain"`
	Reasoning             string                    `json:"reasoning"`
}

type PriorityTargetData struct {
	Zone                  string                    `json:"zone"`
	Priority              string                    `json:"priority"`      // "high", "medium", "low"
	StrategicValue        float64                   `json:"strategic_value"`
	ClearingFrequency     float64                   `json:"clearing_frequency"`
	RecommendedFocus      string                    `json:"recommended_focus"`
}

type ClearingOpportunityData struct {
	Situation             string                    `json:"situation"`
	OpportunityType       string                    `json:"opportunity_type"`
	SafetyLevel           string                    `json:"safety_level"`
	ExpectedReward        float64                   `json:"expected_reward"`
	Tips                  []string                  `json:"tips"`
}

// WardRecommendation represents actionable ward advice
type WardRecommendation struct {
	Priority              string                    `json:"priority"`       // "high", "medium", "low"
	Category              string                    `json:"category"`       // "placement", "clearing", "timing", "positioning"
	Title                 string                    `json:"title"`
	Description           string                    `json:"description"`
	Impact                string                    `json:"impact"`
	GamePhase             []string                  `json:"game_phase"`
	ExpectedImprovement   float64                   `json:"expected_improvement"` // Expected score improvement
	ImplementationDifficulty string                 `json:"implementation_difficulty"` // "easy", "medium", "hard"
}

// MatchWardData represents ward performance in a specific match
type MatchWardData struct {
	MatchID               string                    `json:"match_id"`
	Champion              string                    `json:"champion"`
	Position              string                    `json:"position"`
	WardsPlaced           int                       `json:"wards_placed"`
	WardsKilled           int                       `json:"wards_killed"`
	ControlWardsPlaced    int                       `json:"control_wards_placed"`
	MapControlScore       float64                   `json:"map_control_score"`
	WardEfficiency        float64                   `json:"ward_efficiency"`
	StrategicScore        float64                   `json:"strategic_score"`
	GameDuration          int                       `json:"game_duration"`
	Result                string                    `json:"result"`
	Date                  time.Time                 `json:"date"`
	VisionAdvantage       float64                   `json:"vision_advantage"`
	ObjectiveControl      float64                   `json:"objective_control"`
}

// NewWardAnalyticsService creates a new ward analytics service
func NewWardAnalyticsService(visionService *VisionAnalyticsService, mapService *MapService) *WardAnalyticsService {
	return &WardAnalyticsService{
		visionService: visionService,
		mapService:    mapService,
	}
}

// AnalyzeWards performs comprehensive ward placement and map control analysis
func (was *WardAnalyticsService) AnalyzeWards(ctx context.Context, playerID string, timeRange string, champion string, position string) (*WardAnalysis, error) {
	// Get ward placement data
	matches, err := was.getWardMatches(ctx, playerID, timeRange, champion, position)
	if err != nil {
		return nil, fmt.Errorf("failed to get ward match data: %w", err)
	}

	if len(matches) == 0 {
		return &WardAnalysis{
			PlayerID:  playerID,
			Champion:  champion,
			Position:  position,
			TimeRange: timeRange,
		}, nil
	}

	analysis := &WardAnalysis{
		PlayerID:  playerID,
		Champion:  champion,
		Position:  position,
		TimeRange: timeRange,
	}

	// Calculate core ward metrics
	was.calculateWardBasics(analysis, matches)

	// Analyze map control
	was.analyzeMapControl(analysis, matches)

	// Analyze placement patterns
	was.analyzePlacementPatterns(analysis, matches)

	// Analyze ward types
	was.analyzeWardTypes(analysis, matches)

	// Analyze clearing patterns
	was.analyzeClearingPatterns(analysis, matches)

	// Analyze by game phases
	was.analyzeWardPhases(analysis, matches)

	// Analyze zone-specific control
	was.analyzeZoneControl(analysis, matches)

	// Analyze impact
	was.analyzeWardImpact(analysis, matches)

	// Perform comparative analysis
	err = was.performWardComparison(ctx, analysis, position)
	if err != nil {
		fmt.Printf("Warning: failed to perform ward comparison: %v", err)
	}

	// Perform trend analysis
	was.analyzeWardTrend(analysis, matches)

	// Generate optimization suggestions
	was.generateOptimizationSuggestions(analysis)

	// Generate recommendations
	was.generateWardRecommendations(analysis)

	// Prepare trend data for visualization
	was.generateWardTrendData(analysis, matches)

	return analysis, nil
}

// calculateWardBasics calculates fundamental ward metrics
func (was *WardAnalyticsService) calculateWardBasics(analysis *WardAnalysis, matches []models.MatchData) {
	if len(matches) == 0 {
		return
	}

	totalWardsPlaced := 0.0
	totalWardsKilled := 0.0
	totalSurvivalTime := 0.0
	totalMaxSurvivalTime := 0.0

	for _, match := range matches {
		totalWardsPlaced += float64(match.WardsPlaced)
		totalWardsKilled += float64(match.WardsKilled)
		
		// Estimate survival time (simplified)
		avgWardDuration := float64(match.GameDuration) / math.Max(1, float64(match.WardsPlaced)) * 0.6
		totalSurvivalTime += avgWardDuration * float64(match.WardsPlaced)
		totalMaxSurvivalTime += 180.0 * float64(match.WardsPlaced) // Max ward duration
	}

	analysis.AverageWardsPlaced = totalWardsPlaced / float64(len(matches))
	analysis.AverageWardsKilled = totalWardsKilled / float64(len(matches))
	
	if analysis.AverageWardsPlaced > 0 {
		analysis.WardEfficiency = analysis.AverageWardsKilled / analysis.AverageWardsPlaced
	}
	
	if totalMaxSurvivalTime > 0 {
		analysis.WardSurvivalRate = (totalSurvivalTime / totalMaxSurvivalTime) * 100
	}
}

// analyzeMapControl analyzes overall map control
func (was *WardAnalyticsService) analyzeMapControl(analysis *WardAnalysis, matches []models.MatchData) {
	// Calculate map control score based on vision coverage and strategic control
	baseScore := 50.0
	
	// Factor 1: Ward placement rate
	if analysis.AverageWardsPlaced > 15 {
		baseScore += 15
	} else if analysis.AverageWardsPlaced > 10 {
		baseScore += 10
	} else if analysis.AverageWardsPlaced < 5 {
		baseScore -= 15
	}
	
	// Factor 2: Ward clearing rate
	if analysis.AverageWardsKilled > 8 {
		baseScore += 15
	} else if analysis.AverageWardsKilled > 5 {
		baseScore += 10
	} else if analysis.AverageWardsKilled < 2 {
		baseScore -= 10
	}
	
	// Factor 3: Ward efficiency
	if analysis.WardEfficiency > 0.6 {
		baseScore += 10
	} else if analysis.WardEfficiency < 0.3 {
		baseScore -= 10
	}
	
	// Factor 4: Survival rate
	if analysis.WardSurvivalRate > 70 {
		baseScore += 10
	} else if analysis.WardSurvivalRate < 40 {
		baseScore -= 5
	}
	
	analysis.MapControlScore = math.Max(0, math.Min(100, baseScore))
	analysis.TerritoryControlled = baseScore * 0.8 // Estimate territory controlled
}

// analyzePlacementPatterns analyzes ward placement patterns
func (was *WardAnalyticsService) analyzePlacementPatterns(analysis *WardAnalysis, matches []models.MatchData) {
	// Simplified placement pattern analysis
	analysis.PlacementPatterns = WardPlacementPatterns{
		AggregateScore:      75.0, // How well wards work together
		DiversityScore:      80.0, // Diversity of placement locations
		PredictabilityScore: 60.0, // How predictable placements are
		AdaptabilityScore:   70.0, // Adaptation to game state
		OptimalityScore:     analysis.MapControlScore * 0.8, // Optimality of placements
		
		FavoriteZones: []ZonePlacementData{
			{Zone: "River Brushes", PlacementCount: 45, PlacementPercentage: 25.0, Effectiveness: 85.0},
			{Zone: "Jungle Entrances", PlacementCount: 35, PlacementPercentage: 19.0, Effectiveness: 80.0},
			{Zone: "Dragon Area", PlacementCount: 30, PlacementPercentage: 16.0, Effectiveness: 90.0},
		},
		
		AvoidedZones: []ZonePlacementData{
			{Zone: "Deep Enemy Jungle", PlacementCount: 5, PlacementPercentage: 3.0, Effectiveness: 40.0},
			{Zone: "Lane Brushes", PlacementCount: 8, PlacementPercentage: 4.0, Effectiveness: 45.0},
		},
		
		UnderUtilizedZones: []ZonePlacementData{
			{Zone: "Baron Pit", PlacementCount: 12, PlacementPercentage: 7.0, Effectiveness: 95.0},
			{Zone: "Scuttle Crab", PlacementCount: 15, PlacementPercentage: 8.0, Effectiveness: 75.0},
		},
	}
	
	// Calculate optimal placements
	totalPlacements := int(analysis.AverageWardsPlaced * float64(len(matches)))
	analysis.OptimalPlacements = OptimalPlacementData{
		OptimalPlacements:    int(float64(totalPlacements) * 0.65),
		SuboptimalPlacements: int(float64(totalPlacements) * 0.25),
		WastedPlacements:     int(float64(totalPlacements) * 0.10),
		OptimalityRate:       65.0,
		
		TimingOptimization:   15.0, // % improvement possible through timing
		LocationOptimization: 20.0, // % improvement possible through location
		TypeOptimization:     10.0, // % improvement possible through ward type
	}
	
	// Calculate placement timing
	analysis.PlacementTiming = PlacementTimingData{
		PreObjectiveWarding: 70.0, // % wards placed before objectives
		ReactiveWarding:     40.0, // % wards placed reactively
		ProactiveWarding:    60.0, // % wards placed proactively
		EmergencyWarding:    15.0, // % wards placed in emergency
		
		ExcellentTiming: int(float64(totalPlacements) * 0.30),
		GoodTiming:      int(float64(totalPlacements) * 0.45),
		PoorTiming:      int(float64(totalPlacements) * 0.25),
		TimingScore:     72.0,
	}
}

// analyzeWardTypes analyzes performance by ward type
func (was *WardAnalyticsService) analyzeWardTypes(analysis *WardAnalysis, matches []models.MatchData) {
	// Yellow Wards Analysis
	analysis.YellowWardsAnalysis = WardTypeAnalysis{
		WardType:              "YELLOW",
		AveragePlaced:         analysis.AverageWardsPlaced * 0.7, // ~70% of wards are yellow
		PlacementEfficiency:   75.0,
		SurvivalRate:          analysis.WardSurvivalRate * 0.9,
		StrategicUsage:        60.0,
		TimingQuality:         70.0,
		ImpactScore:           72.0,
		EarlyGameUsage:        40.0,
		MidGameUsage:          45.0,
		LateGameUsage:        15.0,
		OptimizationPotential: 15.0,
		RecommendedUsage:     "Focus on river control and jungle entrances",
	}
	
	// Control Wards Analysis
	analysis.ControlWardsAnalysis = WardTypeAnalysis{
		WardType:              "CONTROL",
		AveragePlaced:         analysis.AverageWardsPlaced * 0.25, // ~25% are control wards
		PlacementEfficiency:   85.0,
		SurvivalRate:          analysis.WardSurvivalRate * 1.2,
		StrategicUsage:        90.0,
		TimingQuality:         75.0,
		ImpactScore:           88.0,
		EarlyGameUsage:        20.0,
		MidGameUsage:          50.0,
		LateGameUsage:        30.0,
		OptimizationPotential: 10.0,
		RecommendedUsage:     "Prioritize objective control and deep vision denial",
	}
	
	// Blue Ward Analysis
	analysis.BlueWardAnalysis = WardTypeAnalysis{
		WardType:              "BLUE_TRINKET",
		AveragePlaced:         analysis.AverageWardsPlaced * 0.05, // ~5% are blue trinkets
		PlacementEfficiency:   95.0,
		SurvivalRate:          0.0, // Blue wards don't persist
		StrategicUsage:        95.0,
		TimingQuality:         80.0,
		ImpactScore:           75.0,
		EarlyGameUsage:        5.0,
		MidGameUsage:          35.0,
		LateGameUsage:        60.0,
		OptimizationPotential: 25.0,
		RecommendedUsage:     "Use for safe objective scouting and baron/elder setups",
	}
}

// analyzeClearingPatterns analyzes ward clearing patterns
func (was *WardAnalyticsService) analyzeClearingPatterns(analysis *WardAnalysis, matches []models.MatchData) {
	totalWardsCleared := analysis.AverageWardsKilled * float64(len(matches))
	
	analysis.WardClearingPatterns = WardClearingData{
		TotalWardsCleared:     int(totalWardsCleared),
		ClearingEfficiency:    analysis.WardEfficiency * 100,
		StrategicClearing:     75.0, // % cleared from strategic areas
		TimingQuality:         70.0,
		SafetyScore:           65.0,
		
		ProactiveClearing:     40.0,
		ReactiveClearing:      45.0,
		OpportunisticClearing: 15.0,
		
		RiverClearing:         80.0,
		JungleClearing:        60.0,
		ObjectiveClearing:     90.0,
	}
	
	// Calculate counter warding score
	baseCounterScore := 50.0
	if analysis.WardEfficiency > 0.5 {
		baseCounterScore += 25
	} else if analysis.WardEfficiency > 0.3 {
		baseCounterScore += 15
	} else if analysis.WardEfficiency < 0.2 {
		baseCounterScore -= 20
	}
	
	analysis.CounterWardingScore = math.Max(0, math.Min(100, baseCounterScore))
}

// analyzeWardPhases analyzes ward performance by game phase
func (was *WardAnalyticsService) analyzeWardPhases(analysis *WardAnalysis, matches []models.MatchData) {
	analysis.EarlyGameWards = WardPhaseData{
		Phase:            "early",
		WardsPlaced:      analysis.AverageWardsPlaced * 0.35,
		WardsKilled:      analysis.AverageWardsKilled * 0.25,
		PlacementQuality: 70.0,
		StrategicFocus:   60.0,
		MapCoverage:      40.0,
		EfficiencyRating: was.rateEfficiency(70.0),
	}
	
	analysis.MidGameWards = WardPhaseData{
		Phase:            "mid",
		WardsPlaced:      analysis.AverageWardsPlaced * 0.45,
		WardsKilled:      analysis.AverageWardsKilled * 0.50,
		PlacementQuality: 80.0,
		StrategicFocus:   85.0,
		MapCoverage:      70.0,
		EfficiencyRating: was.rateEfficiency(80.0),
	}
	
	analysis.LateGameWards = WardPhaseData{
		Phase:            "late",
		WardsPlaced:      analysis.AverageWardsPlaced * 0.20,
		WardsKilled:      analysis.AverageWardsKilled * 0.25,
		PlacementQuality: 90.0,
		StrategicFocus:   95.0,
		MapCoverage:      60.0,
		EfficiencyRating: was.rateEfficiency(90.0),
	}
}

// analyzeZoneControl analyzes control of specific zones
func (was *WardAnalyticsService) analyzeZoneControl(analysis *WardAnalysis, matches []models.MatchData) {
	analysis.ZoneControl = map[string]ZoneControlData{
		"Dragon Pit": {
			ZoneName:          "Dragon Pit",
			ControlPercentage: 75.0,
			ContestLevel:      "high",
			StrategicValue:    95.0,
			WardsPlaced:       30,
			WardsCleared:      12,
			ControlEfficiency: 85.0,
			RecommendedFocus:  "Maintain consistent vision before all dragon spawns",
		},
		"Baron Pit": {
			ZoneName:          "Baron Pit",
			ControlPercentage: 60.0,
			ContestLevel:      "high",
			StrategicValue:    98.0,
			WardsPlaced:       18,
			WardsCleared:      15,
			ControlEfficiency: 70.0,
			RecommendedFocus:  "Improve pre-baron vision setup and maintenance",
		},
		"River": {
			ZoneName:          "River",
			ControlPercentage: 85.0,
			ContestLevel:      "medium",
			StrategicValue:    80.0,
			WardsPlaced:       45,
			WardsCleared:      20,
			ControlEfficiency: 90.0,
			RecommendedFocus:  "Excellent river control - maintain current approach",
		},
	}
	
	// River control analysis
	analysis.RiverControl = RiverControlData{
		OverallRiverControl: 85.0,
		TopRiverControl:     80.0,
		BottomRiverControl:  90.0,
		ScuttleCrabControl:  75.0,
		RiverBrushControl:   95.0,
		CrossingPoints:      70.0,
		RiverScore:          85.0,
	}
	
	// Jungle control analysis
	analysis.JungleControl = JungleControlData{
		OwnJungleControl:    70.0,
		EnemyJungleControl:  35.0,
		BuffControl:        85.0,
		CampTimerTracking:  60.0,
		InvasionDetection:  75.0,
		CounterJungling:    45.0,
		JungleScore:        65.0,
	}
}

// analyzeWardImpact analyzes the impact of warding
func (was *WardAnalyticsService) analyzeWardImpact(analysis *WardAnalysis, matches []models.MatchData) {
	// Strategic coverage analysis
	analysis.StrategicCoverage = StrategicCoverageData{
		DragonPitCoverage:     75.0,
		BaronPitCoverage:      60.0,
		RiverBrushesCoverage:  90.0,
		JungleEntryCoverage:   70.0,
		LaneBrushesCoverage:   50.0,
		OverallStrategicScore: 75.0,
	}
	
	// Vision denial score
	analysis.VisionDeniedScore = analysis.CounterWardingScore * 0.8
	
	// Objective setup analysis
	analysis.ObjectiveSetup = ObjectiveSetupData{
		DragonSetupScore:    80.0,
		BaronSetupScore:     65.0,
		HeraldSetupScore:    70.0,
		ElderSetupScore:     75.0,
		PreObjectiveSetup:   70.0,
		SetupTiming:         45.0, // seconds before objective
		SetupEfficiency:     75.0,
		ApproachCoverage:    80.0,
		EscapeCoverage:      60.0,
		FlankCoverage:       70.0,
	}
	
	// Safety provided analysis
	analysis.SafetyProvided = SafetyData{
		GankPrevention:      75.0,
		InvasionDetection:   70.0,
		RotationTracking:    65.0,
		SafeFarmingProvided: 25.0, // % increase
		SafetyScore:         70.0,
		HighRiskDetection:   80.0,
		MediumRiskDetection: 65.0,
		RiskMitigationScore: 72.0,
	}
	
	// Calculate overall ward impact score
	baseImpact := 50.0
	baseImpact += (analysis.MapControlScore - 50) * 0.3
	baseImpact += (analysis.StrategicCoverage.OverallStrategicScore - 50) * 0.3
	baseImpact += (analysis.SafetyProvided.SafetyScore - 50) * 0.2
	baseImpact += (analysis.CounterWardingScore - 50) * 0.2
	
	analysis.WardImpactScore = math.Max(0, math.Min(100, baseImpact))
}

// performWardComparison compares against benchmarks
func (was *WardAnalyticsService) performWardComparison(ctx context.Context, analysis *WardAnalysis, position string) error {
	// Get role benchmark
	roleBenchmark, err := was.getWardBenchmark(ctx, "role", position, "")
	if err == nil {
		analysis.RoleBenchmark = *roleBenchmark
		analysis.RoleBenchmark.PlayerPercentile = was.calculatePercentile(analysis.MapControlScore, roleBenchmark.AverageMapControl)
	}

	// Get rank benchmark
	rankBenchmark, err := was.getWardBenchmark(ctx, "rank", "GOLD", "")
	if err == nil {
		analysis.RankBenchmark = *rankBenchmark
		analysis.RankBenchmark.PlayerPercentile = was.calculatePercentile(analysis.MapControlScore, rankBenchmark.AverageMapControl)
	}

	// Get global benchmark
	globalBenchmark, err := was.getWardBenchmark(ctx, "global", "", "")
	if err == nil {
		analysis.GlobalBenchmark = *globalBenchmark
		analysis.GlobalBenchmark.PlayerPercentile = was.calculatePercentile(analysis.MapControlScore, globalBenchmark.AverageMapControl)
	}

	return nil
}

// analyzeWardTrend analyzes ward performance trends
func (was *WardAnalyticsService) analyzeWardTrend(analysis *WardAnalysis, matches []models.MatchData) {
	if len(matches) < 5 {
		analysis.TrendDirection = "insufficient_data"
		return
	}

	// Sort matches by date
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Date.Before(matches[j].Date)
	})

	// Extract map control scores (estimated)
	controlScores := make([]float64, len(matches))
	for i, match := range matches {
		// Estimate control score based on wards placed/killed
		controlScore := 50.0
		if match.WardsPlaced > 10 {
			controlScore += 20
		}
		if match.WardsKilled > 5 {
			controlScore += 15
		}
		if match.VisionScore > 30 {
			controlScore += 15
		}
		controlScores[i] = math.Max(0, math.Min(100, controlScore))
	}

	// Calculate trend using linear regression
	slope, confidence := was.calculateLinearRegression(controlScores)

	analysis.TrendSlope = slope
	analysis.TrendConfidence = confidence

	// Determine trend direction
	if slope > 3.0 && confidence > 0.6 {
		analysis.TrendDirection = "improving"
	} else if slope < -3.0 && confidence > 0.6 {
		analysis.TrendDirection = "declining"
	} else {
		analysis.TrendDirection = "stable"
	}
}

// generateOptimizationSuggestions creates optimization advice
func (was *WardAnalyticsService) generateOptimizationSuggestions(analysis *WardAnalysis) {
	// Placement optimization
	analysis.PlacementOptimization = PlacementOptimizationData{
		OptimalSpots: []OptimalSpotData{
			{
				Zone:            "Dragon Pit Entrance",
				Coordinates:     [2]int{10000, 4400},
				StrategicValue:  95.0,
				CurrentUsage:    60.0,
				RecommendedUsage: 85.0,
				Reason:          "Critical for dragon control and team fight positioning",
			},
			{
				Zone:            "Baron Pit Bush",
				Coordinates:     [2]int{5000, 10400},
				StrategicValue:  98.0,
				CurrentUsage:    40.0,
				RecommendedUsage: 90.0,
				Reason:          "Essential for baron control and late game vision",
			},
		},
		
		TimingImprovements: []TimingImprovementData{
			{
				Situation:         "Pre-Dragon Spawn",
				CurrentTiming:     20.0, // seconds before
				OptimalTiming:     45.0,
				ImprovementNeeded: 25.0,
				Tips:              []string{"Ward 45 seconds before dragon spawns", "Clear enemy vision first", "Coordinate with team"},
			},
		},
		
		TypeOptimizations: []TypeOptimizationData{
			{
				Situation:       "Baron Area Control",
				CurrentType:     "YELLOW",
				RecommendedType: "CONTROL",
				EfficiencyGain:  25.0,
				Reasoning:       "Control wards provide permanent vision and deny enemy wards in this critical area",
			},
		},
		
		ExpectedControlGain: 15.0,
		ExpectedSafetyGain:  20.0,
		ImplementationTips: []string{
			"Practice ward placement timing in training mode",
			"Use F-keys to check ward coverage regularly",
			"Coordinate ward placement with team rotations",
		},
	}
	
	// Clearing optimization
	analysis.ClearingOptimization = ClearingOptimizationData{
		PriorityTargets: []PriorityTargetData{
			{
				Zone:              "Dragon Pit",
				Priority:          "high",
				StrategicValue:    95.0,
				ClearingFrequency: 60.0,
				RecommendedFocus:  "Always clear before objectives",
			},
			{
				Zone:              "River Bushes",
				Priority:          "medium",
				StrategicValue:    75.0,
				ClearingFrequency: 80.0,
				RecommendedFocus:  "Good current focus, maintain consistency",
			},
		},
		
		ClearingOpportunities: []ClearingOpportunityData{
			{
				Situation:       "Post Team Fight Victory",
				OpportunityType: "safe_clearing",
				SafetyLevel:     "high",
				ExpectedReward:  85.0,
				Tips:            []string{"Clear all visible wards", "Push for deeper vision", "Coordinate with team"},
			},
		},
		
		SafetyClearingTips: []string{
			"Always clear with team support in dangerous areas",
			"Use sweepers before entering unwarded areas",
			"Prioritize escape routes when clearing deep wards",
		},
		
		ExpectedDenialGain: 20.0,
		ExpectedSafetyGain: 15.0,
	}
}

// generateWardRecommendations creates actionable recommendations
func (was *WardAnalyticsService) generateWardRecommendations(analysis *WardAnalysis) {
	recommendations := []WardRecommendation{}

	// Check map control vs benchmark
	if analysis.RoleBenchmark.PlayerPercentile < 50 {
		recommendations = append(recommendations, WardRecommendation{
			Priority:                 "high",
			Category:                 "placement",
			Title:                    "Improve Map Control",
			Description:              fmt.Sprintf("Your map control score (%.0f) is below role average. Focus on consistent ward placement in strategic areas.", analysis.MapControlScore),
			Impact:                   "Better map control increases team fight success by 15-20%",
			GamePhase:                []string{"mid", "late"},
			ExpectedImprovement:      15.0,
			ImplementationDifficulty: "medium",
		})
	}

	// Check ward efficiency
	if analysis.WardEfficiency < 0.4 {
		recommendations = append(recommendations, WardRecommendation{
			Priority:                 "high",
			Category:                 "clearing",
			Title:                    "Increase Ward Clearing",
			Description:              fmt.Sprintf("Your ward efficiency (%.2f) suggests room for improvement in clearing enemy wards. Focus on denying enemy vision.", analysis.WardEfficiency),
			Impact:                   "Better ward clearing reduces enemy map awareness significantly",
			GamePhase:                []string{"mid", "late"},
			ExpectedImprovement:      20.0,
			ImplementationDifficulty: "easy",
		})
	}

	// Check strategic coverage
	if analysis.StrategicCoverage.OverallStrategicScore < 70 {
		recommendations = append(recommendations, WardRecommendation{
			Priority:                 "medium",
			Category:                 "positioning",
			Title:                    "Focus on Strategic Areas",
			Description:              "Your strategic area coverage can be improved. Prioritize dragon pit, baron pit, and river control.",
			Impact:                   "Strategic vision control is crucial for objective control",
			GamePhase:                []string{"mid", "late"},
			ExpectedImprovement:      12.0,
			ImplementationDifficulty: "medium",
		})
	}

	// Check placement timing
	if analysis.PlacementTiming.TimingScore < 70 {
		recommendations = append(recommendations, WardRecommendation{
			Priority:                 "medium",
			Category:                 "timing",
			Title:                    "Improve Ward Timing",
			Description:              "Your ward placement timing needs improvement. Place wards proactively before objectives and rotations.",
			Impact:                   "Better timing provides crucial early information",
			GamePhase:                []string{"early", "mid", "late"},
			ExpectedImprovement:      10.0,
			ImplementationDifficulty: "hard",
		})
	}

	analysis.Recommendations = recommendations

	// Generate strength and improvement areas
	analysis.StrengthAreas = []string{}
	analysis.ImprovementAreas = []string{}

	if analysis.RoleBenchmark.PlayerPercentile > 75 {
		analysis.StrengthAreas = append(analysis.StrengthAreas, "Map Control")
	} else if analysis.RoleBenchmark.PlayerPercentile < 25 {
		analysis.ImprovementAreas = append(analysis.ImprovementAreas, "Map Control")
	}

	if analysis.WardEfficiency > 0.6 {
		analysis.StrengthAreas = append(analysis.StrengthAreas, "Ward Clearing")
	} else if analysis.WardEfficiency < 0.3 {
		analysis.ImprovementAreas = append(analysis.ImprovementAreas, "Ward Clearing")
	}

	if analysis.StrategicCoverage.OverallStrategicScore > 80 {
		analysis.StrengthAreas = append(analysis.StrengthAreas, "Strategic Vision")
	} else if analysis.StrategicCoverage.OverallStrategicScore < 60 {
		analysis.ImprovementAreas = append(analysis.ImprovementAreas, "Strategic Vision")
	}
}

// generateWardTrendData creates trend visualization data
func (was *WardAnalyticsService) generateWardTrendData(analysis *WardAnalysis, matches []models.MatchData) {
	analysis.TrendData = []WardTrendPoint{}

	// Sort matches by date
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Date.Before(matches[j].Date)
	})

	// Create trend points with moving averages
	windowSize := 5
	for i := range matches {
		point := WardTrendPoint{
			Date:            matches[i].Date,
			WardsPlaced:     float64(matches[i].WardsPlaced),
			WardsKilled:     float64(matches[i].WardsKilled),
			MapControlScore: was.calculateMatchMapControl(matches[i]),
			WardEfficiency:  float64(matches[i].WardsKilled) / math.Max(1, float64(matches[i].WardsPlaced)),
		}

		// Calculate moving average
		if i >= windowSize-1 {
			sum := 0.0
			for j := i - windowSize + 1; j <= i; j++ {
				sum += was.calculateMatchMapControl(matches[j])
			}
			point.MovingAverage = sum / float64(windowSize)
		} else {
			point.MovingAverage = point.MapControlScore
		}

		analysis.TrendData = append(analysis.TrendData, point)
	}
}

// Helper functions

func (was *WardAnalyticsService) calculateMatchMapControl(match models.MatchData) float64 {
	// Estimate map control score for a single match
	baseScore := 50.0
	
	// Factor in wards placed
	if match.WardsPlaced > 15 {
		baseScore += 20
	} else if match.WardsPlaced > 10 {
		baseScore += 10
	}
	
	// Factor in wards killed
	if match.WardsKilled > 8 {
		baseScore += 15
	} else if match.WardsKilled > 5 {
		baseScore += 8
	}
	
	// Factor in vision score
	if match.VisionScore > 40 {
		baseScore += 15
	} else if match.VisionScore > 25 {
		baseScore += 10
	}
	
	return math.Max(0, math.Min(100, baseScore))
}

func (was *WardAnalyticsService) calculateLinearRegression(values []float64) (slope, confidence float64) {
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

func (was *WardAnalyticsService) calculatePercentile(playerValue, benchmarkAverage float64) float64 {
	if benchmarkAverage == 0 {
		return 50.0
	}
	ratio := playerValue / benchmarkAverage
	if ratio >= 1.4 {
		return 90.0
	} else if ratio >= 1.2 {
		return 75.0
	} else if ratio >= 1.0 {
		return 60.0
	} else if ratio >= 0.8 {
		return 40.0
	} else if ratio >= 0.6 {
		return 25.0
	} else {
		return 10.0
	}
}

func (was *WardAnalyticsService) rateEfficiency(score float64) string {
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

// Placeholder functions for data access
func (was *WardAnalyticsService) getWardMatches(ctx context.Context, playerID string, timeRange string, champion string, position string) ([]models.MatchData, error) {
	// This would query the database for matches with ward data
	return []models.MatchData{}, nil
}

func (was *WardAnalyticsService) getWardBenchmark(ctx context.Context, category, filter, champion string) (*WardBenchmark, error) {
	// This would query benchmark data from database
	return &WardBenchmark{
		Category:           category,
		AverageWardsPlaced: 12.0,
		AverageWardsKilled: 6.0,
		AverageMapControl:  70.0,
		Top10Percent:       85.0,
		Top25Percent:       78.0,
		Median:             65.0,
	}, nil
}