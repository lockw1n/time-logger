package company

import "github.com/lockw1n/time-logger/internal/models"

type Repository interface {
	Create(company *models.Company) (*models.Company, error)
	Update(company *models.Company) (*models.Company, error)
	Delete(id uint64) error

	FindByID(id uint64) (*models.Company, error)
	FindByName(name string) (*models.Company, error)
	ListAll() ([]models.Company, error)
}
