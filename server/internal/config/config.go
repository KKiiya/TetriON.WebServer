package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"TetriON.WebServer/server/internal/logging"
)

var (
	config = make(map[string]any)
)

// Configuration keys
var (
	CONFIG_SERVER_PORT 		= "server_port"
	CONFIG_LOG_LEVEL    	= "log_level"
	CONFIG_MAX_CONNECTIONS 	= "max_connections"
	CONFIG_DEBUG_MODE      	= "debug_mode"
	CONFIG_ALLOW_ORIGINS   	= "allow_origins"
	CONFIG_SESSION_TIMEOUT  = "session_timeout"
	CONFIG_EMAIL_SERVICE	= "email_service"
	CONFIG_CACHE_TTL        = "cache_ttl"
)

// Environment variable keys
var (
	ENV_REDIS_ADDRESS       = "REDIS_ADDR"
	ENV_REDIS_PASSWORD      = "REDIS_PASSWORD"
	ENV_POSTGRES_USER       = "POSTGRES_USER"
	ENV_POSTGRES_PASSWORD   = "POSTGRES_PASSWORD"
	ENV_POSTGRES_HOST       = "POSTGRES_HOST"
	ENV_POSTGRES_PORT       = "POSTGRES_PORT"
	ENV_POSTGRES_DBNAME     = "POSTGRES_DBNAME"
	ENV_POSTGRES_SSLMODE    = "POSTGRES_SSLMODE"
	ENV_JWT_SECRET		 	= "JWT_SECRET"
)

func LoadEnv() {
	logging.LogWithTime(logging.Yellow, "INFO", "⚙️  Loading environment variables (.env)...")

	file, err := os.OpenFile("../../.env", os.O_RDONLY, 0644)
	if err != nil {
		logging.LogWithTime(logging.Red, "ERROR", "❌ Failed to load .env file: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
		count++
	}

	if err := scanner.Err(); err != nil {
		logging.LogWithTime(logging.Red, "ERROR", "❌ Error reading .env file: %v", err)
		return
	}

	logging.LogWithTime(logging.Green, "INFO", "✅ Loaded %d environment variables successfully.", count)
	logging.LogLine(logging.White, "")
}

func LoadConfig() {
	logging.LogWithTime(logging.Yellow, "INFO", "🧩 Loading configuration (config.json)...")

	file, err := os.OpenFile("../../config.json", os.O_RDONLY, 0644)
	if err != nil {
		logging.LogWithTime(logging.Red, "ERROR", "❌ Failed to open config.json: %v", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		logging.LogWithTime(logging.Red, "ERROR", "❌ Error decoding config.json: %v", err)
		return
	}

	logging.LogWithTime(logging.Green, "INFO", "✅ Configuration loaded successfully.")
	logging.LogWithTime(logging.White, "INFO", "📋 Configuration details:")
	for key, value := range config {
		logging.LogLine(logging.White, "			" + key, value)
	}
	fmt.Println()
}

func GetConfig(key string) any {
	return config[key]
}

func GetAllConfig() map[string]any {
	return config
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}	