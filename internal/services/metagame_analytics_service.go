package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"lol-match-exporter/internal/models"
)

// MetagameAnalyticsService fournit les analyses de métagame avancées
type MetagameAnalyticsService struct {
	db      *sql.DB
	enabled bool
}

// NewMetagameAnalyticsService crée une nouvelle instance du service d'analytics métagame
func NewMetagameAnalyticsService(db *sql.DB) *MetagameAnalyticsService {
	if db == nil {
		log.Println("Database not available, metagame analytics service disabled")
		return &MetagameAnalyticsService{enabled: false}
	}

	return &MetagameAnalyticsService{
		db:      db,
		enabled: true,
	}
}

// IsEnabled retourne true si le service est activé
func (mas *MetagameAnalyticsService) IsEnabled() bool {
	return mas.enabled && mas.db != nil
}

// AnalyzeChampionMeta analyse les métriques des champions sur une période donnée
func (mas *MetagameAnalyticsService) AnalyzeChampionMeta(days int) ([]models.ChampionMetrics, error) {
	if !mas.IsEnabled() {
		return nil, fmt.Errorf("metagame analytics service is not enabled")
	}

	// Requête SQL pour récupérer les statistiques des champions
	query := `
		SELECT 
			champion_name,
			role,
			COUNT(*) as games_played,
			SUM(CASE WHEN result = 'Victory' THEN 1 ELSE 0 END) as wins,
			COUNT(*) - SUM(CASE WHEN result = 'Victory' THEN 1 ELSE 0 END) as losses,
			ROUND(AVG((kills + assists) / NULLIF(deaths, 0)), 2) as avg_kda,
			ROUND(AVG(damage_dealt), 0) as avg_damage,
			ROUND(AVG(gold_earned), 0) as avg_gold,
			ROUND(AVG(cs_total), 0) as avg_cs
		FROM matches 
		WHERE created_at >= datetime('now', '-' || ? || ' days')
		GROUP BY champion_name, role
		HAVING games_played >= 5
		ORDER BY games_played DESC
	`

	rows, err := mas.db.Query(query, days)
	if err != nil {
		return nil, fmt.Errorf("failed to query champion metrics: %v", err)
	}
	defer rows.Close()

	var champions []models.ChampionMetrics
	for rows.Next() {
		var champ models.ChampionMetrics
		err := rows.Scan(
			&champ.ChampionName,
			&champ.Role,
			&champ.GamesPlayed,
			&champ.Wins,
			&champ.Losses,
			&champ.AvgKDA,
			&champ.AvgDamage,
			&champ.AvgGold,
			&champ.AvgCS,
		)
		if err != nil {
			log.Printf("Error scanning champion metrics: %v", err)
			continue
		}

		// Calculer le winrate et autres métriques
		if champ.GamesPlayed > 0 {
			champ.WinRate = float64(champ.Wins) / float64(champ.GamesPlayed) * 100
		}

		// Calculer un score de performance
		champ.TierScore = mas.calculateTierScore(champ)
		champ.TrendDirection = mas.calculateTrend(champ.ChampionName, days)
		champ.UpdatedAt = time.Now()

		champions = append(champions, champ)
	}

	log.Printf("Analyzed %d champions for %d days period", len(champions), days)
	return champions, nil
}

// GetMetagameSnapshot crée un snapshot du métagame actuel
func (mas *MetagameAnalyticsService) GetMetagameSnapshot(region, tier string) (*models.MetagameSnapshot, error) {
	if !mas.IsEnabled() {
		return nil, fmt.Errorf("metagame analytics service is not enabled")
	}

	snapshot := &models.MetagameSnapshot{
		Timestamp: time.Now(),
		Patch:     "14.15", // Serait récupéré dynamiquement en production
		Region:    region,
		Tier:      tier,
		CreatedAt: time.Now(),
	}

	// Sauvegarder le snapshot en base (optionnel)
	if mas.db != nil {
		query := `
			INSERT INTO metagame_snapshots (timestamp, patch, region, tier, created_at)
			VALUES (?, ?, ?, ?, ?)
		`
		result, err := mas.db.Exec(query, snapshot.Timestamp, snapshot.Patch, snapshot.Region, snapshot.Tier, snapshot.CreatedAt)
		if err != nil {
			log.Printf("Failed to save metagame snapshot: %v", err)
		} else {
			id, _ := result.LastInsertId()
			snapshot.ID = int(id)
		}
	}

	return snapshot, nil
}

// DetectMetagameTrends détecte les tendances du métagame
func (mas *MetagameAnalyticsService) DetectMetagameTrends(timeframe string) ([]models.MetagameTrend, error) {
	if !mas.IsEnabled() {
		return nil, fmt.Errorf("metagame analytics service is not enabled")
	}

	var trends []models.MetagameTrend

	// Exemple de détection de tendance simple
	champions, err := mas.AnalyzeChampionMeta(7) // 7 derniers jours
	if err != nil {
		return nil, err
	}

	for _, champ := range champions {
		var trendType string
		var confidence float64 = 0.7 // Confidence par défaut

		if champ.WinRate > 55 && champ.GamesPlayed > 20 {
			trendType = "rising"
			confidence = 0.9
		} else if champ.WinRate < 45 && champ.GamesPlayed > 20 {
			trendType = "falling"
			confidence = 0.8
		} else {
			trendType = "stable"
		}

		trend := models.MetagameTrend{
			Type:          "champion",
			EntityID:      champ.ChampionID,
			EntityName:    champ.ChampionName,
			TrendType:     trendType,
			ChangePercent: champ.WinRate - 50, // Écart par rapport à 50%
			Timeframe:     timeframe,
			Confidence:    confidence,
			Description:   fmt.Sprintf("%s is %s in %s role with %.1f%% winrate", champ.ChampionName, trendType, champ.Role, champ.WinRate),
			StartDate:     time.Now().AddDate(0, 0, -7),
			EndDate:       time.Now(),
			CreatedAt:     time.Now(),
		}

		trends = append(trends, trend)
	}

	return trends, nil
}

