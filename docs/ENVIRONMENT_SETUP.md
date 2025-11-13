# LUMEN Environment Configuration Guide

**Complete environment variable setup for all deployment scenarios**

---

## 🔑 Your Supabase Credentials

```
Project URL:        https://ocopuqketaddqhjwsdjd.supabase.co
Anon Key:           eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im9jb3B1cWtldGFkZHFoandzZGpkIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjMwNDMwNTQsImV4cCI6MjA3ODYxOTA1NH0.qOGQT48cxLN-F2mjyMP25sN-ch_LSEnQGkceYLwfyaw
Service Role Key:   eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im9jb3B1cWtldGFkZHFoandzZGpkIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTc2MzA0MzA1NCwiZXhwIjoyMDc4NjE5MDU0fQ.KBijGXuYttkQhrSR4ua0rPWyq2oyAU1kngTw-O_M0cE
JWT Secret:         5s8LWi7vS3Z0xGBs35Z4p80Yx4ernAXoB6Sp9EDX9IEAIe4eXcXiCx5zAQtjdBQcsxPWIhycqhWGUEZx1mo4Kg==
Generated JWT:      Z7AO/XN5EERiDwKyrFXvJdU+va9M1HGd8Zx2UzaHs58=
```

**⚠️ CRITICAL: Database Password Missing**
You need to get your database password from Supabase dashboard:
1. Go to https://supabase.com/dashboard
2. Click on your project: ocopuqketaddqhjwsdjd
3. Go to Settings > Database
4. Look for "Database Password" section
5. If you forgot it, click "Reset Database Password"

---

## 📁 Environment Files Created

### ✅ Frontend Files
```
frontend/.env.local           ← Local development (NOT committed to git)
frontend/.env.example         ← Template with your values (committed to git)
```

### ✅ Backend Files
```
backend-go/.env               ← Local development (NOT committed to git)
backend-go/.env.example       ← Template with your values (committed to git)
```

---

## 🚀 Production Deployment Configuration

### 1️⃣ Railway Backend Environment Variables

**Go to Railway Dashboard → lumen_backend → Variables tab**

Copy and paste these EXACTLY (replace [YOUR-DB-PASSWORD] with actual password):

```env
GIN_MODE=release
APP_PORT=8080

DATABASE_URL=postgresql://postgres.ocopuqketaddqhjwsdjd:[YOUR-DB-PASSWORD]@aws-0-us-east-1.pooler.supabase.com:6543/postgres

JWT_SECRET=Z7AO/XN5EERiDwKyrFXvJdU+va9M1HGd8Zx2UzaHs58=
JWT_EXPIRY=24h
REFRESH_TOKEN_EXPIRY=168h

SUPABASE_URL=https://ocopuqketaddqhjwsdjd.supabase.co
SUPABASE_SERVICE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im9jb3B1cWtldGFkZHFoandzZGpkIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTc2MzA0MzA1NCwiZXhwIjoyMDc4NjE5MDU0fQ.KBijGXuYttkQhrSR4ua0rPWyq2oyAU1kngTw-O_M0cE
SUPABASE_JWT_SECRET=5s8LWi7vS3Z0xGBs35Z4p80Yx4ernAXoB6Sp9EDX9IEAIe4eXcXiCx5zAQtjdBQcsxPWIhycqhWGUEZx1mo4Kg==

CORS_ALLOWED_ORIGINS=https://lumen-frontend-theta.vercel.app,https://lumen-frontend-git-main-renatodaps-projects.vercel.app

LOG_LEVEL=info
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
RATE_LIMIT_ENABLED=true
```

**Step-by-step:**
1. Click "Add Variable" for each line
2. Variable name = left side (before `=`)
3. Value = right side (after `=`)
4. Click anywhere to save
5. Railway will auto-redeploy

---

### 2️⃣ Get Railway Backend URL

**Before configuring Vercel, you need your Railway URL:**

1. Go to Railway Dashboard → lumen_backend → Settings
2. Scroll to "Networking" section
3. Click "Generate Domain"
4. Copy the URL: `https://lumen-backend-production-[random].up.railway.app`
5. **Save this URL** - you'll need it for Vercel!

Your Railway URL: _______________________________________________

---

### 3️⃣ Vercel Frontend Environment Variables

**Go to Vercel Dashboard → lumen-frontend → Settings → Environment Variables**

Add these variables (use your Railway URL from step 2):

```env
NEXT_PUBLIC_SUPABASE_URL=https://ocopuqketaddqhjwsdjd.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im9jb3B1cWtldGFkZHFoandzZGpkIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjMwNDMwNTQsImV4cCI6MjA3ODYxOTA1NH0.qOGQT48cxLN-F2mjyMP25sN-ch_LSEnQGkceYLwfyaw
NEXT_PUBLIC_API_URL=[YOUR-RAILWAY-URL]/api
NODE_ENV=production
```

**Step-by-step:**
1. Click "Add New" for each variable
2. Key = variable name (left side)
3. Value = value (right side)
4. Select: ✓ Production ✓ Preview ✓ Development
5. Click "Save"
6. After all 4 variables added, click "Redeploy"

