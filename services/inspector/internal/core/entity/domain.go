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

type Change struct {
	// TODO: add enum
	// Type of change (add, update, delete)
	Type string
	// Path to field that has been changed
	Path []string
	// From - initial value
	From interface{}
	// To - resulting value
	To interface{}
}

type Changelog []Change

// type DomainChangelog struct {
// 	FQDN      string
// 	Timestamp time.Time
// 	Changes   Changelog
// }

type DomainRepository interface {
	GetByFQDN(ctx context.Context, fqdn string) (*Domain, error)
	GetRottenDomainsFQDN(ctx context.Context) ([]string, error)
	Store(ctx context.Context, domain *Domain) error
	SaveChangelog(ctx context.Context, fqdn string, changelog *Changelog) error
}
