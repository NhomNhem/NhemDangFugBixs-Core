# Emergency Deployment Fix Guide

## 🚨 Problem: Backend Crashing on Fly.io

**Symptom**: Machines restarting frequently (503 errors)  
**Machine ID**: 3287d013c99785  
**Cause**: Application crash on startup

---

## 🔧 Fix Without fly CLI

### Option 1: Fly.io Web Dashboard (Easiest)

**Step 1: Access Dashboard**
1. Go to: https://fly.io/dashboard
2. Click on your app: `gamefeel-backend`
3. Go to **Secrets** section (left sidebar)

**Step 2: Check Required Secrets**
Verify these secrets exist:
- ✅ `DATABASE_URL` (Supabase connection string)
- ✅ `JWT_SECRET` (your JWT secret key)
- ✅ `PLAYFAB_TITLE_ID` (186CA8)
- ✅ `ALLOWED_ORIGINS` (*)

**If ANY are missing** → Backend will crash!

---

### Option 2: Check Logs (Find the Real Error)

**Via Web Dashboard:**
1. Go to: https://fly.io/dashboard
2. Click on your app: `gamefeel-backend`
3. Go to **Monitoring** → **Logs**
4. Click **"Logs from Previous Starts"** tab
5. Look for error messages before crash

**Common errors to look for:**
```
❌ DATABASE_URL not set
❌ failed to connect to database
❌ panic: runtime error
❌ missing JWT_SECRET
```

---

## 🎯 Most Likely Issue

When you ran:
```bash
fly secrets set PLAYFAB_TITLE_ID=186CA8
fly secrets set ALLOWED_ORIGINS="*"
```

**Fly.io might have lost other secrets!** (This is a known issue)

Secrets that might be missing:
- `DATABASE_URL`
- `JWT_SECRET`

---

## ✅ Fix: Re-add All Secrets

If you have access to fly CLI (install or use another machine):

```bash
# Re-add all secrets
fly secrets set DATABASE_URL="postgresql://postgres:[PASSWORD]@db.vxbcjtnznelovdeevxur.supabase.co:5432/postgres"

fly secrets set JWT_SECRET="your-jwt-secret-here"

fly secrets set PLAYFAB_TITLE_ID=186CA8

fly secrets set ALLOWED_ORIGINS="*"
```

**Important**: Replace `[PASSWORD]` with your actual Supabase password

---

## 🔄 Alternative: Rollback

If you can't fix secrets, rollback to previous working version:

```bash
# See releases
fly releases

# Rollback to last working version (before security changes)
fly releases rollback <version-before-security>
```

This will undo the security changes and backend will work like before.

---

## 🆘 Quick Workaround: Disable PlayFab Validation

If the issue is PlayFab validation code:

```bash
# Remove PlayFab Title ID
fly secrets unset PLAYFAB_TITLE_ID
```

This makes validation skip (like development mode).

---

## 📋 What You Need From configs/.env

Check your local `configs/.env` file for these values:

```bash
# Backend folder
cd I:\unityVers\GameFeel-Backend

# View .env file
cat configs/.env
```

You need:
- `SUPABASE_DATABASE_URL` (or `DATABASE_URL`)
- `JWT_SECRET`

Then set them on Fly.io.

---

## 🔍 Debugging Steps

### 1. Check Logs (Web Dashboard)
Look for the actual error message before crash.

### 2. Verify Secrets Exist
Make sure DATABASE_URL and JWT_SECRET are still there.

### 3. Test Locally
```bash
# In backend folder
go run cmd/server/main.go
```

If it works locally but not on Fly.io → secrets issue!

### 4. Check GitHub Actions
https://github.com/NhomNhem/HollowWilds-Backend/actions

See if deployment succeeded or failed with errors.

---

## 💡 Temporary Solution

While debugging, you can:

1. **Rollback** to previous working version
2. Backend works normally (without security features)
3. Debug the issue
4. Re-deploy security features later

```bash
fly releases
fly releases rollback <previous-version>
```

---

## 🎯 Next Steps

1. **Check Fly.io Dashboard** → Logs → Find error message
2. **Check Secrets** → Verify DATABASE_URL exists
3. **Re-add secrets** if missing
4. **Or rollback** if can't fix immediately

Then let me know what error you found in logs!

---

## 📞 Need Help?

Share the error message from Fly.io logs and we can fix it together.

Error format example:
```
panic: missing environment variable: DATABASE_URL
```

Or:
```
failed to connect to database: connection refused
```
