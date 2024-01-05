package domain

import (
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service/domain/cache"
)

type Option func(*Service)

func WithDomainDiffer(d Differ) Option {
	return func(service *Service) {
		service.domainDiffer = d
	}
}

func WithDomainCache(c cache.Cache[string, domain.Domain]) Option {
	return func(service *Service) {
		service.domainCache = c
	}
}

func WithDNSCache(c cache.Cache[GetDNSRequest, GetDNSResponse]) Option {
	return func(service *Service) {
		service.dnsCache = c
	}
}

func WithWHOISCache(c cache.Cache[GetWhoisRequest, GetWhoisResponse]) Option {
	return func(service *Service) {
		service.whoisCache = c
	}
}
