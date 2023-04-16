package dns

import (
	"sort"
	"strings"
)

type MX struct {
	Host string
	Pref uint16
}

type MXSlice []MX

func (s MXSlice) Len() int {
	return len(s)
}

func (s MXSlice) Less(i, j int) bool {
	if s[i].Pref < s[j].Pref {
		return true
	}
	if strings.Compare(s[i].Host, s[j].Host) == -1 {
		return true
	}
	return false
}

func (s MXSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type NS struct {
	Host string
}

type NSSlice []NS

func (s NSSlice) Len() int {
	return len(s)
}

func (s NSSlice) Less(i, j int) bool {
	if strings.Compare(s[i].Host, s[j].Host) == -1 {
		return true
	}
	return false
}

func (s NSSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type SRV struct {
	Target   string
	Port     uint16
	Priority uint16
	Weight   uint16
}

type SRVSlice []SRV

func (s SRVSlice) Len() int {
	return len(s)
}

func (s SRVSlice) Less(i, j int) bool {
	if s[i].Priority < s[j].Priority {
		return true
	}
	if s[i].Weight < s[j].Weight {
		return true
	}
	if strings.Compare(s[i].Target, s[j].Target) == -1 {
		return true
	}
	return false
}

func (s SRVSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ResourceRecords struct {
	A     []string
	AAAA  []string
	CNAME string
	MX    MXSlice
	NS    NSSlice
	SRV   SRVSlice
	TXT   []string
}

func (rr *ResourceRecords) Sort() {
	sort.Strings(rr.A)
	sort.Strings(rr.AAAA)
	sort.Sort(rr.MX)
	sort.Sort(rr.NS)
	sort.Sort(rr.SRV)
	sort.Strings(rr.TXT)
}
