package service

import (
	"context"

	"github.com/lockw1n/time-logger/internal/activity/domain"
)

type Service interface {
	CreateActivity(ctx context.Context, input CreateActivityInput) (domain.Activity, error)
	UpdateActivity(ctx context.Context, id uint64, input UpdateActivityInput) (domain.Activity, error)
	DeleteActivity(ctx context.Context, id uint64) error

	GetActivity(ctx context.Context, id uint64) (domain.Activity, error)
	ListActivitiesForCompany(ctx context.Context, companyID uint64) ([]domain.Activity, error)
}
