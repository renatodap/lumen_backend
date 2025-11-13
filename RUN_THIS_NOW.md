# 🚨 URGENT FIX - Run This SQL Now

**Error**: `column "subscription_tier" of relation "users" does not exist`

**Solution**: Run this ONE migration that adds ALL missing columns.

---

## ⚡ Quick Fix (2 minutes)

### Step 1: Run the Complete Migration

**Go to Supabase Dashboard → SQL Editor**

1. Open this file on your computer:
   ```
   C:\Users\pradord\Documents\Projects\LUMEN\supabase\migrations\20251113_add_all_user_columns.sql
   ```

2. **Copy the ENTIRE file** (Ctrl+A, Ctrl+C)

3. **Paste** into Supabase SQL Editor

4. Click **"Run"** (bottom right)

5. ✅ Success message: **"Success. No rows returned"**

---

### Step 2: Verify Columns Exist

**Go to Supabase Dashboard → Table Editor → users**

You should now see these columns:
- ✅ id
- ✅ email
- ✅ timezone
- ✅ **full_name** ← NEW
- ✅ **role** ← NEW
- ✅ **subscription_tier** ← NEW
- ✅ **subscription_status** ← NEW
- ✅ **is_active** ← NEW
- ✅ **email_verified** ← NEW
- ✅ created_at
- ✅ updated_at

---

### Step 3: Now Run Your Seed File

**Your seed file should now work!**

**Supabase Dashboard → SQL Editor**

1. Open your seed.sql file (the one that was failing)

2. Copy entire contents

3. Paste into SQL Editor

4. Click "Run"

5. ✅ Should complete without errors

---

## 📋 What This Migration Adds

### New Columns:
```sql
full_name              TEXT                  -- User's display name
role                   TEXT DEFAULT 'user'   -- user, admin, premium
subscription_tier      TEXT DEFAULT 'free'   -- free, starter, pro, enterprise
subscription_status    TEXT DEFAULT 'active' -- active, cancelled, past_due, trialing
is_active             BOOLEAN DEFAULT true   -- Account active/disabled
email_verified        BOOLEAN DEFAULT false  -- Email confirmation status
```

### Constraints:
- Role must be: user, admin, or premium
- Subscription tier must be: free, starter, pro, or enterprise
- Subscription status must be: active, cancelled, past_due, or trialing

### Indexes:
- Fast queries by role
- Fast queries by subscription tier
- Fast queries by active status

### Automatic Sync:
- Trigger automatically creates user record when someone signs up
- Syncs email verification status from Supabase Auth

---

## ✅ Checklist

```
☐ Ran 20251113_add_all_user_columns.sql
☐ Verified all columns exist in Table Editor
☐ Re-ran your seed.sql file
☐ Verified seed data created successfully
☐ Ready to configure Railway/Vercel
```

---

## 🔧 If You Still Get Errors

### Error: "relation already exists"
```
This is NORMAL - it means migration already partially ran
Action: Safe to ignore, continue to seed file
```

### Error: "constraint already exists"
```
This is NORMAL - migration uses IF EXISTS checks
Action: Safe to ignore, continue to seed file
```

### Error: Different column missing
```
Action: Let me know which column and I'll add it
```

---

## 📚 After This Works

Continue with deployment:
1. **Next Steps**: Open `NEXT_STEPS.md`
2. **Get database password** from Supabase
3. **Configure Railway** environment variables
4. **Configure Vercel** environment variables
5. **Test** everything end-to-end

---

**Status**: Migration updated to include ALL columns your seed file needs!

**Action**: Run the SQL file above, then retry your seed file.
