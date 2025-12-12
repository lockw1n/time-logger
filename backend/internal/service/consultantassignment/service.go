package consultantassignment

import (
	assignmentdto "github.com/lockw1n/time-logger/internal/dto/consultantassignment"
)

type Service interface {
	Create(req assignmentdto.Request) (*assignmentdto.Response, error)
	Update(id uint64, req assignmentdto.Request) (*assignmentdto.Response, error)
	Delete(id uint64) error

	Get(id uint64) (*assignmentdto.Response, error)
	GetForConsultant(consultantID uint64) ([]assignmentdto.Response, error)
	GetForCompany(companyID uint64) ([]assignmentdto.Response, error)
	GetPair(consultantID, companyID uint64) (*assignmentdto.Response, error)
}
