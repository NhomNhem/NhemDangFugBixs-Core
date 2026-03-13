package repository

import (
	"context"
)

// IdentityRepository defines the interface for external identity providers (PlayFab)
type IdentityRepository interface {
	ValidateTicket(ctx context.Context, sessionTicket string) (string, error)
	GetFriends(ctx context.Context, playfabID string) ([]string, error)
}
