package entities

import (
	"github.com/google/uuid"
)

type Environment struct {
	ID         uuid.UUID
	Name       string
	Containers []Container
	Networks   []Network
}

func NewEnvironment(name string, containers []Container, networks []Network) (*Environment, error) {
	// TODO: validate containers and networks
	return &Environment{
		Name:       name,
		Containers: containers,
		Networks:   networks,
	}, nil
}
