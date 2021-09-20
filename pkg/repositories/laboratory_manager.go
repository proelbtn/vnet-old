package repositories

import (
	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/usecases/managers"
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

func (v *LaboratoryManager) Start(lab *entities.Laboratory) error {
	for _, network := range lab.Networks {
		err := v.networkManager.Create(network)
		if err != nil {
			return err
		}
	}

	for _, container := range lab.Containers {
		pid, err := v.containerManager.Create(container)
		if err != nil {
			return err
		}

		err = v.networkManager.AttachPorts(int(pid), container.Ports)
		if err != nil {
			return err
		}

		err = v.containerManager.Start(container)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *LaboratoryManager) Stop(lab *entities.Laboratory) error {
	for _, container := range lab.Containers {
		err := v.containerManager.Stop(container)
		if err != nil {
			return err
		}

		err = v.containerManager.Delete(container)
		if err != nil {
			return err
		}
	}

	for _, network := range lab.Networks {
		err := v.networkManager.Delete(network)
		if err != nil {
			return err
		}
	}

	return nil
}
