package dns

import (
	"context"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/dns"
)

type LookupParams struct {
	FQDN          string
	DNSServerHost string
}

type Client interface {
	LookupRR(ctx context.Context, params LookupParams) (*dns.ResourceRecords, error)
}
