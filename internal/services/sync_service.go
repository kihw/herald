package services

import (
	"database/sql"
	"fmt"
	"time"
	"log"
	"encoding/json"
	"lol-match-exporter/internal/db"
)

// SyncEvent represents a synchronization event
type SyncEvent struct {
	Type      string                 `json:"type"`
	UserID    int                    `json:"user_id"`
	MatchID   string                 `json:"match_id,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// SyncService handles synchronization with Riot API and analytics processing
type SyncService struct {
	db                  *sql.DB
	riotService         *RiotService
	analyticsService    *AnalyticsService
	notificationService *NotificationService
	isAutoSyncEnabled   bool
	syncInterval        time.Duration
	subscribers         []chan SyncEvent
	stopChan            chan bool
	processingQueue     chan AnalyticsJob
}

// AnalyticsJob represents an analytics processing job
type AnalyticsJob struct {
	UserID   int
	MatchID  string
	JobType  string
	Data     map[string]interface{}
}

// NewSyncService creates a new sync service
func NewSyncService(database *db.Database, riotService *RiotService, analyticsService *AnalyticsService, notificationService *NotificationService) *SyncService {
	return &SyncService{
		db:                  database.DB,
		riotService:         riotService,
		analyticsService:    analyticsService,
		notificationService: notificationService,
		isAutoSyncEnabled:   false,
		syncInterval:        30 * time.Minute,
		subscribers:         make([]chan SyncEvent, 0),
		stopChan:            make(chan bool),
		processingQueue:     make(chan AnalyticsJob, 100),
	}
}

// SyncUserMatches synchronizes matches for a user
func (ss *SyncService) SyncUserMatches(userID int, username, tagline, region string) error {
	if !ss.riotService.IsConfigured() {
		return fmt.Errorf("Riot API not configured")
	}

	// Créer un job de synchronisation
	jobID, err := ss.createSyncJob(userID, "manual")
	if err != nil {
		return fmt.Errorf("failed to create sync job: %w", err)
	}

	// Marquer le job comme démarré
	err = ss.updateSyncJobStatus(jobID, "running", "")
	if err != nil {
		return fmt.Errorf("failed to update sync job status: %w", err)
	}

	// Effectuer la synchronisation
	processedCount, err := ss.performSync(userID, username, tagline, region)
	if err != nil {
		// Marquer le job comme échoué
		ss.updateSyncJobStatus(jobID, "failed", err.Error())
		return err
	}

	// Marquer le job comme réussi
	err = ss.updateSyncJobStatus(jobID, "completed", "")
	if err != nil {
		return fmt.Errorf("failed to update sync job status: %w", err)
	}

	// Mettre à jour la date de dernière synchronisation de l'utilisateur
	err = ss.updateUserLastSync(userID)
	if err != nil {
		return fmt.Errorf("failed to update user last sync: %w", err)
	}

	// Process analytics for new matches
	if ss.analyticsService != nil {
		go ss.processUserAnalytics(userID)
	}

	// Notify subscribers
	ss.notifySubscribers(SyncEvent{
		Type:      "sync_completed",
		UserID:    userID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"matches_processed": processedCount,
		},
	})

	return nil
}

// performSync effectue la synchronisation réelle
func (ss *SyncService) performSync(userID int, username, tagline, region string) (int, error) {
	// 1. Récupérer les informations du compte Riot
	account, err := ss.riotService.GetAccountByRiotID(username, tagline, region)
	if err != nil {
		return 0, fmt.Errorf("failed to get Riot account: %w", err)
	}

	// 2. Récupérer les informations de l'invocateur
	summoner, err := ss.riotService.GetSummonerByPUUID(account.PUUID, region)
	if err != nil {
		return 0, fmt.Errorf("failed to get summoner info: %w", err)
	}

	// 3. Mettre à jour les informations de l'utilisateur
	err = ss.updateUserRiotInfo(userID, account.PUUID, summoner.ID, summoner.AccountID, 
		summoner.ProfileIconID, summoner.SummonerLevel)
	if err != nil {
		return 0, fmt.Errorf("failed to update user Riot info: %w", err)
	}

	// 4. Récupérer la liste des derniers matches (20 matches)
	matchIDs, err := ss.riotService.GetMatchListByPUUID(account.PUUID, region, 0, 20)
	if err != nil {
		return 0, fmt.Errorf("failed to get match list: %w", err)
	}

	// 5. Traiter chaque match
	processedCount := 0
	for _, matchID := range matchIDs {
		// Vérifier si le match existe déjà
		exists, err := ss.matchExists(matchID)
		if err != nil {
			return 0, fmt.Errorf("failed to check if match exists: %w", err)
		}

		if exists {
			continue // Skip ce match s'il existe déjà
		}

		// Récupérer les détails du match
		matchInfo, err := ss.riotService.GetMatchByID(matchID, region)
		if err != nil {
			fmt.Printf("Warning: failed to get match %s: %v\n", matchID, err)
			continue // Continue avec le prochain match en cas d'erreur
		}

		// Sauvegarder le match et les participants
		err = ss.saveMatch(matchInfo, account.PUUID)
		if err != nil {
			fmt.Printf("Warning: failed to save match %s: %v\n", matchID, err)
			continue
		}
		
		processedCount++
	}

	return processedCount, nil
}

// createSyncJob crée un nouveau job de synchronisation
func (ss *SyncService) createSyncJob(userID int, triggerType string) (int, error) {
	query := `
		INSERT INTO sync_jobs (user_id, status, trigger_type, started_at, created_at)
		VALUES ($1, 'pending', $2, NOW(), NOW())
		RETURNING id`
	
	var jobID int
	err := ss.db.QueryRow(query, userID, triggerType).Scan(&jobID)
	if err != nil {
		return 0, err
	}

	return jobID, nil
}

// updateSyncJobStatus met à jour le statut d'un job de synchronisation
func (ss *SyncService) updateSyncJobStatus(jobID int, status, errorMsg string) error {
	var query string
	var args []interface{}

	if status == "completed" || status == "failed" {
		query = `
			UPDATE sync_jobs 
			SET status = $1, error_message = $2, completed_at = NOW(), updated_at = NOW()
			WHERE id = $3`
		args = []interface{}{status, errorMsg, jobID}
	} else {
		query = `
			UPDATE sync_jobs 
			SET status = $1, error_message = $2, updated_at = NOW()
			WHERE id = $3`
		args = []interface{}{status, errorMsg, jobID}
	}

	_, err := ss.db.Exec(query, args...)
	return err
}

// updateUserLastSync met à jour la date de dernière synchronisation
func (ss *SyncService) updateUserLastSync(userID int) error {
	query := `UPDATE users SET last_sync = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := ss.db.Exec(query, userID)
	return err
}

