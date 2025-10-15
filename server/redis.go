package main

import (
	"context"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var redisChannel = "websocket_broadcast"

func InitRedis() {
	green.Println("Initializing Redis...")
	redisAddr := os.Getenv("REDIS_ADDR")
	green.Println("Connecting to Redis at:", redisAddr)
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default address
		yellow.Println("Using default Redis address:", redisAddr)
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		red.Println("No Redis password set.")
	}
	redisDBStr := os.Getenv("REDIS_DB")
	redisDB := 0
	if redisDBStr != "" {
		if db, err := strconv.Atoi(redisDBStr); err == nil {
			redisDB = db
		}
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
	

	ctx := context.Background()
	ping, err := redisClient.Ping(ctx).Result()
	if err != nil {
		red.Println("Error connecting to Redis:", err)
		return
	}
	green.Println("Connected to Redis", redisAddr, "successfully. (", ping, ")")
}