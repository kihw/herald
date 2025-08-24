package summoner

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/herald-lol/herald/backend/internal/analytics"
	"github.com/herald-lol/herald/backend/internal/riot"
)

// Herald.lol Gaming Analytics - Summoner Service Helper Methods
// Supporting functions for the summoner analytics service

// Helper methods for SummonerService

func (s *SummonerService) processLiveGame(liveGame *riot.LiveGame) *LiveGameInfo {
	info := &LiveGameInfo{
		GameID:     liveGame.GameID,
		GameType:   liveGame.GameType,
		GameMode:   liveGame.GameMode,
		GameLength: liveGame.GameLength,
		MapID:      liveGame.MapID,
	}

	// Process participants
	var playerTeamParticipants []*LiveParticipant
	var enemyTeamParticipants []*LiveParticipant
	var playerTeamID int

	for _, participant := range liveGame.Participants {
		liveParticipant := &LiveParticipant{
			SummonerName:  participant.SummonerName,
			ChampionID:    participant.ChampionID,
			TeamID:        participant.TeamID,
			Spell1:        participant.Spell1ID,
			Spell2:        participant.Spell2ID,
			ProfileIconID: participant.ProfileIconID,
			IsBot:         participant.Bot,
		}

		// Set player team ID (assuming we know which participant is the target)
		// In real implementation, this would be determined by the request context
		if participant.TeamID == 100 {
			playerTeamID = 100
		}

		if participant.TeamID == playerTeamID {
			playerTeamParticipants = append(playerTeamParticipants, liveParticipant)
		} else {
			enemyTeamParticipants = append(enemyTeamParticipants, liveParticipant)
		}
	}

	// Create team info
	info.PlayerTeam = &LiveTeamInfo{
		TeamID:       playerTeamID,
		Participants: playerTeamParticipants,
	}

	enemyTeamID := 200
	if playerTeamID == 200 {
		enemyTeamID = 100
	}

	info.EnemyTeam = &LiveTeamInfo{
		TeamID:       enemyTeamID,
		Participants: enemyTeamParticipants,
	}

	// Add team composition analysis
	info.PlayerTeam.Composition = s.analyzeTeamComposition(playerTeamParticipants)
	info.EnemyTeam.Composition = s.analyzeTeamComposition(enemyTeamParticipants)

	return info
}

func (s *SummonerService) analyzeTeamComposition(participants []*LiveParticipant) *TeamComposition {
	composition := &TeamComposition{
		CompositionType:   "Balanced", // Default
		StrengthPhases:    []string{"Mid"},
		WeakPhases:        []string{},
		WinConditions:     []string{"Teamfight", "Objective Control"},
		ThreatsToWatch:    []string{},
		StrengthRating:    75.0, // Default rating
		SynergyRating:     70.0,
		FlexibilityRating: 80.0,
	}

	// In a real implementation, this would analyze champion synergies,
	// power spikes, win conditions, etc. based on the champion pool

	return composition
}

