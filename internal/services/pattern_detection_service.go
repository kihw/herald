package services

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

// PatternDetectionService détecte des patterns dans les données de match
type PatternDetectionService struct {
	db      *sql.DB
	enabled bool
}

// MatchPattern représente un pattern détecté dans les matches
type MatchPattern struct {
	ID          int       `json:"id"`
	PatternType string    `json:"pattern_type"` // win_streak, loss_streak, champion_synergy, role_preference, etc.
	UserID      int       `json:"user_id"`
	Description string    `json:"description"`
	Frequency   int       `json:"frequency"`
	Confidence  float64   `json:"confidence"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Metadata    string    `json:"metadata"` // JSON avec des détails supplémentaires
	CreatedAt   time.Time `json:"created_at"`
}

// ChampionSynergy représente une synergie entre champions
type ChampionSynergy struct {
	Champion1   string  `json:"champion1"`
	Champion2   string  `json:"champion2"`
	WinRate     float64 `json:"win_rate"`
	GamesPlayed int     `json:"games_played"`
	Synergy     float64 `json:"synergy_score"`
}

// NewPatternDetectionService crée une nouvelle instance du service
func NewPatternDetectionService(db *sql.DB) *PatternDetectionService {
	if db == nil {
		log.Println("Database not available, pattern detection service disabled")
		return &PatternDetectionService{enabled: false}
	}

	return &PatternDetectionService{
		db:      db,
		enabled: true,
	}
}

// IsEnabled retourne true si le service est activé
func (pds *PatternDetectionService) IsEnabled() bool {
	return pds.enabled && pds.db != nil
}

// DetectWinStreaks détecte les séries de victoires/défaites
func (pds *PatternDetectionService) DetectWinStreaks(userID int, minStreak int) ([]MatchPattern, error) {
	if !pds.IsEnabled() {
		return nil, fmt.Errorf("pattern detection service is not enabled")
	}

	query := `
		SELECT match_id, result, created_at, champion_name
		FROM matches 
		WHERE user_id = ? 
		ORDER BY created_at DESC 
		LIMIT 100
	`

	rows, err := pds.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query matches: %v", err)
	}
	defer rows.Close()

	var matches []struct {
		MatchID   string
		Result    string
		CreatedAt time.Time
		Champion  string
	}

	for rows.Next() {
		var match struct {
			MatchID   string
			Result    string
			CreatedAt time.Time
			Champion  string
		}
		err := rows.Scan(&match.MatchID, &match.Result, &match.CreatedAt, &match.Champion)
		if err != nil {
			continue
		}
		matches = append(matches, match)
	}

	var patterns []MatchPattern
	currentStreak := 1
	currentResult := ""
	streakStart := time.Time{}

	for i, match := range matches {
		if i == 0 {
			currentResult = match.Result
			streakStart = match.CreatedAt
			continue
		}

		if match.Result == currentResult {
			currentStreak++
		} else {
			// Fin de streak
			if currentStreak >= minStreak {
				patternType := "win_streak"
				if currentResult != "Victory" {
					patternType = "loss_streak"
				}

				pattern := MatchPattern{
					PatternType: patternType,
					UserID:      userID,
					Description: fmt.Sprintf("%s of %d games", strings.Title(strings.Replace(patternType, "_", " ", -1)), currentStreak),
					Frequency:   currentStreak,
					Confidence:  0.9,
					StartDate:   matches[i-1].CreatedAt,
					EndDate:     streakStart,
					Metadata:    fmt.Sprintf(`{"streak_length": %d, "result": "%s"}`, currentStreak, currentResult),
					CreatedAt:   time.Now(),
				}
				patterns = append(patterns, pattern)
			}

			// Reset pour nouveau streak
			currentStreak = 1
			currentResult = match.Result
			streakStart = match.CreatedAt
		}
	}

	// Vérifier le dernier streak
	if currentStreak >= minStreak {
		patternType := "win_streak"
		if currentResult != "Victory" {
			patternType = "loss_streak"
		}

		pattern := MatchPattern{
			PatternType: patternType,
			UserID:      userID,
			Description: fmt.Sprintf("%s of %d games (ongoing)", strings.Title(strings.Replace(patternType, "_", " ", -1)), currentStreak),
			Frequency:   currentStreak,
			Confidence:  0.9,
			StartDate:   matches[len(matches)-1].CreatedAt,
			EndDate:     streakStart,
			Metadata:    fmt.Sprintf(`{"streak_length": %d, "result": "%s", "ongoing": true}`, currentStreak, currentResult),
			CreatedAt:   time.Now(),
		}
		patterns = append(patterns, pattern)
	}

	return patterns, nil
}

// DetectChampionPreferences détecte les préférences de champions
func (pds *PatternDetectionService) DetectChampionPreferences(userID int) ([]MatchPattern, error) {
	if !pds.IsEnabled() {
		return nil, fmt.Errorf("pattern detection service is not enabled")
	}

	query := `
		SELECT 
			champion_name,
			COUNT(*) as play_count,
			AVG(CASE WHEN result = 'Victory' THEN 1.0 ELSE 0.0 END) * 100 as win_rate
		FROM matches 
		WHERE user_id = ? 
		AND created_at >= datetime('now', '-30 days')
		GROUP BY champion_name
		HAVING play_count >= 3
		ORDER BY play_count DESC, win_rate DESC
	`

	rows, err := pds.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query champion preferences: %v", err)
	}
	defer rows.Close()

	var patterns []MatchPattern
	rank := 1

	for rows.Next() {
		var champion string
		var playCount int
		var winRate float64

		err := rows.Scan(&champion, &playCount, &winRate)
		if err != nil {
			continue
		}

		var patternType string
		var confidence float64

		if rank <= 3 && playCount >= 5 {
			patternType = "champion_main"
			confidence = 0.9
		} else if winRate >= 70 && playCount >= 3 {
			patternType = "champion_comfort_pick"
			confidence = 0.8
		} else if playCount >= 10 {
			patternType = "champion_frequent_pick"
			confidence = 0.7
		} else {
			rank++
			continue
		}

		pattern := MatchPattern{
			PatternType: patternType,
			UserID:      userID,
			Description: fmt.Sprintf("Plays %s frequently (%d games, %.1f%% winrate)", champion, playCount, winRate),
			Frequency:   playCount,
			Confidence:  confidence,
			StartDate:   time.Now().AddDate(0, 0, -30),
			EndDate:     time.Now(),
			Metadata:    fmt.Sprintf(`{"champion": "%s", "play_count": %d, "win_rate": %.1f, "rank": %d}`, champion, playCount, winRate, rank),
			CreatedAt:   time.Now(),
		}

		patterns = append(patterns, pattern)
		rank++
	}

	return patterns, nil
}

// DetectRolePreferences détecte les préférences de rôles
func (pds *PatternDetectionService) DetectRolePreferences(userID int) ([]MatchPattern, error) {
	if !pds.IsEnabled() {
		return nil, fmt.Errorf("pattern detection service is not enabled")
	}

	query := `
		SELECT 
			role,
			COUNT(*) as play_count,
			AVG(CASE WHEN result = 'Victory' THEN 1.0 ELSE 0.0 END) * 100 as win_rate
		FROM matches 
		WHERE user_id = ? 
		AND created_at >= datetime('now', '-30 days')
		AND role != ''
		GROUP BY role
		ORDER BY play_count DESC
	`

	rows, err := pds.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query role preferences: %v", err)
	}
	defer rows.Close()

	var patterns []MatchPattern
	totalGames := 0
	roles := make(map[string]struct {
		playCount int
		winRate   float64
	})

	// Compter le total de games et collecter les stats par rôle
	for rows.Next() {
		var role string
		var playCount int
		var winRate float64

		err := rows.Scan(&role, &playCount, &winRate)
		if err != nil {
			continue
		}

		roles[role] = struct {
			playCount int
			winRate   float64
		}{playCount, winRate}
		totalGames += playCount
	}

	// Analyser les patterns
	for role, stats := range roles {
		percentage := float64(stats.playCount) / float64(totalGames) * 100

		var patternType string
		var confidence float64

		if percentage >= 60 {
			patternType = "role_one_trick"
			confidence = 0.9
		} else if percentage >= 40 {
			patternType = "role_main"
			confidence = 0.8
		} else if percentage >= 25 {
			patternType = "role_secondary"
			confidence = 0.7
		} else {
			continue
		}

		pattern := MatchPattern{
			PatternType: patternType,
			UserID:      userID,
			Description: fmt.Sprintf("Prefers %s role (%.1f%% of games, %.1f%% winrate)", role, percentage, stats.winRate),
			Frequency:   stats.playCount,
			Confidence:  confidence,
			StartDate:   time.Now().AddDate(0, 0, -30),
			EndDate:     time.Now(),
			Metadata:    fmt.Sprintf(`{"role": "%s", "percentage": %.1f, "win_rate": %.1f}`, role, percentage, stats.winRate),
			CreatedAt:   time.Now(),
		}

		patterns = append(patterns, pattern)
	}

	return patterns, nil
}

// DetectPlayTimePatterns détecte les patterns de temps de jeu
func (pds *PatternDetectionService) DetectPlayTimePatterns(userID int) ([]MatchPattern, error) {
	if !pds.IsEnabled() {
		return nil, fmt.Errorf("pattern detection service is not enabled")
	}

	query := `
		SELECT 
			strftime('%H', created_at) as hour,
			COUNT(*) as game_count,
			AVG(CASE WHEN result = 'Victory' THEN 1.0 ELSE 0.0 END) * 100 as win_rate
		FROM matches 
		WHERE user_id = ? 
		AND created_at >= datetime('now', '-30 days')
		GROUP BY hour
		HAVING game_count >= 3
		ORDER BY game_count DESC
	`

	rows, err := pds.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query playtime patterns: %v", err)
	}
	defer rows.Close()

	var patterns []MatchPattern

	for rows.Next() {
		var hour string
		var gameCount int
		var winRate float64

		err := rows.Scan(&hour, &gameCount, &winRate)
		if err != nil {
			continue
		}

		if gameCount >= 10 {
			var timeOfDay string
			hourInt := 0
			fmt.Sscanf(hour, "%d", &hourInt)

			if hourInt >= 6 && hourInt < 12 {
				timeOfDay = "morning"
			} else if hourInt >= 12 && hourInt < 18 {
				timeOfDay = "afternoon"
			} else if hourInt >= 18 && hourInt < 24 {
				timeOfDay = "evening"
			} else {
				timeOfDay = "night"
			}

			pattern := MatchPattern{
				PatternType: "playtime_preference",
				UserID:      userID,
				Description: fmt.Sprintf("Active in %s (%s:00, %d games, %.1f%% winrate)", timeOfDay, hour, gameCount, winRate),
				Frequency:   gameCount,
				Confidence:  0.7,
				StartDate:   time.Now().AddDate(0, 0, -30),
				EndDate:     time.Now(),
				Metadata:    fmt.Sprintf(`{"hour": "%s", "time_of_day": "%s", "game_count": %d, "win_rate": %.1f}`, hour, timeOfDay, gameCount, winRate),
				CreatedAt:   time.Now(),
			}

			patterns = append(patterns, pattern)
		}
	}

	return patterns, nil
}

// DetectAllPatterns détecte tous les types de patterns pour un utilisateur
func (pds *PatternDetectionService) DetectAllPatterns(userID int) ([]MatchPattern, error) {
	if !pds.IsEnabled() {
		return nil, fmt.Errorf("pattern detection service is not enabled")
	}

	var allPatterns []MatchPattern

	// Détecter les win streaks (minimum 3)
	streaks, err := pds.DetectWinStreaks(userID, 3)
	if err == nil {
		allPatterns = append(allPatterns, streaks...)
	}

	// Détecter les préférences de champions
	champPrefs, err := pds.DetectChampionPreferences(userID)
	if err == nil {
		allPatterns = append(allPatterns, champPrefs...)
	}

	// Détecter les préférences de rôles
	rolePrefs, err := pds.DetectRolePreferences(userID)
	if err == nil {
		allPatterns = append(allPatterns, rolePrefs...)
	}

	// Détecter les patterns de temps de jeu
	timePatterns, err := pds.DetectPlayTimePatterns(userID)
	if err == nil {
		allPatterns = append(allPatterns, timePatterns...)
	}

	log.Printf("Detected %d patterns for user %d", len(allPatterns), userID)
	return allPatterns, nil
}
