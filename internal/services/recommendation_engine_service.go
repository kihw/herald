package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"

	"lol-match-exporter/internal/db"
	"lol-match-exporter/internal/models"
)

// RecommendationEngineService generates intelligent recommendations based on player performance and meta
type RecommendationEngineService struct {
	db       *sql.DB
	settings models.RecommendationSettings
}

// NewRecommendationEngineService creates a new recommendation engine service
func NewRecommendationEngineService(database *db.Database) *RecommendationEngineService {
	var sqlDB *sql.DB
	if database != nil {
		sqlDB = database.DB
	}
	return &RecommendationEngineService{
		db: sqlDB,
		settings: models.RecommendationSettings{
			MetaWeight:           0.7,
			PersonalWeight:       0.8,
			RiskTolerance:        0.5,
			LearningRate:         0.3,
			UpdateFrequency:      6,
			MinConfidence:        0.4,
			MaxRecommendations:   15,
			PersonalizationLevel: "high",
		},
	}
}

// GenerateComprehensiveRecommendations generates comprehensive recommendations for a user
func (res *RecommendationEngineService) GenerateComprehensiveRecommendations(userID int) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation

	// Get user data
	matches, err := res.getMatchesForPeriod(userID, "month")
	if err != nil {
		return nil, fmt.Errorf("failed to get matches: %w", err)
	}

	if len(matches) == 0 {
		return recommendations, nil
	}

	// Generate different types of recommendations
	championRecs, err := res.SuggestChampionsForRoles(userID)
	if err == nil {
		recommendations = append(recommendations, championRecs...)
	}

	performanceRecs, err := res.AnalyzeChampionPerformanceGaps(userID)
	if err == nil {
		recommendations = append(recommendations, performanceRecs...)
	}

	gameplayRecs, err := res.GenerateGameplayTips(userID)
	if err == nil {
		recommendations = append(recommendations, gameplayRecs...)
	}

	banRecs, err := res.RecommendBanPriorities(userID)
	if err == nil {
		recommendations = append(recommendations, banRecs...)
	}

	metaRecs, err := res.SuggestMetaAdaptations(userID)
	if err == nil {
		recommendations = append(recommendations, metaRecs...)
	}

	trainingRecs, err := res.RecommendTrainingFocus(userID)
	if err == nil {
		recommendations = append(recommendations, trainingRecs...)
	}

	// Sort by priority and confidence
	sort.Slice(recommendations, func(i, j int) bool {
		if recommendations[i].Priority == recommendations[j].Priority {
			return recommendations[i].Confidence > recommendations[j].Confidence
		}
		return recommendations[i].Priority < recommendations[j].Priority
	})

	// Apply settings limits
	if len(recommendations) > res.settings.MaxRecommendations {
		recommendations = recommendations[:res.settings.MaxRecommendations]
	}

	// Filter by minimum confidence
	filteredRecs := make([]models.Recommendation, 0)
	for _, rec := range recommendations {
		if rec.Confidence >= res.settings.MinConfidence {
			filteredRecs = append(filteredRecs, rec)
		}
	}

	// Save to database
	if err := res.saveRecommendations(userID, filteredRecs); err != nil {
		log.Printf("Failed to save recommendations: %v", err)
	}

	return filteredRecs, nil
}

