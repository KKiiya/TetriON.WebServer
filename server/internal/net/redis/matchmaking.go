package redis

import (
	"context"
	"fmt"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
)

const matchmakingQueuePrefix = "mm:queue:"

func EnqueuePlayer(ctx context.Context, queue string, userID string, skill int) error {
	if redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	key := matchmakingQueuePrefix + queue
	score := float64(skill)*1_000_000 + float64(time.Now().Unix())

	return redisClient.ZAdd(ctx, key, redisv9.Z{Score: score, Member: userID}).Err()
}

func RemovePlayerFromQueue(ctx context.Context, queue string, userID string) error {
	if redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	key := matchmakingQueuePrefix + queue
	return redisClient.ZRem(ctx, key, userID).Err()
}

func PeekPlayers(ctx context.Context, queue string, limit int64) ([]string, error) {
	if redisClient == nil {
		return nil, fmt.Errorf("redis client is not initialized")
	}
	if limit <= 0 {
		limit = 10
	}

	key := matchmakingQueuePrefix + queue
	return redisClient.ZRange(ctx, key, 0, limit-1).Result()
}

func QueueSize(ctx context.Context, queue string) (int64, error) {
	if redisClient == nil {
		return 0, fmt.Errorf("redis client is not initialized")
	}

	key := matchmakingQueuePrefix + queue
	return redisClient.ZCard(ctx, key).Result()
}
