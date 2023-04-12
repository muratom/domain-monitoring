package entity

import (
	"context"

	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/whois"
)

type Domain struct {
	FQDN  string
	WHOIS whois.Record
	DNS   dns.ResourceRecords
}

type DomainRepository interface {
	GetByFQDN(ctx context.Context, fqdn string) (*Domain, error)
	Store(ctx context.Context, domain Domain) error
}
