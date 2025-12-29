package service

import (
	"context"

	"github.com/lockw1n/time-logger/internal/consultant/domain"
)

type Service interface {
	CreateConsultant(ctx context.Context, input CreateConsultantInput) (domain.Consultant, error)
	UpdateConsultant(ctx context.Context, id uint64, input UpdateConsultantInput) (domain.Consultant, error)
	DeleteConsultant(ctx context.Context, id uint64) error

	GetConsultant(ctx context.Context, id uint64) (domain.Consultant, error)
}
