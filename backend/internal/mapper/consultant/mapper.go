package consultant

import (
	"time"

	consultantdto "github.com/lockw1n/time-logger/internal/dto/consultant"
	"github.com/lockw1n/time-logger/internal/models"
)

/*************** DTO → Model ***************/

func ToModel(d consultantdto.Request) *models.Consultant {
	return &models.Consultant{
		FirstName:    d.FirstName,
		MiddleName:   d.MiddleName,
		LastName:     d.LastName,
		AddressLine1: d.AddressLine1,
		AddressLine2: d.AddressLine2,
		Zip:          d.Zip,
		City:         d.City,
		Region:       d.Region,
		Country:      d.Country,
		TaxNumber:    d.TaxNumber,
		BankName:     d.BankName,
		BankAddress:  d.BankAddress,
		BankCountry:  d.BankCountry,
		BankIBAN:     d.BankIBAN,
		BankBIC:      d.BankBIC,
	}
}

/*************** Model → DTO ***************/

func ToResponse(m *models.Consultant) consultantdto.Response {
	return consultantdto.Response{
		ID:           m.ID,
		FirstName:    m.FirstName,
		MiddleName:   m.MiddleName,
		LastName:     m.LastName,
		AddressLine1: m.AddressLine1,
		AddressLine2: m.AddressLine2,
		Zip:          m.Zip,
		City:         m.City,
		Region:       m.Region,
		Country:      m.Country,
		TaxNumber:    m.TaxNumber,
		BankName:     m.BankName,
		BankAddress:  m.BankAddress,
		BankCountry:  m.BankCountry,
		BankIBAN:     m.BankIBAN,
		BankBIC:      m.BankBIC,
		CreatedAt:    m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    m.UpdatedAt.Format(time.RFC3339),
	}
}
