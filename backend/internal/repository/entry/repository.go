package entry

import (
	"context"
	"time"

	"github.com/lockw1n/time-logger/internal/models"
)

type Repository interface {
	Create(ctx context.Context, entry *models.Entry) (*models.Entry, error)
	Update(ctx context.Context, entry *models.Entry) (*models.Entry, error)
	Delete(ctx context.Context, id uint64) error

	FindByID(ctx context.Context, id uint64) (*models.Entry, error)
	FindForPeriodWithDetails(
		ctx context.Context,
		consultantID uint64,
		companyID uint64,
		start time.Time,
		end time.Time,
	) ([]models.Entry, error)
}
