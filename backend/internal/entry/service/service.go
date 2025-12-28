package service

import (
	"context"

	"github.com/lockw1n/time-logger/internal/entry/domain"
)

type Service interface {
	CreateEntry(ctx context.Context, input CreateEntryInput) (domain.Entry, error)
	UpdateEntry(ctx context.Context, id uint64, input UpdateEntryInput) (domain.Entry, error)
	DeleteEntry(ctx context.Context, id uint64) error

	GetEntry(ctx context.Context, id uint64) (domain.Entry, error)
}
