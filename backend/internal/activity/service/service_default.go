package service

import (
	"context"
	"errors"

	"github.com/lockw1n/time-logger/internal/activity/domain"
	"github.com/lockw1n/time-logger/internal/activity/repository"
)

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateActivity(ctx context.Context, input CreateActivityInput) (domain.Activity, error) {
	if err := input.Validate(); err != nil {
		return domain.Activity{}, ErrActivityInvalid
	}

	activity := domain.Activity{
		CompanyID: input.CompanyID,
		Name:      input.Name,
		Color:     input.Color,
		Billable:  input.Billable,
		Priority:  input.Priority,
	}

	created, err := s.repo.Create(ctx, activity)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return domain.Activity{}, ErrActivityAlreadyExists
		}
		if errors.Is(err, repository.ErrConflict) {
			return domain.Activity{}, ErrActivityConflict
		}
		return domain.Activity{}, err
	}

	return created, nil
}

func (s *service) UpdateActivity(ctx context.Context, id uint64, input UpdateActivityInput) (domain.Activity, error) {
	if input == (UpdateActivityInput{}) {
		return domain.Activity{}, ErrActivityInvalid
	}
	if err := input.Validate(); err != nil {
		return domain.Activity{}, ErrActivityInvalid
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Activity{}, ErrActivityNotFound
		}
		return domain.Activity{}, err
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Color != nil {
		existing.Color = input.Color
	}
	if input.Billable != nil {
		existing.Billable = *input.Billable
	}
	if input.Priority != nil {
		existing.Priority = *input.Priority
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return domain.Activity{}, ErrActivityConflict
		}
		return domain.Activity{}, err
	}

	return updated, nil
}

func (s *service) DeleteActivity(ctx context.Context, id uint64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrActivityNotFound
		}
		return err
	}

	return nil
}

func (s *service) GetActivity(ctx context.Context, id uint64) (domain.Activity, error) {
	activity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Activity{}, ErrActivityNotFound
		}
		return domain.Activity{}, err
	}

	return activity, nil
}

func (s *service) ListActivitiesForCompany(ctx context.Context, companyID uint64) ([]domain.Activity, error) {
	activities, err := s.repo.ListByCompany(ctx, companyID)
	if errors.Is(err, repository.ErrConflict) {
		return nil, ErrActivityConflict
	}

	return activities, nil
}
