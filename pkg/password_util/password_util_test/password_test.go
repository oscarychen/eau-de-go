package password_util_test

import (
	"eau-de-go/pkg/password_util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword_HappyPath(t *testing.T) {
	plainPassword := "StrongPassword123!"
	hashedPassword, err := password_util.HashPassword(plainPassword)
	assert.Nil(t, err)
	assert.NotNil(t, hashedPassword)
}

func TestHashPassword_WeakPassword(t *testing.T) {
	plainPassword := "weak"
	_, err := password_util.HashPassword(plainPassword)
	assert.NotNil(t, err)
}

func TestCheckPassword_HappyPath(t *testing.T) {
	plainPassword := "StrongPassword123!"
	hashedPassword, _ := password_util.HashPassword(plainPassword)
	err := password_util.CheckPassword(plainPassword, hashedPassword)
	assert.Nil(t, err)
}

func TestCheckPassword_WrongPassword(t *testing.T) {
	plainPassword := "StrongPassword123!"
	hashedPassword, _ := password_util.HashPassword(plainPassword)
	err := password_util.CheckPassword("WrongPassword", hashedPassword)
	assert.NotNil(t, err)
}

func TestValidatePassword_HappyPath(t *testing.T) {
	plainPassword := "StrongPassword123!"
	err := password_util.ValidatePassword(plainPassword)
	assert.Nil(t, err)
}

func TestValidatePassword_WeakPassword(t *testing.T) {
	plainPassword := "weak"
	err := password_util.ValidatePassword(plainPassword)
	assert.NotNil(t, err)
}
