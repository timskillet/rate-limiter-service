package limiter

import (
	"testing"
	"time"

	"github.com/timskillet/rate-limiter-service/pkg/limiter"
)

func TestAllowInitialRequests(t *testing.T) {
	redisClient := setupTestRedis()
	lim := limiter.NewFixedWindowLimiter(redisClient, 5, 10*time.Second)

	clientID := "test-client"
	for i := 0; i < 5; i++ {
		allowed := lim.Allow(clientID)
		if !allowed {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}

	if lim.Allow(clientID) {
		t.Error("Expected 6th request to be rate limited")
	}
}

func TestWindowReset(t *testing.T) {
	redisClient := setupTestRedis()
	lim := limiter.NewFixedWindowLimiter(redisClient, 5, 5*time.Second)

	clientID := "test-client"
	for i := 0; i < 5; i++ {
		allowed := lim.Allow(clientID)
		if !allowed {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}

	if lim.Allow(clientID) {
		t.Error("Expected 6th request to be rate limited")
	}

	time.Sleep(5 * time.Second)

	if !lim.Allow(clientID) {
		t.Error("Expected request to be allowed after window reset")
	}
}

func TestMultipleClients(t *testing.T) {
	redisClient := setupTestRedis()
	client1 := "client1"
	client2 := "client2"

	lim := limiter.NewFixedWindowLimiter(redisClient, 2, 5*time.Second)

	if !lim.Allow(client1) || !lim.Allow(client2) {
		t.Error("Both clients should be allowed separately")
	}
}
