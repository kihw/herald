package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Herald.lol Gaming Analytics - Notification Service
// Comprehensive notification system supporting real-time, email, and push notifications

// NotificationService handles all types of notifications
type NotificationService struct {
	config *NotificationConfig
	
	// Real-time connections
	wsConnections map[string]*WebSocketConnection
	connectionsMu sync.RWMutex
	
	// Notification channels
	notificationQueue chan *NotificationJob
	emailQueue        chan *EmailNotification
	pushQueue         chan *PushNotification
	
	// Service components
	emailProvider EmailProvider
	pushProvider  PushProvider
	templateEngine TemplateEngine
	
	// Worker management
	workers int
	shutdown chan bool
	wg      sync.WaitGroup
}

// NotificationConfig contains service configuration
type NotificationConfig struct {
	// Queue settings
	QueueSize           int           `json:"queue_size"`
	WorkerCount         int           `json:"worker_count"`
	ProcessingTimeout   time.Duration `json:"processing_timeout"`
	RetryAttempts       int           `json:"retry_attempts"`
	RetryDelay          time.Duration `json:"retry_delay"`
	
	// Real-time settings
	WebSocketReadTimeout  time.Duration `json:"websocket_read_timeout"`
	WebSocketWriteTimeout time.Duration `json:"websocket_write_timeout"`
	PingInterval         time.Duration `json:"ping_interval"`
	MaxConnections       int           `json:"max_connections"`
	
	// Email settings
	EmailEnabled         bool          `json:"email_enabled"`
	EmailTemplatesPath   string        `json:"email_templates_path"`
	EmailRateLimit       int           `json:"email_rate_limit"`
	
	// Push settings
	PushEnabled          bool          `json:"push_enabled"`
	PushRateLimit        int           `json:"push_rate_limit"`
	
	// Notification settings
	DefaultRetentionDays int           `json:"default_retention_days"`
	EnableBatching       bool          `json:"enable_batching"`
	BatchSize            int           `json:"batch_size"`
	BatchFlushInterval   time.Duration `json:"batch_flush_interval"`
}

// WebSocketConnection represents a WebSocket connection
type WebSocketConnection struct {
	UserID     string          `json:"user_id"`
	Conn       *websocket.Conn `json:"-"`
	Send       chan []byte     `json:"-"`
	CreatedAt  time.Time       `json:"created_at"`
	LastPing   time.Time       `json:"last_ping"`
	ClientInfo *ClientInfo     `json:"client_info"`
}

// ClientInfo contains client information
type ClientInfo struct {
	UserAgent    string `json:"user_agent"`
	IPAddress    string `json:"ip_address"`
	Platform     string `json:"platform"`
	Version      string `json:"version"`
	Capabilities []string `json:"capabilities"`
}

// NotificationJob represents a notification processing job
type NotificationJob struct {
	ID           string                 `json:"id"`
	Type         NotificationType       `json:"type"`
	UserID       string                 `json:"user_id"`
	Channels     []NotificationChannel  `json:"channels"`
	Content      *NotificationContent   `json:"content"`
	Priority     NotificationPriority   `json:"priority"`
	ScheduledAt  *time.Time             `json:"scheduled_at"`
	ExpiresAt    *time.Time             `json:"expires_at"`
	Context      map[string]interface{} `json:"context"`
	CreatedAt    time.Time              `json:"created_at"`
	Status       string                 `json:"status"`
	Attempts     int                    `json:"attempts"`
	LastError    string                 `json:"last_error"`
}

// NotificationType represents different types of notifications
type NotificationType string

const (
	NotificationTypeMatchComplete   NotificationType = "match_complete"
	NotificationTypeRankChange     NotificationType = "rank_change"
	NotificationTypeAchievement    NotificationType = "achievement"
	NotificationTypeCoachingTip    NotificationType = "coaching_tip"
	NotificationTypeSystemAlert    NotificationType = "system_alert"
	NotificationTypeWeeklyReport   NotificationType = "weekly_report"
	NotificationTypeFriendActivity NotificationType = "friend_activity"
	NotificationTypeMatchAlert     NotificationType = "match_alert"
	NotificationTypeCustom         NotificationType = "custom"
)

// NotificationChannel represents delivery channels
type NotificationChannel string

