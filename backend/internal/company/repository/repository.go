package repository

import (
	"context"

	"github.com/lockw1n/time-logger/internal/company/domain"
)

type Repository interface {
	Create(ctx context.Context, company domain.Company) (domain.Company, error)
	Update(ctx context.Context, company domain.Company) (domain.Company, error)
	Delete(ctx context.Context, id uint64) error

	FindByID(ctx context.Context, id uint64) (domain.Company, error)
	ListAll(ctx context.Context) ([]domain.Company, error)
	ListByConsultant(ctx context.Context, consultantID uint64) ([]domain.Company, error)
}
