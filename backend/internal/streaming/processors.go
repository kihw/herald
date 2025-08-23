package streaming

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/herald-lol/backend/internal/analytics"
	"github.com/herald-lol/backend/internal/match"
)

// Herald.lol Gaming Analytics - Streaming Event Processors
// Event processors for different types of real-time gaming events

// LiveMatchEventProcessor processes live match events
type LiveMatchEventProcessor struct {
	analyticsEngine *analytics.AnalyticsEngine
	matchAnalyzer   *match.MatchAnalyzer
}

// Process processes a live match event
func (p *LiveMatchEventProcessor) Process(event *StreamEvent) error {
	liveEvent, ok := event.Data.(*LiveMatchEvent)
	if !ok {
		return fmt.Errorf("invalid data type for live match event")
	}

	log.Printf("🎮 Processing live match event: %s in match %s", liveEvent.Type, event.MatchID)

	// Process different types of live match events
	switch liveEvent.Type {
	case "kill":
		return p.processKillEvent(event.MatchID, liveEvent)
	case "death":
		return p.processDeathEvent(event.MatchID, liveEvent)
	case "objective_taken":
		return p.processObjectiveEvent(event.MatchID, liveEvent)
	case "item_purchased":
		return p.processItemEvent(event.MatchID, liveEvent)
	case "level_up":
		return p.processLevelUpEvent(event.MatchID, liveEvent)
	case "team_fight":
		return p.processTeamFightEvent(event.MatchID, liveEvent)
	case "game_end":
		return p.processGameEndEvent(event.MatchID, liveEvent)
	default:
		return p.processGenericEvent(event.MatchID, liveEvent)
	}
}

func (p *LiveMatchEventProcessor) processKillEvent(matchID string, event *LiveMatchEvent) error {
	// Extract kill details
	killer := ""
	victim := ""
	assists := []string{}
	
	if details := event.Details; details != nil {
		if k, ok := details["killer"].(string); ok {
			killer = k
		}
		if v, ok := details["victim"].(string); ok {
			victim = v
		}
		if a, ok := details["assists"].([]string); ok {
			assists = a
		}
	}

	// Update player statistics
	for _, participant := range event.Participants {
		if participant == killer {
			// Increment kills for killer
			p.updatePlayerStat(participant, "kills", 1)
		} else if participant == victim {
			// Increment deaths for victim
			p.updatePlayerStat(participant, "deaths", 1)
		} else {
			// Check if this participant assisted
			for _, assist := range assists {
				if participant == assist {
					p.updatePlayerStat(participant, "assists", 1)
					break
				}
			}
		}
	}

	// Trigger real-time analytics update
	p.triggerPerformanceUpdate(event.Participants)
	
	log.Printf("✅ Processed kill event: %s killed %s with %d assists", killer, victim, len(assists))
	return nil
}

func (p *LiveMatchEventProcessor) processDeathEvent(matchID string, event *LiveMatchEvent) error {
	// Death events are typically processed as part of kill events
	// This handles standalone death events (e.g., jungle monsters, tower deaths)
	
	if len(event.Participants) > 0 {
		victim := event.Participants[0]
		p.updatePlayerStat(victim, "deaths", 1)
		p.triggerPerformanceUpdate([]string{victim})
	}
	
	log.Printf("✅ Processed death event in match %s", matchID)
	return nil
}

func (p *LiveMatchEventProcessor) processObjectiveEvent(matchID string, event *LiveMatchEvent) error {
	// Extract objective details
	objectiveType := ""
	team := ""
	
	if details := event.Details; details != nil {
		if obj, ok := details["objective_type"].(string); ok {
			objectiveType = obj
		}
		if t, ok := details["team"].(string); ok {
			team = t
		}
	}

	// Update objective statistics for participants
	for _, participant := range event.Participants {
		switch objectiveType {
		case "dragon":
			p.updatePlayerStat(participant, "dragons", 1)
		case "baron":
			p.updatePlayerStat(participant, "barons", 1)
		case "tower":
			p.updatePlayerStat(participant, "towers", 1)
		case "inhibitor":
			p.updatePlayerStat(participant, "inhibitors", 1)
		}
	}

	// Trigger team analytics update
	p.triggerTeamUpdate(matchID, team, objectiveType)
	
	log.Printf("✅ Processed objective event: %s %s taken by %s", objectiveType, matchID, team)
	return nil
}

