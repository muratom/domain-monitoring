package server

import (
	"context"
	"fmt"

	"github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/dns"
	pb "github.com/muratom/domain-monitoring/services/emitter/internal/delivery/grpc/emitter"
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
		return nil, fmt.Errorf("failed to lookup resource records: %w", err)
	}

	return buildResourceRecordsResponse(ctx, resourceRecords), nil
}

func buildResourceRecordsResponse(ctx context.Context, resourceRecords *dns.ResourceRecords) *pb.ResourceRecords {
	a := make([]string, len(resourceRecords.A))
	for i, ip4 := range resourceRecords.A {
		a[i] = string(ip4[:])
	}
	aaaa := make([]string, len(resourceRecords.AAAA))
	for i, ip6 := range resourceRecords.AAAA {
		aaaa[i] = string(ip6[:])
	}
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
		A:     a,
		AAAA:  aaaa,
		CNAME: resourceRecords.CNAME,
		MX:    mx,
		NS:    ns,
		SRV:   srv,
		TXT:   txt,
	}
}

func (e *EmitterServer) GetWhois(_ context.Context, _ *pb.GetWhoisRequest) (*pb.WhoisRecord, error) {
	panic("not implemented") // TODO: Implement
}
