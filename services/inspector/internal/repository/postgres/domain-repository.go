package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/groupcache/lru"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/whois"
	"github.com/muratom/domain-monitoring/services/inspector/internal/repository/postgres/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// TODO: split this file in multiple parts
type DomainRepository struct {
	cache *lru.Cache
	Conn  *sql.DB
}

func NewDomainRepository(dbConnection *sql.DB, cacheSize int) *DomainRepository {
	return &DomainRepository{
		Conn:  dbConnection,
		cache: lru.New(cacheSize),
	}
}

func (r *DomainRepository) GetByFQDN(ctx context.Context, fqdn string) (*entity.Domain, error) {
	if v, ok := r.cache.Get(fqdn); ok {
		return v.(*entity.Domain), nil
	}

	domainEntry, err := r.prepareDomainEntry(ctx, fqdn)
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

	canonicalName := domainEntry.R.CanonicalNames[0]

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

	whoisInfo := domainEntry.R.Registrations[0]

	result := &entity.Domain{
		FQDN: domainEntry.FQDN,
		WHOIS: whois.Records{
			DomainName:  domainEntry.FQDN,
			NameServers: []string{},
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

	r.cache.Add(fqdn, result)

	return result, nil
}

func (r *DomainRepository) GetRottenDomainsFQDN(ctx context.Context) ([]string, error) {
	rottenDomains, err := models.Domains(
		qm.Where("(NOW() - updated_at) >= update_delay"),
	).All(ctx, r.Conn)
	if err != nil {
		return nil, fmt.Errorf("error when fetching rotten domains: %w", err)
	}

	result := make([]string, len(rottenDomains))
	for i, domain := range rottenDomains {
		result[i] = domain.FQDN
	}
	return result, nil
}

func (r *DomainRepository) Store(ctx context.Context, domain *entity.Domain) error {
	domainEntry := models.Domain{
		FQDN:        domain.FQDN,
		UpdatedAt:   time.Now(),
		UpdateDelay: "1W", // TODO: move to function parameters
	}

	tx, err := r.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = domainEntry.Insert(ctx, tx, boil.Infer())
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("rollback of the Domain transaction was failed: %w", rollbackErr)
		}
		return fmt.Errorf("failed to insert Domain into DB: %w", err)
	}

	err = insertDNS(ctx, tx, domainEntry, domain)
	if err != nil {
		return err
	}

	err = insertWhois(ctx, tx, domainEntry, domain)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("rollback of the transaction was failed: %w", rollbackErr)
		}
		return fmt.Errorf("commit of the transcation failed: %w", err)
	}

	return nil
}

func (r *DomainRepository) SaveChangelog(ctx context.Context, fqdn string, changelog *entity.Changelog) error {
	domainEntry, err := r.prepareDomainEntry(ctx, fqdn)
	if err != nil {
		return fmt.Errorf("failed to fetch data from DB for FQDN (%s): %w", fqdn, err)
	}

	rawChangelog, err := json.Marshal(changelog)
	if err != nil {
		return fmt.Errorf("error when making raw changelog: %w", err)
	}
	changelogEntry := &models.Changelog{
		CreatedAt: time.Now(),
		Changes:   rawChangelog,
	}

	return domainEntry.AddChangelogs(ctx, r.Conn, true, changelogEntry)
}
