package auth

import (
	"encoding/json"
	"net/http"

	"TetriON.WebServer/server/internal/logging"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.LogError("Failed to decode register request: %v", err)
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Register user
	user, token, err := Register(req.Username, req.Email, req.Password)
	if err != nil {
		logging.LogWarning("Registration failed: %v", err)
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	logging.LogInfo("User registered successfully: %s (ID: %s)", user.Username, user.ID)

	// Remove password hash from response
	user.PasswordHash = ""

	respondJSON(w, AuthResponse{
		Success: true,
		Message: "Registration successful",
		Token:   token,
		User:    user,
	}, http.StatusCreated)
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.LogError("Failed to decode login request: %v", err)
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Authenticate user
	user, token, err := Login(req.Username, req.Password)
	if err != nil {
		logging.LogWarning("Login failed for user %s: %v", req.Username, err)
		respondError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	logging.LogInfo("User logged in successfully: %s", user.Username)

	// Remove password hash from response
	user.PasswordHash = ""

	respondJSON(w, AuthResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User:    user,
	}, http.StatusOK)
}

// ProfileHandler returns the current user's profile (protected endpoint)
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract token from Authorization header
	token := extractToken(r)
	if token == "" {
		respondError(w, "No authorization token provided", http.StatusUnauthorized)
		return
	}

	// Validate token and get user
	user, err := ValidateToken(token)
	if err != nil {
		logging.LogWarning("Invalid token: %v", err)
		respondError(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Remove password hash from response
	user.PasswordHash = ""

	respondJSON(w, map[string]interface{}{
		"success": true,
		"user":    user,
	}, http.StatusOK)
}

// Helper functions

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Expected format: "Bearer <token>"
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}

	return authHeader
}

func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	respondJSON(w, ErrorResponse{
		Success: false,
		Error:   message,
	}, statusCode)
}
