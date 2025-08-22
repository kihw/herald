-- Migration pour le système de groupes
-- Version: 004
-- Description: Création des tables pour le système de groupes d'amis

-- Table des utilisateurs (étendue)
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    google_id TEXT UNIQUE,
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    picture TEXT,
    riot_id TEXT,
    riot_tag TEXT,
    region TEXT DEFAULT 'euw1',
    rank TEXT,
    lp INTEGER,
    mmr INTEGER,
    preferences TEXT DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT 1
);

-- Index pour la recherche d'utilisateurs
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_riot ON users(riot_id, riot_tag);
CREATE INDEX IF NOT EXISTS idx_users_region ON users(region);

-- Table des groupes
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    owner_id INTEGER NOT NULL,
    privacy TEXT DEFAULT 'private' CHECK (privacy IN ('public', 'private', 'invite_only')),
    invite_code TEXT UNIQUE NOT NULL,
    settings TEXT DEFAULT '{}',
    member_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Index pour les groupes
CREATE INDEX IF NOT EXISTS idx_groups_owner ON groups(owner_id);
CREATE INDEX IF NOT EXISTS idx_groups_privacy ON groups(privacy);
CREATE INDEX IF NOT EXISTS idx_groups_invite_code ON groups(invite_code);
CREATE INDEX IF NOT EXISTS idx_groups_name ON groups(name);

-- Table des membres de groupes
CREATE TABLE IF NOT EXISTS group_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role TEXT DEFAULT 'member' CHECK (role IN ('owner', 'admin', 'member')),
    status TEXT DEFAULT 'active' CHECK (status IN ('active', 'pending', 'banned', 'removed')),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    nickname TEXT,
    permissions TEXT DEFAULT '{}',
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(group_id, user_id)
);

-- Index pour les membres de groupes
CREATE INDEX IF NOT EXISTS idx_group_members_group ON group_members(group_id);
CREATE INDEX IF NOT EXISTS idx_group_members_user ON group_members(user_id);
CREATE INDEX IF NOT EXISTS idx_group_members_status ON group_members(status);
CREATE INDEX IF NOT EXISTS idx_group_members_role ON group_members(role);

