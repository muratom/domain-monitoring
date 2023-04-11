package whois

import (
	"context"
)

type Client interface {
	FetchWhois(context.Context, *Request) (*Response, error)
}
