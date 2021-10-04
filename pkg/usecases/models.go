package usecases

import (
	"net"

	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/errors"
)

type WritableLaboratory struct {
	Name       string
	Containers []*WritableContainer
	Networks   []*WritableNetwork
}

func NewWritableLaboratory(name string, containers []*WritableContainer, networks []*WritableNetwork) *WritableLaboratory {
	return &WritableLaboratory{
		Name:       name,
		Containers: containers,
		Networks:   networks,
	}
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
	for i := range containers {
		container, err := v.Containers[i].ToEntity(networks)
		if err != nil {
			return nil, err
		}
		containers[i] = container
	}

	return entities.NewLaboratory(
		v.Name,
		entities.WithContainers(containers),
		entities.WithNetworks(networks),
	)
}

type WritableContainer struct {
	Name      string
	ImageName string
	Ports     []*WritablePort
	Commands  []string
	Volumes   []*WritableContainerVolume
}

func NewWritableContainer(name string, imageName string, ports []*WritablePort, commands []string, volumes []*WritableContainerVolume) *WritableContainer {
	return &WritableContainer{
		Name:      name,
		ImageName: imageName,
		Ports:     ports,
		Commands:  commands,
		Volumes:   volumes,
	}
}

func (v *WritableContainer) ToEntity(networks []*entities.Network) (*entities.Container, error) {
	ports := make([]*entities.Port, len(v.Ports))
	for i := range v.Ports {
		port, err := v.Ports[i].ToEntity(networks)
		if err != nil {
			return nil, err
		}
		ports[i] = port
	}

	volumes := make([]*entities.ContainerVolume, len(v.Volumes))
	for i := range v.Volumes {
		volumes[i] = v.Volumes[i].ToEntity()
	}

	return entities.NewContainer(v.Name, v.ImageName,
		entities.WithPorts(ports),
		entities.WithCommands(v.Commands),
		entities.WithVolumes(volumes),
	)
}

type WritableContainerVolume struct {
	Source      string
	Destination string
}

func NewWritableContainerVolume(source, destination string) *WritableContainerVolume {
	return &WritableContainerVolume{
		Source:      source,
		Destination: destination,
	}
}

func (v *WritableContainerVolume) ToEntity() *entities.ContainerVolume {
	return &entities.ContainerVolume{
		Source:      v.Source,
		Destination: v.Destination,
	}
}

type WritablePort struct {
	Name      string
	Network   string
	Addresses []*net.IPNet
}

func NewWritablePort(name string, network string, addresses []string) (*WritablePort, error) {
	port := &WritablePort{
		Name:      name,
		Network:   network,
		Addresses: make([]*net.IPNet, len(addresses)),
	}

	for i := range addresses {
		addr, err := entities.NewIPAddress(addresses[i])
		if err != nil {
			return nil, err
		}
		port.Addresses[i] = addr
	}

	return port, nil
}

func (v *WritablePort) ToEntity(networks []*entities.Network) (*entities.Port, error) {
	for _, network := range networks {
		if v.Network == network.Name {
			return entities.NewPort(
				v.Name, network,
				entities.WithIPAddresses(v.Addresses),
			)
		}
	}
	return nil, errors.ErrNotFound
}

type WritableNetwork struct {
	Name string
	Mtu  int
}

func NewWritableNetwork(name string, mtu int) *WritableNetwork {
	return &WritableNetwork{
		Name: name,
		Mtu:  mtu,
	}
}

func (v *WritableNetwork) ToEntity() (*entities.Network, error) {
	return entities.NewNetwork(
		v.Name,
		entities.WithMtu(v.Mtu),
	)
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
