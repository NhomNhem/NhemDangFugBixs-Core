package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresLeaderboardRepository struct {
	db *pgxpool.Pool
}

// NewPostgresLeaderboardRepository creates a new PostgreSQL leaderboard repository
func NewPostgresLeaderboardRepository(db *pgxpool.Pool) repository.LeaderboardRepository {
	return &postgresLeaderboardRepository{db: db}
}

func (r *postgresLeaderboardRepository) GetByCriteria(ctx context.Context, lbType, scope, character string, limit, offset int) ([]models.HollowWildsLeaderboardEntry, int, error) {
	if r.db == nil {
		return []models.HollowWildsLeaderboardEntry{
			{
				Rank:        1,
				PlayerID:    "MOCK_PLAYER",
				DisplayName: "Mock Player",
				Value:       42,
				Character:   character,
				UpdatedAt:   time.Now().Format(time.RFC3339),
			},
		}, 1, nil
	}
	query := `
		SELECT 
			p.playfab_id,
			COALESCE(p.display_name, 'Anonymous'),
			le.value,
			le.character,
			le.world_seed,
			le.combat_build,
			le.updated_at,
			le.run_metadata,
			RANK() OVER (ORDER BY le.value DESC) as rank
		FROM leaderboard_entries le
		JOIN players p ON le.player_id = p.id
		WHERE le.type = $1
		  AND ($2 = 'global' OR le.character = $3)
		ORDER BY le.value DESC
		LIMIT $4 OFFSET $5
	`

	rows, err := r.db.Query(ctx, query, lbType, scope, character, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query leaderboard: %w", err)
	}
	defer rows.Close()

	var entries []models.HollowWildsLeaderboardEntry
	for rows.Next() {
		var entry models.HollowWildsLeaderboardEntry
		var updatedAt time.Time
		err := rows.Scan(
			&entry.PlayerID,
			&entry.DisplayName,
			&entry.Value,
			&entry.Character,
			&entry.WorldSeed,
			&entry.CombatBuild,
			&updatedAt,
			&entry.RunMetadata,
			&entry.Rank,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan entry: %w", err)
		}
		entry.UpdatedAt = updatedAt.Format("2006-01-02T15:04:05Z")
		entries = append(entries, entry)
	}

	// Get total
	var total int
	countQuery := `
		SELECT COUNT(*) FROM leaderboard_entries 
		WHERE type = $1 AND ($2 = 'global' OR character = $3)
	`
	err = r.db.QueryRow(ctx, countQuery, lbType, scope, character).Scan(&total)
	if err != nil {
		total = 0
	}

	return entries, total, nil
}

func (r *postgresLeaderboardRepository) GetPersonalBest(ctx context.Context, playerID uuid.UUID, lbType, character string) (int64, error) {
	if r.db == nil {
		return 0, nil
	}
	var currentBest int64
	err := r.db.QueryRow(ctx, `
		SELECT value FROM leaderboard_entries 
		WHERE player_id = $1 AND type = $2 AND character = $3
	`, playerID, lbType, character).Scan(&currentBest)

	if err == pgx.ErrNoRows {
		return 0, nil
	}
	return currentBest, err
}

func (r *postgresLeaderboardRepository) UpsertEntry(ctx context.Context, entry *models.HollowWildsLeaderboardEntry) error {
	if r.db == nil {
		return nil
	}
	playerUUID, _ := uuid.Parse(entry.PlayerID) // Internal ID expected here in real implementation
	_, err := r.db.Exec(ctx, `
		INSERT INTO leaderboard_entries (player_id, type, value, character, world_seed, combat_build, run_metadata, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (player_id, type, character) DO UPDATE
		SET value = EXCLUDED.value,
		    world_seed = EXCLUDED.world_seed,
		    combat_build = EXCLUDED.combat_build,
		    run_metadata = EXCLUDED.run_metadata,
		    updated_at = NOW()
	`, playerUUID, entry.Type, entry.Value, entry.Character, entry.WorldSeed, entry.CombatBuild, entry.RunMetadata)

	return err
}

