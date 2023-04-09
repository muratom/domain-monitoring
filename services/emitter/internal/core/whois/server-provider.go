package whois

const (
	DefaultWhoisServer = "whois.iana.org"
)

type ServerProvider interface {
	GetServerByFQDN(fqdn string) (string, error)
}
