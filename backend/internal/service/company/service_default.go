package company

import (
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

func (s *service) Create(data companydto.Request) (*companydto.Response, error) {
	// DTO → Model
	model := companymapper.ToModel(data)

	// Save
	if err := s.repo.Create(model); err != nil {
		return nil, err
	}

	// Model → DTO
	out := companymapper.ToResponse(model)
	return &out, nil
}

func (s *service) Update(id uint64, data companydto.Request) (*companydto.Response, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	updated := companymapper.ToModel(data)
	updated.ID = existing.ID
	updated.CreatedAt = existing.CreatedAt

	if err := s.repo.Update(updated); err != nil {
		return nil, err
	}

	out := companymapper.ToResponse(updated)
	return &out, nil
}

func (s *service) Delete(id uint64) error {
	return s.repo.Delete(id)
}

func (s *service) Get(id uint64) (*companydto.Response, error) {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	out := companymapper.ToResponse(m)
	return &out, nil
}

func (s *service) GetByName(name string) (*companydto.Response, error) {
	m, err := s.repo.FindByName(name)
	if err != nil {
		return nil, err
	}

	out := companymapper.ToResponse(m)
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
