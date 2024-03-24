package email_util_test

import (
	"eau-de-go/pkg/email_util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidEmail_HappyPath(t *testing.T) {
	email := "test@example.com"
	_, err := email_util.ValidateEmailAddress(email)
	assert.Nil(t, err)
}

func TestValidEmail_InvalidFormat(t *testing.T) {
	email := "invalid email"
	_, err := email_util.ValidateEmailAddress(email)
	assert.NotNil(t, err)
}

func TestValidEmail_EmptyString(t *testing.T) {
	email := ""
	_, err := email_util.ValidateEmailAddress(email)
	assert.NotNil(t, err)
}
