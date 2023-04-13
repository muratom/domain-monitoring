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
	domainDiffer     domainDiffer
}

func NewDomainService(emitterClient EmitterClient, domainRepo entity.DomainRepository) *DomainService {
	differ := &libDomainDiffer{}
	return &DomainService{
		emitter:          emitterClient,
		domainRepository: domainRepo,
		domainDiffer:     differ,
	}
}

func (s *DomainService) AddDomain(ctx context.Context, fqdn string) error {
	domain, err := s.getUpdatedDomain(ctx, fqdn)
	if err != nil {
		return fmt.Errorf("error when getting updated domain data for FQDN (%v): %w", fqdn, err)
	}

	err = s.domainRepository.Store(ctx, domain)
	if err != nil {
		return fmt.Errorf("domain service fails to store domain in the repository: %w", err)
	}
	return nil
}

func (s *DomainService) GetRottenDomainsFQDN(ctx context.Context) ([]string, error) {
	return s.domainRepository.GetRottenDomainsFQDN(ctx)
}

func (s *DomainService) UpdateDomain(ctx context.Context, fqdn string) (entity.Changelog, error) {
	rottenDomain, err := s.domainRepository.GetByFQDN(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch domain (%v) data from DB: %w", fqdn, err)
	}

	freshDomain, err := s.getUpdatedDomain(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("error when getting updated domain data for FQDN (%v): %w", fqdn, err)
	}

	changelog, err := s.domainDiffer.Diff(freshDomain, rottenDomain)
	if err != nil {
		return nil, fmt.Errorf("error when making diff between domains: %w", err)
	}
	return changelog, nil
}

func (s *DomainService) getUpdatedDomain(ctx context.Context, fqdn string) (*entity.Domain, error) {
	dnsRecords, err := s.emitter.GetDNS(ctx, &GetDNSRequest{FQDN: fqdn})
	if err != nil {
		return nil, fmt.Errorf("failed to get updated DNS records: %w", err)
	}

	whoisRecords, err := s.emitter.GetWhois(ctx, &GetWhoisRequest{FQDN: fqdn})
	if err != nil {
		return nil, fmt.Errorf("failed to get updated Whois records: %w", err)
	}

	domain := &entity.Domain{
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
	return domain, nil
}

func (s *DomainService) DeleteDomain(ctx context.Context, deleteParams any) {}
