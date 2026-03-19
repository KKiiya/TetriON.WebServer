package api

import (
	"net/http"

	"TetriON.WebServer/server/internal/admin"
	"TetriON.WebServer/server/internal/auth"
	"TetriON.WebServer/server/internal/logging"
	"TetriON.WebServer/server/internal/metrics"
	"TetriON.WebServer/server/internal/middleware"
)

// SetupRoutes registers all API routes to the provided mux
func SetupRoutes(mux *http.ServeMux) {
	logging.LogInfo("Setting up API routes...")
	chain := func(h http.Handler) http.Handler {
		return middleware.RateLimitMiddleware(metrics.HTTPMetricsMiddleware(h))
	}

	// Authentication routes
	mux.Handle("/api/auth/register", chain(http.HandlerFunc(auth.RegisterHandler)))
	mux.Handle("/api/auth/login", chain(http.HandlerFunc(auth.LoginHandler)))
	mux.Handle("/api/auth/profile", chain(middleware.RequireAuth(http.HandlerFunc(auth.ProfileHandler))))

	// Health check
	mux.Handle("/api/health", chain(http.HandlerFunc(HealthCheckHandler)))

	// Metrics and admin routes
	mux.Handle("/api/metrics", chain(http.HandlerFunc(metrics.Handler)))
	mux.Handle("/api/admin/health", chain(http.HandlerFunc(admin.HealthHandler)))
	mux.Handle("/api/admin/stats", chain(http.HandlerFunc(admin.StatsHandler)))
	mux.Handle("/api/admin/config", chain(http.HandlerFunc(admin.ConfigHandler)))

	logging.LogInfo("API routes registered successfully")
}

// HealthCheckHandler returns server health status
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"TetriON.WebServer"}`))
}
