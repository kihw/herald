package export

import "time"

// Herald.lol Gaming Analytics - Export Service Configuration
// Configuration settings for multi-format data export

// ExportConfig contains main export service configuration
type ExportConfig struct {
	// General settings
	EnableCompression bool          `json:"enable_compression"`
	EnableEncryption  bool          `json:"enable_encryption"`
	ExportTTL         time.Duration `json:"export_ttl"`
	MaxFileSize       int64         `json:"max_file_size"`       // In bytes
	MaxConcurrentJobs int           `json:"max_concurrent_jobs"`
	CleanupInterval   time.Duration `json:"cleanup_interval"`
	
	// Storage settings
	StoragePath       string        `json:"storage_path"`
	CDNBaseURL        string        `json:"cdn_base_url"`
	SignedURLTTL      time.Duration `json:"signed_url_ttl"`
	
	// Format-specific configurations
	CSV    *CSVConfig    `json:"csv"`
	JSON   *JSONConfig   `json:"json"`
	XLSX   *XLSXConfig   `json:"xlsx"`
	PDF    *PDFConfig    `json:"pdf"`
	Charts *ChartsConfig `json:"charts"`
	
	// Performance settings
	PerformanceTargets *ExportPerformanceTargets `json:"performance_targets"`
	
	// Security settings
	SecuritySettings *ExportSecuritySettings `json:"security_settings"`
}

// Format-specific configurations

// CSVConfig contains CSV export configuration
type CSVConfig struct {
	DefaultDelimiter    string `json:"default_delimiter"`
	MaxRows             int    `json:"max_rows"`
	IncludeHeadersDefault bool `json:"include_headers_default"`
	DateFormatDefault   string `json:"date_format_default"`
	EncodingDefault     string `json:"encoding_default"`
	
	// CSV-specific limits
	MaxColumnWidth      int    `json:"max_column_width"`
	EscapeSpecialChars  bool   `json:"escape_special_chars"`
}

// JSONConfig contains JSON export configuration
type JSONConfig struct {
	PrettyPrintDefault bool   `json:"pretty_print_default"`
	MaxDepth          int    `json:"max_depth"`
	DateFormatDefault string `json:"date_format_default"`
	IncludeNulls      bool   `json:"include_nulls"`
	
	// JSON-specific settings
	CompressArrays    bool   `json:"compress_arrays"`
	StreamLargeData   bool   `json:"stream_large_data"`
}

// XLSXConfig contains Excel export configuration
type XLSXConfig struct {
	MaxRows            int    `json:"max_rows"`
	MaxColumns         int    `json:"max_columns"`
	MaxSheets          int    `json:"max_sheets"`
	DefaultSheetName   string `json:"default_sheet_name"`
	AutoFitColumnsDefault bool `json:"auto_fit_columns_default"`
	
	// Excel-specific formatting
	HeaderStyle        *XLSXStyle `json:"header_style"`
	DataStyle          *XLSXStyle `json:"data_style"`
	ChartStyle         *XLSXStyle `json:"chart_style"`
	
	// Chart settings
	EnableCharts       bool     `json:"enable_charts"`
	DefaultChartTypes  []string `json:"default_chart_types"`
	ChartDimensions    *ChartDimensions `json:"chart_dimensions"`
}

// PDFConfig contains PDF export configuration
type PDFConfig struct {
	DefaultPageSize     string      `json:"default_page_size"`
	DefaultOrientation  string      `json:"default_orientation"`
	DefaultFontFamily   string      `json:"default_font_family"`
	DefaultFontSize     int         `json:"default_font_size"`
	DefaultMargins      *PDFMargins `json:"default_margins"`
	MaxPages            int         `json:"max_pages"`
	
	// PDF-specific settings
	EnableWatermark     bool        `json:"enable_watermark"`
	WatermarkText       string      `json:"watermark_text"`
	WatermarkOpacity    float64     `json:"watermark_opacity"`
	
	// Brand settings
	LogoPath            string      `json:"logo_path"`
	BrandColors         *BrandColors `json:"brand_colors"`
	
	// Chart settings
	ChartDPI            int         `json:"chart_dpi"`
	EmbedCharts         bool        `json:"embed_charts"`
}

