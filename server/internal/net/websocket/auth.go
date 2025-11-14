package websocket

import (
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer conn.CloseNow()

	var payload struct {
		Token string `json:"token"`
	}

	if err := wsjson.Read(r.Context(), conn, &payload); err != nil {
		return
	}

	user, err := auth.VerifyToken(payload.Token)
	if err != nil {
		wsjson.Write(r.Context(), conn, map[string]any{
			"error": "invalid_token",
		})
		return
	}

	wsjson.Write(r.Context(), conn, map[string]any{
		"status": "ok",
		"user":   user,
	})
}
