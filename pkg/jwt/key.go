package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

type inMemoryRsaKeyPair struct {
	signingKey      *rsa.PrivateKey
	verificationKey *rsa.PublicKey
}

var inMemoryRsaKeyPairInstance *inMemoryRsaKeyPair

func GetInMemoryRsaKeyPair() *inMemoryRsaKeyPair {
	if inMemoryRsaKeyPairInstance == nil {
		inMemoryRsaKeyPairInstance = &inMemoryRsaKeyPair{}
	}
	return inMemoryRsaKeyPairInstance
}

func (keyPair *inMemoryRsaKeyPair) makeKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	signingKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to create private key: %s", err))
		return nil, nil, err
	}
	verificationKey := &signingKey.PublicKey

	keyPair.signingKey = signingKey
	keyPair.verificationKey = verificationKey
	fmt.Println("Created new key pair")
	return signingKey, verificationKey, nil
}

func (keyPair *inMemoryRsaKeyPair) GetVerificationKey() (*rsa.PublicKey, error) {
	if keyPair.verificationKey == nil {
		_, _, err := keyPair.makeKeyPair()
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to create key pair: %s", err))
			return nil, nil
		}
	}
	return keyPair.verificationKey, nil
}

func (keyPair *inMemoryRsaKeyPair) GetSigningKey() (*rsa.PrivateKey, error) {
	if keyPair.signingKey == nil {
		_, _, err := keyPair.makeKeyPair()
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to create key pair: %s", err))
			return nil, err
		}
	}
	return keyPair.signingKey, nil
}
