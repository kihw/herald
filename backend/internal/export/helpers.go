package export

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Herald.lol Gaming Analytics - Export Service Helper Methods
// Helper functions for data export operations

// Validation helper methods

func (s *ExportService) validatePlayerExportRequest(request *PlayerExportRequest) error {
	if request.PlayerPUUID == "" {
		return fmt.Errorf("player PUUID is required")
	}
	
	if request.Region == "" {
		return fmt.Errorf("region is required")
	}
	
	if !s.isValidFormat(request.Format) {
		return fmt.Errorf("unsupported format: %s", request.Format)
	}
	
	if request.TimeRange == "" {
		return fmt.Errorf("time range is required")
	}
	
	// Validate subscription limits
	if err := s.validateSubscriptionLimits(request.PlayerPUUID, request.Format); err != nil {
		return fmt.Errorf("subscription limit exceeded: %w", err)
	}
	
	return nil
}

func (s *ExportService) validateMatchExportRequest(request *MatchExportRequest) error {
	if request.MatchID == "" {
		return fmt.Errorf("match ID is required")
	}
	
	if request.PlayerPUUID == "" {
		return fmt.Errorf("player PUUID is required")
	}
	
	if !s.isValidFormat(request.Format) {
		return fmt.Errorf("unsupported format: %s", request.Format)
	}
	
	return nil
}

func (s *ExportService) validateTeamExportRequest(request *TeamExportRequest) error {
	if request.TeamName == "" {
		return fmt.Errorf("team name is required")
	}
	
	if len(request.PlayerPUUIDs) < 2 || len(request.PlayerPUUIDs) > 10 {
		return fmt.Errorf("team must have between 2 and 10 players")
	}
	
	if !s.isValidFormat(request.Format) {
		return fmt.Errorf("unsupported format: %s", request.Format)
	}
	
	if request.TimeRange == "" {
		return fmt.Errorf("time range is required")
	}
	
	return nil
}

func (s *ExportService) validateChampionExportRequest(request *ChampionExportRequest) error {
	if request.PlayerPUUID == "" {
		return fmt.Errorf("player PUUID is required")
	}
	
	if request.ChampionName == "" {
		return fmt.Errorf("champion name is required")
	}
	
	if !s.isValidFormat(request.Format) {
		return fmt.Errorf("unsupported format: %s", request.Format)
	}
	
	if request.TimeRange == "" {
		return fmt.Errorf("time range is required")
	}
	
	return nil
}

func (s *ExportService) validateCustomReportRequest(request *CustomReportRequest) error {
	if request.ReportName == "" {
		return fmt.Errorf("report name is required")
	}
	
	if request.ReportType == "" {
		return fmt.Errorf("report type is required")
	}
	
	if !s.isValidFormat(request.Format) {
		return fmt.Errorf("unsupported format: %s", request.Format)
	}
	
	validReportTypes := []string{"performance_trends", "champion_comparison", "rank_progression", "meta_analysis"}
	if !s.isValidReportType(request.ReportType, validReportTypes) {
		return fmt.Errorf("unsupported report type: %s", request.ReportType)
	}
	
	return nil
}

// Utility helper methods

func (s *ExportService) isValidFormat(format string) bool {
	validFormats := []string{"csv", "json", "xlsx", "pdf", "charts"}
	for _, validFormat := range validFormats {
		if format == validFormat {
			return true
		}
	}
	return false
}

func (s *ExportService) isValidReportType(reportType string, validTypes []string) bool {
	for _, validType := range validTypes {
		if reportType == validType {
			return true
		}
	}
	return false
}

func (s *ExportService) validateSubscriptionLimits(playerPUUID, format string) error {
	// This would check against user subscription and usage limits
	// For now, return nil (no limits enforced)
	return nil
}

// Cache helper methods

func (s *ExportService) generateCacheKey(dataType, identifier, format, timeRange string) string {
	return fmt.Sprintf("%s:%s:%s:%s", dataType, identifier, format, timeRange)
}

