package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/repository/postgres/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (r *DomainRepository) prepareDomainEntry(ctx context.Context, fqdn string) (*models.Domain, error) {
	return models.Domains(
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
}

func addWhois(ctx context.Context, tx *sql.Tx, dbDomain models.Domain, domain *entity.Domain) error {
	whoisRecord := &models.Registration{
		DomainID:  dbDomain.ID,
		Created:   domain.WHOIS.Created,
		PaidTill:  domain.WHOIS.PaidTill,
		Registrar: domain.WHOIS.Registrar,
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

func addDNS(ctx context.Context, tx *sql.Tx, dbDomain models.Domain, domain *entity.Domain) error {
	err := addIPv4(ctx, tx, dbDomain, domain.DNS.A)
	if err != nil {
		return err
	}

	err = addIPv6(ctx, tx, dbDomain, domain.DNS.AAAA)
	if err != nil {
		return err
	}

	err = addCNAME(ctx, tx, dbDomain, domain.DNS.CNAME)
	if err != nil {
		return err
	}

	err = addMX(ctx, tx, dbDomain, domain.DNS.MX)
	if err != nil {
		return err
	}

	err = addNS(ctx, tx, dbDomain, domain.DNS.NS)
	if err != nil {
		return err
	}

	err = addSRV(ctx, tx, dbDomain, domain.DNS.SRV)
	if err != nil {
		return err
	}

	err = addTXT(ctx, tx, dbDomain, domain.DNS.TXT)
	if err != nil {
		return err
	}

	return nil
}

func addIPv4(ctx context.Context, tx *sql.Tx, domain models.Domain, ips []string) error {
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

func addIPv6(ctx context.Context, tx *sql.Tx, domain models.Domain, ips []string) error {
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

func addCNAME(ctx context.Context, tx *sql.Tx, domain models.Domain, cname string) error {
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

func addMX(ctx context.Context, tx *sql.Tx, domain models.Domain, mxs []dns.MX) error {
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

func addNS(ctx context.Context, tx *sql.Tx, domain models.Domain, nss []dns.NS) error {
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

func addSRV(ctx context.Context, tx *sql.Tx, domain models.Domain, srvs []dns.SRV) error {
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

func addTXT(ctx context.Context, tx *sql.Tx, domain models.Domain, txts []string) error {
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

func deleteRelatedEntries(ctx context.Context, tx *sql.Tx, domainEntry models.Domain) error {
	_, err := domainEntry.R.Ipv4Addresses.DeleteAll(ctx, tx)
	if err != nil {
		return fmt.Errorf("unable to delete IPv4 entries for FQDN (%v): %w", domainEntry.FQDN, err)
	}
	_, err = domainEntry.R.Ipv6Addresses.DeleteAll(ctx, tx)
	if err != nil {
		return fmt.Errorf("unable to delete IPv6 entries for FQDN (%v): %w", domainEntry.FQDN, err)
	}
	_, err = domainEntry.R.CanonicalNames.DeleteAll(ctx, tx)
	if err != nil {
		return fmt.Errorf("unable to delete CNAME entries for FQDN (%v): %w", domainEntry.FQDN, err)
	}
	_, err = domainEntry.R.MailExchangers.DeleteAll(ctx, tx)
	if err != nil {
		return fmt.Errorf("unable to delete MX entries for FQDN (%v): %w", domainEntry.FQDN, err)
	}
	_, err = domainEntry.R.NameServers.DeleteAll(ctx, tx)
	if err != nil {
		return fmt.Errorf("unable to delete NS entries for FQDN (%v): %w", domainEntry.FQDN, err)
	}
	_, err = domainEntry.R.Registrations.DeleteAll(ctx, tx)
	if err != nil {
		return fmt.Errorf("unable to delete WHOIS entries for FQDN (%v): %w", domainEntry.FQDN, err)
	}
	_, err = domainEntry.R.ServerSelections.DeleteAll(ctx, tx)
	if err != nil {
		return fmt.Errorf("unable to delete SRV entries for FQDN (%v): %w", domainEntry.FQDN, err)
	}
	_, err = domainEntry.R.TextStrings.DeleteAll(ctx, tx)
	if err != nil {
		return fmt.Errorf("unable to delete TXT entries for FQDN (%v): %w", domainEntry.FQDN, err)
	}

	return nil
}
