# LUMEN Deployment Guide

**Production Deployment: Supabase + Railway + Vercel**

---

## Prerequisites

1. Supabase account (database)
2. Railway account (backend hosting)
3. Vercel account (frontend hosting)
4. GitHub repositories (already created)

---

## Step 1: Supabase Database Setup

### 1.1 Create Supabase Project
```
1. Go to https://supabase.com/dashboard
2. Click "New Project"
3. Project name: "lumen-production"
4. Database password: [SAVE THIS - YOU'LL NEED IT]
5. Region: Choose closest to your users (e.g., us-east-1)
6. Click "Create new project" (takes ~2 minutes)
```

### 1.2 Get Connection Details
```
1. Go to Project Settings > Database
2. Copy these values:

   - Project URL: https://[project-ref].supabase.co
   - Anon/Public Key: eyJh... (starts with eyJ)
   - Service Role Key: eyJh... (different from anon key)
   - Database Password: [the one you set in 1.1]

3. Connection string format:
   postgresql://postgres:[YOUR-PASSWORD]@db.[project-ref].supabase.co:5432/postgres
```

### 1.3 Run Database Migrations
```sql
-- Go to SQL Editor in Supabase Dashboard
-- Run this migration to create all tables:

-- Users table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  timezone TEXT NOT NULL DEFAULT 'UTC',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Areas table
CREATE TABLE areas (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  icon TEXT,
  color TEXT,
  order_index INTEGER DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Goals table
CREATE TABLE goals (
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
CREATE TABLE habits (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  goal_id UUID REFERENCES goals(id) ON DELETE SET NULL,
  name TEXT NOT NULL,
  frequency TEXT,
  reminder_times JSONB,
  icon TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Tasks table
CREATE TABLE tasks (
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
CREATE TABLE acceptance_criteria (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  criteria_text TEXT NOT NULL,
  day_type TEXT DEFAULT 'standard',
  order_index INTEGER DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Daily Logs table
CREATE TABLE daily_logs (
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
CREATE TABLE habit_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  habit_id UUID REFERENCES habits(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  logged_at TIMESTAMPTZ DEFAULT NOW(),
  date DATE NOT NULL,
  completed BOOLEAN DEFAULT TRUE,
  notes TEXT
);

-- Provider Connections table
CREATE TABLE provider_connections (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  provider TEXT NOT NULL,
  access_token TEXT,
  refresh_token TEXT,
  token_expires_at TIMESTAMPTZ,
  connected_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(user_id, provider)
);

-- Indexes for performance
CREATE INDEX idx_habits_user ON habits(user_id);
CREATE INDEX idx_tasks_user_horizon ON tasks(user_id, horizon);
CREATE INDEX idx_daily_logs_user_date ON daily_logs(user_id, date);
CREATE INDEX idx_habit_logs_habit_date ON habit_logs(habit_id, date);
CREATE INDEX idx_areas_user ON areas(user_id);
CREATE INDEX idx_goals_user ON goals(user_id);

-- Enable Row Level Security
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE areas ENABLE ROW LEVEL SECURITY;
ALTER TABLE goals ENABLE ROW LEVEL SECURITY;
ALTER TABLE habits ENABLE ROW LEVEL SECURITY;
ALTER TABLE tasks ENABLE ROW LEVEL SECURITY;
ALTER TABLE acceptance_criteria ENABLE ROW LEVEL SECURITY;
ALTER TABLE daily_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE habit_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE provider_connections ENABLE ROW LEVEL SECURITY;

-- RLS Policies (Users can only access their own data)
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
```

### 1.4 Enable Supabase Auth
```
1. Go to Authentication > Providers
2. Enable Email provider (default)
3. Optional: Enable Google OAuth, Microsoft OAuth
4. Configure email templates (Settings > Auth > Email Templates)
```

---

## Step 2: Railway Backend Deployment

### 2.1 Create Railway Project
```
1. Go to https://railway.app/dashboard
2. Click "New Project"
3. Select "Deploy from GitHub repo"
4. Choose repository: renatodap/lumen_backend
5. Railway auto-detects Dockerfile
6. Click "Deploy"
```

### 2.2 Configure Environment Variables
```
1. In Railway dashboard, go to your service
2. Click "Variables" tab
3. Add these variables (click "Add Variable" for each):

   GIN_MODE=release
   APP_PORT=8080

   DATABASE_URL=postgresql://postgres:[PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres

   JWT_SECRET=[GENERATE A RANDOM 32+ CHARACTER STRING]
   JWT_EXPIRY=24h
   REFRESH_TOKEN_EXPIRY=168h

   SUPABASE_URL=https://[PROJECT-REF].supabase.co
   SUPABASE_SERVICE_KEY=[SERVICE ROLE KEY FROM SUPABASE]
   SUPABASE_JWT_SECRET=[JWT SECRET FROM SUPABASE SETTINGS]

   CORS_ALLOWED_ORIGINS=https://lumen-frontend-theta.vercel.app

   LOG_LEVEL=info
   RATE_LIMIT_REQUESTS=100
   RATE_LIMIT_WINDOW=60s
   RATE_LIMIT_ENABLED=true

4. Click "Deploy" to restart with new variables
```

### 2.3 Get Railway Backend URL
```
1. Go to Settings tab
2. Scroll to "Domains"
3. Click "Generate Domain"
4. Copy the URL: https://lumen-backend-production-xxxx.up.railway.app
5. Save this - you'll need it for frontend
```

