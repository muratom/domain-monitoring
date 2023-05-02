package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync/atomic"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/jellydator/ttlcache/v3"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/whois"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slices"
)

// TODO: set by config
const (
	cacheTTL                = 1 * time.Minute
	workerPoolSize          = 5
	emitterClientTimeout    = 3 * time.Second
	expiringDomainThreshold = 30 * 24 * time.Hour // 1 month
)

type DomainService struct {
	emitters         []EmitterClient
	emitterCounter   atomic.Uint32
	domainRepository entity.DomainRepository
	domainDiffer     domainDiffer
	domainTTLCache   domainTTLCache
}

func NewDomainService(
	emitterClients []EmitterClient,
	domainRepo entity.DomainRepository,
	ttlCache domainTTLCache,
) *DomainService {
	differ := &libDomainDiffer{}

	return &DomainService{
		emitters:         emitterClients,
		domainRepository: domainRepo,
		domainDiffer:     differ,
		domainTTLCache:   ttlCache,
	}
}

func (s *DomainService) Start(_ context.Context) {
	// Enable an automatic removal of expired items
	go s.domainTTLCache.Start()
}

func (s *DomainService) Stop(_ context.Context) {
	s.domainTTLCache.Stop()
}

func (s *DomainService) AddDomain(ctx context.Context, fqdn string) (*entity.Domain, error) {
	ctx, span := otel.Tracer("").Start(ctx, "DomainService.AddDomain", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	domain, err := s.getUpdatedDomain(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("DomainService.AddDomain: error getting updated domain data for FQDN (%v): %w", fqdn, err)
	}

	err = s.domainRepository.Store(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("DomainService.AddDomain: failed to store domain in the repository: %w", err)
	}

	return domain, nil
}

func (s *DomainService) GetDomain(ctx context.Context, fqdn string) (*entity.Domain, error) {
	ctx, span := otel.Tracer("").Start(ctx, "DomainService.GetDomain", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	domain, err := s.domainRepository.GetByFQDN(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("DomainService.GetDomain: failed to get domain from repository: %w", err)
	}

	return domain, err
}

func (s *DomainService) GetAllDomainsFQDN(ctx context.Context) ([]string, error) {
	ctx, span := otel.Tracer("").Start(ctx, "DomainService.GetAllDomainsFQDN")
	defer span.End()

	return s.domainRepository.GetAllDomainsFQDN(ctx)
}

func (s *DomainService) UpdateDomain(ctx context.Context, fqdn string) (*entity.Domain, error) {
	ctx, span := otel.Tracer("").Start(ctx, "DomainService.UpdateDomain", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	domain, err := s.getUpdatedDomain(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("DomainService.UpdateDomain: error getting updated domain data for FQDN (%v): %w", fqdn, err)
	}

	err = s.domainRepository.Update(ctx, domain, fqdn)
	if err != nil {
		return nil, fmt.Errorf("DomainService.UpdateDomain: failed to store domain in the repository: %w", err)
	}
	return domain, nil
}

func (s *DomainService) DeleteDomain(ctx context.Context, fqdn string) error {
	ctx, span := otel.Tracer("").Start(ctx, "DomainService.DeleteDomain", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	err := s.domainRepository.Delete(ctx, fqdn)
	if err != nil {
		return fmt.Errorf("DomainService.DeleteDomain: failed to delete domain: %w", err)
	}

	return nil
}

func (s *DomainService) GetRottenDomainsFQDN(ctx context.Context) ([]string, error) {
	ctx, span := otel.Tracer("").Start(ctx, "DomainService.GetRotternDomainsFQDN")
	defer span.End()

	return s.domainRepository.GetRottenDomainsFQDN(ctx)
}

type dnsResult struct {
	response *GetDNSResponse
	request  *GetDNSRequest
	err      error
}

var (
	errStopServing = errors.New("DNS server stopped serving a domain")
)

func (s *DomainService) CheckDomainNameServers(ctx context.Context, fqdn string) ([]entity.Notification, error) {
	ctx, span := otel.Tracer("").Start(ctx, "DomainService.CheckDomainNameServers", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	domain, err := s.domainRepository.GetByFQDN(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("error getting domain by FQDN: %w", err)
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

	notifications := make([]entity.Notification, 0)
	responses := make([]GetDNSResponse, 0, len(nameServers))
	for res := range results {
		if res.err != nil {
			if errors.Is(res.err, ErrStopServing) {
				notifications = append(notifications, &entity.DomainStoppedBeingServedNotification{
					FQDN:           res.request.FQDN,
					NameServerHost: res.request.DNSServerHost,
				})
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
		logrus.Warnf("DomainService.CheckDNSService. DNS servers is not synchronized for FQND %v", fqdn)
		notifications = append(notifications, &entity.NameServersNotSynchronizedNotification{
			FQDN:                       fqdn,
			NotSynchronizedNameServers: notSychronizedDNSServers,
		})
	}

	return notifications, nil
}

func (s *DomainService) CheckDomainRegistration(ctx context.Context, fqdn string) ([]entity.Notification, error) {
	ctx, span := otel.Tracer("").Start(ctx, "DomainService.CheckDomainRegistration", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	emitter := s.getEmitterClient(ctx)
	whoisResp, err := emitter.GetWhois(ctx, &GetWhoisRequest{FQDN: fqdn})
	if err != nil {
		return nil, fmt.Errorf("error getting registration information: %w", err)
	}

	notifications := make([]entity.Notification, 0, 2)

	expiringSoonTimestamp := whoisResp.Records.PaidTill.Add(-expiringDomainThreshold)
	if time.Now().After(expiringSoonTimestamp) {
		// Domain registration is going to expire
		notification := &entity.RegistrationExpiresSoonNotification{
			FQDN:      fqdn,
			Registrar: whoisResp.Records.Registrar,
			PaidTill:  whoisResp.Records.PaidTill,
		}
		notifications = append(notifications, notification)
	} else if time.Now().After(whoisResp.Records.PaidTill) {
		// Domain has been expired
		notification := &entity.RegistrationExpiredNotification{
			FQDN:      fqdn,
			Registrar: whoisResp.Records.Registrar,
			PaidTill:  whoisResp.Records.PaidTill,
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (s *DomainService) CheckDomainChanges(ctx context.Context, fqdn string) ([]entity.Notification, error) {
	ctx, span := otel.Tracer("").Start(ctx, "DomainService.CheckDomainChanges", trace.WithAttributes(
		attribute.String("FQDN", fqdn),
	))
	defer span.End()

	changelog, err := s.getDomainChanges(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("error getting domain changes: %w", err)
	}

	if len(changelog) != 0 {
		err := s.domainRepository.SaveChangelog(ctx, fqdn, &changelog)
		if err != nil {
			return nil, fmt.Errorf("failed to save changelog: %w", err)
		}
	}

	notifications := make([]entity.Notification, 0)
	for _, change := range changelog {
		switch change.FieldType {
		case entity.FQDN:
			notification := &entity.DomainNameChangedNotification{
				Old: change.From.(string),
				New: change.To.(string),
			}
			notifications = append(notifications, notification)
		case entity.DNS:
			recordType, path := change.Path[0], change.Path[1:]
			notification := &entity.ResourceRecordChangedNotification{
				FQDN:       fqdn,
				RecordType: recordType,
				Path:       path,
				From:       change.From,
				To:         change.To,
			}
			notifications = append(notifications, notification)
		case entity.WHOIS:
			notification := &entity.RegistrationInfoChangedNotification{
				FQDN: fqdn,
				Path: change.Path,
				From: change.From,
				To:   change.To,
			}
			notifications = append(notifications, notification)
		}
	}

	return notifications, nil
}

func (s *DomainService) getUpdatedDomain(ctx context.Context, fqdn string) (*entity.Domain, error) {
	if item := s.domainTTLCache.Get(fqdn); item != nil {
		domain := item.Value()
		return &domain, nil
	}

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

	s.domainTTLCache.Set(fqdn, *domain, ttlcache.DefaultTTL)

	return domain, nil
}

func (s *DomainService) getDomainChanges(ctx context.Context, fqdn string) (entity.Changelog, error) {
	rottenDomain, err := s.domainRepository.GetByFQDN(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch domain (%v) data from DB: %w", fqdn, err)
	}

	freshDomain, err := s.getUpdatedDomain(ctx, fqdn)
	if err != nil {
		return nil, fmt.Errorf("error getting updated domain data for FQDN (%v): %w", fqdn, err)
	}

	changelog, err := s.domainDiffer.Diff(rottenDomain, freshDomain)
	if err != nil {
		return nil, fmt.Errorf("error making diff between domains: %w", err)
	}
	return changelog, nil
}

func (s *DomainService) getEmitterClient(ctx context.Context) EmitterClient {
	index := s.emitterCounter.Add(1) % uint32(len(s.emitters))
	return s.emitters[index]
}

func (s *DomainService) isDNSServersSync(responses []GetDNSResponse) (bool, []string) {
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
