package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDatabase creates a new database connection
func NewDatabase(config Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connection established")
	return &Database{db}, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.DB.Close()
}

// Migrate runs database migrations
func (db *Database) Migrate() error {
	log.Println("🔄 Running database migrations...")

	// Read and execute migration file - Updated for validation system
	migrationSQL := `
-- Users table for validation-based authentication
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    riot_id VARCHAR(50) NOT NULL,
    riot_tag VARCHAR(20) NOT NULL,
    riot_puuid VARCHAR(78) UNIQUE NOT NULL,
    summoner_id VARCHAR(63),
    account_id VARCHAR(56),
    profile_icon_id INTEGER DEFAULT 0,
    summoner_level INTEGER DEFAULT 1,
    region VARCHAR(10) NOT NULL DEFAULT 'EUW1',
    is_validated BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_sync TIMESTAMP WITH TIME ZONE,
    
    UNIQUE(riot_id, riot_tag)
);

CREATE INDEX IF NOT EXISTS idx_users_puuid ON users(riot_puuid);
CREATE INDEX IF NOT EXISTS idx_users_riot_id_tag ON users(riot_id, riot_tag);
CREATE INDEX IF NOT EXISTS idx_users_validation ON users(is_validated);

-- User settings for export configuration
CREATE TABLE IF NOT EXISTS user_settings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    platform VARCHAR(10) DEFAULT 'euw1',
    queue_types INTEGER[] DEFAULT '{420,440}',
    language VARCHAR(10) DEFAULT 'en_US',
    include_timeline BOOLEAN DEFAULT true,
    include_all_data BOOLEAN DEFAULT true,
    light_mode BOOLEAN DEFAULT true,
    auto_sync_enabled BOOLEAN DEFAULT true,
    sync_frequency_hours INTEGER DEFAULT 24,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(user_id)
);

-- Matches table
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    match_id VARCHAR(20) UNIQUE NOT NULL,
    platform VARCHAR(10) NOT NULL,
    game_creation BIGINT NOT NULL,
    game_duration INTEGER NOT NULL,
    game_end_timestamp BIGINT,
    game_mode VARCHAR(20),
    game_type VARCHAR(20),
    game_version VARCHAR(20),
    map_id INTEGER,
    queue_id INTEGER,
    season_id INTEGER,
    tournament_code VARCHAR(255),
    data_version VARCHAR(10),
    raw_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_matches_match_id ON matches(match_id);
CREATE INDEX IF NOT EXISTS idx_matches_platform ON matches(platform);
CREATE INDEX IF NOT EXISTS idx_matches_queue_id ON matches(queue_id);
CREATE INDEX IF NOT EXISTS idx_matches_game_creation ON matches(game_creation);

-- Match participants
CREATE TABLE IF NOT EXISTS match_participants (
    id SERIAL PRIMARY KEY,
    match_id INTEGER REFERENCES matches(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    participant_id INTEGER NOT NULL,
    team_id INTEGER NOT NULL,
    champion_id INTEGER NOT NULL,
    champion_name VARCHAR(50),
    champion_level INTEGER,
    
    kills INTEGER DEFAULT 0,
    deaths INTEGER DEFAULT 0,
    assists INTEGER DEFAULT 0,
    total_damage_dealt INTEGER DEFAULT 0,
    total_damage_dealt_to_champions INTEGER DEFAULT 0,
    total_damage_taken INTEGER DEFAULT 0,
    gold_earned INTEGER DEFAULT 0,
    total_minions_killed INTEGER DEFAULT 0,
    vision_score INTEGER DEFAULT 0,
    
    item0 INTEGER DEFAULT 0,
    item1 INTEGER DEFAULT 0,
    item2 INTEGER DEFAULT 0,
    item3 INTEGER DEFAULT 0,
    item4 INTEGER DEFAULT 0,
    item5 INTEGER DEFAULT 0,
    item6 INTEGER DEFAULT 0,
    
    win BOOLEAN NOT NULL,
    detailed_stats JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(match_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_participants_user_id ON match_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_participants_match_id ON match_participants(match_id);
CREATE INDEX IF NOT EXISTS idx_participants_champion_id ON match_participants(champion_id);

-- Sync jobs
CREATE TABLE IF NOT EXISTS sync_jobs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    job_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    matches_processed INTEGER DEFAULT 0,
    matches_new INTEGER DEFAULT 0,
    matches_updated INTEGER DEFAULT 0,
    error_message TEXT,
    last_match_timestamp BIGINT
);

CREATE INDEX IF NOT EXISTS idx_sync_jobs_user_id ON sync_jobs(user_id);
CREATE INDEX IF NOT EXISTS idx_sync_jobs_status ON sync_jobs(status);
CREATE INDEX IF NOT EXISTS idx_sync_jobs_started_at ON sync_jobs(started_at);

-- System configuration
CREATE TABLE IF NOT EXISTS system_config (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert default configuration
INSERT INTO system_config (key, value, description) VALUES
('sync_batch_size', '100', 'Number of matches to process per batch'),
('rate_limit_requests_per_second', '20', 'API rate limit'),
('rate_limit_requests_per_2_minutes', '100', 'API rate limit per 2 minutes'),
('manual_sync_cooldown_minutes', '15', 'Cooldown between manual syncs'),
('daily_sync_time', '00:00', 'Time for daily automatic sync (HH:MM)'),
('max_matches_per_sync', '1000', 'Maximum matches to sync per job')
ON CONFLICT (key) DO NOTHING;

-- Functions and triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_user_settings_updated_at ON user_settings;
CREATE TRIGGER update_user_settings_updated_at BEFORE UPDATE ON user_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_system_config_updated_at ON system_config;
CREATE TRIGGER update_system_config_updated_at BEFORE UPDATE ON system_config
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
`

	_, err := db.Exec(migrationSQL)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// Health checks the database connection
func (db *Database) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}
