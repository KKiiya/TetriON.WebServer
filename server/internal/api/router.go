package api

import (
	"net/http"

	"TetriON.WebServer/server/internal/auth"
	"TetriON.WebServer/server/internal/logging"
)

// SetupRoutes registers all API routes to the provided mux
func SetupRoutes(mux *http.ServeMux) {
	logging.LogInfo("Setting up API routes...")

	// Authentication routes
	mux.HandleFunc("/api/auth/register", auth.RegisterHandler)
	mux.HandleFunc("/api/auth/login", auth.LoginHandler)
	mux.HandleFunc("/api/auth/profile", auth.ProfileHandler)

	// Health check
	mux.HandleFunc("/api/health", HealthCheckHandler)

	logging.LogInfo("API routes registered successfully")
}

// HealthCheckHandler returns server health status
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"TetriON.WebServer"}`))
}