-- Table des invitations de groupes
CREATE TABLE IF NOT EXISTS group_invites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    inviter_id INTEGER NOT NULL,
    invitee_id INTEGER,
    email TEXT,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined', 'expired')),
    message TEXT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (inviter_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (invitee_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Index pour les invitations
CREATE INDEX IF NOT EXISTS idx_group_invites_group ON group_invites(group_id);
CREATE INDEX IF NOT EXISTS idx_group_invites_inviter ON group_invites(inviter_id);
CREATE INDEX IF NOT EXISTS idx_group_invites_invitee ON group_invites(invitee_id);
CREATE INDEX IF NOT EXISTS idx_group_invites_email ON group_invites(email);
CREATE INDEX IF NOT EXISTS idx_group_invites_status ON group_invites(status);
CREATE INDEX IF NOT EXISTS idx_group_invites_expires ON group_invites(expires_at);

-- Table des statistiques de groupes
CREATE TABLE IF NOT EXISTS group_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    total_members INTEGER DEFAULT 0,
    active_members INTEGER DEFAULT 0,
    average_rank TEXT,
    average_mmr REAL,
    top_champions TEXT DEFAULT '[]',
    popular_roles TEXT DEFAULT '[]',
    winrate_comparison TEXT DEFAULT '{}',
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    UNIQUE(group_id)
);

-- Index pour les stats de groupes
CREATE INDEX IF NOT EXISTS idx_group_stats_group ON group_stats(group_id);
CREATE INDEX IF NOT EXISTS idx_group_stats_updated ON group_stats(last_updated);

-- Table des comparaisons de groupes
CREATE TABLE IF NOT EXISTS group_comparisons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    creator_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    compare_type TEXT NOT NULL CHECK (compare_type IN ('champions', 'roles', 'performance', 'trends')),
    parameters TEXT NOT NULL DEFAULT '{}',
    results TEXT DEFAULT '{}',
    is_public BOOLEAN DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Index pour les comparaisons
CREATE INDEX IF NOT EXISTS idx_group_comparisons_group ON group_comparisons(group_id);
CREATE INDEX IF NOT EXISTS idx_group_comparisons_creator ON group_comparisons(creator_id);
CREATE INDEX IF NOT EXISTS idx_group_comparisons_type ON group_comparisons(compare_type);
CREATE INDEX IF NOT EXISTS idx_group_comparisons_public ON group_comparisons(is_public);
CREATE INDEX IF NOT EXISTS idx_group_comparisons_created ON group_comparisons(created_at);

-- Triggers pour maintenir la cohérence des données

-- Trigger pour mettre à jour le compteur de membres lors de l'ajout
CREATE TRIGGER IF NOT EXISTS update_member_count_add
AFTER INSERT ON group_members
WHEN NEW.status = 'active'
BEGIN
    UPDATE groups 
    SET member_count = member_count + 1,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.group_id;
END;

-- Trigger pour mettre à jour le compteur de membres lors de la modification
CREATE TRIGGER IF NOT EXISTS update_member_count_update
AFTER UPDATE ON group_members
WHEN OLD.status != NEW.status
BEGIN
    UPDATE groups 
    SET member_count = (
        SELECT COUNT(*) FROM group_members 
        WHERE group_id = NEW.group_id AND status = 'active'
    ),
    updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.group_id;
END;

-- Trigger pour mettre à jour le compteur de membres lors de la suppression
CREATE TRIGGER IF NOT EXISTS update_member_count_delete
AFTER DELETE ON group_members
WHEN OLD.status = 'active'
BEGIN
    UPDATE groups 
    SET member_count = member_count - 1,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.group_id;
END;

-- Trigger pour expirer automatiquement les invitations
CREATE TRIGGER IF NOT EXISTS expire_old_invites
AFTER INSERT ON group_invites
BEGIN
    UPDATE group_invites 
    SET status = 'expired', updated_at = CURRENT_TIMESTAMP
    WHERE expires_at < CURRENT_TIMESTAMP AND status = 'pending';
END;

-- Trigger pour mettre à jour last_active des utilisateurs
CREATE TRIGGER IF NOT EXISTS update_user_last_active
AFTER UPDATE ON users
BEGIN
    UPDATE users 
    SET last_active = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;

-- Vues pour les requêtes courantes

-- Vue pour les groupes avec leurs propriétaires
CREATE VIEW IF NOT EXISTS groups_with_owners AS
SELECT 
    g.*,
    u.name as owner_name,
    u.email as owner_email,
    u.riot_id as owner_riot_id,
    u.riot_tag as owner_riot_tag
FROM groups g
LEFT JOIN users u ON g.owner_id = u.id;

-- Vue pour les membres actifs avec leurs infos utilisateur
CREATE VIEW IF NOT EXISTS active_group_members AS
SELECT 
    gm.*,
    u.name as user_name,
    u.email as user_email,
    u.riot_id as user_riot_id,
    u.riot_tag as user_riot_tag,
    u.region as user_region,
    u.rank as user_rank,
    u.lp as user_lp,
    u.mmr as user_mmr
FROM group_members gm
JOIN users u ON gm.user_id = u.id
WHERE gm.status = 'active';

-- Vue pour les invitations en attente avec les détails
CREATE VIEW IF NOT EXISTS pending_invites_detail AS
SELECT 
    gi.*,
    g.name as group_name,
    g.description as group_description,
    g.privacy as group_privacy,
    inviter.name as inviter_name,
    inviter.email as inviter_email,
    invitee.name as invitee_name,
    invitee.email as invitee_email
FROM group_invites gi
JOIN groups g ON gi.group_id = g.id
JOIN users inviter ON gi.inviter_id = inviter.id
LEFT JOIN users invitee ON gi.invitee_id = invitee.id
WHERE gi.status = 'pending' AND gi.expires_at > CURRENT_TIMESTAMP;

-- Données de test (optionnel, pour le développement)
-- INSERT OR IGNORE INTO users (id, email, name, riot_id, riot_tag, region) VALUES
-- (1, 'test1@example.com', 'TestUser1', 'TestPlayer1', 'EUW1', 'euw1'),
-- (2, 'test2@example.com', 'TestUser2', 'TestPlayer2', 'NA1', 'na1'),
-- (3, 'test3@example.com', 'TestUser3', 'TestPlayer3', 'EUW2', 'euw1');

-- INSERT OR IGNORE INTO groups (id, name, description, owner_id, privacy, invite_code) VALUES
-- (1, 'Test Group 1', 'A test group for development', 1, 'private', 'TESTCODE1'),
-- (2, 'Public Test Group', 'A public test group', 2, 'public', 'TESTCODE2');

-- INSERT OR IGNORE INTO group_members (group_id, user_id, role, status) VALUES
-- (1, 1, 'owner', 'active'),
-- (1, 2, 'member', 'active'),
-- (2, 2, 'owner', 'active'),
-- (2, 3, 'member', 'active');