// SuggestChampionsForRoles suggests optimal champions for each role based on current meta and player skill
func (res *RecommendationEngineService) SuggestChampionsForRoles(userID int) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation

	// Get role performance
	monthMatches, err := res.getMatchesForPeriod(userID, "month")
	if err != nil {
		return nil, err
	}

	rolePerformance := res.analyzeRolePerformance(monthMatches)

	for role, matches := range res.groupMatchesByRole(monthMatches) {
		if len(matches) < 3 { // Need minimum games to analyze
			continue
		}

		roleMetrics := rolePerformance[role]

		// Find champions user hasn't played much but are strong in meta
		playedChampions := res.getPlayedChampions(matches)
		roleSuggestions := res.getMetaChampionsForRole(role, playedChampions)

		for i, suggestion := range roleSuggestions {
			if i >= 3 { // Top 3 suggestions per role
				break
			}

			championName := res.getChampionName(suggestion.ChampionID)

			// Calculate expected improvement
			currentWR := roleMetrics.WinRate
			expectedWR := math.Min(currentWR+(suggestion.MetaStrength*15), 85)
			improvement := fmt.Sprintf("+%.1f%% winrate", expectedWR-currentWR)

			priority := 1
			if suggestion.MetaStrength < 0.9 {
				priority = 2
			}

			recommendations = append(recommendations, models.Recommendation{
				Type:                models.ChampionSuggestion,
				Title:               fmt.Sprintf("Essaie %s en %s", championName, role),
				Description:         fmt.Sprintf("%s est très fort dans la méta actuelle (%.0f%% force) et pourrait améliorer tes performances en %s.", championName, suggestion.MetaStrength*100, role),
				Priority:            priority,
				Confidence:          suggestion.MetaStrength * 0.8,
				ExpectedImprovement: improvement,
				ActionItems: []string{
					fmt.Sprintf("Regarde des guides sur %s", championName),
					fmt.Sprintf("Pratique %s en normal d'abord", championName),
					fmt.Sprintf("Focus sur le build optimal pour %s", championName),
				},
				ChampionID: &suggestion.ChampionID,
				Role:       &role,
				TimePeriod: "week",
			})
		}
	}

	return recommendations, nil
}

// AnalyzeChampionPerformanceGaps identifies champions with potential for improvement
func (res *RecommendationEngineService) AnalyzeChampionPerformanceGaps(userID int) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation

	// Get champion stats from database
	championStats, err := res.getChampionStats(userID)
	if err != nil {
		return nil, err
	}

	for _, stats := range championStats {
		if stats.GamesPlayed < 5 { // Need minimum games
			continue
		}

		// Identify specific improvement areas
		var issues []string
		var improvements []string

		if stats.WinRate < 45 {
			issues = append(issues, "winrate faible")
			improvements = append(improvements, "Analyse tes replays pour identifier les erreurs")
		}

		if stats.AvgKDA < 1.5 {
			issues = append(issues, "KDA faible")
			improvements = append(improvements, "Focus sur la survie et le positionnement")
		}

		if stats.AvgCSPerMin < 6 && (stats.Role == "MIDDLE" || stats.Role == "BOTTOM" || stats.Role == "TOP") {
			issues = append(issues, "farm insuffisant")
			improvements = append(improvements, "Pratique le last-hitting et la gestion des vagues")
		}

		if len(issues) > 0 {
			priority := 2
			if len(issues) >= 2 {
				priority = 1
			}

			confidence := math.Min(0.8, float64(stats.GamesPlayed)/20.0) // More games = higher confidence

			recommendations = append(recommendations, models.Recommendation{
				Type:                models.TrainingFocus,
				Title:               fmt.Sprintf("Améliore tes performances sur %s", stats.ChampionName),
				Description:         fmt.Sprintf("Tu as %d games sur %s mais %s. Potentiel d'amélioration important.", stats.GamesPlayed, stats.ChampionName, strings.Join(issues, ", ")),
				Priority:            priority,
				Confidence:          confidence,
				ExpectedImprovement: fmt.Sprintf("+%d%% performance", 10+len(issues)*5),
				ActionItems:         improvements[:int(math.Min(float64(len(improvements)), 3))],
				ChampionID:          &stats.ChampionID,
				Role:                &stats.Role,
				TimePeriod:          "month",
			})
		}
	}

	// Sort by potential impact and return top 5
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Confidence > recommendations[j].Confidence
	})

	if len(recommendations) > 5 {
		recommendations = recommendations[:5]
	}

	return recommendations, nil
}

