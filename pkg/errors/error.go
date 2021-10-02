package errors

import (
	"errors"
	"fmt"
)

var (
	ErrAlreadyExists = fmt.Errorf("already exists")
	ErrNotFound      = fmt.Errorf("not found")
	ErrInvalidName   = fmt.Errorf("invalid name")
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func IsAlreadyExists(err error) bool {
	return Is(err, ErrNotFound)
}

func IsNotFound(err error) bool {
	return Is(err, ErrNotFound)
}

func IsInvalidName(err error) bool {
	return Is(err, ErrNotFound)
}
