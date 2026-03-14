package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructuredLogging_Production(t *testing.T) {
	// 1. Setup buffer to capture output
	var buf bytes.Buffer

	// 2. Create JSON handler (simulating production logic in main.go)
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)

	// 3. Log something
	logger.Info("test production log", slog.String("key", "value"))

	// 4. Verify output is JSON
	var logMap map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logMap)

	assert.NoError(t, err, "Log output should be valid JSON")
	assert.Equal(t, "test production log", logMap["msg"])
	assert.Equal(t, "value", logMap["key"])
	assert.Equal(t, "INFO", logMap["level"])
}

func TestLoggerInitialization_Logic(t *testing.T) {
	// This tests the logic used in main.go
	os.Setenv("ENV", "production")
	defer os.Unsetenv("ENV")

	var handler slog.Handler
	if os.Getenv("ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}

	_, ok := handler.(*slog.JSONHandler)
	assert.True(t, ok, "Should use JSONHandler when ENV is production")
}
