package services

import (
	"context"
	"fmt"
	"math"

	"github.com/NhomNhem/GameFeel-Backend/internal/database"
	"github.com/NhomNhem/GameFeel-Backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// TalentService handles talent upgrade logic
type TalentService struct{}

// NewTalentService creates a new talent service
func NewTalentService() *TalentService {
	return &TalentService{}
}

// GetTalentConfigs returns all talent configurations
// TODO: Store in database, for now use hardcoded values
func (s *TalentService) GetTalentConfigs() map[string]*models.TalentConfig {
	return map[string]*models.TalentConfig{
		models.TalentHealth: {
			ID:            models.TalentHealth,
			Name:          "Health",
			Description:   "Increases max HP",
			MaxLevel:      20,
			BaseCost:      500,
			CostScaling:   1.15, // Cost increases 15% per level
			BonusPerLevel: 5.0,  // +5% HP per level
			StatType:      "hp",
			UnlockMap:     1,
		},
		models.TalentDamage: {
			ID:            models.TalentDamage,
			Name:          "Damage",
			Description:   "Increases attack power",
			MaxLevel:      20,
			BaseCost:      500,
			CostScaling:   1.15,
			BonusPerLevel: 3.0, // +3% Damage per level
			StatType:      "damage",
			UnlockMap:     1,
		},
		models.TalentArmor: {
			ID:            models.TalentArmor,
			Name:          "Armor",
			Description:   "Reduces damage taken",
			MaxLevel:      20,
			BaseCost:      500,
			CostScaling:   1.15,
			BonusPerLevel: 2.0, // +2% Damage Reduction per level
			StatType:      "armor",
			UnlockMap:     1,
		},
		models.TalentSpeed: {
			ID:            models.TalentSpeed,
			Name:          "Speed",
			Description:   "Increases movement speed",
			MaxLevel:      20,
			BaseCost:      500,
			CostScaling:   1.15,
			BonusPerLevel: 2.0, // +2% Move Speed per level
			StatType:      "speed",
			UnlockMap:     1,
		},
		models.TalentDashCooldown: {
			ID:            models.TalentDashCooldown,
			Name:          "Dash Cooldown",
			Description:   "Reduces dash cooldown",
			MaxLevel:      10,
			BaseCost:      1000,
			CostScaling:   1.2,
			BonusPerLevel: 5.0, // -5% Cooldown per level
			StatType:      "dash_cdr",
			UnlockMap:     2, // Unlock at Map 2
		},
	}
}

// CalculateTalentCost calculates the cost to upgrade to next level
func (s *TalentService) CalculateTalentCost(config *models.TalentConfig, currentLevel int) int {
	if currentLevel >= config.MaxLevel {
		return 0 // Already max level
	}
	
	// Formula: baseCost * (scaling ^ currentLevel)
	cost := float64(config.BaseCost) * math.Pow(config.CostScaling, float64(currentLevel))
	return int(math.Round(cost))
}

// GetOrCreateUserTalent gets user's current talent level
func (s *TalentService) GetOrCreateUserTalent(ctx context.Context, userID uuid.UUID, talentID string) (*models.UserTalent, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database not connected")
	}

	var talent models.UserTalent
	err := database.Pool.QueryRow(ctx, `
		SELECT id, user_id, talent_id, current_level, upgraded_at
		FROM user_talents
		WHERE user_id = $1 AND talent_id = $2
	`, userID, talentID).Scan(
		&talent.ID, &talent.UserID, &talent.TalentID, &talent.CurrentLevel, &talent.UpgradedAt,
	)

	if err == pgx.ErrNoRows {
		// Create new talent at level 0
		talent.ID = uuid.New()
		talent.UserID = userID
		talent.TalentID = talentID
		talent.CurrentLevel = 0

		_, err = database.Pool.Exec(ctx, `
			INSERT INTO user_talents (id, user_id, talent_id, current_level, upgraded_at)
			VALUES ($1, $2, $3, 0, NOW())
		`, talent.ID, userID, talentID)

		if err != nil {
			return nil, fmt.Errorf("failed to create talent: %w", err)
		}

		return &talent, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get talent: %w", err)
	}

	return &talent, nil
}

