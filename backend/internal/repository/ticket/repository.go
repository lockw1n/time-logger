package ticket

import "github.com/lockw1n/time-logger/internal/models"

type Repository interface {
	Create(ticket *models.Ticket) (*models.Ticket, error)
	Update(ticket *models.Ticket) (*models.Ticket, error)
	Delete(id uint64) error

	FindByID(id uint64) (*models.Ticket, error)
	FindByCode(companyID uint64, code string) (*models.Ticket, error)
	FindByCompany(companyID uint64) ([]models.Ticket, error)
	ListAll() ([]models.Ticket, error)
}
