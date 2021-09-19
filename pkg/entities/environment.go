package entities

import (
	"github.com/google/uuid"
)

type Environment struct {
	ID         uuid.UUID
	Name       string
	Containers []*Container
	Networks   []*Network
}

func NewEnvironment(name string, containers []*Container, networks []*Network) (*Environment, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	env := &Environment{
		ID:   id,
		Name: name,
	}

	for _, container := range containers {
		container.SetEnvironment(env)
	}

	for _, network := range networks {
		network.SetEnvironment(env)
	}

	env.Containers = containers
	env.Networks = networks

	return env, nil
}
