package admin

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"TetriON.WebServer/server/internal/config"
)

var startedAt = time.Now()

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok","scope":"admin"}`))
}

func StatsHandler(w http.ResponseWriter, _ *http.Request) {
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)

	payload := map[string]any{
		"uptime_seconds": int64(time.Since(startedAt).Seconds()),
		"go_routines":    runtime.NumGoroutine(),
		"alloc_bytes":    mem.Alloc,
		"total_alloc":    mem.TotalAlloc,
		"sys_bytes":      mem.Sys,
	}

	writeJSON(w, payload, http.StatusOK)
}

func ConfigHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{
		"config": config.GetAllConfig(),
	}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
