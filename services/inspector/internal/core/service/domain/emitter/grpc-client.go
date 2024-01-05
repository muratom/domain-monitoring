package emitterclient

import (
	"context"
	"fmt"
	pb "github.com/muratom/domain-monitoring/api/proto/v1/emitter"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain/dns"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain/whois"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/service/domain"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcEmitterClient struct {
	grpcClient pb.EmitterClient
}

func NewGrpcEmitterClient(cc *grpc.ClientConn) *grpcEmitterClient {
	return &grpcEmitterClient{
		grpcClient: pb.NewEmitterClient(cc),
	}
}

func (c *grpcEmitterClient) GetDNS(ctx context.Context, req *domain.GetDNSRequest) (*domain.GetDNSResponse, error) {
	pbRequest := &pb.GetDNSRequest{
		Fqdn: req.FQDN,
		Host: req.DNSServerHost,
	}
	pbResponse, err := c.grpcClient.GetDNS(ctx, pbRequest)
	if err != nil {
		st := status.Convert(err)
		switch st.Code() {
		case codes.NotFound:
			for _, details := range st.Details() {
				switch d := details.(type) {
				case *errdetails.ErrorInfo:
					if d.Domain == "emitter" && d.Reason == "STOP_SERVING" {
						return nil, domain.ErrStopServing
					}
				}
			}
		}

		return nil, fmt.Errorf("emitter gRPC client's call failed: %w", err)
	}

	return buildDNSResponse(ctx, req, pbResponse), nil
}

func buildDNSResponse(_ context.Context, req *domain.GetDNSRequest, dnsResponse *pb.GetDNSResponse) *domain.GetDNSResponse {
	mx := make([]dns.MX, len(dnsResponse.ResourceRecords.MX))
	for i, m := range dnsResponse.ResourceRecords.MX {
		mx[i] = dns.MX{
			Host: m.Host,
			Pref: uint16(m.Pref),
		}
	}
	ns := make([]dns.NS, len(dnsResponse.ResourceRecords.NS))
	for i, n := range dnsResponse.ResourceRecords.NS {
		ns[i] = dns.NS{
			Host: n.Host,
		}
	}
	srv := make([]dns.SRV, len(dnsResponse.ResourceRecords.SRV))
	for i, s := range dnsResponse.ResourceRecords.SRV {
		srv[i] = dns.SRV{
			Target:   s.Target,
			Port:     uint16(s.Port),
			Priority: uint16(s.Priority),
			Weight:   uint16(s.Weight),
		}
	}

	resourceRecords := &dns.ResourceRecords{
		A:     dnsResponse.ResourceRecords.A,
		AAAA:  dnsResponse.ResourceRecords.AAAA,
		CNAME: dnsResponse.ResourceRecords.CNAME,
		MX:    mx,
		NS:    ns,
		SRV:   srv,
		TXT:   dnsResponse.ResourceRecords.TXT,
	}
	// Sort resource records because DNS servers can return same values in different order
	resourceRecords.Sort()

	return &domain.GetDNSResponse{
		Request:         *req,
		ResourceRecords: *resourceRecords,
	}
}

func (c *grpcEmitterClient) GetWhois(ctx context.Context, req *domain.GetWhoisRequest) (*domain.GetWhoisResponse, error) {
	pbRequest := &pb.GetWhoisRequest{
		Fqdn: req.FQDN,
	}
	pbResponse, err := c.grpcClient.GetWhois(ctx, pbRequest)
	if err != nil {
		return nil, fmt.Errorf("emitter gRPC client's call failed: %w", err)
	}

	return &domain.GetWhoisResponse{
		Request: *req,
		Records: whois.Records{
			DomainName:  pbResponse.GetRecords().GetDomainName(),
			NameServers: pbResponse.GetRecords().GetNameServers(),
			Registrar:   pbResponse.GetRecords().GetRegistrar(),
			Created:     pbResponse.GetRecords().GetCreated().AsTime(),
			PaidTill:    pbResponse.GetRecords().GetPaidTill().AsTime(),
		},
	}, nil
}
