package label

import (
	"errors"

	labeldto "github.com/lockw1n/time-logger/internal/dto/label"
	labelmapper "github.com/lockw1n/time-logger/internal/mapper/label"
	labelrepo "github.com/lockw1n/time-logger/internal/repository/label"
)

type service struct {
	repo labelrepo.Repository
}

func NewService(repo labelrepo.Repository) Service {
	return &service{repo: repo}
}

/* =======================
   CREATE
========================== */

func (s *service) Create(req labeldto.Request) (*labeldto.Response, error) {
	model := labelmapper.ToModel(req)
	created, err := s.repo.Create(model)
	if err != nil {
		return nil, err
	}

	out := labelmapper.ToResponse(created)
	return &out, nil
}

/* =======================
   UPDATE
========================== */

func (s *service) Update(id uint64, req labeldto.Request) (*labeldto.Response, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, labelrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	model := labelmapper.ToModel(req)
	model.ID = existing.ID
	model.CreatedAt = existing.CreatedAt

	updated, err := s.repo.Update(model)

	if err != nil {
		return nil, err
	}

	out := labelmapper.ToResponse(updated)
	return &out, nil
}

/* =======================
   DELETE
========================== */

func (s *service) Delete(id uint64) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, labelrepo.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

/* =======================
   GET BY ID
========================== */

func (s *service) Get(id uint64) (*labeldto.Response, error) {
	model, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, labelrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	out := labelmapper.ToResponse(model)
	return &out, nil
}

/* =======================
   LIST BY COMPANY
========================== */

func (s *service) ListByCompany(companyID uint64) ([]labeldto.Response, error) {
	list, err := s.repo.FindByCompany(companyID)
	if err != nil {
		return nil, err
	}

	return labelmapper.ToResponses(list), nil
}

/* =======================
   GET BY NAME
========================== */

func (s *service) GetByName(companyID uint64, name string) (*labeldto.Response, error) {
	model, err := s.repo.FindByName(companyID, name)
	if err != nil {
		if errors.Is(err, labelrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	out := labelmapper.ToResponse(model)
	return &out, nil
}
