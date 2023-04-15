package client

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/dns"
	client "github.com/muratom/domain-monitoring/services/emitter/internal/core/service/dns"
)

const (
	defaultConnectionTimeout = 5 * time.Second
)

type LibraryClient struct {
	resolver *net.Resolver
}

func NewLibraryClient(resolver *net.Resolver) *LibraryClient {
	return &LibraryClient{
		resolver: resolver,
	}
}

func (c *LibraryClient) LookupRR(ctx context.Context, lookupParams client.LookupParams) (*dns.ResourceRecords, error) {
	resolver := c.resolver
	if lookupParams.DNSServerHost != "" {
		resolver.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				// TODO: set from config
				Timeout: defaultConnectionTimeout,
			}
			return d.DialContext(ctx, network, fmt.Sprintf("%v:53", lookupParams.DNSServerHost))
		}
	}

	host := lookupParams.FQDN

	ips, err := resolver.LookupIP(ctx, "ip", host)
	if err != nil {
		// Check fact of serving only here because A and AAAA records are basic
		if dnsError, ok := err.(*net.DNSError); ok {
			if dnsError.IsNotFound || strings.Contains(dnsError.Err, "server misbehaving") {
				return nil, client.ErrStopServing
			}
		}
		return nil, fmt.Errorf("failed to get IP addresses for the host (%s): %w", host, err)
	}

	ipv4s := make([]string, 0, len(ips))
	ipv6s := make([]string, 0, len(ips))
	for _, ip := range ips {
		// To4() returns nil if IP address is not IPv4
		if ipv4 := ip.To4(); ipv4 != nil {
			ipv4s = append(ipv4s, ip.String())
		} else {
			ipv6s = append(ipv6s, ip.String())
		}
	}

	cname, err := c.resolver.LookupCNAME(ctx, host)
	if err != nil && isFatalError(err) {
		return nil, fmt.Errorf("failed to get CNAME for the host (%s): %w", host, err)
	}

	resolvedMXs, err := c.resolver.LookupMX(ctx, host)
	if err != nil && isFatalError(err) {
		return nil, fmt.Errorf("failed to get MX for the host (%s): %w", host, err)
	}
	mxs := make([]dns.MX, len(resolvedMXs))
	for i, mx := range resolvedMXs {
		mxs[i] = dns.MX{
			Host: mx.Host,
			Pref: mx.Pref,
		}
	}

	resolvedNSs, err := c.resolver.LookupNS(ctx, host)
	if err != nil && isFatalError(err) {
		return nil, fmt.Errorf("failed to get NS for the host (%s): %w", host, err)
	}
	nss := make([]dns.NS, len(resolvedNSs))
	for i, ns := range resolvedNSs {
		nss[i] = dns.NS{Host: ns.Host}
	}

	_, resolvedSRVs, err := c.resolver.LookupSRV(ctx, "", "", host)
	if err != nil && isFatalError(err) {
		return nil, fmt.Errorf("failed to get SRV for the host (%s): %w", host, err)
	}
	srvs := make([]dns.SRV, len(resolvedSRVs))
	for i, srv := range resolvedSRVs {
		srvs[i] = dns.SRV{
			Target:   srv.Target,
			Port:     srv.Port,
			Priority: srv.Priority,
			Weight:   srv.Weight,
		}
	}

	resolvedTXTs, err := c.resolver.LookupTXT(ctx, host)
	if err != nil && isFatalError(err) {
		return nil, fmt.Errorf("failed to get TXT for the host (%s): %w", host, err)
	}
	txts := make([]dns.TXT, len(resolvedTXTs))
	for i, txt := range resolvedTXTs {
		txts[i] = dns.TXT(txt)
	}

	return &dns.ResourceRecords{
		A:     ipv4s,
		AAAA:  ipv6s,
		CNAME: cname,
		MX:    mxs,
		NS:    nss,
		SRV:   srvs,
		TXT:   txts,
	}, nil
}

func isFatalError(err error) bool {
	switch e := err.(type) {
	case *net.DNSError:
		return !e.IsNotFound
	default:
		return true
	}
}
