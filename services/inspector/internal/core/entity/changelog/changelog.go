package changelog

type Changelog []Change

type Change struct {
	// Type of change (add, update, delete)
	OperationType OperationType
	// Type of the field that has been changed
	FieldType FieldType
	// Path to the field that has been changed
	Path []string
	// Initial value
	From any
	// Resulting value
	To any
}

type OperationType string

const (
	CREATE OperationType = "create"
	UPDATE OperationType = "update"
	DELETE OperationType = "delete"
)

type FieldType string

const (
	FQDN  FieldType = "fqdn"
	WHOIS FieldType = "whois"
	DNS   FieldType = "dns"
)
