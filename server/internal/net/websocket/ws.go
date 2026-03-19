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
	server      *http.Server
	hub         *Hub
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
	hub = NewHub()
	go hub.Run()

	// WebSocket endpoints
	mux.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWSClient(w, r, timeout)
	})
	mux.HandleFunc("/api/ws/auth", AuthWSHandler)

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
		// This map is applied only during Init().
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

	hub = nil
	initialized = false
	logging.LogInfo("WebSocket server stopped cleanly.")
}

func IsInitialized() bool {
	return initialized
}

func Broadcast(payload any) {
	if hub == nil {
		return
	}
	hub.Broadcast(payload)
}

func handleWSClient(w http.ResponseWriter, r *http.Request, timeout time.Duration) {
	if hub == nil {
		http.Error(w, "websocket hub not initialized", http.StatusServiceUnavailable)
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		logging.LogError("WebSocket accept error: %v", err)
		return
	}

	clientID := fmt.Sprintf("%s-%d", r.RemoteAddr, time.Now().UnixNano())
	client := NewClient(clientID, conn)
	hub.Register(client)

	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()
	defer conn.Close(websocket.StatusNormalClosure, "bye")
	defer hub.Unregister(client)

	_ = wsjson.Write(ctx, conn, map[string]any{
		"type":    "welcome",
		"message": "connected",
	})

	go client.WritePump(ctx)
	client.ReadPump(ctx, func(v any) {
		msg := map[string]any{
			"type":      "ws_message",
			"client_id": clientID,
			"payload":   v,
			"timestamp": time.Now().Unix(),
		}
		hub.Broadcast(msg)
	})
}
