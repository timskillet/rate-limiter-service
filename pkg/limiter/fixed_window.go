package limiter

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	redisPkg "github.com/timskillet/rate-limiter-service/pkg/redis"
)

type FixedWindowLimiter struct {
	redisClient *redis.Client
	limit       int           // Max allowed requests in the current window
	windowSize  time.Duration // Duration of the fixed window
}

func NewFixedWindowLimiter(redisClient *redis.Client, limit int, windowSize time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		redisClient: redisClient,
		limit:       limit,
		windowSize:  windowSize,
	}
}

func (l *FixedWindowLimiter) Allow(clientID string) bool {
	ctx := redisPkg.GetContext()

	// Create a unique key for the client's request count
	windowKey := fmt.Sprintf("fw:%s:%d", clientID, time.Now().Unix()/int64(l.windowSize.Seconds()))

	// Increment the count for this window
	count, err := l.redisClient.Incr(ctx, windowKey).Result()
	if err != nil {
		return false
	}

	// Set expiration for key if it's new
	if count == 1 {
		l.redisClient.Expire(ctx, windowKey, l.windowSize)
	}

	return count <= int64(l.limit)
}