func (p *LiveMatchEventProcessor) processItemEvent(matchID string, event *LiveMatchEvent) error {
	// Extract item details
	itemID := ""
	playerPUUID := ""
	
	if details := event.Details; details != nil {
		if item, ok := details["item_id"].(string); ok {
			itemID = item
		}
		if player, ok := details["player_puuid"].(string); ok {
			playerPUUID = player
		}
	}

	// Track item builds and gold efficiency
	p.updateItemBuild(playerPUUID, itemID, event.GameTime)
	
	log.Printf("✅ Processed item event: Player %s purchased %s at %d seconds", 
		playerPUUID[:8], itemID, event.GameTime)
	return nil
}

func (p *LiveMatchEventProcessor) processLevelUpEvent(matchID string, event *LiveMatchEvent) error {
	// Extract level up details
	playerPUUID := ""
	newLevel := 0
	
	if details := event.Details; details != nil {
		if player, ok := details["player_puuid"].(string); ok {
			playerPUUID = player
		}
		if level, ok := details["level"].(float64); ok {
			newLevel = int(level)
		}
	}

	// Update player level
	p.updatePlayerStat(playerPUUID, "level", newLevel)
	
	// Check for significant level milestones
	if p.isSignificantLevel(newLevel) {
		p.triggerMilestoneUpdate(playerPUUID, "level", newLevel)
	}
	
	log.Printf("✅ Processed level up: Player %s reached level %d", playerPUUID[:8], newLevel)
	return nil
}

func (p *LiveMatchEventProcessor) processTeamFightEvent(matchID string, event *LiveMatchEvent) error {
	// Extract team fight details
	duration := 0
	location := ""
	outcome := ""
	
	if details := event.Details; details != nil {
		if d, ok := details["duration"].(float64); ok {
			duration = int(d)
		}
		if loc, ok := details["location"].(string); ok {
			location = loc
		}
		if out, ok := details["outcome"].(string); ok {
			outcome = out
		}
	}

	// Analyze team fight performance for all participants
	for _, participant := range event.Participants {
		p.analyzeTeamFightPerformance(participant, event, duration, location, outcome)
	}

	// Trigger team synergy analysis
	p.triggerTeamSynergyUpdate(matchID, event.Participants)
	
	log.Printf("✅ Processed team fight: %ds at %s with outcome %s", duration, location, outcome)
	return nil
}

func (p *LiveMatchEventProcessor) processGameEndEvent(matchID string, event *LiveMatchEvent) error {
	// Extract game end details
	winningTeam := ""
	duration := 0
	
	if details := event.Details; details != nil {
		if team, ok := details["winning_team"].(string); ok {
			winningTeam = team
		}
		if dur, ok := details["duration"].(float64); ok {
			duration = int(dur)
		}
	}

	// Process final match statistics
	for _, participant := range event.Participants {
		p.processFinalStats(participant, matchID, winningTeam, duration)
	}

	// Trigger post-game analysis
	p.triggerPostGameAnalysis(matchID, event.Participants)
	
	log.Printf("✅ Processed game end: Match %s ended, winner: %s, duration: %ds", 
		matchID, winningTeam, duration)
	return nil
}

func (p *LiveMatchEventProcessor) processGenericEvent(matchID string, event *LiveMatchEvent) error {
	// Handle other event types
	log.Printf("📝 Processing generic event: %s in match %s", event.Type, matchID)
	
	// Basic event logging and metrics
	p.updateEventMetrics(event.Type, matchID)
	
	return nil
}

// Helper methods for live match event processing

