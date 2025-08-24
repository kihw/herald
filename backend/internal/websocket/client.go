package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 8192 // 8KB for gaming data
)

// Client represents a WebSocket client connection
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
	rooms  map[string]bool // Rooms this client is subscribed to

	// Gaming-specific client state
	watchingMatches map[string]bool
	preferences     ClientPreferences
	lastSeen        time.Time
}

// ClientPreferences stores user's real-time notification preferences
type ClientPreferences struct {
	MatchUpdates        bool `json:"match_updates"`
	RankUpdates         bool `json:"rank_updates"`
	FriendActivity      bool `json:"friend_activity"`
	CoachingSuggestions bool `json:"coaching_suggestions"`
	SystemNotifications bool `json:"system_notifications"`
}

// ClientMessage represents incoming messages from clients
type ClientMessage struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

// Subscription actions
const (
	ActionSubscribe         = "subscribe"
	ActionUnsubscribe       = "unsubscribe"
	ActionJoinRoom          = "join_room"
	ActionLeaveRoom         = "leave_room"
	ActionWatchMatch        = "watch_match"
	ActionUnwatchMatch      = "unwatch_match"
	ActionUpdatePreferences = "update_preferences"
	ActionPong              = "pong"
	ActionGetStats          = "get_stats"
)

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for user %s: %v", c.userID, err)
			}
			break
		}

		c.lastSeen = time.Now()

		// Parse client message
		var clientMessage ClientMessage
		if err := json.Unmarshal(messageBytes, &clientMessage); err != nil {
			log.Printf("Error parsing client message from %s: %v", c.userID, err)
			c.sendError("Invalid message format")
			continue
		}

		// Handle client actions
		c.handleClientMessage(clientMessage)
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current WebSocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleClientMessage processes incoming client messages
func (c *Client) handleClientMessage(msg ClientMessage) {
	switch msg.Action {
	case ActionSubscribe:
		c.handleSubscribe(msg.Data)
	case ActionUnsubscribe:
		c.handleUnsubscribe(msg.Data)
	case ActionJoinRoom:
		c.handleJoinRoom(msg.Data)
	case ActionLeaveRoom:
		c.handleLeaveRoom(msg.Data)
	case ActionWatchMatch:
		c.handleWatchMatch(msg.Data)
	case ActionUnwatchMatch:
		c.handleUnwatchMatch(msg.Data)
	case ActionUpdatePreferences:
		c.handleUpdatePreferences(msg.Data)
	case ActionPong:
		// Client responded to ping - update last seen
		c.lastSeen = time.Now()
	case ActionGetStats:
		c.sendStats()
	default:
		c.sendError("Unknown action: " + msg.Action)
	}
}

// handleSubscribe subscribes client to user-specific updates
func (c *Client) handleSubscribe(data interface{}) {
	c.hub.mu.Lock()
	defer c.hub.mu.Unlock()

	if c.hub.userSubs[c.userID] == nil {
		c.hub.userSubs[c.userID] = make(map[*Client]bool)
	}
	c.hub.userSubs[c.userID][c] = true

	c.sendMessage(Message{
		Type: "subscription_confirmed",
		Data: map[string]string{"user_id": c.userID},
	})
}

// handleJoinRoom joins client to a specific room
func (c *Client) handleJoinRoom(data interface{}) {
	roomData, ok := data.(map[string]interface{})
	if !ok {
		c.sendError("Invalid room data format")
		return
	}

	roomID, ok := roomData["room_id"].(string)
	if !ok {
		c.sendError("Room ID is required")
		return
	}

	c.hub.mu.Lock()
	defer c.hub.mu.Unlock()

	if c.hub.rooms[roomID] == nil {
		c.hub.rooms[roomID] = make(map[*Client]bool)
	}
	c.hub.rooms[roomID][c] = true
	c.rooms[roomID] = true

	c.sendMessage(Message{
		Type: "room_joined",
		Data: map[string]string{"room_id": roomID},
	})
}

// handleLeaveRoom removes client from a room
func (c *Client) handleLeaveRoom(data interface{}) {
	roomData, ok := data.(map[string]interface{})
	if !ok {
		c.sendError("Invalid room data format")
		return
	}

	roomID, ok := roomData["room_id"].(string)
	if !ok {
		c.sendError("Room ID is required")
		return
	}

	c.hub.mu.Lock()
	defer c.hub.mu.Unlock()

	if clients, exists := c.hub.rooms[roomID]; exists {
		delete(clients, c)
		if len(clients) == 0 {
			delete(c.hub.rooms, roomID)
		}
	}
	delete(c.rooms, roomID)

	c.sendMessage(Message{
		Type: "room_left",
		Data: map[string]string{"room_id": roomID},
	})
}

