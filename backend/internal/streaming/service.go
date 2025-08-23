package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/herald-lol/backend/internal/analytics"
	"github.com/herald-lol/backend/internal/match"
	"github.com/herald-lol/backend/internal/riot"
)

// Herald.lol Gaming Analytics - Real-time Streaming Service
// Live gaming data streaming for real-time analytics and notifications

// StreamingService handles real-time gaming data streaming
type StreamingService struct {
	config           *StreamingConfig
	
	// Core services
	riotClient       *riot.RiotClient
	analyticsEngine  *analytics.AnalyticsEngine
	matchAnalyzer    *match.MatchAnalyzer
	
	// Connection management
	clients          map[string]*ClientConnection
	clientMutex      sync.RWMutex
	channels         map[string]*StreamChannel
	channelMutex     sync.RWMutex
	
	// Live match tracking
	liveMatches      map[string]*LiveMatchTracker
	liveMatchMutex   sync.RWMutex
	
	// Event processors
	eventProcessors  map[string]EventProcessor
	eventQueue       chan *StreamEvent
	
	// Performance monitoring
	metrics          *StreamingMetrics
	
	// Cleanup and lifecycle
	shutdown         chan struct{}
	shutdownOnce     sync.Once
}

// NewStreamingService creates a new real-time streaming service
func NewStreamingService(
	config *StreamingConfig,
	riotClient *riot.RiotClient,
	analyticsEngine *analytics.AnalyticsEngine,
	matchAnalyzer *match.MatchAnalyzer,
) *StreamingService {
	service := &StreamingService{
		config:          config,
		riotClient:      riotClient,
		analyticsEngine: analyticsEngine,
		matchAnalyzer:   matchAnalyzer,
		
		clients:         make(map[string]*ClientConnection),
		channels:        make(map[string]*StreamChannel),
		liveMatches:     make(map[string]*LiveMatchTracker),
		eventProcessors: make(map[string]EventProcessor),
		eventQueue:      make(chan *StreamEvent, config.EventQueueSize),
		metrics:         NewStreamingMetrics(),
		shutdown:        make(chan struct{}),
	}
	
	// Initialize event processors
	service.initializeEventProcessors()
	
	// Start background workers
	go service.eventWorker()
	go service.liveMatchWorker()
	go service.metricsWorker()
	go service.cleanupWorker()
	
	return service
}

// HandleWebSocketConnection handles new WebSocket connections for real-time streaming
func (s *StreamingService) HandleWebSocketConnection(ws *websocket.Conn, userID, playerPUUID string) {
	clientID := s.generateClientID()
	
	client := &ClientConnection{
		ID:          clientID,
		UserID:      userID,
		PlayerPUUID: playerPUUID,
		Connection:  ws,
		Channels:    make(map[string]bool),
		LastPing:    time.Now(),
		Connected:   true,
		JoinedAt:    time.Now(),
		MessageCount: 0,
	}
	
	// Register client
	s.clientMutex.Lock()
	s.clients[clientID] = client
	s.clientMutex.Unlock()
	
	s.metrics.IncrementConnections()
	
	// Send welcome message
	welcome := &StreamMessage{
		Type:      "welcome",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"client_id":    clientID,
			"server_time":  time.Now().Unix(),
			"capabilities": s.getClientCapabilities(),
		},
	}
	s.sendToClient(client, welcome)
	
	// Handle client messages
	go s.handleClientMessages(client)
	
	// Start periodic updates for this client
	go s.clientUpdateWorker(client)
	
	log.Printf("🎮 Client connected: %s (User: %s, Player: %s)", clientID, userID, playerPUUID)
}

// SubscribeToLiveMatch subscribes a client to live match updates
func (s *StreamingService) SubscribeToLiveMatch(clientID, matchID string) error {
	s.clientMutex.RLock()
	client, exists := s.clients[clientID]
	s.clientMutex.RUnlock()
	
	if !exists {
		return fmt.Errorf("client not found: %s", clientID)
	}
	
	// Start tracking live match if not already tracked
	s.liveMatchMutex.Lock()
	tracker, exists := s.liveMatches[matchID]
	if !exists {
		tracker = s.createLiveMatchTracker(matchID)
		s.liveMatches[matchID] = tracker
		go s.trackLiveMatch(tracker)
	}
	tracker.Subscribers[clientID] = true
	s.liveMatchMutex.Unlock()
	
	// Add client to live match channel
	channelName := fmt.Sprintf("live_match:%s", matchID)
	s.subscribeClientToChannel(client, channelName)
	
	// Send initial match state
	initialState := &StreamMessage{
		Type:      "live_match_state",
		Channel:   channelName,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"match_id":    matchID,
			"state":       tracker.CurrentState,
			"game_time":   tracker.GameTime,
			"participants": tracker.Participants,
		},
	}
	s.sendToClient(client, initialState)
	
	return nil
}

