package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var redisChannel = "websocket_broadcast"

func InitRedis() {
	fmt.Println("Initializing Redis...")
	redisAddr := os.Getenv("REDIS_ADDR")
	fmt.Println("Connecting to Redis at:", redisAddr)
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default address
		fmt.Println("Using default Redis address:", redisAddr)
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	fmt.Println(redisPassword)
	if redisPassword == "" {
		redisPassword = "yourpassword" // Default password
		fmt.Println("Using default Redis password:", redisPassword)
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
		fmt.Println("Error connecting to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis", redisAddr, "successfully. (", ping, ")")
}