func (s *SummonerService) generateRecommendations(analysis *analytics.PlayerAnalysis, currentRank string) *RecommendationSummary {
	recommendations := &RecommendationSummary{
		ImmediateFocus:  []string{},
		ChampionPool:    []string{},
		SkillPriorities: []SkillPriority{},
		NextRankTarget:  s.getNextRankTarget(currentRank),
		ConfidenceLevel: 0.8,
	}

	// Generate immediate focus areas based on weakest metrics
	if analysis.CoreMetrics.CSPerMinute < 6.0 {
		recommendations.ImmediateFocus = append(recommendations.ImmediateFocus, "Improve farming efficiency - focus on last-hitting practice")
	}

	if analysis.CoreMetrics.AverageKDA < 2.0 {
		recommendations.ImmediateFocus = append(recommendations.ImmediateFocus, "Work on positioning and death avoidance")
	}

	if analysis.CoreMetrics.AverageVision < 15.0 {
		recommendations.ImmediateFocus = append(recommendations.ImmediateFocus, "Increase ward placement and vision control")
	}

	// Generate champion pool recommendations from best performing champions
	for _, champ := range analysis.ChampionMetrics {
		if champ.WinRate > 0.6 && champ.GamesPlayed >= 5 {
			recommendations.ChampionPool = append(recommendations.ChampionPool, champ.ChampionName)
		}
		if len(recommendations.ChampionPool) >= 5 {
			break
		}
	}

	// Generate skill priorities
	recommendations.SkillPriorities = []SkillPriority{
		{
			Skill:        "Farming",
			Priority:     s.getFarmingPriority(analysis.CoreMetrics.CSPerMinute),
			CurrentLevel: analysis.CoreMetrics.CSPerMinute * 10, // Scale to 0-100
			TargetLevel:  70.0,
			Improvement:  70.0 - (analysis.CoreMetrics.CSPerMinute * 10),
		},
		{
			Skill:        "Positioning",
			Priority:     s.getPositioningPriority(analysis.CoreMetrics.AverageKDA),
			CurrentLevel: min(analysis.CoreMetrics.AverageKDA*25, 100), // Scale to 0-100
			TargetLevel:  75.0,
			Improvement:  75.0 - min(analysis.CoreMetrics.AverageKDA*25, 100),
		},
		{
			Skill:        "Vision Control",
			Priority:     s.getVisionPriority(analysis.CoreMetrics.AverageVision),
			CurrentLevel: min(analysis.CoreMetrics.AverageVision*3, 100), // Scale to 0-100
			TargetLevel:  80.0,
			Improvement:  80.0 - min(analysis.CoreMetrics.AverageVision*3, 100),
		},
	}

	// Determine role focus
	recommendations.RoleFocus = s.determineBestRole(analysis.RoleMetrics)

	// Create training plan
	recommendations.TrainingPlan = &TrainingPlanSummary{
		Duration:       "2 weeks",
		DailyTime:      "45 minutes",
		KeyExercises:   s.generateKeyExercises(recommendations.ImmediateFocus),
		Milestones:     []string{"Improve CS/min by 1.0", "Reduce average deaths by 1", "Increase vision score by 5"},
		SuccessMetrics: []string{"Win 3 games in a row", "Achieve 7+ CS/min in ranked", "Place 15+ wards per game"},
	}

	// Estimate time to next rank
	recommendations.EstimatedTimeToRank = s.estimateTimeToRank(analysis, currentRank)

	return recommendations
}

func (s *SummonerService) compareSummoners(analysis1, analysis2 *analytics.PlayerAnalysis) *ComparisonResult {
	result := &ComparisonResult{
		MetricComparisons: make(map[string]*MetricComparison),
		OverallWinner:     "tie",
		WinnerMargin:      0.0,
	}

	// Compare core metrics
	metrics := map[string]struct {
		val1, val2 float64
	}{
		"KDA":          {analysis1.CoreMetrics.AverageKDA, analysis2.CoreMetrics.AverageKDA},
		"CS/min":       {analysis1.CoreMetrics.CSPerMinute, analysis2.CoreMetrics.CSPerMinute},
		"Vision":       {analysis1.CoreMetrics.AverageVision, analysis2.CoreMetrics.AverageVision},
		"Win Rate":     {analysis1.CoreMetrics.WinRate * 100, analysis2.CoreMetrics.WinRate * 100},
		"Damage Share": {analysis1.CoreMetrics.DamageShare * 100, analysis2.CoreMetrics.DamageShare * 100},
	}

	var summoner1Wins, summoner2Wins int
	var totalDifference float64

	for metric, values := range metrics {
		var winner string
		difference := values.val1 - values.val2
		absDiff := difference
		if absDiff < 0 {
			absDiff = -absDiff
		}

		if absDiff < 0.05*values.val1 { // Within 5%
			winner = "tie"
		} else if difference > 0 {
			winner = "summoner1"
			summoner1Wins++
		} else {
			winner = "summoner2"
			summoner2Wins++
		}

		significance := "minor"
		if absDiff > 0.2*values.val1 {
			significance = "major"
		} else if absDiff > 0.1*values.val1 {
			significance = "moderate"
		}

		result.MetricComparisons[metric] = &MetricComparison{
			Metric:       metric,
			Summoner1:    values.val1,
			Summoner2:    values.val2,
			Winner:       winner,
			Difference:   difference,
			Significance: significance,
		}

		totalDifference += absDiff
	}

	// Determine overall winner
	if summoner1Wins > summoner2Wins {
		result.OverallWinner = "summoner1"
		result.WinnerMargin = float64(summoner1Wins-summoner2Wins) / float64(len(metrics)) * 100
	} else if summoner2Wins > summoner1Wins {
		result.OverallWinner = "summoner2"
		result.WinnerMargin = float64(summoner2Wins-summoner1Wins) / float64(len(metrics)) * 100
	}

	// Generate strength areas
	result.StrengthAreas = s.generateComparisonStrengths(result.MetricComparisons)

	// Generate summary
	result.Summary = s.generateComparisonSummary(result)

	return result
}

