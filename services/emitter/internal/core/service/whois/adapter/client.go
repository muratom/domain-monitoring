package adapter

import (
	"context"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/whois"
)

type Client interface {
	FetchWhois(context.Context, *whois.Request) (*whois.Response, error)
}
