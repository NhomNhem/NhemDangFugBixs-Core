## 1. Setup & Infrastructure

- [x] 1.1 Install `testify` and `mockery` dependencies: `go get github.com/stretchr/testify`
- [x] 1.2 Create a `Makefile` or `scripts/generate_mocks.sh` to automate mock generation
- [x] 1.3 Configure `mockery` to target `internal/domain/repository` and `internal/domain/usecase`

## 2. Unit Testing (Usecase Layer)

- [x] 2.1 Generate mocks for `PlayerRepository`, `IdentityRepository`, and `TokenRepository`
- [x] 2.2 Implement unit tests for `AuthUsecase` (Login, Refresh, Logout)
- [x] 2.3 Implement unit tests for `PlayerUsecase` (Save/Load logic)
- [x] 2.4 Implement unit tests for `LeaderboardUsecase`

## 3. Integration Testing (Delivery Layer)

- [x] 3.1 Create `internal/delivery/http/hollow_wilds_handler_test.go`
- [x] 3.2 Implement integration tests for the Login flow using Fiber's `app.Test`
- [x] 3.3 Implement integration tests for Save/Load endpoints
- [x] 3.4 Implement integration tests for Leaderboard submission

## 4. CI/CD Integration

- [x] 4.1 Update `.github/workflows/deploy.yml` to include a `Run Tests` step before deployment
- [x] 4.2 Verify the workflow fails correctly if a test is broken
