package repository

import (
	"context"
	"time"

	"github.com/lockw1n/time-logger/internal/contract/domain"
)

type Repository interface {
	Create(ctx context.Context, contract domain.Contract) (domain.Contract, error)
	Update(ctx context.Context, contract domain.Contract) (domain.Contract, error)
	Delete(ctx context.Context, id uint64) error

	FindByID(ctx context.Context, id uint64) (domain.Contract, error)
	ListByConsultant(ctx context.Context, consultantID uint64) ([]domain.Contract, error)
	FindActiveByConsultantCompany(
		ctx context.Context,
		consultantID uint64,
		companyID uint64,
		at time.Time,
	) (domain.Contract, error)
}