// GenerateGameplayTips generates specific gameplay improvement tips
func (res *RecommendationEngineService) GenerateGameplayTips(userID int) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation

	// Analyze recent performance trends
	weekMatches, err := res.getMatchesForPeriod(userID, "week")
	if err != nil {
		return nil, err
	}

	if len(weekMatches) < 5 {
		return recommendations, nil
	}

	// Early game analysis
	earlyGameIssues := res.analyzeEarlyGamePerformance(weekMatches)
	if earlyGameIssues != nil {
		recommendations = append(recommendations, models.Recommendation{
			Type:                models.GameplayTip,
			Title:               "Améliore ton early game",
			Description:         earlyGameIssues.Description,
			Priority:            1,
			Confidence:          0.7,
			ExpectedImprovement: "+8% winrate",
			ActionItems:         earlyGameIssues.Tips,
			TimePeriod:          "week",
		})
	}

	// Team fighting analysis
	teamfightIssues := res.analyzeTeamfightPerformance(weekMatches)
	if teamfightIssues != nil {
		recommendations = append(recommendations, models.Recommendation{
			Type:                models.GameplayTip,
			Title:               "Optimise tes teamfights",
			Description:         teamfightIssues.Description,
			Priority:            2,
			Confidence:          0.6,
			ExpectedImprovement: "+6% winrate",
			ActionItems:         teamfightIssues.Tips,
			TimePeriod:          "week",
		})
	}

	// Vision control analysis
	visionIssues := res.analyzeVisionControl(weekMatches)
	if visionIssues != nil {
		recommendations = append(recommendations, models.Recommendation{
			Type:                models.GameplayTip,
			Title:               "Améliore ton contrôle de vision",
			Description:         visionIssues.Description,
			Priority:            2,
			Confidence:          0.6,
			ExpectedImprovement: "+5% winrate",
			ActionItems:         visionIssues.Tips,
			TimePeriod:          "week",
		})
	}

	return recommendations, nil
}

// RecommendBanPriorities recommends champions to ban based on current meta and personal weaknesses
func (res *RecommendationEngineService) RecommendBanPriorities(userID int) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation

	// Get recent matches to analyze enemy champions that caused problems
	recentMatches, err := res.getMatchesForPeriod(userID, "week")
	if err != nil {
		return nil, err
	}

	// Analyze losses to identify problematic enemy champions
	problemChampions := res.identifyProblemChampions(recentMatches)

	// Combine with current meta threats
	metaThreats := res.getCurrentMetaThreats()

	var banPriorities []models.BanPriority

	// Personal problem champions (higher priority)
	for championID, data := range problemChampions {
		if data.GamesAgainst >= 3 && data.WinRateAgainst < 35 {
			championName := res.getChampionName(championID)
			banPriorities = append(banPriorities, models.BanPriority{
				ChampionID:   championID,
				ChampionName: championName,
				Priority:     1,
				Reason:       fmt.Sprintf("Tu perds %.0f%% contre ce champion", 100-data.WinRateAgainst),
				Confidence:   math.Min(float64(data.GamesAgainst)/10.0, 0.9),
				ThreatLevel:  "personal",
			})
		}
	}

	// Meta threats
	for championID, metaStrength := range metaThreats {
		if metaStrength > 0.92 {
			championName := res.getChampionName(championID)
			banPriorities = append(banPriorities, models.BanPriority{
				ChampionID:   championID,
				ChampionName: championName,
				Priority:     2,
				Reason:       fmt.Sprintf("Pick/ban très élevé dans la méta (%.0f%%)", metaStrength*100),
				Confidence:   0.7,
				ThreatLevel:  "meta",
			})
		}
	}

	// Sort by priority and confidence
	sort.Slice(banPriorities, func(i, j int) bool {
		if banPriorities[i].Priority == banPriorities[j].Priority {
			return banPriorities[i].Confidence > banPriorities[j].Confidence
		}
		return banPriorities[i].Priority < banPriorities[j].Priority
	})

	// Create recommendations for top 3 ban priorities
	for i, banData := range banPriorities {
		if i >= 3 {
			break
		}

		recommendations = append(recommendations, models.Recommendation{
			Type:                models.BanSuggestion,
			Title:               fmt.Sprintf("Ban %s", banData.ChampionName),
			Description:         banData.Reason,
			Priority:            banData.Priority,
			Confidence:          banData.Confidence,
			ExpectedImprovement: "+3% winrate par évitement",
			ActionItems: []string{
				fmt.Sprintf("Priorise le ban de %s", banData.ChampionName),
				fmt.Sprintf("Étudie les contres à %s si non banni", banData.ChampionName),
			},
			ChampionID: &banData.ChampionID,
			TimePeriod: "week",
		})
	}

	return recommendations, nil
}