const (
	ChannelRealTime NotificationChannel = "realtime"
	ChannelEmail    NotificationChannel = "email"
	ChannelPush     NotificationChannel = "push"
	ChannelSMS      NotificationChannel = "sms"
	ChannelInApp    NotificationChannel = "in_app"
)

// NotificationPriority represents priority levels
type NotificationPriority string

const (
	PriorityLow      NotificationPriority = "low"
	PriorityNormal   NotificationPriority = "normal"
	PriorityHigh     NotificationPriority = "high"
	PriorityCritical NotificationPriority = "critical"
)

// NotificationContent contains notification content
type NotificationContent struct {
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Data        map[string]interface{} `json:"data"`
	Actions     []*NotificationAction  `json:"actions"`
	Icon        string                 `json:"icon"`
	Image       string                 `json:"image"`
	Sound       string                 `json:"sound"`
	Badge       int                    `json:"badge"`
	Category    string                 `json:"category"`
	ThreadID    string                 `json:"thread_id"`
	Tags        []string               `json:"tags"`
	URL         string                 `json:"url"`
	ClickAction string                 `json:"click_action"`
}

// NotificationAction represents an action button
type NotificationAction struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
	URL   string `json:"url"`
}

// EmailNotification represents an email notification
type EmailNotification struct {
	ID          string            `json:"id"`
	To          []string          `json:"to"`
	CC          []string          `json:"cc"`
	BCC         []string          `json:"bcc"`
	Subject     string            `json:"subject"`
	TextBody    string            `json:"text_body"`
	HTMLBody    string            `json:"html_body"`
	Template    string            `json:"template"`
	TemplateData map[string]interface{} `json:"template_data"`
	Attachments []*EmailAttachment `json:"attachments"`
	Priority    NotificationPriority `json:"priority"`
	ScheduledAt *time.Time        `json:"scheduled_at"`
	CreatedAt   time.Time         `json:"created_at"`
}

// EmailAttachment represents an email attachment
type EmailAttachment struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
	MimeType string `json:"mime_type"`
}

// PushNotification represents a push notification
type PushNotification struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	DeviceToken string                 `json:"device_token"`
	Platform    string                 `json:"platform"`
	Title       string                 `json:"title"`
	Body        string                 `json:"body"`
	Data        map[string]interface{} `json:"data"`
	Badge       int                    `json:"badge"`
	Sound       string                 `json:"sound"`
	Icon        string                 `json:"icon"`
	ClickAction string                 `json:"click_action"`
	TTL         int                    `json:"ttl"`
	Priority    NotificationPriority   `json:"priority"`
	CreatedAt   time.Time              `json:"created_at"`
}

// EmailProvider interface for email providers
type EmailProvider interface {
	SendEmail(ctx context.Context, email *EmailNotification) error
	ValidateEmail(email string) bool
	GetDeliveryStatus(emailID string) (*DeliveryStatus, error)
}

// PushProvider interface for push notification providers
type PushProvider interface {
	SendPush(ctx context.Context, push *PushNotification) error
	ValidateDeviceToken(token string, platform string) bool
	GetDeliveryStatus(pushID string) (*DeliveryStatus, error)
}

// TemplateEngine interface for template rendering
type TemplateEngine interface {
	RenderEmailTemplate(templateName string, data map[string]interface{}) (string, string, error)
	RenderPushTemplate(templateName string, data map[string]interface{}) (string, string, error)
}

// DeliveryStatus represents delivery status
type DeliveryStatus struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"`
	DeliveredAt *time.Time `json:"delivered_at"`
	Error       string    `json:"error"`
}

// NotificationPreferences represents user notification preferences
type NotificationPreferences struct {
	UserID              string                                 `json:"user_id"`
	Channels            map[NotificationChannel]bool           `json:"channels"`
	Types               map[NotificationType]bool              `json:"types"`
	QuietHours          *QuietHours                            `json:"quiet_hours"`
	Frequency           map[NotificationType]string            `json:"frequency"`
	CustomSettings      map[string]interface{}                 `json:"custom_settings"`
	LastUpdated         time.Time                              `json:"last_updated"`
}

// QuietHours represents quiet hours settings
type QuietHours struct {
	Enabled   bool   `json:"enabled"`
	StartTime string `json:"start_time"` // "22:00"
	EndTime   string `json:"end_time"`   // "08:00"
	Timezone  string `json:"timezone"`
}

