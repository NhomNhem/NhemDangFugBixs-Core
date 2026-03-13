package persistence

import (
	"context"
	"fmt"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresTalentRepository struct {
	db *pgxpool.Pool
}

// NewPostgresTalentRepository creates a new PostgreSQL talent repository
func NewPostgresTalentRepository(db *pgxpool.Pool) repository.TalentRepository {
	return &postgresTalentRepository{db: db}
}

func (r *postgresTalentRepository) GetConfigs() map[string]*models.TalentConfig {
	// TODO: Move to database, for now keep the logic from TalentService
	return map[string]*models.TalentConfig{
		models.TalentHealth: {
			ID:            models.TalentHealth,
			Name:          "Health",
			Description:   "Increases max HP",
			MaxLevel:      20,
			BaseCost:      500,
			CostScaling:   1.15,
			BonusPerLevel: 5.0,
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
			BonusPerLevel: 3.0,
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
			BonusPerLevel: 2.0,
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
			BonusPerLevel: 2.0,
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
			BonusPerLevel: 5.0,
			StatType:      "dash_cdr",
			UnlockMap:     2,
		},
	}
}

func (r *postgresTalentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.UserTalent, error) {
	rows, err := r.db.Query(ctx, `
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

func (r *postgresTalentRepository) GetOrCreate(ctx context.Context, userID uuid.UUID, talentID string) (*models.UserTalent, error) {
	var talent models.UserTalent
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, talent_id, current_level, upgraded_at
		FROM user_talents
		WHERE user_id = $1 AND talent_id = $2
	`, userID, talentID).Scan(
		&talent.ID, &talent.UserID, &talent.TalentID, &talent.CurrentLevel, &talent.UpgradedAt,
	)

	if err == pgx.ErrNoRows {
		talent.ID = uuid.New()
		talent.UserID = userID
		talent.TalentID = talentID
		talent.CurrentLevel = 0

		_, err = r.db.Exec(ctx, `
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

func (r *postgresTalentRepository) UpdateLevel(ctx context.Context, id uuid.UUID, newLevel int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE user_talents
		SET current_level = $1, upgraded_at = NOW()
		WHERE id = $2
	`, newLevel, id)

	if err != nil {
		return fmt.Errorf("failed to update talent level: %w", err)
	}
	return nil
}

func (r *postgresTalentRepository) UpdateUserGold(ctx context.Context, userID uuid.UUID, goldChange int) (int, error) {
	var newTotalGold int
	err := r.db.QueryRow(ctx, `
		UPDATE users
		SET gold = gold + $1
		WHERE id = $2
		RETURNING gold
	`, goldChange, userID).Scan(&newTotalGold)

	if err != nil {
		return 0, fmt.Errorf("failed to update user gold: %w", err)
	}
	return newTotalGold, nil
}
