package inspector

import (
	"context"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/notification"
)

type checker interface {
	CheckDomainNameServers(ctx context.Context, fqdn string) ([]notification.Notification, error)
	CheckDomainRegistration(ctx context.Context, fqdn string) ([]notification.Notification, error)
	CheckDomainChanges(ctx context.Context, fqdn string) ([]notification.Notification, error)
}