func (p *LiveMatchEventProcessor) updatePlayerStat(playerPUUID, stat string, value interface{}) {
	// In a real implementation, this would update player statistics in real-time
	log.Printf("📊 Updating %s for player %s: %v", stat, playerPUUID[:8], value)
}

func (p *LiveMatchEventProcessor) updateItemBuild(playerPUUID, itemID string, gameTime int) {
	// Track item build progression
	log.Printf("🛒 Item build update: %s bought %s at %d seconds", playerPUUID[:8], itemID, gameTime)
}

func (p *LiveMatchEventProcessor) triggerPerformanceUpdate(participants []string) {
	// Trigger real-time performance analytics update
	for _, participant := range participants {
		log.Printf("📈 Triggering performance update for %s", participant[:8])
	}
}

func (p *LiveMatchEventProcessor) triggerTeamUpdate(matchID, team, objectiveType string) {
	log.Printf("🏆 Team %s objective update: %s", team, objectiveType)
}

func (p *LiveMatchEventProcessor) triggerMilestoneUpdate(playerPUUID, milestoneType string, value int) {
	log.Printf("🎯 Milestone reached: %s level %d for %s", milestoneType, value, playerPUUID[:8])
}

func (p *LiveMatchEventProcessor) analyzeTeamFightPerformance(participant string, event *LiveMatchEvent, duration int, location, outcome string) {
	log.Printf("⚔️ Team fight analysis: %s performance in %ds fight at %s", participant[:8], duration, location)
}

func (p *LiveMatchEventProcessor) triggerTeamSynergyUpdate(matchID string, participants []string) {
	log.Printf("🤝 Team synergy update for match %s with %d participants", matchID, len(participants))
}

func (p *LiveMatchEventProcessor) processFinalStats(participant, matchID, winningTeam string, duration int) {
	log.Printf("🏁 Final stats for %s in match %s (%ds)", participant[:8], matchID, duration)
}

func (p *LiveMatchEventProcessor) triggerPostGameAnalysis(matchID string, participants []string) {
	log.Printf("📊 Post-game analysis triggered for match %s", matchID)
}

func (p *LiveMatchEventProcessor) updateEventMetrics(eventType, matchID string) {
	log.Printf("📋 Event metrics: %s in %s", eventType, matchID)
}

func (p *LiveMatchEventProcessor) isSignificantLevel(level int) bool {
	// Levels 6, 11, 16 are significant (ultimate levels)
	significantLevels := []int{6, 11, 16}
	for _, sig := range significantLevels {
		if level == sig {
			return true
		}
	}
	return false
}

func (p *LiveMatchEventProcessor) GetEventType() string {
	return "live_match_event"
}

func (p *LiveMatchEventProcessor) GetPriority() int {
	return 8 // High priority for live match events
}

// PlayerUpdateProcessor processes player update events
type PlayerUpdateProcessor struct {
	analyticsEngine *analytics.AnalyticsEngine
}

func (p *PlayerUpdateProcessor) Process(event *StreamEvent) error {
	playerUpdate, ok := event.Data.(*PlayerUpdate)
	if !ok {
		return fmt.Errorf("invalid data type for player update event")
	}

	log.Printf("👤 Processing player update: %s for %s", 
		playerUpdate.UpdateType, event.PlayerPUUID[:8])

	switch playerUpdate.UpdateType {
	case "status":
		return p.processStatusUpdate(playerUpdate)
	case "rank":
		return p.processRankUpdate(playerUpdate)
	case "match_start":
		return p.processMatchStartUpdate(playerUpdate)
	case "match_end":
		return p.processMatchEndUpdate(playerUpdate)
	case "achievement":
		return p.processAchievementUpdate(playerUpdate)
	default:
		return p.processGenericUpdate(playerUpdate)
	}
}

