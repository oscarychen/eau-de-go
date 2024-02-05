package password_util

import "fmt"

type WeakPasswordError struct {
	Key string
}

func (e *WeakPasswordError) Error() string {
	return fmt.Sprintf("Password is weak: %v", e.Key)
}
