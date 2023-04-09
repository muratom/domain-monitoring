package whois

import (
	"context"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/whois"
)

type Adapter interface {
	MakeRequest(ctx context.Context, req *whois.Request) (*whois.Response, error)
	PrepareRequest(ctx context.Context, fqdn string) (*whois.Request, error)
	ParseResponse(ctx context.Context, resp *whois.Response) (*whois.Record, error)
}