// Notification statistics and metrics
type NotificationMetrics struct {
	TotalSent      int                            `json:"total_sent"`
	TotalDelivered int                            `json:"total_delivered"`
	TotalFailed    int                            `json:"total_failed"`
	DeliveryRate   float64                        `json:"delivery_rate"`
	ByChannel      map[NotificationChannel]int    `json:"by_channel"`
	ByType         map[NotificationType]int       `json:"by_type"`
	ByPriority     map[NotificationPriority]int   `json:"by_priority"`
	AverageLatency time.Duration                  `json:"average_latency"`
	Timestamp      time.Time                      `json:"timestamp"`
}

// NewNotificationService creates a new notification service
func NewNotificationService(config *NotificationConfig, emailProvider EmailProvider, pushProvider PushProvider, templateEngine TemplateEngine) *NotificationService {
	if config == nil {
		config = DefaultNotificationConfig()
	}

	service := &NotificationService{
		config:            config,
		wsConnections:     make(map[string]*WebSocketConnection),
		notificationQueue: make(chan *NotificationJob, config.QueueSize),
		emailQueue:        make(chan *EmailNotification, config.QueueSize),
		pushQueue:         make(chan *PushNotification, config.QueueSize),
		emailProvider:     emailProvider,
		pushProvider:      pushProvider,
		templateEngine:    templateEngine,
		workers:           config.WorkerCount,
		shutdown:          make(chan bool),
	}

	// Start worker goroutines
	service.startWorkers()

	return service
}

// SendNotification sends a notification through specified channels
func (s *NotificationService) SendNotification(ctx context.Context, notification *NotificationJob) error {
	if notification == nil {
		return fmt.Errorf("notification cannot be nil")
	}

	// Set default values
	if notification.ID == "" {
		notification.ID = fmt.Sprintf("notif_%d", time.Now().UnixNano())
	}
	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = time.Now()
	}
	if notification.Priority == "" {
		notification.Priority = PriorityNormal
	}
	if notification.Status == "" {
		notification.Status = "pending"
	}

	// Validate notification
	if err := s.validateNotification(notification); err != nil {
		return fmt.Errorf("invalid notification: %w", err)
	}

	// Check if notification is scheduled for future
	if notification.ScheduledAt != nil && notification.ScheduledAt.After(time.Now()) {
		return s.scheduleNotification(notification)
	}

	// Queue notification for immediate processing
	select {
	case s.notificationQueue <- notification:
		log.Printf("Notification queued: %s", notification.ID)
		return nil
	case <-ctx.Done():
		return fmt.Errorf("context cancelled while queueing notification")
	default:
		return fmt.Errorf("notification queue is full")
	}
}

// SendRealTimeNotification sends a real-time notification via WebSocket
func (s *NotificationService) SendRealTimeNotification(userID string, data map[string]interface{}) error {
	s.connectionsMu.RLock()
	conn, exists := s.wsConnections[userID]
	s.connectionsMu.RUnlock()

	if !exists {
		return fmt.Errorf("user %s not connected", userID)
	}

	message, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	select {
	case conn.Send <- message:
		return nil
	default:
		return fmt.Errorf("failed to send notification to user %s", userID)
	}
}

// RegisterWebSocketConnection registers a new WebSocket connection
func (s *NotificationService) RegisterWebSocketConnection(userID string, conn *websocket.Conn, clientInfo *ClientInfo) error {
	s.connectionsMu.Lock()
	defer s.connectionsMu.Unlock()

	// Check connection limit
	if len(s.wsConnections) >= s.config.MaxConnections {
		return fmt.Errorf("maximum WebSocket connections reached")
	}

	// Close existing connection if any
	if existingConn, exists := s.wsConnections[userID]; exists {
		close(existingConn.Send)
		existingConn.Conn.Close()
	}

	wsConn := &WebSocketConnection{
		UserID:     userID,
		Conn:       conn,
		Send:       make(chan []byte, 256),
		CreatedAt:  time.Now(),
		LastPing:   time.Now(),
		ClientInfo: clientInfo,
	}

	s.wsConnections[userID] = wsConn

	// Start connection handler
	go s.handleWebSocketConnection(wsConn)

	log.Printf("WebSocket connection registered for user: %s", userID)
	return nil
}

