package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/mocks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthUsecase_Login(t *testing.T) {
	playerRepo := new(repo_mock.PlayerRepository)
	identityRepo := new(repo_mock.IdentityRepository)
	tokenRepo := new(repo_mock.TokenRepository)
	usecase := NewAuthUsecase(playerRepo, identityRepo, tokenRepo)

	ctx := context.Background()
	ticket := "valid-ticket"
	playfabID := "PLAYFAB_123"
	playerID := uuid.New()

	t.Run("successful login for existing player", func(t *testing.T) {
		identityRepo.On("ValidateTicket", ctx, ticket).Return(playfabID, nil).Once()
		playerRepo.On("GetByPlayFabID", ctx, playfabID).Return(&models.Player{
			ID:        playerID,
			PlayFabID: playfabID,
		}, nil).Once()
		playerRepo.On("UpdateLastSeen", ctx, playerID).Return(nil).Once()
		tokenRepo.On("StoreRefreshToken", ctx, mock.Anything, playerID.String(), mock.Anything).Return(nil).Once()

		resp, err := usecase.Login(ctx, ticket, "")

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, playerID.String(), resp.PlayerID)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
		
		identityRepo.AssertExpectations(t)
		playerRepo.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})

	t.Run("failed validation", func(t *testing.T) {
		identityRepo.On("ValidateTicket", ctx, "invalid").Return("", errors.New("invalid ticket")).Once()

		resp, err := usecase.Login(ctx, "invalid", "")

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "authentication failed")
	})
}

func TestAuthUsecase_Logout(t *testing.T) {
	playerRepo := new(repo_mock.PlayerRepository)
	identityRepo := new(repo_mock.IdentityRepository)
	tokenRepo := new(repo_mock.TokenRepository)
	usecase := NewAuthUsecase(playerRepo, identityRepo, tokenRepo)

	ctx := context.Background()
	refreshToken := "ref-123"
	jti := "jti-123"

	t.Run("successful logout", func(t *testing.T) {
		tokenRepo.On("DeleteRefreshToken", ctx, refreshToken).Return(nil).Once()
		tokenRepo.On("BlacklistJWT", ctx, jti, mock.Anything).Return(nil).Once()

		err := usecase.Logout(ctx, refreshToken, jti)

		assert.NoError(t, err)
		tokenRepo.AssertExpectations(t)
	})
}
