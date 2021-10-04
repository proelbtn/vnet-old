package managers

import (
	"context"

	"github.com/proelbtn/vnet/pkg/entities"
)

type NetworkManager interface {
	Create(ctx context.Context, network *entities.Network) error
	Delete(ctx context.Context, network *entities.Network) error
	CreatePorts(ctx context.Context, pid int, ports []*entities.Port) error
	DeletePorts(ctx context.Context, ports []*entities.Port) error
	GetBridgeName(labName, networkName string) string
	GetPortName(labName, containerName, portName string) string
}
