package usecases

import (
	"context"

	"github.com/proelbtn/vnet/pkg/usecases/managers"
)

// Create laboratory instantly. This usecase doesn'nt need LaboratoryGateway.
// This usecase expects to be used by standalone vnet client.
type InstantLaboratoryUsecase struct {
	laboratoryManager managers.LaboratoryManager
}

func NewInstantLaboratoryUsecase(laboratoryManager managers.LaboratoryManager) *InstantLaboratoryUsecase {
	return &InstantLaboratoryUsecase{
		laboratoryManager: laboratoryManager,
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
