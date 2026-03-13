package identity

import (
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
