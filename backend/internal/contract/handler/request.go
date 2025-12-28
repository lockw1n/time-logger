package handler

type CreateContractRequest struct {
	ConsultantID uint64  `json:"consultant_id"`
	CompanyID    uint64  `json:"company_id"`
	HourlyRate   float64 `json:"hourly_rate"`
	Currency     string  `json:"currency"`
	OrderNumber  string  `json:"order_number"`
	PaymentTerms *string `json:"payment_terms"`
	StartDate    string  `json:"start_date"`
	EndDate      *string `json:"end_date"`
}

type UpdateContractRequest struct {
	HourlyRate   *float64 `json:"hourly_rate"`
	Currency     *string  `json:"currency"`
	OrderNumber  *string  `json:"order_number"`
	PaymentTerms *string  `json:"payment_terms"`
	StartDate    *string  `json:"start_date"`
	EndDate      *string  `json:"end_date"`
}