// UnregisterWebSocketConnection removes a WebSocket connection
func (s *NotificationService) UnregisterWebSocketConnection(userID string) {
	s.connectionsMu.Lock()
	defer s.connectionsMu.Unlock()

	if conn, exists := s.wsConnections[userID]; exists {
		close(conn.Send)
		conn.Conn.Close()
		delete(s.wsConnections, userID)
		log.Printf("WebSocket connection unregistered for user: %s", userID)
	}
}

// SendEmailNotification sends an email notification
func (s *NotificationService) SendEmailNotification(ctx context.Context, email *EmailNotification) error {
	if !s.config.EmailEnabled || s.emailProvider == nil {
		return fmt.Errorf("email notifications are disabled")
	}

	if email.ID == "" {
		email.ID = fmt.Sprintf("email_%d", time.Now().UnixNano())
	}
	if email.CreatedAt.IsZero() {
		email.CreatedAt = time.Now()
	}

	// Process template if specified
	if email.Template != "" && s.templateEngine != nil {
		textBody, htmlBody, err := s.templateEngine.RenderEmailTemplate(email.Template, email.TemplateData)
		if err != nil {
			return fmt.Errorf("failed to render email template: %w", err)
		}
		email.TextBody = textBody
		email.HTMLBody = htmlBody
	}

	select {
	case s.emailQueue <- email:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("context cancelled while queueing email")
	default:
		return fmt.Errorf("email queue is full")
	}
}

// SendPushNotification sends a push notification
func (s *NotificationService) SendPushNotification(ctx context.Context, push *PushNotification) error {
	if !s.config.PushEnabled || s.pushProvider == nil {
		return fmt.Errorf("push notifications are disabled")
	}

	if push.ID == "" {
		push.ID = fmt.Sprintf("push_%d", time.Now().UnixNano())
	}
	if push.CreatedAt.IsZero() {
		push.CreatedAt = time.Now()
	}

	select {
	case s.pushQueue <- push:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("context cancelled while queueing push notification")
	default:
		return fmt.Errorf("push notification queue is full")
	}
}

// GetNotificationPreferences retrieves user notification preferences
func (s *NotificationService) GetNotificationPreferences(userID string) (*NotificationPreferences, error) {
	// Mock implementation - would fetch from database in real system
	return &NotificationPreferences{
		UserID: userID,
		Channels: map[NotificationChannel]bool{
			ChannelRealTime: true,
			ChannelEmail:    true,
			ChannelPush:     true,
			ChannelInApp:    true,
		},
		Types: map[NotificationType]bool{
			NotificationTypeMatchComplete:   true,
			NotificationTypeRankChange:     true,
			NotificationTypeAchievement:    true,
			NotificationTypeCoachingTip:    false,
			NotificationTypeSystemAlert:    true,
			NotificationTypeWeeklyReport:   true,
			NotificationTypeFriendActivity: false,
			NotificationTypeMatchAlert:     true,
		},
		QuietHours: &QuietHours{
			Enabled:   true,
			StartTime: "22:00",
			EndTime:   "08:00",
			Timezone:  "UTC",
		},
		Frequency: map[NotificationType]string{
			NotificationTypeMatchComplete:   "immediate",
			NotificationTypeCoachingTip:    "daily",
			NotificationTypeWeeklyReport:   "weekly",
		},
		LastUpdated: time.Now(),
	}, nil
}

// UpdateNotificationPreferences updates user notification preferences
func (s *NotificationService) UpdateNotificationPreferences(userID string, preferences *NotificationPreferences) error {
	// Mock implementation - would update database in real system
	preferences.UserID = userID
	preferences.LastUpdated = time.Now()
	
	log.Printf("Updated notification preferences for user: %s", userID)
	return nil
}

// GetMetrics returns notification metrics
func (s *NotificationService) GetMetrics(since time.Time) (*NotificationMetrics, error) {
	// Mock implementation - would calculate from actual data
	return &NotificationMetrics{
		TotalSent:      15420,
		TotalDelivered: 14892,
		TotalFailed:    528,
		DeliveryRate:   96.6,
		ByChannel: map[NotificationChannel]int{
			ChannelRealTime: 8250,
			ChannelEmail:    4180,
			ChannelPush:     2990,
		},
		ByType: map[NotificationType]int{
			NotificationTypeMatchComplete:   6830,
			NotificationTypeRankChange:     1250,
			NotificationTypeAchievement:    2180,
			NotificationTypeCoachingTip:    1890,
			NotificationTypeSystemAlert:    340,
			NotificationTypeWeeklyReport:   980,
			NotificationTypeFriendActivity: 1950,
		},
		ByPriority: map[NotificationPriority]int{
			PriorityLow:      3580,
			PriorityNormal:   9840,
			PriorityCritical: 2000,
		},
		AverageLatency: 250 * time.Millisecond,
		Timestamp:      time.Now(),
	}, nil
}

