package repositories

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
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

func hashForNetworkManager(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetBridgeName(network *entities.Network) string {
	idStr := hashForNetworkManager(network.GetUniqueName())
	return fmt.Sprintf("br-%s", idStr[:12])
}

func GetPortName(port *entities.Port) string {
	idStr := hashForNetworkManager(port.GetUniqueName())
	return fmt.Sprintf("po-%s", idStr[:12])
}

func (v *NetworkManager) getLogger(network *entities.Network) *zap.Logger {
	return zap.L().With(
		zap.String("Name", network.Name),
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
		logger.Debug("creating port", zap.String("Name", port.Name))

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

		for _, addr := range port.IPAddrs {
			containerNsHandler.AddrAdd(peer, &netlink.Addr{
				IPNet: addr,
			})
		}
	}

	return nil
}

func (v *NetworkManager) createBridge(network *entities.Network) error {
	logger := v.getLogger(network)

	logger.Debug("creating bridge")

	logger.Debug("checking bridge is already created")
	link, err := netlink.LinkByName(GetBridgeName(network))
	if err != nil {
		switch err.(type) {
		case netlink.LinkNotFoundError:
			// link may not exist, skipped
		default:
			return err
		}
	}

	if link != nil {
		if link.Type() != "bridge" {
			return errors.New("link found, but not brdige")
		}
		logger.Debug("bridge is already created, skipped")
		return nil
	}

	logger.Debug("bridge is not already created, creating")

	attrs := netlink.NewLinkAttrs()
	attrs.Name = GetBridgeName(network)
	attrs.Flags = attrs.Flags | net.FlagUp

	bridge := &netlink.Bridge{
		LinkAttrs: attrs,
	}

	err = netlink.LinkAdd(bridge)
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
		switch err.(type) {
		case netlink.LinkNotFoundError:
			logger.Warn("link not found")
			return nil
		default:
			return err
		}
	}

	if link.Type() != "bridge" {
		return fmt.Errorf("link found, but not bridge")
	}

	err = netlink.LinkDel(link)
	if err != nil {
		return err
	}

	logger.Debug("deleted bridge")
	return nil
}
