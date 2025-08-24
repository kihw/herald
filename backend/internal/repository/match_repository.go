package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/herald-lol/herald/backend/internal/models"
)

// MatchRepository handles match data persistence
type MatchRepository struct {
	db *sql.DB
}

// NewMatchRepository creates a new match repository
func NewMatchRepository(db *sql.DB) *MatchRepository {
	return &MatchRepository{
		db: db,
	}
}

// GetMatchByID retrieves a match by ID
func (r *MatchRepository) GetMatchByID(matchID string) (*models.Match, error) {
	// Placeholder implementation
	match := &models.Match{
		ID:           uuid.New(),
		MatchID:      matchID,
		GameMode:     "CLASSIC",
		GameDuration: 25 * 60, // 25 minutes in seconds
		CreatedAt:    time.Now(),
	}
	return match, nil
}

// SaveMatch saves a match to the database
func (r *MatchRepository) SaveMatch(match *models.Match) error {
	// Placeholder implementation
	return nil
}

// GetPlayerMatches retrieves matches for a player
func (r *MatchRepository) GetPlayerMatches(playerID string, limit int) ([]*models.Match, error) {
	// Placeholder implementation
	matches := make([]*models.Match, 0, limit)
	for i := 0; i < limit && i < 10; i++ {
		matches = append(matches, &models.Match{
			ID:           uuid.New(),
			MatchID:      "match-" + playerID + "-" + string(rune(48+i)),
			GameMode:     "CLASSIC",
			GameDuration: (20 + i*2) * 60, // minutes converted to seconds
			CreatedAt:    time.Now().Add(-time.Duration(i*24) * time.Hour),
		})
	}
	return matches, nil
}
