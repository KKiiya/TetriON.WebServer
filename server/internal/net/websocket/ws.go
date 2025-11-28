package websocket

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"TetriON.WebServer/server/internal/api"
	"TetriON.WebServer/server/internal/config"
	"TetriON.WebServer/server/internal/logging"
)

var (
	initialized = false
	endpoints   = make(map[string]func(http.ResponseWriter, *http.Request))
	server      *http.Server // <--- keep reference to server
)

func Init() {
	if initialized {
		return
	}

	port := fmt.Sprint(config.GetConfig(config.CONFIG_SERVER_PORT))
	logging.LogInfo("Starting WebSocket on http://127.0.0.1:%s", port)

	timeoutVal, _ := strconv.Atoi(fmt.Sprint(config.GetConfig(config.CONFIG_SESSION_TIMEOUT)))
	timeout := time.Duration(timeoutVal) * time.Second

	// Create a custom multiplexer
	mux := http.NewServeMux()

	// WebSocket main endpoint
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			logging.LogError("WebSocket accept error: %v", err)
			return
		}
		defer c.CloseNow()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		var v any
		if err := wsjson.Read(ctx, c, &v); err != nil {
			logging.LogError("Read error: %v", err)
			return
		}

		logging.LogDebug("Received: %v", v)
		c.Close(websocket.StatusNormalClosure, "")
	})

	// Register API routes (authentication, health check, etc.)
	api.SetupRoutes(mux)

	// Register other endpoints
	LoadEndpoints(mux)

	// Create and store the server
	server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	initialized = true

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.LogError("WebSocket server error: %v", err)
		}
	}()

	logging.LogInfo("WebSocket successfully initialized!")
}

func HandleFunction(path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, exists := endpoints[path]; exists {
		logging.LogError("WebSocket handler for path %s already exists.", path)
		return
	}
	endpoints[path] = handler
	if initialized && server != nil {
		// note: we need to rebuild mux if adding dynamically after init
		logging.LogInfoC(logging.Yellow, "Registered WebSocket endpoint at %s (pending mux reload)", path)
	}
}

func LoadEndpoints(mux *http.ServeMux) {
	logging.LogInfoC(logging.Yellow, "Loading WebSocket endpoints...")
	for path, handler := range endpoints {
		mux.HandleFunc(path, handler)
		logging.Log(logging.Gray, "Registered WebSocket endpoint: %s", path)
	}
	logging.LogInfo("All WebSocket endpoints loaded.")
}

func Stop() {
	if !initialized || server == nil {
		logging.LogWarning("WebSocket server not running or already stopped.")
		return
	}

	logging.LogInfo("Stopping WebSocket server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logging.LogError("Error shutting down WebSocket server: %v", err)
		return
	}

	initialized = false
	logging.LogInfo("WebSocket server stopped cleanly.")
}

func IsInitialized() bool {
	return initialized
}
