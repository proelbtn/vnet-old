package entities

import (
	"fmt"
)

type Container struct {
	Name      string
	ImageName string

	Laboratory *Laboratory

	EnvironmentVariables map[string]string
	Ports                []*Port
	Commands             []string
	Volumes              []*ContainerVolume
}

type ContainerVolume struct {
	Source      string
	Destination string
}

type NewContainerOpts func(*Container) error

func NewContainer(name string, imageName string, options ...NewContainerOpts) (*Container, error) {
	container := &Container{
		Name:      name,
		ImageName: imageName,
	}

	for _, option := range options {
		if err := option(container); err != nil {
			return nil, err
		}
	}

	err := validateName(name)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func WithEnvironmentVariable(key, value string) NewContainerOpts {
	return func(container *Container) error {
		container.EnvironmentVariables[key] = value
		return nil
	}
}

func WithEnvironmentVariables(variables map[string]string) NewContainerOpts {
	return func(container *Container) error {
		for key, value := range variables {
			if err := WithEnvironmentVariable(key, value)(container); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithPort(port *Port) NewContainerOpts {
	return func(container *Container) error {
		container.Ports = append(container.Ports, port)
		return nil
	}
}

func WithPorts(ports []*Port) NewContainerOpts {
	return func(container *Container) error {
		for _, port := range ports {
			port.SetContainer(container)
			if err := WithPort(port)(container); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithCommand(cmd string) NewContainerOpts {
	return func(container *Container) error {
		container.Commands = append(container.Commands, cmd)
		return nil
	}
}

func WithCommands(cmds []string) NewContainerOpts {
	return func(container *Container) error {
		for _, cmd := range cmds {
			if err := WithCommand(cmd)(container); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithVolume(volume *ContainerVolume) NewContainerOpts {
	return func(container *Container) error {
		container.Volumes = append(container.Volumes, volume)
		return nil
	}
}

func WithVolumes(volumes []*ContainerVolume) NewContainerOpts {
	return func(container *Container) error {
		for _, volume := range volumes {
			if err := WithVolume(volume)(container); err != nil {
				return err
			}
		}
		return nil
	}
}

func (v *Container) setLaboratory(env *Laboratory) {
	v.Laboratory = env
}

func (v *Container) GetUniqueName() string {
	return fmt.Sprintf("%s/%s", v.Laboratory.GetUniqueName(), v.Name)
}
