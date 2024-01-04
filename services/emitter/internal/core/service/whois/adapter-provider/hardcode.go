package adapterprovider

import (
	"strings"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois/adapter"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois/adapter/ru"
)

type Hardcode struct {
	adapterByDomain map[string]whois.Adapter
}

func NewHardcodeAdapterProvider(whoisClient adapter.Client, whoisServerProvider whois.ServerProvider) *Hardcode {
	ruAdapter := ru.NewAdapter(
		whoisClient,
		whoisServerProvider,
	)
	adapterByDomain := map[string]whois.Adapter{
		"ru":                 ruAdapter,
		"su":                 ruAdapter,
		"xn--p1ai" /* рф */ : ruAdapter,
	}

	return &Hardcode{
		adapterByDomain: adapterByDomain,
	}
}

func (p *Hardcode) GetAdapterByFQDN(fqdn string) whois.Adapter {
	// Remove root zone
	sfx := strings.TrimRight(fqdn, ".")
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
