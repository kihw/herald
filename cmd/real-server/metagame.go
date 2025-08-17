package main

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"
)

// MetaGameAnalytics provides advanced meta-game analysis
type MetaGameAnalytics struct {
	database *Database
}

// ChampionMetrics represents champion performance in the meta
type ChampionMetrics struct {
	ChampionName    string  `json:"champion_name"`
	ChampionID      int     `json:"champion_id"`
	TotalPicks      int     `json:"total_picks"`
	WinRate         float64 `json:"win_rate"`
	PickRate        float64 `json:"pick_rate"`
	BanRate         float64 `json:"ban_rate,omitempty"`
	AverageKDA      float64 `json:"average_kda"`
	TrendDirection  string  `json:"trend_direction"` // "rising", "stable", "falling"
	TrendStrength   float64 `json:"trend_strength"`  // 0.0 to 1.0
	MetaScore       float64 `json:"meta_score"`      // Composite score 0-100
	LastUpdated     time.Time `json:"last_updated"`
}

// GameModeMetrics represents meta-game data by game mode
type GameModeMetrics struct {
	GameMode        string            `json:"game_mode"`
	QueueID         int               `json:"queue_id"`
	TotalGames      int               `json:"total_games"`
	AverageGameTime float64           `json:"average_game_time_minutes"`
	TopChampions    []ChampionMetrics `json:"top_champions"`
	WinRatesByRole  map[string]float64 `json:"winrates_by_role"`
	LastUpdated     time.Time         `json:"last_updated"`
}

// ItemBuildPattern represents popular item builds
type ItemBuildPattern struct {
	ChampionName   string    `json:"champion_name"`
	BuildHash      string    `json:"build_hash"`
	Items          []string  `json:"items"`
	CoreItems      []string  `json:"core_items"`
	Popularity     float64   `json:"popularity"`     // 0.0 to 1.0
	WinRate        float64   `json:"win_rate"`
	AverageKDA     float64   `json:"average_kda"`
	SampleSize     int       `json:"sample_size"`
	GameMode       string    `json:"game_mode"`
	LastSeen       time.Time `json:"last_seen"`
}

// TeamComposition represents team composition analysis
type TeamComposition struct {
	CompositionID   string            `json:"composition_id"`
	Roles           map[string]string `json:"roles"` // role -> champion
	WinRate         float64           `json:"win_rate"`
	PickFrequency   int               `json:"pick_frequency"`
	Synergies       []ChampionSynergy `json:"synergies"`
	Counters        []ChampionCounter `json:"counters"`
	StrengthPhases  []string          `json:"strength_phases"` // "early", "mid", "late"
	LastAnalyzed    time.Time         `json:"last_analyzed"`
}

// ChampionSynergy represents champion synergy data
type ChampionSynergy struct {
	Champion1       string  `json:"champion1"`
	Champion2       string  `json:"champion2"`
	SynergyScore    float64 `json:"synergy_score"` // 0.0 to 1.0
	WinRateTogether float64 `json:"winrate_together"`
	WinRateSeparate float64 `json:"winrate_separate"`
	SampleSize      int     `json:"sample_size"`
}

// ChampionCounter represents counter relationships
type ChampionCounter struct {
	Champion        string  `json:"champion"`
	CounterChampion string  `json:"counter_champion"`
	CounterStrength float64 `json:"counter_strength"` // 0.0 to 1.0
	WinRateAgainst  float64 `json:"winrate_against"`
	LanePhase       string  `json:"lane_phase"` // "early", "mid", "late"
	SampleSize      int     `json:"sample_size"`
}

// TrendAnalysis represents meta trends over time
type TrendAnalysis struct {
	Period          string               `json:"period"` // "daily", "weekly", "monthly"
	StartDate       time.Time            `json:"start_date"`
	EndDate         time.Time            `json:"end_date"`
	RisingChampions []ChampionTrendData  `json:"rising_champions"`
	FallingChampions []ChampionTrendData `json:"falling_champions"`
	EmergingBuilds  []ItemBuildPattern   `json:"emerging_builds"`
	MetaShifts      []MetaShiftEvent     `json:"meta_shifts"`
	LastCalculated  time.Time            `json:"last_calculated"`
}

