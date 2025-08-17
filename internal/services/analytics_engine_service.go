package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"time"

	"lol-match-exporter/internal/db"
	"lol-match-exporter/internal/models"
)

// AnalyticsEngineService handles performance analysis and statistics
type AnalyticsEngineService struct {
	db *sql.DB
}

// NewAnalyticsEngineService creates a new analytics engine service
func NewAnalyticsEngineService(database *db.Database) *AnalyticsEngineService {
	var sqlDB *sql.DB
	if database != nil {
		sqlDB = database.DB
	}
	return &AnalyticsEngineService{
		db: sqlDB,
	}
}

// GeneratePeriodStats generates comprehensive statistics for a time period
func (aes *AnalyticsEngineService) GeneratePeriodStats(userID int, period string) (*models.PeriodStats, error) {
	matches, err := aes.getMatchesForPeriod(userID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches: %w", err)
	}

	if len(matches) == 0 {
		return &models.PeriodStats{
			Period:          period,
			TotalGames:      0,
			WinRate:         0.0,
			AvgKDA:          0.0,
			BestRole:        "",
			WorstRole:       "",
			TopChampions:    []models.ChampionPerformance{},
			RolePerformance: map[string]models.PerformanceMetrics{},
			RecentTrend:     "stable",
			Suggestions:     []string{},
		}, nil
	}

	// Calculate basic stats
	totalGames := len(matches)
	wins := 0
	for _, match := range matches {
		if aes.extractWinStatus(match) {
			wins++
		}
	}
	winRate := float64(wins) / float64(totalGames) * 100

	// Calculate KDA
	totalKills, totalDeaths, totalAssists := aes.calculateTotalKDA(matches)
	avgKDA := (totalKills + totalAssists) / math.Max(totalDeaths, 1)

	// Analyze by role
	rolePerformance := aes.analyzeRolePerformance(matches)
	
	var bestRole, worstRole string
	if len(rolePerformance) > 0 {
		bestScore := -1.0
		worstScore := 999.0
		
		for role, metrics := range rolePerformance {
			if metrics.PerformanceScore > bestScore {
				bestScore = metrics.PerformanceScore
				bestRole = role
			}
			if metrics.PerformanceScore < worstScore {
				worstScore = metrics.PerformanceScore
				worstRole = role
			}
		}
	}

	// Top champions
	topChampions := aes.getTopChampions(matches, 5)

	// Recent trend
	recentTrend := aes.calculateRecentTrend(matches)

	// Generate suggestions
	suggestions := aes.generatePeriodSuggestions(matches, period, rolePerformance)

	return &models.PeriodStats{
		Period:          period,
		TotalGames:      totalGames,
		WinRate:         winRate,
		AvgKDA:          avgKDA,
		BestRole:        bestRole,
		WorstRole:       worstRole,
		TopChampions:    topChampions,
		RolePerformance: rolePerformance,
		RecentTrend:     recentTrend,
		Suggestions:     suggestions,
	}, nil
}

// CalculateRolePerformance calculates detailed performance metrics for a specific role
func (aes *AnalyticsEngineService) CalculateRolePerformance(userID int, role, period string) (*models.PerformanceMetrics, error) {
	matches, err := aes.getMatchesForPeriod(userID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches: %w", err)
	}

	// Filter matches by role
	var roleMatches []map[string]interface{}
	for _, match := range matches {
		if aes.extractRole(match) == role {
			roleMatches = append(roleMatches, match)
		}
	}

	if len(roleMatches) == 0 {
		return &models.PerformanceMetrics{}, nil
	}

	metrics := aes.calculateMetricsFromMatches(roleMatches)
	metrics.TrendDirection = models.TrendDirection(aes.calculateRoleTrend(roleMatches))

	return &metrics, nil
}

// AnalyzeChampionMastery analyzes mastery and performance for a specific champion
func (aes *AnalyticsEngineService) AnalyzeChampionMastery(userID, championID int, period string) (*models.ChampionMasteryAnalysis, error) {
	matches, err := aes.getMatchesForPeriod(userID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches: %w", err)
	}

	// Filter matches by champion
	var championMatches []map[string]interface{}
	for _, match := range matches {
		if aes.extractChampionID(match) == championID {
			championMatches = append(championMatches, match)
		}
	}

	if len(championMatches) == 0 {
		return nil, fmt.Errorf("no matches found for this champion")
	}

	metrics := aes.calculateMetricsFromMatches(championMatches)

	// Calculate mastery score (0-100)
	masteryScore := aes.calculateMasteryScore(championMatches)

	// Find best and worst performances
	bestGame := aes.findBestGame(championMatches)
	worstGame := aes.findWorstGame(championMatches)

	// Generate improvement suggestions
	suggestions := aes.generateChampionSuggestions(championMatches, championID)

	// Analyze skill progression
	skillProgression := aes.analyzeSkillProgression(championMatches)

	championName := aes.extractChampionName(championMatches[0])

	return &models.ChampionMasteryAnalysis{
		ChampionID:             championID,
		ChampionName:           championName,
		GamesPlayed:            metrics.GamesPlayed,
		PerformanceMetrics:     metrics,
		MasteryScore:           masteryScore,
		BestGame:               bestGame,
		WorstGame:              worstGame,
		ImprovementSuggestions: suggestions,
		SkillProgression:       skillProgression,
	}, nil
}

