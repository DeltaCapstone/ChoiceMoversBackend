package mailer

import (
	"fmt"
	"net/smtp"

	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
	"github.com/jordan-wright/email"
)

// config

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

func SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("ChoiceMovers <%s>", utils.ServerConfig.EmailSenderAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth(
		"",
		utils.ServerConfig.EmailSenderAddress,
		utils.ServerConfig.EmailSenderPassword,
		smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}
