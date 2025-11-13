# LUMEN Setup Checklist

**Complete these steps in order to get LUMEN fully operational**

---

## Current Status

✅ Frontend deployed to Vercel
✅ Backend deployed to Railway
✅ Code pushed to GitHub
✅ Build errors fixed
⏳ Database needs setup
⏳ Environment variables need configuration

---

## Step-by-Step Setup

### Phase 1: Supabase Database (10 minutes)

#### ☐ 1.1 Create Supabase Project
```
1. Go to: https://supabase.com/dashboard
2. Click "New Project"
3. Fill in:
   - Name: lumen-production
   - Database Password: [STRONG PASSWORD - SAVE THIS!]
   - Region: us-east-1 (or closest)
4. Click "Create new project"
5. Wait ~2 minutes for provisioning
```

#### ☐ 1.2 Collect Supabase Credentials
```
Go to Project Settings > API:
[ ] Copy Project URL: https://____________.supabase.co
[ ] Copy Anon public key: eyJhbG_______________
[ ] Copy Service role key: eyJhbG_______________ (different!)

Go to Project Settings > API > JWT Settings:
[ ] Copy JWT Secret: _______________

Go to Project Settings > Database:
[ ] Get connection string format:
    postgresql://postgres:[YOUR-PASSWORD]@db.____________.supabase.co:5432/postgres
```

#### ☐ 1.3 Run Database Migration
```
1. In Supabase Dashboard, go to SQL Editor
2. On your computer, open:
   C:\Users\pradord\Documents\Projects\LUMEN\supabase\migrations\20251113_initial_schema.sql
3. Copy entire file contents (Ctrl+A, Ctrl+C)
4. Paste into Supabase SQL Editor
5. Click "Run" (bottom right corner)
6. Wait for success message
7. Verify: Go to Table Editor, should see 9 tables:
   - users, areas, goals, habits, tasks
   - acceptance_criteria, daily_logs, habit_logs, provider_connections
```

#### ☐ 1.4 Enable Supabase Authentication
```
1. In Supabase Dashboard, go to Authentication > Providers
2. Verify "Email" is enabled (should be by default)
3. Optional: Enable Google OAuth, Microsoft OAuth
4. Go to Authentication > Email Templates
5. Customize confirmation email (optional)
```

---

### Phase 2: Railway Backend Configuration (5 minutes)

#### ☐ 2.1 Get Railway Backend URL
```
1. Go to: https://railway.app/dashboard
2. Click on your "lumen_backend" project
3. Click on the service
4. Go to "Settings" tab
5. Scroll to "Networking" section
6. Click "Generate Domain" if not already generated
7. Copy URL: https://lumen-backend-production-________.up.railway.app
   Save this: _______________________________________________
```

#### ☐ 2.2 Generate JWT Secret
```
Option 1 - Use online generator:
1. Go to: https://generate-secret.now.sh/32
2. Copy the generated secret

Option 2 - Use command line (if you have OpenSSL):
1. Run: openssl rand -base64 32
2. Copy output

Your JWT Secret: _______________________________________________
```

#### ☐ 2.3 Configure Railway Environment Variables
```
In Railway Dashboard > lumen_backend > Variables tab:

Add each variable (click "Add Variable" button for each):

[ ] GIN_MODE = release

[ ] DATABASE_URL = postgresql://postgres:[PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres
    (Replace [PASSWORD] and [PROJECT-REF] with your values)

[ ] JWT_SECRET = [Your generated secret from 2.2]

[ ] JWT_EXPIRY = 24h

[ ] REFRESH_TOKEN_EXPIRY = 168h

[ ] SUPABASE_URL = https://[PROJECT-REF].supabase.co

[ ] SUPABASE_SERVICE_KEY = [Service role key from step 1.2]

[ ] SUPABASE_JWT_SECRET = [JWT Secret from step 1.2]

[ ] CORS_ALLOWED_ORIGINS = https://lumen-frontend-theta.vercel.app,https://lumen-frontend-git-main-renatodaps-projects.vercel.app

[ ] LOG_LEVEL = info

[ ] RATE_LIMIT_REQUESTS = 100

[ ] RATE_LIMIT_WINDOW = 60s

[ ] RATE_LIMIT_ENABLED = true

After adding all variables:
- Railway will automatically redeploy
- Wait ~2 minutes for deployment
- Check logs for "Starting server" message
```

---

### Phase 3: Vercel Frontend Configuration (5 minutes)

#### ☐ 3.1 Configure Vercel Environment Variables
```
1. Go to: https://vercel.com/dashboard
2. Click on "lumen-frontend" project
3. Go to Settings > Environment Variables
4. Add these variables:

[ ] NEXT_PUBLIC_SUPABASE_URL
    Value: https://[PROJECT-REF].supabase.co
    Apply to: ✓ Production ✓ Preview ✓ Development

[ ] NEXT_PUBLIC_SUPABASE_ANON_KEY
    Value: [Anon key from step 1.2]
    Apply to: ✓ Production ✓ Preview ✓ Development

[ ] NEXT_PUBLIC_API_URL
    Value: [Railway backend URL from step 2.1]/api
    Example: https://lumen-backend-production-xxxx.up.railway.app/api
    Apply to: ✓ Production ✓ Preview ✓ Development

[ ] NODE_ENV
    Value: production
    Apply to: ✓ Production only
```

