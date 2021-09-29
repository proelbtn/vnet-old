package errors

import "fmt"

var (
	ErrAlreadyExists = fmt.Errorf("already exists")
	ErrNotFound      = fmt.Errorf("not found")
	ErrInvalidName   = fmt.Errorf("invalid name")
)
