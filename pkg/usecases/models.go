package usecases

import (
	"github.com/proelbtn/vnet/pkg/entities"
)

type WritableLaboratory struct {
	Name       string
	Containers []*WritableContainer
	Networks   []*WritableNetwork
}

func (v *WritableLaboratory) ToEntity() (*entities.Laboratory, error) {
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

	return entities.NewLaboratory(v.Name, containers, networks)
}

type WritableContainer struct {
	Name      string
	ImageName string
}

func (v *WritableContainer) ToEntity() (*entities.Container, error) {
	return entities.NewContainer(v.Name, v.ImageName, nil)
}

type WritableNetwork struct {
	Name string
}

func (v *WritableNetwork) ToEntity() (*entities.Network, error) {
	return entities.NewNetwork(v.Name)
}

type Laboratory struct {
	Name       string
	Containers []*Container
	Networks   []*Network
}

func NewLaboratory(laboratory *entities.Laboratory) *Laboratory {
	networks := make([]*Network, len(laboratory.Networks))
	for i := range laboratory.Networks {
		networks[i] = NewNetwork(laboratory.Networks[i])
	}

	containers := make([]*Container, len(laboratory.Containers))
	for i := range laboratory.Containers {
		containers[i] = NewContainer(laboratory.Containers[i])
	}

	return &Laboratory{
		Name:       laboratory.Name,
		Containers: containers,
		Networks:   networks,
	}
}

type Container struct {
	Name      string
	ImageName string
}

func NewContainer(container *entities.Container) *Container {
	return &Container{
		Name:      container.Name,
		ImageName: container.ImageName,
	}
}

type Network struct {
	Name string
}

func NewNetwork(network *entities.Network) *Network {
	return &Network{
		Name: network.Name,
	}
}
