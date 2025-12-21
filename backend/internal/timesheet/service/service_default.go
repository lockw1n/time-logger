package service

import (
	"context"

	entryrepo "github.com/lockw1n/time-logger/internal/repository/entry"
	"github.com/lockw1n/time-logger/internal/timesheet/domain"
)

type timesheet struct {
	entryRepo entryrepo.Repository
}

func NewTimesheet(entryRepo entryrepo.Repository) Timesheet {
	return &timesheet{entryRepo: entryRepo}
}

func (s *timesheet) GenerateReport(ctx context.Context, cmd GenerateReportCommand) (*domain.Report, error) {
	if err := validateReportScope(cmd.ConsultantID, cmd.CompanyID); err != nil {
		return nil, err
	}

	start, err := parseDate(cmd.Start)
	if err != nil {
		return nil, err
	}

	end, err := parseDate(cmd.End)
	if err != nil {
		return nil, err
	}

	if err := validateDateRange(start, end); err != nil {
		return nil, err
	}

	entries, err := s.entryRepo.FindForPeriodWithDetails(ctx, cmd.ConsultantID, cmd.CompanyID, start, end)
	if err != nil {
		return nil, err
	}

	rows := groupEntries(entries)
	sortRowsByTicketCode(rows)

	return &domain.Report{
		ConsultantID: cmd.ConsultantID,
		CompanyID:    cmd.CompanyID,
		Start:        cmd.Start,
		End:          cmd.End,
		Rows:         rows,
	}, nil
}
