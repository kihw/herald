package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/herald/internal/models"
	"github.com/herald/internal/websocket"
)

// RealtimeService manages real-time gaming data updates
type RealtimeService struct {
	hub          *websocket.Hub
	riotService  *RiotService
	analyticsService *AnalyticsService
	userService  *UserService
	
	// Active live matches being tracked
	liveMatches  map[string]*LiveMatchTracker
	matchesMu    sync.RWMutex
	
	// Performance update intervals
	performanceUpdateInterval time.Duration
	matchUpdateInterval       time.Duration
	
	// Channels for processing updates
	matchUpdates       chan MatchUpdate
	performanceUpdates chan PerformanceUpdate
	rankUpdates        chan RankUpdate
	
	// Context for graceful shutdown
	ctx    context.Context
	cancel context.CancelFunc
}

// LiveMatchTracker tracks live match data and participants
type LiveMatchTracker struct {
	GameID       string                    `json:"game_id"`
	StartTime    time.Time                 `json:"start_time"`
	GameMode     string                    `json:"game_mode"`
	Participants []LiveParticipantData     `json:"participants"`
	LastUpdate   time.Time                 `json:"last_update"`
	Watchers     map[string]bool          `json:"-"` // UserIDs watching this match
	UpdateCount  int                      `json:"update_count"`
	mu           sync.RWMutex
}

// LiveParticipantData represents live participant data
type LiveParticipantData struct {
	SummonerName   string    `json:"summoner_name"`
	SummonerID     string    `json:"summoner_id"`
	ChampionName   string    `json:"champion_name"`
	ChampionID     int       `json:"champion_id"`
	Team           int       `json:"team"`
	Position       string    `json:"position"`
	Level          int       `json:"level"`
	Kills          int       `json:"kills"`
	Deaths         int       `json:"deaths"`
	Assists        int       `json:"assists"`
	CS             int       `json:"cs"`
	Gold           int       `json:"gold"`
	Items          []int     `json:"items"`
	Runes          []int     `json:"runes"`
	SummonerSpells []int     `json:"summoner_spells"`
	KDA            float64   `json:"kda"`
	LastUpdate     time.Time `json:"last_update"`
}

// Update event types
type MatchUpdate struct {
	GameID   string
	UserID   string
	Data     websocket.MatchUpdateData
	Priority int // 1=low, 2=medium, 3=high
}

type PerformanceUpdate struct {
	UserID string
	Data   websocket.PerformanceUpdateData
}

type RankUpdate struct {
	UserID   string
	OldRank  string
	NewRank  string
	LP       int
	Series   *RankedSeries
}

type RankedSeries struct {
	Target  string `json:"target"`
	Wins    int    `json:"wins"`
	Losses  int    `json:"losses"`
	Progress string `json:"progress"`
}

// NewRealtimeService creates a new realtime service
func NewRealtimeService(hub *websocket.Hub, riotService *RiotService, analyticsService *AnalyticsService, userService *UserService) *RealtimeService {
	ctx, cancel := context.WithCancel(context.Background())
	
	service := &RealtimeService{
		hub:                      hub,
		riotService:             riotService,
		analyticsService:        analyticsService,
		userService:             userService,
		liveMatches:             make(map[string]*LiveMatchTracker),
		performanceUpdateInterval: 30 * time.Second,
		matchUpdateInterval:       15 * time.Second,
		matchUpdates:             make(chan MatchUpdate, 1000),
		performanceUpdates:       make(chan PerformanceUpdate, 500),
		rankUpdates:              make(chan RankUpdate, 100),
		ctx:                      ctx,
		cancel:                   cancel,
	}
	
	return service
}

// Start begins the realtime service processing
func (rs *RealtimeService) Start() {
	log.Println("Starting Herald.lol Realtime Service...")
	
	// Start background processors
	go rs.processMatchUpdates()
	go rs.processPerformanceUpdates()
	go rs.processRankUpdates()
	go rs.trackLiveMatches()
	go rs.periodicPerformanceUpdates()
	
	log.Println("Realtime Service started successfully")
}

// Stop gracefully shuts down the realtime service
func (rs *RealtimeService) Stop() {
	log.Println("Stopping Realtime Service...")
	rs.cancel()
	
	// Close channels
	close(rs.matchUpdates)
	close(rs.performanceUpdates)
	close(rs.rankUpdates)
}

