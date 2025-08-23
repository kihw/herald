package services

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/herald/internal/models"
	"github.com/herald/internal/repository"
)

// AnalyticsService provides comprehensive gaming analytics
type AnalyticsService struct {
	db           *sql.DB
	matchRepo    *repository.MatchRepository
	playerRepo   *repository.PlayerRepository
	redisService *RedisService
}

// KDAAnalysis represents KDA statistical analysis
type KDAAnalysis struct {
	PlayerID     string              `json:"player_id"`
	Champion     string              `json:"champion,omitempty"`
	TimeRange    string              `json:"time_range"`
	
	// Core KDA Metrics
	TotalKills   int                 `json:"total_kills"`
	TotalDeaths  int                 `json:"total_deaths"`
	TotalAssists int                 `json:"total_assists"`
	
	// Statistical Analysis
	AverageKDA   float64             `json:"average_kda"`
	MedianKDA    float64             `json:"median_kda"`
	BestKDA      float64             `json:"best_kda"`
	WorstKDA     float64             `json:"worst_kda"`
	StandardDev  float64             `json:"standard_deviation"`
	
	// Trend Analysis
	TrendDirection string             `json:"trend_direction"` // "improving", "declining", "stable"
	TrendSlope     float64            `json:"trend_slope"`
	TrendConfidence float64           `json:"trend_confidence"`
	
	// Performance Distribution
	KDADistribution map[string]int    `json:"kda_distribution"`
	
	// Comparative Analysis
	RankPercentile  float64           `json:"rank_percentile"`
	GlobalPercentile float64          `json:"global_percentile"`
	
	// Contextual Metrics
	WinRateByKDA    map[string]float64 `json:"winrate_by_kda"`
	KDAByGameLength map[string]float64 `json:"kda_by_game_length"`
	KDAByPosition   map[string]float64 `json:"kda_by_position"`
	
	// Recent Performance
	Last7Days       KDASnapshot        `json:"last_7_days"`
	Last30Days      KDASnapshot        `json:"last_30_days"`
	CurrentStreak   StreakAnalysis     `json:"current_streak"`
	
	// Historical Data
	TrendData       []KDATrendPoint    `json:"trend_data"`
	MatchHistory    []MatchKDAData     `json:"recent_matches"`
}

// KDASnapshot represents KDA performance in a specific timeframe
type KDASnapshot struct {
	Matches      int     `json:"matches"`
	AverageKDA   float64 `json:"average_kda"`
	WinRate      float64 `json:"win_rate"`
	Improvement  float64 `json:"improvement_percent"`
}

// StreakAnalysis represents current performance streak
type StreakAnalysis struct {
	Type         string  `json:"type"` // "positive", "negative", "mixed"
	Length       int     `json:"length"`
	AverageKDA   float64 `json:"average_kda"`
	StartDate    time.Time `json:"start_date"`
}

// KDATrendPoint represents a point in KDA trend over time
type KDATrendPoint struct {
	Date       time.Time `json:"date"`
	KDA        float64   `json:"kda"`
	Matches    int       `json:"matches"`
	MovingAvg  float64   `json:"moving_average"`
}

// MatchKDAData represents KDA data for a specific match
type MatchKDAData struct {
	MatchID      string    `json:"match_id"`
	Champion     string    `json:"champion"`
	Position     string    `json:"position"`
	Kills        int       `json:"kills"`
	Deaths       int       `json:"deaths"`
	Assists      int       `json:"assists"`
	KDA          float64   `json:"kda"`
	GameDuration int       `json:"game_duration_minutes"`
	Result       string    `json:"result"` // "victory", "defeat"
	Date         time.Time `json:"date"`
}

