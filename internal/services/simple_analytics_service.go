package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// SimpleAnalyticsService provides basic analytics using SQLite
type SimpleAnalyticsService struct {
	db *sql.DB
}

// NewSimpleAnalyticsService creates a new simple analytics service for SQLite
func NewSimpleAnalyticsService(db *sql.DB) *SimpleAnalyticsService {
	return &SimpleAnalyticsService{
		db: db,
	}
}

// PeriodStats represents analytics for a specific time period
type SimplePeriodStats struct {
	Period           string                       `json:"period"`
	TotalGames       int                         `json:"total_games"`
	WinRate          float64                     `json:"win_rate"`
	AvgKDA           float64                     `json:"avg_kda"`
	BestRole         string                      `json:"best_role"`
	WorstRole        string                      `json:"worst_role"`
	TopChampions     []SimpleChampionPerformance `json:"top_champions"`
	RolePerformance  map[string]interface{}      `json:"role_performance"`
	RecentTrend      string                      `json:"recent_trend"`
	Suggestions      []string                    `json:"suggestions"`
}

// SimpleChampionPerformance represents champion performance data
type SimpleChampionPerformance struct {
	ChampionID       int     `json:"champion_id"`
	ChampionName     string  `json:"champion_name"`
	Games            int     `json:"games"`
	WinRate          float64 `json:"win_rate"`
	PerformanceScore float64 `json:"performance_score"`
	AvgKDA           float64 `json:"avg_kda"`
}

