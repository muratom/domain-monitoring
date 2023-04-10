package server

import (
	"context"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/dns"
)

type dnsService interface {
	LookupResourceRecords(ctx context.Context, fqdn string) (*dns.ResourceRecords, error)
}
