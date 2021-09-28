package entities

import (
	"github.com/google/uuid"
)

type Network struct {
	ID         uuid.UUID
	Name       string
	Laboratory *Laboratory
}

func NewNetwork(name string) (*Network, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	err = validateName(name)
	if err != nil {
		return nil, err
	}

	return &Network{
		ID:         id,
		Name:       name,
		Laboratory: nil,
	}, nil
}

func (v *Network) SetLaboratory(env *Laboratory) {
	v.Laboratory = env
}
