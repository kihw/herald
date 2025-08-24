package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Herald.lol Gaming Analytics - Gaming Audit Logger
// Comprehensive audit logging and compliance tracking for gaming platform

// GamingAuditLogger handles gaming-specific audit logging
type GamingAuditLogger struct {
	redis             *redis.Client
	config            *AuditConfig
	complianceTracker *ComplianceTracker
}

// AuditConfig contains audit logging configuration
type AuditConfig struct {
	// Logging levels
	LogLevel          string   `json:"log_level"`          // debug, info, warn, error
	EnabledCategories []string `json:"enabled_categories"` // auth, gaming, api, compliance

	// Storage settings
	RetentionDays   int `json:"retention_days"`
	MaxEventsPerDay int `json:"max_events_per_day"`

	// Gaming-specific settings
	LogGamingOperations bool `json:"log_gaming_operations"`
	LogRiotAPIAccess    bool `json:"log_riot_api_access"`
	LogDataExports      bool `json:"log_data_exports"`
	LogUserAccess       bool `json:"log_user_access"`
	LogTeamOperations   bool `json:"log_team_operations"`

	// Compliance settings
	GDPRLogging       bool `json:"gdpr_logging"`
	RiotComplianceLog bool `json:"riot_compliance_log"`
	SecurityEventLog  bool `json:"security_event_log"`

	// Alert settings
	RealTimeAlerts  bool           `json:"real_time_alerts"`
	AlertThresholds map[string]int `json:"alert_thresholds"`
}

// AuditEvent represents a gaming audit event
type AuditEvent struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Category  string    `json:"category"` // auth, gaming, api, compliance, security
	Action    string    `json:"action"`
	UserID    string    `json:"user_id,omitempty"`
	ClientIP  string    `json:"client_ip"`
	UserAgent string    `json:"user_agent,omitempty"`

	// Gaming-specific fields
	GamingContext *GamingAuditContext `json:"gaming_context,omitempty"`

	// Request/Response data
	RequestPath   string        `json:"request_path,omitempty"`
	RequestMethod string        `json:"request_method,omitempty"`
	ResponseCode  int           `json:"response_code,omitempty"`
	ResponseTime  time.Duration `json:"response_time,omitempty"`

	// Additional data
	Data map[string]interface{} `json:"data,omitempty"`
	Tags []string               `json:"tags,omitempty"`

	// Security fields
	ThreatLevel     string   `json:"threat_level,omitempty"` // low, medium, high, critical
	ComplianceFlags []string `json:"compliance_flags,omitempty"`

	// Session info
	SessionID        string `json:"session_id,omitempty"`
	APIKey           string `json:"api_key,omitempty"`
	SubscriptionTier string `json:"subscription_tier,omitempty"`
}

// GamingAuditContext contains gaming-specific audit context
type GamingAuditContext struct {
	Region           string `json:"region,omitempty"`
	GameType         string `json:"game_type,omitempty"` // LoL, TFT
	SummonerName     string `json:"summoner_name,omitempty"`
	TeamID           string `json:"team_id,omitempty"`
	MatchID          string `json:"match_id,omitempty"`
	AnalyticsType    string `json:"analytics_type,omitempty"` // basic, advanced, real-time
	DataSize         string `json:"data_size,omitempty"`      // small, medium, large
	ExportFormat     string `json:"export_format,omitempty"`  // json, csv, xlsx, pdf
	RiotAPIEndpoint  string `json:"riot_api_endpoint,omitempty"`
	ProcessingTimeMs int64  `json:"processing_time_ms,omitempty"`
}

// NewGamingAuditLogger creates new gaming audit logger
func NewGamingAuditLogger(redis *redis.Client, config *AuditConfig) *GamingAuditLogger {
	return &GamingAuditLogger{
		redis:             redis,
		config:            config,
		complianceTracker: NewComplianceTracker(redis),
	}
}

// LogEvent logs a gaming audit event
func (g *GamingAuditLogger) LogEvent(ctx context.Context, event *AuditEvent) error {
	// Set defaults
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Check if this category is enabled
	if !g.isCategoryEnabled(event.Category) {
		return nil
	}

	// Serialize event
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize audit event: %w", err)
	}

	// Store in Redis with TTL
	eventKey := fmt.Sprintf("audit:event:%s", event.ID)
	err = g.redis.Set(ctx, eventKey, eventData, time.Duration(g.config.RetentionDays)*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store audit event: %w", err)
	}

	// Add to time-series indices
	if err := g.addToTimeSeriesIndex(ctx, event); err != nil {
		// Log error but don't fail the main operation
		fmt.Printf("Failed to add to time series index: %v\n", err)
	}

	// Add to category index
	if err := g.addToCategoryIndex(ctx, event); err != nil {
		fmt.Printf("Failed to add to category index: %v\n", err)
	}

	// Check for compliance tracking
	if g.config.GDPRLogging || g.config.RiotComplianceLog {
		if err := g.complianceTracker.TrackEvent(ctx, event); err != nil {
			fmt.Printf("Failed to track compliance: %v\n", err)
		}
	}

	// Check for real-time alerts
	if g.config.RealTimeAlerts {
		if err := g.checkAlerts(ctx, event); err != nil {
			fmt.Printf("Failed to check alerts: %v\n", err)
		}
	}

	return nil
}

