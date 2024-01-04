package ttl

import (
	"github.com/jellydator/ttlcache/v3"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service"
	"time"
)

const (
	cacheTTL = 1 * time.Minute
)

type Cache[K comparable, V any] struct {
	cache *ttlcache.Cache[K, V]
}

func New[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		cache: ttlcache.New(
			ttlcache.WithTTL[K, V](cacheTTL),
		),
	}
}

func (c *Cache[K, V]) Get(key K) *V {
	item := c.cache.Get(key)
	if item == nil {
		return nil
	}
	result := item.Value()
	return &result
}

func (c *Cache[K, V]) Set(key K, value V, opts ...service.CacheOption) {
	config := new(service.CacheConfig)
	for _, opt := range opts {
		opt(config)
	}

	ttl := ttlcache.DefaultTTL
	if config.TTL != nil {
		ttl = *config.TTL
	}
	c.cache.Set(key, value, ttl)
}

func (c *Cache[K, V]) Start() {
	go c.cache.Start()
}

func (c *Cache[K, V]) Stop() {
	c.cache.Stop()
}
