package whois

type AdapterProvider interface {
	GetAdapterByFQDN(fqdn string) Adapter
}
