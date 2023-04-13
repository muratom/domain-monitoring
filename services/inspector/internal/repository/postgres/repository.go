package postgres

import (
	"context"
	"database/sql"
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

	domainName, err := models.Domains(
		models.DomainWhere.FQDN.EQ(fqdn),
		qm.Load(models.DomainRels.Ipv4Addresses),
		qm.Load(models.DomainRels.Ipv6Addresses),
		qm.Load(models.DomainRels.CanonicalNames),
		qm.Load(models.DomainRels.MailExchangers),
		qm.Load(models.DomainRels.NameServers),
		qm.Load(models.DomainRels.ServerSelections),
		qm.Load(models.DomainRels.TextStrings),
		qm.Load(models.DomainRels.Registrations),
	).One(ctx, r.Conn)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from DB for FQDN (%s): %w", fqdn, err)
	}

	ip4 := make([]string, len(domainName.R.Ipv4Addresses))
	for i, ip := range domainName.R.Ipv4Addresses {
		ip4[i] = ip.IP
	}

	ip6 := make([]string, len(domainName.R.Ipv6Addresses))
	for i, ip := range domainName.R.Ipv6Addresses {
		ip6[i] = ip.IP
	}

	canonicalName := domainName.R.CanonicalNames[0]

	mx := make([]dns.MX, len(domainName.R.MailExchangers))
	for i, m := range domainName.R.MailExchangers {
		mx[i] = dns.MX{
			Host: m.Host,
			Pref: uint16(m.Pref),
		}
	}

	ns := make([]dns.NS, len(domainName.R.NameServers))
	for i, n := range domainName.R.NameServers {
		ns[i] = dns.NS{
			Host: n.NameServer,
		}
	}

	srv := make([]dns.SRV, len(domainName.R.ServerSelections))
	for i, s := range domainName.R.ServerSelections {
		srv[i] = dns.SRV{
			Target:   s.Target,
			Port:     uint16(s.Port),
			Priority: uint16(s.Priority),
			Weight:   uint16(s.Weight),
		}
	}

	txt := make([]string, len(domainName.R.TextStrings))
	for i, t := range domainName.R.TextStrings {
		txt[i] = t.Text
	}

	whoisInfo := domainName.R.Registrations[0]

	result := &entity.Domain{
		FQDN: domainName.FQDN,
		WHOIS: whois.Record{
			DomainName:  domainName.FQDN,
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

func (r *DomainRepository) Store(ctx context.Context, domain entity.Domain) error {
	domainName := models.Domain{
		FQDN:        domain.FQDN,
		UpdateAt:    time.Now(),
		UpdateDelay: "1W", // TODO: move to function parameters
	}

	tx, err := r.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = domainName.Insert(ctx, tx, boil.Infer())
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("rollback of the Domain transaction was failed: %w", rollbackErr)
		}
		return fmt.Errorf("failed to insert Domain into DB: %w", err)
	}

	err = insertDNS(ctx, tx, domainName, domain)
	if err != nil {
		return err
	}

	err = insertWhois(ctx, tx, domainName, domain)
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

func insertWhois(ctx context.Context, tx *sql.Tx, dbDomain models.Domain, domain entity.Domain) error {
	whoisRecord := &models.Registration{
		DomainID: dbDomain.ID,
		Created:  domain.WHOIS.Created,
		PaidTill: domain.WHOIS.PaidTill,
	}
	err := dbDomain.AddRegistrations(ctx, tx, true, whoisRecord)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("rollback of the WHOIS transaction was failed: %w", rollbackErr)
		}
		return fmt.Errorf("failed to insert WHOIS record into DB: %w", err)
	}
	return nil
}

func insertDNS(ctx context.Context, tx *sql.Tx, dbDomain models.Domain, domain entity.Domain) error {
	err := insertIPv4(ctx, tx, dbDomain, domain.DNS.A)
	if err != nil {
		return err
	}

	err = insertIPv6(ctx, tx, dbDomain, domain.DNS.AAAA)
	if err != nil {
		return err
	}

	err = insertCNAME(ctx, tx, dbDomain, domain.DNS.CNAME)
	if err != nil {
		return err
	}

	err = insertMX(ctx, tx, dbDomain, domain.DNS.MX)
	if err != nil {
		return err
	}

	err = insertNS(ctx, tx, dbDomain, domain.DNS.NS)
	if err != nil {
		return err
	}

	err = insertSRV(ctx, tx, dbDomain, domain.DNS.SRV)
	if err != nil {
		return err
	}

	err = insertTXT(ctx, tx, dbDomain, domain.DNS.TXT)
	if err != nil {
		return err
	}

	return nil
}

