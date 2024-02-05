package email_util

import (
	"net/mail"
)

func ValidateEmail(email string) (string, error) {
	parsed, err := mail.ParseAddress(email)
	if err != nil {
		return "", &InvalidEmailError{Key: err.Error()}
	}
	return parsed.Address, nil
}
