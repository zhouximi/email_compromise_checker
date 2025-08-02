package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/zhouximi/email_compromise_checker/data_model"
	"log"
	"os"
)

type RemoteCache struct {
	redisCache *redis.Client
}

const redisConfigPath = "config/redis_config.json"

func NewRemoteCache() (*RemoteCache, error) {
	data, err := os.ReadFile(redisConfigPath)
	if err != nil {
		log.Printf("[NewRemoteCache] Failed to read config file: %v", err)
		return nil, err
	}

	var cfg data_model.RedisConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Printf("[NewRemoteCache] Failed to unmarshal config: %v", err)
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Printf("[NewRemoteCache] Redis ping failed: %v", err)
		return nil, err
	}

	log.Println("[NewRemoteCache] Redis client successfully initialized")
	return &RemoteCache{redisCache: rdb}, nil
}

func (c *RemoteCache) Get(cacheKey string) (interface{}, error) {
	if c == nil || c.redisCache == nil {
		return nil, errors.New("RemoteCache is not initialized")
	}

	data, err := c.redisCache.Get(context.Background(), cacheKey).Bytes()
	if err == redis.Nil {
		log.Printf("[RemoteCache.Get] Key %s not found", cacheKey)
		return nil, errors.New("key not found in remote cache")
	} else if err != nil {
		log.Printf("[RemoteCache.Get] Failed to get key %s: %v", cacheKey, err)
		return nil, err
	}

	return data, nil
}

func (c *RemoteCache) Set(cacheKey string, value interface{}) error {
	if c == nil || c.redisCache == nil {
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
