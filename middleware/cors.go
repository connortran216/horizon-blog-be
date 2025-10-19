package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles CORS preflight requests and sets appropriate headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers with specific frontend origin
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173") // Specific frontend URL
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400") // Cache preflight response for 24 hours

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No Content response for preflight
			return
		}

		c.Next()
	}
}