// CSAnalysis represents Creep Score analysis
type CSAnalysis struct {
	PlayerID           string              `json:"player_id"`
	Champion           string              `json:"champion,omitempty"`
	Position           string              `json:"position,omitempty"`
	TimeRange          string              `json:"time_range"`
	
	// Core CS Metrics
	TotalCS            int                 `json:"total_cs"`
	AverageCS          float64             `json:"average_cs"`
	AverageCSPerMin    float64             `json:"average_cs_per_minute"`
	BestCSPerMin       float64             `json:"best_cs_per_minute"`
	MedianCSPerMin     float64             `json:"median_cs_per_minute"`
	
	// Benchmark Comparison
	RoleAverage        float64             `json:"role_average_cs_per_min"`
	RankAverage        float64             `json:"rank_average_cs_per_min"`
	ProAverage         float64             `json:"pro_average_cs_per_min"`
	
	// Performance Analysis
	CSEfficiency       float64             `json:"cs_efficiency"` // % of possible CS obtained
	EarlyGameCS        float64             `json:"early_game_cs"` // 0-15 min avg
	MidGameCS          float64             `json:"mid_game_cs"`   // 15-25 min avg
	LateGameCS         float64             `json:"late_game_cs"`  // 25+ min avg
	
	// Trend Analysis
	TrendDirection     string              `json:"trend_direction"`
	TrendSlope         float64             `json:"trend_slope"`
	ImprovementRate    float64             `json:"improvement_rate"`
	
	// Contextual Analysis
	CSByMatchLength    map[string]float64  `json:"cs_by_match_length"`
	CSByResult         map[string]float64  `json:"cs_by_result"`
	CSByChampion       map[string]float64  `json:"cs_by_champion"`
	
	// Recent Performance
	Last7Days          CSSnapshot          `json:"last_7_days"`
	Last30Days         CSSnapshot          `json:"last_30_days"`
	
	// Historical Trend
	TrendData          []CSTrendPoint      `json:"trend_data"`
	Recommendations    []string            `json:"recommendations"`
}

// CSSnapshot represents CS performance snapshot
type CSSnapshot struct {
	Matches         int     `json:"matches"`
	AverageCSPerMin float64 `json:"average_cs_per_minute"`
	Improvement     float64 `json:"improvement_percent"`
	Efficiency      float64 `json:"efficiency_percent"`
}

// CSTrendPoint represents CS trend over time
type CSTrendPoint struct {
	Date         time.Time `json:"date"`
	CSPerMinute  float64   `json:"cs_per_minute"`
	Efficiency   float64   `json:"efficiency"`
	Matches      int       `json:"matches"`
}

// PerformanceComparison represents comparative performance analysis
type PerformanceComparison struct {
	PlayerMetrics    PlayerMetricsSummary `json:"player_metrics"`
	RankBenchmarks   RankBenchmarks       `json:"rank_benchmarks"`
	RoleBenchmarks   RoleBenchmarks       `json:"role_benchmarks"`
	GlobalBenchmarks GlobalBenchmarks     `json:"global_benchmarks"`
	
	// Performance Scoring
	OverallScore     float64              `json:"overall_score"` // 0-100
	CategoryScores   map[string]float64   `json:"category_scores"`
	
	// Improvement Areas
	StrengthAreas    []string             `json:"strength_areas"`
	ImprovementAreas []string             `json:"improvement_areas"`
	
	// Ranking
	RankPercentile   float64              `json:"rank_percentile"`
	GlobalPercentile float64              `json:"global_percentile"`
}

// Performance metric structures
type PlayerMetricsSummary struct {
	KDA              float64 `json:"kda"`
	CSPerMinute      float64 `json:"cs_per_minute"`
	VisionScore      float64 `json:"vision_score"`
	DamageShare      float64 `json:"damage_share"`
	GoldEfficiency   float64 `json:"gold_efficiency"`
	WinRate          float64 `json:"win_rate"`
}

type RankBenchmarks struct {
	Rank             string              `json:"rank"`
	AverageKDA       float64             `json:"average_kda"`
	AverageCSPerMin  float64             `json:"average_cs_per_minute"`
	AverageVision    float64             `json:"average_vision_score"`
	AverageDamage    float64             `json:"average_damage_share"`
	AverageWinRate   float64             `json:"average_win_rate"`
}

type RoleBenchmarks struct {
	Role             string              `json:"role"`
	AverageKDA       float64             `json:"average_kda"`
	AverageCSPerMin  float64             `json:"average_cs_per_minute"`
	AverageVision    float64             `json:"average_vision_score"`
	AverageDamage    float64             `json:"average_damage_share"`
	ExpectedGold     float64             `json:"expected_gold_per_minute"`
}

