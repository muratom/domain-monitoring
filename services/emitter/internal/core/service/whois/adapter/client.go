package adapter

import (
	"context"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois"
)

type Client interface {
	FetchWhois(context.Context, *whois.Request) (*whois.Response, error)
}
