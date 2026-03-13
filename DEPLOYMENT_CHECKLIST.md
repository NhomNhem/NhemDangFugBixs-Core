# Production Deployment Checklist

## 🚀 Security Configuration (Manual Steps)

### Your Situation:
- ✅ PLAYFAB_TITLE_ID: `186CA8`
- ✅ Platform: WebGL + PC/Mobile (both)
- ⚠️ WebGL domain: Not decided yet (use `*` for now)

---

## Step 1: Set Fly.io Secrets (Required)

### Minimal Setup (Recommended Now):

```bash
# Required: Enable PlayFab token validation
fly secrets set PLAYFAB_TITLE_ID=186CA8

# For now: Allow all origins (safe for PC/Mobile, convenient for testing)
fly secrets set ALLOWED_ORIGINS="*"
```

**Explanation:**
- `PLAYFAB_TITLE_ID=186CA8`: Your PlayFab title ID (enables token validation)
- `ALLOWED_ORIGINS="*"`: Allow all domains (safe because you're building PC/Mobile native apps)
  - Native apps (Windows, Android, iOS) don't have CORS restrictions
  - WebGL will work from any domain for testing
  - Update later when you have a specific WebGL hosting domain

---

## Step 2: (Optional) Restrict WebGL Origins Later

When you deploy WebGL to a specific platform, update:

```bash
# Itch.io example
fly secrets set ALLOWED_ORIGINS="https://yourusername.itch.io"

# Own domain example
fly secrets set ALLOWED_ORIGINS="https://yourgame.com,https://www.yourgame.com"

# Multiple platforms
fly secrets set ALLOWED_ORIGINS="https://yourgame.itch.io,https://poki.com,https://yourgame.com"
```

---

## Step 3: Set Fly.io Secrets (Commands)

**Copy and run these commands:**

```bash
# Navigate to backend folder
cd I:\unityVers\GameFeel-Backend

# Set PlayFab Title ID
fly secrets set PLAYFAB_TITLE_ID=186CA8

# Allow all origins (for now)
fly secrets set ALLOWED_ORIGINS="*"
```

**What happens next:**
1. Fly.io will automatically redeploy (~2 minutes)
2. Backend restarts with new secrets
3. PlayFab validation is now active

---

## Step 4: Verify Deployment

```bash
# Check deployment status
fly releases

# View logs
fly logs

# Test health endpoint
curl https://gamefeel-backend.fly.dev/health
```

**Expected response:**
```json
{
  "status": "ok",
  "message": "Hollow Wilds Backend is running",
  "version": "1.0.0",
  "database": "connected"
}
```

---

## Step 5: Test Security Features

### Test 1: PlayFab Token Validation (Should Fail)

```bash
curl -X POST https://gamefeel-backend.fly.dev/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-PlayFab-SessionToken: FAKE_TOKEN" \
  -d '{"playfabId": "TEST123"}'
```

**Expected (401 Unauthorized):**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid PlayFab session token"
  }
}
```

✅ **Good!** This means token validation is working.

---

### Test 2: Rate Limiting

```bash
# Send 101 requests (Windows PowerShell)
1..101 | ForEach-Object { Invoke-WebRequest https://gamefeel-backend.fly.dev/health }
```

**101st request should return (429):**
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests. Please try again later."
  }
}
```

✅ **Good!** Rate limiting is working.

---

## 📊 What Changed

### Code Changes:
- ✅ PlayFab token validation (auth_service.go)
- ✅ Rate limiting middleware (main.go)
- ✅ CORS configuration (main.go)
- ✅ Security documentation (docs/SECURITY.md)

### Deployed to GitHub:
- Commit: `feat: Add security hardening`
- Commit: `docs: Add security guide`
- Branch: `main` (auto-deployed by GitHub Actions)

---

## ⚠️ Important Notes

1. **ALLOWED_ORIGINS="*"**:
   - ✅ Safe for PC/Mobile native apps (no CORS)
   - ✅ Convenient for WebGL testing
   - ⚠️ Update when you have a specific WebGL domain
   - 💡 PC/Mobile builds don't need CORS restrictions

2. **PlayFab Validation**: 
   - ✅ Enabled (with PLAYFAB_TITLE_ID=186CA8)
   - Only valid PlayFab session tokens can login
   - Invalid tokens get 401 Unauthorized

3. **Rate Limiting**:
   - ✅ Active (100 requests/minute per IP)
   - Applies to ALL endpoints
   - Prevents brute force attacks

---

## 🎯 Quick Start Commands

**Just copy-paste this into PowerShell:**

```powershell
# Navigate to backend
cd I:\unityVers\GameFeel-Backend

# Set secrets
fly secrets set PLAYFAB_TITLE_ID=186CA8
fly secrets set ALLOWED_ORIGINS="*"

# Wait 2-3 minutes for redeploy

# Test (should return OK)
curl https://gamefeel-backend.fly.dev/health

# Test token validation (should return 401)
curl -X POST https://gamefeel-backend.fly.dev/api/v1/auth/login -H "Content-Type: application/json" -H "X-PlayFab-SessionToken: FAKE" -d '{\"playfabId\":\"TEST\"}'
```

---

## 🔍 Monitoring

After deployment, monitor:
```bash
# Watch logs in real-time
fly logs

# Check for errors
fly logs --level error

# View all secrets (names only, not values)
fly secrets list
```

Look for:
- `401 Unauthorized` - Failed PlayFab validation ✅ Expected for fake tokens
- `429 Too Many Requests` - Rate limit hits ✅ Expected when exceeded
- `200 OK` - Successful requests

---

## 📝 Unity Client Testing

Sau khi deploy xong, test trong Unity:

1. Open `BackendDirectTest.cs` (Assets/Scripts/Runtime/Services/Tests/)
2. Set `Environment = Production`
3. Run `Test 2: Mock Login`
   - ❌ Should fail (401) because mock token is invalid
   - ✅ This means validation is working!
4. To test real login:
   - Use PlayFab login in game first
   - Get real session token
   - Pass to backend login

---

## ✅ Completion Checklist

- [ ] Run: `fly secrets set PLAYFAB_TITLE_ID=186CA8`
- [ ] Run: `fly secrets set ALLOWED_ORIGINS="*"`
- [ ] Wait 2-3 minutes for auto-redeploy
- [ ] Test: `curl https://gamefeel-backend.fly.dev/health`
- [ ] Test: Token validation (should return 401 for fake token)
- [ ] Test in Unity: BackendDirectTest (Production mode)
- [ ] ✅ Done!

---

## 🆘 Troubleshooting

### Error: "Invalid PlayFab session token"
- ✅ **Expected!** Token validation is working correctly
- Use real PlayFab token from Unity login

### Error: "Too many requests"
- ✅ **Expected!** Rate limiting is working
- Wait 60 seconds and try again

### Error: Can't connect
- Check Fly.io deployment: `fly status`
- Check logs: `fly logs`
- Verify GitHub Actions completed: https://github.com/NhomNhem/HollowWilds-Backend/actions

---

## 🎯 When to Update ALLOWED_ORIGINS

Update when you publish WebGL to a specific platform:

**Itch.io:**
```bash
fly secrets set ALLOWED_ORIGINS="https://yourusername.itch.io"
```

**Own Domain:**
```bash
fly secrets set ALLOWED_ORIGINS="https://yourgame.com,https://www.yourgame.com"
```

**Multiple Platforms:**
```bash
fly secrets set ALLOWED_ORIGINS="https://yourgame.itch.io,https://poki.com,https://newgrounds.com"
```

**PC/Mobile only (no WebGL):**
- Keep `ALLOWED_ORIGINS="*"` (doesn't matter for native apps)

---

**Status**: ✅ Ready for deployment! Just run the commands above.