// GenerateImprovementSuggestions generates personalized improvement suggestions
func (aes *AnalyticsEngineService) GenerateImprovementSuggestions(userID int) ([]models.ImprovementSuggestion, error) {
	var suggestions []models.ImprovementSuggestion

	// Analyze recent performance (last 7 days)
	recentMatches, err := aes.getMatchesForPeriod(userID, "week")
	if err != nil {
		return nil, err
	}

	monthMatches, err := aes.getMatchesForPeriod(userID, "month")
	if err != nil {
		return nil, err
	}

	if len(recentMatches) < 3 {
		return suggestions, nil
	}

	// Performance trend analysis
	trendSuggestions := aes.analyzePerformanceTrends(recentMatches, monthMatches)
	suggestions = append(suggestions, trendSuggestions...)

	// Role-specific suggestions
	roleSuggestions := aes.analyzeRoleWeaknesses(monthMatches)
	suggestions = append(suggestions, roleSuggestions...)

	// Champion pool suggestions
	poolSuggestions := aes.analyzeChampionPool(monthMatches)
	suggestions = append(suggestions, poolSuggestions...)

	// Meta adaptation suggestions
	metaSuggestions := aes.generateMetaSuggestions(recentMatches)
	suggestions = append(suggestions, metaSuggestions...)

	// Sort by priority and return top 10
	sort.Slice(suggestions, func(i, j int) bool {
		if suggestions[i].Priority == suggestions[j].Priority {
			return suggestions[i].Type < suggestions[j].Type
		}
		return suggestions[i].Priority < suggestions[j].Priority
	})

	if len(suggestions) > 10 {
		suggestions = suggestions[:10]
	}

	return suggestions, nil
}

// CalculatePerformanceTrends calculates detailed performance trends over time
func (aes *AnalyticsEngineService) CalculatePerformanceTrends(userID int) (*models.PerformanceTrends, error) {
	// Get matches for different periods
	todayMatches, _ := aes.getMatchesForPeriod(userID, "today")
	weekMatches, _ := aes.getMatchesForPeriod(userID, "week")
	monthMatches, _ := aes.getMatchesForPeriod(userID, "month")
	seasonMatches, _ := aes.getMatchesForPeriod(userID, "season")

	return &models.PerformanceTrends{
		DailyTrend:         aes.calculateTrendMetrics(todayMatches),
		WeeklyTrend:        aes.calculateTrendMetrics(weekMatches),
		MonthlyTrend:       aes.calculateTrendMetrics(monthMatches),
		SeasonalTrend:      aes.calculateTrendMetrics(seasonMatches),
		ImprovementVelocity: aes.calculateImprovementVelocity(weekMatches, monthMatches),
		ConsistencyScore:   aes.calculateConsistencyScore(monthMatches),
		PeakPerformance:    aes.findPeakPerformancePeriod(seasonMatches),
	}, nil
}

// UpdateChampionStats updates champion stats in database for all champions
func (aes *AnalyticsEngineService) UpdateChampionStats(userID int, period string) error {
	matches, err := aes.getMatchesForPeriod(userID, period)
	if err != nil {
		return fmt.Errorf("failed to get matches: %w", err)
	}

	// Group matches by champion and role
	championGroups := make(map[string][]map[string]interface{})
	
	for _, match := range matches {
		championID := aes.extractChampionID(match)
		role := aes.extractRole(match)
		season := aes.extractSeason(match)
		
		key := fmt.Sprintf("%d-%s-%s", championID, role, season)
		championGroups[key] = append(championGroups[key], match)
	}

	// Calculate and save stats for each champion/role combination
	for key, champMatches := range championGroups {
		parts := strings.Split(key, "-")
		if len(parts) != 3 {
			continue
		}
		
		championID := aes.extractChampionID(champMatches[0])
		championName := aes.extractChampionName(champMatches[0])
		role := parts[1]
		season := parts[2]
		
		metrics := aes.calculateMetricsFromMatches(champMatches)
		lastPlayed := aes.getLastPlayedDate(champMatches)
		
		err := aes.saveChampionStats(userID, championID, championName, role, season, period, metrics, lastPlayed)
		if err != nil {
			log.Printf("Failed to save champion stats: %v", err)
		}
	}

	return nil
}

// Helper methods

