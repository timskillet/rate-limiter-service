package limiter

import (
	"testing"
	"time"

	"github.com/timskillet/rate-limiter-service/pkg/limiter"
)

func TestAllowInitialTokens(t *testing.T) {
	redisClient := setupTestRedis()
	lim := limiter.NewTokenBucketLimiter(redisClient, 5, 1)

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

func TestRefill(t *testing.T) {
	redisClient := setupTestRedis()
	lim := limiter.NewTokenBucketLimiter(redisClient, 2, 1)

	clientID := "refill-client"
	lim.Allow(clientID)
	lim.Allow(clientID) // exhausted tokens

	time.Sleep(2 * time.Second)

	if !lim.Allow(clientID) {
		t.Error("Expected request to be allowed after refill")
	}
}
