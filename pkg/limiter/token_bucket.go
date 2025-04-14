package limiter

import (
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	redisPkg "github.com/timskillet/rate-limiter-service/pkg/redis"
)

type TokenBucketLimiter struct {
	redisClient *redis.Client
	maxTokens   int
	refillRate  int // tokens per second
}

func NewTokenBucketLimiter(redisClient *redis.Client, maxTokens, refillRate int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		redisClient: redisClient,
		maxTokens:   maxTokens,
		refillRate:  refillRate,
	}
}

func (l *TokenBucketLimiter) Allow(clientID string) bool {
	ctx := redisPkg.GetContext()
	key := fmt.Sprintf("token_bucket:%s", clientID)
	lastRefillKey := fmt.Sprintf("token_bucket_last_refill:%s", clientID)

	// Initialize bucket if not exists
	exists, err := l.redisClient.Exists(ctx, key).Result()
	if err != nil {
		log.Println("Redis EXISTS error:", err)
		return false
	}
	if exists == 0 {
		// New client, initialize tokens
		err := l.redisClient.Set(ctx, key, l.maxTokens, 0).Err()
		if err != nil {
			log.Println("Redis SET error:", err)
			return false
		}
		l.redisClient.Set(ctx, lastRefillKey, time.Now().Unix(), 0)
	}

	// Refill tokens if enough time has passed
	now := time.Now().Unix()
	lastRefill, _ := l.redisClient.Get(ctx, lastRefillKey).Int64()
	elapsed := now - lastRefill
	if elapsed > 0 {
		refillAmount := int(elapsed) * l.refillRate
		if refillAmount > 0 {
			currTokens, _ := l.redisClient.Get(ctx, key).Int()
			newTokens := min(l.maxTokens, currTokens+refillAmount)
			l.redisClient.Set(ctx, key, newTokens, 0)
			l.redisClient.Set(ctx, lastRefillKey, now, 0)
		}
	}

	// Try to consume a token
	tokens, err := l.redisClient.Get(ctx, key).Int()
	if err != nil {
		log.Println("Redis GET error:", err)
		return false
	}

	if tokens > 0 {
		l.redisClient.Decr(ctx, key)
		fmt.Printf("Client %s allowed. Tokens left: %d\n", clientID, tokens-1)
		return true
	}

	fmt.Printf("Client %s denied. Tokens left: %d\n", clientID, tokens)
	return false
}

// Utility function for refilling bucket
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
