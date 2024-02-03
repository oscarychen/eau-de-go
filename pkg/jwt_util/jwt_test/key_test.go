package jwt_test

import (
	"eau-de-go/pkg/jwt_util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInMemoryRsaKeyPair_Singleton(t *testing.T) {
	keyPair1 := jwt_util.GetInMemoryRsaKeyPair()
	keyPair2 := jwt_util.GetInMemoryRsaKeyPair()

	assert.Equal(t, keyPair1, keyPair2, "GetInMemoryRsaKeyPair should always return the same instance")
}

func TestInMemoryRsaKeyPair_GetVerificationKey(t *testing.T) {
	keyPair := jwt_util.GetInMemoryRsaKeyPair()

	verificationKey1, err1 := keyPair.GetVerificationKey()
	assert.Nil(t, err1, "GetVerificationKey should not return an error")

	verificationKey2, err2 := keyPair.GetVerificationKey()
	assert.Nil(t, err2, "GetVerificationKey should not return an error on subsequent calls")

	assert.Equal(t, verificationKey1, verificationKey2, "GetVerificationKey should always return the same key")
}

func TestInMemoryRsaKeyPair_GetSigningKey(t *testing.T) {
	keyPair := jwt_util.GetInMemoryRsaKeyPair()

	signingKey1, err1 := keyPair.GetSigningKey()
	assert.Nil(t, err1, "GetSigningKey should not return an error")

	signingKey2, err2 := keyPair.GetSigningKey()
	assert.Nil(t, err2, "GetSigningKey should not return an error on subsequent calls")

	assert.Equal(t, signingKey1, signingKey2, "GetSigningKey should always return the same key")
}
