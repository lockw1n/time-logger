package entry

import (
	"time"

	"github.com/lockw1n/time-logger/internal/models"
)

type Repository interface {
	Create(entry *models.Entry) (*models.Entry, error)
	Update(entry *models.Entry) (*models.Entry, error)
	Delete(id uint64) error

	FindByID(id uint64) (*models.Entry, error)
	FindByCompany(companyID uint64) ([]models.Entry, error)
	FindByConsultant(consultantID uint64) ([]models.Entry, error)
	ListAll() ([]models.Entry, error)
	FindWithDetails(consultantID uint64, companyID uint64, start time.Time, end time.Time) ([]models.Entry, error)
}
