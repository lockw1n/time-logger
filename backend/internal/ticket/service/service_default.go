package service

import (
	"context"
	"errors"

	"github.com/lockw1n/time-logger/internal/ticket/domain"
	"github.com/lockw1n/time-logger/internal/ticket/repository"
)

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateTicket(ctx context.Context, input CreateTicketInput) (domain.Ticket, error) {
	if err := input.Validate(); err != nil {
		return domain.Ticket{}, ErrTicketInvalid
	}

	ticket := domain.Ticket{
		CompanyID:   input.CompanyID,
		Code:        input.Code,
		Title:       input.Title,
		Label:       input.Label,
		Description: input.Description,
	}

	created, err := s.repo.Create(ctx, ticket)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return domain.Ticket{}, ErrTicketAlreadyExists
		}
		if errors.Is(err, repository.ErrConflict) {
			return domain.Ticket{}, ErrTicketConflict
		}
		return domain.Ticket{}, err
	}

	return created, nil
}

func (s *service) UpdateTicket(ctx context.Context, id uint64, input UpdateTicketInput) (domain.Ticket, error) {
	if input == (UpdateTicketInput{}) {
		return domain.Ticket{}, ErrTicketInvalid
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Ticket{}, ErrTicketNotFound
		}
		return domain.Ticket{}, err
	}

	if input.Title != nil {
		existing.Title = input.Title
	}
	if input.Label != nil {
		existing.Label = input.Label
	}
	if input.Description != nil {
		existing.Description = input.Description
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return domain.Ticket{}, ErrTicketConflict
		}
		return domain.Ticket{}, err
	}

	return updated, nil
}

func (s *service) DeleteTicket(ctx context.Context, id uint64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTicketNotFound
		}
		return err
	}

	return nil
}

func (s *service) GetTicket(ctx context.Context, id uint64) (domain.Ticket, error) {
	ticket, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Ticket{}, ErrTicketNotFound
		}
		return domain.Ticket{}, err
	}

	return ticket, nil
}

func (s *service) ListTicketsForCompany(ctx context.Context, companyID uint64) ([]domain.Ticket, error) {
	tickets, err := s.repo.ListByCompany(ctx, companyID)
	if errors.Is(err, repository.ErrConflict) {
		return nil, ErrTicketConflict
	}

	return tickets, nil
}
