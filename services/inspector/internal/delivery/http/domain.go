package http

import (
	"context"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/changelog"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain"
)

type DomainService interface {
	AddDomain(ctx context.Context, fqdn string) (*domain.Domain, error)
	GetDomain(ctx context.Context, fqdn string) (*domain.Domain, error)
	UpdateDomain(ctx context.Context, fqdn string) (*domain.Domain, error)
	DeleteDomain(ctx context.Context, fqdn string) error
	GetAllDomainsFQDN(ctx context.Context) ([]string, error)
	GetChangelogs(ctx context.Context, fqdn string) ([]changelog.Changelog, error)
}
