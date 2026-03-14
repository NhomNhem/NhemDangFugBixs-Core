package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type modeType string

const (
	modePreDeploy  modeType = "predeploy"
	modePostDeploy modeType = "postdeploy"
	modeFull       modeType = "full"
)

func main() {
	var (
		mode    string
		baseURL string
		appName string
	)

	flag.StringVar(&mode, "mode", string(modeFull), "Verification mode: predeploy|postdeploy|full")
	flag.StringVar(&baseURL, "base-url", "https://gamefeel-backend.fly.dev", "Base URL for post-deploy smoke tests")
	flag.StringVar(&appName, "app", "gamefeel-backend", "Fly.io app name")
	flag.Parse()

	selectedMode := modeType(strings.ToLower(strings.TrimSpace(mode)))
	if selectedMode != modePreDeploy && selectedMode != modePostDeploy && selectedMode != modeFull {
		exitf("invalid -mode value: %s", mode)
	}

	fmt.Printf("Starting deployment verification (mode=%s, app=%s)\n", selectedMode, appName)

	if selectedMode == modePreDeploy || selectedMode == modeFull {
		if err := runPreDeployChecks(appName); err != nil {
			exitf("pre-deploy checks failed: %v", err)
		}
		fmt.Println("✅ Pre-deploy checks passed")
	}

	if selectedMode == modePostDeploy || selectedMode == modeFull {
		if err := runPostDeployChecks(baseURL, appName); err != nil {
			exitf("post-deploy checks failed: %v", err)
		}
		fmt.Println("✅ Post-deploy checks passed")
	}

	fmt.Println("🎉 Deployment verification complete")
}

func runPreDeployChecks(appName string) error {
	fmt.Println("[PreDeploy] Checking fly.toml runtime port configuration")
	flyToml, err := os.ReadFile("fly.toml")
	if err != nil {
		return fmt.Errorf("read fly.toml: %w", err)
	}
	if !strings.Contains(string(flyToml), "internal_port = 8080") {
		return errors.New("fly.toml must define internal_port = 8080")
	}

	if strings.TrimSpace(os.Getenv("FLY_API_TOKEN")) == "" {
		return errors.New("FLY_API_TOKEN is required for Fly CLI verification")
	}

	fmt.Println("[PreDeploy] Building Docker image")
	if _, err := runCommand("docker", "build", "-t", "gamefeel-backend:verify", "."); err != nil {
		return fmt.Errorf("docker build failed: %w", err)
	}

	fmt.Println("[PreDeploy] Verifying required Fly secrets")
	secretOutput, err := runCommand("flyctl", "secrets", "list", "--app", appName)
	if err != nil {
		return fmt.Errorf("flyctl secrets list failed: %w", err)
	}
	secretSet := parseSecretNames(secretOutput)

	required := []string{"JWT_SECRET", "PLAYFAB_TITLE_ID", "ALLOWED_ORIGINS"}
	for _, key := range required {
		if _, ok := secretSet[key]; !ok {
			return fmt.Errorf("required Fly secret missing: %s", key)
		}
	}

	if !hasAny(secretSet, "DATABASE_URL", "SUPABASE_DATABASE_URL") {
		return errors.New("required DB secret missing: need DATABASE_URL or SUPABASE_DATABASE_URL")
	}
	if !hasAny(secretSet, "REDIS_URL", "UPSTASH_REDIS_URL") {
		return errors.New("required Redis secret missing: need REDIS_URL or UPSTASH_REDIS_URL")
	}

	return nil
}

