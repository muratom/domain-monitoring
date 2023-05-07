package service

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
)

type Item struct {
	key   string
	value entity.Domain
}

func (i *Item) Key() string {
	return i.key
}

func (i *Item) Value() entity.Domain {
	return i.value
}

type domainTTLCache interface {
	Start()
	Stop()
	Get(key string) *Item
	Set(key string, value entity.Domain, ttl time.Duration)
}

type libDomainTTLCache struct {
	cache *ttlcache.Cache[string, entity.Domain]
}

func NewLibDomainTTLCache() *libDomainTTLCache {
	return &libDomainTTLCache{
		cache: ttlcache.New(
			ttlcache.WithTTL[string, entity.Domain](cacheTTL),
		),
	}
}

func (c *libDomainTTLCache) Start() {
	go c.cache.Start()
}

func (c *libDomainTTLCache) Stop() {
	c.cache.Stop()
}

func (c *libDomainTTLCache) Get(key string) *Item {
	item := c.cache.Get(key)
	if item == nil {
		return nil
	}
	return &Item{
		key:   item.Key(),
		value: item.Value(),
	}
}

func (c *libDomainTTLCache) Set(key string, value entity.Domain, ttl time.Duration) {
	c.cache.Set(key, value, ttl)
}
