package gateways

import (
	"github.com/google/uuid"
	"github.com/proelbtn/vnet/pkg/entities"
)

type EnvironmentGateway interface {
	Save(env *entities.Environment) error
	GetAll() ([]entities.Environment, error)
	FindByID(id uuid.UUID) (*entities.Environment, error)
	FindByName(name string) (*entities.Environment, error)
	Delete(env *entities.Environment) error
}