func (p *PlayerUpdateProcessor) processStatusUpdate(update *PlayerUpdate) error {
	if update.StatusUpdate == nil {
		return fmt.Errorf("missing status update data")
	}

	status := update.StatusUpdate
	log.Printf("🟢 Status update: %s is %s (last seen: %s)", 
		update.PlayerPUUID[:8], status.OnlineStatus, status.LastSeen.Format("15:04:05"))

	// Update player cache
	p.updatePlayerStatus(update.PlayerPUUID, status)

	// Notify friends if status changed to online
	if status.OnlineStatus == "online" {
		p.notifyFriends(update.PlayerPUUID, "came_online")
	}

	return nil
}

func (p *PlayerUpdateProcessor) processRankUpdate(update *PlayerUpdate) error {
	if update.RankUpdate == nil {
		return fmt.Errorf("missing rank update data")
	}

	rankUpdate := update.RankUpdate
	log.Printf("📊 Rank update: %s %s → %s (%+d LP)", 
		update.PlayerPUUID[:8], rankUpdate.OldRank, rankUpdate.NewRank, rankUpdate.LPChange)

	// Update rank history
	p.updateRankHistory(update.PlayerPUUID, rankUpdate)

	// Check for rank up/down notifications
	if p.isRankPromotion(rankUpdate) {
		p.sendRankUpNotification(update.PlayerPUUID, rankUpdate)
	} else if p.isRankDemotion(rankUpdate) {
		p.sendRankDownNotification(update.PlayerPUUID, rankUpdate)
	}

	return nil
}

func (p *PlayerUpdateProcessor) processMatchStartUpdate(update *PlayerUpdate) error {
	if update.MatchUpdate == nil {
		return fmt.Errorf("missing match update data")
	}

	matchUpdate := update.MatchUpdate
	log.Printf("🎮 Match start: %s playing %s (%s) in %s", 
		update.PlayerPUUID[:8], matchUpdate.Champion, matchUpdate.Role, matchUpdate.Queue)

	// Update current match status
	p.updateCurrentMatch(update.PlayerPUUID, matchUpdate)

	// Notify friends about match start
	p.notifyFriends(update.PlayerPUUID, "match_started")

	return nil
}

func (p *PlayerUpdateProcessor) processMatchEndUpdate(update *PlayerUpdate) error {
	if update.MatchUpdate == nil {
		return fmt.Errorf("missing match update data")
	}

	matchUpdate := update.MatchUpdate
	log.Printf("🏁 Match end: %s finished with %s (Rating: %.1f)", 
		update.PlayerPUUID[:8], matchUpdate.Result, matchUpdate.Performance.Rating)

	// Update match history
	p.updateMatchHistory(update.PlayerPUUID, matchUpdate)

	// Clear current match status
	p.clearCurrentMatch(update.PlayerPUUID)

	// Check for performance milestones
	if matchUpdate.Performance != nil {
		p.checkPerformanceMilestones(update.PlayerPUUID, matchUpdate.Performance)
	}

	return nil
}

func (p *PlayerUpdateProcessor) processAchievementUpdate(update *PlayerUpdate) error {
	if update.AchievementUpdate == nil {
		return fmt.Errorf("missing achievement update data")
	}

	achievement := update.AchievementUpdate
	log.Printf("🏆 Achievement unlocked: %s earned '%s' (%s)", 
		update.PlayerPUUID[:8], achievement.Name, achievement.Rarity)

	// Update achievement collection
	p.updateAchievements(update.PlayerPUUID, achievement)

	// Send achievement notification
	p.sendAchievementNotification(update.PlayerPUUID, achievement)

	return nil
}

func (p *PlayerUpdateProcessor) processGenericUpdate(update *PlayerUpdate) error {
	log.Printf("📝 Generic player update: %s for %s", 
		update.UpdateType, update.PlayerPUUID[:8])
	return nil
}

// Helper methods for player update processing

func (p *PlayerUpdateProcessor) updatePlayerStatus(playerPUUID string, status *PlayerStatusUpdate) {
	log.Printf("💾 Updating status cache for %s", playerPUUID[:8])
}

