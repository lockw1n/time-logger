package label

import (
	labeldto "github.com/lockw1n/time-logger/internal/dto/label"
)

type Service interface {
	Create(data labeldto.Request) (*labeldto.Response, error)
	Update(id uint64, data labeldto.Request) (*labeldto.Response, error)
	Delete(id uint64) error

	Get(id uint64) (*labeldto.Response, error)
	ListByCompany(companyID uint64) ([]labeldto.Response, error)
	GetByName(companyID uint64, name string) (*labeldto.Response, error)
}
