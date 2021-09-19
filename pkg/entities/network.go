package entities

import (
	"github.com/google/uuid"
)

type Network struct {
	ID   uuid.UUID
	Name string
}

func NewNetwork(name string) (*Network, error) {
	return &Network{
		Name: name,
	}, nil
}