type GlobalBenchmarks struct {
	TotalPlayers     int64               `json:"total_players"`
	AverageKDA       float64             `json:"average_kda"`
	AverageCSPerMin  float64             `json:"average_cs_per_minute"`
	AverageVision    float64             `json:"average_vision_score"`
	AverageDamage    float64             `json:"average_damage_share"`
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(db *sql.DB, matchRepo *repository.MatchRepository, playerRepo *repository.PlayerRepository, redisService *RedisService) *AnalyticsService {
	return &AnalyticsService{
		db:           db,
		matchRepo:    matchRepo,
		playerRepo:   playerRepo,
		redisService: redisService,
	}
}

// AnalyzeKDA performs comprehensive KDA analysis
func (as *AnalyticsService) AnalyzeKDA(ctx context.Context, playerID string, timeRange string, champion string) (*KDAAnalysis, error) {
	// Define time range
	startDate, endDate := as.parseTimeRange(timeRange)
	
	// Get match data
	matches, err := as.getPlayerMatches(ctx, playerID, startDate, endDate, champion)
	if err != nil {
		return nil, fmt.Errorf("failed to get player matches: %w", err)
	}
	
	if len(matches) == 0 {
		return &KDAAnalysis{
			PlayerID:  playerID,
			Champion:  champion,
			TimeRange: timeRange,
		}, nil
	}
	
	analysis := &KDAAnalysis{
		PlayerID:  playerID,
		Champion:  champion,
		TimeRange: timeRange,
	}
	
	// Calculate basic statistics
	as.calculateKDABasics(analysis, matches)
	
	// Perform trend analysis
	as.analyzeKDATrend(analysis, matches)
	
	// Calculate distribution
	as.calculateKDADistribution(analysis, matches)
	
	// Comparative analysis
	err = as.performKDAComparison(ctx, analysis, playerID)
	if err != nil {
		// Log error but don't fail the entire analysis
		fmt.Printf("Warning: failed to perform KDA comparison: %v", err)
	}
	
	// Contextual analysis
	as.analyzeKDAContext(analysis, matches)
	
	// Calculate snapshots
	as.calculateKDASnapshots(analysis, matches)
	
	// Streak analysis
	as.analyzeKDAStreak(analysis, matches)
	
	// Generate trend data for visualization
	as.generateKDATrendData(analysis, matches)
	
	// Cache results
	as.cacheKDAAnalysis(ctx, analysis)
	
	return analysis, nil
}

// AnalyzeCS performs comprehensive Creep Score analysis
func (as *AnalyticsService) AnalyzeCS(ctx context.Context, playerID string, timeRange string, position string, champion string) (*CSAnalysis, error) {
	startDate, endDate := as.parseTimeRange(timeRange)
	
	matches, err := as.getPlayerMatches(ctx, playerID, startDate, endDate, champion)
	if err != nil {
		return nil, fmt.Errorf("failed to get player matches: %w", err)
	}
	
	if len(matches) == 0 {
		return &CSAnalysis{
			PlayerID:  playerID,
			Champion:  champion,
			Position:  position,
			TimeRange: timeRange,
		}, nil
	}
	
	analysis := &CSAnalysis{
		PlayerID:  playerID,
		Champion:  champion,
		Position:  position,
		TimeRange: timeRange,
	}
	
	// Calculate CS basics
	as.calculateCSBasics(analysis, matches)
	
	// Get benchmarks
	err = as.getCSBenchmarks(ctx, analysis, position)
	if err != nil {
		fmt.Printf("Warning: failed to get CS benchmarks: %v", err)
	}
	
	// Analyze CS efficiency
	as.analyzeCSEfficiency(analysis, matches)
	
	// Trend analysis
	as.analyzeCSTimings(analysis, matches)
	
	// Contextual analysis
	as.analyzeCSContext(analysis, matches)
	
	// Generate recommendations
	as.generateCSRecommendations(analysis)
	
	// Cache results
	as.cacheCSAnalysis(ctx, analysis)
	
	return analysis, nil
}

// ComparePerformance provides comprehensive performance comparison
func (as *AnalyticsService) ComparePerformance(ctx context.Context, playerID string, timeRange string) (*PerformanceComparison, error) {
	// Get player metrics
	playerMetrics, err := as.getPlayerMetrics(ctx, playerID, timeRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get player metrics: %w", err)
	}
	
	// Get benchmarks
	rankBenchmarks, err := as.getRankBenchmarks(ctx, playerMetrics.Rank)
	if err != nil {
		return nil, fmt.Errorf("failed to get rank benchmarks: %w", err)
	}
	
	roleBenchmarks, err := as.getRoleBenchmarks(ctx, playerMetrics.MainRole)
	if err != nil {
		return nil, fmt.Errorf("failed to get role benchmarks: %w", err)
	}
	
	globalBenchmarks, err := as.getGlobalBenchmarks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get global benchmarks: %w", err)
	}
	
	comparison := &PerformanceComparison{
		PlayerMetrics:    playerMetrics.Summary,
		RankBenchmarks:   *rankBenchmarks,
		RoleBenchmarks:   *roleBenchmarks,
		GlobalBenchmarks: *globalBenchmarks,
		CategoryScores:   make(map[string]float64),
	}
	
	// Calculate performance scores
	as.calculatePerformanceScores(comparison)
	
	// Identify strengths and improvement areas
	as.identifyPerformanceAreas(comparison)
	
	// Calculate percentiles
	as.calculatePercentiles(comparison)
	
	return comparison, nil
}

