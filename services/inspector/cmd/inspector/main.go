package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/whois"
	"github.com/muratom/domain-monitoring/services/inspector/internal/repository/postgres"
	"github.com/sirupsen/logrus"
)

func main() {
	db, err := sql.Open("postgres", "dbname=domain user=user sslmode=disable password=root")
	if err != nil {
		logrus.Panicf("connection to DB was failed: %v", err)
	}

	repo := postgres.NewDomainRepository(db, 10)
	domain, err := repo.GetByFQDN(context.Background(), "yandex.ru.")
	if err != nil {
		logrus.Panicf("failed to get domain by FQDN: %v", err)
	}

	fmt.Printf("%+v", domain)

	t := time.Now()
	newDomain := entity.Domain{
		FQDN: "tutu.ru.",
		WHOIS: whois.Record{
			DomainName:  "tutu.ru.",
			NameServers: []string{"abracadabra.ru."},
			Created:     t,
			PaidTill:    t.Add(24 * time.Hour),
		},
		DNS: dns.ResourceRecords{
			A:     []string{"1.1.1.1", "2.2.2.2"},
			AAAA:  []string{"f::f"},
			CNAME: "aaaa.ru",
			MX:    []dns.MX{{Host: "tatar.ru", Pref: 42}},
			NS:    []dns.NS{{Host: "bashkiriya.ru"}},
			SRV:   []dns.SRV{{Target: "someservice", Port: 73, Priority: 0, Weight: 10}},
			TXT:   []string{"secret"},
		},
	}
	err = repo.Store(context.Background(), newDomain)
	if err != nil {
		logrus.Panicf("failed to store domain: %v", err)
	}
}
