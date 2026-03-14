package persistence

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresLevelRepository_GetConfig_Mock(t *testing.T) {
	repo := NewPostgresLevelRepository(nil)
	ctx := context.Background()

	config, err := repo.GetConfig(ctx, "level_1", "map_1")

	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "level_1", config.LevelID)
	assert.Equal(t, 10.0, config.MinTimeSeconds)
}

func TestPostgresTalentRepository_GetConfigs_Mock(t *testing.T) {
	repo := NewPostgresTalentRepository(nil)
	ctx := context.Background()

	configs, err := repo.GetConfigs(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, configs)
	assert.Greater(t, len(configs), 0)
	assert.NotNil(t, configs["health"])
}
