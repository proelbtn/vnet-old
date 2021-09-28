package managers

import (
	"context"

	"github.com/proelbtn/vnet/pkg/entities"
)

type LaboratoryManager interface {
	Start(ctx context.Context, lab *entities.Laboratory) error
	Stop(ctx context.Context, lab *entities.Laboratory) error
}