// ChartsConfig contains chart export configuration
type ChartsConfig struct {
	DefaultWidth        int      `json:"default_width"`
	DefaultHeight       int      `json:"default_height"`
	DefaultTheme        string   `json:"default_theme"`
	DefaultColorScheme  []string `json:"default_color_scheme"`
	MaxDataPoints       int      `json:"max_data_points"`
	
	// Chart libraries
	EnableD3            bool     `json:"enable_d3"`
	EnableChartJS       bool     `json:"enable_chartjs"`
	EnablePlotly        bool     `json:"enable_plotly"`
	
	// Interactive settings
	EnableInteractivity bool     `json:"enable_interactivity"`
	EnableZoom          bool     `json:"enable_zoom"`
	EnableTooltips      bool     `json:"enable_tooltips"`
	
	// Export formats
	SupportedFormats    []string `json:"supported_formats"` // PNG, SVG, HTML, etc.
}

// Supporting configuration structures

// XLSXStyle contains Excel styling configuration
type XLSXStyle struct {
	FontName      string  `json:"font_name"`
	FontSize      int     `json:"font_size"`
	FontBold      bool    `json:"font_bold"`
	FontItalic    bool    `json:"font_italic"`
	FontColor     string  `json:"font_color"`
	BackgroundColor string `json:"background_color"`
	BorderStyle   string  `json:"border_style"`
	Alignment     string  `json:"alignment"`
}

// ChartDimensions contains chart size configuration
type ChartDimensions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// BrandColors contains brand color configuration
type BrandColors struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Accent    string `json:"accent"`
	Text      string `json:"text"`
	Background string `json:"background"`
}

// Performance configuration

// ExportPerformanceTargets contains performance targets for exports
type ExportPerformanceTargets struct {
	// Export generation times by format (in seconds)
	CSVGeneration    time.Duration `json:"csv_generation"`
	JSONGeneration   time.Duration `json:"json_generation"`
	XLSXGeneration   time.Duration `json:"xlsx_generation"`
	PDFGeneration    time.Duration `json:"pdf_generation"`
	ChartsGeneration time.Duration `json:"charts_generation"`
	
	// Memory usage limits (in MB)
	MaxMemoryPerExport int `json:"max_memory_per_export"`
	MaxMemoryTotal     int `json:"max_memory_total"`
	
	// Throughput targets
	MaxExportsPerMinute int `json:"max_exports_per_minute"`
	MaxConcurrentUsers  int `json:"max_concurrent_users"`
	
	// File size limits (in MB)
	MaxFileSizeCSV   int `json:"max_file_size_csv"`
	MaxFileSizeJSON  int `json:"max_file_size_json"`
	MaxFileSizeXLSX  int `json:"max_file_size_xlsx"`
	MaxFileSizePDF   int `json:"max_file_size_pdf"`
	MaxFileSizeCharts int `json:"max_file_size_charts"`
}

// Security configuration

// ExportSecuritySettings contains security settings for exports
type ExportSecuritySettings struct {
	// Access control
	RequireAuthentication bool     `json:"require_authentication"`
	AllowedRoles         []string `json:"allowed_roles"`
	RateLimitPerUser     int      `json:"rate_limit_per_user"`     // Per hour
	RateLimitPerIP       int      `json:"rate_limit_per_ip"`       // Per hour
	
	// Data protection
	EnableEncryption     bool     `json:"enable_encryption"`
	EncryptionAlgorithm  string   `json:"encryption_algorithm"`
	KeyRotationInterval  time.Duration `json:"key_rotation_interval"`
	
	// Content security
	DisallowSensitiveData bool    `json:"disallow_sensitive_data"`
	SensitiveFields      []string `json:"sensitive_fields"`
	DataAnonymization    bool     `json:"data_anonymization"`
	
	// Audit settings
	EnableAuditLogging   bool     `json:"enable_audit_logging"`
	AuditLogRetention    time.Duration `json:"audit_log_retention"`
	
	// Download security
	SignedURLs           bool     `json:"signed_urls"`
	DownloadTimeLimit    time.Duration `json:"download_time_limit"`
	MaxDownloadAttempts  int      `json:"max_download_attempts"`
}

