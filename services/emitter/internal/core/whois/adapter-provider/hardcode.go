package adapterprovider

import (
	"strings"
	"time"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/whois"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/whois/ru"
	serverprovider "github.com/muratom/domain-monitoring/services/emitter/internal/core/whois/server-provider"
)

type HardcodeAdapterProvider struct {
	adapterByDomain map[string]whois.Adapter
}

func NewHardcodeAdapterProvider() *HardcodeAdapterProvider {
	ruAdapter := ru.NewAdapter(
		whois.NewWhoisClient(1*time.Minute),
		serverprovider.NewZoneDBServerProvider(),
	)
	adapterByDomain := map[string]whois.Adapter{
		"ru":                 ruAdapter,
		"su":                 ruAdapter,
		"xn--p1ai" /* рф */ : ruAdapter,
	}

	return &HardcodeAdapterProvider{
		adapterByDomain: adapterByDomain,
	}
}

func (p *HardcodeAdapterProvider) GetAdapterByFQDN(fqdn string) whois.Adapter {
	sfx := fqdn
	for {
		if a, ok := p.adapterByDomain[sfx]; ok {
			return a
		}
		if i := strings.Index(sfx, "."); i >= 0 {
			sfx = sfx[i+1:]
		} else {
			break
		}
	}
	return nil
}
