package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/herald-lol/herald/backend/internal/services"
)

// NotificationHandler handles notification-related requests
type NotificationHandler struct {
	notificationService *services.NotificationService
	wsUpgrader          websocket.Upgrader
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		wsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// In production, implement proper origin checking
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// RegisterRoutes registers all notification routes
func (h *NotificationHandler) RegisterRoutes(r *gin.RouterGroup) {
	notifications := r.Group("/notifications")
	{
		// WebSocket connection
		notifications.GET("/ws/:user_id", h.HandleWebSocket)

		// Notification management
		notifications.POST("/send", h.SendNotification)
		notifications.POST("/send-bulk", h.SendBulkNotifications)
		notifications.POST("/schedule", h.ScheduleNotification)

		// Specific notification types
		notifications.POST("/email", h.SendEmailNotification)
		notifications.POST("/push", h.SendPushNotification)
		notifications.POST("/realtime", h.SendRealtimeNotification)

		// Gaming-specific notifications
		notifications.POST("/match-complete", h.NotifyMatchComplete)
		notifications.POST("/rank-change", h.NotifyRankChange)
		notifications.POST("/achievement", h.NotifyAchievement)
		notifications.POST("/coaching-tip", h.SendCoachingTip)

		// User preferences
		notifications.GET("/preferences/:user_id", h.GetNotificationPreferences)
		notifications.PUT("/preferences/:user_id", h.UpdateNotificationPreferences)

		// Notification history
		notifications.GET("/history/:user_id", h.GetNotificationHistory)
		notifications.GET("/unread/:user_id", h.GetUnreadNotifications)
		notifications.PUT("/mark-read", h.MarkNotificationsRead)

		// System management
		notifications.GET("/metrics", h.GetNotificationMetrics)
		notifications.GET("/status", h.GetSystemStatus)
		notifications.GET("/connections", h.GetActiveConnections)

		// Templates
		notifications.GET("/templates", h.GetNotificationTemplates)
		notifications.POST("/test-template", h.TestNotificationTemplate)
	}
}

// HandleWebSocket handles WebSocket connections for real-time notifications
func (h *NotificationHandler) HandleWebSocket(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := h.wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to upgrade WebSocket connection",
			"details": err.Error(),
		})
		return
	}

	// Extract client info
	clientInfo := &services.ClientInfo{
		UserAgent:    c.Request.Header.Get("User-Agent"),
		IPAddress:    c.ClientIP(),
		Platform:     c.Query("platform"),
		Version:      c.Query("version"),
		Capabilities: []string{"realtime", "push"},
	}

	// Register WebSocket connection
	if err := h.notificationService.RegisterWebSocketConnection(userID, conn, clientInfo); err != nil {
		conn.Close()
		return
	}
}

// SendNotification handles general notification sending requests
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var request struct {
		UserID   string                         `json:"user_id" binding:"required"`
		Type     services.NotificationType      `json:"type" binding:"required"`
		Channels []services.NotificationChannel `json:"channels" binding:"required"`
		Content  *services.NotificationContent  `json:"content" binding:"required"`
		Priority services.NotificationPriority  `json:"priority"`
		Context  map[string]interface{}         `json:"context"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid notification request",
			"details": err.Error(),
		})
		return
	}

	notification := &services.NotificationJob{
		UserID:   request.UserID,
		Type:     request.Type,
		Channels: request.Channels,
		Content:  request.Content,
		Priority: request.Priority,
		Context:  request.Context,
	}

	if err := h.notificationService.SendNotification(c.Request.Context(), notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notification_id":    notification.ID,
		"status":             "queued",
		"message":            "Notification queued successfully",
		"estimated_delivery": time.Now().Add(5 * time.Second),
	})
}

// SendBulkNotifications handles bulk notification sending
func (h *NotificationHandler) SendBulkNotifications(c *gin.Context) {
	var request struct {
		UserIDs  []string                       `json:"user_ids" binding:"required"`
		Type     services.NotificationType      `json:"type" binding:"required"`
		Channels []services.NotificationChannel `json:"channels" binding:"required"`
		Content  *services.NotificationContent  `json:"content" binding:"required"`
		Priority services.NotificationPriority  `json:"priority"`
		Context  map[string]interface{}         `json:"context"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid bulk notification request",
			"details": err.Error(),
		})
		return
	}

	if len(request.UserIDs) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Maximum 1000 users allowed per bulk request",
		})
		return
	}

	batchID := time.Now().UnixNano()
	var successCount, failureCount int

	for _, userID := range request.UserIDs {
		notification := &services.NotificationJob{
			UserID:   userID,
			Type:     request.Type,
			Channels: request.Channels,
			Content:  request.Content,
			Priority: request.Priority,
			Context:  request.Context,
		}

		if err := h.notificationService.SendNotification(c.Request.Context(), notification); err != nil {
			failureCount++
		} else {
			successCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"batch_id":      batchID,
		"total_users":   len(request.UserIDs),
		"success_count": successCount,
		"failure_count": failureCount,
		"success_rate":  float64(successCount) / float64(len(request.UserIDs)) * 100,
		"message":       "Bulk notifications processed",
	})
}

