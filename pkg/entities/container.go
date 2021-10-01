package entities

import "fmt"

type Container struct {
	Name       string
	Laboratory *Laboratory

	ImageName            string
	EnvironmentVariables map[string]string
	Ports                []*Port
	Commands             []string
	Volumes              []*ContainerVolume
}

type ContainerVolume struct {
	Source      string
	Destination string
}

func NewContainer(name string, imageName string, ports []*Port, commands []string, volumes []*ContainerVolume) (*Container, error) {
	err := validateName(name)
	if err != nil {
		return nil, err
	}

	con := &Container{
		Name:                 name,
		Laboratory:           nil,
		ImageName:            imageName,
		EnvironmentVariables: make(map[string]string),
		Ports:                make([]*Port, 0),
		Commands:             commands,
		Volumes:              volumes,
	}

	for _, port := range ports {
		port.Container = con
	}
	con.Ports = ports

	return con, nil
}

func (v *Container) SetLaboratory(env *Laboratory) {
	v.Laboratory = env
}

func (v *Container) GetUniqueName() string {
	return fmt.Sprintf("%s/%s", v.Laboratory.GetUniqueName(), v.Name)
}
