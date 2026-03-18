-- Migration: Admin Password Login
-- Description: Adds password_hash to users table for admin dashboard login
-- Date: 2026-03-17

ALTER TABLE users
ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);

COMMENT ON COLUMN users.password_hash IS 'bcrypt hash for admin dashboard login. NULL = PlayFab-only account';
