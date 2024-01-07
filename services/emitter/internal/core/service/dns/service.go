package dns

import (
	"context"
	"fmt"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/dns"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	dnsClient Client
}

func NewService(dnsClient Client) *Service {
	return &Service{
		dnsClient: dnsClient,
	}
}

func (s *Service) LookupResourceRecords(ctx context.Context, lookupParams LookupParams) (*dns.ResourceRecords, error) {
	ctx, span := otel.Tracer("").Start(ctx, "DNSService.LookupResourceRecords", trace.WithAttributes(
		attribute.String("FQDN", lookupParams.FQDN),
		attribute.String("DNS server", lookupParams.DNSServerAddress),
	))
	defer span.End()

	span.AddEvent("dnsClient.LookupRR")
	rr, err := s.dnsClient.LookupRR(ctx, lookupParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource records for FQDN (%s): %w", lookupParams.FQDN, err)
	}

	return rr, nil
}
