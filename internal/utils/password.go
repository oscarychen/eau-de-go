package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plainPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
}

func ComparePassword(plainPassword string, hashedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(plainPassword))
}
