package whois

import (
	"context"
	"fmt"
)

type Response struct {
	Raw []byte
}

type Record struct {
}

type Adapter interface {
	MakeRequest(ctx context.Context, req Request) (*Response, error)
	PrepareRequest(ctx context.Context, fqdn string) (*Request, error)
	ParseResponse(ctx context.Context, resp Response) (*Record, error)
}

type DefaultAdapter struct {
	whoisClient    Client
	whoisProvaider ServerProvider
}

func NewDefaultAdapter(whoisClient Client, whoisProvider ServerProvider) *DefaultAdapter {
	return &DefaultAdapter{
		whoisClient:    whoisClient,
		whoisProvaider: whoisProvider,
	}
}

func (a *DefaultAdapter) PrepareRequest(ctx context.Context, fqdn string) (*Request, error) {
	body := []byte(fmt.Sprintf("%s\r\n", fqdn))
	whoisServer, err := a.whoisProvaider.GetServerByFQDN(fqdn)
	if err != nil {
		return nil, fmt.Errorf("failed to find a corresponding WHOIS server: %w", err)
	}

	return &Request{
		WhoisServer: whoisServer,
		Body:        body,
	}, nil
}

func (a *DefaultAdapter) MakeRequest(ctx context.Context, req Request) (*Response, error) {
	resp, err := a.whoisClient.FetchWhois(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch a response from WHOIS server: %w", err)
	}

	return &Response{
		Raw: resp,
	}, nil
}

func (a *DefaultAdapter) ParseResponse(context.Context, Response) (*Record, error) {
	// We don't know the response format of an abstract WHOIS server
	panic("not implemented")
}
