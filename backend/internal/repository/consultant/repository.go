package consultant

import (
	"context"

	"github.com/lockw1n/time-logger/internal/models"
)

type Repository interface {
	Create(consultant *models.Consultant) (*models.Consultant, error)
	Update(consultant *models.Consultant) (*models.Consultant, error)
	Delete(id uint64) error

	FindByID(ctx context.Context, id uint64) (*models.Consultant, error)
	FindByLastName(lastName string) ([]models.Consultant, error)
	ListAll() ([]models.Consultant, error)
}