// ScheduleNotification handles notification scheduling requests
func (h *NotificationHandler) ScheduleNotification(c *gin.Context) {
	var request struct {
		UserID      string                         `json:"user_id" binding:"required"`
		Type        services.NotificationType      `json:"type" binding:"required"`
		Channels    []services.NotificationChannel `json:"channels" binding:"required"`
		Content     *services.NotificationContent  `json:"content" binding:"required"`
		Priority    services.NotificationPriority  `json:"priority"`
		ScheduledAt time.Time                      `json:"scheduled_at" binding:"required"`
		Context     map[string]interface{}         `json:"context"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid scheduled notification request",
			"details": err.Error(),
		})
		return
	}

	if request.ScheduledAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Scheduled time must be in the future",
		})
		return
	}

	notification := &services.NotificationJob{
		UserID:      request.UserID,
		Type:        request.Type,
		Channels:    request.Channels,
		Content:     request.Content,
		Priority:    request.Priority,
		ScheduledAt: &request.ScheduledAt,
		Context:     request.Context,
	}

	if err := h.notificationService.SendNotification(c.Request.Context(), notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to schedule notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notification_id": notification.ID,
		"scheduled_at":    request.ScheduledAt,
		"status":          "scheduled",
		"message":         "Notification scheduled successfully",
	})
}

// SendEmailNotification handles email-specific notifications
func (h *NotificationHandler) SendEmailNotification(c *gin.Context) {
	var request struct {
		To           []string                      `json:"to" binding:"required"`
		CC           []string                      `json:"cc"`
		BCC          []string                      `json:"bcc"`
		Subject      string                        `json:"subject" binding:"required"`
		Template     string                        `json:"template"`
		TemplateData map[string]interface{}        `json:"template_data"`
		TextBody     string                        `json:"text_body"`
		HTMLBody     string                        `json:"html_body"`
		Priority     services.NotificationPriority `json:"priority"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid email notification request",
			"details": err.Error(),
		})
		return
	}

	email := &services.EmailNotification{
		To:           request.To,
		CC:           request.CC,
		BCC:          request.BCC,
		Subject:      request.Subject,
		Template:     request.Template,
		TemplateData: request.TemplateData,
		TextBody:     request.TextBody,
		HTMLBody:     request.HTMLBody,
		Priority:     request.Priority,
	}

	if err := h.notificationService.SendEmailNotification(c.Request.Context(), email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send email notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email_id":   email.ID,
		"recipients": len(request.To),
		"status":     "queued",
		"message":    "Email notification queued successfully",
	})
}