func (s *ExportService) isCacheExpired(cached *CachedExport) bool {
	return time.Now().After(cached.ExpiresAt)
}

// ID generation

func (s *ExportService) generateExportID() string {
	return uuid.New().String()
}

// Compression and encryption

func (s *ExportService) compressData(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}
	
	var compressed strings.Builder
	writer := gzip.NewWriter(&compressed)
	
	_, err := writer.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write data for compression: %w", err)
	}
	
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close compression writer: %w", err)
	}
	
	return []byte(compressed.String()), nil
}

func (s *ExportService) encryptData(data []byte) ([]byte, error) {
	// Create a new AES cipher block
	key := make([]byte, 32) // 256-bit key
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}
	
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}
	
	// Create GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}
	
	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	
	// Encrypt data
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Storage helper methods

func (s *ExportService) storeExport(exportID, fileName string, data []byte) (string, error) {
	// In a real implementation, this would store to S3, GCS, or similar
	// For now, return a mock download URL
	downloadURL := fmt.Sprintf("%s/%s/%s", s.config.CDNBaseURL, exportID, fileName)
	return downloadURL, nil
}

func (s *ExportService) getStoredExportInfo(exportID string) (*StoredExportInfo, error) {
	// In a real implementation, this would query storage metadata
	// For now, return a mock response
	return &StoredExportInfo{
		ExportID:    exportID,
		Status:      "completed",
		Progress:    100,
		FileSize:    1024,
		DownloadURL: fmt.Sprintf("%s/%s", s.config.CDNBaseURL, exportID),
		CreatedAt:   time.Now().Add(-1 * time.Hour),
		ExpiresAt:   time.Now().Add(s.config.ExportTTL),
	}, nil
}

func (s *ExportService) getUserExports(userID string, limit int) ([]*UserExport, error) {
	// In a real implementation, this would query user export history
	// For now, return empty slice
	return []*UserExport{}, nil
}

func (s *ExportService) deleteStoredExport(exportID string) error {
	// In a real implementation, this would delete from storage
	return nil
}

// Team metrics calculation

func (s *ExportService) calculateTeamMetrics(players []*PlayerExportData) *TeamMetrics {
	if len(players) == 0 {
		return &TeamMetrics{}
	}
	
	totalGames := 0
	totalWins := 0
	totalKDA := 0.0
	totalDuration := 0
	
	for _, player := range players {
		totalGames += player.TotalGames
		
		// Calculate wins and KDA from matches
		for _, match := range player.Matches {
			if match.Result == "Victory" {
				totalWins++
			}
			if match.Performance != nil {
				totalKDA += match.Performance.KDA
			}
			totalDuration += match.Duration
		}
	}
	
	avgKDA := 0.0
	avgDuration := 0
	winRate := 0.0
	
	if totalGames > 0 {
		avgKDA = totalKDA / float64(totalGames)
		avgDuration = totalDuration / totalGames
		winRate = float64(totalWins) / float64(totalGames)
	}
	
	return &TeamMetrics{
		TeamWinRate:         winRate,
		AverageTeamKDA:      avgKDA,
		AverageGameDuration: avgDuration,
		ObjectiveControl:    75.0, // Placeholder
		TeamFightRating:     80.0, // Placeholder
		MacroPlay:          70.0, // Placeholder
		PlayerSynergy:      make(map[string]float64),
	}
}

// Custom report data generation

func (s *ExportService) generatePerformanceTrendsData(request *CustomReportRequest) []map[string]interface{} {
	// Generate mock performance trends data
	data := []map[string]interface{}{}
	
	// Sample data points
	dates := []string{"2024-01-01", "2024-01-02", "2024-01-03", "2024-01-04", "2024-01-05"}
	
	for i, date := range dates {
		row := map[string]interface{}{
			"date":        date,
			"win_rate":    0.6 + float64(i)*0.02,
			"kda":         2.1 + float64(i)*0.1,
			"cs_per_min":  6.8 + float64(i)*0.05,
			"vision":      18 + i*2,
			"damage":      25000 + i*1000,
			"rating":      75.0 + float64(i)*1.5,
		}
		data = append(data, row)
	}
	
	return data
}

