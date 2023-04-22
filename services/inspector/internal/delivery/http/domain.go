package http

import (
	"context"

	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
)

type DomainService interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
	AddDomain(ctx context.Context, fqdn string) (*entity.Domain, error)
	GetDomain(ctx context.Context, fqdn string) (*entity.Domain, error)
	UpdateDomain(ctx context.Context, fqdn string) (*entity.Domain, error)
	DeleteDomain(ctx context.Context, fqdn string) error
	GetRottenDomainsFQDN(ctx context.Context) ([]string, error)
	CheckDomainNameServers(ctx context.Context, fqdn string) ([]entity.Notification, error)
	CheckDomainRegistration(ctx context.Context, fqdn string) ([]entity.Notification, error)
	CheckDomainChanges(ctx context.Context, fqdn string) ([]entity.Notification, error)
}
