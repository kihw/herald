package analytics

import (
	"fmt"
	"math"
	"strings"

	"github.com/herald-lol/herald/backend/internal/riot"
)

// Herald.lol Gaming Analytics - Engine Helper Functions
// Supporting functions for the core analytics engine

// Helper functions for analytics engine calculations

// findPlayerInMatch finds the participant data for a specific player in a match
func (a *AnalyticsEngine) findPlayerInMatch(match *riot.Match, playerPUUID string) *riot.Participant {
	for i := range match.Info.Participants {
		if match.Info.Participants[i].PUUID == playerPUUID {
			return &match.Info.Participants[i]
		}
	}
	return nil
}

// calculateTeamDamage calculates total team damage for damage share calculation
func (a *AnalyticsEngine) calculateTeamDamage(match *riot.Match, teamID int) int64 {
	var totalDamage int64
	for _, participant := range match.Info.Participants {
		if participant.TeamID == teamID {
			totalDamage += int64(participant.TotalDamageDealtToChampions)
		}
	}
	return totalDamage
}

// normalizeRole converts various role representations to standard roles
func (a *AnalyticsEngine) normalizeRole(role string) string {
	role = strings.ToUpper(strings.TrimSpace(role))

	roleMap := map[string]string{
		"TOP":     "TOP",
		"JUNGLE":  "JUNGLE",
		"MIDDLE":  "MIDDLE",
		"MID":     "MIDDLE",
		"BOTTOM":  "BOTTOM",
		"BOT":     "BOTTOM",
		"ADC":     "BOTTOM",
		"UTILITY": "SUPPORT",
		"SUPPORT": "SUPPORT",
		"SUPP":    "SUPPORT",
	}

	if normalized, exists := roleMap[role]; exists {
		return normalized
	}
	return "UNKNOWN"
}

// calculateRoleRating calculates performance rating for a specific role
func (a *AnalyticsEngine) calculateRoleRating(performance *RolePerformance, role string) float64 {
	expectations := a.config.RoleExpectations[role]
	if expectations == nil {
		return 50.0 // Default rating
	}

	var score float64
	var factors int

	// KDA comparison
	if expectations.ExpectedKDA > 0 {
		score += math.Min(performance.AverageKDA/expectations.ExpectedKDA*100, 200)
		factors++
	}

	// CS comparison
	if expectations.ExpectedCS > 0 {
		score += math.Min(performance.AverageCS/expectations.ExpectedCS*100, 200)
		factors++
	}

	// Damage comparison
	if expectations.ExpectedDamage > 0 {
		score += math.Min(float64(performance.AverageDamage)/expectations.ExpectedDamage*100, 200)
		factors++
	}

	// Vision comparison
	if expectations.ExpectedVision > 0 {
		score += math.Min(performance.AverageVision/expectations.ExpectedVision*100, 200)
		factors++
	}

	if factors > 0 {
		return math.Max(0, math.Min(100, score/float64(factors)))
	}
	return 50.0
}

// calculateMasteryLevel estimates mastery level based on performance
func (a *AnalyticsEngine) calculateMasteryLevel(performance *ChampionPerformance) int {
	// Estimate mastery level based on games played and performance
	baseLevel := performance.GamesPlayed / 5 // ~5 games per mastery level initially

	// Adjust based on performance
	if performance.WinRate > 0.6 {
		baseLevel += 2
	} else if performance.WinRate > 0.5 {
		baseLevel += 1
	}

	if performance.AverageKDA > 3.0 {
		baseLevel += 2
	} else if performance.AverageKDA > 2.0 {
		baseLevel += 1
	}

	return int(math.Max(1, math.Min(float64(baseLevel), 7)))
}

// calculateChampionTrend calculates performance trend for a specific champion
func (a *AnalyticsEngine) calculateChampionTrend(matches []*riot.Match, playerPUUID, champion string) string {
	var championMatches []*riot.Match

	// Filter matches for this champion
	for _, match := range matches {
		participant := a.findPlayerInMatch(match, playerPUUID)
		if participant != nil && participant.ChampionName == champion {
			championMatches = append(championMatches, match)
		}
	}

	if len(championMatches) < 4 {
		return "stable" // Not enough data
	}

	// Compare recent vs older performance
	mid := len(championMatches) / 2
	recent := championMatches[:mid]
	older := championMatches[mid:]

	recentWins := 0
	olderWins := 0

	for _, match := range recent {
		participant := a.findPlayerInMatch(match, playerPUUID)
		if participant != nil && participant.Win {
			recentWins++
		}
	}

	for _, match := range older {
		participant := a.findPlayerInMatch(match, playerPUUID)
		if participant != nil && participant.Win {
			olderWins++
		}
	}

	recentWinRate := float64(recentWins) / float64(len(recent))
	olderWinRate := float64(olderWins) / float64(len(older))

	return a.calculateTrendDirection(olderWinRate, recentWinRate)
}

