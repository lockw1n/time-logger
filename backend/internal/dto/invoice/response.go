package invoice

import (
	ilDTO "github.com/lockw1n/time-logger/internal/dto/invoiceline"
)

type Response struct {
	ID           uint64 `json:"id"`
	ConsultantID uint64 `json:"consultant_id"`
	CompanyID    uint64 `json:"company_id"`

	PeriodStart string `json:"period_start"`
	PeriodEnd   string `json:"period_end"`

	TotalAmount float64 `json:"total_amount"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`

	InvoiceNumber *string `json:"invoice_number,omitempty"`
	OrderNumber   *string `json:"order_number,omitempty"`

	Metadata map[string]any `json:"metadata,omitempty"`

	Lines []ilDTO.Line `json:"lines"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
