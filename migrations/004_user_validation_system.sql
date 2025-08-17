-- Migration: 004_user_validation_system.sql
-- Description: Update user table for validation-based authentication system
-- Date: 2025-08-16

-- First, let's see what exists and drop OAuth-specific columns
ALTER TABLE users DROP COLUMN IF EXISTS access_token;
ALTER TABLE users DROP COLUMN IF EXISTS refresh_token;
ALTER TABLE users DROP COLUMN IF EXISTS token_expires_at;
ALTER TABLE users DROP COLUMN IF EXISTS last_region_used;

-- Add validation-specific columns
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_validated BOOLEAN DEFAULT FALSE;

-- Update existing data to be validated (if any)
UPDATE users SET is_validated = TRUE WHERE riot_puuid IS NOT NULL;

-- Ensure proper indexes
CREATE INDEX IF NOT EXISTS idx_users_validation ON users(is_validated);
CREATE INDEX IF NOT EXISTS idx_users_riot_id_tag ON users(riot_id, riot_tag);

-- Update user settings to ensure they work with new system
-- (The table should already exist from previous migrations)

-- Create or update the users table structure to match our new model
CREATE TABLE IF NOT EXISTS users_new (
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
    last_sync TIMESTAMP WITH TIME ZONE
);

-- Copy data if the old table exists and has data
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users') THEN
        INSERT INTO users_new (riot_id, riot_tag, riot_puuid, summoner_id, account_id, 
                              profile_icon_id, summoner_level, region, is_validated, 
                              created_at, updated_at, last_sync)
        SELECT riot_id, riot_tag, riot_puuid, summoner_id, account_id, 
               COALESCE(profile_icon_id, 0), COALESCE(summoner_level, 1), 
               COALESCE(region, 'EUW1'), COALESCE(is_validated, FALSE),
               COALESCE(created_at, NOW()), COALESCE(updated_at, NOW()), last_sync
        FROM users
        ON CONFLICT (riot_puuid) DO NOTHING;
    END IF;
END $$;

-- Replace the old table with the new one
DROP TABLE IF EXISTS users CASCADE;
ALTER TABLE users_new RENAME TO users;

-- Recreate indexes
CREATE UNIQUE INDEX idx_users_riot_puuid ON users(riot_puuid);
CREATE INDEX idx_users_riot_id_tag ON users(riot_id, riot_tag);
CREATE INDEX idx_users_validation ON users(is_validated);
CREATE INDEX idx_users_region ON users(region);

-- Ensure user_settings table exists and is properly linked
CREATE TABLE IF NOT EXISTS user_settings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    platform VARCHAR(10) NOT NULL DEFAULT 'EUW1',
    queue_types INTEGER[] DEFAULT '{420,440}',
    language VARCHAR(10) DEFAULT 'en_US',
    include_timeline BOOLEAN DEFAULT TRUE,
    include_all_data BOOLEAN DEFAULT TRUE,
    light_mode BOOLEAN DEFAULT TRUE,
    auto_sync_enabled BOOLEAN DEFAULT TRUE,
    sync_frequency_hours INTEGER DEFAULT 24,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Ensure proper indexes on user_settings
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_settings_user_id ON user_settings(user_id);

-- Create or update matches table to ensure it works with our system
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    match_id VARCHAR(50) UNIQUE NOT NULL,
    platform VARCHAR(10) NOT NULL,
    game_creation BIGINT NOT NULL,
    game_duration INTEGER NOT NULL,
    game_end_timestamp BIGINT,
    game_mode VARCHAR(50),
    game_type VARCHAR(50),
    game_version VARCHAR(20),
    map_id INTEGER,
    queue_id INTEGER,
    season_id INTEGER,
    tournament_code VARCHAR(100),
    data_version VARCHAR(10),
    raw_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create or update match_participants table
CREATE TABLE IF NOT EXISTS match_participants (
    id SERIAL PRIMARY KEY,
    match_id INTEGER NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
    win BOOLEAN DEFAULT FALSE,
    detailed_stats JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_matches_match_id ON matches(match_id);
CREATE INDEX IF NOT EXISTS idx_matches_platform ON matches(platform);
CREATE INDEX IF NOT EXISTS idx_matches_queue_id ON matches(queue_id);
CREATE INDEX IF NOT EXISTS idx_match_participants_user_id ON match_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_match_participants_match_id ON match_participants(match_id);
CREATE INDEX IF NOT EXISTS idx_match_participants_champion ON match_participants(champion_id);

-- Create sync_jobs table for tracking synchronization
CREATE TABLE IF NOT EXISTS sync_jobs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    job_type VARCHAR(50) NOT NULL DEFAULT 'manual',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
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

-- Create system_config table for system settings
CREATE TABLE IF NOT EXISTS system_config (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert default system config if not exists
INSERT INTO system_config (key, value, description) VALUES 
('sync_cooldown_minutes', '15', 'Minimum minutes between manual syncs')
ON CONFLICT (key) DO NOTHING;

COMMIT;
