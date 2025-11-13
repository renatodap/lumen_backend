-- Final user schema migration (handles existing constraints)
-- Created: 2025-11-13
-- Run this if you get "constraint already exists" errors

-- Add columns (IF NOT EXISTS handles re-runs)
ALTER TABLE users ADD COLUMN IF NOT EXISTS full_name TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS role TEXT DEFAULT 'user';
ALTER TABLE users ADD COLUMN IF NOT EXISTS subscription_tier TEXT DEFAULT 'free';
ALTER TABLE users ADD COLUMN IF NOT EXISTS subscription_status TEXT DEFAULT 'active';
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified BOOLEAN DEFAULT false;

-- Drop old constraints first (safe - uses IF EXISTS)
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_subscription_tier_check;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_subscription_status_check;

-- Re-create constraints (fresh, no conflicts)
ALTER TABLE users ADD CONSTRAINT users_role_check
  CHECK (role IN ('user', 'admin', 'premium'));

ALTER TABLE users ADD CONSTRAINT users_subscription_tier_check
  CHECK (subscription_tier IN ('free', 'starter', 'pro', 'enterprise'));

ALTER TABLE users ADD CONSTRAINT users_subscription_status_check
  CHECK (subscription_status IN ('active', 'cancelled', 'past_due', 'trialing'));

-- Create indexes (IF NOT EXISTS handles re-runs)
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_subscription_tier ON users(subscription_tier);
CREATE INDEX IF NOT EXISTS idx_users_email_verified ON users(email_verified);

-- Update trigger (drop and recreate for clean state)
DROP TRIGGER IF EXISTS on_auth_user_created ON auth.users;
DROP FUNCTION IF EXISTS public.handle_new_user();

CREATE OR REPLACE FUNCTION public.handle_new_user()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO public.users (id, email, email_verified, created_at)
  VALUES (
    NEW.id,
    NEW.email,
    NEW.email_confirmed_at IS NOT NULL,
    NEW.created_at
  )
  ON CONFLICT (id) DO UPDATE SET
    email = NEW.email,
    email_verified = NEW.email_confirmed_at IS NOT NULL,
    updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE TRIGGER on_auth_user_created
  AFTER INSERT OR UPDATE ON auth.users
  FOR EACH ROW
  EXECUTE FUNCTION public.handle_new_user();

-- Verify schema
DO $$
DECLARE
  col_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO col_count
  FROM information_schema.columns
  WHERE table_name = 'users'
  AND column_name IN ('full_name', 'role', 'subscription_tier', 'subscription_status', 'is_active', 'email_verified');

  IF col_count = 6 THEN
    RAISE NOTICE 'SUCCESS: All 6 new columns exist in users table';
  ELSE
    RAISE NOTICE 'WARNING: Only % of 6 columns exist', col_count;
  END IF;
END $$;

-- Migration complete - Safe to re-run multiple times
