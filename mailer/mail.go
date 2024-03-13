package mailer

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

// config
const CMEmail = "auto@choicemovers.com"
const CMEmailPw = "secret"

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

func sendMail(
	subject string,
	to []string,
	message string,
	attachments []string,
) error {
	m := email.NewEmail()
	m.From = CMEmail
	m.To = to
	m.Text = []byte(message)

	for _, f := range attachments {
		_, err := m.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", CMEmail, CMEmailPw, smtpAuthAddress)
	return m.Send(smtpServerAddress, smtpAuth)
}
