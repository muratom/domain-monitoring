package service

import (
	"context"
	"errors"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain/whois"
)

var (
	ErrStopServing = errors.New("DNS server stopped serving a domain")
)

type EmitterClient interface {
	GetDNS(ctx context.Context, req *GetDNSRequest) (*GetDNSResponse, error)
	GetWhois(ctx context.Context, req *GetWhoisRequest) (*GetWhoisResponse, error)
}

type GetDNSRequest struct {
	FQDN          string
	DNSServerHost string
}

type GetDNSResponse struct {
	Request         GetDNSRequest
	ResourceRecords dns.ResourceRecords
}

type GetWhoisRequest struct {
	FQDN string
}

type GetWhoisResponse struct {
	Request GetWhoisRequest
	Records whois.Records
}
