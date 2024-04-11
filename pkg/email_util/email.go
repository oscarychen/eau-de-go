package email_util

import (
	"eau-de-go/settings"
	"fmt"
	"net/smtp"
	"strings"
)

type EmailSender interface {
	SendSingleEmail(recipientEmail string, mailSubject string, mailBody string) error
	SendMassEmail(recipientEmails []string, mailSubject string, mailBody string) error
}

type emailSender struct {
	EmailHost         string
	EmailPort         string
	EmailHostUser     string
	EmailHostPassword string
}

// NewEmailSender creates a new EmailSender
func NewEmailSender() *emailSender {
	return &emailSender{
		EmailHost:         settings.EmailHost,
		EmailPort:         settings.EmailPort,
		EmailHostUser:     settings.EmailHostUser,
		EmailHostPassword: settings.EmailHostPassword,
	}
}

func (e *emailSender) makeMailBytes(toAddresses []string, mailSubject string, mailMessage string) []byte {
	mail := fmt.Sprintf("Subject: %s\r\n\r\n%s", mailSubject, mailMessage)

	if len(toAddresses) > 0 {
		mail = fmt.Sprintf("To: %s\r\n", strings.Join(toAddresses, ", ")) + mail
	}

	return []byte(mail)
}

// SendSingleEmail sends an email to a single user
func (e *emailSender) SendSingleEmail(recipientEmail string, mailSubject string, mailBody string) error {

	recipients := []string{recipientEmail}

	mailBytes := e.makeMailBytes(recipients, mailSubject, mailBody)

	fullServerAddress := e.EmailHost + ":" + e.EmailPort
	auth := smtp.PlainAuth("", e.EmailHostUser, e.EmailHostPassword, e.EmailHost)
	err := smtp.SendMail(fullServerAddress, auth, e.EmailHostUser, recipients, mailBytes)

	return err
}

// SendMassEmail sends an email to multiple users, where user emails are not included in the email body to avoid recipients from seeing each other's email addresses
func (e *emailSender) SendMassEmail(recipientEmails []string, mailSubject string, mailBody string) error {
	mailBytes := e.makeMailBytes([]string{}, mailSubject, mailBody)

	fullServerAddress := e.EmailHost + ":" + e.EmailPort
	auth := smtp.PlainAuth("", e.EmailHostUser, e.EmailHostPassword, e.EmailHost)
	err := smtp.SendMail(fullServerAddress, auth, e.EmailHostUser, recipientEmails, mailBytes)

	return err
}
