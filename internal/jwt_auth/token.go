package jwt_auth

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

func getClaimsWithStandardClaims(claims map[string]interface{}, tokenLife time.Duration) map[string]interface{} {

	newClaims := make(map[string]interface{})

	for key, value := range claims {
		newClaims[key] = value
	}

	now := NowFunc()
	newClaims["iat"] = now.Unix()
	newClaims["jti"] = uuid.New().String()

	if _, ok := newClaims["exp"]; ok == false {
		newClaims["exp"] = now.Add(tokenLife).Unix()
	}
	return newClaims
}

func CreateToken(tokenType TokenType, claims ...map[string]interface{}) (string, map[string]interface{}, error) {
	var tokenClaims map[string]interface{}
	if len(claims) > 0 {
		tokenClaims = claims[0]
	} else {
		tokenClaims = make(map[string]interface{})
	}
	tokenClaims["token_type"] = string(tokenType)

	var tokenLife time.Duration
	if tokenType == Refresh {
		tokenLife = settings.RefreshTokenLife
	} else if tokenType == Access {
		tokenLife = settings.AccessTokenLife
	}

	tokenClaims = getClaimsWithStandardClaims(tokenClaims, tokenLife)

	token := jwt.NewWithClaims(jwt.SigningMethodPS256, jwt.MapClaims(tokenClaims))

	signingKey, err := GetInMemoryRsaKeyPair().GetSigningKey()
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", nil, err
	}

	return tokenString, tokenClaims, nil
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
