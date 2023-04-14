package emitterclient

import (
	"context"
	"fmt"
	"time"

	pb "github.com/muratom/domain-monitoring/services/emitter/api/proto/gen/go/emitter"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/whois"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcEmitterClient struct {
	grpcClient pb.EmitterClient
}

func NewGrpcEmitterClient(cc *grpc.ClientConn, requestTimeout time.Duration) *grpcEmitterClient {
	return &grpcEmitterClient{
		grpcClient: pb.NewEmitterClient(cc),
	}
}

func (c *grpcEmitterClient) GetDNS(ctx context.Context, req *service.GetDNSRequest) (*service.GetDNSResponse, error) {
	pbRequest := &pb.GetDNSRequest{
		Fqdn: req.FQDN,
		Host: req.DNSServerHost,
	}
	pbResponse, err := c.grpcClient.GetDNS(ctx, pbRequest)
	if err != nil {
		st := status.Convert(err)
		switch st.Code() {
		case codes.NotFound:
			return nil, service.ErrStopServing
			// for _, details := range st.Details() {
			// 	switch d := details.(type) {
			// 	case *errdetails.ErrorInfo:
			// 		if d.Reason == "NOT_SERVING" {
			// 			return nil, service.ErrStopServing
			// 		}
			// 	}
			// }
		}

		return nil, fmt.Errorf("emitter gRPC client's call failed: %w", err)
	}

	return buildResourceRecordsResponse(ctx, pbResponse), nil
}

func buildResourceRecordsResponse(ctx context.Context, resourceRecords *pb.ResourceRecords) *service.GetDNSResponse {
	mx := make([]dns.MX, len(resourceRecords.MX))
	for i, m := range resourceRecords.MX {
		mx[i] = dns.MX{
			Host: m.Host,
			Pref: uint16(m.Pref),
		}
	}
	ns := make([]dns.NS, len(resourceRecords.NS))
	for i, n := range resourceRecords.NS {
		ns[i] = dns.NS{
			Host: n.Host,
		}
	}
	srv := make([]dns.SRV, len(resourceRecords.SRV))
	for i, s := range resourceRecords.SRV {
		srv[i] = dns.SRV{
			Target:   s.Target,
			Port:     uint16(s.Port),
			Priority: uint16(s.Priority),
			Weight:   uint16(s.Weight),
		}
	}

	return &service.GetDNSResponse{
		ResourceRecords: dns.ResourceRecords{
			A:     resourceRecords.A,
			AAAA:  resourceRecords.AAAA,
			CNAME: resourceRecords.CNAME,
			MX:    mx,
			NS:    ns,
			SRV:   srv,
			TXT:   resourceRecords.TXT,
		},
	}
}

func (c *grpcEmitterClient) GetWhois(ctx context.Context, req *service.GetWhoisRequest) (*service.GetWhoisResponse, error) {
	pbRequest := &pb.GetWhoisRequest{
		Fqdn: req.FQDN,
	}
	pbResponse, err := c.grpcClient.GetWhois(ctx, pbRequest)
	if err != nil {
		return nil, fmt.Errorf("emitter gRPC client's call failed: %w", err)
	}

	return &service.GetWhoisResponse{
		Record: whois.Record{
			DomainName:  pbResponse.GetDomainName(),
			NameServers: pbResponse.GetNameServers(),
			Created:     pbResponse.GetCreated().AsTime(),
			PaidTill:    pbResponse.GetPaidTill().AsTime(),
		},
	}, nil
}
