package service

import (
	"context"
	"errors"
	"time"

	activitydomain "github.com/lockw1n/time-logger/internal/activity/domain"
	activityrepo "github.com/lockw1n/time-logger/internal/activity/repository"
	companyrepo "github.com/lockw1n/time-logger/internal/company/repository"
	consultantrepo "github.com/lockw1n/time-logger/internal/consultant/repository"
	contractrepo "github.com/lockw1n/time-logger/internal/contract/repository"
	entryrepo "github.com/lockw1n/time-logger/internal/entry/repository"
	"github.com/lockw1n/time-logger/internal/invoice/domain"
	ticketdomain "github.com/lockw1n/time-logger/internal/ticket/domain"
	ticketrepo "github.com/lockw1n/time-logger/internal/ticket/repository"
)

type service struct {
	consultantRepo consultantrepo.Repository
	companyRepo    companyrepo.Repository
	contractRepo   contractrepo.Repository
	activityRepo   activityrepo.Repository
	entryRepo      entryrepo.Repository
	ticketRepo     ticketrepo.Repository
}

func NewService(
	consultantRepo consultantrepo.Repository,
	companyRepo companyrepo.Repository,
	contractRepo contractrepo.Repository,
	activityRepo activityrepo.Repository,
	entryRepo entryrepo.Repository,
	ticketRepo ticketrepo.Repository,
) Service {
	return &service{
		consultantRepo: consultantRepo,
		companyRepo:    companyRepo,
		contractRepo:   contractRepo,
		activityRepo:   activityRepo,
		entryRepo:      entryRepo,
		ticketRepo:     ticketRepo,
	}
}

func (s *service) GenerateInvoice(ctx context.Context, input GenerateInvoiceInput) (domain.Invoice, error) {
	if err := input.Validate(); err != nil {
		return domain.Invoice{}, ErrInvoiceInvalid
	}

	consultant, err := s.consultantRepo.FindByID(ctx, input.ConsultantID)
	if err != nil {
		return domain.Invoice{}, ErrInvoiceConflict
	}

	company, err := s.companyRepo.FindByID(ctx, input.CompanyID)
	if err != nil {
		return domain.Invoice{}, ErrInvoiceConflict
	}

	contract, err := s.contractRepo.FindActiveByConsultantCompany(ctx, input.ConsultantID, input.CompanyID, input.End)
	if err != nil {
		if errors.Is(err, contractrepo.ErrNotFound) {
			return domain.Invoice{}, ErrInvoiceConflict
		}

		return domain.Invoice{}, err
	}

	entries, err := s.entryRepo.ListForConsultantPeriod(
		ctx,
		input.ConsultantID,
		input.CompanyID,
		input.Start,
		input.End,
	)
	if err != nil {
		return domain.Invoice{}, ErrInvoiceConflict
	}

	issuedAt := time.Now()
	number := generateInvoiceNumber(input.End)

	if len(entries) == 0 {
		return domain.Invoice{
			Number:     number,
			IssuedAt:   issuedAt,
			Start:      input.Start,
			End:        input.End,
			Consultant: consultant,
			Company:    company,
			Contract:   contract,
			Activities: []domain.InvoiceActivity{},
			Totals:     domain.InvoiceTotals{},
		}, nil
	}

	ticketIDSet := make(map[uint64]struct{})
	activityIDSet := make(map[uint64]struct{})

	for _, entry := range entries {
		ticketIDSet[entry.TicketID] = struct{}{}
		activityIDSet[entry.ActivityID] = struct{}{}
	}

	tickets, err := s.ticketRepo.ListByIDs(ctx, keys(ticketIDSet))
	if err != nil {
		return domain.Invoice{}, ErrInvoiceConflict
	}

	activities, err := s.activityRepo.ListByIDs(ctx, keys(activityIDSet))
	if err != nil {
		return domain.Invoice{}, ErrInvoiceConflict
	}

	ticketsByID := make(map[uint64]ticketdomain.Ticket, len(tickets))
	for _, ticket := range tickets {
		ticketsByID[ticket.ID] = ticket
	}

	activitiesByID := make(map[uint64]activitydomain.Activity, len(activities))
	for _, activity := range activities {
		activitiesByID[activity.ID] = activity
	}

	invoiceActivities := groupActivities(
		entries,
		activitiesByID,
		ticketsByID,
		contract.HourlyRate,
	)

	totals := calculateTotals(invoiceActivities)

	return domain.Invoice{
		Number:     number,
		IssuedAt:   issuedAt,
		Start:      input.Start,
		End:        input.End,
		Consultant: consultant,
		Company:    company,
		Contract:   contract,
		Activities: invoiceActivities,
		Totals:     totals,
	}, nil
}

func keys[K comparable](m map[K]struct{}) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func generateInvoiceNumber(periodEnd time.Time) string {
	invoiceDate := periodEnd.AddDate(0, 0, 1)
	return invoiceDate.Format("20060102")
}
