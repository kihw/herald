package export

import (
	"context"
	"fmt"
	"time"

	"github.com/herald-lol/herald/backend/internal/analytics"
	"github.com/herald-lol/herald/backend/internal/match"
	"github.com/herald-lol/herald/backend/internal/riot"
	"github.com/herald-lol/herald/backend/internal/summoner"
)

// Herald.lol Gaming Analytics - Data Export Service
// Multi-format data export service for gaming analytics data

// ExportService handles exporting gaming data in various formats
type ExportService struct {
	config          *ExportConfig
	analyticsEngine *analytics.AnalyticsEngine
	matchAnalyzer   *match.MatchAnalyzer
	summonerService *summoner.SummonerService

	// Export processors
	csvProcessor   *CSVProcessor
	jsonProcessor  *JSONProcessor
	xlsxProcessor  *XLSXProcessor
	pdfProcessor   *PDFProcessor
	chartProcessor *ChartProcessor

	// Cache and storage
	exportCache        map[string]*CachedExport
	compressionEnabled bool
	encryptionEnabled  bool
}

// NewExportService creates a new export service
func NewExportService(
	config *ExportConfig,
	analyticsEngine *analytics.AnalyticsEngine,
	matchAnalyzer *match.MatchAnalyzer,
	summonerService *summoner.SummonerService,
) *ExportService {
	service := &ExportService{
		config:             config,
		analyticsEngine:    analyticsEngine,
		matchAnalyzer:      matchAnalyzer,
		summonerService:    summonerService,
		exportCache:        make(map[string]*CachedExport),
		compressionEnabled: config.EnableCompression,
		encryptionEnabled:  config.EnableEncryption,
	}

	// Initialize processors
	service.csvProcessor = NewCSVProcessor(config.CSV)
	service.jsonProcessor = NewJSONProcessor(config.JSON)
	service.xlsxProcessor = NewXLSXProcessor(config.XLSX)
	service.pdfProcessor = NewPDFProcessor(config.PDF)
	service.chartProcessor = NewChartProcessor(config.Charts)

	return service
}

// ExportPlayerAnalytics exports comprehensive player analytics data
func (s *ExportService) ExportPlayerAnalytics(ctx context.Context, request *PlayerExportRequest) (*ExportResult, error) {
	// Validate request
	if err := s.validatePlayerExportRequest(request); err != nil {
		return nil, fmt.Errorf("invalid export request: %w", err)
	}

	// Check cache first
	cacheKey := s.generateCacheKey("player", request.PlayerPUUID, request.Format, request.TimeRange)
	if cached, exists := s.exportCache[cacheKey]; exists && !s.isCacheExpired(cached) {
		return &ExportResult{
			ExportID:    cached.ExportID,
			Format:      cached.Format,
			FileSize:    cached.FileSize,
			Status:      "completed",
			DownloadURL: cached.DownloadURL,
			CreatedAt:   cached.CreatedAt,
			ExpiresAt:   cached.ExpiresAt,
		}, nil
	}

	// Generate export ID
	exportID := s.generateExportID()

	// Collect player data
	playerData, err := s.collectPlayerData(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to collect player data: %w", err)
	}

	// Export data in requested format
	var exportedData []byte
	var fileName string

	switch request.Format {
	case "csv":
		exportedData, fileName, err = s.csvProcessor.ExportPlayerData(playerData, request)
	case "json":
		exportedData, fileName, err = s.jsonProcessor.ExportPlayerData(playerData, request)
	case "xlsx":
		exportedData, fileName, err = s.xlsxProcessor.ExportPlayerData(playerData, request)
	case "pdf":
		exportedData, fileName, err = s.pdfProcessor.ExportPlayerData(playerData, request)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", request.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to export data: %w", err)
	}

	// Apply compression if enabled
	if s.compressionEnabled {
		exportedData, err = s.compressData(exportedData)
		if err != nil {
			return nil, fmt.Errorf("failed to compress data: %w", err)
		}
	}

	// Apply encryption if enabled
	if s.encryptionEnabled {
		exportedData, err = s.encryptData(exportedData)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt data: %w", err)
		}
	}

	// Store export
	downloadURL, err := s.storeExport(exportID, fileName, exportedData)
	if err != nil {
		return nil, fmt.Errorf("failed to store export: %w", err)
	}

	// Cache result
	result := &ExportResult{
		ExportID:    exportID,
		Format:      request.Format,
		FileSize:    len(exportedData),
		Status:      "completed",
		DownloadURL: downloadURL,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(s.config.ExportTTL),
		Metadata: &ExportMetadata{
			PlayerPUUID: request.PlayerPUUID,
			TimeRange:   request.TimeRange,
			DataPoints:  len(playerData.Matches),
			Compressed:  s.compressionEnabled,
			Encrypted:   s.encryptionEnabled,
		},
	}

	s.exportCache[cacheKey] = &CachedExport{
		ExportID:    exportID,
		Format:      request.Format,
		FileSize:    len(exportedData),
		DownloadURL: downloadURL,
		CreatedAt:   result.CreatedAt,
		ExpiresAt:   result.ExpiresAt,
	}

	return result, nil
}

