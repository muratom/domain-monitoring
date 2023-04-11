package whois

import (
	"context"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/whois"
)

type Request struct {
	WhoisServer string
	Body        []byte
}

type Response struct {
	Request Request
	Raw     []byte
}

type Adapter interface {
	MakeRequest(ctx context.Context, req *Request) (*Response, error)
	PrepareRequest(ctx context.Context, fqdn string) (*Request, error)
	ParseResponse(ctx context.Context, resp *Response) (*whois.Record, error)
}
