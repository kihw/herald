package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type Database struct {
	*sql.DB
}

type Match struct {
	ID           string    `json:"id"`
	UserID       int       `json:"user_id"`
	GameCreation int64     `json:"game_creation"`
	GameDuration int       `json:"game_duration"`
	GameMode     string    `json:"game_mode"`
	QueueID      int       `json:"queue_id"`
	ChampionName string    `json:"champion_name"`
	ChampionID   int       `json:"champion_id"`
	Kills        int       `json:"kills"`
	Deaths       int       `json:"deaths"`
	Assists      int       `json:"assists"`
	Win          bool      `json:"win"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	ID           int       `json:"id"`
	RiotID       string    `json:"riot_id"`
	RiotTag      string    `json:"riot_tag"`
	RiotPUUID    string    `json:"riot_puuid"`
	Region       string    `json:"region"`
	SummonerLevel int      `json:"summoner_level"`
	LastSync     *time.Time `json:"last_sync"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewSQLiteDatabase creates a new SQLite database connection
func NewSQLiteDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{db}
	
	// Create tables
	if err := database.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("✅ SQLite database initialized successfully")
	return database, nil
}

// createTables creates the necessary tables
func (db *Database) createTables() error {
	// Users table
	usersSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		riot_id TEXT NOT NULL,
		riot_tag TEXT NOT NULL,
		riot_puuid TEXT UNIQUE NOT NULL,
		region TEXT NOT NULL,
		summoner_level INTEGER DEFAULT 0,
		last_sync DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Matches table
	matchesSQL := `
	CREATE TABLE IF NOT EXISTS matches (
		id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		game_creation INTEGER NOT NULL,
		game_duration INTEGER NOT NULL,
		game_mode TEXT NOT NULL,
		queue_id INTEGER NOT NULL,
		champion_name TEXT NOT NULL,
		champion_id INTEGER NOT NULL,
		kills INTEGER DEFAULT 0,
		deaths INTEGER DEFAULT 0,
		assists INTEGER DEFAULT 0,
		win BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id)
	);`

	// Indexes for faster queries
	indexQueries := []string{
		`CREATE INDEX IF NOT EXISTS idx_matches_user_id ON matches(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_matches_game_mode ON matches(game_mode);`,
		`CREATE INDEX IF NOT EXISTS idx_matches_game_creation ON matches(game_creation);`,
		`CREATE INDEX IF NOT EXISTS idx_matches_champion_name ON matches(champion_name);`,
		`CREATE INDEX IF NOT EXISTS idx_matches_user_game_creation ON matches(user_id, game_creation);`,
		`CREATE INDEX IF NOT EXISTS idx_matches_user_game_mode ON matches(user_id, game_mode);`,
	}

	allQueries := append([]string{usersSQL, matchesSQL}, indexQueries...)
	for _, query := range allQueries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}

// UpsertUser inserts or updates a user
func (db *Database) UpsertUser(user User) (int, error) {
	query := `
	INSERT INTO users (riot_id, riot_tag, riot_puuid, region, summoner_level, updated_at)
	VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(riot_puuid) DO UPDATE SET
		riot_id = excluded.riot_id,
		riot_tag = excluded.riot_tag,
		region = excluded.region,
		summoner_level = excluded.summoner_level,
		updated_at = CURRENT_TIMESTAMP
	RETURNING id`

	var userID int
	err := db.QueryRow(query, user.RiotID, user.RiotTag, user.RiotPUUID, user.Region, user.SummonerLevel).Scan(&userID)
	
	if err != nil {
		// Fallback for SQLite versions that don't support RETURNING
		query = `
		INSERT OR REPLACE INTO users (riot_id, riot_tag, riot_puuid, region, summoner_level, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`
		
		result, err := db.Exec(query, user.RiotID, user.RiotTag, user.RiotPUUID, user.Region, user.SummonerLevel)
		if err != nil {
			return 0, fmt.Errorf("failed to upsert user: %w", err)
		}
		
		id, err := result.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("failed to get user ID: %w", err)
		}
		userID = int(id)
	}

	return userID, nil
}

