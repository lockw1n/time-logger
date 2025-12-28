package repository

import "github.com/lockw1n/time-logger/internal/contract/domain"

func toModel(d domain.Contract) contractModel {
	return contractModel{
		ID:           d.ID,
		ConsultantID: d.ConsultantID,
		CompanyID:    d.CompanyID,
		HourlyRate:   d.HourlyRate,
		Currency:     d.Currency,
		OrderNumber:  d.OrderNumber,
		PaymentTerms: d.PaymentTerms,
		StartDate:    d.StartDate,
		EndDate:      d.EndDate,
	}
}

func toDomain(m contractModel) domain.Contract {
	return domain.Contract{
		ID:           m.ID,
		ConsultantID: m.ConsultantID,
		CompanyID:    m.CompanyID,
		HourlyRate:   m.HourlyRate,
		Currency:     m.Currency,
		OrderNumber:  m.OrderNumber,
		PaymentTerms: m.PaymentTerms,
		StartDate:    m.StartDate,
		EndDate:      m.EndDate,
	}
}

func toDomainSlice(models []contractModel) []domain.Contract {
	out := make([]domain.Contract, 0, len(models))
	for _, m := range models {
		out = append(out, toDomain(m))
	}
	return out
}
