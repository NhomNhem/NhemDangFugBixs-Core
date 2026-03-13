package models

// HollowWildsLoginRequest represents the login request for Hollow Wilds
type HollowWildsLoginRequest struct {
	PlayfabSessionTicket string `json:"playfab_session_ticket" validate:"required,min=1"`
}

// HollowWildsAuthResponse represents the authentication response with refresh token
type HollowWildsAuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	PlayerID     string `json:"player_id"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents the refresh token response
type RefreshTokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

// SaveGameRequest represents a save game request
type SaveGameRequest struct {
	World          WorldData       `json:"world"`
	Player         PlayerState     `json:"player"`
	Inventory      InventoryData   `json:"inventory"`
	Sebilah        SebilahData     `json:"sebilah"`
	Base           BaseData        `json:"base"`
	DiscoveredPOIs []string        `json:"discovered_pois"`
	QuestFlags     map[string]bool `json:"quest_flags"`
}

// SaveGameResponse represents a save game response
type SaveGameResponse struct {
	Success     bool   `json:"success"`
	SaveVersion int    `json:"save_version"`
	UpdatedAt   string `json:"updated_at"`
}

// LoadGameResponse represents a load game response
type LoadGameResponse struct {
	PlayerID       string          `json:"player_id"`
	SaveVersion    int             `json:"save_version"`
	UpdatedAt      string          `json:"updated_at"`
	World          WorldData       `json:"world"`
	Player         PlayerState     `json:"player"`
	Inventory      InventoryData   `json:"inventory"`
	Sebilah        SebilahData     `json:"sebilah"`
	Base           BaseData        `json:"base"`
	DiscoveredPOIs []string        `json:"discovered_pois"`
	QuestFlags     map[string]bool `json:"quest_flags"`
}

// BackupResponse represents a backup creation response
type BackupResponse struct {
	Success   bool   `json:"success"`
	BackupID  string `json:"backup_id"`
	CreatedAt string `json:"created_at"`
}

// BackupListResponse represents a backup list response
type BackupListResponse struct {
	Backups []BackupInfo `json:"backups"`
}

// BackupInfo represents a single backup entry
type BackupInfo struct {
	BackupID    string `json:"backup_id"`
	SaveVersion int    `json:"save_version"`
	CreatedAt   string `json:"created_at"`
}

// RestoreRequest represents a restore request
type RestoreRequest struct {
	BackupID string `json:"backup_id" validate:"required"`
}

// VersionConflictError represents a version conflict error response
type VersionConflictError struct {
	ErrorCode     string `json:"error"`
	ServerVersion int    `json:"server_version"`
	Message       string `json:"message"`
}

func (e *VersionConflictError) Error() string {
	return e.Message
}

// AnalyticsEvent represents a single analytics event
type AnalyticsEvent struct {
	UserID    interface{}            `json:"user_id"`
	EventName string                 `json:"event_name" validate:"required"`
	Timestamp string                 `json:"timestamp" validate:"required"`
	SessionID string                 `json:"session_id"`
	Payload   map[string]interface{} `json:"payload"`
}

// AnalyticsEventsRequest represents a batch analytics submission
type AnalyticsEventsRequest struct {
	Events []AnalyticsEvent `json:"events" validate:"required"`
}

// AnalyticsEventsResponse represents analytics submission response
type AnalyticsEventsResponse struct {
	Accepted int `json:"accepted"`
	Rejected int `json:"rejected"`
}