// GetDefaultExportConfig returns default configuration for export service
func GetDefaultExportConfig() *ExportConfig {
	return &ExportConfig{
		EnableCompression: true,
		EnableEncryption:  false,
		ExportTTL:         24 * time.Hour,
		MaxFileSize:       100 * 1024 * 1024, // 100MB
		MaxConcurrentJobs: 10,
		CleanupInterval:   1 * time.Hour,
		
		StoragePath:   "/var/herald/exports",
		CDNBaseURL:    "https://cdn.herald.lol/exports",
		SignedURLTTL:  4 * time.Hour,
		
		CSV: &CSVConfig{
			DefaultDelimiter:      ",",
			MaxRows:              100000,
			IncludeHeadersDefault: true,
			DateFormatDefault:     "2006-01-02 15:04:05",
			EncodingDefault:       "UTF-8",
			MaxColumnWidth:        1000,
			EscapeSpecialChars:    true,
		},
		
		JSON: &JSONConfig{
			PrettyPrintDefault: false,
			MaxDepth:          10,
			DateFormatDefault: "2006-01-02T15:04:05Z07:00",
			IncludeNulls:      false,
			CompressArrays:    true,
			StreamLargeData:   true,
		},
		
		XLSX: &XLSXConfig{
			MaxRows:               65536,
			MaxColumns:            256,
			MaxSheets:             10,
			DefaultSheetName:      "Herald Analytics",
			AutoFitColumnsDefault: true,
			EnableCharts:          true,
			DefaultChartTypes:     []string{"line", "bar", "pie"},
			
			HeaderStyle: &XLSXStyle{
				FontName:        "Calibri",
				FontSize:        12,
				FontBold:        true,
				FontColor:       "#FFFFFF",
				BackgroundColor: "#4472C4",
				BorderStyle:     "thin",
				Alignment:       "center",
			},
			
			DataStyle: &XLSXStyle{
				FontName:        "Calibri",
				FontSize:        11,
				FontBold:        false,
				FontColor:       "#000000",
				BackgroundColor: "#FFFFFF",
				BorderStyle:     "thin",
				Alignment:       "left",
			},
			
			ChartDimensions: &ChartDimensions{
				Width:  600,
				Height: 400,
			},
		},
		
		PDF: &PDFConfig{
			DefaultPageSize:    "A4",
			DefaultOrientation: "portrait",
			DefaultFontFamily:  "Arial",
			DefaultFontSize:    11,
			MaxPages:          200,
			
			DefaultMargins: &PDFMargins{
				Top:    25,
				Bottom: 25,
				Left:   25,
				Right:  25,
			},
			
			EnableWatermark:  false,
			WatermarkText:    "Herald.lol Gaming Analytics",
			WatermarkOpacity: 0.3,
			
			LogoPath: "/assets/herald-logo.png",
			BrandColors: &BrandColors{
				Primary:    "#4472C4",
				Secondary:  "#70AD47",
				Accent:     "#FFC000",
				Text:       "#2F2F2F",
				Background: "#FFFFFF",
			},
			
			ChartDPI:    300,
			EmbedCharts: true,
		},
		
		Charts: &ChartsConfig{
			DefaultWidth:       800,
			DefaultHeight:      600,
			DefaultTheme:       "herald",
			DefaultColorScheme: []string{"#4472C4", "#70AD47", "#FFC000", "#E1575A", "#7030A0"},
			MaxDataPoints:      10000,
			
			EnableD3:            true,
			EnableChartJS:       true,
			EnablePlotly:        false,
			
			EnableInteractivity: true,
			EnableZoom:          true,
			EnableTooltips:      true,
			
			SupportedFormats:    []string{"HTML", "PNG", "SVG", "PDF"},
		},
		
		PerformanceTargets: &ExportPerformanceTargets{
			CSVGeneration:      5 * time.Second,
			JSONGeneration:     3 * time.Second,
			XLSXGeneration:     15 * time.Second,
			PDFGeneration:      20 * time.Second,
			ChartsGeneration:   10 * time.Second,
			
			MaxMemoryPerExport:  256, // MB
			MaxMemoryTotal:      2048, // MB
			
			MaxExportsPerMinute: 100,
			MaxConcurrentUsers:  500,
			
			MaxFileSizeCSV:   50,  // MB
			MaxFileSizeJSON:  30,  // MB
			MaxFileSizeXLSX:  100, // MB
			MaxFileSizePDF:   150, // MB
			MaxFileSizeCharts: 25, // MB
		},
		
		SecuritySettings: &ExportSecuritySettings{
			RequireAuthentication: true,
			AllowedRoles:         []string{"user", "premium", "pro", "enterprise"},
			RateLimitPerUser:     50,  // Per hour
			RateLimitPerIP:       200, // Per hour
			
			EnableEncryption:     false,
			EncryptionAlgorithm:  "AES-256-GCM",
			KeyRotationInterval:  30 * 24 * time.Hour, // 30 days
			
			DisallowSensitiveData: true,
			SensitiveFields:      []string{"email", "ip_address", "device_id"},
			DataAnonymization:    false,
			
			EnableAuditLogging:  true,
			AuditLogRetention:   90 * 24 * time.Hour, // 90 days
			
			SignedURLs:          true,
			DownloadTimeLimit:   4 * time.Hour,
			MaxDownloadAttempts: 10,
		},
	}
}

