package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// ExportFormat represents different export formats
type ExportFormat string

const (
	ExportFormatCSV     ExportFormat = "csv"
	ExportFormatJSON    ExportFormat = "json"
	ExportFormatParquet ExportFormat = "parquet"
	ExportFormatXLSX    ExportFormat = "xlsx"
)

// ExportRequest represents an export request
type ExportRequest struct {
	UserID          int          `json:"user_id"`
	Format          ExportFormat `json:"format"`
	IncludeMetadata bool         `json:"include_metadata"`
	DateRange       *DateRange   `json:"date_range,omitempty"`
	Filters         *ExportFilters `json:"filters,omitempty"`
	Compression     bool         `json:"compression"`
	RequestedBy     string       `json:"requested_by"`
	RequestedAt     time.Time    `json:"requested_at"`
}

// DateRange represents a date range filter
type DateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// ExportFilters represents export filtering options
type ExportFilters struct {
	ChampionNames []string `json:"champion_names,omitempty"`
	GameModes     []string `json:"game_modes,omitempty"`
	QueueIDs      []int    `json:"queue_ids,omitempty"`
	MinKDA        *float64 `json:"min_kda,omitempty"`
	MaxKDA        *float64 `json:"max_kda,omitempty"`
	WinsOnly      bool     `json:"wins_only"`
	LossesOnly    bool     `json:"losses_only"`
}

// ExportMetadata represents metadata included with exports
type ExportMetadata struct {
	ExportInfo    ExportInfo    `json:"export_info"`
	UserInfo      UserInfo      `json:"user_info"`
	DataSummary   DataSummary   `json:"data_summary"`
	SystemInfo    SystemInfo    `json:"system_info"`
}

// ExportInfo contains export-specific information
type ExportInfo struct {
	Format           string    `json:"format"`
	ExportedAt       time.Time `json:"exported_at"`
	RequestedBy      string    `json:"requested_by"`
	TotalRecords     int       `json:"total_records"`
	FiltersApplied   bool      `json:"filters_applied"`
	CompressionUsed  bool      `json:"compression_used"`
	ProcessingTimeMS int64     `json:"processing_time_ms"`
}

// UserInfo contains user information
type UserInfo struct {
	RiotID       string `json:"riot_id"`
	RiotTag      string `json:"riot_tag"`
	Region       string `json:"region"`
	LastSync     string `json:"last_sync"`
	TotalMatches int    `json:"total_matches"`
}

// DataSummary contains summary statistics
type DataSummary struct {
	DateRange         DateRange              `json:"date_range"`
	ChampionCounts    map[string]int         `json:"champion_counts"`
	GameModeCounts    map[string]int         `json:"game_mode_counts"`
	OverallStats      OverallStats           `json:"overall_stats"`
	PerformanceMetrics PerformanceMetrics     `json:"performance_metrics"`
}

// OverallStats contains overall statistics
type OverallStats struct {
	TotalMatches     int     `json:"total_matches"`
	WinRate          float64 `json:"win_rate"`
	AverageKDA       float64 `json:"average_kda"`
	AverageGameTime  float64 `json:"average_game_time_minutes"`
	FavoriteChampion string  `json:"favorite_champion"`
	LongestWinStreak int     `json:"longest_win_streak"`
}

// PerformanceMetrics contains performance analysis
type PerformanceMetrics struct {
	KillsPerGame    float64 `json:"kills_per_game"`
	DeathsPerGame   float64 `json:"deaths_per_game"`
	AssistsPerGame  float64 `json:"assists_per_game"`
	BestKDAGame     float64 `json:"best_kda_game"`
	WorstKDAGame    float64 `json:"worst_kda_game"`
	ConsistencyScore float64 `json:"consistency_score"`
}

// SystemInfo contains system information
type SystemInfo struct {
	ExporterVersion   string `json:"exporter_version"`
	DatabaseVersion   string `json:"database_version"`
	APIVersion        string `json:"api_version"`
	GeneratedBy       string `json:"generated_by"`
	ProcessingNode    string `json:"processing_node"`
}

// AdvancedExporter handles advanced export functionality
type AdvancedExporter struct {
	database *Database
}

// NewAdvancedExporter creates a new advanced exporter
func NewAdvancedExporter(db *Database) *AdvancedExporter {
	return &AdvancedExporter{
		database: db,
	}
}

