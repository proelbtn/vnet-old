package managers

import (
	"context"

	"github.com/proelbtn/vnet/pkg/entities"
)

type ContainerManager interface {
	Create(ctx context.Context, container *entities.Container) (uint32, error)
	Start(ctx context.Context, container *entities.Container) error
	Stop(ctx context.Context, container *entities.Container) error
	Delete(ctx context.Context, container *entities.Container) error
}
