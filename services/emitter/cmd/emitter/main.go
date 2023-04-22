package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/muratom/domain-monitoring/api/proto/v1/emitter"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/dns"
	dnsclient "github.com/muratom/domain-monitoring/services/emitter/internal/core/service/dns/client"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois"
	adapterprovider "github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois/adapter-provider"
	whoisclient "github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois/adapter/client"
	serverprovider "github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois/server-provider"
	server "github.com/muratom/domain-monitoring/services/emitter/internal/delivery/grpc"
	"github.com/muratom/domain-monitoring/tools/tracing"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	tp := tracing.InitTracer("emitter", "http://jaeger:14268/api/traces")
	defer func() {
		err := tp.Shutdown(context.Background())
		if err != nil {
			logrus.Fatalf("faield to shutdown tracer provider: %v", err)
		}
	}()

	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalln("failed to listen:", err)
	}

	logger := logrus.New()
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall, logging.PayloadReceived, logging.PayloadSent),
		logging.WithDurationField(logging.DefaultDurationToFields),
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(),
			logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
		),
	)

	dnsService := dns.NewService(dnsclient.NewLibraryClient(5 * time.Second))

	// TODO: set timeout by config
	whoisClient := whoisclient.NewWhoisClient(1 * time.Second)
	whoisService := whois.NewService(adapterprovider.NewHardcodeAdapterProvider(whoisClient, serverprovider.NewZoneDBServerProvider()))
	emitterServer := server.NewEmitterServer(dnsService, whoisService)
	pb.RegisterEmitterServer(grpcServer, emitterServer)

	log.Printf("Serving gRPC on http://0.0.0.0:8080")
	go func() {
		log.Fatal(grpcServer.Serve(listener))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterEmitterHandler(
		context.Background(),
		gwmux,
		conn,
	)
	if err != nil {
		log.Fatalln("failed to register gateway:", err)
	}

	gwServer := http.Server{
		Addr:    "0.0.0.0:8090",
		Handler: gwmux,
	}
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}

func InterceptorLogger(l logrus.FieldLogger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make(map[string]any, len(fields)/2)
		i := logging.Fields(fields).Iterator()
		if i.Next() {
			k, v := i.At()
			f[k] = v
		}
		l = l.WithFields(f)

		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg)
		case logging.LevelInfo:
			l.Info(msg)
		case logging.LevelWarn:
			l.Warn(msg)
		case logging.LevelError:
			l.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
