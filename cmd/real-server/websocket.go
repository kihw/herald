package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketMessage represents a message sent through WebSocket
type WebSocketMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
	UserID    int         `json:"user_id,omitempty"`
}

// WebSocketClient represents a connected client
type WebSocketClient struct {
	ID         string
	UserID     int
	Connection *websocket.Conn
	Send       chan WebSocketMessage
	Hub        *WebSocketHub
}

// WebSocketHub manages all WebSocket connections
type WebSocketHub struct {
	// Registered clients by user ID
	clients map[int]map[*WebSocketClient]bool
	
	// Channel for broadcasting messages
	broadcast chan WebSocketMessage
	
	// Channel for registering clients
	register chan *WebSocketClient
	
	// Channel for unregistering clients
	unregister chan *WebSocketClient
	
	// Mutex for thread safety
	mutex sync.RWMutex
	
	// Statistics
	stats WebSocketStats
}

// WebSocketStats tracks WebSocket performance
type WebSocketStats struct {
	TotalConnections    int64              `json:"total_connections"`
	ActiveConnections   int                `json:"active_connections"`
	MessagesSent        int64              `json:"messages_sent"`
	MessagesReceived    int64              `json:"messages_received"`
	ConnectionsByUser   map[int]int        `json:"connections_by_user"`
	LastActivity        time.Time          `json:"last_activity"`
	AverageLatency      float64            `json:"average_latency_ms"`
	ErrorCount          int64              `json:"error_count"`
}

// Global WebSocket hub
var wsHub *WebSocketHub

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from our frontend origins
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"http://localhost:5173",
			"http://localhost:5174",
			"http://localhost:3000",
			"http://localhost:8004",
		}
		
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				return true
			}
		}
		return false
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// NewWebSocketHub creates a new WebSocket hub
func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients:    make(map[int]map[*WebSocketClient]bool),
		broadcast:  make(chan WebSocketMessage, 256),
		register:   make(chan *WebSocketClient),
		unregister: make(chan *WebSocketClient),
		stats: WebSocketStats{
			ConnectionsByUser: make(map[int]int),
			LastActivity:      time.Now(),
		},
	}
}

// Run starts the WebSocket hub
func (h *WebSocketHub) Run() {
	log.Println("🔌 WebSocket hub started")
	
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
			
		case client := <-h.unregister:
			h.unregisterClient(client)
			
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient registers a new WebSocket client
func (h *WebSocketHub) registerClient(client *WebSocketClient) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	if h.clients[client.UserID] == nil {
		h.clients[client.UserID] = make(map[*WebSocketClient]bool)
	}
	
	h.clients[client.UserID][client] = true
	h.stats.TotalConnections++
	h.stats.ActiveConnections++
	h.stats.ConnectionsByUser[client.UserID]++
	h.stats.LastActivity = time.Now()
	
	log.Printf("🔌 WebSocket client registered: %s (User: %d, Total active: %d)", 
		client.ID, client.UserID, h.stats.ActiveConnections)
	
	// Send welcome message
	welcomeMsg := WebSocketMessage{
		Type:      "welcome",
		Data:      gin.H{"message": "WebSocket connected successfully", "client_id": client.ID},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	
	select {
	case client.Send <- welcomeMsg:
	default:
		close(client.Send)
		delete(h.clients[client.UserID], client)
	}
}

// unregisterClient unregisters a WebSocket client
func (h *WebSocketHub) unregisterClient(client *WebSocketClient) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	if clients, ok := h.clients[client.UserID]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.Send)
			
			h.stats.ActiveConnections--
			h.stats.ConnectionsByUser[client.UserID]--
			
			if h.stats.ConnectionsByUser[client.UserID] <= 0 {
				delete(h.stats.ConnectionsByUser, client.UserID)
			}
			
			if len(clients) == 0 {
				delete(h.clients, client.UserID)
			}
			
			log.Printf("🔌 WebSocket client disconnected: %s (User: %d, Total active: %d)", 
				client.ID, client.UserID, h.stats.ActiveConnections)
		}
	}
}

// broadcastMessage sends a message to relevant clients
func (h *WebSocketHub) broadcastMessage(message WebSocketMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	h.stats.MessagesSent++
	h.stats.LastActivity = time.Now()
	
	// If message has a specific user ID, send only to that user
	if message.UserID > 0 {
		if clients, ok := h.clients[message.UserID]; ok {
			for client := range clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(clients, client)
				}
			}
		}
		return
	}
	
	// Broadcast to all connected clients
	for userID, clients := range h.clients {
		for client := range clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(clients, client)
				h.stats.ActiveConnections--
				h.stats.ConnectionsByUser[userID]--
			}
		}
	}
}

