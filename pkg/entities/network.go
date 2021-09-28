package entities

import "fmt"

type Network struct {
	Name       string
	Laboratory *Laboratory
}

func NewNetwork(name string) (*Network, error) {
	err := validateName(name)
	if err != nil {
		return nil, err
	}

	return &Network{
		Name:       name,
		Laboratory: nil,
	}, nil
}

func (v *Network) SetLaboratory(env *Laboratory) {
	v.Laboratory = env
}

func (v *Network) GetUniqueName() string {
	return fmt.Sprintf("%s/%s", v.Laboratory.GetUniqueName(), v.Name)
}
