package whois

import (
	"context"
	"fmt"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/whois"
)

type Service struct {
	adapterProvider AdapterProvider
}

func NewService(adapterProvider AdapterProvider) *Service {
	return &Service{
		adapterProvider: adapterProvider,
	}
}

func (s *Service) FetchWhois(ctx context.Context, fqdn string) (*whois.Record, error) {
	adapter := s.adapterProvider.GetAdapterByFQDN(fqdn)
	if adapter == nil {
		return nil, fmt.Errorf("unable to select an WHOIS adapter for FQDN (%s)", fqdn)
	}

	req, err := adapter.PrepareRequest(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("faield to prepare a WHOIS request for FQDN (%s): %w", fqdn, err)
	}

	resp, err := adapter.MakeRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("faield to get a WHOIS response for FQDN (%s): %w", fqdn, err)
	}

	record, err := adapter.ParseResponse(ctx, resp)
	if err != nil {
		return nil, fmt.Errorf("faield to parse a WHOIS response for FQDN (%s): %w", fqdn, err)
	}

	return record, nil
}