func (r *postgresLeaderboardRepository) GetPlayerStats(ctx context.Context, playerID uuid.UUID) ([]models.PlayerLeaderboardEntry, error) {
	if r.db == nil {
		return []models.PlayerLeaderboardEntry{
			{
				Type:          "longest_run_days",
				GlobalRank:    10,
				CharacterRank: 5,
				Character:     "RIMBA",
				Value:         15,
				PersonalBest:  true,
			},
		}, nil
	}
	query := `
		WITH ranked_global AS (
			SELECT type, player_id, value, RANK() OVER (PARTITION BY type ORDER BY value DESC) as rank
			FROM leaderboard_entries
		),
		ranked_character AS (
			SELECT type, player_id, character, value, RANK() OVER (PARTITION BY type, character ORDER BY value DESC) as rank
			FROM leaderboard_entries
		)
		SELECT 
			le.type,
			rg.rank as global_rank,
			rc.rank as character_rank,
			le.character,
			le.value
		FROM leaderboard_entries le
		JOIN ranked_global rg ON le.player_id = rg.player_id AND le.type = rg.type
		JOIN ranked_character rc ON le.player_id = rc.player_id AND le.type = rc.type AND le.character = rc.character
		WHERE le.player_id = $1
	`

	rows, err := r.db.Query(ctx, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query player stats: %w", err)
	}
	defer rows.Close()

	var entries []models.PlayerLeaderboardEntry
	for rows.Next() {
		var entry models.PlayerLeaderboardEntry
		err := rows.Scan(
			&entry.Type,
			&entry.GlobalRank,
			&entry.CharacterRank,
			&entry.Character,
			&entry.Value,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan entry: %w", err)
		}
		entry.PersonalBest = true
		entries = append(entries, entry)
	}

	return entries, nil
}

// Legacy Implementation

