package cache

import (
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

type ICache interface {
	Set(cacheKey string, value interface{}) error
	Get(cacheKey string) (interface{}, error)
}

type MultiLayerCache struct {
	localCache  ICache
	remoteCache ICache
}

func NewMultiLayerCache(localCache ICache, remoteCache ICache) *MultiLayerCache {
	return &MultiLayerCache{
		localCache:  localCache,
		remoteCache: remoteCache,
	}
}

func (c *MultiLayerCache) Set(cacheKey string, value interface{}) error {
	if err := c.localCache.Set(cacheKey, value); err != nil {
		log.Printf("[MultiLayerCache.Set] Local cache set failed: %v", err)
	}

	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("[MultiLayerCache.Set] Failed to marshal value: %v", err)
		return err
	}
	err = c.remoteCache.Set(cacheKey, data)
	if err != nil {
		log.Printf("[MultiLayerCache.Set] Remote cache set failed: %v", err)
		return err
	}

	return nil
}

func (c *MultiLayerCache) Get(cacheKey string) (interface{}, error) {
	// Try local cache first
	if value, err := c.localCache.Get(cacheKey); err == nil {
		log.Printf("[MultiLayerCache.Get] Found key %s in local cache", cacheKey)
		return value, nil
	}

	// Fallback to remote cache
	data, err := c.remoteCache.Get(cacheKey)
	if err == redis.Nil {
		log.Printf("[MultiLayerCache.Get] Key %s not found in remote cache", cacheKey)
		return nil, err
	} else if err != nil {
		log.Printf("[MultiLayerCache.Get] Remote cache get failed: %v", err)
		return nil, err
	}

	// Store back into local cache,
	c.localCache.Set(cacheKey, data)

	return data, nil
}
