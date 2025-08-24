package services

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/herald-lol/herald/backend/internal/models"
)

// VisionAnalyticsService handles vision-related analytics and heatmap generation
type VisionAnalyticsService struct {
	analyticsService *AnalyticsService
	mapService       *MapService
}

// VisionAnalysis represents comprehensive vision analysis
type VisionAnalysis struct {
	PlayerID  string `json:"player_id"`
	Champion  string `json:"champion,omitempty"`
	Position  string `json:"position,omitempty"`
	TimeRange string `json:"time_range"`

	// Core Vision Metrics
	AverageVisionScore float64 `json:"average_vision_score"`
	MedianVisionScore  float64 `json:"median_vision_score"`
	BestVisionScore    float64 `json:"best_vision_score"`
	WorstVisionScore   float64 `json:"worst_vision_score"`
	VisionScoreStdDev  float64 `json:"vision_score_std_dev"`

	// Ward Statistics
	AverageWardsPlaced  float64 `json:"average_wards_placed"`
	AverageWardsKilled  float64 `json:"average_wards_killed"`
	ControlWardsPerGame float64 `json:"control_wards_per_game"`
	WardEfficiency      float64 `json:"ward_efficiency"` // Wards killed / Wards placed

	// Vision Control Metrics
	EarlyGameVision VisionPhaseData `json:"early_game_vision"`
	MidGameVision   VisionPhaseData `json:"mid_game_vision"`
	LateGameVision  VisionPhaseData `json:"late_game_vision"`

	// Comparative Analysis
	RoleBenchmark   VisionBenchmark `json:"role_benchmark"`
	RankBenchmark   VisionBenchmark `json:"rank_benchmark"`
	GlobalBenchmark VisionBenchmark `json:"global_benchmark"`

	// Performance Analysis
	VisionImpactScore float64       `json:"vision_impact_score"` // 0-100
	VisionRanking     VisionRanking `json:"vision_ranking"`

	// Trend Analysis
	TrendDirection  string  `json:"trend_direction"`
	TrendSlope      float64 `json:"trend_slope"`
	TrendConfidence float64 `json:"trend_confidence"`

	// Heatmap Data
	WardHeatmaps VisionHeatmaps `json:"ward_heatmaps"`

	// Insights and Recommendations
	StrengthAreas    []string               `json:"strength_areas"`
	ImprovementAreas []string               `json:"improvement_areas"`
	Recommendations  []VisionRecommendation `json:"recommendations"`

	// Match Performance
	RecentMatches []MatchVisionData  `json:"recent_matches"`
	TrendData     []VisionTrendPoint `json:"trend_data"`
}

// VisionPhaseData represents vision metrics for a game phase
type VisionPhaseData struct {
	Phase              string  `json:"phase"` // "early", "mid", "late"
	AverageVisionScore float64 `json:"average_vision_score"`
	AverageWardsPlaced float64 `json:"average_wards_placed"`
	AverageWardsKilled float64 `json:"average_wards_killed"`
	ControlWardUsage   float64 `json:"control_ward_usage"`
	EfficiencyRating   string  `json:"efficiency_rating"` // "excellent", "good", "average", "poor"
}

// VisionBenchmark represents comparative vision benchmarks
type VisionBenchmark struct {
	Category           string  `json:"category"` // "role", "rank", "global"
	AverageVisionScore float64 `json:"average_vision_score"`
	Top10Percent       float64 `json:"top_10_percent"`
	Top25Percent       float64 `json:"top_25_percent"`
	Median             float64 `json:"median"`
	PlayerPercentile   float64 `json:"player_percentile"`
}

// VisionRanking represents player's vision ranking
type VisionRanking struct {
	OverallRank       int     `json:"overall_rank"`
	RoleRank          int     `json:"role_rank"`
	TierRank          int     `json:"tier_rank"`
	PercentileOverall float64 `json:"percentile_overall"`
	PercentileInRole  float64 `json:"percentile_in_role"`
	PercentileInTier  float64 `json:"percentile_in_tier"`
}

// VisionHeatmaps contains heatmap data for different ward types
type VisionHeatmaps struct {
	YellowWards    HeatmapData `json:"yellow_wards"`
	ControlWards   HeatmapData `json:"control_wards"`
	WardKills      HeatmapData `json:"ward_kills"`
	DeathLocations HeatmapData `json:"death_locations"`
	VisionDenied   HeatmapData `json:"vision_denied"`
}

