package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"TetriON.WebServer/server/internal/redis"
	"github.com/fatih/color"
)

var (
	cyan   = color.New(color.FgCyan).Add(color.Bold)
	green  = color.New(color.FgGreen).Add(color.Bold)
	yellow = color.New(color.FgYellow).Add(color.Bold)
	red    = color.New(color.FgRed).Add(color.Bold)
	white  = color.New(color.FgWhite)
)

var config any // replace with your actual config struct

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

func main() {
	cyan.Println("======================================================================")
	cyan.Println("		 ______    __      _ ____  _  ____")
	cyan.Println("		/_  __/__ / /_____(_) __ \\/ |/ / /")
	cyan.Println("		 / / / -_) __/ __/ / /_/ /    /_/ ")
	cyan.Println("		/_/  \\__/\\__/_/ /_/\\____/_/|_(_)  ")
	cyan.Println("		  							  ")
	cyan.Println("======================================================================")

	fmt.Println()
	LogWithTime(green, "INFO", "🚀 Starting server initialization...")
	fmt.Println()

	LoadEnv()
	LoadConfig()
	redis.InitRedis()
	fmt.Println()

	LogWithTime(green, "INFO", "✅ All systems initialized successfully!")
	fmt.Println("======================================================================")
}

func LoadEnv() {
	LogWithTime(yellow, "INFO", "⚙️  Loading environment variables (.env)...")

	file, err := os.OpenFile("../.env", os.O_RDONLY, 0644)
	if err != nil {
		LogWithTime(red, "ERROR", "❌ Failed to load .env file: %v", err)
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
		LogWithTime(red, "ERROR", "❌ Error reading .env file: %v", err)
		return
	}

	LogWithTime(green, "INFO", "✅ Loaded %d environment variables successfully.", count)
	fmt.Println()
}

func LoadConfig() {
	LogWithTime(yellow, "INFO", "🧩 Loading configuration (config.json)...")

	file, err := os.OpenFile("../config.json", os.O_RDONLY, 0644)
	if err != nil {
		LogWithTime(red, "ERROR", "❌ Failed to open config.json: %v", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		LogWithTime(red, "ERROR", "❌ Error decoding config.json: %v", err)
		return
	}

	LogWithTime(green, "INFO", "✅ Configuration loaded successfully.")
	LogWithTime(white, "INFO", "📋 Configuration details:")
	for key, value := range config.(map[string]any) {
		white.Printf("			%s: %+v\n", key, value)
	}
	fmt.Println()
}