package entities

import (
	"github.com/google/uuid"
)

type Network struct {
	ID          uuid.UUID
	Name        string
	Environment *Environment
}

func NewNetwork(name string) (*Network, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Network{
		ID:          id,
		Name:        name,
		Environment: nil,
	}, nil
}

func (v *Network) SetEnvironment(env *Environment) {
	v.Environment = env
}
