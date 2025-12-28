package service

import (
	"context"

	"github.com/lockw1n/time-logger/internal/timesheet/domain"
)

type Service interface {
	GenerateTimesheet(ctx context.Context, input GenerateTimesheetInput) (domain.Timesheet, error)
}
