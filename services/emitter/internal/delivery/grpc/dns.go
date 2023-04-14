package grpc

import (
	"context"

	dnsentity "github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/dns"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/dns"
)

type dnsService interface {
	LookupResourceRecords(ctx context.Context, lookupParams dns.LookupParams) (*dnsentity.ResourceRecords, error)
}
