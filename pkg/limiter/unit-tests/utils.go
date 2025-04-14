package limiter

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func setupTestRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1,
	})
	client.FlushDB(context.Background()) // ensure Redis database is empty
	return client
}
