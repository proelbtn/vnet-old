package usecases

import (
	"errors"

	"github.com/google/uuid"
	"github.com/proelbtn/vnet/pkg/usecases/gateways"
)

var (
	ErrLaboratoryAlreadyExists = errors.New("laboratory already exists")
	ErrLaboratoryNotFound      = errors.New("laboratory not found")
)

type LaboratoryUsecase struct {
	laboratoryGateway gateways.LaboratoryGateway
}

func NewLaboratoryUsecase(laboratoryGateway gateways.LaboratoryGateway) *LaboratoryUsecase {
	return &LaboratoryUsecase{
		laboratoryGateway: laboratoryGateway,
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

func (v *LaboratoryUsecase) DeleteLaboratory(identifier string) error {
	if id, err := uuid.Parse(identifier); err == nil {
		return v.deleteLaboratoryWithID(id)
	} else {
		return v.deleteLaboratoryWithName(identifier)
	}
}

func (v *LaboratoryUsecase) deleteLaboratoryWithID(id uuid.UUID) error {
	laboratory, err := v.laboratoryGateway.FindByID(id)
	if err != nil {
		return err
	}

	return v.laboratoryGateway.Delete(laboratory)
}

func (v *LaboratoryUsecase) deleteLaboratoryWithName(name string) error {
	laboratory, err := v.laboratoryGateway.FindByName(name)
	if err != nil {
		return ErrLaboratoryNotFound
	}

	return v.laboratoryGateway.Delete(laboratory)
}