// HeatmapData represents spatial heatmap information
type HeatmapData struct {
	MapSide    string         `json:"map_side"` // "blue", "red", "both"
	DataPoints []HeatmapPoint `json:"data_points"`
	Intensity  map[string]int `json:"intensity"` // Zone -> frequency
	Coverage   float64        `json:"coverage_percent"`
}

// HeatmapPoint represents a point on the heatmap
type HeatmapPoint struct {
	X         int     `json:"x"`
	Y         int     `json:"y"`
	Frequency int     `json:"frequency"`
	Weight    float64 `json:"weight"`
	Zone      string  `json:"zone"` // "jungle", "river", "lane", "baron", "dragon"
}

// VisionRecommendation represents actionable vision advice
type VisionRecommendation struct {
	Priority    string   `json:"priority"` // "high", "medium", "low"
	Category    string   `json:"category"` // "warding", "dewarding", "positioning", "timing"
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Impact      string   `json:"impact"`
	GamePhase   []string `json:"game_phase"`           // ["early", "mid", "late"]
	VisualAid   string   `json:"visual_aid,omitempty"` // URL to visual guide
}

// MatchVisionData represents vision data from a specific match
type MatchVisionData struct {
	MatchID          string    `json:"match_id"`
	Champion         string    `json:"champion"`
	Position         string    `json:"position"`
	VisionScore      int       `json:"vision_score"`
	WardsPlaced      int       `json:"wards_placed"`
	WardsKilled      int       `json:"wards_killed"`
	ControlWards     int       `json:"control_wards"`
	GameDuration     int       `json:"game_duration"`
	Result           string    `json:"result"`
	Date             time.Time `json:"date"`
	VisionEfficiency float64   `json:"vision_efficiency"`
	MapControl       float64   `json:"map_control_percent"`
}

// VisionTrendPoint represents vision performance over time
type VisionTrendPoint struct {
	Date          time.Time `json:"date"`
	VisionScore   float64   `json:"vision_score"`
	WardsPlaced   float64   `json:"wards_placed"`
	WardsKilled   float64   `json:"wards_killed"`
	MovingAverage float64   `json:"moving_average"`
	Efficiency    float64   `json:"efficiency"`
}

// Ward placement zones for analysis
type WardZone struct {
	Name        string   `json:"name"`
	Coordinates [][]int  `json:"coordinates"` // Polygon coordinates
	Strategic   bool     `json:"strategic"`   // High-value zone
	GamePhase   []string `json:"game_phase"`  // When this zone is most important
}

// Map zones definition for League of Legends
var MapZones = []WardZone{
	{
		Name:        "Dragon Pit",
		Coordinates: [][]int{{9800, 4200}, {10200, 4200}, {10200, 4600}, {9800, 4600}},
		Strategic:   true,
		GamePhase:   []string{"mid", "late"},
	},
	{
		Name:        "Baron Pit",
		Coordinates: [][]int{{4800, 10200}, {5200, 10200}, {5200, 10600}, {4800, 10600}},
		Strategic:   true,
		GamePhase:   []string{"late"},
	},
	{
		Name:        "Blue Side Blue Buff",
		Coordinates: [][]int{{3800, 8000}, {4200, 8000}, {4200, 8400}, {3800, 8400}},
		Strategic:   true,
		GamePhase:   []string{"early", "mid"},
	},
	{
		Name:        "Red Side Red Buff",
		Coordinates: [][]int{{10800, 6600}, {11200, 6600}, {11200, 7000}, {10800, 7000}},
		Strategic:   true,
		GamePhase:   []string{"early", "mid"},
	},
	{
		Name:        "River Bushes",
		Coordinates: [][]int{{6000, 6000}, {9000, 6000}, {9000, 9000}, {6000, 9000}},
		Strategic:   true,
		GamePhase:   []string{"early", "mid", "late"},
	},
	// Add more zones...
}

// NewVisionAnalyticsService creates a new vision analytics service
func NewVisionAnalyticsService(analyticsService *AnalyticsService, mapService *MapService) *VisionAnalyticsService {
	return &VisionAnalyticsService{
		analyticsService: analyticsService,
		mapService:       mapService,
	}
}

