package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/redis/go-redis/v9"
	"time"
)

type ICache interface {
	Set(cacheKey string, value interface{}) error
	Get(cacheKey string) (interface{}, error)
}

type MultiLayerCache struct {
	localCache  *ristretto.Cache
	remoteCache *redis.Client
	ttl         time.Duration // TTL for Redis cache in seconds
}

func NewMultiLayerCache(localCache *ristretto.Cache, remoteCache *redis.Client) *MultiLayerCache {
	return &MultiLayerCache{
		localCache:  localCache,
		remoteCache: remoteCache,
	}
}

func (c *MultiLayerCache) Get(cacheKey string) (interface{}, error) {
	if c.localCache != nil {
		if val, ok := c.localCache.Get(cacheKey); ok {
			return val, nil
		}
	}

	val, err := c.remoteCache.Get(context.Background(), cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // not found in Redis
	} else if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}

	// 4. Set to local cache
	if c.localCache != nil {
		c.localCache.Set(cacheKey, result, 1)
		c.localCache.Wait()
	}

	return result, nil
}

func (c *MultiLayerCache) Set(cacheKey string, value interface{}) error {
	// 1. Set to local cache
	if c.localCache != nil {
		c.localCache.Set(cacheKey, value, 1)
	}

	// 2. Marshal to JSON for Redis
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	// 3. Set in Redis with TTL
	err = c.remoteCache.Set(context.Background(), cacheKey, data, c.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set value in Redis: %w", err)
	}

	return nil
}
