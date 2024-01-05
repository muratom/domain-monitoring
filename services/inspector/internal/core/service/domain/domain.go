package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/changelog"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain/whois"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/notification"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service/cache/ttl"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service/differ/lib"
	"sort"
	"sync/atomic"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slices"
)

// TODO: set by config
const (
	workerPoolSize          = 5
	emitterClientTimeout    = 3 * time.Second
	expiringDomainThreshold = 30 * 24 * time.Hour // 1 month
)

type Service struct {
	emitters       []service.EmitterClient
	emitterCounter atomic.Uint32

	domainRepository    domain.CRUDRepository
	changelogRepository changelog.Repository

	domainDiffer service.DomainDiffer

	domainCache service.Cache[string, domain.Domain]
	dnsCache    service.Cache[service.GetDNSRequest, service.GetDNSResponse]
	whoisCache  service.Cache[service.GetWhoisRequest, service.GetWhoisResponse]
}

func New(
	emitterClients []service.EmitterClient,
	domainRepo domain.CRUDRepository,
	changelogRepo changelog.Repository,
	opts ...Option,
) *Service {
	s := &Service{
		emitters:            emitterClients,
		domainRepository:    domainRepo,
		changelogRepository: changelogRepo,
		domainDiffer:        &lib.Differ{},
		domainCache:         ttl.New[string, domain.Domain](),
		dnsCache:            ttl.New[service.GetDNSRequest, service.GetDNSResponse](),
		whoisCache:          ttl.New[service.GetWhoisRequest, service.GetWhoisResponse](),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Service) Start(ctx context.Context) {
	// Enable an automatic removal of expired items
	s.domainCache.Start(ctx)
	s.dnsCache.Start(ctx)
	s.whoisCache.Start(ctx)
}

func (s *Service) Stop(ctx context.Context) {
	s.domainCache.Stop(ctx)
	s.dnsCache.Stop(ctx)
	s.whoisCache.Stop(ctx)
}

func (s *Service) AddDomain(ctx context.Context, fqdn string) (*domain.Domain, error) {
	ctx, span := otel.Tracer("").Start(ctx, "Service.AddDomain", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	updatedDomain, err := s.getUpdatedDomain(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("Service.AddDomain: error getting updated domain data for FQDN (%v): %w", fqdn, err)
	}

	err = s.domainRepository.Store(ctx, updatedDomain)
	if err != nil {
		return nil, fmt.Errorf("Service.AddDomain: failed to store domain in the repository: %w", err)
	}

	return updatedDomain, nil
}

func (s *Service) GetDomain(ctx context.Context, fqdn string) (*domain.Domain, error) {
	ctx, span := otel.Tracer("").Start(ctx, "Service.GetDomain", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	retrievedDomain, err := s.domainRepository.ByFQDN(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("Service.GetDomain: failed to get domain from repository: %w", err)
	}

	return retrievedDomain, err
}

func (s *Service) GetAllDomainsFQDN(ctx context.Context) ([]string, error) {
	ctx, span := otel.Tracer("").Start(ctx, "Service.RetrieveAllDomainsFQDN")
	defer span.End()

	return s.domainRepository.AllDomainsFQDN(ctx)
}

func (s *Service) UpdateDomain(ctx context.Context, fqdn string) (*domain.Domain, error) {
	ctx, span := otel.Tracer("").Start(ctx, "Service.UpdateDomain", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	updatedDomain, err := s.getUpdatedDomain(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("Service.UpdateDomain: error getting updated domain data for FQDN (%v): %w", fqdn, err)
	}

	err = s.domainRepository.Update(ctx, updatedDomain, fqdn)
	if err != nil {
		return nil, fmt.Errorf("Service.UpdateDomain: failed to store domain in the repository: %w", err)
	}
	return updatedDomain, nil
}

func (s *Service) DeleteDomain(ctx context.Context, fqdn string) error {
	ctx, span := otel.Tracer("").Start(ctx, "Service.DeleteDomain", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	err := s.domainRepository.Delete(ctx, fqdn)
	if err != nil {
		return fmt.Errorf("Service.DeleteDomain: failed to delete domain: %w", err)
	}

	return nil
}

func (s *Service) RetrieveRottenDomainsFQDN(ctx context.Context) ([]string, error) {
	ctx, span := otel.Tracer("").Start(ctx, "Service.RetrieveRottenDomainsFQDN")
	defer span.End()

	return s.domainRepository.RottenDomainsFQDN(ctx)
}

func (s *Service) GetChangelogs(ctx context.Context, fqdn string) ([]changelog.Changelog, error) {
	ctx, span := otel.Tracer("").Start(ctx, "Service.Retrieve")
	defer span.End()

	return s.changelogRepository.Retrieve(ctx, fqdn)
}

func (s *Service) CheckDomainNameServers(ctx context.Context, fqdn string) ([]notification.Notification, error) {
	ctx, span := otel.Tracer("").Start(ctx, "Service.CheckDomainNameServers", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	retrievedDomain, err := s.domainRepository.ByFQDN(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("error getting domain by FQDN: %w", err)
	}

	nameServers := retrievedDomain.DNS.NS
	requests := make([]service.GetDNSRequest, len(nameServers))
	for i := range requests {
		requests[i] = service.GetDNSRequest{
			FQDN:          fqdn,
			DNSServerHost: nameServers[i].Host,
		}
	}

	wp := workerpool.New(workerPoolSize)
	results := make(chan dnsResult, len(nameServers))

	for _, req := range requests {
		req := req
		wp.Submit(func() {
			ctx, cancel := context.WithCancel(ctx)

			ctx, span := otel.Tracer("").Start(ctx, "Service.CheckDomainNameServers.worker", trace.WithAttributes(
				attribute.String("FQDN", req.FQDN),
				attribute.String("DNS server", req.DNSServerHost),
			))
			defer span.End()
			defer cancel()

			logrus.Infof("worker: starting DNS request to NS %v for FQDN %v", req.DNSServerHost, req.FQDN)
			resp, err := s.getDNS(ctx, req)
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

	notifications := make([]notification.Notification, 0)
	responses := make([]service.GetDNSResponse, 0, len(nameServers))
	for res := range results {
		if res.err != nil {
			if errors.Is(res.err, service.ErrStopServing) {
				notifications = append(notifications, &notification.DomainStoppedBeingServedNotification{
					FQDN:           res.request.FQDN,
					NameServerHost: res.request.DNSServerHost,
				})
				continue
			} else {
				return nil, fmt.Errorf("error making DNS request to NS %v for FQDN %v: %w", res.request.DNSServerHost, res.request.FQDN, res.err)
			}
		}
		responses = append(responses, *res.response)
	}

	if len(responses) == 0 {
		return nil, fmt.Errorf("no responses w/o error from emitters")
	}

	ok, notSychronizedDNSServers := s.isDNSServersSync(responses)
	if !ok {
		logrus.Warnf("Service.CheckDNSService. DNS servers is not synchronized for FQND %v", fqdn)
		notifications = append(notifications, &notification.NameServersNotSynchronizedNotification{
			FQDN:                       fqdn,
			NotSynchronizedNameServers: notSychronizedDNSServers,
		})
	}

	return notifications, nil
}

func (s *Service) CheckDomainRegistration(ctx context.Context, fqdn string) ([]notification.Notification, error) {
	ctx, span := otel.Tracer("").Start(ctx, "Service.CheckDomainRegistration", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	emitter := s.getEmitterClient(ctx)
	whoisResp, err := emitter.GetWhois(ctx, &service.GetWhoisRequest{FQDN: fqdn})
	if err != nil {
		return nil, fmt.Errorf("error getting registration information: %w", err)
	}

	notifications := make([]notification.Notification, 0, 2)

	expiringSoonTimestamp := whoisResp.Records.PaidTill.Add(-expiringDomainThreshold)
	if time.Now().After(expiringSoonTimestamp) {
		// Domain registration is going to expire
		n := &notification.RegistrationExpiresSoonNotification{
			FQDN:      fqdn,
			Registrar: whoisResp.Records.Registrar,
			PaidTill:  whoisResp.Records.PaidTill,
		}
		notifications = append(notifications, n)
	} else if time.Now().After(whoisResp.Records.PaidTill) {
		// Domain has been expired
		n := &notification.RegistrationExpiredNotification{
			FQDN:      fqdn,
			Registrar: whoisResp.Records.Registrar,
			PaidTill:  whoisResp.Records.PaidTill,
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (s *Service) CheckDomainChanges(ctx context.Context, fqdn string) ([]notification.Notification, error) {
	ctx, span := otel.Tracer("").Start(ctx, "Service.CheckDomainChanges", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	changes, err := s.getDomainChanges(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("error getting domain changes: %w", err)
	}

	if len(changes) != 0 {
		err := s.changelogRepository.Store(ctx, fqdn, changes)
		if err != nil {
			return nil, fmt.Errorf("failed to save changelog: %w", err)
		}
	}

	notifications := make([]notification.Notification, 0)
	for _, change := range changes {
		switch change.FieldType {
		case changelog.FQDN:
			n := &notification.DomainNameChangedNotification{
				Old: change.From.(string),
				New: change.To.(string),
			}
			notifications = append(notifications, n)
		case changelog.DNS:
			recordType, path := change.Path[0], change.Path[1:]
			n := &notification.ResourceRecordChangedNotification{
				FQDN:       fqdn,
				RecordType: recordType,
				Path:       path,
				From:       change.From,
				To:         change.To,
			}
			notifications = append(notifications, n)
		case changelog.WHOIS:
			n := &notification.RegistrationInfoChangedNotification{
				FQDN: fqdn,
				Path: change.Path,
				From: change.From,
				To:   change.To,
			}
			notifications = append(notifications, n)
		}
	}

	return notifications, nil
}

func (s *Service) getUpdatedDomain(ctx context.Context, fqdn string) (*domain.Domain, error) {
	ctx, span := otel.Tracer("").Start(ctx, "Service.getUpdatedDomain", trace.WithAttributes(
		attribute.Bool("from_cache", false),
	))
	defer span.End()

	if cachedDomain := s.domainCache.Get(fqdn); cachedDomain != nil {
		span.SetAttributes(attribute.Bool("from_cache", true))
		return cachedDomain, nil
	}

	dnsChan := make(chan dnsResponse)
	go func() {
		dnsRecords, err := s.getDNS(ctx, service.GetDNSRequest{FQDN: fqdn})
		dnsChan <- dnsResponse{
			response: dnsRecords,
			err:      err,
		}
	}()

	whoisChan := make(chan whoisResponse)
	go func() {
		whoisRecords, err := s.getWhois(ctx, service.GetWhoisRequest{FQDN: fqdn})
		whoisChan <- whoisResponse{
			response: whoisRecords,
			err:      err,
		}
	}()

	dnsResp := <-dnsChan
	if dnsResp.err != nil {
		return nil, fmt.Errorf("failed to get updated DNS records: %v", dnsResp.err)
	}
	dnsRecords := dnsResp.response

	whoisResp := <-whoisChan
	if whoisResp.err != nil {
		return nil, fmt.Errorf("failed to get updated Whois records: %v", whoisResp.err)
	}
	whoisRecords := whoisResp.response

	updatedDomain := &domain.Domain{
		FQDN: fqdn,
		WHOIS: whois.Records{
			DomainName:  whoisRecords.Records.DomainName,
			NameServers: whoisRecords.Records.NameServers,
			Registrar:   whoisRecords.Records.Registrar,
			Created:     whoisRecords.Records.Created,
			PaidTill:    whoisRecords.Records.PaidTill,
		},
		DNS: dns.ResourceRecords{
			A:     dnsRecords.ResourceRecords.A,
			AAAA:  dnsRecords.ResourceRecords.AAAA,
			CNAME: dnsRecords.ResourceRecords.CNAME,
			MX:    dnsRecords.ResourceRecords.MX,
			NS:    dnsRecords.ResourceRecords.NS,
			SRV:   dnsRecords.ResourceRecords.SRV,
			TXT:   dnsRecords.ResourceRecords.TXT,
		},
	}

	s.domainCache.Set(fqdn, *updatedDomain)

	return updatedDomain, nil
}

func (s *Service) getDomainChanges(ctx context.Context, fqdn string) (changelog.Changelog, error) {
	rottenDomain, err := s.domainRepository.ByFQDN(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch domain (%v) data from DB: %w", fqdn, err)
	}

	freshDomain, err := s.getUpdatedDomain(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("error getting updated domain data for FQDN (%v): %w", fqdn, err)
	}

	changes, err := s.domainDiffer.Diff(rottenDomain, freshDomain)
	if err != nil {
		return nil, fmt.Errorf("error making diff between domains: %w", err)
	}
	return changes, nil
}

func (s *Service) getEmitterClient(ctx context.Context) service.EmitterClient {
	index := s.emitterCounter.Add(1) % uint32(len(s.emitters))
	return s.emitters[index]
}

func (s *Service) getDNS(ctx context.Context, req service.GetDNSRequest) (*service.GetDNSResponse, error) {
	if resp := s.dnsCache.Get(req); resp != nil {
		logrus.Info("retrieve DNS response from cache for req (%+v)", req)
		return resp, nil
	}

	emitter := s.getEmitterClient(ctx)
	resp, err := emitter.GetDNS(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated DNS records: %w", err)
	}

	s.dnsCache.Set(req, *resp)

	return resp, nil
}

func (s *Service) getWhois(ctx context.Context, req service.GetWhoisRequest) (*service.GetWhoisResponse, error) {
	if resp := s.whoisCache.Get(req); resp != nil {
		logrus.Info("retrieve Whois response from cache for req (%+v)", req)
		return resp, nil
	}

	emitter := s.getEmitterClient(ctx)
	resp, err := emitter.GetWhois(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated Whois records: %w", err)
	}

	s.whoisCache.Set(req, *resp)

	return resp, nil
}

func (s *Service) isDNSServersSync(responses []service.GetDNSResponse) (bool, []string) {
	if len(responses) == 1 {
		return true, nil
	}
	baseResponse := responses[0]
	baseResponseResourceRecords := baseResponse.ResourceRecords
	syncWithBase := []string{baseResponse.Request.DNSServerHost}
	notSyncWithBase := []string{}
	for _, resp := range responses[1:] {
		// if !reflect.DeepEqual(resp.ResourceRecords, baseResponseResourceRecords) {
		if !compareResourceRecords(resp.ResourceRecords, baseResponseResourceRecords) {
			notSyncWithBase = append(notSyncWithBase, resp.Request.DNSServerHost)
		} else {
			syncWithBase = append(syncWithBase, resp.Request.DNSServerHost)
		}
	}
	if len(notSyncWithBase) == 0 {
		return true, nil
	}

	// Tell about a minority of servers
	if len(syncWithBase) >= len(notSyncWithBase) {
		return false, notSyncWithBase
	} else {
		return false, syncWithBase
	}
}

type dnsResult struct {
	response *service.GetDNSResponse
	request  *service.GetDNSRequest
	err      error
}

type dnsResponse struct {
	response *service.GetDNSResponse
	err      error
}

type whoisResponse struct {
	response *service.GetWhoisResponse
	err      error
}

func compareResourceRecords(x, y dns.ResourceRecords) bool {
	sort.Strings(x.A)
	sort.Strings(y.A)
	if !slices.Equal(x.A, y.A) {
		return false
	}

	sort.Strings(x.AAAA)
	sort.Strings(y.AAAA)
	if !slices.Equal(x.AAAA, y.AAAA) {
		return false
	}

	if x.CNAME != y.CNAME {
		return false
	}

	sort.Sort(x.MX)
	sort.Sort(y.MX)
	if !slices.Equal(x.MX, y.MX) {
		return false
	}

	sort.Sort(x.NS)
	sort.Sort(y.NS)
	if !slices.Equal(x.NS, y.NS) {
		return false
	}

	sort.Sort(x.SRV)
	sort.Sort(y.SRV)
	if !slices.Equal(x.SRV, y.SRV) {
		return false
	}

	sort.Strings(x.TXT)
	sort.Strings(y.TXT)
	return slices.Equal(x.TXT, y.TXT)
}