// updateUserRiotInfo met à jour les informations Riot de l'utilisateur
func (ss *SyncService) updateUserRiotInfo(userID int, puuid, summonerID, accountID string, profileIconID, summonerLevel int) error {
	query := `
		UPDATE users 
		SET riot_puuid = $1, summoner_id = $2, account_id = $3, 
		    profile_icon_id = $4, summoner_level = $5, updated_at = NOW()
		WHERE id = $6`
	
	_, err := ss.db.Exec(query, puuid, summonerID, accountID, profileIconID, summonerLevel, userID)
	return err
}

// matchExists vérifie si un match existe déjà dans la base de données
func (ss *SyncService) matchExists(matchID string) (bool, error) {
	query := `SELECT COUNT(*) FROM matches WHERE match_id = $1`
	var count int
	err := ss.db.QueryRow(query, matchID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// saveMatch sauvegarde un match et ses participants
func (ss *SyncService) saveMatch(matchInfo *MatchInfo, userPUUID string) error {
	// Commencer une transaction
	tx, err := ss.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insérer le match
	matchQuery := `
		INSERT INTO matches (
			match_id, game_creation, game_duration, game_mode, game_type, 
			game_version, map_id, platform_id, queue_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		RETURNING id`
	
	var matchDBID int
	gameCreation := time.Unix(matchInfo.GameCreation/1000, 0)
	
	err = tx.QueryRow(matchQuery, matchInfo.MatchID, gameCreation, matchInfo.GameDuration,
		matchInfo.GameMode, matchInfo.GameType, matchInfo.GameVersion,
		matchInfo.MapID, matchInfo.PlatformID, matchInfo.QueueID).Scan(&matchDBID)
	if err != nil {
		return err
	}

	// Insérer les participants
	participantQuery := `
		INSERT INTO match_participants (
			match_id, participant_id, puuid, champion_id, champion_name, team_id,
			position, kills, deaths, assists, gold_earned, total_minions_killed,
			vision_score, damage_dealt_to_champions, win, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW())`

	for _, participant := range matchInfo.Participants {
		_, err = tx.Exec(participantQuery,
			matchDBID, participant.ParticipantID, participant.PUUID,
			participant.ChampionID, participant.ChampionName, participant.TeamID,
			participant.TeamPosition, participant.Kills, participant.Deaths, participant.Assists,
			participant.GoldEarned, participant.TotalMinionsKilled, participant.VisionScore,
			participant.TotalDamageDealtToChampions, participant.Win)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	// Trigger analytics processing for this match
	go ss.queueAnalyticsProcessing(userPUUID, matchInfo.MatchID, matchInfo)

	return nil
}

// Analytics Processing Methods

// StartAnalyticsProcessor starts the background analytics processing worker
func (ss *SyncService) StartAnalyticsProcessor() {
	if ss.analyticsService == nil {
		return
	}

	go func() {
		log.Println("Analytics processor started")
		for {
			select {
			case job := <-ss.processingQueue:
				ss.processAnalyticsJob(job)
			case <-ss.stopChan:
				log.Println("Analytics processor stopped")
				return
			}
		}
	}()
}

// StopAnalyticsProcessor stops the background analytics processing
func (ss *SyncService) StopAnalyticsProcessor() {
	ss.stopChan <- true
}

// queueAnalyticsProcessing queues a match for analytics processing
func (ss *SyncService) queueAnalyticsProcessing(userPUUID, matchID string, matchInfo *MatchInfo) {
	if ss.analyticsService == nil {
		return
	}

	// Get user ID from PUUID
	userID, err := ss.getUserIDByPUUID(userPUUID)
	if err != nil {
		log.Printf("Failed to get user ID for PUUID %s: %v", userPUUID, err)
		return
	}

	// Create analytics job
	job := AnalyticsJob{
		UserID:  userID,
		MatchID: matchID,
		JobType: "match_analysis",
		Data: map[string]interface{}{
			"match_info": matchInfo,
			"user_puuid": userPUUID,
		},
	}

	// Queue the job
	select {
	case ss.processingQueue <- job:
		log.Printf("Queued analytics processing for match %s", matchID)
	default:
		log.Printf("Analytics processing queue full, skipping match %s", matchID)
	}
}

// processAnalyticsJob processes a single analytics job
func (ss *SyncService) processAnalyticsJob(job AnalyticsJob) {
	log.Printf("Processing analytics job for user %d, match %s", job.UserID, job.MatchID)

	switch job.JobType {
	case "match_analysis":
		ss.processMatchAnalytics(job)
	case "user_summary":
		ss.processUserSummaryAnalytics(job)
	default:
		log.Printf("Unknown analytics job type: %s", job.JobType)
	}
}

// processMatchAnalytics processes analytics for a single match
func (ss *SyncService) processMatchAnalytics(job AnalyticsJob) {
	// Update champion stats
	err := ss.updateChampionStats(job.UserID, job.MatchID)
	if err != nil {
		log.Printf("Failed to update champion stats for match %s: %v", job.MatchID, err)
	}

	// Update MMR estimate
	err = ss.updateMMREstimate(job.UserID, job.MatchID)
	if err != nil {
		log.Printf("Failed to update MMR estimate for match %s: %v", job.MatchID, err)
	}

	// Generate recommendations if needed
	err = ss.updateRecommendations(job.UserID)
	if err != nil {
		log.Printf("Failed to update recommendations for user %d: %v", job.UserID, err)
	}

	// Generate insights if notification service is available
	if ss.notificationService != nil {
		ss.generateMatchInsights(job)
	}

	// Notify subscribers
	ss.notifySubscribers(SyncEvent{
		Type:      "match_analyzed",
		UserID:    job.UserID,
		MatchID:   job.MatchID,
		Timestamp: time.Now(),
	})
}

// processUserAnalytics processes analytics for all user matches
func (ss *SyncService) processUserAnalytics(userID int) {
	log.Printf("Processing user analytics for user %d", userID)

	// Queue user summary job
	job := AnalyticsJob{
		UserID:  userID,
		JobType: "user_summary",
	}

	select {
	case ss.processingQueue <- job:
		log.Printf("Queued user summary analytics for user %d", userID)
	default:
		log.Printf("Analytics processing queue full, skipping user summary for user %d", userID)
	}
}

// processUserSummaryAnalytics processes summary analytics for a user
func (ss *SyncService) processUserSummaryAnalytics(job AnalyticsJob) {
	// Update performance metrics
	err := ss.updatePerformanceMetrics(job.UserID)
	if err != nil {
		log.Printf("Failed to update performance metrics for user %d: %v", job.UserID, err)
	}

	// Update recommendations
	err = ss.updateRecommendations(job.UserID)
	if err != nil {
		log.Printf("Failed to update recommendations for user %d: %v", job.UserID, err)
	}

	// Notify subscribers
	ss.notifySubscribers(SyncEvent{
		Type:      "user_analytics_updated",
		UserID:    job.UserID,
		Timestamp: time.Now(),
	})
}

// Helper methods for analytics processing

// getUserIDByPUUID gets user ID from PUUID
func (ss *SyncService) getUserIDByPUUID(puuid string) (int, error) {
	query := `SELECT id FROM users WHERE riot_puuid = $1`
	var userID int
	err := ss.db.QueryRow(query, puuid).Scan(&userID)
	return userID, err
}

// updateChampionStats updates champion statistics
func (ss *SyncService) updateChampionStats(userID int, matchID string) error {
	// Get match participant data
	query := `
		SELECT mp.champion_id, mp.champion_name, mp.kills, mp.deaths, mp.assists,
			   mp.gold_earned, mp.total_minions_killed, mp.vision_score,
			   mp.damage_dealt_to_champions, mp.win
		FROM match_participants mp
		JOIN matches m ON m.id = mp.match_id
		WHERE m.match_id = $1 AND mp.puuid = (
			SELECT riot_puuid FROM users WHERE id = $2
		)`

	rows, err := ss.db.Query(query, matchID, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var championID int
		var championName string
		var kills, deaths, assists, goldEarned, cs, visionScore, damage int
		var win bool

		err = rows.Scan(&championID, &championName, &kills, &deaths, &assists,
			&goldEarned, &cs, &visionScore, &damage, &win)
		if err != nil {
			continue
		}

		// Update or create champion stats
		err = ss.upsertChampionStat(userID, championID, championName, kills, deaths, assists,
			goldEarned, cs, visionScore, damage, win)
		if err != nil {
			log.Printf("Failed to upsert champion stat: %v", err)
		}
	}

	return nil
}

// upsertChampionStat updates or inserts champion statistics
func (ss *SyncService) upsertChampionStat(userID, championID int, championName string,
	kills, deaths, assists, goldEarned, cs, visionScore, damage int, win bool) error {

	query := `
		INSERT INTO champion_stats (
			user_id, champion_id, champion_name, games_played, wins, losses,
			kills, deaths, assists, cs_total, gold_earned, damage_dealt,
			vision_score, last_played, created_at, updated_at
		) VALUES ($1, $2, $3, 1, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW(), NOW())
		ON CONFLICT (user_id, champion_id) DO UPDATE SET
			games_played = champion_stats.games_played + 1,
			wins = champion_stats.wins + $4,
			losses = champion_stats.losses + $5,
			kills = champion_stats.kills + $6,
			deaths = champion_stats.deaths + $7,
			assists = champion_stats.assists + $8,
			cs_total = champion_stats.cs_total + $9,
			gold_earned = champion_stats.gold_earned + $10,
			damage_dealt = champion_stats.damage_dealt + $11,
			vision_score = champion_stats.vision_score + $12,
			last_played = NOW(),
			updated_at = NOW()`

	winCount := 0
	lossCount := 1
	if win {
		winCount = 1
		lossCount = 0
	}

	_, err := ss.db.Exec(query, userID, championID, championName, winCount, lossCount,
		kills, deaths, assists, cs, goldEarned, damage, visionScore)
	return err
}

// updateMMREstimate updates MMR estimate using analytics service
func (ss *SyncService) updateMMREstimate(userID int, matchID string) error {
	if ss.analyticsService == nil {
		return nil
	}

	// Call MMR calculator
	mmrData, err := ss.analyticsService.GetMMRTrajectory(userID, 30)
	if err != nil {
		return err
	}

	// Extract latest MMR estimate
	if len(mmrData.MMRHistory) == 0 {
		return nil
	}

	latest := mmrData.MMRHistory[len(mmrData.MMRHistory)-1]
	
	// Insert into MMR history
	query := `
		INSERT INTO mmr_history (user_id, estimated_mmr, mmr_change, rank_estimate, confidence, match_id)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = ss.db.Exec(query, userID, latest.EstimatedMMR, latest.MMRChange,
		latest.RankEstimate, latest.Confidence, latest.MatchID)
	return err
}

// updateRecommendations updates recommendations using analytics service
func (ss *SyncService) updateRecommendations(userID int) error {
	if ss.analyticsService == nil {
		return nil
	}

	// Get recommendations from analytics service
	recommendations, err := ss.analyticsService.GetRecommendations(userID)
	if err != nil {
		return err
	}

	// Clear old recommendations
	_, err = ss.db.Exec(`DELETE FROM recommendations WHERE user_id = $1 AND expires_at < NOW()`, userID)
	if err != nil {
		return err
	}

	// Insert new recommendations
	for _, rec := range recommendations {
		query := `
			INSERT INTO recommendations (
				user_id, type, title, description, priority, confidence,
				expected_improvement, action_items, champion_id, role,
				time_period, expires_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW() + INTERVAL '7 days')
			ON CONFLICT (user_id, type, title) DO UPDATE SET
				description = EXCLUDED.description,
				priority = EXCLUDED.priority,
				confidence = EXCLUDED.confidence,
				expected_improvement = EXCLUDED.expected_improvement,
				action_items = EXCLUDED.action_items,
				expires_at = EXCLUDED.expires_at`

		actionItemsJSON, _ := json.Marshal(rec.ActionItems)
		
		_, err = ss.db.Exec(query, userID, rec.Type, rec.Title, rec.Description,
			rec.Priority, rec.Confidence, rec.ExpectedImprovement, actionItemsJSON,
			rec.ChampionID, rec.Role, rec.TimePeriod)
		if err != nil {
			log.Printf("Failed to insert recommendation: %v", err)
		}
	}

	return nil
}

// updatePerformanceMetrics updates performance metrics
func (ss *SyncService) updatePerformanceMetrics(userID int) error {
	periods := []string{"today", "week", "month", "season"}
	
	for _, period := range periods {
		err := ss.updatePerformanceMetricForPeriod(userID, period)
		if err != nil {
			log.Printf("Failed to update performance metrics for period %s: %v", period, err)
		}
	}
	
	return nil
}

// updatePerformanceMetricForPeriod updates performance metrics for a specific period
func (ss *SyncService) updatePerformanceMetricForPeriod(userID int, period string) error {
	// Get period stats from analytics service
	if ss.analyticsService == nil {
		return nil
	}

	stats, err := ss.analyticsService.GetPeriodStats(userID, period)
	if err != nil {
		return err
	}

	// Upsert performance metrics
	query := `
		INSERT INTO performance_metrics (
			user_id, period, total_games, win_rate, avg_kda, avg_cs_per_min,
			avg_gold_per_min, avg_damage_per_min, avg_vision_score,
			best_role, worst_role, performance_score, trend_direction
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (user_id, period) DO UPDATE SET
			total_games = EXCLUDED.total_games,
			win_rate = EXCLUDED.win_rate,
			avg_kda = EXCLUDED.avg_kda,
			avg_cs_per_min = EXCLUDED.avg_cs_per_min,
			avg_gold_per_min = EXCLUDED.avg_gold_per_min,
			avg_damage_per_min = EXCLUDED.avg_damage_per_min,
			avg_vision_score = EXCLUDED.avg_vision_score,
			best_role = EXCLUDED.best_role,
			worst_role = EXCLUDED.worst_role,
			performance_score = EXCLUDED.performance_score,
			trend_direction = EXCLUDED.trend_direction,
			recorded_at = NOW()`

	_, err = ss.db.Exec(query, userID, period, stats.TotalGames, stats.WinRate,
		stats.AvgKDA, 6.5, 350.0, 500.0, // Placeholder values for CS, Gold, Damage per min
		20.0, stats.BestRole, stats.WorstRole, 75.0, // Placeholder vision score and performance
		stats.RecentTrend)
	
	return err
}

// Event notification methods

// Subscribe adds a subscriber to sync events
func (ss *SyncService) Subscribe() chan SyncEvent {
	ch := make(chan SyncEvent, 10)
	ss.subscribers = append(ss.subscribers, ch)
	return ch
}

// notifySubscribers notifies all subscribers of an event
func (ss *SyncService) notifySubscribers(event SyncEvent) {
	for _, subscriber := range ss.subscribers {
		select {
		case subscriber <- event:
		default:
			// Skip if channel is full
		}
	}
}

// Auto-sync methods

// EnableAutoSync enables automatic synchronization
func (ss *SyncService) EnableAutoSync(interval time.Duration) {
	ss.isAutoSyncEnabled = true
	ss.syncInterval = interval
	
	go ss.autoSyncWorker()
}

// DisableAutoSync disables automatic synchronization
func (ss *SyncService) DisableAutoSync() {
	ss.isAutoSyncEnabled = false
}

// autoSyncWorker runs automatic synchronization
func (ss *SyncService) autoSyncWorker() {
	ticker := time.NewTicker(ss.syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !ss.isAutoSyncEnabled {
				return
			}
			ss.performAutoSync()
		case <-ss.stopChan:
			return
		}
	}
}

// performAutoSync performs automatic sync for active users
func (ss *SyncService) performAutoSync() {
	log.Println("Starting automatic sync")

	// Get users who need sync (last_sync > 1 hour ago)
	query := `
		SELECT id, riot_id, riot_tag, region 
		FROM users 
		WHERE last_sync < NOW() - INTERVAL '1 hour' OR last_sync IS NULL
		LIMIT 10`

	rows, err := ss.db.Query(query)
	if err != nil {
		log.Printf("Failed to get users for auto sync: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		var riotID, riotTag, region string

		err = rows.Scan(&userID, &riotID, &riotTag, &region)
		if err != nil {
			continue
		}

		// Sync user matches
		go func(userID int, riotID, riotTag, region string) {
			err := ss.SyncUserMatches(userID, riotID, riotTag, region)
			if err != nil {
				log.Printf("Auto sync failed for user %d: %v", userID, err)
			}
		}(userID, riotID, riotTag, region)
	}
}

// Insight generation methods

// generateMatchInsights generates insights after match processing
func (ss *SyncService) generateMatchInsights(job AnalyticsJob) {
	// Get match analysis data
	matchData := ss.getMatchAnalysisData(job.UserID, job.MatchID)
	
	// Process match insights
	if ss.notificationService != nil {
		ss.notificationService.ProcessMatchInsights(job.UserID, job.MatchID, matchData)
	}
}

// getMatchAnalysisData retrieves match analysis data for insight generation
func (ss *SyncService) getMatchAnalysisData(userID int, matchID string) map[string]interface{} {
	data := make(map[string]interface{})
	
	// Get recent performance trend
	performanceData := ss.getPerformanceTrend(userID)
	if performanceData != nil {
		data["performance"] = performanceData
	}
	
	// Get streak information
	streakData := ss.getStreakData(userID)
	if streakData != nil {
		data["streak"] = streakData
	}
	
	// Get champion performance
	championData := ss.getChampionPerformance(userID, matchID)
	if championData != nil {
		data["champion"] = championData
	}
	
	// Get MMR changes
	mmrData := ss.getMMRChanges(userID)
	if mmrData != nil {
		data["mmr"] = mmrData
		
		// Process MMR insights separately
		if ss.notificationService != nil {
			ss.notificationService.ProcessMMRInsights(userID, mmrData)
		}
	}
	
	return data
}

// getPerformanceTrend calculates recent performance trends
func (ss *SyncService) getPerformanceTrend(userID int) map[string]interface{} {
	// Get recent performance metrics
	query := `
		SELECT performance_score, recorded_at
		FROM performance_metrics 
		WHERE user_id = $1 AND period = 'week'
		ORDER BY recorded_at DESC 
		LIMIT 5`
	
	rows, err := ss.db.Query(query, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	
	var scores []float64
	for rows.Next() {
		var score float64
		var recordedAt time.Time
		
		err = rows.Scan(&score, &recordedAt)
		if err != nil {
			continue
		}
		
		scores = append(scores, score)
	}
	
	if len(scores) < 2 {
		return nil
	}
	
	// Calculate trend
	latest := scores[0]
	previous := scores[1]
	change := (latest - previous) / previous
	
	result := map[string]interface{}{
		"latest_score": latest,
		"previous_score": previous,
		"change": change,
	}
	
	if change > 0.15 {
		result["score_improvement"] = change
	} else if change < -0.20 {
		result["score_drop"] = -change
	}
	
	return result
}

// getStreakData calculates win/loss streaks
func (ss *SyncService) getStreakData(userID int) map[string]interface{} {
	// Get recent match results
	query := `
		SELECT mp.win
		FROM match_participants mp
		JOIN matches m ON m.id = mp.match_id
		WHERE mp.puuid = (SELECT riot_puuid FROM users WHERE id = $1)
		ORDER BY m.game_creation DESC
		LIMIT 10`
	
	rows, err := ss.db.Query(query, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	
	var results []bool
	for rows.Next() {
		var win bool
		err = rows.Scan(&win)
		if err != nil {
			continue
		}
		results = append(results, win)
	}
	
	if len(results) == 0 {
		return nil
	}
	
	// Calculate current streak
	currentWin := results[0]
	streak := 1
	
	for i := 1; i < len(results); i++ {
		if results[i] == currentWin {
			streak++
		} else {
			break
		}
	}
	
	data := map[string]interface{}{}
	
	if currentWin {
		data["win_streak"] = streak
	} else {
		data["loss_streak"] = streak
	}
	
	return data
}

// getChampionPerformance gets champion performance data for the match
func (ss *SyncService) getChampionPerformance(userID int, matchID string) map[string]interface{} {
	// Get champion info from the match
	query := `
		SELECT mp.champion_name, cs.win_rate
		FROM match_participants mp
		JOIN matches m ON m.id = mp.match_id
		LEFT JOIN (
			SELECT champion_id, 
				   CASE WHEN games_played > 0 THEN CAST(wins AS FLOAT) / games_played ELSE 0 END as win_rate
			FROM champion_stats 
			WHERE user_id = $1
		) cs ON cs.champion_id = mp.champion_id
		WHERE m.match_id = $2 AND mp.puuid = (
			SELECT riot_puuid FROM users WHERE id = $1
		)`
	
	var championName string
	var winRate sql.NullFloat64
	
	err := ss.db.QueryRow(query, userID, matchID).Scan(&championName, &winRate)
	if err != nil {
		return nil
	}
	
	data := map[string]interface{}{
		"name": championName,
	}
	
	if winRate.Valid {
		data["win_rate"] = winRate.Float64
	}
	
	return data
}

// getMMRChanges gets recent MMR changes
func (ss *SyncService) getMMRChanges(userID int) map[string]interface{} {
	// Get recent MMR history
	query := `
		SELECT estimated_mmr, mmr_change, rank_estimate, recorded_at
		FROM mmr_history 
		WHERE user_id = $1
		ORDER BY recorded_at DESC 
		LIMIT 5`
	
	rows, err := ss.db.Query(query, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	
	var mmrEntries []map[string]interface{}
	var totalChange float64
	var latestRank string
	
	for rows.Next() {
		var mmr, change int
		var rank string
		var recordedAt time.Time
		
		err = rows.Scan(&mmr, &change, &rank, &recordedAt)
		if err != nil {
			continue
		}
		
		entry := map[string]interface{}{
			"mmr": mmr,
			"change": change,
			"rank": rank,
			"timestamp": recordedAt,
		}
		
		mmrEntries = append(mmrEntries, entry)
		totalChange += float64(change)
		
		if latestRank == "" {
			latestRank = rank
		}
	}
	
	if len(mmrEntries) == 0 {
		return nil
	}
	
	data := map[string]interface{}{
		"recent_change": totalChange,
		"current_rank": latestRank,
		"history": mmrEntries,
	}
	
	// Check for rank promotion (simplified logic)
	if len(mmrEntries) >= 2 {
		currentRank := mmrEntries[0]["rank"].(string)
		previousRank := mmrEntries[1]["rank"].(string)
		
		if currentRank != previousRank {
			data["rank_promotion"] = true
			data["new_rank"] = currentRank
		}
	}
	
	return data
}