func (aes *AnalyticsEngineService) extractStat(match map[string]interface{}, stat string) float64 {
	participantDataStr, ok := match["participant_data"].(string)
	if !ok {
		return 0.0
	}

	var participantData map[string]interface{}
	if err := json.Unmarshal([]byte(participantDataStr), &participantData); err != nil {
		return 0.0
	}

	if value, ok := participantData[stat]; ok {
		if floatVal, ok := value.(float64); ok {
			return floatVal
		}
	}
	return 0.0
}

func (aes *AnalyticsEngineService) extractWinStatus(match map[string]interface{}) bool {
	participantDataStr, ok := match["participant_data"].(string)
	if !ok {
		return false
	}

	var participantData map[string]interface{}
	if err := json.Unmarshal([]byte(participantDataStr), &participantData); err != nil {
		return false
	}

	win, _ := participantData["win"].(bool)
	return win
}

func (aes *AnalyticsEngineService) extractRole(match map[string]interface{}) string {
	participantDataStr, ok := match["participant_data"].(string)
	if !ok {
		return "UNKNOWN"
	}

	var participantData map[string]interface{}
	if err := json.Unmarshal([]byte(participantDataStr), &participantData); err != nil {
		return "UNKNOWN"
	}

	role, _ := participantData["teamPosition"].(string)
	return role
}

func (aes *AnalyticsEngineService) extractChampionID(match map[string]interface{}) int {
	participantDataStr, ok := match["participant_data"].(string)
	if !ok {
		return 0
	}

	var participantData map[string]interface{}
	if err := json.Unmarshal([]byte(participantDataStr), &participantData); err != nil {
		return 0
	}

	if championID, ok := participantData["championId"].(float64); ok {
		return int(championID)
	}
	return 0
}

func (aes *AnalyticsEngineService) extractChampionName(match map[string]interface{}) string {
	participantDataStr, ok := match["participant_data"].(string)
	if !ok {
		return "Unknown"
	}

	var participantData map[string]interface{}
	if err := json.Unmarshal([]byte(participantDataStr), &participantData); err != nil {
		return "Unknown"
	}

	name, _ := participantData["championName"].(string)
	if name == "" {
		return "Unknown"
	}
	return name
}

func (aes *AnalyticsEngineService) extractSeason(match map[string]interface{}) string {
	if season, ok := match["season"].(string); ok {
		return season
	}
	return "S14" // Default current season
}

func (aes *AnalyticsEngineService) calculateTotalKDA(matches []map[string]interface{}) (float64, float64, float64) {
	var totalKills, totalDeaths, totalAssists float64
	
	for _, match := range matches {
		totalKills += aes.extractStat(match, "kills")
		totalDeaths += aes.extractStat(match, "deaths")
		totalAssists += aes.extractStat(match, "assists")
	}
	
	return totalKills, totalDeaths, totalAssists
}

func (aes *AnalyticsEngineService) calculateMetricsFromMatches(matches []map[string]interface{}) models.PerformanceMetrics {
	if len(matches) == 0 {
		return models.PerformanceMetrics{}
	}

	totalGames := len(matches)
	wins := 0
	var totalKills, totalDeaths, totalAssists float64
	var totalCS, totalGold, totalDamage, totalVision float64
	var totalDuration float64

	for _, match := range matches {
		if aes.extractWinStatus(match) {
			wins++
		}
		
		totalKills += aes.extractStat(match, "kills")
		totalDeaths += aes.extractStat(match, "deaths")
		totalAssists += aes.extractStat(match, "assists")
		
		totalCS += aes.extractStat(match, "totalMinionsKilled") + aes.extractStat(match, "neutralMinionsKilled")
		totalGold += aes.extractStat(match, "goldEarned")
		totalDamage += aes.extractStat(match, "totalDamageDealtToChampions")
		totalVision += aes.extractStat(match, "visionScore")
		
		if duration, ok := match["game_duration"].(int); ok {
			totalDuration += float64(duration)
		}
	}

	losses := totalGames - wins
	winRate := float64(wins) / float64(totalGames) * 100
	
	avgKills := totalKills / float64(totalGames)
	avgDeaths := totalDeaths / float64(totalGames)
	avgAssists := totalAssists / float64(totalGames)
	avgKDA := (avgKills + avgAssists) / math.Max(avgDeaths, 1)
	
	avgCSPerMin := 0.0
	avgGoldPerMin := 0.0
	avgDamagePerMin := 0.0
	
	if totalDuration > 0 {
		avgCSPerMin = totalCS / (totalDuration / 60)
		avgGoldPerMin = totalGold / (totalDuration / 60)
		avgDamagePerMin = totalDamage / (totalDuration / 60)
	}
	
	avgVisionScore := totalVision / float64(totalGames)
	
	performanceScore := aes.calculatePerformanceScore(winRate, avgKDA, avgCSPerMin, avgVisionScore, avgDamagePerMin)

	return models.PerformanceMetrics{
		GamesPlayed:      totalGames,
		Wins:             wins,
		Losses:           losses,
		WinRate:          winRate,
		AvgKills:         avgKills,
		AvgDeaths:        avgDeaths,
		AvgAssists:       avgAssists,
		AvgKDA:           avgKDA,
		AvgCSPerMin:      avgCSPerMin,
		AvgGoldPerMin:    avgGoldPerMin,
		AvgDamagePerMin:  avgDamagePerMin,
		AvgVisionScore:   avgVisionScore,
		PerformanceScore: performanceScore,
		TrendDirection:   models.TrendStable,
	}
}

