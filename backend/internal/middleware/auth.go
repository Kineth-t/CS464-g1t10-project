package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := parseToken(r)
		if err != nil {
			http.Error(w, "missing or invalid token", http.StatusUnauthorized)
			return
		}
		userID := int(claims["user_id"].(float64))
		role := model.Role(claims["role"].(string))
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := parseToken(r)
		if err != nil {
			http.Error(w, "missing or invalid token", http.StatusUnauthorized)
			return
		}
		if model.Role(claims["role"].(string)) != model.RoleAdmin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		userID := int(claims["user_id"].(float64))
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, model.RoleAdmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseToken(r *http.Request) (jwt.MapClaims, error) {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return nil, errors.New("missing token")
	}
	tokenStr := strings.TrimPrefix(header, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "change-me-in-production"
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}