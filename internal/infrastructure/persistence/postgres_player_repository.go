package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresPlayerRepository struct {
	db *pgxpool.Pool
}

// NewPostgresPlayerRepository creates a new PostgreSQL player repository
func NewPostgresPlayerRepository(db *pgxpool.Pool) repository.PlayerRepository {
	return &postgresPlayerRepository{db: db}
}

func (r *postgresPlayerRepository) GetByPlayFabID(ctx context.Context, playfabID string) (*models.Player, error) {
	if r.db == nil {
		// Mock player for development
		return &models.Player{
			ID:         uuid.New(),
			PlayFabID:  playfabID,
			CreatedAt:  time.Now(),
			LastSeenAt: time.Now(),
		}, nil
	}

	var player models.Player
	err := r.db.QueryRow(ctx, `
		SELECT id, playfab_id, display_name, created_at, last_seen_at
		FROM players
		WHERE playfab_id = $1
	`, playfabID).Scan(
		&player.ID, &player.PlayFabID, &player.DisplayName, &player.CreatedAt, &player.LastSeenAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get player by playfab id: %w", err)
	}

	return &player, nil
}

func (r *postgresPlayerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Player, error) {
	if r.db == nil {
		return &models.Player{
			ID:         id,
			PlayFabID:  "MOCK_PLAYFAB_ID",
			CreatedAt:  time.Now(),
			LastSeenAt: time.Now(),
		}, nil
	}

	var player models.Player
	err := r.db.QueryRow(ctx, `
		SELECT id, playfab_id, display_name, created_at, last_seen_at
		FROM players
		WHERE id = $1
	`, id).Scan(
		&player.ID, &player.PlayFabID, &player.DisplayName, &player.CreatedAt, &player.LastSeenAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get player by id: %w", err)
	}

	return &player, nil
}

func (r *postgresPlayerRepository) Create(ctx context.Context, player *models.Player) error {
	if r.db == nil {
		return nil // Success in mock mode
	}
	err := r.db.QueryRow(ctx, `
		INSERT INTO players (playfab_id, display_name, created_at, last_seen_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, player.PlayFabID, player.DisplayName, player.CreatedAt, player.LastSeenAt).Scan(&player.ID)

	if err != nil {
		return fmt.Errorf("failed to create player: %w", err)
	}

	return nil
}

func (r *postgresPlayerRepository) UpdateLastSeen(ctx context.Context, playerID uuid.UUID) error {
	if r.db == nil {
		return nil // Success in mock mode
	}
	_, err := r.db.Exec(ctx, `UPDATE players SET last_seen_at = NOW() WHERE id = $1`, playerID)
	if err != nil {
		return fmt.Errorf("failed to update player last seen: %w", err)
	}
	return nil
}
