package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/database"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// AuthService handles authentication logic
type AuthService struct{}

// NewAuthService creates a new auth service
func NewAuthService() *AuthService {
	return &AuthService{}
}

// PlayFabValidationResponse represents PlayFab API response
type PlayFabValidationResponse struct {
	Code   int                    `json:"code"`
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
	Error  string                 `json:"error,omitempty"`
}

// ValidatePlayFabToken validates a PlayFab session token
func (s *AuthService) ValidatePlayFabToken(sessionToken string, playfabID string) error {
	if sessionToken == "" {
		return fmt.Errorf("session token is required")
	}

	if playfabID == "" {
		return fmt.Errorf("playfab ID is required")
	}

	// Get PlayFab Title ID from environment
	titleID := os.Getenv("PLAYFAB_TITLE_ID")
	if titleID == "" {
		// If not configured, skip validation (development mode)
		return nil
	}

	// Call PlayFab Client API to validate the session token
	// We use GetAccountInfo which requires valid session token
	url := fmt.Sprintf("https://%s.playfabapi.com/Client/GetAccountInfo", titleID)

	reqBody := fmt.Sprintf(`{"PlayFabId": "%s"}`, playfabID)
	req, err := http.NewRequest("POST", url, strings.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add session token to headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization", sessionToken)

	// Send request
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var pfResp PlayFabValidationResponse
	if err := json.Unmarshal(body, &pfResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if validation succeeded
	if pfResp.Code != 200 {
		return fmt.Errorf("invalid PlayFab token: %s", pfResp.Error)
	}

	return nil
}

// GetOrCreateUser gets existing user or creates new one
func (s *AuthService) GetOrCreateUser(ctx context.Context, playfabID string, displayName *string) (*models.User, error) {
	if database.Pool == nil {
		return nil, fmt.Errorf("database not connected")
	}

	// Try to get existing user
	var user models.User
	err := database.Pool.QueryRow(ctx, `
		SELECT id, playfab_id, display_name, gold, diamonds, 
		       max_map_unlocked, total_stars_collected,
		       created_at, last_login_at, last_played_at, total_play_time_seconds,
		       facebook_id, google_id, is_banned, ban_reason, banned_at, deleted_at
		FROM users 
		WHERE playfab_id = $1 AND deleted_at IS NULL
	`, playfabID).Scan(
		&user.ID, &user.PlayFabID, &user.DisplayName, &user.Gold, &user.Diamonds,
		&user.MaxMapUnlocked, &user.TotalStarsCollected,
		&user.CreatedAt, &user.LastLoginAt, &user.LastPlayedAt, &user.TotalPlayTimeSeconds,
		&user.FacebookID, &user.GoogleID, &user.IsBanned, &user.BanReason, &user.BannedAt, &user.DeletedAt,
	)

	if err == nil {
		// User exists - update last login
		_, err = database.Pool.Exec(ctx, `
			UPDATE users SET last_login_at = NOW() WHERE id = $1
		`, user.ID)

		if err != nil {
			return nil, fmt.Errorf("failed to update last login: %w", err)
		}

		user.LastLoginAt = time.Now()
		return &user, nil
	}

	if err != pgx.ErrNoRows {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// User doesn't exist - create new user
	user.ID = uuid.New()
	user.PlayFabID = playfabID
	user.DisplayName = displayName
	user.Gold = 0
	user.Diamonds = 0
	user.MaxMapUnlocked = 1
	user.TotalStarsCollected = 0
	user.CreatedAt = time.Now()
	user.LastLoginAt = time.Now()
	user.IsBanned = false

	_, err = database.Pool.Exec(ctx, `
		INSERT INTO users (id, playfab_id, display_name, gold, diamonds, 
		                   max_map_unlocked, total_stars_collected, 
		                   created_at, last_login_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, user.ID, user.PlayFabID, user.DisplayName, user.Gold, user.Diamonds,
		user.MaxMapUnlocked, user.TotalStarsCollected,
		user.CreatedAt, user.LastLoginAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// GenerateJWT generates a JWT token for the user
func (s *AuthService) GenerateJWT(userID uuid.UUID, playfabID string) (string, int, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", 0, fmt.Errorf("JWT_SECRET not configured")
	}

	expiresIn := 24 * 60 * 60 // 24 hours
	now := time.Now()
	expiresAt := now.Add(time.Duration(expiresIn) * time.Second)

	claims := jwt.MapClaims{
		"userId":    userID.String(),
		"playfabId": playfabID,
		"iat":       now.Unix(),
		"exp":       expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresIn, nil
}

// VerifyJWT verifies and parses a JWT token
func (s *AuthService) VerifyJWT(tokenString string) (*models.JWTClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &models.JWTClaims{
		UserID:    claims["userId"].(string),
		PlayFabID: claims["playfabId"].(string),
		IssuedAt:  int64(claims["iat"].(float64)),
		ExpiresAt: int64(claims["exp"].(float64)),
	}, nil
}
