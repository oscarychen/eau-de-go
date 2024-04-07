package email_util_test

import (
	"eau-de-go/pkg/email_util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEmailVerificationTokenCreationHappyPath(t *testing.T) {
	email := "test@example.com"
	token, err := email_util.CreateEmailVerificationToken(email)
	assert.NoError(t, err)
	assert.NotNil(t, token)
}

func TestEmailVerificationTokenVerificationHappyPath(t *testing.T) {
	email := "test@example.com"
	token, _ := email_util.CreateEmailVerificationToken(email)
	returnedEmail, err := email_util.VerifyEmailVerificationToken(token)
	assert.NoError(t, err)
	assert.Equal(t, email, returnedEmail)
}

func TestEmailVerificationTokenVerificationInvalidToken(t *testing.T) {
	token := "invalid token"
	email, err := email_util.VerifyEmailVerificationToken(token)
	assert.Error(t, err)
	assert.Empty(t, email)
}

func TestEmailVerificationTokenVerificationExpiredToken(t *testing.T) {
	email := "test@example.com"
	token, _ := email_util.CreateEmailVerificationToken(email)
	email_util.NowFunc = func() time.Time {
		return time.Now().Add(time.Hour * 13)
	}
	email, err := email_util.VerifyEmailVerificationToken(token)
	assert.Error(t, err)
	assert.Empty(t, email)
}
