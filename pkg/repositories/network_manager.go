package repositories

import (
	"fmt"
	"net"

	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/usecases/managers"
	"github.com/vishvananda/netlink"
)

type NetworkManager struct {
}

var _ managers.NetworkManager = (*NetworkManager)(nil)

func NewNetworkManger() *NetworkManager {
	return &NetworkManager{}
}

func GetBridgeName(network *entities.Network) string {
	idStr := network.ID.String()
	return fmt.Sprintf("br-%s%s", idStr[:8], idStr[9:13])
}

func (v *NetworkManager) Create(network *entities.Network) error {
	return v.createBridge(network)
}

func (v *NetworkManager) Delete(network *entities.Network) error {
	return v.deleteBridge(network)
}

func (v *NetworkManager) createBridge(network *entities.Network) error {
	attrs := netlink.NewLinkAttrs()
	attrs.Name = GetBridgeName(network)
	attrs.Flags = attrs.Flags | net.FlagUp

	bridge := &netlink.Bridge{
		LinkAttrs: attrs,
	}

	return netlink.LinkAdd(bridge)
}

func (v *NetworkManager) deleteBridge(network *entities.Network) error {
	name := GetBridgeName(network)
	bridge, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	// TODO: add warning if the type of `bridge` is not "bridge"

	return netlink.LinkDel(bridge)
}
