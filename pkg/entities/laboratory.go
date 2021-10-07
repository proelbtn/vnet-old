package entities

type Laboratory struct {
	Name          string
	PreRequisites LaboratoryPreRequisites
	Containers    []*Container
	Networks      []*Network
}

type LaboratoryPreRequisites struct {
	KernelVersion string
	Modules       []string
	Configs       []string
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

func WithRequiredKernelVersion(version string) NewLaboratoryOpts {
	return func(laboratory *Laboratory) error {
		if laboratory.PreRequisites.KernelVersion != "" {
			return nil
		}
		laboratory.PreRequisites.KernelVersion = version
		return nil
	}
}

func WithRequiredKernelModule(module string) NewLaboratoryOpts {
	return func(laboratory *Laboratory) error {
		laboratory.PreRequisites.Modules = append(laboratory.PreRequisites.Modules, module)
		return nil
	}
}

func WithRequiredKernelModules(modules []string) NewLaboratoryOpts {
	return func(laboratory *Laboratory) error {
		for _, module := range modules {
			if err := WithRequiredKernelModule(module)(laboratory); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithRequiredKernelConfig(config string) NewLaboratoryOpts {
	return func(laboratory *Laboratory) error {
		laboratory.PreRequisites.Configs = append(laboratory.PreRequisites.Configs, config)
		return nil
	}
}

func WithRequiredKernelConfigs(configs []string) NewLaboratoryOpts {
	return func(laboratory *Laboratory) error {
		for _, config := range configs {
			if err := WithRequiredKernelConfig(config)(laboratory); err != nil {
				return err
			}
		}
		return nil
	}
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
