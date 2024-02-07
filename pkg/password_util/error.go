package password_util

import "fmt"

type WeakPasswordError struct {
	Key string
}

func (e *WeakPasswordError) Error() string {
	return fmt.Sprintf("Password is weak: %v", e.Key)
}

type SamePasswordError struct{}

func (e *SamePasswordError) Error() string {
	return "Old and new password are the same"
}