// AnalyzeVision performs comprehensive vision analysis
func (vas *VisionAnalyticsService) AnalyzeVision(ctx context.Context, playerID string, timeRange string, champion string, position string) (*VisionAnalysis, error) {
	// Get match data with vision information
	matches, err := vas.getVisionMatches(ctx, playerID, timeRange, champion, position)
	if err != nil {
		return nil, fmt.Errorf("failed to get vision match data: %w", err)
	}

	if len(matches) == 0 {
		return &VisionAnalysis{
			PlayerID:  playerID,
			Champion:  champion,
			Position:  position,
			TimeRange: timeRange,
		}, nil
	}

	analysis := &VisionAnalysis{
		PlayerID:  playerID,
		Champion:  champion,
		Position:  position,
		TimeRange: timeRange,
	}

	// Calculate basic vision statistics
	vas.calculateVisionBasics(analysis, matches)

	// Analyze vision by game phases
	vas.analyzeVisionPhases(analysis, matches)

	// Perform comparative analysis
	err = vas.performVisionComparison(ctx, analysis, position)
	if err != nil {
		fmt.Printf("Warning: failed to perform vision comparison: %v", err)
	}

	// Calculate vision impact score
	vas.calculateVisionImpactScore(analysis)

	// Perform trend analysis
	vas.analyzeVisionTrend(analysis, matches)

	// Generate heatmaps
	err = vas.generateVisionHeatmaps(ctx, analysis, matches)
	if err != nil {
		fmt.Printf("Warning: failed to generate vision heatmaps: %v", err)
	}

	// Generate insights and recommendations
	vas.generateVisionRecommendations(analysis)

	// Prepare trend data for visualization
	vas.generateVisionTrendData(analysis, matches)

	// Cache results
	vas.cacheVisionAnalysis(ctx, analysis)

	return analysis, nil
}

// GenerateVisionHeatmap creates detailed heatmap for ward placements
func (vas *VisionAnalyticsService) GenerateVisionHeatmap(ctx context.Context, playerID string, timeRange string, wardType string) (*HeatmapData, error) {
	// Get detailed ward placement data
	wardData, err := vas.getWardPlacementData(ctx, playerID, timeRange, wardType)
	if err != nil {
		return nil, fmt.Errorf("failed to get ward placement data: %w", err)
	}

	heatmap := &HeatmapData{
		MapSide:    "both",
		DataPoints: []HeatmapPoint{},
		Intensity:  make(map[string]int),
	}

	// Process ward placements into heatmap points
	for _, ward := range wardData {
		point := HeatmapPoint{
			X:         ward.X,
			Y:         ward.Y,
			Frequency: 1,
			Weight:    vas.calculateWardWeight(ward),
			Zone:      vas.identifyZone(ward.X, ward.Y),
		}
		heatmap.DataPoints = append(heatmap.DataPoints, point)
		heatmap.Intensity[point.Zone]++
	}

	// Aggregate nearby points
	heatmap.DataPoints = vas.aggregateHeatmapPoints(heatmap.DataPoints)

	// Calculate coverage percentage
	heatmap.Coverage = vas.calculateMapCoverage(heatmap.DataPoints)

	return heatmap, nil
}