// SimpleRecommendation represents a simple recommendation
type SimpleRecommendation struct {
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

// GetPeriodStats récupère les analytics pour une période donnée
func (s *SimpleAnalyticsService) GetPeriodStats(userID int, period string) (*SimplePeriodStats, error) {
	// Déterminer la période de temps
	startTime := s.getPeriodStartTime(period)
	
	// Récupérer les matches de la période
	query := `
		SELECT 
			m.game_creation,
			mp.champion_name,
			mp.kills,
			mp.deaths,
			mp.assists,
			mp.total_damage_dealt_to_champions,
			mp.gold_earned,
			mp.total_minions_killed,
			mp.win
		FROM match_participants mp
		JOIN matches m ON mp.match_id = m.id
		WHERE mp.user_id = ? AND m.game_creation >= ?
		ORDER BY m.game_creation DESC
	`
	
	rows, err := s.db.Query(query, userID, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query matches: %w", err)
	}
	defer rows.Close()

	var matches []SimpleMatchData
	totalGames := 0
	totalWins := 0
	totalKills := 0
	totalDeaths := 0
	totalAssists := 0
	championStats := make(map[string]*SimpleChampionPerformance)

	for rows.Next() {
		var match SimpleMatchData
		var championName sql.NullString
		
		err := rows.Scan(
			&match.GameCreation,
			&championName,
			&match.Kills,
			&match.Deaths,
			&match.Assists,
			&match.Damage,
			&match.Gold,
			&match.CS,
			&match.Win,
		)
		if err != nil {
			log.Printf("Error scanning match data: %v", err)
			continue
		}

		match.ChampionName = championName.String
		matches = append(matches, match)
		
		totalGames++
		if match.Win {
			totalWins++
		}
		totalKills += match.Kills
		totalDeaths += match.Deaths
		totalAssists += match.Assists

		// Agrégation par champion
		if championName.Valid && championName.String != "" {
			champName := championName.String
			if stats, exists := championStats[champName]; exists {
				stats.Games++
				if match.Win {
					stats.WinRate = (stats.WinRate*float64(stats.Games-1) + 100) / float64(stats.Games)
				} else {
					stats.WinRate = (stats.WinRate * float64(stats.Games-1)) / float64(stats.Games)
				}
				stats.AvgKDA = (stats.AvgKDA*float64(stats.Games-1) + match.getKDA()) / float64(stats.Games)
			} else {
				winRate := 0.0
				if match.Win {
					winRate = 100.0
				}
				championStats[champName] = &SimpleChampionPerformance{
					ChampionID:   0, // Pas d'ID champion pour l'instant
					ChampionName: champName,
					Games:        1,
					WinRate:      winRate,
					AvgKDA:      match.getKDA(),
				}
			}
		}
	}

	if totalGames == 0 {
		return &SimplePeriodStats{
			Period:          period,
			TotalGames:      0,
			WinRate:         0,
			AvgKDA:         0,
			BestRole:        "Unknown",
			WorstRole:       "Unknown",
			TopChampions:    []SimpleChampionPerformance{},
			RolePerformance: map[string]interface{}{},
			RecentTrend:     "stable",
			Suggestions:     []string{"Jouez plus de parties pour obtenir des analyses détaillées"},
		}, nil
	}

	// Calculer les statistiques générales
	winRate := float64(totalWins) / float64(totalGames) * 100
	avgKDA := float64(totalKills+totalAssists) / float64(max(totalDeaths, 1))

	// Calculer les scores de performance pour les champions
	for _, stats := range championStats {
		stats.PerformanceScore = s.calculatePerformanceScore(stats)
	}

	// Convertir en slice et trier par performance
	topChampions := make([]SimpleChampionPerformance, 0, len(championStats))
	for _, stats := range championStats {
		topChampions = append(topChampions, *stats)
	}

	// Trier par score de performance décroissant
	for i := 0; i < len(topChampions)-1; i++ {
		for j := i + 1; j < len(topChampions); j++ {
			if topChampions[i].PerformanceScore < topChampions[j].PerformanceScore {
				topChampions[i], topChampions[j] = topChampions[j], topChampions[i]
			}
		}
	}

	// Limiter aux 10 meilleurs champions
	if len(topChampions) > 10 {
		topChampions = topChampions[:10]
	}

	// Déterminer la tendance
	trend := s.calculateTrend(matches)

	// Générer des suggestions
	suggestions := s.generateSuggestions(winRate, avgKDA, topChampions, totalGames)

	return &SimplePeriodStats{
		Period:          period,
		TotalGames:      totalGames,
		WinRate:         winRate,
		AvgKDA:         avgKDA,
		BestRole:        "ADC", // TODO: Calculer à partir des données réelles
		WorstRole:       "Support", // TODO: Calculer à partir des données réelles
		TopChampions:    topChampions,
		RolePerformance: map[string]interface{}{},
		RecentTrend:     trend,
		Suggestions:     suggestions,
	}, nil
}

// SimpleMatchData structure for internal processing
type SimpleMatchData struct {
	GameCreation  int64
	ChampionName  string
	Kills         int
	Deaths        int
	Assists       int
	Damage        int
	Gold          int
	CS            int
	Win           bool
}

func (m SimpleMatchData) getKDA() float64 {
	if m.Deaths == 0 {
		return float64(m.Kills + m.Assists)
	}
	return float64(m.Kills+m.Assists) / float64(m.Deaths)
}

// getPeriodStartTime retourne le timestamp de début pour une période donnée
func (s *SimpleAnalyticsService) getPeriodStartTime(period string) int64 {
	now := time.Now()
	
	switch period {
	case "today":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	case "week":
		// Début de la semaine (lundi)
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Dimanche = 7
		}
		return now.AddDate(0, 0, -weekday+1).Unix()
	case "month":
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Unix()
	case "season":
		// Approximation: début de l'année
		return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()).Unix()
	default:
		// Par défaut, dernière semaine
		return now.AddDate(0, 0, -7).Unix()
	}
}

// calculatePerformanceScore calcule un score de performance pour un champion
func (s *SimpleAnalyticsService) calculatePerformanceScore(stats *SimpleChampionPerformance) float64 {
	// Score basé sur le taux de victoire, KDA, et nombre de parties
	winRateScore := stats.WinRate
	kdaScore := minFloat(stats.AvgKDA*20, 100) // Cap à 100
	gamesWeight := minFloat(float64(stats.Games)*5, 25) // Plus de parties = plus de poids
	
	score := (winRateScore*0.5 + kdaScore*0.3) * (1 + gamesWeight/100)
	return minFloat(score, 100)
}

