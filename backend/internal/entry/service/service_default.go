package service

import (
	"context"
	"errors"
	"strings"
	"time"

	contractrepo "github.com/lockw1n/time-logger/internal/contract/repository"
	"github.com/lockw1n/time-logger/internal/entry/domain"
	"github.com/lockw1n/time-logger/internal/entry/repository"
	ticketdomain "github.com/lockw1n/time-logger/internal/ticket/domain"
	ticketrepo "github.com/lockw1n/time-logger/internal/ticket/repository"
)

type service struct {
	repo         repository.Repository
	contractRepo contractrepo.Repository
	ticketRepo   ticketrepo.Repository
}

func NewService(
	repo repository.Repository,
	contractRepo contractrepo.Repository,
	ticketRepo ticketrepo.Repository,
) Service {
	return &service{
		repo:         repo,
		contractRepo: contractRepo,
		ticketRepo:   ticketRepo,
	}
}

func (s *service) CreateEntry(ctx context.Context, consultantID uint64, input CreateEntryInput) (domain.Entry, error) {
	if err := input.Validate(); err != nil {
		return domain.Entry{}, ErrEntryInvalid
	}

	date := normalizeDate(input.Date)
	contract, err := s.contractRepo.FindActiveByConsultantCompany(ctx, consultantID, input.CompanyID, date)
	if err != nil {
		if errors.Is(err, contractrepo.ErrNotFound) {
			return domain.Entry{}, ErrEntryConflict
		}
		return domain.Entry{}, err
	}

	ticket, err := s.resolveTicket(ctx, input.CompanyID, input.TicketCode)
	if err != nil {
		return domain.Entry{}, err
	}

	entry := domain.Entry{
		ConsultantID:    consultantID,
		CompanyID:       input.CompanyID,
		ContractID:      contract.ID,
		TicketID:        ticket.ID,
		ActivityID:      input.ActivityID,
		Date:            date,
		DurationMinutes: input.DurationMinutes,
		Comment:         input.Comment,
	}

	created, err := s.repo.Create(ctx, entry)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return domain.Entry{}, ErrEntryAlreadyExists
		}
		if errors.Is(err, repository.ErrConflict) {
			return domain.Entry{}, ErrEntryConflict
		}
		return domain.Entry{}, err
	}

	return created, nil
}

func (s *service) UpdateEntry(
	ctx context.Context,
	consultantID uint64,
	entryID uint64,
	input UpdateEntryInput,
) (domain.Entry, error) {
	if err := input.Validate(); err != nil {
		return domain.Entry{}, ErrEntryInvalid
	}

	existing, err := s.repo.FindByID(ctx, entryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Entry{}, ErrEntryNotFound
		}
		return domain.Entry{}, err
	}

	if existing.ConsultantID != consultantID {
		return domain.Entry{}, ErrForbidden
	}

	if input.DurationMinutes != nil {
		existing.DurationMinutes = *input.DurationMinutes
	}
	if input.Comment != nil {
		existing.Comment = input.Comment
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return domain.Entry{}, ErrEntryConflict
		}
		return domain.Entry{}, err
	}

	return updated, nil
}

func (s *service) DeleteEntry(ctx context.Context, consultantID uint64, entryID uint64) error {
	entry, err := s.repo.FindByID(ctx, entryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrEntryNotFound
		}
		return err
	}

	if entry.ConsultantID != consultantID {
		return ErrForbidden
	}

	if err := s.repo.Delete(ctx, entryID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrEntryNotFound
		}
		return err
	}

	return nil
}

func (s *service) GetEntry(ctx context.Context, consultantID uint64, entryID uint64) (domain.Entry, error) {
	entry, err := s.repo.FindByID(ctx, entryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Entry{}, ErrEntryNotFound
		}
		return domain.Entry{}, err
	}

	if entry.ConsultantID != consultantID {
		return domain.Entry{}, ErrForbidden
	}

	return entry, nil
}

func (s *service) resolveTicket(ctx context.Context, companyID uint64, code string) (ticketdomain.Ticket, error) {

	code = strings.TrimSpace(code)

	ticket, err := s.ticketRepo.FindByCompanyAndCode(ctx, companyID, code)
	if err == nil {
		return ticket, nil
	}

	if !errors.Is(err, ticketrepo.ErrNotFound) {
		return ticketdomain.Ticket{}, err
	}

	ticket, err = s.ticketRepo.Create(ctx, ticketdomain.Ticket{
		CompanyID: companyID,
		Code:      code,
	})
	if err != nil {
		if errors.Is(err, ticketrepo.ErrAlreadyExists) {
			return s.ticketRepo.FindByCompanyAndCode(ctx, companyID, code)
		}
		return ticketdomain.Ticket{}, err
	}

	return ticket, nil
}

func normalizeDate(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0, 0, 0, 0,
		t.Location(),
	)
}