// ExportMatches exports matches in the specified format
func (ae *AdvancedExporter) ExportMatches(request ExportRequest) ([]byte, string, error) {
	start := time.Now()
	log.Printf("📤 Starting export: format=%s, user=%d, filters=%v", 
		request.Format, request.UserID, request.Filters != nil)
	
	// Get filtered matches
	matches, err := ae.getFilteredMatches(request)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get filtered matches: %w", err)
	}
	
	log.Printf("📊 Retrieved %d matches for export", len(matches))
	
	// Generate metadata if requested
	var metadata *ExportMetadata
	if request.IncludeMetadata {
		metadata, err = ae.generateMetadata(request, matches, start)
		if err != nil {
			log.Printf("⚠️ Failed to generate metadata: %v", err)
		}
	}
	
	// Export based on format
	var data []byte
	var filename string
	
	switch request.Format {
	case ExportFormatJSON:
		data, filename, err = ae.exportJSON(matches, metadata, request)
	case ExportFormatCSV:
		data, filename, err = ae.exportCSV(matches, metadata, request)
	case ExportFormatParquet:
		data, filename, err = ae.exportParquet(matches, metadata, request)
	case ExportFormatXLSX:
		data, filename, err = ae.exportXLSX(matches, metadata, request)
	default:
		return nil, "", fmt.Errorf("unsupported export format: %s", request.Format)
	}
	
	if err != nil {
		return nil, "", err
	}
	
	processingTime := time.Since(start)
	log.Printf("✅ Export completed: %s (%d bytes) in %v", 
		filename, len(data), processingTime)
	
	return data, filename, nil
}

// getFilteredMatches retrieves matches based on filters
func (ae *AdvancedExporter) getFilteredMatches(request ExportRequest) ([]Match, error) {
	// Build dynamic query based on filters
	query := `SELECT id, user_id, game_creation, game_duration, game_mode, queue_id,
				     champion_name, champion_id, kills, deaths, assists, win, created_at
			  FROM matches WHERE user_id = ?`
	
	args := []interface{}{request.UserID}
	
	// Apply date range filter
	if request.DateRange != nil {
		query += " AND game_creation >= ? AND game_creation <= ?"
		args = append(args, request.DateRange.StartDate.Unix()*1000)
		args = append(args, request.DateRange.EndDate.Unix()*1000)
	}
	
	// Apply other filters
	if request.Filters != nil {
		if len(request.Filters.ChampionNames) > 0 {
			placeholders := make([]string, len(request.Filters.ChampionNames))
			for i, champion := range request.Filters.ChampionNames {
				placeholders[i] = "?"
				args = append(args, champion)
			}
			query += fmt.Sprintf(" AND champion_name IN (%s)", 
				fmt.Sprintf("%s", placeholders))
		}
		
		if len(request.Filters.GameModes) > 0 {
			placeholders := make([]string, len(request.Filters.GameModes))
			for i, mode := range request.Filters.GameModes {
				placeholders[i] = "?"
				args = append(args, mode)
			}
			query += fmt.Sprintf(" AND game_mode IN (%s)", 
				fmt.Sprintf("%s", placeholders))
		}
		
		if len(request.Filters.QueueIDs) > 0 {
			placeholders := make([]string, len(request.Filters.QueueIDs))
			for i, queueID := range request.Filters.QueueIDs {
				placeholders[i] = "?"
				args = append(args, queueID)
			}
			query += fmt.Sprintf(" AND queue_id IN (%s)", 
				fmt.Sprintf("%s", placeholders))
		}
		
		if request.Filters.WinsOnly {
			query += " AND win = 1"
		} else if request.Filters.LossesOnly {
			query += " AND win = 0"
		}
		
		// KDA filters (requires calculation)
		if request.Filters.MinKDA != nil || request.Filters.MaxKDA != nil {
			if request.Filters.MinKDA != nil {
				query += " AND (CASE WHEN deaths > 0 THEN (kills + assists) / CAST(deaths AS FLOAT) ELSE (kills + assists) END) >= ?"
				args = append(args, *request.Filters.MinKDA)
			}
			if request.Filters.MaxKDA != nil {
				query += " AND (CASE WHEN deaths > 0 THEN (kills + assists) / CAST(deaths AS FLOAT) ELSE (kills + assists) END) <= ?"
				args = append(args, *request.Filters.MaxKDA)
			}
		}
	}
	
	query += " ORDER BY game_creation DESC"
	
	// Execute query
	rows, err := ae.database.Query(query, args...)
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
			continue
		}
		matches = append(matches, match)
	}
	
	return matches, nil
}

