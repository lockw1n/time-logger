package service

import (
	"context"
	"errors"
	"time"

	"github.com/lockw1n/time-logger/internal/contract/domain"
	"github.com/lockw1n/time-logger/internal/contract/repository"
)

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateContract(ctx context.Context, input CreateContractInput) (domain.Contract, error) {
	if err := input.Validate(); err != nil {
		return domain.Contract{}, ErrContractInvalid
	}

	contract := domain.Contract{
		ConsultantID: input.ConsultantID,
		CompanyID:    input.CompanyID,
		HourlyRate:   input.HourlyRate,
		Currency:     input.Currency,
		OrderNumber:  input.OrderNumber,
		PaymentTerms: input.PaymentTerms,
		StartDate:    normalizeDate(input.StartDate),
		EndDate:      normalizeDatePtr(input.EndDate),
	}

	created, err := s.repo.Create(ctx, contract)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return domain.Contract{}, ErrContractAlreadyExists
		}
		if errors.Is(err, repository.ErrConflict) {
			return domain.Contract{}, ErrContractConflict
		}
		return domain.Contract{}, err
	}

	return created, nil
}

func (s *service) UpdateContract(ctx context.Context, id uint64, input UpdateContractInput) (domain.Contract, error) {
	if input == (UpdateContractInput{}) {
		return domain.Contract{}, ErrContractInvalid
	}

	if err := input.Validate(); err != nil {
		return domain.Contract{}, ErrContractInvalid
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Contract{}, ErrContractNotFound
		}
		return domain.Contract{}, err
	}

	effectiveStart := existing.StartDate
	if input.StartDate != nil {
		effectiveStart = *input.StartDate
	}

	if input.EndDate != nil && input.EndDate.Before(effectiveStart) {
		return domain.Contract{}, ErrContractInvalid
	}

	if input.HourlyRate != nil {
		existing.HourlyRate = *input.HourlyRate
	}
	if input.Currency != nil {
		existing.Currency = *input.Currency
	}
	if input.OrderNumber != nil {
		existing.OrderNumber = *input.OrderNumber
	}
	if input.PaymentTerms != nil {
		existing.PaymentTerms = input.PaymentTerms
	}
	if input.StartDate != nil {
		existing.StartDate = normalizeDate(*input.StartDate)
	}

	normalizedEndDate := normalizeDatePtr(input.EndDate)
	if normalizedEndDate != nil {
		existing.EndDate = normalizedEndDate
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return domain.Contract{}, ErrContractConflict
		}
		return domain.Contract{}, err
	}

	return updated, nil
}

func (s *service) DeleteContract(ctx context.Context, id uint64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrContractNotFound
		}
		return err
	}

	return nil
}

func (s *service) GetContract(ctx context.Context, id uint64) (domain.Contract, error) {
	contract, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Contract{}, ErrContractNotFound
		}
		return domain.Contract{}, err
	}

	return contract, nil
}

func (s *service) ListContractsForConsultant(ctx context.Context, consultantID uint64) ([]domain.Contract, error) {
	contracts, err := s.repo.ListByConsultant(ctx, consultantID)
	if errors.Is(err, repository.ErrConflict) {
		return nil, ErrContractConflict
	}

	return contracts, nil
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

func normalizeDatePtr(t *time.Time) *time.Time {
	if t == nil || t.IsZero() {
		return nil
	}

	normalized := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0, 0, 0, 0,
		t.Location(),
	)

	return &normalized
}
