package service

import (
	"context"
	"fmt"

	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/whois"
)

type DomainService struct {
	emitter EmitterClient

	domainRepository entity.DomainRepository
}

func (s *DomainService) AddDomain(ctx context.Context, fqdn string) error {
	dnsRecords, err := s.emitter.GetDNS(ctx, &GetDNSRequest{FQDN: fqdn})
	if err != nil {
		return fmt.Errorf("domain service fails to get DNS records: %w", err)
	}

	whoisRecords, err := s.emitter.GetWhois(ctx, &GetWhoisRequest{FQDN: fqdn})
	if err != nil {
		return fmt.Errorf("domain service fails to get Whois records: %w", err)
	}

	domain := entity.Domain{
		FQDN: fqdn,
		WHOIS: whois.Record{
			DomainName:  whoisRecords.DomainName,
			NameServers: whoisRecords.NameServers,
			Created:     whoisRecords.Created,
			PaidTill:    whoisRecords.PaidTill,
		},
		DNS: dns.ResourceRecords{
			A:     dnsRecords.A,
			AAAA:  dnsRecords.AAAA,
			CNAME: dnsRecords.CNAME,
			MX:    dnsRecords.MX,
			NS:    dnsRecords.NS,
			SRV:   dnsRecords.SRV,
			TXT:   dnsRecords.TXT,
		},
	}
	err = s.domainRepository.Store(ctx, domain)
	if err != nil {
		return fmt.Errorf("domain service fails to store domain in the repository: %w", err)
	}
	return nil
}

func (s *DomainService) UpdateDomain(ctx context.Context, updateParams any) {}
func (s *DomainService) DeleteDomain(ctx context.Context, deleteParams any) {}
