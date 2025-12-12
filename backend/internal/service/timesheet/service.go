package timesheet

import (
	timesheetdto "github.com/lockw1n/time-logger/internal/dto/timesheet"
)

type Service interface {
	GenerateReport(req timesheetdto.Request) (*timesheetdto.Report, error)
}
