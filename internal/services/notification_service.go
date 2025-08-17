package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"database/sql"
	"lol-match-exporter/internal/db"
)

// InsightType represents different types of insights
type InsightType string

const (
	PerformanceInsight  InsightType = "performance"
	RecommendationInsight InsightType = "recommendation"
	MMRInsight         InsightType = "mmr"
	ChampionInsight    InsightType = "champion"
	StreakInsight      InsightType = "streak"
)

// NotificationLevel represents the urgency level of a notification
type NotificationLevel string

const (
	LevelInfo    NotificationLevel = "info"
	LevelWarning NotificationLevel = "warning"
	LevelSuccess NotificationLevel = "success"
	LevelCritical NotificationLevel = "critical"
)

// Insight represents a real-time insight notification
type Insight struct {
	ID          int               `json:"id"`
	UserID      int               `json:"user_id"`
	Type        InsightType       `json:"type"`
	Level       NotificationLevel `json:"level"`
	Title       string            `json:"title"`
	Message     string            `json:"message"`
	Data        map[string]interface{} `json:"data,omitempty"`
	ActionURL   string            `json:"action_url,omitempty"`
	IsRead      bool              `json:"is_read"`
	CreatedAt   time.Time         `json:"created_at"`
	ExpiresAt   *time.Time        `json:"expires_at,omitempty"`
}

// NotificationService handles real-time insights and notifications
type NotificationService struct {
	db          *sql.DB
	subscribers map[int][]chan Insight // userID -> channels
}

// NewNotificationService creates a new notification service
func NewNotificationService(database *db.Database) *NotificationService {
	return &NotificationService{
		db:          database.DB,
		subscribers: make(map[int][]chan Insight),
	}
}

// Subscribe subscribes a user to real-time notifications
func (ns *NotificationService) Subscribe(userID int) chan Insight {
	ch := make(chan Insight, 20)
	
	if ns.subscribers[userID] == nil {
		ns.subscribers[userID] = make([]chan Insight, 0)
	}
	ns.subscribers[userID] = append(ns.subscribers[userID], ch)
	
	log.Printf("User %d subscribed to notifications", userID)
	return ch
}

// Unsubscribe removes a user subscription
func (ns *NotificationService) Unsubscribe(userID int, ch chan Insight) {
	if channels, exists := ns.subscribers[userID]; exists {
		for i, channel := range channels {
			if channel == ch {
				// Remove channel from slice
				ns.subscribers[userID] = append(channels[:i], channels[i+1:]...)
				close(ch)
				break
			}
		}
	}
}

// CreateInsight creates a new insight and notifies subscribers
func (ns *NotificationService) CreateInsight(insight Insight) error {
	// Store insight in database
	query := `
		INSERT INTO insights (
			user_id, type, level, title, message, data, action_url, expires_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING id`
	
	dataJSON, _ := json.Marshal(insight.Data)
	
	err := ns.db.QueryRow(query, insight.UserID, insight.Type, insight.Level,
		insight.Title, insight.Message, dataJSON, insight.ActionURL, insight.ExpiresAt).Scan(&insight.ID)
	if err != nil {
		return fmt.Errorf("failed to create insight: %w", err)
	}
	
	insight.CreatedAt = time.Now()
	
	// Notify subscribers
	ns.notifySubscribers(insight)
	
	log.Printf("Created insight for user %d: %s", insight.UserID, insight.Title)
	return nil
}

// notifySubscribers sends insight to all subscribers
func (ns *NotificationService) notifySubscribers(insight Insight) {
	if channels, exists := ns.subscribers[insight.UserID]; exists {
		for _, ch := range channels {
			select {
			case ch <- insight:
			default:
				// Channel is full, skip
			}
		}
	}
}

// GetUserInsights retrieves insights for a user
func (ns *NotificationService) GetUserInsights(userID int, limit int, onlyUnread bool) ([]Insight, error) {
	query := `
		SELECT id, user_id, type, level, title, message, data, action_url, is_read, created_at, expires_at
		FROM insights 
		WHERE user_id = $1`
	
	args := []interface{}{userID}
	
	if onlyUnread {
		query += " AND is_read = false"
	}
	
	query += " AND (expires_at IS NULL OR expires_at > NOW())"
	query += " ORDER BY created_at DESC"
	
	if limit > 0 {
		query += " LIMIT $2"
		args = append(args, limit)
	}
	
	rows, err := ns.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var insights []Insight
	for rows.Next() {
		var insight Insight
		var dataJSON []byte
		
		err = rows.Scan(&insight.ID, &insight.UserID, &insight.Type, &insight.Level,
			&insight.Title, &insight.Message, &dataJSON, &insight.ActionURL,
			&insight.IsRead, &insight.CreatedAt, &insight.ExpiresAt)
		if err != nil {
			continue
		}
		
		if len(dataJSON) > 0 {
			json.Unmarshal(dataJSON, &insight.Data)
		}
		
		insights = append(insights, insight)
	}
	
	return insights, nil
}

