package cache

import (
	"github.com/redis/go-redis/v9"
	"log"
)

func NewRemoteCache() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // replace with your Redis address
		DB:   0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}