// GetVisionRecommendations provides personalized vision improvement tips
func (vas *VisionAnalyticsService) GetVisionRecommendations(ctx context.Context, analysis *VisionAnalysis) []VisionRecommendation {
	recommendations := []VisionRecommendation{}

	// Analyze vision score vs role average
	if analysis.AverageVisionScore < analysis.RoleBenchmark.AverageVisionScore*0.8 {
		recommendations = append(recommendations, VisionRecommendation{
			Priority:    "high",
			Category:    "warding",
			Title:       "Increase Ward Placement",
			Description: fmt.Sprintf("Your average vision score (%.1f) is significantly below the role average (%.1f). Focus on placing more wards throughout the game.", analysis.AverageVisionScore, analysis.RoleBenchmark.AverageVisionScore),
			Impact:      "Improving vision score by 20% can increase win rate by 8-12%",
			GamePhase:   []string{"early", "mid", "late"},
		})
	}

	// Analyze ward efficiency
	if analysis.WardEfficiency < 0.3 {
		recommendations = append(recommendations, VisionRecommendation{
			Priority:    "medium",
			Category:    "dewarding",
			Title:       "Improve Ward Clearing",
			Description: fmt.Sprintf("Your ward efficiency (%.2f) suggests you're not clearing enough enemy wards. Look for opportunities to deny enemy vision.", analysis.WardEfficiency),
			Impact:      "Better dewarding reduces enemy map control and creates pick opportunities",
			GamePhase:   []string{"mid", "late"},
		})
	}

	// Analyze control ward usage
	if analysis.ControlWardsPerGame < 3.0 {
		recommendations = append(recommendations, VisionRecommendation{
			Priority:    "medium",
			Category:    "warding",
			Title:       "Use More Control Wards",
			Description: fmt.Sprintf("You average %.1f control wards per game. Aim for 4-6 control wards per game for better vision control.", analysis.ControlWardsPerGame),
			Impact:      "Control wards provide permanent vision and deny enemy wards in key areas",
			GamePhase:   []string{"mid", "late"},
		})
	}

	// Phase-specific recommendations
	if analysis.EarlyGameVision.AverageVisionScore < analysis.RoleBenchmark.AverageVisionScore*0.7 {
		recommendations = append(recommendations, VisionRecommendation{
			Priority:    "high",
			Category:    "timing",
			Title:       "Improve Early Game Warding",
			Description: "Your early game vision is weak. Focus on river wards and jungle entrances during laning phase.",
			Impact:      "Early vision prevents ganks and enables aggressive plays",
			GamePhase:   []string{"early"},
		})
	}

	// Sort recommendations by priority
	sort.Slice(recommendations, func(i, j int) bool {
		priorityOrder := map[string]int{"high": 3, "medium": 2, "low": 1}
		return priorityOrder[recommendations[i].Priority] > priorityOrder[recommendations[j].Priority]
	})

	return recommendations
}

// Helper functions for vision analysis

func (vas *VisionAnalyticsService) calculateVisionBasics(analysis *VisionAnalysis, matches []models.MatchData) {
	if len(matches) == 0 {
		return
	}

	visionScores := make([]float64, 0, len(matches))
	wardsPlaced := make([]float64, 0, len(matches))
	wardsKilled := make([]float64, 0, len(matches))
	controlWards := make([]float64, 0, len(matches))

	for _, match := range matches {
		visionScores = append(visionScores, float64(match.VisionScore))
		wardsPlaced = append(wardsPlaced, float64(match.WardsPlaced))
		wardsKilled = append(wardsKilled, float64(match.WardsKilled))
		controlWards = append(controlWards, float64(match.ControlWardsPlaced))
	}

	// Calculate statistics
	analysis.AverageVisionScore = vas.calculateMean(visionScores)
	analysis.MedianVisionScore = vas.calculateMedian(visionScores)
	analysis.BestVisionScore = vas.calculateMax(visionScores)
	analysis.WorstVisionScore = vas.calculateMin(visionScores)
	analysis.VisionScoreStdDev = vas.calculateStandardDeviation(visionScores)

	analysis.AverageWardsPlaced = vas.calculateMean(wardsPlaced)
	analysis.AverageWardsKilled = vas.calculateMean(wardsKilled)
	analysis.ControlWardsPerGame = vas.calculateMean(controlWards)

	// Calculate ward efficiency
	if analysis.AverageWardsPlaced > 0 {
		analysis.WardEfficiency = analysis.AverageWardsKilled / analysis.AverageWardsPlaced
	}
}

