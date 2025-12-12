package entry

type Request struct {
	ConsultantID    uint64  `json:"consultant_id"`
	CompanyID       uint64  `json:"company_id"`
	TicketCode      string  `json:"ticket_code"`
	LabelID         *uint64 `json:"label_id"`
	Date            string  `json:"date"`
	DurationMinutes int     `json:"duration_minutes"`
	Comment         *string `json:"comment"`
}
