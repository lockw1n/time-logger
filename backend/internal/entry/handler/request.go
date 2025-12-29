package handler

type CreateEntryRequest struct {
	ConsultantID    uint64  `json:"consultant_id"`
	CompanyID       uint64  `json:"company_id"`
	TicketCode      string  `json:"ticket_code"`
	ActivityID      uint64  `json:"activity_id"`
	Date            string  `json:"date"`
	DurationMinutes int     `json:"duration_minutes"`
	Comment         *string `json:"comment"`
}

type UpdateEntryRequest struct {
	DurationMinutes *int    `json:"duration_minutes"`
	Comment         *string `json:"comment"`
}
