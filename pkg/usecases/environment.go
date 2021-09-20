package usecases

import (
	"errors"

	"github.com/google/uuid"
	"github.com/proelbtn/vnet/pkg/usecases/gateways"
)

var (
	ErrEnvironmentAlreadyExists = errors.New("environment already exists")
	ErrEnvironmentNotFound      = errors.New("environment not found")
)

type EnvironmentUsecase struct {
	environmentGateway gateways.EnvironmentGateway
}

func NewEnvironmentUsecase(environmentGateway gateways.EnvironmentGateway) *EnvironmentUsecase {
	return &EnvironmentUsecase{
		environmentGateway: environmentGateway,
	}
}

func (v *EnvironmentUsecase) CreateEnvironment(req WritableEnvironment) (*Environment, error) {
	saved, err := v.environmentGateway.FindByName(req.Name)
	if err != nil {
		return nil, err
	}
	if saved != nil {
		return nil, ErrEnvironmentAlreadyExists
	}

	environment, err := req.ToEntity()
	if err != nil {
		return nil, err
	}

	err = v.environmentGateway.Save(environment)
	if err != nil {
		return nil, err
	}

	return NewEnvironment(environment), nil
}

func (v *EnvironmentUsecase) DeleteEnvironment(identifier string) error {
	if id, err := uuid.Parse(identifier); err == nil {
		return v.deleteEnvironmentWithID(id)
	} else {
		return v.deleteEnvironmentWithName(identifier)
	}
}

func (v *EnvironmentUsecase) deleteEnvironmentWithID(id uuid.UUID) error {
	environment, err := v.environmentGateway.FindByID(id)
	if err != nil {
		return err
	}

	return v.environmentGateway.Delete(environment)
}

func (v *EnvironmentUsecase) deleteEnvironmentWithName(name string) error {
	environment, err := v.environmentGateway.FindByName(name)
	if err != nil {
		return ErrEnvironmentNotFound
	}

	return v.environmentGateway.Delete(environment)
}
