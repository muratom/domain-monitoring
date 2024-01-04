package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain/whois"
	"time"

	"github.com/muratom/domain-monitoring/services/inspector/internal/repository/postgres/models"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TODO: split this file in multiple parts
type DomainRepository struct {
	Conn *sql.DB
}

func NewDomainRepository(dbConnection *sql.DB) *DomainRepository {
	return &DomainRepository{
		Conn: dbConnection,
	}
}

func (r *DomainRepository) ByFQDN(ctx context.Context, fqdn string) (*domain.Domain, error) {
	domainEntry, err := prepareDomainEntry(ctx, fqdn, r.Conn)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from DB for FQDN (%s): %w", fqdn, err)
	}

	ip4 := make([]string, len(domainEntry.R.Ipv4Addresses))
	for i, ip := range domainEntry.R.Ipv4Addresses {
		ip4[i] = ip.IP
	}

	ip6 := make([]string, len(domainEntry.R.Ipv6Addresses))
	for i, ip := range domainEntry.R.Ipv6Addresses {
		ip6[i] = ip.IP
	}

	var canonicalName models.CanonicalName
	if len(domainEntry.R.CanonicalNames) > 0 {
		canonicalName = *domainEntry.R.CanonicalNames[0]
	}

	mx := make([]dns.MX, len(domainEntry.R.MailExchangers))
	for i, m := range domainEntry.R.MailExchangers {
		mx[i] = dns.MX{
			Host: m.Host,
			Pref: uint16(m.Pref),
		}
	}

	ns := make([]dns.NS, len(domainEntry.R.NameServers))
	for i, n := range domainEntry.R.NameServers {
		ns[i] = dns.NS{
			Host: n.NameServer,
		}
	}

	srv := make([]dns.SRV, len(domainEntry.R.ServerSelections))
	for i, s := range domainEntry.R.ServerSelections {
		srv[i] = dns.SRV{
			Target:   s.Target,
			Port:     uint16(s.Port),
			Priority: uint16(s.Priority),
			Weight:   uint16(s.Weight),
		}
	}

	txt := make([]string, len(domainEntry.R.TextStrings))
	for i, t := range domainEntry.R.TextStrings {
		txt[i] = t.Text
	}

	var whoisInfo models.Registration
	if len(domainEntry.R.Registrations) > 0 {
		whoisInfo = *domainEntry.R.Registrations[0]
	}

	result := &domain.Domain{
		FQDN: domainEntry.FQDN,
		WHOIS: whois.Records{
			DomainName:  domainEntry.FQDN,
			NameServers: []string{},
			Registrar:   whoisInfo.Registrar,
			Created:     whoisInfo.Created,
			PaidTill:    whoisInfo.PaidTill,
		},
		DNS: dns.ResourceRecords{
			A:     ip4,
			AAAA:  ip6,
			CNAME: canonicalName.CanonicalName,
			MX:    mx,
			NS:    ns,
			SRV:   srv,
			TXT:   txt,
		},
	}

	return result, nil
}

func (r *DomainRepository) AllDomainsFQDN(ctx context.Context) ([]string, error) {
	domainEntities, err := models.Domains().All(ctx, r.Conn)
	if err != nil {
		return nil, fmt.Errorf("error fetching all domains FQDN: %w", err)
	}

	result := make([]string, len(domainEntities))
	for i, d := range domainEntities {
		result[i] = d.FQDN
	}
	return result, nil
}

func (r *DomainRepository) RottenDomainsFQDN(ctx context.Context) ([]string, error) {
	rottenDomains, err := models.Domains(
		qm.Where("(NOW() - updated_at) >= update_delay"),
	).All(ctx, r.Conn)
	if err != nil {
		return nil, fmt.Errorf("error when fetching rotten domains: %w", err)
	}

	result := make([]string, len(rottenDomains))
	for i, d := range rottenDomains {
		result[i] = d.FQDN
	}
	return result, nil
}

func (r *DomainRepository) Store(ctx context.Context, domain *domain.Domain) (err error) {
	ctx, span := otel.Tracer("").Start(ctx, "DomainRepository.Store", trace.WithAttributes(
		attribute.String("FQDN", domain.FQDN),
	))
	defer span.End()

	domainEntry := models.Domain{
		FQDN:        domain.FQDN,
		UpdatedAt:   time.Now(),
		UpdateDelay: "1W", // TODO: move to function parameters
	}

	tx, err := r.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			logrus.Error("error stroing domain: %v", err)
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				panic(fmt.Sprintf("rollback of the transaction was failed: %v", rollbackErr))
			}
		}
	}()

	err = domainEntry.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to insert Domain into DB: %w", err)
	}

	err = addDNS(ctx, tx, domainEntry, domain)
	if err != nil {
		return err
	}

	err = addWhois(ctx, tx, domainEntry, domain)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit of the transcation failed: %w", err)
	}

	return nil
}

func (r *DomainRepository) Update(ctx context.Context, domain *domain.Domain, storedFQDN string) (err error) {
	// Domain's FQDN may have changed
	domainEntry, err := prepareDomainEntry(ctx, storedFQDN, r.Conn)
	if err != nil {
		return fmt.Errorf("failed to fetch data from DB for FQDN (%s): %w", storedFQDN, err)
	}

	tx, err := r.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			logrus.Error("error updating domain: %v", err)
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				panic(fmt.Sprintf("rollback of the transaction was failed: %v", rollbackErr))
			}
		}
	}()

	err = deleteRelatedEntries(ctx, tx, *domainEntry)
	if err != nil {
		return fmt.Errorf("failed to delete related entries: %w", err)
	}

	domainEntry.UpdatedAt = time.Now()
	rowsUpdated, err := domainEntry.Update(ctx, tx, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to update domain entry: %w", err)
	}
	logrus.Infof("%v rows updated for FQDN (%v)", rowsUpdated, storedFQDN)

	err = addDNS(ctx, tx, *domainEntry, domain)
	if err != nil {
		return err
	}

	err = addWhois(ctx, tx, *domainEntry, domain)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit of the transcation failed: %w", err)
	}

	return nil
}

func (r *DomainRepository) Delete(ctx context.Context, fqdn string) error {
	domainEntry, err := prepareDomainEntry(ctx, fqdn, r.Conn)
	if err != nil {
		return fmt.Errorf("failed to fetch data from DB for FQDN (%v): %w", fqdn, err)
	}

	rowsDeleted, err := domainEntry.Delete(ctx, r.Conn)
	if err != nil {
		return fmt.Errorf("failed to delete from DB data for FQDN (%v): %w", fqdn, err)
	}
	logrus.Infof("%v rows deleted for FQDN (%v)", rowsDeleted, fqdn)

	return nil
}
