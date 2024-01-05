package inspector

import (
	"context"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/domain"
)

type updater interface {
	UpdateDomain(ctx context.Context, fqdn string) (*domain.Domain, error)
}
