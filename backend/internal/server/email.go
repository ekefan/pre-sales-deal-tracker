package server

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender interface {
	sendEmail(
		subject string,
		content string,
		to []string,
		// cc []string,
		// bcc []string,
		// attachedFiles []string,
	) error
}

const (
	smtpAuthAddress  = "smtp.gmail.com"
	smtpServerAdress = "smtp.gmail.com:587"
)

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name, fromEmailAddress, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) sendEmail(
	subject string,
	content string,
	to []string,
	// cc []string,
	// bcc []string,
	// attachedFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	// e.Cc = cc
	// e.Bcc = bcc

	// for _, f := range attachedFiles {
	// 	_, err := e.AttachFile(f)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to attach file %s: %w", f, err)
	// 	}
	// }
	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAdress, smtpAuth)
}
