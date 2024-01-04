package domain

import (
	"context"
	"github.com/friendsofgo/errors"
)

var (
	ErrDomainNotFound = errors.New("domain not found")
)

type CRUDRepository interface {
	QueryRepository
	MutationRepository
}

type QueryRepository interface {
	ByFQDN(ctx context.Context, fqdn string) (*Domain, error)
	AllDomainsFQDN(ctx context.Context) ([]string, error)
	RottenDomainsFQDN(ctx context.Context) ([]string, error)
}

type MutationRepository interface {
	Store(ctx context.Context, domain *Domain) error
	Update(ctx context.Context, domain *Domain, storedFQDN string) error
	Delete(ctx context.Context, fqdn string) error
}