// SubscribeToPlayerUpdates subscribes a client to player performance updates
func (s *StreamingService) SubscribeToPlayerUpdates(clientID, playerPUUID string) error {
	s.clientMutex.RLock()
	client, exists := s.clients[clientID]
	s.clientMutex.RUnlock()
	
	if !exists {
		return fmt.Errorf("client not found: %s", clientID)
	}
	
	channelName := fmt.Sprintf("player:%s", playerPUUID)
	s.subscribeClientToChannel(client, channelName)
	
	// Send current player status
	status := s.getPlayerCurrentStatus(playerPUUID)
	statusMessage := &StreamMessage{
		Type:      "player_status",
		Channel:   channelName,
		Timestamp: time.Now(),
		Data:      status,
	}
	s.sendToClient(client, statusMessage)
	
	return nil
}

// SubscribeToAnalyticsUpdates subscribes a client to analytics updates
func (s *StreamingService) SubscribeToAnalyticsUpdates(clientID string, updateTypes []string) error {
	s.clientMutex.RLock()
	client, exists := s.clients[clientID]
	s.clientMutex.RUnlock()
	
	if !exists {
		return fmt.Errorf("client not found: %s", clientID)
	}
	
	for _, updateType := range updateTypes {
		channelName := fmt.Sprintf("analytics:%s", updateType)
		s.subscribeClientToChannel(client, channelName)
	}
	
	return nil
}

// BroadcastLiveMatchEvent broadcasts a live match event to all subscribers
func (s *StreamingService) BroadcastLiveMatchEvent(matchID string, event *LiveMatchEvent) {
	channelName := fmt.Sprintf("live_match:%s", matchID)
	
	message := &StreamMessage{
		Type:      "live_match_event",
		Channel:   channelName,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"match_id": matchID,
			"event":    event,
		},
	}
	
	s.broadcastToChannel(channelName, message)
	
	// Queue for processing
	streamEvent := &StreamEvent{
		Type:      "live_match_event",
		MatchID:   matchID,
		Data:      event,
		Timestamp: time.Now(),
	}
	
	select {
	case s.eventQueue <- streamEvent:
		s.metrics.IncrementEvents()
	default:
		s.metrics.IncrementDroppedEvents()
		log.Printf("⚠️ Event queue full, dropping event for match %s", matchID)
	}
}

// BroadcastPlayerUpdate broadcasts a player performance update
func (s *StreamingService) BroadcastPlayerUpdate(playerPUUID string, update *PlayerUpdate) {
	channelName := fmt.Sprintf("player:%s", playerPUUID)
	
	message := &StreamMessage{
		Type:      "player_update",
		Channel:   channelName,
		Timestamp: time.Now(),
		Data:      update,
	}
	
	s.broadcastToChannel(channelName, message)
}

// BroadcastAnalyticsUpdate broadcasts an analytics update
func (s *StreamingService) BroadcastAnalyticsUpdate(updateType string, data interface{}) {
	channelName := fmt.Sprintf("analytics:%s", updateType)
	
	message := &StreamMessage{
		Type:      "analytics_update",
		Channel:   channelName,
		Timestamp: time.Now(),
		Data:      data,
	}
	
	s.broadcastToChannel(channelName, message)
}

// SendNotification sends a notification to specific clients
func (s *StreamingService) SendNotification(userIDs []string, notification *Notification) {
	message := &StreamMessage{
		Type:      "notification",
		Timestamp: time.Now(),
		Data:      notification,
	}
	
	s.clientMutex.RLock()
	defer s.clientMutex.RUnlock()
	
	for _, client := range s.clients {
		for _, targetUserID := range userIDs {
			if client.UserID == targetUserID {
				s.sendToClient(client, message)
				break
			}
		}
	}
}