func (p *PlayerUpdateProcessor) notifyFriends(playerPUUID, eventType string) {
	log.Printf("👥 Notifying friends of %s: %s", playerPUUID[:8], eventType)
}

func (p *PlayerUpdateProcessor) updateRankHistory(playerPUUID string, rankUpdate *PlayerRankUpdate) {
	log.Printf("📈 Updating rank history: %s %s", playerPUUID[:8], rankUpdate.NewRank)
}

func (p *PlayerUpdateProcessor) isRankPromotion(update *PlayerRankUpdate) bool {
	// Simple check - in production would parse rank tiers properly
	return len(update.NewRank) > len(update.OldRank) || update.LPChange > 0
}

func (p *PlayerUpdateProcessor) isRankDemotion(update *PlayerRankUpdate) bool {
	return len(update.NewRank) < len(update.OldRank) || update.LPChange < -50
}

func (p *PlayerUpdateProcessor) sendRankUpNotification(playerPUUID string, update *PlayerRankUpdate) {
	log.Printf("🎉 Rank up notification: %s promoted to %s", playerPUUID[:8], update.NewRank)
}

func (p *PlayerUpdateProcessor) sendRankDownNotification(playerPUUID string, update *PlayerRankUpdate) {
	log.Printf("😔 Rank down notification: %s demoted to %s", playerPUUID[:8], update.NewRank)
}

func (p *PlayerUpdateProcessor) updateCurrentMatch(playerPUUID string, match *PlayerMatchUpdate) {
	log.Printf("🎯 Current match update: %s in %s", playerPUUID[:8], match.MatchID)
}

func (p *PlayerUpdateProcessor) updateMatchHistory(playerPUUID string, match *PlayerMatchUpdate) {
	log.printf("📊 Match history update: %s finished %s", playerPUUID[:8], match.MatchID)
}

func (p *PlayerUpdateProcessor) clearCurrentMatch(playerPUUID string) {
	log.Printf("🔄 Clearing current match for %s", playerPUUID[:8])
}

func (p *PlayerUpdateProcessor) checkPerformanceMilestones(playerPUUID string, perf *MatchPerformanceSummary) {
	// Check for performance milestones like first pentakill, 10+ KDA, etc.
	if perf.KDA >= 10.0 {
		log.Printf("🌟 Performance milestone: %s achieved %.1f KDA", playerPUUID[:8], perf.KDA)
	}
}

func (p *PlayerUpdateProcessor) updateAchievements(playerPUUID string, achievement *PlayerAchievementUpdate) {
	log.Printf("🏆 Achievement update: %s unlocked %s", playerPUUID[:8], achievement.Name)
}

func (p *PlayerUpdateProcessor) sendAchievementNotification(playerPUUID string, achievement *PlayerAchievementUpdate) {
	log.Printf("🔔 Achievement notification: %s earned %s", playerPUUID[:8], achievement.Name)
}

func (p *PlayerUpdateProcessor) GetEventType() string {
	return "player_update"
}

func (p *PlayerUpdateProcessor) GetPriority() int {
	return 6 // Medium-high priority for player updates
}

// AnalyticsUpdateProcessor processes analytics update events
type AnalyticsUpdateProcessor struct {
	analyticsEngine *analytics.AnalyticsEngine
}

func (p *AnalyticsUpdateProcessor) Process(event *StreamEvent) error {
	analyticsUpdate, ok := event.Data.(*AnalyticsUpdate)
	if !ok {
		return fmt.Errorf("invalid data type for analytics update event")
	}

	log.Printf("📊 Processing analytics update: %s (%s)", 
		analyticsUpdate.Type, analyticsUpdate.Category)

	switch analyticsUpdate.Type {
	case "trend_alert":
		return p.processTrendAlert(analyticsUpdate)
	case "milestone":
		return p.processMilestone(analyticsUpdate)
	case "performance_change":
		return p.processPerformanceChange(analyticsUpdate)
	case "meta_update":
		return p.processMetaUpdate(analyticsUpdate)
	default:
		return p.processGenericAnalytics(analyticsUpdate)
	}
}

