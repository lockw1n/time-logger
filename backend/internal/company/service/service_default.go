package service

import (
	"context"
	"errors"

	"github.com/lockw1n/time-logger/internal/company/domain"
	"github.com/lockw1n/time-logger/internal/company/repository"
)

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateCompany(ctx context.Context, input CreateCompanyInput) (domain.Company, error) {
	if err := input.Validate(); err != nil {
		return domain.Company{}, ErrCompanyInvalid
	}

	company := domain.Company{
		Name:         input.Name,
		NameShort:    input.NameShort,
		TaxNumber:    input.TaxNumber,
		AddressLine1: input.AddressLine1,
		AddressLine2: input.AddressLine2,
		Zip:          input.Zip,
		City:         input.City,
		Region:       input.Region,
		Country:      input.Country,
	}

	created, err := s.repo.Create(ctx, company)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return domain.Company{}, ErrCompanyAlreadyExists
		}
		if errors.Is(err, repository.ErrConflict) {
			return domain.Company{}, ErrCompanyConflict
		}
		return domain.Company{}, err
	}

	return created, nil
}

func (s *service) UpdateCompany(ctx context.Context, id uint64, input UpdateCompanyInput) (domain.Company, error) {
	if input == (UpdateCompanyInput{}) {
		return domain.Company{}, ErrCompanyInvalid
	}
	if err := input.Validate(); err != nil {
		return domain.Company{}, ErrCompanyInvalid
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Company{}, ErrCompanyNotFound
		}
		return domain.Company{}, err
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.NameShort != nil {
		existing.NameShort = input.NameShort
	}
	if input.TaxNumber != nil {
		existing.TaxNumber = *input.TaxNumber
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

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return domain.Company{}, ErrCompanyConflict
		}
		return domain.Company{}, err
	}

	return updated, nil
}

func (s *service) DeleteCompany(ctx context.Context, id uint64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrCompanyNotFound
		}
		return err
	}

	return nil
}

func (s *service) GetCompany(ctx context.Context, id uint64) (domain.Company, error) {
	company, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Company{}, ErrCompanyNotFound
		}
		return domain.Company{}, err
	}

	return company, nil
}

func (s *service) ListCompanies(ctx context.Context) ([]domain.Company, error) {
	companies, err := s.repo.ListAll(ctx)
	if errors.Is(err, repository.ErrConflict) {
		return nil, ErrCompanyConflict
	}

	return companies, nil
}

func (s *service) ListCompaniesForConsultant(ctx context.Context, consultantID uint64) ([]domain.Company, error) {
	companies, err := s.repo.ListByConsultant(ctx, consultantID)
	if errors.Is(err, repository.ErrConflict) {
		return nil, ErrCompanyConflict
	}

	return companies, nil
}
