package mail

import (
	"fmt"
	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity/notification"
	"gopkg.in/gomail.v2"
)

type Notifier struct {
	from     string
	to       string
	username string
	password string
	smtpHost string
	smtpPort int
}

func New(from, to, username, password, smtpHost string, smtpPort int) *Notifier {
	return &Notifier{
		from:     from,
		to:       to,
		username: username,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (n *Notifier) Name() string {
	return "mail"
}

func (n *Notifier) Notify(notifications []notification.Notification) error {
	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":    {n.from},
		"To":      {n.to},
		"Subject": {"Notification from domain-monitoring"},
	})
	var body string
	for _, n := range notifications {
		body += fmt.Sprintf("%v\n", n.AsHumanReadable())
	}
	if len(body) == 0 {
		return nil
	}
	m.SetBody("text/plain", body)

	dialer := gomail.NewDialer(n.smtpHost, n.smtpPort, n.username, n.password)

	if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