// StartTrackingMatch begins tracking a live match
func (rs *RealtimeService) StartTrackingMatch(gameID, userID string) error {
	rs.matchesMu.Lock()
	defer rs.matchesMu.Unlock()
	
	// Check if match is already being tracked
	if tracker, exists := rs.liveMatches[gameID]; exists {
		tracker.mu.Lock()
		tracker.Watchers[userID] = true
		tracker.mu.Unlock()
		return nil
	}
	
	// Get live match data from Riot API
	liveGameInfo, err := rs.riotService.GetLiveGame(userID)
	if err != nil {
		return fmt.Errorf("failed to get live game info: %w", err)
	}
	
	// Create new tracker
	tracker := &LiveMatchTracker{
		GameID:     liveGameInfo.GameID,
		StartTime:  time.Unix(liveGameInfo.StartTime/1000, 0),
		GameMode:   liveGameInfo.GameMode,
		LastUpdate: time.Now(),
		Watchers:   map[string]bool{userID: true},
	}
	
	// Convert participants
	for _, participant := range liveGameInfo.Participants {
		tracker.Participants = append(tracker.Participants, LiveParticipantData{
			SummonerName:   participant.SummonerName,
			SummonerID:     participant.SummonerID,
			ChampionName:   participant.ChampionName,
			ChampionID:     participant.ChampionID,
			Team:           participant.TeamID,
			SummonerSpells: participant.SummonerSpells,
			Runes:          participant.Runes,
			LastUpdate:     time.Now(),
		})
	}
	
	rs.liveMatches[gameID] = tracker
	
	log.Printf("Started tracking live match %s for user %s", gameID, userID)
	return nil
}

// StopTrackingMatch stops tracking a match for a specific user
func (rs *RealtimeService) StopTrackingMatch(gameID, userID string) {
	rs.matchesMu.Lock()
	defer rs.matchesMu.Unlock()
	
	tracker, exists := rs.liveMatches[gameID]
	if !exists {
		return
	}
	
	tracker.mu.Lock()
	delete(tracker.Watchers, userID)
	hasWatchers := len(tracker.Watchers) > 0
	tracker.mu.Unlock()
	
	// If no one is watching, remove the tracker
	if !hasWatchers {
		delete(rs.liveMatches, gameID)
		log.Printf("Stopped tracking match %s - no more watchers", gameID)
	}
}

// SendPerformanceUpdate sends real-time performance update to user
func (rs *RealtimeService) SendPerformanceUpdate(userID string) {
	// Get user's recent match performance
	recentMatches, err := rs.riotService.GetMatchHistory(userID, 5)
	if err != nil {
		log.Printf("Error getting recent matches for performance update: %v", err)
		return
	}
	
	if len(recentMatches.Matches) == 0 {
		return
	}
	
	// Calculate current session performance
	var totalKDA, totalCS, totalVision, totalDamage, totalGold float64
	matchCount := float64(len(recentMatches.Matches))
	
	for _, match := range recentMatches.Matches {
		// Get detailed match data
		matchDetail, err := rs.riotService.GetMatchDetails(match.MatchID)
		if err != nil {
			continue
		}
		
		// Find user's participant data
		for _, participant := range matchDetail.Info.Participants {
			if participant.SummonerID == userID {
				kda := float64(participant.Kills + participant.Assists)
				if participant.Deaths > 0 {
					kda = kda / float64(participant.Deaths)
				}
				
				totalKDA += kda
				totalCS += float64(participant.TotalMinionsKilled + participant.NeutralMinionsKilled)
				totalVision += float64(participant.VisionScore)
				totalDamage += float64(participant.TotalDamageDealtToChampions)
				totalGold += float64(participant.GoldEarned)
				break
			}
		}
	}
	
	// Calculate averages
	avgKDA := totalKDA / matchCount
	avgCS := totalCS / matchCount
	avgVision := totalVision / matchCount
	avgDamage := totalDamage / matchCount
	avgGold := totalGold / matchCount
	
	// Get historical averages for comparison
	historicalStats, err := rs.analyticsService.GetPlayerStats(userID, 30) // 30 days
	if err != nil {
		log.Printf("Error getting historical stats: %v", err)
		return
	}
	
	// Generate improvement suggestion
	improvement := rs.generateImprovementSuggestion(avgKDA, avgCS, avgVision, historicalStats)
	
	// Create performance update
	update := PerformanceUpdate{
		UserID: userID,
		Data: websocket.PerformanceUpdateData{
			UserID:         userID,
			CurrentKDA:     avgKDA,
			AverageKDA:     historicalStats.AverageKDA,
			CSPerMinute:    avgCS / 30.0, // Assuming 30min average game
			VisionScore:    avgVision,
			DamageShare:    (avgDamage / (avgDamage + 40000)) * 100, // Rough team damage estimate
			GoldEfficiency: (avgGold / 15000) * 100, // Rough efficiency calculation
			Improvement:    improvement,
		},
	}
	
	// Send update
	select {
	case rs.performanceUpdates <- update:
	default:
		log.Printf("Performance update channel full for user %s", userID)
	}
}

// SendRankUpdate sends rank update notification
func (rs *RealtimeService) SendRankUpdate(userID, oldRank, newRank string, lp int) {
	update := RankUpdate{
		UserID:  userID,
		OldRank: oldRank,
		NewRank: newRank,
		LP:      lp,
	}
	
	select {
	case rs.rankUpdates <- update:
	default:
		log.Printf("Rank update channel full for user %s", userID)
	}
}

// Background processors

