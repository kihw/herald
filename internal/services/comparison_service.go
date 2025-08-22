package services

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"lol-match-exporter/internal/models"
)

// ComparisonService gère les comparaisons de statistiques entre membres
type ComparisonService struct {
	db             *sql.DB
	analyticsService *AnalyticsService
}

// NewComparisonService crée une nouvelle instance du service de comparaison
func NewComparisonService(db *sql.DB, analyticsService *AnalyticsService) *ComparisonService {
	return &ComparisonService{
		db:               db,
		analyticsService: analyticsService,
	}
}

// CreateComparison crée une nouvelle comparaison
func (cs *ComparisonService) CreateComparison(groupID, creatorID int, name, description, compareType string, parameters models.ComparisonParameters) (*models.GroupComparison, error) {
	// Valider les paramètres
	if name == "" {
		return nil, fmt.Errorf("comparison name cannot be empty")
	}
	
	validTypes := []string{"champions", "roles", "performance", "trends"}
	if !contains(validTypes, compareType) {
		return nil, fmt.Errorf("invalid comparison type: %s", compareType)
	}
	
	if len(parameters.MemberIDs) < 2 {
		return nil, fmt.Errorf("need at least 2 members to compare")
	}
	
	// Vérifier que tous les membres appartiennent au groupe
	err := cs.validateMembersInGroup(groupID, parameters.MemberIDs)
	if err != nil {
		return nil, fmt.Errorf("invalid members: %w", err)
	}
	
	now := time.Now()
	query := `
		INSERT INTO group_comparisons (group_id, creator_id, name, description, compare_type, parameters, is_public, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, false, ?, ?)
	`
	
	result, err := cs.db.Exec(query, groupID, creatorID, name, description, compareType, parameters, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create comparison: %w", err)
	}
	
	comparisonID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get comparison ID: %w", err)
	}
	
	// Générer les résultats de comparaison plus tard
	_ = comparisonID // éviter l'erreur unused
	
	err = cs.GenerateComparisonResults(int(comparisonID))
	if err != nil {
		return nil, fmt.Errorf("failed to generate results: %w", err)
	}
	
	return cs.GetComparisonByID(int(comparisonID))
}

// GetComparisonByID récupère une comparaison par son ID
func (cs *ComparisonService) GetComparisonByID(comparisonID int) (*models.GroupComparison, error) {
	query := `
		SELECT gc.id, gc.group_id, gc.creator_id, gc.name, gc.description, 
		       gc.compare_type, gc.parameters, gc.results, gc.is_public, 
		       gc.created_at, gc.updated_at,
		       g.name, g.description,
		       u.riot_id, u.riot_tag, u.region, u.profile_icon_id
		FROM group_comparisons gc
		JOIN groups g ON gc.group_id = g.id
		JOIN users u ON gc.creator_id = u.id
		WHERE gc.id = ?
	`
	
	row := cs.db.QueryRow(query, comparisonID)
	
	var comparison models.GroupComparison
	var group models.Group
	var creator models.User
	
	err := row.Scan(
		&comparison.ID, &comparison.GroupID, &comparison.CreatorID, &comparison.Name, &comparison.Description,
		&comparison.CompareType, &comparison.Parameters, &comparison.Results, &comparison.IsPublic,
		&comparison.CreatedAt, &comparison.UpdatedAt,
		&group.Name, &group.Description,
		&creator.RiotID, &creator.RiotTag, &creator.Region, &creator.ProfileIconID,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comparison not found")
		}
		return nil, fmt.Errorf("failed to get comparison: %w", err)
	}
	
	comparison.Group = &group
	comparison.Creator = &creator
	
	return &comparison, nil
}

// GetGroupComparisons récupère toutes les comparaisons d'un groupe
func (cs *ComparisonService) GetGroupComparisons(groupID int, limit int) ([]models.GroupComparison, error) {
	if limit <= 0 {
		limit = 20
	}
	
	query := `
		SELECT gc.id, gc.group_id, gc.creator_id, gc.name, gc.description,
		       gc.compare_type, gc.parameters, gc.results, gc.is_public,
		       gc.created_at, gc.updated_at,
		       u.riot_id, u.riot_tag, u.region, u.profile_icon_id
		FROM group_comparisons gc
		JOIN users u ON gc.creator_id = u.id
		WHERE gc.group_id = ?
		ORDER BY gc.created_at DESC
		LIMIT ?
	`
	
	rows, err := cs.db.Query(query, groupID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get group comparisons: %w", err)
	}
	defer rows.Close()
	
	var comparisons []models.GroupComparison
	
	for rows.Next() {
		var comparison models.GroupComparison
		var creator models.User
		
		err := rows.Scan(
			&comparison.ID, &comparison.GroupID, &comparison.CreatorID, &comparison.Name, &comparison.Description,
			&comparison.CompareType, &comparison.Parameters, &comparison.Results, &comparison.IsPublic,
			&comparison.CreatedAt, &comparison.UpdatedAt,
			&creator.RiotID, &creator.RiotTag, &creator.Region, &creator.ProfileIconID,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan comparison: %w", err)
		}
		
		comparison.Creator = &creator
		comparisons = append(comparisons, comparison)
	}
	
	return comparisons, nil
}

