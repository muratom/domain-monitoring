package client

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/foxcpp/go-mockdns"
	dnsentity "github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/dns"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/dns"
	"github.com/stretchr/testify/suite"
)

type LibraryClientTestSuite struct {
	suite.Suite
	netResolver   *net.Resolver
	dnsServerMock *mockdns.Server
}

func (s *LibraryClientTestSuite) SetupSuite() {
	s.dnsServerMock, _ = mockdns.NewServer(map[string]mockdns.Zone{
		"www.example.com.": {
			CNAME: "example.com.",
		},
		"example.com.": {
			A:    []string{"1.2.3.4", "42.73.7.2"},
			AAAA: []string{"1234:6b8:a::a"},
			MX: []net.MX{
				{Host: "mail.example.com.", Pref: 10},
			},
			NS: []net.NS{
				{Host: "ns1.example.com."},
				{Host: "ns2.example.com."},
			},
			SRV: []net.SRV{
				{
					Target:   "sipserver.example.com.",
					Port:     72,
					Priority: 0,
				},
			},
			TXT: []string{"abracadabra"},
		},
	}, false)

	s.netResolver = net.DefaultResolver
	s.dnsServerMock.PatchNet(s.netResolver)
}

func (s *LibraryClientTestSuite) TearDownSuite() {
	s.dnsServerMock.Close()
	mockdns.UnpatchNet(net.DefaultResolver)
}

func (s *LibraryClientTestSuite) TestAllResourceRecords() {
	dnsClient := NewLibraryClient(1 * time.Second)
	rr, err := dnsClient.LookupRR(context.Background(), dns.LookupParams{FQDN: "www.example.com"})
	s.Require().NoError(err)
	s.Require().ElementsMatch([]string{"1.2.3.4", "42.73.7.2"}, rr.A)
	s.Require().Equal("example.com.", rr.CNAME)
	s.Require().ElementsMatch([]dnsentity.MX{{Host: "mail.example.com.", Pref: 10}}, rr.MX)
	s.Require().ElementsMatch([]dnsentity.TXT{"abracadabra"}, rr.TXT)

	expectedNS := []dnsentity.NS{
		{Host: "ns1.example.com."},
		{Host: "ns2.example.com."},
	}
	s.Require().ElementsMatch(expectedNS, rr.NS)

	expectedSRV := []dnsentity.SRV{
		{
			Target:   "sipserver.example.com.",
			Port:     72,
			Priority: 0,
		},
	}
	s.Require().ElementsMatch(expectedSRV, rr.SRV)
}

func (s *LibraryClientTestSuite) TestAllResourceRecordsWithDNSServerSet() {
	dnsClient := NewLibraryClient(1 * time.Second)
	s.dnsServerMock.Authoritative = true
	lookupParams := dns.LookupParams{
		FQDN:          "www.example.com",
		DNSServerHost: s.dnsServerMock.LocalAddr().String(),
	}
	rr, err := dnsClient.LookupRR(context.Background(), lookupParams)
	s.Require().NoError(err)
	s.Require().ElementsMatch([]string{"1.2.3.4", "42.73.7.2"}, rr.A)
	s.Require().Equal("example.com.", rr.CNAME)
	s.Require().ElementsMatch([]dnsentity.MX{{Host: "mail.example.com.", Pref: 10}}, rr.MX)
	s.Require().ElementsMatch([]dnsentity.TXT{"abracadabra"}, rr.TXT)

	expectedNS := []dnsentity.NS{
		{Host: "ns1.example.com."},
		{Host: "ns2.example.com."},
	}
	s.Require().ElementsMatch(expectedNS, rr.NS)

	expectedSRV := []dnsentity.SRV{
		{
			Target:   "sipserver.example.com.",
			Port:     72,
			Priority: 0,
		},
	}
	s.Require().ElementsMatch(expectedSRV, rr.SRV)
}

func (s *LibraryClientTestSuite) TestNotServing() {
	dnsClient := NewLibraryClient(1 * time.Second)
	s.dnsServerMock.Authoritative = true
	lookupParams := dns.LookupParams{
		FQDN:          "hotstuff.com",
		DNSServerHost: s.dnsServerMock.LocalAddr().String(),
	}
	_, err := dnsClient.LookupRR(context.Background(), lookupParams)
	s.Require().Error(err)
	s.Require().True(errors.Is(err, dns.ErrStopServing))
}

func TestLibraryClientTestSuite(t *testing.T) {
	suite.Run(t, new(LibraryClientTestSuite))
}
