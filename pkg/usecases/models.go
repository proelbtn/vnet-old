package usecases

import (
	"github.com/google/uuid"
	"github.com/proelbtn/vnet/pkg/entities"
)

type WritableEnvironment struct {
	Name       string
	Containers []*WritableContainer
	Networks   []*WritableNetwork
}

func (v *WritableEnvironment) ToEntity() (*entities.Environment, error) {
	networks := make([]*entities.Network, len(v.Networks))
	for i := range networks {
		network, err := v.Networks[i].ToEntity()
		if err != nil {
			return nil, err
		}
		networks[i] = network
	}

	containers := make([]*entities.Container, len(v.Containers))
	for i := range networks {
		container, err := v.Containers[i].ToEntity()
		if err != nil {
			return nil, err
		}
		containers[i] = container
	}

	return entities.NewEnvironment(v.Name, containers, networks)
}

type WritableContainer struct {
	Name string
}

func (v *WritableContainer) ToEntity() (*entities.Container, error) {
	return entities.NewContainer(v.Name)
}

type WritableNetwork struct {
	Name string
}

func (v *WritableNetwork) ToEntity() (*entities.Network, error) {
	return entities.NewNetwork(v.Name)
}

type Environment struct {
	ID         uuid.UUID
	Name       string
	Containers []*Container
	Networks   []*Network
}

func NewEnvironment(environment *entities.Environment) *Environment {
	networks := make([]*Network, len(environment.Networks))
	for i := range environment.Networks {
		networks[i] = NewNetwork(environment.Networks[i])
	}

	containers := make([]*Container, len(environment.Containers))
	for i := range environment.Containers {
		containers[i] = NewContainer(environment.Containers[i])
	}

	return &Environment{
		ID:         environment.ID,
		Name:       environment.Name,
		Containers: containers,
		Networks:   networks,
	}
}

type Container struct {
	ID   uuid.UUID
	Name string
}

func NewContainer(container *entities.Container) *Container {
	return &Container{
		ID:   container.ID,
		Name: container.Name,
	}
}

type Network struct {
	ID   uuid.UUID
	Name string
}

func NewNetwork(network *entities.Network) *Network {
	return &Network{
		ID:   network.ID,
		Name: network.Name,
	}
}
