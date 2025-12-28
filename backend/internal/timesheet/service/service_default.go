package service

import (
	"context"

	activitydomain "github.com/lockw1n/time-logger/internal/activity/domain"
	activityrepo "github.com/lockw1n/time-logger/internal/activity/repository"
	entryrepo "github.com/lockw1n/time-logger/internal/entry/repository"
	ticketdomain "github.com/lockw1n/time-logger/internal/ticket/domain"
	ticketrepo "github.com/lockw1n/time-logger/internal/ticket/repository"
	"github.com/lockw1n/time-logger/internal/timesheet/domain"
)

type service struct {
	activityRepo activityrepo.Repository
	entryRepo    entryrepo.Repository
	ticketRepo   ticketrepo.Repository
}

func NewService(
	activityRepo activityrepo.Repository,
	entryRepo entryrepo.Repository,
	ticketRepo ticketrepo.Repository,
) Service {
	return &service{
		activityRepo: activityRepo,
		entryRepo:    entryRepo,
		ticketRepo:   ticketRepo,
	}
}

func (s *service) GenerateTimesheet(ctx context.Context, input GenerateTimesheetInput) (domain.Timesheet, error) {
	if err := input.Validate(); err != nil {
		return domain.Timesheet{}, ErrTimesheetInvalid
	}

	entries, err := s.entryRepo.ListForConsultantPeriod(
		ctx,
		input.ConsultantID,
		input.CompanyID,
		input.Start,
		input.End,
	)
	if err != nil {
		return domain.Timesheet{}, ErrTimesheetConflict
	}

	if len(entries) == 0 {
		return domain.Timesheet{
			ConsultantID: input.ConsultantID,
			CompanyID:    input.CompanyID,
			Start:        input.Start,
			End:          input.End,
			Rows:         []domain.TimesheetRow{},
		}, nil
	}

	ticketIDSet := make(map[uint64]struct{})
	activityIDSet := make(map[uint64]struct{})

	for _, e := range entries {
		ticketIDSet[e.TicketID] = struct{}{}
		activityIDSet[e.ActivityID] = struct{}{}
	}

	tickets, err := s.ticketRepo.ListByIDs(ctx, keys(ticketIDSet))
	if err != nil {
		return domain.Timesheet{}, ErrTimesheetConflict
	}

	activities, err := s.activityRepo.ListByIDs(ctx, keys(activityIDSet))
	if err != nil {
		return domain.Timesheet{}, ErrTimesheetConflict
	}

	ticketsByID := make(map[uint64]ticketdomain.Ticket, len(tickets))
	for _, ticket := range tickets {
		ticketsByID[ticket.ID] = ticket
	}

	activitiesByID := make(map[uint64]activitydomain.Activity, len(activities))
	for _, activity := range activities {
		activitiesByID[activity.ID] = activity
	}

	rows := groupEntries(entries, ticketsByID, activitiesByID)
	sortRowsByTicketCode(rows)

	return domain.Timesheet{
		ConsultantID: input.ConsultantID,
		CompanyID:    input.CompanyID,
		Start:        input.Start,
		End:          input.End,
		Rows:         rows,
	}, nil
}

func keys[K comparable](m map[K]struct{}) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
