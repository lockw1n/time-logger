package consultant

import (
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

func (s *service) Create(data consultantdto.Request) (*consultantdto.Response, error) {
	model := consultantmapper.ToModel(data)

	if err := s.repo.Create(model); err != nil {
		return nil, err
	}

	out := consultantmapper.ToResponse(model)
	return &out, nil
}

/* ------------------ UPDATE ------------------ */

func (s *service) Update(id uint64, data consultantdto.Request) (*consultantdto.Response, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	updated := consultantmapper.ToModel(data)

	// preserve identity + timestamps
	updated.ID = existing.ID
	updated.CreatedAt = existing.CreatedAt

	if err := s.repo.Update(updated); err != nil {
		return nil, err
	}

	out := consultantmapper.ToResponse(updated)
	return &out, nil
}

/* ------------------ DELETE ------------------ */

func (s *service) Delete(id uint64) error {
	return s.repo.Delete(id)
}

/* ------------------ GET ------------------ */

func (s *service) Get(id uint64) (*consultantdto.Response, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	out := consultantmapper.ToResponse(c)
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