// GetExportConfigByProfile returns configuration optimized for different user profiles
func GetExportConfigByProfile(profile string) *ExportConfig {
	baseConfig := GetDefaultExportConfig()
	
	switch profile {
	case "free":
		// Limited configuration for free users
		baseConfig.MaxFileSize = 10 * 1024 * 1024 // 10MB
		baseConfig.CSV.MaxRows = 1000
		baseConfig.XLSX.MaxRows = 1000
		baseConfig.XLSX.EnableCharts = false
		baseConfig.PDF.MaxPages = 10
		baseConfig.Charts.EnableInteractivity = false
		baseConfig.SecuritySettings.RateLimitPerUser = 5 // Per hour
		
	case "premium":
		// Enhanced configuration for premium users
		baseConfig.MaxFileSize = 50 * 1024 * 1024 // 50MB
		baseConfig.CSV.MaxRows = 10000
		baseConfig.XLSX.MaxRows = 10000
		baseConfig.PDF.MaxPages = 50
		baseConfig.SecuritySettings.RateLimitPerUser = 20 // Per hour
		
	case "pro":
		// Professional configuration
		baseConfig.MaxFileSize = 100 * 1024 * 1024 // 100MB
		baseConfig.CSV.MaxRows = 50000
		baseConfig.XLSX.MaxRows = 50000
		baseConfig.PDF.MaxPages = 100
		baseConfig.Charts.EnablePlotly = true
		baseConfig.SecuritySettings.RateLimitPerUser = 100 // Per hour
		
	case "enterprise":
		// Enterprise configuration with maximum limits
		baseConfig.MaxFileSize = 500 * 1024 * 1024 // 500MB
		baseConfig.CSV.MaxRows = 1000000
		baseConfig.XLSX.MaxRows = 65536
		baseConfig.PDF.MaxPages = 500
		baseConfig.EnableEncryption = true
		baseConfig.SecuritySettings.EnableEncryption = true
		baseConfig.SecuritySettings.RateLimitPerUser = 1000 // Per hour
	}
	
	return baseConfig
}

