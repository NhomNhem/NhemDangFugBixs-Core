package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
)

var baseURL = "https://gamefeel-backend.fly.dev/api/v1"

func main() {
	// Check for environment overrides
	if envURL := os.Getenv("TEST_BASE_URL"); envURL != "" {
		baseURL = envURL
	}

	testPlayFabID := "TEST_PLAYER_1"
	if envID := os.Getenv("TEST_PLAYFAB_ID"); envID != "" {
		testPlayFabID = envID
	}

	fmt.Printf("🧪 Starting Hollow Wilds Phase 1 Integration Tests against %s...\n", baseURL)

	// 1. Auth: Login
	fmt.Println("\nStep 1: Testing Login...")
	loginReq := models.HollowWildsLoginRequest{
		PlayfabSessionTicket: "test-ticket",
	}

	if envTicket := os.Getenv("TEST_PLAYFAB_TICKET"); envTicket != "" {
		loginReq.PlayfabSessionTicket = envTicket
	}

	// Create request with X-PlayFab-ID header
	client := &http.Client{}
	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", baseURL+"/auth/hw/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-PlayFab-ID", testPlayFabID)


	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Login failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var authResp models.HollowWildsAuthResponse
	json.NewDecoder(resp.Body).Decode(&authResp)

	if authResp.Token == "" {
		fmt.Println("❌ Failed to get JWT token")
		return
	}
	fmt.Printf("✅ Login successful. PlayerID: %s\n", authResp.PlayerID)
	token := authResp.Token

	// 2. Save/Load: Update Save
	fmt.Println("\nStep 2: Testing Update Save...")
	saveReq := models.SaveGameRequest{
		World: models.WorldData{
			Seed:            123456,
			PlayTimeSeconds: 3600,
			DayCount:        5,
		},
		Player: models.PlayerState{
			Character: "RIMBA",
			Health:    100,
			Hunger:    80,
		},
	}

	body, _ = json.Marshal(saveReq)
	req, _ = http.NewRequest("PUT", baseURL+"/player/save", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("❌ Save failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var saveResp models.SaveGameResponse
	json.NewDecoder(resp.Body).Decode(&saveResp)

	if !saveResp.Success {
		fmt.Println("❌ Failed to update save")
		return
	}
	fmt.Printf("✅ Save successful. Version: %d\n", saveResp.SaveVersion)

	// 3. Save/Load: Get Save
	fmt.Println("\nStep 3: Testing Get Save...")
	req, _ = http.NewRequest("GET", baseURL+"/player/save", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("❌ Load failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var loadResp models.LoadGameResponse
	json.NewDecoder(resp.Body).Decode(&loadResp)

	if loadResp.Player.Character != "RIMBA" {
		fmt.Println("❌ Load data mismatch")
		return
	}
	fmt.Println("✅ Load successful")

	// 4. Leaderboard: Submit
	fmt.Println("\nStep 4: Testing Leaderboard Submission...")
	lbReq := models.LeaderboardSubmitRequest{
		Type:      "longest_run_days",
		Value:     15,
		Character: "RIMBA",
		WorldSeed: 123456,
	}

	body, _ = json.Marshal(lbReq)
	req, _ = http.NewRequest("POST", baseURL+"/leaderboard/submit", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("❌ Leaderboard submission failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var lbResp models.LeaderboardSubmitResponse
	json.NewDecoder(resp.Body).Decode(&lbResp)

	if !lbResp.Success {
		fmt.Println("❌ Failed to submit leaderboard")
		return
	}
	fmt.Printf("✅ Leaderboard submitted. Global Rank: %d\n", lbResp.GlobalRank)

	// 5. Analytics: Track Events
	fmt.Println("\nStep 5: Testing Analytics...")
	analyticsReq := models.AnalyticsEventsRequest{
		Events: []models.AnalyticsEvent{
			{
				EventName: "player_death",
				Timestamp: time.Now().Format(time.RFC3339),
				Payload: map[string]interface{}{
					"cause": "spirit_attack",
				},
			},
		},
	}

	body, _ = json.Marshal(analyticsReq)
	req, _ = http.NewRequest("POST", baseURL+"/analytics/events", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("❌ Analytics submission failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var analyticsResp models.AnalyticsEventsResponse
	json.NewDecoder(resp.Body).Decode(&analyticsResp)

	fmt.Printf("✅ Analytics test: Accepted %d, Rejected %d\n", analyticsResp.Accepted, analyticsResp.Rejected)

	fmt.Println("\n🎉 All tests completed successfully!")
}
