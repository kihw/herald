package services

import (
	"fmt"
	"log"
	"time"

	"lol-match-exporter/internal/db"
)

// AnalyticsService provides comprehensive analytics using Go native services
type AnalyticsService struct {
	analyticsEngine     *AnalyticsEngineService
	mmrCalculation      *MMRCalculationService
	recommendationEngine *RecommendationEngineService
}

// NewAnalyticsService creates a new analytics service with Go native engines
func NewAnalyticsService(database *db.Database) *AnalyticsService {
	return &AnalyticsService{
		analyticsEngine:     NewAnalyticsEngineService(database),
		mmrCalculation:      NewMMRCalculationService(database),
		recommendationEngine: NewRecommendationEngineService(database),
	}
}

// PeriodStats represents analytics for a specific time period
type PeriodStats struct {
	Period           string                    `json:"period"`
	TotalGames       int                      `json:"total_games"`
	WinRate          float64                  `json:"win_rate"`
	AvgKDA           float64                  `json:"avg_kda"`
	BestRole         string                   `json:"best_role"`
	WorstRole        string                   `json:"worst_role"`
	TopChampions     []ChampionPerformance    `json:"top_champions"`
	RolePerformance  map[string]RoleMetrics   `json:"role_performance"`
	RecentTrend      string                   `json:"recent_trend"`
	Suggestions      []string                 `json:"suggestions"`
}

// ChampionPerformance represents champion performance data
type ChampionPerformance struct {
	ChampionID       int     `json:"champion_id"`
	ChampionName     string  `json:"champion_name"`
	Games            int     `json:"games"`
	WinRate          float64 `json:"win_rate"`
	PerformanceScore float64 `json:"performance_score"`
	AvgKDA           float64 `json:"avg_kda"`
}

// RoleMetrics represents performance metrics for a role
type RoleMetrics struct {
	GamesPlayed      int     `json:"games_played"`
	Wins             int     `json:"wins"`
	Losses           int     `json:"losses"`
	WinRate          float64 `json:"win_rate"`
	AvgKills         float64 `json:"avg_kills"`
	AvgDeaths        float64 `json:"avg_deaths"`
	AvgAssists       float64 `json:"avg_assists"`
	AvgKDA           float64 `json:"avg_kda"`
	AvgCSPerMin      float64 `json:"avg_cs_per_min"`
	AvgGoldPerMin    float64 `json:"avg_gold_per_min"`
	AvgDamagePerMin  float64 `json:"avg_damage_per_min"`
	AvgVisionScore   float64 `json:"avg_vision_score"`
	PerformanceScore float64 `json:"performance_score"`
	TrendDirection   string  `json:"trend_direction"`
}

// MMRTrajectory represents MMR analysis over time
type MMRTrajectory struct {
	MMRHistory      []MMRDataPoint `json:"mmr_history"`
	CurrentMMR      int            `json:"current_mmr"`
	CurrentRank     string         `json:"current_rank"`
	MMRRange        MMRRange       `json:"mmr_range"`
	Volatility      float64        `json:"volatility"`
	Trend           string         `json:"trend"`
	ConfidenceGrade string         `json:"confidence_grade"`
}

// MMRDataPoint represents a single MMR data point
type MMRDataPoint struct {
	Date          time.Time `json:"date"`
	MatchID       string    `json:"match_id"`
	EstimatedMMR  int       `json:"estimated_mmr"`
	MMRChange     int       `json:"mmr_change"`
	Confidence    float64   `json:"confidence"`
	RankEstimate  string    `json:"rank_estimate"`
}

