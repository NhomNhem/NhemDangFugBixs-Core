package models

import "time"

// LeaderboardEntry represents a single leaderboard entry
type LeaderboardEntry struct {
	Rank            int        `json:"rank"`
	PlayerID        string     `json:"playerId"`
	DisplayName     *string    `json:"displayName,omitempty"`
	TotalStars      int        `json:"totalStars,omitempty"`
	BestTime        float64    `json:"bestTime,omitempty"`
	Stars           int        `json:"stars,omitempty"`
	PlayCount       int        `json:"playCount,omitempty"`
	LevelsCompleted int        `json:"levelsCompleted,omitempty"`
	MaxMapUnlocked  int        `json:"maxMapUnlocked,omitempty"`
	FirstCompleted  *time.Time `json:"firstCompleted,omitempty"`
}

// GlobalLeaderboardResponse for global rankings
type GlobalLeaderboardResponse struct {
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
	Total       int                `json:"total"`
	Page        int                `json:"page"`
	PerPage     int                `json:"perPage"`
}

// LevelLeaderboardResponse for per-level rankings
type LevelLeaderboardResponse struct {
	LevelID     string             `json:"levelId"`
	MapID       string             `json:"mapId,omitempty"`
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
	Total       int                `json:"total"`
}

// PlayerStatsResponse for player position & stats
type PlayerStatsResponse struct {
	PlayerID           string             `json:"playerId"`
	DisplayName        *string            `json:"displayName,omitempty"`
	GlobalRank         int                `json:"globalRank"`
	TotalStars         int                `json:"totalStars"`
	MaxMapUnlocked     int                `json:"maxMapUnlocked"`
	LevelsCompleted    int                `json:"levelsCompleted"`
	AverageStars       float64            `json:"averageStars"`
	BestTime           float64            `json:"bestTime,omitempty"`
	Stars              int                `json:"stars,omitempty"`
	SurroundingPlayers []LeaderboardEntry `json:"surroundingPlayers,omitempty"`
}

// LevelStatsResponse for level analytics
type LevelStatsResponse struct {
	LevelID        string  `json:"levelId"`
	MapID          string  `json:"mapId,omitempty"`
	UniquePlayers  int     `json:"uniquePlayers"`
	AverageTime    float64 `json:"averageTime"`
	BestTime       float64 `json:"bestTime"`
	AverageStars   float64 `json:"averageStars"`
	TotalPlays     int     `json:"totalPlays"`
	CompletionRate float64 `json:"completionRate,omitempty"`
}

// HollowWildsLeaderboardEntry represents a single entry in Hollow Wilds leaderboard
type HollowWildsLeaderboardEntry struct {
	Rank        int                    `json:"rank"`
	PlayerID    string                 `json:"player_id"`
	DisplayName string                 `json:"display_name"`
	Value       int64                  `json:"value"`
	Character   string                 `json:"character"`
	WorldSeed   int64                  `json:"world_seed"`
	CombatBuild string                 `json:"combat_build"`
	UpdatedAt   string                 `json:"updated_at"`
	RunMetadata map[string]interface{} `json:"run_metadata,omitempty"`
}

// HollowWildsLeaderboardResponse represents a Hollow Wilds leaderboard
type HollowWildsLeaderboardResponse struct {
	Type      string                        `json:"type"`
	Scope     string                        `json:"scope"`
	Character string                        `json:"character,omitempty"`
	Total     int                           `json:"total"`
	Entries   []HollowWildsLeaderboardEntry `json:"entries"`
}

// LeaderboardSubmitRequest represents a request to submit a leaderboard entry
type LeaderboardSubmitRequest struct {
	Type        string                 `json:"type" validate:"required,oneof=longest_run_days sebilah_soul_level bosses_killed"`
	Value       int64                  `json:"value" validate:"required,min=0"`
	Character   string                 `json:"character" validate:"required,oneof=RIMBA DARA BAYU SARI"`
	WorldSeed   int64                  `json:"world_seed"`
	CombatBuild string                 `json:"combat_build" validate:"omitempty,oneof=balanced berserker shade_walker"`
	RunMetadata map[string]interface{} `json:"run_metadata"`
}

// LeaderboardSubmitResponse represents the response after submitting a leaderboard entry
type LeaderboardSubmitResponse struct {
	Success            bool `json:"success"`
	GlobalRank         int  `json:"global_rank"`
	CharacterRank      int  `json:"character_rank"`
	PreviousGlobalRank int  `json:"previous_global_rank"`
	IsPersonalBest     bool `json:"is_personal_best"`
}

// PlayerLeaderboardEntry represents a single entry for a player across types
type PlayerLeaderboardEntry struct {
	Type          string `json:"type"`
	GlobalRank    int    `json:"global_rank"`
	CharacterRank int    `json:"character_rank"`
	Character     string `json:"character"`
	Value         int64  `json:"value"`
	PersonalBest  bool   `json:"personal_best"`
}

// PlayerLeaderboardResponse represents all leaderboard entries for a player
type PlayerLeaderboardResponse struct {
	Entries []PlayerLeaderboardEntry `json:"entries"`
}
