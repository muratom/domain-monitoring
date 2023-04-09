package whois

type Request struct {
	WhoisServer string
	Body        []byte
}