func (s *ExportService) generateChampionComparisonData(request *CustomReportRequest) []map[string]interface{} {
	// Generate mock champion comparison data
	data := []map[string]interface{}{}
	
	champions := []string{"Jinx", "Caitlyn", "Kai'Sa", "Ezreal", "Vayne"}
	
	for i, champion := range champions {
		row := map[string]interface{}{
			"champion":    champion,
			"games":       20 + i*5,
			"win_rate":    0.55 + float64(i)*0.03,
			"kda":         2.0 + float64(i)*0.15,
			"cs_per_min":  7.2 + float64(i)*0.1,
			"damage":      28000 + i*2000,
			"rating":      72.0 + float64(i)*2.0,
		}
		data = append(data, row)
	}
	
	return data
}

func (s *ExportService) generateRankProgressionData(request *CustomReportRequest) []map[string]interface{} {
	// Generate mock rank progression data
	data := []map[string]interface{}{}
	
	ranks := []string{"Silver 3", "Silver 2", "Silver 1", "Gold 4", "Gold 3"}
	
	for i, rank := range ranks {
		row := map[string]interface{}{
			"date":     fmt.Sprintf("2024-01-%02d", i+1),
			"rank":     rank,
			"lp":       50 + i*20,
			"games":    i*10 + 25,
			"win_rate": 0.58 + float64(i)*0.02,
			"kda":      2.1 + float64(i)*0.05,
		}
		data = append(data, row)
	}
	
	return data
}

func (s *ExportService) generateMetaAnalysisData(request *CustomReportRequest) []map[string]interface{} {
	// Generate mock meta analysis data
	data := []map[string]interface{}{}
	
	champions := []map[string]interface{}{
		{"champion": "Jinx", "pick_rate": 0.15, "win_rate": 0.52, "ban_rate": 0.08, "tier": "S"},
		{"champion": "Caitlyn", "pick_rate": 0.12, "win_rate": 0.51, "ban_rate": 0.06, "tier": "A"},
		{"champion": "Kai'Sa", "pick_rate": 0.10, "win_rate": 0.53, "ban_rate": 0.12, "tier": "S"},
		{"champion": "Ezreal", "pick_rate": 0.18, "win_rate": 0.48, "ban_rate": 0.03, "tier": "B"},
		{"champion": "Vayne", "pick_rate": 0.08, "win_rate": 0.54, "ban_rate": 0.15, "tier": "A"},
	}
	
	return champions
}

// Gaming-specific helper methods for Herald.lol

func (s *ExportService) calculateGamingMetrics(matches []*MatchExportData) *GamingMetrics {
	if len(matches) == 0 {
		return &GamingMetrics{}
	}
	
	totalKills := 0
	totalDeaths := 0
	totalAssists := 0
	totalCS := 0
	totalDamage := 0
	totalVision := 0
	totalDuration := 0
	wins := 0
	
	for _, match := range matches {
		if match.Performance != nil {
			totalKills += match.Performance.Kills
			totalDeaths += match.Performance.Deaths
			totalAssists += match.Performance.Assists
			totalCS += int(match.Performance.CSPerMinute * float64(match.Duration) / 60)
			totalDamage += match.Performance.TotalDamage
			totalVision += match.Performance.VisionScore
		}
		totalDuration += match.Duration
		if match.Result == "Victory" {
			wins++
		}
	}
	
	gameCount := len(matches)
	avgGameDuration := totalDuration / gameCount
	
	return &GamingMetrics{
		GamesPlayed:         gameCount,
		WinRate:            float64(wins) / float64(gameCount),
		AverageKDA:         s.calculateKDA(totalKills, totalDeaths, totalAssists),
		AverageKills:       float64(totalKills) / float64(gameCount),
		AverageDeaths:      float64(totalDeaths) / float64(gameCount),
		AverageAssists:     float64(totalAssists) / float64(gameCount),
		AverageCSPerMin:    float64(totalCS) / float64(totalDuration) * 60,
		AverageDamage:      totalDamage / gameCount,
		AverageVision:      totalVision / gameCount,
		AverageGameLength: avgGameDuration,
	}
}

