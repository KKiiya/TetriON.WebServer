package websocket

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"TetriON.WebServer/server/internal/config"
	"TetriON.WebServer/server/internal/logging"
)

var initialized = false
var endpoints = make(map[string]func(http.ResponseWriter, *http.Request))

func Init() {
	if initialized {
		return
	}
	port := fmt.Sprint(config.GetConfig(config.CONFIG_SERVER_PORT))
	timeoutVal := int(reflect.ValueOf(config.GetConfig(config.CONFIG_SESSION_TIMEOUT)).Int())
	timeout := time.Duration(timeoutVal) * time.Second

	logging.LogInfo("Starting WebSocket on  http://127.0.0.1:%s", port)

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			logging.LogError("WebSocket accept error: %v", err)
			return
		}
		defer c.CloseNow()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
	
	LoadEndpoints()

	initialized = true
	http.ListenAndServe(":"+port, nil)
	logging.LogInfo("WebSocket successfully initialized!")
}

func HandleFunction(path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, exists := endpoints[path]; exists {
		logging.LogError("WebSocket handler for path %s already exists.", path)
		return
	}
	endpoints[path] = handler
	if initialized {
		http.HandleFunc(path, handler)
		logging.LogInfoC(logging.Yellow, "Registered WebSocket endpoint at %s", path)
	}
}

func LoadEndpoints() {
	logging.LogInfoC(logging.Yellow, "Loading WebSocket endpoints...")
	for path, handler := range endpoints {
		http.HandleFunc(path, handler)
		logging.Log(logging.Gray,"Registered WebSocket endpoint: %s", path)
	}
	logging.LogInfo("All WebSocket endpoints loaded.")
}

func IsInitialized() bool {
	return initialized
}