func insertIPv4(ctx context.Context, tx *sql.Tx, domain models.Domain, ips []string) error {
	bulkIPv4 := make([]*models.Ipv4Address, len(ips))
	for i, ip := range ips {
		bulkIPv4[i] = &models.Ipv4Address{
			DomainID: domain.ID,
			IP:       ip,
		}
	}
	err := domain.AddIpv4Addresses(ctx, tx, true, bulkIPv4...)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("rollback of the IPv4 transaction was failed: %w", rollbackErr)
		}
		return fmt.Errorf("failed to insert IPv4 addresses into DB: %w", err)
	}
	return nil
}

func insertIPv6(ctx context.Context, tx *sql.Tx, domain models.Domain, ips []string) error {
	bulkIPv6 := make([]*models.Ipv6Address, len(ips))
	for i, ip := range ips {
		bulkIPv6[i] = &models.Ipv6Address{
			DomainID: domain.ID,
			IP:       ip,
		}
	}
	err := domain.AddIpv6Addresses(ctx, tx, true, bulkIPv6...)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("rollback of the IPv6 transaction was failed: %w", rollbackErr)
		}
		return fmt.Errorf("failed to insert IPv6 addresses into DB: %w", err)
	}
	return nil
}

func insertCNAME(ctx context.Context, tx *sql.Tx, domain models.Domain, cname string) error {
	canonicalName := &models.CanonicalName{
		DomainID:      domain.ID,
		CanonicalName: cname,
	}
	err := domain.AddCanonicalNames(ctx, tx, true, canonicalName)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("rollback of the CNAME transaction was failed: %w", rollbackErr)
		}
		return fmt.Errorf("failed to insert CNAME into DB: %w", err)
	}
	return nil
}

func insertMX(ctx context.Context, tx *sql.Tx, domain models.Domain, mxs []dns.MX) error {
	bulkMX := make([]*models.MailExchanger, len(mxs))
	for i, mx := range mxs {
		bulkMX[i] = &models.MailExchanger{
			DomainID: domain.ID,
			Host:     mx.Host,
			Pref:     int(mx.Pref),
		}
	}
	err := domain.AddMailExchangers(ctx, tx, true, bulkMX...)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("rollback of the MX transaction was failed: %w", rollbackErr)
		}
		return fmt.Errorf("failed to insert MX into DB: %w", err)
	}
	return nil
}

func insertNS(ctx context.Context, tx *sql.Tx, domain models.Domain, nss []dns.NS) error {
	bulkNS := make([]*models.NameServer, len(nss))
	for i, ns := range nss {
		bulkNS[i] = &models.NameServer{
			DomainID:   domain.ID,
			NameServer: ns.Host,
		}
	}
	err := domain.AddNameServers(ctx, tx, true, bulkNS...)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("rollback of the NS transaction was failed: %w", rollbackErr)
		}
		return fmt.Errorf("failed to insert NS into DB: %w", err)
	}
	return nil
}

func insertSRV(ctx context.Context, tx *sql.Tx, domain models.Domain, srvs []dns.SRV) error {
	bulkSRV := make([]*models.ServerSelection, len(srvs))
	for i, srv := range srvs {
		bulkSRV[i] = &models.ServerSelection{
			DomainID: domain.ID,
			Target:   srv.Target,
			Port:     int(srv.Port),
			Priority: int(srv.Priority),
			Weight:   int(srv.Weight),
		}
	}
	err := domain.AddServerSelections(ctx, tx, true, bulkSRV...)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("rollback of the SRV transaction was failed: %w", rollbackErr)
		}
		return fmt.Errorf("failed to insert SRV into DB: %w", err)
	}
	return nil
}

func insertTXT(ctx context.Context, tx *sql.Tx, domain models.Domain, txts []string) error {
	bulkTXT := make([]*models.TextString, len(txts))
	for i, txt := range txts {
		bulkTXT[i] = &models.TextString{
			DomainID: domain.ID,
			Text:     txt,
		}
	}
	err := domain.AddTextStrings(ctx, tx, true, bulkTXT...)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("rollback of the TXT transaction was failed: %w", rollbackErr)
		}
		return fmt.Errorf("failed to insert TXT addresses into DB: %w", err)
	}
	return nil
}
