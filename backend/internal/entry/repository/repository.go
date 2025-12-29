package repository

import (
	"context"
	"time"

	"github.com/lockw1n/time-logger/internal/entry/domain"
)

type Repository interface {
	Create(ctx context.Context, entry domain.Entry) (domain.Entry, error)
	Update(ctx context.Context, entry domain.Entry) (domain.Entry, error)
	Delete(ctx context.Context, id uint64) error

	FindByID(ctx context.Context, id uint64) (domain.Entry, error)
	ListForConsultantPeriod(ctx context.Context,
		consultantID uint64,
		companyID uint64,
		start time.Time,
		end time.Time,
	) ([]domain.Entry, error)
}