// SaveMatch saves a match to the database
func (db *Database) SaveMatch(match Match) error {
	query := `
	INSERT OR REPLACE INTO matches 
	(id, user_id, game_creation, game_duration, game_mode, queue_id, champion_name, champion_id, kills, deaths, assists, win)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query, 
		match.ID, match.UserID, match.GameCreation, match.GameDuration, 
		match.GameMode, match.QueueID, match.ChampionName, match.ChampionID,
		match.Kills, match.Deaths, match.Assists, match.Win)

	if err != nil {
		return fmt.Errorf("failed to save match: %w", err)
	}

	return nil
}

// GetMatchesByUser retrieves matches for a user with pagination
func (db *Database) GetMatchesByUser(userID int, limit, offset int) ([]Match, error) {
	query := `
	SELECT id, user_id, game_creation, game_duration, game_mode, queue_id, 
		   champion_name, champion_id, kills, deaths, assists, win, created_at
	FROM matches 
	WHERE user_id = ? 
	ORDER BY game_creation DESC 
	LIMIT ? OFFSET ?`

	rows, err := db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query matches: %w", err)
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		err := rows.Scan(
			&match.ID, &match.UserID, &match.GameCreation, &match.GameDuration,
			&match.GameMode, &match.QueueID, &match.ChampionName, &match.ChampionID,
			&match.Kills, &match.Deaths, &match.Assists, &match.Win, &match.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan match: %w", err)
		}
		matches = append(matches, match)
	}

	return matches, nil
}

// CountMatchesByUser counts matches for a user
func (db *Database) CountMatchesByUser(userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM matches WHERE user_id = ?`
	err := db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count matches: %w", err)
	}
	return count, nil
}

