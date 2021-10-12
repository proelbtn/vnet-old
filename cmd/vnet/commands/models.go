package commands

import (
	"errors"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/proelbtn/vnet/pkg/usecases"
)

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
	Name       string   `yaml:"name"`
	Network    string   `yaml:"network"`
	MacAddress string   `yaml:"mac"`
	Addresses  []string `yaml:"addresses"`
}

func (v *Port) parseMacAddress() ([]byte, error) {
	matched, err := regexp.MatchString("[[:xdigit:]]{2}(:[[:xdigit:]]{2}){5}", v.MacAddress)
	if err != nil {
		return nil, err
	}

	if !matched {
		return nil, errors.New("")
	}

	addr := []byte{}
	for _, p := range strings.Split(v.MacAddress, ":") {
		v, err := strconv.ParseUint(p, 16, 8)
		if err != nil {
			return nil, err
		}

		addr = append(addr, byte(v))
	}

	return addr, nil
}

func (v *Port) ToWritablePort() (*usecases.WritablePort, error) {
	var mac []byte = nil
	if v.MacAddress != "" {
		addr, err := v.parseMacAddress()
		if err != nil {
			return nil, err
		}
		mac = addr
	}

	return usecases.NewWritablePort(v.Name, v.Network, mac, v.Addresses)
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
	Mtu  int    `yaml:"mtu"`
}

func (v *Network) ToWritableNetwork() *usecases.WritableNetwork {
	mtu := 1500
	if v.Mtu != 0 {
		mtu = v.Mtu
	}

	return usecases.NewWritableNetwork(v.Name, mtu)
}