// MarkAsRead marks insights as read
func (ns *NotificationService) MarkAsRead(userID int, insightIDs []int) error {
	if len(insightIDs) == 0 {
		return nil
	}
	
	query := `UPDATE insights SET is_read = true WHERE user_id = $1 AND id = ANY($2)`
	_, err := ns.db.Exec(query, userID, insightIDs)
	return err
}

// Analytics-based insight generators

// ProcessMatchInsights generates insights after a match analysis
func (ns *NotificationService) ProcessMatchInsights(userID int, matchID string, matchData map[string]interface{}) {
	// Performance insights
	if performance, ok := matchData["performance"].(map[string]interface{}); ok {
		ns.generatePerformanceInsights(userID, matchID, performance)
	}
	
	// Streak insights
	if streak, ok := matchData["streak"].(map[string]interface{}); ok {
		ns.generateStreakInsights(userID, streak)
	}
	
	// Champion insights
	if champion, ok := matchData["champion"].(map[string]interface{}); ok {
		ns.generateChampionInsights(userID, champion)
	}
}

// generatePerformanceInsights creates performance-based insights
func (ns *NotificationService) generatePerformanceInsights(userID int, matchID string, performance map[string]interface{}) {
	// Check for significant performance improvements
	if scoreImprovement, ok := performance["score_improvement"].(float64); ok && scoreImprovement > 0.15 {
		insight := Insight{
			UserID: userID,
			Type:   PerformanceInsight,
			Level:  LevelSuccess,
			Title:  "🚀 Performance Boost!",
			Message: fmt.Sprintf("Your performance score improved by %.1f%% in your last match!", scoreImprovement*100),
			Data: map[string]interface{}{
				"match_id": matchID,
				"improvement": scoreImprovement,
			},
			ActionURL: "/analytics/performance",
		}
		ns.CreateInsight(insight)
	}
	
	// Check for concerning performance drops
	if scoreDrop, ok := performance["score_drop"].(float64); ok && scoreDrop > 0.20 {
		insight := Insight{
			UserID: userID,
			Type:   PerformanceInsight,
			Level:  LevelWarning,
			Title:  "⚠️ Performance Dip Detected",
			Message: fmt.Sprintf("Your performance dropped by %.1f%% recently. Check recommendations for improvement tips.", scoreDrop*100),
			Data: map[string]interface{}{
				"match_id": matchID,
				"drop": scoreDrop,
			},
			ActionURL: "/analytics/recommendations",
		}
		ns.CreateInsight(insight)
	}
}

// generateStreakInsights creates streak-based insights
func (ns *NotificationService) generateStreakInsights(userID int, streak map[string]interface{}) {
	if winStreak, ok := streak["win_streak"].(int); ok {
		if winStreak >= 5 {
			insight := Insight{
				UserID: userID,
				Type:   StreakInsight,
				Level:  LevelSuccess,
				Title:  "🔥 Win Streak Alert!",
				Message: fmt.Sprintf("Amazing! You're on a %d game win streak! Keep up the momentum!", winStreak),
				Data: map[string]interface{}{
					"streak_length": winStreak,
					"type": "win",
				},
				ActionURL: "/analytics/performance",
			}
			ns.CreateInsight(insight)
		}
	}
	
	if lossStreak, ok := streak["loss_streak"].(int); ok {
		if lossStreak >= 3 {
			insight := Insight{
				UserID: userID,
				Type:   StreakInsight,
				Level:  LevelWarning,
				Title:  "💪 Turn It Around",
				Message: fmt.Sprintf("You've had %d losses recently. Check our recommendations to get back on track!", lossStreak),
				Data: map[string]interface{}{
					"streak_length": lossStreak,
					"type": "loss",
				},
				ActionURL: "/analytics/recommendations",
			}
			ns.CreateInsight(insight)
		}
	}
}

