package leaderboard

import (
	"context"
	"fmt"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/usecase"
	"github.com/google/uuid"
)

type leaderboardUsecase struct {
	repo repository.LeaderboardRepository
}

// NewLeaderboardUsecase creates a new leaderboard usecase
func NewLeaderboardUsecase(repo repository.LeaderboardRepository) usecase.LeaderboardUsecase {
	return &leaderboardUsecase{repo: repo}
}

func (u *leaderboardUsecase) GetLeaderboard(ctx context.Context, lbType, scope, character string, limit, offset int) (*models.HollowWildsLeaderboardResponse, error) {
	entries, total, err := u.repo.GetByCriteria(ctx, lbType, scope, character, limit, offset)
	if err != nil {
		return nil, err
	}

	return &models.HollowWildsLeaderboardResponse{
		Type:      lbType,
		Scope:     scope,
		Character: character,
		Total:     total,
		Entries:   entries,
	}, nil
}

func (u *leaderboardUsecase) SubmitEntry(ctx context.Context, playerID uuid.UUID, req models.LeaderboardSubmitRequest) (*models.LeaderboardSubmitResponse, error) {
	// 1. Get current personal best
	currentBest, err := u.repo.GetPersonalBest(ctx, playerID, req.Type, req.Character)
	if err != nil {
		return nil, err
	}

	if req.Value <= currentBest {
		return nil, fmt.Errorf("value_too_low: Submitted value does not beat personal best")
	}

	// 2. Submit entry
	entry := &models.HollowWildsLeaderboardEntry{
		PlayerID:    playerID.String(),
		Value:       req.Value,
		Character:   req.Character,
		WorldSeed:   req.WorldSeed,
		CombatBuild: req.CombatBuild,
		RunMetadata: req.RunMetadata,
	}

	if err := u.repo.UpsertEntry(ctx, entry); err != nil {
		return nil, err
	}

	// 3. Get updated ranks (simplified for now as in old implementation)
	// In a real app, the repository should return these or we calculate them.
	// For now, we'll return a mock success response to maintain parity.
	return &models.LeaderboardSubmitResponse{
		Success:        true,
		IsPersonalBest: true,
		GlobalRank:     1, // Mock
		CharacterRank:  1, // Mock
	}, nil
}

func (u *leaderboardUsecase) GetPlayerStats(ctx context.Context, playerID uuid.UUID) (*models.PlayerLeaderboardResponse, error) {
	entries, err := u.repo.GetPlayerStats(ctx, playerID)
	if err != nil {
		return nil, err
	}

	return &models.PlayerLeaderboardResponse{
		Entries: entries,
	}, nil
}
