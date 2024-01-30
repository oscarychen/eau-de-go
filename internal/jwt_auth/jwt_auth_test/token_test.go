package jwt_auth_test

import (
	"eau-de-go/internal/jwt_auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateTokenReturnsValidTokenAndClaims(t *testing.T) {
	tokenType := jwt_auth.Access
	token, claims, err := jwt_auth.CreateToken(tokenType)

	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.NotNil(t, claims)
	assert.Equal(t, string(tokenType), claims["token_type"])
	assert.NotNil(t, claims["iat"])
	assert.NotNil(t, claims["exp"])
	assert.NotNil(t, claims["jti"])
}

func TestCreateTokenWithExplicitExpClaim(t *testing.T) {
	tokenType := jwt_auth.Access
	claims := map[string]interface{}{
		"exp": 12345,
	}
	token, tokenClaims, err := jwt_auth.CreateToken(tokenType, claims)

	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, 12345, tokenClaims["exp"])
}

func TestCreateTokenReturnsDifferentTokens(t *testing.T) {
	tokenType := jwt_auth.Access
	token1, _, _ := jwt_auth.CreateToken(tokenType)
	token2, _, _ := jwt_auth.CreateToken(tokenType)

	assert.NotEqual(t, token1, token2)
}

func TestDecodeTokenReturnsValidClaims(t *testing.T) {
	tokenType := jwt_auth.Access
	token, _, _ := jwt_auth.CreateToken(tokenType)
	claims, err := jwt_auth.DecodeToken(tokenType, token)

	assert.Nil(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, string(tokenType), claims["token_type"])
}

func TestDecodeTokenReturnsErrorForInvalidToken(t *testing.T) {
	tokenType := jwt_auth.Access
	_, err := jwt_auth.DecodeToken(tokenType, "invalid")

	assert.NotNil(t, err)
}

func TestDecodeTokenReturnsErrorForMismatchedTokenType(t *testing.T) {
	token, _, _ := jwt_auth.CreateToken(jwt_auth.Access)
	_, err := jwt_auth.DecodeToken(jwt_auth.Refresh, token)

	assert.NotNil(t, err)
}

func TestCreateTokenDoesNotModifyClaims(t *testing.T) {
	tokenType := jwt_auth.Access
	claims := map[string]interface{}{}
	_, tokenClaims, _ := jwt_auth.CreateToken(tokenType, claims)

	_, expInClaims := claims["exp"]
	assert.False(t, expInClaims)

	_, expInTokenClaims := tokenClaims["exp"]
	assert.True(t, expInTokenClaims)
}