func (aes *AnalyticsEngineService) calculatePerformanceScore(winRate, kda, csPerMin, visionScore, damagePerMin float64) float64 {
	// Normalize metrics to 0-100 scale
	winRateNorm := math.Min(winRate, 100)
	kdaNorm := math.Min((kda/4.0)*100, 100)      // 4.0 KDA = 100 score
	csNorm := math.Min((csPerMin/10.0)*100, 100) // 10 CS/min = 100 score
	visionNorm := math.Min((visionScore/2.0)*10, 100) // 20 vision score = 100
	damageNorm := math.Min((damagePerMin/1000)*100, 100) // 1000 DPM = 100

	// Weighted average
	score := winRateNorm*0.4 + kdaNorm*0.25 + csNorm*0.15 + visionNorm*0.1 + damageNorm*0.1

	return math.Round(score*100) / 100
}

func (aes *AnalyticsEngineService) analyzeRolePerformance(matches []map[string]interface{}) map[string]models.PerformanceMetrics {
	roleGroups := make(map[string][]map[string]interface{})
	
	for _, match := range matches {
		role := aes.extractRole(match)
		roleGroups[role] = append(roleGroups[role], match)
	}

	result := make(map[string]models.PerformanceMetrics)
	for role, roleMatches := range roleGroups {
		result[role] = aes.calculateMetricsFromMatches(roleMatches)
	}

	return result
}

func (aes *AnalyticsEngineService) getTopChampions(matches []map[string]interface{}, limit int) []models.ChampionPerformance {
	championGroups := make(map[int][]map[string]interface{})
	championNames := make(map[int]string)
	
	for _, match := range matches {
		championID := aes.extractChampionID(match)
		championName := aes.extractChampionName(match)
		
		championGroups[championID] = append(championGroups[championID], match)
		championNames[championID] = championName
	}

	var performances []models.ChampionPerformance
	for championID, champMatches := range championGroups {
		metrics := aes.calculateMetricsFromMatches(champMatches)
		
		performances = append(performances, models.ChampionPerformance{
			ChampionID:       championID,
			ChampionName:     championNames[championID],
			Games:            metrics.GamesPlayed,
			WinRate:          metrics.WinRate,
			PerformanceScore: metrics.PerformanceScore,
			AvgKDA:           metrics.AvgKDA,
		})
	}

	// Sort by performance score
	sort.Slice(performances, func(i, j int) bool {
		return performances[i].PerformanceScore > performances[j].PerformanceScore
	})

	if len(performances) > limit {
		performances = performances[:limit]
	}

	return performances
}

func (aes *AnalyticsEngineService) calculateRecentTrend(matches []map[string]interface{}) string {
	if len(matches) < 6 {
		return "stable"
	}

	// Sort by date
	sort.Slice(matches, func(i, j int) bool {
		dateI, _ := matches[i]["game_creation"].(time.Time)
		dateJ, _ := matches[j]["game_creation"].(time.Time)
		return dateI.Before(dateJ)
	})

	// Split into first and second half
	midPoint := len(matches) / 2
	firstHalf := matches[:midPoint]
	secondHalf := matches[midPoint:]

	firstWR := aes.calculateWinRate(firstHalf)
	secondWR := aes.calculateWinRate(secondHalf)

	diff := secondWR - firstWR

	if diff > 10 {
		return "improving"
	} else if diff < -10 {
		return "declining"
	}
	return "stable"
}

func (aes *AnalyticsEngineService) calculateWinRate(matches []map[string]interface{}) float64 {
	if len(matches) == 0 {
		return 0
	}
	
	wins := 0
	for _, match := range matches {
		if aes.extractWinStatus(match) {
			wins++
		}
	}
	
	return float64(wins) / float64(len(matches)) * 100
}

