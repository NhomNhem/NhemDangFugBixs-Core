# Hollow Wilds Backend Server

Go backend server for GameFeel Unity game.

## Tech Stack
- **Framework**: Fiber (Go web framework)
- **Database**: Supabase (PostgreSQL)
- **Auth**: JWT (PlayFab session token validation)
- **Payments**: Stripe
- **Cache**: Upstash Redis

## Project Structure

```
.
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── api/              # HTTP handlers & routes
│   ├── models/           # Data models
│   ├── services/         # Business logic
│   ├── database/         # Database queries
│   ├── middleware/       # Auth, logging, rate limiting
│   └── validation/       # Request validation
├── pkg/
│   └── utils/            # Shared utilities
├── configs/              # Configuration files
├── deployments/          # Docker & deployment configs
└── docs/                 # Documentation
```

## Prerequisites

- Go 1.22+
- Supabase account
- Stripe account (test mode)

## Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Configure Environment

```bash
cp configs/.env.example .env
# Edit .env with your credentials
```

### 3. Run Server

```bash
go run cmd/server/main.go
```

Server will start on `http://localhost:8080`

## Development

### Run locally
```bash
go run cmd/server/main.go
```

### Run with hot reload (air)
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run
air
```

### Run tests
```bash
go test ./...
```

### Build Verification
To ensure everything is correct before committing, run the build verification script:
```powershell
./scripts/verify_build.ps1
```

For detailed testing documentation, see [docs/TESTING.md](docs/TESTING.md).

### Build for production
```bash
go build -o bin/server cmd/server/main.go
```

## API Endpoints

See [docs/api.md](docs/api.md) for full API documentation.

### Quick Reference

- `GET /health` - Health check
- `POST /api/v1/auth/login` - PlayFab auth → JWT
- `POST /api/v1/levels/complete` - Submit level completion
- `POST /api/v1/talents/upgrade` - Upgrade talent
- `POST /api/v1/payments/create-session` - Create Stripe checkout
- `POST /api/v1/analytics/events` - Submit analytics events

## Deployment

### Docker

```bash
docker build -t gamefeel-backend .
docker run -p 8080:8080 --env-file .env gamefeel-backend
```

### Fly.io

```bash
fly launch
fly deploy
```

## Environment Variables

See `configs/.env.example` for all required environment variables.

Key variables:
- `SUPABASE_DATABASE_URL` - PostgreSQL connection string
- `SUPABASE_SERVICE_ROLE_KEY` - Supabase admin key
- `JWT_SECRET` - JWT signing secret
- `STRIPE_SECRET_KEY` - Stripe API key
- `STRIPE_WEBHOOK_SECRET` - Stripe webhook secret

## License

MIT