// Worker implementations

func (s *NotificationService) startWorkers() {
	// Start notification workers
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.notificationWorker(fmt.Sprintf("notif-worker-%d", i))
	}

	// Start email workers
	s.wg.Add(1)
	go s.emailWorker()

	// Start push workers
	s.wg.Add(1)
	go s.pushWorker()
}

func (s *NotificationService) notificationWorker(workerID string) {
	defer s.wg.Done()

	for {
		select {
		case notification := <-s.notificationQueue:
			s.processNotification(notification, workerID)
		case <-s.shutdown:
			return
		}
	}
}

func (s *NotificationService) emailWorker() {
	defer s.wg.Done()

	for {
		select {
		case email := <-s.emailQueue:
			s.processEmail(email)
		case <-s.shutdown:
			return
		}
	}
}

func (s *NotificationService) pushWorker() {
	defer s.wg.Done()

	for {
		select {
		case push := <-s.pushQueue:
			s.processPush(push)
		case <-s.shutdown:
			return
		}
	}
}

func (s *NotificationService) processNotification(notification *NotificationJob, workerID string) {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ProcessingTimeout)
	defer cancel()

	notification.Status = "processing"
	
	// Process each channel
	for _, channel := range notification.Channels {
		switch channel {
		case ChannelRealTime:
			s.processRealTimeNotification(ctx, notification)
		case ChannelEmail:
			s.processEmailFromNotification(ctx, notification)
		case ChannelPush:
			s.processPushFromNotification(ctx, notification)
		case ChannelInApp:
			s.processInAppNotification(ctx, notification)
		}
	}

	notification.Status = "completed"
	log.Printf("Notification processed: %s by %s", notification.ID, workerID)
}

func (s *NotificationService) processRealTimeNotification(ctx context.Context, notification *NotificationJob) {
	data := map[string]interface{}{
		"id":       notification.ID,
		"type":     notification.Type,
		"content":  notification.Content,
		"priority": notification.Priority,
		"timestamp": time.Now(),
	}

	if err := s.SendRealTimeNotification(notification.UserID, data); err != nil {
		log.Printf("Failed to send real-time notification: %v", err)
	}
}

func (s *NotificationService) processEmailFromNotification(ctx context.Context, notification *NotificationJob) {
	if notification.Content == nil {
		return
	}

	email := &EmailNotification{
		To:       []string{notification.UserID}, // Assume userID is email for simplicity
		Subject:  notification.Content.Title,
		TextBody: notification.Content.Message,
		Priority: notification.Priority,
	}

	if err := s.SendEmailNotification(ctx, email); err != nil {
		log.Printf("Failed to send email notification: %v", err)
	}
}

func (s *NotificationService) processPushFromNotification(ctx context.Context, notification *NotificationJob) {
	if notification.Content == nil {
		return
	}

	push := &PushNotification{
		UserID:   notification.UserID,
		Title:    notification.Content.Title,
		Body:     notification.Content.Message,
		Data:     notification.Content.Data,
		Priority: notification.Priority,
	}

	if err := s.SendPushNotification(ctx, push); err != nil {
		log.Printf("Failed to send push notification: %v", err)
	}
}

func (s *NotificationService) processInAppNotification(ctx context.Context, notification *NotificationJob) {
	// Store in-app notification for later retrieval
	log.Printf("In-app notification stored for user %s: %s", notification.UserID, notification.Content.Title)
}

func (s *NotificationService) processEmail(email *EmailNotification) {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ProcessingTimeout)
	defer cancel()

	if err := s.emailProvider.SendEmail(ctx, email); err != nil {
		log.Printf("Failed to send email %s: %v", email.ID, err)
	} else {
		log.Printf("Email sent successfully: %s", email.ID)
	}
}

