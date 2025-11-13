# LUMEN Database Setup Guide

**Fix the schema mismatch and populate with sample data**

---

## 🔴 Issue You Encountered

You ran the initial migration successfully, but the seed.sql file expected columns that don't exist:
- `full_name`
- `role`
- `subscription_tier`
- `subscription_status`
- `is_active`
- `email_verified`

**Solution**: Run the additional migration to add these columns, then use the new seed file.

---

## ✅ Step-by-Step Fix

### 1️⃣ Run the Additional Migration (2 minutes)

**In Supabase Dashboard → SQL Editor:**

1. Open this file on your computer:
   ```
   C:\Users\pradord\Documents\Projects\LUMEN\supabase\migrations\20251113_add_user_columns.sql
   ```

2. Copy the ENTIRE file (Ctrl+A, Ctrl+C)

3. Paste into Supabase SQL Editor

4. Click "Run" (bottom right)

5. You should see: **Success. No rows returned**

**This migration adds:**
- `full_name` column (for display names)
- `role` column (user, admin, premium)
- `is_active` column (account status)
- `email_verified` column (sync with auth)
- Automatic trigger to sync with Supabase Auth

---

### 2️⃣ Verify the Migration (1 minute)

**In Supabase Dashboard → Table Editor:**

1. Click on `users` table

2. You should now see these columns:
   - id
   - email
   - timezone
   - **full_name** ← NEW
   - **role** ← NEW
   - **is_active** ← NEW
   - **email_verified** ← NEW
   - created_at
   - updated_at

---

### 3️⃣ Run the Seed Data (2 minutes)

**Important: You need at least one user account first!**

#### Option A: Create User via Frontend
```
1. Go to: https://lumen-frontend-theta.vercel.app
2. Click "Sign Up"
3. Create an account with email/password
4. Verify email (check inbox/spam)
5. Log in once to confirm user is created
```

#### Option B: Create User Manually in Supabase
```
1. Supabase Dashboard → Authentication → Users
2. Click "Add User"
3. Enter email and password
4. Click "Create User"
5. Mark as "Email Confirmed"
```

#### After User Exists:

**In Supabase Dashboard → SQL Editor:**

1. Open this file on your computer:
   ```
   C:\Users\pradord\Documents\Projects\LUMEN\supabase\seeds\seed.sql
   ```

2. Copy the ENTIRE file (Ctrl+A, Ctrl+C)

3. Paste into Supabase SQL Editor

4. Click "Run"

5. You should see: **Success. No rows returned**

---

### 4️⃣ Verify Seed Data (1 minute)

**In Supabase Dashboard → Table Editor:**

Check each table has data:

- **areas**: Should have 3 rows
  - Health & Fitness 💪
  - Career & Learning 📚
  - Personal Growth 🌱

- **goals**: Should have 2 rows
  - Winter Arc Transformation
  - Master Full-Stack Development

- **habits**: Should have 3 rows
  - Morning Workout 🏋️
  - Track Calories 🍎
  - Code for 2 hours 💻

- **acceptance_criteria**: Should have 5 rows
  - Complete morning workout
  - Track all meals and hit calorie target
  - Code for at least 2 hours
  - Plan tomorrow before 10 PM
  - No junk food

- **tasks**: Should have 3 rows
  - Finish LUMEN authentication flow
  - Meal prep for the week
  - Deploy LUMEN to production

- **daily_logs**: Should have 1 row (yesterday's log)

- **habit_logs**: Should have 3 rows (last 3 days of workouts)

---

## 📋 Complete Migration Order

If you need to rebuild from scratch:

```sql
-- 1. Initial schema (creates all tables)
Run: supabase/migrations/20251113_initial_schema.sql

-- 2. Add user columns (adds missing columns)
Run: supabase/migrations/20251113_add_user_columns.sql

-- 3. Seed data (creates sample data)
Run: supabase/seeds/seed.sql
```

---

## 🔧 Database Schema Reference

### Users Table (Final)
```sql
users (
  id UUID PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  timezone TEXT DEFAULT 'UTC',
  full_name TEXT,
  role TEXT DEFAULT 'user',
  is_active BOOLEAN DEFAULT true,
  email_verified BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ
)
```

### All Tables
1. **users** - User accounts (synced with auth.users)
2. **areas** - Life areas (PARA method)
3. **goals** - Projects and goals
4. **habits** - Tracked habits with reminders
5. **tasks** - Tasks with time horizons (2-day, 7-day, future)
6. **acceptance_criteria** - Daily win criteria
7. **daily_logs** - Daily reflections and wins
8. **habit_logs** - Habit completion history
9. **provider_connections** - OAuth integrations (Google, Microsoft)

---

## 🎯 What the Seed Data Creates

Sample data for testing the full LUMEN experience:

### Health & Fitness Area
- **Goal**: Winter Arc Transformation
- **Habits**:
  - Morning Workout (daily, 6 AM / 6 PM reminders)
  - Track Calories (daily, 12 PM / 8 PM reminders)
- **Tasks**:
  - Meal prep for the week (2-day horizon)

### Career & Learning Area
- **Goal**: Master Full-Stack Development
- **Habit**:
  - Code for 2 hours (daily, 9 AM / 9 PM reminders)
- **Tasks**:
  - Finish LUMEN authentication flow (2-day horizon)
  - Deploy LUMEN to production (7-day horizon)

### Daily Win Criteria
1. Complete morning workout
2. Track all meals and hit calorie target
3. Code for at least 2 hours
4. Plan tomorrow before 10 PM
5. No junk food

### Sample Data
- 1 daily log (yesterday - marked as "won")
- 3 habit logs (last 3 days of morning workouts)

---

## 🧪 Testing the Seeded Data

### Test in Frontend

Once environment variables are configured:

1. **Log in** with your account

2. **Check Today View**:
   - Should see 5 acceptance criteria
   - Should show progress bars

3. **Check Habits**:
   - Should see 3 habits
   - Each with reminder times
   - Some with completion history

4. **Check Tasks**:
   - Should see 3 tasks
   - Filtered by horizon (2-day, 7-day)

5. **Check Stats**:
   - Should show streak for morning workout (3 days)
   - Should show yesterday as a "win"

---

## 🔄 Automatic User Sync

The migration includes a trigger that automatically:
- Creates a `users` row when someone signs up via Supabase Auth
- Syncs email and verification status
- Sets default values (role: user, is_active: true)

**You don't need to manually create users!** Just sign up via the frontend and the trigger handles it.

---

## 🐛 Troubleshooting

### Error: "column does not exist"
```
Cause: Migration not run yet
Fix: Run 20251113_add_user_columns.sql first
```

### Error: "null value in column violates not-null constraint"
```
Cause: No users exist yet
Fix: Create a user via Auth first, then run seed
```

### Seed creates no data
```
Cause: No matching users found
Fix:
1. Check auth.users has at least one user
2. Check public.users has matching row
3. Verify trigger is working (sign up test user)
```

### Duplicate key errors when re-running seed
```
Cause: Seed already run
Fix: Normal - seed uses ON CONFLICT DO NOTHING
```

---

## 📝 Next Steps After Database Setup

1. ✅ Initial schema migration - **DONE**
2. ✅ Add user columns migration - **RUN NOW**
3. ✅ Seed sample data - **RUN AFTER SIGNUP**
4. ⏳ Configure Railway environment variables
5. ⏳ Configure Vercel environment variables
6. ⏳ Test authentication end-to-end

**Continue with**: `NEXT_STEPS.md` for deployment configuration

---

**Status**: Database schema fixed and ready for seeding!

**Action**: Run the two migrations above, then seed with sample data.
