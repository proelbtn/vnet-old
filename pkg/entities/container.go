package entities

import (
	"github.com/google/uuid"
)

type Container struct {
<<<<<<< HEAD
	ID         uuid.UUID
	Name       string
	Laboratory *Laboratory

	ImageName string
=======
	ID          uuid.UUID
	Name        string
	ImageName   string
	Environment *Environment
>>>>>>> 6ad4b19 (add ImageName to Container)
}

func NewContainer(name string, imageName string) (*Container, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Container{
<<<<<<< HEAD
		ID:         id,
		Name:       name,
		Laboratory: nil,
		ImageName:  imageName,
=======
		ID:          id,
		Name:        name,
		ImageName:   imageName,
		Environment: nil,
>>>>>>> 6ad4b19 (add ImageName to Container)
	}, nil
}

func (v *Container) SetLaboratory(env *Laboratory) {
	v.Laboratory = env
}