func (vas *VisionAnalyticsService) analyzeVisionPhases(analysis *VisionAnalysis, matches []models.MatchData) {
	// This would analyze vision performance by game phases
	// For now, we'll use simplified calculations

	analysis.EarlyGameVision = VisionPhaseData{
		Phase:              "early",
		AverageVisionScore: analysis.AverageVisionScore * 0.6, // Early game typically has lower vision scores
		AverageWardsPlaced: analysis.AverageWardsPlaced * 0.4,
		AverageWardsKilled: analysis.AverageWardsKilled * 0.3,
		ControlWardUsage:   analysis.ControlWardsPerGame * 0.2,
		EfficiencyRating:   vas.rateEfficiency(analysis.AverageVisionScore * 0.6),
	}

	analysis.MidGameVision = VisionPhaseData{
		Phase:              "mid",
		AverageVisionScore: analysis.AverageVisionScore * 1.2,
		AverageWardsPlaced: analysis.AverageWardsPlaced * 0.4,
		AverageWardsKilled: analysis.AverageWardsKilled * 0.5,
		ControlWardUsage:   analysis.ControlWardsPerGame * 0.6,
		EfficiencyRating:   vas.rateEfficiency(analysis.AverageVisionScore * 1.2),
	}

	analysis.LateGameVision = VisionPhaseData{
		Phase:              "late",
		AverageVisionScore: analysis.AverageVisionScore * 1.4,
		AverageWardsPlaced: analysis.AverageWardsPlaced * 0.2,
		AverageWardsKilled: analysis.AverageWardsKilled * 0.2,
		ControlWardUsage:   analysis.ControlWardsPerGame * 0.2,
		EfficiencyRating:   vas.rateEfficiency(analysis.AverageVisionScore * 1.4),
	}
}

func (vas *VisionAnalyticsService) performVisionComparison(ctx context.Context, analysis *VisionAnalysis, position string) error {
	// Get role benchmark
	roleBenchmark, err := vas.getVisionBenchmark(ctx, "role", position, "")
	if err == nil {
		analysis.RoleBenchmark = *roleBenchmark
		analysis.RoleBenchmark.PlayerPercentile = vas.calculatePercentile(analysis.AverageVisionScore, roleBenchmark.AverageVisionScore)
	}

	// Get rank benchmark (would need player's rank)
	rankBenchmark, err := vas.getVisionBenchmark(ctx, "rank", "GOLD", "") // Simplified
	if err == nil {
		analysis.RankBenchmark = *rankBenchmark
		analysis.RankBenchmark.PlayerPercentile = vas.calculatePercentile(analysis.AverageVisionScore, rankBenchmark.AverageVisionScore)
	}

	// Get global benchmark
	globalBenchmark, err := vas.getVisionBenchmark(ctx, "global", "", "")
	if err == nil {
		analysis.GlobalBenchmark = *globalBenchmark
		analysis.GlobalBenchmark.PlayerPercentile = vas.calculatePercentile(analysis.AverageVisionScore, globalBenchmark.AverageVisionScore)
	}

	return nil
}

func (vas *VisionAnalyticsService) calculateVisionImpactScore(analysis *VisionAnalysis) {
	// Calculate impact score based on multiple factors
	baseScore := 50.0 // Starting point

	// Factor 1: Vision score vs role average
	if analysis.RoleBenchmark.AverageVisionScore > 0 {
		visionRatio := analysis.AverageVisionScore / analysis.RoleBenchmark.AverageVisionScore
		baseScore += (visionRatio - 1.0) * 25.0 // ±25 points based on role performance
	}

	// Factor 2: Ward efficiency
	efficiencyBonus := (analysis.WardEfficiency - 0.5) * 20.0 // ±10 points
	baseScore += efficiencyBonus

	// Factor 3: Control ward usage
	controlWardBonus := math.Min((analysis.ControlWardsPerGame-2.0)*5.0, 15.0) // Up to +15 points
	baseScore += controlWardBonus

	// Factor 4: Consistency (lower std dev is better)
	if analysis.AverageVisionScore > 0 {
		consistencyBonus := math.Max(10.0-(analysis.VisionScoreStdDev/analysis.AverageVisionScore)*50.0, -10.0)
		baseScore += consistencyBonus
	}

	// Clamp between 0-100
	analysis.VisionImpactScore = math.Max(0.0, math.Min(100.0, baseScore))
}

func (vas *VisionAnalyticsService) analyzeVisionTrend(analysis *VisionAnalysis, matches []models.MatchData) {
	if len(matches) < 5 {
		analysis.TrendDirection = "insufficient_data"
		return
	}

	// Sort matches by date
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Date.Before(matches[j].Date)
	})

	// Extract vision scores
	visionScores := make([]float64, len(matches))
	for i, match := range matches {
		visionScores[i] = float64(match.VisionScore)
	}

	// Calculate trend using linear regression
	slope, confidence := vas.calculateLinearRegression(visionScores)

	analysis.TrendSlope = slope
	analysis.TrendConfidence = confidence

	// Determine trend direction
	if slope > 2.0 && confidence > 0.6 {
		analysis.TrendDirection = "improving"
	} else if slope < -2.0 && confidence > 0.6 {
		analysis.TrendDirection = "declining"
	} else {
		analysis.TrendDirection = "stable"
	}
}

