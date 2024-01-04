package service

import (
	"time"
)

type Cache[K comparable, V any] interface {
	Runnable
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
