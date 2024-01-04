package service

import "github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain"

type Option func(*DomainService)

func WithDomainDiffer(d DomainDiffer) Option {
	return func(service *DomainService) {
		service.domainDiffer = d
	}
}

func WithDomainCache(c Cache[string, domain.Domain]) Option {
	return func(service *DomainService) {
		service.domainCache = c
	}
}

func WithDNSCache(c Cache[GetDNSRequest, GetDNSResponse]) Option {
	return func(service *DomainService) {
		service.dnsCache = c
	}
}

func WithWHOISCache(c Cache[GetWhoisRequest, GetWhoisResponse]) Option {
	return func(service *DomainService) {
		service.whoisCache = c
	}
}
