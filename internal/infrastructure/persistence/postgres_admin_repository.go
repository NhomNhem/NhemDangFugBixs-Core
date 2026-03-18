package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresAdminRepository struct {
	db *pgxpool.Pool
}

// NewPostgresAdminRepository creates a new PostgreSQL admin repository
func NewPostgresAdminRepository(db *pgxpool.Pool) repository.AdminRepository {
	return &postgresAdminRepository{db: db}
}

func (r *postgresAdminRepository) SearchUsers(ctx context.Context, query string, page, perPage int) (*models.UserSearchResponse, error) {
	offset := (page - 1) * perPage
	searchPattern := "%" + query + "%"

	var totalCount int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM users
		WHERE playfab_id ILIKE $1
		   OR COALESCE(display_name, '') ILIKE $1
	`, searchPattern).Scan(&totalCount)

	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT
			u.id,
			u.playfab_id,
			'' as email,
			COALESCE(u.display_name, '') as username,
			u.gold,
			u.total_stars_collected as total_stars,
			u.is_admin,
			u.created_at,
			u.last_login_at,
			COALESCE(
				(SELECT COUNT(DISTINCT level_id) FROM level_completions WHERE user_id = u.id),
				0
			) as levels_completed,
			EXISTS(SELECT 1 FROM user_bans WHERE user_id = u.id AND is_active = true) as is_banned
		FROM users u
		WHERE u.playfab_id ILIKE $1
		   OR u.display_name ILIKE $1
		ORDER BY u.created_at DESC
		LIMIT $2 OFFSET $3
	`, searchPattern, perPage, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []models.UserProfile
	for rows.Next() {
		var user models.UserProfile
		err := rows.Scan(
			&user.ID,
			&user.PlayFabID,
			&user.Email,
			&user.Username,
			&user.TotalGold,
			&user.TotalStars,
			&user.IsAdmin,
			&user.CreatedAt,
			&user.LastLoginAt,
			&user.LevelsCompleted,
			&user.IsBanned,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return &models.UserSearchResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    perPage,
	}, nil
}

func (r *postgresAdminRepository) GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := r.db.QueryRow(ctx, `
		SELECT
			u.id,
			u.playfab_id,
			'' as email,
			COALESCE(u.display_name, '') as username,
			u.gold,
			u.total_stars_collected as total_stars,
			u.is_admin,
			u.created_at,
			u.last_login_at,
			COALESCE(
				(SELECT COUNT(DISTINCT level_id) FROM level_completions WHERE user_id = u.id),
				0
			) as levels_completed,
			EXISTS(SELECT 1 FROM user_bans WHERE user_id = u.id AND is_active = true) as is_banned
		FROM users u
		WHERE u.id = $1
	`, userID).Scan(
		&profile.ID,
		&profile.PlayFabID,
		&profile.Email,
		&profile.Username,
		&profile.TotalGold,
		&profile.TotalStars,
		&profile.IsAdmin,
		&profile.CreatedAt,
		&profile.LastLoginAt,
		&profile.LevelsCompleted,
		&profile.IsBanned,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return &profile, nil
}

func (r *postgresAdminRepository) AdjustGold(ctx context.Context, adminID, userID uuid.UUID, amount int, reason, ipAddress string) (*models.AdjustGoldResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var oldBalance int
	err = tx.QueryRow(ctx, "SELECT gold FROM users WHERE id = $1", userID).Scan(&oldBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to get current balance: %w", err)
	}

	newBalance := oldBalance + amount
	if newBalance < 0 {
		return nil, fmt.Errorf("insufficient gold: cannot reduce below 0")
	}

	_, err = tx.Exec(ctx, "UPDATE users SET gold = $1 WHERE id = $2", newBalance, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update gold: %w", err)
	}

	details := map[string]any{
		"reason":      reason,
		"amount":      amount,
		"old_balance": oldBalance,
		"new_balance": newBalance,
	}
	detailsJSON, _ := json.Marshal(details)

	var actionID uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO admin_actions (admin_id, action_type, target_user_id, details, ip_address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, adminID, "adjust_gold", userID, detailsJSON, ipAddress).Scan(&actionID)

	if err != nil {
		return nil, fmt.Errorf("failed to log admin action: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.AdjustGoldResponse{
		UserID:     userID,
		OldBalance: oldBalance,
		NewBalance: newBalance,
		Adjustment: amount,
		Reason:     reason,
		ActionID:   actionID,
	}, nil
}

func (r *postgresAdminRepository) BanUser(ctx context.Context, adminID, userID uuid.UUID, reason string, bannedUntil *time.Time, ipAddress string) (*models.BanUserResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "UPDATE user_bans SET is_active = false WHERE user_id = $1 AND is_active = true", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to deactivate existing bans: %w", err)
	}

	var banID uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO user_bans (user_id, banned_by, reason, banned_until, is_active)
		VALUES ($1, $2, $3, $4, true)
		RETURNING id
	`, userID, adminID, reason, bannedUntil).Scan(&banID)

	if err != nil {
		return nil, fmt.Errorf("failed to create ban: %w", err)
	}

	details := map[string]any{
		"reason":       reason,
		"banned_until": bannedUntil,
	}
	detailsJSON, _ := json.Marshal(details)

	_, err = tx.Exec(ctx, `
		INSERT INTO admin_actions (admin_id, action_type, target_user_id, details, ip_address)
		VALUES ($1, $2, $3, $4, $5)
	`, adminID, "ban_user", userID, detailsJSON, ipAddress)

	if err != nil {
		return nil, fmt.Errorf("failed to log admin action: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.BanUserResponse{
		BanID:       banID,
		UserID:      userID,
		Reason:      reason,
		BannedUntil: bannedUntil,
		IsPermanent: bannedUntil == nil,
	}, nil
}

func (r *postgresAdminRepository) UnbanUser(ctx context.Context, adminID, userID uuid.UUID, reason, ipAddress string) (*models.UnbanUserResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, `
		UPDATE user_bans
		SET is_active = false,
		    unbanned_at = CURRENT_TIMESTAMP,
		    unbanned_by = $1,
		    unban_reason = $2
		WHERE user_id = $3 AND is_active = true
	`, adminID, reason, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to unban user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("no active ban found for user")
	}

	details := map[string]any{
		"reason": reason,
	}
	detailsJSON, _ := json.Marshal(details)

	_, err = tx.Exec(ctx, `
		INSERT INTO admin_actions (admin_id, action_type, target_user_id, details, ip_address)
		VALUES ($1, $2, $3, $4, $5)
	`, adminID, "unban_user", userID, detailsJSON, ipAddress)

	if err != nil {
		return nil, fmt.Errorf("failed to log admin action: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.UnbanUserResponse{
		UserID:      userID,
		UnbannedAt:  time.Now(),
		UnbanReason: reason,
	}, nil
}

func (r *postgresAdminRepository) GetBanHistory(ctx context.Context, userID uuid.UUID) ([]models.UserBan, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, banned_by, reason, banned_until, is_active,
		       created_at, unbanned_at, unbanned_by, COALESCE(unban_reason, '')
		FROM user_bans
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to get ban history: %w", err)
	}
	defer rows.Close()

	var bans []models.UserBan
	for rows.Next() {
		var ban models.UserBan
		err := rows.Scan(
			&ban.ID,
			&ban.UserID,
			&ban.BannedBy,
			&ban.Reason,
			&ban.BannedUntil,
			&ban.IsActive,
			&ban.CreatedAt,
			&ban.UnbannedAt,
			&ban.UnbannedBy,
			&ban.UnbanReason,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ban: %w", err)
		}
		bans = append(bans, ban)
	}

	return bans, nil
}

func (r *postgresAdminRepository) GetActions(ctx context.Context, page, perPage int) (*models.AdminActionsResponse, error) {
	offset := (page - 1) * perPage

	var totalCount int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM admin_actions").Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count admin actions: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, admin_id, action_type, target_user_id, details,
		       COALESCE(ip_address, ''), created_at
		FROM admin_actions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, perPage, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get admin actions: %w", err)
	}
	defer rows.Close()

	var actions []models.AdminAction
	for rows.Next() {
		var action models.AdminAction
		var detailsJSON []byte
		err := rows.Scan(
			&action.ID,
			&action.AdminID,
			&action.ActionType,
			&action.TargetUserID,
			&detailsJSON,
			&action.IPAddress,
			&action.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan admin action: %w", err)
		}

		if err := json.Unmarshal(detailsJSON, &action.Details); err != nil {
			action.Details = map[string]any{}
		}

		actions = append(actions, action)
	}

	return &models.AdminActionsResponse{
		Actions:    actions,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    perPage,
	}, nil
}

func (r *postgresAdminRepository) GetSystemStats(ctx context.Context) (*models.SystemStatsResponse, error) {
	stats := &models.SystemStatsResponse{
		LastUpdated: time.Now(),
	}

	r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE last_login_at > NOW() - INTERVAL '24 hours'`).Scan(&stats.ActiveToday)
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE last_login_at > NOW() - INTERVAL '7 days'`).Scan(&stats.ActiveThisWeek)
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE last_login_at > NOW() - INTERVAL '30 days'`).Scan(&stats.ActiveThisMonth)
	r.db.QueryRow(ctx, `SELECT COUNT(DISTINCT user_id) FROM user_bans WHERE is_active = true`).Scan(&stats.BannedUsers)
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM admin_actions WHERE created_at > NOW() - INTERVAL '24 hours'`).Scan(&stats.AdminActionsToday)
	r.db.QueryRow(ctx, "SELECT COALESCE(SUM(gold), 0) FROM users").Scan(&stats.TotalGoldInGame)
	r.db.QueryRow(ctx, "SELECT COALESCE(SUM(total_stars_collected), 0) FROM users").Scan(&stats.TotalStarsEarned)
	r.db.QueryRow(ctx, "SELECT COUNT(*) FROM level_completions").Scan(&stats.LevelsCompleted)

	return stats, nil
}

func (r *postgresAdminRepository) GetLevelCompletions(ctx context.Context, userID uuid.UUID) ([]map[string]any, error) {
	var levelCompletions []map[string]any
	rows, err := r.db.Query(ctx, `
		SELECT level_id, map_id, stars_earned, best_time_seconds, last_played_at
		FROM level_completions
		WHERE user_id = $1
		ORDER BY last_played_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get level completions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var levelID, mapID string
		var starsEarned int
		var bestTime float64
		var lastPlayedAt time.Time
		rows.Scan(&levelID, &mapID, &starsEarned, &bestTime, &lastPlayedAt)
		levelCompletions = append(levelCompletions, map[string]any{
			"level_id":          levelID,
			"map_id":            mapID,
			"stars_earned":      starsEarned,
			"best_time_seconds": bestTime,
			"last_played_at":    lastPlayedAt,
		})
	}
	return levelCompletions, nil
}

