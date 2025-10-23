package db

import (
	"context"
	"fmt"
	"time"

	"strconv"

	"TetriON.WebServer/server/internal/config"
	"TetriON.WebServer/server/internal/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Init() {
	logging.LogWithTime(logging.Yellow, "INFO", "🔧 Initializing PostgreSQL database connection...")
	port, err := strconv.Atoi(config.GetEnv(config.ENV_POSTGRES_PORT))
	if err != nil {
		logging.LogWithTime(logging.Red, "ERROR", "❌ Invalid port number: %v", err)
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.GetEnv(config.ENV_POSTGRES_USER),   
		config.GetEnv(config.ENV_POSTGRES_PASSWORD),
		config.GetEnv(config.ENV_POSTGRES_HOST),    
		port,
		config.GetEnv(config.ENV_POSTGRES_DBNAME),
		config.GetEnv(config.ENV_POSTGRES_SSLMODE),
	)
	logging.LogWithTime(logging.White, "INFO", "📡 Connecting to PostgreSQL at: %s", dsn)
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		logging.LogWithTime(logging.Red, "ERROR", "❌ Unable to connect to database: %v", err)
		return
	}
	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = DB.Ping(ctx)
	if err != nil {
		logging.LogWithTime(logging.Red, "ERROR", "❌ Unable to ping the database: %v", err)
		return
	}
	logging.LogWithTime(logging.Green, "INFO", "✅ Connected to PostgreSQL database successfully.")
}