// Helper functions for KDA analysis

func (as *AnalyticsService) calculateKDABasics(analysis *KDAAnalysis, matches []models.MatchData) {
	if len(matches) == 0 {
		return
	}
	
	kdaValues := make([]float64, 0, len(matches))
	
	for _, match := range matches {
		analysis.TotalKills += match.Kills
		analysis.TotalDeaths += match.Deaths
		analysis.TotalAssists += match.Assists
		
		kda := as.calculateKDA(match.Kills, match.Deaths, match.Assists)
		kdaValues = append(kdaValues, kda)
	}
	
	// Calculate statistics
	analysis.AverageKDA = as.calculateMean(kdaValues)
	analysis.MedianKDA = as.calculateMedian(kdaValues)
	analysis.BestKDA = as.calculateMax(kdaValues)
	analysis.WorstKDA = as.calculateMin(kdaValues)
	analysis.StandardDev = as.calculateStandardDeviation(kdaValues)
}

func (as *AnalyticsService) analyzeKDATrend(analysis *KDAAnalysis, matches []models.MatchData) {
	if len(matches) < 5 {
		analysis.TrendDirection = "insufficient_data"
		return
	}
	
	// Sort matches by date
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Date.Before(matches[j].Date)
	})
	
	// Calculate moving averages to smooth the trend
	windowSize := 5
	if len(matches) < windowSize {
		windowSize = len(matches)
	}
	
	movingAverages := make([]float64, 0)
	for i := windowSize-1; i < len(matches); i++ {
		sum := 0.0
		for j := i-windowSize+1; j <= i; j++ {
			kda := as.calculateKDA(matches[j].Kills, matches[j].Deaths, matches[j].Assists)
			sum += kda
		}
		movingAverages = append(movingAverages, sum/float64(windowSize))
	}
	
	// Calculate trend slope using linear regression
	slope, confidence := as.calculateLinearRegression(movingAverages)
	
	analysis.TrendSlope = slope
	analysis.TrendConfidence = confidence
	
	// Determine trend direction
	if slope > 0.05 && confidence > 0.6 {
		analysis.TrendDirection = "improving"
	} else if slope < -0.05 && confidence > 0.6 {
		analysis.TrendDirection = "declining"
	} else {
		analysis.TrendDirection = "stable"
	}
}

func (as *AnalyticsService) calculateKDADistribution(analysis *KDAAnalysis, matches []models.MatchData) {
	distribution := map[string]int{
		"excellent": 0, // KDA >= 3.0
		"good":      0, // KDA >= 2.0
		"average":   0, // KDA >= 1.0
		"poor":      0, // KDA < 1.0
	}
	
	for _, match := range matches {
		kda := as.calculateKDA(match.Kills, match.Deaths, match.Assists)
		
		if kda >= 3.0 {
			distribution["excellent"]++
		} else if kda >= 2.0 {
			distribution["good"]++
		} else if kda >= 1.0 {
			distribution["average"]++
		} else {
			distribution["poor"]++
		}
	}
	
	analysis.KDADistribution = distribution
}