// GetStreamingStats returns current streaming service statistics
func (s *StreamingService) GetStreamingStats() *StreamingStats {
	s.clientMutex.RLock()
	clientCount := len(s.clients)
	s.clientMutex.RUnlock()
	
	s.channelMutex.RLock()
	channelCount := len(s.channels)
	s.channelMutex.RUnlock()
	
	s.liveMatchMutex.RLock()
	liveMatchCount := len(s.liveMatches)
	s.liveMatchMutex.RUnlock()
	
	return &StreamingStats{
		ConnectedClients:   clientCount,
		ActiveChannels:     channelCount,
		LiveMatches:        liveMatchCount,
		EventsProcessed:    s.metrics.EventsProcessed,
		MessagesDelivered:  s.metrics.MessagesDelivered,
		DroppedEvents:     s.metrics.DroppedEvents,
		AverageLatency:    s.metrics.AverageLatency,
		Uptime:           time.Since(s.metrics.StartTime),
	}
}

// Internal helper methods

func (s *StreamingService) handleClientMessages(client *ClientConnection) {
	defer s.disconnectClient(client.ID)
	
	for {
		select {
		case <-s.shutdown:
			return
		default:
			var message ClientMessage
			err := client.Connection.ReadJSON(&message)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error for client %s: %v", client.ID, err)
				}
				return
			}
			
			client.LastPing = time.Now()
			client.MessageCount++
			
			s.handleClientMessage(client, &message)
		}
	}
}

func (s *StreamingService) handleClientMessage(client *ClientConnection, message *ClientMessage) {
	switch message.Type {
	case "ping":
		s.sendToClient(client, &StreamMessage{
			Type:      "pong",
			Timestamp: time.Now(),
		})
		
	case "subscribe":
		if channelName, ok := message.Data["channel"].(string); ok {
			s.subscribeClientToChannel(client, channelName)
		}
		
	case "unsubscribe":
		if channelName, ok := message.Data["channel"].(string); ok {
			s.unsubscribeClientFromChannel(client, channelName)
		}
		
	case "live_match_subscribe":
		if matchID, ok := message.Data["match_id"].(string); ok {
			s.SubscribeToLiveMatch(client.ID, matchID)
		}
		
	case "player_subscribe":
		if playerPUUID, ok := message.Data["player_puuid"].(string); ok {
			s.SubscribeToPlayerUpdates(client.ID, playerPUUID)
		}
		
	default:
		log.Printf("Unknown message type from client %s: %s", client.ID, message.Type)
	}
}

func (s *StreamingService) subscribeClientToChannel(client *ClientConnection, channelName string) {
	s.channelMutex.Lock()
	channel, exists := s.channels[channelName]
	if !exists {
		channel = &StreamChannel{
			Name:        channelName,
			Subscribers: make(map[string]*ClientConnection),
			CreatedAt:   time.Now(),
			MessageCount: 0,
		}
		s.channels[channelName] = channel
	}
	channel.Subscribers[client.ID] = client
	s.channelMutex.Unlock()
	
	client.Channels[channelName] = true
	
	// Send subscription confirmation
	s.sendToClient(client, &StreamMessage{
		Type:      "subscribed",
		Channel:   channelName,
		Timestamp: time.Now(),
	})
}

func (s *StreamingService) unsubscribeClientFromChannel(client *ClientConnection, channelName string) {
	s.channelMutex.Lock()
	if channel, exists := s.channels[channelName]; exists {
		delete(channel.Subscribers, client.ID)
		if len(channel.Subscribers) == 0 {
			delete(s.channels, channelName)
		}
	}
	s.channelMutex.Unlock()
	
	delete(client.Channels, channelName)
	
	// Send unsubscription confirmation
	s.sendToClient(client, &StreamMessage{
		Type:      "unsubscribed",
		Channel:   channelName,
		Timestamp: time.Now(),
	})
}

func (s *StreamingService) broadcastToChannel(channelName string, message *StreamMessage) {
	s.channelMutex.RLock()
	channel, exists := s.channels[channelName]
	if !exists {
		s.channelMutex.RUnlock()
		return
	}
	
	subscribers := make([]*ClientConnection, 0, len(channel.Subscribers))
	for _, client := range channel.Subscribers {
		if client.Connected {
			subscribers = append(subscribers, client)
		}
	}
	channel.MessageCount++
	s.channelMutex.RUnlock()
	
	// Send to all subscribers
	for _, client := range subscribers {
		s.sendToClient(client, message)
	}
	
	s.metrics.IncrementMessages(len(subscribers))
}

func (s *StreamingService) sendToClient(client *ClientConnection, message *StreamMessage) {
	if !client.Connected {
		return
	}
	
	startTime := time.Now()
	
	err := client.Connection.WriteJSON(message)
	if err != nil {
		log.Printf("Failed to send message to client %s: %v", client.ID, err)
		s.disconnectClient(client.ID)
		return
	}
	
	latency := time.Since(startTime)
	s.metrics.RecordLatency(latency)
}