// GenerateComparisonResults génère les résultats d'une comparaison
func (cs *ComparisonService) GenerateComparisonResults(comparisonID int) error {
	comparison, err := cs.GetComparisonByID(comparisonID)
	if err != nil {
		return fmt.Errorf("failed to get comparison: %w", err)
	}
	
	var results models.ComparisonResults
	
	switch comparison.CompareType {
	case "champions":
		results, err = cs.generateChampionComparison(comparison)
	case "roles":
		results, err = cs.generateRoleComparison(comparison)
	case "performance":
		results, err = cs.generatePerformanceComparison(comparison)
	case "trends":
		results, err = cs.generateTrendComparison(comparison)
	default:
		return fmt.Errorf("unsupported comparison type: %s", comparison.CompareType)
	}
	
	if err != nil {
		return fmt.Errorf("failed to generate %s comparison: %w", comparison.CompareType, err)
	}
	
	results.GeneratedAt = time.Now()
	
	// Sauvegarder les résultats
	query := `
		UPDATE group_comparisons 
		SET results = ?, updated_at = ?
		WHERE id = ?
	`
	
	_, err = cs.db.Exec(query, results, time.Now(), comparisonID)
	if err != nil {
		return fmt.Errorf("failed to save comparison results: %w", err)
	}
	
	return nil
}

// generateChampionComparison génère une comparaison de champions
func (cs *ComparisonService) generateChampionComparison(comparison *models.GroupComparison) (models.ComparisonResults, error) {
	var results models.ComparisonResults
	
	// Simuler des données de champions pour les membres
	memberStats := make(map[string]interface{})
	var rankings []models.MemberRanking
	
	// Récupérer les statistiques de chaque membre
	for i, memberID := range comparison.Parameters.MemberIDs {
		// Ici, on simulerait un appel à l'API Riot ou à la base de données
		// Pour l'instant, on génère des données de test
		stats := cs.generateMemberChampionStats(memberID, comparison.Parameters.Champions)
		memberStats[fmt.Sprintf("member_%d", memberID)] = stats
		
		// Calculer un score pour le ranking
		score := cs.calculateChampionScore(stats)
		rankings = append(rankings, models.MemberRanking{
			UserID:   memberID,
			Username: fmt.Sprintf("Player%d", memberID),
			Rank:     i + 1,
			Score:    score,
			Metric:   "Champion Mastery",
			Change:   "same",
		})
	}
	
	// Trier les rankings par score
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})
	
	// Mettre à jour les rangs
	for i := range rankings {
		rankings[i].Rank = i + 1
	}
	
	// Créer les graphiques
	charts := []models.ChartData{
		{
			Type:  "radar",
			Title: "Champion Performance Comparison",
			Labels: []string{"Win Rate", "KDA", "CS/Min", "Damage", "Vision"},
			Datasets: []models.ChartDataset{
				{
					Label: "Member 1",
					Data:  []float64{65.5, 2.1, 7.2, 22500, 1.2},
					BackgroundColor: []string{"rgba(54, 162, 235, 0.2)"},
					BorderColor:     []string{"rgba(54, 162, 235, 1)"},
				},
				{
					Label: "Member 2",
					Data:  []float64{58.3, 1.8, 6.9, 20800, 1.0},
					BackgroundColor: []string{"rgba(255, 99, 132, 0.2)"},
					BorderColor:     []string{"rgba(255, 99, 132, 1)"},
				},
			},
		},
	}
	
	// Générer des insights
	insights := []string{
		"Player1 shows superior CS/min performance across all champions",
		"Player2 has more consistent KDA ratios",
		"Champion pool diversity varies significantly between players",
	}
	
	results = models.ComparisonResults{
		Summary: models.ComparisonSummary{
			TopPerformer:       rankings[0].Username,
			BestMetric:         "Champion Mastery",
			AverageWinRate:     61.9,
			TotalGamesCompared: 150,
			TimeSpan:           comparison.Parameters.TimeRange,
		},
		MemberStats:  memberStats,
		Charts:       charts,
		Rankings:     rankings,
		Insights:     insights,
	}
	
	return results, nil
}