// AuditMiddleware creates Gin middleware for automatic audit logging
func (g *GamingAuditLogger) AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Store request data
		requestPath := c.Request.URL.Path
		requestMethod := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		// Process request
		c.Next()

		// Calculate response time
		responseTime := time.Since(start)
		responseCode := c.Writer.Status()

		// Determine if we should log this request
		if !g.shouldLogRequest(requestPath, requestMethod, responseCode) {
			return
		}

		// Create audit event
		event := &AuditEvent{
			Category:         g.getRequestCategory(requestPath),
			Action:           g.getActionFromRequest(requestMethod, requestPath),
			UserID:           g.getUserID(c),
			ClientIP:         clientIP,
			UserAgent:        userAgent,
			RequestPath:      requestPath,
			RequestMethod:    requestMethod,
			ResponseCode:     responseCode,
			ResponseTime:     responseTime,
			SessionID:        g.getSessionID(c),
			APIKey:           g.getAPIKey(c),
			SubscriptionTier: g.getSubscriptionTier(c),
		}

		// Add gaming context if applicable
		if g.isGamingRequest(requestPath) {
			event.GamingContext = &GamingAuditContext{
				Region:           c.Param("region"),
				SummonerName:     c.Param("summonerName"),
				TeamID:           c.Param("teamId"),
				MatchID:          c.Param("matchId"),
				AnalyticsType:    g.getAnalyticsType(requestPath),
				DataSize:         g.getDataSize(requestPath),
				ExportFormat:     g.getExportFormat(requestPath),
				RiotAPIEndpoint:  g.getRiotEndpoint(requestPath),
				ProcessingTimeMs: responseTime.Milliseconds(),
			}
		}

		// Add threat level if error occurred
		if responseCode >= 400 {
			event.ThreatLevel = g.getThreatLevel(responseCode, requestPath)
		}

		// Add compliance flags
		if g.needsComplianceLogging(c, event) {
			event.ComplianceFlags = g.getComplianceFlags(c, event)
		}

		// Log the event
		if err := g.LogEvent(c.Request.Context(), event); err != nil {
			fmt.Printf("Failed to log audit event: %v\n", err)
		}
	}
}

// LogGamingOperation logs specific gaming operations
func (g *GamingAuditLogger) LogGamingOperation(ctx context.Context, operation *GamingOperation) error {
	if !g.config.LogGamingOperations {
		return nil
	}

	event := &AuditEvent{
		Category: "gaming",
		Action:   operation.Type,
		UserID:   operation.UserID,
		ClientIP: operation.ClientIP,
		GamingContext: &GamingAuditContext{
			Region:           operation.Region,
			GameType:         operation.GameType,
			SummonerName:     operation.SummonerName,
			TeamID:           operation.TeamID,
			AnalyticsType:    operation.AnalyticsType,
			ProcessingTimeMs: operation.ProcessingTimeMs,
		},
		Data: operation.Metadata,
		Tags: operation.Tags,
	}

	return g.LogEvent(ctx, event)
}

// LogDataExport logs data export operations with enhanced security tracking
func (g *GamingAuditLogger) LogDataExport(ctx context.Context, export *DataExportOperation) error {
	if !g.config.LogDataExports {
		return nil
	}

	event := &AuditEvent{
		Category: "data_export",
		Action:   "export_gaming_data",
		UserID:   export.UserID,
		ClientIP: export.ClientIP,
		GamingContext: &GamingAuditContext{
			Region:       export.Region,
			SummonerName: export.SummonerName,
			DataSize:     export.DataSize,
			ExportFormat: export.Format,
		},
		Data: map[string]interface{}{
			"export_type":   export.Type,
			"file_size_mb":  export.FileSizeMB,
			"record_count":  export.RecordCount,
			"mfa_verified":  export.MFAVerified,
			"export_reason": export.Reason,
		},
		ThreatLevel:      "medium", // Data exports are always medium threat
		ComplianceFlags:  []string{"gdpr_data_export", "gaming_data_export"},
		SubscriptionTier: export.SubscriptionTier,
	}

	return g.LogEvent(ctx, event)
}

