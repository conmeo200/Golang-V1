package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var allowedKeysStr = os.Getenv("ALLOWED_API_KEYS")

// RequireAPIKey checks for the X-API-Key header and validates it against ALLOWED_API_KEYS
func RequireAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Missing X-API-Key header"})
			return
		}

		if allowedKeysStr == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Server configuration error: No API keys configured"})
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Invalid API Key"})
			return
		}

		c.Next()
	}
}