---

## Step 3: Vercel Frontend Deployment

### 3.1 Create Vercel Project
```
1. Go to https://vercel.com/dashboard
2. Click "Add New" > "Project"
3. Import repository: renatodap/lumen_frontend
4. Framework: Next.js (auto-detected)
5. Configure project (DON'T deploy yet)
```

### 3.2 Configure Environment Variables
```
1. In project settings, go to "Environment Variables"
2. Add these variables:

   NEXT_PUBLIC_SUPABASE_URL=https://[PROJECT-REF].supabase.co
   NEXT_PUBLIC_SUPABASE_ANON_KEY=[ANON KEY FROM SUPABASE]
   NEXT_PUBLIC_API_URL=https://[RAILWAY-BACKEND-URL]/api
   NODE_ENV=production

3. Apply to: Production, Preview, and Development
4. Click "Deploy"
```

### 3.3 Update Backend CORS
```
1. Go back to Railway backend
2. Update CORS_ALLOWED_ORIGINS variable to include Vercel URL:

   CORS_ALLOWED_ORIGINS=https://lumen-frontend-theta.vercel.app,https://lumen-frontend-git-main-renatodaps-projects.vercel.app

3. Redeploy backend
```

---

## Step 4: Verification

### 4.1 Test Backend Health
```bash
curl https://[RAILWAY-URL]/health

# Expected response:
{
  "status": "healthy",
  "timestamp": "2025-11-13T17:30:00Z"
}
```

### 4.2 Test Frontend
```
1. Open https://lumen-frontend-theta.vercel.app
2. Should see LUMEN homepage
3. Try to sign up/login (Supabase Auth)
4. Open DevTools > Network tab
5. Check API calls go to Railway backend
```

### 4.3 Test Database Connection
```
1. Try creating a habit in frontend
2. Check Supabase dashboard > Table Editor
3. Verify data appears in habits table
```

---

## Step 5: Post-Deployment Configuration

### 5.1 Generate JWT Secret (if not done)
```bash
# Use this command to generate secure secrets:
openssl rand -base64 32
# Copy output and use as JWT_SECRET
```

### 5.2 Supabase JWT Secret Location
```
1. Go to Supabase Dashboard > Settings > API
2. Scroll to "JWT Settings"
3. Copy "JWT Secret" value
4. Use this for SUPABASE_JWT_SECRET in Railway
```

### 5.3 Update Frontend API URL (if backend URL changes)
```
1. Vercel Dashboard > Project > Settings > Environment Variables
2. Edit NEXT_PUBLIC_API_URL
3. Redeploy frontend
```

---

## Environment Variables Summary

### Frontend (Vercel)
```env
NEXT_PUBLIC_SUPABASE_URL=https://xxxxx.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbG...
NEXT_PUBLIC_API_URL=https://backend.railway.app/api
NODE_ENV=production
```

### Backend (Railway)
```env
GIN_MODE=release
APP_PORT=8080
DATABASE_URL=postgresql://postgres:pass@db.xxxxx.supabase.co:5432/postgres
JWT_SECRET=your-32-char-secret
SUPABASE_URL=https://xxxxx.supabase.co
SUPABASE_SERVICE_KEY=eyJhbG...
SUPABASE_JWT_SECRET=your-jwt-secret
CORS_ALLOWED_ORIGINS=https://lumen-frontend-theta.vercel.app
```

---

## Troubleshooting

### Backend won't connect to database
```
- Check DATABASE_URL format is correct
- Verify password has no special characters that need URL encoding
- Test connection from Railway logs
- Check Supabase firewall allows Railway IP
```

### Frontend can't reach backend
```
- Verify NEXT_PUBLIC_API_URL is correct
- Check CORS_ALLOWED_ORIGINS includes Vercel domain
- Check Railway backend is running (logs)
- Test backend /health endpoint directly
```

### Authentication not working
```
- Verify SUPABASE_URL and keys match frontend/backend
- Check Supabase Auth is enabled
- Review Supabase Auth logs
- Verify RLS policies are correct
```

### Build failures
```
- Frontend: Check TypeScript errors, missing dependencies
- Backend: Check Go module errors, Dockerfile syntax
- Review build logs in Railway/Vercel dashboards
```

---

## Deployment Checklist

- [ ] Supabase project created
- [ ] Database schema migrated
- [ ] RLS policies enabled
- [ ] Railway backend deployed
- [ ] Backend environment variables configured
- [ ] Backend health check passing
- [ ] Vercel frontend deployed
- [ ] Frontend environment variables configured
- [ ] CORS configured correctly
- [ ] Test user signup/login
- [ ] Test habit creation
- [ ] Test data persistence
- [ ] Monitor logs for errors

---

## Next Steps

1. **Set up monitoring**: Add Sentry or LogRocket
2. **Configure custom domain**: Point your domain to Vercel
3. **Enable analytics**: Vercel Analytics, Supabase Analytics
4. **Set up backups**: Supabase automated backups
5. **CI/CD**: GitHub Actions for automated testing before deploy

---

**Status**: Both frontend and backend are deployed and running!
**Frontend**: https://lumen-frontend-theta.vercel.app
**Backend**: https://[railway-url] (check Railway dashboard for exact URL)
