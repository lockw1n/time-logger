package company

import (
	companydto "github.com/lockw1n/time-logger/internal/dto/company"
)

type Service interface {
	Create(data companydto.Request) (*companydto.Response, error)
	Update(id uint64, data companydto.Request) (*companydto.Response, error)
	Delete(id uint64) error

	Get(id uint64) (*companydto.Response, error)
	GetByName(name string) (*companydto.Response, error)
	ListAll() ([]companydto.Response, error)
}
