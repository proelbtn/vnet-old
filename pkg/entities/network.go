package entities

import "fmt"

type Network struct {
	Name       string
	Mtu        int
	Laboratory *Laboratory
}

type NewNetworkOpts func(*Network) error

func NewNetwork(name string, options ...NewNetworkOpts) (*Network, error) {
	network := &Network{
		Name: name,
		Mtu:  1500,
	}

	for _, option := range options {
		if err := option(network); err != nil {
			return nil, err
		}
	}

	err := validateName(name)
	if err != nil {
		return nil, err
	}

	return network, nil
}

func WithJumboframe(network *Network) error {
	network.Mtu = 9000
	return nil
}

func WithMtu(mtu int) NewNetworkOpts {
	return func(network *Network) error {
		network.Mtu = mtu
		return nil
	}
}

func (v *Network) setLaboratory(env *Laboratory) {
	v.Laboratory = env
}

func (v *Network) GetUniqueName() string {
	return fmt.Sprintf("%s/%s", v.Laboratory.GetUniqueName(), v.Name)
}