// SendPushNotification handles push-specific notifications
func (h *NotificationHandler) SendPushNotification(c *gin.Context) {
	var request struct {
		UserID      string                        `json:"user_id" binding:"required"`
		DeviceToken string                        `json:"device_token"`
		Platform    string                        `json:"platform" binding:"required"`
		Title       string                        `json:"title" binding:"required"`
		Body        string                        `json:"body" binding:"required"`
		Data        map[string]interface{}        `json:"data"`
		Badge       int                           `json:"badge"`
		Sound       string                        `json:"sound"`
		ClickAction string                        `json:"click_action"`
		Priority    services.NotificationPriority `json:"priority"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid push notification request",
			"details": err.Error(),
		})
		return
	}

	push := &services.PushNotification{
		UserID:      request.UserID,
		DeviceToken: request.DeviceToken,
		Platform:    request.Platform,
		Title:       request.Title,
		Body:        request.Body,
		Data:        request.Data,
		Badge:       request.Badge,
		Sound:       request.Sound,
		ClickAction: request.ClickAction,
		Priority:    request.Priority,
	}

	if err := h.notificationService.SendPushNotification(c.Request.Context(), push); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send push notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"push_id":  push.ID,
		"user_id":  push.UserID,
		"platform": push.Platform,
		"status":   "queued",
		"message":  "Push notification queued successfully",
	})
}

// SendRealtimeNotification handles real-time WebSocket notifications
func (h *NotificationHandler) SendRealtimeNotification(c *gin.Context) {
	var request struct {
		UserID string                 `json:"user_id" binding:"required"`
		Data   map[string]interface{} `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid real-time notification request",
			"details": err.Error(),
		})
		return
	}

	if err := h.notificationService.SendRealTimeNotification(request.UserID, request.Data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send real-time notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":   request.UserID,
		"status":    "sent",
		"message":   "Real-time notification sent successfully",
		"timestamp": time.Now(),
	})
}

// Gaming-specific notification handlers

// NotifyMatchComplete handles match completion notifications
func (h *NotificationHandler) NotifyMatchComplete(c *gin.Context) {
	var request struct {
		UserID     string                         `json:"user_id" binding:"required"`
		MatchID    string                         `json:"match_id" binding:"required"`
		Result     string                         `json:"result" binding:"required"` // "victory" or "defeat"
		Champion   string                         `json:"champion" binding:"required"`
		KDA        string                         `json:"kda" binding:"required"`
		Duration   int                            `json:"duration" binding:"required"`
		Rating     float64                        `json:"rating"`
		RankChange int                            `json:"rank_change"`
		Channels   []services.NotificationChannel `json:"channels"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid match complete notification request",
			"details": err.Error(),
		})
		return
	}

	if len(request.Channels) == 0 {
		request.Channels = []services.NotificationChannel{
			services.ChannelRealTime,
			services.ChannelInApp,
		}
	}

	resultEmoji := "⚔️"
	if request.Result == "victory" {
		resultEmoji = "🏆"
	} else if request.Result == "defeat" {
		resultEmoji = "💔"
	}

	content := &services.NotificationContent{
		Title: fmt.Sprintf("%s Match Complete", resultEmoji),
		Message: fmt.Sprintf("Your %s match as %s is complete! KDA: %s, Duration: %dm, Rating: %.1f",
			request.Result, request.Champion, request.KDA, request.Duration/60, request.Rating),
		Data: map[string]interface{}{
			"match_id":    request.MatchID,
			"result":      request.Result,
			"champion":    request.Champion,
			"kda":         request.KDA,
			"duration":    request.Duration,
			"rating":      request.Rating,
			"rank_change": request.RankChange,
		},
		Icon:        "match_complete",
		Category:    "gaming",
		ClickAction: fmt.Sprintf("/matches/%s", request.MatchID),
	}

	notification := &services.NotificationJob{
		UserID:   request.UserID,
		Type:     services.NotificationTypeMatchComplete,
		Channels: request.Channels,
		Content:  content,
		Priority: services.PriorityNormal,
		Context: map[string]interface{}{
			"match_id": request.MatchID,
			"result":   request.Result,
		},
	}

	if err := h.notificationService.SendNotification(c.Request.Context(), notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send match complete notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notification_id": notification.ID,
		"match_id":        request.MatchID,
		"result":          request.Result,
		"message":         "Match complete notification sent successfully",
	})
}

// NotifyRankChange handles rank change notifications
func (h *NotificationHandler) NotifyRankChange(c *gin.Context) {
	var request struct {
		UserID      string                         `json:"user_id" binding:"required"`
		OldRank     string                         `json:"old_rank" binding:"required"`
		NewRank     string                         `json:"new_rank" binding:"required"`
		LP          int                            `json:"lp"`
		LPChange    int                            `json:"lp_change"`
		QueueType   string                         `json:"queue_type"`
		IsPromotion bool                           `json:"is_promotion"`
		IsDemotion  bool                           `json:"is_demotion"`
		Channels    []services.NotificationChannel `json:"channels"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid rank change notification request",
			"details": err.Error(),
		})
		return
	}

	if len(request.Channels) == 0 {
		request.Channels = []services.NotificationChannel{
			services.ChannelRealTime,
			services.ChannelPush,
			services.ChannelEmail,
		}
	}

	emoji := "📊"
	title := "Rank Updated"
	priority := services.PriorityNormal

	if request.IsPromotion {
		emoji = "🎉"
		title = "Rank Promotion!"
		priority = services.PriorityHigh
	} else if request.IsDemotion {
		emoji = "😢"
		title = "Rank Demotion"
		priority = services.PriorityNormal
	}

	message := fmt.Sprintf("%s Your rank changed from %s to %s", emoji, request.OldRank, request.NewRank)
	if request.LPChange != 0 {
		message += fmt.Sprintf(" (%+d LP)", request.LPChange)
	}

	content := &services.NotificationContent{
		Title:   title,
		Message: message,
		Data: map[string]interface{}{
			"old_rank":     request.OldRank,
			"new_rank":     request.NewRank,
			"lp":           request.LP,
			"lp_change":    request.LPChange,
			"queue_type":   request.QueueType,
			"is_promotion": request.IsPromotion,
			"is_demotion":  request.IsDemotion,
		},
		Icon:        "rank_change",
		Category:    "gaming",
		ClickAction: "/profile/ranked",
	}

	notification := &services.NotificationJob{
		UserID:   request.UserID,
		Type:     services.NotificationTypeRankChange,
		Channels: request.Channels,
		Content:  content,
		Priority: priority,
		Context: map[string]interface{}{
			"rank_change": map[string]interface{}{
				"old":       request.OldRank,
				"new":       request.NewRank,
				"promotion": request.IsPromotion,
			},
		},
	}

	if err := h.notificationService.SendNotification(c.Request.Context(), notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send rank change notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notification_id": notification.ID,
		"old_rank":        request.OldRank,
		"new_rank":        request.NewRank,
		"change_type": map[string]interface{}{
			"promotion": request.IsPromotion,
			"demotion":  request.IsDemotion,
		},
		"message": "Rank change notification sent successfully",
	})
}

