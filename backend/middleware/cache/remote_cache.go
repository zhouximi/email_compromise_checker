package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type RemoteCache struct {
	redisCache *redis.Client
}

func NewRemoteCache() (*RemoteCache, error) {
	redisConfig := getRedisConfig()
	rdb := redis.NewClient(redisConfig)

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Printf("[NewRemoteCache] Redis ping failed: %v, redis_config=%#v", err, redisConfig)
		return nil, err
	}

	log.Println("[NewRemoteCache] Redis client successfully initialized")
	return &RemoteCache{redisCache: rdb}, nil
}

func (c *RemoteCache) Get(cacheKey string) (interface{}, error) {
	if c.redisCache == nil {
		return nil, errors.New("RemoteCache is not initialized")
	}

	data, err := c.redisCache.Get(context.Background(), cacheKey).Bytes()
	if err == redis.Nil {
		log.Printf("[RemoteCache.Get] Key %s not found", cacheKey)
		return nil, errors.New("[RemoteCache.Get] Key not found in remote cache")
	} else if err != nil {
		log.Printf("[RemoteCache.Get] Failed to get key %s: %v", cacheKey, err)
		return nil, err
	}

	return data, nil
}

func (c *RemoteCache) Set(cacheKey string, value interface{}) error {
	if c.redisCache == nil {
		return errors.New("RemoteCache is not initialized")
	}

	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("[RemoteCache.Set] Failed to marshal value: %v", err)
		return err
	}

	err = c.redisCache.Set(context.Background(), cacheKey, data, 0).Err() // 0 means no expiration
	if err != nil {
		log.Printf("[RemoteCache.Set] Failed to set key %s: %v", cacheKey, err)
		return err
	}

	return nil
}

func getRedisConfig() *redis.Options {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	if host == "" {
		host = "redis"
	}
	if port == "" {
		port = "6379"
	}

	addr := host + ":" + port
	log.Printf("Connecting Redis at %s", addr)

	return &redis.Options{
		Addr: addr,
	}
}