func (r *postgresAdminRepository) GetUserTalents(ctx context.Context, userID uuid.UUID) ([]map[string]any, error) {
	var talents []map[string]any
	rows, err := r.db.Query(ctx, `
		SELECT talent_id, current_level, upgraded_at
		FROM user_talents
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user talents: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var talentID string
		var currentLevel int
		var upgradedAt time.Time
		rows.Scan(&talentID, &currentLevel, &upgradedAt)
		talents = append(talents, map[string]any{
			"talent_id":     talentID,
			"current_level": currentLevel,
			"upgraded_at":   upgradedAt,
		})
	}
	return talents, nil
}

func (r *postgresAdminRepository) LogAction(ctx context.Context, adminID uuid.UUID, actionType string, targetUserID *uuid.UUID, details map[string]any, ipAddress string) error {
	detailsJSON, _ := json.Marshal(details)
	_, err := r.db.Exec(ctx, `
		INSERT INTO admin_actions (admin_id, action_type, target_user_id, details, ip_address)
		VALUES ($1, $2, $3, $4, $5)
	`, adminID, actionType, targetUserID, detailsJSON, ipAddress)
	return err
}

func (r *postgresAdminRepository) GetAnalyticsSummary(ctx context.Context) (*models.AnalyticsSummaryResponse, error) {
	summary := &models.AnalyticsSummaryResponse{
		LastUpdated: time.Now(),
		TopEvents:   []models.EventCount{},
		DAULast7d:   []models.DAUDataPoint{},
	}

	// Total events last 24h
	r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM analytics_events WHERE created_at > NOW() - INTERVAL '24 hours'
	`).Scan(&summary.TotalEventsLast24h)

	// Total events last 7d
	r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM analytics_events WHERE created_at > NOW() - INTERVAL '7 days'
	`).Scan(&summary.TotalEventsLast7d)

	// Top 10 event types last 7d
	rows, err := r.db.Query(ctx, `
		SELECT event_type, COUNT(*) as cnt
		FROM analytics_events
		WHERE created_at > NOW() - INTERVAL '7 days'
		GROUP BY event_type
		ORDER BY cnt DESC
		LIMIT 10
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ec models.EventCount
			if rows.Scan(&ec.EventName, &ec.Count) == nil {
				summary.TopEvents = append(summary.TopEvents, ec)
			}
		}
	}

	// DAU last 7 days
	dauRows, err := r.db.Query(ctx, `
		SELECT DATE(created_at) as day, COUNT(DISTINCT user_id) as dau
		FROM analytics_events
		WHERE created_at > NOW() - INTERVAL '7 days'
		GROUP BY day
		ORDER BY day ASC
	`)
	if err == nil {
		defer dauRows.Close()
		for dauRows.Next() {
			var dp models.DAUDataPoint
			if dauRows.Scan(&dp.Date, &dp.Count) == nil {
				summary.DAULast7d = append(summary.DAULast7d, dp)
			}
		}
	}

	return summary, nil
}

