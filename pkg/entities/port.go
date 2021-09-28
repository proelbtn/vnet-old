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

func NewPort(name string, network *Network, addrs []*net.IPNet) (*Port, error) {
	err := validateName(name)
	if err != nil {
		return nil, err
	}

	return &Port{
		Name:    name,
		Network: network,
		IPAddrs: addrs,
	}, nil
}

func (v *Port) SetContainer(con *Container) {
	v.Container = con
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
