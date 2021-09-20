package managers

import "github.com/proelbtn/vnet/pkg/entities"

type LaboratoryManager interface {
	Start(lab *entities.Laboratory) error
	Stop(lab *entities.Laboratory) error
}