// calculateTrend détermine la tendance récente des performances
func (s *SimpleAnalyticsService) calculateTrend(matches []SimpleMatchData) string {
	if len(matches) < 5 {
		return "stable"
	}

	// Analyser les 5 dernières parties vs les 5 précédentes
	recent := matches[:minInt(5, len(matches))]
	older := matches[minInt(5, len(matches)):minInt(10, len(matches))]

	if len(older) == 0 {
		return "stable"
	}

	recentWins := 0
	for _, match := range recent {
		if match.Win {
			recentWins++
		}
	}

	olderWins := 0
	for _, match := range older {
		if match.Win {
			olderWins++
		}
	}

	recentWinRate := float64(recentWins) / float64(len(recent))
	olderWinRate := float64(olderWins) / float64(len(older))

	if recentWinRate > olderWinRate+0.2 {
		return "improving"
	} else if recentWinRate < olderWinRate-0.2 {
		return "declining"
	}
	return "stable"
}

// generateSuggestions génère des suggestions basées sur les performances
func (s *SimpleAnalyticsService) generateSuggestions(winRate, avgKDA float64, topChampions []SimpleChampionPerformance, totalGames int) []string {
	var suggestions []string

	if totalGames < 10 {
		suggestions = append(suggestions, "Jouez plus de parties pour obtenir des analyses plus précises")
	}

	if winRate < 45 {
		suggestions = append(suggestions, "Concentrez-vous sur vos champions les plus performants pour améliorer votre taux de victoire")
	}

	if avgKDA < 2.0 {
		suggestions = append(suggestions, "Travaillez sur votre positionnement pour réduire vos morts et augmenter vos participations")
	}

	if len(topChampions) > 0 {
		bestChamp := topChampions[0]
		if bestChamp.Games >= 3 && bestChamp.WinRate > 60 {
			suggestions = append(suggestions, fmt.Sprintf("Continuez à jouer %s, c'est votre champion le plus performant (%.1f%% de victoires)", bestChamp.ChampionName, bestChamp.WinRate))
		}
	}

	if winRate > 60 && avgKDA > 2.5 {
		suggestions = append(suggestions, "Excellentes performances ! Continuez sur cette lancée")
	}

	if len(suggestions) == 0 {
		suggestions = append(suggestions, "Maintenez votre niveau de jeu actuel")
	}

	return suggestions
}

// GetRecommendations génère des recommandations simples
func (s *SimpleAnalyticsService) GetRecommendations(userID int) ([]SimpleRecommendation, error) {
	// Pour l'instant, retournons des recommandations statiques
	recommendations := []SimpleRecommendation{
		{
			Type:                "performance",
			Title:               "Améliorer votre KDA",
			Description:         "Votre ratio KDA pourrait être amélioré en travaillant sur votre positionnement",
			Priority:            2,
			Confidence:          0.75,
			ExpectedImprovement: "+10% winrate",
			ActionItems:         []string{"Restez derrière vos tanks", "Évitez les trades défavorables", "Placez plus de wards"},
			TimePeriod:          "2 semaines",
		},
		{
			Type:                "champion",
			Title:               "Diversifier votre pool de champions",
			Description:         "Jouez avec 2-3 champions pour avoir plus d'options",
			Priority:            3,
			Confidence:          0.65,
			ExpectedImprovement: "+5% winrate",
			ActionItems:         []string{"Choisissez 2 nouveaux champions", "Pratiquez en normal draft", "Regardez des guides"},
			TimePeriod:          "1 mois",
		},
	}

	return recommendations, nil
}

// GetPerformanceTrends génère des tendances de performance simples
func (s *SimpleAnalyticsService) GetPerformanceTrends(userID int) (map[string]interface{}, error) {
	// Pour l'instant, retournons des tendances statiques
	trends := map[string]interface{}{
		"daily_trend": map[string]interface{}{
			"trend":    "stable",
			"win_rate": 52.5,
		},
		"weekly_trend": map[string]interface{}{
			"trend":    "improving",
			"win_rate": 58.2,
		},
		"consistency_score":   78.5,
		"peak_performance": map[string]interface{}{
			"performance": 85.0,
		},
	}

	return trends, nil
}

// UpdateChampionStats met à jour les statistiques (placeholder)
func (s *SimpleAnalyticsService) UpdateChampionStats(userID int, period string) error {
	// Pour l'instant, ne fait rien - juste un placeholder
	log.Printf("Updating champion stats for user %d (period: %s)", userID, period)
	return nil
}

// Fonction utilitaire minFloat
func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Fonction utilitaire minInt
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Fonction utilitaire max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}