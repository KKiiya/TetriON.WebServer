package websocket

import (
	"net/http"

	"TetriON.WebServer/server/internal/auth"
	"TetriON.WebServer/server/internal/logging"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// AuthWSHandler handles WebSocket authentication
func AuthWSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		logging.LogError("WebSocket accept error: %v", err)
		return
	}
	defer conn.CloseNow()

	var payload struct {
		Token string `json:"token"`
	}

	if err := wsjson.Read(r.Context(), conn, &payload); err != nil {
		logging.LogError("Failed to read WebSocket auth payload: %v", err)
		return
	}

	// Verify token using auth package
	user, err := auth.ValidateToken(payload.Token)
	if err != nil {
		logging.LogWarning("Invalid token in WebSocket auth: %v", err)
		wsjson.Write(r.Context(), conn, map[string]any{
			"success": false,
			"error":   "invalid_token",
		})
		return
	}

	logging.LogInfo("User authenticated via WebSocket: %s", user.Username)

	// Remove sensitive data
	user.PasswordHash = ""

	wsjson.Write(r.Context(), conn, map[string]any{
		"success": true,
		"user":    user,
	})
}
