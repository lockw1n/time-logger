package repository

import (
	"context"

	"github.com/lockw1n/time-logger/internal/ticket/domain"
)

type Repository interface {
	Create(ctx context.Context, ticket domain.Ticket) (domain.Ticket, error)
	Update(ctx context.Context, ticket domain.Ticket) (domain.Ticket, error)
	Delete(ctx context.Context, id uint64) error

	FindByID(ctx context.Context, id uint64) (domain.Ticket, error)
	ListByIDs(ctx context.Context, ids []uint64) ([]domain.Ticket, error)
	ListByCompany(ctx context.Context, companyID uint64) ([]domain.Ticket, error)
	FindByCompanyAndCode(ctx context.Context, companyID uint64, code string) (domain.Ticket, error)
}
