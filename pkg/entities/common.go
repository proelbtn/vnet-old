package entities

import (
	"regexp"

	"github.com/proelbtn/vnet/pkg/errors"
)

func validateName(name string) error {
	matched, err := regexp.Match("^[a-zA-Z0-9][a-zA-Z0-9_-]{1,32}$", []byte(name))
	if err != nil {
		return err
	}
	if !matched {
		return errors.ErrInvalidName
	}
	return nil
}