func (r *postgresAdminRepository) GetActionsFiltered(ctx context.Context, page, perPage int, actionType string) (*models.AdminActionsResponse, error) {
	offset := (page - 1) * perPage

	var totalCount int
	var countErr error
	if actionType != "" {
		countErr = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM admin_actions WHERE action_type = $1", actionType).Scan(&totalCount)
	} else {
		countErr = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM admin_actions").Scan(&totalCount)
	}
	if countErr != nil {
		return nil, fmt.Errorf("failed to count admin actions: %w", countErr)
	}

	var actions []models.AdminAction
	var queryErr error

	if actionType != "" {
		rows, err := r.db.Query(ctx, `
			SELECT id, admin_id, action_type, target_user_id, details,
			       COALESCE(ip_address, ''), created_at
			FROM admin_actions
			WHERE action_type = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`, actionType, perPage, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to get admin actions: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var action models.AdminAction
			var detailsJSON []byte
			if err := rows.Scan(&action.ID, &action.AdminID, &action.ActionType, &action.TargetUserID, &detailsJSON, &action.IPAddress, &action.CreatedAt); err != nil {
				return nil, fmt.Errorf("failed to scan admin action: %w", err)
			}
			if err := json.Unmarshal(detailsJSON, &action.Details); err != nil {
				action.Details = map[string]any{}
			}
			actions = append(actions, action)
		}
		queryErr = rows.Err()
	} else {
		rows, err := r.db.Query(ctx, `
			SELECT id, admin_id, action_type, target_user_id, details,
			       COALESCE(ip_address, ''), created_at
			FROM admin_actions
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`, perPage, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to get admin actions: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var action models.AdminAction
			var detailsJSON []byte
			if err := rows.Scan(&action.ID, &action.AdminID, &action.ActionType, &action.TargetUserID, &detailsJSON, &action.IPAddress, &action.CreatedAt); err != nil {
				return nil, fmt.Errorf("failed to scan admin action: %w", err)
			}
			if err := json.Unmarshal(detailsJSON, &action.Details); err != nil {
				action.Details = map[string]any{}
			}
			actions = append(actions, action)
		}
		queryErr = rows.Err()
	}

	if queryErr != nil {
		return nil, fmt.Errorf("failed to iterate admin actions: %w", queryErr)
	}

	return &models.AdminActionsResponse{
		Actions:    actions,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    perPage,
	}, nil
}

