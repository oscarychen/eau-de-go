package password_util

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plainPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
}

func CheckPassword(plainPassword string, hashedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(plainPassword))
}
