package services

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/NhomNhem/HollowWilds-Backend/internal/database"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// LevelService handles level completion logic
type LevelService struct{}

// NewLevelService creates a new level service
func NewLevelService() *LevelService {
	return &LevelService{}
}

// GetLevelConfig returns configuration for a level
// TODO: Store in database, for now use hardcoded values
func (s *LevelService) GetLevelConfig(levelID string, mapID string) (*models.LevelConfig, error) {
	// Hardcoded config for testing (TODO: move to database)
	return &models.LevelConfig{
		LevelID:        levelID,
		MapID:          mapID,
		MinTimeSeconds: 10.0, // Minimum 10 seconds (anti-cheat)
		BaseGold:       60,   // Base reward for map 1
		Objectives: []models.LevelObjective{
			{Type: "completion", Threshold: 1, Operator: "gte"}, // Always star 1
			{Type: "health", Threshold: 50, Operator: "gte"},    // HP > 50%
			{Type: "time", Threshold: 60, Operator: "lte"},      // Time <= 60s
		},
	}, nil
}

// ValidateCompletion validates level completion with anti-cheat
func (s *LevelService) ValidateCompletion(req *models.LevelCompletionRequest, config *models.LevelConfig) *models.AntiCheatResult {
	result := &models.AntiCheatResult{
		Level:   models.AntiCheatNone,
		Reasons: []string{},
	}

	// Check 1: Time too fast (impossible)
	if req.TimeSeconds < config.MinTimeSeconds {
		result.Level = models.AntiCheatSevere
		result.Reasons = append(result.Reasons, fmt.Sprintf("Completion time %.2fs is below minimum %.2fs", req.TimeSeconds, config.MinTimeSeconds))
		return result
	}

	// Check 2: HP above 100%
	if req.FinalHP > 100 {
		result.Level = models.AntiCheatSevere
		result.Reasons = append(result.Reasons, fmt.Sprintf("Final HP %.2f exceeds maximum 100", req.FinalHP))
		return result
	}

	// Check 3: Suspiciously high action counts (moderate)
	if req.DashCount > 200 || req.CounterCount > 200 {
		result.Level = models.AntiCheatSuspicious
		result.Reasons = append(result.Reasons, "Unusually high action counts")
	}

	// Check 4: Perfect time (suspicious for some levels)
	if req.TimeSeconds == config.MinTimeSeconds {
		result.Level = models.AntiCheatSuspicious
		result.Reasons = append(result.Reasons, "Exactly minimum time")
	}

	return result
}

// CalculateStars calculates stars earned based on objectives
func (s *LevelService) CalculateStars(req *models.LevelCompletionRequest, config *models.LevelConfig) int {
	starsEarned := 0

	for _, obj := range config.Objectives {
		var value float64

		switch obj.Type {
		case "completion":
			value = 1 // Always true if level completed
		case "health":
			value = req.FinalHP
		case "time":
			value = req.TimeSeconds
		case "dash_count":
			value = float64(req.DashCount)
		case "counter_count":
			value = float64(req.CounterCount)
		case "vulnerable_kills":
			value = float64(req.VulnerableKills)
		default:
			continue
		}

		// Check if objective met
		objectiveMet := false
		switch obj.Operator {
		case "gte":
			objectiveMet = value >= obj.Threshold
		case "gt":
			objectiveMet = value > obj.Threshold
		case "lte":
			objectiveMet = value <= obj.Threshold
		case "lt":
			objectiveMet = value < obj.Threshold
		case "eq":
			objectiveMet = value == obj.Threshold
		}

		if objectiveMet {
			starsEarned++
		}
	}

	return starsEarned
}

// CalculateGoldReward calculates gold reward based on stars
func (s *LevelService) CalculateGoldReward(starsEarned int, baseGold int) int {
	// Base gold + (stars * 20 gold per star)
	return baseGold + (starsEarned * 20)
}

