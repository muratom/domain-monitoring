package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"github.com/muratom/domain-monitoring/api/rpc/v1/inspector"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service"
	emitterclient "github.com/muratom/domain-monitoring/services/inspector/internal/core/service/emitter-client"
	inspectorserver "github.com/muratom/domain-monitoring/services/inspector/internal/delivery/http"
	"github.com/muratom/domain-monitoring/services/inspector/internal/repository/postgres"
	"github.com/muratom/domain-monitoring/services/inspector/tools/tracing"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	emitterClientTimeout = 10 * time.Second
)

func main() {
	tp := tracing.InitTracer("inspector", "http://jaeger:14268/api/traces")
	defer func() {
		err := tp.Shutdown(context.Background())
		if err != nil {
			logrus.Fatalf("faield to shutdown tracer provider: %v", err)
		}
	}()
	tracer := otel.Tracer("")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	emitterAddresses := []string{
		"emitter_1:8080",
		"emitter_2:8080",
	}
	emitters := make([]service.EmitterClient, 0, 2)
	for _, address := range emitterAddresses {
		conn, err := grpc.DialContext(
			ctx,
			address,
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(
				timeout.UnaryClientInterceptor(emitterClientTimeout),
				otelgrpc.UnaryClientInterceptor()),
		)
		if err != nil {
			logrus.Fatalf("unable to connect to the emitter: %v", err)
		}
		logrus.Infof("successfully connect to an emitter at address %v", address)
		emitters = append(emitters, emitterclient.NewGrpcEmitterClient(conn))
	}

	dbConn, err := sql.Open("postgres", "host=db port=5432 dbname=domain user=user sslmode=disable password=root")
	if err != nil {
		logrus.Fatalf("connection to DB was failed: %v", err)
	}
	logrus.Infof("successfully connect to a database")
	repo := postgres.NewDomainRepository(dbConn)

	domainService := service.NewDomainService(emitters, repo)
	domainService.Start(ctx)
	defer domainService.Stop(ctx)

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
				ctx, span := tracer.Start(ctx, "tick")
				logrus.Infof("Tick")
				rottenFQDNs, err := domainService.GetRottenDomainsFQDN(ctx)
				if err != nil {
					logrus.Warnf("failed to get rotten domains' FQDNs: %v", err)
					continue
				}

				wp := workerpool.New(len(rottenFQDNs))
				checkResults := make(chan checkResult, len(rottenFQDNs))
				for _, fqdn := range rottenFQDNs {
					fqdn := fqdn
					wp.Submit(func() {
						ctx, cancel := context.WithCancel(ctx)
						defer cancel()

						ctx, span := tracer.Start(ctx, "worker", trace.WithAttributes(
							attribute.String("FQDN", fqdn),
						))

						var notifications []entity.Notification
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

						_, err = domainService.UpdateDomain(ctx, fqdn)
						if err != nil {
							logrus.Warnf("error updating domain (%v): %v", fqdn, err)
						}

						checkResults <- checkResult{
							notifications: notifications,
							err:           err,
						}
						span.End()
					})
				}

				wp.StopWait()
				close(checkResults)

				for _, notifier := range notifiers {
					for result := range checkResults {
						notifier.Notify(result.notifications)
					}
				}
				span.End()
			case <-ctx.Done():
				logrus.Infof("ticker is stopping...")
				return
			}
		}
	}()

	server := inspectorserver.NewInspectorServer(domainService)
	e := getEchoServer()
	inspector.RegisterHandlers(e, server)

	address := "0.0.0.0:8000"
	logrus.Infof(fmt.Sprintf("start serving at %v", address))
	if err := e.Start(address); err != http.ErrServerClosed {
		logrus.Fatal(err)
	}

	logrus.Infof("exiting...")
}

type checkResult struct {
	notifications []entity.Notification
	err           error
}

func getEchoServer() *echo.Echo {
	e := echo.New()
	e.GET("/debug/pprof", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	e.GET("/debug/pprof/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	e.GET("/debug/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	e.GET("/debug/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	e.GET("/debug/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
	e.GET("/debug/pprof/goroutine", echo.WrapHandler(pprof.Handler("goroutine")))
	e.GET("/debug/pprof/heap", echo.WrapHandler(pprof.Handler("heap")))
	e.GET("/debug/pprof/allocs", echo.WrapHandler(pprof.Handler("allocs")))
	e.GET("/debug/pprof/block", echo.WrapHandler(pprof.Handler("block")))
	e.GET("/debug/pprof/mutex", echo.WrapHandler(pprof.Handler("mutex")))
	return e
}
