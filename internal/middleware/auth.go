package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/conmeo200/Golang-V1/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Check if Kong already verified and passed User ID
		kongUserID := r.Header.Get("X-User-ID")
		if kongUserID != "" {
			ctx := context.WithValue(r.Context(), "user_id", kongUserID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// 2. Fallback to manual JWT verification (Double check or local dev)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := auth.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			userID, _ := claims["user_id"].(string)
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// WebAuthMiddleware looks for a token in the cookies instead of Authorization header
func WebAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			// Redirect to login page if no token is found
			// For now, we just deny access since we don't have a web login page yet
			http.Error(w, "Web Access Denied - Please Login to receive cookie", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		token, err := auth.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Web Access Denied - Invalid Token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			userID, _ := claims["user_id"].(string)
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}