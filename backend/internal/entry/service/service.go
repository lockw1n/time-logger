package service

import (
	"context"

	"github.com/lockw1n/time-logger/internal/entry/domain"
)

type Service interface {
	CreateEntry(ctx context.Context, consultantID uint64, input CreateEntryInput) (domain.Entry, error)
	UpdateEntry(ctx context.Context, consultantID uint64, entryID uint64, input UpdateEntryInput) (domain.Entry, error)
	DeleteEntry(ctx context.Context, consultantID uint64, entryID uint64) error

	GetEntry(ctx context.Context, consultantID uint64, entryID uint64) (domain.Entry, error)
}
