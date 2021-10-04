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

func GetBridgeName(labName, networkName string) string {
	uniqueName := fmt.Sprintf("%s/%s", labName, networkName)
	idStr := hashForNetworkManager(uniqueName)
	return fmt.Sprintf("br-%s", idStr[:12])
}

func GetPortName(labName, containerName, portName string) string {
	uniqueName := fmt.Sprintf("%s/%s/%s", labName, containerName, portName)
	idStr := hashForNetworkManager(uniqueName)
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

func (v *NetworkManager) CreatePorts(ctx context.Context, pid int, ports []*entities.Port) error {
	return v.createPorts(ctx, pid, ports)
}

func (v *NetworkManager) DeletePorts(ctx context.Context, ports []*entities.Port) error {
	return v.deletePorts(ctx, ports)
}

func (v *NetworkManager) GetBridgeName(labName, networkName string) string {
	return GetBridgeName(labName, networkName)
}

func (v *NetworkManager) GetPortName(labName, containerName, portName string) string {
	return GetPortName(labName, containerName, portName)
}

func (v *NetworkManager) getBridgeName(network *entities.Network) string {
	return v.GetBridgeName(network.Laboratory.Name, network.Name)
}

func (v *NetworkManager) getPortName(port *entities.Port) string {
	return v.GetPortName(port.Container.Laboratory.Name, port.Container.Name, port.Name)
}

func (v *NetworkManager) findLinkWithHandler(h *netlink.Handle, name string) (netlink.Link, error) {
	link, err := h.LinkByName(name)
	if err != nil {
		switch err.(type) {
		case netlink.LinkNotFoundError:
			return nil, errors.ErrNotFound
		default:
			return nil, err
		}
	}

	return link, nil
}

func (v *NetworkManager) findLink(name string) (netlink.Link, error) {
	h, err := netlink.NewHandle()
	if err != nil {
		return nil, err
	}

	return v.findLinkWithHandler(h, name)
}

func (v *NetworkManager) findBridgeWithHandler(h *netlink.Handle, name string) (netlink.Link, error) {
	link, err := v.findLinkWithHandler(h, name)
	if err != nil {
		return nil, err
	}

	if link.Type() != "bridge" {
		return nil, errors.ErrInvalidType
	}
	return link, nil
}

func (v *NetworkManager) findBridge(name string) (netlink.Link, error) {
	h, err := netlink.NewHandle()
	if err != nil {
		return nil, err
	}

	return v.findBridgeWithHandler(h, name)
}

func (v *NetworkManager) ensureBridgeExists(spec *entities.Network) error {
	logger := zap.L().With(zap.String("name", spec.Name))

	logger.Debug("ensuring bridge exists")

	name := v.getBridgeName(spec)

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
	attrs.MTU = spec.Mtu
	attrs.Name = name
	attrs.Flags = attrs.Flags | net.FlagUp

	return netlink.LinkAdd(&netlink.Bridge{
		LinkAttrs: attrs,
	})
}

func (v *NetworkManager) ensureBridgeNotExist(spec *entities.Network) error {
	logger := zap.L().With(zap.String("name", spec.Name))

	logger.Debug("ensuring bridge not exist")

	name := v.getBridgeName(spec)

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

func (v *NetworkManager) ensurePortAttached(ctx context.Context, pid int, port *entities.Port) error {
	logger := zap.L().With(zap.Int("pid", pid)).With(zap.String("name", port.Name))

	logger.Debug("ensuring port is attached")

	nsHandle, err := netns.GetFromPid(pid)
	if err != nil {
		return err
	}

	containerNetlinkHandler, err := netlink.NewHandleAt(nsHandle)
	if err != nil {
		return err
	}

	bridge, err := v.findBridge(v.getBridgeName(port.Network))
	if err != nil {
		return err
	}

	link, err := v.findLink(v.getPortName(port))
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if link == nil {
		attrs := netlink.NewLinkAttrs()
		attrs.MTU = bridge.Attrs().MTU
		attrs.Name = v.getPortName(port)
		attrs.Flags = attrs.Flags | net.FlagUp
		attrs.MasterIndex = bridge.Attrs().Index

		link = &netlink.Veth{
			LinkAttrs:     attrs,
			PeerName:      port.Name,
			PeerNamespace: netlink.NsPid(pid),
		}

		if err := netlink.LinkAdd(link); err != nil {
			return err
		}
	}

	peer, err := v.findLinkWithHandler(containerNetlinkHandler, port.Name)
	if err != nil {
		return err
	}

	if err := containerNetlinkHandler.LinkSetUp(peer); err != nil {
		return err
	}

	if err := containerNetlinkHandler.LinkSetMTU(peer, bridge.Attrs().MTU); err != nil {
		return err
	}

	for _, addr := range port.IPAddrs {
		err = containerNetlinkHandler.AddrReplace(peer, &netlink.Addr{
			IPNet: addr,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *NetworkManager) ensurePortNotExist(ctx context.Context, port *entities.Port) error {
	logger := zap.L().With(zap.String("name", port.Name))

	logger.Debug("ensuring port not exist")

	link, err := v.findLink(v.getPortName(port))
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if link == nil {
		return nil
	}

	return netlink.LinkDel(link)
}

func (v *NetworkManager) create(spec *entities.Network) error {
	logger := v.getLogger(spec)

	logger.Debug("creating bridge")

	err := v.ensureBridgeExists(spec)
	if err != nil {
		return err
	}

	logger.Debug("created bridge")
	return nil
}

func (v *NetworkManager) delete(spec *entities.Network) error {
	logger := v.getLogger(spec)

	logger.Debug("deleting bridge")

	err := v.ensureBridgeNotExist(spec)
	if err != nil {
		return err
	}

	logger.Debug("deleted bridge")
	return nil
}

func (v *NetworkManager) createPorts(ctx context.Context, pid int, ports []*entities.Port) error {
	logger := zap.L()

	logger.Debug("creating ports")

	for _, port := range ports {
		err := v.ensurePortAttached(ctx, pid, port)
		if err != nil {
			return err
		}
	}

	logger.Debug("created ports")
	return nil
}

func (v *NetworkManager) deletePorts(ctx context.Context, ports []*entities.Port) error {
	logger := zap.L()

	logger.Debug("deleting ports")

	for _, port := range ports {
		err := v.ensurePortNotExist(ctx, port)
		if err != nil {
			return err
		}
	}

	logger.Debug("deleted ports")
	return nil
}
