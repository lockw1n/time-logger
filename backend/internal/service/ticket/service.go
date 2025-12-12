package ticket

import (
	ticketdto "github.com/lockw1n/time-logger/internal/dto/ticket"
)

type Service interface {
	Create(req ticketdto.Request) (*ticketdto.Response, error)
	Update(id uint64, req ticketdto.Request) (*ticketdto.Response, error)
	Delete(id uint64) error

	Get(id uint64) (*ticketdto.Response, error)
	GetByCode(companyID uint64, code string) (*ticketdto.Response, error)
	ListByCompany(companyID uint64) ([]ticketdto.Response, error)
	ListAll() ([]ticketdto.Response, error)
}
