package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service"
	emitterclient "github.com/muratom/domain-monitoring/services/inspector/internal/core/service/emitter-client"
	"github.com/muratom/domain-monitoring/services/inspector/internal/repository/postgres"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("unable to connect to the emitter: %v", err)
	}
	logrus.Infof("successfully connect to an emitter")
	emitterClient := emitterclient.NewGrpcEmitterClient(conn, 3*time.Second)

	dbConn, err := sql.Open("postgres", "host=localhost port=5432 dbname=domain user=user sslmode=disable password=root")
	if err != nil {
		logrus.Fatalf("connection to DB was failed: %v", err)
	}
	logrus.Infof("successfully connect to a database")
	repo := postgres.NewDomainRepository(dbConn)

	domainService := service.NewDomainService([]service.EmitterClient{emitterClient}, repo)

	// mailNotifier := service.NewMailNotifier("<from>", "<to>", "<username>", "<password>", "<smtp_host>", 42)
	stdoutNotifier := &service.StdoutNotifier{}
	notifiers := []service.Notifier{
		stdoutNotifier,
	}

	ticker := time.After(2 * time.Second)
	go func() {
		logrus.Infof("starting cron")
		for {
			select {
			case <-ticker:
				rottenFQDNs, err := domainService.GetRottenDomainsFQDN(ctx)
				if err != nil {
					logrus.Warnf("failed to get rotten domains' FQDNs: %v", err)
					continue
				}
				var notifications []entity.Notification
				for _, fqdn := range rottenFQDNs {
					nots, err := domainService.CheckDomainNameServers(ctx, fqdn)
					if err != nil {
						logrus.Warnf("failed to check name servers for FQDN (%v): %v", fqdn, err)
					}
					notifications = append(notifications, nots...)

					nots, err = domainService.CheckDomainRegistration(ctx, fqdn)
					if err != nil {
						logrus.Warnf("failed to check registration for FQDN (%v): %v", fqdn, err)
					}
					notifications = append(notifications, nots...)

					nots, err = domainService.CheckDomainChanges(ctx, fqdn)
					if err != nil {
						logrus.Warnf("failed to check changes for FQDN (%v): %v", fqdn, err)
					}
					notifications = append(notifications, nots...)

					err = domainService.UpdateDomain(ctx, fqdn)
					if err != nil {
						logrus.Warnf("error updating domain (%v): %v", fqdn, err)
					}
				}

				for _, notifier := range notifiers {
					notifier.Notify(notifications)
				}
			case <-ctx.Done():
				logrus.Infof("ticker is stopping...")
				return
			}
		}
	}()

	time.Sleep(10 * time.Second)
	cancel()
	logrus.Infof("exiting...")
}