func (vas *VisionAnalyticsService) generateVisionHeatmaps(ctx context.Context, analysis *VisionAnalysis, matches []models.MatchData) error {
	// This would generate actual heatmap data from ward placement coordinates
	// For now, we'll create placeholder heatmap structure

	analysis.WardHeatmaps = VisionHeatmaps{
		YellowWards: HeatmapData{
			MapSide:    "both",
			DataPoints: []HeatmapPoint{},
			Coverage:   75.0,
			Intensity:  make(map[string]int),
		},
		ControlWards: HeatmapData{
			MapSide:    "both",
			DataPoints: []HeatmapPoint{},
			Coverage:   45.0,
			Intensity:  make(map[string]int),
		},
		WardKills: HeatmapData{
			MapSide:    "both",
			DataPoints: []HeatmapPoint{},
			Coverage:   30.0,
			Intensity:  make(map[string]int),
		},
	}

	return nil
}

func (vas *VisionAnalyticsService) generateVisionRecommendations(analysis *VisionAnalysis) {
	analysis.Recommendations = vas.GetVisionRecommendations(context.Background(), analysis)

	// Identify strength and improvement areas
	analysis.StrengthAreas = []string{}
	analysis.ImprovementAreas = []string{}

	if analysis.RoleBenchmark.PlayerPercentile > 75 {
		analysis.StrengthAreas = append(analysis.StrengthAreas, "Vision Score")
	} else if analysis.RoleBenchmark.PlayerPercentile < 25 {
		analysis.ImprovementAreas = append(analysis.ImprovementAreas, "Vision Score")
	}

	if analysis.WardEfficiency > 0.4 {
		analysis.StrengthAreas = append(analysis.StrengthAreas, "Ward Clearing")
	} else if analysis.WardEfficiency < 0.2 {
		analysis.ImprovementAreas = append(analysis.ImprovementAreas, "Ward Clearing")
	}

	if analysis.ControlWardsPerGame > 4.0 {
		analysis.StrengthAreas = append(analysis.StrengthAreas, "Control Ward Usage")
	} else if analysis.ControlWardsPerGame < 2.0 {
		analysis.ImprovementAreas = append(analysis.ImprovementAreas, "Control Ward Usage")
	}
}

func (vas *VisionAnalyticsService) generateVisionTrendData(analysis *VisionAnalysis, matches []models.MatchData) {
	analysis.TrendData = []VisionTrendPoint{}

	// Sort matches by date
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Date.Before(matches[j].Date)
	})

	// Create trend points with moving averages
	windowSize := 5
	for i := range matches {
		point := VisionTrendPoint{
			Date:        matches[i].Date,
			VisionScore: float64(matches[i].VisionScore),
			WardsPlaced: float64(matches[i].WardsPlaced),
			WardsKilled: float64(matches[i].WardsKilled),
		}

		// Calculate moving average
		if i >= windowSize-1 {
			sum := 0.0
			for j := i - windowSize + 1; j <= i; j++ {
				sum += float64(matches[j].VisionScore)
			}
			point.MovingAverage = sum / float64(windowSize)
		} else {
			point.MovingAverage = point.VisionScore
		}

		// Calculate efficiency
		if point.WardsPlaced > 0 {
			point.Efficiency = point.WardsKilled / point.WardsPlaced
		}

		analysis.TrendData = append(analysis.TrendData, point)
	}
}

// Utility functions

func (vas *VisionAnalyticsService) calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (vas *VisionAnalyticsService) calculateMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}

