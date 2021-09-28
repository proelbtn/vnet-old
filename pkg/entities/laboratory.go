package entities

type Laboratory struct {
	Name       string
	Containers []*Container
	Networks   []*Network
}

func NewLaboratory(name string, containers []*Container, networks []*Network) (*Laboratory, error) {
	err := validateName(name)
	if err != nil {
		return nil, err
	}

	lab := &Laboratory{
		Name: name,
	}

	for _, container := range containers {
		container.SetLaboratory(lab)
	}

	for _, network := range networks {
		network.SetLaboratory(lab)
	}

	lab.Containers = containers
	lab.Networks = networks

	return lab, nil
}

func (v *Laboratory) GetUniqueName() string {
	return v.Name
}
