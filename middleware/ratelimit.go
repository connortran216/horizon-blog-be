package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	limiters map[string]*Visitor
	mu       sync.RWMutex
	rate     time.Duration
	burst    int
}

// Visitor tracks rate limiting for individual IP addresses
type Visitor struct {
	lastSeen time.Time
	requests int
}

// NewRateLimiter creates a new rate limiter
// rate: minimum time between requests
// burst: maximum burst requests allowed
func NewRateLimiter(rate time.Duration, burst int) *RateLimiter {
	limiter := &RateLimiter{
		limiters: make(map[string]*Visitor),
		rate:     rate,
		burst:    burst,
	}

	// Clean up old entries every minute
	go limiter.cleanup()
	return limiter
}

// cleanup removes old entries to prevent memory leaks
func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		for ip, visitor := range rl.limiters {
			if time.Since(visitor.lastSeen) > 5*time.Minute {
				delete(rl.limiters, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// isAllowed checks if a request is allowed for the given IP
func (rl *RateLimiter) isAllowed(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	visitor, exists := rl.limiters[ip]
	if !exists {
		// First request from this IP
		rl.limiters[ip] = &Visitor{
			lastSeen: time.Now(),
			requests: 1,
		}
		return true
	}

	// Check if enough time has passed since last request
	if time.Since(visitor.lastSeen) < rl.rate {
		// Too soon, check if under burst limit
		if visitor.requests >= rl.burst {
			return false
		}
		visitor.requests++
	} else {
		// Reset counter if enough time has passed
		visitor.requests = 1
	}

	visitor.lastSeen = time.Now()
	return true
}

// RateLimitMiddleware creates a rate limiting middleware
// rate: minimum time between requests (e.g., 100 * time.Millisecond)
// burst: maximum burst requests allowed (e.g., 10)
func RateLimitMiddleware(rate time.Duration, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.isAllowed(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many requests",
				"message": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
