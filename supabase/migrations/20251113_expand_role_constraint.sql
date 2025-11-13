-- Expand role constraint to allow more user types
-- Created: 2025-11-13
-- Fixes: new row violates check constraint "users_role_check"

-- Drop existing role constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;

-- Add expanded role constraint with more options
ALTER TABLE users ADD CONSTRAINT users_role_check
  CHECK (role IN ('user', 'admin', 'premium', 'educator', 'instructor', 'student', 'moderator', 'contributor'));

-- Success message
DO $$
BEGIN
  RAISE NOTICE 'SUCCESS: Role constraint expanded to allow: user, admin, premium, educator, instructor, student, moderator, contributor';
END $$;

-- Migration complete - Your seed file should now work
