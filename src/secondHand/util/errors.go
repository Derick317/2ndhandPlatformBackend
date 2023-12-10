package util

import (
	"errors"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailExists   = errors.New("email has already exists")
	ErrOrderNoExists = errors.New("order does not exist")
	ErrItemNotFound  = errors.New("item not found")
	ErrGCS           = errors.New("errors from GCS")
	ErrUnexpected    = errors.New("unexpected error because of bugs")
)
