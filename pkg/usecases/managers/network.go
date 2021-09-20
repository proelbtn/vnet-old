package managers

import "github.com/proelbtn/vnet/pkg/entities"

type NetworkManager interface {
	Create(network *entities.Network) error
	Delete(network *entities.Network) error
	AttachPorts(pid uint32, ports []*entities.Port) error
}
