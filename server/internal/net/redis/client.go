package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"TetriON.WebServer/server/internal/net/websocket"
	"github.com/fatih/color"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var redisChannel = "websocket_broadcast"

var (
	cyan   = color.New(color.FgCyan).Add(color.Bold)
	green  = color.New(color.FgGreen).Add(color.Bold)
	yellow = color.New(color.FgYellow).Add(color.Bold)
	red    = color.New(color.FgRed).Add(color.Bold)
	white  = color.New(color.FgWhite)
)

// LogWithTime prints a timestamped message with the given color.
func LogWithTime(c *color.Color, level string, msg string, a ...any) {
	timestamp := time.Now().Format("15:04:05")

	// If there are arguments but the msg contains no formatting verbs, append a `%v`
	if len(a) > 0 && !strings.ContainsAny(msg, "%") {
		msg = msg + " %v"
	}

	formatted := fmt.Sprintf(msg, a...)
	// Print: [HH:MM:SS][LEVEL] EMOJI formatted-message
	c.Printf("[%s] [%s] %s\n", timestamp, level, formatted)
}

// --- Redis initialization ---
func Init() {
	LogWithTime(yellow, "INFO", "🔧 Initializing Redis...")

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
		LogWithTime(yellow, "WARN", "⚠️ No REDIS_ADDR found, using default: %s", redisAddr)
	} else {
		LogWithTime(white, "INFO", "📡 Connecting to Redis at: %s", redisAddr)
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		LogWithTime(yellow, "WARN", "🔒 No Redis password set.")
	}

	redisDBStr := os.Getenv("REDIS_DB")
	redisDB := 0
	if redisDBStr != "" {
		if db, err := strconv.Atoi(redisDBStr); err == nil {
			redisDB = db
			LogWithTime(white, "INFO", "📁 Using Redis DB index: %d", redisDB)
		} else {
			LogWithTime(yellow, "WARN", "⚠️ Invalid REDIS_DB value, defaulting to 0")
		}
	} else {
		LogWithTime(white, "INFO", "📁 Using default Redis DB index: 0")
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	ctx := context.Background()
	ping, err := redisClient.Ping(ctx).Result()
	if err != nil {
		LogWithTime(red, "ERROR", "❌ Failed to connect to Redis: %v", err)
		return
	}

	LogWithTime(green, "SUCCESS", "✅ Connected to Redis at %s successfully. (%s)", redisAddr, ping)
}

func PublishMessage(ctx context.Context, message string) error {
	if redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	err := redisClient.Publish(ctx, redisChannel, message).Err()
	if err != nil {
		LogWithTime(red, "ERROR", "❌ Failed to publish message: %v", err)
		return err
	}
	LogWithTime(green, "INFO", "📤 Published to channel '%s': %s", redisChannel, message)
	return nil
}

// --- Subscribe and listen for messages ---
func SubscribeMessages(ctx context.Context) {
	if redisClient == nil {
		LogWithTime(red, "ERROR", "❌ SubscribeMessages called before Redis initialization")
		return
	}

	pubsub := redisClient.Subscribe(ctx, redisChannel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	LogWithTime(cyan, "INFO", "📡 Subscribed to channel '%s'", redisChannel)

	for msg := range ch {
		LogWithTime(white, "RECV", "📨 Received message: %s", msg.Payload)
		websocket.Broadcast(map[string]any{
			"type":      "redis_broadcast",
			"channel":   msg.Channel,
			"payload":   msg.Payload,
			"timestamp": time.Now().Unix(),
		})
	}

	LogWithTime(yellow, "INFO", "❌ Subscription closed for channel '%s'", redisChannel)
}

func Close() {
	if redisClient == nil {
		return
	}
	if err := redisClient.Close(); err != nil {
		LogWithTime(red, "ERROR", "❌ Error closing Redis client: %v", err)
		return
	}
	LogWithTime(white, "INFO", "🔒 Redis connection closed.")
	redisClient = nil
}
