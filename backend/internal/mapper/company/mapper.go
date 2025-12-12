package company

import (
	"time"

	companydto "github.com/lockw1n/time-logger/internal/dto/company"
	"github.com/lockw1n/time-logger/internal/models"
)

/*************** DTO → Model ***************/

func ToModel(d companydto.Request) *models.Company {
	return &models.Company{
		Name:         d.Name,
		NameShort:    d.NameShort,
		TaxNumber:    d.TaxNumber,
		AddressLine1: d.AddressLine1,
		AddressLine2: d.AddressLine2,
		Zip:          d.Zip,
		City:         d.City,
		Region:       d.Region,
		Country:      d.Country,
	}
}

/*************** Model → DTO ***************/

func ToResponse(m *models.Company) companydto.Response {
	return companydto.Response{
		ID:           m.ID,
		Name:         m.Name,
		NameShort:    m.NameShort,
		TaxNumber:    m.TaxNumber,
		AddressLine1: m.AddressLine1,
		AddressLine2: m.AddressLine2,
		Zip:          m.Zip,
		City:         m.City,
		Region:       m.Region,
		Country:      m.Country,
		CreatedAt:    m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    m.UpdatedAt.Format(time.RFC3339),
	}
}
