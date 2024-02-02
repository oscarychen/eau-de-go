package repository

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
