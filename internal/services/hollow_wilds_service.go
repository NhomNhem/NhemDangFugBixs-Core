package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/NhomNhem/GameFeel-Backend/internal/database"
	"github.com/NhomNhem/GameFeel-Backend/internal/models"
	"github.com/NhomNhem/GameFeel-Backend/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// HollowWildsService handles Hollow Wilds game operations
type HollowWildsService struct{}

// NewHollowWildsService creates a new Hollow Wilds service
func NewHollowWildsService() *HollowWildsService {
	return &HollowWildsService{}
}

// ValidatePlayFabTicket validates a PlayFab session ticket
func (s *HollowWildsService) ValidatePlayFabTicket(sessionTicket string, playfabID string) error {
	if sessionTicket == "" {
		return fmt.Errorf("session ticket is required")
	}

	if playfabID == "" {
		return fmt.Errorf("playfab ID is required")
	}

	titleID := os.Getenv("PLAYFAB_TITLE_ID")
	if titleID == "" || titleID == "DEV" {
		return nil // Skip validation in development mode
	}

	url := fmt.Sprintf("https://%s.playfabapi.com/Client/GetAccountInfo", titleID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization", sessionTicket)
	req.Body = nil

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to validate ticket: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid PlayFab session ticket")
	}

	return nil
}

// GetOrCreatePlayer gets existing player or creates new one
func (s *HollowWildsService) GetOrCreatePlayer(ctx context.Context, playfabID string, displayName *string) (*models.Player, error) {
	if database.Pool == nil {
		// Mock player for development
		return &models.Player{
			ID:          uuid.New(),
			PlayFabID:   playfabID,
			DisplayName: displayName,
			CreatedAt:   time.Now(),
			LastSeenAt:  time.Now(),
		}, nil
	}

	var player models.Player
	err := database.Pool.QueryRow(ctx, `
		SELECT id, playfab_id, display_name, created_at, last_seen_at
		FROM players
		WHERE playfab_id = $1
	`, playfabID).Scan(
		&player.ID, &player.PlayFabID, &player.DisplayName, &player.CreatedAt, &player.LastSeenAt,
	)

	if err == nil {
		// Update last seen
		_, err = database.Pool.Exec(ctx, `UPDATE players SET last_seen_at = NOW() WHERE id = $1`, player.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to update last seen: %w", err)
		}
		player.LastSeenAt = time.Now()
		return &player, nil
	}

	if err != pgx.ErrNoRows {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Create new player
	player.ID = uuid.New()
	player.PlayFabID = playfabID
	player.DisplayName = displayName
	player.CreatedAt = time.Now()
	player.LastSeenAt = time.Now()

	_, err = database.Pool.Exec(ctx, `
		INSERT INTO players (id, playfab_id, display_name, created_at, last_seen_at)
		VALUES ($1, $2, $3, $4, $5)
	`, player.ID, player.PlayFabID, player.DisplayName, player.CreatedAt, player.LastSeenAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	return &player, nil
}

// GenerateJWT generates a JWT token for the player
func (s *HollowWildsService) GenerateJWT(playerID uuid.UUID, playfabID string) (string, int, error) {
	return generateJWT(playerID.String(), playfabID)
}

// GenerateRefreshToken generates a refresh token and stores it in Redis
func (s *HollowWildsService) GenerateRefreshToken(ctx context.Context, playerID uuid.UUID) (string, error) {
	refreshToken := uuid.New().String()
	ttl := 7 * 24 * time.Hour // 7 days

	if utils.RedisClient == nil {
		return refreshToken, nil // Development mock
	}

	err := utils.StoreRefreshToken(ctx, refreshToken, playerID.String(), ttl)
	if err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return refreshToken, nil
}

// ValidateRefreshToken validates a refresh token
func (s *HollowWildsService) ValidateRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	if utils.RedisClient == nil {
		return uuid.New().String(), nil // Development mock
	}

	playerID, err := utils.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to validate refresh token: %w", err)
	}

	if playerID == "" {
		return "", fmt.Errorf("invalid or expired refresh token")
	}

	return playerID, nil
}

// RevokeRefreshToken revokes a refresh token
func (s *HollowWildsService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	return utils.DeleteRefreshToken(ctx, refreshToken)
}

// BlacklistJWT adds a JWT to the blacklist
func (s *HollowWildsService) BlacklistJWT(ctx context.Context, jti string) error {
	return utils.BlacklistJWT(ctx, jti, 24*time.Hour)
}

