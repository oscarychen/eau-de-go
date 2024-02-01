package jwt_test

import (
	"eau-de-go/internal/settings"
	"eau-de-go/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateRefreshToken(t *testing.T) {
	claims := map[string]interface{}{
		"username": "testuser",
	}
	token, tokenClaims, err := jwt.CreateRefreshToken(claims)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if token == "" {
		t.Error("Expected token to be non-empty")
	}
	assert.Equal(t, claims["username"], tokenClaims["username"], "Expected username to be 'testuser'")
	assert.NotNil(t, tokenClaims["iat"], "Expected iat to be set")
	assert.NotNil(t, tokenClaims["jti"], "Expected jti to be set")
	assert.NotNil(t, tokenClaims["exp"], "Expected exp to be set")

	assert.Nil(t, claims["iat"], "Expected claims argument to not be modified")
	assert.Nil(t, claims["jti"], "Expected claims argument to not be modified")
	assert.Nil(t, claims["exp"], "Expected claims argument to not be modified")
}

func TestCreateAccessTokenFromRefreshToken(t *testing.T) {
	claims := map[string]interface{}{
		"username": "testuser",
	}
	refreshToken, _, err := jwt.CreateRefreshToken(claims)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	accessToken, accessTokenClaims, err := jwt.CreateAccessTokenFromRefreshToken(refreshToken)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if accessToken == "" {
		t.Error("Expected token to be non-empty")
	}
	assert.Equal(t, claims["username"], accessTokenClaims["username"], "Expected username to be 'testuser'")
	assert.NotNil(t, accessTokenClaims["iat"], "Expected iat to be set")
	assert.NotNil(t, accessTokenClaims["jti"], "Expected jti to be set")
	assert.NotNil(t, accessTokenClaims["exp"], "Expected exp to be set")

	assert.Nil(t, claims["iat"], "Expected claims argument to not be modified")
	assert.Nil(t, claims["jti"], "Expected claims argument to not be modified")
	assert.Nil(t, claims["exp"], "Expected claims argument to not be modified")

}

func TestCreateAccessTokenFromInvalidRefreshToken(t *testing.T) {
	_, _, err := jwt.CreateAccessTokenFromRefreshToken("invalidToken")
	if err == nil {
		t.Error("Expected error for invalid refresh token")
	}
}

func TestDecodeToken(t *testing.T) {
	claims := map[string]interface{}{
		"username": "testuser",
	}
	token, _, err := jwt.CreateRefreshToken(claims)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	decodedClaims, err := jwt.DecodeToken(jwt.Refresh, token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if decodedClaims["username"] != "testuser" {
		t.Errorf("Expected username to be 'testuser', got '%v'", decodedClaims["username"])
	}

}

func TestDecodeInvalidToken(t *testing.T) {
	_, err := jwt.DecodeToken(jwt.Refresh, "invalidToken")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestTokenExpiry(t *testing.T) {
	claims := map[string]interface{}{
		"username": "testuser",
	}

	jwt.NowFunc = func() time.Time {
		return time.Now().Add(-settings.RefreshTokenLife - time.Minute)
	}

	token, _, err := jwt.CreateRefreshToken(claims)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	_, err = jwt.DecodeToken(jwt.Refresh, token)
	if err == nil {
		t.Error("Expected error for expired token")
	}
	jwt.NowFunc = time.Now
}
