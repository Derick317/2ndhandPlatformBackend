package util

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailExists   = errors.New("email has already exists")
	ErrOrderNoExists = errors.New("order does not exist")
	ErrItemNotFound  = errors.New("item not found")
)

func ErrUnexpected(format string, a ...any) error {
	return fmt.Errorf("unexpected error because of bugs: %s", fmt.Sprintf(format, a...))
}
