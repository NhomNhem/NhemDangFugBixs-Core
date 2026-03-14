## 1. Environment Setup

- [x] 1.1 Add `github.com/stretchr/testify` to `go.mod`
- [x] 1.2 Create `test_helpers.go` for common testing utilities (e.g. Fiber setup)

## 2. Unit Testing (Usecase Layer)

- [x] 2.1 Implement unit tests for `AuthUsecase` using mock repositories
- [x] 2.2 Implement unit tests for `PlayerUsecase` (Save/Load logic)
- [x] 2.3 Implement unit tests for `LeaderboardUsecase`

## 3. API Integration Testing

- [x] 3.1 Create integration tests for Auth endpoints (`/login`, `/refresh`)
- [x] 3.2 Create integration tests for Save/Load endpoints (`/save`)
- [x] 3.3 Create integration tests for Leaderboard endpoints (`/leaderboard`)
- [x] 3.4 Create integration tests for Analytics endpoints (`/events`)

## 4. Build Verification

- [x] 4.1 Create a `verify_build.ps1` script to ensure code compiles and basic tests pass
- [x] 4.2 Validate Docker build completes successfully

## 5. Documentation & Finalization

- [x] 5.1 Update `README.md` with instructions on how to run the new test suite
- [x] 5.2 Document the mocking pattern for future developers
