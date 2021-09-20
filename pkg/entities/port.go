package entities

import (
	"net"

	"github.com/google/uuid"
)

type Port struct {
	ID      uuid.UUID
	Name    string
	Network *Network
	IPAddrs []*Address
}

type Address struct {
	Addr *net.IP
	Net  *net.IPNet
}

func NewPort(name string, network *Network, addrs []*Address) (*Port, error) {
	id, err := uuid.NewRandom()
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
