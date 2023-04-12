package whois

import "time"

type Record struct {
	// Доменное имя
	DomainName string
	// Список DNS-серверов, указанных для домена
	NameServers []string
	// Дата и время регистрации домена в формате UTC
	Created time.Time
	// Дата и время окончания срока регистрации домена в формате UTC
	PaidTill time.Time
}
