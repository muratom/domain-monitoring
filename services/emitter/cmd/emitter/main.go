package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/dns"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/dns/client"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois"
	adapterprovider "github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois/adapter-provider"
	serverprovider "github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois/server-provider"
	"github.com/muratom/domain-monitoring/services/emitter/internal/delivery/grpc/emitter"
	pb "github.com/muratom/domain-monitoring/services/emitter/internal/delivery/grpc/emitter"
	"github.com/muratom/domain-monitoring/services/emitter/internal/delivery/grpc/emitter/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	dnsService := dns.NewService(client.NewLibraryClient(net.DefaultResolver))
	whoisService := whois.NewService(adapterprovider.NewHardcodeAdapterProvider(serverprovider.NewZoneDBServerProvider()))
	emitterServer := server.NewEmitterServer(dnsService, whoisService)
	emitter.RegisterEmitterServer(grpcServer, emitterServer)

	log.Printf("Serving gRPC on http://0.0.0.0:8080")
	go func() {
		log.Fatal(grpcServer.Serve(listener))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:8080",
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
		Addr:    "localhost:8090",
		Handler: gwmux,
	}
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
