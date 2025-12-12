package entry

import (
	entrydto "github.com/lockw1n/time-logger/internal/dto/entry"
)

type Service interface {
	Create(req entrydto.Request) (*entrydto.Response, error)
	Update(id uint64, req entrydto.Request) (*entrydto.Response, error)
	Delete(id uint64) error

	Get(id uint64) (*entrydto.Response, error)
	ListByCompany(companyID uint64) ([]entrydto.Response, error)
	ListByConsultant(consultantID uint64) ([]entrydto.Response, error)
	ListAll() ([]entrydto.Response, error)
}
