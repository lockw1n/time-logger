package handler

import "github.com/lockw1n/time-logger/internal/consultant/domain"

type ConsultantResponse struct {
	ID           uint64  `json:"id"`
	FirstName    string  `json:"first_name"`
	MiddleName   *string `json:"middle_name"`
	LastName     string  `json:"last_name"`
	AddressLine1 string  `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
	Zip          string  `json:"zip"`
	City         string  `json:"city"`
	Region       *string `json:"region"`
	Country      string  `json:"country"`
	TaxNumber    string  `json:"tax_number"`
	BankName     string  `json:"bank_name"`
	BankAddress  string  `json:"bank_address"`
	BankCountry  string  `json:"bank_country"`
	BankIBAN     string  `json:"bank_iban"`
	BankBIC      string  `json:"bank_bic"`
}

func toResponse(consultant domain.Consultant) ConsultantResponse {
	return ConsultantResponse{
		ID:           consultant.ID,
		FirstName:    consultant.FirstName,
		MiddleName:   consultant.MiddleName,
		LastName:     consultant.LastName,
		AddressLine1: consultant.AddressLine1,
		AddressLine2: consultant.AddressLine2,
		Zip:          consultant.Zip,
		City:         consultant.City,
		Region:       consultant.Region,
		Country:      consultant.Country,
		TaxNumber:    consultant.TaxNumber,
		BankName:     consultant.BankName,
		BankAddress:  consultant.BankAddress,
		BankCountry:  consultant.BankCountry,
		BankIBAN:     consultant.BankIBAN,
		BankBIC:      consultant.BankBIC,
	}
}