// generateChampionInsights creates champion-based insights
func (ns *NotificationService) generateChampionInsights(userID int, champion map[string]interface{}) {
	if championName, ok := champion["name"].(string); ok {
		if winRate, ok := champion["win_rate"].(float64); ok && winRate >= 0.8 {
			insight := Insight{
				UserID: userID,
				Type:   ChampionInsight,
				Level:  LevelSuccess,
				Title:  "⭐ Champion Mastery!",
				Message: fmt.Sprintf("Incredible! You have a %.0f%% win rate with %s. You've mastered this champion!", winRate*100, championName),
				Data: map[string]interface{}{
					"champion": championName,
					"win_rate": winRate,
				},
				ActionURL: "/analytics/champions",
			}
			ns.CreateInsight(insight)
		}
	}
}

// ProcessMMRInsights generates MMR-based insights
func (ns *NotificationService) ProcessMMRInsights(userID int, mmrData map[string]interface{}) {
	if mmrChange, ok := mmrData["recent_change"].(float64); ok {
		if mmrChange > 50 {
			insight := Insight{
				UserID: userID,
				Type:   MMRInsight,
				Level:  LevelSuccess,
				Title:  "📈 MMR Climbing!",
				Message: fmt.Sprintf("Great progress! Your MMR increased by %.0f points recently!", mmrChange),
				Data: map[string]interface{}{
					"mmr_change": mmrChange,
				},
				ActionURL: "/analytics/mmr",
			}
			ns.CreateInsight(insight)
		} else if mmrChange < -50 {
			insight := Insight{
				UserID: userID,
				Type:   MMRInsight,
				Level:  LevelWarning,
				Title:  "📉 MMR Fluctuation",
				Message: fmt.Sprintf("Your MMR dropped by %.0f points. Don't worry, check our tips to recover!", -mmrChange),
				Data: map[string]interface{}{
					"mmr_change": mmrChange,
				},
				ActionURL: "/analytics/recommendations",
			}
			ns.CreateInsight(insight)
		}
	}
	
	// Rank promotion insights
	if promoted, ok := mmrData["rank_promotion"].(bool); ok && promoted {
		if newRank, ok := mmrData["new_rank"].(string); ok {
			insight := Insight{
				UserID: userID,
				Type:   MMRInsight,
				Level:  LevelSuccess,
				Title:  "🎉 Rank Promotion!",
				Message: fmt.Sprintf("Congratulations! You've been promoted to %s!", newRank),
				Data: map[string]interface{}{
					"new_rank": newRank,
					"promoted": true,
				},
				ActionURL: "/analytics/mmr",
				ExpiresAt: func() *time.Time { t := time.Now().Add(7 * 24 * time.Hour); return &t }(),
			}
			ns.CreateInsight(insight)
		}
	}
}

// ProcessRecommendationInsights generates insights for new recommendations
func (ns *NotificationService) ProcessRecommendationInsights(userID int, recommendations []map[string]interface{}) {
	if len(recommendations) == 0 {
		return
	}
	
	// High priority recommendations
	highPriorityCount := 0
	for _, rec := range recommendations {
		if priority, ok := rec["priority"].(int); ok && priority >= 8 {
			highPriorityCount++
		}
	}
	
	if highPriorityCount > 0 {
		insight := Insight{
			UserID: userID,
			Type:   RecommendationInsight,
			Level:  LevelInfo,
			Title:  "💡 New High-Priority Recommendations",
			Message: fmt.Sprintf("We've identified %d high-priority areas for improvement. Check them out!", highPriorityCount),
			Data: map[string]interface{}{
				"high_priority_count": highPriorityCount,
				"total_count": len(recommendations),
			},
			ActionURL: "/analytics/recommendations",
		}
		ns.CreateInsight(insight)
	}
}

// CleanupExpiredInsights removes expired insights
func (ns *NotificationService) CleanupExpiredInsights() error {
	query := `DELETE FROM insights WHERE expires_at IS NOT NULL AND expires_at < NOW()`
	_, err := ns.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired insights: %w", err)
	}
	
	log.Println("Cleaned up expired insights")
	return nil
}

// StartCleanupWorker starts a background worker to clean up expired insights
func (ns *NotificationService) StartCleanupWorker() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		
		for range ticker.C {
			err := ns.CleanupExpiredInsights()
			if err != nil {
				log.Printf("Error cleaning up insights: %v", err)
			}
		}
	}()
}