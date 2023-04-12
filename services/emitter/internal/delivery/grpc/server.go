package grpc

import (
	"context"
	"fmt"

	pb "github.com/muratom/domain-monitoring/services/emitter/api/proto/gen/go/emitter"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/dns"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (e *EmitterServer) GetDNS(ctx context.Context, req *pb.GetDNSRequest) (*pb.ResourceRecords, error) {
	resourceRecords, err := e.dnsService.LookupResourceRecords(ctx, req.GetFqdn())
	if err != nil {
		return nil, fmt.Errorf("failed to lookup resource records for FQDN (%v): %w", req.GetFqdn(), err)
	}

	return buildResourceRecordsResponse(ctx, resourceRecords), nil
}

func buildResourceRecordsResponse(ctx context.Context, resourceRecords *dns.ResourceRecords) *pb.ResourceRecords {
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

	return &pb.ResourceRecords{
		A:     resourceRecords.A,
		AAAA:  resourceRecords.AAAA,
		CNAME: resourceRecords.CNAME,
		MX:    mx,
		NS:    ns,
		SRV:   srv,
		TXT:   txt,
	}
}

func (e *EmitterServer) GetWhois(ctx context.Context, req *pb.GetWhoisRequest) (*pb.WhoisRecord, error) {
	whoisRecord, err := e.whoisService.FetchWhois(ctx, req.GetFqdn())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch WHOIS for FQDN (%v): %w", req.GetFqdn(), err)
	}

	return &pb.WhoisRecord{
		DomainName:  whoisRecord.DomainName,
		NameServers: whoisRecord.NameServers,
		Created:     timestamppb.New(whoisRecord.Created),
		PaidTill:    timestamppb.New(whoisRecord.PaidTill),
	}, nil
}
