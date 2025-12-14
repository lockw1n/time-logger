package company

import (
	"context"
	"errors"

	companydto "github.com/lockw1n/time-logger/internal/dto/company"
	companymapper "github.com/lockw1n/time-logger/internal/mapper/company"
	companyrepo "github.com/lockw1n/time-logger/internal/repository/company"
)

type service struct {
	repo companyrepo.Repository
}

func NewService(repo companyrepo.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(req companydto.Request) (*companydto.Response, error) {
	model := companymapper.ToModel(req)
	created, err := s.repo.Create(model)

	if err != nil {
		return nil, err
	}

	out := companymapper.ToResponse(created)
	return &out, nil
}

func (s *service) Update(ctx context.Context, id uint64, req companydto.Request) (*companydto.Response, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, companyrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	model := companymapper.ToModel(req)
	model.ID = existing.ID
	model.CreatedAt = existing.CreatedAt

	updated, err := s.repo.Update(model)

	if err != nil {
		return nil, err
	}

	out := companymapper.ToResponse(updated)
	return &out, nil
}

func (s *service) Delete(id uint64) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, companyrepo.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *service) Get(ctx context.Context, id uint64) (*companydto.Response, error) {
	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, companyrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	out := companymapper.ToResponse(model)
	return &out, nil
}

func (s *service) GetByName(name string) (*companydto.Response, error) {
	model, err := s.repo.FindByName(name)
	if err != nil {
		if errors.Is(err, companyrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	out := companymapper.ToResponse(model)
	return &out, nil
}

func (s *service) ListAll() ([]companydto.Response, error) {
	list, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}

	out := make([]companydto.Response, len(list))
	for i := range list {
		out[i] = companymapper.ToResponse(&list[i])
	}

	return out, nil
}