// BroadcastToUser sends a message to a specific user
func (h *WebSocketHub) BroadcastToUser(userID int, msgType string, data interface{}) {
	message := WebSocketMessage{
		Type:      msgType,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
		UserID:    userID,
	}
	
	select {
	case h.broadcast <- message:
	default:
		h.stats.ErrorCount++
		log.Printf("⚠️ Failed to queue WebSocket message for user %d", userID)
	}
}

// BroadcastToAll sends a message to all connected clients
func (h *WebSocketHub) BroadcastToAll(msgType string, data interface{}) {
	message := WebSocketMessage{
		Type:      msgType,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	
	select {
	case h.broadcast <- message:
	default:
		h.stats.ErrorCount++
		log.Printf("⚠️ Failed to queue WebSocket broadcast message")
	}
}

// GetStats returns WebSocket statistics
func (h *WebSocketHub) GetStats() WebSocketStats {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	stats := h.stats
	stats.ConnectionsByUser = make(map[int]int)
	for userID, count := range h.stats.ConnectionsByUser {
		stats.ConnectionsByUser[userID] = count
	}
	
	return stats
}

// writeMessage handles writing messages to WebSocket connection
func (c *WebSocketClient) writeMessage() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Connection.Close()
	}()
	
	for {
		select {
		case message, ok := <-c.Send:
			c.Connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
			
			if !ok {
				c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			if err := c.Connection.WriteJSON(message); err != nil {
				log.Printf("⚠️ WebSocket write error for client %s: %v", c.ID, err)
				c.Hub.stats.ErrorCount++
				return
			}
			
		case <-ticker.C:
			c.Connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("⚠️ WebSocket ping error for client %s: %v", c.ID, err)
				return
			}
		}
	}
}

// readMessage handles reading messages from WebSocket connection
func (c *WebSocketClient) readMessage() {
	defer func() {
		c.Hub.unregister <- c
		c.Connection.Close()
	}()
	
	c.Connection.SetReadLimit(512)
	c.Connection.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Connection.SetPongHandler(func(string) error {
		c.Connection.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	
	for {
		var message WebSocketMessage
		err := c.Connection.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("⚠️ WebSocket read error for client %s: %v", c.ID, err)
				c.Hub.stats.ErrorCount++
			}
			break
		}
		
		c.Hub.stats.MessagesReceived++
		c.Hub.stats.LastActivity = time.Now()
		
		// Handle incoming messages (echo back for now)
		log.Printf("📨 Received WebSocket message from client %s: %s", c.ID, message.Type)
		
		// Echo the message back to sender
		response := WebSocketMessage{
			Type:      "echo",
			Data:      gin.H{"original": message, "echo": true},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		
		select {
		case c.Send <- response:
		default:
			close(c.Send)
			return
		}
	}
}

// InitializeWebSocket initializes the WebSocket system
func InitializeWebSocket() {
	wsHub = NewWebSocketHub()
	go wsHub.Run()
	log.Println("🚀 WebSocket system initialized")
}

// NotifyMatchSync sends real-time notification when matches are synced
func NotifyMatchSync(userID, newMatches, totalMatches int) {
	if wsHub == nil {
		return
	}
	
	data := gin.H{
		"new_matches":   newMatches,
		"total_matches": totalMatches,
		"synced_at":     time.Now().Format(time.RFC3339),
		"message":       "New matches synchronized successfully!",
	}
	
	wsHub.BroadcastToUser(userID, "match_sync", data)
	log.Printf("📡 WebSocket notification sent to user %d: %d new matches", userID, newMatches)
}

// NotifyStatsUpdate sends real-time notification when stats are updated
func NotifyStatsUpdate(userID int, stats interface{}) {
	if wsHub == nil {
		return
	}
	
	data := gin.H{
		"stats":      stats,
		"updated_at": time.Now().Format(time.RFC3339),
		"message":    "Statistics updated",
	}
	
	wsHub.BroadcastToUser(userID, "stats_update", data)
	log.Printf("📡 WebSocket stats update sent to user %d", userID)
}

// NotifySystemUpdate sends system-wide notifications
func NotifySystemUpdate(message string, data interface{}) {
	if wsHub == nil {
		return
	}
	
	payload := gin.H{
		"message":    message,
		"data":       data,
		"updated_at": time.Now().Format(time.RFC3339),
	}
	
	wsHub.BroadcastToAll("system_update", payload)
	log.Printf("📡 System-wide WebSocket notification sent: %s", message)
}