// LogSecurityEvent logs security-related events
func (g *GamingAuditLogger) LogSecurityEvent(ctx context.Context, secEvent *SecurityEvent) error {
	if !g.config.SecurityEventLog {
		return nil
	}

	event := &AuditEvent{
		Category:    "security",
		Action:      secEvent.Type,
		UserID:      secEvent.UserID,
		ClientIP:    secEvent.ClientIP,
		UserAgent:   secEvent.UserAgent,
		ThreatLevel: secEvent.ThreatLevel,
		Data: map[string]interface{}{
			"event_details":      secEvent.Details,
			"blocked":            secEvent.Blocked,
			"alert_triggered":    secEvent.AlertTriggered,
			"automated_response": secEvent.AutomatedResponse,
		},
		Tags: secEvent.Tags,
	}

	return g.LogEvent(ctx, event)
}

// GetAuditEvents retrieves audit events with filtering
func (g *GamingAuditLogger) GetAuditEvents(ctx context.Context, filter *AuditFilter) ([]*AuditEvent, error) {
	// Build query based on filter
	var events []*AuditEvent

	// Get events from time range
	startTime := filter.StartTime
	endTime := filter.EndTime
	if endTime.IsZero() {
		endTime = time.Now()
	}

	// Query time series index
	timeSeriesKey := fmt.Sprintf("audit:time_series:%s", startTime.Format("2006-01-02"))
	eventIDs, err := g.redis.ZRangeByScore(ctx, timeSeriesKey, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", startTime.Unix()),
		Max: fmt.Sprintf("%d", endTime.Unix()),
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to query time series: %w", err)
	}

	// Retrieve events
	for _, eventID := range eventIDs {
		if len(events) >= filter.Limit {
			break
		}

		eventKey := fmt.Sprintf("audit:event:%s", eventID)
		eventData, err := g.redis.Get(ctx, eventKey).Result()
		if err != nil {
			continue
		}

		var event AuditEvent
		if err := json.Unmarshal([]byte(eventData), &event); err != nil {
			continue
		}

		// Apply filters
		if g.matchesFilter(&event, filter) {
			events = append(events, &event)
		}
	}

	return events, nil
}

// GetComplianceReport generates compliance report
func (g *GamingAuditLogger) GetComplianceReport(ctx context.Context, reportType string, timeRange time.Duration) (*ComplianceReport, error) {
	return g.complianceTracker.GenerateReport(ctx, reportType, timeRange)
}

// Helper methods

func (g *GamingAuditLogger) isCategoryEnabled(category string) bool {
	if len(g.config.EnabledCategories) == 0 {
		return true // All categories enabled by default
	}

	for _, enabled := range g.config.EnabledCategories {
		if enabled == category {
			return true
		}
	}
	return false
}

func (g *GamingAuditLogger) shouldLogRequest(path, method string, responseCode int) bool {
	// Always log errors
	if responseCode >= 400 {
		return true
	}

	// Always log gaming operations
	if g.isGamingRequest(path) && g.config.LogGamingOperations {
		return true
	}

	// Log API access
	if g.config.LogUserAccess && method != "GET" {
		return true
	}

	return false
}

func (g *GamingAuditLogger) isGamingRequest(path string) bool {
	gamingPaths := []string{"/gaming/", "/riot/", "/analytics/", "/teams/", "/matches/"}
	for _, gamingPath := range gamingPaths {
		if strings.Contains(path, gamingPath) {
			return true
		}
	}
	return false
}

func (g *GamingAuditLogger) getRequestCategory(path string) string {
	switch {
	case strings.Contains(path, "/auth/"):
		return "auth"
	case strings.Contains(path, "/gaming/") || strings.Contains(path, "/riot/"):
		return "gaming"
	case strings.Contains(path, "/export"):
		return "data_export"
	case strings.Contains(path, "/teams/"):
		return "team"
	default:
		return "api"
	}
}

func (g *GamingAuditLogger) getActionFromRequest(method, path string) string {
	switch method {
	case "GET":
		if strings.Contains(path, "/analytics/") {
			return "view_analytics"
		}
		return "view"
	case "POST":
		if strings.Contains(path, "/export") {
			return "export_data"
		}
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return method
	}
}

func (g *GamingAuditLogger) addToTimeSeriesIndex(ctx context.Context, event *AuditEvent) error {
	timeSeriesKey := fmt.Sprintf("audit:time_series:%s", event.Timestamp.Format("2006-01-02"))
	score := float64(event.Timestamp.Unix())

	err := g.redis.ZAdd(ctx, timeSeriesKey, &redis.Z{
		Score:  score,
		Member: event.ID,
	}).Err()
	if err != nil {
		return err
	}

	// Set expiration on time series key
	g.redis.Expire(ctx, timeSeriesKey, time.Duration(g.config.RetentionDays)*24*time.Hour)

	return nil
}

