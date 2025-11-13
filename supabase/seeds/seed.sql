-- LUMEN Seed Data
-- Created: 2025-11-13
-- Description: Sample data for testing and development

-- Note: This seed file assumes you have at least one user created via Supabase Auth
-- If not, sign up first via the frontend, then run this seed

-- Get the first user ID (or create a test user first)
-- For testing, you can replace this UUID with your actual user ID from auth.users

-- Sample Areas
INSERT INTO areas (id, user_id, name, icon, color, order_index)
SELECT
  gen_random_uuid(),
  id,
  'Health & Fitness',
  '💪',
  '#10B981',
  1
FROM users
LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO areas (id, user_id, name, icon, color, order_index)
SELECT
  gen_random_uuid(),
  id,
  'Career & Learning',
  '📚',
  '#3B82F6',
  2
FROM users
LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO areas (id, user_id, name, icon, color, order_index)
SELECT
  gen_random_uuid(),
  id,
  'Personal Growth',
  '🌱',
  '#8B5CF6',
  3
FROM users
LIMIT 1
ON CONFLICT DO NOTHING;

-- Sample Goals
INSERT INTO goals (id, user_id, area_id, title, timeframe, end_date, win_condition, description, status)
SELECT
  gen_random_uuid(),
  u.id,
  a.id,
  'Winter Arc Transformation',
  '3-6 months',
  CURRENT_DATE + INTERVAL '90 days',
  'Daily caloric deficit and consistent training',
  'Complete fitness transformation focusing on strength and conditioning',
  'active'
FROM users u
JOIN areas a ON a.user_id = u.id AND a.name = 'Health & Fitness'
LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO goals (id, user_id, area_id, title, timeframe, end_date, win_condition, description, status)
SELECT
  gen_random_uuid(),
  u.id,
  a.id,
  'Master Full-Stack Development',
  '2-year vision',
  CURRENT_DATE + INTERVAL '730 days',
  'Build and deploy 5 production applications',
  'Become proficient in React, Go, and cloud infrastructure',
  'active'
FROM users u
JOIN areas a ON a.user_id = u.id AND a.name = 'Career & Learning'
LIMIT 1
ON CONFLICT DO NOTHING;

-- Sample Habits
INSERT INTO habits (id, user_id, goal_id, name, frequency, reminder_times, icon)
SELECT
  gen_random_uuid(),
  u.id,
  g.id,
  'Morning Workout',
  'daily',
  '["06:00", "18:00"]'::jsonb,
  '🏋️'
FROM users u
JOIN goals g ON g.user_id = u.id AND g.title = 'Winter Arc Transformation'
LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO habits (id, user_id, goal_id, name, frequency, reminder_times, icon)
SELECT
  gen_random_uuid(),
  u.id,
  g.id,
  'Track Calories',
  'daily',
  '["12:00", "20:00"]'::jsonb,
  '🍎'
FROM users u
JOIN goals g ON g.user_id = u.id AND g.title = 'Winter Arc Transformation'
LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO habits (id, user_id, goal_id, name, frequency, reminder_times, icon)
SELECT
  gen_random_uuid(),
  u.id,
  g.id,
  'Code for 2 hours',
  'daily',
  '["09:00", "21:00"]'::jsonb,
  '💻'
FROM users u
JOIN goals g ON g.user_id = u.id AND g.title = 'Master Full-Stack Development'
LIMIT 1
ON CONFLICT DO NOTHING;

-- Sample Acceptance Criteria
INSERT INTO acceptance_criteria (id, user_id, criteria_text, day_type, order_index)
SELECT
  gen_random_uuid(),
  id,
  'Complete morning workout',
  'standard',
  1
FROM users
LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO acceptance_criteria (id, user_id, criteria_text, day_type, order_index)
SELECT
  gen_random_uuid(),
  id,
  'Track all meals and hit calorie target',
  'standard',
  2
FROM users
LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO acceptance_criteria (id, user_id, criteria_text, day_type, order_index)
SELECT
  gen_random_uuid(),
  id,
  'Code for at least 2 hours',
  'standard',
  3
FROM users
LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO acceptance_criteria (id, user_id, criteria_text, day_type, order_index)
SELECT
  gen_random_uuid(),
  id,
  'Plan tomorrow before 10 PM',
  'standard',
  4