func (aes *AnalyticsEngineService) generatePeriodSuggestions(matches []map[string]interface{}, period string, rolePerformance map[string]models.PerformanceMetrics) []string {
	var suggestions []string

	if len(matches) == 0 {
		return suggestions
	}

	switch period {
	case "today":
		recentWins := 0
		recentGames := math.Min(float64(len(matches)), 5)
		for i := len(matches) - int(recentGames); i < len(matches); i++ {
			if aes.extractWinStatus(matches[i]) {
				recentWins++
			}
		}
		
		if recentWins >= 3 {
			suggestions = append(suggestions, "Tu es en forme ! Continue sur ta lancée.")
		} else if recentWins <= 1 {
			suggestions = append(suggestions, "Prends une pause, reviens plus tard.")
		}

	case "week":
		if len(rolePerformance) > 0 {
			bestWR := 0.0
			bestRole := ""
			for role, metrics := range rolePerformance {
				if metrics.WinRate > bestWR {
					bestWR = metrics.WinRate
					bestRole = role
				}
			}
			suggestions = append(suggestions, fmt.Sprintf("Focus sur %s, ton meilleur rôle cette semaine (%.1f%% WR)", bestRole, bestWR))
		}

	case "month":
		totalGames := len(matches)
		if totalGames > 50 {
			avgPerformance := 0.0
			for _, match := range matches {
				avgPerformance += aes.calculateGamePerformanceScore(match)
			}
			avgPerformance /= float64(totalGames)
			
			if avgPerformance > 70 {
				suggestions = append(suggestions, "Excellente consistance ! Prêt pour le climb.")
			} else {
				suggestions = append(suggestions, "Focus sur l'amélioration, pas sur la quantité de games.")
			}
		}
	}

	return suggestions
}

func (aes *AnalyticsEngineService) calculateGamePerformanceScore(match map[string]interface{}) float64 {
	kills := aes.extractStat(match, "kills")
	deaths := aes.extractStat(match, "deaths")
	assists := aes.extractStat(match, "assists")
	
	kda := (kills + assists) / math.Max(deaths, 1)
	
	totalCS := aes.extractStat(match, "totalMinionsKilled") + aes.extractStat(match, "neutralMinionsKilled")
	duration, _ := match["game_duration"].(int)
	if duration == 0 {
		duration = 1800 // Default 30 minutes
	}
	csPerMin := totalCS / (float64(duration) / 60)

	winBonus := 0.0
	if aes.extractWinStatus(match) {
		winBonus = 30
	}
	
	kdaScore := math.Min(kda*15, 40)
	csScore := math.Min(csPerMin*2, 30)

	return winBonus + kdaScore + csScore
}

// Additional helper methods for remaining functionality...

func (aes *AnalyticsEngineService) getMatchesForPeriod(userID int, period string) ([]map[string]interface{}, error) {
	// This would implement database query to get matches for a specific period
	// For now, return empty slice
	return []map[string]interface{}{}, nil
}

func (aes *AnalyticsEngineService) calculateRoleTrend(matches []map[string]interface{}) string {
	// Implement role trend calculation
	return "stable"
}

func (aes *AnalyticsEngineService) calculateMasteryScore(matches []map[string]interface{}) float64 {
	if len(matches) == 0 {
		return 0
	}

	// Factors: games played, consistency, recent performance
	gamesPlayed := len(matches)
	experienceScore := math.Min(float64(gamesPlayed)*2, 40) // Max 40 for experience

	// Consistency (standard deviation of performance)
	var performances []float64
	for _, match := range matches {
		performances = append(performances, aes.calculateGamePerformanceScore(match))
	}
	
	avgPerformance := 0.0
	for _, perf := range performances {
		avgPerformance += perf
	}
	avgPerformance /= float64(len(performances))

	variance := 0.0
	for _, perf := range performances {
		variance += math.Pow(perf-avgPerformance, 2)
	}
	variance /= float64(len(performances))
	stdDev := math.Sqrt(variance)

	consistency := 100 - stdDev
	consistencyScore := math.Max(math.Min(consistency/2, 30), 0) // Max 30 for consistency

	// Recent trend (last 5 games vs overall)
	recentCount := 5
	if len(performances) < recentCount {
		recentCount = len(performances)
	}
	
	recentPerformances := performances[len(performances)-recentCount:]
	recentAvg := 0.0
	for _, perf := range recentPerformances {
		recentAvg += perf
	}
	recentAvg /= float64(len(recentPerformances))
	
	trendScore := math.Min(math.Max((recentAvg-avgPerformance)/2+15, 0), 30) // Max 30 for trend

	return math.Min(experienceScore+consistencyScore+trendScore, 100)
}

func (aes *AnalyticsEngineService) findBestGame(matches []map[string]interface{}) models.GameSummary {
	if len(matches) == 0 {
		return models.GameSummary{}
	}

	bestScore := -1.0
	bestMatch := matches[0]
	
	for _, match := range matches {
		score := aes.calculateGamePerformanceScore(match)
		if score > bestScore {
			bestScore = score
			bestMatch = match
		}
	}

	return aes.formatGameSummary(bestMatch)
}

func (aes *AnalyticsEngineService) findWorstGame(matches []map[string]interface{}) models.GameSummary {
	if len(matches) == 0 {
		return models.GameSummary{}
	}

	worstScore := 999.0
	worstMatch := matches[0]
	
	for _, match := range matches {
		score := aes.calculateGamePerformanceScore(match)
		if score < worstScore {
			worstScore = score
			worstMatch = match
		}
	}

	return aes.formatGameSummary(worstMatch)
}

