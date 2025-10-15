package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs important request information
// Only logs essential data to avoid over-logging
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Only log important information
		if raw != "" {
			path = path + "?" + raw
		}

		// Log format: [METHOD] PATH - STATUS_CODE - LATENCY
		// Only log for non-health endpoints and errors to avoid over-logging
		if path != "/health" && c.Writer.Status() >= 400 {
			log.Printf("[%s] %s - %d - %v",
				c.Request.Method,
				path,
				c.Writer.Status(),
				latency,
			)
		} else if path == "/health" && c.Writer.Status() != 200 {
			// Log health check failures
			log.Printf("[HEALTH] %s - %d - %v",
				path,
				c.Writer.Status(),
				latency,
			)
		}
	}
}

// ErrorLoggerMiddleware logs errors with more detail
func ErrorLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log 5xx errors with more detail
		if c.Writer.Status() >= 500 {
			log.Printf("[ERROR] %s %s - %d - Error: %v",
				c.Request.Method,
				c.Request.URL.Path,
				c.Writer.Status(),
				c.Errors.String(),
			)
		}
	}
}
