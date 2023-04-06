package serverprovider

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/muratom/domain-monitoring/services/emitter/internal/whois"
	"github.com/zonedb/zonedb"
)

type ZoneDBServerProvider struct{}

func NewZoneDBServerProvider() *ZoneDBServerProvider {
	return &ZoneDBServerProvider{}
}

func (p *ZoneDBServerProvider) GetServerByFQDN(fqdn string) (string, error) {
	if !strings.Contains(fqdn, ".") {
		return whois.DefaultWhoisServer, nil
	}
	z := zonedb.PublicZone(fqdn)
	if z == nil {
		return "", fmt.Errorf("failed to get a public DNS zone for %s", fqdn)
	}

	// Try whois URL first (these are relatively rare)
	whoisURL := z.WhoisURL()
	if whoisURL != "" {
		parsedURL, err := url.Parse(whoisURL)
		if err == nil && parsedURL.Host != "" {
			return parsedURL.Host, nil
		}
	}

	// Then try host (more common)
	whoisServer := z.WhoisServer()
	if whoisServer != "" {
		return whoisServer, nil
	}

	return "", fmt.Errorf("no WHOIS server found for %s", fqdn)
}