func (rs *RealtimeService) processMatchUpdates() {
	for {
		select {
		case <-rs.ctx.Done():
			return
		case update := <-rs.matchUpdates:
			message := websocket.Message{
				Type:    websocket.MessageTypeMatchUpdate,
				UserID:  update.UserID,
				MatchID: update.Data.GameID,
				Data:    update.Data,
			}
			
			// Send to specific user
			if update.UserID != "" {
				rs.hub.BroadcastToUser(update.UserID, message)
			}
			
			// Also send to anyone watching this match
			rs.hub.BroadcastToMatch(update.GameID, message)
		}
	}
}

func (rs *RealtimeService) processPerformanceUpdates() {
	for {
		select {
		case <-rs.ctx.Done():
			return
		case update := <-rs.performanceUpdates:
			message := websocket.Message{
				Type:   websocket.MessageTypePerformanceUpdate,
				UserID: update.UserID,
				Data:   update.Data,
			}
			
			rs.hub.BroadcastToUser(update.UserID, message)
		}
	}
}

func (rs *RealtimeService) processRankUpdates() {
	for {
		select {
		case <-rs.ctx.Done():
			return
		case update := <-rs.rankUpdates:
			message := websocket.Message{
				Type:   websocket.MessageTypeRankUpdate,
				UserID: update.UserID,
				Data: map[string]interface{}{
					"old_rank": update.OldRank,
					"new_rank": update.NewRank,
					"lp":       update.LP,
					"series":   update.Series,
				},
			}
			
			rs.hub.BroadcastToUser(update.UserID, message)
		}
	}
}

func (rs *RealtimeService) trackLiveMatches() {
	ticker := time.NewTicker(rs.matchUpdateInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-rs.ctx.Done():
			return
		case <-ticker.C:
			rs.updateLiveMatches()
		}
	}
}

func (rs *RealtimeService) periodicPerformanceUpdates() {
	ticker := time.NewTicker(rs.performanceUpdateInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-rs.ctx.Done():
			return
		case <-ticker.C:
			// Send performance updates to all connected users
			stats := rs.hub.GetStats()
			log.Printf("Sending periodic performance updates to %d users", stats.ActiveConnections)
			
			// This would typically be triggered by actual game events
			// For now, we'll skip automatic updates to avoid spam
		}
	}
}

func (rs *RealtimeService) updateLiveMatches() {
	rs.matchesMu.RLock()
	matchesToUpdate := make(map[string]*LiveMatchTracker)
	for gameID, tracker := range rs.liveMatches {
		matchesToUpdate[gameID] = tracker
	}
	rs.matchesMu.RUnlock()
	
	for gameID, tracker := range matchesToUpdate {
		// Skip if updated recently
		if time.Since(tracker.LastUpdate) < rs.matchUpdateInterval {
			continue
		}
		
		// Get updated match data (this would typically come from Riot's spectator API)
		// For now, we'll simulate updates
		rs.simulateMatchUpdate(gameID, tracker)
	}
}

func (rs *RealtimeService) simulateMatchUpdate(gameID string, tracker *LiveMatchTracker) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()
	
	// Simulate some updates
	tracker.UpdateCount++
	tracker.LastUpdate = time.Now()
	
	// Create match update data
	updateData := websocket.MatchUpdateData{
		GameID:   gameID,
		Status:   "in_progress",
		GameTime: int(time.Since(tracker.StartTime).Minutes()),
		Participants: make([]websocket.ParticipantData, len(tracker.Participants)),
	}
	
	// Convert participants
	for i, participant := range tracker.Participants {
		updateData.Participants[i] = websocket.ParticipantData{
			SummonerName: participant.SummonerName,
			ChampionName: participant.ChampionName,
			Level:        participant.Level,
			Kills:        participant.Kills,
			Deaths:       participant.Deaths,
			Assists:      participant.Assists,
			CS:           participant.CS,
			Gold:         participant.Gold,
			Items:        participant.Items,
			KDA:          participant.KDA,
		}
	}
	
	// Send update to all watchers
	for userID := range tracker.Watchers {
		update := MatchUpdate{
			GameID:   gameID,
			UserID:   userID,
			Data:     updateData,
			Priority: 2, // Medium priority
		}
		
		select {
		case rs.matchUpdates <- update:
		default:
			log.Printf("Match update channel full for game %s", gameID)
		}
	}
}

func (rs *RealtimeService) generateImprovementSuggestion(currentKDA, currentCS, currentVision float64, historical *models.PlayerStats) string {
	suggestions := []string{}
	
	if currentKDA < historical.AverageKDA * 0.9 {
		suggestions = append(suggestions, "Focus on positioning and avoid risky trades")
	}
	
	if currentCS < historical.AverageCS * 0.9 {
		suggestions = append(suggestions, "Practice last-hitting minions and farming efficiency")
	}
	
	if currentVision < historical.AverageVisionScore * 0.9 {
		suggestions = append(suggestions, "Place more wards and buy control wards")
	}
	
	if len(suggestions) == 0 {
		return "Great performance! Keep up the consistency"
	}
	
	return suggestions[0] // Return first suggestion for now
}