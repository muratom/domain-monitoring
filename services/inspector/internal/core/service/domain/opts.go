package domain

import (
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service"
)

type Option func(*Service)

func WithDomainDiffer(d service.DomainDiffer) Option {
	return func(service *Service) {
		service.domainDiffer = d
	}
}

func WithDomainCache(c service.Cache[string, domain.Domain]) Option {
	return func(service *Service) {
		service.domainCache = c
	}
}

func WithDNSCache(c service.Cache[service.GetDNSRequest, service.GetDNSResponse]) Option {
	return func(service *Service) {
		service.dnsCache = c
	}
}

func WithWHOISCache(c service.Cache[service.GetWhoisRequest, service.GetWhoisResponse]) Option {
	return func(service *Service) {
		service.whoisCache = c
	}
}
