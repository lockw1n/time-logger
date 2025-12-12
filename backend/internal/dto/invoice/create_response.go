package invoice

import (
	ilDTO "github.com/lockw1n/time-logger/internal/dto/invoiceline"
)

type GenerateResponse struct {
	ConsultantID uint64 `json:"consultant_id"`
	CompanyID    uint64 `json:"company_id"`

	PeriodStart string `json:"period_start"`
	PeriodEnd   string `json:"period_end"`

	TotalAmount float64 `json:"total_amount"`
	Currency    string  `json:"currency"`

	Lines []ilDTO.Line `json:"lines"`
}
