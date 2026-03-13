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

type postgresLevelRepository struct {
	db *pgxpool.Pool
}

// NewPostgresLevelRepository creates a new PostgreSQL level repository
func NewPostgresLevelRepository(db *pgxpool.Pool) repository.LevelRepository {
	return &postgresLevelRepository{db: db}
}

func (r *postgresLevelRepository) GetConfig(levelID string, mapID string) (*models.LevelConfig, error) {
	// TODO: Move to database, for now keep the logic from LevelService
	return &models.LevelConfig{
		LevelID:        levelID,
		MapID:          mapID,
		MinTimeSeconds: 10.0,
		BaseGold:       60,
		Objectives: []models.LevelObjective{
			{Type: "completion", Threshold: 1, Operator: "gte"},
			{Type: "health", Threshold: 50, Operator: "gte"},
			{Type: "time", Threshold: 60, Operator: "lte"},
		},
	}, nil
}

func (r *postgresLevelRepository) GetCompletion(ctx context.Context, userID uuid.UUID, levelID string) (*models.LevelCompletion, error) {
	var completion models.LevelCompletion
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, level_id, map_id, stars_earned, best_time_seconds, play_count,
		       last_final_hp, last_dash_count, last_counter_count, last_vulnerable_kills,
		       first_completed_at, last_played_at
		FROM level_completions
		WHERE user_id = $1 AND level_id = $2
	`, userID, levelID).Scan(
		&completion.ID,
		&completion.UserID,
		&completion.LevelID,
		&completion.MapID,
		&completion.StarsEarned,
		&completion.BestTimeSeconds,
		&completion.PlayCount,
		&completion.LastFinalHP,
		&completion.LastDashCount,
		&completion.LastCounterCount,
		&completion.LastVulnerableKills,
		&completion.FirstCompletedAt,
		&completion.LastPlayedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get completion: %w", err)
	}

	return &completion, nil
}

func (r *postgresLevelRepository) CreateCompletion(ctx context.Context, completion *models.LevelCompletion) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO level_completions (id, user_id, level_id, map_id, stars_earned, best_time_seconds, play_count,
		                               last_final_hp, last_dash_count, last_counter_count, last_vulnerable_kills,
		                               first_completed_at, last_played_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, completion.ID, completion.UserID, completion.LevelID, completion.MapID, completion.StarsEarned, completion.BestTimeSeconds,
		completion.PlayCount, completion.LastFinalHP, completion.LastDashCount, completion.LastCounterCount, completion.LastVulnerableKills,
		completion.FirstCompletedAt, completion.LastPlayedAt)

	if err != nil {
		return fmt.Errorf("failed to create completion: %w", err)
	}
	return nil
}

func (r *postgresLevelRepository) UpdateCompletion(ctx context.Context, completion *models.LevelCompletion) error {
	_, err := r.db.Exec(ctx, `
		UPDATE level_completions
		SET stars_earned = $1,
		    best_time_seconds = $2,
		    play_count = $3,
		    last_final_hp = $4,
		    last_dash_count = $5,
		    last_counter_count = $6,
		    last_vulnerable_kills = $7,
		    last_played_at = $8
		WHERE id = $9
	`, completion.StarsEarned, completion.BestTimeSeconds, completion.PlayCount, completion.LastFinalHP, completion.LastDashCount,
		completion.LastCounterCount, completion.LastVulnerableKills, completion.LastPlayedAt, completion.ID)

	if err != nil {
		return fmt.Errorf("failed to update completion: %w", err)
	}
	return nil
}

func (r *postgresLevelRepository) UpdateUserStats(ctx context.Context, userID uuid.UUID, goldEarned, starsEarned int) (int, int, error) {
	var newTotalGold, newTotalStars int
	err := r.db.QueryRow(ctx, `
		UPDATE users
		SET gold = gold + $1,
		    total_stars_collected = total_stars_collected + $2,
		    last_played_at = NOW()
		WHERE id = $3
		RETURNING gold, total_stars_collected
	`, goldEarned, starsEarned, userID).Scan(&newTotalGold, &newTotalStars)

	if err != nil {
		return 0, 0, fmt.Errorf("failed to update user stats: %w", err)
	}
	return newTotalGold, newTotalStars, nil
}
