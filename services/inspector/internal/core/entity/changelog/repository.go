package changelog

import "context"

type Repository interface {
	Store(ctx context.Context, fqdn string, changelog Changelog) error
	Retrieve(ctx context.Context, fqdn string) ([]Changelog, error)
}
