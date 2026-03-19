package main

import (
	"context"

	"TetriON.WebServer/server/internal/config"
	"TetriON.WebServer/server/internal/db"
	"TetriON.WebServer/server/internal/logging"
	"TetriON.WebServer/server/internal/net/redis"
	"TetriON.WebServer/server/internal/net/websocket"
	"TetriON.WebServer/server/internal/worker"
)

func main() {
	logging.LogLine(logging.Cyan, "======================================================================")
	logging.LogLine(logging.Cyan, "		 ______    __      _ ____  _  ____")
	logging.LogLine(logging.Cyan, "		/_  __/__ / /_____(_) __ \\/ |/ / /")
	logging.LogLine(logging.Cyan, "		 / / / -_) __/ __/ / /_/ /    /_/ ")
	logging.LogLine(logging.Cyan, "		/_/  \\__/\\__/_/ /_/\\____/_/|_(_)  ")
	logging.LogLine(logging.Cyan, "		  							  ")
	logging.LogLine(logging.Cyan, "======================================================================")

	logging.LogLine(logging.White, "")
	logging.LogWithTime(logging.White, "DEBUG", "🚀 Starting server initialization...")
	logging.LogLine(logging.White, "")
	if err := logging.Init(); err != nil {
		logging.LogWarning("Log file output is disabled: %v", err)
	}

	config.LoadEnv()
	config.LoadConfig()

	redis.Init()
	db.Init()
	websocket.Init()

	logging.LogWithTime(logging.Green, "INFO", "✅ All systems initialized successfully!")
	logging.LogLine(logging.White, "======================================================================")
	logging.LogLine(logging.White, "")

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	keyspaceSub := worker.NewKeyspaceSubscriber()
	keyspaceSub.Start(rootCtx)

	if err := redis.PublishMessage(context.Background(), "REDIS ON!"); err != nil {
		logging.LogWarning("Unable to publish startup message to Redis: %v", err)
	}

	Listen()

	cancel()
	keyspaceSub.Stop()
	websocket.Stop()
	db.Close()
	redis.Close()
}