// generateMetadata generates comprehensive metadata for the export
func (ae *AdvancedExporter) generateMetadata(request ExportRequest, matches []Match, startTime time.Time) (*ExportMetadata, error) {
	// Get user information
	userQuery := `SELECT riot_id, riot_tag, riot_puuid, region, last_sync FROM users WHERE id = ?`
	var userInfo UserInfo
	var riotPUUID string
	var lastSync *time.Time
	
	err := ae.database.QueryRow(userQuery, request.UserID).Scan(
		&userInfo.RiotID, &userInfo.RiotTag, &riotPUUID, &userInfo.Region, &lastSync)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	
	if lastSync != nil {
		userInfo.LastSync = lastSync.Format(time.RFC3339)
	}
	
	// Count total matches for user
	var totalMatches int
	ae.database.QueryRow("SELECT COUNT(*) FROM matches WHERE user_id = ?", request.UserID).Scan(&totalMatches)
	userInfo.TotalMatches = totalMatches
	
	// Calculate data summary
	dataSummary := ae.calculateDataSummary(matches)
	
	// Create metadata
	metadata := &ExportMetadata{
		ExportInfo: ExportInfo{
			Format:           string(request.Format),
			ExportedAt:       time.Now(),
			RequestedBy:      request.RequestedBy,
			TotalRecords:     len(matches),
			FiltersApplied:   request.Filters != nil,
			CompressionUsed:  request.Compression,
			ProcessingTimeMS: time.Since(startTime).Milliseconds(),
		},
		UserInfo:    userInfo,
		DataSummary: dataSummary,
		SystemInfo: SystemInfo{
			ExporterVersion: "2.0.0",
			DatabaseVersion: "SQLite 3.x",
			APIVersion:      "Riot API v5",
			GeneratedBy:     "LoL Match Exporter Advanced",
			ProcessingNode:  "localhost",
		},
	}
	
	return metadata, nil
}

// calculateDataSummary calculates summary statistics for the exported data
func (ae *AdvancedExporter) calculateDataSummary(matches []Match) DataSummary {
	if len(matches) == 0 {
		return DataSummary{}
	}
	
	// Initialize counters
	championCounts := make(map[string]int)
	gameModeCounts := make(map[string]int)
	
	var totalKills, totalDeaths, totalAssists int
	var totalGameTime int64
	var wins int
	var earliestGame, latestGame int64 = matches[0].GameCreation, matches[0].GameCreation
	var favoriteChampion string
	var maxChampionCount int
	
	// Calculate statistics
	for _, match := range matches {
		// Champion counts
		championCounts[match.ChampionName]++
		if championCounts[match.ChampionName] > maxChampionCount {
			maxChampionCount = championCounts[match.ChampionName]
			favoriteChampion = match.ChampionName
		}
		
		// Game mode counts
		gameModeCounts[match.GameMode]++
		
		// Performance stats
		totalKills += match.Kills
		totalDeaths += match.Deaths
		totalAssists += match.Assists
		totalGameTime += int64(match.GameDuration)
		
		if match.Win {
			wins++
		}
		
		// Date range
		if match.GameCreation < earliestGame {
			earliestGame = match.GameCreation
		}
		if match.GameCreation > latestGame {
			latestGame = match.GameCreation
		}
	}
	
	// Calculate derived metrics
	totalMatches := len(matches)
	winRate := float64(wins) / float64(totalMatches) * 100
	averageKDA := float64(totalKills+totalAssists) / float64(totalDeaths)
	if totalDeaths == 0 {
		averageKDA = float64(totalKills + totalAssists)
	}
	averageGameTime := float64(totalGameTime) / float64(totalMatches) / 60.0 // Convert to minutes
	
	killsPerGame := float64(totalKills) / float64(totalMatches)
	deathsPerGame := float64(totalDeaths) / float64(totalMatches)
	assistsPerGame := float64(totalAssists) / float64(totalMatches)
	
	// Calculate consistency score (simplified)
	consistencyScore := 50.0 // Placeholder calculation
	if averageKDA > 0 {
		consistencyScore = math.Min(100, averageKDA*20)
	}
	
	return DataSummary{
		DateRange: DateRange{
			StartDate: time.Unix(earliestGame/1000, 0),
			EndDate:   time.Unix(latestGame/1000, 0),
		},
		ChampionCounts: championCounts,
		GameModeCounts: gameModeCounts,
		OverallStats: OverallStats{
			TotalMatches:     totalMatches,
			WinRate:          winRate,
			AverageKDA:       averageKDA,
			AverageGameTime:  averageGameTime,
			FavoriteChampion: favoriteChampion,
			LongestWinStreak: 0, // Would need additional calculation
		},
		PerformanceMetrics: PerformanceMetrics{
			KillsPerGame:     killsPerGame,
			DeathsPerGame:    deathsPerGame,
			AssistsPerGame:   assistsPerGame,
			BestKDAGame:      0,  // Would need additional calculation
			WorstKDAGame:     0,  // Would need additional calculation
			ConsistencyScore: consistencyScore,
		},
	}
}