// SuggestMetaAdaptations suggests adaptations to current meta trends
func (res *RecommendationEngineService) SuggestMetaAdaptations(userID int) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation

	// Analyze current champion pool vs meta
	recentMatches, err := res.getMatchesForPeriod(userID, "month")
	if err != nil {
		return nil, err
	}

	playedChampions := make(map[int]struct {
		Games int
		Roles map[string]bool
	})

	for _, match := range recentMatches {
		champID := res.extractChampionID(match)
		role := res.extractRole(match)

		if data, exists := playedChampions[champID]; exists {
			data.Games++
			data.Roles[role] = true
			playedChampions[champID] = data
		} else {
			playedChampions[champID] = struct {
				Games int
				Roles map[string]bool
			}{
				Games: 1,
				Roles: map[string]bool{role: true},
			}
		}
	}

	// Check if user is playing off-meta champions too much
	var offMetaChampions []models.OffMetaChampion
	for champID, data := range playedChampions {
		metaStrength, exists := models.ChampionMetaStrength[champID]
		if !exists {
			metaStrength = 0.5 // Default strength
		}

		if data.Games >= 5 && metaStrength < 0.7 {
			championName := res.getChampionName(champID)
			var roles []string
			for role := range data.Roles {
				roles = append(roles, role)
			}

			offMetaChampions = append(offMetaChampions, models.OffMetaChampion{
				ChampionID:   champID,
				ChampionName: championName,
				GamesPlayed:  data.Games,
				MetaStrength: metaStrength,
				Recommendation: "reduce",
			})
		}
	}

	if len(offMetaChampions) > 0 {
		// Sort by games played (biggest impact)
		sort.Slice(offMetaChampions, func(i, j int) bool {
			return offMetaChampions[i].GamesPlayed > offMetaChampions[j].GamesPlayed
		})

		topOffMeta := offMetaChampions[0]

		recommendations = append(recommendations, models.Recommendation{
			Type:                models.MetaAdaptationType,
			Title:               fmt.Sprintf("Réduis tes games sur %s", topOffMeta.ChampionName),
			Description:         fmt.Sprintf("%s est actuellement faible dans la méta (%.0f%% force). Considère des picks plus forts.", topOffMeta.ChampionName, topOffMeta.MetaStrength*100),
			Priority:            2,
			Confidence:          0.8,
			ExpectedImprovement: "+4% winrate par adaptation méta",
			ActionItems: []string{
				fmt.Sprintf("Limite %s aux matchups favorables uniquement", topOffMeta.ChampionName),
				"Apprends des champions plus forts dans la méta",
				"Suis les patch notes pour les buffs/nerfs",
			},
			ChampionID: &topOffMeta.ChampionID,
			TimePeriod: "month",
		})
	}

	return recommendations, nil
}

// RecommendTrainingFocus recommends specific areas to focus training on
func (res *RecommendationEngineService) RecommendTrainingFocus(userID int) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation

	// Analyze performance across different metrics
	recentMatches, err := res.getMatchesForPeriod(userID, "month")
	if err != nil {
		return nil, err
	}

	if len(recentMatches) == 0 {
		return recommendations, nil
	}

	// Calculate percentile performance across key metrics
	performanceAnalysis := res.analyzeSkillPercentiles(recentMatches)

	// Identify weakest areas
	var weakAreas []string
	for area, percentile := range performanceAnalysis {
		if percentile < 40 {
			weakAreas = append(weakAreas, area)
		}
	}

	// Sort by weakest first
	sort.Slice(weakAreas, func(i, j int) bool {
		return performanceAnalysis[weakAreas[i]] < performanceAnalysis[weakAreas[j]]
	})

	trainingRecommendations := map[string]struct {
		Title       string
		Description string
		Tips        []string
	}{
		"cs_per_min": {
			Title:       "Améliore ton farming",
			Description: "Ton CS/min est en dessous de la moyenne. Le farming est crucial pour l'avantage économique.",
			Tips: []string{
				"Pratique le last-hitting en partie personnalisée",
				"Apprends la gestion des vagues minions",
				"Optimise tes recalls pour minimiser les CS perdus",
			},
		},
		"vision_score": {
			Title:       "Améliore ton contrôle de vision",
			Description: "Ton vision score est faible. La vision donne de l'information cruciale.",
			Tips: []string{
				"Place plus de wards dans les zones clés",
				"Achète plus de wards de contrôle",
				"Clear plus les wards ennemies",
			},
		},
		"kda": {
			Title:       "Améliore ton positionnement",
			Description: "Ta KDA indique des problèmes de positionnement ou de prise de décision.",
			Tips: []string{
				"Focus sur la survie plutôt que les kills",
				"Améliore ton positionnement en teamfight",
				"Prends moins de risques inutiles",
			},
		},
		"damage_per_min": {
			Title:       "Augmente ton impact damage",
			Description: "Tes dégâts par minute sont faibles. Tu peux être plus agressif.",
			Tips: []string{
				"Cherche plus d'opportunités de trade",
				"Améliore ton spacing en lane",
				"Utilise mieux tes fenêtres de puissance",
			},
		},
	}

	// Create recommendations for weakest areas
	for i, area := range weakAreas {
		if i >= 3 { // Top 3 weakest areas
			break
		}

		if recData, exists := trainingRecommendations[area]; exists {
			priority := 1
			if i > 0 {
				priority = 2 // Highest priority for weakest area
			}

			recommendations = append(recommendations, models.Recommendation{
				Type:                models.TrainingFocus,
				Title:               recData.Title,
				Description:         recData.Description,
				Priority:            priority,
				Confidence:          0.9, // High confidence in skill analysis
				ExpectedImprovement: fmt.Sprintf("+%d%% performance dans ce domaine", 15-i*3),
				ActionItems:         recData.Tips,
				TimePeriod:          "month",
			})
		}
	}

	return recommendations, nil
}