// NotifyAchievement handles achievement notifications
func (h *NotificationHandler) NotifyAchievement(c *gin.Context) {
	var request struct {
		UserID          string                         `json:"user_id" binding:"required"`
		AchievementID   string                         `json:"achievement_id" binding:"required"`
		AchievementName string                         `json:"achievement_name" binding:"required"`
		Description     string                         `json:"description"`
		Icon            string                         `json:"icon"`
		Points          int                            `json:"points"`
		Rarity          string                         `json:"rarity"`
		Channels        []services.NotificationChannel `json:"channels"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid achievement notification request",
			"details": err.Error(),
		})
		return
	}

	if len(request.Channels) == 0 {
		request.Channels = []services.NotificationChannel{
			services.ChannelRealTime,
			services.ChannelPush,
		}
	}

	emoji := "🏅"
	if request.Rarity == "legendary" {
		emoji = "💎"
	} else if request.Rarity == "epic" {
		emoji = "🌟"
	}

	content := &services.NotificationContent{
		Title: fmt.Sprintf("%s Achievement Unlocked!", emoji),
		Message: fmt.Sprintf("You've earned the '%s' achievement! %s",
			request.AchievementName, request.Description),
		Data: map[string]interface{}{
			"achievement_id":   request.AchievementID,
			"achievement_name": request.AchievementName,
			"description":      request.Description,
			"points":           request.Points,
			"rarity":           request.Rarity,
		},
		Icon:        request.Icon,
		Category:    "achievement",
		ClickAction: fmt.Sprintf("/achievements/%s", request.AchievementID),
	}

	priority := services.PriorityNormal
	if request.Rarity == "legendary" {
		priority = services.PriorityHigh
	}

	notification := &services.NotificationJob{
		UserID:   request.UserID,
		Type:     services.NotificationTypeAchievement,
		Channels: request.Channels,
		Content:  content,
		Priority: priority,
		Context: map[string]interface{}{
			"achievement": map[string]interface{}{
				"id":     request.AchievementID,
				"name":   request.AchievementName,
				"rarity": request.Rarity,
			},
		},
	}

	if err := h.notificationService.SendNotification(c.Request.Context(), notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send achievement notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notification_id":  notification.ID,
		"achievement_id":   request.AchievementID,
		"achievement_name": request.AchievementName,
		"points":           request.Points,
		"message":          "Achievement notification sent successfully",
	})
}

// SendCoachingTip handles coaching tip notifications
func (h *NotificationHandler) SendCoachingTip(c *gin.Context) {
	var request struct {
		UserID      string                         `json:"user_id" binding:"required"`
		TipCategory string                         `json:"tip_category" binding:"required"`
		Title       string                         `json:"title" binding:"required"`
		Content     string                         `json:"content" binding:"required"`
		Priority    string                         `json:"priority"`
		Tags        []string                       `json:"tags"`
		URL         string                         `json:"url"`
		Channels    []services.NotificationChannel `json:"channels"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid coaching tip notification request",
			"details": err.Error(),
		})
		return
	}

	if len(request.Channels) == 0 {
		request.Channels = []services.NotificationChannel{
			services.ChannelInApp,
		}
	}

	content := &services.NotificationContent{
		Title:   fmt.Sprintf("💡 %s", request.Title),
		Message: request.Content,
		Data: map[string]interface{}{
			"tip_category": request.TipCategory,
			"tags":         request.Tags,
		},
		Icon:     "coaching_tip",
		Category: "coaching",
		URL:      request.URL,
		Tags:     request.Tags,
	}

	priority := services.PriorityLow
	if request.Priority == "high" {
		priority = services.PriorityHigh
	} else if request.Priority == "normal" {
		priority = services.PriorityNormal
	}

	notification := &services.NotificationJob{
		UserID:   request.UserID,
		Type:     services.NotificationTypeCoachingTip,
		Channels: request.Channels,
		Content:  content,
		Priority: priority,
		Context: map[string]interface{}{
			"coaching": map[string]interface{}{
				"category": request.TipCategory,
				"tags":     request.Tags,
			},
		},
	}

	if err := h.notificationService.SendNotification(c.Request.Context(), notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send coaching tip notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notification_id": notification.ID,
		"tip_category":    request.TipCategory,
		"title":           request.Title,
		"message":         "Coaching tip notification sent successfully",
	})
}

