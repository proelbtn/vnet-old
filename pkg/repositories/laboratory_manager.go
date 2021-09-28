package repositories

import (
	"context"

	"github.com/proelbtn/vnet/pkg/entities"
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

func (v *LaboratoryManager) Start(ctx context.Context, lab *entities.Laboratory) error {
	logger := v.getLogger(lab)

	logger.Debug("starting Laboratory")

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

		err = v.networkManager.AttachPorts(ctx, int(pid), container.Ports)
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