// GetPlayerSave retrieves a player's save data with Redis caching
func (s *HollowWildsService) GetPlayerSave(ctx context.Context, playerID uuid.UUID) (*models.PlayerSave, error) {
	if database.Pool == nil {
		// Mock save for development
		return &models.PlayerSave{
			ID:          uuid.New(),
			PlayerID:    playerID,
			SaveVersion: 1,
			UpdatedAt:   time.Now(),
			SaveData: models.GameSaveData{
				World:  models.WorldData{Seed: 12345, DayCount: 1},
				Player: models.PlayerState{Character: "RIMBA", Health: 100},
			},
		}, nil
	}

	// Try cache first
	cached, err := utils.GetCachedSaveData(ctx, playerID.String())
	if err == nil && cached != "" {
		var save models.PlayerSave
		if json.Unmarshal([]byte(cached), &save) == nil {
			return &save, nil
		}
	}

	// Cache miss - query database
	var save models.PlayerSave
	var saveDataJSON []byte
	err = database.Pool.QueryRow(ctx, `
		SELECT id, player_id, save_version, save_data, updated_at
		FROM player_saves
		WHERE player_id = $1
	`, playerID).Scan(
		&save.ID, &save.PlayerID, &save.SaveVersion, &saveDataJSON, &save.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil // No save found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get save: %w", err)
	}

	if err := json.Unmarshal(saveDataJSON, &save.SaveData); err != nil {
		return nil, fmt.Errorf("failed to parse save data: %w", err)
	}

	// Cache the result
	if saveDataStr, err := json.Marshal(save); err == nil {
		utils.CacheSaveData(ctx, playerID.String(), string(saveDataStr), 5*time.Minute)
	}

	return &save, nil
}

// SavePlayerSave saves a player's game data with version control
func (s *HollowWildsService) SavePlayerSave(ctx context.Context, playerID uuid.UUID, saveData models.GameSaveData, expectedVersion int) (*models.PlayerSave, error) {
	if database.Pool == nil {
		// Mock save for development
		return &models.PlayerSave{
			ID:          uuid.New(),
			PlayerID:    playerID,
			SaveVersion: expectedVersion + 1,
			UpdatedAt:   time.Now(),
			SaveData:    saveData,
		}, nil
	}

	// Validate character
	if !isValidCharacter(saveData.Player.Character) {
		return nil, fmt.Errorf("invalid character: %s", saveData.Player.Character)
	}

	tx, err := database.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get current version
	var currentVersion int
	err = tx.QueryRow(ctx, `SELECT save_version FROM player_saves WHERE player_id = $1`, playerID).Scan(&currentVersion)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to get current version: %w", err)
	}

	// Check version conflict
	if err == nil && expectedVersion != currentVersion {
		return nil, &models.VersionConflictError{
			ErrorCode:     "version_conflict",
			ServerVersion: currentVersion,
			Message:       "Save is outdated, fetch latest first",
		}
	}

	newVersion := currentVersion + 1
	if newVersion == 1 {
		newVersion = 1
	}

	// Marshal save data
	saveDataJSON, err := json.Marshal(saveData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal save data: %w", err)
	}

	// Create or update save
	var save models.PlayerSave
	save.PlayerID = playerID
	save.SaveVersion = newVersion
	save.SaveData = saveData

	if err == pgx.ErrNoRows {
		// Insert new save
		err = tx.QueryRow(ctx, `
			INSERT INTO player_saves (player_id, save_version, save_data, updated_at)
			VALUES ($1, $2, $3, NOW())
			RETURNING id, player_id, save_version, updated_at
		`, playerID, newVersion, saveDataJSON).Scan(
			&save.ID, &save.PlayerID, &save.SaveVersion, &save.UpdatedAt,
		)
	} else {
		// Update existing save
		err = tx.QueryRow(ctx, `
			UPDATE player_saves
			SET save_version = $2, save_data = $3, updated_at = NOW()
			WHERE player_id = $1
			RETURNING id, player_id, save_version, updated_at
		`, playerID, newVersion, saveDataJSON).Scan(
			&save.ID, &save.PlayerID, &save.SaveVersion, &save.UpdatedAt,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to save: %w", err)
	}

	// Check if we need to create automatic backup (major version threshold)
	if newVersion > 1 && newVersion%10 == 0 {
		s.createBackupInternal(ctx, tx, playerID, currentVersion, saveDataJSON)
	}

	// Update player's last_seen_at
	_, err = tx.Exec(ctx, `UPDATE players SET last_seen_at = NOW() WHERE id = $1`, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to update last seen: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	// Invalidate cache
	utils.InvalidateSaveCache(ctx, playerID.String())

	return &save, nil
}

// CreateBackup creates a manual backup of the player's save
func (s *HollowWildsService) CreateBackup(ctx context.Context, playerID uuid.UUID) (*models.PlayerSaveBackup, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database not connected")
	}

	// Get current save
	save, err := s.GetPlayerSave(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get save: %w", err)
	}
	if save == nil {
		return nil, fmt.Errorf("no save data to backup")
	}

	tx, err := database.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Marshal save data
	saveDataJSON, err := json.Marshal(save.SaveData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal save data: %w", err)
	}

	backup, err := s.createBackupInternal(ctx, tx, playerID, save.SaveVersion, saveDataJSON)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return backup, nil
}

