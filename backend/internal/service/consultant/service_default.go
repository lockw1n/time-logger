package consultant

import (
	"context"
	"errors"

	consultantdto "github.com/lockw1n/time-logger/internal/dto/consultant"
	consultantmapper "github.com/lockw1n/time-logger/internal/mapper/consultant"
	consultantrepo "github.com/lockw1n/time-logger/internal/repository/consultant"
)

type service struct {
	repo consultantrepo.Repository
}

func NewService(repo consultantrepo.Repository) Service {
	return &service{repo: repo}
}

/* ------------------ CREATE ------------------ */

func (s *service) Create(req consultantdto.Request) (*consultantdto.Response, error) {
	model := consultantmapper.ToModel(req)
	created, err := s.repo.Create(model)

	if err != nil {
		return nil, err
	}

	out := consultantmapper.ToResponse(created)
	return &out, nil
}

/* ------------------ UPDATE ------------------ */

func (s *service) Update(ctx context.Context, id uint64, req consultantdto.Request) (*consultantdto.Response, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, consultantrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	model := consultantmapper.ToModel(req)
	model.ID = existing.ID
	model.CreatedAt = existing.CreatedAt

	updated, err := s.repo.Update(model)

	if err != nil {
		return nil, err
	}

	out := consultantmapper.ToResponse(updated)
	return &out, nil
}

/* ------------------ DELETE ------------------ */

func (s *service) Delete(id uint64) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, consultantrepo.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

/* ------------------ GET ------------------ */

func (s *service) Get(ctx context.Context, id uint64) (*consultantdto.Response, error) {
	consultant, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, consultantrepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	out := consultantmapper.ToResponse(consultant)
	return &out, nil
}

/* ------------------ LIST BY LAST NAME ------------------ */

func (s *service) ListByLastName(lastName string) ([]consultantdto.Response, error) {
	list, err := s.repo.FindByLastName(lastName)
	if err != nil {
		return nil, err
	}

	result := make([]consultantdto.Response, len(list))
	for i := range list {
		result[i] = consultantmapper.ToResponse(&list[i])
	}

	return result, nil
}

/* ------------------ LIST ALL ------------------ */

func (s *service) ListAll() ([]consultantdto.Response, error) {
	list, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}

	result := make([]consultantdto.Response, len(list))
	for i := range list {
		result[i] = consultantmapper.ToResponse(&list[i])
	}

	return result, nil
}
