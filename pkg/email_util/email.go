package email_util

import (
	"eau-de-go/settings"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

func makeMailBytes(toAddresses []string, mailSubject string, mailMessage string) []byte {
	mail := fmt.Sprintf("Subject: %s\r\n\r\n%s", mailSubject, mailMessage)

	if len(toAddresses) > 0 {
		mail = fmt.Sprintf("To: %s\r\n", strings.Join(toAddresses, ", ")) + mail
	}

	return []byte(mail)
}

// SendSingleEmail sends an email to a single user
func SendSingleEmail(recipientEmail string, mailSubject string, mailBody string) {

	recipients := []string{recipientEmail}

	mailBytes := makeMailBytes(recipients, mailSubject, mailBody)

	fullServerAddress := settings.EmailHost + ":" + settings.EmailPort
	auth := smtp.PlainAuth("", settings.EmailHostUser, settings.EmailHostPassword, settings.EmailHost)
	err := smtp.SendMail(fullServerAddress, auth, settings.EmailHostUser, recipients, mailBytes)

	if err != nil {
		log.Fatal(err)
	}
}

// SendMassEmail sends an email to multiple users, where user emails are not included in the email body to avoid recipients from seeing each other's email addresses
func SendMassEmail(recipientEmails []string, mailSubject string, mailBody string) {
	mailBytes := makeMailBytes([]string{}, mailSubject, mailBody)

	fullServerAddress := settings.EmailHost + ":" + settings.EmailPort
	auth := smtp.PlainAuth("", settings.EmailHostUser, settings.EmailHostPassword, settings.EmailHost)
	err := smtp.SendMail(fullServerAddress, auth, settings.EmailHostUser, recipientEmails, mailBytes)

	if err != nil {
		log.Fatal(err)
	}
}
