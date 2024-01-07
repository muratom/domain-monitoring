package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service/domain"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service/domain/emitter"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service/inspector"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service/inspector/notifier/stdout"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	rpc_inspector "github.com/muratom/domain-monitoring/api/rpc/v1/inspector"
	inspectorserver "github.com/muratom/domain-monitoring/services/inspector/internal/delivery/http"
	"github.com/muratom/domain-monitoring/services/inspector/internal/repository/postgres"
	"github.com/muratom/domain-monitoring/services/inspector/tools/tracing"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	emitterClientTimeout = 10 * time.Second
	workerPoolSize       = 5
)

func main() {
	tp := tracing.InitTracer("inspector", "http://jaeger:14268/api/traces")
	defer func() {
		err := tp.Shutdown(context.Background())
		if err != nil {
			logrus.Fatalf("faield to shutdown tracer provider: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	emitterAddresses := []string{
		"emitter_1:8080",
		"emitter_2:8080",
	}
	emitters := make([]domain.EmitterClient, 0, len(emitterAddresses))
	for _, address := range emitterAddresses {
		conn, err := grpc.DialContext(
			ctx,
			address,
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(
				timeout.UnaryClientInterceptor(emitterClientTimeout),
				otelgrpc.UnaryClientInterceptor(),
				retry.UnaryClientInterceptor(retry.WithBackoff(retry.BackoffExponential(100*time.Millisecond))),
			),
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
	domainRepository := postgres.NewDomainRepository(dbConn)
	changelogRepository := postgres.NewChangelogRepository(dbConn)

	domainService := domain.New(emitters, domainRepository, changelogRepository)
	domainService.Start(ctx)
	defer domainService.Stop(ctx)

	// mailNotifier := service.NewMailNotifier("<from>", "<to>", "<username>", "<password>", "<smtp_host>", 42)
	stdoutNotifier := stdout.New()
	notifiers := []inspector.Notifier{
		stdoutNotifier,
		// mailNotifier,
	}

	var domainInspector DomainInspector
	domainInspector = inspector.New(
		domainService,
		domainService,
		domainService,
		inspector.WithTick(1*time.Minute),
		inspector.WithWorkerNumber(5),
		inspector.WithNotifiers(notifiers),
	)
	domainInspector.Start(ctx)
	defer domainInspector.Stop(ctx)

	server := inspectorserver.NewInspectorServer(domainService)
	e := getEchoServer()
	rpc_inspector.RegisterHandlers(e, server)

	address := "0.0.0.0:8000"
	logrus.Infof(fmt.Sprintf("start serving at %v", address))
	if err := e.Start(address); err != http.ErrServerClosed {
		logrus.Fatal(err)
	}

	logrus.Infof("exiting...")
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

type DomainInspector interface {
	service.Runnable
}
