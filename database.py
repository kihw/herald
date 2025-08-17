"""
Database module for LoL Analytics Platform.
Handles SQLite database operations, migrations, and data models.
"""

import sqlite3
import json
from datetime import datetime, timedelta
from typing import List, Dict, Any, Optional, Tuple
from pathlib import Path
import logging

class DatabaseManager:
    """Manages SQLite database operations for the LoL Analytics Platform."""
    
    def __init__(self, db_path: str = "data/lol_analytics.db"):
        self.db_path = Path(db_path)
        self.db_path.parent.mkdir(exist_ok=True)
        self.init_database()
    
    def get_connection(self) -> sqlite3.Connection:
        """Get database connection with proper configuration."""
        conn = sqlite3.connect(self.db_path)
        conn.row_factory = sqlite3.Row  # Enable dict-like access
        conn.execute("PRAGMA foreign_keys = ON")  # Enable foreign keys
        return conn
    
    def init_database(self):
        """Initialize database with all required tables."""
        with self.get_connection() as conn:
            # Create tables in dependency order
            self._create_users_table(conn)
            self._create_matches_table(conn)
            self._create_scan_history_table(conn)
            self._create_champion_stats_table(conn)
            self._create_role_performance_table(conn)
            self._create_mmr_history_table(conn)
            self._create_performance_insights_table(conn)
            logging.info("Database initialized successfully")
    
    def _create_users_table(self, conn: sqlite3.Connection):
        """Create users table."""
        conn.execute("""
            CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                puuid TEXT UNIQUE NOT NULL,
                riot_id TEXT NOT NULL,
                platform TEXT NOT NULL,
                summoner_name TEXT,
                summoner_level INTEGER,
                last_scan_date TIMESTAMP,
                total_matches INTEGER DEFAULT 0,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
        """)
    
    def _create_matches_table(self, conn: sqlite3.Connection):
        """Create matches table for storing raw match data."""
        conn.execute("""
            CREATE TABLE IF NOT EXISTS matches (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                match_id TEXT UNIQUE NOT NULL,
                user_id INTEGER NOT NULL,
                game_creation TIMESTAMP NOT NULL,
                season INTEGER NOT NULL,
                queue_id INTEGER NOT NULL,
                game_duration INTEGER,
                game_mode TEXT,
                patch_version TEXT,
                participant_data JSON NOT NULL,
                team_data JSON,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
            )
        """)
        
        # Create indexes for performance
        conn.execute("CREATE INDEX IF NOT EXISTS idx_matches_user_id ON matches(user_id)")
        conn.execute("CREATE INDEX IF NOT EXISTS idx_matches_game_creation ON matches(game_creation)")
        conn.execute("CREATE INDEX IF NOT EXISTS idx_matches_season ON matches(season)")
        conn.execute("CREATE INDEX IF NOT EXISTS idx_matches_queue_id ON matches(queue_id)")
    
    def _create_scan_history_table(self, conn: sqlite3.Connection):
        """Create scan history table."""
        conn.execute("""
            CREATE TABLE IF NOT EXISTS scan_history (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                user_id INTEGER NOT NULL,
                season INTEGER,
                scan_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                matches_found INTEGER DEFAULT 0,
                new_matches INTEGER DEFAULT 0,
                scan_type TEXT DEFAULT 'auto', -- 'auto', 'manual', 'incremental'
                success BOOLEAN DEFAULT TRUE,
                error_message TEXT,
                FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
            )
        """)
    
    def _create_champion_stats_table(self, conn: sqlite3.Connection):
        """Create champion statistics table."""
        conn.execute("""
            CREATE TABLE IF NOT EXISTS champion_stats (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                user_id INTEGER NOT NULL,
                champion_id INTEGER NOT NULL,
                champion_name TEXT NOT NULL,
                role TEXT NOT NULL,
                season INTEGER NOT NULL,
                time_period TEXT NOT NULL, -- 'today', 'week', 'month', 'season', 'all'
                games_played INTEGER DEFAULT 0,
                wins INTEGER DEFAULT 0,
                losses INTEGER DEFAULT 0,
                win_rate REAL DEFAULT 0.0,
                avg_kills REAL DEFAULT 0.0,
                avg_deaths REAL DEFAULT 0.0,
                avg_assists REAL DEFAULT 0.0,
                avg_kda REAL DEFAULT 0.0,
                avg_cs_per_min REAL DEFAULT 0.0,
                avg_gold_per_min REAL DEFAULT 0.0,
                avg_damage_per_min REAL DEFAULT 0.0,
                avg_vision_score REAL DEFAULT 0.0,
                performance_score REAL DEFAULT 0.0,
                trend_direction TEXT DEFAULT 'stable', -- 'improving', 'declining', 'stable'
                last_played TIMESTAMP,
                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
                UNIQUE(user_id, champion_id, role, time_period, season)
            )
        """)
        
        conn.execute("CREATE INDEX IF NOT EXISTS idx_champion_stats_user_role ON champion_stats(user_id, role)")
        conn.execute("CREATE INDEX IF NOT EXISTS idx_champion_stats_performance ON champion_stats(performance_score DESC)")
    
    def _create_role_performance_table(self, conn: sqlite3.Connection):
        """Create role performance table."""
        conn.execute("""
            CREATE TABLE IF NOT EXISTS role_performance (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                user_id INTEGER NOT NULL,
                role TEXT NOT NULL,
                time_period TEXT NOT NULL, -- 'today', 'week', 'month', 'season'
                season INTEGER NOT NULL,
                games INTEGER DEFAULT 0,
                wins INTEGER DEFAULT 0,
                win_rate REAL DEFAULT 0.0,
                avg_performance REAL DEFAULT 0.0,
                avg_kda REAL DEFAULT 0.0,
                avg_cs_per_min REAL DEFAULT 0.0,
                trend_direction TEXT DEFAULT 'stable',
                best_champion TEXT,
                worst_champion TEXT,
                improvement_potential REAL DEFAULT 0.0,
                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
                UNIQUE(user_id, role, time_period, season)
            )
        """)
    
    def _create_mmr_history_table(self, conn: sqlite3.Connection):
        """Create MMR history table."""
        conn.execute("""
            CREATE TABLE IF NOT EXISTS mmr_history (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                user_id INTEGER NOT NULL,
                match_id TEXT NOT NULL,
                estimated_mmr INTEGER NOT NULL,
                mmr_change INTEGER DEFAULT 0,
                confidence_score REAL DEFAULT 0.0,
                game_date TIMESTAMP NOT NULL,
                rank_estimate TEXT,
                lp_estimate INTEGER,
                streak_count INTEGER DEFAULT 0,
                performance_modifier REAL DEFAULT 0.0,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
                FOREIGN KEY(match_id) REFERENCES matches(match_id)
            )
        """)
        
        conn.execute("CREATE INDEX IF NOT EXISTS idx_mmr_history_user_date ON mmr_history(user_id, game_date)")
    
    def _create_performance_insights_table(self, conn: sqlite3.Connection):
        """Create performance insights table."""
        conn.execute("""
            CREATE TABLE IF NOT EXISTS performance_insights (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                user_id INTEGER NOT NULL,
                insight_type TEXT NOT NULL, -- 'suggestion', 'warning', 'achievement', 'prediction'
                category TEXT NOT NULL, -- 'champion', 'role', 'gameplay', 'mmr', 'meta'
                title TEXT NOT NULL,
                description TEXT NOT NULL,
                priority INTEGER DEFAULT 1, -- 1=high, 2=medium, 3=low
                confidence REAL DEFAULT 0.0,
                time_period TEXT NOT NULL,
                expected_improvement TEXT,
                action_items JSON,
                is_active BOOLEAN DEFAULT TRUE,
                expires_at TIMESTAMP,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
            )
        """)
        
        conn.execute("CREATE INDEX IF NOT EXISTS idx_insights_user_active ON performance_insights(user_id, is_active)")
        conn.execute("CREATE INDEX IF NOT EXISTS idx_insights_priority ON performance_insights(priority, created_at)")


