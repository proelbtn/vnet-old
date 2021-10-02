package commands

import (
	"io/ioutil"
	"path/filepath"

	"github.com/proelbtn/vnet/pkg/repositories"
	"github.com/proelbtn/vnet/pkg/usecases"
	"gopkg.in/yaml.v3"
)

// TODO: refactoring
var usecase *usecases.InstantLaboratoryUsecase = nil

func getUsecase() (*usecases.InstantLaboratoryUsecase, error) {
	if usecase != nil {
		return usecase, nil
	}

	networkManager := repositories.NewNetworkManger()
	containerManager, err := repositories.NewContainerManager()
	if err != nil {
		return nil, err
	}

	laboratoryManager := repositories.NewLaboratoryManager(containerManager, networkManager)
	usecase := usecases.NewInstantLaboratoryUsecase(laboratoryManager, containerManager, networkManager)

	return usecase, nil
}

type Laboratory struct {
	Name       string       `yaml:"name"`
	Containers []*Container `yaml:"containers"`
	Networks   []*Network   `yaml:"networks"`
}

func (v *Laboratory) ToWritableLaboratory() (*usecases.WritableLaboratory, error) {
	networks := make([]*usecases.WritableNetwork, len(v.Networks))
	for i := range networks {
		networks[i] = v.Networks[i].ToWritableNetwork()
	}

	containers := make([]*usecases.WritableContainer, len(v.Containers))
	for i := range containers {
		container, err := v.Containers[i].ToWritableContainer()
		if err != nil {
			return nil, err
		}
		containers[i] = container
	}

	return usecases.NewWritableLaboratory(v.Name, containers, networks), nil
}

type Container struct {
	Name      string             `yaml:"name"`
	ImageName string             `yaml:"image"`
	Ports     []*Port            `yaml:"ports"`
	Commands  []string           `yaml:"commands"`
	Volumes   []*ContainerVolume `yaml:"volumes"`
}

func (v *Container) ToWritableContainer() (*usecases.WritableContainer, error) {
	ports := make([]*usecases.WritablePort, len(v.Ports))
	for i := range v.Ports {
		port, err := v.Ports[i].ToWritablePort()
		if err != nil {
			return nil, err
		}
		ports[i] = port
	}

	volumes := make([]*usecases.WritableContainerVolume, len(v.Volumes))
	for i := range v.Volumes {
		volume, err := v.Volumes[i].ToWritableContainerVolume()
		if err != nil {
			return nil, err
		}
		volumes[i] = volume
	}

	return usecases.NewWritableContainer(v.Name, v.ImageName, ports, v.Commands, volumes), nil
}

type Port struct {
	Name      string   `yaml:"name"`
	Network   string   `yaml:"network"`
	Addresses []string `yaml:"addresses"`
}

func (v *Port) ToWritablePort() (*usecases.WritablePort, error) {
	return usecases.NewWritablePort(v.Name, v.Network, v.Addresses)
}

type ContainerVolume struct {
	Source      string
	Destination string
}

func (v *ContainerVolume) ToWritableContainerVolume() (*usecases.WritableContainerVolume, error) {
	source, err := filepath.Abs(v.Source)
	if err != nil {
		return nil, err
	}

	return usecases.NewWritableContainerVolume(source, v.Destination), nil
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
