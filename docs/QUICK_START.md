# LUMEN Quick Start Guide

**Get LUMEN running in 10 minutes**

---

## ✅ Current Status

- **Frontend**: ✅ Deployed on Vercel (https://lumen-frontend-theta.vercel.app)
- **Backend**: ✅ Deployed on Railway (Go 1.23, health check passing)
- **Database**: ⏳ Needs Supabase setup

---

## 🚀 Next Steps (In Order)

### 1. Set Up Supabase (5 minutes)

```bash
# Go to https://supabase.com/dashboard
1. Click "New Project"
2. Name: "lumen-production"
3. Set strong database password (SAVE THIS!)
4. Region: us-east-1 (or closest to you)
5. Wait ~2 minutes for project creation
```

**Get your credentials:**
```
Project Settings > API:
  - Project URL: https://xxxxx.supabase.co
  - Anon public key: eyJhbG...
  - Service role key: eyJhbG... (different key!)

Project Settings > Database:
  - Connection string: postgresql://postgres:[password]@db.xxxxx.supabase.co:5432/postgres
```

### 2. Run Database Migration (1 minute)

```bash
# In Supabase Dashboard:
1. Go to SQL Editor
2. Open: LUMEN/supabase/migrations/20251113_initial_schema.sql
3. Copy entire contents
4. Paste into SQL Editor
5. Click "Run" (bottom right)
6. Verify: Check "Table Editor" - should see 9 tables
```

### 3. Configure Backend Environment Variables (2 minutes)

```bash
# In Railway Dashboard (https://railway.app):
1. Go to your lumen_backend service
2. Click "Variables" tab
3. Add these variables:

GIN_MODE=release
DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres
JWT_SECRET=[Generate with: openssl rand -base64 32]
SUPABASE_URL=https://[PROJECT-REF].supabase.co
SUPABASE_SERVICE_KEY=[Service role key from Supabase]
SUPABASE_JWT_SECRET=[From Supabase Settings > API > JWT Secret]
CORS_ALLOWED_ORIGINS=https://lumen-frontend-theta.vercel.app

4. Click anywhere to save
5. Railway auto-redeploys with new variables
```

### 4. Configure Frontend Environment Variables (2 minutes)

```bash
# In Vercel Dashboard (https://vercel.com):
1. Go to lumen-frontend project
2. Settings > Environment Variables
3. Add these variables:

NEXT_PUBLIC_SUPABASE_URL=https://[PROJECT-REF].supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=[Anon key from Supabase]
NEXT_PUBLIC_API_URL=https://[RAILWAY-BACKEND-URL]/api

4. Apply to: Production, Preview, Development
5. Deployments > Latest > Redeploy
```

### 5. Get Railway Backend URL

```bash
# In Railway Dashboard:
1. Go to your lumen_backend service
2. Click "Settings" tab
3. Scroll to "Networking" section
4. Click "Generate Domain"
5. Copy the URL: https://lumen-backend-production-xxxx.up.railway.app
6. Use this URL in frontend NEXT_PUBLIC_API_URL variable
```

### 6. Test Everything (Optional)

```bash
# Test backend health:
curl https://[RAILWAY-URL]/health

# Expected response:
{"status":"healthy","timestamp":"2025-11-13T..."}

# Test frontend:
1. Open https://lumen-frontend-theta.vercel.app
2. Click "Sign Up" (Supabase Auth)
3. Create account with email
4. Should receive verification email
5. After verification, you can log in
```

---

## 📋 Environment Variables Checklist

### Backend (Railway) - 7 variables
- [ ] `GIN_MODE=release`
- [ ] `DATABASE_URL` (from Supabase)
- [ ] `JWT_SECRET` (generate with OpenSSL)
- [ ] `SUPABASE_URL` (from Supabase)
- [ ] `SUPABASE_SERVICE_KEY` (from Supabase)
- [ ] `SUPABASE_JWT_SECRET` (from Supabase Settings > API)
- [ ] `CORS_ALLOWED_ORIGINS` (Vercel frontend URL)

### Frontend (Vercel) - 3 variables
- [ ] `NEXT_PUBLIC_SUPABASE_URL` (from Supabase)
- [ ] `NEXT_PUBLIC_SUPABASE_ANON_KEY` (from Supabase)
- [ ] `NEXT_PUBLIC_API_URL` (Railway backend URL + /api)

---

## 🔧 Useful Commands

```bash
# Generate JWT secret
openssl rand -base64 32

# Test backend health
curl https://[RAILWAY-URL]/health

# View backend logs
# Go to Railway dashboard > service > Logs tab

# View frontend logs
# Go to Vercel dashboard > project > Deployment > Runtime Logs
```

---

## 🐛 Troubleshooting

### Backend won't start
```
1. Check Railway logs for errors
2. Verify DATABASE_URL is correct
3. Test database connection from Supabase dashboard
4. Check all environment variables are set
```

### Frontend can't connect to backend
```
1. Verify NEXT_PUBLIC_API_URL includes /api at the end
2. Check CORS_ALLOWED_ORIGINS in backend includes exact Vercel URL
3. Test backend /health endpoint directly
4. Check browser console for CORS errors
```

### Database connection errors
```
1. Verify DATABASE_URL password is correct
2. Check Supabase project is "Active" (not paused)
3. Run migration again if tables are missing
4. Check RLS policies are enabled
```

### Authentication not working
```
1. Verify Supabase Auth is enabled (Authentication > Providers)
2. Check email templates are configured
3. Verify SUPABASE_URL and keys match between frontend and backend
4. Check user is verified (Supabase Authentication > Users)
```

---

## 📚 Next Steps After Setup

1. **Create your first habit**: Test the core functionality
2. **Set up night planning**: Define acceptance criteria
3. **Track "Did I Win My Day?"**: Complete criteria and mark day as won
4. **Explore stats**: View streaks and completion patterns
5. **Customize**: Add areas, goals, and tasks

---

## 🎯 What You Built

✅ **Production-ready architecture**
- Frontend: Next.js 15 + React 19 + TypeScript
- Backend: Go 1.23 + Gin + PostgreSQL
- Database: Supabase (managed PostgreSQL)
- Hosting: Vercel (frontend) + Railway (backend)

✅ **Features**
- User authentication (Supabase Auth)
- Habit tracking with streaks
- Task management with horizons
- Night planning workflow
- "Did I Win My Day?" system
- Offline-first PWA (service workers)
- Mobile-responsive design

✅ **Security**
- Row Level Security (RLS) enabled
- JWT authentication
- Environment variables secured
- CORS configured
- HTTPS enforced

---

**Status**: Almost ready! Just need to configure environment variables and you're live!

**Questions?** Check `/docs/DEPLOYMENT.md` for detailed step-by-step guide.