// Helper methods

func (res *RecommendationEngineService) getMatchesForPeriod(userID int, period string) ([]map[string]interface{}, error) {
	// This would implement database query to get matches for a specific period
	// For now, return empty slice
	return []map[string]interface{}{}, nil
}

func (res *RecommendationEngineService) extractChampionID(match map[string]interface{}) int {
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

func (res *RecommendationEngineService) extractRole(match map[string]interface{}) string {
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

func (res *RecommendationEngineService) getChampionName(championID int) string {
	// Champion ID to name mapping (simplified)
	championNames := map[int]string{
		266: "Aatrox", 103: "Ahri", 84: "Akali", 12: "Alistar", 32: "Amumu",
		34: "Anivia", 1: "Annie", 22: "Ashe", 136: "Aurelion Sol", 268: "Azir",
		432: "Bard", 53: "Blitzcrank", 63: "Brand", 201: "Braum", 51: "Caitlyn",
		164: "Camille", 69: "Cassiopeia", 31: "Cho'Gath", 42: "Corki", 122: "Darius",
		131: "Diana", 119: "Draven", 36: "Dr. Mundo", 245: "Ekko", 60: "Elise",
		28: "Evelynn", 81: "Ezreal", 9: "Fiddlesticks", 114: "Fiora", 105: "Fizz",
		3: "Galio", 41: "Gangplank", 86: "Garen", 150: "Gnar", 79: "Gragas",
		104: "Graves", 120: "Hecarim", 74: "Heimerdinger", 420: "Illaoi", 39: "Irelia",
		// ... (add more champions as needed)
	}

	if name, exists := championNames[championID]; exists {
		return name
	}
	return fmt.Sprintf("Champion %d", championID)
}

func (res *RecommendationEngineService) analyzeRolePerformance(matches []map[string]interface{}) map[string]models.PerformanceMetrics {
	roleGroups := make(map[string][]map[string]interface{})

	for _, match := range matches {
		role := res.extractRole(match)
		roleGroups[role] = append(roleGroups[role], match)
	}

	result := make(map[string]models.PerformanceMetrics)
	for role, roleMatches := range roleGroups {
		result[role] = res.calculateMetricsFromMatches(roleMatches)
	}

	return result
}

func (res *RecommendationEngineService) groupMatchesByRole(matches []map[string]interface{}) map[string][]map[string]interface{} {
	roleGroups := make(map[string][]map[string]interface{})

	for _, match := range matches {
		role := res.extractRole(match)
		roleGroups[role] = append(roleGroups[role], match)
	}

	return roleGroups
}

func (res *RecommendationEngineService) getPlayedChampions(matches []map[string]interface{}) map[int]bool {
	played := make(map[int]bool)
	for _, match := range matches {
		champID := res.extractChampionID(match)
		played[champID] = true
	}
	return played
}

func (res *RecommendationEngineService) getMetaChampionsForRole(role string, excludedChampions map[int]bool) []models.MetaChampion {
	// Role to champion mapping (simplified)
	roleChampions := map[string][]int{
		"TOP":     {266, 122, 86, 150, 79, 114, 420, 39, 240, 54, 57, 75, 516, 58, 92, 14, 27, 83, 106, 19},
		"JUNGLE":  {32, 245, 60, 120, 104, 427, 59, 141, 121, 203, 64, 76, 56, 20, 2, 421, 107, 113, 35, 102, 72, 77, 154},
		"MIDDLE":  {103, 84, 1, 136, 268, 69, 42, 131, 28, 105, 3, 74, 38, 55, 10, 7, 127, 99, 90, 82, 61, 80, 246, 13, 517, 134, 163, 91, 4, 112, 8, 142, 238, 115, 26},
		"BOTTOM":  {22, 51, 119, 81, 202, 222, 145, 429, 96, 236, 21, 15, 18, 110, 67, 29, 498},
		"UTILITY": {12, 432, 53, 63, 201, 40, 43, 89, 117, 25, 267, 111, 516, 78, 555, 497, 44, 16, 50, 223, 412, 37, 143, 350},
	}

	var candidates []models.MetaChampion
	if championIDs, exists := roleChampions[role]; exists {
		for _, champID := range championIDs {
			if !excludedChampions[champID] {
				metaStrength, exists := models.ChampionMetaStrength[champID]
				if !exists {
					metaStrength = 0.5
				}

				candidates = append(candidates, models.MetaChampion{
					ChampionID:   champID,
					ChampionName: res.getChampionName(champID),
					MetaStrength: metaStrength,
					Role:         role,
				})
			}
		}
	}

	// Sort by meta strength
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].MetaStrength > candidates[j].MetaStrength
	})

	if len(candidates) > 10 {
		candidates = candidates[:10]
	}

	return candidates
}