// UpdateUserLastSync updates the last sync time for a user
func (db *Database) UpdateUserLastSync(userID int) error {
	query := `UPDATE users SET last_sync = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last sync: %w", err)
	}
	return nil
}

// UserStats represents calculated user statistics
type UserStats struct {
	TotalMatches     int     `json:"total_matches"`
	WinRate          float64 `json:"win_rate"`
	AverageKDA       float64 `json:"average_kda"`
	FavoriteChampion string  `json:"favorite_champion"`
	TotalKills       int     `json:"total_kills"`
	TotalDeaths      int     `json:"total_deaths"`
	TotalAssists     int     `json:"total_assists"`
}

// ChampionStats represents champion-specific statistics
type ChampionStats struct {
	ChampionName string  `json:"champion_name"`
	Matches      int     `json:"matches"`
	Wins         int     `json:"wins"`
	WinRate      float64 `json:"win_rate"`
	AvgKDA       float64 `json:"avg_kda"`
	TotalKills   int     `json:"total_kills"`
	TotalDeaths  int     `json:"total_deaths"`
	TotalAssists int     `json:"total_assists"`
}

// GetUserStats calculates user statistics from database
func (db *Database) GetUserStats(userID int) (UserStats, error) {
	var stats UserStats
	
	// Get basic match statistics
	query := `
	SELECT 
		COUNT(*) as total_matches,
		SUM(CASE WHEN win = 1 THEN 1 ELSE 0 END) as wins,
		SUM(kills) as total_kills,
		SUM(deaths) as total_deaths,
		SUM(assists) as total_assists
	FROM matches 
	WHERE user_id = ?`
	
	var wins int
	err := db.QueryRow(query, userID).Scan(
		&stats.TotalMatches, &wins, &stats.TotalKills, &stats.TotalDeaths, &stats.TotalAssists)
	
	if err != nil {
		return stats, fmt.Errorf("failed to get basic stats: %w", err)
	}
	
	// Calculate win rate
	if stats.TotalMatches > 0 {
		stats.WinRate = float64(wins) / float64(stats.TotalMatches) * 100
	}
	
	// Calculate average KDA
	if stats.TotalDeaths > 0 {
		stats.AverageKDA = float64(stats.TotalKills+stats.TotalAssists) / float64(stats.TotalDeaths)
	} else if stats.TotalKills > 0 || stats.TotalAssists > 0 {
		stats.AverageKDA = float64(stats.TotalKills + stats.TotalAssists) // Perfect KDA
	}
	
	// Get favorite champion (most played)
	championQuery := `
	SELECT champion_name 
	FROM matches 
	WHERE user_id = ? 
	GROUP BY champion_name 
	ORDER BY COUNT(*) DESC 
	LIMIT 1`
	
	err = db.QueryRow(championQuery, userID).Scan(&stats.FavoriteChampion)
	if err != nil && err != sql.ErrNoRows {
		return stats, fmt.Errorf("failed to get favorite champion: %w", err)
	}
	
	return stats, nil
}

// GetChampionStats gets statistics for each champion played by user
func (db *Database) GetChampionStats(userID int) ([]ChampionStats, error) {
	query := `
	SELECT 
		champion_name,
		COUNT(*) as matches,
		SUM(CASE WHEN win = 1 THEN 1 ELSE 0 END) as wins,
		SUM(kills) as total_kills,
		SUM(deaths) as total_deaths,
		SUM(assists) as total_assists
	FROM matches 
	WHERE user_id = ? 
	GROUP BY champion_name 
	ORDER BY matches DESC`
	
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query champion stats: %w", err)
	}
	defer rows.Close()
	
	var champions []ChampionStats
	for rows.Next() {
		var champ ChampionStats
		err := rows.Scan(
			&champ.ChampionName, &champ.Matches, &champ.Wins, 
			&champ.TotalKills, &champ.TotalDeaths, &champ.TotalAssists)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan champion stats: %w", err)
		}
		
		// Calculate win rate
		if champ.Matches > 0 {
			champ.WinRate = float64(champ.Wins) / float64(champ.Matches) * 100
		}
		
		// Calculate average KDA
		if champ.TotalDeaths > 0 {
			champ.AvgKDA = float64(champ.TotalKills+champ.TotalAssists) / float64(champ.TotalDeaths)
		} else if champ.TotalKills > 0 || champ.TotalAssists > 0 {
			champ.AvgKDA = float64(champ.TotalKills + champ.TotalAssists)
		}
		
		champions = append(champions, champ)
	}
	
	return champions, nil
}

// GetRecentPerformance gets performance statistics for recent matches
func (db *Database) GetRecentPerformance(userID int, days int) (UserStats, error) {
	var stats UserStats
	
	// Calculate timestamp for X days ago
	daysAgo := time.Now().AddDate(0, 0, -days).Unix() * 1000
	
	query := `
	SELECT 
		COUNT(*) as total_matches,
		SUM(CASE WHEN win = 1 THEN 1 ELSE 0 END) as wins,
		SUM(kills) as total_kills,
		SUM(deaths) as total_deaths,
		SUM(assists) as total_assists
	FROM matches 
	WHERE user_id = ? AND game_creation >= ?`
	
	var wins int
	err := db.QueryRow(query, userID, daysAgo).Scan(
		&stats.TotalMatches, &wins, &stats.TotalKills, &stats.TotalDeaths, &stats.TotalAssists)
	
	if err != nil {
		return stats, fmt.Errorf("failed to get recent performance: %w", err)
	}
	
	// Calculate win rate
	if stats.TotalMatches > 0 {
		stats.WinRate = float64(wins) / float64(stats.TotalMatches) * 100
	}
	
	// Calculate average KDA
	if stats.TotalDeaths > 0 {
		stats.AverageKDA = float64(stats.TotalKills+stats.TotalAssists) / float64(stats.TotalDeaths)
	} else if stats.TotalKills > 0 || stats.TotalAssists > 0 {
		stats.AverageKDA = float64(stats.TotalKills + stats.TotalAssists)
	}
	
	return stats, nil
}

// GameModeStats represents statistics by game mode
type GameModeStats struct {
	GameMode     string  `json:"game_mode"`
	TotalMatches int     `json:"total_matches"`
	Wins         int     `json:"wins"`
	WinRate      float64 `json:"win_rate"`
	AvgKDA       float64 `json:"avg_kda"`
	TotalKills   int     `json:"total_kills"`
	TotalDeaths  int     `json:"total_deaths"`
	TotalAssists int     `json:"total_assists"`
	AvgDuration  float64 `json:"avg_duration_minutes"`
}

// PerformanceTrend represents performance over time
type PerformanceTrend struct {
	Date     string  `json:"date"`
	Matches  int     `json:"matches"`
	Wins     int     `json:"wins"`
	WinRate  float64 `json:"win_rate"`
	AvgKDA   float64 `json:"avg_kda"`
}

// GetStatsByGameMode gets statistics grouped by game mode
func (db *Database) GetStatsByGameMode(userID int) ([]GameModeStats, error) {
	query := `
	SELECT 
		game_mode,
		COUNT(*) as total_matches,
		SUM(CASE WHEN win = 1 THEN 1 ELSE 0 END) as wins,
		SUM(kills) as total_kills,
		SUM(deaths) as total_deaths,
		SUM(assists) as total_assists,
		AVG(CAST(game_duration AS FLOAT) / 60) as avg_duration_minutes
	FROM matches 
	WHERE user_id = ? 
	GROUP BY game_mode 
	ORDER BY total_matches DESC`
	
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query game mode stats: %w", err)
	}
	defer rows.Close()
	
	var gameModes []GameModeStats
	for rows.Next() {
		var mode GameModeStats
		err := rows.Scan(
			&mode.GameMode, &mode.TotalMatches, &mode.Wins,
			&mode.TotalKills, &mode.TotalDeaths, &mode.TotalAssists,
			&mode.AvgDuration)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan game mode stats: %w", err)
		}
		
		// Calculate win rate
		if mode.TotalMatches > 0 {
			mode.WinRate = float64(mode.Wins) / float64(mode.TotalMatches) * 100
		}
		
		// Calculate average KDA
		if mode.TotalDeaths > 0 {
			mode.AvgKDA = float64(mode.TotalKills+mode.TotalAssists) / float64(mode.TotalDeaths)
		} else if mode.TotalKills > 0 || mode.TotalAssists > 0 {
			mode.AvgKDA = float64(mode.TotalKills + mode.TotalAssists)
		}
		
		gameModes = append(gameModes, mode)
	}
	
	return gameModes, nil
}

// GetPerformanceTrend gets daily performance trends
func (db *Database) GetPerformanceTrend(userID int, days int) ([]PerformanceTrend, error) {
	query := `
	SELECT 
		DATE(game_creation / 1000, 'unixepoch') as match_date,
		COUNT(*) as matches,
		SUM(CASE WHEN win = 1 THEN 1 ELSE 0 END) as wins,
		SUM(kills) as total_kills,
		SUM(deaths) as total_deaths,
		SUM(assists) as total_assists
	FROM matches 
	WHERE user_id = ? 
		AND game_creation >= ? 
	GROUP BY DATE(game_creation / 1000, 'unixepoch')
	ORDER BY match_date DESC`
	
	daysAgo := time.Now().AddDate(0, 0, -days).Unix() * 1000
	
	rows, err := db.Query(query, userID, daysAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to query performance trend: %w", err)
	}
	defer rows.Close()
	
	var trends []PerformanceTrend
	for rows.Next() {
		var trend PerformanceTrend
		var totalKills, totalDeaths, totalAssists int
		
		err := rows.Scan(
			&trend.Date, &trend.Matches, &trend.Wins,
			&totalKills, &totalDeaths, &totalAssists)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan performance trend: %w", err)
		}
		
		// Calculate win rate
		if trend.Matches > 0 {
			trend.WinRate = float64(trend.Wins) / float64(trend.Matches) * 100
		}
		
		// Calculate average KDA
		if totalDeaths > 0 {
			trend.AvgKDA = float64(totalKills+totalAssists) / float64(totalDeaths)
		} else if totalKills > 0 || totalAssists > 0 {
			trend.AvgKDA = float64(totalKills + totalAssists)
		}
		
		trends = append(trends, trend)
	}
	
	return trends, nil
}

// GetMatchesForExport gets all matches with detailed info for CSV export
func (db *Database) GetMatchesForExport(userID int) ([]Match, error) {
	query := `
	SELECT id, user_id, game_creation, game_duration, game_mode, queue_id,
		   champion_name, champion_id, kills, deaths, assists, win, created_at
	FROM matches 
	WHERE user_id = ? 
	ORDER BY game_creation DESC`
	
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query matches for export: %w", err)
	}
	defer rows.Close()
	
	var matches []Match
	for rows.Next() {
		var match Match
		err := rows.Scan(
			&match.ID, &match.UserID, &match.GameCreation, &match.GameDuration,
			&match.GameMode, &match.QueueID, &match.ChampionName, &match.ChampionID,
			&match.Kills, &match.Deaths, &match.Assists, &match.Win, &match.CreatedAt)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan match for export: %w", err)
		}
		
		matches = append(matches, match)
	}
	
	return matches, nil
}

// Recommendation represents an AI-generated recommendation
type Recommendation struct {
	Type        string  `json:"type"`        // "champion", "gamemode", "improvement", "warning"
	Priority    string  `json:"priority"`    // "high", "medium", "low"
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Confidence  float64 `json:"confidence"`  // 0.0 to 1.0
	DataPoints  int     `json:"data_points"` // Number of matches analyzed
	Suggestion  string  `json:"suggestion"`
}

// PerformanceAnalysis represents detailed performance analysis
type PerformanceAnalysis struct {
	OverallTrend    string          `json:"overall_trend"`    // "improving", "stable", "declining"
	StreakInfo      StreakInfo      `json:"streak_info"`
	ChampionInsight ChampionInsight `json:"champion_insight"`
	GameModeInsight GameModeInsight `json:"gamemode_insight"`
	PlayTimeInsight PlayTimeInsight `json:"playtime_insight"`
}

type StreakInfo struct {
	CurrentStreak int    `json:"current_streak"`
	StreakType    string `json:"streak_type"` // "win", "loss", "mixed"
	LongestWin    int    `json:"longest_win"`
	LongestLoss   int    `json:"longest_loss"`
}

type ChampionInsight struct {
	BestChampion  string  `json:"best_champion"`
	BestWinRate   float64 `json:"best_winrate"`
	WorstChampion string  `json:"worst_champion"`
	WorstWinRate  float64 `json:"worst_winrate"`
	Diversity     float64 `json:"diversity"` // Champion pool diversity 0-1
}

type GameModeInsight struct {
	BestMode     string  `json:"best_mode"`
	BestModeWR   float64 `json:"best_mode_wr"`
	PreferredMode string `json:"preferred_mode"` // Most played
}

type PlayTimeInsight struct {
	AverageSessionLength float64 `json:"avg_session_length_hours"`
	PeakPerformanceHour  int     `json:"peak_performance_hour"` // 0-23
	TotalPlaytime        float64 `json:"total_playtime_hours"`
}

// GenerateRecommendations creates AI-powered recommendations based on performance data
func (db *Database) GenerateRecommendations(userID int) ([]Recommendation, error) {
	var recommendations []Recommendation
	
	// Get user stats for analysis
	stats, err := db.GetUserStats(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}
	
	// Get champion stats
	championStats, err := db.GetChampionStats(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get champion stats: %w", err)
	}
	
	// Get game mode stats
	gameModeStats, err := db.GetStatsByGameMode(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game mode stats: %w", err)
	}
	
	// Get recent performance for trend analysis
	recentPerf, err := db.GetRecentPerformance(userID, 7)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent performance: %w", err)
	}
	
	// Champion recommendations
	if len(championStats) >= 2 {
		bestChamp := championStats[0]
		
		// Recommend playing best champion more
		if bestChamp.WinRate > 70 && bestChamp.Matches >= 3 {
			recommendations = append(recommendations, Recommendation{
				Type:        "champion",
				Priority:    "high",
				Title:       fmt.Sprintf("Focus on %s", bestChamp.ChampionName),
				Description: fmt.Sprintf("You have an excellent %.1f%% winrate with %s over %d games", bestChamp.WinRate, bestChamp.ChampionName, bestChamp.Matches),
				Confidence:  0.85,
				DataPoints:  bestChamp.Matches,
				Suggestion:  fmt.Sprintf("Consider playing %s more often to climb ranks efficiently", bestChamp.ChampionName),
			})
		}
		
		// Find underperforming champions
		for _, champ := range championStats {
			if champ.Matches >= 3 && champ.WinRate < 40 {
				recommendations = append(recommendations, Recommendation{
					Type:        "warning",
					Priority:    "medium",
					Title:       fmt.Sprintf("Consider alternatives to %s", champ.ChampionName),
					Description: fmt.Sprintf("Your winrate with %s is %.1f%% over %d games", champ.ChampionName, champ.WinRate, champ.Matches),
					Confidence:  0.75,
					DataPoints:  champ.Matches,
					Suggestion:  fmt.Sprintf("Try practicing %s in normals or consider switching to a similar champion", champ.ChampionName),
				})
			}
		}
	}
	
	// Game mode recommendations
	if len(gameModeStats) >= 2 {
		var bestMode GameModeStats
		for _, mode := range gameModeStats {
			if mode.WinRate > bestMode.WinRate && mode.TotalMatches >= 2 {
				bestMode = mode
			}
		}
		
		if bestMode.WinRate > 70 {
			recommendations = append(recommendations, Recommendation{
				Type:        "gamemode",
				Priority:    "medium",
				Title:       fmt.Sprintf("Excel in %s mode", bestMode.GameMode),
				Description: fmt.Sprintf("You perform best in %s with %.1f%% winrate", bestMode.GameMode, bestMode.WinRate),
				Confidence:  0.70,
				DataPoints:  bestMode.TotalMatches,
				Suggestion:  fmt.Sprintf("Focus on %s games to maximize your climb potential", bestMode.GameMode),
			})
		}
	}
	
	// Performance trend recommendations
	if stats.TotalMatches >= 5 {
		// KDA improvement
		if stats.AverageKDA < 1.5 {
			recommendations = append(recommendations, Recommendation{
				Type:        "improvement",
				Priority:    "high",
				Title:       "Focus on KDA improvement",
				Description: fmt.Sprintf("Your average KDA is %.2f, which could be improved", stats.AverageKDA),
				Confidence:  0.80,
				DataPoints:  stats.TotalMatches,
				Suggestion:  "Focus on positioning, avoid risky plays, and prioritize objectives over kills",
			})
		}
		
		// Recent performance analysis
		if recentPerf.WinRate < 40 && recentPerf.TotalMatches >= 3 {
			recommendations = append(recommendations, Recommendation{
				Type:        "warning",
				Priority:    "high",
				Title:       "Recent performance decline",
				Description: fmt.Sprintf("Your recent winrate is %.1f%% over %d games", recentPerf.WinRate, recentPerf.TotalMatches),
				Confidence:  0.90,
				DataPoints:  recentPerf.TotalMatches,
				Suggestion:  "Consider taking a break, reviewing replays, or focusing on fundamentals",
			})
		} else if recentPerf.WinRate > 70 && recentPerf.TotalMatches >= 3 {
			recommendations = append(recommendations, Recommendation{
				Type:        "champion",
				Priority:    "low",
				Title:       "Great recent performance!",
				Description: fmt.Sprintf("You're on fire with %.1f%% winrate recently", recentPerf.WinRate),
				Confidence:  0.95,
				DataPoints:  recentPerf.TotalMatches,
				Suggestion:  "Keep up the momentum and continue with your current strategy",
			})
		}
	}
	
	return recommendations, nil
}

// GetPerformanceAnalysis provides detailed AI analysis of player performance
func (db *Database) GetPerformanceAnalysis(userID int) (PerformanceAnalysis, error) {
	var analysis PerformanceAnalysis
	
	// Get basic stats
	stats, err := db.GetUserStats(userID)
	if err != nil {
		return analysis, fmt.Errorf("failed to get user stats: %w", err)
	}
	
	// Get champion stats
	championStats, err := db.GetChampionStats(userID)
	if err != nil {
		return analysis, fmt.Errorf("failed to get champion stats: %w", err)
	}
	
	// Get game mode stats
	gameModeStats, err := db.GetStatsByGameMode(userID)
	if err != nil {
		return analysis, fmt.Errorf("failed to get game mode stats: %w", err)
	}
	
	// Analyze overall trend
	recent7, _ := db.GetRecentPerformance(userID, 7)
	recent14, _ := db.GetRecentPerformance(userID, 14)
	
	if recent7.WinRate > recent14.WinRate+10 {
		analysis.OverallTrend = "improving"
	} else if recent7.WinRate < recent14.WinRate-10 {
		analysis.OverallTrend = "declining"
	} else {
		analysis.OverallTrend = "stable"
	}
	
	// Analyze current streak
	analysis.StreakInfo = db.analyzeStreak(userID)
	
	// Champion insights
	if len(championStats) > 0 {
		analysis.ChampionInsight.BestChampion = championStats[0].ChampionName
		analysis.ChampionInsight.BestWinRate = championStats[0].WinRate
		analysis.ChampionInsight.Diversity = float64(len(championStats)) / 10.0 // Normalize to 0-1
		if analysis.ChampionInsight.Diversity > 1.0 {
			analysis.ChampionInsight.Diversity = 1.0
		}
		
		// Find worst performing champion (min 3 games)
		for _, champ := range championStats {
			if champ.Matches >= 3 {
				analysis.ChampionInsight.WorstChampion = champ.ChampionName
				analysis.ChampionInsight.WorstWinRate = champ.WinRate
			}
		}
	}
	
	// Game mode insights
	if len(gameModeStats) > 0 {
		analysis.GameModeInsight.PreferredMode = gameModeStats[0].GameMode
		
		var bestMode GameModeStats
		for _, mode := range gameModeStats {
			if mode.WinRate > bestMode.WinRate && mode.TotalMatches >= 2 {
				bestMode = mode
			}
		}
		analysis.GameModeInsight.BestMode = bestMode.GameMode
		analysis.GameModeInsight.BestModeWR = bestMode.WinRate
	}
	
	// Playtime insights (simplified calculation)
	totalHours := 0.0
	if stats.TotalMatches > 0 {
		// Estimate based on average game duration (assuming 25 min average)
		totalHours = float64(stats.TotalMatches) * 25.0 / 60.0
	}
	analysis.PlayTimeInsight.TotalPlaytime = totalHours
	analysis.PlayTimeInsight.AverageSessionLength = totalHours / 7.0 // Weekly average
	analysis.PlayTimeInsight.PeakPerformanceHour = 20 // Default to 8 PM
	
	return analysis, nil
}

// analyzeStreak analyzes win/loss streaks
func (db *Database) analyzeStreak(userID int) StreakInfo {
	var streak StreakInfo
	
	// Get recent matches ordered by date
	query := `
	SELECT win 
	FROM matches 
	WHERE user_id = ? 
	ORDER BY game_creation DESC 
	LIMIT 20`
	
	rows, err := db.Query(query, userID)
	if err != nil {
		return streak
	}
	defer rows.Close()
	
	var matches []bool
	for rows.Next() {
		var win bool
		rows.Scan(&win)
		matches = append(matches, win)
	}
	
	if len(matches) == 0 {
		return streak
	}
	
	// Calculate current streak
	currentType := matches[0]
	currentCount := 1
	
	for i := 1; i < len(matches); i++ {
		if matches[i] == currentType {
			currentCount++
		} else {
			break
		}
	}
	
	streak.CurrentStreak = currentCount
	if currentType {
		streak.StreakType = "win"
	} else {
		streak.StreakType = "loss"
	}
	
	// Calculate longest streaks
	longestWin := 0
	longestLoss := 0
	currentWinStreak := 0
	currentLossStreak := 0
	
	for _, win := range matches {
		if win {
			currentWinStreak++
			currentLossStreak = 0
			if currentWinStreak > longestWin {
				longestWin = currentWinStreak
			}
		} else {
			currentLossStreak++
			currentWinStreak = 0
			if currentLossStreak > longestLoss {
				longestLoss = currentLossStreak
			}
		}
	}
	
	streak.LongestWin = longestWin
	streak.LongestLoss = longestLoss
	
	return streak
}