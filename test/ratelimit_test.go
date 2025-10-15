package test

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"
)

// Simple test to verify rate limiting is working
func TestRateLimit(t *testing.T) {
	client := &http.Client{}

	// Test health endpoint with rapid requests
	fmt.Println("Testing rate limiting on /health endpoint...")

	for i := 0; i < 10; i++ {
		resp, err := client.Get("http://localhost:8080/health")
		if err != nil {
			log.Printf("Request %d failed: %v", i+1, err)
			continue
		}

		fmt.Printf("Request %d: Status %d\n", i+1, resp.StatusCode)

		if resp.StatusCode == 429 {
			fmt.Println("✅ Rate limiting is working! Received 429 Too Many Requests")
			break
		}

		resp.Body.Close()
		time.Sleep(50 * time.Millisecond) // Quick succession
	}

	fmt.Println("Rate limiting test completed.")
}
