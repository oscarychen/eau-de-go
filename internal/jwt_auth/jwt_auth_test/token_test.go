package jwt_auth_test

import (
	"eau-de-go/internal/jwt_auth"
	"eau-de-go/internal/settings"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateToken_SuccessfulCreation(t *testing.T) {
	tokenType := jwt_auth.Access
	token, err := jwt_auth.CreateToken(tokenType)

	assert.Nil(t, err, "CreateToken should not return an error")
	assert.NotNil(t, token, "CreateToken should return a valid token")
}

func TestDecodeToken_SuccessfulDecoding(t *testing.T) {
	tokenType := jwt_auth.Access
	token, _ := jwt_auth.CreateToken(tokenType)
	claims, err := jwt_auth.DecodeToken(tokenType, token)

	assert.Nil(t, err, "DecodeToken should not return an error")
	assert.NotNil(t, claims, "DecodeToken should return valid claims")
}

func TestDecodeToken_InvalidToken(t *testing.T) {
	tokenType := jwt_auth.Access
	_, err := jwt_auth.DecodeToken(tokenType, "invalid")

	assert.NotNil(t, err, "DecodeToken should return an error for invalid token")
}

func TestDecodeToken_InvalidTokenType(t *testing.T) {
	tokenType := jwt_auth.Access
	token, _ := jwt_auth.CreateToken(tokenType)
	_, err := jwt_auth.DecodeToken(jwt_auth.TokenType("invalid"), token)

	assert.NotNil(t, err, "DecodeToken should return an error for invalid token type")
}

func TestDecodeToken_ExpiredToken(t *testing.T) {

	jwt_auth.NowFunc = func() time.Time {
		return time.Now().Add(-settings.AccessTokenLife - time.Second)
	}
	tokenType := jwt_auth.Access
	token, _ := jwt_auth.CreateToken(tokenType)

	_, err := jwt_auth.DecodeToken(tokenType, token)

	assert.NotNil(t, err, "DecodeToken should return an error for expired token")
}
