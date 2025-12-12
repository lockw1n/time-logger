package ticket

import (
	"errors"

	ticketdto "github.com/lockw1n/time-logger/internal/dto/ticket"
	ticketmapper "github.com/lockw1n/time-logger/internal/mapper/ticket"
	ticketrepo "github.com/lockw1n/time-logger/internal/repository/ticket"
)

type service struct {
	repo ticketrepo.Repository
}

func NewService(repo ticketrepo.Repository) Service {
	return &service{repo: repo}
}

/* =======================
   CREATE
========================== */

func (s *service) Create(req ticketdto.Request) (*ticketdto.Response, error) {
	model := ticketmapper.ToModel(req)
	created, err := s.repo.Create(model)
	if err != nil {
		return nil, err
	}

	out := ticketmapper.ToResponse(created)
	return &out, nil
}

/* =======================
   UPDATE
========================== */

func (s *service) Update(id uint64, req ticketdto.Request) (*ticketdto.Response, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, ticketrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	model := ticketmapper.ToModel(req)
	model.ID = existing.ID
	model.CreatedAt = existing.CreatedAt

	updated, err := s.repo.Update(model)

	if err != nil {
		return nil, err
	}

	out := ticketmapper.ToResponse(updated)
	return &out, nil
}

/* =======================
   DELETE
========================== */

func (s *service) Delete(id uint64) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, ticketrepo.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

/* =======================
   GET BY ID
========================== */

func (s *service) Get(id uint64) (*ticketdto.Response, error) {
	model, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, ticketrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	out := ticketmapper.ToResponse(model)
	return &out, nil
}

/* =======================
   GET BY CODE
========================== */

func (s *service) GetByCode(companyID uint64, code string) (*ticketdto.Response, error) {
	model, err := s.repo.FindByCode(companyID, code)
	if err != nil {
		if errors.Is(err, ticketrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	out := ticketmapper.ToResponse(model)
	return &out, nil
}

/* =======================
   LIST BY COMPANY
========================== */

func (s *service) ListByCompany(companyID uint64) ([]ticketdto.Response, error) {
	list, err := s.repo.FindByCompany(companyID)
	if err != nil {
		return nil, err
	}

	return ticketmapper.ToResponses(list), nil
}

/* =======================
   LIST ALL
========================== */

func (s *service) ListAll() ([]ticketdto.Response, error) {
	list, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}

	return ticketmapper.ToResponses(list), nil
}
