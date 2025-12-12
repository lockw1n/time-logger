package consultantassignment

type Response struct {
	ID           uint64 `json:"id"`
	ConsultantID uint64 `json:"consultant_id"`
	CompanyID    uint64 `json:"company_id"`

	HourlyRate  float64 `json:"hourly_rate"`
	Currency    string  `json:"currency"`
	OrderNumber string  `json:"order_number"`

	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
