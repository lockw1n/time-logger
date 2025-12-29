package repository

import (
	"context"

	"github.com/lockw1n/time-logger/internal/activity/domain"
)

type Repository interface {
	Create(ctx context.Context, activity domain.Activity) (domain.Activity, error)
	Update(ctx context.Context, activity domain.Activity) (domain.Activity, error)
	Delete(ctx context.Context, id uint64) error

	FindByID(ctx context.Context, id uint64) (domain.Activity, error)
	ListByIDs(ctx context.Context, ids []uint64) ([]domain.Activity, error)
	ListByCompany(ctx context.Context, companyID uint64) ([]domain.Activity, error)
}
