package dns

import (
	"context"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/dns"
)

type Client interface {
	LookupRR(ctx context.Context, fqdn string) (*dns.ResourceRecords, error)
}
