package repositories

import (
	"fmt"
	"net"

	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/usecases/managers"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
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

func GetPortName(port *entities.Port) string {
	idStr := port.ID.String()
	return fmt.Sprintf("po-%s%s", idStr[:8], idStr[9:13])
}

func (v *NetworkManager) Create(network *entities.Network) error {
	return v.createBridge(network)
}

func (v *NetworkManager) Delete(network *entities.Network) error {
	return v.deleteBridge(network)
}

func (v *NetworkManager) AttachPorts(pid int, ports []*entities.Port) error {
	nsHandle, err := netns.GetFromPid(pid)
	if err != nil {
		return err
	}

	defaultNsHandler, err := netlink.NewHandle()
	if err != nil {
		return err
	}

	containerNsHandler, err := netlink.NewHandleAt(nsHandle)
	if err != nil {
		return err
	}

	for _, port := range ports {
		attrs := netlink.NewLinkAttrs()
		attrs.Name = GetPortName(port)
		attrs.Flags = attrs.Flags | net.FlagUp

		veth := &netlink.Veth{
			LinkAttrs:     attrs,
			PeerName:      port.Name,
			PeerNamespace: netlink.NsPid(pid),
		}

		err := defaultNsHandler.LinkAdd(veth)
		if err != nil {
			return err
		}

		bridge, err := defaultNsHandler.LinkByName(GetBridgeName(port.Network))
		if err != nil {
			return err
		}

		err = defaultNsHandler.LinkSetMaster(veth, bridge)
		if err != nil {
			return err
		}

		peer, err := containerNsHandler.LinkByName(port.Name)
		if err != nil {
			return err
		}

		err = containerNsHandler.LinkSetUp(peer)
		if err != nil {
			return err
		}
	}

	return nil
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