// exportJSON exports matches to JSON format
func (ae *AdvancedExporter) exportJSON(matches []Match, metadata *ExportMetadata, request ExportRequest) ([]byte, string, error) {
	exportData := map[string]interface{}{
		"matches": matches,
	}
	
	if metadata != nil {
		exportData["metadata"] = metadata
	}
	
	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("lol_matches_%d_%s.json", request.UserID, timestamp)
	
	return data, filename, nil
}

// exportCSV exports matches to CSV format (reusing existing CSV export)
func (ae *AdvancedExporter) exportCSV(matches []Match, metadata *ExportMetadata, request ExportRequest) ([]byte, string, error) {
	// For now, use simplified CSV export
	// In a real implementation, you'd create a proper CSV with metadata header
	
	var csvData [][]string
	
	// Add metadata as comments if included
	if metadata != nil {
		csvData = append(csvData, []string{
			fmt.Sprintf("# Export generated at: %s", metadata.ExportInfo.ExportedAt.Format(time.RFC3339)),
			fmt.Sprintf("# User: %s#%s", metadata.UserInfo.RiotID, metadata.UserInfo.RiotTag),
			fmt.Sprintf("# Total records: %d", metadata.ExportInfo.TotalRecords),
			fmt.Sprintf("# Format: %s", metadata.ExportInfo.Format),
		})
	}
	
	// CSV headers
	headers := []string{
		"match_id", "game_creation", "game_duration_minutes", "game_mode", "queue_id",
		"champion_name", "champion_id", "kills", "deaths", "assists", "kda", "win", "created_at",
	}
	csvData = append(csvData, headers)
	
	// Add match data
	for _, match := range matches {
		kda := float64(match.Kills + match.Assists)
		if match.Deaths > 0 {
			kda = float64(match.Kills+match.Assists) / float64(match.Deaths)
		}
		
		row := []string{
			match.ID,
			fmt.Sprintf("%d", match.GameCreation),
			fmt.Sprintf("%.1f", float64(match.GameDuration)/60),
			match.GameMode,
			fmt.Sprintf("%d", match.QueueID),
			match.ChampionName,
			fmt.Sprintf("%d", match.ChampionID),
			fmt.Sprintf("%d", match.Kills),
			fmt.Sprintf("%d", match.Deaths),
			fmt.Sprintf("%d", match.Assists),
			fmt.Sprintf("%.2f", kda),
			fmt.Sprintf("%t", match.Win),
			match.CreatedAt.Format(time.RFC3339),
		}
		csvData = append(csvData, row)
	}
	
	// Convert to CSV bytes
	var csvBytes []byte
	for _, row := range csvData {
		line := ""
		for i, field := range row {
			if i > 0 {
				line += ","
			}
			line += fmt.Sprintf("\"%s\"", field)
		}
		line += "\n"
		csvBytes = append(csvBytes, []byte(line)...)
	}
	
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("lol_matches_%d_%s.csv", request.UserID, timestamp)
	
	return csvBytes, filename, nil
}

// exportParquet exports matches to Parquet format (placeholder)
func (ae *AdvancedExporter) exportParquet(matches []Match, metadata *ExportMetadata, request ExportRequest) ([]byte, string, error) {
	// For now, return JSON as placeholder for Parquet
	// In a real implementation, you'd use a library like github.com/xitongsys/parquet-go
	log.Println("⚠️ Parquet export not yet implemented, returning JSON")
	return ae.exportJSON(matches, metadata, request)
}

// exportXLSX exports matches to Excel format (placeholder)
func (ae *AdvancedExporter) exportXLSX(matches []Match, metadata *ExportMetadata, request ExportRequest) ([]byte, string, error) {
	// For now, return CSV as placeholder for XLSX
	// In a real implementation, you'd use a library like github.com/360EntSecGroup-Skylar/excelize
	log.Println("⚠️ XLSX export not yet implemented, returning CSV")
	return ae.exportCSV(matches, metadata, request)
}

// Global advanced exporter instance
var advancedExporter *AdvancedExporter

// InitializeAdvancedExporter initializes the advanced export system
func InitializeAdvancedExporter() {
	if database != nil {
		advancedExporter = NewAdvancedExporter(database)
		log.Println("📤 Advanced export system initialized")
	}
}