// ChampionTrendData represents champion trend information
type ChampionTrendData struct {
	ChampionName    string    `json:"champion_name"`
	CurrentWinRate  float64   `json:"current_winrate"`
	PreviousWinRate float64   `json:"previous_winrate"`
	ChangePercent   float64   `json:"change_percent"`
	PickRateChange  float64   `json:"pickrate_change"`
	TrendConfidence float64   `json:"trend_confidence"` // 0.0 to 1.0
	LastUpdated     time.Time `json:"last_updated"`
}

// MetaShiftEvent represents significant meta changes
type MetaShiftEvent struct {
	EventType       string    `json:"event_type"` // "champion_nerf", "item_change", "role_shift"
	Description     string    `json:"description"`
	AffectedChampions []string `json:"affected_champions"`
	Impact          string    `json:"impact"` // "major", "moderate", "minor"
	DetectedAt      time.Time `json:"detected_at"`
	Confidence      float64   `json:"confidence"`
}

// NewMetaGameAnalytics creates a new meta-game analytics instance
func NewMetaGameAnalytics(db *Database) *MetaGameAnalytics {
	return &MetaGameAnalytics{
		database: db,
	}
}

// AnalyzeChampionMeta analyzes champion performance in the current meta
func (mga *MetaGameAnalytics) AnalyzeChampionMeta(days int) ([]ChampionMetrics, error) {
	log.Printf("🔍 Analyzing champion meta for last %d days", days)
	
	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	startTimestamp := startDate.Unix() * 1000
	
	// Query champion performance data
	query := `
	SELECT 
		champion_name,
		champion_id,
		COUNT(*) as total_picks,
		SUM(CASE WHEN win = 1 THEN 1 ELSE 0 END) as wins,
		SUM(kills) as total_kills,
		SUM(deaths) as total_deaths,
		SUM(assists) as total_assists
	FROM matches 
	WHERE game_creation >= ? 
	GROUP BY champion_name, champion_id
	HAVING total_picks >= 3
	ORDER BY total_picks DESC`
	
	rows, err := mga.database.Query(query, startTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to query champion data: %w", err)
	}
	defer rows.Close()
	
	var totalGames int
	var champions []ChampionMetrics
	
	// First pass: collect data and calculate total games
	for rows.Next() {
		var champ ChampionMetrics
		var wins, totalKills, totalDeaths, totalAssists int
		
		err := rows.Scan(
			&champ.ChampionName, &champ.ChampionID, &champ.TotalPicks,
			&wins, &totalKills, &totalDeaths, &totalAssists)
		if err != nil {
			continue
		}
		
		// Calculate basic metrics
		champ.WinRate = float64(wins) / float64(champ.TotalPicks) * 100
		if totalDeaths > 0 {
			champ.AverageKDA = float64(totalKills+totalAssists) / float64(totalDeaths)
		} else {
			champ.AverageKDA = float64(totalKills + totalAssists)
		}
		
		champ.LastUpdated = time.Now()
		champions = append(champions, champ)
		totalGames += champ.TotalPicks
	}
	
	// Second pass: calculate pick rates and meta scores
	for i := range champions {
		champions[i].PickRate = float64(champions[i].TotalPicks) / float64(totalGames) * 100
		champions[i].MetaScore = mga.calculateMetaScore(champions[i])
		champions[i].TrendDirection, champions[i].TrendStrength = mga.calculateTrend(champions[i], days)
	}
	
	// Sort by meta score
	sort.Slice(champions, func(i, j int) bool {
		return champions[i].MetaScore > champions[j].MetaScore
	})
	
	log.Printf("📊 Analyzed %d champions in current meta", len(champions))
	return champions, nil
}

// calculateMetaScore computes a composite meta score for a champion
func (mga *MetaGameAnalytics) calculateMetaScore(champ ChampionMetrics) float64 {
	// Weighted scoring: WinRate (40%), PickRate (30%), KDA (20%), Trend (10%)
	winRateScore := math.Min(champ.WinRate, 100) / 100 * 40
	pickRateScore := math.Min(champ.PickRate*10, 100) / 100 * 30 // Scale pick rate
	kdaScore := math.Min(champ.AverageKDA*20, 100) / 100 * 20    // Scale KDA
	
	var trendScore float64 = 5 // Neutral trend
	if champ.TrendDirection == "rising" {
		trendScore = 5 + (champ.TrendStrength * 5)
	} else if champ.TrendDirection == "falling" {
		trendScore = 5 - (champ.TrendStrength * 5)
	}
	
	return winRateScore + pickRateScore + kdaScore + trendScore
}