func (vas *VisionAnalyticsService) calculateMax(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func (vas *VisionAnalyticsService) calculateMin(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

func (vas *VisionAnalyticsService) calculateStandardDeviation(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	mean := vas.calculateMean(values)
	sumSquaredDiffs := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiffs += diff * diff
	}
	variance := sumSquaredDiffs / float64(len(values))
	return math.Sqrt(variance)
}

func (vas *VisionAnalyticsService) calculateLinearRegression(values []float64) (slope, confidence float64) {
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

func (vas *VisionAnalyticsService) rateEfficiency(score float64) string {
	switch {
	case score >= 45:
		return "excellent"
	case score >= 35:
		return "good"
	case score >= 25:
		return "average"
	default:
		return "poor"
	}
}

func (vas *VisionAnalyticsService) calculatePercentile(playerValue, benchmarkAverage float64) float64 {
	if benchmarkAverage == 0 {
		return 50.0
	}
	// Simplified percentile calculation
	ratio := playerValue / benchmarkAverage
	if ratio >= 1.5 {
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

// Placeholder functions that would be implemented with real data access

func (vas *VisionAnalyticsService) getVisionMatches(ctx context.Context, playerID string, timeRange string, champion string, position string) ([]models.MatchData, error) {
	// This would query the database for matches with vision data
	return []models.MatchData{}, nil
}

func (vas *VisionAnalyticsService) getVisionBenchmark(ctx context.Context, category, filter, champion string) (*VisionBenchmark, error) {
	// This would query benchmark data from database
	return &VisionBenchmark{
		Category:           category,
		AverageVisionScore: 32.0, // Placeholder
		Top10Percent:       50.0,
		Top25Percent:       42.0,
		Median:             30.0,
	}, nil
}

func (vas *VisionAnalyticsService) getWardPlacementData(ctx context.Context, playerID string, timeRange string, wardType string) ([]models.WardPlacement, error) {
	// This would query detailed ward placement data
	return []models.WardPlacement{}, nil
}

func (vas *VisionAnalyticsService) identifyZone(x, y int) string {
	// Identify which zone of the map this coordinate belongs to
	for _, zone := range MapZones {
		if vas.isPointInPolygon(x, y, zone.Coordinates) {
			return zone.Name
		}
	}
	return "unknown"
}

func (vas *VisionAnalyticsService) isPointInPolygon(x, y int, polygon [][]int) bool {
	// Ray casting algorithm to determine if point is in polygon
	n := len(polygon)
	inside := false

	j := n - 1
	for i := 0; i < n; i++ {
		if ((polygon[i][1] > y) != (polygon[j][1] > y)) &&
			(x < (polygon[j][0]-polygon[i][0])*(y-polygon[i][1])/(polygon[j][1]-polygon[i][1])+polygon[i][0]) {
			inside = !inside
		}
		j = i
	}

	return inside
}

func (vas *VisionAnalyticsService) calculateWardWeight(ward models.WardPlacement) float64 {
	// Calculate strategic weight of ward placement
	weight := 1.0

	// Bonus for strategic zones
	zone := vas.identifyZone(ward.X, ward.Y)
	for _, mapZone := range MapZones {
		if mapZone.Name == zone && mapZone.Strategic {
			weight *= 1.5
			break
		}
	}

	return weight
}

func (vas *VisionAnalyticsService) aggregateHeatmapPoints(points []HeatmapPoint) []HeatmapPoint {
	// Aggregate nearby points to reduce noise
	aggregated := []HeatmapPoint{}
	threshold := 200 // Distance threshold in game units

	for _, point := range points {
		merged := false
		for i, existing := range aggregated {
			distance := math.Sqrt(float64((point.X-existing.X)*(point.X-existing.X) + (point.Y-existing.Y)*(point.Y-existing.Y)))
			if distance < float64(threshold) {
				// Merge points
				aggregated[i].Frequency += point.Frequency
				aggregated[i].Weight = (aggregated[i].Weight + point.Weight) / 2
				merged = true
				break
			}
		}
		if !merged {
			aggregated = append(aggregated, point)
		}
	}

	return aggregated
}

func (vas *VisionAnalyticsService) calculateMapCoverage(points []HeatmapPoint) float64 {
	// Calculate what percentage of strategic areas are covered
	coveredZones := make(map[string]bool)

	for _, point := range points {
		if point.Zone != "unknown" {
			coveredZones[point.Zone] = true
		}
	}

	strategicZones := 0
	for _, zone := range MapZones {
		if zone.Strategic {
			strategicZones++
		}
	}

	if strategicZones == 0 {
		return 0
	}

	return float64(len(coveredZones)) / float64(strategicZones) * 100
}

func (vas *VisionAnalyticsService) cacheVisionAnalysis(ctx context.Context, analysis *VisionAnalysis) {
	// Cache analysis results for faster subsequent requests
	// Implementation would use Redis or similar caching system
}
