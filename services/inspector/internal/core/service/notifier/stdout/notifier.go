package stdout

import (
	"fmt"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/notification"
)

type Notifier struct{}

func New() *Notifier {
	return new(Notifier)
}

func (s *Notifier) Notify(notifications []notification.Notification) error {
	for _, n := range notifications {
		fmt.Println(n.AsHumanReadable())
	}
	return nil
}
