package redis

import (
	"context"
	"fmt"
	"time"
)

const presencePrefix = "presence:user:"

func SetUserOnline(ctx context.Context, userID string, ttl time.Duration) error {
	if redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	if ttl <= 0 {
		ttl = 60 * time.Second
	}

	key := presencePrefix + userID
	return redisClient.Set(ctx, key, "online", ttl).Err()
}

func HeartbeatUser(ctx context.Context, userID string, ttl time.Duration) error {
	if redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	if ttl <= 0 {
		ttl = 60 * time.Second
	}

	key := presencePrefix + userID
	return redisClient.Expire(ctx, key, ttl).Err()
}

func SetUserOffline(ctx context.Context, userID string) error {
	if redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	key := presencePrefix + userID
	return redisClient.Del(ctx, key).Err()
}

func IsUserOnline(ctx context.Context, userID string) (bool, error) {
	if redisClient == nil {
		return false, fmt.Errorf("redis client is not initialized")
	}

	key := presencePrefix + userID
	exists, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
