package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/usecases/gateways"
	"github.com/proelbtn/vnet/pkg/usecases/managers"
)

var (
	ErrLaboratoryAlreadyExists = errors.New("laboratory already exists")
	ErrLaboratoryNotFound      = errors.New("laboratory not found")
)

type LaboratoryUsecase struct {
	laboratoryGateway gateways.LaboratoryGateway
	laboratoryManager managers.LaboratoryManager
}

func NewLaboratoryUsecase(laboratoryGateway gateways.LaboratoryGateway, laboratoryManager managers.LaboratoryManager) *LaboratoryUsecase {
	return &LaboratoryUsecase{
		laboratoryGateway: laboratoryGateway,
		laboratoryManager: laboratoryManager,
	}
}

func (v *LaboratoryUsecase) CreateLaboratory(req WritableLaboratory) (*Laboratory, error) {
	saved, err := v.laboratoryGateway.FindByName(req.Name)
	if err != nil {
		return nil, err
	}
	if saved != nil {
		return nil, ErrLaboratoryAlreadyExists
	}

	laboratory, err := req.ToEntity()
	if err != nil {
		return nil, err
	}

	err = v.laboratoryGateway.Save(laboratory)
	if err != nil {
		return nil, err
	}

	return NewLaboratory(laboratory), nil
}

func (v *LaboratoryUsecase) findLaboratoryByIdentifier(identifier string) (*entities.Laboratory, error) {
	if id, err := uuid.Parse(identifier); err == nil {
		return v.laboratoryGateway.FindByID(id)
	} else {
		return v.laboratoryGateway.FindByName(identifier)
	}
}

func (v *LaboratoryUsecase) DeleteLaboratory(identifier string) error {
	lab, err := v.findLaboratoryByIdentifier(identifier)
	if err != nil {
		return err
	}

	return v.laboratoryGateway.Delete(lab)
}

func (v *LaboratoryUsecase) StartLaboratory(identifier string) error {
	ctx := context.Background()
	lab, err := v.findLaboratoryByIdentifier(identifier)
	if err != nil {
		return err
	}

	return v.laboratoryManager.Start(ctx, lab)
}

func (v *LaboratoryUsecase) StopLaboratory(identifier string) error {
	ctx := context.Background()
	lab, err := v.findLaboratoryByIdentifier(identifier)
	if err != nil {
		return err
	}

	return v.laboratoryManager.Stop(ctx, lab)
}
