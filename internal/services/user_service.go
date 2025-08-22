package services

import (
	"database/sql"
	"fmt"
	"time"

	"lol-match-exporter/internal/models"
)

// UserService gère les opérations liées aux utilisateurs
type UserService struct {
	db *sql.DB
}

// NewUserService crée une nouvelle instance du service utilisateur
func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// GetUserByRiotID récupère un utilisateur par son Riot ID et tag
func (us *UserService) GetUserByRiotID(riotID, riotTag string) (*models.User, error) {
	query := `
		SELECT id, riot_id, riot_tag, riot_puuid, summoner_id, account_id, 
			   profile_icon_id, summoner_level, region, is_validated, 
			   created_at, updated_at, last_sync
		FROM users 
		WHERE riot_id = ? AND riot_tag = ?
	`
	
	var user models.User
	var summonerID, accountID sql.NullString
	var lastSync sql.NullTime
	
	err := us.db.QueryRow(query, riotID, riotTag).Scan(
		&user.ID, &user.RiotID, &user.RiotTag, &user.RiotPUUID,
		&summonerID, &accountID, &user.ProfileIconID, &user.SummonerLevel,
		&user.Region, &user.IsValidated, &user.CreatedAt, &user.UpdatedAt,
		&lastSync,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Utilisateur non trouvé
		}
		return nil, fmt.Errorf("erreur lors de la récupération de l'utilisateur: %w", err)
	}
	
	// Gérer les valeurs nullables
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

// CreateUser crée un nouvel utilisateur
func (us *UserService) CreateUser(riotID, riotTag, riotPUUID, region string, profileIconID, summonerLevel int) (*models.User, error) {
	query := `
		INSERT INTO users (riot_id, riot_tag, riot_puuid, region, profile_icon_id, 
						  summoner_level, is_validated, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, true, ?, ?)
	`
	
	now := time.Now()
	result, err := us.db.Exec(query, riotID, riotTag, riotPUUID, region, 
							  profileIconID, summonerLevel, now, now)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de l'utilisateur: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de l'ID: %w", err)
	}
	
	user := &models.User{
		ID:              int(id),
		RiotID:          riotID,
		RiotTag:         riotTag,
		RiotPUUID:       riotPUUID,
		ProfileIconID:   profileIconID,
		SummonerLevel:   summonerLevel,
		Region:          region,
		IsValidated:     true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	
	return user, nil
}

// UpdateUserSummonerInfo met à jour les informations d'invocateur d'un utilisateur
func (us *UserService) UpdateUserSummonerInfo(userID int, summonerID, accountID string, profileIconID, summonerLevel int) error {
	query := `
		UPDATE users 
		SET summoner_id = ?, account_id = ?, profile_icon_id = ?, 
			summoner_level = ?, updated_at = ?
		WHERE id = ?
	`
	
	_, err := us.db.Exec(query, summonerID, accountID, profileIconID, 
						 summonerLevel, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour de l'utilisateur: %w", err)
	}
	
	return nil
}

// UpdateLastSync met à jour le timestamp de dernière synchronisation
func (us *UserService) UpdateLastSync(userID int) error {
	query := `UPDATE users SET last_sync = ? WHERE id = ?`
	
	_, err := us.db.Exec(query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour de last_sync: %w", err)
	}
	
	return nil
}

// GetUserSettings récupère les paramètres d'un utilisateur
func (us *UserService) GetUserSettings(userID int) (*models.UserSettings, error) {
	query := `
		SELECT id, user_id, platform, queue_types, language, include_timeline,
			   include_all_data, light_mode, auto_sync_enabled, sync_frequency_hours,
			   created_at, updated_at
		FROM user_settings 
		WHERE user_id = ?
	`
	
	var settings models.UserSettings
	err := us.db.QueryRow(query, userID).Scan(
		&settings.ID, &settings.UserID, &settings.Platform, &settings.QueueTypes,
		&settings.Language, &settings.IncludeTimeline, &settings.IncludeAllData,
		&settings.LightMode, &settings.AutoSyncEnabled, &settings.SyncFrequencyHours,
		&settings.CreatedAt, &settings.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			// Créer des paramètres par défaut
			return us.CreateDefaultSettings(userID)
		}
		return nil, fmt.Errorf("erreur lors de la récupération des paramètres: %w", err)
	}
	
	return &settings, nil
}

// CreateDefaultSettings crée des paramètres par défaut pour un utilisateur
func (us *UserService) CreateDefaultSettings(userID int) (*models.UserSettings, error) {
	query := `
		INSERT INTO user_settings (user_id, platform, language, include_timeline,
								  include_all_data, light_mode, auto_sync_enabled,
								  sync_frequency_hours, created_at, updated_at)
		VALUES (?, 'euw1', 'fr', true, true, false, true, 24, ?, ?)
	`
	
	now := time.Now()
	result, err := us.db.Exec(query, userID, now, now)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création des paramètres: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de l'ID: %w", err)
	}
	
	settings := &models.UserSettings{
		ID:                 int(id),
		UserID:             userID,
		Platform:           "euw1",
		Language:           "fr",
		IncludeTimeline:    true,
		IncludeAllData:     true,
		LightMode:          false,
		AutoSyncEnabled:    true,
		SyncFrequencyHours: 24,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	
	return settings, nil
}

// UpdateUserSettings met à jour les paramètres d'un utilisateur
func (us *UserService) UpdateUserSettings(settings *models.UserSettings) error {
	query := `
		UPDATE user_settings 
		SET platform = ?, queue_types = ?, language = ?, include_timeline = ?,
			include_all_data = ?, light_mode = ?, auto_sync_enabled = ?,
			sync_frequency_hours = ?, updated_at = ?
		WHERE user_id = ?
	`
	
	_, err := us.db.Exec(query, settings.Platform, settings.QueueTypes,
						 settings.Language, settings.IncludeTimeline,
						 settings.IncludeAllData, settings.LightMode,
						 settings.AutoSyncEnabled, settings.SyncFrequencyHours,
						 time.Now(), settings.UserID)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour des paramètres: %w", err)
	}
	
	return nil
}