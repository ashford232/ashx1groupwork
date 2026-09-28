package errork

import (
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("Invalid Credentials")
	ErrUserExists = errors.New("User already exist, try signing in.")
)
