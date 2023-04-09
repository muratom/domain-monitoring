package dns

import "context"

type IPv4 [4]byte
type IPv6 [16]byte

type CNAME string

type MX struct {
	Host string
	Pref uint16
}

type NS struct {
	Host string
}

type SRV struct {
	Target   string
	Port     uint16
	Priority uint16
	Weight   uint16
}

type TXT string

type ResourceRecords struct {
	A     []IPv4
	AAAA  []IPv6
	CNAME string
	MX    []MX
	NS    []NS
	SRV   []SRV
	TXT   []TXT
}

type Client interface {
	LookupRR(ctx context.Context, host string) (*ResourceRecords, error)
}