// ExportMatchAnalytics exports match analysis data
func (s *ExportService) ExportMatchAnalytics(ctx context.Context, request *MatchExportRequest) (*ExportResult, error) {
	// Validate request
	if err := s.validateMatchExportRequest(request); err != nil {
		return nil, fmt.Errorf("invalid match export request: %w", err)
	}

	// Generate export ID
	exportID := s.generateExportID()

	// Collect match data
	matchData, err := s.collectMatchData(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to collect match data: %w", err)
	}

	// Export in requested format
	var exportedData []byte
	var fileName string

	switch request.Format {
	case "csv":
		exportedData, fileName, err = s.csvProcessor.ExportMatchData(matchData, request)
	case "json":
		exportedData, fileName, err = s.jsonProcessor.ExportMatchData(matchData, request)
	case "xlsx":
		exportedData, fileName, err = s.xlsxProcessor.ExportMatchData(matchData, request)
	case "pdf":
		exportedData, fileName, err = s.pdfProcessor.ExportMatchData(matchData, request)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", request.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to export match data: %w", err)
	}

	// Store export
	downloadURL, err := s.storeExport(exportID, fileName, exportedData)
	if err != nil {
		return nil, fmt.Errorf("failed to store export: %w", err)
	}

	return &ExportResult{
		ExportID:    exportID,
		Format:      request.Format,
		FileSize:    len(exportedData),
		Status:      "completed",
		DownloadURL: downloadURL,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(s.config.ExportTTL),
		Metadata: &ExportMetadata{
			MatchID:    request.MatchID,
			DataPoints: 1,
			Compressed: s.compressionEnabled,
			Encrypted:  s.encryptionEnabled,
		},
	}, nil
}

// ExportTeamAnalytics exports team performance analytics
func (s *ExportService) ExportTeamAnalytics(ctx context.Context, request *TeamExportRequest) (*ExportResult, error) {
	// Validate request
	if err := s.validateTeamExportRequest(request); err != nil {
		return nil, fmt.Errorf("invalid team export request: %w", err)
	}

	// Generate export ID
	exportID := s.generateExportID()

	// Collect team data
	teamData, err := s.collectTeamData(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to collect team data: %w", err)
	}

	// Export in requested format
	var exportedData []byte
	var fileName string

	switch request.Format {
	case "csv":
		exportedData, fileName, err = s.csvProcessor.ExportTeamData(teamData, request)
	case "json":
		exportedData, fileName, err = s.jsonProcessor.ExportTeamData(teamData, request)
	case "xlsx":
		exportedData, fileName, err = s.xlsxProcessor.ExportTeamData(teamData, request)
	case "pdf":
		exportedData, fileName, err = s.pdfProcessor.ExportTeamData(teamData, request)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", request.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to export team data: %w", err)
	}

	// Store export
	downloadURL, err := s.storeExport(exportID, fileName, exportedData)
	if err != nil {
		return nil, fmt.Errorf("failed to store export: %w", err)
	}

	return &ExportResult{
		ExportID:    exportID,
		Format:      request.Format,
		FileSize:    len(exportedData),
		Status:      "completed",
		DownloadURL: downloadURL,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(s.config.ExportTTL),
		Metadata: &ExportMetadata{
			TeamName:   request.TeamName,
			DataPoints: len(teamData.Players),
			Compressed: s.compressionEnabled,
			Encrypted:  s.encryptionEnabled,
		},
	}, nil
}