class UserManager:
    """Manages user-related database operations."""
    
    def __init__(self, db_manager: DatabaseManager):
        self.db = db_manager
    
    def get_or_create_user(self, puuid: str, riot_id: str, platform: str, 
                          summoner_name: str = None, summoner_level: int = None) -> int:
        """Get existing user or create new one. Returns user_id."""
        with self.db.get_connection() as conn:
            # Try to get existing user
            cursor = conn.execute(
                "SELECT id FROM users WHERE puuid = ?", (puuid,)
            )
            row = cursor.fetchone()
            
            if row:
                # Update user info
                conn.execute("""
                    UPDATE users 
                    SET riot_id = ?, platform = ?, summoner_name = ?, 
                        summoner_level = ?, updated_at = CURRENT_TIMESTAMP
                    WHERE puuid = ?
                """, (riot_id, platform, summoner_name, summoner_level, puuid))
                return row[0]
            else:
                # Create new user
                cursor = conn.execute("""
                    INSERT INTO users (puuid, riot_id, platform, summoner_name, summoner_level)
                    VALUES (?, ?, ?, ?, ?)
                """, (puuid, riot_id, platform, summoner_name, summoner_level))
                return cursor.lastrowid
    
    def get_user_by_puuid(self, puuid: str) -> Optional[Dict[str, Any]]:
        """Get user by PUUID."""
        with self.db.get_connection() as conn:
            cursor = conn.execute(
                "SELECT * FROM users WHERE puuid = ?", (puuid,)
            )
            row = cursor.fetchone()
            return dict(row) if row else None
    
    def update_last_scan_date(self, user_id: int):
        """Update last scan date for user."""
        with self.db.get_connection() as conn:
            conn.execute("""
                UPDATE users 
                SET last_scan_date = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
                WHERE id = ?
            """, (user_id,))
    
    def get_scan_history(self, user_id: int, limit: int = 10) -> List[Dict[str, Any]]:
        """Get recent scan history for user."""
        with self.db.get_connection() as conn:
            cursor = conn.execute("""
                SELECT * FROM scan_history 
                WHERE user_id = ? 
                ORDER BY scan_date DESC 
                LIMIT ?
            """, (user_id, limit))
            return [dict(row) for row in cursor.fetchall()]


