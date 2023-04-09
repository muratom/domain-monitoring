package dns

import (
	"context"
	"fmt"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/dns"
)

type Service struct {
	dnsClient Client
}

func NewService(dnsClient Client) *Service {
	return &Service{
		dnsClient: dnsClient,
	}
}

func (s *Service) LookupResourceRecords(ctx context.Context, fqdn string) (*dns.ResourceRecords, error) {
	rr, err := s.dnsClient.LookupRR(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource records for FQDN (%s): %w", fqdn, err)
	}

	return rr, nil
}
