-- LUMEN Initial Database Schema
-- Created: 2025-11-13
-- Description: Complete schema for LUMEN personal OS enforcer

-- Users table with timezone support
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  timezone TEXT NOT NULL DEFAULT 'UTC',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Areas table (PARA method)
CREATE TABLE IF NOT EXISTS areas (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  icon TEXT,
  color TEXT,
  order_index INTEGER DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Goals table (Projects in PARA)
CREATE TABLE IF NOT EXISTS goals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  area_id UUID REFERENCES areas(id) ON DELETE SET NULL,
  title TEXT NOT NULL,
  timeframe TEXT,
  end_date DATE,
  win_condition TEXT,
  description TEXT,
  status TEXT DEFAULT 'active',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Habits table
CREATE TABLE IF NOT EXISTS habits (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  goal_id UUID REFERENCES goals(id) ON DELETE SET NULL,
  name TEXT NOT NULL,
  frequency TEXT,
  reminder_times JSONB,
  icon TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Tasks table with horizon system
CREATE TABLE IF NOT EXISTS tasks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  goal_id UUID REFERENCES goals(id) ON DELETE SET NULL,
  title TEXT NOT NULL,
  due_date TIMESTAMPTZ,
  horizon TEXT,
  notes TEXT,
  completed BOOLEAN DEFAULT FALSE,
  completed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Acceptance Criteria table
CREATE TABLE IF NOT EXISTS acceptance_criteria (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  criteria_text TEXT NOT NULL,
  day_type TEXT DEFAULT 'standard',
  order_index INTEGER DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Daily Logs table
CREATE TABLE IF NOT EXISTS daily_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  date DATE NOT NULL,
  goal_id UUID REFERENCES goals(id) ON DELETE SET NULL,
  criteria_met JSONB,
  day_won BOOLEAN DEFAULT FALSE,
  win_condition_met BOOLEAN,
  reflection TEXT,
  planned_next_day BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(user_id, date, goal_id)
);

-- Habit Logs table
CREATE TABLE IF NOT EXISTS habit_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  habit_id UUID REFERENCES habits(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  logged_at TIMESTAMPTZ DEFAULT NOW(),
  date DATE NOT NULL,
  completed BOOLEAN DEFAULT TRUE,
  notes TEXT
);

-- Provider Connections table (for OAuth integrations)
CREATE TABLE IF NOT EXISTS provider_connections (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  provider TEXT NOT NULL,
  access_token TEXT,
  refresh_token TEXT,
  token_expires_at TIMESTAMPTZ,
  connected_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(user_id, provider)
);

-- Performance indexes
CREATE INDEX IF NOT EXISTS idx_habits_user ON habits(user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_user_horizon ON tasks(user_id, horizon);
CREATE INDEX IF NOT EXISTS idx_daily_logs_user_date ON daily_logs(user_id, date);
CREATE INDEX IF NOT EXISTS idx_habit_logs_habit_date ON habit_logs(habit_id, date);
CREATE INDEX IF NOT EXISTS idx_areas_user ON areas(user_id);
CREATE INDEX IF NOT EXISTS idx_goals_user ON goals(user_id);
CREATE INDEX IF NOT EXISTS idx_goals_status ON goals(user_id, status);
CREATE INDEX IF NOT EXISTS idx_tasks_completed ON tasks(user_id, completed);

-- Enable Row Level Security (RLS) on all tables
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE areas ENABLE ROW LEVEL SECURITY;
ALTER TABLE goals ENABLE ROW LEVEL SECURITY;
ALTER TABLE habits ENABLE ROW LEVEL SECURITY;
ALTER TABLE tasks ENABLE ROW LEVEL SECURITY;
ALTER TABLE acceptance_criteria ENABLE ROW LEVEL SECURITY;
ALTER TABLE daily_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE habit_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE provider_connections ENABLE ROW LEVEL SECURITY;

-- Drop existing policies if they exist (for re-running migration)
DROP POLICY IF EXISTS "Users can CRUD their own data" ON users;
DROP POLICY IF EXISTS "Users can CRUD their own areas" ON areas;
DROP POLICY IF EXISTS "Users can CRUD their own goals" ON goals;
DROP POLICY IF EXISTS "Users can CRUD their own habits" ON habits;
DROP POLICY IF EXISTS "Users can CRUD their own tasks" ON tasks;
DROP POLICY IF EXISTS "Users can CRUD their own criteria" ON acceptance_criteria;
DROP POLICY IF EXISTS "Users can CRUD their own logs" ON daily_logs;
DROP POLICY IF EXISTS "Users can CRUD their own habit logs" ON habit_logs;
DROP POLICY IF EXISTS "Users can CRUD their own connections" ON provider_connections;

-- RLS Policies: Users can only access their own data
CREATE POLICY "Users can CRUD their own data" ON users
  FOR ALL USING (auth.uid() = id);

CREATE POLICY "Users can CRUD their own areas" ON areas
  FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can CRUD their own goals" ON goals
  FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can CRUD their own habits" ON habits
  FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can CRUD their own tasks" ON tasks
  FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can CRUD their own criteria" ON acceptance_criteria
  FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can CRUD their own logs" ON daily_logs
  FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can CRUD their own habit logs" ON habit_logs
  FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can CRUD their own connections" ON provider_connections
  FOR ALL USING (auth.uid() = user_id);

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for updated_at
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_goals_updated_at ON goals;
CREATE TRIGGER update_goals_updated_at
  BEFORE UPDATE ON goals
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- Migration complete
-- Run this in Supabase SQL Editor to set up the complete database schema
