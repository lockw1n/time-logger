package service

import "time"

type CreateContractInput struct {
	ConsultantID uint64
	CompanyID    uint64
	HourlyRate   float64
	Currency     string
	OrderNumber  string
	PaymentTerms *string
	StartDate    time.Time
	EndDate      *time.Time
}

type UpdateContractInput struct {
	HourlyRate   *float64
	Currency     *string
	OrderNumber  *string
	PaymentTerms *string
	StartDate    *time.Time
	EndDate      *time.Time
}
