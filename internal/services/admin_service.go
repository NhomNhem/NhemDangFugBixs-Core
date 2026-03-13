package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/database"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

// SearchUsers searches for users by PlayFab ID, email, or username
func (s *AdminService) SearchUsers(ctx context.Context, query string, page, perPage int) (*models.UserSearchResponse, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	offset := (page - 1) * perPage
	searchPattern := "%" + query + "%"

	// Count total matches
	var totalCount int
	err := db.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM users 
		WHERE playfab_id ILIKE $1 
		   OR COALESCE(display_name, '') ILIKE $1
	`, searchPattern).Scan(&totalCount)

	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users with stats
	rows, err := db.Query(ctx, `
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
		   OR COALESCE(u.display_name, '') ILIKE $1
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

// GetUserProfile gets detailed profile for a specific user
func (s *AdminService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	var profile models.UserProfile
	err := db.QueryRow(ctx, `
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

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return &profile, nil
}

// AdjustGold manually adjusts user's gold balance
func (s *AdminService) AdjustGold(ctx context.Context, adminID, userID uuid.UUID, amount int, reason, ipAddress string) (*models.AdjustGoldResponse, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	// Start transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get current balance
	var oldBalance int
	err = tx.QueryRow(ctx, "SELECT gold FROM users WHERE id = $1", userID).Scan(&oldBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to get current balance: %w", err)
	}

	// Calculate new balance
	newBalance := oldBalance + amount
	if newBalance < 0 {
		return nil, fmt.Errorf("insufficient gold: cannot reduce below 0")
	}

	// Update user's gold
	_, err = tx.Exec(ctx, "UPDATE users SET gold = $1 WHERE id = $2", newBalance, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update gold: %w", err)
	}

	// Log admin action
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

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Admin %s adjusted gold for user %s: %+d (reason: %s)", adminID, userID, amount, reason)

	return &models.AdjustGoldResponse{
		UserID:     userID,
		OldBalance: oldBalance,
		NewBalance: newBalance,
		Adjustment: amount,
		Reason:     reason,
		ActionID:   actionID,
	}, nil
}

// BanUser bans a user with reason and optional duration
func (s *AdminService) BanUser(ctx context.Context, adminID, userID uuid.UUID, reason string, bannedUntil *time.Time, ipAddress string) (*models.BanUserResponse, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	// Start transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if user already has active ban
	var hasActiveBan bool
	err = tx.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM user_bans WHERE user_id = $1 AND is_active = true)",
		userID,
	).Scan(&hasActiveBan)

	if err != nil {
		return nil, fmt.Errorf("failed to check existing bans: %w", err)
	}

	if hasActiveBan {
		return nil, fmt.Errorf("user already has an active ban")
	}

	// Create ban record
	var banID uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO user_bans (user_id, banned_by, reason, banned_until, is_active)
		VALUES ($1, $2, $3, $4, true)
		RETURNING id
	`, userID, adminID, reason, bannedUntil).Scan(&banID)

	if err != nil {
		return nil, fmt.Errorf("failed to create ban: %w", err)
	}

	// Log admin action
	details := map[string]any{
		"reason":       reason,
		"banned_until": bannedUntil,
		"is_permanent": bannedUntil == nil,
	}
	detailsJSON, _ := json.Marshal(details)

	_, err = tx.Exec(ctx, `
		INSERT INTO admin_actions (admin_id, action_type, target_user_id, details, ip_address)
		VALUES ($1, $2, $3, $4, $5)
	`, adminID, "ban_user", userID, detailsJSON, ipAddress)

	if err != nil {
		return nil, fmt.Errorf("failed to log admin action: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	isPermanent := bannedUntil == nil
	log.Printf("Admin %s banned user %s (permanent: %v, reason: %s)", adminID, userID, isPermanent, reason)

	return &models.BanUserResponse{
		BanID:       banID,
		UserID:      userID,
		Reason:      reason,
		BannedUntil: bannedUntil,
		IsPermanent: isPermanent,
	}, nil
}

// UnbanUser unbans a user
func (s *AdminService) UnbanUser(ctx context.Context, adminID, userID uuid.UUID, reason, ipAddress string) (*models.UnbanUserResponse, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	// Start transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update active bans
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

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no active ban found for user")
	}

	// Log admin action
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

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Admin %s unbanned user %s (reason: %s)", adminID, userID, reason)

	return &models.UnbanUserResponse{
		UserID:      userID,
		UnbannedAt:  time.Now(),
		UnbanReason: reason,
	}, nil
}

// GetBanHistory gets ban history for a user
func (s *AdminService) GetBanHistory(ctx context.Context, userID uuid.UUID) ([]models.UserBan, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	rows, err := db.Query(ctx, `
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

// GetAdminActions gets audit log of admin actions
func (s *AdminService) GetAdminActions(ctx context.Context, page, perPage int) (*models.AdminActionsResponse, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	offset := (page - 1) * perPage

	// Count total actions
	var totalCount int
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM admin_actions").Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count admin actions: %w", err)
	}

	// Get actions
	rows, err := db.Query(ctx, `
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

		// Parse JSON details
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

// GetSystemStats gets system overview statistics
func (s *AdminService) GetSystemStats(ctx context.Context) (*models.SystemStatsResponse, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	stats := &models.SystemStatsResponse{
		LastUpdated: time.Now(),
	}

	// Total users
	db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)

	// Active users (last 24h, 7d, 30d)
	db.QueryRow(ctx, `
		SELECT COUNT(*) FROM users 
		WHERE last_login_at > NOW() - INTERVAL '24 hours'
	`).Scan(&stats.ActiveToday)

	db.QueryRow(ctx, `
		SELECT COUNT(*) FROM users 
		WHERE last_login_at > NOW() - INTERVAL '7 days'
	`).Scan(&stats.ActiveThisWeek)

	db.QueryRow(ctx, `
		SELECT COUNT(*) FROM users 
		WHERE last_login_at > NOW() - INTERVAL '30 days'
	`).Scan(&stats.ActiveThisMonth)

	// Banned users
	db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT user_id) FROM user_bans WHERE is_active = true
	`).Scan(&stats.BannedUsers)

	// Admin actions today
	db.QueryRow(ctx, `
		SELECT COUNT(*) FROM admin_actions 
		WHERE created_at > NOW() - INTERVAL '24 hours'
	`).Scan(&stats.AdminActionsToday)

	// Total gold and stars
	db.QueryRow(ctx, "SELECT COALESCE(SUM(gold), 0) FROM users").Scan(&stats.TotalGoldInGame)
	db.QueryRow(ctx, "SELECT COALESCE(SUM(total_stars_collected), 0) FROM users").Scan(&stats.TotalStarsEarned)

	// Levels completed
	db.QueryRow(ctx, "SELECT COUNT(*) FROM level_completions").Scan(&stats.LevelsCompleted)

	return stats, nil
}

// ExportUserData exports all user data for GDPR compliance
func (s *AdminService) ExportUserData(ctx context.Context, userID uuid.UUID) (*models.ExportUserDataResponse, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	// Get user profile
	profile, err := s.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get level completions (optional - may not exist yet)
	var levelCompletions []map[string]any
	levelRows, err := db.Query(ctx, `
		SELECT level_id, map_id, time_seconds, stars_earned, completed_at
		FROM level_completions
		WHERE user_id = $1
		ORDER BY completed_at DESC
	`, userID)
	if err == nil {
		defer levelRows.Close()
		for levelRows.Next() {
			var levelID, mapID string
			var timeSeconds int
			var starsEarned int
			var completedAt time.Time
			levelRows.Scan(&levelID, &mapID, &timeSeconds, &starsEarned, &completedAt)
			levelCompletions = append(levelCompletions, map[string]any{
				"level_id":     levelID,
				"map_id":       mapID,
				"time_seconds": timeSeconds,
				"stars_earned": starsEarned,
				"completed_at": completedAt,
			})
		}
	}

	// Get talents (optional - may not exist yet)
	var talents []map[string]any
	talentRows, err := db.Query(ctx, `
		SELECT talent_id, current_level, times_upgraded, last_upgraded_at
		FROM user_talents
		WHERE user_id = $1
	`, userID)
	if err == nil {
		defer talentRows.Close()
		for talentRows.Next() {
			var talentID string
			var currentLevel, timesUpgraded int
			var lastUpgradedAt time.Time
			talentRows.Scan(&talentID, &currentLevel, &timesUpgraded, &lastUpgradedAt)
			talents = append(talents, map[string]any{
				"talent_id":        talentID,
				"current_level":    currentLevel,
				"times_upgraded":   timesUpgraded,
				"last_upgraded_at": lastUpgradedAt,
			})
		}
	}

	// Get ban history
	banHistory, err := s.GetBanHistory(ctx, userID)
	if err != nil {
		banHistory = []models.UserBan{}
	}

	return &models.ExportUserDataResponse{
		User:             *profile,
		LevelCompletions: levelCompletions,
		Talents:          talents,
		BanHistory:       banHistory,
		ExportedAt:       time.Now(),
	}, nil
}
