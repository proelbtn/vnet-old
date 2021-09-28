package main

import (
	"github.com/proelbtn/vnet/pkg/usecases"
)

type Laboratory struct {
	Name       string       `yaml:"name"`
	Containers []*Container `yaml:"containers"`
	Networks   []*Network   `yaml:"networks"`
}

func (v *Laboratory) ToWritableLaboratory() *usecases.WritableLaboratory {
	networks := make([]*usecases.WritableNetwork, len(v.Networks))
	for i := range networks {
		networks[i] = v.Networks[i].ToWritableNetwork()
	}

	containers := make([]*usecases.WritableContainer, len(v.Containers))
	for i := range containers {
		containers[i] = v.Containers[i].ToWritableContainer()
	}

	return usecases.NewWritableLaboratory(v.Name, containers, networks)
}

type Container struct {
	Name      string `yaml:"name"`
	ImageName string `yaml:"image"`
}

func (v *Container) ToWritableContainer() *usecases.WritableContainer {
	return usecases.NewWritableContainer(v.Name, v.ImageName)
}

type Network struct {
	Name string `yaml:"name"`
}

func (v *Network) ToWritableNetwork() *usecases.WritableNetwork {
	return usecases.NewWritableNetwork(v.Name)
}
