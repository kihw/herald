package analytics

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"herald.lol/internal/riot"
)

// Herald.lol Gaming Analytics - Core Analytics Engine
// Advanced gaming metrics calculation engine for League of Legends

// AnalyticsEngine handles core gaming metrics calculations
type AnalyticsEngine struct {
	config *AnalyticsConfig
}

// AnalyticsConfig contains engine configuration
type AnalyticsConfig struct {
	// Performance thresholds for different ranks
	RankThresholds   map[string]*RankThresholds `json:"rank_thresholds"`
	
	// Metric weights for overall performance scoring
	MetricWeights    *MetricWeights             `json:"metric_weights"`
	
	// Analysis settings
	MinMatchesRequired    int     `json:"min_matches_required"`
	RecentMatchesWindow   int     `json:"recent_matches_window"`
	TrendConfidenceMin    float64 `json:"trend_confidence_min"`
	
	// Champion role mappings and expectations
	ChampionRoles        map[string]string       `json:"champion_roles"`
	RoleExpectations     map[string]*RoleMetrics `json:"role_expectations"`
	
	// Advanced analytics settings
	EnableAIInsights     bool    `json:"enable_ai_insights"`
	EnablePredictions    bool    `json:"enable_predictions"`
	PerformanceDecay     float64 `json:"performance_decay"`    // How much old matches decay
}

// RankThresholds defines performance expectations by rank
type RankThresholds struct {
	MinKDA           float64 `json:"min_kda"`
	MinCSPerMin      float64 `json:"min_cs_per_min"`
	MinVisionScore   float64 `json:"min_vision_score"`
	MinDamageShare   float64 `json:"min_damage_share"`
	MinGoldEff       float64 `json:"min_gold_efficiency"`
	MinWinRate       float64 `json:"min_win_rate"`
}

// MetricWeights defines importance of each metric
type MetricWeights struct {
	KDA              float64 `json:"kda"`
	CSPerMinute      float64 `json:"cs_per_minute"`
	VisionScore      float64 `json:"vision_score"`
	DamageShare      float64 `json:"damage_share"`
	GoldEfficiency   float64 `json:"gold_efficiency"`
	WinRate          float64 `json:"win_rate"`
	ObjectiveControl float64 `json:"objective_control"`
	Positioning      float64 `json:"positioning"`
}

// RoleMetrics defines role-specific expectations
type RoleMetrics struct {
	ExpectedKDA        float64 `json:"expected_kda"`
	ExpectedCS         float64 `json:"expected_cs"`
	ExpectedDamage     float64 `json:"expected_damage"`
	ExpectedVision     float64 `json:"expected_vision"`
	ExpectedGold       int     `json:"expected_gold"`
	PriorityStats      []string `json:"priority_stats"`
}

// NewAnalyticsEngine creates new analytics engine
func NewAnalyticsEngine(config *AnalyticsConfig) *AnalyticsEngine {
	if config == nil {
		config = DefaultAnalyticsConfig()
	}
	
	return &AnalyticsEngine{
		config: config,
	}
}

// AnalyzePlayer performs comprehensive player analysis
func (a *AnalyticsEngine) AnalyzePlayer(ctx context.Context, request *PlayerAnalysisRequest) (*PlayerAnalysis, error) {
	if len(request.Matches) < a.config.MinMatchesRequired {
		return nil, fmt.Errorf("insufficient matches: need at least %d, got %d", 
			a.config.MinMatchesRequired, len(request.Matches))
	}

	analysis := &PlayerAnalysis{
		SummonerID:   request.SummonerID,
		SummonerName: request.SummonerName,
		Region:       request.Region,
		AnalyzedAt:   time.Now(),
		TotalMatches: len(request.Matches),
	}

	// Calculate core metrics
	coreMetrics, err := a.calculateCoreMetrics(request.Matches, request.PlayerPUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate core metrics: %w", err)
	}
	analysis.CoreMetrics = coreMetrics

	// Calculate role-specific metrics
	roleMetrics, err := a.calculateRoleMetrics(request.Matches, request.PlayerPUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate role metrics: %w", err)
	}
	analysis.RoleMetrics = roleMetrics

	// Calculate champion-specific metrics
	championMetrics, err := a.calculateChampionMetrics(request.Matches, request.PlayerPUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate champion metrics: %w", err)
	}
	analysis.ChampionMetrics = championMetrics

	// Perform trend analysis
	trends, err := a.calculateTrends(request.Matches, request.PlayerPUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate trends: %w", err)
	}
	analysis.Trends = trends

	// Calculate overall performance score
	analysis.PerformanceScore = a.calculatePerformanceScore(coreMetrics, request.CurrentRank)

	// Generate insights
	if a.config.EnableAIInsights {
		insights, err := a.generateInsights(analysis, request.CurrentRank)
		if err != nil {
			return nil, fmt.Errorf("failed to generate insights: %w", err)
		}
		analysis.Insights = insights
	}

	return analysis, nil
}

