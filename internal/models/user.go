package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID          uuid.UUID  `json:"id"`
	PlayFabID   string     `json:"playfabId"`
	DisplayName *string    `json:"displayName,omitempty"`
	
	// Currency
	Gold     int `json:"gold"`
	Diamonds int `json:"diamonds"`
	
	// Progression
	MaxMapUnlocked       int `json:"maxMapUnlocked"`
	TotalStarsCollected  int `json:"totalStarsCollected"`
	
	// Metadata
	CreatedAt       time.Time  `json:"createdAt"`
	LastLoginAt     time.Time  `json:"lastLoginAt"`
	LastPlayedAt    *time.Time `json:"lastPlayedAt,omitempty"`
	TotalPlayTimeSeconds *int  `json:"totalPlayTimeSeconds,omitempty"`
	
	// Social
	FacebookID *string `json:"facebookId,omitempty"`
	GoogleID   *string `json:"googleId,omitempty"`
	
	// Flags
	IsBanned  bool       `json:"isBanned"`
	BanReason *string    `json:"banReason,omitempty"`
	BannedAt  *time.Time `json:"bannedAt,omitempty"`
	
	// GDPR
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

// AuthRequest represents the login request
type AuthRequest struct {
	PlayFabID   string  `json:"playfabId" validate:"required"`
	DisplayName *string `json:"displayName,omitempty"`
}

// AuthResponse represents the login response
type AuthResponse struct {
	JWT       string `json:"jwt"`
	User      User   `json:"user"`
	ExpiresIn int    `json:"expiresIn"` // seconds
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID    string `json:"userId"`
	PlayFabID string `json:"playfabId"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}
