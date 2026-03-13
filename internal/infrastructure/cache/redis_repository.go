package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	client *redis.Client
}

// NewRedisRepository creates a new Redis repository
func NewRedisRepository(client *redis.Client) interface {
	repository.CacheRepository
	repository.TokenRepository
} {
	return &redisRepository{client: client}
}

// CacheRepository implementation

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	if r.client == nil {
		return "", nil
	}
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *redisRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if r.client == nil {
		return nil
	}
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redisRepository) Delete(ctx context.Context, key string) error {
	if r.client == nil {
		return nil
	}
	return r.client.Del(ctx, key).Err()
}

// TokenRepository implementation

func (r *redisRepository) StoreRefreshToken(ctx context.Context, token string, userID string, ttl time.Duration) error {
	if r.client == nil {
		return nil
	}
	key := fmt.Sprintf("refresh_token:%s", token)
	return r.client.Set(ctx, key, userID, ttl).Err()
}

func (r *redisRepository) GetRefreshToken(ctx context.Context, token string) (string, error) {
	if r.client == nil {
		return "MOCK_USER_ID", nil
	}
	key := fmt.Sprintf("refresh_token:%s", token)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *redisRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	if r.client == nil {
		return nil
	}
	key := fmt.Sprintf("refresh_token:%s", token)
	return r.client.Del(ctx, key).Err()
}

func (r *redisRepository) BlacklistJWT(ctx context.Context, jti string, ttl time.Duration) error {
	if r.client == nil {
		return nil
	}
	key := fmt.Sprintf("session:%s:blacklist", jti)
	return r.client.Set(ctx, key, "1", ttl).Err()
}

func (r *redisRepository) IsJWTBlacklisted(ctx context.Context, jti string) (bool, error) {
	if r.client == nil {
		return false, nil
	}
	key := fmt.Sprintf("session:%s:blacklist", jti)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}
