package whois

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type Request struct {
	WhoisServer string
	Body        []byte
}

type Client interface {
	FetchWhois(context.Context, Request) ([]byte, error)
}

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

func (c *DialClient) FetchWhois(ctx context.Context, req Request) ([]byte, error) {
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

	return rawResponse, nil
}
