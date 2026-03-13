package identity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
)

type playfabRepository struct {
	httpClient *http.Client
}

// NewPlayFabRepository creates a new PlayFab identity repository
func NewPlayFabRepository() repository.IdentityRepository {
	return &playfabRepository{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (r *playfabRepository) ValidateTicket(ctx context.Context, sessionTicket string) (string, error) {
	titleID := os.Getenv("PLAYFAB_TITLE_ID")
	if titleID == "" || titleID == "DEV" {
		return "MOCK_PLAYFAB_ID", nil
	}

	url := fmt.Sprintf("https://%s.playfabapi.com/Client/GetAccountInfo", titleID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization", sessionTicket)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to validate ticket: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid PlayFab session ticket")
	}

	var result struct {
		Data struct {
			AccountInfo struct {
				PlayFabId string `json:"PlayFabId"`
			} `json:"AccountInfo"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode PlayFab response: %w", err)
	}

	return result.Data.AccountInfo.PlayFabId, nil
}

func (r *playfabRepository) GetFriends(ctx context.Context, playfabID string) ([]string, error) {
	titleID := os.Getenv("PLAYFAB_TITLE_ID")
	secretKey := os.Getenv("PLAYFAB_SECRET_KEY")

	if titleID == "" || titleID == "DEV" || secretKey == "" {
		// Mock friends for dev
		return []string{"FRIEND_1", "FRIEND_2"}, nil
	}

	url := fmt.Sprintf("https://%s.playfabapi.com/Server/GetFriendsList", titleID)

	body := map[string]string{
		"PlayFabId": playfabID,
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-SecretKey", secretKey)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Friends []struct {
				FriendPlayFabId string `json:"FriendPlayFabId"`
			} `json:"Friends"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	friends := make([]string, len(result.Data.Friends))
	for i, f := range result.Data.Friends {
		friends[i] = f.FriendPlayFabId
	}

	return friends, nil
}
