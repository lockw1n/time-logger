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

func (s *service) CreateEntry(ctx context.Context, input CreateEntryInput) (domain.Entry, error) {
	if err := input.Validate(); err != nil {
		return domain.Entry{}, ErrEntryInvalid
	}

	date := normalizeDate(input.Date)
	contract, err := s.contractRepo.FindActiveByConsultantCompany(ctx, input.ConsultantID, input.CompanyID, date)
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
		ConsultantID:    input.ConsultantID,
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

func (s *service) UpdateEntry(ctx context.Context, id uint64, input UpdateEntryInput) (domain.Entry, error) {
	if input == (UpdateEntryInput{}) {
		return domain.Entry{}, ErrEntryInvalid
	}

	if err := input.Validate(); err != nil {
		return domain.Entry{}, ErrEntryInvalid
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Entry{}, ErrEntryNotFound
		}
		return domain.Entry{}, err
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

func (s *service) DeleteEntry(ctx context.Context, id uint64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrEntryNotFound
		}
		return err
	}

	return nil
}

func (s *service) GetEntry(ctx context.Context, id uint64) (domain.Entry, error) {
	entry, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Entry{}, ErrEntryNotFound
		}
		return domain.Entry{}, err
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