func (p *AnalyticsUpdateProcessor) processTrendAlert(update *AnalyticsUpdate) error {
	if update.TrendAlert == nil {
		return fmt.Errorf("missing trend alert data")
	}

	alert := update.TrendAlert
	log.Printf("📈 Trend alert: %s %s %s by %.2f (was %.2f, now %.2f)", 
		alert.PlayerPUUID[:8], alert.Metric, alert.Direction, 
		alert.Change, alert.OldValue, alert.NewValue)

	// Send trend notification if significant
	if alert.Magnitude == "significant" || alert.Magnitude == "dramatic" {
		p.sendTrendNotification(alert)
	}

	return nil
}

func (p *AnalyticsUpdateProcessor) processMilestone(update *AnalyticsUpdate) error {
	if update.Milestone == nil {
		return fmt.Errorf("missing milestone data")
	}

	milestone := update.Milestone
	log.Printf("🎯 Milestone: %s reached %.1f%% progress on %s (%s)", 
		milestone.PlayerPUUID[:8], milestone.Progress*100, 
		milestone.Description, milestone.Rarity)

	// Send milestone notification if near completion or completed
	if milestone.Progress >= 0.9 {
		p.sendMilestoneNotification(milestone)
	}

	return nil
}

func (p *AnalyticsUpdateProcessor) processPerformanceChange(update *AnalyticsUpdate) error {
	if update.PerformanceChange == nil {
		return fmt.Errorf("missing performance change data")
	}

	perfChange := update.PerformanceChange
	log.Printf("🔄 Performance change: %s %s changed from %.1f to %.1f (%s %s)", 
		perfChange.PlayerPUUID[:8], perfChange.PerformanceMetric,
		perfChange.OldRating, perfChange.NewRating, 
		perfChange.ChangeType, perfChange.Significance)

	// Send performance change notification for major changes
	if perfChange.Significance == "major" {
		p.sendPerformanceChangeNotification(perfChange)
	}

	return nil
}

func (p *AnalyticsUpdateProcessor) processMetaUpdate(update *AnalyticsUpdate) error {
	if update.MetaUpdate == nil {
		return fmt.Errorf("missing meta update data")
	}

	metaUpdate := update.MetaUpdate
	log.Printf("🌍 Meta update: %s %s patch %s", 
		metaUpdate.Region, metaUpdate.Queue, metaUpdate.Patch)

	// Broadcast meta changes to interested players
	p.broadcastMetaUpdate(metaUpdate)

	return nil
}

func (p *AnalyticsUpdateProcessor) processGenericAnalytics(update *AnalyticsUpdate) error {
	log.Printf("📝 Generic analytics update: %s (%s)", update.Type, update.Category)
	return nil
}

// Helper methods for analytics processing

func (p *AnalyticsUpdateProcessor) sendTrendNotification(alert *TrendAlert) {
	log.Printf("🔔 Trend notification: %s trend alert sent", alert.PlayerPUUID[:8])
}

func (p *AnalyticsUpdateProcessor) sendMilestoneNotification(milestone *MilestoneUpdate) {
	log.Printf("🎯 Milestone notification: %s milestone alert sent", milestone.PlayerPUUID[:8])
}

func (p *AnalyticsUpdateProcessor) sendPerformanceChangeNotification(change *PerformanceChangeUpdate) {
	log.Printf("📊 Performance notification: %s performance alert sent", change.PlayerPUUID[:8])
}

func (p *AnalyticsUpdateProcessor) broadcastMetaUpdate(metaUpdate *MetaAnalyticsUpdate) {
	log.Printf("🌍 Broadcasting meta update for %s %s", metaUpdate.Region, metaUpdate.Queue)
}

func (p *AnalyticsUpdateProcessor) GetEventType() string {
	return "analytics_update"
}

func (p *AnalyticsUpdateProcessor) GetPriority() int {
	return 4 // Medium priority for analytics updates
}

// Fix missing import for sync
import "sync"