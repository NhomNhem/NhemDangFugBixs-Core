package models

import (
	"time"

	"github.com/google/uuid"
)

// TalentUpgradeRequest represents the request to upgrade a talent
type TalentUpgradeRequest struct {
	TalentID string `json:"talentId" validate:"required"`
}

// TalentUpgradeResponse represents the response after upgrading
type TalentUpgradeResponse struct {
	Success       bool   `json:"success"`
	TalentID      string `json:"talentId"`
	NewLevel      int    `json:"newLevel"`
	GoldSpent     int    `json:"goldSpent"`
	NewTotalGold  int    `json:"newTotalGold"`
	StatBonus     string `json:"statBonus"` // e.g., "+5% HP", "+3% Damage"
	NextLevelCost int    `json:"nextLevelCost,omitempty"`
}

// UserTalent represents a user's talent progress
type UserTalent struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"userId"`
	TalentID     string     `json:"talentId"`
	CurrentLevel int        `json:"currentLevel"`
	UpgradedAt   time.Time  `json:"upgradedAt"`
}

// TalentConfig represents talent configuration
type TalentConfig struct {
	ID           string  `json:"id"`           // "health", "damage", "armor", "speed"
	Name         string  `json:"name"`         // Display name
	Description  string  `json:"description"`  // "Increases max HP"
	MaxLevel     int     `json:"maxLevel"`     // 20 for primary, 10 for advanced
	BaseCost     int     `json:"baseCost"`     // Starting cost
	CostScaling  float64 `json:"costScaling"`  // Multiplier per level
	BonusPerLevel float64 `json:"bonusPerLevel"` // Stat bonus per level (e.g., 5 for +5%)
	StatType     string  `json:"statType"`     // "hp", "damage", "armor", "speed"
	UnlockMap    int     `json:"unlockMap"`    // 1 for primary, 2+ for advanced
}

// Talent IDs
const (
	TalentHealth            = "health"
	TalentDamage            = "damage"
	TalentArmor             = "armor"
	TalentSpeed             = "speed"
	TalentDashCooldown      = "dash_cooldown"
	TalentCounterWindow     = "counter_window"
	TalentVulnerableDuration = "vulnerable_duration"
)
