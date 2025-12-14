package consultantassignment

import (
	"context"
	"time"

	"github.com/lockw1n/time-logger/internal/models"
)

type Repository interface {
	Create(assignment *models.ConsultantAssignment) (*models.ConsultantAssignment, error)
	Update(assignment *models.ConsultantAssignment) (*models.ConsultantAssignment, error)
	Delete(id uint64) error

	FindByID(id uint64) (*models.ConsultantAssignment, error)
	FindByPair(consultantID uint64, companyID uint64) (*models.ConsultantAssignment, error)
	FindByConsultant(consultantID uint64) ([]models.ConsultantAssignment, error)
	FindByCompany(companyID uint64) ([]models.ConsultantAssignment, error)
	FindActiveForPeriod(
		ctx context.Context,
		consultantID uint64,
		companyID uint64,
		start time.Time,
		end time.Time,
	) ([]models.ConsultantAssignment, error)
}