// createBackupInternal creates a backup within a transaction
func (s *HollowWildsService) createBackupInternal(ctx context.Context, tx pgx.Tx, playerID uuid.UUID, saveVersion int, saveDataJSON []byte) (*models.PlayerSaveBackup, error) {
	// Count existing backups
	var count int
	err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM player_save_backups WHERE player_id = $1`, playerID).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to count backups: %w", err)
	}

	// Delete oldest if at limit (10 backups)
	if count >= 10 {
		_, err = tx.Exec(ctx, `
			DELETE FROM player_save_backups
			WHERE player_id = $1
			AND id = (
				SELECT id FROM player_save_backups
				WHERE player_id = $1
				ORDER BY created_at ASC
				LIMIT 1
			)
		`, playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete old backup: %w", err)
		}
	}

	// Create backup
	var backup models.PlayerSaveBackup
	backup.PlayerID = playerID
	backup.SaveVersion = saveVersion

	err = tx.QueryRow(ctx, `
		INSERT INTO player_save_backups (player_id, save_version, save_data, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, player_id, save_version, created_at
	`, playerID, saveVersion, saveDataJSON).Scan(
		&backup.ID, &backup.PlayerID, &backup.SaveVersion, &backup.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	return &backup, nil
}

// GetBackups lists all backups for a player
func (s *HollowWildsService) GetBackups(ctx context.Context, playerID uuid.UUID) ([]models.PlayerSaveBackup, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database not connected")
	}

	rows, err := database.Pool.Query(ctx, `
		SELECT id, player_id, save_version, created_at
		FROM player_save_backups
		WHERE player_id = $1
		ORDER BY created_at DESC
	`, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get backups: %w", err)
	}
	defer rows.Close()

	var backups []models.PlayerSaveBackup
	for rows.Next() {
		var backup models.PlayerSaveBackup
		if err := rows.Scan(&backup.ID, &backup.PlayerID, &backup.SaveVersion, &backup.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan backup: %w", err)
		}
		backups = append(backups, backup)
	}

	return backups, nil
}

// RestoreFromBackup restores a player's save from a backup
func (s *HollowWildsService) RestoreFromBackup(ctx context.Context, playerID uuid.UUID, backupID uuid.UUID) (*models.PlayerSave, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database not connected")
	}

	tx, err := database.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get backup
	var backup models.PlayerSaveBackup
	var saveDataJSON []byte
	err = tx.QueryRow(ctx, `
		SELECT id, player_id, save_version, save_data, created_at
		FROM player_save_backups
		WHERE id = $1 AND player_id = $2
	`, backupID, playerID).Scan(
		&backup.ID, &backup.PlayerID, &backup.SaveVersion, &saveDataJSON, &backup.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("backup not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get backup: %w", err)
	}

	// Get or create current save to get next version
	var currentVersion int
	err = tx.QueryRow(ctx, `SELECT save_version FROM player_saves WHERE player_id = $1`, playerID).Scan(&currentVersion)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to get current version: %w", err)
	}

	newVersion := currentVersion + 1

	// Update player save
	var save models.PlayerSave
	save.PlayerID = playerID
	save.SaveVersion = newVersion

	err = tx.QueryRow(ctx, `
		INSERT INTO player_saves (player_id, save_version, save_data, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (player_id) DO UPDATE
		SET save_version = $2, save_data = $3, updated_at = NOW()
		RETURNING id, player_id, save_version, updated_at
	`, playerID, newVersion, saveDataJSON).Scan(
		&save.ID, &save.PlayerID, &save.SaveVersion, &save.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to restore save: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	// Invalidate cache
	utils.InvalidateSaveCache(ctx, playerID.String())

	return &save, nil
}

// isValidCharacter checks if the character name is valid
func isValidCharacter(character string) bool {
	validChars := []string{"RIMBA", "DARA", "BAYU", "SARI"}
	for _, c := range validChars {
		if c == character {
			return true
		}
	}
	return false
}

// RecordAnalyticsEvents records a batch of analytics events
func (s *HollowWildsService) RecordAnalyticsEvents(ctx context.Context, playerID *uuid.UUID, events []models.AnalyticsEvent) (int, int) {
	if database.Pool == nil {
		return len(events), 0 // Mock all accepted in development
	}

	accepted := 0
	rejected := 0

	for _, event := range events {
		if event.EventName == "" || event.Timestamp == "" {
			rejected++
			continue
		}

		payloadJSON, _ := json.Marshal(event.Payload)

		var playerIDStr *string
		if playerID != nil {
			s := playerID.String()
			playerIDStr = &s
		}

		_, err := database.Pool.Exec(ctx, `
			INSERT INTO analytics_events (user_id, session_id, event_type, event_properties, created_at)
			VALUES ($1, $2, $3, $4, NOW())
		`, playerIDStr, event.SessionID, event.EventName, payloadJSON)

		if err != nil {
			rejected++
		} else {
			accepted++
		}
	}

	return accepted, rejected
}

// generateJWT is a helper to generate a JWT token
func generateJWT(playerID, playfabID string) (string, int, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key-123" // Default for development
	}

	expiresIn := 3600 // 1 hour as per design doc
	now := time.Now()
	expiresAt := now.Add(time.Duration(expiresIn) * time.Second)

	claims := jwt.MapClaims{
		"userId":    playerID,
		"sub":       playerID,
		"playfabId": playfabID,
		"iat":       now.Unix(),
		"exp":       expiresAt.Unix(),
		"jti":       uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresIn, nil
}