// CompleteLevel handles full level completion logic
func (s *LevelService) CompleteLevel(ctx context.Context, userID uuid.UUID, req *models.LevelCompletionRequest) (*models.LevelCompletionResponse, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database not connected")
	}

	// Get level configuration
	config, err := s.GetLevelConfig(req.LevelID, req.MapID)
	if err != nil {
		return nil, fmt.Errorf("failed to get level config: %w", err)
	}

	// Anti-cheat validation
	antiCheat := s.ValidateCompletion(req, config)
	if antiCheat.Level >= models.AntiCheatSevere {
		log.Printf("Anti-cheat SEVERE for user %s: %v", userID, antiCheat.Reasons)
		return nil, fmt.Errorf("invalid completion data")
	}

	if antiCheat.Level == models.AntiCheatModerate {
		log.Printf("Anti-cheat MODERATE for user %s: %v", userID, antiCheat.Reasons)
		// Continue but log
	}

	// Calculate stars and rewards
	starsEarned := s.CalculateStars(req, config)
	goldEarned := s.CalculateGoldReward(starsEarned, config.BaseGold)

	// Start transaction
	tx, err := database.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check existing completion
	var existingCompletion models.LevelCompletion
	var isFirstCompletion bool
	var newBestTime bool

	err = tx.QueryRow(ctx, `
		SELECT id, stars_earned, best_time_seconds, play_count
		FROM level_completions
		WHERE user_id = $1 AND level_id = $2
	`, userID, req.LevelID).Scan(
		&existingCompletion.ID,
		&existingCompletion.StarsEarned,
		&existingCompletion.BestTimeSeconds,
		&existingCompletion.PlayCount,
	)

	if err == pgx.ErrNoRows {
		// First completion - insert
		isFirstCompletion = true
		newBestTime = true

		_, err = tx.Exec(ctx, `
			INSERT INTO level_completions (id, user_id, level_id, map_id, stars_earned, best_time_seconds, play_count,
			                               last_final_hp, last_dash_count, last_counter_count, last_vulnerable_kills,
			                               first_completed_at, last_played_at)
			VALUES ($1, $2, $3, $4, $5, $6, 1, $7, $8, $9, $10, NOW(), NOW())
		`, uuid.New(), userID, req.LevelID, req.MapID, starsEarned, req.TimeSeconds,
			req.FinalHP, req.DashCount, req.CounterCount, req.VulnerableKills)

		if err != nil {
			return nil, fmt.Errorf("failed to insert completion: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to query completion: %w", err)
	} else {
		// Update existing - only if better
		isFirstCompletion = false
		newBestTime = req.TimeSeconds < existingCompletion.BestTimeSeconds
		bestTime := existingCompletion.BestTimeSeconds
		if newBestTime {
			bestTime = req.TimeSeconds
		}

		bestStars := existingCompletion.StarsEarned
		if starsEarned > bestStars {
			bestStars = starsEarned
		}

		_, err = tx.Exec(ctx, `
			UPDATE level_completions
			SET stars_earned = $1,
			    best_time_seconds = $2,
			    play_count = play_count + 1,
			    last_final_hp = $3,
			    last_dash_count = $4,
			    last_counter_count = $5,
			    last_vulnerable_kills = $6,
			    last_played_at = NOW()
			WHERE id = $7
		`, bestStars, bestTime, req.FinalHP, req.DashCount, req.CounterCount, req.VulnerableKills, existingCompletion.ID)

		if err != nil {
			return nil, fmt.Errorf("failed to update completion: %w", err)
		}

		// Only give gold for first completion
		if !isFirstCompletion {
			goldEarned = 0
		}
	}

	// Update user gold and stars
	var newTotalGold, newTotalStars int
	err = tx.QueryRow(ctx, `
		UPDATE users
		SET gold = gold + $1,
		    total_stars_collected = total_stars_collected + $2,
		    last_played_at = NOW()
		WHERE id = $3
		RETURNING gold, total_stars_collected
	`, goldEarned, starsEarned, userID).Scan(&newTotalGold, &newTotalStars)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// TODO: Check if next level/map unlocked

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

// CanUnlockMap checks if user has enough stars to unlock next map
func (s *LevelService) CanUnlockMap(playerStars, previousMapLevels int) bool {
	maxStars := previousMapLevels * 3
	requiredStars := int(math.Ceil(float64(maxStars) * 0.8))
	return playerStars >= requiredStars
}
