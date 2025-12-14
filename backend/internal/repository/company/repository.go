package company

import (
	"context"

	"github.com/lockw1n/time-logger/internal/models"
)

type Repository interface {
	Create(company *models.Company) (*models.Company, error)
	Update(company *models.Company) (*models.Company, error)
	Delete(id uint64) error

	FindByID(ctx context.Context, id uint64) (*models.Company, error)
	FindByName(name string) (*models.Company, error)
	ListAll() ([]models.Company, error)
}
