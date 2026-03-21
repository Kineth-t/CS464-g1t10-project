package middleware

import (
	"context"   // Used to store values across request lifecycle
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5" // JWT library

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

// Custom type for context keys (avoids collisions)
type contextKey string

// Keys used to store values in request context
const (
	UserIDKey contextKey = "user_id" // stores logged-in user ID
	RoleKey   contextKey = "role"    // stores user role (admin/user)
)

// RequireAuth ensures the user is authenticated (valid JWT required)
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Parse and validate JWT token
		claims, err := parseToken(r)
		if err != nil {
			http.Error(w, "missing or invalid token", http.StatusUnauthorized)
			return
		}

		// Extract user ID from token claims
		// NOTE: JWT numbers are float64 by default → convert to int
		userID := int(claims["user_id"].(float64))

		// Extract user role from claims
		role := model.Role(claims["role"].(string))

		// Store values in request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)

		// Pass modified request (with context) to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin ensures the user is an admin
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Parse JWT token
		claims, err := parseToken(r)
		if err != nil {
			http.Error(w, "missing or invalid token", http.StatusUnauthorized)
			return
		}

		// Check if role is admin
		if model.Role(claims["role"].(string)) != model.RoleAdmin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		// Extract user ID
		userID := int(claims["user_id"].(float64))

		// Store admin info in context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, model.RoleAdmin)

		// Continue request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// parseToken extracts and validates JWT from Authorization header
func parseToken(r *http.Request) (jwt.MapClaims, error) {

	// Get Authorization header
	header := r.Header.Get("Authorization")

	// Expect format: "Bearer <token>"
	if !strings.HasPrefix(header, "Bearer ") {
		return nil, errors.New("missing token")
	}

	// Remove "Bearer " prefix to get token string
	tokenStr := strings.TrimPrefix(header, "Bearer ")

	// Parse and validate token
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {

		// Ensure token uses HMAC signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		// Get secret key from environment variable
		secret := os.Getenv("JWT_SECRET")

		// Fallback (NOT safe for production)
		if secret == "" {
			secret = "Just-for-developement-phase-change-in-production"
		}

		return []byte(secret), nil
	})

	// If parsing fails or token is invalid
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims (payload data inside JWT)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}