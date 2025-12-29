package service

import (
	"context"

	"github.com/lockw1n/time-logger/internal/contract/domain"
)

type Service interface {
	CreateContract(ctx context.Context, input CreateContractInput) (domain.Contract, error)
	UpdateContract(ctx context.Context, id uint64, input UpdateContractInput) (domain.Contract, error)
	DeleteContract(ctx context.Context, id uint64) error

	GetContract(ctx context.Context, id uint64) (domain.Contract, error)
	ListContractsForConsultant(ctx context.Context, consultantID uint64) ([]domain.Contract, error)
}