func (as *AnalyticsService) performKDAComparison(ctx context.Context, analysis *KDAAnalysis, playerID string) error {
	// Get player's current rank
	player, err := as.playerRepo.GetByID(ctx, playerID)
	if err != nil {
		return err
	}
	
	// Calculate rank percentile
	rankPercentile, err := as.calculateRankPercentile(ctx, player.CurrentRank, analysis.AverageKDA, "kda")
	if err == nil {
		analysis.RankPercentile = rankPercentile
	}
	
	// Calculate global percentile
	globalPercentile, err := as.calculateGlobalPercentile(ctx, analysis.AverageKDA, "kda")
	if err == nil {
		analysis.GlobalPercentile = globalPercentile
	}
	
	return nil
}

// Helper functions for CS analysis

func (as *AnalyticsService) calculateCSBasics(analysis *CSAnalysis, matches []models.MatchData) {
	if len(matches) == 0 {
		return
	}
	
	csPerMinValues := make([]float64, 0, len(matches))
	
	for _, match := range matches {
		analysis.TotalCS += match.TotalCS
		
		csPerMin := float64(match.TotalCS) / (float64(match.GameDuration) / 60.0)
		csPerMinValues = append(csPerMinValues, csPerMin)
	}
	
	analysis.AverageCS = float64(analysis.TotalCS) / float64(len(matches))
	analysis.AverageCSPerMin = as.calculateMean(csPerMinValues)
	analysis.BestCSPerMin = as.calculateMax(csPerMinValues)
	analysis.MedianCSPerMin = as.calculateMedian(csPerMinValues)
}

func (as *AnalyticsService) getCSBenchmarks(ctx context.Context, analysis *CSAnalysis, position string) error {
	// Get role-specific benchmarks from database
	query := `
		SELECT 
			AVG(total_cs::float / (game_duration::float / 60)) as avg_cs_per_min
		FROM match_participants mp 
		JOIN matches m ON mp.match_id = m.id 
		WHERE mp.position = $1 
		AND m.created_at > NOW() - INTERVAL '30 days'
	`
	
	var roleAvg sql.NullFloat64
	err := as.db.QueryRowContext(ctx, query, position).Scan(&roleAvg)
	if err == nil && roleAvg.Valid {
		analysis.RoleAverage = roleAvg.Float64
	}
	
	// Set benchmark values (these would typically come from database)
	switch position {
	case "ADC", "MID":
		analysis.RankAverage = 7.5
		analysis.ProAverage = 9.2
	case "TOP":
		analysis.RankAverage = 7.0
		analysis.ProAverage = 8.8
	case "JUNGLE":
		analysis.RankAverage = 6.5
		analysis.ProAverage = 7.8
	case "SUPPORT":
		analysis.RankAverage = 2.5
		analysis.ProAverage = 3.2
	default:
		analysis.RankAverage = 6.0
		analysis.ProAverage = 7.5
	}
	
	return nil
}

func (as *AnalyticsService) analyzeCSEfficiency(analysis *CSAnalysis, matches []models.MatchData) {
	if len(matches) == 0 {
		return
	}
	
	efficiencySum := 0.0
	
	for _, match := range matches {
		// Calculate theoretical maximum CS based on game duration
		theoreticalMaxCS := as.calculateTheoreticalMaxCS(match.GameDuration, match.Position)
		
		efficiency := float64(match.TotalCS) / theoreticalMaxCS * 100
		if efficiency > 100 {
			efficiency = 100 // Cap at 100%
		}
		
		efficiencySum += efficiency
	}
	
	analysis.CSEfficiency = efficiencySum / float64(len(matches))
}

