package entities

import (
	"github.com/google/uuid"
)

type Laboratory struct {
	ID         uuid.UUID
	Name       string
	Containers []*Container
	Networks   []*Network
}

func NewLaboratory(name string, containers []*Container, networks []*Network) (*Laboratory, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	lab := &Laboratory{
		ID:   id,
		Name: name,
	}

	for _, container := range containers {
		container.SetLaboratory(lab)
	}

	for _, network := range networks {
		network.SetLaboratory(lab)
	}

	lab.Containers = containers
	lab.Networks = networks

	return lab, nil
}
