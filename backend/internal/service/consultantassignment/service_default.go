package consultantassignment

import (
	"errors"

	assignmentdto "github.com/lockw1n/time-logger/internal/dto/consultantassignment"
	assignmentmapper "github.com/lockw1n/time-logger/internal/mapper/consultantassignment"
	assignmentrepo "github.com/lockw1n/time-logger/internal/repository/consultantassignment"
)

type service struct {
	repo assignmentrepo.Repository
}

func NewService(repo assignmentrepo.Repository) Service {
	return &service{repo: repo}
}

/* --------------------- CREATE --------------------- */

func (s *service) Create(data assignmentdto.Request) (*assignmentdto.Response, error) {
	model, err := assignmentmapper.ToModel(data)
	if err != nil {
		return nil, err
	}

	created, err := s.repo.Create(model)

	if err != nil {
		return nil, err
	}

	out := assignmentmapper.ToResponse(created)
	return &out, nil
}

/* --------------------- UPDATE --------------------- */

func (s *service) Update(id uint64, data assignmentdto.Request) (*assignmentdto.Response, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, assignmentrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	model, err := assignmentmapper.ToModel(data)
	if err != nil {
		return nil, err
	}

	model.ID = existing.ID
	model.CreatedAt = existing.CreatedAt

	updated, err := s.repo.Update(model)

	if err != nil {
		return nil, err
	}

	out := assignmentmapper.ToResponse(updated)
	return &out, nil
}

/* --------------------- DELETE --------------------- */

func (s *service) Delete(id uint64) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, assignmentrepo.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

/* --------------------- GET by ID --------------------- */

func (s *service) Get(id uint64) (*assignmentdto.Response, error) {
	m, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, assignmentrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	out := assignmentmapper.ToResponse(m)
	return &out, nil
}

/* --------------------- LIST for Consultant --------------------- */

func (s *service) GetForConsultant(consultantID uint64) ([]assignmentdto.Response, error) {
	list, err := s.repo.FindByConsultant(consultantID)
	if err != nil {
		return nil, err
	}

	result := make([]assignmentdto.Response, len(list))
	for i := range list {
		result[i] = assignmentmapper.ToResponse(&list[i])
	}
	return result, nil
}

/* --------------------- LIST for Company --------------------- */

func (s *service) GetForCompany(companyID uint64) ([]assignmentdto.Response, error) {
	list, err := s.repo.FindByCompany(companyID)
	if err != nil {
		return nil, err
	}

	result := make([]assignmentdto.Response, len(list))
	for i := range list {
		result[i] = assignmentmapper.ToResponse(&list[i])
	}
	return result, nil
}

/* --------------------- GET Pair (Consultant + Company) --------------------- */

func (s *service) GetPair(consultantID, companyID uint64) (*assignmentdto.Response, error) {
	m, err := s.repo.FindByPair(consultantID, companyID)
	if err != nil {
		if errors.Is(err, assignmentrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	out := assignmentmapper.ToResponse(m)
	return &out, nil
}
