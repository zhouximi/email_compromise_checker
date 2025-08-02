package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type ICache interface {
	Set(cacheKey string, value interface{}) error
	Get(cacheKey string) (interface{}, error)
}

type MultiLayerCache struct {
	localCache  *LocalCache
	remoteCache *RemoteCache
	ttl         time.Duration // TTL for Redis cache in seconds
}

func NewMultiLayerCache(localCache *LocalCache, remoteCache *RemoteCache) *MultiLayerCache {
	return &MultiLayerCache{
		localCache:  localCache,
		remoteCache: remoteCache,
	}
}

func (c *MultiLayerCache) Set(cacheKey string, value interface{}) error {
	if c.localCache != nil {
		if err := c.localCache.Set(cacheKey, value); err != nil {
			log.Printf("[MultiLayerCache.Set] Local cache set failed: %v", err)
		}
	}

	if c.remoteCache != nil {
		data, err := json.Marshal(value)
		if err != nil {
			log.Printf("[MultiLayerCache.Set] Failed to marshal value: %v", err)
			return err
		}
		err = c.remoteCache.redisCache.Set(context.Background(), cacheKey, data, c.ttl).Err()
		if err != nil {
			log.Printf("[MultiLayerCache.Set] Remote cache set failed: %v", err)
			return err
		}
	}

	return nil
}

func (c *MultiLayerCache) Get(cacheKey string) (interface{}, error) {
	// Try local cache first
	if c.localCache != nil {
		if value, err := c.localCache.Get(cacheKey); err == nil {
			log.Printf("[MultiLayerCache.Get] Found key %s in local cache", cacheKey)
			return value, nil
		}
	}

	// Fallback to remote cache
	if c.remoteCache != nil {
		data, err := c.remoteCache.redisCache.Get(context.Background(), cacheKey).Bytes()
		if err == redis.Nil {
			log.Printf("[MultiLayerCache.Get] Key %s not found in remote cache", cacheKey)
			return nil, errors.New("key not found in cache")
		} else if err != nil {
			log.Printf("[MultiLayerCache.Get] Remote cache get failed: %v", err)
			return nil, err
		}

		// Store back into local cache
		var result interface{}
		if err := json.Unmarshal(data, &result); err == nil && c.localCache != nil {
			_ = c.localCache.Set(cacheKey, result)
		}

		return data, nil
	}

	return nil, errors.New("no cache layer available")
}