// ExportChampionAnalytics exports champion-specific performance data
func (s *ExportService) ExportChampionAnalytics(ctx context.Context, request *ChampionExportRequest) (*ExportResult, error) {
	// Validate request
	if err := s.validateChampionExportRequest(request); err != nil {
		return nil, fmt.Errorf("invalid champion export request: %w", err)
	}

	// Generate export ID
	exportID := s.generateExportID()

	// Collect champion data
	championData, err := s.collectChampionData(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to collect champion data: %w", err)
	}

	// Export in requested format with charts for champion data
	var exportedData []byte
	var fileName string

	switch request.Format {
	case "csv":
		exportedData, fileName, err = s.csvProcessor.ExportChampionData(championData, request)
	case "json":
		exportedData, fileName, err = s.jsonProcessor.ExportChampionData(championData, request)
	case "xlsx":
		exportedData, fileName, err = s.xlsxProcessor.ExportChampionData(championData, request)
	case "pdf":
		// PDF with charts for champion analytics
		exportedData, fileName, err = s.pdfProcessor.ExportChampionDataWithCharts(championData, request)
	case "charts":
		// Generate interactive charts
		exportedData, fileName, err = s.chartProcessor.ExportChampionCharts(championData, request)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", request.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to export champion data: %w", err)
	}

	// Store export
	downloadURL, err := s.storeExport(exportID, fileName, exportedData)
	if err != nil {
		return nil, fmt.Errorf("failed to store export: %w", err)
	}

	return &ExportResult{
		ExportID:    exportID,
		Format:      request.Format,
		FileSize:    len(exportedData),
		Status:      "completed",
		DownloadURL: downloadURL,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(s.config.ExportTTL),
		Metadata: &ExportMetadata{
			ChampionName: request.ChampionName,
			DataPoints:   len(championData.PerformanceHistory),
			Compressed:   s.compressionEnabled,
			Encrypted:    s.encryptionEnabled,
		},
	}, nil
}

// ExportCustomReport exports custom analytics reports
func (s *ExportService) ExportCustomReport(ctx context.Context, request *CustomReportRequest) (*ExportResult, error) {
	// Validate request
	if err := s.validateCustomReportRequest(request); err != nil {
		return nil, fmt.Errorf("invalid custom report request: %w", err)
	}

	// Generate export ID
	exportID := s.generateExportID()

	// Build custom report based on specifications
	reportData, err := s.buildCustomReport(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to build custom report: %w", err)
	}

	// Export in requested format
	var exportedData []byte
	var fileName string

	switch request.Format {
	case "csv":
		exportedData, fileName, err = s.csvProcessor.ExportCustomReport(reportData, request)
	case "json":
		exportedData, fileName, err = s.jsonProcessor.ExportCustomReport(reportData, request)
	case "xlsx":
		exportedData, fileName, err = s.xlsxProcessor.ExportCustomReport(reportData, request)
	case "pdf":
		exportedData, fileName, err = s.pdfProcessor.ExportCustomReport(reportData, request)
	case "charts":
		exportedData, fileName, err = s.chartProcessor.ExportCustomCharts(reportData, request)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", request.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to export custom report: %w", err)
	}

	// Store export
	downloadURL, err := s.storeExport(exportID, fileName, exportedData)
	if err != nil {
		return nil, fmt.Errorf("failed to store export: %w", err)
	}

	return &ExportResult{
		ExportID:    exportID,
		Format:      request.Format,
		FileSize:    len(exportedData),
		Status:      "completed",
		DownloadURL: downloadURL,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(s.config.ExportTTL),
		Metadata: &ExportMetadata{
			ReportName:  request.ReportName,
			DataPoints:  len(reportData.DataRows),
			Compressed:  s.compressionEnabled,
			Encrypted:   s.encryptionEnabled,
			CustomQuery: request.Query,
		},
	}, nil
}

// GetExportStatus returns the status of an export job
func (s *ExportService) GetExportStatus(ctx context.Context, exportID string) (*ExportStatus, error) {
	// Check cache first
	for _, cached := range s.exportCache {
		if cached.ExportID == exportID {
			status := "completed"
			if s.isCacheExpired(cached) {
				status = "expired"
			}

			return &ExportStatus{
				ExportID:    exportID,
				Status:      status,
				Progress:    100,
				FileSize:    cached.FileSize,
				DownloadURL: cached.DownloadURL,
				CreatedAt:   cached.CreatedAt,
				ExpiresAt:   cached.ExpiresAt,
			}, nil
		}
	}

	// Check persistent storage
	exportInfo, err := s.getStoredExportInfo(exportID)
	if err != nil {
		return nil, fmt.Errorf("export not found: %w", err)
	}

	return &ExportStatus{
		ExportID:     exportID,
		Status:       exportInfo.Status,
		Progress:     exportInfo.Progress,
		FileSize:     exportInfo.FileSize,
		DownloadURL:  exportInfo.DownloadURL,
		CreatedAt:    exportInfo.CreatedAt,
		ExpiresAt:    exportInfo.ExpiresAt,
		ErrorMessage: exportInfo.ErrorMessage,
	}, nil
}