// calculateTrendDirection determines trend direction between two values
func (a *AnalyticsEngine) calculateTrendDirection(oldValue, newValue float64) string {
	difference := newValue - oldValue
	threshold := 0.05 // 5% change threshold

	if difference > threshold {
		return "improving"
	} else if difference < -threshold {
		return "declining"
	}
	return "stable"
}

// getRankThresholds returns performance thresholds for a given rank
func (a *AnalyticsEngine) getRankThresholds(rank string) *RankThresholds {
	rank = strings.ToUpper(rank)

	if thresholds, exists := a.config.RankThresholds[rank]; exists {
		return thresholds
	}

	// Default to Silver thresholds if rank not found
	if thresholds, exists := a.config.RankThresholds["SILVER"]; exists {
		return thresholds
	}

	// Fallback thresholds
	return &RankThresholds{
		MinKDA:         1.5,
		MinCSPerMin:    6.0,
		MinVisionScore: 15.0,
		MinDamageShare: 0.20,
		MinGoldEff:     0.85,
		MinWinRate:     0.50,
	}
}

// generateInsights generates AI-powered insights based on analysis
func (a *AnalyticsEngine) generateInsights(analysis *PlayerAnalysis, currentRank string) (*GameInsights, error) {
	insights := &GameInsights{
		StrengthAreas:    []string{},
		ImprovementAreas: []string{},
		CoachingTips:     []string{},
		NextGoals:        []string{},
	}

	metrics := analysis.CoreMetrics
	threshold := a.getRankThresholds(currentRank)

	// Identify strength areas
	if metrics.AverageKDA >= threshold.MinKDA*1.2 {
		insights.StrengthAreas = append(insights.StrengthAreas, "Excellent KDA ratio - good at avoiding deaths while contributing to kills")
	}
	if metrics.CSPerMinute >= threshold.MinCSPerMin*1.1 {
		insights.StrengthAreas = append(insights.StrengthAreas, "Strong farming skills - consistently good CS per minute")
	}
	if metrics.AverageVision >= threshold.MinVisionScore*1.1 {
		insights.StrengthAreas = append(insights.StrengthAreas, "Great vision control - contributing well to team vision")
	}
	if metrics.DamageShare >= threshold.MinDamageShare*1.1 {
		insights.StrengthAreas = append(insights.StrengthAreas, "High damage contribution - carrying team fights well")
	}

	// Identify improvement areas
	if metrics.AverageKDA < threshold.MinKDA {
		insights.ImprovementAreas = append(insights.ImprovementAreas, "KDA needs improvement - focus on positioning and death reduction")
		insights.CoachingTips = append(insights.CoachingTips, "Practice safer positioning in team fights and avoid overextending in lane")
	}
	if metrics.CSPerMinute < threshold.MinCSPerMin {
		insights.ImprovementAreas = append(insights.ImprovementAreas, "CS per minute below expected - improve farming efficiency")
		insights.CoachingTips = append(insights.CoachingTips, "Focus on last-hitting practice and wave management fundamentals")
	}
	if metrics.AverageVision < threshold.MinVisionScore {
		insights.ImprovementAreas = append(insights.ImprovementAreas, "Vision score is low - increase ward placement and control ward usage")
		insights.CoachingTips = append(insights.CoachingTips, "Buy more control wards and focus on strategic vision placement near objectives")
	}

	// Determine playstyle profile
	insights.PlaystyleProfile = a.determinePlaystyle(analysis)

	// Generate champion recommendations
	insights.RecommendedChamps = a.generateChampionRecommendations(analysis)

	// Set next goals
	insights.NextGoals = a.generateNextGoals(analysis, currentRank)

	// Assess skill level
	insights.SkillLevel = a.assessSkillLevel(analysis, currentRank)

	// Calculate confidence based on sample size
	insights.Confidence = math.Min(float64(analysis.TotalMatches)/50.0, 1.0)

	return insights, nil
}

// determinePlaystyle analyzes play pattern to determine playstyle
func (a *AnalyticsEngine) determinePlaystyle(analysis *PlayerAnalysis) string {
	metrics := analysis.CoreMetrics

	// Analyze aggression level
	killParticipation := metrics.AverageKills + metrics.AverageAssists
	deathRate := metrics.AverageDeaths

	aggressionScore := killParticipation / math.Max(deathRate, 0.1)

	if aggressionScore > 4.0 {
		return "Aggressive"
	} else if aggressionScore < 2.0 {
		return "Passive"
	}
	return "Balanced"
}

