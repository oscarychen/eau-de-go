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

func addStandardClaims(claims *map[string]interface{}, tokenLife time.Duration) {
	if *claims == nil {
		*claims = make(map[string]interface{})
	}

	now := NowFunc()
	(*claims)["iat"] = now.Unix()
	(*claims)["exp"] = now.Add(tokenLife).Unix()
	(*claims)["jti"] = uuid.New().String()
}

func CreateToken(tokenType TokenType, claims ...map[string]interface{}) (string, error) {
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

	addStandardClaims(&tokenClaims, tokenLife)

	token := jwt.NewWithClaims(jwt.SigningMethodPS256, jwt.MapClaims(tokenClaims))

	signingKey, err := GetInMemoryRsaKeyPair().GetSigningKey()
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
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
