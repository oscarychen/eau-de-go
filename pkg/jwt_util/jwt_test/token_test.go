package jwt_test

import (
	"eau-de-go/pkg/jwt_util"
	"eau-de-go/settings"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateRefreshToken(t *testing.T) {
	claims := map[string]interface{}{
		"username": "testuser",
	}

	jwtUtil := jwt_util.NewJwtUtil()

	token, tokenClaims, err := jwtUtil.CreateRefreshToken(claims)
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

func TestDecodeToken(t *testing.T) {
	claims := map[string]interface{}{
		"username": "testuser",
	}

	jwtUtil := jwt_util.NewJwtUtil()

	token, _, err := jwtUtil.CreateRefreshToken(claims)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	decodedClaims, err := jwtUtil.DecodeToken(jwt_util.Refresh, token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if decodedClaims["username"] != "testuser" {
		t.Errorf("Expected username to be 'testuser', got '%v'", decodedClaims["username"])
	}
}

func TestDecodeInvalidToken(t *testing.T) {
	jwtUtil := jwt_util.NewJwtUtil()

	_, err := jwtUtil.DecodeToken(jwt_util.Refresh, "invalidToken")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestTokenExpiry(t *testing.T) {
	claims := map[string]interface{}{
		"username": "testuser",
	}
	jwtUtil := jwt_util.NewJwtUtil()

	jwt_util.NowFunc = func() time.Time {
		return time.Now().Add(-settings.RefreshTokenLife - time.Minute)
	}

	token, _, err := jwtUtil.CreateRefreshToken(claims)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	_, err = jwtUtil.DecodeToken(jwt_util.Refresh, token)
	if err == nil {
		t.Error("Expected error for expired token")
	}
	jwt_util.NowFunc = time.Now
}