// MMRRange represents min/max MMR values
type MMRRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// Recommendation represents an AI-powered suggestion
type Recommendation struct {
	Type                string    `json:"type"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	Priority            int       `json:"priority"`
	Confidence          float64   `json:"confidence"`
	ExpectedImprovement string    `json:"expected_improvement"`
	ActionItems         []string  `json:"action_items"`
	ChampionID          *int      `json:"champion_id,omitempty"`
	Role                *string   `json:"role,omitempty"`
	TimePeriod          string    `json:"time_period"`
	ExpiresAt           *time.Time `json:"expires_at,omitempty"`
}

// GetPeriodStats retrieves analytics for a specific time period using native Go service
func (s *AnalyticsService) GetPeriodStats(userID int, period string) (*PeriodStats, error) {
	// Use native Go analytics engine
	nativeStats, err := s.analyticsEngine.GeneratePeriodStats(userID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to generate period stats: %w", err)
	}

	// Convert native models to API models
	apiStats := &PeriodStats{
		Period:      nativeStats.Period,
		TotalGames:  nativeStats.TotalGames,
		WinRate:     nativeStats.WinRate,
		AvgKDA:      nativeStats.AvgKDA,
		BestRole:    nativeStats.BestRole,
		WorstRole:   nativeStats.WorstRole,
		RecentTrend: nativeStats.RecentTrend,
		Suggestions: nativeStats.Suggestions,
	}

	// Convert top champions
	apiStats.TopChampions = make([]ChampionPerformance, len(nativeStats.TopChampions))
	for i, champ := range nativeStats.TopChampions {
		apiStats.TopChampions[i] = ChampionPerformance{
			ChampionID:       champ.ChampionID,
			ChampionName:     champ.ChampionName,
			Games:            champ.Games,
			WinRate:          champ.WinRate,
			PerformanceScore: champ.PerformanceScore,
			AvgKDA:           champ.AvgKDA,
		}
	}

	// Convert role performance
	apiStats.RolePerformance = make(map[string]RoleMetrics)
	for role, metrics := range nativeStats.RolePerformance {
		apiStats.RolePerformance[role] = RoleMetrics{
			GamesPlayed:      metrics.GamesPlayed,
			Wins:             metrics.Wins,
			Losses:           metrics.Losses,
			WinRate:          metrics.WinRate,
			AvgKills:         metrics.AvgKills,
			AvgDeaths:        metrics.AvgDeaths,
			AvgAssists:       metrics.AvgAssists,
			AvgKDA:           metrics.AvgKDA,
			AvgCSPerMin:      metrics.AvgCSPerMin,
			AvgGoldPerMin:    metrics.AvgGoldPerMin,
			AvgDamagePerMin:  metrics.AvgDamagePerMin,
			AvgVisionScore:   metrics.AvgVisionScore,
			PerformanceScore: metrics.PerformanceScore,
			TrendDirection:   string(metrics.TrendDirection),
		}
	}

	return apiStats, nil
}

// GetMMRTrajectory calculates MMR trajectory over time using native Go service
func (s *AnalyticsService) GetMMRTrajectory(userID int, days int) (*MMRTrajectory, error) {
	// Use native Go MMR calculation service
	nativeTrajectory, err := s.mmrCalculation.CalculateMMRTrajectory(userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate MMR trajectory: %w", err)
	}

	// Convert native models to API models
	apiTrajectory := &MMRTrajectory{
		CurrentMMR:      nativeTrajectory.CurrentMMR,
		CurrentRank:     nativeTrajectory.CurrentRank,
		Volatility:      nativeTrajectory.Volatility,
		Trend:           nativeTrajectory.Trend,
		ConfidenceGrade: fmt.Sprintf("%.1f", nativeTrajectory.ConfidenceGrade*100),
		MMRRange: MMRRange{
			Min: nativeTrajectory.MMRRange.Min,
			Max: nativeTrajectory.MMRRange.Max,
		},
	}

	// Convert MMR history
	apiTrajectory.MMRHistory = make([]MMRDataPoint, len(nativeTrajectory.MMRHistory))
	for i, entry := range nativeTrajectory.MMRHistory {
		apiTrajectory.MMRHistory[i] = MMRDataPoint{
			Date:         entry.Date,
			MatchID:      entry.MatchID,
			EstimatedMMR: entry.EstimatedMMR,
			MMRChange:    entry.MMRChange,
			Confidence:   entry.Confidence,
			RankEstimate: entry.RankEstimate,
		}
	}

	return apiTrajectory, nil
}

// GetRecommendations retrieves AI-powered recommendations using native Go service
func (s *AnalyticsService) GetRecommendations(userID int) ([]Recommendation, error) {
	// Use native Go recommendation engine
	nativeRecs, err := s.recommendationEngine.GenerateComprehensiveRecommendations(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendations: %w", err)
	}

	// Convert native models to API models
	apiRecs := make([]Recommendation, len(nativeRecs))
	for i, rec := range nativeRecs {
		apiRecs[i] = Recommendation{
			Type:                string(rec.Type),
			Title:               rec.Title,
			Description:         rec.Description,
			Priority:            rec.Priority,
			Confidence:          rec.Confidence,
			ExpectedImprovement: rec.ExpectedImprovement,
			ActionItems:         rec.ActionItems,
			ChampionID:          rec.ChampionID,
			Role:                rec.Role,
			TimePeriod:          rec.TimePeriod,
			ExpiresAt:           rec.ExpiresAt,
		}
	}

	return apiRecs, nil
}

// GetChampionAnalysis analyzes performance for a specific champion using native Go service
func (s *AnalyticsService) GetChampionAnalysis(userID int, championID int, period string) (map[string]interface{}, error) {
	// Use native Go analytics engine
	analysis, err := s.analyticsEngine.AnalyzeChampionMastery(userID, championID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze champion mastery: %w", err)
	}

	// Convert to map for API compatibility
	result := map[string]interface{}{
		"champion_id":            analysis.ChampionID,
		"champion_name":          analysis.ChampionName,
		"games_played":           analysis.GamesPlayed,
		"mastery_score":          analysis.MasteryScore,
		"performance_metrics":    analysis.PerformanceMetrics,
		"best_game":              analysis.BestGame,
		"worst_game":             analysis.WorstGame,
		"improvement_suggestions": analysis.ImprovementSuggestions,
		"skill_progression":      analysis.SkillProgression,
	}

	return result, nil
}

// UpdateChampionStats triggers champion stats calculation and database update using native Go service
func (s *AnalyticsService) UpdateChampionStats(userID int, period string) error {
	// Use native Go analytics engine to recalculate stats
	_, err := s.analyticsEngine.GeneratePeriodStats(userID, period)
	if err != nil {
		return fmt.Errorf("failed to update champion stats: %w", err)
	}

	log.Printf("Successfully updated champion stats for user %d (period: %s)", userID, period)
	return nil
}

// GetPerformanceTrends calculates detailed performance trends using native Go service
func (s *AnalyticsService) GetPerformanceTrends(userID int) (map[string]interface{}, error) {
	// This would use a new method in the analytics engine for trends
	// For now, return basic analytics as trends
	stats, err := s.analyticsEngine.GeneratePeriodStats(userID, "month")
	if err != nil {
		return nil, fmt.Errorf("failed to calculate performance trends: %w", err)
	}

	trends := map[string]interface{}{
		"daily_trend":         map[string]interface{}{"trend": stats.RecentTrend},
		"weekly_trend":        map[string]interface{}{"trend": stats.RecentTrend},
		"monthly_trend":       map[string]interface{}{"trend": stats.RecentTrend},
		"improvement_velocity": 0.5, // Placeholder
		"consistency_score":   85.0, // Placeholder
	}

	return trends, nil
}

// Worker Interface Adapter Methods
// These methods adapt the typed methods to interface{} for worker compatibility

// GetPeriodStatsAsInterface returns period stats as interface{} for worker compatibility
func (s *AnalyticsService) GetPeriodStatsAsInterface(userID int, period string) (interface{}, error) {
	return s.GetPeriodStats(userID, period)
}

// GetMMRTrajectoryAsInterface returns MMR trajectory as interface{} for worker compatibility
func (s *AnalyticsService) GetMMRTrajectoryAsInterface(userID int, days int) (interface{}, error) {
	return s.GetMMRTrajectory(userID, days)
}

// GetRecommendationsAsInterface returns recommendations as interface{} for worker compatibility
func (s *AnalyticsService) GetRecommendationsAsInterface(userID int) (interface{}, error) {
	return s.GetRecommendations(userID)
}

// GetChampionAnalysisAsInterface returns champion analysis as interface{} for worker compatibility
func (s *AnalyticsService) GetChampionAnalysisAsInterface(userID int, championID int, period string) (interface{}, error) {
	return s.GetChampionAnalysis(userID, championID, period)
}

// ProcessNewMatches triggers analytics processing for new match data using native Go service
func (s *AnalyticsService) ProcessNewMatches(userID int, matchIDs []string) error {
	if len(matchIDs) == 0 {
		return nil // No matches to process
	}

	// Process through native analytics engine
	// This would typically analyze the new matches and update all relevant statistics
	err := s.UpdateChampionStats(userID, "all")
	if err != nil {
		return fmt.Errorf("failed to process new matches: %w", err)
	}

	log.Printf("Successfully processed %d new matches for user %d using Go native service", len(matchIDs), userID)
	return nil
}

// ValidateEnvironment checks if Go native services are properly initialized
func (s *AnalyticsService) ValidateEnvironment() error {
	// Check if all native services are available
	if s.analyticsEngine == nil {
		return fmt.Errorf("analytics engine service not initialized")
	}
	
	if s.mmrCalculation == nil {
		return fmt.Errorf("MMR calculation service not initialized")
	}
	
	if s.recommendationEngine == nil {
		return fmt.Errorf("recommendation engine service not initialized")
	}

	log.Println("Go native analytics services validated successfully")
	return nil
}

// Additional methods for advanced analytics

// GetMMRPrediction predicts rank changes using native Go service
func (s *AnalyticsService) GetMMRPrediction(userID int, targetRank string) (map[string]interface{}, error) {
	prediction, err := s.mmrCalculation.PredictRankChanges(userID, targetRank)
	if err != nil {
		return nil, fmt.Errorf("failed to predict rank changes: %w", err)
	}

	result := map[string]interface{}{
		"current_rank":     prediction.CurrentRank,
		"predicted_rank":   prediction.PredictedRank,
		"lp_needed":        prediction.LPNeeded,
		"games_needed":     prediction.GamesNeeded,
		"win_rate_required": prediction.WinRateRequired,
		"confidence":       prediction.Confidence,
		"timeline_days":    prediction.TimelineDays,
	}

	return result, nil
}

// GetVolatilityAnalysis analyzes MMR volatility using native Go service
func (s *AnalyticsService) GetVolatilityAnalysis(userID int) (map[string]interface{}, error) {
	analysis, err := s.mmrCalculation.AnalyzeMMRVolatility(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze MMR volatility: %w", err)
	}

	result := map[string]interface{}{
		"volatility":        analysis.Volatility,
		"consistency_score": analysis.ConsistencyScore,
		"stability_rating":  analysis.StabilityRating,
		"streak_analysis":   analysis.StreakAnalysis,
		"risk_assessment":   analysis.RiskAssessment,
		"recommendations":   analysis.Recommendations,
	}

	return result, nil
}

// GetSkillCeiling calculates estimated skill ceiling using native Go service
func (s *AnalyticsService) GetSkillCeiling(userID int, role string) (map[string]interface{}, error) {
	ceiling, err := s.mmrCalculation.CalculateSkillCeiling(userID, role)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate skill ceiling: %w", err)
	}

	result := map[string]interface{}{
		"current_skill_level": ceiling.CurrentSkillLevel,
		"estimated_ceiling":   ceiling.EstimatedCeiling,
		"peak_performances":   ceiling.PeakPerformances,
		"improvement_rate":    ceiling.ImprovementRate,
		"time_to_ceiling":     ceiling.TimeToCeiling,
		"confidence":          ceiling.Confidence,
	}

	return result, nil
}