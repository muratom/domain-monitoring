package dns

import (
	"context"
	"fmt"
	"net"
)

type LibraryClient struct {
	resolver *net.Resolver
}

func NewLibraryClient(resolver *net.Resolver) *LibraryClient {
	resolver.StrictErrors = true
	return &LibraryClient{
		resolver: resolver,
	}
}

func (c *LibraryClient) LookupRR(ctx context.Context, host string) (*ResourceRecords, error) {
	ips, err := c.resolver.LookupIP(ctx, "ip", host)
	if err != nil {
		return nil, fmt.Errorf("failed to get IP addresses for the host (%s): %w", host, err)
	}

	ipv4s := make([]IPv4, 0, len(ips))
	ipv6s := make([]IPv6, 0, len(ips))
	for _, ip := range ips {
		// To4() returns nil if IP address is not IPv4
		if ipv4 := ip.To4(); ipv4 != nil {
			ipv4s = append(ipv4s, *(*IPv4)(ipv4))
		} else {
			ipv6s = append(ipv6s, *(*IPv6)(ip))
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
	mxs := make([]MX, len(resolvedMXs))
	for i, mx := range resolvedMXs {
		mxs[i] = MX{
			Host: mx.Host,
			Pref: mx.Pref,
		}
	}

	resolvedNSs, err := c.resolver.LookupNS(ctx, host)
	if err != nil && isFatalError(err) {
		return nil, fmt.Errorf("failed to get NS for the host (%s): %w", host, err)
	}
	nss := make([]NS, len(resolvedNSs))
	for i, ns := range resolvedNSs {
		nss[i] = NS{Host: ns.Host}
	}

	_, resolvedSRVs, err := c.resolver.LookupSRV(ctx, "", "", host)
	if err != nil && isFatalError(err) {
		return nil, fmt.Errorf("failed to get SRV for the host (%s): %w", host, err)
	}
	srvs := make([]SRV, len(resolvedSRVs))
	for i, srv := range resolvedSRVs {
		srvs[i] = SRV{
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
	txts := make([]TXT, len(resolvedTXTs))
	for i, txt := range resolvedTXTs {
		txts[i] = TXT(txt)
	}

	return &ResourceRecords{
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
