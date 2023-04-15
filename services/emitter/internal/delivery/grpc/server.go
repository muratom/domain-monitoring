package grpc

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/muratom/domain-monitoring/services/emitter/api/proto/gen/go/emitter"
	dnsentity "github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/dns"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/dns"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	stopServing = "STOP_SERVING"
)

type EmitterServer struct {
	// All implementations must embed UnimplementedEmitterServer
	// for forward compatibility
	pb.UnimplementedEmitterServer

	dnsService   dnsService
	whoisService whoisService
}

func NewEmitterServer(dnsService dnsService, whoisService whoisService) *EmitterServer {
	return &EmitterServer{
		dnsService:   dnsService,
		whoisService: whoisService,
	}
}

func (e *EmitterServer) GetDNS(ctx context.Context, req *pb.GetDNSRequest) (*pb.GetDNSResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "EmitterServer.GetDNS cancelled")
	}
	lookupParams := dns.LookupParams{
		FQDN:          req.Fqdn,
		DNSServerHost: req.Host,
	}
	resourceRecords, err := e.dnsService.LookupResourceRecords(ctx, lookupParams)
	if err != nil {
		if errors.Is(err, dns.ErrStopServing) {
			st := status.New(codes.NotFound, fmt.Sprintf("DNS server stop serving domain %s", req.Fqdn))
			br := &errdetails.ErrorInfo{
				Reason: stopServing,
				Domain: "emitter",
				Metadata: map[string]string{
					"fqdn":       req.Fqdn,
					"dns_server": req.Host,
				},
			}
			st, err = st.WithDetails(br)
			if err != nil {
				panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
			}
			return nil, st.Err()
		}
		return nil, fmt.Errorf("failed to lookup resource records for FQDN (%v): %w", req.GetFqdn(), err)
	}

	return buildDNSResponse(ctx, req, resourceRecords), nil
}

func buildDNSResponse(ctx context.Context, req *pb.GetDNSRequest, resourceRecords *dnsentity.ResourceRecords) *pb.GetDNSResponse {
	mx := make([]*pb.MX, len(resourceRecords.MX))
	for i, m := range resourceRecords.MX {
		mx[i] = &pb.MX{
			Host: m.Host,
			Pref: uint32(m.Pref),
		}
	}
	ns := make([]*pb.NS, len(resourceRecords.NS))
	for i, n := range resourceRecords.NS {
		ns[i] = &pb.NS{
			Host: n.Host,
		}
	}
	srv := make([]*pb.SRV, len(resourceRecords.SRV))
	for i, s := range resourceRecords.SRV {
		srv[i] = &pb.SRV{
			Target:   s.Target,
			Port:     uint32(s.Port),
			Priority: uint32(s.Priority),
			Weight:   uint32(s.Weight),
		}
	}
	txt := make([]string, len(resourceRecords.TXT))
	for i, t := range resourceRecords.TXT {
		txt[i] = string(t)
	}

	return &pb.GetDNSResponse{
		Request: req,
		ResourceRecords: &pb.ResourceRecords{
			A:     resourceRecords.A,
			AAAA:  resourceRecords.AAAA,
			CNAME: resourceRecords.CNAME,
			MX:    mx,
			NS:    ns,
			SRV:   srv,
			TXT:   txt,
		},
	}
}

func (e *EmitterServer) GetWhois(ctx context.Context, req *pb.GetWhoisRequest) (*pb.GetWhoisResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "EmitterServer.GetWhois cancelled")
	}
	whoisRecord, err := e.whoisService.FetchWhois(ctx, req.GetFqdn())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch WHOIS for FQDN (%v): %w", req.GetFqdn(), err)
	}

	return &pb.GetWhoisResponse{
		Request: req,
		Records: &pb.WhoisRecords{
			DomainName:  whoisRecord.DomainName,
			NameServers: whoisRecord.NameServers,
			Registrar:   whoisRecord.Registrar,
			Created:     timestamppb.New(whoisRecord.Created),
			PaidTill:    timestamppb.New(whoisRecord.PaidTill),
		},
	}, nil
}
