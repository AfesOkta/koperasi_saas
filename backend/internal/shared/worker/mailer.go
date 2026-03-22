package worker

import (
	"crypto/tls"

	"github.com/koperasi-gresik/backend/config"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	config config.SMTPConfig
}

func NewMailer(cfg config.SMTPConfig) *Mailer {
	return &Mailer{config: cfg}
}

func (m *Mailer) SendEmail(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.config.FromName+" <"+m.config.FromEmail+">")
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	d := gomail.NewDialer(m.config.Host, m.config.Port, m.config.Username, m.config.Password)
	// For production, a valid TLS config should be used.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(msg)
}
