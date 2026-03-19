package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type visitor struct {
	count       int
	windowStart time.Time
	lastSeen    time.Time
}

var (
	visitorsMu sync.Mutex
	visitors   = map[string]*visitor{}

	rateLimitWindow  = 1 * time.Minute
	rateLimitMaxHits = 120
)

func init() {
	go cleanupVisitors()
}

// RateLimitMiddleware applies a simple per-IP fixed-window limiter.
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/health" {
			next.ServeHTTP(w, r)
			return
		}

		if !allowRequest(clientIP(r)) {
			writeRateLimitError(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func allowRequest(ip string) bool {
	now := time.Now()

	visitorsMu.Lock()
	defer visitorsMu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		visitors[ip] = &visitor{count: 1, windowStart: now, lastSeen: now}
		return true
	}

	v.lastSeen = now
	if now.Sub(v.windowStart) >= rateLimitWindow {
		v.count = 1
		v.windowStart = now
		return true
	}

	v.count++
	return v.count <= rateLimitMaxHits
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for now := range ticker.C {
		visitorsMu.Lock()
		for ip, v := range visitors {
			if now.Sub(v.lastSeen) > 10*time.Minute {
				delete(visitors, ip)
			}
		}
		visitorsMu.Unlock()
	}
}

func writeRateLimitError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error":   "rate limit exceeded",
	})
}
