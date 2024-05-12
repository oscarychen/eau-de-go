package jwt_util

import (
	"eau-de-go/pkg/keys"
	"eau-de-go/settings"
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

type JwtUtil interface {
	CreateRefreshToken(claims map[string]interface{}) (string, map[string]interface{}, error)
	CreateAccessToken(claims map[string]interface{}) (string, map[string]interface{}, error)
	DecodeToken(tokenType TokenType, tokenString string) (map[string]interface{}, error)
	CopyTokenClaims(claims map[string]interface{}) map[string]interface{}
}

type jwtUtil struct {
}

func NewJwtUtil() *jwtUtil {
	return &jwtUtil{}
}

func (j *jwtUtil) injectStandardClaims(claims map[string]interface{}) {

	now := NowFunc()
	claims["iat"] = now.Unix()
	claims["jti"] = uuid.New().String()
}

func (j *jwtUtil) createToken(claims map[string]interface{}) (string, map[string]interface{}, error) {

	j.injectStandardClaims(claims)
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, jwt.MapClaims(claims))

	signingKey, err := keys.GetInMemoryRsaKeyStore().GetSigningKey()
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", nil, err
	}

	return tokenString, claims, nil
}

func (j *jwtUtil) DecodeToken(tokenType TokenType, tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return keys.GetInMemoryRsaKeyStore().GetVerificationKey()
	})

	if err != nil {
		msg := err.Error()
		return nil, &InvalidTokenError{token: tokenString, msg: &msg}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if err := j.validateTokenTypes(claims, tokenType); err != nil {
			return nil, err
		}
		if err := j.validateJti(claims); err != nil {
			return nil, err
		}
		return claims, nil
	} else {
		return nil, &InvalidTokenError{token: tokenString}
	}
}

func (j *jwtUtil) validateTokenTypes(claims map[string]interface{}, tokenType TokenType) error {
	if claims["token_type"] != string(tokenType) {
		msg := "Invalid token type."
		return &InvalidTokenError{msg: &msg}
	}
	return nil
}

func (j *jwtUtil) validateJti(claims map[string]interface{}) error {
	return nil
}

func (j *jwtUtil) CreateRefreshToken(claims map[string]interface{}) (string, map[string]interface{}, error) {
	tokenClaims := j.CopyTokenClaims(claims)
	tokenClaims["token_type"] = Refresh
	tokenClaims["exp"] = NowFunc().Add(settings.RefreshTokenLife).Unix()
	return j.createToken(tokenClaims)
}

func (j *jwtUtil) CreateAccessToken(claims map[string]interface{}) (string, map[string]interface{}, error) {
	tokenClaims := j.CopyTokenClaims(claims)
	tokenClaims["token_type"] = Access
	tokenClaims["exp"] = NowFunc().Add(settings.AccessTokenLife).Unix()
	return j.createToken(tokenClaims)
}

func (j *jwtUtil) CopyTokenClaims(claims map[string]interface{}) map[string]interface{} {
	copiedClaimns := make(map[string]interface{})
	for key, value := range claims {
		if key != "exp" && key != "iat" && key != "jti" && key != "token_type" {
			copiedClaimns[key] = value
		}
	}
	return copiedClaimns
}