func (aes *AnalyticsEngineService) formatGameSummary(match map[string]interface{}) models.GameSummary {
	matchID, _ := match["match_id"].(string)
	gameCreation, _ := match["game_creation"].(time.Time)
	duration, _ := match["game_duration"].(int)
	
	kills := int(aes.extractStat(match, "kills"))
	deaths := int(aes.extractStat(match, "deaths"))
	assists := int(aes.extractStat(match, "assists"))
	
	return models.GameSummary{
		MatchID:          matchID,
		Date:             gameCreation,
		Champion:         aes.extractChampionName(match),
		Role:             aes.extractRole(match),
		KDA:              fmt.Sprintf("%d/%d/%d", kills, deaths, assists),
		CS:               int(aes.extractStat(match, "totalMinionsKilled") + aes.extractStat(match, "neutralMinionsKilled")),
		Gold:             int(aes.extractStat(match, "goldEarned")),
		Damage:           int(aes.extractStat(match, "totalDamageDealtToChampions")),
		VisionScore:      int(aes.extractStat(match, "visionScore")),
		GameDuration:     duration,
		Win:              aes.extractWinStatus(match),
		PerformanceScore: aes.calculateGamePerformanceScore(match),
	}
}

func (aes *AnalyticsEngineService) generateChampionSuggestions(matches []map[string]interface{}, championID int) []string {
	var suggestions []string

	if len(matches) == 0 {
		return suggestions
	}

	metrics := aes.calculateMetricsFromMatches(matches)

	// KDA suggestions
	if metrics.AvgKDA < 1.5 {
		suggestions = append(suggestions, "Focus sur la survie et l'assistance plutôt que les kills")
	} else if metrics.AvgKDA > 3.0 {
		suggestions = append(suggestions, "Excellente KDA ! Tu maîtrises bien ce champion")
	}

	// CS suggestions
	if metrics.AvgCSPerMin < 6 {
		suggestions = append(suggestions, "Améliore ton farming pour plus d'impact économique")
	}

	// Win rate suggestions
	if metrics.WinRate < 45 {
		suggestions = append(suggestions, "Considère étudier des guides ou VODs pour ce champion")
	} else if metrics.WinRate > 65 {
		suggestions = append(suggestions, "Champion forte ! Continue à le jouer pour le climb")
	}

	return suggestions
}

func (aes *AnalyticsEngineService) analyzeSkillProgression(matches []map[string]interface{}) models.SkillProgression {
	if len(matches) < 5 {
		return models.SkillProgression{InsufficientData: true}
	}

	// Sort by date
	sort.Slice(matches, func(i, j int) bool {
		dateI, _ := matches[i]["game_creation"].(time.Time)
		dateJ, _ := matches[j]["game_creation"].(time.Time)
		return dateI.Before(dateJ)
	})

	// Calculate progression in chunks
	chunkSize := len(matches) / 4
	if chunkSize < 2 {
		chunkSize = 2
	}

	var progression []models.ProgressionPeriod
	for i := 0; i < len(matches); i += chunkSize {
		end := i + chunkSize
		if end > len(matches) {
			end = len(matches)
		}
		
		chunk := matches[i:end]
		if len(chunk) >= 2 {
			metrics := aes.calculateMetricsFromMatches(chunk)
			progression = append(progression, models.ProgressionPeriod{
				Period:           len(progression) + 1,
				Games:            len(chunk),
				WinRate:          metrics.WinRate,
				AvgKDA:           metrics.AvgKDA,
				PerformanceScore: metrics.PerformanceScore,
			})
		}
	}

	// Calculate overall trend
	trend := "stable"
	improvement := 0.0
	if len(progression) >= 2 {
		firstScore := progression[0].PerformanceScore
		lastScore := progression[len(progression)-1].PerformanceScore
		improvement = lastScore - firstScore
		
		if improvement > 5 {
			trend = "improving"
		} else if improvement < -5 {
			trend = "declining"
		}
	}

	return models.SkillProgression{
		ProgressionData:    progression,
		OverallImprovement: improvement,
		Trend:              trend,
	}
}

// Implement remaining methods...
func (aes *AnalyticsEngineService) analyzePerformanceTrends(recentMatches, monthMatches []map[string]interface{}) []models.ImprovementSuggestion {
	var suggestions []models.ImprovementSuggestion
	
	if len(recentMatches) < 3 || len(monthMatches) < 10 {
		return suggestions
	}

	recentWR := aes.calculateWinRate(recentMatches)
	monthlyWR := aes.calculateWinRate(monthMatches)
	wrDiff := recentWR - monthlyWR

	if wrDiff > 10 {
		suggestions = append(suggestions, models.ImprovementSuggestion{
			Type:                "performance",
			Title:               "Forme excellente !",
			Description:         fmt.Sprintf("Winrate récent (%.1f%%) bien supérieur au mensuel (%.1f%%)", recentWR, monthlyWR),
			Priority:            1,
			ExpectedImprovement: "+5% winrate stable",
			ActionItems:         []string{"Continue sur cette lancée", "Analyse ce qui fonctionne bien"},
			TimePeriod:          "week",
		})
	} else if wrDiff < -10 {
		suggestions = append(suggestions, models.ImprovementSuggestion{
			Type:                "performance",
			Title:               "Baisse de forme détectée",
			Description:         fmt.Sprintf("Winrate récent (%.1f%%) en baisse par rapport au mensuel (%.1f%%)", recentWR, monthlyWR),
			Priority:            1,
			ExpectedImprovement: "+8% winrate recovery",
			ActionItems:         []string{"Prends une pause", "Revois tes replays récents", "Retour aux basics"},
			TimePeriod:          "week",
		})
	}

	return suggestions
}

