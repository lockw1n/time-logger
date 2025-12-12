package consultant

import (
	consultantdto "github.com/lockw1n/time-logger/internal/dto/consultant"
)

type Service interface {
	Create(req consultantdto.Request) (*consultantdto.Response, error)
	Update(id uint64, req consultantdto.Request) (*consultantdto.Response, error)
	Delete(id uint64) error

	Get(id uint64) (*consultantdto.Response, error)
	ListByLastName(lastName string) ([]consultantdto.Response, error)
	ListAll() ([]consultantdto.Response, error)
}
