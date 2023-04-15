package entity

import (
	"context"
	"time"

	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/whois"
)

type Domain struct {
	FQDN  string
	WHOIS whois.Records
	DNS   dns.ResourceRecords
}

func (d *Domain) GetFQDN() string {
	return d.FQDN
}

func (d *Domain) GetA() []string {
	return d.DNS.A
}

func (d *Domain) GetAAAA() []string {
	return d.DNS.AAAA
}

func (d *Domain) GetCNAME() string {
	return d.DNS.CNAME
}

func (d *Domain) GetMX() []dns.MX {
	return d.DNS.MX
}

func (d *Domain) GetNS() []dns.NS {
	return d.DNS.NS
}

func (d *Domain) GetSRV() []dns.SRV {
	return d.DNS.SRV
}

func (d *Domain) GetTXT() []string {
	return d.DNS.TXT
}

func (d *Domain) GetRegistrar() string {
	return d.WHOIS.Registrar
}

func (d *Domain) GetCreatedDatetime() time.Time {
	return d.WHOIS.Created
}

func (d *Domain) GetPaidTill() time.Time {
	return d.WHOIS.PaidTill
}

type FieldType string

const (
	FQDN  FieldType = "fqdn"
	WHOIS FieldType = "whois"
	DNS   FieldType = "dns"
)

type OperationType string

const (
	CREATE OperationType = "create"
	UPDATE OperationType = "update"
	DELETE OperationType = "delete"
)

type Change struct {
	// TODO: add enum
	// Type of change (add, update, delete)
	OperationType OperationType
	// Type of the field that has been changed
	FieldType FieldType
	// Path to the field that has been changed
	Path []string
	// Initial value
	From interface{}
	// Resulting value
	To interface{}
}

type Changelog []Change

type DomainRepository interface {
	GetByFQDN(ctx context.Context, fqdn string) (*Domain, error)
	GetRottenDomainsFQDN(ctx context.Context) ([]string, error)
	Store(ctx context.Context, domain *Domain) error
	SaveChangelog(ctx context.Context, fqdn string, changelog *Changelog) error
}