func (s *StreamingService) disconnectClient(clientID string) {
	s.clientMutex.Lock()
	client, exists := s.clients[clientID]
	if exists {
		client.Connected = false
		delete(s.clients, clientID)
	}
	s.clientMutex.Unlock()
	
	if exists {
		// Remove from all channels
		s.channelMutex.Lock()
		for channelName := range client.Channels {
			if channel, exists := s.channels[channelName]; exists {
				delete(channel.Subscribers, clientID)
				if len(channel.Subscribers) == 0 {
					delete(s.channels, channelName)
				}
			}
		}
		s.channelMutex.Unlock()
		
		// Remove from live matches
		s.liveMatchMutex.Lock()
		for _, tracker := range s.liveMatches {
			delete(tracker.Subscribers, clientID)
		}
		s.liveMatchMutex.Unlock()
		
		client.Connection.Close()
		s.metrics.DecrementConnections()
		
		log.Printf("🔌 Client disconnected: %s", clientID)
	}
}

func (s *StreamingService) createLiveMatchTracker(matchID string) *LiveMatchTracker {
	return &LiveMatchTracker{
		MatchID:      matchID,
		StartTime:    time.Now(),
		GameTime:     0,
		CurrentState: "in_progress",
		Participants: make(map[string]*LiveParticipant),
		Subscribers:  make(map[string]bool),
		Events:       []*LiveMatchEvent{},
		LastUpdate:   time.Now(),
	}
}

func (s *StreamingService) trackLiveMatch(tracker *LiveMatchTracker) {
	ticker := time.NewTicker(s.config.LiveMatchUpdateInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-s.shutdown:
			return
		case <-ticker.C:
			if len(tracker.Subscribers) == 0 {
				// No subscribers, stop tracking
				s.liveMatchMutex.Lock()
				delete(s.liveMatches, tracker.MatchID)
				s.liveMatchMutex.Unlock()
				return
			}
			
			// Update match state from Riot API
			s.updateLiveMatchState(tracker)
		}
	}
}

func (s *StreamingService) updateLiveMatchState(tracker *LiveMatchTracker) {
	// In a real implementation, this would call Riot's live game API
	// For now, simulate updates
	tracker.GameTime += int(s.config.LiveMatchUpdateInterval.Seconds())
	tracker.LastUpdate = time.Now()
	
	// Simulate random events
	if tracker.GameTime%300 == 0 { // Every 5 minutes
		event := &LiveMatchEvent{
			Type:        "objective_taken",
			GameTime:    tracker.GameTime,
			Description: "Dragon taken by Blue Team",
			Participants: []string{"player1", "player2", "player3", "player4", "player5"},
			Impact:      "medium",
			Timestamp:   time.Now(),
		}
		
		tracker.Events = append(tracker.Events, event)
		s.BroadcastLiveMatchEvent(tracker.MatchID, event)
	}
}

func (s *StreamingService) getPlayerCurrentStatus(playerPUUID string) map[string]interface{} {
	// In a real implementation, this would fetch current player status
	return map[string]interface{}{
		"player_puuid":   playerPUUID,
		"online_status":  "in_game",
		"current_rank":   "Gold III",
		"lp":            65,
		"current_match":  nil,
		"last_seen":     time.Now().Unix(),
	}
}

func (s *StreamingService) getClientCapabilities() map[string]interface{} {
	return map[string]interface{}{
		"live_matches":      s.config.EnableLiveMatches,
		"player_updates":    s.config.EnablePlayerUpdates,
		"analytics_stream":  s.config.EnableAnalyticsStream,
		"notifications":     s.config.EnableNotifications,
		"max_channels":      s.config.MaxChannelsPerClient,
		"update_frequency":  s.config.LiveMatchUpdateInterval.String(),
	}
}

func (s *StreamingService) generateClientID() string {
	return fmt.Sprintf("client_%d_%d", time.Now().UnixNano(), len(s.clients))
}

// Background workers

func (s *StreamingService) eventWorker() {
	for {
		select {
		case <-s.shutdown:
			return
		case event := <-s.eventQueue:
			s.processStreamEvent(event)
		}
	}
}

