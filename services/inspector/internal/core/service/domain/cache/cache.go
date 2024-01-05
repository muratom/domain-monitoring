package cache

import (
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service"
	"time"
)

type Cache[K comparable, V any] interface {
	service.Runnable
	Get(key K) *V
	Set(key K, value V, opts ...CacheOption)
}

type CacheOption func(*CacheConfig)

type CacheConfig struct {
	TTL *time.Duration
}

func WithTTL(ttl time.Duration) CacheOption {
	return func(c *CacheConfig) {
		c.TTL = &ttl
	}
}
