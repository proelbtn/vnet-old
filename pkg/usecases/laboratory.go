package usecases

import (
	"context"

	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/errors"
	"github.com/proelbtn/vnet/pkg/usecases/managers"
)

// Create laboratory instantly. This usecase doesn'nt need LaboratoryGateway.
// This usecase expects to be used by standalone vnet client.
type LaboratoryUsecase struct {
	laboratoryManager managers.LaboratoryManager
	containerManager  managers.ContainerManager
	networkManager    managers.NetworkManager
}

func NewLaboratoryUsecase(laboratoryManager managers.LaboratoryManager, containerManager managers.ContainerManager, networkManager managers.NetworkManager) *LaboratoryUsecase {
	return &LaboratoryUsecase{
		laboratoryManager: laboratoryManager,
		containerManager:  containerManager,
		networkManager:    networkManager,
	}
}

func (v *LaboratoryUsecase) findNetwork(lab *entities.Laboratory, name string) (*entities.Network, error) {
	for _, container := range lab.Networks {
		if container.Name == name {
			return container, nil
		}
	}

	return nil, errors.ErrNotFound
}

func (v *LaboratoryUsecase) findContainer(lab *entities.Laboratory, name string) (*entities.Container, error) {
	for _, container := range lab.Containers {
		if container.Name == name {
			return container, nil
		}
	}

	return nil, errors.ErrNotFound
}

func (v *LaboratoryUsecase) findPort(container *entities.Container, name string) (*entities.Port, error) {
	for _, port := range container.Ports {
		if port.Name == name {
			return port, nil
		}
	}

	return nil, errors.ErrNotFound
}

func (v *LaboratoryUsecase) StartLaboratory(req *WritableLaboratory) error {
	ctx := context.Background()
	lab, err := req.ToEntity()
	if err != nil {
		return err
	}

	return v.laboratoryManager.Start(ctx, lab)
}

func (v *LaboratoryUsecase) StopLaboratory(req *WritableLaboratory) error {
	ctx := context.Background()
	lab, err := req.ToEntity()
	if err != nil {
		return err
	}

	return v.laboratoryManager.Stop(ctx, lab)
}

func (v *LaboratoryUsecase) Exec(req *WritableLaboratory, name string, args []string) error {
	ctx := context.Background()
	lab, err := req.ToEntity()
	if err != nil {
		return err
	}

	container, err := v.findContainer(lab, name)
	if err != nil {
		return err
	}

	execArgs := managers.ExecArgs{
		Args: args,
	}
	return v.containerManager.Exec(ctx, container, execArgs)
}

func (v *LaboratoryUsecase) GetPortName(labName, containerName, portName string) string {
	return v.networkManager.GetPortName(labName, containerName, portName)
}

func (v *LaboratoryUsecase) GetBridgeName(labName, networkName string) string {
	return v.networkManager.GetBridgeName(labName, networkName)
}

func (v *LaboratoryUsecase) GetTopology(req *WritableLaboratory) (string, error) {
	lab, err := req.ToEntity()
	if err != nil {
		return "", err
	}

	return renderTopology(lab)
}
