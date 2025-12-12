package consultantassignment

type Request struct {
	ConsultantID uint64  `json:"consultant_id"`
	CompanyID    uint64  `json:"company_id"`
	HourlyRate   float64 `json:"hourly_rate"`
	Currency     string  `json:"currency"`
	OrderNumber  string  `json:"order_number"`
	StartDate    *string `json:"start_date"` // YYYY-MM-DD
	EndDate      *string `json:"end_date"`   // YYYY-MM-DD
}
