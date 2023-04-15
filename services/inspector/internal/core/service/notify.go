package service

import (
	"fmt"

	"github.com/muratom/domain-monitoring/services/inspector/internal/core/entity"
	"gopkg.in/gomail.v2"
)

// TODO: move to App.New
type Notifier interface {
	Notify(notifications []entity.Notification) error
}

type StdoutNotifier struct{}

func (s *StdoutNotifier) Notify(notifications []entity.Notification) error {
	for _, notification := range notifications {
		fmt.Println(notification.AsHumanReadable())
	}
	return nil
}

type MailNotifier struct {
	from     string
	to       string
	username string
	password string
	smtpHost string
	smtpPort int
}

func NewMailNotifier(from, to, username, password, smtpHost string, smtpPort int) *MailNotifier {
	return &MailNotifier{
		from:     from,
		to:       to,
		username: username,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (s *MailNotifier) Notify(notifications []entity.Notification) error {
	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":    {s.from},
		"To":      {s.to},
		"Subject": {"Notification from domain-monitoring"},
	})
	var body string
	for _, notification := range notifications {
		body += fmt.Sprintf("%s\n", notification.AsHumanReadable())
	}
	m.SetBody("text/plain", body)

	dialer := gomail.NewDialer(s.smtpHost, s.smtpPort, s.username, s.password)

	if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
