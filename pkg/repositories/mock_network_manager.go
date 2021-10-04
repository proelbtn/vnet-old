package repositories

import (
	"context"

	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/usecases/managers"
)

type MockNetworkManager struct{}

var _ managers.NetworkManager = (*MockNetworkManager)(nil)

func NewMockNetworkManger() *MockNetworkManager {
	return &MockNetworkManager{}
}

func (v *MockNetworkManager) Create(ctx context.Context, spec *entities.Network) error {
	return nil
}

func (v *MockNetworkManager) Delete(ctx context.Context, spec *entities.Network) error {
	return nil
}

func (v *MockNetworkManager) CreatePorts(ctx context.Context, pid int, ports []*entities.Port) error {
	return nil
}

func (v *MockNetworkManager) DeletePorts(ctx context.Context, ports []*entities.Port) error {
	return nil
}

func (v *MockNetworkManager) GetBridgeName(labName, networkName string) string {
	return ""
}

func (v *MockNetworkManager) GetPortName(labName, containerName, portName string) string {
	return ""
}
