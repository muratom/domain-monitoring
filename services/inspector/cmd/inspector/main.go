package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/muratom/domain-monitoring/api/rpc/v1/inspector"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service"
	emitterclient "github.com/muratom/domain-monitoring/services/inspector/internal/core/service/emitter-client"
	inspectorserver "github.com/muratom/domain-monitoring/services/inspector/internal/delivery/http"
	"github.com/muratom/domain-monitoring/services/inspector/internal/repository/postgres"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	emitterAddresses := []string{
		"emitter_1:8080",
		"emitter_2:8080",
	}
	emitters := make([]service.EmitterClient, 0, 2)
	for _, address := range emitterAddresses {
		conn, err := grpc.DialContext(ctx, address, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logrus.Fatalf("unable to connect to the emitter: %v", err)
		}
		logrus.Infof("successfully connect to an emitter at address %v", address)
		emitters = append(emitters, emitterclient.NewGrpcEmitterClient(conn, 10*time.Second))
	}

	dbConn, err := sql.Open("postgres", "host=db port=5432 dbname=domain user=user sslmode=disable password=root")
	if err != nil {
		logrus.Fatalf("connection to DB was failed: %v", err)
	}
	logrus.Infof("successfully connect to a database")
	repo := postgres.NewDomainRepository(dbConn)

	domainService := service.NewDomainService(emitters, repo)

	// mailNotifier := service.NewMailNotifier("<from>", "<to>", "<username>", "<password>", "<smtp_host>", 42)
	stdoutNotifier := &service.StdoutNotifier{}
	notifiers := []service.Notifier{
		stdoutNotifier,
		// mailNotifier,
	}

	ticker := time.After(2 * time.Second)
	go func() {
		logrus.Infof("starting cron")
		for {
			select {
			case <-ticker:
				logrus.Infof("Tick")
				rottenFQDNs, err := domainService.GetRottenDomainsFQDN(ctx)
				if err != nil {
					logrus.Warnf("failed to get rotten domains' FQDNs: %v", err)
					continue
				}
				var notifications []entity.Notification
				for _, fqdn := range rottenFQDNs {
					logrus.Infof("check domain (%v) name servers", fqdn)
					nots, err := domainService.CheckDomainNameServers(ctx, fqdn)
					if err != nil {
						logrus.Warnf("failed to check name servers for FQDN (%v): %v", fqdn, err)
					}
					notifications = append(notifications, nots...)

					logrus.Infof("check domain (%v) registration", fqdn)
					nots, err = domainService.CheckDomainRegistration(ctx, fqdn)
					if err != nil {
						logrus.Warnf("failed to check registration for FQDN (%v): %v", fqdn, err)
					}
					notifications = append(notifications, nots...)

					logrus.Infof("check domain (%v) changes", fqdn)
					nots, err = domainService.CheckDomainChanges(ctx, fqdn)
					if err != nil {
						logrus.Warnf("failed to check changes for FQDN (%v): %v", fqdn, err)
					}
					notifications = append(notifications, nots...)

					logrus.Infof("updating domain (%v)", fqdn)
					_, err = domainService.UpdateDomain(ctx, fqdn)
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

	server := inspectorserver.NewInspectorServer(domainService)
	e := echo.New()
	inspector.RegisterHandlers(e, server)

	address := "0.0.0.0:8000"
	logrus.Infof(fmt.Sprintf("start serving at %v", address))
	if err := e.Start(address); err != http.ErrServerClosed {
		logrus.Fatal(err)
	}

	logrus.Infof("exiting...")
}
