package websocket

import (
	"context"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"TetriON.WebServer/server/internal/config"
	"TetriON.WebServer/server/internal/logging"
)

var initialized = false

func Init() {
	if initialized {
		return
	}
	port := config.GetEnv(config.CONFIG_SERVER_PORT)
	logging.LogInfo(port)

	logging.LogInfo("Starting WebSocket on port %s", port)

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			logging.LogError("WebSocket accept error: %v", err)
			return
		}
		defer c.CloseNow()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var v any
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			logging.LogError("Read error: %v", err)
			return
		}

		logging.LogDebug("Received: %v", v)

		c.Close(websocket.StatusNormalClosure, "")
	})
	initialized = true
	logging.LogInfo("WebSocket successfully initialized!")
}

func HandleFunction(path string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(path, handler)
}

func IsInitialized() bool {
	return initialized
}
