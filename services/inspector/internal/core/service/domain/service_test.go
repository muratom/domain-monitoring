package domain

import (
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain/dns"
	"testing"
)

func Test_compareResourceRecords(t *testing.T) {
	type args struct {
		x dns.ResourceRecords
		y dns.ResourceRecords
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "All fields are equal",
			args: args{
				x: dns.ResourceRecords{
					A:     []string{"1.2.3.4"},
					AAAA:  []string{"2001:0db8:11a3:09d7:1f34:8a2e:07a0:765d"},
					CNAME: "foo.example.com.",
					MX: dns.MXSlice{
						dns.MX{
							Host: "mail.example.com.",
							Pref: 42,
						},
					},
					NS: dns.NSSlice{
						dns.NS{
							Host: "ns.example.com.",
						},
					},
					SRV: dns.SRVSlice{
						dns.SRV{
							Target:   "bigbox.example.com.",
							Port:     1234,
							Priority: 10,
							Weight:   20,
						},
					},
					TXT: []string{"I need a tatarka"},
				},
				y: dns.ResourceRecords{
					A:     []string{"1.2.3.4"},
					AAAA:  []string{"2001:0db8:11a3:09d7:1f34:8a2e:07a0:765d"},
					CNAME: "foo.example.com.",
					MX: dns.MXSlice{
						dns.MX{
							Host: "mail.example.com.",
							Pref: 42,
						},
					},
					NS: dns.NSSlice{
						dns.NS{
							Host: "ns.example.com.",
						},
					},
					SRV: dns.SRVSlice{
						dns.SRV{
							Target:   "bigbox.example.com.",
							Port:     1234,
							Priority: 10,
							Weight:   20,
						},
					},
					TXT: []string{"I need a tatarka"},
				},
			},
			want: true,
		},
		{
			name: "All fields are diverge",
			args: args{
				x: dns.ResourceRecords{
					A:     []string{"4.3.2.1"},
					AAAA:  []string{"1234:0db8:11a3:09d7:1f34:8a2e:07a0:765d"},
					CNAME: "bar.example.com.",
					MX: dns.MXSlice{
						dns.MX{
							Host: "mail1.example.com.",
							Pref: 43,
						},
					},
					NS: dns.NSSlice{
						dns.NS{
							Host: "ns2.example.com.",
						},
					},
					SRV: dns.SRVSlice{
						dns.SRV{
							Target:   "bigbox1.example.com.",
							Port:     1234,
							Priority: 10,
							Weight:   20,
						},
					},
					TXT: []string{"I need a gubadiya"},
				},
				y: dns.ResourceRecords{
					A:     []string{"1.2.3.4"},
					AAAA:  []string{"2001:0db8:11a3:09d7:1f34:8a2e:07a0:765d"},
					CNAME: "foo.example.com.",
					MX: dns.MXSlice{
						dns.MX{
							Host: "mail.example.com.",
							Pref: 42,
						},
					},
					NS: dns.NSSlice{
						dns.NS{
							Host: "ns.example.com.",
						},
					},
					SRV: dns.SRVSlice{
						dns.SRV{
							Target:   "bigbox.example.com.",
							Port:     1234,
							Priority: 10,
							Weight:   20,
						},
					},
					TXT: []string{"I need a tatarka"},
				},
			},
			want: false,
		},
		{
			name: "One empty record",
			args: args{
				x: dns.ResourceRecords{},
				y: dns.ResourceRecords{
					A:     []string{"1.2.3.4"},
					AAAA:  []string{"2001:0db8:11a3:09d7:1f34:8a2e:07a0:765d"},
					CNAME: "foo.example.com.",
					MX: dns.MXSlice{
						dns.MX{
							Host: "mail.example.com.",
							Pref: 42,
						},
					},
					NS: dns.NSSlice{
						dns.NS{
							Host: "ns.example.com.",
						},
					},
					SRV: dns.SRVSlice{
						dns.SRV{
							Target:   "bigbox.example.com.",
							Port:     1234,
							Priority: 10,
							Weight:   20,
						},
					},
					TXT: []string{"I need a tatarka"},
				},
			},
			want: false,
		},
		{
			name: "Both records are empty",
			args: args{
				x: dns.ResourceRecords{},
				y: dns.ResourceRecords{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareResourceRecords(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("compareResourceRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}
