package service

import (
	"context"
	"errors"

	"github.com/lockw1n/time-logger/internal/consultant/domain"
	"github.com/lockw1n/time-logger/internal/consultant/repository"
)

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateConsultant(ctx context.Context, input CreateConsultantInput) (domain.Consultant, error) {
	if err := input.Validate(); err != nil {
		return domain.Consultant{}, ErrConsultantInvalid
	}

	consultant := domain.Consultant{
		FirstName:    input.FirstName,
		MiddleName:   input.MiddleName,
		LastName:     input.LastName,
		AddressLine1: input.AddressLine1,
		AddressLine2: input.AddressLine2,
		Zip:          input.Zip,
		City:         input.City,
		Region:       input.Region,
		Country:      input.Country,
		TaxNumber:    input.TaxNumber,
		BankName:     input.BankName,
		BankAddress:  input.BankAddress,
		BankCountry:  input.BankCountry,
		BankIBAN:     input.BankIBAN,
		BankBIC:      input.BankBIC,
	}

	created, err := s.repo.Create(ctx, consultant)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return domain.Consultant{}, ErrConsultantAlreadyExists
		}
		if errors.Is(err, repository.ErrConflict) {
			return domain.Consultant{}, ErrConsultantConflict
		}
		return domain.Consultant{}, err
	}

	return created, nil
}

func (s *service) UpdateConsultant(
	ctx context.Context,
	id uint64,
	input UpdateConsultantInput,
) (domain.Consultant, error) {
	if input == (UpdateConsultantInput{}) {
		return domain.Consultant{}, ErrConsultantInvalid
	}
	if err := input.Validate(); err != nil {
		return domain.Consultant{}, ErrConsultantInvalid
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Consultant{}, ErrConsultantNotFound
		}
		return domain.Consultant{}, err
	}

	if input.FirstName != nil {
		existing.FirstName = *input.FirstName
	}
	if input.MiddleName != nil {
		existing.MiddleName = input.MiddleName
	}
	if input.LastName != nil {
		existing.LastName = *input.LastName
	}
	if input.AddressLine1 != nil {
		existing.AddressLine1 = *input.AddressLine1
	}
	if input.AddressLine2 != nil {
		existing.AddressLine2 = input.AddressLine2
	}
	if input.Zip != nil {
		existing.Zip = *input.Zip
	}
	if input.City != nil {
		existing.City = *input.City
	}
	if input.Region != nil {
		existing.Region = input.Region
	}
	if input.Country != nil {
		existing.Country = *input.Country
	}
	if input.TaxNumber != nil {
		existing.TaxNumber = *input.TaxNumber
	}
	if input.BankName != nil {
		existing.BankName = *input.BankName
	}
	if input.BankAddress != nil {
		existing.BankAddress = *input.BankAddress
	}
	if input.BankCountry != nil {
		existing.BankCountry = *input.BankCountry
	}
	if input.BankIBAN != nil {
		existing.BankIBAN = *input.BankIBAN
	}
	if input.BankBIC != nil {
		existing.BankBIC = *input.BankBIC
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return domain.Consultant{}, ErrConsultantConflict
		}
		return domain.Consultant{}, err
	}

	return updated, nil
}

func (s *service) DeleteConsultant(ctx context.Context, id uint64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrConsultantNotFound
		}
		return err
	}

	return nil
}

func (s *service) GetConsultant(ctx context.Context, id uint64) (domain.Consultant, error) {
	consultant, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Consultant{}, ErrConsultantNotFound
		}
		return domain.Consultant{}, err
	}

	return consultant, nil
}
