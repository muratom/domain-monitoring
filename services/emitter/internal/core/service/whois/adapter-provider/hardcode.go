package adapterprovider

import (
	"strings"
	"time"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois"
	whoisclient "github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois/adapter/client"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois/adapter/ru"
)

type HardcodeAdapterProvider struct {
	adapterByDomain map[string]whois.Adapter
}

func NewHardcodeAdapterProvider(whoisServerProvider whois.ServerProvider) *HardcodeAdapterProvider {
	ruAdapter := ru.NewAdapter(
		whoisclient.NewWhoisClient(1*time.Minute),
		whoisServerProvider,
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
