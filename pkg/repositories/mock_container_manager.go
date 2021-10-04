package repositories

import (
	"context"

	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/usecases/managers"
)

type MockContainerManager struct{}

var _ managers.ContainerManager = (*MockContainerManager)(nil)

func NewMockContainerManager() *MockContainerManager {
	return &MockContainerManager{}
}

func (v *MockContainerManager) Create(ctx context.Context, spec *entities.Container) (uint32, error) {
	return 0, nil
}

func (v *MockContainerManager) Start(ctx context.Context, spec *entities.Container) error {
	return nil
}

func (v *MockContainerManager) Stop(ctx context.Context, spec *entities.Container) error {
	return nil
}

func (v *MockContainerManager) Delete(ctx context.Context, spec *entities.Container) error {
	return nil
}

func (v *MockContainerManager) Exec(ctx context.Context, spec *entities.Container, args managers.ExecArgs) error {
	return nil
}
