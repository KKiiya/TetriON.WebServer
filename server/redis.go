package main

import (
	"context"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var redisChannel = "websocket_broadcast"

// --- Redis initialization ---
func InitRedis() {
	logWithTime(yellow, "INFO", "🔧 Initializing Redis...")

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
		logWithTime(yellow, "WARN", "⚠️ No REDIS_ADDR found, using default: %s", redisAddr)
	} else {
		logWithTime(white, "INFO", "📡 Connecting to Redis at: %s", redisAddr)
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		logWithTime(yellow, "WARN", "🔒 No Redis password set.")
	}

	redisDBStr := os.Getenv("REDIS_DB")
	redisDB := 0
	if redisDBStr != "" {
		if db, err := strconv.Atoi(redisDBStr); err == nil {
			redisDB = db
			logWithTime(white, "INFO", "📁 Using Redis DB index: %d", redisDB)
		} else {
			logWithTime(yellow, "WARN", "⚠️ Invalid REDIS_DB value, defaulting to 0")
		}
	} else {
		logWithTime(white, "INFO", "📁 Using default Redis DB index: 0")
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	ctx := context.Background()
	ping, err := redisClient.Ping(ctx).Result()
	if err != nil {
		logWithTime(red, "ERROR", "❌ Failed to connect to Redis: %v", err)
		return
	}

	logWithTime(green, "SUCCESS", "✅ Connected to Redis at %s successfully. (%s)", redisAddr, ping)
}