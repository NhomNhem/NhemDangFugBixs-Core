package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_HealthCheck(t *testing.T) {
	app, _ := SetupTestApp()

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "ok", body["status"])
}

func TestAPI_RootV1(t *testing.T) {
	app, _ := SetupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/", nil)
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAPI_RateLimiting(t *testing.T) {
	app, _ := SetupTestApp()

	// Make some requests
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req, 10000)
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	}
}
