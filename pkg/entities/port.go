package entities

import (
	"net"

	"github.com/google/uuid"
)

type Port struct {
	ID      uuid.UUID
	Name    string
	Network *Network
	IPAddrs []*net.IPNet
}

func NewPort(name string, network *Network, addrs []*net.IPNet) (*Port, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	err = validateName(name)
	if err != nil {
		return nil, err
	}

	return &Port{
		ID:      id,
		Name:    name,
		Network: network,
		IPAddrs: addrs,
	}, nil
}

func NewIPAddress(cidr string) (*net.IPNet, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	ipNet.IP = ip
	return ipNet, nil
}