---

## 🧪 Local Development Setup

### Frontend Local Development

**File: `frontend/.env.local` (already created)**
```env
NEXT_PUBLIC_SUPABASE_URL=https://ocopuqketaddqhjwsdjd.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im9jb3B1cWtldGFkZHFoandzZGpkIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjMwNDMwNTQsImV4cCI6MjA3ODYxOTA1NH0.qOGQT48cxLN-F2mjyMP25sN-ch_LSEnQGkceYLwfyaw
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NODE_ENV=development
```

**To run locally:**
```bash
cd frontend
npm install
npm run dev
# Frontend: http://localhost:3000
```

---

### Backend Local Development

**File: `backend-go/.env` (already created)**

⚠️ **YOU MUST ADD YOUR DATABASE PASSWORD** in this file:
```env
DATABASE_URL=postgresql://postgres.ocopuqketaddqhjwsdjd:[YOUR-DB-PASSWORD]@aws-0-us-east-1.pooler.supabase.com:6543/postgres
```

**To run locally:**
```bash
cd backend-go
go mod download
go run cmd/server/main.go
# Backend: http://localhost:8080
```

---

## 📋 Quick Copy-Paste Checklist

### ☐ Step 1: Get Database Password
```
1. Go to Supabase Dashboard
2. Project Settings > Database
3. Get/Reset database password
4. Save password: _______________
```

### ☐ Step 2: Configure Railway (13 variables)
```
1. Railway Dashboard > lumen_backend > Variables
2. Copy all variables from section 1️⃣ above
3. Replace [YOUR-DB-PASSWORD] with actual password
4. Wait for auto-redeploy (~2 minutes)
```

### ☐ Step 3: Get Railway URL
```
1. Railway > Settings > Networking > Generate Domain
2. Copy URL: _______________
```

### ☐ Step 4: Configure Vercel (4 variables)
```
1. Vercel Dashboard > lumen-frontend > Settings > Env Variables
2. Copy all variables from section 3️⃣ above
3. Replace [YOUR-RAILWAY-URL] with actual Railway URL
4. Click "Redeploy"
```

### ☐ Step 5: Update Local .env Files
```
1. Edit: backend-go/.env
2. Add your database password to DATABASE_URL
3. Save file
4. Run: go run cmd/server/main.go
```

---

## 🧪 Testing Your Configuration

### Test Backend Health
```bash
# Local:
curl http://localhost:8080/health

# Production:
curl https://[YOUR-RAILWAY-URL]/health

# Expected response:
{"status":"healthy","timestamp":"2025-11-13T..."}
```

### Test Frontend Connection
```bash
# Local:
Open http://localhost:3000
Check browser console (F12) for errors

# Production:
Open https://lumen-frontend-theta.vercel.app
Check browser console for errors
```

### Test Database Connection
```bash
# In Supabase Dashboard:
1. Go to SQL Editor
2. Run: SELECT NOW();
3. Should return current timestamp
```

---

## 🔐 Security Notes

**✅ Files that ARE committed to git:**
- `frontend/.env.example` - Template with your Supabase values
- `backend-go/.env.example` - Template with your Supabase values

**🚫 Files that are NOT committed (in .gitignore):**
- `frontend/.env.local` - Your local development config
- `backend-go/.env` - Your local development config

**⚠️ Never commit:**
- Database passwords
- Service role keys (except in .env.example for your reference)
- JWT secrets

---

## 🐛 Troubleshooting

### Backend won't connect to database
```
Error: "connection refused" or "authentication failed"

Fix:
1. Verify DATABASE_URL password is correct
2. Check Supabase project is active (not paused)
3. Try resetting database password in Supabase
4. Use connection pooler URL (port 6543, not 5432)
```

### Frontend can't reach backend
```
Error: "CORS policy" or "Network Error"

Fix:
1. Verify CORS_ALLOWED_ORIGINS includes your Vercel URL
2. Check NEXT_PUBLIC_API_URL ends with /api
3. Ensure Railway backend is running (check logs)
4. Test backend /health endpoint directly
```

### Environment variables not updating
```
Railway:
- Click "..." menu > "Restart"
- Check Deployments tab for new deployment

Vercel:
- Deployments tab > Click latest > "Redeploy"
- Verify variables are in correct environment (Production)
```

---

## 📚 Reference

**Database Connection Formats:**
```
Development (Direct):
postgresql://postgres.ocopuqketaddqhjwsdjd:[password]@aws-0-us-east-1.pooler.supabase.com:5432/postgres

Production (Pooler - RECOMMENDED):
postgresql://postgres.ocopuqketaddqhjwsdjd:[password]@aws-0-us-east-1.pooler.supabase.com:6543/postgres
```

**CORS Origins:**
```
Development:
http://localhost:3000

Production:
https://lumen-frontend-theta.vercel.app
https://lumen-frontend-git-main-renatodaps-projects.vercel.app
```

---

**Status**: All environment files configured! Just need to:
1. Get your database password from Supabase
2. Add it to Railway and local .env files
3. Deploy!
