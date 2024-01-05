package inspector

import (
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/notification"
)

type Notifier interface {
	Notify(notifications []notification.Notification) error
	Name() string
}
