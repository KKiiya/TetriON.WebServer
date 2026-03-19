package metrics

import (
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	startedAt      = time.Now()
	totalRequests  uint64
	totalErrors    uint64
	activeRequests int64
	pathCountersMu sync.Mutex
	pathCounters   = map[string]uint64{}
	methodCounters = map[string]uint64{}
	statusCounters = map[int]uint64{}
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// HTTPMetricsMiddleware records request counters and basic response statistics.
func HTTPMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&totalRequests, 1)
		atomic.AddInt64(&activeRequests, 1)
		defer atomic.AddInt64(&activeRequests, -1)

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		if rec.status >= 400 {
			atomic.AddUint64(&totalErrors, 1)
		}

		pathCountersMu.Lock()
		pathCounters[r.URL.Path]++
		methodCounters[r.Method]++
		statusCounters[rec.status]++
		pathCountersMu.Unlock()
	})
}

// Handler serves current in-memory metrics as JSON.
func Handler(w http.ResponseWriter, _ *http.Request) {
	pathCountersMu.Lock()
	paths := make(map[string]uint64, len(pathCounters))
	methods := make(map[string]uint64, len(methodCounters))
	statuses := make(map[int]uint64, len(statusCounters))
	for k, v := range pathCounters {
		paths[k] = v
	}
	for k, v := range methodCounters {
		methods[k] = v
	}
	for k, v := range statusCounters {
		statuses[k] = v
	}
	pathCountersMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"uptime_seconds":  int64(time.Since(startedAt).Seconds()),
		"active_requests": atomic.LoadInt64(&activeRequests),
		"total_requests":  atomic.LoadUint64(&totalRequests),
		"total_errors":    atomic.LoadUint64(&totalErrors),
		"by_path":         paths,
		"by_method":       methods,
		"by_status":       statuses,
	})
}