func (g *GamingAuditLogger) addToCategoryIndex(ctx context.Context, event *AuditEvent) error {
	categoryKey := fmt.Sprintf("audit:category:%s", event.Category)
	score := float64(event.Timestamp.Unix())

	err := g.redis.ZAdd(ctx, categoryKey, &redis.Z{
		Score:  score,
		Member: event.ID,
	}).Err()
	if err != nil {
		return err
	}

	g.redis.Expire(ctx, categoryKey, time.Duration(g.config.RetentionDays)*24*time.Hour)

	return nil
}

func (g *GamingAuditLogger) getUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

func (g *GamingAuditLogger) getSessionID(c *gin.Context) string {
	if sessionID, exists := c.Get("session_id"); exists {
		if id, ok := sessionID.(string); ok {
			return id
		}
	}
	return ""
}

func (g *GamingAuditLogger) getAPIKey(c *gin.Context) string {
	if apiKey, exists := c.Get("api_key"); exists {
		if key, ok := apiKey.(string); ok {
			return key
		}
	}
	return ""
}

func (g *GamingAuditLogger) getSubscriptionTier(c *gin.Context) string {
	if tier, exists := c.Get("subscription_tier"); exists {
		if t, ok := tier.(string); ok {
			return t
		}
	}
	return ""
}

func (g *GamingAuditLogger) matchesFilter(event *AuditEvent, filter *AuditFilter) bool {
	// Apply category filter
	if filter.Category != "" && event.Category != filter.Category {
		return false
	}

	// Apply user filter
	if filter.UserID != "" && event.UserID != filter.UserID {
		return false
	}

	// Apply action filter
	if filter.Action != "" && event.Action != filter.Action {
		return false
	}

	return true
}

// Data structures for operations

type GamingOperation struct {
	Type             string                 `json:"type"`
	UserID           string                 `json:"user_id"`
	ClientIP         string                 `json:"client_ip"`
	Region           string                 `json:"region,omitempty"`
	GameType         string                 `json:"game_type,omitempty"`
	SummonerName     string                 `json:"summoner_name,omitempty"`
	TeamID           string                 `json:"team_id,omitempty"`
	AnalyticsType    string                 `json:"analytics_type,omitempty"`
	ProcessingTimeMs int64                  `json:"processing_time_ms,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	Tags             []string               `json:"tags,omitempty"`
}

type DataExportOperation struct {
	UserID           string `json:"user_id"`
	ClientIP         string `json:"client_ip"`
	Region           string `json:"region,omitempty"`
	SummonerName     string `json:"summoner_name,omitempty"`
	Type             string `json:"type"`
	Format           string `json:"format"`
	DataSize         string `json:"data_size"`
	FileSizeMB       int    `json:"file_size_mb"`
	RecordCount      int    `json:"record_count"`
	MFAVerified      bool   `json:"mfa_verified"`
	Reason           string `json:"reason,omitempty"`
	SubscriptionTier string `json:"subscription_tier"`
}

type SecurityEvent struct {
	Type              string                 `json:"type"`
	UserID            string                 `json:"user_id,omitempty"`
	ClientIP          string                 `json:"client_ip"`
	UserAgent         string                 `json:"user_agent,omitempty"`
	ThreatLevel       string                 `json:"threat_level"`
	Details           map[string]interface{} `json:"details"`
	Blocked           bool                   `json:"blocked"`
	AlertTriggered    bool                   `json:"alert_triggered"`
	AutomatedResponse string                 `json:"automated_response,omitempty"`
	Tags              []string               `json:"tags,omitempty"`
}

type AuditFilter struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Category  string    `json:"category,omitempty"`
	Action    string    `json:"action,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	ClientIP  string    `json:"client_ip,omitempty"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}

// DefaultAuditConfig returns default audit configuration
func DefaultAuditConfig() *AuditConfig {
	return &AuditConfig{
		LogLevel:            "info",
		EnabledCategories:   []string{"auth", "gaming", "api", "compliance", "security"},
		RetentionDays:       90,
		MaxEventsPerDay:     1000000,
		LogGamingOperations: true,
		LogRiotAPIAccess:    true,
		LogDataExports:      true,
		LogUserAccess:       true,
		LogTeamOperations:   true,
		GDPRLogging:         true,
		RiotComplianceLog:   true,
		SecurityEventLog:    true,
		RealTimeAlerts:      true,
		AlertThresholds: map[string]int{
			"failed_logins":   10,
			"data_exports":    5,
			"security_events": 3,
		},
	}
}
