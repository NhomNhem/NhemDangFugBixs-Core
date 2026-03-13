package usecase

import (
	"context"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
)

// AuthUsecase defines the interface for authentication business logic
type AuthUsecase interface {
	Login(ctx context.Context, sessionTicket string, overridePlayFabID string) (*models.HollowWildsAuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.RefreshTokenResponse, error)
	Logout(ctx context.Context, refreshToken string, jti string) error
	LegacyLogin(ctx context.Context, playfabID, displayName, sessionToken string) (*models.AuthResponse, error)
}
