package cache

import (
	"encoding/json"
	"errors"
	"github.com/dgraph-io/ristretto"
	"github.com/zhouximi/email_compromise_checker/data_model"
	"log"
	"os"
)

var GlobalLocalCache *LocalCache

type LocalCache struct {
	localCache *ristretto.Cache
}

func NewLocalCache() (*LocalCache, error) {
	localCacheConfig, err := os.ReadFile("./config/local_cache_config.json")
	if err != nil {
		log.Fatalf("[InitLocalCache]Failed to read config file: %v", err)
	}

	var cfg data_model.LocalCacheConfig
	if err := json.Unmarshal(localCacheConfig, &cfg); err != nil {
		log.Fatalf("[InitLocalCache]Failed to parse config: %v", err)
	}

	localCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: cfg.NumCounters,
		MaxCost:     cfg.MaxCost,
		BufferItems: cfg.BufferItems,
	})
	if err != nil {
		log.Fatalf("[InitLocalCache]Failed to create cache: %v", err)
		return nil, err
	}
	return &LocalCache{
		localCache: localCache,
	}, nil
}

func (c *LocalCache) Get(cacheKey string) (interface{}, error) {
	if c == nil {
		log.Println("[GetFromCache]Cache is not initialized")
		return nil, errors.New("LocalCache is not initialized")
	}

	value, found := c.localCache.Get(cacheKey)
	if !found {
		log.Printf("[GetFromCache]Key %s not found in cache", cacheKey)
		return nil, errors.New("key not found in cache")
	}

	return value, nil
}

func (c *LocalCache) Set(cacheKey string, value interface{}) error {
	if c == nil {
		log.Println("[SetToCache]Cache is not initialized")
		return errors.New("LocalCache is not initialized")
	}

	success := c.localCache.Set(cacheKey, value, 1)
	if !success {
		log.Printf("[SetToCache]Failed to set key %s in cache", cacheKey)
		return errors.New("failed to set value in cache")
	}

	// Ensure the item is stored before returning
	c.localCache.Wait()

	log.Printf("[SetToCache]Successfully set key %s in cache", cacheKey)
	return nil
}