func (s *SummonerService) analyzeTrendsByPeriod(matches []*riot.Match, playerPUUID string, timeWindows []string) map[string]*TrendPeriod {
	trends := make(map[string]*TrendPeriod)
	now := time.Now()

	for _, window := range timeWindows {
		var startTime time.Time
		switch window {
		case "7d":
			startTime = now.AddDate(0, 0, -7)
		case "30d":
			startTime = now.AddDate(0, 0, -30)
		case "90d":
			startTime = now.AddDate(0, 0, -90)
		case "season":
			startTime = now.AddDate(0, -6, 0) // Approximate season length
		default:
			continue
		}

		// Filter matches for this time period
		var periodMatches []*riot.Match
		for _, match := range matches {
			matchTime := time.Unix(match.Info.GameStartTimestamp/1000, 0)
			if matchTime.After(startTime) {
				periodMatches = append(periodMatches, match)
			}
		}

		if len(periodMatches) < 5 { // Need minimum matches for trend
			continue
		}

		// Analyze trends for this period
		trends[window] = s.analyzePeriodTrends(periodMatches, playerPUUID, startTime, now)
	}

	return trends
}

func (s *SummonerService) analyzePeriodTrends(matches []*riot.Match, playerPUUID string, startTime, endTime time.Time) *TrendPeriod {
	period := &TrendPeriod{
		Period:      fmt.Sprintf("%dd", int(endTime.Sub(startTime).Hours()/24)),
		StartDate:   startTime,
		EndDate:     endTime,
		GamesPlayed: len(matches),
		Metrics:     &TrendMetrics{},
	}

	// Split matches into early and late periods
	midPoint := startTime.Add(endTime.Sub(startTime) / 2)
	var earlyMatches, lateMatches []*riot.Match

	for _, match := range matches {
		matchTime := time.Unix(match.Info.GameStartTimestamp/1000, 0)
		if matchTime.Before(midPoint) {
			earlyMatches = append(earlyMatches, match)
		} else {
			lateMatches = append(lateMatches, match)
		}
	}

	if len(earlyMatches) > 0 && len(lateMatches) > 0 {
		// Calculate metrics for both periods
		earlyMetrics := s.calculatePeriodMetrics(earlyMatches, playerPUUID)
		lateMetrics := s.calculatePeriodMetrics(lateMatches, playerPUUID)

		// Create trend analysis
		period.Metrics.WinRate = s.createMetricTrend(earlyMetrics.WinRate, lateMetrics.WinRate)
		period.Metrics.KDA = s.createMetricTrend(earlyMetrics.AverageKDA, lateMetrics.AverageKDA)
		period.Metrics.CSPerMinute = s.createMetricTrend(earlyMetrics.CSPerMinute, lateMetrics.CSPerMinute)
		period.Metrics.VisionScore = s.createMetricTrend(earlyMetrics.AverageVision, lateMetrics.AverageVision)
		period.Metrics.DamageShare = s.createMetricTrend(earlyMetrics.DamageShare, lateMetrics.DamageShare)

		// Determine trend strength
		period.TrendStrength = s.calculateTrendStrength(period.Metrics)
	}

	return period
}