class MatchManager:
    """Manages match-related database operations."""
    
    def __init__(self, db_manager: DatabaseManager):
        self.db = db_manager
    
    def save_matches(self, user_id: int, matches_data: List[Dict[str, Any]]) -> int:
        """Save multiple matches for a user. Returns number of new matches."""
        new_matches = 0
        with self.db.get_connection() as conn:
            for match_data in matches_data:
                try:
                    # Extract match info
                    match_id = match_data['metadata']['matchId']
                    info = match_data['info']
                    game_creation = datetime.fromtimestamp(info['gameCreation'] / 1000)
                    
                    # Find participant data for this user
                    participant_data = None
                    for p in info['participants']:
                        if p.get('puuid') == self.get_user_puuid(user_id):
                            participant_data = p
                            break
                    
                    if not participant_data:
                        continue
                    
                    # Determine season from game date
                    season = game_creation.year
                    
                    # Insert match (ignore if already exists)
                    cursor = conn.execute("""
                        INSERT OR IGNORE INTO matches 
                        (match_id, user_id, game_creation, season, queue_id, 
                         game_duration, game_mode, patch_version, participant_data, team_data)
                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                    """, (
                        match_id, user_id, game_creation, season, info['queueId'],
                        info.get('gameDuration', 0), info.get('gameMode', ''),
                        info.get('gameVersion', ''), json.dumps(participant_data),
                        json.dumps(info.get('teams', []))
                    ))
                    
                    if cursor.rowcount > 0:
                        new_matches += 1
                        
                except Exception as e:
                    logging.error(f"Error saving match {match_data.get('metadata', {}).get('matchId', 'unknown')}: {e}")
                    continue
        
        return new_matches
    
    def get_user_puuid(self, user_id: int) -> str:
        """Get PUUID for user_id."""
        with self.db.get_connection() as conn:
            cursor = conn.execute("SELECT puuid FROM users WHERE id = ?", (user_id,))
            row = cursor.fetchone()
            return row[0] if row else None
    
    def get_matches_for_period(self, user_id: int, period: str) -> List[Dict[str, Any]]:
        """Get matches for specific time period."""
        now = datetime.now()
        
        if period == 'today':
            start_date = now.replace(hour=0, minute=0, second=0, microsecond=0)
        elif period == 'week':
            start_date = now - timedelta(days=7)
        elif period == 'month':
            start_date = now - timedelta(days=30)
        elif period == 'season':
            start_date = datetime(now.year, 1, 1)
        else:  # all
            start_date = datetime(2020, 1, 1)
        
        with self.db.get_connection() as conn:
            cursor = conn.execute("""
                SELECT * FROM matches 
                WHERE user_id = ? AND game_creation >= ?
                ORDER BY game_creation DESC
            """, (user_id, start_date))
            return [dict(row) for row in cursor.fetchall()]
    
    def get_latest_match_date(self, user_id: int) -> Optional[datetime]:
        """Get date of latest match for user."""
        with self.db.get_connection() as conn:
            cursor = conn.execute("""
                SELECT MAX(game_creation) FROM matches WHERE user_id = ?
            """, (user_id,))
            row = cursor.fetchone()
            if row and row[0]:
                return datetime.fromisoformat(row[0])
            return None


# Initialize global database manager
db_manager = DatabaseManager()
user_manager = UserManager(db_manager)
match_manager = MatchManager(db_manager)