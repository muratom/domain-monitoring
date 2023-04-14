package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/whois"
	"github.com/sirupsen/logrus"
)

const (
	workerPoolSize       = 5
	emitterClientTimeout = 3 * time.Second
)

type DomainService struct {
	emitters         []EmitterClient
	emitterCounter   atomic.Uint32
	domainRepository entity.DomainRepository
	domainDiffer     domainDiffer
}

func NewDomainService(emitterClients []EmitterClient, domainRepo entity.DomainRepository) *DomainService {
	differ := &libDomainDiffer{}
	return &DomainService{
		emitters:         emitterClients,
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

type dnsResult struct {
	response *GetDNSResponse
	request  *GetDNSRequest
	err      error
}

var (
	errNotSync = errors.New("DNS servers for domain don't synchronized")
)

func (s *DomainService) CheckDNSServers(ctx context.Context, fqdn string) error {
	domain, err := s.domainRepository.GetByFQDN(ctx, fqdn)
	if err != nil {
		return fmt.Errorf("error getting domain by FQDN: %w", err)
	}

	nameServers := domain.DNS.NS
	requests := make([]GetDNSRequest, len(nameServers))
	for i := range requests {
		requests[i] = GetDNSRequest{
			FQDN:          fqdn,
			DNSServerHost: nameServers[i].Host,
		}
	}

	wp := workerpool.New(workerPoolSize)
	results := make(chan dnsResult, len(nameServers))

	for _, req := range requests {
		req := req
		// For load-balancing get next emitter to make request
		emitter := s.getEmitterClient(ctx)
		wp.Submit(func() {
			ctx, cancel := context.WithTimeout(ctx, emitterClientTimeout)
			defer cancel()

			logrus.Infof("worker: starting DNS request to NS %v for FQDN %v", req.DNSServerHost, req.FQDN)
			resp, err := emitter.GetDNS(ctx, &req)
			results <- dnsResult{
				response: resp,
				request:  &req,
				err:      err,
			}
			logrus.Infof("worker: finished DNS request to NS %v for FQDN %v", req.DNSServerHost, req.FQDN)
		})
	}

	wp.StopWait()
	close(results)
	responses := make([]GetDNSResponse, 0, len(nameServers))
	for res := range results {
		if res.err != nil {
			return fmt.Errorf("error making DNS request to NS %v for FQDN %v: %w", res.request.DNSServerHost, res.request.DNSServerHost, err)
		}
		responses = append(responses, *res.response)
	}

	if len(responses) == 0 {
		return fmt.Errorf("no responses from emitters")
	}

	ok := s.isDNSServersSync(responses)
	if !ok {
		logrus.Errorf("DomainService.CheckDNSService. DNS servers is not synchronized for FQND %v", fqdn)
		return errNotSync
	}

	return nil
}

func (s *DomainService) DeleteDomain(ctx context.Context, deleteParams any) {}

func (s *DomainService) getUpdatedDomain(ctx context.Context, fqdn string) (*entity.Domain, error) {
	emitter := s.getEmitterClient(ctx)
	dnsRecords, err := emitter.GetDNS(ctx, &GetDNSRequest{FQDN: fqdn})
	if err != nil {
		return nil, fmt.Errorf("failed to get updated DNS records: %w", err)
	}

	whoisRecords, err := emitter.GetWhois(ctx, &GetWhoisRequest{FQDN: fqdn})
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

func (s *DomainService) getEmitterClient(ctx context.Context) EmitterClient {
	index := s.emitterCounter.Add(1) % uint32(len(s.emitters))
	return s.emitters[index]
}

func (s *DomainService) isDNSServersSync(responses []GetDNSResponse) bool {
	if len(responses) == 1 {
		return true
	}
	baseResponse := responses[0]
	for _, resp := range responses[1:] {
		if !reflect.DeepEqual(resp, baseResponse) {
			return false
		}
	}
	return true
}