func runPostDeployChecks(baseURL, appName string) error {
	fmt.Println("[PostDeploy] Verifying deployment completion via Fly.io")
	if _, err := runCommand("flyctl", "status", "--app", appName); err != nil {
		return fmt.Errorf("flyctl status failed: %w", err)
	}
	if _, err := runCommand("flyctl", "releases", "--app", appName); err != nil {
		return fmt.Errorf("flyctl releases failed: %w", err)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	baseURL = strings.TrimRight(baseURL, "/")
	time.Sleep(5 * time.Second)

	fmt.Println("[PostDeploy] Smoke test: /health")
	healthBody, healthCode, err := request(client, "GET", baseURL+"/health", nil, nil)
	if err != nil {
		return fmt.Errorf("request /health: %w", err)
	}
	if healthCode != http.StatusOK {
		return fmt.Errorf("/health returned %d", healthCode)
	}
	var health map[string]any
	if err := json.Unmarshal(healthBody, &health); err != nil {
		return fmt.Errorf("parse /health JSON: %w", err)
	}
	if fmt.Sprintf("%v", health["status"]) != "ok" {
		return fmt.Errorf("/health status expected 'ok', got '%v'", health["status"])
	}

	fmt.Println("[PostDeploy] Smoke test: /swagger/index.html")
	swaggerBody, swaggerCode, err := request(client, "GET", baseURL+"/swagger/index.html", nil, nil)
	if err != nil {
		return fmt.Errorf("request /swagger/index.html: %w", err)
	}
	if swaggerCode != http.StatusOK {
		return fmt.Errorf("/swagger/index.html returned %d", swaggerCode)
	}
	if !containsIgnoreCase(string(swaggerBody), "swagger") {
		return errors.New("/swagger/index.html response does not look like Swagger UI")
	}

	fmt.Println("[PostDeploy] Smoke test: auth endpoint should reject fake token")
	authPayload := []byte(`{"playfabId":"TEST123"}`)
	authHeaders := map[string]string{
		"Content-Type":           "application/json",
		"X-PlayFab-SessionToken": "FAKE_TOKEN",
	}
	_, authCode, err := request(client, "POST", baseURL+"/api/v1/auth/login", authPayload, authHeaders)
	if err != nil {
		return fmt.Errorf("request /api/v1/auth/login: %w", err)
	}
	if authCode != http.StatusUnauthorized && authCode != http.StatusForbidden {
		return fmt.Errorf("expected auth failure (401/403), got %d", authCode)
	}

	fmt.Println("[PostDeploy] Smoke test: player endpoint requires auth")
	_, playerCode, err := request(client, "GET", baseURL+"/api/v1/player/save", nil, nil)
	if err != nil {
		return fmt.Errorf("request /api/v1/player/save: %w", err)
	}
	if playerCode != http.StatusUnauthorized && playerCode != http.StatusForbidden {
		return fmt.Errorf("expected player auth failure (401/403), got %d", playerCode)
	}

	fmt.Println("[PostDeploy] Smoke test: leaderboard endpoint is available")
	_, leaderboardCode, err := request(client, "GET", baseURL+"/api/v1/leaderboard?type=longest_run_days&limit=1", nil, nil)
	if err != nil {
		return fmt.Errorf("request /api/v1/leaderboard: %w", err)
	}
	if leaderboardCode != http.StatusOK {
		return fmt.Errorf("/api/v1/leaderboard returned %d", leaderboardCode)
	}

	return nil
}

func request(client *http.Client, method, url string, payload []byte, headers map[string]string) ([]byte, int, error) {
	var bodyReader io.Reader
	if payload != nil {
		bodyReader = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, 0, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	out := string(output)
	if err != nil {
		return out, fmt.Errorf("%w: %s", err, strings.TrimSpace(out))
	}
	return out, nil
}

func parseSecretNames(output string) map[string]struct{} {
	names := make(map[string]struct{})
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "NAME") || strings.HasPrefix(trimmed, "-") {
			continue
		}
		fields := strings.Fields(trimmed)
		if len(fields) == 0 {
			continue
		}
		names[fields[0]] = struct{}{}
	}
	return names
}

func hasAny(items map[string]struct{}, keys ...string) bool {
	for _, k := range keys {
		if _, ok := items[k]; ok {
			return true
		}
	}
	return false
}

func containsIgnoreCase(text, needle string) bool {
	return strings.Contains(strings.ToLower(text), strings.ToLower(needle))
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "❌ "+format+"\n", args...)
	os.Exit(1)
}
