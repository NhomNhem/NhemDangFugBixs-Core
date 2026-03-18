# AGENTS.md - Agent Coding Guidelines for GameFeel Backend

This document provides guidelines for AI agents working on this codebase.

## Project Overview
- **Language**: Go 1.25.1
- **Framework**: Fiber v2 (HTTP server)
- **Database**: Supabase (PostgreSQL via pgx)
- **Cache**: Upstash Redis
- **Auth**: JWT + PlayFab
- **Testing**: testify + mockery

## Build & Test Commands
```bash
go run cmd/server/main.go           # Run server
go build -o bin/server cmd/server/main.go  # Build binary
go test ./...                       # Run all tests
go test -run TestName ./path/       # Run single test
go test -v ./...                    # Verbose output
golangci-lint run                   # Linting
./scripts/check-quality.ps1         # Quality check (lint + tests)
./scripts/verify_build.ps1          # Full build verification
./scripts/generate-mocks.ps1        # Generate mocks
```

## Project Structure
```
cmd/server/           # Entry point
internal/
  api/               # HTTP handlers (delivery layer)
  domain/models/     # Data models/DTOs
  domain/repository/ # Repository interfaces
  domain/usecase/    # Business logic interfaces
  infrastructure/    # PostgreSQL, Redis, PlayFab implementations
  usecase/           # Business logic implementations
  middleware/        # Auth, logging, rate limiting
  database/          # DB connection management
pkg/utils/           # Shared utilities
```

## Naming Conventions
- **Files**: snake_case.go (auth_handler.go, player_repository.go)
- **Packages**: lowercase (auth, player, models)
- **Structs**: PascalCase (AuthHandler, PlayerRepository)
- **Interfaces**: PascalCase + suffix (PlayerRepository, AuthUsecase)
- **Variables**: camelCase (playerRepo, authUsecase)
- **Constants**: PascalCase or SCREAMING_SNAKE for error codes

## Import Organization
Order imports with blank lines between groups: stdlib, external (github.com/), local (github.com/NhomNhem/...).
```go
import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
)
```

## Error Handling
Use `fmt.Errorf("context: %w", err)` for wrapped errors. HTTP handlers return errors with appropriate status codes.
```go
// In usecases
if err != nil {
    return nil, fmt.Errorf("authentication failed: %w", err)
}
// In handlers
if err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
        Success: false,
        Error: &models.APIError{
            Code:    models.ErrCodeInvalidRequest,
            Message: err.Error(),
        },
    })
}
```

## Response Format
All API responses use: `APIResponse{Success bool, Data interface{}, Error *APIError}`. Use error codes from models package:
```go
const (
    ErrCodeInvalidRequest = "INVALID_REQUEST"
    ErrCodeUnauthorized  = "UNAUTHORIZED"
    ErrCodeInvalidToken  = "INVALID_TOKEN"
    ErrCodeInternalError = "INTERNAL_ERROR"
)
```

## Handler Pattern
Handlers receive fiber.Ctx and return error:
```go
func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var req models.AuthRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
            Success: false,
            Error: &models.APIError{Code: models.ErrCodeInvalidRequest, Message: "Invalid request body"},
        })
    }
    return c.JSON(models.APIResponse{Success: true, Data: resp})
}
```

## Usecase Pattern
Constructor injection returning domain models:
```go
func NewAuthUsecase(playerRepo repository.PlayerRepository, identityRepo repository.IdentityRepository) usecase.AuthUsecase {
    return &authUsecase{playerRepo: playerRepo, identityRepo: identityRepo}
}
```

## Testing Patterns
- **Unit tests**: Use mockery mocks from `internal/mocks/`, testify assertions. Test file `*_test.go` next to implementation.
- **Integration tests**: Use `internal/api/test_helpers.go` -> `SetupTestApp()` for full HTTP cycle.
```go
func TestAuthUsecase_Login(t *testing.T) {
    playerRepo := new(repository_mock.MockPlayerRepository)
    identityRepo.On("ValidateTicket", ctx, ticket).Return(playfabID, nil).Once()
    resp, err := usecase.Login(ctx, ticket)
    assert.NoError(t, err)
}
```

## Swagger Documentation
Add annotations above handlers:
```go
// @Summary Login with PlayFab
// @Tags Authentication
// @Success 200 {object} models.APIResponse{data=models.AuthResponse}
// @Router /auth/login [post]
```

## General Guidelines
- Keep functions focused and small
- Use context.Context as first parameter for all operations
- Validate input at handler layer
- Use meaningful variable names
- No comments unless explaining complex logic
