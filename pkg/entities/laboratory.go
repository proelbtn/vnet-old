package entities

type Laboratory struct {
	Name       string
	Containers []*Container
	Networks   []*Network
}

type NewLaboratoryOpts func(*Laboratory) error

func NewLaboratory(name string, options ...NewLaboratoryOpts) (*Laboratory, error) {
	laboratory := &Laboratory{
		Name: name,
	}

	for _, option := range options {
		if err := option(laboratory); err != nil {
			return nil, err
		}
	}

	err := validateName(name)
	if err != nil {
		return nil, err
	}

	return laboratory, nil
}

func WithContainer(container *Container) NewLaboratoryOpts {
	return func(laboratory *Laboratory) error {
		container.setLaboratory(laboratory)
		laboratory.Containers = append(laboratory.Containers, container)
		return nil
	}
}

func WithContainers(containers []*Container) NewLaboratoryOpts {
	return func(laboratory *Laboratory) error {
		for _, container := range containers {
			if err := WithContainer(container)(laboratory); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithNetwork(network *Network) NewLaboratoryOpts {
	return func(laboratory *Laboratory) error {
		network.setLaboratory(laboratory)
		laboratory.Networks = append(laboratory.Networks, network)
		return nil
	}
}

func WithNetworks(networks []*Network) NewLaboratoryOpts {
	return func(laboratory *Laboratory) error {
		for _, network := range networks {
			if err := WithNetwork(network)(laboratory); err != nil {
				return err
			}
		}
		return nil
	}
}

func (v *Laboratory) GetUniqueName() string {
	return v.Name
}
