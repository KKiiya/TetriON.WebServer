package api

import (
	"encoding/json"
	"net/http"
)

// JSONResponse is a helper to send JSON responses
func JSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// ErrorResponse sends an error response in JSON format
func ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	JSONResponse(w, map[string]interface{}{
		"success": false,
		"error":   message,
	}, statusCode)
}
