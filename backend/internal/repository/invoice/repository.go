package invoice

import "github.com/lockw1n/time-logger/internal/models"

type Repository interface {
	Create(invoice *models.Invoice) error
	Update(invoice *models.Invoice) error
	Delete(id uint64) error

	FindByID(id uint64) (*models.Invoice, error)
	FindByCompany(companyID uint64) ([]models.Invoice, error)
	FindByConsultant(consultantID uint64) ([]models.Invoice, error)
	ListAll() ([]models.Invoice, error)
}