// GetFormatCapabilities returns capabilities for each export format
func GetFormatCapabilities() map[string]*FormatCapabilities {
	return map[string]*FormatCapabilities{
		"csv": {
			Name:              "CSV",
			SupportsCharts:    false,
			SupportsFormatting: false,
			SupportsMultiSheet: false,
			SupportsImages:    false,
			MaxFileSize:       50 * 1024 * 1024,  // 50MB
			MaxRows:           100000,
			MaxColumns:        1000,
			StreamingSupport:  true,
			CompressionRatio:  0.3, // 30% of original size when compressed
			AverageGenTime:    2,   // seconds
		},
		"json": {
			Name:              "JSON",
			SupportsCharts:    false,
			SupportsFormatting: false,
			SupportsMultiSheet: false,
			SupportsImages:    false,
			MaxFileSize:       30 * 1024 * 1024,  // 30MB
			MaxRows:           50000,
			MaxColumns:        500,
			StreamingSupport:  true,
			CompressionRatio:  0.4, // 40% of original size when compressed
			AverageGenTime:    1,   // seconds
		},
		"xlsx": {
			Name:              "Excel",
			SupportsCharts:    true,
			SupportsFormatting: true,
			SupportsMultiSheet: true,
			SupportsImages:    true,
			MaxFileSize:       100 * 1024 * 1024, // 100MB
			MaxRows:           65536,
			MaxColumns:        256,
			StreamingSupport:  false,
			CompressionRatio:  0.5, // 50% of original size when compressed
			AverageGenTime:    8,   // seconds
		},
		"pdf": {
			Name:              "PDF Report",
			SupportsCharts:    true,
			SupportsFormatting: true,
			SupportsMultiSheet: false,
			SupportsImages:    true,
			MaxFileSize:       150 * 1024 * 1024, // 150MB
			MaxRows:           10000,
			MaxColumns:        50,
			StreamingSupport:  false,
			CompressionRatio:  0.6, // 60% of original size when compressed
			AverageGenTime:    12,  // seconds
		},
		"charts": {
			Name:              "Interactive Charts",
			SupportsCharts:    true,
			SupportsFormatting: true,
			SupportsMultiSheet: false,
			SupportsImages:    true,
			MaxFileSize:       25 * 1024 * 1024,  // 25MB
			MaxRows:           10000,
			MaxColumns:        20,
			StreamingSupport:  false,
			CompressionRatio:  0.7, // 70% of original size when compressed
			AverageGenTime:    6,   // seconds
		},
	}
}

// FormatCapabilities describes capabilities of an export format
type FormatCapabilities struct {
	Name              string  `json:"name"`
	SupportsCharts    bool    `json:"supports_charts"`
	SupportsFormatting bool   `json:"supports_formatting"`
	SupportsMultiSheet bool   `json:"supports_multi_sheet"`
	SupportsImages    bool    `json:"supports_images"`
	MaxFileSize       int64   `json:"max_file_size"`
	MaxRows           int     `json:"max_rows"`
	MaxColumns        int     `json:"max_columns"`
	StreamingSupport  bool    `json:"streaming_support"`
	CompressionRatio  float64 `json:"compression_ratio"`
	AverageGenTime    int     `json:"average_gen_time"` // seconds
}

// GetSubscriptionExportLimits returns export limits by subscription tier
func GetSubscriptionExportLimits() map[string]*SubscriptionLimits {
	return map[string]*SubscriptionLimits{
		"free": {
			Tier:                "Free",
			MaxExportsPerDay:    5,
			MaxExportsPerMonth: 50,
			MaxFileSize:        10 * 1024 * 1024, // 10MB
			AllowedFormats:     []string{"csv", "json"},
			MaxDataRows:        1000,
			ChartsEnabled:      false,
			Priority:          "low",
			RetentionDays:     1,
		},
		"premium": {
			Tier:                "Premium",
			MaxExportsPerDay:    25,
			MaxExportsPerMonth: 500,
			MaxFileSize:        50 * 1024 * 1024, // 50MB
			AllowedFormats:     []string{"csv", "json", "xlsx"},
			MaxDataRows:        10000,
			ChartsEnabled:      true,
			Priority:          "normal",
			RetentionDays:     7,
		},
		"pro": {
			Tier:                "Pro",
			MaxExportsPerDay:    100,
			MaxExportsPerMonth: 2000,
			MaxFileSize:        100 * 1024 * 1024, // 100MB
			AllowedFormats:     []string{"csv", "json", "xlsx", "pdf"},
			MaxDataRows:        50000,
			ChartsEnabled:      true,
			Priority:          "high",
			RetentionDays:     30,
		},
		"enterprise": {
			Tier:                "Enterprise",
			MaxExportsPerDay:    1000,
			MaxExportsPerMonth: 20000,
			MaxFileSize:        500 * 1024 * 1024, // 500MB
			AllowedFormats:     []string{"csv", "json", "xlsx", "pdf", "charts"},
			MaxDataRows:        1000000,
			ChartsEnabled:      true,
			Priority:          "highest",
			RetentionDays:     90,
		},
	}
}

