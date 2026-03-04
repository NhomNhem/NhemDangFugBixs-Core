# GameFeel Backend Deployment Guide

## 🚀 Deploy to Fly.io

### Prerequisites
1. Install Fly.io CLI: https://fly.io/docs/hands-on/install-flyctl/
2. Create Fly.io account (free): https://fly.io/app/sign-up

### Step 1: Install Fly CLI

**Windows (PowerShell):**
```powershell
iwr https://fly.io/install.ps1 -useb | iex
```

**Verify installation:**
```bash
flyctl version
```

### Step 2: Login to Fly.io

```bash
flyctl auth login
```

This will open a browser for authentication.

### Step 3: Create Fly.io App

```bash
cd I:\unityVers\GameFeel-Backend
flyctl launch --no-deploy
```

When asked:
- Choose app name: `gamefeel-backend` (or your preferred name)
- Choose region: **Singapore (sin)** - closest to Vietnam
- Would you like to set up a PostgreSQL database? **No** (we use Supabase)
- Would you like to set up an Upstash Redis database? **No** (optional later)

### Step 4: Set Environment Secrets

**IMPORTANT:** Never commit secrets to Git!

```bash
# Set database URL (from Supabase)
flyctl secrets set DATABASE_URL="postgresql://postgres.xxx:password@xxx.supabase.co:5432/postgres"

# Set JWT secret (generate a strong random string)
flyctl secrets set JWT_SECRET="your-super-secret-jwt-key-change-this-in-production"
```

**Generate a strong JWT secret:**
```powershell
# PowerShell
[System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes((New-Guid).ToString() + (New-Guid).ToString()))
```

### Step 5: Deploy!

```bash
flyctl deploy
```

This will:
1. Build Docker image
2. Push to Fly.io registry
3. Deploy to production
4. Start the server

### Step 6: Verify Deployment

```bash
# Check status
flyctl status

# View logs
flyctl logs

# Open in browser
flyctl open
```

Test your deployed API:
```bash
curl https://gamefeel-backend.fly.dev/health
```

Expected response:
```json
{
  "status": "ok",
  "message": "GameFeel Backend is running",
  "database": "connected",
  "version": "1.0.0"
}
```

### Step 7: Update Unity Client

In Unity, update `BackendClient.cs`:

```csharp
// Change from localhost to production URL
private const string BACKEND_URL = "https://gamefeel-backend.fly.dev";
```

Or use environment-based URL:
```csharp
#if UNITY_EDITOR
    private const string BACKEND_URL = "http://localhost:8080";
#else
    private const string BACKEND_URL = "https://gamefeel-backend.fly.dev";
#endif
```

---

## 📊 Monitoring

### View Logs (Real-time)
```bash
flyctl logs -a gamefeel-backend
```

### Check App Status
```bash
flyctl status -a gamefeel-backend
```

### View Metrics
```bash
flyctl dashboard -a gamefeel-backend
```

Opens Fly.io dashboard in browser with:
- CPU usage
- Memory usage
- Request metrics
- Response times

---

## 🔧 Troubleshooting

### Database Connection Issues

If you see "database unreachable" errors:

1. **Check Supabase IP allowlist:**
   - Supabase → Settings → Database → Connection pooling
   - Make sure "Restrict access to trusted IPs" is OFF for free tier
   - Or add Fly.io IP ranges

2. **Test connection from Fly.io:**
   ```bash
   flyctl ssh console
   # Inside the container:
   apk add postgresql-client
   psql "$DATABASE_URL"
   ```

### Port Issues

If health check fails:
- Ensure `PORT=8080` in fly.toml matches your Go server
- Check server logs: `flyctl logs`

### Memory Issues

Free tier has 256MB RAM. If OOM (Out of Memory):
```toml
# Edit fly.toml
[[vm]]
  memory = '512mb'  # Upgrade to 512MB (still free tier)
```

Then redeploy:
```bash
flyctl deploy
```

---

## 🔄 CI/CD Setup (Optional)

### GitHub Actions Auto-Deploy

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to Fly.io

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - run: flyctl deploy --remote-only
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

Get your Fly.io token:
```bash
flyctl auth token
```

Add to GitHub:
- Repository → Settings → Secrets → New secret
- Name: `FLY_API_TOKEN`
- Value: (paste token)

Now every push to `main` auto-deploys! 🚀

---

## 💰 Cost Estimate

**Fly.io Free Tier:**
- Up to 3 shared-cpu VMs
- 256MB RAM per VM
- 160GB bandwidth/month
- **Cost: $0/month** ✅

**If you exceed free tier:**
- Shared CPU VM: ~$2/month
- 256MB RAM: included
- Bandwidth: $0.02/GB after 160GB

**For your game (<10k DAU):** Should stay FREE! 🎉

---

## 🌐 Custom Domain (Optional)

Want a custom domain like `api.yourgame.com`?

1. **Add certificate:**
   ```bash
   flyctl certs add api.yourgame.com
   ```

2. **Add DNS records** (in your domain provider):
   - Type: `CNAME`
   - Name: `api`
   - Value: `gamefeel-backend.fly.dev`

3. **Wait for certificate** (5-10 minutes)

4. **Update Unity:**
   ```csharp
   private const string BACKEND_URL = "https://api.yourgame.com";
   ```

---

## 🔐 Security Checklist

Before going live:

- [ ] Change JWT_SECRET to a strong random value
- [ ] Enable HTTPS only (force_https = true in fly.toml) ✅
- [ ] Review CORS settings in Go backend
- [ ] Enable rate limiting (TODO: add Redis)
- [ ] Monitor logs for suspicious activity
- [ ] Backup database regularly (Supabase auto-backups enabled)
- [ ] Don't log sensitive data (passwords, tokens)

---

## 📚 Useful Commands

```bash
# Restart app
flyctl apps restart

# Scale machines
flyctl scale count 2  # Run 2 instances

# SSH into container
flyctl ssh console

# View environment variables
flyctl secrets list

# Change secrets
flyctl secrets set KEY=VALUE

# Destroy app (careful!)
flyctl apps destroy gamefeel-backend
```

---

## 🆘 Getting Help

- Fly.io Docs: https://fly.io/docs/
- Fly.io Community: https://community.fly.io/
- Check logs: `flyctl logs`
- GitHub Issues: Create issue in your repo

---

**Happy Deploying! 🚀**
