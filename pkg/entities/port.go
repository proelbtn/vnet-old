package entities

import (
	"fmt"
	"net"
)

type Port struct {
	Name      string
	Network   *Network
	Container *Container
	IPAddrs   []*net.IPNet
}

type NewPortOpts func(*Port) error

func NewPort(name string, network *Network, options ...NewPortOpts) (*Port, error) {
	port := &Port{
		Name:    name,
		Network: network,
	}

	for _, option := range options {
		if err := option(port); err != nil {
			return nil, err
		}
	}

	err := validateName(name)
	if err != nil {
		return nil, err
	}

	return port, nil
}

func WithIPAddress(addr *net.IPNet) NewPortOpts {
	return func(port *Port) error {
		port.IPAddrs = append(port.IPAddrs, addr)
		return nil
	}
}

func WithIPAddresses(addrs []*net.IPNet) NewPortOpts {
	return func(port *Port) error {
		for _, addr := range addrs {
			if err := WithIPAddress(addr)(port); err != nil {
				return err
			}
		}
		return nil
	}
}

func (v *Port) SetContainer(container *Container) {
	v.Container = container
}

func NewIPAddress(cidr string) (*net.IPNet, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	ipNet.IP = ip
	return ipNet, nil
}

func (v *Port) GetUniqueName() string {
	return fmt.Sprintf("%s/%s", v.Container.GetUniqueName(), v.Name)
}
