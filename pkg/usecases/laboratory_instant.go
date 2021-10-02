package usecases

import (
	"context"

	"github.com/proelbtn/vnet/pkg/errors"
	"github.com/proelbtn/vnet/pkg/usecases/managers"
)

// Create laboratory instantly. This usecase doesn'nt need LaboratoryGateway.
// This usecase expects to be used by standalone vnet client.
type InstantLaboratoryUsecase struct {
	laboratoryManager managers.LaboratoryManager
	containerManager  managers.ContainerManager
	networkManager    managers.NetworkManager
}

func NewInstantLaboratoryUsecase(laboratoryManager managers.LaboratoryManager, containerManager managers.ContainerManager, networkManager managers.NetworkManager) *InstantLaboratoryUsecase {
	return &InstantLaboratoryUsecase{
		laboratoryManager: laboratoryManager,
		containerManager:  containerManager,
		networkManager:    networkManager,
	}
}

func (v *InstantLaboratoryUsecase) StartLaboratory(req *WritableLaboratory) error {
	ctx := context.Background()
	lab, err := req.ToEntity()
	if err != nil {
		return err
	}

	return v.laboratoryManager.Start(ctx, lab)
}

func (v *InstantLaboratoryUsecase) StopLaboratory(req *WritableLaboratory) error {
	ctx := context.Background()
	lab, err := req.ToEntity()
	if err != nil {
		return err
	}

	return v.laboratoryManager.Stop(ctx, lab)
}

func (v *InstantLaboratoryUsecase) Exec(req *WritableLaboratory, name string, args []string) error {
	ctx := context.Background()
	lab, err := req.ToEntity()
	if err != nil {
		return err
	}

	for _, container := range lab.Containers {
		if container.Name == name {
			execArgs := managers.ExecArgs{
				Args: args,
			}
			if err := v.containerManager.Exec(ctx, container, execArgs); err != nil {
				return err
			}
			return nil
		}
	}

	return errors.ErrNotFound
}
