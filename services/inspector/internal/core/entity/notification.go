package entity

import (
	"fmt"
	"time"
)

const (
	timeLayout = time.RFC822
)

type Notification interface {
	AsHumanReadable() string
}

type TempNotification struct{}

func (n *TempNotification) AsHumanReadable() string {
	return "Hello"
}

type RegistrationExpireSoonNotification struct {
	FQDN      string
	Registrar string
	PaidTill  time.Time
}

func (n *RegistrationExpireSoonNotification) AsHumanReadable() string {
	return fmt.Sprintf(
		"A registration of a domain %v is going to expire soon. "+
			"Till %v contact the registrar %v to pay for prolongation.",
		n.FQDN, n.PaidTill.Format(timeLayout), n.Registrar,
	)
}

type RegistrationExpiredNotification struct {
	FQDN      string
	Registrar string
	PaidTill  time.Time
}

func (n *RegistrationExpiredNotification) AsHumanReadable() string {
	return fmt.Sprintf(
		"A registration of a domain %v was expired on %v. "+
			"Contact the registrar %v to find out if it is possible to register the domain again",
		n.FQDN, n.PaidTill.Format(timeLayout), n.Registrar,
	)
}

type DomainNameChangedNotification struct {
	Old string
	New string
}

func (n *DomainNameChangedNotification) AsHumanReadable() string {
	return fmt.Sprintf(
		"A domain name has changed from %v to %v",
		n.Old, n.New,
	)
}

// TODO: add enum for RecordType
type ResourceRecordChangedNotification struct {
	FQDN       string
	RecordType string
	Path       []string
	From       interface{}
	To         interface{}
}

func (n *ResourceRecordChangedNotification) AsHumanReadable() string {
	n.From = nilToEmpty(n.From)
	n.To = nilToEmpty(n.To)
	return fmt.Sprintf(
		"A DNS resource record %v for a domain %v has changed from %v to %v",
		n.RecordType, n.FQDN, n.From, n.To,
	)
}

type RegistrationInfoChangedNotification struct {
	FQDN string
	Path []string
	From interface{}
	To   interface{}
}

func (n *RegistrationInfoChangedNotification) AsHumanReadable() string {
	n.From = nilToEmpty(n.From)
	n.To = nilToEmpty(n.To)
	return fmt.Sprintf(
		"A domain %v registration information (%v) has changed from %v to %v",
		n.FQDN, n.Path, n.From, n.To,
	)
}

type DomainStoppedBeingServedNotification struct {
	FQDN           string
	NameServerHost string
}

func (n *DomainStoppedBeingServedNotification) AsHumanReadable() string {
	return fmt.Sprintf(
		"A domain %v has stopped being served by DNS server %v",
		n.FQDN, n.NameServerHost,
	)
}

type NameServersNotSynchronizedNotification struct {
	FQDN                       string
	NotSynchronizedNameServers []string
}

func (n *NameServersNotSynchronizedNotification) AsHumanReadable() string {
	return fmt.Sprintf(
		"Name servers %v not synchronized for a domain %v",
		n.NotSynchronizedNameServers, n.FQDN,
	)
}

func nilToEmpty(v interface{}) interface{} {
	if v == nil {
		return "<empty>"
	}
	return v
}
