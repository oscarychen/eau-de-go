package keys_test

import (
	"eau-de-go/pkg/keys"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Tests for InMemoryRsaKeyStore
//

func TestGetInMemoryRsaKeyPair_Singleton(t *testing.T) {
	keyStore1 := keys.GetInMemoryRsaKeyStore()
	keyStore2 := keys.GetInMemoryRsaKeyStore()

	assert.Equal(t, keyStore1, keyStore2, "GetInMemoryRsaKeyStore should always return the same instance")
}

func TestInMemoryRsaKeyPair_GetVerificationKey(t *testing.T) {
	keyStore := keys.GetInMemoryRsaKeyStore()

	verificationKey1, err1 := keyStore.GetVerificationKey()
	assert.Nil(t, err1, "GetVerificationKey should not return an error")

	verificationKey2, err2 := keyStore.GetVerificationKey()

	assert.Nil(t, err2, "GetVerificationKey should not return an error on subsequent calls")

	assert.Equal(t, verificationKey1, verificationKey2, "GetVerificationKey should always return the same key")
}

func TestInMemoryRsaKeyPair_GetSigningKey(t *testing.T) {
	keyStore := keys.GetInMemoryRsaKeyStore()

	signingKey1, err1 := keyStore.GetSigningKey()
	assert.Nil(t, err1, "GetSigningKey should not return an error")

	signingKey2, err2 := keyStore.GetSigningKey()
	assert.Nil(t, err2, "GetSigningKey should not return an error on subsequent calls")

	assert.Equal(t, signingKey1, signingKey2, "GetSigningKey should always return the same key")
}
