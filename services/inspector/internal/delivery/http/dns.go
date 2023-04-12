package http

// TODO: move to entity
type CheckWorkResult struct{}

type CheckSyncResult struct{}

type dnsService interface {
	GetResourceRecords(fqdn string)
	CheckServersWork(fqdn string) (*CheckWorkResult, error)
	CheckServersSync(fqdn string) (*CheckSyncResult, error)
}
