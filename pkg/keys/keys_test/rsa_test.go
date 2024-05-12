package keys_test

import (
	"eau-de-go/pkg/keys"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInMemoryRsaKeyPair_Singleton(t *testing.T) {
	keyPair1 := keys.GetInMemoryRsaKeyStore()
	keyPair2 := keys.GetInMemoryRsaKeyStore()

	assert.Equal(t, keyPair1, keyPair2, "GetInMemoryRsaKeyStore should always return the same instance")
}

func TestInMemoryRsaKeyPair_GetVerificationKey(t *testing.T) {
	keyPair := keys.GetInMemoryRsaKeyStore()

	verificationKey1, err1 := keyPair.GetVerificationKey()
	assert.Nil(t, err1, "GetVerificationKey should not return an error")

	verificationKey2, err2 := keyPair.GetVerificationKey()
	assert.Nil(t, err2, "GetVerificationKey should not return an error on subsequent calls")

	assert.Equal(t, verificationKey1, verificationKey2, "GetVerificationKey should always return the same key")
}

func TestInMemoryRsaKeyPair_GetSigningKey(t *testing.T) {
	keyPair := keys.GetInMemoryRsaKeyStore()

	signingKey1, err1 := keyPair.GetSigningKey()
	assert.Nil(t, err1, "GetSigningKey should not return an error")

	signingKey2, err2 := keyPair.GetSigningKey()
	assert.Nil(t, err2, "GetSigningKey should not return an error on subsequent calls")

	assert.Equal(t, signingKey1, signingKey2, "GetSigningKey should always return the same key")
}