// calculateCoreMetrics calculates KDA, CS/min, Vision, Damage Share, Gold Efficiency
func (a *AnalyticsEngine) calculateCoreMetrics(matches []*riot.Match, playerPUUID string) (*CoreMetrics, error) {
	if len(matches) == 0 {
		return nil, fmt.Errorf("no matches provided")
	}

	var (
		totalKills    int
		totalDeaths   int
		totalAssists  int
		totalCS       int
		totalDuration int
		totalGold     int
		totalDamage   int64
		totalVision   int
		totalMatches  = len(matches)
		wins         int
	)

	var teamDamages []int64 // For damage share calculation
	var goldPerMinutes []float64 // For gold efficiency

	for _, match := range matches {
		participant := a.findPlayerInMatch(match, playerPUUID)
		if participant == nil {
			continue
		}

		// Basic stats
		totalKills += participant.Kills
		totalDeaths += participant.Deaths
		totalAssists += participant.Assists
		totalCS += participant.TotalMinionsKilled + participant.NeutralMinionsKilled
		totalDuration += match.Info.GameDuration
		totalGold += participant.GoldEarned
		totalDamage += int64(participant.TotalDamageDealtToChampions)
		totalVision += participant.VisionScore

		if participant.Win {
			wins++
		}

		// Calculate team damage for damage share
		teamDamage := a.calculateTeamDamage(match, participant.TeamID)
		if teamDamage > 0 {
			teamDamages = append(teamDamages, teamDamage)
		}

		// Calculate gold per minute
		if match.Info.GameDuration > 0 {
			goldPerMin := float64(participant.GoldEarned) / (float64(match.Info.GameDuration) / 60.0)
			goldPerMinutes = append(goldPerMinutes, goldPerMin)
		}
	}

	// Avoid division by zero
	if totalDeaths == 0 {
		totalDeaths = 1
	}

	metrics := &CoreMetrics{
		// KDA Calculation
		AverageKDA: float64(totalKills+totalAssists) / float64(totalDeaths),
		AverageKills: float64(totalKills) / float64(totalMatches),
		AverageDeaths: float64(totalDeaths) / float64(totalMatches),
		AverageAssists: float64(totalAssists) / float64(totalMatches),

		// CS Metrics
		AverageCS: float64(totalCS) / float64(totalMatches),
		CSPerMinute: float64(totalCS) / (float64(totalDuration) / 60.0),

		// Vision Metrics
		AverageVision: float64(totalVision) / float64(totalMatches),

		// Economic Metrics
		AverageGold: totalGold / totalMatches,
		AverageDamage: int(totalDamage / int64(totalMatches)),

		// Win Rate
		WinRate: float64(wins) / float64(totalMatches),
	}

	// Calculate Damage Share
	if len(teamDamages) > 0 {
		var totalDamageShare float64
		for i, match := range matches {
			if i < len(teamDamages) {
				participant := a.findPlayerInMatch(match, playerPUUID)
				if participant != nil && teamDamages[i] > 0 {
					damageShare := float64(participant.TotalDamageDealtToChampions) / float64(teamDamages[i])
					totalDamageShare += damageShare
				}
			}
		}
		metrics.DamageShare = totalDamageShare / float64(len(teamDamages))
	}

	// Calculate Gold Efficiency
	if len(goldPerMinutes) > 0 {
		var totalGoldEff float64
		for _, gpm := range goldPerMinutes {
			// Gold efficiency based on expected gold per minute for the game duration
			// Higher GPM = better gold efficiency
			goldEff := math.Min(gpm / 400.0, 2.0) // Cap at 2.0 for very high efficiency
			totalGoldEff += goldEff
		}
		metrics.GoldEfficiency = totalGoldEff / float64(len(goldPerMinutes))
	}

	return metrics, nil
}