FROM users
LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO acceptance_criteria (id, user_id, criteria_text, day_type, order_index)
SELECT
  gen_random_uuid(),
  id,
  'No junk food',
  'standard',
  5
FROM users
LIMIT 1
ON CONFLICT DO NOTHING;

-- Sample Tasks (2-day horizon)
INSERT INTO tasks (id, user_id, goal_id, title, due_date, horizon, notes, completed)
SELECT
  gen_random_uuid(),
  u.id,
  g.id,
  'Finish LUMEN authentication flow',
  CURRENT_TIMESTAMP + INTERVAL '1 day',
  '2-day',
  'Complete user signup, login, and email verification',
  false
FROM users u
JOIN goals g ON g.user_id = u.id AND g.title = 'Master Full-Stack Development'
LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO tasks (id, user_id, goal_id, title, due_date, horizon, notes, completed)
SELECT
  gen_random_uuid(),
  u.id,
  g.id,
  'Meal prep for the week',
  CURRENT_TIMESTAMP + INTERVAL '2 days',
  '2-day',
  'Prepare healthy meals for next 3 days',
  false
FROM users u
JOIN goals g ON g.user_id = u.id AND g.title = 'Winter Arc Transformation'
LIMIT 1
ON CONFLICT DO NOTHING;

-- Sample Tasks (7-day horizon)
INSERT INTO tasks (id, user_id, goal_id, title, due_date, horizon, notes, completed)
SELECT
  gen_random_uuid(),
  u.id,
  g.id,
  'Deploy LUMEN to production',
  CURRENT_TIMESTAMP + INTERVAL '7 days',
  '7-day',
  'Configure all environment variables and test end-to-end',
  false
FROM users u
JOIN goals g ON g.user_id = u.id AND g.title = 'Master Full-Stack Development'
LIMIT 1
ON CONFLICT DO NOTHING;

-- Sample Daily Log (yesterday - marked as won)
INSERT INTO daily_logs (id, user_id, date, goal_id, criteria_met, day_won, win_condition_met, reflection, planned_next_day)
SELECT
  gen_random_uuid(),
  u.id,
  CURRENT_DATE - INTERVAL '1 day',
  g.id,
  (SELECT jsonb_agg(id) FROM acceptance_criteria WHERE user_id = u.id LIMIT 4),
  true,
  true,
  'Great day! Hit all my targets and felt energized.',
  true
FROM users u
JOIN goals g ON g.user_id = u.id AND g.title = 'Winter Arc Transformation'
LIMIT 1
ON CONFLICT (user_id, date, goal_id) DO NOTHING;

-- Sample Habit Logs (last 7 days for morning workout)
INSERT INTO habit_logs (id, habit_id, user_id, logged_at, date, completed, notes)
SELECT
  gen_random_uuid(),
  h.id,
  h.user_id,
  (CURRENT_DATE - INTERVAL '1 day' + TIME '06:30:00')::timestamp,
  CURRENT_DATE - INTERVAL '1 day',
  true,
  'Full body workout - felt strong'
FROM habits h
WHERE h.name = 'Morning Workout'
LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO habit_logs (id, habit_id, user_id, logged_at, date, completed, notes)
SELECT
  gen_random_uuid(),
  h.id,
  h.user_id,
  (CURRENT_DATE - INTERVAL '2 days' + TIME '06:30:00')::timestamp,
  CURRENT_DATE - INTERVAL '2 days',
  true,
  'Upper body focus'
FROM habits h
WHERE h.name = 'Morning Workout'
LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO habit_logs (id, habit_id, user_id, logged_at, date, completed, notes)
SELECT
  gen_random_uuid(),
  h.id,
  h.user_id,
  (CURRENT_DATE - INTERVAL '3 days' + TIME '06:30:00')::timestamp,
  CURRENT_DATE - INTERVAL '3 days',
  true,
  'Legs day - challenging but completed'
FROM habits h
WHERE h.name = 'Morning Workout'
LIMIT 1
ON CONFLICT DO NOTHING;

-- Seed data complete
-- You now have sample areas, goals, habits, tasks, criteria, logs for testing
