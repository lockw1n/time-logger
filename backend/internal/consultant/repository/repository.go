package repository

import (
	"context"

	"github.com/lockw1n/time-logger/internal/consultant/domain"
)

type Repository interface {
	Create(ctx context.Context, consultant domain.Consultant) (domain.Consultant, error)
	Update(ctx context.Context, consultant domain.Consultant) (domain.Consultant, error)
	Delete(ctx context.Context, id uint64) error

	FindByID(ctx context.Context, id uint64) (domain.Consultant, error)
}
