package email_util_test

import (
	"eau-de-go/pkg/email_util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEmailVerificationTokenCreationHappyPath(t *testing.T) {
	email := "test@example.com"
	emailVerifier := email_util.NewEmailTokenVerifier()
	token, err := emailVerifier.CreateToken(email)
	assert.NoError(t, err)
	assert.NotNil(t, token)
}

func TestEmailVerificationTokenVerificationHappyPath(t *testing.T) {
	email := "test@example.com"
	emailVerifier := email_util.NewEmailTokenVerifier()
	token, err := emailVerifier.CreateToken(email)
	returnedEmail, err := emailVerifier.VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, email, returnedEmail)
}

func TestEmailVerificationTokenVerificationInvalidToken(t *testing.T) {
	token := "invalid token"
	emailVerifier := email_util.NewEmailTokenVerifier()
	email, err := emailVerifier.VerifyToken(token)
	assert.Error(t, err)
	assert.Empty(t, email)
}

func TestEmailVerificationTokenVerificationExpiredToken(t *testing.T) {
	email := "test@example.com"
	emailVerifier := email_util.NewEmailTokenVerifier()
	token, err := emailVerifier.CreateToken(email)
	email_util.NowFunc = func() time.Time {
		return time.Now().Add(time.Hour * 13)
	}
	email, err = emailVerifier.VerifyToken(token)
	assert.Error(t, err)
	assert.Empty(t, email)
}