// generateRoleComparison génère une comparaison de rôles
func (cs *ComparisonService) generateRoleComparison(comparison *models.GroupComparison) (models.ComparisonResults, error) {
	var results models.ComparisonResults
	
	memberStats := make(map[string]interface{})
	var rankings []models.MemberRanking
	
	// Générer des stats par rôle pour chaque membre
	for i, memberID := range comparison.Parameters.MemberIDs {
		stats := cs.generateMemberRoleStats(memberID, comparison.Parameters.Roles)
		memberStats[fmt.Sprintf("member_%d", memberID)] = stats
		
		score := cs.calculateRoleScore(stats)
		rankings = append(rankings, models.MemberRanking{
			UserID:   memberID,
			Username: fmt.Sprintf("Player%d", memberID),
			Rank:     i + 1,
			Score:    score,
			Metric:   "Role Performance",
			Change:   "same",
		})
	}
	
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})
	
	for i := range rankings {
		rankings[i].Rank = i + 1
	}
	
	charts := []models.ChartData{
		{
			Type:   "bar",
			Title:  "Win Rate by Role",
			Labels: []string{"TOP", "JUNGLE", "MID", "ADC", "SUPPORT"},
			Datasets: []models.ChartDataset{
				{
					Label: "Player 1",
					Data:  []float64{62.5, 58.0, 70.2, 55.8, 48.3},
					BackgroundColor: []string{"rgba(75, 192, 192, 0.6)"},
				},
				{
					Label: "Player 2", 
					Data:  []float64{55.1, 65.7, 52.4, 68.9, 60.2},
					BackgroundColor: []string{"rgba(255, 159, 64, 0.6)"},
				},
			},
		},
	}
	
	insights := []string{
		"Player1 excels in Mid lane with 70.2% win rate",
		"Player2 shows versatility across multiple roles", 
		"Support role needs improvement for both players",
	}
	
	results = models.ComparisonResults{
		Summary: models.ComparisonSummary{
			TopPerformer:       rankings[0].Username,
			BestMetric:         "Role Versatility",
			AverageWinRate:     59.8,
			TotalGamesCompared: 280,
			TimeSpan:           comparison.Parameters.TimeRange,
		},
		MemberStats:  memberStats,
		Charts:       charts,
		Rankings:     rankings,
		Insights:     insights,
	}
	
	return results, nil
}

// generatePerformanceComparison génère une comparaison de performance
func (cs *ComparisonService) generatePerformanceComparison(comparison *models.GroupComparison) (models.ComparisonResults, error) {
	var results models.ComparisonResults
	
	memberStats := make(map[string]interface{})
	var rankings []models.MemberRanking
	
	for i, memberID := range comparison.Parameters.MemberIDs {
		stats := cs.generateMemberPerformanceStats(memberID, comparison.Parameters.Metrics)
		memberStats[fmt.Sprintf("member_%d", memberID)] = stats
		
		score := cs.calculatePerformanceScore(stats, comparison.Parameters.Metrics)
		rankings = append(rankings, models.MemberRanking{
			UserID:   memberID,
			Username: fmt.Sprintf("Player%d", memberID),
			Rank:     i + 1,
			Score:    score,
			Metric:   "Overall Performance",
			Change:   "same",
		})
	}
	
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})
	
	for i := range rankings {
		rankings[i].Rank = i + 1
	}
	
	charts := []models.ChartData{
		{
			Type:   "line",
			Title:  "Performance Trends",
			Labels: []string{"Week 1", "Week 2", "Week 3", "Week 4"},
			Datasets: []models.ChartDataset{
				{
					Label: "Player 1 Win Rate",
					Data:  []float64{58.5, 62.1, 65.8, 63.2},
					BorderColor: []string{"rgba(54, 162, 235, 1)"},
				},
				{
					Label: "Player 2 Win Rate",
					Data:  []float64{52.3, 55.7, 59.4, 61.8},
					BorderColor: []string{"rgba(255, 99, 132, 1)"},
				},
			},
		},
	}
	
	insights := []string{
		fmt.Sprintf("%s shows the most consistent improvement trend", rankings[0].Username),
		"Average team damage has increased by 15% this month",
		"Vision control varies significantly between players",
	}
	
	results = models.ComparisonResults{
		Summary: models.ComparisonSummary{
			TopPerformer:       rankings[0].Username,
			BestMetric:         "Overall Performance",
			AverageWinRate:     60.5,
			TotalGamesCompared: 320,
			TimeSpan:           comparison.Parameters.TimeRange,
		},
		MemberStats:  memberStats,
		Charts:       charts,
		Rankings:     rankings,
		Insights:     insights,
	}
	
	return results, nil
}