func (s *SummonerService) calculatePeriodMetrics(matches []*riot.Match, playerPUUID string) *analytics.CoreMetrics {
	// Use analytics engine to calculate metrics for this period
	analyticsRequest := &analytics.PlayerAnalysisRequest{
		PlayerPUUID: playerPUUID,
		Matches:     matches,
	}

	// Simplified metrics calculation for trend analysis
	var totalKills, totalDeaths, totalAssists, totalCS int
	var totalDuration int
	var wins int

	for _, match := range matches {
		participant := s.findPlayerInMatch(match, playerPUUID)
		if participant == nil {
			continue
		}

		totalKills += participant.Kills
		totalDeaths += participant.Deaths
		totalAssists += participant.Assists
		totalCS += participant.TotalMinionsKilled + participant.NeutralMinionsKilled
		totalDuration += match.Info.GameDuration

		if participant.Win {
			wins++
		}
	}

	if totalDeaths == 0 {
		totalDeaths = 1
	}

	return &analytics.CoreMetrics{
		AverageKDA:    float64(totalKills+totalAssists) / float64(totalDeaths),
		CSPerMinute:   float64(totalCS) / (float64(totalDuration) / 60.0),
		AverageVision: 15.0, // Placeholder - would calculate from match data
		WinRate:       float64(wins) / float64(len(matches)),
		DamageShare:   0.25, // Placeholder - would calculate from match data
	}
}

func (s *SummonerService) findPlayerInMatch(match *riot.Match, playerPUUID string) *riot.Participant {
	for i := range match.Info.Participants {
		if match.Info.Participants[i].PUUID == playerPUUID {
			return &match.Info.Participants[i]
		}
	}
	return nil
}

func (s *SummonerService) createMetricTrend(startValue, endValue float64) *MetricTrend {
	change := endValue - startValue
	percentChange := 0.0
	if startValue != 0 {
		percentChange = (change / startValue) * 100
	}

	direction := "stable"
	significance := "minor"

	if change > 0.05*startValue {
		direction = "up"
		if change > 0.15*startValue {
			significance = "major"
		} else if change > 0.10*startValue {
			significance = "moderate"
		}
	} else if change < -0.05*startValue {
		direction = "down"
		if change < -0.15*startValue {
			significance = "major"
		} else if change < -0.10*startValue {
			significance = "moderate"
		}
	}

	return &MetricTrend{
		StartValue:    startValue,
		EndValue:      endValue,
		Change:        change,
		PercentChange: percentChange,
		Direction:     direction,
		Significance:  significance,
	}
}

func (s *SummonerService) calculateTrendStrength(metrics *TrendMetrics) string {
	positiveCount := 0
	negativeCount := 0

	trends := []*MetricTrend{
		metrics.WinRate, metrics.KDA, metrics.CSPerMinute,
		metrics.VisionScore, metrics.DamageShare,
	}

	for _, trend := range trends {
		if trend.Direction == "up" {
			positiveCount++
		} else if trend.Direction == "down" {
			negativeCount++
		}
	}

	if positiveCount >= 4 {
		return "strong_positive"
	} else if positiveCount >= 3 {
		return "positive"
	} else if negativeCount >= 4 {
		return "strong_negative"
	} else if negativeCount >= 3 {
		return "negative"
	}
	return "stable"
}

