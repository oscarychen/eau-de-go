package jwt_util

import (
	"eau-de-go/internal/settings"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type TokenType string

const (
	Refresh TokenType = "refresh"
	Access  TokenType = "access"
)

var NowFunc = time.Now

func injectStandardClaims(claims map[string]interface{}) {

	now := NowFunc()
	claims["iat"] = now.Unix()
	claims["jti"] = uuid.New().String()
}

func createToken(claims map[string]interface{}) (string, map[string]interface{}, error) {

	injectStandardClaims(claims)
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, jwt.MapClaims(claims))

	signingKey, err := GetInMemoryRsaKeyPair().GetSigningKey()
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", nil, err
	}

	return tokenString, claims, nil
}

func DecodeToken(tokenType TokenType, tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return GetInMemoryRsaKeyPair().GetVerificationKey()
	})

	if err != nil {
		msg := err.Error()
		return nil, &InvalidTokenError{token: tokenString, msg: &msg}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if err := validateTokenTypes(claims, tokenType); err != nil {
			return nil, err
		}
		if err := validateJti(claims); err != nil {
			return nil, err
		}
		return claims, nil
	} else {
		return nil, &InvalidTokenError{token: tokenString}
	}
}

func validateTokenTypes(claims map[string]interface{}, tokenType TokenType) error {
	if claims["token_type"] != string(tokenType) {
		msg := "Invalid token type."
		return &InvalidTokenError{msg: &msg}
	}
	return nil
}

func validateJti(claims map[string]interface{}) error {
	return nil
}

func CreateRefreshToken(claims map[string]interface{}) (string, map[string]interface{}, error) {
	tokenClaims := CopyTokenClaims(claims)
	tokenClaims["token_type"] = Refresh
	tokenClaims["exp"] = NowFunc().Add(settings.RefreshTokenLife).Unix()
	return createToken(tokenClaims)
}

func CreateAccessToken(claims map[string]interface{}) (string, map[string]interface{}, error) {
	tokenClaims := CopyTokenClaims(claims)
	tokenClaims["token_type"] = Access
	tokenClaims["exp"] = NowFunc().Add(settings.AccessTokenLife).Unix()
	return createToken(tokenClaims)
}

func CopyTokenClaims(claims map[string]interface{}) map[string]interface{} {
	copiedClaimns := make(map[string]interface{})
	for key, value := range claims {
		if key != "exp" && key != "iat" && key != "jti" && key != "token_type" {
			copiedClaimns[key] = value
		}
	}
	return copiedClaimns
}