// generateChampionRecommendations suggests champions based on performance patterns
func (a *AnalyticsEngine) generateChampionRecommendations(analysis *PlayerAnalysis) []string {
	recommendations := []string{}

	// Find best performing champions
	if len(analysis.ChampionMetrics) > 0 {
		// Sort by combined performance score
		type championScore struct {
			name  string
			score float64
		}

		var scores []championScore
		for _, champ := range analysis.ChampionMetrics {
			// Calculate composite score
			score := champ.WinRate*0.5 + (champ.AverageKDA/5.0)*0.3 + (champ.CSPerMinute/10.0)*0.2
			scores = append(scores, championScore{champ.ChampionName, score})
		}

		// Add top performing champions to recommendations
		for i, score := range scores {
			if i < 3 && score.score > 0.6 { // Top 3 if score > 60%
				recommendations = append(recommendations, fmt.Sprintf("Continue playing %s - showing strong performance", score.name))
			}
		}
	}

	// Role-based recommendations
	if len(analysis.RoleMetrics) > 0 {
		for role, rolePerf := range analysis.RoleMetrics {
			if rolePerf.PerformanceRating > 70 {
				recommendations = append(recommendations, fmt.Sprintf("Focus on %s role - showing strong performance", strings.ToLower(role)))
			}
		}
	}

	return recommendations
}

// generateNextGoals creates specific improvement goals
func (a *AnalyticsEngine) generateNextGoals(analysis *PlayerAnalysis, currentRank string) []string {
	goals := []string{}
	metrics := analysis.CoreMetrics
	threshold := a.getRankThresholds(currentRank)

	// KDA improvement goal
	if metrics.AverageKDA < threshold.MinKDA {
		targetKDA := threshold.MinKDA + 0.3
		goals = append(goals, fmt.Sprintf("Improve KDA to %.1f (currently %.1f)", targetKDA, metrics.AverageKDA))
	}

	// CS improvement goal
	if metrics.CSPerMinute < threshold.MinCSPerMin {
		targetCS := threshold.MinCSPerMin + 0.5
		goals = append(goals, fmt.Sprintf("Increase CS/min to %.1f (currently %.1f)", targetCS, metrics.CSPerMinute))
	}

	// Win rate goal
	if metrics.WinRate < 0.55 {
		goals = append(goals, fmt.Sprintf("Achieve 55%+ win rate (currently %.0f%%)", metrics.WinRate*100))
	}

	// Vision goal
	if metrics.AverageVision < threshold.MinVisionScore {
		targetVision := threshold.MinVisionScore + 5
		goals = append(goals, fmt.Sprintf("Improve vision score to %.0f per game", targetVision))
	}

	return goals
}

// assessSkillLevel provides skill assessment based on performance
func (a *AnalyticsEngine) assessSkillLevel(analysis *PlayerAnalysis, currentRank string) string {
	score := analysis.PerformanceScore

	// Map performance score to skill assessment
	switch {
	case score >= 90:
		return "Challenger Level"
	case score >= 80:
		return "Master Level"
	case score >= 70:
		return "Diamond Level"
	case score >= 60:
		return "Platinum Level"
	case score >= 50:
		return "Gold Level"
	case score >= 40:
		return "Silver Level"
	case score >= 30:
		return "Bronze Level"
	default:
		return "Iron Level"
	}
}

// Accumulator structs for calculations
type roleAccumulator struct {
	matches       int
	wins          int
	totalKills    int
	totalDeaths   int
	totalAssists  int
	totalCS       int
	totalDamage   int64
	totalGold     int
	totalVision   int
	totalDuration int
}

type championAccumulator struct {
	championName  string
	matches       int
	wins          int
	totalKills    int
	totalDeaths   int
	totalAssists  int
	totalCS       int
	totalDuration int
}

// Performance calculation utilities

// CalculateSkillGap calculates skill gap between current and target rank
func (a *AnalyticsEngine) CalculateSkillGap(currentMetrics *CoreMetrics, targetRank string) *SkillGap {
	_ = a.getRankThresholds("SILVER") // Default current - unused for now
	targetThreshold := a.getRankThresholds(targetRank)

	return &SkillGap{
		KDAGap:         math.Max(0, targetThreshold.MinKDA-currentMetrics.AverageKDA),
		CSGap:          math.Max(0, targetThreshold.MinCSPerMin-currentMetrics.CSPerMinute),
		VisionGap:      math.Max(0, targetThreshold.MinVisionScore-currentMetrics.AverageVision),
		DamageGap:      math.Max(0, targetThreshold.MinDamageShare-currentMetrics.DamageShare),
		WinRateGap:     math.Max(0, targetThreshold.MinWinRate-currentMetrics.WinRate),
		OverallGap:     a.calculateOverallGap(currentMetrics, targetRank),
		EstimatedGames: a.estimateGamesNeeded(currentMetrics, targetRank),
	}
}

func (a *AnalyticsEngine) calculateOverallGap(metrics *CoreMetrics, targetRank string) float64 {
	targetScore := a.calculatePerformanceScore(metrics, targetRank)
	currentScore := a.calculatePerformanceScore(metrics, "SILVER")
	return math.Max(0, targetScore-currentScore)
}

func (a *AnalyticsEngine) estimateGamesNeeded(metrics *CoreMetrics, targetRank string) int {
	gap := a.calculateOverallGap(metrics, targetRank)
	// Rough estimate: 1 point improvement per 2 games with focused improvement
	return int(gap * 2)
}
