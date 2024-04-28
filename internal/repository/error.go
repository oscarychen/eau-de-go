package repository

import "fmt"

type DuplicateKeyError struct {
	Key string
}

func (e *DuplicateKeyError) Error() string {
	return e.Key
}

type IncorrectUserCredentialError struct{}

func (e *IncorrectUserCredentialError) Error() string {
	return "Incorrect credentials"
}

type InactiveUserError struct {
	Username string
}

func (e *InactiveUserError) Error() string {
	return fmt.Sprintf("User %s is inactive", e.Username)
}
