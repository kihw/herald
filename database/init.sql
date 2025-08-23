-- Herald.lol Database Initialization Script
-- This script sets up the PostgreSQL database for development

-- Create database (this will be handled by Docker environment variables)
-- CREATE DATABASE herald_dev;

-- Connect to the database
-- \c herald_dev;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create indexes for common queries (will be created by GORM auto-migration)
-- These are here for reference and can be manually created if needed

-- Users table indexes
-- CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
-- CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
-- CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Riot accounts table indexes  
-- CREATE INDEX IF NOT EXISTS idx_riot_accounts_puuid ON riot_accounts(puuid);
-- CREATE INDEX IF NOT EXISTS idx_riot_accounts_user_id ON riot_accounts(user_id);
-- CREATE INDEX IF NOT EXISTS idx_riot_accounts_region ON riot_accounts(region);

-- Matches table indexes
-- CREATE INDEX IF NOT EXISTS idx_matches_match_id ON matches(match_id);
-- CREATE INDEX IF NOT EXISTS idx_matches_platform_id ON matches(platform_id);
-- CREATE INDEX IF NOT EXISTS idx_matches_queue_id ON matches(queue_id);
-- CREATE INDEX IF NOT EXISTS idx_matches_game_start_timestamp ON matches(game_start_timestamp);

-- Match participants table indexes
-- CREATE INDEX IF NOT EXISTS idx_match_participants_match_id ON match_participants(match_id);
-- CREATE INDEX IF NOT EXISTS idx_match_participants_puuid ON match_participants(puuid);
-- CREATE INDEX IF NOT EXISTS idx_match_participants_champion_id ON match_participants(champion_id);

-- TFT matches and participants indexes
-- CREATE INDEX IF NOT EXISTS idx_tft_matches_match_id ON tft_matches(match_id);
-- CREATE INDEX IF NOT EXISTS idx_tft_participants_tft_match_id ON tft_participants(tft_match_id);
-- CREATE INDEX IF NOT EXISTS idx_tft_participants_puuid ON tft_participants(puuid);

-- Insert some initial data for development
-- This will be handled by the application

-- Development user for testing (password: "password123")
-- INSERT INTO users (
--     id, email, username, password_hash, display_name, 
--     is_active, is_premium, created_at, updated_at
-- ) VALUES (
--     uuid_generate_v4(),
--     'dev@herald.lol',
--     'developer',
--     crypt('password123', gen_salt('bf')),
--     'Developer',
--     true,
--     false,
--     NOW(),
--     NOW()
-- ) ON CONFLICT (email) DO NOTHING;

-- Log successful initialization
DO $$
BEGIN
    RAISE NOTICE 'Herald.lol database initialized successfully!';
    RAISE NOTICE 'Extensions created: uuid-ossp, pgcrypto';
    RAISE NOTICE 'Ready for GORM auto-migration';
END $$;