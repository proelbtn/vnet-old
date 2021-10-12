package repositories

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/Masterminds/semver"
	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/errors"
	"github.com/proelbtn/vnet/pkg/usecases/managers"
	"go.uber.org/zap"
)

type LaboratoryManager struct {
	containerManager managers.ContainerManager
	networkManager   managers.NetworkManager
}

var _ managers.LaboratoryManager = (*LaboratoryManager)(nil)

func NewLaboratoryManager(containerManager managers.ContainerManager, networkManager managers.NetworkManager) *LaboratoryManager {
	return &LaboratoryManager{
		containerManager: containerManager,
		networkManager:   networkManager,
	}
}

func (v *LaboratoryManager) getLogger(lab *entities.Laboratory) *zap.Logger {
	return zap.L().With(
		zap.String("Name", lab.Name),
	)
}

func (v *LaboratoryManager) getKernelVersion() (*semver.Version, error) {
	uname := &syscall.Utsname{}
	syscall.Uname(uname)

	var release string
	for _, c := range uname.Release {
		release += string(int(c))
	}

	return semver.NewVersion(release)
}

func (v *LaboratoryManager) checkKernelVersion(constraint string) (bool, error) {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return false, err
	}

	version, err := v.getKernelVersion()
	if err != nil {
		return false, err
	}

	return c.Check(version), nil
}

func (v *LaboratoryManager) checkKernelModule(module string) (bool, error) {
	path := fmt.Sprintf("/sys/module/%s", module)

	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, errors.ErrNotFound
		}
		return false, err
	}

	if info.Mode()&os.ModeDir != 0 {
		return true, nil
	}

	return false, nil
}

func (v *LaboratoryManager) checkKernelModules(modules []string) (bool, error) {
	for _, module := range modules {
		if verdict, err := v.checkKernelModule(module); !verdict {
			return verdict, err
		}
	}
	return true, nil
}

func (v *LaboratoryManager) checkPreRequisites(lab *entities.Laboratory) (bool, error) {
	if verdict, err := v.checkKernelVersion(lab.PreRequisites.KernelVersion); !verdict {
		return verdict, err
	}

	if verdict, err := v.checkKernelModules(lab.PreRequisites.Modules); !verdict {
		return verdict, err
	}

	return true, nil
}

func (v *LaboratoryManager) ensureKernelParameters() error {
	params := []struct {
		key   string
		value []byte
	}{
		{
			key:   "/proc/sys/net/ipv4/conf/all/forwarding",
			value: []byte("1"),
		},
		{
			key:   "/proc/sys/net/ipv6/conf/all/forwarding",
			value: []byte("1"),
		},
	}

	for _, param := range params {
		err := ioutil.WriteFile(param.key, param.value, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *LaboratoryManager) Start(ctx context.Context, lab *entities.Laboratory) error {
	logger := v.getLogger(lab)

	logger.Debug("starting Laboratory")

	logger.Debug("checking prerequisites")
	if verdict, err := v.checkPreRequisites(lab); err != nil {
		return err
	} else if !verdict {
		return errors.ErrInvalidType
	}

	logger.Debug("checking kernel parameters")
	if err := v.ensureKernelParameters(); err != nil {
		return err
	}

	logger.Debug("creating Networks")
	for _, network := range lab.Networks {
		err := v.networkManager.Create(ctx, network)
		if err != nil {
			return err
		}
	}

	logger.Debug("creating Containers")
	for _, container := range lab.Containers {
		pid, err := v.containerManager.Create(ctx, container)
		if err != nil {
			return err
		}

		err = v.networkManager.CreatePorts(ctx, int(pid), container.Ports)
		if err != nil {
			return err
		}

		err = v.containerManager.Start(ctx, container)
		if err != nil {
			return err
		}
	}

	logger.Debug("started Laboratory")
	return nil
}

func (v *LaboratoryManager) Stop(ctx context.Context, lab *entities.Laboratory) error {
	logger := v.getLogger(lab)

	logger.Debug("stopping Laboratory")

	logger.Debug("stopping Containers")
	for _, container := range lab.Containers {
		err := v.containerManager.Stop(ctx, container)
		if err != nil {
			return err
		}

		err = v.networkManager.DeletePorts(ctx, container.Ports)
		if err != nil {
			return err
		}

		err = v.containerManager.Delete(ctx, container)
		if err != nil {
			return err
		}
	}

	logger.Debug("stopping Networks")
	for _, network := range lab.Networks {
		err := v.networkManager.Delete(ctx, network)
		if err != nil {
			return err
		}
	}

	logger.Debug("stopped Laboratory")
	return nil
}