func (r *postgresLeaderboardRepository) GetLegacyGlobal(ctx context.Context, levelID string, limit, offset int) ([]models.LeaderboardEntry, int, error) {
	if r.db == nil {
		return nil, 0, nil
	}

	query := `
		SELECT 
			u.id,
			COALESCE(u.display_name, 'Anonymous'),
			lr.best_time_seconds,
			lr.stars_earned,
			RANK() OVER (ORDER BY lr.best_time_seconds ASC) as rank
		FROM level_rankings lr
		JOIN users u ON lr.user_id = u.id
		WHERE lr.level_id = $1
		ORDER BY lr.best_time_seconds ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, levelID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var entries []models.LeaderboardEntry
	for rows.Next() {
		var entry models.LeaderboardEntry
		var userID uuid.UUID
		err := rows.Scan(
			&userID,
			&entry.DisplayName,
			&entry.BestTime,
			&entry.Stars,
			&entry.Rank,
		)
		if err != nil {
			return nil, 0, err
		}
		entry.PlayerID = userID.String()
		entries = append(entries, entry)
	}

	var total int
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM level_rankings WHERE level_id = $1", levelID).Scan(&total)
	if err != nil {
		total = 0
	}

	return entries, total, nil
}

func (r *postgresLeaderboardRepository) GetLegacyPlayerRank(ctx context.Context, userID uuid.UUID, levelID string) (int, float64, int, error) {
	if r.db == nil {
		return 0, 0, 0, nil
	}

	query := `
		WITH ranked AS (
			SELECT 
				user_id,
				best_time_seconds,
				stars_earned,
				RANK() OVER (ORDER BY best_time_seconds ASC) as rank
			FROM level_rankings
			WHERE level_id = $1
		)
		SELECT rank, best_time_seconds, stars_earned
		FROM ranked
		WHERE user_id = $2
	`

	var rank int
	var bestTime float64
	var stars int
	err := r.db.QueryRow(ctx, query, levelID, userID).Scan(&rank, &bestTime, &stars)
	if err == pgx.ErrNoRows {
		return 0, 0, 0, nil
	}
	if err != nil {
		return 0, 0, 0, err
	}

	return rank, bestTime, stars, nil
}

func (r *postgresLeaderboardRepository) UpsertLegacyEntry(ctx context.Context, userID uuid.UUID, levelID string, timeSeconds float64, stars int) error {
	if r.db == nil {
		return nil
	}

	query := `
		INSERT INTO level_rankings (user_id, level_id, best_time_seconds, stars_earned, completed_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (user_id, level_id) DO UPDATE
		SET best_time_seconds = LEAST(level_rankings.best_time_seconds, EXCLUDED.best_time_seconds),
		    stars_earned = GREATEST(level_rankings.stars_earned, EXCLUDED.stars_earned),
		    updated_at = NOW()
	`

	_, err := r.db.Exec(ctx, query, userID, levelID, timeSeconds, stars)
	return err
}

func (r *postgresLeaderboardRepository) GetLegacyFriends(ctx context.Context, friendIDs []string, levelID string) ([]models.LeaderboardEntry, error) {
	if r.db == nil {
		return nil, nil
	}

	if len(friendIDs) == 0 {
		return []models.LeaderboardEntry{}, nil
	}

	query := `
		SELECT 
			u.id,
			COALESCE(u.display_name, 'Anonymous'),
			lr.best_time_seconds,
			lr.stars_earned,
			RANK() OVER (ORDER BY lr.best_time_seconds ASC) as rank
		FROM level_rankings lr
		JOIN users u ON lr.user_id = u.id
		WHERE lr.level_id = $1
		  AND u.playfab_id = ANY($2)
		ORDER BY lr.best_time_seconds ASC
	`

	rows, err := r.db.Query(ctx, query, levelID, friendIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.LeaderboardEntry
	for rows.Next() {
		var entry models.LeaderboardEntry
		var userID uuid.UUID
		err := rows.Scan(
			&userID,
			&entry.DisplayName,
			&entry.BestTime,
			&entry.Stars,
			&entry.Rank,
		)
		if err != nil {
			return nil, err
		}
		entry.PlayerID = userID.String()
		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *postgresLeaderboardRepository) ResetLegacyLeaderboard(ctx context.Context, levelID string) error {
	if r.db == nil {
		return nil
	}
	_, err := r.db.Exec(ctx, "DELETE FROM level_rankings WHERE level_id = $1", levelID)
	return err
}

func (r *postgresLeaderboardRepository) GetLegacyStats(ctx context.Context) (*models.LeaderboardStatsResponse, error) {
	if r.db == nil {
		return &models.LeaderboardStatsResponse{
			TotalEntries:  100,
			UniquePlayers: 50,
			LevelsTracked: 5,
			LastUpdated:   time.Now(),
		}, nil
	}

	stats := &models.LeaderboardStatsResponse{
		LastUpdated: time.Now(),
		LevelStats:  []models.LevelStatsInfo{},
	}

	// Basic stats
	err := r.db.QueryRow(ctx, `
		SELECT 
			COUNT(*), 
			COUNT(DISTINCT user_id), 
			COUNT(DISTINCT level_id),
			AVG(best_time_seconds)
		FROM level_rankings
	`).Scan(&stats.TotalEntries, &stats.UniquePlayers, &stats.LevelsTracked, &stats.AverageTime)
	if err != nil {
		// Table may not exist yet — return empty stats
		stats.LevelStats = []models.LevelStatsInfo{}
		return stats, nil
	}

	// Per-level stats
	rows, err := r.db.Query(ctx, `
		SELECT 
			level_id, 
			COUNT(*), 
			COUNT(DISTINCT user_id), 
			AVG(best_time_seconds),
			MIN(best_time_seconds)
		FROM level_rankings
		GROUP BY level_id
		ORDER BY COUNT(*) DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var info models.LevelStatsInfo
		err := rows.Scan(&info.LevelID, &info.TotalEntries, &info.UniquePlayers, &info.AverageTime, &info.BestTime)
		if err != nil {
			continue
		}
		stats.LevelStats = append(stats.LevelStats, info)
	}

	if len(stats.LevelStats) > 0 {
		stats.TopLevelID = stats.LevelStats[0].LevelID
		stats.TopLevelCount = stats.LevelStats[0].TotalEntries
	}

	return stats, nil
}
