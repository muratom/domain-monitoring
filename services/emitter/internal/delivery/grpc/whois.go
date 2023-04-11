package grpc

import (
	"context"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/whois"
)

type whoisService interface {
	FetchWhois(ctx context.Context, fqdn string) (*whois.Record, error)
}
