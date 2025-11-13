# LUMEN Backend Deployment Guide

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database (via Supabase)
- Environment variables configured

## Local Development

### 1. Setup Environment

Copy the example environment file:
```bash
cp .env.example .env
```

Update `.env` with your credentials:
```env
# Server
PORT=8080
GIN_MODE=debug
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# Database (Supabase)
DB_HOST=db.your-project.supabase.co
DB_PORT=5432
DB_NAME=postgres
DB_USER=postgres
DB_PASSWORD=your-db-password
DB_SSLMODE=require

# Supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-role-key

# JWT
JWT_SECRET=your-secure-jwt-secret-minimum-32-characters

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Run Development Server

```bash
# Using Go directly
go run cmd/server/main.go

# Or using Make
make dev
```

The server will start on `http://localhost:8080`

### 4. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Expected response:
# {
#   "status": "ok",
#   "timestamp": "2025-11-13T10:00:00Z",
#   "service": "lumen-api",
#   "version": "1.0.0",
#   "database": "healthy"
# }
```

## Supabase Setup

### 1. Create Supabase Project

1. Go to https://supabase.com
2. Click "New Project"
3. Choose organization and region
4. Set database password
5. Wait for project to be created

### 2. Get Connection Details

From Supabase Dashboard:

1. Go to Settings → Database
2. Find "Connection string" section
3. Copy the connection parameters:
   - Host: `db.xxxxx.supabase.co`
   - Port: `5432`
   - Database: `postgres`
   - User: `postgres`
   - Password: (your database password)

4. Go to Settings → API
5. Copy:
   - Project URL (SUPABASE_URL)
   - Anon public key (SUPABASE_KEY)
   - Service role key (SUPABASE_SERVICE_KEY)

### 3. Create Database Schema

Run the following SQL in Supabase SQL Editor:

```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email TEXT UNIQUE NOT NULL,
  name TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Habits table
