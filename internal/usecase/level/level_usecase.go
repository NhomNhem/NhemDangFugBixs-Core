package level

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/usecase"
	"github.com/google/uuid"
)

type levelUsecase struct {
	levelRepo repository.LevelRepository
	lbUsecase usecase.LeaderboardUsecase
}

// NewLevelUsecase creates a new level usecase
func NewLevelUsecase(levelRepo repository.LevelRepository, lbUsecase usecase.LeaderboardUsecase) usecase.LevelUsecase {
	return &levelUsecase{
		levelRepo: levelRepo,
		lbUsecase: lbUsecase,
	}
}

func (u *levelUsecase) CompleteLevel(ctx context.Context, userID uuid.UUID, req *models.LevelCompletionRequest) (*models.LevelCompletionResponse, error) {
	// 1. Get level configuration
	config, err := u.levelRepo.GetConfig(ctx, req.LevelID, req.MapID)
	if err != nil {
		return nil, fmt.Errorf("failed to get level config: %w", err)
	}

	// 2. Anti-cheat validation
	antiCheat := u.validateCompletion(req, config)
	if antiCheat.Level >= models.AntiCheatSevere {
		log.Printf("Anti-cheat SEVERE for user %s: %v", userID, antiCheat.Reasons)
		return nil, fmt.Errorf("invalid completion data")
	}

	// 3. Calculate stars and rewards
	starsEarned := u.calculateStars(req, config)
	goldEarned := u.calculateGoldReward(starsEarned, config.BaseGold)

	// 4. Check existing completion
	existingCompletion, err := u.levelRepo.GetCompletion(ctx, userID, req.LevelID)
	if err != nil {
		return nil, err
	}

	var isFirstCompletion bool
	var newBestTime bool

	if existingCompletion == nil {
		isFirstCompletion = true
		newBestTime = true

		newCompletion := &models.LevelCompletion{
			ID:                  uuid.New(),
			UserID:              userID,
			LevelID:             req.LevelID,
			MapID:               req.MapID,
			StarsEarned:         starsEarned,
			BestTimeSeconds:     req.TimeSeconds,
			PlayCount:           1,
			LastFinalHP:         &req.FinalHP,
			LastDashCount:       &req.DashCount,
			LastCounterCount:    &req.CounterCount,
			LastVulnerableKills: &req.VulnerableKills,
			FirstCompletedAt:    time.Now(),
			LastPlayedAt:        time.Now(),
		}

		if err := u.levelRepo.CreateCompletion(ctx, newCompletion); err != nil {
			return nil, err
		}
	} else {
		isFirstCompletion = false
		newBestTime = req.TimeSeconds < existingCompletion.BestTimeSeconds
		if newBestTime {
			existingCompletion.BestTimeSeconds = req.TimeSeconds
		}

		if starsEarned > existingCompletion.StarsEarned {
			existingCompletion.StarsEarned = starsEarned
		}

		existingCompletion.PlayCount++
		existingCompletion.LastFinalHP = &req.FinalHP
		existingCompletion.LastDashCount = &req.DashCount
		existingCompletion.LastCounterCount = &req.CounterCount
		existingCompletion.LastVulnerableKills = &req.VulnerableKills
		existingCompletion.LastPlayedAt = time.Now()

		if err := u.levelRepo.UpdateCompletion(ctx, existingCompletion); err != nil {
			return nil, err
		}

		// Only give gold for first completion
		goldEarned = 0
	}

	// 5. Update user gold and stars
	newTotalGold, newTotalStars, err := u.levelRepo.UpdateUserStats(ctx, userID, goldEarned, starsEarned)
	if err != nil {
		return nil, err
	}

	// 6. Update Leaderboard (Legacy)
	if u.lbUsecase != nil {
		if err := u.lbUsecase.UpdateEntry(ctx, userID, req.LevelID, req.TimeSeconds, starsEarned); err != nil {
			// Log but don't fail the completion
			log.Printf("Failed to update leaderboard for user %s: %v", userID, err)
		}
	}

	return &models.LevelCompletionResponse{
		Success:           true,
		StarsEarned:       starsEarned,
		GoldEarned:        goldEarned,
		NewTotalGold:      newTotalGold,
		NewTotalStars:     newTotalStars,
		IsFirstCompletion: isFirstCompletion,
		NewBestTime:       newBestTime,
	}, nil
}

func (u *levelUsecase) validateCompletion(req *models.LevelCompletionRequest, config *models.LevelConfig) *models.AntiCheatResult {
	result := &models.AntiCheatResult{
		Level:   models.AntiCheatNone,
		Reasons: []string{},
	}

	if req.TimeSeconds < config.MinTimeSeconds {
		result.Level = models.AntiCheatSevere
		result.Reasons = append(result.Reasons, fmt.Sprintf("Completion time %.2fs is below minimum %.2fs", req.TimeSeconds, config.MinTimeSeconds))
	}

	if req.FinalHP > 100 {
		result.Level = models.AntiCheatSevere
		result.Reasons = append(result.Reasons, fmt.Sprintf("Final HP %.2f exceeds maximum 100", req.FinalHP))
	}

	return result
}

func (u *levelUsecase) calculateStars(req *models.LevelCompletionRequest, config *models.LevelConfig) int {
	starsEarned := 0
	for _, obj := range config.Objectives {
		var value float64
		switch obj.Type {
		case "completion":
			value = 1
		case "health":
			value = req.FinalHP
		case "time":
			value = req.TimeSeconds
		default:
			continue
		}

		objectiveMet := false
		switch obj.Operator {
		case "gte":
			objectiveMet = value >= obj.Threshold
		case "lte":
			objectiveMet = value <= obj.Threshold
		}

		if objectiveMet {
			starsEarned++
		}
	}
	return starsEarned
}

func (u *levelUsecase) calculateGoldReward(starsEarned int, baseGold int) int {
	return baseGold + (starsEarned * 20)
}
