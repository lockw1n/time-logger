package invoice

import (
	invoicedto "github.com/lockw1n/time-logger/internal/dto/invoice"
)

type Service interface {
	Generate(
		consultantID uint64,
		companyID uint64,
		periodStart string,
		periodEnd string,
	) (*invoicedto.GenerateResponse, error)

	ListByCompany(companyID uint64) ([]invoicedto.ListItem, error)
	ListByConsultant(consultantID uint64) ([]invoicedto.ListItem, error)

	Get(id uint64) (*invoicedto.Response, error)
}
