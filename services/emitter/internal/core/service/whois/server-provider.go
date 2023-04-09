package whois

const (
	DefaultWhoisServer = "whois.iana.org"
)

type ServerProvider interface {
	GetServerByDomain(domain string) (string, error)
}
