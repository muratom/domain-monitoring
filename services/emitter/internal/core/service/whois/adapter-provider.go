package whois

type adapterProvider interface {
	GetAdapterByFQDN(fqdn string) Adapter
}