// handleWatchMatch subscribes client to match updates
func (c *Client) handleWatchMatch(data interface{}) {
	matchData, ok := data.(map[string]interface{})
	if !ok {
		c.sendError("Invalid match data format")
		return
	}

	matchID, ok := matchData["match_id"].(string)
	if !ok {
		c.sendError("Match ID is required")
		return
	}

	c.hub.mu.Lock()
	defer c.hub.mu.Unlock()

	if c.hub.matchSubs[matchID] == nil {
		c.hub.matchSubs[matchID] = make(map[*Client]bool)
	}
	c.hub.matchSubs[matchID][c] = true

	if c.watchingMatches == nil {
		c.watchingMatches = make(map[string]bool)
	}
	c.watchingMatches[matchID] = true

	c.sendMessage(Message{
		Type: "match_subscribed",
		Data: map[string]string{"match_id": matchID},
	})
}

// handleUnwatchMatch unsubscribes client from match updates
func (c *Client) handleUnwatchMatch(data interface{}) {
	matchData, ok := data.(map[string]interface{})
	if !ok {
		c.sendError("Invalid match data format")
		return
	}

	matchID, ok := matchData["match_id"].(string)
	if !ok {
		c.sendError("Match ID is required")
		return
	}

	c.hub.mu.Lock()
	defer c.hub.mu.Unlock()

	if clients, exists := c.hub.matchSubs[matchID]; exists {
		delete(clients, c)
		if len(clients) == 0 {
			delete(c.hub.matchSubs, matchID)
		}
	}

	if c.watchingMatches != nil {
		delete(c.watchingMatches, matchID)
	}

	c.sendMessage(Message{
		Type: "match_unsubscribed",
		Data: map[string]string{"match_id": matchID},
	})
}

// handleUpdatePreferences updates client's notification preferences
func (c *Client) handleUpdatePreferences(data interface{}) {
	prefsData, ok := data.(map[string]interface{})
	if !ok {
		c.sendError("Invalid preferences data format")
		return
	}

	// Update preferences
	if val, exists := prefsData["match_updates"]; exists {
		if b, ok := val.(bool); ok {
			c.preferences.MatchUpdates = b
		}
	}
	if val, exists := prefsData["rank_updates"]; exists {
		if b, ok := val.(bool); ok {
			c.preferences.RankUpdates = b
		}
	}
	if val, exists := prefsData["friend_activity"]; exists {
		if b, ok := val.(bool); ok {
			c.preferences.FriendActivity = b
		}
	}
	if val, exists := prefsData["coaching_suggestions"]; exists {
		if b, ok := val.(bool); ok {
			c.preferences.CoachingSuggestions = b
		}
	}
	if val, exists := prefsData["system_notifications"]; exists {
		if b, ok := val.(bool); ok {
			c.preferences.SystemNotifications = b
		}
	}

	c.sendMessage(Message{
		Type: "preferences_updated",
		Data: c.preferences,
	})
}

// sendMessage sends a message to the client
func (c *Client) sendMessage(message Message) {
	message.Timestamp = time.Now()
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message for client %s: %v", c.userID, err)
		return
	}

	select {
	case c.send <- data:
	default:
		// Client's send channel is full - disconnect
		close(c.send)
	}
}

// sendError sends an error message to the client
func (c *Client) sendError(errorMsg string) {
	c.sendMessage(Message{
		Type: MessageTypeError,
		Data: map[string]string{"error": errorMsg},
	})
}

// sendStats sends current hub statistics to the client
func (c *Client) sendStats() {
	stats := c.hub.GetStats()
	c.sendMessage(Message{
		Type: "stats",
		Data: stats,
	})
}

// handleUnsubscribe unsubscribes client from user-specific updates
func (c *Client) handleUnsubscribe(data interface{}) {
	c.hub.mu.Lock()
	defer c.hub.mu.Unlock()

	if clients, exists := c.hub.userSubs[c.userID]; exists {
		delete(clients, c)
		if len(clients) == 0 {
			delete(c.hub.userSubs, c.userID)
		}
	}

	c.sendMessage(Message{
		Type: "subscription_cancelled",
		Data: map[string]string{"user_id": c.userID},
	})
}

// IsConnected checks if the client connection is still active
func (c *Client) IsConnected() bool {
	return time.Since(c.lastSeen) < pongWait*2
}
