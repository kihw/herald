-- LoL Analytics Database Schema
-- PostgreSQL Migration Script

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    riot_id VARCHAR(255) NOT NULL,
    riot_tag VARCHAR(50) NOT NULL,
    riot_puuid VARCHAR(100) UNIQUE NOT NULL,
    region VARCHAR(10) NOT NULL,
    summoner_id VARCHAR(100),
    account_id VARCHAR(100),
    summoner_name VARCHAR(255),
    summoner_level INTEGER,
    profile_icon_id INTEGER,
    revision_date BIGINT,
    last_sync TIMESTAMP,
    is_validated BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(riot_id, riot_tag)
);

-- Matches table
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    match_id VARCHAR(100) NOT NULL,
    game_creation BIGINT,
    game_duration INTEGER,
    queue_id INTEGER,
    season_id INTEGER,
    participant_data JSONB,
    team_data JSONB,
    timeline_data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, match_id)
);

-- Champion stats table
CREATE TABLE IF NOT EXISTS champion_stats (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    champion_id INTEGER NOT NULL,
    champion_name VARCHAR(100),
    games_played INTEGER DEFAULT 0,
    wins INTEGER DEFAULT 0,
    losses INTEGER DEFAULT 0,
    kills INTEGER DEFAULT 0,
    deaths INTEGER DEFAULT 0,
    assists INTEGER DEFAULT 0,
    cs_total INTEGER DEFAULT 0,
    gold_earned INTEGER DEFAULT 0,
    damage_dealt INTEGER DEFAULT 0,
    damage_taken INTEGER DEFAULT 0,
    vision_score INTEGER DEFAULT 0,
    performance_score FLOAT,
    last_played TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, champion_id)
);

-- MMR history table
CREATE TABLE IF NOT EXISTS mmr_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    estimated_mmr INTEGER,
    mmr_change INTEGER,
    rank_estimate VARCHAR(50),
    confidence FLOAT,
    match_id VARCHAR(100),
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Analytics cache table
CREATE TABLE IF NOT EXISTS analytics_cache (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    cache_key VARCHAR(255) NOT NULL,
    cache_value JSONB,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, cache_key)
);

-- Recommendations table
CREATE TABLE IF NOT EXISTS recommendations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50),
    title VARCHAR(255),
    description TEXT,
    priority INTEGER,
    confidence FLOAT,
    expected_improvement VARCHAR(100),
    action_items JSONB,
    champion_id INTEGER,
    role VARCHAR(20),
    time_period VARCHAR(20),
    expires_at TIMESTAMP,
    is_dismissed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Performance metrics table
CREATE TABLE IF NOT EXISTS performance_metrics (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    period VARCHAR(20),
    total_games INTEGER,
    win_rate FLOAT,
    avg_kda FLOAT,
    avg_cs_per_min FLOAT,
    avg_gold_per_min FLOAT,
    avg_damage_per_min FLOAT,
    avg_vision_score FLOAT,
    best_role VARCHAR(20),
    worst_role VARCHAR(20),
    performance_score FLOAT,
    trend_direction VARCHAR(20),
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Jobs table for async processing
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    job_type VARCHAR(50),
    status VARCHAR(20) DEFAULT 'pending',
    progress INTEGER DEFAULT 0,
    result JSONB,
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_riot_puuid ON users(riot_puuid);
CREATE INDEX IF NOT EXISTS idx_users_region ON users(region);
CREATE INDEX IF NOT EXISTS idx_matches_user_id ON matches(user_id);
CREATE INDEX IF NOT EXISTS idx_matches_match_id ON matches(match_id);
CREATE INDEX IF NOT EXISTS idx_matches_game_creation ON matches(game_creation DESC);
CREATE INDEX IF NOT EXISTS idx_champion_stats_user_champion ON champion_stats(user_id, champion_id);
CREATE INDEX IF NOT EXISTS idx_mmr_history_user_date ON mmr_history(user_id, recorded_at DESC);
CREATE INDEX IF NOT EXISTS idx_analytics_cache_user_key ON analytics_cache(user_id, cache_key);
CREATE INDEX IF NOT EXISTS idx_analytics_cache_expires ON analytics_cache(expires_at);
CREATE INDEX IF NOT EXISTS idx_recommendations_user ON recommendations(user_id, is_dismissed, expires_at);
CREATE INDEX IF NOT EXISTS idx_performance_metrics_user_period ON performance_metrics(user_id, period, recorded_at DESC);
CREATE INDEX IF NOT EXISTS idx_jobs_user_status ON jobs(user_id, status);

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
CREATE TRIGGER update_champion_stats_updated_at BEFORE UPDATE ON champion_stats
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insights table for real-time notifications
CREATE TABLE IF NOT EXISTS insights (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    level VARCHAR(20) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB,
    action_url VARCHAR(255),
    is_read BOOLEAN DEFAULT false,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Additional indexes for insights
CREATE INDEX IF NOT EXISTS idx_insights_user_read ON insights(user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_insights_expires ON insights(expires_at);
CREATE INDEX IF NOT EXISTS idx_insights_created ON insights(created_at DESC);

-- Function to clean expired cache and insights
CREATE OR REPLACE FUNCTION clean_expired_cache()
RETURNS void AS $$
BEGIN
    DELETE FROM analytics_cache WHERE expires_at < CURRENT_TIMESTAMP;
    DELETE FROM recommendations WHERE expires_at < CURRENT_TIMESTAMP AND is_dismissed = false;
    DELETE FROM insights WHERE expires_at IS NOT NULL AND expires_at < CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;