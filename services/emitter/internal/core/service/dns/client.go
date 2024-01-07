package dns

import (
	"context"
	"errors"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/dns"
)

var (
	ErrStopServing = errors.New("DNS server stopped serving a domain")
)

type Client interface {
	LookupRR(ctx context.Context, params LookupParams) (*dns.ResourceRecords, error)
}

type LookupParams struct {
	FQDN string
	// IP and port of DNS server
	DNSServerAddress string
}