// UpgradeTalent upgrades a talent for the user
func (s *TalentService) UpgradeTalent(ctx context.Context, userID uuid.UUID, talentID string) (*models.TalentUpgradeResponse, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database not connected")
	}

	// Get talent config
	configs := s.GetTalentConfigs()
	config, ok := configs[talentID]
	if !ok {
		return nil, fmt.Errorf("invalid talent ID: %s", talentID)
	}

	// Start transaction
	tx, err := database.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get user's current gold
	var currentGold int
	err = tx.QueryRow(ctx, `SELECT gold FROM users WHERE id = $1`, userID).Scan(&currentGold)
	if err != nil {
		return nil, fmt.Errorf("failed to get user gold: %w", err)
	}

	// Get or create talent
	talent, err := s.GetOrCreateUserTalent(ctx, userID, talentID)
	if err != nil {
		return nil, err
	}

	// Check if already max level
	if talent.CurrentLevel >= config.MaxLevel {
		return nil, fmt.Errorf("talent already at max level")
	}

	// Calculate upgrade cost
	upgradeCost := s.CalculateTalentCost(config, talent.CurrentLevel)

	// Check if user has enough gold
	if currentGold < upgradeCost {
		return nil, fmt.Errorf("insufficient gold: need %d, have %d", upgradeCost, currentGold)
	}

	// Deduct gold
	var newTotalGold int
	err = tx.QueryRow(ctx, `
		UPDATE users 
		SET gold = gold - $1 
		WHERE id = $2 
		RETURNING gold
	`, upgradeCost, userID).Scan(&newTotalGold)
	
	if err != nil {
		return nil, fmt.Errorf("failed to deduct gold: %w", err)
	}

	// Upgrade talent
	newLevel := talent.CurrentLevel + 1
	err = tx.QueryRow(ctx, `
		UPDATE user_talents 
		SET current_level = $1, upgraded_at = NOW() 
		WHERE id = $2
		RETURNING current_level
	`, newLevel, talent.ID).Scan(&newLevel)
	
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade talent: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Calculate stat bonus
	totalBonus := config.BonusPerLevel * float64(newLevel)
	statBonus := fmt.Sprintf("+%.0f%% %s", totalBonus, config.StatType)

	// Calculate next level cost (if not max)
	var nextLevelCost int
	if newLevel < config.MaxLevel {
		nextLevelCost = s.CalculateTalentCost(config, newLevel)
	}

	return &models.TalentUpgradeResponse{
		Success:       true,
		TalentID:      talentID,
		NewLevel:      newLevel,
		GoldSpent:     upgradeCost,
		NewTotalGold:  newTotalGold,
		StatBonus:     statBonus,
		NextLevelCost: nextLevelCost,
	}, nil
}

// GetUserTalents gets all talents for a user
func (s *TalentService) GetUserTalents(ctx context.Context, userID uuid.UUID) ([]models.UserTalent, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database not connected")
	}

	rows, err := database.Pool.Query(ctx, `
		SELECT id, user_id, talent_id, current_level, upgraded_at
		FROM user_talents
		WHERE user_id = $1
	`, userID)
	
	if err != nil {
		return nil, fmt.Errorf("failed to query talents: %w", err)
	}
	defer rows.Close()

	var talents []models.UserTalent
	for rows.Next() {
		var talent models.UserTalent
		err := rows.Scan(&talent.ID, &talent.UserID, &talent.TalentID, &talent.CurrentLevel, &talent.UpgradedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan talent: %w", err)
		}
		talents = append(talents, talent)
	}

	return talents, nil
}