// GetMetagameStats calcule les statistiques globales du métagame
func (mas *MetagameAnalyticsService) GetMetagameStats() (*models.MetagameStats, error) {
	if !mas.IsEnabled() {
		return nil, fmt.Errorf("metagame analytics service is not enabled")
	}

	stats := &models.MetagameStats{
		CreatedAt: time.Now(),
	}

	// Compter le total de games
	err := mas.db.QueryRow("SELECT COUNT(*) FROM matches WHERE created_at >= datetime('now', '-7 days')").Scan(&stats.TotalGames)
	if err != nil {
		log.Printf("Failed to get total games: %v", err)
	}

	// Compter les champions uniques
	err = mas.db.QueryRow("SELECT COUNT(DISTINCT champion_name) FROM matches WHERE created_at >= datetime('now', '-7 days')").Scan(&stats.UniqueChampions)
	if err != nil {
		log.Printf("Failed to get unique champions: %v", err)
	}

	// Durée moyenne des games
	err = mas.db.QueryRow("SELECT AVG(game_duration) FROM matches WHERE created_at >= datetime('now', '-7 days')").Scan(&stats.AvgGameDuration)
	if err != nil {
		log.Printf("Failed to get avg game duration: %v", err)
	}

	// Champion le plus joué
	err = mas.db.QueryRow(`
		SELECT champion_name FROM matches 
		WHERE created_at >= datetime('now', '-7 days') 
		GROUP BY champion_name 
		ORDER BY COUNT(*) DESC 
		LIMIT 1
	`).Scan(&stats.MostPickedChampion)
	if err != nil {
		log.Printf("Failed to get most picked champion: %v", err)
	}

	// Winrate le plus élevé
	err = mas.db.QueryRow(`
		SELECT MAX(winrate) FROM (
			SELECT 
				champion_name,
				ROUND(AVG(CASE WHEN result = 'Victory' THEN 1.0 ELSE 0.0 END) * 100, 2) as winrate
			FROM matches 
			WHERE created_at >= datetime('now', '-7 days')
			GROUP BY champion_name
			HAVING COUNT(*) >= 10
		)
	`).Scan(&stats.HighestWinRate)
	if err != nil {
		log.Printf("Failed to get highest winrate: %v", err)
	}

	// Calculer l'index de diversité (simple)
	stats.DiversityIndex = float64(stats.UniqueChampions) / float64(stats.TotalGames) * 100

	// Niveau de puissance du métagame
	if stats.HighestWinRate > 60 {
		stats.PowerLevel = "high"
	} else if stats.HighestWinRate > 55 {
		stats.PowerLevel = "medium"
	} else {
		stats.PowerLevel = "low"
	}

	return stats, nil
}

// calculateTierScore calcule un score de tier basé sur les performances
func (mas *MetagameAnalyticsService) calculateTierScore(champ models.ChampionMetrics) float64 {
	// Formule simple: (winrate * 0.4) + (games_played_weight * 0.3) + (kda_weight * 0.3)
	winrateScore := champ.WinRate / 100 * 40

	gamesWeight := float64(champ.GamesPlayed) / 100 * 30
	if gamesWeight > 30 {
		gamesWeight = 30
	}

	kdaWeight := champ.AvgKDA / 3 * 30
	if kdaWeight > 30 {
		kdaWeight = 30
	}

	return winrateScore + gamesWeight + kdaWeight
}

// calculateTrend calcule la tendance d'un champion
func (mas *MetagameAnalyticsService) calculateTrend(championName string, days int) string {
	// Implémentation simple - comparer les 3 derniers jours vs les 3 précédents
	recentQuery := `
		SELECT AVG(CASE WHEN result = 'Victory' THEN 1.0 ELSE 0.0 END) * 100
		FROM matches 
		WHERE champion_name = ? AND created_at >= datetime('now', '-3 days')
	`

	previousQuery := `
		SELECT AVG(CASE WHEN result = 'Victory' THEN 1.0 ELSE 0.0 END) * 100
		FROM matches 
		WHERE champion_name = ? 
		AND created_at >= datetime('now', '-6 days') 
		AND created_at < datetime('now', '-3 days')
	`

	var recentWinrate, previousWinrate float64

	err1 := mas.db.QueryRow(recentQuery, championName).Scan(&recentWinrate)
	err2 := mas.db.QueryRow(previousQuery, championName).Scan(&previousWinrate)

	if err1 != nil || err2 != nil {
		return "stable"
	}

	diff := recentWinrate - previousWinrate
	if diff > 3 {
		return "rising"
	} else if diff < -3 {
		return "falling"
	}

	return "stable"
}
