package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/olad5/AfriHacks2023-stressless-backend/config"
	"go.uber.org/zap"
)

type RedisCache struct {
	Client *redis.Client
	logger *zap.Logger
}

var ttl = time.Minute * 30

func New(ctx context.Context, configurations *config.Configurations, logger *zap.Logger) (*RedisCache, error) {
	opt, err := redis.ParseURL(configurations.CacheAddress)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		Client: client,
		logger: logger,
	}, nil
}

func (r *RedisCache) SetOne(ctx context.Context, key, value string) error {
	_, err := r.Client.Set(ctx, key, value, ttl).Result()
	if err != nil {
		return fmt.Errorf("Error setting value in cache: %w", err)
	}
	return nil
}

func (r *RedisCache) GetOne(ctx context.Context, key string) (string, error) {
	result, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("Error getting value from cache: %w", err)
	}
	return result, nil
}

func (r *RedisCache) DeleteOne(ctx context.Context, key string) error {
	_, err := r.Client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("Error deleting key in cache: %w", err)
	}
	return nil
}

func (r *RedisCache) Ping(ctx context.Context) error {
	if err := r.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Failed to Ping Redis Cache: %v", err)
	}
	return nil
}
