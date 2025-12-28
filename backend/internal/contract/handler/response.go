package handler

import (
	"github.com/lockw1n/time-logger/internal/constants"
	"github.com/lockw1n/time-logger/internal/contract/domain"
)

type ContractResponse struct {
	ID           uint64  `json:"id"`
	ConsultantID uint64  `json:"consultant_id"`
	CompanyID    uint64  `json:"company_id"`
	HourlyRate   float64 `json:"hourly_rate"`
	Currency     string  `json:"currency"`
	OrderNumber  string  `json:"order_number"`
	PaymentTerms *string `json:"payment_terms"`
	StartDate    string  `json:"start_date"`
	EndDate      *string `json:"end_date"`
}

func toResponse(contract domain.Contract) ContractResponse {
	var endDate *string
	if contract.EndDate != nil {
		formatted := contract.EndDate.Format(constants.ResponseDateFormat)
		endDate = &formatted
	}

	return ContractResponse{
		ID:           contract.ID,
		ConsultantID: contract.ConsultantID,
		CompanyID:    contract.CompanyID,
		HourlyRate:   contract.HourlyRate,
		Currency:     contract.Currency,
		OrderNumber:  contract.OrderNumber,
		PaymentTerms: contract.PaymentTerms,
		StartDate:    contract.StartDate.Format(constants.ResponseDateFormat),
		EndDate:      endDate,
	}
}
