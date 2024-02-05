package email_util

import "fmt"

type InvalidEmailError struct {
	Key string
}

func (e *InvalidEmailError) Error() string {
	return fmt.Sprintf("Invalid email: %v", e.Key)
}