// ListExports returns a list of exports for a user
func (s *ExportService) ListExports(ctx context.Context, userID string, limit int) ([]*ExportSummary, error) {
	exports, err := s.getUserExports(userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user exports: %w", err)
	}

	summaries := make([]*ExportSummary, len(exports))
	for i, export := range exports {
		summaries[i] = &ExportSummary{
			ExportID:    export.ExportID,
			Format:      export.Format,
			Status:      export.Status,
			FileSize:    export.FileSize,
			CreatedAt:   export.CreatedAt,
			ExpiresAt:   export.ExpiresAt,
			DataType:    export.DataType,
			Description: export.Description,
		}
	}

	return summaries, nil
}

// DeleteExport removes an export from storage
func (s *ExportService) DeleteExport(ctx context.Context, exportID string) error {
	// Remove from cache
	for key, cached := range s.exportCache {
		if cached.ExportID == exportID {
			delete(s.exportCache, key)
			break
		}
	}

	// Remove from persistent storage
	return s.deleteStoredExport(exportID)
}

// GetSupportedFormats returns the list of supported export formats
func (s *ExportService) GetSupportedFormats() []ExportFormat {
	return []ExportFormat{
		{
			Name:        "CSV",
			Key:         "csv",
			Description: "Comma-separated values for spreadsheet applications",
			Extensions:  []string{".csv"},
			MimeType:    "text/csv",
			Features:    []string{"lightweight", "universal", "data-only"},
		},
		{
			Name:        "JSON",
			Key:         "json",
			Description: "JavaScript Object Notation for API integration",
			Extensions:  []string{".json"},
			MimeType:    "application/json",
			Features:    []string{"structured", "api-friendly", "nested-data"},
		},
		{
			Name:        "Excel",
			Key:         "xlsx",
			Description: "Microsoft Excel format with formatting and charts",
			Extensions:  []string{".xlsx"},
			MimeType:    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			Features:    []string{"formatted", "charts", "multiple-sheets", "business-friendly"},
		},
		{
			Name:        "PDF Report",
			Key:         "pdf",
			Description: "Professional PDF reports with visualizations",
			Extensions:  []string{".pdf"},
			MimeType:    "application/pdf",
			Features:    []string{"professional", "charts", "formatted", "print-ready"},
		},
		{
			Name:        "Interactive Charts",
			Key:         "charts",
			Description: "Interactive HTML charts and visualizations",
			Extensions:  []string{".html"},
			MimeType:    "text/html",
			Features:    []string{"interactive", "visualizations", "web-based", "responsive"},
		},
	}
}