// Level Config methods

func (r *postgresAdminRepository) ListLevelConfigs(ctx context.Context, page, perPage int) (*models.LevelConfigListResponse, error) {
	offset := (page - 1) * perPage

	var totalCount int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM level_configs").Scan(&totalCount); err != nil {
		// Table may not exist yet — return empty list
		return &models.LevelConfigListResponse{Levels: []models.AdminLevelConfig{}, TotalCount: 0, Page: page, PerPage: perPage}, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT level_id, map_id, name, difficulty, min_time_seconds,
		       base_gold, reward_stars, objectives, is_active, created_at, updated_at
		FROM level_configs
		ORDER BY level_id
		LIMIT $1 OFFSET $2
	`, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list level configs: %w", err)
	}
	defer rows.Close()

	var levels []models.AdminLevelConfig
	for rows.Next() {
		var lc models.AdminLevelConfig
		var objectivesJSON []byte
		if err := rows.Scan(
			&lc.LevelID, &lc.MapID, &lc.Name, &lc.Difficulty, &lc.MinTimeSeconds,
			&lc.BaseGold, &lc.RewardStars, &objectivesJSON, &lc.IsActive, &lc.CreatedAt, &lc.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan level config: %w", err)
		}
		if err := json.Unmarshal(objectivesJSON, &lc.Objectives); err != nil {
			lc.Objectives = []models.LevelObjective{}
		}
		levels = append(levels, lc)
	}

	return &models.LevelConfigListResponse{
		Levels:     levels,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    perPage,
	}, nil
}

func (r *postgresAdminRepository) GetLevelConfig(ctx context.Context, levelID string) (*models.AdminLevelConfig, error) {
	var lc models.AdminLevelConfig
	var objectivesJSON []byte
	err := r.db.QueryRow(ctx, `
		SELECT level_id, map_id, name, difficulty, min_time_seconds,
		       base_gold, reward_stars, objectives, is_active, created_at, updated_at
		FROM level_configs WHERE level_id = $1
	`, levelID).Scan(
		&lc.LevelID, &lc.MapID, &lc.Name, &lc.Difficulty, &lc.MinTimeSeconds,
		&lc.BaseGold, &lc.RewardStars, &objectivesJSON, &lc.IsActive, &lc.CreatedAt, &lc.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("level config not found: %s", levelID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get level config: %w", err)
	}
	if err := json.Unmarshal(objectivesJSON, &lc.Objectives); err != nil {
		lc.Objectives = []models.LevelObjective{}
	}
	return &lc, nil
}

func (r *postgresAdminRepository) CreateLevelConfig(ctx context.Context, config *models.AdminLevelConfig) error {
	objectivesJSON, err := json.Marshal(config.Objectives)
	if err != nil {
		return fmt.Errorf("failed to marshal objectives: %w", err)
	}
	_, err = r.db.Exec(ctx, `
		INSERT INTO level_configs (level_id, map_id, name, difficulty, min_time_seconds,
		                           base_gold, reward_stars, objectives, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
	`, config.LevelID, config.MapID, config.Name, config.Difficulty, config.MinTimeSeconds,
		config.BaseGold, config.RewardStars, objectivesJSON, config.IsActive)
	if err != nil {
		return fmt.Errorf("failed to create level config: %w", err)
	}
	return nil
}

func (r *postgresAdminRepository) UpdateLevelConfig(ctx context.Context, config *models.AdminLevelConfig) error {
	objectivesJSON, err := json.Marshal(config.Objectives)
	if err != nil {
		return fmt.Errorf("failed to marshal objectives: %w", err)
	}
	result, err := r.db.Exec(ctx, `
		UPDATE level_configs
		SET map_id = $2, name = $3, difficulty = $4, min_time_seconds = $5,
		    base_gold = $6, reward_stars = $7, objectives = $8, is_active = $9, updated_at = NOW()
		WHERE level_id = $1
	`, config.LevelID, config.MapID, config.Name, config.Difficulty, config.MinTimeSeconds,
		config.BaseGold, config.RewardStars, objectivesJSON, config.IsActive)
	if err != nil {
		return fmt.Errorf("failed to update level config: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("level config not found: %s", config.LevelID)
	}
	return nil
}

func (r *postgresAdminRepository) DeleteLevelConfig(ctx context.Context, levelID string) error {
	result, err := r.db.Exec(ctx, "DELETE FROM level_configs WHERE level_id = $1", levelID)
	if err != nil {
		return fmt.Errorf("failed to delete level config: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("level config not found: %s", levelID)
	}
	return nil
}

func (r *postgresAdminRepository) LevelHasLeaderboardEntries(ctx context.Context, levelID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM leaderboard_entries WHERE level_id = $1
	`, levelID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count leaderboard entries: %w", err)
	}
	return count, nil
}

