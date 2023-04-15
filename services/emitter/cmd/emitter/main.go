package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/muratom/domain-monitoring/api/proto/v1/emitter"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/dns"
	dnsclient "github.com/muratom/domain-monitoring/services/emitter/internal/core/service/dns/client"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois"
	adapterprovider "github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois/adapter-provider"
	whoisclient "github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois/adapter/client"
	serverprovider "github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois/server-provider"
	server "github.com/muratom/domain-monitoring/services/emitter/internal/delivery/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalln("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	dnsResolver := &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{
				// TODO: set timeout by config
				Timeout: 1 * time.Second,
			}
			return dialer.DialContext(ctx, network, address)
		},
	}
	dnsService := dns.NewService(dnsclient.NewLibraryClient(dnsResolver))

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
