# Hollow Wilds Backend - Git Workflow

## Branch Strategy

```
main (production)
  ├── Protected branch
  ├── Auto-deploy to Fly.io production
  └── Only merge via Pull Request

develop (staging)
  ├── Development branch
  ├── Merge features here first
  └── Test before merging to main

feature/* (new features)
  └── Create from develop, merge back to develop

hotfix/* (urgent fixes)
  └── Create from main, merge to both main & develop
```

## Workflow for New Features

### 1. Create Feature Branch
```bash
# Start from develop
git checkout develop
git pull origin develop

# Create feature branch
git checkout -b feature/your-feature-name
```

### 2. Develop & Commit
```bash
# Make changes
git add .
git commit -m "feat: add your feature"

# Push to GitHub
git push origin feature/your-feature-name
```

### 3. Create Pull Request
- Go to GitHub
- Create PR: `feature/your-feature-name` → `develop`
- Review & merge

### 4. Test on Develop
```bash
# Switch to develop
git checkout develop
git pull origin develop

# Test locally
go run cmd/server/main.go
```

### 5. Deploy to Production
When ready for production:
```bash
# Create PR: develop → main on GitHub
# After merge, Fly.io auto-deploys to production
```

## Workflow for Hotfixes

### 1. Create Hotfix Branch
```bash
# Start from main (production)
git checkout main
git pull origin main

# Create hotfix branch
git checkout -b hotfix/fix-description
```

### 2. Fix & Test
```bash
# Make changes
git add .
git commit -m "fix: urgent bug fix"

# Push to GitHub
git push origin hotfix/fix-description
```

### 3. Merge to Main
- Create PR: `hotfix/fix-description` → `main`
- Review & merge
- Fly.io auto-deploys

### 4. Backport to Develop
```bash
# Also merge to develop
git checkout develop
git merge hotfix/fix-description
git push origin develop
```

## Current Branches

- **main**: Production code (deploy to Fly.io)
- **develop**: Staging/development code
- **feature/\***: Feature branches
- **hotfix/\***: Hotfix branches

## Commit Message Convention

Use conventional commits:

```
feat: add new feature
fix: bug fix
docs: documentation changes
style: formatting, no code change
refactor: code restructuring
test: add tests
chore: build, dependencies, etc.
```

Examples:
```bash
git commit -m "feat: add analytics endpoint"
git commit -m "fix: resolve database connection timeout"
git commit -m "docs: update deployment guide"
```

## Deployment Targets

| Branch  | Deploy To | URL |
|---------|-----------|-----|
| main    | Production | https://gamefeel-backend.fly.dev |
| develop | Local/Manual | http://localhost:8080 |

## Protected Branch Rules (Setup on GitHub)

For `main` branch:
- ✅ Require pull request reviews (1 reviewer)
- ✅ Require status checks to pass
- ✅ Require branches to be up to date
- ✅ No direct pushes

## Quick Commands

```bash
# Switch to develop
git checkout develop

# Create feature
git checkout -b feature/my-feature

# Update from remote
git pull origin develop

# Push changes
git push origin feature/my-feature

# Back to main
git checkout main
```

---

**Happy coding! 🚀**
