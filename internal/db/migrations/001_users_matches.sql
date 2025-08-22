-- Migration 001: Tables principales pour les utilisateurs et matchs
-- Conversion de la structure PostgreSQL vers SQLite

-- Table users - Structure principale pour l'authentification Riot
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    riot_id TEXT NOT NULL,
    riot_tag TEXT NOT NULL,
    riot_puuid TEXT UNIQUE NOT NULL,
    region TEXT NOT NULL,
    summoner_id TEXT,
    account_id TEXT,
    summoner_name TEXT,
    summoner_level INTEGER,
    profile_icon_id INTEGER DEFAULT 0,
    revision_date INTEGER,
    last_sync TIMESTAMP,
    is_validated BOOLEAN DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(riot_id, riot_tag)
);

-- Index pour les utilisateurs
CREATE INDEX IF NOT EXISTS idx_users_riot_puuid ON users(riot_puuid);
CREATE INDEX IF NOT EXISTS idx_users_region ON users(region);
CREATE INDEX IF NOT EXISTS idx_users_riot_id_tag ON users(riot_id, riot_tag);

-- Table user_settings - Paramètres utilisateur
CREATE TABLE IF NOT EXISTS user_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    platform TEXT DEFAULT 'euw1',
    queue_types TEXT DEFAULT '[]', -- JSON array
    language TEXT DEFAULT 'fr',
    include_timeline BOOLEAN DEFAULT 1,
    include_all_data BOOLEAN DEFAULT 1,
    light_mode BOOLEAN DEFAULT 0,
    auto_sync_enabled BOOLEAN DEFAULT 1,
    sync_frequency_hours INTEGER DEFAULT 24,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id)
);

-- Table matches - Stockage permanent des parties
CREATE TABLE IF NOT EXISTS matches (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    match_id TEXT NOT NULL,
    platform TEXT NOT NULL,
    game_creation INTEGER,
    game_duration INTEGER,
    game_end_timestamp INTEGER,
    game_mode TEXT,
    game_type TEXT,
    game_version TEXT,
    map_id INTEGER,
    queue_id INTEGER,
    season_id INTEGER,
    tournament_code TEXT,
    data_version TEXT,
    raw_data TEXT, -- JSON string for complete match data
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(match_id)
);

-- Index pour les matches
CREATE INDEX IF NOT EXISTS idx_matches_match_id ON matches(match_id);
CREATE INDEX IF NOT EXISTS idx_matches_game_creation ON matches(game_creation DESC);
CREATE INDEX IF NOT EXISTS idx_matches_queue_id ON matches(queue_id);
CREATE INDEX IF NOT EXISTS idx_matches_platform ON matches(platform);

-- Table match_participants - Participation des utilisateurs aux matches
CREATE TABLE IF NOT EXISTS match_participants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    match_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    participant_id INTEGER NOT NULL,
    team_id INTEGER NOT NULL,
    champion_id INTEGER NOT NULL,
    champion_name TEXT,
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
    win BOOLEAN DEFAULT 0,
    detailed_stats TEXT, -- JSON string for additional stats
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(match_id, user_id)
);

-- Index pour les participants
CREATE INDEX IF NOT EXISTS idx_match_participants_match ON match_participants(match_id);
CREATE INDEX IF NOT EXISTS idx_match_participants_user ON match_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_match_participants_champion ON match_participants(champion_id);
CREATE INDEX IF NOT EXISTS idx_match_participants_win ON match_participants(win);

-- Table sync_jobs - Gestion des synchronisations
CREATE TABLE IF NOT EXISTS sync_jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    job_type TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    matches_processed INTEGER DEFAULT 0,
    matches_new INTEGER DEFAULT 0,
    matches_updated INTEGER DEFAULT 0,
    error_message TEXT,
    last_match_timestamp INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Index pour les jobs de sync
CREATE INDEX IF NOT EXISTS idx_sync_jobs_user ON sync_jobs(user_id);
CREATE INDEX IF NOT EXISTS idx_sync_jobs_status ON sync_jobs(status);
CREATE INDEX IF NOT EXISTS idx_sync_jobs_started ON sync_jobs(started_at DESC);

-- Table system_config - Configuration système
CREATE TABLE IF NOT EXISTS system_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT UNIQUE NOT NULL,
    value TEXT,
    description TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Triggers pour mettre à jour updated_at
CREATE TRIGGER IF NOT EXISTS update_users_updated_at
AFTER UPDATE ON users
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_user_settings_updated_at
AFTER UPDATE ON user_settings
BEGIN
    UPDATE user_settings SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_system_config_updated_at
AFTER UPDATE ON system_config
BEGIN
    UPDATE system_config SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;