// calculateRoleMetrics calculates role-specific performance metrics
func (a *AnalyticsEngine) calculateRoleMetrics(matches []*riot.Match, playerPUUID string) (map[string]*RolePerformance, error) {
	roleStats := make(map[string]*roleAccumulator)

	for _, match := range matches {
		participant := a.findPlayerInMatch(match, playerPUUID)
		if participant == nil {
			continue
		}

		role := a.normalizeRole(participant.TeamPosition)
		if role == "" {
			continue
		}

		if roleStats[role] == nil {
			roleStats[role] = &roleAccumulator{}
		}

		acc := roleStats[role]
		acc.matches++
		acc.totalKills += participant.Kills
		acc.totalDeaths += participant.Deaths
		acc.totalAssists += participant.Assists
		acc.totalCS += participant.TotalMinionsKilled + participant.NeutralMinionsKilled
		acc.totalDamage += int64(participant.TotalDamageDealtToChampions)
		acc.totalGold += participant.GoldEarned
		acc.totalVision += participant.VisionScore
		acc.totalDuration += match.Info.GameDuration

		if participant.Win {
			acc.wins++
		}
	}

	roleMetrics := make(map[string]*RolePerformance)
	for role, acc := range roleStats {
		if acc.matches < 3 { // Need minimum matches for reliable stats
			continue
		}

		performance := &RolePerformance{
			Role:         role,
			GamesPlayed:  acc.matches,
			WinRate:      float64(acc.wins) / float64(acc.matches),
			AverageKDA:   float64(acc.totalKills+acc.totalAssists) / math.Max(float64(acc.totalDeaths), 1),
			AverageCS:    float64(acc.totalCS) / float64(acc.matches),
			CSPerMinute:  float64(acc.totalCS) / (float64(acc.totalDuration) / 60.0),
			AverageDamage: int(acc.totalDamage / int64(acc.matches)),
			AverageGold:  acc.totalGold / acc.matches,
			AverageVision: float64(acc.totalVision) / float64(acc.matches),
		}

		// Calculate performance rating compared to role expectations
		performance.PerformanceRating = a.calculateRoleRating(performance, role)
		
		roleMetrics[role] = performance
	}

	return roleMetrics, nil
}

// calculateChampionMetrics calculates champion-specific performance
func (a *AnalyticsEngine) calculateChampionMetrics(matches []*riot.Match, playerPUUID string) ([]*ChampionPerformance, error) {
	championStats := make(map[string]*championAccumulator)

	for _, match := range matches {
		participant := a.findPlayerInMatch(match, playerPUUID)
		if participant == nil {
			continue
		}

		champion := participant.ChampionName
		if championStats[champion] == nil {
			championStats[champion] = &championAccumulator{
				championName: champion,
			}
		}

		acc := championStats[champion]
		acc.matches++
		acc.totalKills += participant.Kills
		acc.totalDeaths += participant.Deaths
		acc.totalAssists += participant.Assists
		acc.totalCS += participant.TotalMinionsKilled + participant.NeutralMinionsKilled
		acc.totalDuration += match.Info.GameDuration

		if participant.Win {
			acc.wins++
		}
	}

	var championMetrics []*ChampionPerformance
	for champion, acc := range championStats {
		if acc.matches < 2 { // Need minimum matches
			continue
		}

		performance := &ChampionPerformance{
			ChampionName: champion,
			GamesPlayed:  acc.matches,
			WinRate:      float64(acc.wins) / float64(acc.matches),
			AverageKDA:   float64(acc.totalKills+acc.totalAssists) / math.Max(float64(acc.totalDeaths), 1),
			AverageCS:    float64(acc.totalCS) / float64(acc.matches),
			CSPerMinute:  float64(acc.totalCS) / (float64(acc.totalDuration) / 60.0),
		}

		// Calculate mastery level and performance trend
		performance.MasteryLevel = a.calculateMasteryLevel(performance)
		performance.PerformanceTrend = a.calculateChampionTrend(matches, playerPUUID, champion)
		
		championMetrics = append(championMetrics, performance)
	}

	// Sort by games played (most played first)
	sort.Slice(championMetrics, func(i, j int) bool {
		return championMetrics[i].GamesPlayed > championMetrics[j].GamesPlayed
	})

	return championMetrics, nil
}

