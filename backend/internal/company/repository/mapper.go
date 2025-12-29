package repository

import "github.com/lockw1n/time-logger/internal/company/domain"

func toModel(d domain.Company) companyModel {
	return companyModel{
		ID:           d.ID,
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

func toDomain(m companyModel) domain.Company {
	return domain.Company{
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
	}
}

func toDomainSlice(models []companyModel) []domain.Company {
	out := make([]domain.Company, 0, len(models))
	for _, m := range models {
		out = append(out, toDomain(m))
	}
	return out
}