func (s *SummonerService) generateEnhancedInsights(analysis *analytics.PlayerAnalysis, insightTypes []string) *EnhancedInsights {
	insights := &EnhancedInsights{}

	for _, insightType := range insightTypes {
		switch insightType {
		case "performance":
			insights.PerformanceInsights = s.generatePerformanceInsights(analysis)
		case "coaching":
			insights.CoachingInsights = s.generateCoachingInsights(analysis)
		case "predictions":
			insights.PredictiveInsights = s.generatePredictiveInsights(analysis)
		case "meta":
			insights.MetaInsights = s.generateMetaInsights(analysis)
		}
	}

	// Generate personalized advice
	insights.PersonalizedAdvice = s.generatePersonalizedAdvice(analysis)

	return insights
}

// Cache methods

func (s *SummonerService) getCachedAnalysis(ctx context.Context, request *SummonerAnalysisRequest) (*SummonerAnalysisResponse, error) {
	key := s.getCacheKey(request)
	data, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var response SummonerAnalysisResponse
	if err := json.Unmarshal([]byte(data), &response); err != nil {
		return nil, err
	}

	response.CacheHit = true
	response.DataFreshness = "cached"
	return &response, nil
}

func (s *SummonerService) cacheAnalysis(ctx context.Context, request *SummonerAnalysisRequest, response *SummonerAnalysisResponse) {
	key := s.getCacheKey(request)
	data, err := json.Marshal(response)
	if err != nil {
		return
	}

	s.redis.Set(ctx, key, data, s.config.CacheAnalyticsTTL)
}

func (s *SummonerService) getCacheKey(request *SummonerAnalysisRequest) string {
	return fmt.Sprintf("summoner_analysis:%s:%s:%s:%s",
		request.Region, strings.ToLower(request.SummonerName),
		request.AnalysisType, request.TimeFrame)
}

func (s *SummonerService) createProgressTracker(requestID string) *AnalysisProgressTracker {
	return &AnalysisProgressTracker{
		RequestID:  requestID,
		StartedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Progress:   0,
		IsComplete: false,
	}
}

// Utility functions

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func (s *SummonerService) getNextRankTarget(currentRank string) string {
	rankProgression := map[string]string{
		"IRON IV": "IRON III", "IRON III": "IRON II", "IRON II": "IRON I", "IRON I": "BRONZE IV",
		"BRONZE IV": "BRONZE III", "BRONZE III": "BRONZE II", "BRONZE II": "BRONZE I", "BRONZE I": "SILVER IV",
		"SILVER IV": "SILVER III", "SILVER III": "SILVER II", "SILVER II": "SILVER I", "SILVER I": "GOLD IV",
		"GOLD IV": "GOLD III", "GOLD III": "GOLD II", "GOLD II": "GOLD I", "GOLD I": "PLATINUM IV",
		"PLATINUM IV": "PLATINUM III", "PLATINUM III": "PLATINUM II", "PLATINUM II": "PLATINUM I", "PLATINUM I": "EMERALD IV",
		"EMERALD IV": "EMERALD III", "EMERALD III": "EMERALD II", "EMERALD II": "EMERALD I", "EMERALD I": "DIAMOND IV",
		"DIAMOND IV": "DIAMOND III", "DIAMOND III": "DIAMOND II", "DIAMOND II": "DIAMOND I", "DIAMOND I": "MASTER I",
		"MASTER I": "GRANDMASTER I", "GRANDMASTER I": "CHALLENGER I",
	}

	if next, exists := rankProgression[currentRank]; exists {
		return next
	}
	return "GOLD IV" // Default target
}

func (s *SummonerService) getFarmingPriority(csPerMin float64) int {
	if csPerMin < 5.0 {
		return 5 // Highest priority
	} else if csPerMin < 6.5 {
		return 3
	}
	return 1
}

func (s *SummonerService) getPositioningPriority(kda float64) int {
	if kda < 1.5 {
		return 5
	} else if kda < 2.5 {
		return 3
	}
	return 2
}

func (s *SummonerService) getVisionPriority(vision float64) int {
	if vision < 12.0 {
		return 4
	} else if vision < 18.0 {
		return 2
	}
	return 1
}

