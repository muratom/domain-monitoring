package service

import (
	"context"

	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/whois"
)

type GetDNSRequest struct {
	FQDN string
}

type GetDNSResponse struct {
	dns.ResourceRecords
}

type GetWhoisRequest struct {
	FQDN string
}

type GetWhoisResponse struct {
	whois.Record
}

type EmitterClient interface {
	GetDNS(ctx context.Context, req *GetDNSRequest) (*GetDNSResponse, error)
	GetWhois(ctx context.Context, req *GetWhoisRequest) (*GetWhoisResponse, error)
}