package repositories

import (
	"context"
	"fmt"
	"net"

	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/usecases/managers"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"go.uber.org/zap"
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

func (v *NetworkManager) getLogger(network *entities.Network) *zap.Logger {
	return zap.L().With(
		zap.String("ID", network.ID.String()),
		zap.String("Name", network.Name),
		zap.String("Laboratory.ID", network.Laboratory.ID.String()),
		zap.String("Laboratory.Name", network.Laboratory.Name),
	)
}

func (v *NetworkManager) Create(ctx context.Context, network *entities.Network) error {
	return v.createBridge(network)
}

func (v *NetworkManager) Delete(ctx context.Context, network *entities.Network) error {
	return v.deleteBridge(network)
}

func (v *NetworkManager) AttachPorts(ctx context.Context, pid int, ports []*entities.Port) error {
	logger := zap.L()

	logger.Debug("creating bridge")

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

	logger.Debug("creating ports")
	for _, port := range ports {
		logger.Debug("creating port", zap.String("ID", port.ID.String()), zap.String("Name", port.Name))

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
	logger := v.getLogger(network)

	logger.Debug("creating bridge")

	attrs := netlink.NewLinkAttrs()
	attrs.Name = GetBridgeName(network)
	attrs.Flags = attrs.Flags | net.FlagUp

	bridge := &netlink.Bridge{
		LinkAttrs: attrs,
	}

	err := netlink.LinkAdd(bridge)
	if err != nil {
		return err
	}

	logger.Debug("created bridge")
	return nil
}

func (v *NetworkManager) deleteBridge(network *entities.Network) error {
	logger := v.getLogger(network)

	logger.Debug("deleting bridge")

	name := GetBridgeName(network)
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	if link.Type() != "bridge" {
		logger.Warn("selected link isn't bridge", zap.Any("link", link))
	}

	err = netlink.LinkDel(link)
	if err != nil {
		return err
	}

	logger.Debug("deleted bridge")
	return nil
}
