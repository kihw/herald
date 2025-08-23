package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Hub maintains active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Room-based subscriptions for gaming features
	rooms map[string]map[*Client]bool

	// User-specific subscriptions
	userSubs map[string]map[*Client]bool

	// Match-specific subscriptions
	matchSubs map[string]map[*Client]bool

	// Statistics
	stats HubStats

	mu sync.RWMutex
}

// HubStats tracks WebSocket connection statistics
type HubStats struct {
	TotalConnections    int64     `json:"total_connections"`
	ActiveConnections   int       `json:"active_connections"`
	MessagesPerSecond   float64   `json:"messages_per_second"`
	LastMessageTime     time.Time `json:"last_message_time"`
	AverageResponseTime float64   `json:"average_response_time_ms"`
}

// Message types for Herald.lol gaming platform
const (
	MessageTypeMatchUpdate     = "match_update"
	MessageTypePerformanceUpdate = "performance_update"
	MessageTypeRankUpdate      = "rank_update"
	MessageTypeFriendActivity  = "friend_activity"
	MessageTypeLiveMatch       = "live_match"
	MessageTypeCoachingSuggestion = "coaching_suggestion"
	MessageTypeChampionMastery = "champion_mastery"
	MessageTypeSystemNotification = "system_notification"
	MessageTypeError           = "error"
	MessageTypePing            = "ping"
	MessageTypePong            = "pong"
)

