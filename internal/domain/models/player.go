package models

import (
	"time"

	"github.com/google/uuid"
)

// Player represents a Hollow Wilds player account
type Player struct {
	ID          uuid.UUID `json:"id"`
	PlayFabID   string    `json:"playfabId"`
	DisplayName *string   `json:"displayName,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	LastSeenAt  time.Time `json:"lastSeenAt"`
}

// PlayerSave represents a player's game save data
type PlayerSave struct {
	ID          uuid.UUID    `json:"id"`
	PlayerID    uuid.UUID    `json:"playerId"`
	SaveVersion int          `json:"saveVersion"`
	SaveData    GameSaveData `json:"saveData"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

// PlayerSaveBackup represents a historical save backup
type PlayerSaveBackup struct {
	ID          uuid.UUID    `json:"id"`
	PlayerID    uuid.UUID    `json:"playerId"`
	SaveVersion int          `json:"saveVersion"`
	SaveData    GameSaveData `json:"saveData"`
	CreatedAt   time.Time    `json:"createdAt"`
}

// GameSaveData represents the complete game state structure
type GameSaveData struct {
	World          WorldData       `json:"world"`
	Player         PlayerState     `json:"player"`
	Inventory      InventoryData   `json:"inventory"`
	Sebilah        SebilahData     `json:"sebilah"`
	Base           BaseData        `json:"base"`
	DiscoveredPOIs []string        `json:"discovered_pois"`
	QuestFlags     map[string]bool `json:"quest_flags"`
}

// WorldData represents world state
type WorldData struct {
	Seed            int64 `json:"seed"`
	PlayTimeSeconds int   `json:"play_time_seconds"`
	DayCount        int   `json:"day_count"`
}

// PlayerState represents player character state
type PlayerState struct {
	Character string   `json:"character" validate:"required,oneof=RIMBA DARA BAYU SARI"` // RIMBA, DARA, BAYU, SARI
	Position  Vector2D `json:"position"`
	Health    float64  `json:"health" validate:"min=0,max=100"`
	Hunger    float64  `json:"hunger" validate:"min=0,max=100"`
	Sanity    float64  `json:"sanity" validate:"min=0,max=100"`
	Warmth    float64  `json:"warmth" validate:"min=0,max=100"`
}

// Vector2D represents a 2D position
type Vector2D struct {
	X float64 `json:"x"`
	Z float64 `json:"z"`
}

// InventoryData represents player inventory
type InventoryData struct {
	Slots          []InventorySlot `json:"slots"`
	EquippedWeapon string          `json:"equipped_weapon"`
}

// InventorySlot represents a single inventory slot
type InventorySlot struct {
	Slot     int    `json:"slot"`
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

// SebilahData represents the Sebilah weapon state
type SebilahData struct {
	WeaponID       string `json:"weapon_id"`
	SoulLevel      int    `json:"soul_level"`
	InfusionPoints int    `json:"infusion_points"`
}

// BaseData represents player's base structures
type BaseData struct {
	PlacedObjects []PlacedObject `json:"placed_objects"`
}

// PlacedObject represents a placed structure in the world
type PlacedObject struct {
	ObjectID string  `json:"object_id"`
	X        float64 `json:"x"`
	Z        float64 `json:"z"`
}

// ValidCharacters defines the allowed character names
var ValidCharacters = []string{"RIMBA", "DARA", "BAYU", "SARI"}