// calculateTrend determines if a champion is trending up/down
func (mga *MetaGameAnalytics) calculateTrend(champ ChampionMetrics, days int) (string, float64) {
	// For now, use a simplified trend calculation
	// In a real implementation, we'd compare with previous periods
	
	// Champions with high win rate and decent pick rate are "rising"
	if champ.WinRate > 55 && champ.PickRate > 2 {
		strength := math.Min((champ.WinRate-50)/50, 1.0)
		return "rising", strength
	}
	
	// Champions with low win rate are "falling"
	if champ.WinRate < 45 {
		strength := math.Min((50-champ.WinRate)/50, 1.0)
		return "falling", strength
	}
	
	return "stable", 0.0
}

// AnalyzeGameModeMeta analyzes meta by game mode
func (mga *MetaGameAnalytics) AnalyzeGameModeMeta(days int) ([]GameModeMetrics, error) {
	log.Printf("🎮 Analyzing game mode meta for last %d days", days)
	
	startTimestamp := time.Now().AddDate(0, 0, -days).Unix() * 1000
	
	// Query game mode data
	query := `
	SELECT 
		game_mode,
		queue_id,
		COUNT(*) as total_games,
		AVG(CAST(game_duration AS FLOAT) / 60) as avg_duration
	FROM matches 
	WHERE game_creation >= ? 
	GROUP BY game_mode, queue_id
	ORDER BY total_games DESC`
	
	rows, err := mga.database.Query(query, startTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to query game mode data: %w", err)
	}
	defer rows.Close()
	
	var gameModes []GameModeMetrics
	
	for rows.Next() {
		var mode GameModeMetrics
		err := rows.Scan(&mode.GameMode, &mode.QueueID, &mode.TotalGames, &mode.AverageGameTime)
		if err != nil {
			continue
		}
		
		// Get top champions for this game mode
		mode.TopChampions, _ = mga.getTopChampionsByMode(mode.GameMode, mode.QueueID, days, 10)
		mode.WinRatesByRole = make(map[string]float64)
		mode.LastUpdated = time.Now()
		
		gameModes = append(gameModes, mode)
	}
	
	log.Printf("📊 Analyzed %d game modes", len(gameModes))
	return gameModes, nil
}

