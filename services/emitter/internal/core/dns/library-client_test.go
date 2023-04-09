package dns

import (
	"context"
	"net"
	"testing"

	"github.com/foxcpp/go-mockdns"
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
	dnsClient := NewLibraryClient(s.netResolver)
	rr, err := dnsClient.LookupRR(context.Background(), "www.example.com")
	s.Require().NoError(err)
	s.Require().ElementsMatch([]IPv4{{1, 2, 3, 4}, {42, 73, 7, 2}}, rr.A)
	s.Require().Equal("example.com.", rr.CNAME)
	s.Require().ElementsMatch([]MX{{Host: "mail.example.com.", Pref: 10}}, rr.MX)
	s.Require().ElementsMatch([]TXT{"abracadabra"}, rr.TXT)

	expectedNS := []NS{
		{Host: "ns1.example.com."},
		{Host: "ns2.example.com."},
	}
	s.Require().ElementsMatch(expectedNS, rr.NS)

	expectedSRV := []SRV{
		{
			Target:   "sipserver.example.com.",
			Port:     72,
			Priority: 0,
		},
	}
	s.Require().ElementsMatch(expectedSRV, rr.SRV)
}

func TestLibraryClientTestSuite(t *testing.T) {
	suite.Run(t, new(LibraryClientTestSuite))
}
