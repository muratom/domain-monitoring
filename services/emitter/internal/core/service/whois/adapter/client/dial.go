package client

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois"
)

type DialClient struct {
	dialer net.Dialer
}

func NewWhoisClient(timeout time.Duration) *DialClient {
	return &DialClient{
		dialer: net.Dialer{
			Timeout: timeout,
		},
	}
}

func (c *DialClient) FetchWhois(ctx context.Context, req *whois.Request) (*whois.Response, error) {
	address := fmt.Sprintf("%s:43", req.WhoisServer)
	conn, err := c.dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write(req.Body); err != nil {
		return nil, fmt.Errorf("failed to send a request body: %w", err)
	}

	rawResponse, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read a response from a WHOIS server: %w", err)
	}

	return &whois.Response{
		Request: *req,
		Raw:     rawResponse,
	}, nil
}