// Message represents a WebSocket message
type Message struct {
	Type      string      `json:"type"`
	UserID    string      `json:"user_id,omitempty"`
	MatchID   string      `json:"match_id,omitempty"`
	RoomID    string      `json:"room_id,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	ID        string      `json:"id"`
}

// Gaming-specific message data structures
type MatchUpdateData struct {
	GameID        string            `json:"game_id"`
	Status        string            `json:"status"`
	GameTime      int               `json:"game_time"`
	Participants  []ParticipantData `json:"participants"`
	TeamStats     TeamStatsData     `json:"team_stats"`
	EventData     interface{}       `json:"event_data,omitempty"`
}

type ParticipantData struct {
	SummonerName string  `json:"summoner_name"`
	ChampionName string  `json:"champion_name"`
	Level        int     `json:"level"`
	Kills        int     `json:"kills"`
	Deaths       int     `json:"deaths"`
	Assists      int     `json:"assists"`
	CS           int     `json:"cs"`
	Gold         int     `json:"gold"`
	Items        []int   `json:"items"`
	KDA          float64 `json:"kda"`
}

type TeamStatsData struct {
	BlueTeam TeamData `json:"blue_team"`
	RedTeam  TeamData `json:"red_team"`
}

type TeamData struct {
	Kills           int `json:"kills"`
	Deaths          int `json:"deaths"`
	Assists         int `json:"assists"`
	Gold            int `json:"gold"`
	Dragons         int `json:"dragons"`
	Barons          int `json:"barons"`
	Towers          int `json:"towers"`
	Inhibitors      int `json:"inhibitors"`
}

type PerformanceUpdateData struct {
	UserID        string  `json:"user_id"`
	CurrentKDA    float64 `json:"current_kda"`
	AverageKDA    float64 `json:"average_kda"`
	CSPerMinute   float64 `json:"cs_per_minute"`
	VisionScore   float64 `json:"vision_score"`
	DamageShare   float64 `json:"damage_share"`
	GoldEfficiency float64 `json:"gold_efficiency"`
	Improvement   string  `json:"improvement_suggestion"`
}

// WebSocket upgrader with gaming-specific configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   4096, // Larger write buffer for gaming data
	HandshakeTimeout:  45 * time.Second,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from Herald.lol domains
		origin := r.Header.Get("Origin")
		return origin == "https://herald.lol" ||
			   origin == "https://www.herald.lol" ||
			   origin == "http://localhost:3000" || // Development
			   origin == "http://localhost:5173"   // Vite dev server
	},
}

// NewHub creates a new WebSocket hub for Herald.lol
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 1000), // Large buffer for high-throughput gaming data
		register:   make(chan *Client, 100),
		unregister: make(chan *Client, 100),
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		userSubs:   make(map[string]map[*Client]bool),
		matchSubs:  make(map[string]map[*Client]bool),
		stats:      HubStats{},
	}
}

// Run starts the WebSocket hub
func (h *Hub) Run(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("WebSocket hub shutting down...")
			return

		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.stats.TotalConnections++
			h.stats.ActiveConnections = len(h.clients)
			h.mu.Unlock()
			
			log.Printf("Client connected: %s (total: %d)", client.userID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				h.stats.ActiveConnections = len(h.clients)
				close(client.send)
				
				// Remove from all subscriptions
				h.removeClientFromSubscriptions(client)
			}
			h.mu.Unlock()
			
			log.Printf("Client disconnected: %s (total: %d)", client.userID, len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			clients := make([]*Client, 0, len(h.clients))
			for client := range h.clients {
				clients = append(clients, client)
			}
			h.mu.RUnlock()

			// Broadcast to all clients
			for _, client := range clients {
				select {
				case client.send <- message:
				default:
					h.mu.Lock()
					delete(h.clients, client)
					close(client.send)
					h.mu.Unlock()
				}
			}

		case <-ticker.C:
			// Send periodic ping to keep connections alive
			h.BroadcastPing()
		}
	}
}

// HandleWebSocket handles WebSocket connection upgrades
func (h *Hub) HandleWebSocket(c *gin.Context) {
	// Extract user authentication
	userID := c.Query("user_id")
	token := c.Query("token")
	
	if userID == "" || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// TODO: Validate JWT token
	// For now, accept any non-empty token
	if !h.validateToken(token, userID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
		rooms:  make(map[string]bool),
	}

	client.hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// BroadcastToUser sends a message to all connections of a specific user
func (h *Hub) BroadcastToUser(userID string, message Message) error {
	message.Timestamp = time.Now()
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	h.mu.RLock()
	clients, exists := h.userSubs[userID]
	h.mu.RUnlock()

	if !exists {
		return nil // User not connected
	}

	for client := range clients {
		select {
		case client.send <- data:
		default:
			// Client buffer is full, remove it
			h.mu.Lock()
			delete(h.userSubs[userID], client)
			h.mu.Unlock()
		}
	}

	return nil
}

// BroadcastToMatch sends a message to all clients watching a specific match
func (h *Hub) BroadcastToMatch(matchID string, message Message) error {
	message.Timestamp = time.Now()
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	h.mu.RLock()
	clients, exists := h.matchSubs[matchID]
	h.mu.RUnlock()

	if !exists {
		return nil // No clients watching this match
	}

	for client := range clients {
		select {
		case client.send <- data:
		default:
			h.mu.Lock()
			delete(h.matchSubs[matchID], client)
			h.mu.Unlock()
		}
	}

	return nil
}

// BroadcastToRoom sends a message to all clients in a room
func (h *Hub) BroadcastToRoom(roomID string, message Message) error {
	message.Timestamp = time.Now()
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	h.mu.RLock()
	clients, exists := h.rooms[roomID]
	h.mu.RUnlock()

	if !exists {
		return nil // Room doesn't exist
	}

	for client := range clients {
		select {
		case client.send <- data:
		default:
			h.mu.Lock()
			delete(h.rooms[roomID], client)
			h.mu.Unlock()
		}
	}

	return nil
}

// BroadcastPing sends ping to all connected clients
func (h *Hub) BroadcastPing() {
	message := Message{
		Type:      MessageTypePing,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling ping message: %v", err)
		return
	}

	h.broadcast <- data
}

// SubscribeUserToMatch subscribes a user to match updates
func (h *Hub) SubscribeUserToMatch(userID, matchID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Find user's clients
	userClients, exists := h.userSubs[userID]
	if !exists {
		return
	}

	// Initialize match subscription if not exists
	if h.matchSubs[matchID] == nil {
		h.matchSubs[matchID] = make(map[*Client]bool)
	}

	// Subscribe all user's clients to the match
	for client := range userClients {
		h.matchSubs[matchID][client] = true
	}
}

// GetStats returns current hub statistics
func (h *Hub) GetStats() HubStats {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	stats := h.stats
	stats.ActiveConnections = len(h.clients)
	return stats
}

// validateToken validates JWT token (placeholder implementation)
func (h *Hub) validateToken(token, userID string) bool {
	// TODO: Implement actual JWT validation
	// For now, accept any non-empty token
	return token != "" && userID != ""
}

// removeClientFromSubscriptions removes client from all subscriptions
func (h *Hub) removeClientFromSubscriptions(client *Client) {
	// Remove from user subscriptions
	for userID, clients := range h.userSubs {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			if len(clients) == 0 {
				delete(h.userSubs, userID)
			}
		}
	}

	// Remove from match subscriptions
	for matchID, clients := range h.matchSubs {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			if len(clients) == 0 {
				delete(h.matchSubs, matchID)
			}
		}
	}

	// Remove from room subscriptions
	for roomID, clients := range h.rooms {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			if len(clients) == 0 {
				delete(h.rooms, roomID)
			}
		}
	}
}