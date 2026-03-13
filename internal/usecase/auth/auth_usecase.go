package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/usecase"
	"github.com/NhomNhem/HollowWilds-Backend/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authUsecase struct {
	playerRepo   repository.PlayerRepository
	identityRepo repository.IdentityRepository
	tokenRepo    repository.TokenRepository
}

// NewAuthUsecase creates a new authentication usecase
func NewAuthUsecase(
	playerRepo repository.PlayerRepository,
	identityRepo repository.IdentityRepository,
	tokenRepo repository.TokenRepository,
) usecase.AuthUsecase {
	return &authUsecase{
		playerRepo:   playerRepo,
		identityRepo: identityRepo,
		tokenRepo:    tokenRepo,
	}
}

func (u *authUsecase) Login(ctx context.Context, sessionTicket string, overridePlayFabID string) (*models.HollowWildsAuthResponse, error) {
	// 1. Validate PlayFab ticket
	playfabID, err := u.identityRepo.ValidateTicket(ctx, sessionTicket)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// 2. Allow override in dev mode
	if playfabID == "MOCK_PLAYFAB_ID" && overridePlayFabID != "" {
		playfabID = overridePlayFabID
	}

	// 3. Get or create player
	player, err := u.playerRepo.GetByPlayFabID(ctx, playfabID)
	if err != nil {
		return nil, err
	}

	if player == nil {
		player = &models.Player{
			ID:         uuid.New(),
			PlayFabID:  playfabID,
			CreatedAt:  time.Now(),
			LastSeenAt: time.Now(),
		}
		if err := u.playerRepo.Create(ctx, player); err != nil {
			return nil, err
		}
	} else {
		// Update last seen
		if err := u.playerRepo.UpdateLastSeen(ctx, player.ID); err != nil {
			// Non-critical, just log or continue
		}
	}

	// 4. Generate JWT
	token, expiresIn, err := u.generateJWT(player.ID.String(), player.PlayFabID)
	if err != nil {
		return nil, err
	}

	// 5. Generate Refresh Token
	refreshToken := uuid.New().String()
	if err := u.tokenRepo.StoreRefreshToken(ctx, refreshToken, player.ID.String(), 7*24*time.Hour); err != nil {
		// Log error but return access token anyway? Or fail?
		// For consistency with old implementation, we continue
	}

	return &models.HollowWildsAuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		PlayerID:     player.ID.String(),
	}, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (*models.RefreshTokenResponse, error) {
	playerIDStr, err := u.tokenRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil || playerIDStr == "" {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	// We need PlayFabID for JWT, but it's not in the refresh token.
	// For now, we'll use a placeholder as in the old implementation,
	// or we could look up the player.
	token, expiresIn, err := u.generateJWT(playerIDStr, "REFRESHED")
	if err != nil {
		return nil, err
	}

	return &models.RefreshTokenResponse{
		Token:     token,
		ExpiresIn: expiresIn,
	}, nil
}

func (u *authUsecase) Logout(ctx context.Context, refreshToken string, jti string) error {
	if refreshToken != "" {
		u.tokenRepo.DeleteRefreshToken(ctx, refreshToken)
	}
	if jti != "" {
		u.tokenRepo.BlacklistJWT(ctx, jti, 24*time.Hour)
	}
	return nil
}

func (u *authUsecase) generateJWT(playerID, playfabID string) (string, int, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key-123"
	}

	expiresIn := 3600
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