func (aes *AnalyticsEngineService) analyzeRoleWeaknesses(matches []map[string]interface{}) []models.ImprovementSuggestion {
	var suggestions []models.ImprovementSuggestion
	
	rolePerformance := aes.analyzeRolePerformance(matches)
	if len(rolePerformance) == 0 {
		return suggestions
	}

	// Find worst performing role
	worstWR := 100.0
	worstRole := ""
	for role, metrics := range rolePerformance {
		if metrics.WinRate < worstWR && metrics.GamesPlayed >= 5 {
			worstWR = metrics.WinRate
			worstRole = role
		}
	}

	if worstWR < 40 && worstRole != "" {
		suggestions = append(suggestions, models.ImprovementSuggestion{
			Type:                "role",
			Title:               fmt.Sprintf("Amélioration nécessaire en %s", worstRole),
			Description:         fmt.Sprintf("Winrate de %.1f%% en %s nécessite du travail", worstWR, worstRole),
			Priority:            2,
			ExpectedImprovement: "+15% winrate en " + worstRole,
			ActionItems:         []string{"Étude des guides pour " + worstRole, "Practice en normal", "Analyse des erreurs communes"},
			Role:                &worstRole,
			TimePeriod:          "month",
		})
	}

	return suggestions
}

func (aes *AnalyticsEngineService) analyzeChampionPool(matches []map[string]interface{}) []models.ImprovementSuggestion {
	var suggestions []models.ImprovementSuggestion
	
	// Count unique champions
	champions := make(map[int]bool)
	for _, match := range matches {
		championID := aes.extractChampionID(match)
		champions[championID] = true
	}

	poolSize := len(champions)

	if poolSize < 3 && len(matches) > 20 {
		suggestions = append(suggestions, models.ImprovementSuggestion{
			Type:                "champion_pool",
			Title:               "Pool de champions trop restreint",
			Description:         fmt.Sprintf("Seulement %d champions joués. Élargis ton pool pour plus de flexibilité", poolSize),
			Priority:            2,
			ExpectedImprovement: "+5% adaptabilité draft",
			ActionItems:         []string{"Apprends 2-3 nouveaux champions", "Focus sur la méta actuelle", "Diversifie par rôle"},
			TimePeriod:          "month",
		})
	} else if poolSize > 15 && len(matches) > 30 {
		suggestions = append(suggestions, models.ImprovementSuggestion{
			Type:                "champion_pool",
			Title:               "Pool de champions trop large",
			Description:         fmt.Sprintf("%d champions différents. Focus sur 3-5 champions pour mieux les maîtriser", poolSize),
			Priority:            2,
			ExpectedImprovement: "+10% maîtrise champions",
			ActionItems:         []string{"Sélectionne 3-5 champions principaux", "Focus sur la maîtrise approfondie", "Abandonne les picks situationnels"},
			TimePeriod:          "month",
		})
	}

	return suggestions
}

func (aes *AnalyticsEngineService) generateMetaSuggestions(matches []map[string]interface{}) []models.ImprovementSuggestion {
	var suggestions []models.ImprovementSuggestion
	
	if len(matches) >= 5 {
		suggestions = append(suggestions, models.ImprovementSuggestion{
			Type:                "meta",
			Title:               "Adapte-toi au meta actuel",
			Description:         "Considère jouer des champions tier S pour maximiser tes chances",
			Priority:            3,
			ExpectedImprovement: "+3% winrate méta",
			ActionItems:         []string{"Consulte les tier lists actuelles", "Regarde les picks pros", "Teste en normal d'abord"},
			TimePeriod:          "week",
		})
	}

	return suggestions
}

