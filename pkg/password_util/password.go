package password_util

import (
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plainPassword string) ([]byte, error) {
	err := ValidatePassword(plainPassword)
	if err != nil {
		return nil, &WeakPasswordError{Key: err.Error()}
	}
	return bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
}

func CheckPassword(plainPassword string, hashedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(plainPassword))
}

func ValidatePassword(plainPassword string) error {
	const minEntropy = 60
	return passwordvalidator.Validate(plainPassword, minEntropy)
}
