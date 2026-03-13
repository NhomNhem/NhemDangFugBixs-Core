## 1. Validation Infrastructure

- [x] 1.1 Install validator package: `go get github.com/go-playground/validator/v10`
- [x] 1.2 Create `pkg/utils/validator.go` to initialize a global validator instance
- [x] 1.3 Add standard validation error response to `internal/domain/models/response.go`

## 2. Model & DTO Hardening

- [x] 2.1 Update `PlayerState` in `internal/domain/models/player.go` with `min/max` and `oneof` tags
- [x] 2.2 Update `HollowWildsLoginRequest` with mandatory field tags
- [x] 2.3 Add validation tags to `LeaderboardSubmitRequest`

## 3. Usecase & Cache Implementation

- [x] 3.1 Integrate validation check in `PlayerUsecase.SaveGame`
- [x] 3.2 Update `RedisCacheRepository` to support and enforce 300s TTL on save data
- [x] 3.3 Add validation logic to `AuthUsecase` for initial character selection

## 4. Testing & Verification

- [x] 4.1 Create unit test `internal/usecase/player/validation_test.go` to verify bound enforcement
- [x] 4.2 Verify Redis TTL manually using `redis-cli TTL` or logs
- [x] 4.3 Run integration tests to ensure valid saves still work
