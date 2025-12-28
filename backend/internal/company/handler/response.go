package handler

import "github.com/lockw1n/time-logger/internal/company/domain"

type CompanyResponse struct {
	ID           uint64  `json:"id"`
	Name         string  `json:"name"`
	NameShort    *string `json:"name_short"`
	TaxNumber    string  `json:"tax_number"`
	AddressLine1 string  `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
	Zip          string  `json:"zip"`
	City         string  `json:"city"`
	Region       *string `json:"region"`
	Country      string  `json:"country"`
}

func toResponse(company domain.Company) CompanyResponse {
	return CompanyResponse{
		ID:           company.ID,
		Name:         company.Name,
		NameShort:    company.NameShort,
		TaxNumber:    company.TaxNumber,
		AddressLine1: company.AddressLine1,
		AddressLine2: company.AddressLine2,
		Zip:          company.Zip,
		City:         company.City,
		Region:       company.Region,
		Country:      company.Country,
	}
}
