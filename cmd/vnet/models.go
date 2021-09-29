package main

import (
	"io/ioutil"

	"github.com/proelbtn/vnet/pkg/usecases"
	"gopkg.in/yaml.v3"
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

func loadManifest(manifestPath string) (*Laboratory, error) {
	var lab Laboratory

	manifest, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(manifest, &lab)
	if err != nil {
		return nil, err
	}

	return &lab, nil
}
