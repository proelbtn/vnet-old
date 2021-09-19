package entities

import (
	"github.com/google/uuid"
)

type Container struct {
	ID          uuid.UUID
	Name        string
	Environment *Environment
}

func NewContainer(name string) (*Container, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Container{
		ID:          id,
		Name:        name,
		Environment: nil,
	}, nil
}

func (v *Container) SetEnvironment(env *Environment) {
	v.Environment = env
}
