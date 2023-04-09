package whois

import (
	"fmt"
	"net/url"

	"github.com/zonedb/zonedb"
)

type ZoneDBServerProvider struct{}

func NewZoneDBServerProvider() *ZoneDBServerProvider {
	return &ZoneDBServerProvider{}
}

func (p *ZoneDBServerProvider) GetServerByDomain(domain string) (string, error) {
	z := zonedb.PublicZone(domain)
	if z == nil {
		return "", fmt.Errorf("failed to get a public DNS zone for the domain (%s)", domain)
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

	return "", fmt.Errorf("no WHOIS server found for %s", domain)
}
