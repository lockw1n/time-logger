package timesheet

import (
	"context"

	timesheetdto "github.com/lockw1n/time-logger/internal/dto/timesheet"
)

type Service interface {
	GenerateReport(ctx context.Context, req timesheetdto.Request) (*timesheetdto.Report, error)
}