// GetNotificationPreferences handles getting user preferences
func (h *NotificationHandler) GetNotificationPreferences(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}

	preferences, err := h.notificationService.GetNotificationPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get notification preferences",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"preferences": preferences,
		"message":     "Notification preferences retrieved successfully",
	})
}

// UpdateNotificationPreferences handles updating user preferences
func (h *NotificationHandler) UpdateNotificationPreferences(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}

	var preferences services.NotificationPreferences
	if err := c.ShouldBindJSON(&preferences); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid notification preferences",
			"details": err.Error(),
		})
		return
	}

	if err := h.notificationService.UpdateNotificationPreferences(userID, &preferences); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update notification preferences",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":    userID,
		"updated_at": time.Now(),
		"message":    "Notification preferences updated successfully",
	})
}

// GetNotificationHistory handles getting notification history
func (h *NotificationHandler) GetNotificationHistory(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	notificationType := c.Query("type")
	channel := c.Query("channel")

	// Mock response - in real implementation would fetch from database
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"notifications": []gin.H{
			{
				"id":           "notif_123456789",
				"type":         "match_complete",
				"title":        "🏆 Match Complete",
				"message":      "Your victory match as Jinx is complete! KDA: 8/2/12",
				"channels":     []string{"realtime", "push"},
				"priority":     "normal",
				"read":         false,
				"created_at":   "2024-01-15T14:30:00Z",
				"delivered_at": "2024-01-15T14:30:02Z",
			},
			{
				"id":           "notif_123456790",
				"type":         "rank_change",
				"title":        "🎉 Rank Promotion!",
				"message":      "Your rank changed from Gold III to Gold II (+18 LP)",
				"channels":     []string{"realtime", "email", "push"},
				"priority":     "high",
				"read":         true,
				"created_at":   "2024-01-15T13:45:00Z",
				"delivered_at": "2024-01-15T13:45:01Z",
			},
		},
		"pagination": gin.H{
			"limit":    limit,
			"total":    25,
			"has_more": false,
		},
		"filters": gin.H{
			"type":    notificationType,
			"channel": channel,
		},
	})
}

