package whois

import (
	"context"
	"fmt"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/whois"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := otel.Tracer("").Start(ctx, "WhoisService.FetchWhois", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	span.AddEvent("GetAdapterByFQDN")
	adapter := s.adapterProvider.GetAdapterByFQDN(fqdn)
	if adapter == nil {
		return nil, fmt.Errorf("unable to select an WHOIS adapter for FQDN (%s)", fqdn)
	}

	span.AddEvent("PrepareRequest")
	req, err := adapter.PrepareRequest(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("faield to prepare a WHOIS request for FQDN (%s): %w", fqdn, err)
	}

	span.AddEvent("MakeRequest")
	ctxMakeRequest, spanMakeRequest := otel.Tracer("").Start(ctx, "WhoisService.FetchWhois.Adapter.MakeRequest", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	resp, err := adapter.MakeRequest(ctxMakeRequest, req)
	if err != nil {
		spanMakeRequest.End()
		return nil, fmt.Errorf("faield to get a WHOIS response for FQDN (%s): %w", fqdn, err)
	}
	spanMakeRequest.End()

	span.AddEvent("ParseResponse")
	record, err := adapter.ParseResponse(ctx, resp)
	if err != nil {
		return nil, fmt.Errorf("faield to parse a WHOIS response for FQDN (%s): %w", fqdn, err)
	}

	return record, nil
}
