package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresTalentRepository struct {
	db         *pgxpool.Pool
	cache      map[string]*models.TalentConfig
	cacheMutex sync.RWMutex
	lastFetch  time.Time
}

// NewPostgresTalentRepository creates a new PostgreSQL talent repository
func NewPostgresTalentRepository(db *pgxpool.Pool) repository.TalentRepository {
	return &postgresTalentRepository{
		db:    db,
		cache: make(map[string]*models.TalentConfig),
	}
}

func (r *postgresTalentRepository) GetConfigs(ctx context.Context) (map[string]*models.TalentConfig, error) {
	r.cacheMutex.RLock()
	if len(r.cache) > 0 && time.Since(r.lastFetch) < 5*time.Minute {
		r.cacheMutex.RUnlock()
		return r.cache, nil
	}
	r.cacheMutex.RUnlock()

	if r.db == nil {
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
		}, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT talent_id, config_json
		FROM talent_configs
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query talent configs: %w", err)
	}
	defer rows.Close()

	newConfigs := make(map[string]*models.TalentConfig)
	for rows.Next() {
		var talentID string
		var configJSON []byte
		if err := rows.Scan(&talentID, &configJSON); err != nil {
			return nil, fmt.Errorf("failed to scan talent config: %w", err)
		}

		var config models.TalentConfig
		if err := json.Unmarshal(configJSON, &config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal talent config: %w", err)
		}
		config.ID = talentID
		newConfigs[talentID] = &config
	}

	r.cacheMutex.Lock()
	r.cache = newConfigs
	r.lastFetch = time.Now()
	r.cacheMutex.Unlock()

	return newConfigs, nil
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
