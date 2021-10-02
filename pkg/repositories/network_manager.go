package repositories

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"

	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/errors"
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

func (v *NetworkManager) Create(ctx context.Context, spec *entities.Network) error {
	return v.create(spec)
}

func (v *NetworkManager) Delete(ctx context.Context, spec *entities.Network) error {
	return v.delete(spec)
}

func (v *NetworkManager) findBridge(name string) (netlink.Link, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		switch err.(type) {
		case netlink.LinkNotFoundError:
		default:
			return nil, err
		}
	}

	if link != nil {
		if link.Type() != "bridge" {
			return nil, errors.ErrInvalidType
		}
		return link, nil
	}

	return nil, errors.ErrNotFound
}

func (v *NetworkManager) ensureBridgeExists(name string) error {
	logger := zap.L().With(zap.String("name", name))

	logger.Debug("ensuring bridge exists")

	logger.Debug("finding bridge")
	bridge, err := v.findBridge(name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if bridge != nil {
		logger.Warn("bridge found")
		return nil
	}

	logger.Debug("bridge not found, creating")

	attrs := netlink.NewLinkAttrs()
	attrs.Name = name
	attrs.Flags = attrs.Flags | net.FlagUp

	return netlink.LinkAdd(&netlink.Bridge{
		LinkAttrs: attrs,
	})
}

func (v *NetworkManager) ensureBridgeNotExist(name string) error {
	logger := zap.L().With(zap.String("name", name))

	logger.Debug("ensuring bridge not exist")

	bridge, err := v.findBridge(name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if bridge == nil {
		logger.Warn("bridge not found")
		return nil
	}

	return netlink.LinkDel(bridge)
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

func (v *NetworkManager) create(spec *entities.Network) error {
	logger := v.getLogger(spec)

	logger.Debug("creating bridge")

	name := GetBridgeName(spec)
	err := v.ensureBridgeExists(name)
	if err != nil {
		return err
	}

	logger.Debug("created bridge")
	return nil
}

func (v *NetworkManager) delete(network *entities.Network) error {
	logger := v.getLogger(network)

	logger.Debug("deleting bridge")

	name := GetBridgeName(network)
	err := v.ensureBridgeNotExist(name)
	if err != nil {
		return err
	}

	logger.Debug("deleted bridge")
	return nil
}