// Talent Config methods

func (r *postgresAdminRepository) ListTalentConfigs(ctx context.Context, page, perPage int) (*models.TalentConfigListResponse, error) {
	offset := (page - 1) * perPage

	var totalCount int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM talent_configs").Scan(&totalCount); err != nil {
		// Table may not exist yet — return empty list
		return &models.TalentConfigListResponse{Talents: []models.AdminTalentConfig{}, TotalCount: 0, Page: page, PerPage: perPage}, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT talent_id, name, description, max_level, base_cost,
		       cost_scaling, bonus_per_level, stat_type, unlock_map, is_active, created_at, updated_at
		FROM talent_configs
		ORDER BY talent_id
		LIMIT $1 OFFSET $2
	`, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list talent configs: %w", err)
	}
	defer rows.Close()

	var talents []models.AdminTalentConfig
	for rows.Next() {
		var tc models.AdminTalentConfig
		if err := rows.Scan(
			&tc.TalentID, &tc.Name, &tc.Description, &tc.MaxLevel, &tc.BaseCost,
			&tc.CostScaling, &tc.BonusPerLevel, &tc.StatType, &tc.UnlockMap, &tc.IsActive, &tc.CreatedAt, &tc.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan talent config: %w", err)
		}
		talents = append(talents, tc)
	}

	return &models.TalentConfigListResponse{
		Talents:    talents,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    perPage,
	}, nil
}

func (r *postgresAdminRepository) GetTalentConfig(ctx context.Context, talentID string) (*models.AdminTalentConfig, error) {
	var tc models.AdminTalentConfig
	err := r.db.QueryRow(ctx, `
		SELECT talent_id, name, description, max_level, base_cost,
		       cost_scaling, bonus_per_level, stat_type, unlock_map, is_active, created_at, updated_at
		FROM talent_configs WHERE talent_id = $1
	`, talentID).Scan(
		&tc.TalentID, &tc.Name, &tc.Description, &tc.MaxLevel, &tc.BaseCost,
		&tc.CostScaling, &tc.BonusPerLevel, &tc.StatType, &tc.UnlockMap, &tc.IsActive, &tc.CreatedAt, &tc.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("talent config not found: %s", talentID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get talent config: %w", err)
	}
	return &tc, nil
}

func (r *postgresAdminRepository) CreateTalentConfig(ctx context.Context, config *models.AdminTalentConfig) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO talent_configs (talent_id, name, description, max_level, base_cost,
		                            cost_scaling, bonus_per_level, stat_type, unlock_map, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	`, config.TalentID, config.Name, config.Description, config.MaxLevel, config.BaseCost,
		config.CostScaling, config.BonusPerLevel, config.StatType, config.UnlockMap, config.IsActive)
	if err != nil {
		return fmt.Errorf("failed to create talent config: %w", err)
	}
	return nil
}

