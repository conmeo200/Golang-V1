package middleware

import (
	"net/http"
	"os"
	"strings"
)

// RequireAPIKey checks for the X-API-Key header and validates it against ALLOWED_API_KEYS
func RequireAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		apiKey := r.Header.Get("X-API-Key")
		
		if apiKey == "" {
			http.Error(w, `{"status": false, "message": "Missing X-API-Key header"}`, http.StatusUnauthorized)
			return
		}

		// Read allowed keys from environment
		allowedKeysStr := os.Getenv("ALLOWED_API_KEYS")
		if allowedKeysStr == "" {
			// Fail-safe: If no keys are configured, deny all requests to protected routes
			http.Error(w, `{"status": false, "message": "Server configuration error: No API keys configured"}`, http.StatusInternalServerError)
			return
		}

		allowedKeys := strings.Split(allowedKeysStr, ",")
		
		isValid := false
		for _, key := range allowedKeys {
			if strings.TrimSpace(key) == apiKey {
				isValid = true
				break
			}
		}

		if !isValid {
			http.Error(w, `{"status": false, "message": "Invalid API Key"}`, http.StatusUnauthorized)
			return
		}

		// Valid key, proceed to next handler
		next.ServeHTTP(w, r)
	})
}