// getTopChampionsByMode gets top champions for a specific game mode
func (mga *MetaGameAnalytics) getTopChampionsByMode(gameMode string, queueID, days, limit int) ([]ChampionMetrics, error) {
	startTimestamp := time.Now().AddDate(0, 0, -days).Unix() * 1000
	
	query := `
	SELECT 
		champion_name,
		champion_id,
		COUNT(*) as total_picks,
		SUM(CASE WHEN win = 1 THEN 1 ELSE 0 END) as wins,
		SUM(kills) as total_kills,
		SUM(deaths) as total_deaths,
		SUM(assists) as total_assists
	FROM matches 
	WHERE game_creation >= ? AND game_mode = ? AND queue_id = ?
	GROUP BY champion_name, champion_id
	HAVING total_picks >= 2
	ORDER BY total_picks DESC
	LIMIT ?`
	
	rows, err := mga.database.Query(query, startTimestamp, gameMode, queueID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var champions []ChampionMetrics
	
	for rows.Next() {
		var champ ChampionMetrics
		var wins, totalKills, totalDeaths, totalAssists int
		
		err := rows.Scan(
			&champ.ChampionName, &champ.ChampionID, &champ.TotalPicks,
			&wins, &totalKills, &totalDeaths, &totalAssists)
		if err != nil {
			continue
		}
		
		champ.WinRate = float64(wins) / float64(champ.TotalPicks) * 100
		if totalDeaths > 0 {
			champ.AverageKDA = float64(totalKills+totalAssists) / float64(totalDeaths)
		} else {
			champ.AverageKDA = float64(totalKills + totalAssists)
		}
		
		champ.LastUpdated = time.Now()
		champions = append(champions, champ)
	}
	
	return champions, nil
}

// DetectMetaShifts detects significant changes in the meta
func (mga *MetaGameAnalytics) DetectMetaShifts(compareDay1, compareDay2 int) ([]MetaShiftEvent, error) {
	log.Printf("🔄 Detecting meta shifts between %d and %d days ago", compareDay2, compareDay1)
	
	// Get champion data for both periods
	period1Champions, err := mga.AnalyzeChampionMeta(compareDay1)
	if err != nil {
		return nil, err
	}
	
	period2Champions, err := mga.AnalyzeChampionMeta(compareDay2)
	if err != nil {
		return nil, err
	}
	
	// Create maps for easier comparison
	period1Map := make(map[string]ChampionMetrics)
	period2Map := make(map[string]ChampionMetrics)
	
	for _, champ := range period1Champions {
		period1Map[champ.ChampionName] = champ
	}
	for _, champ := range period2Champions {
		period2Map[champ.ChampionName] = champ
	}
	
	var shifts []MetaShiftEvent
	
	// Detect significant winrate changes
	for champName, current := range period1Map {
		if previous, exists := period2Map[champName]; exists {
			winRateChange := current.WinRate - previous.WinRate
			pickRateChange := current.PickRate - previous.PickRate
			
			// Significant change thresholds
			if math.Abs(winRateChange) > 5 || math.Abs(pickRateChange) > 2 {
				impact := "minor"
				if math.Abs(winRateChange) > 10 || math.Abs(pickRateChange) > 5 {
					impact = "moderate"
				}
				if math.Abs(winRateChange) > 15 || math.Abs(pickRateChange) > 8 {
					impact = "major"
				}
				
				eventType := "performance_change"
				description := fmt.Sprintf("%s winrate changed by %.1f%% (%.1f%% → %.1f%%)", 
					champName, winRateChange, previous.WinRate, current.WinRate)
				
				if pickRateChange > 0 {
					description += fmt.Sprintf(", pick rate increased by %.1f%%", pickRateChange)
				} else if pickRateChange < 0 {
					description += fmt.Sprintf(", pick rate decreased by %.1f%%", math.Abs(pickRateChange))
				}
				
				shifts = append(shifts, MetaShiftEvent{
					EventType:         eventType,
					Description:       description,
					AffectedChampions: []string{champName},
					Impact:            impact,
					DetectedAt:        time.Now(),
					Confidence:        math.Min(math.Abs(winRateChange)/20, 1.0),
				})
			}
		}
	}
	
	log.Printf("🔍 Detected %d meta shifts", len(shifts))
	return shifts, nil
}

// GenerateMetaReport generates a comprehensive meta analysis report
func (mga *MetaGameAnalytics) GenerateMetaReport(days int) (map[string]interface{}, error) {
	log.Printf("📋 Generating comprehensive meta report for last %d days", days)
	
	// Analyze champion meta
	championMeta, err := mga.AnalyzeChampionMeta(days)
	if err != nil {
		return nil, err
	}
	
	// Analyze game mode meta
	gameModeMeta, err := mga.AnalyzeGameModeMeta(days)
	if err != nil {
		return nil, err
	}
	
	// Detect recent meta shifts
	metaShifts, err := mga.DetectMetaShifts(3, 7) // Compare last 3 days vs previous 7 days
	if err != nil {
		log.Printf("⚠️ Failed to detect meta shifts: %v", err)
		metaShifts = []MetaShiftEvent{} // Continue without meta shifts
	}
	
	// Create comprehensive report
	report := map[string]interface{}{
		"meta_analysis": map[string]interface{}{
			"period_days":        days,
			"analysis_date":      time.Now().Format(time.RFC3339),
			"total_champions":    len(championMeta),
			"total_game_modes":   len(gameModeMeta),
		},
		"champion_meta": map[string]interface{}{
			"top_performers":    mga.getTopPerformers(championMeta, 10),
			"rising_champions":  mga.getByTrend(championMeta, "rising", 5),
			"falling_champions": mga.getByTrend(championMeta, "falling", 5),
			"balanced_champions": mga.getBalancedChampions(championMeta, 5),
		},
		"game_mode_analysis": gameModeMeta,
		"meta_shifts":        metaShifts,
		"recommendations": map[string]interface{}{
			"strong_picks":    mga.getStrongPicks(championMeta),
			"avoid_picks":     mga.getAvoidPicks(championMeta),
			"trend_picks":     mga.getTrendPicks(championMeta),
		},
		"statistics": map[string]interface{}{
			"average_winrate":   mga.calculateAverageWinRate(championMeta),
			"meta_diversity":    mga.calculateMetaDiversity(championMeta),
			"balance_score":     mga.calculateBalanceScore(championMeta),
		},
	}
	
	log.Printf("✅ Meta report generated successfully")
	return report, nil
}

// Helper functions for report generation
func (mga *MetaGameAnalytics) getTopPerformers(champions []ChampionMetrics, limit int) []ChampionMetrics {
	if len(champions) > limit {
		return champions[:limit]
	}
	return champions
}

func (mga *MetaGameAnalytics) getByTrend(champions []ChampionMetrics, trend string, limit int) []ChampionMetrics {
	var filtered []ChampionMetrics
	for _, champ := range champions {
		if champ.TrendDirection == trend {
			filtered = append(filtered, champ)
		}
	}
	
	// Sort by trend strength
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].TrendStrength > filtered[j].TrendStrength
	})
	
	if len(filtered) > limit {
		return filtered[:limit]
	}
	return filtered
}