func (s *StreamingService) processStreamEvent(event *StreamEvent) {
	processor, exists := s.eventProcessors[event.Type]
	if !exists {
		log.Printf("No processor for event type: %s", event.Type)
		return
	}
	
	err := processor.Process(event)
	if err != nil {
		log.Printf("Failed to process event %s: %v", event.Type, err)
		s.metrics.IncrementFailedEvents()
	} else {
		s.metrics.IncrementProcessedEvents()
	}
}

func (s *StreamingService) liveMatchWorker() {
	// Periodically check for new live matches to track
	ticker := time.NewTicker(s.config.LiveMatchScanInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-s.shutdown:
			return
		case <-ticker.C:
			s.scanForNewLiveMatches()
		}
	}
}

func (s *StreamingService) scanForNewLiveMatches() {
	// In a real implementation, this would scan for new live matches
	// from active players and start tracking them
	log.Printf("🔍 Scanning for new live matches...")
}

func (s *StreamingService) clientUpdateWorker(client *ClientConnection) {
	ticker := time.NewTicker(s.config.ClientUpdateInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-s.shutdown:
			return
		case <-ticker.C:
			if !client.Connected {
				return
			}
			
			// Send periodic updates if client has subscriptions
			if len(client.Channels) > 0 {
				s.sendPeriodicUpdates(client)
			}
			
			// Check for client timeout
			if time.Since(client.LastPing) > s.config.ClientTimeout {
				log.Printf("Client %s timed out", client.ID)
				s.disconnectClient(client.ID)
				return
			}
		}
	}
}

func (s *StreamingService) sendPeriodicUpdates(client *ClientConnection) {
	// Send heartbeat
	heartbeat := &StreamMessage{
		Type:      "heartbeat",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"server_time": time.Now().Unix(),
			"channels":    len(client.Channels),
		},
	}
	s.sendToClient(client, heartbeat)
}

func (s *StreamingService) metricsWorker() {
	ticker := time.NewTicker(s.config.MetricsUpdateInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-s.shutdown:
			return
		case <-ticker.C:
			s.updateMetrics()
		}
	}
}

func (s *StreamingService) updateMetrics() {
	s.clientMutex.RLock()
	clientCount := len(s.clients)
	s.clientMutex.RUnlock()
	
	s.channelMutex.RLock()
	channelCount := len(s.channels)
	s.channelMutex.RUnlock()
	
	s.liveMatchMutex.RLock()
	liveMatchCount := len(s.liveMatches)
	s.liveMatchMutex.RUnlock()
	
	s.metrics.UpdateStats(clientCount, channelCount, liveMatchCount)
}

func (s *StreamingService) cleanupWorker() {
	ticker := time.NewTicker(s.config.CleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-s.shutdown:
			return
		case <-ticker.C:
			s.performCleanup()
		}
	}
}

func (s *StreamingService) performCleanup() {
	now := time.Now()
	
	// Clean up empty channels
	s.channelMutex.Lock()
	for name, channel := range s.channels {
		if len(channel.Subscribers) == 0 || time.Since(channel.CreatedAt) > s.config.ChannelTTL {
			delete(s.channels, name)
		}
	}
	s.channelMutex.Unlock()
	
	// Clean up inactive live match trackers
	s.liveMatchMutex.Lock()
	for matchID, tracker := range s.liveMatches {
		if len(tracker.Subscribers) == 0 && time.Since(tracker.LastUpdate) > s.config.LiveMatchTTL {
			delete(s.liveMatches, matchID)
		}
	}
	s.liveMatchMutex.Unlock()
	
	log.Printf("🧹 Cleanup completed at %s", now.Format("15:04:05"))
}

func (s *StreamingService) initializeEventProcessors() {
	s.eventProcessors["live_match_event"] = &LiveMatchEventProcessor{
		analyticsEngine: s.analyticsEngine,
		matchAnalyzer:   s.matchAnalyzer,
	}
	s.eventProcessors["player_update"] = &PlayerUpdateProcessor{
		analyticsEngine: s.analyticsEngine,
	}
	s.eventProcessors["analytics_update"] = &AnalyticsUpdateProcessor{
		analyticsEngine: s.analyticsEngine,
	}
}

// Shutdown gracefully shuts down the streaming service
func (s *StreamingService) Shutdown(ctx context.Context) error {
	s.shutdownOnce.Do(func() {
		close(s.shutdown)
		
		// Close all client connections
		s.clientMutex.Lock()
		for _, client := range s.clients {
			client.Connection.Close()
		}
		s.clients = make(map[string]*ClientConnection)
		s.clientMutex.Unlock()
		
		log.Printf("🔄 Streaming service shutdown completed")
	})
	
	return nil
}