func (s *SummonerService) determineBestRole(roleMetrics map[string]*analytics.RolePerformance) string {
	bestRole := "MIDDLE" // Default
	bestRating := 0.0

	for role, performance := range roleMetrics {
		if performance.GamesPlayed >= 3 && performance.PerformanceRating > bestRating {
			bestRole = role
			bestRating = performance.PerformanceRating
		}
	}

	return strings.ToLower(bestRole)
}

func (s *SummonerService) generateKeyExercises(focuses []string) []string {
	exercises := []string{}

	for _, focus := range focuses {
		if strings.Contains(strings.ToLower(focus), "farm") {
			exercises = append(exercises, "30 minutes CS practice in training mode daily")
		}
		if strings.Contains(strings.ToLower(focus), "position") {
			exercises = append(exercises, "Review 3 deaths per game - analyze positioning mistakes")
		}
		if strings.Contains(strings.ToLower(focus), "vision") {
			exercises = append(exercises, "Place 2+ wards every back, buy 1 control ward per game")
		}
	}

	return exercises
}

func (s *SummonerService) estimateTimeToRank(analysis *analytics.PlayerAnalysis, currentRank string) string {
	// Simplified estimation based on performance gap
	gap := 75.0 - analysis.PerformanceScore // Target 75+ for next rank

	if gap <= 5 {
		return "1-2 weeks"
	} else if gap <= 15 {
		return "3-4 weeks"
	} else if gap <= 25 {
		return "1-2 months"
	}
	return "2-3 months"
}

// Placeholder methods for advanced insights
func (s *SummonerService) generateComparisonStrengths(comparisons map[string]*MetricComparison) *ComparisonStrengths {
	return &ComparisonStrengths{
		Summoner1Strengths: []string{"KDA Control", "Farming"},
		Summoner2Strengths: []string{"Vision Control", "Damage Output"},
		SharedStrengths:    []string{"Consistent Performance"},
	}
}

func (s *SummonerService) generateComparisonSummary(result *ComparisonResult) string {
	if result.OverallWinner == "tie" {
		return "Both summoners show similar overall performance with different strengths"
	}
	return fmt.Sprintf("Summoner %s shows stronger performance with %.1f%% advantage",
		result.OverallWinner, result.WinnerMargin)
}

func (s *SummonerService) generatePerformanceInsights(analysis *analytics.PlayerAnalysis) *PerformanceInsights {
	return &PerformanceInsights{
		CurrentForm:       "Stable",
		ConsistencyRating: 75.0,
		ClutchFactor:      60.0,
		AdaptabilityScore: 70.0,
	}
}

func (s *SummonerService) generateCoachingInsights(analysis *analytics.PlayerAnalysis) *CoachingInsights {
	return &CoachingInsights{
		LearningStyle:      "Practice-focused",
		OptimalPlayTimes:   []string{"Evening", "Weekend"},
		RecommendedRoutine: "3 ranked games with 15min practice between",
	}
}

func (s *SummonerService) generatePredictiveInsights(analysis *analytics.PlayerAnalysis) *PredictiveInsights {
	return &PredictiveInsights{
		RankPrediction: &RankPrediction{
			PredictedRank: "GOLD I",
			Confidence:    0.75,
			TimeFrame:     "1 month",
		},
	}
}

func (s *SummonerService) generateMetaInsights(analysis *analytics.PlayerAnalysis) *MetaInsights {
	return &MetaInsights{
		MetaAlignment: "Good",
		MetaChampions: []string{"Jinx", "Caitlyn", "Jhin"},
	}
}

func (s *SummonerService) generatePersonalizedAdvice(analysis *analytics.PlayerAnalysis) []PersonalizedTip {
	return []PersonalizedTip{
		{
			Tip:       "Focus on farming improvements to increase gold income",
			Priority:  5,
			Category:  "Economy",
			Reasoning: "Your CS/min is below your rank average",
		},
	}
}
