package email_util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"eau-de-go/pkg/keys"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"time"
)

type EmailVerificationTokenClaims struct {
	Email  string
	Expiry time.Time
}

var NowFunc = time.Now

func CreateEmailVerificationToken(email string) (string, error) {

	key := keys.GetInMemoryAesKey()
	fmt.Println(key)
	tokenClaims := EmailVerificationTokenClaims{
		Email:  email,
		Expiry: NowFunc().Add(time.Hour * 12),
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(tokenClaims)
	if err != nil {
		return "", err
	}

	data := buf.Bytes()
	block, err := aes.NewCipher(*key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	token := gcm.Seal(nonce, nonce, data, nil)
	tokenString := base64.StdEncoding.EncodeToString(token)
	return tokenString, nil
}

func VerifyEmailVerificationToken(tokenString string) (string, error) {

	key := keys.GetInMemoryAesKey()

	block, err := aes.NewCipher(*key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	token, err := base64.StdEncoding.DecodeString(tokenString)
	if len(token) < gcm.NonceSize() {
		return "", errors.New("token too short")
	}

	nonce, cipherText := token[:gcm.NonceSize()], token[gcm.NonceSize():]

	data, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	var tokenClaims EmailVerificationTokenClaims
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(&tokenClaims)
	if err != nil {
		return "", err
	}
	if tokenClaims.Expiry.Before(NowFunc()) {
		return "", errors.New("token expired")
	}

	return tokenClaims.Email, nil
}
