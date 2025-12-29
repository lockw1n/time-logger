package service

import (
	"context"

	"github.com/lockw1n/time-logger/internal/company/domain"
)

type Service interface {
	CreateCompany(ctx context.Context, input CreateCompanyInput) (domain.Company, error)
	UpdateCompany(ctx context.Context, id uint64, input UpdateCompanyInput) (domain.Company, error)
	DeleteCompany(ctx context.Context, id uint64) error

	GetCompany(ctx context.Context, id uint64) (domain.Company, error)
	ListCompanies(ctx context.Context) ([]domain.Company, error)
	ListCompaniesForConsultant(ctx context.Context, consultantID uint64) ([]domain.Company, error)
}