CREATE TABLE habits (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
  name TEXT NOT NULL,
  color TEXT NOT NULL,
  icon TEXT NOT NULL,
  frequency TEXT NOT NULL CHECK (frequency IN ('daily', 'weekly', 'monthly')),
  target_count INTEGER NOT NULL CHECK (target_count > 0),
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Tasks table
CREATE TABLE tasks (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
  title TEXT NOT NULL,
  description TEXT,
  horizon TEXT NOT NULL CHECK (horizon IN ('now', 'next', 'later', 'someday')),
  priority TEXT NOT NULL CHECK (priority IN ('low', 'medium', 'high', 'urgent')),
  status TEXT NOT NULL DEFAULT 'todo' CHECK (status IN ('todo', 'in_progress', 'done', 'archived')),
  due_date TIMESTAMPTZ,
  completed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Daily logs table
CREATE TABLE daily_logs (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
  date DATE NOT NULL,
  morning_routine BOOLEAN DEFAULT FALSE,
  evening_routine BOOLEAN DEFAULT FALSE,
  water_intake INTEGER DEFAULT 0 CHECK (water_intake >= 0 AND water_intake <= 20),
  sleep_hours DECIMAL(3,1) DEFAULT 0 CHECK (sleep_hours >= 0 AND sleep_hours <= 24),
  energy_level INTEGER CHECK (energy_level >= 1 AND energy_level <= 5),
  mood_rating INTEGER CHECK (mood_rating >= 1 AND mood_rating <= 5),
  productivity_rating INTEGER CHECK (productivity_rating >= 1 AND productivity_rating <= 5),
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(user_id, date)
);

-- Indexes for performance
CREATE INDEX idx_habits_user_id ON habits(user_id);
CREATE INDEX idx_habits_user_active ON habits(user_id, is_active);
CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_tasks_user_horizon ON tasks(user_id, horizon);
CREATE INDEX idx_tasks_user_status ON tasks(user_id, status);
CREATE INDEX idx_daily_logs_user_date ON daily_logs(user_id, date);

-- Row Level Security (RLS)
ALTER TABLE habits ENABLE ROW LEVEL SECURITY;
ALTER TABLE tasks ENABLE ROW LEVEL SECURITY;
ALTER TABLE daily_logs ENABLE ROW LEVEL SECURITY;

-- RLS Policies
CREATE POLICY "Users can view their own habits"
  ON habits FOR SELECT
  USING (auth.uid() = user_id);

CREATE POLICY "Users can insert their own habits"
  ON habits FOR INSERT
  WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update their own habits"
  ON habits FOR UPDATE
  USING (auth.uid() = user_id);

CREATE POLICY "Users can delete their own habits"
  ON habits FOR DELETE
  USING (auth.uid() = user_id);

-- Similar policies for tasks and daily_logs
CREATE POLICY "Users can view their own tasks"
  ON tasks FOR SELECT
  USING (auth.uid() = user_id);

CREATE POLICY "Users can insert their own tasks"
  ON tasks FOR INSERT
  WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update their own tasks"
  ON tasks FOR UPDATE
  USING (auth.uid() = user_id);

CREATE POLICY "Users can delete their own tasks"
  ON tasks FOR DELETE
  USING (auth.uid() = user_id);

CREATE POLICY "Users can view their own daily logs"
  ON daily_logs FOR SELECT
  USING (auth.uid() = user_id);

CREATE POLICY "Users can insert their own daily logs"
  ON daily_logs FOR INSERT
  WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update their own daily logs"
  ON daily_logs FOR UPDATE
  USING (auth.uid() = user_id);

CREATE POLICY "Users can delete their own daily logs"
  ON daily_logs FOR DELETE
  USING (auth.uid() = user_id);

-- Updated at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to all tables
CREATE TRIGGER update_habits_updated_at
  BEFORE UPDATE ON habits
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tasks_updated_at
  BEFORE UPDATE ON tasks
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_daily_logs_updated_at
  BEFORE UPDATE ON daily_logs
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();
```

## Building for Production

### 1. Build Binary

```bash
# Build for current platform
make build

# Or manually
go build -o bin/server cmd/server/main.go

# Cross-compile for Linux (from Windows/Mac)
GOOS=linux GOARCH=amd64 go build -o bin/server-linux cmd/server/main.go
```

### 2. Run Production Server

```bash
# Set production environment
export GIN_MODE=release
export LOG_LEVEL=info

# Run the binary
./bin/server
```

## Docker Deployment

### 1. Build Docker Image

```bash
make docker-build

# Or manually
docker build -t lumen-api:latest .
```

### 2. Run Container

```bash
# Create .env file with production values
docker run -d \
  --name lumen-api \
  -p 8080:8080 \
  --env-file .env \
  lumen-api:latest
```

### 3. Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - GIN_MODE=release
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_NAME=${DB_NAME}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_SSLMODE=require
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_KEY=${SUPABASE_KEY}
      - SUPABASE_SERVICE_KEY=${SUPABASE_SERVICE_KEY}
      - JWT_SECRET=${JWT_SECRET}
      - LOG_LEVEL=info
      - LOG_FORMAT=json
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

Run with:
```bash
docker-compose up -d
```

## Cloud Deployment Options

### Option 1: Railway

1. Install Railway CLI:
```bash
npm install -g @railway/cli
```

2. Login and initialize:
```bash
railway login
railway init
```

3. Add environment variables:
```bash
railway variables set PORT=8080
railway variables set GIN_MODE=release
railway variables set DB_HOST=your-db-host
# ... add all other variables
```

4. Deploy:
```bash
railway up
```

### Option 2: Render

1. Create `render.yaml`:

```yaml
services:
  - type: web
    name: lumen-api
    env: go
    buildCommand: go build -o server cmd/server/main.go
    startCommand: ./server
    envVars:
      - key: PORT
        value: 8080
      - key: GIN_MODE
        value: release
      - key: DB_HOST
        sync: false
      - key: DB_PASSWORD
        sync: false
      # Add other environment variables
```

2. Connect GitHub repo to Render
3. Add environment variables in Render dashboard
4. Deploy

### Option 3: Google Cloud Run

1. Build and push to Google Container Registry:
```bash
gcloud builds submit --tag gcr.io/PROJECT-ID/lumen-api
```

2. Deploy to Cloud Run:
```bash
gcloud run deploy lumen-api \
  --image gcr.io/PROJECT-ID/lumen-api \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars "GIN_MODE=release,PORT=8080,DB_HOST=..."
```

### Option 4: AWS ECS

1. Create ECR repository:
```bash
aws ecr create-repository --repository-name lumen-api
```

2. Build and push:
```bash
docker tag lumen-api:latest AWS_ACCOUNT_ID.dkr.ecr.REGION.amazonaws.com/lumen-api:latest
docker push AWS_ACCOUNT_ID.dkr.ecr.REGION.amazonaws.com/lumen-api:latest
```

3. Create ECS task definition and service through AWS Console or CLI

## Environment Variables Reference

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_HOST` | Database host | `db.xxxxx.supabase.co` |
| `DB_PORT` | Database port | `5432` |
| `DB_NAME` | Database name | `postgres` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `your-secure-password` |
| `DB_SSLMODE` | SSL mode | `require` |
| `SUPABASE_URL` | Supabase project URL | `https://xxxxx.supabase.co` |
| `SUPABASE_KEY` | Supabase anon key | `eyJhbGci...` |
| `JWT_SECRET` | JWT signing secret | `min-32-chars-secret` |

### Optional Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GIN_MODE` | Gin mode | `release` |
| `LOG_LEVEL` | Log level | `info` |
| `LOG_FORMAT` | Log format | `json` |
| `ALLOWED_ORIGINS` | CORS origins | `*` |
| `SUPABASE_SERVICE_KEY` | Service role key | - |

## Monitoring & Health Checks

### Health Check Endpoint

```bash
curl http://your-domain.com/health
```

Expected response:
```json
{
  "status": "ok",
  "timestamp": "2025-11-13T10:00:00Z",
  "service": "lumen-api",
  "version": "1.0.0",
  "database": "healthy"
}
```

### Readiness Check

```bash
curl http://your-domain.com/ready
```

### Metrics

```bash
curl http://your-domain.com/metrics
```

Returns database connection pool statistics.

## Security Checklist

- [ ] JWT_SECRET is at least 32 characters
- [ ] Database password is strong
- [ ] DB_SSLMODE is set to `require` in production
- [ ] GIN_MODE is set to `release`
- [ ] CORS origins are configured (not `*`)
- [ ] Rate limiting is enabled (default: 100 req/min)
- [ ] All secrets are in environment variables, not code
- [ ] Supabase RLS policies are enabled
- [ ] HTTPS is configured on hosting platform
- [ ] API is behind a firewall/security group

## Troubleshooting

### Database Connection Failed

1. Check database credentials in `.env`
2. Verify Supabase project is active
3. Check network connectivity
4. Ensure SSL mode is correct

```bash
# Test connection
psql "postgresql://postgres:password@db.xxx.supabase.co:5432/postgres?sslmode=require"
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### JWT Authentication Failing

1. Verify JWT_SECRET is set
2. Check token format: `Bearer <token>`
3. Ensure token hasn't expired
4. Validate token contains `user_id`

### High Memory Usage

1. Check database connection pool settings
2. Monitor for connection leaks
3. Review query efficiency

```go
// Adjust pool settings in config
config.MaxConns = 10
config.MinConns = 2
```

## Performance Optimization

### Database Indexes

Ensure these indexes exist:
```sql
CREATE INDEX IF NOT EXISTS idx_habits_user_id ON habits(user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_user_horizon ON tasks(user_id, horizon);
CREATE INDEX IF NOT EXISTS idx_daily_logs_user_date ON daily_logs(user_id, date);
```

### Connection Pool Tuning

In `config/config.go`:
```go
config.MaxConns = 25      // Max connections
config.MinConns = 5       // Min idle connections
config.MaxConnLifetime = time.Hour
config.MaxConnIdleTime = 30 * time.Minute
```

### Caching

Consider adding Redis for:
- Session storage
- API response caching
- Rate limit counters

## Backup & Recovery

### Database Backup

Supabase provides automatic backups. To manual backup:

```bash
# Export database
pg_dump "postgresql://postgres:password@db.xxx.supabase.co:5432/postgres" > backup.sql

# Restore
psql "postgresql://postgres:password@db.xxx.supabase.co:5432/postgres" < backup.sql
```

### Application State

All state is in the database. No application-level state to backup.

## CI/CD Pipeline

### GitHub Actions Example

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21

      - name: Run tests
        run: go test -v ./...

      - name: Build
        run: go build -o server cmd/server/main.go

      - name: Deploy to Railway
        run: |
          npm install -g @railway/cli
          railway up
        env:
          RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN }}
```

## Support

For issues or questions:
- Check API documentation: `/docs/API.md`
- Review logs: `docker logs lumen-api`
- Verify environment variables
- Test database connectivity
