-- Migration: Convert users table to OAuth-based authentication
-- This migration updates the users table to work with Riot OAuth instead of email/password

-- First, add new columns for OAuth
ALTER TABLE users 
ADD COLUMN riot_id VARCHAR(50),
ADD COLUMN riot_tag VARCHAR(10),
ADD COLUMN region VARCHAR(10) DEFAULT 'euw1',
ADD COLUMN last_region_used VARCHAR(10),
ADD COLUMN access_token TEXT,
ADD COLUMN refresh_token TEXT,
ADD COLUMN token_expires_at TIMESTAMP WITH TIME ZONE;

-- Make riot_puuid NOT NULL (it was optional before)
ALTER TABLE users ALTER COLUMN riot_puuid SET NOT NULL;

-- Drop old authentication columns (only if they exist)
DO $$ 
BEGIN
    -- Check if email column exists before dropping
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'email') THEN
        ALTER TABLE users DROP COLUMN email;
    END IF;
    
    -- Check if password_hash column exists before dropping
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'password_hash') THEN
        ALTER TABLE users DROP COLUMN password_hash;
    END IF;
    
    -- Check if username column exists before dropping
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'username') THEN
        ALTER TABLE users DROP COLUMN username;
    END IF;
    
    -- Check if tagline column exists before dropping
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'tagline') THEN
        ALTER TABLE users DROP COLUMN tagline;
    END IF;
END $$;

-- Create unique index on riot_puuid (primary identifier)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_riot_puuid ON users(riot_puuid);

-- Create index on riot_id and riot_tag for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_riot_id_tag ON users(riot_id, riot_tag);

-- Create index on region for filtering
CREATE INDEX IF NOT EXISTS idx_users_region ON users(region);

-- Update existing data if any (this is a one-time migration, so existing users will need to re-authenticate)
-- Since we're changing the authentication system completely, we'll clean existing test data
TRUNCATE TABLE users CASCADE;

-- Update user_settings to reference correct platform format
-- Make sure platform values are valid region codes
UPDATE user_settings SET platform = 'euw1' WHERE platform NOT IN ('br1', 'eun1', 'euw1', 'jp1', 'kr', 'la1', 'la2', 'na1', 'oc1', 'ph2', 'ru', 'sg2', 'th2', 'tr1', 'tw2', 'vn2');
