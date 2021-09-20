package gateways

import (
	"github.com/google/uuid"
	"github.com/proelbtn/vnet/pkg/entities"
)

type LaboratoryGateway interface {
	Save(env *entities.Laboratory) error
	GetAll() ([]entities.Laboratory, error)
	FindByID(id uuid.UUID) (*entities.Laboratory, error)
	FindByName(name string) (*entities.Laboratory, error)
	Delete(env *entities.Laboratory) error
}
