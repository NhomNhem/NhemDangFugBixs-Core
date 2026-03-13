package repository

import (
	"context"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
)

// PlayerRepository defines the interface for player data access
type PlayerRepository interface {
	GetByPlayFabID(ctx context.Context, playfabID string) (*models.Player, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Player, error)
	Create(ctx context.Context, player *models.Player) error
	UpdateLastSeen(ctx context.Context, playerID uuid.UUID) error
}
