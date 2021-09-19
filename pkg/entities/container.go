package entities

import (
	"github.com/google/uuid"
)

type Container struct {
	ID   uuid.UUID
	Name string
}

func NewContainer(name string) (*Container, error) {
	return &Container{
		Name: name,
	}, nil
}