// SubscriptionLimits contains limits for different subscription tiers
type SubscriptionLimits struct {
	Tier                string   `json:"tier"`
	MaxExportsPerDay    int      `json:"max_exports_per_day"`
	MaxExportsPerMonth  int      `json:"max_exports_per_month"`
	MaxFileSize         int64    `json:"max_file_size"`
	AllowedFormats      []string `json:"allowed_formats"`
	MaxDataRows         int      `json:"max_data_rows"`
	ChartsEnabled       bool     `json:"charts_enabled"`
	Priority            string   `json:"priority"`
	RetentionDays       int      `json:"retention_days"`
}

// GetExportPerformanceTargets returns performance targets for Herald.lol gaming platform
func GetExportPerformanceTargets() *HeraldExportTargets {
	return &HeraldExportTargets{
		// Gaming-specific performance requirements
		MaxAnalyticsExportTime:  30 * time.Second, // <30s for match analysis exports
		MaxPlayerExportTime:     60 * time.Second, // <1min for comprehensive player data
		MaxTeamExportTime:      120 * time.Second, // <2min for team analytics
		
		// Concurrent user support
		ConcurrentExports:       100, // Support 100+ simultaneous exports
		PeakConcurrentUsers:    1000, // Support 1K+ concurrent users
		
		// Gaming data volume targets
		MaxMatchesPerExport:    1000,  // Up to 1000 matches per export
		MaxChampionDataPoints: 5000,  // Up to 5000 champion performance points
		MaxTeamMembers:        10,    // Support up to 10-player team analytics
		
		// Performance metrics
		ExportSuccessRate:      99.5, // 99.5% success rate
		AverageQueueTime:       3,    // 3 seconds average queue time
		MaxMemoryUsage:         512,  // 512MB max memory per export
		
		// Gaming-specific features
		ChartGenerationTime:    15 * time.Second, // <15s for gaming charts
		InteractiveChartsTime:  20 * time.Second, // <20s for interactive visualizations
		RealTimeDataLatency:    5 * time.Second,  // <5s for real-time data inclusion
	}
}

// HeraldExportTargets contains Herald.lol specific performance targets
type HeraldExportTargets struct {
	MaxAnalyticsExportTime  time.Duration `json:"max_analytics_export_time"`
	MaxPlayerExportTime     time.Duration `json:"max_player_export_time"`
	MaxTeamExportTime       time.Duration `json:"max_team_export_time"`
	
	ConcurrentExports       int           `json:"concurrent_exports"`
	PeakConcurrentUsers     int           `json:"peak_concurrent_users"`
	
	MaxMatchesPerExport     int           `json:"max_matches_per_export"`
	MaxChampionDataPoints   int           `json:"max_champion_data_points"`
	MaxTeamMembers          int           `json:"max_team_members"`
	
	ExportSuccessRate       float64       `json:"export_success_rate"`
	AverageQueueTime        int           `json:"average_queue_time"`
	MaxMemoryUsage          int           `json:"max_memory_usage"`
	
	ChartGenerationTime     time.Duration `json:"chart_generation_time"`
	InteractiveChartsTime   time.Duration `json:"interactive_charts_time"`
	RealTimeDataLatency     time.Duration `json:"real_time_data_latency"`
}