package dns

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

type ResourceRecords struct {
	A     []string
	AAAA  []string
	CNAME string
	MX    []MX
	NS    []NS
	SRV   []SRV
	TXT   []string
}
