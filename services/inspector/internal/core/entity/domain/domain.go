package domain

import (
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain/whois"
)

type Domain struct {
	FQDN  string
	WHOIS whois.Records
	DNS   dns.ResourceRecords
}
