# 🔧 Seed File Fix - Role Constraint Error

**Error**: `new row for relation "users" violates check constraint "users_role_check"`

**Problem**: Your seed file uses `role='educator'` but the constraint only allows: `user`, `admin`, `premium`

---

## ⚡ Quick Fix - Option 1 (Recommended - 30 seconds)

### Expand the Role Constraint

**Supabase Dashboard → SQL Editor**

1. **Copy this SQL**:
   ```sql
   -- Expand role constraint to allow more user types
   ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;

   ALTER TABLE users ADD CONSTRAINT users_role_check
     CHECK (role IN ('user', 'admin', 'premium', 'educator', 'instructor', 'student', 'moderator', 'contributor'));
   ```

2. **Paste** into SQL Editor

3. Click **"Run"**

4. ✅ Now run your seed file - it will work!

---

## 🔓 Alternative - Option 2 (Even Simpler)

### Remove Role Constraint Entirely

If you want to use ANY role values without restrictions:

**Supabase Dashboard → SQL Editor**

```sql
-- Remove role constraint entirely
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;

-- Remove subscription tier constraint (if you want flexibility)
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_subscription_tier_check;

-- Remove subscription status constraint (if you want flexibility)
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_subscription_status_check;
```

✅ This allows any text value for role, subscription_tier, and subscription_status

---

## 📋 What Your Seed File Has

From the error, your seed file is trying to insert:
```sql
role = 'educator'           ❌ Not allowed by current constraint
subscription_tier = 'premium' ✅ This is allowed
```

---

## ✅ After Running Fix

**Your seed file should now work!**

1. Run one of the SQL fixes above
2. Re-run your seed.sql file
3. Should complete without errors

---

## 🎯 Recommended: Option 1

**Why**: Keeps some data validation while allowing common roles

**Allowed roles after fix**:
- ✅ user
- ✅ admin
- ✅ premium
- ✅ **educator** ← NEW
- ✅ **instructor** ← NEW
- ✅ **student** ← NEW
- ✅ moderator
- ✅ contributor

---

## 🔍 If You Get Different Constraint Errors

The constraint names are:
- `users_role_check` - Role validation
- `users_subscription_tier_check` - Subscription tier validation
- `users_subscription_status_check` - Subscription status validation

**To see what values are currently allowed**:
```sql
-- Check constraint definitions
SELECT conname, pg_get_constraintdef(oid)
FROM pg_constraint
WHERE conrelid = 'users'::regclass
AND contype = 'c';
```

**To remove ANY constraint**:
```sql
ALTER TABLE users DROP CONSTRAINT constraint_name_here;
```

---

## 📝 Next Steps

After seed file works:
1. ✅ Database fully set up
2. ⏳ Get database password from Supabase
3. ⏳ Configure Railway environment variables
4. ⏳ Configure Vercel environment variables
5. ⏳ Test authentication

**Continue with**: `NEXT_STEPS.md`

---

**Status**: Role constraint issue identified and fix provided!

**Action**: Run Option 1 SQL above → Re-run your seed file → Should work!
