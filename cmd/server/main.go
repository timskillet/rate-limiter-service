package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/timskillet/rate-limiter-service/pkg/api"
	"github.com/timskillet/rate-limiter-service/pkg/limiter"
	"github.com/timskillet/rate-limiter-service/pkg/metrics"
	"github.com/timskillet/rate-limiter-service/pkg/redis"
)

// Define rate limiter
var rateLimiter api.LimiterInterface

func rateLimiterHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client")

	if rateLimiter.Allow(clientID) {
		metrics.RecordRateLimitHit("success")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Request allowed\n")
	} else {
		metrics.RecordRateLimitHit("failure\n")
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(w, "Rate limit exceeded\n")
	}
}

func main() {
	// Setup Prometheus metrics
	metrics.SetupMetrics()

	// Create Redis client
	redisClient := redis.NewRedisClient()

	// Initialize Redis-backed token bucket limiter (e.g., 10 tokens, refill 1 per second)
	// rateLimiter := limiter.NewTokenBucketLimiter(redisClient, 10, 1)

	// Initialize fixed window limiter: allow 5 requests per 10 seconds
	rateLimiter = limiter.NewFixedWindowLimiter(redisClient, 5, 10*time.Second)

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		clientID := r.URL.Query().Get("client")

		if rateLimiter.Allow(clientID) {
			metrics.RecordRateLimitHit("success")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Request allowed")
		} else {
			metrics.RecordRateLimitHit("failure")
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintf(w, "Rate limit exceeded: %s", clientID)
		}
	})

	log.Println("Starting rate limiter service on :8080")
	http.ListenAndServe(":8080", nil)
}
