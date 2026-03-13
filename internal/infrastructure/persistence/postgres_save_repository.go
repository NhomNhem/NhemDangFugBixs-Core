package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresSaveRepository struct {
	db *pgxpool.Pool
}

// NewPostgresSaveRepository creates a new PostgreSQL save repository
func NewPostgresSaveRepository(db *pgxpool.Pool) repository.SaveRepository {
	return &postgresSaveRepository{db: db}
}

func (r *postgresSaveRepository) GetByPlayerID(ctx context.Context, playerID uuid.UUID) (*models.PlayerSave, error) {
	if r.db == nil {
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
	var save models.PlayerSave
	var saveDataJSON []byte
	err := r.db.QueryRow(ctx, `
		SELECT id, player_id, save_version, save_data, updated_at
		FROM player_saves
		WHERE player_id = $1
	`, playerID).Scan(
		&save.ID, &save.PlayerID, &save.SaveVersion, &saveDataJSON, &save.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get save: %w", err)
	}

	if err := json.Unmarshal(saveDataJSON, &save.SaveData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal save data: %w", err)
	}

	return &save, nil
}

func (r *postgresSaveRepository) Upsert(ctx context.Context, save *models.PlayerSave) error {
	if r.db == nil {
		save.ID = uuid.New()
		save.UpdatedAt = time.Now()
		return nil
	}
	saveDataJSON, err := json.Marshal(save.SaveData)
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %w", err)
	}

	err = r.db.QueryRow(ctx, `
		INSERT INTO player_saves (player_id, save_version, save_data, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (player_id) DO UPDATE
		SET save_version = EXCLUDED.save_version,
		    save_data = EXCLUDED.save_data,
		    updated_at = NOW()
		RETURNING id, updated_at
	`, save.PlayerID, save.SaveVersion, saveDataJSON).Scan(&save.ID, &save.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to upsert save: %w", err)
	}

	return nil
}

func (r *postgresSaveRepository) CreateBackup(ctx context.Context, backup *models.PlayerSaveBackup) error {
	if r.db == nil {
		backup.ID = uuid.New()
		backup.CreatedAt = time.Now()
		return nil
	}
	saveDataJSON, err := json.Marshal(backup.SaveData)
	if err != nil {
		return fmt.Errorf("failed to marshal backup data: %w", err)
	}

	err = r.db.QueryRow(ctx, `
		INSERT INTO player_save_backups (player_id, save_version, save_data, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at
	`, backup.PlayerID, backup.SaveVersion, saveDataJSON).Scan(&backup.ID, &backup.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	return nil
}

func (r *postgresSaveRepository) GetBackupsByPlayerID(ctx context.Context, playerID uuid.UUID) ([]models.PlayerSaveBackup, error) {
	if r.db == nil {
		return []models.PlayerSaveBackup{
			{
				ID:          uuid.New(),
				PlayerID:    playerID,
				SaveVersion: 1,
				CreatedAt:   time.Now(),
			},
		}, nil
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, player_id, save_version, save_data, created_at
		FROM player_save_backups
		WHERE player_id = $1
		ORDER BY created_at DESC
	`, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query backups: %w", err)
	}
	defer rows.Close()

	var backups []models.PlayerSaveBackup
	for rows.Next() {
		var backup models.PlayerSaveBackup
		var saveDataJSON []byte
		err := rows.Scan(&backup.ID, &backup.PlayerID, &backup.SaveVersion, &saveDataJSON, &backup.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan backup: %w", err)
		}
		if err := json.Unmarshal(saveDataJSON, &backup.SaveData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal backup data: %w", err)
		}
		backups = append(backups, backup)
	}

	return backups, nil
}

func (r *postgresSaveRepository) CountBackups(ctx context.Context, playerID uuid.UUID) (int, error) {
	if r.db == nil {
		return 1, nil
	}
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM player_save_backups WHERE player_id = $1`, playerID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count backups: %w", err)
	}
	return count, nil
}

func (r *postgresSaveRepository) DeleteOldestBackup(ctx context.Context, playerID uuid.UUID) error {
	if r.db == nil {
		return nil
	}
	_, err := r.db.Exec(ctx, `
		DELETE FROM player_save_backups
		WHERE id = (
			SELECT id FROM player_save_backups
			WHERE player_id = $1
			ORDER BY created_at ASC
			LIMIT 1
		)
	`, playerID)
	if err != nil {
		return fmt.Errorf("failed to delete oldest backup: %w", err)
	}
	return nil
}