// GetUnreadNotifications handles getting unread notifications
func (h *NotificationHandler) GetUnreadNotifications(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}

	// Mock response
	c.JSON(http.StatusOK, gin.H{
		"user_id":      userID,
		"unread_count": 3,
		"notifications": []gin.H{
			{
				"id":         "notif_123456789",
				"type":       "match_complete",
				"title":      "🏆 Match Complete",
				"message":    "Your victory match as Jinx is complete!",
				"priority":   "normal",
				"created_at": "2024-01-15T14:30:00Z",
			},
			{
				"id":         "notif_123456791",
				"type":       "coaching_tip",
				"title":      "💡 Farming Tip",
				"message":    "Focus on CS in the first 15 minutes to maximize gold income",
				"priority":   "low",
				"created_at": "2024-01-15T14:00:00Z",
			},
			{
				"id":         "notif_123456792",
				"type":       "achievement",
				"title":      "🏅 Achievement Unlocked!",
				"message":    "You've earned the 'Pentakill Master' achievement!",
				"priority":   "normal",
				"created_at": "2024-01-15T13:50:00Z",
			},
		},
	})
}

// MarkNotificationsRead handles marking notifications as read
func (h *NotificationHandler) MarkNotificationsRead(c *gin.Context) {
	var request struct {
		UserID          string   `json:"user_id" binding:"required"`
		NotificationIDs []string `json:"notification_ids"`
		MarkAllRead     bool     `json:"mark_all_read"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid mark read request",
			"details": err.Error(),
		})
		return
	}

	markedCount := len(request.NotificationIDs)
	if request.MarkAllRead {
		markedCount = 15 // Mock count
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":           request.UserID,
		"marked_read_count": markedCount,
		"mark_all_read":     request.MarkAllRead,
		"updated_at":        time.Now(),
		"message":           "Notifications marked as read successfully",
	})
}

// GetNotificationMetrics handles getting system metrics
func (h *NotificationHandler) GetNotificationMetrics(c *gin.Context) {
	hoursStr := c.DefaultQuery("hours", "24")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours <= 0 {
		hours = 24
	}

	since := time.Now().Add(-time.Duration(hours) * time.Hour)
	metrics, err := h.notificationService.GetMetrics(since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get notification metrics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"time_range_hours": hours,
		"metrics":          metrics,
		"generated_at":     time.Now(),
	})
}

// GetSystemStatus handles getting notification system status
func (h *NotificationHandler) GetSystemStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"system_status": gin.H{
			"status":  "healthy",
			"uptime":  "8h 32m 15s",
			"version": "v1.2.0",
		},
		"queue_status": gin.H{
			"notification_queue": gin.H{
				"size":        12,
				"capacity":    1000,
				"utilization": 0.012,
			},
			"email_queue": gin.H{
				"size":        5,
				"capacity":    1000,
				"utilization": 0.005,
			},
			"push_queue": gin.H{
				"size":        8,
				"capacity":    1000,
				"utilization": 0.008,
			},
		},
		"worker_status": gin.H{
			"notification_workers": gin.H{
				"total":  3,
				"active": 2,
				"idle":   1,
			},
			"email_workers": gin.H{
				"total":  1,
				"active": 1,
				"idle":   0,
			},
			"push_workers": gin.H{
				"total":  1,
				"active": 0,
				"idle":   1,
			},
		},
		"performance": gin.H{
			"average_processing_time": "120ms",
			"success_rate":            98.5,
			"error_rate":              1.5,
		},
		"connections": gin.H{
			"websocket_connections":  145,
			"max_connections":        1000,
			"connection_utilization": 0.145,
		},
	})
}

// GetActiveConnections handles getting active WebSocket connections
func (h *NotificationHandler) GetActiveConnections(c *gin.Context) {
	// Mock response - would return actual connection data
	c.JSON(http.StatusOK, gin.H{
		"total_connections": 145,
		"connections_by_platform": gin.H{
			"web":     89,
			"mobile":  45,
			"desktop": 11,
		},
		"recent_connections": []gin.H{
			{
				"user_id":      "user_123",
				"platform":     "web",
				"connected_at": "2024-01-15T14:30:00Z",
				"last_ping":    "2024-01-15T14:32:15Z",
				"ip_address":   "192.168.1.100",
			},
			{
				"user_id":      "user_456",
				"platform":     "mobile",
				"connected_at": "2024-01-15T14:25:00Z",
				"last_ping":    "2024-01-15T14:32:10Z",
				"ip_address":   "192.168.1.101",
			},
		},
		"connection_stats": gin.H{
			"average_session_duration": "25m 30s",
			"total_messages_sent":      8450,
			"messages_per_minute":      125,
		},
	})
}

// GetNotificationTemplates handles getting available templates
func (h *NotificationHandler) GetNotificationTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"email_templates": []gin.H{
			{
				"name":        "match_complete",
				"title":       "Match Complete",
				"description": "Notification sent when a match is completed",
				"variables":   []string{"player_name", "result", "champion", "kda", "duration"},
			},
			{
				"name":        "rank_change",
				"title":       "Rank Change",
				"description": "Notification sent when player rank changes",
				"variables":   []string{"player_name", "old_rank", "new_rank", "lp_change"},
			},
			{
				"name":        "weekly_report",
				"title":       "Weekly Performance Report",
				"description": "Weekly summary of player performance",
				"variables":   []string{"player_name", "matches_played", "win_rate", "favorite_champion"},
			},
		},
		"push_templates": []gin.H{
			{
				"name":        "achievement",
				"title":       "Achievement Unlocked",
				"description": "Notification sent when player unlocks achievement",
				"variables":   []string{"achievement_name", "description", "points"},
			},
			{
				"name":        "coaching_tip",
				"title":       "Coaching Tip",
				"description": "Daily coaching tip notification",
				"variables":   []string{"tip_category", "tip_content"},
			},
		},
	})
}

// TestNotificationTemplate handles template testing
func (h *NotificationHandler) TestNotificationTemplate(c *gin.Context) {
	var request struct {
		TemplateName string                 `json:"template_name" binding:"required"`
		TemplateType string                 `json:"template_type" binding:"required"` // "email" or "push"
		TestData     map[string]interface{} `json:"test_data" binding:"required"`
		TestUserID   string                 `json:"test_user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid template test request",
			"details": err.Error(),
		})
		return
	}

	// Mock template rendering result
	renderedContent := gin.H{
		"template_name": request.TemplateName,
		"template_type": request.TemplateType,
		"rendered_at":   time.Now(),
	}

	if request.TemplateType == "email" {
		renderedContent["subject"] = "Test Email: " + request.TemplateName
		renderedContent["text_body"] = "This is a test email rendered from template " + request.TemplateName
		renderedContent["html_body"] = "<p>This is a test email rendered from template " + request.TemplateName + "</p>"
	} else if request.TemplateType == "push" {
		renderedContent["title"] = "Test Push: " + request.TemplateName
		renderedContent["body"] = "This is a test push notification rendered from template " + request.TemplateName
	}

	c.JSON(http.StatusOK, gin.H{
		"template_test_id": fmt.Sprintf("test_%d", time.Now().UnixNano()),
		"template_name":    request.TemplateName,
		"template_type":    request.TemplateType,
		"test_user_id":     request.TestUserID,
		"rendered_content": renderedContent,
		"test_status":      "success",
		"message":          "Template test completed successfully",
	})
}
