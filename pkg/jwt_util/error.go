package jwt_util

import "fmt"

type InvalidTokenError struct {
	token string
	msg   *string
}

func (e *InvalidTokenError) Error() string {
	if e.msg != nil {
		return fmt.Sprintf("Invalid token: %s", *e.msg)
	}
	return "Invalid token."
}
