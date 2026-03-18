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

// AdminLevelConfig represents level configuration for admin management
type AdminLevelConfig struct {
	LevelID        string           `json:"level_id"`
	MapID          string           `json:"map_id"`
	Name           string           `json:"name"`
	Difficulty     int              `json:"difficulty"`
	MinTimeSeconds float64          `json:"min_time_seconds"`
	BaseGold       int              `json:"base_gold"`
	RewardStars    int              `json:"reward_stars"`
	Objectives     []LevelObjective `json:"objectives"`
	IsActive       bool             `json:"is_active"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

// CreateLevelConfigRequest for creating a new level config
type CreateLevelConfigRequest struct {
	LevelID        string           `json:"level_id" validate:"required"`
	MapID          string           `json:"map_id" validate:"required"`
	Name           string           `json:"name" validate:"required"`
	Difficulty     int              `json:"difficulty" validate:"required,min=1,max=5"`
	MinTimeSeconds float64          `json:"min_time_seconds" validate:"required,gt=0"`
	BaseGold       int              `json:"base_gold" validate:"required,gte=0"`
	RewardStars    int              `json:"reward_stars" validate:"required,min=1,max=3"`
	Objectives     []LevelObjective `json:"objectives"`
	IsActive       bool             `json:"is_active"`
}

// UpdateLevelConfigRequest for partial updates to a level config
type UpdateLevelConfigRequest struct {
	Name           *string          `json:"name,omitempty"`
	Difficulty     *int             `json:"difficulty,omitempty"`
	MinTimeSeconds *float64         `json:"min_time_seconds,omitempty"`
	BaseGold       *int             `json:"base_gold,omitempty"`
	RewardStars    *int             `json:"reward_stars,omitempty"`
	Objectives     []LevelObjective `json:"objectives,omitempty"`
	IsActive       *bool            `json:"is_active,omitempty"`
}

// LevelConfigListResponse for paginated level config list
type LevelConfigListResponse struct {
	Levels     []AdminLevelConfig `json:"levels"`
	TotalCount int                `json:"total_count"`
	Page       int                `json:"page"`
	PerPage    int                `json:"per_page"`
}

// AdminTalentConfig represents talent configuration for admin management
type AdminTalentConfig struct {
	TalentID      string    `json:"talent_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	MaxLevel      int       `json:"max_level"`
	BaseCost      int       `json:"base_cost"`
	CostScaling   float64   `json:"cost_scaling"`
	BonusPerLevel float64   `json:"bonus_per_level"`
	StatType      string    `json:"stat_type"`
	UnlockMap     int       `json:"unlock_map"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateTalentConfigRequest for creating a new talent config
type CreateTalentConfigRequest struct {
	TalentID      string  `json:"talent_id" validate:"required"`
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description" validate:"required"`
	MaxLevel      int     `json:"max_level" validate:"required,min=1"`
	BaseCost      int     `json:"base_cost" validate:"required,gte=0"`
	CostScaling   float64 `json:"cost_scaling" validate:"required,gt=0"`
	BonusPerLevel float64 `json:"bonus_per_level" validate:"required,gt=0"`
	StatType      string  `json:"stat_type" validate:"required"`
	UnlockMap     int     `json:"unlock_map" validate:"required,min=1"`
	IsActive      bool    `json:"is_active"`
}

// UpdateTalentConfigRequest for partial updates to a talent config
type UpdateTalentConfigRequest struct {
	Name          *string  `json:"name,omitempty"`
	Description   *string  `json:"description,omitempty"`
	MaxLevel      *int     `json:"max_level,omitempty"`
	BaseCost      *int     `json:"base_cost,omitempty"`
	CostScaling   *float64 `json:"cost_scaling,omitempty"`
	BonusPerLevel *float64 `json:"bonus_per_level,omitempty"`
	StatType      *string  `json:"stat_type,omitempty"`
	UnlockMap     *int     `json:"unlock_map,omitempty"`
	IsActive      *bool    `json:"is_active,omitempty"`
}

// TalentConfigListResponse for paginated talent config list
type TalentConfigListResponse struct {
	Talents    []AdminTalentConfig `json:"talents"`
	TotalCount int                 `json:"total_count"`
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
}

// AdminLoginRequest for admin dashboard login with username/password
type AdminLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AdminLoginResponse for admin dashboard login response
type AdminLoginResponse struct {
	JWT      string      `json:"jwt"`
	Admin    UserProfile `json:"admin"`
}

// AnalyticsSummaryResponse for analytics dashboard
type AnalyticsSummaryResponse struct {
	TotalEventsLast24h int64          `json:"total_events_last_24h"`
	TotalEventsLast7d  int64          `json:"total_events_last_7d"`
	TopEvents          []EventCount   `json:"top_events"`
	DAULast7d          []DAUDataPoint `json:"dau_last_7d"`
	LastUpdated        time.Time      `json:"last_updated"`
}

// EventCount represents an event name and its occurrence count
type EventCount struct {
	EventName string `json:"event_name"`
	Count     int64  `json:"count"`
}

// DAUDataPoint represents daily active users for a specific date
type DAUDataPoint struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}
