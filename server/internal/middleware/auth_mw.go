package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userContextKey contextKey = "authenticated_user"

type AuthUser struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type tokenClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// RequireAuth validates a bearer token and forwards user claims in request context.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r.Header.Get("Authorization"))
		if token == "" {
			writeJSONError(w, "missing authorization token", http.StatusUnauthorized)
			return
		}

		user, err := validateToken(token)
		if err != nil {
			writeJSONError(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserFromContext returns the authenticated user set by RequireAuth.
func UserFromContext(ctx context.Context) (*AuthUser, bool) {
	user, ok := ctx.Value(userContextKey).(*AuthUser)
	return user, ok
}

func extractBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	const prefix = "Bearer "
	if len(authHeader) > len(prefix) && authHeader[:len(prefix)] == prefix {
		return authHeader[len(prefix):]
	}
	return authHeader
}

func validateToken(tokenString string) (*AuthUser, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET not configured")
	}

	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return &AuthUser{
		UserID:   claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
	}, nil
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error":   message,
	})
}
