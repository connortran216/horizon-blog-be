package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles CORS preflight requests and sets appropriate headers
func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := []string{
		"http://localhost:5173",             // Development frontend
		"https://blog.connortran.io.vn",     // Production frontend
		"https://blog-api.connortran.io.vn", // API domain (in case of direct calls)
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if the origin is in our allowed list
		allowOrigin := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowOrigin = true
				break
			}
		}

		if allowOrigin {
			c.Header("Access-Control-Allow-Origin", origin)
		}

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
