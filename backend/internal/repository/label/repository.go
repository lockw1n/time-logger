package label

import "github.com/lockw1n/time-logger/internal/models"

type Repository interface {
	Create(label *models.Label) error
	Update(label *models.Label) error
	Delete(id uint64) error

	FindByID(id uint64) (*models.Label, error)
	FindByCompany(companyID uint64) ([]models.Label, error)
	FindByName(companyID uint64, name string) (*models.Label, error)
}