#### ☐ 3.2 Redeploy Frontend
```
1. In Vercel dashboard, go to "Deployments" tab
2. Click on the latest deployment
3. Click "..." menu (three dots)
4. Click "Redeploy"
5. Wait ~2 minutes for deployment
6. Click "Visit" to open the site
```

---

### Phase 4: Verification (5 minutes)

#### ☐ 4.1 Test Backend Health
```
1. Open browser
2. Go to: [Railway backend URL]/health
   Example: https://lumen-backend-production-xxxx.up.railway.app/health
3. Should see:
   {
     "status": "healthy",
     "timestamp": "2025-11-13T..."
   }

[ ] Backend health check passed
```

#### ☐ 4.2 Test Frontend Loads
```
1. Go to: https://lumen-frontend-theta.vercel.app
2. Should see LUMEN homepage with:
   - Petroleum blue background (#0A0E1A)
   - Golden accent colors
   - No errors in browser console (F12)

[ ] Frontend loads successfully
```

#### ☐ 4.3 Test Authentication
```
1. On frontend, click "Sign Up" or "Get Started"
2. Enter email and password
3. Should receive confirmation email
4. Check spam folder if not received
5. Click confirmation link in email
6. Should be able to log in

[ ] User registration works
[ ] Email confirmation received
[ ] Login successful
```

#### ☐ 4.4 Test Database Connection
```
1. After logging in, try to create a habit
2. Check Supabase Dashboard > Table Editor > habits table
3. Should see new row with your habit

[ ] Data saves to database successfully
```

---

## Troubleshooting

### ❌ Backend Health Check Fails
```
Check Railway logs:
1. Railway Dashboard > lumen_backend > Logs tab
2. Look for errors related to:
   - DATABASE_URL connection
   - Missing environment variables
   - Port binding issues

Common fixes:
- Verify DATABASE_URL password is correct
- Check all env variables are set
- Ensure Supabase project is active
```

### ❌ Frontend Shows CORS Error
```
Check browser console (F12):
1. Look for "CORS policy" error
2. Verify CORS_ALLOWED_ORIGINS includes exact Vercel URL
3. Redeploy backend after changing CORS settings

Common fixes:
- Add all Vercel URLs to CORS_ALLOWED_ORIGINS
- Include protocol (https://) in URLs
- No trailing slashes in URLs
```

### ❌ Authentication Not Working
```
Check:
1. Supabase Dashboard > Authentication > Providers
2. Verify Email provider is enabled
3. Check email templates are configured
4. Verify SUPABASE_URL and keys match between frontend/backend

Common fixes:
- Enable email confirmation in Supabase settings
- Check spam folder for confirmation email
- Verify RLS policies are enabled
```

### ❌ Database Connection Fails
```
Test connection:
1. Supabase Dashboard > SQL Editor
2. Run: SELECT NOW();
3. Should return current timestamp

Common fixes:
- Check Supabase project status (not paused)
- Verify password has no special chars needing URL encoding
- Ensure migration ran successfully
```

---

## Post-Setup Tasks

After everything works:

#### ☐ Security Checklist
```
[ ] Change all default passwords
[ ] Verify RLS policies are enabled
[ ] Check no secrets in git history
[ ] Enable 2FA on Supabase, Railway, Vercel accounts
[ ] Review CORS settings
[ ] Enable rate limiting
```

#### ☐ Monitoring Setup
```
[ ] Set up error tracking (Sentry - planned)
[ ] Enable Vercel Analytics
[ ] Monitor Railway logs
[ ] Check Supabase database usage
```

#### ☐ Documentation
```
[ ] Save all credentials in password manager
[ ] Document custom configuration
[ ] Update team on deployment status
```

---

## Quick Reference

### Important URLs
```
Frontend: https://lumen-frontend-theta.vercel.app
Backend: https://[your-railway-url]
Supabase: https://[your-project].supabase.co

Dashboards:
- Vercel: https://vercel.com/dashboard
- Railway: https://railway.app/dashboard
- Supabase: https://supabase.com/dashboard
```

### Required Credentials (Save These!)
```
Supabase:
- Project URL: _______________
- Anon Key: _______________
- Service Key: _______________
- JWT Secret: _______________
- Database Password: _______________

Railway:
- Backend URL: _______________

Generated:
- JWT Secret: _______________
```

---

## Completion Checklist

Mark when complete:

### Database Setup
- [ ] Supabase project created
- [ ] Credentials collected and saved
- [ ] Migration ran successfully
- [ ] 9 tables created
- [ ] RLS policies enabled
- [ ] Auth provider enabled

### Backend Configuration
- [ ] Railway URL generated
- [ ] JWT secret generated
- [ ] All 13 environment variables set
- [ ] Backend redeployed
- [ ] Health check passing
- [ ] No errors in logs

### Frontend Configuration
- [ ] All 4 environment variables set
- [ ] Frontend redeployed
- [ ] Site loads without errors
- [ ] API calls connect to backend

### Verification
- [ ] Backend health check passes
- [ ] Frontend loads successfully
- [ ] User registration works
- [ ] Email confirmation received
- [ ] Login successful
- [ ] Data saves to database

---

## Time Estimate

- Phase 1 (Database): 10 minutes
- Phase 2 (Backend): 5 minutes
- Phase 3 (Frontend): 5 minutes
- Phase 4 (Verification): 5 minutes

**Total: ~25 minutes** for complete setup

---

**Status**: Ready to configure! Follow this checklist step by step.

**Need Help?** See `/docs/DEPLOYMENT.md` for detailed explanations.
