package db

import (
	"context"
	"fmt"
	"time"

	"strconv"

	"TetriON.WebServer/server/internal/config"
	"TetriON.WebServer/server/internal/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func Close() {
	if DB != nil {
		DB.Close()
		logging.LogWithTime(logging.White, "INFO", "🔒 PostgreSQL database connection closed.")
	}
}

func GetData(query string, args ...interface{}) (pgx.Rows, error) {
	return DB.Query(context.Background(), query, args...)
}

func ExecCommand(command string, args ...interface{}) (pgconn.CommandTag, error) {
	return DB.Exec(context.Background(), command, args...)
}

func GetRow(query string, args ...interface{}) pgx.Row {
	return DB.QueryRow(context.Background(), query, args...)
}

func GetColumn(query string, args ...interface{}) ([]interface{}, error) {
	rows, err := DB.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []any
	for rows.Next() {
		var value interface{}
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		results = append(results, value)
	}
	return results, nil
}