// calculateTrends analyzes performance trends over time
func (a *AnalyticsEngine) calculateTrends(matches []*riot.Match, playerPUUID string) (*TrendAnalysis, error) {
	if len(matches) < 10 {
		return &TrendAnalysis{
			TrendConfidence: 0.0,
			TrendPeriod:     "insufficient_data",
		}, nil
	}

	// Sort matches by time (most recent first)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Info.GameStartTimestamp > matches[j].Info.GameStartTimestamp
	})

	// Analyze recent vs older performance
	recentCount := a.config.RecentMatchesWindow
	if recentCount > len(matches)/2 {
		recentCount = len(matches) / 2
	}

	recentMatches := matches[:recentCount]
	olderMatches := matches[recentCount:]

	recentMetrics, _ := a.calculateCoreMetrics(recentMatches, playerPUUID)
	olderMetrics, _ := a.calculateCoreMetrics(olderMatches, playerPUUID)

	trends := &TrendAnalysis{
		TrendPeriod: fmt.Sprintf("last_%d_games", recentCount),
	}

	// Calculate trend directions
	trends.WinRateTrend = a.calculateTrendDirection(olderMetrics.WinRate, recentMetrics.WinRate)
	trends.KDATrend = a.calculateTrendDirection(olderMetrics.AverageKDA, recentMetrics.AverageKDA)
	trends.CSPerMinTrend = a.calculateTrendDirection(olderMetrics.CSPerMinute, recentMetrics.CSPerMinute)
	trends.VisionTrend = a.calculateTrendDirection(olderMetrics.AverageVision, recentMetrics.AverageVision)
	trends.DamageTrend = a.calculateTrendDirection(float64(olderMetrics.AverageDamage), float64(recentMetrics.AverageDamage))

	// Calculate overall performance trend
	trendScore := 0
	if trends.WinRateTrend == "improving" { trendScore++ }
	if trends.KDATrend == "improving" { trendScore++ }
	if trends.CSPerMinTrend == "improving" { trendScore++ }
	if trends.VisionTrend == "improving" { trendScore++ }
	if trends.DamageTrend == "improving" { trendScore++ }

	if trendScore >= 3 {
		trends.PerformanceTrend = "improving"
	} else if trendScore <= 1 {
		trends.PerformanceTrend = "declining"
	} else {
		trends.PerformanceTrend = "stable"
	}

	// Calculate trend confidence based on sample size and consistency
	trends.TrendConfidence = math.Min(float64(recentCount)/20.0, 1.0)
	trends.RecentWinRate = recentMetrics.WinRate

	return trends, nil
}

// calculatePerformanceScore calculates overall performance score (0-100)
func (a *AnalyticsEngine) calculatePerformanceScore(metrics *CoreMetrics, currentRank string) float64 {
	weights := a.config.MetricWeights
	threshold := a.getRankThresholds(currentRank)

	// Normalize each metric to 0-100 scale
	kdaScore := math.Min(metrics.AverageKDA/threshold.MinKDA*100, 100)
	csScore := math.Min(metrics.CSPerMinute/threshold.MinCSPerMin*100, 100)
	visionScore := math.Min(metrics.AverageVision/threshold.MinVisionScore*100, 100)
	damageScore := math.Min(metrics.DamageShare/threshold.MinDamageShare*100, 100)
	goldScore := math.Min(metrics.GoldEfficiency/threshold.MinGoldEff*100, 100)
	winRateScore := math.Min(metrics.WinRate/threshold.MinWinRate*100, 100)

	// Calculate weighted average
	totalWeight := weights.KDA + weights.CSPerMinute + weights.VisionScore + 
		weights.DamageShare + weights.GoldEfficiency + weights.WinRate

	weightedScore := (kdaScore*weights.KDA + 
		csScore*weights.CSPerMinute +
		visionScore*weights.VisionScore +
		damageScore*weights.DamageShare +
		goldScore*weights.GoldEfficiency +
		winRateScore*weights.WinRate) / totalWeight

	return math.Max(0, math.Min(100, weightedScore))
}

// Helper functions continue in next part...