func (s *NotificationService) processPush(push *PushNotification) {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ProcessingTimeout)
	defer cancel()

	if err := s.pushProvider.SendPush(ctx, push); err != nil {
		log.Printf("Failed to send push notification %s: %v", push.ID, err)
	} else {
		log.Printf("Push notification sent successfully: %s", push.ID)
	}
}

// WebSocket connection handling

func (s *NotificationService) handleWebSocketConnection(wsConn *WebSocketConnection) {
	defer func() {
		wsConn.Conn.Close()
		s.UnregisterWebSocketConnection(wsConn.UserID)
	}()

	// Set connection timeouts
	wsConn.Conn.SetReadDeadline(time.Now().Add(s.config.WebSocketReadTimeout))
	wsConn.Conn.SetWriteDeadline(time.Now().Add(s.config.WebSocketWriteTimeout))

	// Start ping routine
	go s.pingWebSocket(wsConn)

	// Message sending routine
	go func() {
		for {
			select {
			case message, ok := <-wsConn.Send:
				if !ok {
					wsConn.Conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				
				wsConn.Conn.SetWriteDeadline(time.Now().Add(s.config.WebSocketWriteTimeout))
				if err := wsConn.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("WebSocket write error for user %s: %v", wsConn.UserID, err)
					return
				}
			}
		}
	}()

	// Message reading routine
	for {
		_, _, err := wsConn.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for user %s: %v", wsConn.UserID, err)
			}
			break
		}
		
		wsConn.LastPing = time.Now()
		wsConn.Conn.SetReadDeadline(time.Now().Add(s.config.WebSocketReadTimeout))
	}
}

func (s *NotificationService) pingWebSocket(wsConn *WebSocketConnection) {
	ticker := time.NewTicker(s.config.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			wsConn.Conn.SetWriteDeadline(time.Now().Add(s.config.WebSocketWriteTimeout))
			if err := wsConn.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Helper methods

func (s *NotificationService) validateNotification(notification *NotificationJob) error {
	if notification.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if len(notification.Channels) == 0 {
		return fmt.Errorf("at least one channel is required")
	}
	if notification.Content == nil {
		return fmt.Errorf("content is required")
	}
	if notification.Content.Title == "" && notification.Content.Message == "" {
		return fmt.Errorf("title or message is required")
	}
	return nil
}

func (s *NotificationService) scheduleNotification(notification *NotificationJob) error {
	// In a real implementation, this would use a scheduler like cron or a job queue
	// For now, we'll use a simple goroutine with timer
	go func() {
		timer := time.NewTimer(time.Until(*notification.ScheduledAt))
		defer timer.Stop()

		<-timer.C
		
		ctx, cancel := context.WithTimeout(context.Background(), s.config.ProcessingTimeout)
		defer cancel()

		if err := s.SendNotification(ctx, notification); err != nil {
			log.Printf("Failed to send scheduled notification %s: %v", notification.ID, err)
		}
	}()

	log.Printf("Notification scheduled: %s for %s", notification.ID, notification.ScheduledAt.Format(time.RFC3339))
	return nil
}

// Shutdown gracefully shuts down the notification service
func (s *NotificationService) Shutdown(ctx context.Context) error {
	log.Println("Shutting down notification service...")
	
	// Close all WebSocket connections
	s.connectionsMu.Lock()
	for userID, conn := range s.wsConnections {
		close(conn.Send)
		conn.Conn.Close()
		delete(s.wsConnections, userID)
	}
	s.connectionsMu.Unlock()
	
	close(s.shutdown)
	
	// Wait for workers to finish
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		log.Println("Notification service shut down successfully")
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown timeout exceeded")
	}
}

// DefaultNotificationConfig returns default configuration
func DefaultNotificationConfig() *NotificationConfig {
	return &NotificationConfig{
		QueueSize:             1000,
		WorkerCount:           3,
		ProcessingTimeout:     10 * time.Second,
		RetryAttempts:         3,
		RetryDelay:            5 * time.Second,
		WebSocketReadTimeout:  60 * time.Second,
		WebSocketWriteTimeout: 10 * time.Second,
		PingInterval:          30 * time.Second,
		MaxConnections:        1000,
		EmailEnabled:          true,
		EmailRateLimit:        100, // per minute
		PushEnabled:           true,
		PushRateLimit:         200, // per minute
		DefaultRetentionDays:  30,
		EnableBatching:        true,
		BatchSize:             50,
		BatchFlushInterval:    5 * time.Minute,
	}
}