func (res *RecommendationEngineService) calculateMetricsFromMatches(matches []map[string]interface{}) models.PerformanceMetrics {
	// This would use the same logic as AnalyticsEngineService.calculateMetricsFromMatches
	// For now, return empty metrics
	return models.PerformanceMetrics{}
}

func (res *RecommendationEngineService) getChampionStats(userID int) ([]models.ChampionStats, error) {
	query := `
		SELECT champion_id, champion_name, role, games_played, win_rate, avg_kda, avg_cs_per_min
		FROM champion_stats 
		WHERE user_id = ? AND games_played >= 5 AND time_period = 'season'
		ORDER BY games_played DESC`

	rows, err := res.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.ChampionStats
	for rows.Next() {
		var stat models.ChampionStats
		err := rows.Scan(&stat.ChampionID, &stat.ChampionName, &stat.Role, 
			&stat.GamesPlayed, &stat.WinRate, &stat.AvgKDA, &stat.AvgCSPerMin)
		if err != nil {
			continue
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

func (res *RecommendationEngineService) analyzeEarlyGamePerformance(matches []map[string]interface{}) *models.GameplayIssue {
	if len(matches) == 0 {
		return nil
	}

	// Simplified early game analysis
	var earlyDeaths []int
	var earlyCS []float64

	for _, match := range matches {
		participantDataStr, ok := match["participant_data"].(string)
		if !ok {
			continue
		}

		var participantData map[string]interface{}
		if err := json.Unmarshal([]byte(participantDataStr), &participantData); err != nil {
			continue
		}

		// Early deaths (approximation)
		if deaths, ok := participantData["deaths"].(float64); ok && deaths > 3 {
			earlyDeaths = append(earlyDeaths, int(deaths))
		}

		// CS efficiency
		totalCS, _ := participantData["totalMinionsKilled"].(float64)
		neutralCS, _ := participantData["neutralMinionsKilled"].(float64)
		duration, _ := match["game_duration"].(int)
		if duration == 0 {
			duration = 1800
		}
		csPerMin := (totalCS + neutralCS) / (float64(duration) / 60.0)
		earlyCS = append(earlyCS, csPerMin)
	}

	if len(earlyCS) == 0 {
		return nil
	}

	avgCSPerMin := 0.0
	for _, cs := range earlyCS {
		avgCSPerMin += cs
	}
	avgCSPerMin /= float64(len(earlyCS))

	highDeathRate := float64(len(earlyDeaths)) / float64(len(matches))

	var issues []string
	var tips []string

	if avgCSPerMin < 6 {
		issues = append(issues, "CS faible")
		tips = append(tips, "Focus sur le last-hitting parfait")
		tips = append(tips, "Évite les trades qui font rater des CS")
	}

	if highDeathRate > 0.3 { // 30% of games with high deaths
		issues = append(issues, "trop de morts early")
		tips = append(tips, "Play plus safe en early game")
		tips = append(tips, "Ward plus pour éviter les ganks")
	}

	if len(issues) > 0 {
		return &models.GameplayIssue{
			Type:        "early_game",
			Description: fmt.Sprintf("Problèmes détectés: %s. L'early game est crucial pour établir un avantage.", strings.Join(issues, ", ")),
			Tips:        tips[:int(math.Min(float64(len(tips)), 3))],
			Severity:    "medium",
		}
	}

	return nil
}

func (res *RecommendationEngineService) analyzeTeamfightPerformance(matches []map[string]interface{}) *models.GameplayIssue {
	if len(matches) == 0 {
		return nil
	}

	// Simplified teamfight analysis based on KDA and damage share
	lowImpactGames := 0

	for _, match := range matches {
		participantDataStr, ok := match["participant_data"].(string)
		if !ok {
			continue
		}

		var participantData map[string]interface{}
		if err := json.Unmarshal([]byte(participantDataStr), &participantData); err != nil {
			continue
		}

		kills, _ := participantData["kills"].(float64)
		deaths, _ := participantData["deaths"].(float64)
		assists, _ := participantData["assists"].(float64)

		if deaths == 0 {
			deaths = 1
		}
		kda := (kills + assists) / deaths

		if kda < 1.5 { // Low impact games
			lowImpactGames++
		}
	}

	lowImpactRate := float64(lowImpactGames) / float64(len(matches))

	if lowImpactRate > 0.4 { // 40% of games with low impact
		return &models.GameplayIssue{
			Type:        "teamfight",
			Description: fmt.Sprintf("Tu as un faible impact dans %.0f%% de tes games. Cela indique des problèmes de teamfight.", lowImpactRate*100),
			Tips: []string{
				"Améliore ton positionnement en teamfight",
				"Focus les bonnes cibles (ADC/Mid)",
				"Ne va pas frontline si tu n'es pas tank",
			},
			Severity: "medium",
		}
	}

	return nil
}

func (res *RecommendationEngineService) analyzeVisionControl(matches []map[string]interface{}) *models.GameplayIssue {
	if len(matches) == 0 {
		return nil
	}

	var visionScores []float64

	for _, match := range matches {
		participantDataStr, ok := match["participant_data"].(string)
		if !ok {
			continue
		}

		var participantData map[string]interface{}
		if err := json.Unmarshal([]byte(participantDataStr), &participantData); err != nil {
			continue
		}

		visionScore, _ := participantData["visionScore"].(float64)
		duration, _ := match["game_duration"].(int)
		if duration == 0 {
			duration = 1800
		}
		visionPerMin := visionScore / (float64(duration) / 60.0)
		visionScores = append(visionScores, visionPerMin)
	}

	if len(visionScores) == 0 {
		return nil
	}

	avgVisionPerMin := 0.0
	for _, vision := range visionScores {
		avgVisionPerMin += vision
	}
	avgVisionPerMin /= float64(len(visionScores))

	// Role-based vision expectations
	if avgVisionPerMin < 1.5 { // General low vision
		return &models.GameplayIssue{
			Type:        "vision",
			Description: fmt.Sprintf("Ton vision score (%.1f/min) est faible. La vision donne un avantage stratégique énorme.", avgVisionPerMin),
			Tips: []string{
				"Achète plus de wards de contrôle",
				"Ward les objectives avant qu'ils spawn",
				"Clear les wards ennemies quand possible",
			},
			Severity: "low",
		}
	}

	return nil
}

func (res *RecommendationEngineService) identifyProblemChampions(matches []map[string]interface{}) map[int]models.ProblemChampion {
	problemChampions := make(map[int]models.ProblemChampion)

	// This would analyze enemy team data to identify problematic champions
	// For now, return empty map as we don't have enemy data readily available
	// In a full implementation, this would require enemy team data from matches

	return problemChampions
}

func (res *RecommendationEngineService) getCurrentMetaThreats() map[int]float64 {
	// Return top meta threats (high pick/ban rate champions)
	threats := make(map[int]float64)
	for champID, strength := range models.ChampionMetaStrength {
		if strength > 0.9 {
			threats[champID] = strength
		}
	}

	return threats
}

func (res *RecommendationEngineService) analyzeSkillPercentiles(matches []map[string]interface{}) map[string]float64 {
	if len(matches) == 0 {
		return map[string]float64{}
	}

	// Calculate metrics
	var csValues, visionValues, kdaValues, damageValues []float64

	for _, match := range matches {
		participantDataStr, ok := match["participant_data"].(string)
		if !ok {
			continue
		}

		var participantData map[string]interface{}
		if err := json.Unmarshal([]byte(participantDataStr), &participantData); err != nil {
			continue
		}

		duration, _ := match["game_duration"].(int)
		if duration == 0 {
			duration = 1800
		}
		durationMin := float64(duration) / 60.0

		// CS per minute
		totalCS, _ := participantData["totalMinionsKilled"].(float64)
		neutralCS, _ := participantData["neutralMinionsKilled"].(float64)
		csPerMin := (totalCS + neutralCS) / durationMin
		csValues = append(csValues, csPerMin)

		// Vision per minute
		visionScore, _ := participantData["visionScore"].(float64)
		visionPerMin := visionScore / durationMin
		visionValues = append(visionValues, visionPerMin)

		// KDA
		kills, _ := participantData["kills"].(float64)
		deaths, _ := participantData["deaths"].(float64)
		assists, _ := participantData["assists"].(float64)
		if deaths == 0 {
			deaths = 1
		}
		kda := (kills + assists) / deaths
		kdaValues = append(kdaValues, kda)

		// Damage per minute
		damage, _ := participantData["totalDamageDealtToChampions"].(float64)
		damagePerMin := damage / durationMin
		damageValues = append(damageValues, damagePerMin)
	}

	// Calculate percentiles (simplified - compare to rough benchmarks)
	results := make(map[string]float64)

	if len(csValues) > 0 {
		avgCS := 0.0
		for _, cs := range csValues {
			avgCS += cs
		}
		avgCS /= float64(len(csValues))
		results["cs_per_min"] = math.Min((avgCS/8.0)*100, 100) // 8 CS/min = 100th percentile
	}

	if len(visionValues) > 0 {
		avgVision := 0.0
		for _, vision := range visionValues {
			avgVision += vision
		}
		avgVision /= float64(len(visionValues))
		results["vision_score"] = math.Min((avgVision/2.5)*100, 100) // 2.5 vision/min = 100th percentile
	}

	if len(kdaValues) > 0 {
		avgKDA := 0.0
		for _, kda := range kdaValues {
			avgKDA += kda
		}
		avgKDA /= float64(len(kdaValues))
		results["kda"] = math.Min((avgKDA/3.0)*100, 100) // 3.0 KDA = 100th percentile
	}

	if len(damageValues) > 0 {
		avgDamage := 0.0
		for _, damage := range damageValues {
			avgDamage += damage
		}
		avgDamage /= float64(len(damageValues))
		results["damage_per_min"] = math.Min((avgDamage/800)*100, 100) // 800 DPM = 100th percentile
	}

	return results
}

func (res *RecommendationEngineService) saveRecommendations(userID int, recommendations []models.Recommendation) error {
	// Clear old recommendations
	_, err := res.db.Exec("DELETE FROM performance_insights WHERE user_id = ?", userID)
	if err != nil {
		return err
	}

	// Insert new recommendations
	for _, rec := range recommendations {
		actionItemsJSON, _ := json.Marshal(rec.ActionItems)

		_, err := res.db.Exec(`
			INSERT INTO performance_insights 
			(user_id, insight_type, category, title, description, priority, 
			 confidence, time_period, expected_improvement, action_items, is_active)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			userID, string(rec.Type), string(rec.Type), rec.Title, rec.Description,
			rec.Priority, rec.Confidence, rec.TimePeriod, rec.ExpectedImprovement,
			string(actionItemsJSON), true)

		if err != nil {
			log.Printf("Failed to save recommendation: %v", err)
		}
	}

	return nil
}