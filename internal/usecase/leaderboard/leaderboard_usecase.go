package leaderboard

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/usecase"
	"github.com/google/uuid"
)

type leaderboardUsecase struct {
	repo         repository.LeaderboardRepository
	playerRepo   repository.PlayerRepository
	identityRepo repository.IdentityRepository
	cacheRepo    repository.CacheRepository
}

// NewLeaderboardUsecase creates a new leaderboard usecase
func NewLeaderboardUsecase(
	repo repository.LeaderboardRepository,
	playerRepo repository.PlayerRepository,
	identityRepo repository.IdentityRepository,
	cacheRepo repository.CacheRepository,
) usecase.LeaderboardUsecase {
	return &leaderboardUsecase{
		repo:         repo,
		playerRepo:   playerRepo,
		identityRepo: identityRepo,
		cacheRepo:    cacheRepo,
	}
}

// ... existing Hollow Wilds methods ...

func (u *leaderboardUsecase) GetGlobalLeaderboard(ctx context.Context, levelID string, page, perPage int) (*models.GlobalLeaderboardResponse, error) {
	// 1. Get current version for this level
	versionKey := fmt.Sprintf("leaderboard:%s:version", levelID)
	version, _ := u.cacheRepo.Get(ctx, versionKey)
	if version == "" {
		version = "0"
	}

	// 2. Try cache with version
	cacheKey := fmt.Sprintf("leaderboard:%s:v%s:%d:%d", levelID, version, page, perPage)
	cached, err := u.cacheRepo.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var resp models.GlobalLeaderboardResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return &resp, nil
		}
	}

	// 3. Query DB
	limit := perPage
	offset := (page - 1) * perPage

	entries, total, err := u.repo.GetLegacyGlobal(ctx, levelID, limit, offset)
	if err != nil {
		return nil, err
	}

	resp := &models.GlobalLeaderboardResponse{
		Leaderboard: entries,
		Total:       total,
		Page:        page,
		PerPage:     perPage,
	}

	// 4. Cache (30s TTL as per design)
	if data, err := json.Marshal(resp); err == nil {
		u.cacheRepo.Set(ctx, cacheKey, string(data), 30*time.Second)
	}

	return resp, nil
}

func (u *leaderboardUsecase) GetPlayerRank(ctx context.Context, userID uuid.UUID, levelID string) (*models.PlayerStatsResponse, error) {
	rank, bestTime, stars, err := u.repo.GetLegacyPlayerRank(ctx, userID, levelID)
	if err != nil {
		return nil, err
	}

	// Fetch surrounding players (e.g., 3 above, 3 below)
	// If rank is 1, offset is 0. If rank is 10, offset is 10-4=6 (ranks 7-13).
	offset := rank - 4
	if offset < 0 {
		offset = 0
	}
	limit := 7

	surrounding, _, err := u.repo.GetLegacyGlobal(ctx, levelID, limit, offset)
	if err != nil {
		// Log error but don't fail the whole request
		fmt.Printf("Failed to fetch surrounding players: %v\n", err)
	}

	return &models.PlayerStatsResponse{
		PlayerID:           userID.String(),
		GlobalRank:         rank,
		BestTime:           bestTime,
		Stars:              stars,
		SurroundingPlayers: surrounding,
	}, nil
}

func (u *leaderboardUsecase) GetFriendsLeaderboard(ctx context.Context, userID uuid.UUID, levelID string) (*models.LevelLeaderboardResponse, error) {
	// 1. Get player's PlayFab ID (needed for friends lookup)
	player, err := u.playerRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, fmt.Errorf("player not found")
	}

	friends, err := u.identityRepo.GetFriends(ctx, player.PlayFabID)
	if err != nil {
		return nil, err
	}

	entries, err := u.repo.GetLegacyFriends(ctx, friends, levelID)
	if err != nil {
		return nil, err
	}

	return &models.LevelLeaderboardResponse{
		LevelID:     levelID,
		Leaderboard: entries,
		Total:       len(entries),
	}, nil
}

func (u *leaderboardUsecase) UpdateEntry(ctx context.Context, userID uuid.UUID, levelID string, timeSeconds float64, stars int) error {
	err := u.repo.UpsertLegacyEntry(ctx, userID, levelID, timeSeconds, stars)
	if err != nil {
		return err
	}

	// Invalidate cache by incrementing version
	versionKey := fmt.Sprintf("leaderboard:%s:version", levelID)
	version, _ := u.cacheRepo.Get(ctx, versionKey)
	v := 0
	if version != "" {
		fmt.Sscanf(version, "%d", &v)
	}
	u.cacheRepo.Set(ctx, versionKey, fmt.Sprintf("%d", v+1), 24*time.Hour)

	return nil
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
