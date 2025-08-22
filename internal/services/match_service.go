package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"lol-match-exporter/internal/models"
)

// MatchService gère les opérations liées aux matches
type MatchService struct {
	db             *sql.DB
	riotAPIService *RiotAPIService
}

// NewMatchService crée une nouvelle instance du service match
func NewMatchService(db *sql.DB, riotAPIService *RiotAPIService) *MatchService {
	return &MatchService{
		db:             db,
		riotAPIService: riotAPIService,
	}
}

// SyncUserMatches synchronise les matches d'un utilisateur avec l'API Riot
func (ms *MatchService) SyncUserMatches(userID int, count int) (*models.SyncJob, error) {
	// Récupérer l'utilisateur
	user, err := ms.getUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Créer un job de synchronisation
	syncJob, err := ms.createSyncJob(userID, "match_sync")
	if err != nil {
		return nil, fmt.Errorf("failed to create sync job: %w", err)
	}

	// Lancer la synchronisation en arrière-plan
	go func() {
		err := ms.performMatchSync(syncJob, user, count)
		if err != nil {
			log.Printf("Match sync failed for user %d: %v", userID, err)
			ms.updateSyncJobStatus(syncJob.ID, "failed", err.Error())
		}
	}()

	return syncJob, nil
}

// performMatchSync effectue la synchronisation des matches
func (ms *MatchService) performMatchSync(syncJob *models.SyncJob, user *models.User, count int) error {
	// Marquer le job comme démarré
	ms.updateSyncJobStatus(syncJob.ID, "running", "")

	if ms.riotAPIService == nil {
		return fmt.Errorf("Riot API service not available")
	}

	// Récupérer les IDs des matches depuis l'API Riot
	matchIDs, err := ms.riotAPIService.GetMatchHistory(user.RiotPUUID, count, nil)
	if err != nil {
		return fmt.Errorf("failed to get match history: %w", err)
	}

	newMatches := 0
	updatedMatches := 0
	processedMatches := 0

	for _, matchID := range matchIDs {
		// Vérifier si le match existe déjà
		exists, err := ms.matchExists(matchID)
		if err != nil {
			log.Printf("Error checking match existence for %s: %v", matchID, err)
			continue
		}

		if !exists {
			// Récupérer les détails du match
			riotMatch, err := ms.riotAPIService.GetMatch(matchID)
			if err != nil {
				log.Printf("Error fetching match %s: %v", matchID, err)
				continue
			}

			// Sauvegarder le match en base
			err = ms.saveMatch(riotMatch, user)
			if err != nil {
				log.Printf("Error saving match %s: %v", matchID, err)
				continue
			}

			newMatches++
		}

		processedMatches++

		// Pause pour respecter le rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	// Mettre à jour le job avec les résultats
	err = ms.completeSyncJob(syncJob.ID, processedMatches, newMatches, updatedMatches)
	if err != nil {
		return fmt.Errorf("failed to complete sync job: %w", err)
	}

	// Mettre à jour last_sync de l'utilisateur
	ms.updateUserLastSync(user.ID)

	return nil
}

// saveMatch sauvegarde un match et ses participants en base
func (ms *MatchService) saveMatch(riotMatch *RiotMatch, user *models.User) error {
	tx, err := ms.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Sérialiser les données complètes du match
	rawDataJSON, _ := json.Marshal(riotMatch)

	// Insérer le match
	matchQuery := `
		INSERT OR IGNORE INTO matches (
			match_id, platform, game_creation, game_duration, game_end_timestamp,
			game_mode, game_type, game_version, map_id, queue_id, season_id,
			tournament_code, data_version, raw_data
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var gameEndTimestamp *int64
	if riotMatch.Info.GameDuration > 0 {
		endTime := riotMatch.Info.GameCreation + int64(riotMatch.Info.GameDuration*1000)
		gameEndTimestamp = &endTime
	}

	result, err := tx.Exec(matchQuery,
		riotMatch.Metadata.MatchID, "EUW1", riotMatch.Info.GameCreation/1000,
		riotMatch.Info.GameDuration, gameEndTimestamp, riotMatch.Info.GameMode,
		riotMatch.Info.GameType, "14.21.1", 11, riotMatch.Info.QueueID, 14,
		nil, "2", string(rawDataJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to insert match: %w", err)
	}

	matchDBID, err := result.LastInsertId()
	if err != nil {
		// Le match existe déjà, récupérer son ID
		var existingID int64
		err = tx.QueryRow("SELECT id FROM matches WHERE match_id = ?", riotMatch.Metadata.MatchID).Scan(&existingID)
		if err != nil {
			return fmt.Errorf("failed to get existing match ID: %w", err)
		}
		matchDBID = existingID
	}

	// Trouver le participant correspondant à l'utilisateur
	var userParticipant *RiotParticipant
	for _, participant := range riotMatch.Info.Participants {
		if participant.PUUID == user.RiotPUUID {
			userParticipant = &participant
			break
		}
	}

	if userParticipant == nil {
		return fmt.Errorf("user not found in match participants")
	}

	// Sérialiser les stats détaillées
	detailedStatsJSON, _ := json.Marshal(userParticipant)

	// Insérer le participant
	participantQuery := `
		INSERT OR IGNORE INTO match_participants (
			match_id, user_id, participant_id, team_id, champion_id, champion_name,
			champion_level, kills, deaths, assists, total_damage_dealt,
			total_damage_dealt_to_champions, total_damage_taken, gold_earned,
			total_minions_killed, vision_score, item0, item1, item2, item3,
			item4, item5, item6, win, detailed_stats
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Déterminer l'équipe (Red=200, Blue=100)
	teamID := 100
	if len(riotMatch.Info.Participants) > 5 {
		for i, p := range riotMatch.Info.Participants {
			if p.PUUID == user.RiotPUUID && i >= 5 {
				teamID = 200
				break
			}
		}
	}

	_, err = tx.Exec(participantQuery,
		matchDBID, user.ID, 1, teamID, userParticipant.ChampionID,
		userParticipant.ChampionName, 18, userParticipant.Kills,
		userParticipant.Deaths, userParticipant.Assists, userParticipant.TotalDamageDealtToChampions,
		userParticipant.TotalDamageDealtToChampions, 0, userParticipant.GoldEarned,
		userParticipant.TotalMinionsKilled+userParticipant.NeutralMinionsKilled,
		userParticipant.VisionScore, userParticipant.Item0, userParticipant.Item1,
		userParticipant.Item2, userParticipant.Item3, userParticipant.Item4,
		userParticipant.Item5, 0, userParticipant.Win, string(detailedStatsJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to insert participant: %w", err)
	}

	return tx.Commit()
}

// GetUserMatches récupère les matches d'un utilisateur avec pagination
func (ms *MatchService) GetUserMatches(userID int, limit, offset int) ([]models.MatchSummary, error) {
	query := `
		SELECT m.id, m.match_id, m.platform, m.game_creation, m.game_duration,
			   m.game_mode, m.game_type, m.queue_id, m.created_at,
			   mp.champion_id, mp.champion_name, mp.kills, mp.deaths, mp.assists,
			   mp.total_damage_dealt_to_champions, mp.gold_earned, 
			   mp.total_minions_killed, mp.vision_score, mp.win
		FROM matches m
		JOIN match_participants mp ON m.id = mp.match_id
		WHERE mp.user_id = ?
		ORDER BY m.game_creation DESC
		LIMIT ? OFFSET ?
	`

	rows, err := ms.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query matches: %w", err)
	}
	defer rows.Close()

	var matches []models.MatchSummary
	for rows.Next() {
		var match models.Match
		var participant models.MatchParticipant
		var gameCreation int64

		err := rows.Scan(
			&match.ID, &match.MatchID, &match.Platform, &gameCreation,
			&match.GameDuration, &match.GameMode, &match.GameType, &match.QueueID,
			&match.CreatedAt, &participant.ChampionID, &participant.ChampionName,
			&participant.Kills, &participant.Deaths, &participant.Assists,
			&participant.TotalDamageDealtToChampions, &participant.GoldEarned,
			&participant.TotalMinionsKilled, &participant.VisionScore, &participant.Win,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan match: %w", err)
		}

		match.GameCreation = gameCreation

		matches = append(matches, models.MatchSummary{
			Match:       match,
			Participant: participant,
		})
	}

	return matches, nil
}

// Helper functions

func (ms *MatchService) getUserByID(userID int) (*models.User, error) {
	query := `
		SELECT id, riot_id, riot_tag, riot_puuid, summoner_id, account_id,
			   profile_icon_id, summoner_level, region, is_validated,
			   created_at, updated_at, last_sync
		FROM users WHERE id = ?
	`

	var user models.User
	var summonerID, accountID sql.NullString
	var lastSync sql.NullTime

	err := ms.db.QueryRow(query, userID).Scan(
		&user.ID, &user.RiotID, &user.RiotTag, &user.RiotPUUID,
		&summonerID, &accountID, &user.ProfileIconID, &user.SummonerLevel,
		&user.Region, &user.IsValidated, &user.CreatedAt, &user.UpdatedAt,
		&lastSync,
	)
	if err != nil {
		return nil, err
	}

	if summonerID.Valid {
		user.SummonerID = &summonerID.String
	}
	if accountID.Valid {
		user.AccountID = &accountID.String
	}
	if lastSync.Valid {
		user.LastSync = &lastSync.Time
	}

	return &user, nil
}

func (ms *MatchService) createSyncJob(userID int, jobType string) (*models.SyncJob, error) {
	query := `
		INSERT INTO sync_jobs (user_id, job_type, status, started_at)
		VALUES (?, ?, 'pending', ?)
	`

	result, err := ms.db.Exec(query, userID, jobType, time.Now())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.SyncJob{
		ID:      int(id),
		UserID:  userID,
		JobType: jobType,
		Status:  "pending",
	}, nil
}

func (ms *MatchService) updateSyncJobStatus(jobID int, status, errorMessage string) error {
	query := `UPDATE sync_jobs SET status = ?, error_message = ? WHERE id = ?`
	_, err := ms.db.Exec(query, status, errorMessage, jobID)
	return err
}

func (ms *MatchService) completeSyncJob(jobID, processed, newMatches, updated int) error {
	query := `
		UPDATE sync_jobs 
		SET status = 'completed', completed_at = ?, matches_processed = ?,
			matches_new = ?, matches_updated = ?
		WHERE id = ?
	`
	_, err := ms.db.Exec(query, time.Now(), processed, newMatches, updated, jobID)
	return err
}

func (ms *MatchService) matchExists(matchID string) (bool, error) {
	var count int
	err := ms.db.QueryRow("SELECT COUNT(*) FROM matches WHERE match_id = ?", matchID).Scan(&count)
	return count > 0, err
}

func (ms *MatchService) updateUserLastSync(userID int) error {
	query := `UPDATE users SET last_sync = ? WHERE id = ?`
	_, err := ms.db.Exec(query, time.Now(), userID)
	return err
}

// GetUserStats calcule les statistiques d'un utilisateur
func (ms *MatchService) GetUserStats(userID int) (*models.DashboardStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_matches,
			AVG(CASE WHEN mp.win THEN 1.0 ELSE 0.0 END) as win_rate,
			AVG(CAST(mp.kills + mp.assists AS FLOAT) / CASE WHEN mp.deaths = 0 THEN 1 ELSE mp.deaths END) as avg_kda,
			mp.champion_name,
			u.last_sync
		FROM match_participants mp
		JOIN matches m ON mp.match_id = m.id
		JOIN users u ON mp.user_id = u.id
		WHERE mp.user_id = ?
		GROUP BY mp.champion_name, u.last_sync
		ORDER BY COUNT(*) DESC
		LIMIT 1
	`

	var stats models.DashboardStats
	var lastSync sql.NullTime

	err := ms.db.QueryRow(query, userID).Scan(
		&stats.TotalMatches, &stats.WinRate, &stats.AverageKDA,
		&stats.FavoriteChampion, &lastSync,
	)
	
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	if lastSync.Valid {
		stats.LastSyncAt = &lastSync.Time
	}

	return &stats, nil
}