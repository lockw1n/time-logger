package service

import (
	"context"

	"github.com/lockw1n/time-logger/internal/timesheet/domain"
)

type Timesheet interface {
	GenerateReport(ctx context.Context, cmd GenerateReportCommand) (*domain.Report, error)
}