func (as *AnalyticsService) generateCSRecommendations(analysis *CSAnalysis) {
	recommendations := make([]string, 0)
	
	if analysis.CSEfficiency < 60 {
		recommendations = append(recommendations, "Focus on last-hitting minions - practice in training mode")
		recommendations = append(recommendations, "Avoid missing CS while harassing enemies")
	}
	
	if analysis.AverageCSPerMin < analysis.RoleAverage*0.9 {
		recommendations = append(recommendations, "Work on farming patterns specific to your role")
	}
	
	if analysis.EarlyGameCS < analysis.MidGameCS {
		recommendations = append(recommendations, "Improve early game laning phase farming")
	}
	
	if analysis.LateGameCS < analysis.MidGameCS {
		recommendations = append(recommendations, "Focus on side lane farming in late game")
	}
	
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Excellent CS performance! Keep up the consistency")
	}
	
	analysis.Recommendations = recommendations
}

// Mathematical utility functions

func (as *AnalyticsService) calculateKDA(kills, deaths, assists int) float64 {
	if deaths == 0 {
		return float64(kills + assists)
	}
	return float64(kills+assists) / float64(deaths)
}

func (as *AnalyticsService) calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (as *AnalyticsService) calculateMedian(values []float64) float64 {
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

func (as *AnalyticsService) calculateMax(values []float64) float64 {
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

func (as *AnalyticsService) calculateMin(values []float64) float64 {
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

func (as *AnalyticsService) calculateStandardDeviation(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	mean := as.calculateMean(values)
	sumSquaredDiffs := 0.0
	
	for _, v := range values {
		diff := v - mean
		sumSquaredDiffs += diff * diff
	}
	
	variance := sumSquaredDiffs / float64(len(values))
	return math.Sqrt(variance)
}

func (as *AnalyticsService) calculateLinearRegression(values []float64) (slope, confidence float64) {
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
	
	// Calculate slope using least squares
	denominator := n*sumXX - sumX*sumX
	if denominator == 0 {
		return 0, 0
	}
	
	slope = (n*sumXY - sumX*sumY) / denominator
	
	// Calculate R-squared for confidence
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

// Cache operations
func (as *AnalyticsService) cacheKDAAnalysis(ctx context.Context, analysis *KDAAnalysis) {
	if as.redisService == nil {
		return
	}
	
	key := fmt.Sprintf("kda_analysis:%s:%s:%s", analysis.PlayerID, analysis.TimeRange, analysis.Champion)
	as.redisService.SetJSON(ctx, key, analysis, 30*time.Minute)
}

func (as *AnalyticsService) cacheCSAnalysis(ctx context.Context, analysis *CSAnalysis) {
	if as.redisService == nil {
		return
	}
	
	key := fmt.Sprintf("cs_analysis:%s:%s:%s", analysis.PlayerID, analysis.TimeRange, analysis.Champion)
	as.redisService.SetJSON(ctx, key, analysis, 30*time.Minute)
}

// Additional helper functions would be implemented here...

func (as *AnalyticsService) parseTimeRange(timeRange string) (time.Time, time.Time) {
	endDate := time.Now()
	var startDate time.Time
	
	switch timeRange {
	case "7d":
		startDate = endDate.AddDate(0, 0, -7)
	case "30d":
		startDate = endDate.AddDate(0, 0, -30)
	case "90d":
		startDate = endDate.AddDate(0, 0, -90)
	default:
		startDate = endDate.AddDate(0, 0, -30) // Default to 30 days
	}
	
	return startDate, endDate
}

func (as *AnalyticsService) getPlayerMatches(ctx context.Context, playerID string, startDate, endDate time.Time, champion string) ([]models.MatchData, error) {
	// This would typically query the database
	// For now, return empty slice
	return []models.MatchData{}, nil
}

func (as *AnalyticsService) calculateTheoreticalMaxCS(gameDuration int, position string) float64 {
	// Simplified calculation - in reality this would be more complex
	minutes := float64(gameDuration) / 60.0
	
	switch position {
	case "ADC", "MID":
		return minutes * 12.0 // ~12 CS per minute theoretical max
	case "TOP":
		return minutes * 11.0
	case "JUNGLE":
		return minutes * 8.0 // Different calculation for jungle
	case "SUPPORT":
		return minutes * 4.0 // Support typically has lower CS
	default:
		return minutes * 10.0
	}
}