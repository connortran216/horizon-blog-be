package views

import (
	"go-crud/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	UserContextKey = "user_id"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check for Bearer token format
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Use 'Bearer <token>'"})
			c.Abort()
			return
		}

		tokenString := bearerToken[1]

		// Validate token
		userID, err := services.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store user ID in context for use in handlers
		c.Set(UserContextKey, userID)
		c.Next()
	}
}

// GetUserIDFromContext retrieves the authenticated user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(UserContextKey)
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	return id, ok
}
