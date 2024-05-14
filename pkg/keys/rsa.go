package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

type RsaKeyStore interface {
	GetVerificationKey() (*rsa.PublicKey, error)
	GetSigningKey() (*rsa.PrivateKey, error)
}

// In-memory RSA key store, for monolithic deployment and development.
// RSA key pair is generated on first access and kept only in memory.
type inMemoryRsaKeyStore struct {
	signingKey      *rsa.PrivateKey
	verificationKey *rsa.PublicKey
}

var inMemoryRsaKeyStoreInstance *inMemoryRsaKeyStore

func GetInMemoryRsaKeyStore() RsaKeyStore {
	if inMemoryRsaKeyStoreInstance == nil {
		inMemoryRsaKeyStoreInstance = &inMemoryRsaKeyStore{}
	}
	return inMemoryRsaKeyStoreInstance
}

func (keyStore *inMemoryRsaKeyStore) makeKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	signingKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to create private key: %s", err))
		return nil, nil, err
	}
	verificationKey := &signingKey.PublicKey

	keyStore.signingKey = signingKey
	keyStore.verificationKey = verificationKey
	fmt.Println("Created new key pair")
	return signingKey, verificationKey, nil
}

func (keyStore *inMemoryRsaKeyStore) GetVerificationKey() (*rsa.PublicKey, error) {
	if keyStore.verificationKey == nil {
		_, _, err := keyStore.makeKeyPair()
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to create key pair: %s", err))
			return nil, nil
		}
	}
	return keyStore.verificationKey, nil
}

func (keyStore *inMemoryRsaKeyStore) GetSigningKey() (*rsa.PrivateKey, error) {
	if keyStore.signingKey == nil {
		_, _, err := keyStore.makeKeyPair()
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to create key pair: %s", err))
			return nil, err
		}
	}
	return keyStore.signingKey, nil
}