func (aes *AnalyticsEngineService) calculateTrendMetrics(matches []map[string]interface{}) models.TrendMetrics {
	if len(matches) == 0 {
		return models.TrendMetrics{Trend: "stable", Games: 0, WinRate: 0}
	}

	wins := 0
	for _, match := range matches {
		if aes.extractWinStatus(match) {
			wins++
		}
	}
	
	winRate := float64(wins) / float64(len(matches)) * 100
	losses := len(matches) - wins

	// Calculate trend based on chronological order
	trend := "stable"
	if len(matches) >= 6 {
		sort.Slice(matches, func(i, j int) bool {
			dateI, _ := matches[i]["game_creation"].(time.Time)
			dateJ, _ := matches[j]["game_creation"].(time.Time)
			return dateI.Before(dateJ)
		})

		midPoint := len(matches) / 2
		firstHalf := matches[:midPoint]
		secondHalf := matches[midPoint:]

		firstWR := aes.calculateWinRate(firstHalf)
		secondWR := aes.calculateWinRate(secondHalf)

		if secondWR-firstWR > 10 {
			trend = "improving"
		} else if firstWR-secondWR > 10 {
			trend = "declining"
		}
	}

	return models.TrendMetrics{
		Trend:   trend,
		Games:   len(matches),
		WinRate: winRate,
		Wins:    wins,
		Losses:  losses,
	}
}

func (aes *AnalyticsEngineService) calculateImprovementVelocity(weekMatches, monthMatches []map[string]interface{}) float64 {
	if len(weekMatches) == 0 || len(monthMatches) == 0 {
		return 0.0
	}

	weekWR := aes.calculateWinRate(weekMatches)
	monthWR := aes.calculateWinRate(monthMatches)

	// Velocity = change per week
	return weekWR - monthWR
}

func (aes *AnalyticsEngineService) calculateConsistencyScore(matches []map[string]interface{}) float64 {
	if len(matches) < 5 {
		return 50.0
	}

	// Calculate game-by-game performance scores
	var performances []float64
	for _, match := range matches {
		score := aes.calculateGamePerformanceScore(match)
		performances = append(performances, score)
	}

	// Lower variance = higher consistency
	meanPerf := 0.0
	for _, p := range performances {
		meanPerf += p
	}
	meanPerf /= float64(len(performances))

	variance := 0.0
	for _, p := range performances {
		variance += math.Pow(p-meanPerf, 2)
	}
	variance /= float64(len(performances))
	stdDev := math.Sqrt(variance)

	// Convert to 0-100 scale (lower std dev = higher consistency)
	consistency := math.Max(0, 100-(stdDev/2))
	return math.Min(consistency, 100)
}

func (aes *AnalyticsEngineService) findPeakPerformancePeriod(matches []map[string]interface{}) models.PeakPerformancePeriod {
	if len(matches) < 10 {
		return models.PeakPerformancePeriod{Period: "insufficient_data", Performance: 0}
	}

	// Sort matches chronologically
	sort.Slice(matches, func(i, j int) bool {
		dateI, _ := matches[i]["game_creation"].(time.Time)
		dateJ, _ := matches[j]["game_creation"].(time.Time)
		return dateI.Before(dateJ)
	})

	// Sliding window to find best 10-game period
	windowSize := 10
	if len(matches) < windowSize {
		windowSize = len(matches)
	}

	bestWinRate := 0.0
	bestStart := 0
	bestEnd := windowSize

	for i := 0; i <= len(matches)-windowSize; i++ {
		window := matches[i : i+windowSize]
		winRate := aes.calculateWinRate(window)

		if winRate > bestWinRate {
			bestWinRate = winRate
			bestStart = i
			bestEnd = i + windowSize
		}
	}

	return models.PeakPerformancePeriod{
		Period:      fmt.Sprintf("Games %d-%d", bestStart+1, bestEnd),
		Performance: bestWinRate,
		Games:       windowSize,
	}
}

func (aes *AnalyticsEngineService) getLastPlayedDate(matches []map[string]interface{}) time.Time {
	if len(matches) == 0 {
		return time.Now()
	}

	latest := time.Time{}
	for _, match := range matches {
		if date, ok := match["game_creation"].(time.Time); ok {
			if date.After(latest) {
				latest = date
			}
		}
	}

	if latest.IsZero() {
		return time.Now()
	}
	return latest
}

func (aes *AnalyticsEngineService) saveChampionStats(userID, championID int, championName, role, season, period string, metrics models.PerformanceMetrics, lastPlayed time.Time) error {
	query := `
		INSERT OR REPLACE INTO champion_stats 
		(user_id, champion_id, champion_name, role, season, time_period,
		 games_played, wins, losses, win_rate, avg_kills, avg_deaths, avg_assists,
		 avg_kda, avg_cs_per_min, avg_gold_per_min, avg_damage_per_min,
		 avg_vision_score, performance_score, trend_direction, last_played, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := aes.db.Exec(query,
		userID, championID, championName, role, season, period,
		metrics.GamesPlayed, metrics.Wins, metrics.Losses, metrics.WinRate,
		metrics.AvgKills, metrics.AvgDeaths, metrics.AvgAssists, metrics.AvgKDA,
		metrics.AvgCSPerMin, metrics.AvgGoldPerMin, metrics.AvgDamagePerMin,
		metrics.AvgVisionScore, metrics.PerformanceScore, string(metrics.TrendDirection),
		lastPlayed, time.Now())

	return err
}