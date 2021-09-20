package managers

import "github.com/proelbtn/vnet/pkg/entities"

type ContainerManager interface {
	Create(container *entities.Container) error
	Start(container *entities.Container) error
	Stop(container *entities.Container) error
	Delete(container *entities.Container) error
}
