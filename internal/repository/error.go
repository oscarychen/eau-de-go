package repository

type DuplicateKeyError struct {
	Key string
}

func (e *DuplicateKeyError) Error() string {
	return e.Key
}