func (r *postgresAdminRepository) UpdateTalentConfig(ctx context.Context, config *models.AdminTalentConfig) error {
	result, err := r.db.Exec(ctx, `
		UPDATE talent_configs
		SET name = $2, description = $3, max_level = $4, base_cost = $5,
		    cost_scaling = $6, bonus_per_level = $7, stat_type = $8, unlock_map = $9,
		    is_active = $10, updated_at = NOW()
		WHERE talent_id = $1
	`, config.TalentID, config.Name, config.Description, config.MaxLevel, config.BaseCost,
		config.CostScaling, config.BonusPerLevel, config.StatType, config.UnlockMap, config.IsActive)
	if err != nil {
		return fmt.Errorf("failed to update talent config: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("talent config not found: %s", config.TalentID)
	}
	return nil
}

func (r *postgresAdminRepository) DeleteTalentConfig(ctx context.Context, talentID string) error {
	result, err := r.db.Exec(ctx, "DELETE FROM talent_configs WHERE talent_id = $1", talentID)
	if err != nil {
		return fmt.Errorf("failed to delete talent config: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("talent config not found: %s", talentID)
	}
	return nil
}

func (r *postgresAdminRepository) TalentHasPlayers(ctx context.Context, talentID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM user_talents WHERE talent_id = $1
	`, talentID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count players with talent: %w", err)
	}
	return count, nil
}

func (r *postgresAdminRepository) GetAdminByUsername(ctx context.Context, username string) (*models.UserProfile, string, error) {
	var profile models.UserProfile
	var passwordHash string
	var lastLoginAt *time.Time

	err := r.db.QueryRow(ctx, `
		SELECT
			u.id, u.playfab_id,
			COALESCE(u.display_name, '') as username,
			u.gold, u.total_stars_collected,
			u.is_admin,
			u.created_at, u.last_login_at,
			COALESCE(u.password_hash, '') as password_hash,
			EXISTS(SELECT 1 FROM user_bans WHERE user_id = u.id AND is_active = true) as is_banned
		FROM users u
		WHERE u.display_name = $1 AND u.is_admin = true
		LIMIT 1
	`, username).Scan(
		&profile.ID, &profile.PlayFabID,
		&profile.Username,
		&profile.TotalGold, &profile.TotalStars,
		&profile.IsAdmin,
		&profile.CreatedAt, &lastLoginAt,
		&passwordHash,
		&profile.IsBanned,
	)
	if err != nil {
		return nil, "", fmt.Errorf("admin not found: %w", err)
	}
	profile.LastLoginAt = lastLoginAt
	return &profile, passwordHash, nil
}

func (r *postgresAdminRepository) SetAdminPassword(ctx context.Context, adminID uuid.UUID, passwordHash string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET password_hash = $1 WHERE id = $2 AND is_admin = true`,
		passwordHash, adminID,
	)
	if err != nil {
		return fmt.Errorf("failed to set admin password: %w", err)
	}
	return nil
}