func (mga *MetaGameAnalytics) getBalancedChampions(champions []ChampionMetrics, limit int) []ChampionMetrics {
	var balanced []ChampionMetrics
	for _, champ := range champions {
		// Balanced: 48-52% win rate, decent pick rate
		if champ.WinRate >= 48 && champ.WinRate <= 52 && champ.PickRate >= 1 {
			balanced = append(balanced, champ)
		}
	}
	
	sort.Slice(balanced, func(i, j int) bool {
		return balanced[i].PickRate > balanced[j].PickRate
	})
	
	if len(balanced) > limit {
		return balanced[:limit]
	}
	return balanced
}

func (mga *MetaGameAnalytics) getStrongPicks(champions []ChampionMetrics) []string {
	var strong []string
	for _, champ := range champions {
		if champ.WinRate > 55 && champ.PickRate > 2 && champ.TotalPicks >= 10 {
			strong = append(strong, champ.ChampionName)
		}
	}
	return strong
}

func (mga *MetaGameAnalytics) getAvoidPicks(champions []ChampionMetrics) []string {
	var avoid []string
	for _, champ := range champions {
		if champ.WinRate < 45 && champ.PickRate > 1 {
			avoid = append(avoid, champ.ChampionName)
		}
	}
	return avoid
}

func (mga *MetaGameAnalytics) getTrendPicks(champions []ChampionMetrics) []string {
	var trend []string
	for _, champ := range champions {
		if champ.TrendDirection == "rising" && champ.TrendStrength > 0.3 {
			trend = append(trend, champ.ChampionName)
		}
	}
	return trend
}

func (mga *MetaGameAnalytics) calculateAverageWinRate(champions []ChampionMetrics) float64 {
	if len(champions) == 0 {
		return 0
	}
	
	var total float64
	for _, champ := range champions {
		total += champ.WinRate
	}
	return total / float64(len(champions))
}

func (mga *MetaGameAnalytics) calculateMetaDiversity(champions []ChampionMetrics) float64 {
	// Higher diversity = more champions with reasonable pick rates
	viableChampions := 0
	for _, champ := range champions {
		if champ.PickRate >= 1.0 { // 1% pick rate threshold
			viableChampions++
		}
	}
	return float64(viableChampions) / float64(len(champions)) * 100
}

func (mga *MetaGameAnalytics) calculateBalanceScore(champions []ChampionMetrics) float64 {
	// Balance score: lower variance in win rates = more balanced
	if len(champions) == 0 {
		return 0
	}
	
	avg := mga.calculateAverageWinRate(champions)
	var variance float64
	
	for _, champ := range champions {
		variance += math.Pow(champ.WinRate-avg, 2)
	}
	variance /= float64(len(champions))
	
	// Convert to 0-100 score (lower variance = higher score)
	standardDev := math.Sqrt(variance)
	return math.Max(0, 100-(standardDev*2))
}

// Global meta-game analytics instance
var metaGameAnalytics *MetaGameAnalytics

// InitializeMetaGameAnalytics initializes the meta-game analytics system
func InitializeMetaGameAnalytics() {
	if database != nil {
		metaGameAnalytics = NewMetaGameAnalytics(database)
		log.Println("🔍 Meta-game analytics system initialized")
	}
}