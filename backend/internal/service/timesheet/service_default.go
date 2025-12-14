package timesheet

import (
	"context"

	timesheetdto "github.com/lockw1n/time-logger/internal/dto/timesheet"
	entryrepo "github.com/lockw1n/time-logger/internal/repository/entry"
)

type service struct {
	entryRepo entryrepo.Repository
}

func NewService(entryRepo entryrepo.Repository) Service {
	return &service{entryRepo: entryRepo}
}

func (s *service) GenerateReport(ctx context.Context, req timesheetdto.Request) (*timesheetdto.Report, error) {
	if err := validateReportScope(req.ConsultantID, req.CompanyID); err != nil {
		return nil, err
	}

	start, err := parseDate(req.Start)
	if err != nil {
		return nil, err
	}

	end, err := parseDate(req.End)
	if err != nil {
		return nil, err
	}

	if err := validateDateRange(start, end); err != nil {
		return nil, err
	}

	entries, err := s.entryRepo.FindForPeriodWithDetails(ctx, req.ConsultantID, req.CompanyID, start, end)
	if err != nil {
		return nil, err
	}

	rows := groupEntries(entries)
	sortRowsByTicketCode(rows)

	return &timesheetdto.Report{
		ConsultantID: req.ConsultantID,
		CompanyID:    req.CompanyID,
		Start:        req.Start,
		End:          req.End,
		Rows:         rows,
	}, nil
}