func (s *ExportService) calculateKDA(kills, deaths, assists int) float64 {
	if deaths == 0 {
		return float64(kills + assists)
	}
	return float64(kills+assists) / float64(deaths)
}

func (s *ExportService) identifyPlaystyle(metrics *GamingMetrics) string {
	if metrics.AverageKills > 8 && metrics.AverageDeaths < 5 {
		return "Aggressive Carry"
	} else if metrics.AverageAssists > 10 && metrics.AverageVision > 20 {
		return "Supportive Team Player"
	} else if metrics.AverageCSPerMin > 7.5 && metrics.AverageDamage > 25000 {
		return "Farming Carry"
	} else if metrics.AverageDeaths < 4 && metrics.WinRate > 0.6 {
		return "Consistent Performer"
	} else if metrics.AverageKills > 6 && metrics.AverageDeaths > 7 {
		return "High Risk High Reward"
	}
	return "Balanced Player"
}

func (s *ExportService) calculatePerformanceRating(metrics *GamingMetrics) float64 {
	// Weighted performance calculation
	kdaScore := s.normalizeKDA(metrics.AverageKDA) * 0.3
	winRateScore := metrics.WinRate * 100 * 0.25
	csScore := s.normalizeCS(metrics.AverageCSPerMin) * 0.2
	damageScore := s.normalizeDamage(metrics.AverageDamage) * 0.15
	visionScore := s.normalizeVision(metrics.AverageVision) * 0.1
	
	return kdaScore + winRateScore + csScore + damageScore + visionScore
}

func (s *ExportService) normalizeKDA(kda float64) float64 {
	// Normalize KDA to 0-100 scale
	if kda >= 4.0 {
		return 100
	} else if kda >= 3.0 {
		return 80 + (kda-3.0)*20
	} else if kda >= 2.0 {
		return 60 + (kda-2.0)*20
	} else if kda >= 1.0 {
		return 30 + (kda-1.0)*30
	}
	return kda * 30
}

func (s *ExportService) normalizeCS(csPerMin float64) float64 {
	// Normalize CS/min to 0-100 scale
	if csPerMin >= 8.0 {
		return 100
	} else if csPerMin >= 6.0 {
		return 60 + (csPerMin-6.0)*20
	} else if csPerMin >= 4.0 {
		return 30 + (csPerMin-4.0)*15
	}
	return csPerMin * 7.5
}

func (s *ExportService) normalizeDamage(damage int) float64 {
	// Normalize damage to 0-100 scale
	damageFloat := float64(damage)
	if damageFloat >= 35000 {
		return 100
	} else if damageFloat >= 25000 {
		return 70 + (damageFloat-25000)*3/1000
	} else if damageFloat >= 15000 {
		return 40 + (damageFloat-15000)*3/1000
	}
	return damageFloat * 40 / 15000
}

func (s *ExportService) normalizeVision(vision int) float64 {
	// Normalize vision score to 0-100 scale
	visionFloat := float64(vision)
	if visionFloat >= 30 {
		return 100
	} else if visionFloat >= 20 {
		return 70 + (visionFloat-20)*3
	} else if visionFloat >= 10 {
		return 35 + (visionFloat-10)*3.5
	}
	return visionFloat * 3.5
}

// GamingMetrics contains gaming-specific performance metrics
type GamingMetrics struct {
	GamesPlayed         int     `json:"games_played"`
	WinRate            float64 `json:"win_rate"`
	AverageKDA         float64 `json:"average_kda"`
	AverageKills       float64 `json:"average_kills"`
	AverageDeaths      float64 `json:"average_deaths"`
	AverageAssists     float64 `json:"average_assists"`
	AverageCSPerMin    float64 `json:"average_cs_per_min"`
	AverageDamage      int     `json:"average_damage"`
	AverageVision      int     `json:"average_vision"`
	AverageGameLength  int     `json:"average_game_length"`
}