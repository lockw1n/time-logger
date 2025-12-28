package domain

import "time"

type Contract struct {
	ID           uint64
	ConsultantID uint64
	CompanyID    uint64
	HourlyRate   float64
	Currency     string
	OrderNumber  string
	PaymentTerms *string
	StartDate    time.Time
	EndDate      *time.Time
}

func (c Contract) IsActive(at time.Time) bool {
	if c.StartDate.After(at) {
		return false
	}

	if c.EndDate == nil {
		return true
	}

	return !c.EndDate.Before(at)
}
