package models

import (
	"time"

	"github.com/google/uuid"
)

// AdminAction represents an admin audit log entry
type AdminAction struct {
	ID           uuid.UUID      `json:"id"`
	AdminID      uuid.UUID      `json:"admin_id"`
	ActionType   string         `json:"action_type"`
	TargetUserID *uuid.UUID     `json:"target_user_id,omitempty"`
	Details      map[string]any `json:"details"`
	IPAddress    string         `json:"ip_address,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

// UserBan represents a user ban record
type UserBan struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	BannedBy    uuid.UUID  `json:"banned_by"`
	Reason      string     `json:"reason"`
	BannedUntil *time.Time `json:"banned_until,omitempty"` // NULL = permanent
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UnbannedAt  *time.Time `json:"unbanned_at,omitempty"`
	UnbannedBy  *uuid.UUID `json:"unbanned_by,omitempty"`
	UnbanReason string     `json:"unban_reason,omitempty"`
}

// UserSearchResponse for user search endpoint
type UserSearchResponse struct {
	Users      []UserProfile `json:"users"`
	TotalCount int           `json:"total_count"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
}

// UserProfile detailed user info for admins
type UserProfile struct {
	ID              uuid.UUID  `json:"id"`
	PlayFabID       string     `json:"playfab_id"`
	Email           string     `json:"email,omitempty"`
	Username        string     `json:"username"`
	TotalGold       int        `json:"total_gold"`
	TotalStars      int        `json:"total_stars"`
	LevelsCompleted int        `json:"levels_completed"`
	IsAdmin         bool       `json:"is_admin"`
	IsBanned        bool       `json:"is_banned"`
	CreatedAt       time.Time  `json:"created_at"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
}

// AdjustGoldRequest for manual gold adjustment
type AdjustGoldRequest struct {
	Amount int    `json:"amount" validate:"required"` // Can be negative
	Reason string `json:"reason" validate:"required,min=10"`
}

// AdjustGoldResponse for gold adjustment result
type AdjustGoldResponse struct {
	UserID     uuid.UUID `json:"user_id"`
	OldBalance int       `json:"old_balance"`
	NewBalance int       `json:"new_balance"`
	Adjustment int       `json:"adjustment"`
	Reason     string    `json:"reason"`
	ActionID   uuid.UUID `json:"action_id"`
}

// BanUserRequest for banning a user
type BanUserRequest struct {
	Reason      string     `json:"reason" validate:"required,min=10"`
	BannedUntil *time.Time `json:"banned_until,omitempty"` // NULL = permanent
}

// BanUserResponse for ban result
type BanUserResponse struct {
	BanID       uuid.UUID  `json:"ban_id"`
	UserID      uuid.UUID  `json:"user_id"`
	Reason      string     `json:"reason"`
	BannedUntil *time.Time `json:"banned_until,omitempty"`
	IsPermanent bool       `json:"is_permanent"`
}

// UnbanUserRequest for unbanning a user
type UnbanUserRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// UnbanUserResponse for unban result
type UnbanUserResponse struct {
	UserID      uuid.UUID `json:"user_id"`
	UnbannedAt  time.Time `json:"unbanned_at"`
	UnbanReason string    `json:"unban_reason"`
}

// AdminActionsResponse for audit log
type AdminActionsResponse struct {
	Actions    []AdminAction `json:"actions"`
	TotalCount int           `json:"total_count"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
}

// SystemStatsResponse for system overview
type SystemStatsResponse struct {
	TotalUsers        int       `json:"total_users"`
	ActiveToday       int       `json:"active_today"`
	ActiveThisWeek    int       `json:"active_this_week"`
	ActiveThisMonth   int       `json:"active_this_month"`
	BannedUsers       int       `json:"banned_users"`
	AdminActionsToday int       `json:"admin_actions_today"`
	TotalGoldInGame   int64     `json:"total_gold_in_game"`
	TotalStarsEarned  int64     `json:"total_stars_earned"`
	LevelsCompleted   int64     `json:"levels_completed"`
	LastUpdated       time.Time `json:"last_updated"`
}

// ExportUserDataResponse for GDPR export
type ExportUserDataResponse struct {
	User             UserProfile      `json:"user"`
	LevelCompletions []map[string]any `json:"level_completions"`
	Talents          []map[string]any `json:"talents"`
	BanHistory       []UserBan        `json:"ban_history"`
	ExportedAt       time.Time        `json:"exported_at"`
}

// ResetLeaderboardRequest for resetting a leaderboard
type ResetLeaderboardRequest struct {
	Reason string `json:"reason" validate:"required,min=10"`
}

// LeaderboardStatsResponse for leaderboard analytics
type LeaderboardStatsResponse struct {
	TotalEntries  int64            `json:"total_entries"`
	UniquePlayers int64            `json:"unique_players"`
	LevelsTracked int              `json:"levels_tracked"`
	TopLevelID    string           `json:"top_level_id,omitempty"`
	TopLevelCount int              `json:"top_level_count,omitempty"`
	AverageTime   float64          `json:"average_time,omitempty"`
	LastUpdated   time.Time        `json:"last_updated"`
	LevelStats    []LevelStatsInfo `json:"level_stats"`
}

// LevelStatsInfo for individual level stats
type LevelStatsInfo struct {
	LevelID       string  `json:"level_id"`
	TotalEntries  int     `json:"total_entries"`
	UniquePlayers int     `json:"unique_players"`
	AverageTime   float64 `json:"average_time"`
	BestTime      float64 `json:"best_time"`
}
