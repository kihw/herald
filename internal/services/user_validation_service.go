package services

import (
	"database/sql"
	"fmt"
	"time"

	"lol-match-exporter/internal/models"
)

type UserValidationService struct {
	DB                   *sql.DB
	RiotValidationService *RiotValidationService
}

func NewUserValidationService(db *sql.DB, riotValidationService *RiotValidationService) *UserValidationService {
	return &UserValidationService{
		DB:                   db,
		RiotValidationService: riotValidationService,
	}
}

// ValidateAndCreateUser validates a Riot account and creates or updates the user
func (s *UserValidationService) ValidateAndCreateUser(riotID, riotTag, region string) (*models.User, error) {
	// Step 1: Validate with Riot API
	validatedUser, err := s.RiotValidationService.ValidateRiotAccount(riotID, riotTag, region)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Step 2: Check if user already exists
	existingUser, err := s.GetUserByPUUID(validatedUser.RiotPUUID)
	if err != nil && err.Error() != "user not found" {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		// Update existing user
		return s.updateExistingUser(existingUser, validatedUser)
	} else {
		// Create new user
		return s.createNewUser(validatedUser)
	}
}

// GetUserByPUUID gets a user by their Riot PUUID
func (s *UserValidationService) GetUserByPUUID(puuid string) (*models.User, error) {
	var user models.User
	
	query := `
		SELECT id, riot_id, riot_tag, riot_puuid, summoner_id, account_id, 
		       profile_icon_id, summoner_level, region, is_validated,
		       created_at, updated_at, last_sync
		FROM users 
		WHERE riot_puuid = $1
	`
	
	err := s.DB.QueryRow(query, puuid).Scan(
		&user.ID, &user.RiotID, &user.RiotTag, &user.RiotPUUID,
		&user.SummonerID, &user.AccountID, &user.ProfileIconID, 
		&user.SummonerLevel, &user.Region, &user.IsValidated,
		&user.CreatedAt, &user.UpdatedAt, &user.LastSync,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return &user, nil
}

// GetUserByID gets a user by their database ID
func (s *UserValidationService) GetUserByID(id int) (*models.User, error) {
	var user models.User
	
	query := `
		SELECT id, riot_id, riot_tag, riot_puuid, summoner_id, account_id, 
		       profile_icon_id, summoner_level, region, is_validated,
		       created_at, updated_at, last_sync
		FROM users 
		WHERE id = $1
	`
	
	err := s.DB.QueryRow(query, id).Scan(
		&user.ID, &user.RiotID, &user.RiotTag, &user.RiotPUUID,
		&user.SummonerID, &user.AccountID, &user.ProfileIconID, 
		&user.SummonerLevel, &user.Region, &user.IsValidated,
		&user.CreatedAt, &user.UpdatedAt, &user.LastSync,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return &user, nil
}

// createNewUser creates a new user in the database
func (s *UserValidationService) createNewUser(user *models.User) (*models.User, error) {
	insertQuery := `
		INSERT INTO users (riot_id, riot_tag, riot_puuid, summoner_id, account_id,
		                  profile_icon_id, summoner_level, region, is_validated, 
		                  created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	
	now := time.Now()
	err := s.DB.QueryRow(
		insertQuery,
		user.RiotID, user.RiotTag, user.RiotPUUID, user.SummonerID, user.AccountID,
		user.ProfileIconID, user.SummonerLevel, user.Region, user.IsValidated,
		now, now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	// Create default user settings
	err = s.createDefaultUserSettings(user.ID, user.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create user settings: %w", err)
	}
	
	return user, nil
}

// updateExistingUser updates an existing user with new validation data
func (s *UserValidationService) updateExistingUser(existingUser, validatedUser *models.User) (*models.User, error) {
	updateQuery := `
		UPDATE users 
		SET riot_id = $1, riot_tag = $2, summoner_id = $3, account_id = $4,
		    profile_icon_id = $5, summoner_level = $6, region = $7, 
		    is_validated = $8, updated_at = $9
		WHERE riot_puuid = $10
	`
	
	now := time.Now()
	_, err := s.DB.Exec(
		updateQuery,
		validatedUser.RiotID, validatedUser.RiotTag, validatedUser.SummonerID, validatedUser.AccountID,
		validatedUser.ProfileIconID, validatedUser.SummonerLevel, validatedUser.Region,
		validatedUser.IsValidated, now, validatedUser.RiotPUUID,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	
	// Update the existing user fields
	existingUser.RiotID = validatedUser.RiotID
	existingUser.RiotTag = validatedUser.RiotTag
	existingUser.SummonerID = validatedUser.SummonerID
	existingUser.AccountID = validatedUser.AccountID
	existingUser.ProfileIconID = validatedUser.ProfileIconID
	existingUser.SummonerLevel = validatedUser.SummonerLevel
	existingUser.Region = validatedUser.Region
	existingUser.IsValidated = validatedUser.IsValidated
	existingUser.UpdatedAt = now
	
	return existingUser, nil
}

// createDefaultUserSettings creates default settings for a new user
func (s *UserValidationService) createDefaultUserSettings(userID int, region string) error {
	insertQuery := `
		INSERT INTO user_settings (user_id, platform, queue_types, language, include_timeline,
		                          include_all_data, light_mode, auto_sync_enabled, 
		                          sync_frequency_hours, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	
	now := time.Now()
	queueTypes := "{420,440}" // Ranked Solo/Duo, Ranked Flex as PostgreSQL array
	
	_, err := s.DB.Exec(
		insertQuery,
		userID, region, queueTypes, "en_US", true,
		true, true, true, 24, now, now,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create user settings: %w", err)
	}
	
	return nil
}

// UpdateLastSync updates the user's last sync timestamp
func (s *UserValidationService) UpdateLastSync(userID int) error {
	query := `UPDATE users SET last_sync = $1 WHERE id = $2`
	now := time.Now()
	_, err := s.DB.Exec(query, now, userID)
	if err != nil {
		return fmt.Errorf("failed to update last sync: %w", err)
	}
	return nil
}