// generateTrendComparison génère une comparaison de tendances
func (cs *ComparisonService) generateTrendComparison(comparison *models.GroupComparison) (models.ComparisonResults, error) {
	// Implémentation similaire aux autres types
	// Pour l'instant, utilisons la comparaison de performance comme base
	return cs.generatePerformanceComparison(comparison)
}

// Helper functions pour générer des données de test

func (cs *ComparisonService) generateMemberChampionStats(memberID int, champions []int) map[string]interface{} {
	return map[string]interface{}{
		"total_games":     50 + (memberID * 10),
		"win_rate":        55.0 + float64(memberID*5),
		"avg_kda":         1.8 + float64(memberID)*0.2,
		"avg_cs_per_min":  6.5 + float64(memberID)*0.5,
		"avg_damage":      18000 + (memberID * 2000),
		"favorite_champ":  fmt.Sprintf("Champion%d", memberID+1),
		"mastery_points":  150000 + (memberID * 25000),
	}
}

func (cs *ComparisonService) generateMemberRoleStats(memberID int, roles []string) map[string]interface{} {
	roleStats := make(map[string]map[string]float64)
	
	for _, role := range roles {
		roleStats[role] = map[string]float64{
			"games":      float64(20 + memberID*5),
			"win_rate":   50.0 + float64(memberID*8) + float64(len(role)),
			"avg_kda":    1.5 + float64(memberID)*0.3,
			"avg_damage": float64(15000 + memberID*3000),
		}
	}
	
	return map[string]interface{}{
		"roles": roleStats,
		"main_role": roles[0],
		"versatility_score": float64(memberID * 15),
	}
}

func (cs *ComparisonService) generateMemberPerformanceStats(memberID int, metrics []string) map[string]interface{} {
	stats := make(map[string]float64)
	
	for _, metric := range metrics {
		switch strings.ToLower(metric) {
		case "winrate":
			stats[metric] = 50.0 + float64(memberID*7)
		case "kda":
			stats[metric] = 1.2 + float64(memberID)*0.4
		case "cs":
			stats[metric] = 120 + float64(memberID*15)
		case "damage":
			stats[metric] = float64(18000 + memberID*4000)
		case "vision":
			stats[metric] = 0.8 + float64(memberID)*0.3
		default:
			stats[metric] = float64(memberID * 20)
		}
	}
	
	return map[string]interface{}{
		"metrics": stats,
		"trend": "improving",
		"consistency": 75.0 + float64(memberID*5),
	}
}

func (cs *ComparisonService) calculateChampionScore(stats map[string]interface{}) float64 {
	// Score basé sur les stats de champion
	winRate := getFloatFromMap(stats, "win_rate", 50.0)
	kda := getFloatFromMap(stats, "avg_kda", 1.0)
	return winRate*0.4 + kda*20*0.3 + 30*0.3
}

func (cs *ComparisonService) calculateRoleScore(stats map[string]interface{}) float64 {
	versatility := getFloatFromMap(stats, "versatility_score", 0.0)
	return versatility + 50.0
}

func (cs *ComparisonService) calculatePerformanceScore(stats map[string]interface{}, metrics []string) float64 {
	metricsMap, ok := stats["metrics"].(map[string]float64)
	if !ok {
		return 50.0
	}
	
	var totalScore float64
	for _, metric := range metrics {
		if value, exists := metricsMap[metric]; exists {
			// Normaliser selon le metric
			switch strings.ToLower(metric) {
			case "winrate":
				totalScore += value
			case "kda":
				totalScore += value * 20
			case "cs":
				totalScore += value / 3
			default:
				totalScore += value / 1000
			}
		}
	}
	
	return totalScore / float64(len(metrics))
}

func (cs *ComparisonService) validateMembersInGroup(groupID int, memberIDs []int) error {
	if len(memberIDs) == 0 {
		return fmt.Errorf("no members provided")
	}
	
	// Créer une liste de placeholders pour la query IN
	placeholders := strings.Repeat("?,", len(memberIDs)-1) + "?"
	query := fmt.Sprintf(`
		SELECT COUNT(*) FROM group_members 
		WHERE group_id = ? AND user_id IN (%s) AND status = 'active'
	`, placeholders)
	
	args := make([]interface{}, len(memberIDs)+1)
	args[0] = groupID
	for i, memberID := range memberIDs {
		args[i+1] = memberID
	}
	
	var count int
	err := cs.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to validate members: %w", err)
	}
	
	if count != len(memberIDs) {
		return fmt.Errorf("some members are not in the group")
	}
	
	return nil
}

func getFloatFromMap(m map[string]interface{}, key string, defaultValue float64) float64 {
	if value, ok := m[key]; ok {
		if fValue, ok := value.(float64); ok {
			return fValue
		}
		if iValue, ok := value.(int); ok {
			return float64(iValue)
		}
	}
	return defaultValue
}