// Data collection helper methods
func (s *ExportService) collectPlayerData(ctx context.Context, request *PlayerExportRequest) (*PlayerExportData, error) {
	// Collect summoner analytics
	summonerAnalysis, err := s.summonerService.GetSummonerAnalysis(ctx, &summoner.SummonerAnalysisRequest{
		Region:       request.Region,
		SummonerName: request.SummonerName,
		PlayerPUUID:  request.PlayerPUUID,
		TimeRange:    request.TimeRange,
		GameModes:    request.GameModes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get summoner analysis: %w", err)
	}

	// Collect match history and analysis
	matches := []*MatchExportData{}
	for _, matchID := range request.MatchIDs {
		// Get match data and analysis
		matchAnalysis, err := s.matchAnalyzer.AnalyzeMatch(ctx, &match.MatchAnalysisRequest{
			PlayerPUUID:   request.PlayerPUUID,
			AnalysisDepth: "standard",
		})
		if err != nil {
			continue // Skip failed matches
		}

		matches = append(matches, &MatchExportData{
			MatchID:     matchID,
			Champion:    matchAnalysis.MatchInfo.Champion,
			Role:        matchAnalysis.MatchInfo.Role,
			Result:      matchAnalysis.MatchInfo.Result,
			Duration:    matchAnalysis.MatchInfo.Duration,
			Performance: matchAnalysis.Performance,
			KeyMoments:  matchAnalysis.KeyMoments,
		})
	}

	return &PlayerExportData{
		PlayerInfo: &PlayerInfo{
			PUUID:        request.PlayerPUUID,
			SummonerName: request.SummonerName,
			Region:       request.Region,
			Rank:         summonerAnalysis.CurrentRank.Tier,
			LP:           summonerAnalysis.CurrentRank.LeaguePoints,
		},
		Summary:    summonerAnalysis.PerformanceSummary,
		Matches:    matches,
		TimeRange:  request.TimeRange,
		ExportedAt: time.Now(),
		TotalGames: len(matches),
	}, nil
}

func (s *ExportService) collectMatchData(ctx context.Context, request *MatchExportRequest) (*MatchExportData, error) {
	// Get detailed match analysis
	matchAnalysis, err := s.matchAnalyzer.AnalyzeMatch(ctx, &match.MatchAnalysisRequest{
		PlayerPUUID:             request.PlayerPUUID,
		AnalysisDepth:           "detailed",
		IncludePhaseAnalysis:    true,
		IncludeKeyMoments:       true,
		IncludeTeamAnalysis:     true,
		IncludeOpponentAnalysis: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to analyze match: %w", err)
	}

	return &MatchExportData{
		MatchID:               request.MatchID,
		Champion:              matchAnalysis.MatchInfo.Champion,
		Role:                  matchAnalysis.MatchInfo.Role,
		Result:                matchAnalysis.MatchInfo.Result,
		Duration:              matchAnalysis.MatchInfo.Duration,
		Performance:           matchAnalysis.Performance,
		PhaseAnalysis:         matchAnalysis.PhaseAnalysis,
		KeyMoments:            matchAnalysis.KeyMoments,
		TeamAnalysis:          matchAnalysis.TeamAnalysis,
		Insights:              matchAnalysis.Insights,
		LearningOpportunities: matchAnalysis.LearningOpportunities,
		OverallRating:         matchAnalysis.OverallRating,
	}, nil
}

func (s *ExportService) collectTeamData(ctx context.Context, request *TeamExportRequest) (*TeamExportData, error) {
	// Collect data for each team member
	players := []*PlayerExportData{}

	for _, playerPUUID := range request.PlayerPUUIDs {
		playerData, err := s.collectPlayerData(ctx, &PlayerExportRequest{
			PlayerPUUID: playerPUUID,
			TimeRange:   request.TimeRange,
			GameModes:   request.GameModes,
			MatchIDs:    request.SharedMatchIDs,
		})
		if err != nil {
			continue // Skip failed players
		}

		players = append(players, playerData)
	}

	// Calculate team metrics
	teamMetrics := s.calculateTeamMetrics(players)

	return &TeamExportData{
		TeamName:    request.TeamName,
		Players:     players,
		TeamMetrics: teamMetrics,
		TimeRange:   request.TimeRange,
		ExportedAt:  time.Now(),
	}, nil
}

func (s *ExportService) collectChampionData(ctx context.Context, request *ChampionExportRequest) (*ChampionExportData, error) {
	// Get champion performance history
	championAnalysis, err := s.analyticsEngine.AnalyzeChampionPerformance(ctx, &analytics.ChampionAnalysisRequest{
		PlayerPUUID:  request.PlayerPUUID,
		ChampionName: request.ChampionName,
		TimeRange:    request.TimeRange,
		GameModes:    request.GameModes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to analyze champion: %w", err)
	}

	return &ChampionExportData{
		ChampionName:       request.ChampionName,
		PlayerPUUID:        request.PlayerPUUID,
		PerformanceHistory: championAnalysis.PerformanceHistory,
		Statistics:         championAnalysis.Statistics,
		Trends:             championAnalysis.Trends,
		Comparisons:        championAnalysis.Comparisons,
		Recommendations:    championAnalysis.Recommendations,
		TimeRange:          request.TimeRange,
		ExportedAt:         time.Now(),
	}, nil
}

func (s *ExportService) buildCustomReport(ctx context.Context, request *CustomReportRequest) (*CustomReportData, error) {
	// Execute custom query based on report specifications
	data := &CustomReportData{
		ReportName:   request.ReportName,
		Description:  request.Description,
		Parameters:   request.Parameters,
		Columns:      request.Columns,
		DataRows:     []map[string]interface{}{},
		Filters:      request.Filters,
		Aggregations: request.Aggregations,
		GeneratedAt:  time.Now(),
	}

	// This would typically execute against a data warehouse or analytics database
	// For now, we'll create a placeholder implementation
	switch request.ReportType {
	case "performance_trends":
		data.DataRows = s.generatePerformanceTrendsData(request)
	case "champion_comparison":
		data.DataRows = s.generateChampionComparisonData(request)
	case "rank_progression":
		data.DataRows = s.generateRankProgressionData(request)
	case "meta_analysis":
		data.DataRows = s.generateMetaAnalysisData(request)
	default:
		return nil, fmt.Errorf("unsupported report type: %s", request.ReportType)
	}

	return data, nil
}

// Helper methods continue in next part...