func pseudomain() {
	// resolver := &net.Resolver{
	// 	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
	// 		d := net.Dialer{
	// 			Timeout: 1 * time.Second,
	// 		}
	// 		return d.DialContext(ctx, network, "ns1.google.com:53")
	// 	},
	// }

	// host := "hotstuff.com"
	// // ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// // defer cancel()
	// ips, err := resolver.LookupIP(context.Background(), "ip", host)
	// if err != nil {
	// 	e := err.(*net.DNSError)
	// 	fmt.Printf("%+v\n", e)
	// 	fmt.Printf("failed to get IP addresses for the host (%s): %+v", host, err.(*net.DNSError))
	// }
	// fmt.Println(ips)

	// t := time.Now()
	// d1 := entity.Domain{
	// 	FQDN: "a.ru",
	// 	WHOIS: whois.Records{
	// 		DomainName:  "a.ru",
	// 		NameServers: []string{},
	// 		Created:     t,
	// 		PaidTill:    t.Add(5 * time.Minute),
	// 	},
	// 	DNS: dns.ResourceRecords{
	// 		A:     []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
	// 		AAAA:  []string{},
	// 		CNAME: "t.ru",
	// 		MX:    []dns.MX{{Host: "mx.a.ru", Pref: 10}},
	// 		NS:    []dns.NS{{Host: "ns.a.ru"}},
	// 		SRV: []dns.SRV{{
	// 			Target:   "service.a.ru",
	// 			Port:     80,
	// 			Priority: 1,
	// 			Weight:   2,
	// 		}},
	// 		TXT: []string{"abracadabra"},
	// 	},
	// }

	// d2 := entity.Domain{
	// 	FQDN: "b.ru",
	// 	WHOIS: whois.Records{
	// 		DomainName:  "b.ru",
	// 		NameServers: []string{},
	// 		Created:     t,
	// 		PaidTill:    t.Add(5 * time.Minute),
	// 	},
	// 	DNS: dns.ResourceRecords{
	// 		A:     []string{"3.3.3.3", "1.1.1.1", "4.4.4.4"},
	// 		AAAA:  []string{},
	// 		CNAME: "t.ru",
	// 		MX:    []dns.MX{{Host: "mx.a.ru", Pref: 10}, {Host: "mx2.a.ru", Pref: 20}},
	// 		NS:    []dns.NS{{Host: "ns1.a.ru"}},
	// 		SRV: []dns.SRV{
	// 			{
	// 				Target:   "service.a.ru",
	// 				Port:     80,
	// 				Priority: 1,
	// 				Weight:   3,
	// 			},
	// 			{
	// 				Target:   "service2.a.ru",
	// 				Port:     81,
	// 				Priority: 1,
	// 				Weight:   0,
	// 			},
	// 		},
	// 		TXT: []string{},
	// 	},
	// }
	// opt := func(d *diff.Differ) error {
	// 	d.SliceOrdering = true
	// 	d.DisableStructValues = true
	// 	return nil
	// }
	// differ, _ := diff.NewDiffer(opt)
	// diffResult, err := differ.Diff(d1, d2)
	// if err != nil {
	// 	logrus.Panic("AAAA!!!")
	// }
	// fmt.Printf("%+v\n", diffResult)
	// return

	// db, err := sql.Open("postgres", "dbname=domain user=user sslmode=disable password=root")
	// if err != nil {
	// 	logrus.Panicf("connection to DB was failed: %v", err)
	// }

	// repo := postgres.NewDomainRepository(db, 10)

	// changelog := &entity.Changelog{
	// 	entity.Change{
	// 		Type: "create",
	// 		Path: []string{"abracadabra"},
	// 		From: "biba",
	// 		To:   "boba",
	// 	},
	// }
	// err = repo.SaveChangelog(context.Background(), "yandex.ru.", changelog)
	// if err != nil {
	// 	logrus.Panicf("failed to save changelog: %v", err)
	// }

	// db, err := sql.Open("postgres", "dbname=domain user=user sslmode=disable password=root")
	// if err != nil {
	// 	logrus.Panicf("connection to DB was failed: %v", err)
	// }

	// repo := postgres.NewDomainRepository(db, 10)
	// domain, err := repo.GetByFQDN(context.Background(), "yandex.ru.")
	// if err != nil {
	// 	logrus.Panicf("failed to get domain by FQDN: %v", err)
	// }

	// fmt.Printf("%+v", domain)

	// t = time.Now()
	// newDomain := entity.Domain{
	// 	FQDN: "tutu.ru.",
	// 	WHOIS: whois.Record{
	// 		DomainName:  "tutu.ru.",
	// 		NameServers: []string{"abracadabra.ru."},
	// 		Created:     t,
	// 		PaidTill:    t.Add(24 * time.Hour),
	// 	},
	// 	DNS: dns.ResourceRecords{
	// 		A:     []string{"1.1.1.1", "2.2.2.2"},
	// 		AAAA:  []string{"f::f"},
	// 		CNAME: "aaaa.ru",
	// 		MX:    []dns.MX{{Host: "tatar.ru", Pref: 42}},
	// 		NS:    []dns.NS{{Host: "bashkiriya.ru"}},
	// 		SRV:   []dns.SRV{{Target: "someservice", Port: 73, Priority: 0, Weight: 10}},
	// 		TXT:   []string{"secret"},
	// 	},
	// }
	// err = repo.Store(context.Background(), newDomain)
	// if err != nil {
	// 	logrus.Panicf("failed to store domain: %v", err)
	// }
}
