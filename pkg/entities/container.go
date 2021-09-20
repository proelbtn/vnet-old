package entities

import (
	"github.com/google/uuid"
)

type Container struct {
	ID         uuid.UUID
	Name       string
	Laboratory *Laboratory

	ImageName string
}

func NewContainer(name string, imageName string) (*Container, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Container{
		ID:         id,
		Name:       name,
		Laboratory: nil,
		ImageName:  imageName,
	}, nil
}

func (v *Container) SetLaboratory(env *Laboratory) {
	v.Laboratory = env
}
