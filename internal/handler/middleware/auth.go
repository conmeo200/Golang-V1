package middleware

import (
	"strings"

	"github.com/conmeo200/Golang-V1/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Check if Kong already verified and passed User ID
		kongUserID := c.GetHeader("X-User-ID")
		if kongUserID != "" {
			c.Set("user_id", kongUserID)
			c.Next()
			return
		}

		// 2. Fallback to manual JWT verification
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Fallback to cookie
			cookie, err := c.Cookie("access_token")
			if err == nil && cookie != "" {
				authHeader = "Bearer " + cookie
			} else {
				c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
				return
			}
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := auth.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			userID, _ := claims["user_id"].(string)
			c.Set("user_id", userID)
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid claims"})
		}
	}
}

// WebAuthMiddleware looks for a token in the cookies instead of Authorization header
func WebAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("access_token")
		if err != nil || cookie == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Web Access Denied - Please Login to receive cookie"})
			return
		}

		tokenString := cookie
		token, err := auth.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Web Access Denied - Invalid Token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			userID, _ := claims["user_id"].(string)
			c.Set("user_id", userID)
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid claims"})
		}
	}
}
