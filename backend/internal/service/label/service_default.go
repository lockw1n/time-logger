package label

import (
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

func (s *service) Create(data labeldto.Request) (*labeldto.Response, error) {
	model := labelmapper.ToModel(data)

	if err := s.repo.Create(model); err != nil {
		return nil, err
	}

	out := labelmapper.ToResponse(model)
	return &out, nil
}

/* =======================
   UPDATE
========================== */

func (s *service) Update(id uint64, data labeldto.Request) (*labeldto.Response, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	updated := labelmapper.ToModel(data)
	updated.ID = existing.ID
	updated.CreatedAt = existing.CreatedAt

	if err := s.repo.Update(updated); err != nil {
		return nil, err
	}

	out := labelmapper.ToResponse(updated)
	return &out, nil
}

/* =======================
   DELETE
========================== */

func (s *service) Delete(id uint64) error {
	return s.repo.Delete(id)
}

/* =======================
   GET BY ID
========================== */

func (s *service) Get(id uint64) (*labeldto.Response, error) {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	out := labelmapper.ToResponse(m)
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
	m, err := s.repo.FindByName(companyID, name)
	if err != nil {
		return nil, err
	}

	out := labelmapper.ToResponse(m)
	return &out, nil
}
